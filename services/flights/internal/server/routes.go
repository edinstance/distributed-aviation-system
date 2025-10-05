package server

import (
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/config"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/metrics"
	v1connect "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1/flightsv1connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/health"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

func NewMux(pool *pgxpool.Pool, client *redis.Client) *http.ServeMux {
	traceInterceptor, _ := otelconnect.NewInterceptor(
		otelconnect.WithTracerProvider(otel.GetTracerProvider()),
	)

	mux := http.NewServeMux()

	// Register Connect/gRPC/gRPC-Web handlers
	grpcFlightsServer := NewGrpcFlightsServer(pool, client)
	flightPath, flightHandler := v1connect.NewFlightsServiceHandler(
		grpcFlightsServer,
		connect.WithInterceptors(traceInterceptor, metrics.GrpcMetricsInterceptor{}),
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
