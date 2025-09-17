package get

import (
	"context"
	"errors"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/google/uuid"
)

func (r *FlightResolver) GetFlight(
	ctx context.Context,
	id string,
) (*models.Flight, error) {
	logger.Debug("GetFlight GraphQL request", "id", id)

	if r.service == nil {
		logger.Error("GetFlight service not configured")
		return nil, errors.New("service not configured")
	}

	flightID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid flight ID format", "id", id, "err", err)
		return nil, errors.New("invalid flight ID format")
	}

	flight, err := r.service.GetFlightByID(ctx, flightID)
	if err != nil {
		logger.Error("Failed to get flight", "id", id, "err", err)
		return nil, err
	}

	logger.Debug("GetFlight GraphQL response retrieved", "id", flight.ID)
	return flight, nil
}