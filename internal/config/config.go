package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN string
	PORT        string
	JWTSecret   string
	RedisHost   string
}

func NewConfig(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback != "" {
		log.Printf("Fallback env variable is %s", fallback)
		return fallback
	}
	log.Fatal("Fallback env variable is required")
	return ""
}

func NewConfigFromEnv() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		PostgresDSN: NewConfig("POSTGRES_DB", ""),
		PORT:        NewConfig("PORT", "8080"),
		JWTSecret:   NewConfig("JWT_SECRET", ""),
		RedisHost:   NewConfig("REDIS_HOST", "localhost"),
	}
}
