package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"sensor-data-service.backend/config" // đổi thành tên module của bạn
)

func InitDB(cfg config.DBConfig) (*sqlx.DB, error) {
	// Validate config
	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" || cfg.Name == "" {
		return nil, fmt.Errorf("missing database config values")
	}

	// Build connection string
	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	// Connect
	db, err := sqlx.Connect(cfg.Driver, dbURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
		return nil, err
	}

	db.SetMaxOpenConns(25)
	log.Println("Connected to database")
	return db, nil
}
