package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/server"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, relying on environment variables")
	}

	port := os.Getenv("PORT")
	env := os.Getenv("ENVIRONMENT")

	logger.Init(env)
	mux := server.NewMux()

	if port == "" {
		port = "8081"
	}
	addr := ":" + port

	logger.Info("FlightsService listening", "addr", addr)
	if err := http.ListenAndServe(addr, h2c.NewHandler(mux, &http2.Server{})); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server error", "err", err)
	}
}
