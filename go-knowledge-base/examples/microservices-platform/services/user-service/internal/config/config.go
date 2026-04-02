package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port       string
	GinMode    string
	LogLevel   string
	DatabaseURL string
	RedisURL   string
	JWTSecret  string
	JWTExpiry  time.Duration
}

func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8081"),
		GinMode:    getEnv("GIN_MODE", "release"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/users?sslmode=disable"),
		RedisURL:   getEnv("REDIS_URL", "localhost:6379"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry:  getDurationEnv("JWT_EXPIRY", 24*time.Hour),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
