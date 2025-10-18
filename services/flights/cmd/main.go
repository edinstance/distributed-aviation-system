package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/cache"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/config"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/kafka"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/metrics"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/server"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/tracing"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		logger.Warn("No .env file found, relying on environment variables")
	}

	ctx := context.Background()

	config.Init()

	provider, err := logger.Init(config.App.Environment)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = provider.Shutdown(context.Background()) }()

	shutdownTracing, err := tracing.Init(ctx, "flights-service", config.App.OtlpGrpcUrl)
	if err != nil {
		logger.Error("failed to init tracing", "err", err)
		os.Exit(1)
	}
	defer func() {
		_ = shutdownTracing(ctx)
	}()

	shutdownMetrics, err := metrics.Init(ctx, "flights-service", config.App.OtlpGrpcUrl)
	if err != nil {
		logger.Error("failed to init metrics: %v", err)
		os.Exit(1)
	}
	defer func() {
		_ = shutdownMetrics(ctx)
	}()

	pool, err := database.Init(config.App.DatabaseURL)
	if err != nil {
		logger.Error("Failed to initialise database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	cacheClient, err := cache.Init(config.App.CacheURL)
	if err != nil {
		logger.Error("Failed to initialise cache", "err", err)
		cacheClient = nil
	}

	defer func() {
		if cacheClient != nil {
			if err := cacheClient.Close(); err != nil {
				logger.Warn("Error closing Redis client", "err", err)
			}
		}
	}()

	kafkaPublisher, err := kafka.NewPublisher(config.App.KafkaBrokerURL, config.App.KafkaSchemaRegistryURL, config.App.KafkaFlightsTopic)
	if err != nil {
		logger.Error("Failed to initialise Kafka publisher", "err", err)
		os.Exit(1)
	}
	defer kafkaPublisher.Close()

	mux := server.NewMux(pool, cacheClient, kafkaPublisher)

	port := config.App.Port
	if port == "" {
		port = "8081"
	}
	addr := ":" + port

	srv := &http.Server{
		Addr:         addr,
		Handler:      h2c.NewHandler(mux, &http2.Server{}),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("FlightsService listening", "addr", addr)
	logger.Debug("Environment", "env", config.App.Environment)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server error", "err", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "err", err)
	}
}
