package flights

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *flightCache) GetFlight(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	tracer := otel.Tracer("flights-service")
	ctx, span := tracer.Start(ctx, "cache.get_flight")
	defer span.End()

	key := fmt.Sprintf("flight:%s", id.String())
	span.SetAttributes(
		attribute.String("cache.operation", "get"),
		attribute.String("cache.key", key),
		attribute.String("flight.id", id.String()),
	)

	val, err := r.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		span.SetAttributes(attribute.String("cache.result", "miss"))
		return nil, nil
	}

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("cache.result", "error"))
		return nil, fmt.Errorf("error getting data from the cache: %w", err)
	}

	var flight models.Flight

	if err := json.Unmarshal([]byte(val), &flight); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("cache.result", "unmarshal_error"))
		return nil, fmt.Errorf("error converting cache data: %w", err)
	}

	span.SetAttributes(
		attribute.String("cache.result", "hit"),
		attribute.String("flight.number", flight.Number),
	)

	return &flight, nil
}
