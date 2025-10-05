package flights

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (r *flightCache) SetFlight(ctx context.Context, flight *models.Flight) error {
	tracer := otel.Tracer("flights-service")
	ctx, span := tracer.Start(ctx, "cache.set_flight")
	defer span.End()

	key := fmt.Sprintf("flight:%s", flight.ID.String())
	span.SetAttributes(
		attribute.String("cache.operation", "set"),
		attribute.String("cache.key", key),
		attribute.String("flight.id", flight.ID.String()),
		attribute.String("flight.number", flight.Number),
	)

	data, err := json.Marshal(flight)

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("cache.result", "marshal_error"))
		return fmt.Errorf("error converting data to json: %w", err)
	}

	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("cache.result", "error"))
		return err
	}

	span.SetAttributes(attribute.String("cache.result", "success"))
	return nil
}
