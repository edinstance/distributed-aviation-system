package server

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/config"
	v1connect "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1/flightsv1connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/health"
	"github.com/edinstance/distributed-aviation-system/services/flights/metrics"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func NewMux(pool *pgxpool.Pool, client *redis.Client) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	// Register Connect/gRPC/gRPC-Web handlers
	grpcFlightsServer := NewGrpcFlightsServer(pool, client)
	flightPath, flightHandler := v1connect.NewFlightsServiceHandler(
		grpcFlightsServer,
		connect.WithInterceptors(metrics.GrpcMetricsInterceptor{}),
	)

	mux.Handle(flightPath, flightHandler)

	// GraphQL handlers
	mux.Handle("/graphql", newGraphQLHandler(pool, client))

	if config.App.Environment != "prod" {
		mux.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
	}

	// Health check
	mux.HandleFunc("/health", health.HealthHandler)

	return mux
}
