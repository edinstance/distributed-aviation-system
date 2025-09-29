package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GrpcRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total gRPC/Connect requests",
		},
		[]string{"service", "procedure", "status"},
	)

	GrpcDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Latency of gRPC/Connect requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "procedure"},
	)
)
