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
			logger.WarnContext(ctx, "Cache error during flight retrieval", "flight_id", id, "err", err)
		} else if flight != nil {
			logger.DebugContext(ctx, "Flight found in cache", "flight_id", id)
			return flight, nil
		}
	}

	flight, err := service.Repo.GetFlightByID(ctx, id)

	if err != nil {
		return nil, err
	}

	logger.DebugContext(ctx, "Flight found in db", "flight_id", id)

	if flight != nil {
		if service.Cache != nil {
			if cacheErr := service.Cache.SetFlight(ctx, flight); cacheErr != nil {
				logger.WarnContext(ctx, "Failed to cache flight",
					"flight_id", id,
					"err", cacheErr)
			}
		}

		logger.InfoContext(ctx, "Flight retrieved",
			"flight_id", id,
			"number", flight.Number,
			"origin", flight.Origin,
			"destination", flight.Destination,
			"departure_time", flight.DepartureTime,
			"arrival_time", flight.ArrivalTime,
			"aircraft_id", flight.AircraftID,
			"airline", flight.Airline)

		return flight, nil
	}

	logger.WarnContext(ctx, "Flight not found", "flight_id", id)
	return nil, nil
}
