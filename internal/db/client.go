package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"sensor-data-service.backend/config" // đổi thành tên module của bạn
)

func InitDB(ctx context.Context, cfg config.DBConfig) (*pgx.Conn, error) {
	// Validate config
	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" || cfg.Name == "" {
		return nil, fmt.Errorf("missing database config values")
	}

	// Build connection string
	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	// Connect
	db, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
		return nil, err
	}

	log.Println("Connected to database")
	return db, nil
}
