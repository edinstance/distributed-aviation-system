package server

import (
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/config"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/metrics"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/middleware"
	v1connect "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1/flightsv1connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/health"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

func NewMux(pool *pgxpool.Pool, client *redis.Client) *http.ServeMux {
	traceInterceptor, err := otelconnect.NewInterceptor(
		otelconnect.WithTracerProvider(otel.GetTracerProvider()),
	)

	if err != nil {
		logger.Warn("Failed to create trace interceptor", "err", err)
		traceInterceptor = nil
	}

	mux := http.NewServeMux()

	interceptors := []connect.Interceptor{metrics.GrpcMetricsInterceptor{}}
	if traceInterceptor != nil {
		interceptors = append([]connect.Interceptor{traceInterceptor}, interceptors...)
	}

	// Register Connect/gRPC/gRPC-Web handlers
	grpcFlightsServer := NewGrpcFlightsServer(pool, client)
	flightPath, flightHandler := v1connect.NewFlightsServiceHandler(
		grpcFlightsServer,
		connect.WithInterceptors(interceptors...),
	)

	mux.Handle(flightPath, flightHandler)

	// GraphQL handlers
	mux.Handle("/graphql", middleware.UserContextMiddleware(newGraphQLHandler(pool, client)))

	if config.App.Environment != "prod" {
		mux.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
	}

	// Health check
	mux.HandleFunc("/health", health.HealthHandler)

	return mux
}
