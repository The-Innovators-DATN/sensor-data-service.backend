package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"sensor-data-service.backend/config"
)

func InitClickHouse(cfg config.ClickhouseConfig) (clickhouse.Conn, error) {
	opts := &clickhouse.Options{
		Addr: []string{cfg.Addr},
		Auth: clickhouse.Auth{
			Database: cfg.Name,
			Username: cfg.User,
			Password: cfg.Password,
		},
		Settings: map[string]interface{}{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
	}

	conn, err := clickhouse.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("clickhouse open error: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("clickhouse ping error: %w", err)
	}

	log.Println("Connected to ClickHouse")
	return conn, nil
}
