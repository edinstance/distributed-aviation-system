package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	port := os.Getenv("PORT")
	env := os.Getenv("ENVIRONMENT")

	logger.Init(env)
	mux := server.NewMux()

	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	logger.Info("FlightsService listening on %s", addr)
	if err := http.ListenAndServe(addr, h2c.NewHandler(mux, &http2.Server{})); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server error: %v", err)
	}
}
