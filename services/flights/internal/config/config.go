package config

import (
	"os"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
)

type Config struct {
	Port        string
	Environment string
	DatabaseURL string
	CacheURL    string
	CacheTTL    time.Duration
}

var App Config

// Init initialises the package-level App configuration from environment variables.
// PORT and ENVIRONMENT default to "8081" and "development" respectively if unset; DATABASE_URL is mandatory and the process will exit if it is missing or empty.
func Init() {
	App = Config{
		Port:        getEnv("PORT", "8081"),
		Environment: getEnv("ENVIRONMENT", "development"),
		DatabaseURL: mustGetEnv("DATABASE_URL"),
		CacheURL:    mustGetEnv("CACHE_URL"),
		CacheTTL:    15 * time.Minute,
	}
}

// getEnv returns the value of the environment variable named by key, or
// fallback if the variable is not present. If the variable is present but
// empty, the empty string is returned (i.e. presence is determined via
// os.LookupEnv).
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// mustGetEnv returns the value of the environment variable named by key.
// If the variable is unset or empty it logs an error and terminates the process
// with exit status 1. The function does not return on error.
func mustGetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	logger.Error("Required environment variable '%s' is missing or empty", key)
	os.Exit(1)
	return ""
}
