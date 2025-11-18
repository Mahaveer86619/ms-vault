package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	// DB config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Auth
	JWTSecret string
}

var GConfig *Config

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GConfig = &Config{
		Port: getEnv("PORT"),

		// DB config
		DBHost:     getEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME"),

		// JWT config
		JWTSecret:  getEnv("JWT_SECRET"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Key %s not found in .env file", key)
	}

	return val
}
