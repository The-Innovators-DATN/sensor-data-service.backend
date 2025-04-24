package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	metricpb "sensor-data-service.backend/api/pb/metricdatapb"
	paramterpb "sensor-data-service.backend/api/pb/parameterpb"
	stationpb "sensor-data-service.backend/api/pb/stationpb"
	"sensor-data-service.backend/config"
	"sensor-data-service.backend/pkg/logger"

	"sensor-data-service.backend/internal/common/middleware"
	"sensor-data-service.backend/internal/domain/repository"
	"sensor-data-service.backend/internal/domain/service"
	"sensor-data-service.backend/internal/infrastructure/cache"
	"sensor-data-service.backend/internal/infrastructure/db"
	"sensor-data-service.backend/internal/infrastructure/metric"
	"sensor-data-service.backend/internal/transport/handler"

	dashboardpb "sensor-data-service.backend/api/pb/dashboardpb"
)

func main() {
	// Load config
	ctx := context.Background()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg, err := config.LoadAllConfigs("config")
	// if err != nil {
	// 	logger.Fatalf("failed to load config: ", err)
	// }
	log.Printf("Config loaded: %+v", cfg)
	logger := logger.New(cfg.Log)
	logger.Info("⇢ initializing backend...")

	// Init Postgres
	postgresDB, err := db.InitDB(ctx, cfg.DB)
	if err != nil {
		logger.Fatalf("failed to connect to Postgres:", err)
	}
	defer postgresDB.Close()

	PGStore := db.NewPostgresStore(postgresDB)

	rows, err := PGStore.ExecQuery(ctx, "SELECT * FROM station LIMIT 5")
	if err != nil {
		logger.Fatal("Query failed:", err)
	}

	for _, row := range rows {
		fmt.Println(row)
	}

	// Init ClickHouse
	clickhouseDB, err := metric.InitClickHouse(cfg.Clickhouse)
	if err != nil {
		logger.Fatalf("failed to connect to ClickHouse: %v", err)
	}
	defer clickhouseDB.Close()

	CHStore := metric.NewClickhouseStore(clickhouseDB)

	// sql := `SELECT * FROM messages_sharded LIMIT 10`
	// rows, err = CHStore.ExecQuery(ctx, sql)
	// if err != nil {
	// 	log.Fatal("Query failed:", err)
	// }

	// for _, row := range rows {
	// 	log.Println(row)
	// }
	// Init Redis
	redisClient, err := cache.InitRedis(cfg.Redis)
	if err != nil {
		logger.Fatalf("failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	RedisStore := cache.NewRedisStore(redisClient)

	// demo redis cache
	// ctx := context.Background()
	_ = RedisStore.Set(ctx, "station:123:last_reading", `{"pH": 6.5}`, 300)
	val, _ := RedisStore.Get(ctx, "station:123:last_reading")
	log.Println("Cached value:", val)
	defer RedisStore.Delete(ctx, "station:123:last_reading")
	// Print the value
	// // TCP listener
	// listener, err := net.Listen("tcp", cfg.App.HostPort)
	// if err != nil {
	// 	log.Fatalf("failed to listen on %s: %v", cfg.App.HostPort, err)
	// }
	// defer listener.Close()

	// log.Printf("Server started on %s", cfg.App.HostPort)

	// // Graceful shutdown handler
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	// 	<-quit
	// 	log.Println("Shutting down server...")
	// 	listener.Close()
	// 	os.Exit(0)
	// }()

	// // Start accepting connections
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		log.Printf("accept error: %v", err)
	// 		continue
	// 	}
	// 	go handleConnection(conn)
	// }

	// TCP listener cho gRPC
	grpcListener, err := net.Listen("tcp", cfg.App.HostPort)
	if err != nil {
		logger.Fatalf("failed to listen on %s: %v", cfg.App.HostPort, err)
	}
	defer grpcListener.Close()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor()),
	)

	// Tạo repository & service & handler cho parameter
	paramRepo := repository.NewParameterRepository(PGStore, RedisStore)

	paramService := service.NewParameterService(paramRepo)
	paramGrpcHandler := handler.NewGrpcParameterHandler(paramService)
	// Register gRPC server
	paramterpb.RegisterParameterServiceServer(grpcServer, paramGrpcHandler)

	metricRepo := repository.NewMetricDataRepository(CHStore, PGStore, RedisStore)
	metricService := service.NewMetricDataService(metricRepo)
	metricGrpcHandler := handler.NewMetricDataHandler(metricService)
	metricpb.RegisterMetricDataServiceServer(grpcServer, metricGrpcHandler)

	stationRepo := repository.NewStationDataRepository(PGStore, RedisStore)
	stationService := service.NewStationService(stationRepo)

	stationGrpcHandler := handler.NewStationHandler(stationService)
	stationpb.RegisterStationServiceServer(grpcServer, stationGrpcHandler)

	dashboardRepo := repository.NewDashboardDataRepository(PGStore, RedisStore)
	dashboardService := service.NewDashboardService(dashboardRepo)
	dashboardGrpcHandler := handler.NewDashboardHandler(dashboardService)
	dashboardpb.RegisterDashboardServiceServer(grpcServer, dashboardGrpcHandler)

	// (Optional) enable reflection để dùng grpcurl debug
	reflection.Register(grpcServer)

	// Graceful shutdown
	go func() {
		log.Printf("gRPC server listening at %s", cfg.App.HostPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	// === REST Gateway ===
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}

		err := paramterpb.RegisterParameterServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)
		if err != nil {
			logger.Fatalf("Failed to start HTTP gateway: %v", err)
		}
		err = metricpb.RegisterMetricDataServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)
		if err != nil {
			logger.Fatalf("Failed to start HTTP gateway: %v", err)
		}
		err = stationpb.RegisterStationServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)
		if err != nil {
			logger.Fatalf("Failed to start HTTP gateway: %v", err)
		}

		err = dashboardpb.RegisterDashboardServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)
		if err != nil {
			logger.Fatalf("Failed to start HTTP gateway: %v", err)
		}
		httpMux := http.NewServeMux()

		httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[HTTP] %s %s", r.Method, r.URL.Path)
			mux.ServeHTTP(w, r)
		}))

		httpMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received ping: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "pong")
		})

		log.Println("REST gateway listening on :8081")
		if err := http.ListenAndServe(":8081", httpMux); err != nil {
			logger.Fatalf("Failed to serve HTTP gateway: %v", err)
		}

	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// Your handler logic
	log.Printf("New connection from %s", conn.RemoteAddr())
}
