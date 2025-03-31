package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"sensor-data-service.backend/config"
	"sensor-data-service.backend/internal/db"
)

func main() {
	// Load config
	cfg, err := config.LoadAllConfigs("config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init Postgres
	postgresDB, err := db.InitDB(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}
	defer postgresDB.Close()

	// Init ClickHouse
	clickhouseDB, err := db.InitClickHouse(cfg.Clickhouse)
	if err != nil {
		log.Fatalf("failed to connect to ClickHouse: %v", err)
	}
	defer clickhouseDB.Close()

	// Init Redis
	redisClient, err := db.InitRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// TCP listener
	listener, err := net.Listen("tcp", cfg.App.HostPort)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", cfg.App.HostPort, err)
	}
	defer listener.Close()

	log.Printf("🚀 Server started on %s", cfg.App.HostPort)

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
