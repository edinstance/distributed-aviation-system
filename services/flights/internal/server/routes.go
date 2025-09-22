package server

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/config"
	v1connect "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1/flightsv1connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/health"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewMux creates and returns an *http.ServeMux pre-configured with the service's HTTP endpoints.
// It registers the Flights connect/gRPC (and gRPC-Web) handler, the GraphQL endpoint at "/graphql",
// a health check at "/health" and, when the environment is not "prod", the GraphQL Playground at "/playground".
// The pool argument is the Postgres connection pool used to construct the Flights server and GraphQL handler.
func NewMux(pool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	// Register Connect/gRPC/gRPC-Web handlers
	flightsServer := NewFlightsServer(pool)
	flightPath, flightHandler := v1connect.NewFlightsServiceHandler(flightsServer)
	mux.Handle(flightPath, flightHandler)

	// GraphQL handlers
	mux.Handle("/graphql", newGraphQLHandler(pool))

	if config.App.Environment != "prod" {
		mux.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
	}

	// Health check
	mux.HandleFunc("/health", health.HealthHandler)

	return mux
}
