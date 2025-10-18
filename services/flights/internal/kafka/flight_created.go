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

// FlightCreatedEvent represents the Avro structure (same as flight_created.avsc)
type FlightCreatedEvent struct {
	FlightID      string `avro:"flightId"`
	Number        string `avro:"number"`
	Origin        string `avro:"origin"`
	Destination   string `avro:"destination"`
	DepartureTime int64  `avro:"departureTime"`
	ArrivalTime   int64  `avro:"arrivalTime"`
	Airline       string `avro:"airline"`
	Status        string `avro:"status"`
}

// PublishFlightCreated serializes FlightCreated event as Avro and sends it
func (p *Publisher) PublishFlightCreated(ctx context.Context, flight *models.Flight) error {
	// Start tracing span
	ctx, span := p.tracer.Start(ctx, "kafka.publish",
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination.name", p.topic),
			attribute.String("messaging.operation", "publish"),
			attribute.String("flight.id", flight.ID.String()),
			attribute.String("flight.number", flight.Number),
			attribute.String("flight.origin", flight.Origin),
			attribute.String("flight.destination", flight.Destination),
		))
	defer span.End()

	// Record message production metric
	metrics.KafkaMessagesProduced.Add(ctx, 1, metric.WithAttributes(
		attribute.String("topic", p.topic),
		attribute.String("event_type", "FlightCreated"),
	))

	logger.InfoContext(ctx, "Publishing FlightCreated event", "flight_id", flight.ID, "number", flight.Number)

	event := FlightCreatedEvent{
		FlightID:      flight.ID.String(),
		Number:        flight.Number,
		Origin:        flight.Origin,
		Destination:   flight.Destination,
		DepartureTime: flight.DepartureTime.Unix(),
		ArrivalTime:   flight.ArrivalTime.Unix(),
		Airline:       flight.Number[:2],
		Status:        string(flight.Status),
	}

	// Serialize the event to Avro binary payload with timing
	serializationStart := time.Now()
	valueBytes, err := p.serializer.Serialize(p.topic, &event)
	serializationDuration := time.Since(serializationStart)
	metrics.KafkaSerializationTime.Record(ctx, serializationDuration.Seconds(),
		metric.WithAttributes(
			attribute.String("topic", p.topic),
			attribute.String("event_type", "FlightCreated"),
		))

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Avro serialization failed")
		metrics.KafkaMessagesErrors.Add(ctx, 1, metric.WithAttributes(
			attribute.String("topic", p.topic),
			attribute.String("error_type", "serialization_error"),
		))
		logger.ErrorContext(ctx, "Avro serialization failed", "err", err, "flight_id", flight.ID)
		return fmt.Errorf("avro serialization failed: %w", err)
	}

	logger.DebugContext(ctx, "Event serialized to Avro", "size_bytes", len(valueBytes), "flight_id", flight.ID)

	// Inject trace context into Kafka message headers
	headers := []kafka.Header{{Key: "eventType", Value: []byte("FlightCreated")}}
	carrier := make(propagation.MapCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	for key, value := range carrier {
		headers = append(headers, kafka.Header{Key: key, Value: []byte(value)})
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Key:            []byte(flight.ID.String()),
		Value:          valueBytes,
		Headers:        headers,
		Timestamp:      time.Now(),
	}

	produceStart := time.Now()
	if err := p.producer.Produce(msg, nil); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Kafka produce failed")
		metrics.KafkaMessagesErrors.Add(ctx, 1, metric.WithAttributes(
			attribute.String("topic", p.topic),
			attribute.String("error_type", "produce_error"),
		))
		logger.ErrorContext(ctx, "Failed to produce message to Kafka", "err", err, "flight_id", flight.ID)
		return fmt.Errorf("kafka produce failed: %w", err)
	}

	logger.DebugContext(ctx, "Message sent to Kafka, waiting for delivery confirmation", "flight_id", flight.ID)

	select {
	case e := <-p.producer.Events():
		producerDuration := time.Since(produceStart)
		metrics.KafkaProducerLatency.Record(ctx, producerDuration.Seconds(),
			metric.WithAttributes(
				attribute.String("topic", p.topic),
				attribute.String("event_type", "FlightCreated"),
			))

		switch m := e.(type) {
		case *kafka.Message:
			if m.TopicPartition.Error != nil {
				span.RecordError(m.TopicPartition.Error)
				span.SetStatus(codes.Error, "Delivery failed")
				metrics.KafkaMessagesErrors.Add(ctx, 1, metric.WithAttributes(
					attribute.String("topic", p.topic),
					attribute.String("error_type", "delivery_error"),
				))
				logger.ErrorContext(ctx, "Delivery failed", "err", m.TopicPartition.Error)
				return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
			}

			// Record successful delivery
			span.SetAttributes(
				attribute.Int("messaging.kafka.partition", int(m.TopicPartition.Partition)),
				attribute.Int64("messaging.kafka.offset", int64(m.TopicPartition.Offset)),
			)
			span.SetStatus(codes.Ok, "Message delivered successfully")
			metrics.KafkaMessagesSent.Add(ctx, 1, metric.WithAttributes(
				attribute.String("topic", p.topic),
				attribute.String("event_type", "FlightCreated"),
				attribute.Int("partition", int(m.TopicPartition.Partition)),
			))

			logger.InfoContext(ctx, "FlightCreated Avro published",
				"flight_id", flight.ID,
				"topic", *m.TopicPartition.Topic,
				"partition", m.TopicPartition.Partition,
				"offset", m.TopicPartition.Offset)
		case kafka.Error:
			span.RecordError(m)
			span.SetStatus(codes.Error, "Kafka producer error")
			metrics.KafkaMessagesErrors.Add(ctx, 1, metric.WithAttributes(
				attribute.String("topic", p.topic),
				attribute.String("error_type", "producer_error"),
			))
			return fmt.Errorf("kafka producer error: %v", m)
		}
	case <-time.After(5 * time.Second):
		span.SetStatus(codes.Error, "Timeout waiting for delivery confirmation")
		metrics.KafkaMessagesErrors.Add(ctx, 1, metric.WithAttributes(
			attribute.String("topic", p.topic),
			attribute.String("error_type", "delivery_timeout"),
		))
		logger.WarnContext(ctx, "Timeout waiting for Kafka delivery confirmation", "flight_id", flight.ID)
		return fmt.Errorf("timeout waiting for delivery confirmation")
	case <-ctx.Done():
		span.SetStatus(codes.Error, "Context cancelled")
		metrics.KafkaMessagesErrors.Add(ctx, 1, metric.WithAttributes(
			attribute.String("topic", p.topic),
			attribute.String("error_type", "context_cancelled"),
		))
		logger.WarnContext(ctx, "Context cancelled while waiting for Kafka delivery", "flight_id", flight.ID)
		return ctx.Err()
	}
	return nil
}
