package aircraft_client

import (
	"context"
	"fmt"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	aircraftv1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/aircraft/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *AircraftClient) ValidateAircraftExists(ctx context.Context, aircraftID uuid.UUID) error {
	tracer := otel.Tracer("flights-service")
	ctx, span := tracer.Start(ctx, "aircraft_client.validate_aircraft_exists")
	defer span.End()

	span.SetAttributes(
		attribute.String("aircraft.id", aircraftID.String()),
		attribute.String("client.service", "aircraft-service"),
		attribute.String("grpc.method", "GetAircraftById"),
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &aircraftv1.GetAircraftByIdRequest{
		Id: aircraftID.String(),
	}

	resp, err := c.client.GetAircraftById(ctx, req)

	if err != nil {
		span.RecordError(err)
		if s, ok := status.FromError(err); ok {
			span.SetAttributes(attribute.String("grpc.status_code", s.Code().String()))
			if s.Code() == codes.NotFound {
				span.SetAttributes(attribute.String("client.result", "not_found"))
				logger.InfoContext(ctx, "Aircraft not found", "aircraft_id", aircraftID)
				return exceptions.AircraftNotFound(aircraftID)
			}
			span.SetAttributes(attribute.String("client.result", "grpc_error"))
			logger.ErrorContext(ctx, "Downstream aircraft service error", "aircraft_id", aircraftID, "err", s.Message())
			return exceptions.ErrDownstreamClientDown
		}
		span.SetAttributes(attribute.String("client.result", "error"))
		return fmt.Errorf("%w: %v", exceptions.ErrDownstreamClientDown, err)
	}

	if resp.Aircraft == nil {
		span.SetAttributes(attribute.String("client.result", "nil_aircraft"))
		logger.InfoContext(ctx, "Aircraft not found (nil response)", "aircraft_id", aircraftID)
		return exceptions.AircraftNotFound(aircraftID)
	}

	span.SetAttributes(
		attribute.String("client.result", "success"),
		attribute.String("aircraft.model", resp.Aircraft.Model),
	)

	return nil
}

var _ AircraftValidator = (*AircraftClient)(nil)
