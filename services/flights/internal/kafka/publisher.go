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
	done       chan struct{}
}

// NewPublisher initializes Kafka producer and Avro serializer
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
		logger.ErrorContext(context.Background(), "Failed to create Kafka producer",
			"err", err, "broker_url", brokerURL)
		return nil, fmt.Errorf("create producer: %w", err)
	}
	defer func() {
		if err != nil {
			prod.Close()
		}
	}()

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
	sConfig.AutoRegisterSchemas = true
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
		done:       make(chan struct{}),
	}

	// Start background handler to consume delivery events (avoids buildup)
	go pub.handleDeliveryEvents()

	return pub, nil
}

// Close flushes pending messages and stops the background handler
func (p *Publisher) Close() {
	if p.producer == nil {
		return
	}

	remaining := p.producer.Flush(5000)
	if remaining > 0 {
		logger.Error("Failed to flush all messages", "remaining", remaining)
	}
	p.producer.Close()

	<-p.done
}

// recordKafkaError tracks Kafka-related errors with metrics and logs
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
	logger.ErrorContext(ctx, message,
		"err", err,
		"topic", topic,
		"event_type", eventType,
		"error_type", kind)
}

// handleDeliveryEvents consumes and classifies Kafka delivery / error events
func (p *Publisher) handleDeliveryEvents() {
	defer close(p.done)

	for e := range p.producer.Events() {
		switch m := e.(type) {
		case *kafka.Message:
			topic := "unknown"
			if m.TopicPartition.Topic != nil {
				topic = *m.TopicPartition.Topic
			}

			if m.TopicPartition.Error != nil {
				recordKafkaError(
					context.Background(),
					m.TopicPartition.Error,
					"delivery_error",
					"Kafka message delivery failed",
					topic,
					"",
				)
				continue
			}

			metrics.KafkaMessagesSent.Add(context.Background(), 1,
				metric.WithAttributes(
					attribute.String("topic", topic),
					attribute.Int("partition", int(m.TopicPartition.Partition)),
				))

			logger.Info("Kafka message delivered",
				"topic", topic,
				"partition", m.TopicPartition.Partition,
				"offset", m.TopicPartition.Offset)

		case kafka.Error:
			ctx := context.Background()
			errCode := m.Code().String()

			attrs := []attribute.KeyValue{
				attribute.String("topic", p.topic),
				attribute.String("error_code", errCode),
				attribute.Bool("is_fatal", m.IsFatal()),
				attribute.Bool("is_retriable", m.IsRetriable()),
			}

			switch {
			case m.IsFatal():
				recordKafkaError(ctx, m,
					"fatal_error", "Kafka producer fatal error", p.topic, "")
				logger.Error("Kafka fatal error encountered — producer must be recreated",
					"code", errCode, "err", m)
			case m.IsRetriable():
				recordKafkaError(ctx, m,
					"retriable_error", "Kafka retriable transient error", p.topic, "")
				logger.Warn("Kafka transient retriable error", "code", errCode, "err", m)
			default:
				recordKafkaError(ctx, m,
					"nonfatal_error", "Kafka non-fatal or local error", p.topic, "")
				logger.Warn("Kafka local/non-fatal error", "code", errCode, "err", m)
			}

			metrics.KafkaMessagesErrors.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
	}
}
