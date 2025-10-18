package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Publisher handles Avro-based message publishing
type Publisher struct {
	producer   *kafka.Producer
	serializer *avro.GenericSerializer
	topic      string
	tracer     trace.Tracer
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

	// Initialize OpenTelemetry tracer
	tracer := otel.Tracer("kafka-publisher")

	return &Publisher{
		producer:   prod,
		serializer: serializer,
		topic:      topic,
		tracer:     tracer,
	}, nil
}

func (p *Publisher) Close() {
	if p.producer != nil {
		p.producer.Flush(5000)
		p.producer.Close()
	}
}
