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

// NewConfig gets environment variable with fallback, fatal if required and missing
func NewConfig(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback != "" {
		log.Printf("Using fallback for %s: %s", key, fallback)
		return fallback
	}
	log.Fatalf("Required environment variable %s is missing", key)
	return ""
}

// GetConfigOptional gets environment variable with fallback, returns empty string if missing
func GetConfigOptional(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
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
		// Required variables
		PostgresDSN: NewConfig("POSTGRES_DB", ""),
		PORT:        NewConfig("PORT", "8080"),
		JWTSecret:   NewConfig("JWT_SECRET", ""),

		// Redis - optional credentials (can be empty for local Redis)
		RedisHost: GetConfigOptional("REDIS_HOST", "localhost:6379"),
		RedisUser: GetConfigOptional("REDIS_USER", ""),
		RedisPass: GetConfigOptional("REDIS_PASSWORD", ""),

		Email: EmailConfig{
			// SMTP - required for email functionality
			SenderEmail:   NewConfig("SENDER_EMAIL", ""),
			SMTPHost:      NewConfig("SMTP_HOST", ""),
			SMTPUser:      NewConfig("SMTP_USER", ""),
			SMTPPass:      NewConfig("SMTP_PASS", ""),
			SMTPPort:      getEnvAsInt("SMTP_PORT", 587),
			// BaseVerifyURL is optional (not used in OTP flow)
			BaseVerifyURL: GetConfigOptional("BASE_VERIFY_URL", ""),
		},
	}
}
