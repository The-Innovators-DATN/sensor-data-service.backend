package main

import (
	"sensor-data-service.backend/pkg/db"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	_, err := db.InitDB()
	if err != nil {
		panic(err)
	}
}