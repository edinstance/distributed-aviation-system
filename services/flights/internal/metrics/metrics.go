package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter = otel.GetMeterProvider().Meter("flights-service")

	GrpcRequests    metric.Int64Counter
	GrpcDuration    metric.Float64Histogram
	GraphQLRequests metric.Int64Counter
	GraphQLDuration metric.Float64Histogram
)

func InitInstruments() error {
	var err error

	GrpcRequests, err = meter.Int64Counter(
		"flights.grpc.requests",
		metric.WithDescription("Total gRPC requests for the flights service"),
	)
	if err != nil {
		return err
	}

	GrpcDuration, err = meter.Float64Histogram(
		"flights.grpc.duration.seconds",
		metric.WithDescription("Latency of gRPC requests for the flights service"),
	)
	if err != nil {
		return err
	}

	GraphQLRequests, err = meter.Int64Counter(
		"flights.graphql.requests",
		metric.WithDescription("Total GraphQL requests for the flights service"),
	)
	if err != nil {
		return err
	}

	GraphQLDuration, err = meter.Float64Histogram(
		"flights.graphql.duration.seconds",
		metric.WithDescription("Duration of GraphQL requests for the flights service"),
	)
	if err != nil {
		return err
	}

	return nil
}
