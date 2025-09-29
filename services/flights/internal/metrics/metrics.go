package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GrpcRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flights_grpc_requests_total",
			Help: "Total gRPC/Connect requests for the flights service",
		},
		[]string{"direction", "procedure", "service", "status"},
	)

	GrpcDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "flights_grpc_request_duration_seconds",
			Help:    "Latency of gRPC/Connect requests for the flights service",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"direction", "procedure", "service"},
	)

	GraphQLRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flights_graphql_requests_total",
			Help: "Total GraphQL requests for the flights service",
		},
		[]string{"operation", "type", "success"},
	)

	GraphQLDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "flights_graphql_request_duration_seconds",
			Help:    "Duration of GraphQL requests for the flights service",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "type"},
	)
)
