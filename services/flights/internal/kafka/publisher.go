package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
)

// FlightEvent represents the Avro structure (same as flight_created.avsc)
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

// Publisher handles Avro-based message publishing
type Publisher struct {
	producer   *kafka.Producer
	serializer *avro.GenericSerializer
	topic      string
}

// NewPublisher initializes Kafka producer + Avro serializer
func NewPublisher(brokerURL, schemaRegistryURL, topic string) (*Publisher, error) {
	logger.InfoContext(context.Background(), "Initializing Kafka publisher",
		"broker_url", brokerURL,
		"schema_registry_url", schemaRegistryURL,
		"topic", topic)

	// Kafka Producer config
	prod, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":       brokerURL,
		"client.id":               "flights-service",
		"acks":                    "all",
		"go.delivery.reports":     true,
		"retries":                 5,
		"retry.backoff.ms":        100,
		"request.timeout.ms":      30000,
		"delivery.timeout.ms":     60000,
		"socket.keepalive.enable": true,
	})
	if err != nil {
		logger.ErrorContext(context.Background(), "Failed to create Kafka producer", "err", err, "broker_url", brokerURL)
		return nil, fmt.Errorf("create producer: %w", err)
	}

	logger.InfoContext(context.Background(), "Kafka producer created successfully")

	// Schema Registry client
	logger.InfoContext(context.Background(), "Connecting to Schema Registry", "url", schemaRegistryURL)
	srClient, err := schemaregistry.NewClient(schemaregistry.NewConfig(schemaRegistryURL))
	if err != nil {
		logger.ErrorContext(context.Background(), "Failed to create Schema Registry client", "err", err, "url", schemaRegistryURL)
		return nil, fmt.Errorf("create schema registry client: %w", err)
	}

	// Avro Serializer configuration
	sConfig := avro.NewSerializerConfig()
	sConfig.AutoRegisterSchemas = true // auto-register the schema on first use
	sConfig.UseLatestVersion = true

	serializer, err := avro.NewGenericSerializer(srClient, serde.ValueSerde, sConfig)
	if err != nil {
		logger.ErrorContext(context.Background(), "Failed to create Avro serializer", "err", err)
		return nil, fmt.Errorf("create Avro serializer: %w", err)
	}

	logger.InfoContext(context.Background(), "Kafka publisher initialized successfully")
	return &Publisher{
		producer:   prod,
		serializer: serializer,
		topic:      topic,
	}, nil
}

// PublishFlightCreated serializes FlightCreated event as Avro and sends it
func (p *Publisher) PublishFlightCreated(ctx context.Context, flight *models.Flight) error {
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

	// Serialize the event to Avro binary payload
	valueBytes, err := p.serializer.Serialize(p.topic, &event)
	if err != nil {
		logger.ErrorContext(ctx, "Avro serialization failed", "err", err, "flight_id", flight.ID)
		return fmt.Errorf("avro serialization failed: %w", err)
	}

	logger.DebugContext(ctx, "Event serialized to Avro", "size_bytes", len(valueBytes), "flight_id", flight.ID)

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Key:            []byte(flight.ID.String()),
		Value:          valueBytes,
		Headers:        []kafka.Header{{Key: "eventType", Value: []byte("FlightCreated")}},
		Timestamp:      time.Now(),
	}

	if err := p.producer.Produce(msg, nil); err != nil {
		logger.ErrorContext(ctx, "Failed to produce message to Kafka", "err", err, "flight_id", flight.ID)
		return fmt.Errorf("kafka produce failed: %w", err)
	}

	logger.DebugContext(ctx, "Message sent to Kafka, waiting for delivery confirmation", "flight_id", flight.ID)

	select {
	case e := <-p.producer.Events():
		switch m := e.(type) {
		case *kafka.Message:
			if m.TopicPartition.Error != nil {
				logger.ErrorContext(ctx, "Delivery failed", "err", m.TopicPartition.Error)
				return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
			}
			logger.InfoContext(ctx, "FlightCreated Avro published",
				"flight_id", flight.ID,
				"topic", *m.TopicPartition.Topic,
				"partition", m.TopicPartition.Partition,
				"offset", m.TopicPartition.Offset)
		case kafka.Error:
			return fmt.Errorf("kafka producer error: %v", m)
		}
	case <-time.After(5 * time.Second):
		logger.WarnContext(ctx, "Timeout waiting for Kafka delivery confirmation", "flight_id", flight.ID)
		return fmt.Errorf("timeout waiting for delivery confirmation")
	case <-ctx.Done():
		logger.WarnContext(ctx, "Context cancelled while waiting for Kafka delivery", "flight_id", flight.ID)
		return ctx.Err()
	}
	return nil
}

func (p *Publisher) Close() {
	if p.producer != nil {
		p.producer.Flush(5000)
		p.producer.Close()
	}
}
