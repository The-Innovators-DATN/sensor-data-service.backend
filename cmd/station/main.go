package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"sensor-data-service.backend/config"
	"sensor-data-service.backend/internal/cache"
	"sensor-data-service.backend/internal/db"
	"sensor-data-service.backend/internal/metric"
)

func main() {
	// Load config
	ctx := context.Background()

	cfg, err := config.LoadAllConfigs("config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init Postgres
	postgresDB, err := db.InitDB(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}
	defer postgresDB.Close(ctx)

	PGStore := db.NewPostgresStore(postgresDB)

	rows, err := PGStore.ExecQuery(ctx, "SELECT * FROM station LIMIT 5")
	if err != nil {
		log.Fatal("Query failed:", err)
	}

	for _, row := range rows {
		fmt.Println(row)
	}

	// Init ClickHouse
	clickhouseDB, err := metric.InitClickHouse(cfg.Clickhouse)
	if err != nil {
		log.Fatalf("failed to connect to ClickHouse: %v", err)
	}
	defer clickhouseDB.Close()

	CHStore := metric.NewClickhouseStore(clickhouseDB)

	sql := `SELECT * FROM messages_local LIMIT 10`
	rows, err = CHStore.ExecQuery(ctx, sql)
	if err != nil {
		log.Fatal("Query failed:", err)
	}

	for _, row := range rows {
		log.Println(row)
	}
	// Init Redis
	redisClient, err := cache.InitRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
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
	// TCP listener
	listener, err := net.Listen("tcp", cfg.App.HostPort)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", cfg.App.HostPort, err)
	}
	defer listener.Close()

	log.Printf("Server started on %s", cfg.App.HostPort)

	// Graceful shutdown handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		listener.Close()
		os.Exit(0)
	}()

	// Start accepting connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// Your handler logic
	log.Printf("New connection from %s", conn.RemoteAddr())
}
