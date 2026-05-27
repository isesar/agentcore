package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	// Server settings
	ServerPort string

	// Database settings
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Logging settings
	LogLevel string

	// Router settings
	RouterEnableLLM                bool
	RouterRulesConfidenceThreshold float64
)

func Load() error {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load server settings
	ServerPort = getEnv("SERVER_PORT", "8080")

	// Load database settings
	DBHost = getEnv("DB_HOST", "localhost")
	DBPort = getEnv("DB_PORT", "5432")
	DBUser = getEnv("DB_USER", "postgres")
	DBPassword = getEnv("DB_PASSWORD", "password")
	DBName = getEnv("DB_NAME", "agentcore")
	DBSSLMode = getEnv("DB_SSL_MODE", "disable")

	// Load logging settings
	LogLevel = getEnv("LOG_LEVEL", "info")

	// Load router settings
	RouterEnableLLM = getEnv("ROUTER_ENABLE_LLM", "false") == "true"
	RouterRulesConfidenceThreshold = getEnvAsFloat("ROUTER_RULES_CONFIDENCE_THRESHOLD", 0.80)

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return parsed
}
