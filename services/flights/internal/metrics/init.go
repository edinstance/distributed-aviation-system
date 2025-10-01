package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init(ctx context.Context, serviceName, otlpAddr string) (func(context.Context) error, error) {
	conn, err := grpc.NewClient(otlpAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(exporter,
				metric.WithInterval(5*time.Second)),
		),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(provider)

	if err := InitInstruments(); err != nil {
		return nil, fmt.Errorf("failed to init instruments: %w", err)
	}
	logger.Info("OTel Metrics started for %s", serviceName)

	return provider.Shutdown, nil
}
