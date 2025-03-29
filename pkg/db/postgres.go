package db

import (
	"fmt"
	"log"
	"os"

	_"github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/jmoiron/sqlx" // SQL Helper
)
// DB is a wrapper around sqlx.DB
func InitDB() (*sqlx.DB, error) {
	// Get the database URL from the environment variable
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if host == "" || port == "" || user == "" || password == "" || dbName == "" {
		log.Fatal("Database connection details are not set in environment variables")
		return nil, fmt.Errorf("database connection details are not set in environment variables")
	}

	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	// Connect to the PostgreSQL database
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil, err
	}

	// Set the maximum number of open connections
	db.SetMaxOpenConns(25)

	log.Println("Connected to the database successfully")
	return db, nil
}