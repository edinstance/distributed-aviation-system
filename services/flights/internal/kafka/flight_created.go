package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// FlightCreated represents the Avro structure for a created flight
type FlightCreated struct {
	FlightId      string `avro:"flightId"`
	Number        string `avro:"number"`
	Origin        string `avro:"origin"`
	Destination   string `avro:"destination"`
	DepartureTime string `avro:"departureTime"`
	ArrivalTime   string `avro:"arrivalTime"`
	Airline       string `avro:"airline"`
	Status        string `avro:"status"`
}

// PublishFlightCreated serializes the FlightCreated event as Avro and sends it to Kafka
func (p *Publisher) PublishFlightCreated(ctx context.Context, flight *models.Flight) error {
	const eventType = "FlightCreated"

	ctx, span := p.tracer.Start(ctx, "kafka.publish",
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination.name", p.topic),
			attribute.String("messaging.destination_kind", "topic"),
			attribute.String("messaging.operation", "publish"),
			attribute.String("event.type", eventType),
			attribute.String("flight.id", flight.ID.String()),
			attribute.String("flight.number", flight.Number),
		))
	defer span.End()

	logFields := []any{
		"topic", p.topic,
		"flight_id", flight.ID,
		"number", flight.Number,
		"origin", flight.Origin,
		"destination", flight.Destination,
	}

	logger.InfoContext(ctx, "Publishing FlightCreated event", logFields...)

	metrics.KafkaMessagesProduced.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("topic", p.topic),
			attribute.String("event_type", eventType),
		))

	event := FlightCreated{
		FlightId:      flight.ID.String(),
		Number:        flight.Number,
		Origin:        flight.Origin,
		Destination:   flight.Destination,
		DepartureTime: flight.DepartureTime.Format(time.RFC3339),
		ArrivalTime:   flight.ArrivalTime.Format(time.RFC3339),
		Airline:       flight.Number[:2],
		Status:        string(flight.Status),
	}

	start := time.Now()
	valueBytes, err := p.serializer.Serialize(p.topic, &event)
	serializationDuration := time.Since(start)

	logger.DebugContext(ctx, "Avro serialization result",
		"bytes_length", len(valueBytes),
		"serialization_duration_ms", serializationDuration.Milliseconds())

	metrics.KafkaSerializationTime.Record(
		ctx,
		float64(serializationDuration.Microseconds())/1000.0,
		metric.WithAttributes(
			attribute.String("topic", p.topic),
			attribute.String("event_type", eventType),
			attribute.String("unit", "ms"),
		),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Avro serialization failed")
		recordKafkaError(ctx, err, "serialization_error", "Avro serialization failed", p.topic, eventType)
		return fmt.Errorf("avro serialization failed: %w", err)
	}

	headers := []kafka.Header{{Key: "eventType", Value: []byte(eventType)}}
	carrier := make(propagation.MapCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	for k, v := range carrier {
		headers = append(headers, kafka.Header{Key: k, Value: []byte(v)})
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Key:       []byte(flight.ID.String()),
		Value:     valueBytes,
		Headers:   headers,
		Timestamp: time.Now(),
	}

	// --- Non-blocking attempt using goroutine and timeout ---
	produceStart := time.Now()
	done := make(chan error, 1)

	go func() {
		err := p.producer.Produce(msg, nil)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Produce failed")
			recordKafkaError(ctx, err, "produce_error", "Kafka produce failed", p.topic, eventType)
			return fmt.Errorf("produce failed: %w", err)
		}
	case <-time.After(2 * time.Second):
		err := fmt.Errorf("kafka produce timeout (buffer full)")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Producer queue full")
		recordKafkaError(ctx, err, "produce_timeout", "Kafka produce timeout (buffer full)", p.topic, eventType)
		return err
	case <-ctx.Done():
		recordKafkaError(ctx, ctx.Err(), "context_cancelled", "Context cancelled producing Kafka message", p.topic, eventType)
		return ctx.Err()
	}

	enqLatency := time.Since(produceStart)
	metrics.KafkaProducerLatency.Record(
		ctx,
		enqLatency.Seconds(),
		metric.WithAttributes(
			attribute.String("topic", p.topic),
			attribute.String("event_type", eventType),
		),
	)

	logger.DebugContext(ctx, "Message accepted by Kafka producer buffer", logFields...)
	return nil
}
