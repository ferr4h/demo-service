package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort   string
	DBPath       string
	JWTSecret    string
	JWTExpiry    time.Duration
	RateLimitRPS int
	LogLevel     string
	LogFormat    string
}

var AppConfig *Config

func Load() error {
	_ = godotenv.Load()

	AppConfig = &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DBPath:       getEnv("DB_PATH", "./data/demo.db"),
		JWTSecret:    getEnv("JWT_SECRET", "1"),
		JWTExpiry:    parseDuration(getEnv("JWT_EXPIRY", "24h")),
		RateLimitRPS: parseInt(getEnv("RATE_LIMIT_RPS", "10")),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		LogFormat:    getEnv("LOG_FORMAT", "text"),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 10 // default
	}
	return val
}

func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour // default
	}
	return duration
}
