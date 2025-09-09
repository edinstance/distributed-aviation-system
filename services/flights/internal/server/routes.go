package server

import (
	"net/http"

	v1connect "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1/flightsv1connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/health"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewMux(pool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	// Register Connect/gRPC/gRPC-Web handlers
	flightsServer := NewFlightsServer(pool)
	flightPath, flightHandler := v1connect.NewFlightsServiceHandler(flightsServer)
	mux.Handle(flightPath, flightHandler)

	// Health check
	mux.HandleFunc("/health", health.HealthHandler)

	return mux
}
