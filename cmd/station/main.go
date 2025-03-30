package main

import (
	"log"

	"sensor-data-service.backend/config"
	"sensor-data-service.backend/internal/db"
)

func main() {
	cfg, err := config.LoadAllConfigs("config")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Println("App:", cfg.App)
	log.Println("Clickhouse:", cfg.Clickhouse)
	database, err := db.InitDB(cfg.DB) // inject config.DB
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	clickhouse, err := db.InitClickHouse(cfg.Clickhouse) // inject config.Clickhouse
	if err != nil {
		log.Fatal(err)
	}
	defer clickhouse.Close()
	// if err := godotenv.Load("../.env"); err != nil {
	// 	panic("Error loading .env file")
	// }

	// cfg, err := config.LoadAllConfigs("../configs")
	// if err != nil {
	// 	panic(err)
	// }

	// dbConn, err := db.InitDB()
	// if err != nil {
	// 	panic(err)
	// }
	// defer dbConn.Close()

	// listener, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	panic(err)
	// }
	// defer listener.Close()

	// // Start the server
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	go handleConnection(conn)
	// }
}
