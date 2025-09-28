package aircraft_client

import (
	"context"
	"fmt"
	"time"

	aircraftv1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/aircraft/v1"
	"github.com/google/uuid"
)

func (c *AircraftClient) ValidateAircraftExists(ctx context.Context, aircraftID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &aircraftv1.GetAircraftByIdRequest{
		Id: aircraftID.String(),
	}

	resp, err := c.client.GetAircraftById(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to validate aircraft existence: %w", err)
	}

	if resp.Aircraft == nil {
		return fmt.Errorf("aircraft with id %s not found", aircraftID.String())
	}

	return nil
}
