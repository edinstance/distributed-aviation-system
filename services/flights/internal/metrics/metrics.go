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

	KafkaMessagesProduced  metric.Int64Counter
	KafkaMessagesSent      metric.Int64Counter
	KafkaMessagesErrors    metric.Int64Counter
	KafkaSerializationTime metric.Float64Histogram
	KafkaProducerLatency   metric.Float64Histogram
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

	KafkaMessagesProduced, err = meter.Int64Counter(
		"flights.kafka.messages.produced",
		metric.WithDescription("Total number of messages produced to Kafka"),
	)
	if err != nil {
		return err
	}

	KafkaMessagesSent, err = meter.Int64Counter(
		"flights.kafka.messages.sent",
		metric.WithDescription("Total number of messages successfully sent to Kafka"),
	)
	if err != nil {
		return err
	}

	KafkaMessagesErrors, err = meter.Int64Counter(
		"flights.kafka.messages.errors",
		metric.WithDescription("Total number of Kafka message errors"),
	)
	if err != nil {
		return err
	}

	KafkaSerializationTime, err = meter.Float64Histogram(
		"flights.kafka.serialization.duration.seconds",
		metric.WithDescription("Time spent serializing messages to Avro"),
	)
	if err != nil {
		return err
	}

	KafkaProducerLatency, err = meter.Float64Histogram(
		"flights.kafka.producer.latency.seconds",
		metric.WithDescription("Kafka producer latency from produce to delivery confirmation"),
	)
	if err != nil {
		return err
	}

	return nil
}
