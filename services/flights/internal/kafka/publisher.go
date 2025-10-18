package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

	pub := &Publisher{
		producer:   prod,
		serializer: serializer,
		topic:      topic,
		tracer:     tracer,
	}

	go pub.handleDeliveryEvents()

	return pub, nil
}

func (p *Publisher) Close() {
	if p.producer != nil {
		p.producer.Flush(5000)
		p.producer.Close()
	}
}

// recordKafkaError centralizes error metric + logging
func recordKafkaError(
	ctx context.Context,
	err error,
	kind string,
	message string,
	topic string,
	eventType string,
) {
	metrics.KafkaMessagesErrors.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("topic", topic),
			attribute.String("event_type", eventType),
			attribute.String("error_type", kind),
		))
	logger.ErrorContext(ctx, message, "err", err, "topic", topic, "event_type", eventType)
}

// handleDeliveryEvents listens asynchronously for confirmation events
func (p *Publisher) handleDeliveryEvents() {
	for e := range p.producer.Events() {
		switch m := e.(type) {
		case *kafka.Message:
			if m.TopicPartition.Error != nil {
				recordKafkaError(context.Background(), m.TopicPartition.Error,
					"delivery_error", "Kafka message delivery failed",
					*m.TopicPartition.Topic, "")
				continue
			}

			metrics.KafkaMessagesSent.Add(context.Background(), 1,
				metric.WithAttributes(
					attribute.String("topic", *m.TopicPartition.Topic),
					attribute.Int("partition", int(m.TopicPartition.Partition)),
				))
			logger.Info("Kafka message delivered",
				"topic", *m.TopicPartition.Topic,
				"partition", m.TopicPartition.Partition,
				"offset", m.TopicPartition.Offset)
		case kafka.Error:
			recordKafkaError(context.Background(), m, "producer_error", "Kafka producer fatal error", "unknown", "")
		}
	}
}
