package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"sensor-data-service.backend/config" // đổi thành tên module của bạn
)

func InitDB(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %w", err)
	}

	log.Println("Connected to Postgres via pgxpool")
	return pool, nil
}
