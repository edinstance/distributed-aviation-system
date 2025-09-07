package server

import (
	"net/http"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/health"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", health.HealthHandler)

	return mux
}
