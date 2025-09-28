package aircraft_client

import (
	"context"
	"fmt"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	aircraftv1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/aircraft/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *AircraftClient) ValidateAircraftExists(ctx context.Context, aircraftID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &aircraftv1.GetAircraftByIdRequest{
		Id: aircraftID.String(),
	}

	resp, err := c.client.GetAircraftById(ctx, req)

	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.NotFound {
				logger.Info("Aircraft not found", "aircraft_id", aircraftID)
				return exceptions.AircraftNotFound(aircraftID)
			}
			logger.Error("Downstream aircraft service error", "aircraft_id", aircraftID, "err", s.Message())
			return exceptions.ErrDownstreamClientDown
		}
		return fmt.Errorf("%w: %v", exceptions.ErrDownstreamClientDown, err)
	}

	if resp.Aircraft == nil {
		logger.Info("Aircraft not found (nil response)", "aircraft_id", aircraftID)
		return exceptions.AircraftNotFound(aircraftID)
	}

	return nil
}

var _ AircraftValidator = (*AircraftClient)(nil)
