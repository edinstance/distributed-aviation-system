package flights

import (
	"context"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/google/uuid"
)

func (service *Service) GetFlightByID(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	if service.Cache != nil {
		flight, err := service.Cache.GetFlight(ctx, id)
		if err != nil {
			logger.Warn("Cache error during flight retrieval", "flight_id", id, "err", err)
		} else if flight != nil {
			logger.Debug("Flight found in cache", "flight_id", id)
			return flight, nil
		}
	}

	flight, err := service.Repo.GetFlightByID(ctx, id)
	logger.Debug("Flight found in db", "flight_id", id)
	if err != nil {
		return nil, err
	}

	if flight != nil && service.Cache != nil {
		if cacheErr := service.Cache.SetFlight(ctx, flight); cacheErr != nil {
			logger.Warn("Failed to cache flight", "flight_id", id, "err", cacheErr)
		}
	}

	return flight, nil
}
