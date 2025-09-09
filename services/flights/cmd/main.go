package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/server"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		logger.Warn("No .env file found, relying on environment variables")
	}

	port := os.Getenv("PORT")
	env := os.Getenv("ENVIRONMENT")
	databaseURL := os.Getenv("DATABASE_URL")

	logger.Init(env)

	pool, err := database.Init(databaseURL)
	if err != nil {
		logger.Error("Failed to initialise database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	mux := server.NewMux(pool)

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
