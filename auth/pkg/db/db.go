package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/Mahaveer86619/ms/auth/pkg/config"
)

var DB *sql.DB

const (
	maxRetries = 5
	retryDelay = 6 * time.Second
)

func InitDB() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.GConfig.DBUser,
		config.GConfig.DBPassword,
		config.GConfig.DBHost,
		config.GConfig.DBPort,
		config.GConfig.DBName,
	)

	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	for i := 1; i <= maxRetries; i++ {
		log.Printf("Attempt %d of %d: Pinging database...", i, maxRetries)

		err = DB.Ping()
		if err == nil {
			log.Println("Successfully connected to the database!")
			return
		}

		log.Printf("Failed to connect to database (attempt %d): %v", i, err)

		if i < maxRetries {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	totalTime := (maxRetries - 1) * retryDelay
	log.Fatalf("Failed to connect to the database after %d attempts (over ~%v). Last error: %v", maxRetries, totalTime.Round(time.Second), err)
}
