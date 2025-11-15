package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN string
	PORT        string
	JWTSecret   string
	RedisHost   string
	RedisUser   string
	RedisPass   string
	Email       EmailConfig
}

type EmailConfig struct {
	SenderEmail   string
	SMTPHost      string
	SMTPUser      string
	SMTPPass      string
	SMTPPort      int
	BaseVerifyURL string
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

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
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
		RedisUser:   NewConfig("REDIS_USER", ""),
		RedisPass:   NewConfig("REDIS_PASSWORD", ""),

		Email: EmailConfig{
			SenderEmail:   NewConfig("SENDER_EMAIL", "default@example.com"),
			SMTPHost:      NewConfig("SMTP_HOST", "localhost"),
			SMTPUser:      NewConfig("SMTP_USER", ""),
			SMTPPass:      NewConfig("SMTP_PASS", ""),
			SMTPPort:      getEnvAsInt("SMTP_PORT", 25),
			BaseVerifyURL: NewConfig("BASE_VERIFY_URL", ""),
		},
	}
}
