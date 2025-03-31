package db

// Package db provides functions to initialize and manage Redis connections.

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"sensor-data-service.backend/config"
)

func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		Protocol: cfg.Protocol,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	fmt.Println("Connected to Redis")
	return rdb, nil
}
