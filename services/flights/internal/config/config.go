package config

import (
	"os"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
)

type Config struct {
	Port        string
	Environment string
	DatabaseURL string
}

var App Config

func Init() {
	App = Config{
		Port:        getEnv("PORT", "8081"),
		Environment: getEnv("ENVIRONMENT", "development"),
		DatabaseURL: mustGetEnv("DATABASE_URL"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func mustGetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	logger.Error("Required environment variable '%s' is missing or empty", key)
	os.Exit(1)
	return ""
}
