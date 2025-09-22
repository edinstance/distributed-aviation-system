package flights

import (
	"context"
	"fmt"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/validation"
	"github.com/google/uuid"
)

func (s *Service) CreateFlight(
	ctx context.Context,
	number string,
	origin string,
	destination string,
	departure time.Time,
	arrival time.Time,
) (*models.Flight, error) {

	if !arrival.After(departure) {
		return nil, fmt.Errorf("%w: departure=%v, arrival=%v",
			exceptions.ErrInvalidTimes, departure, arrival)
	}

	normalizedNumber, err := validation.ValidateAndNormalizeFlightNumber(number)
	if err != nil {
		return nil, err
	}

	normalizedOrigin, err := validation.ValidateAndNormalizeIATACode(origin)
	if err != nil {
		return nil, err
	}

	normalizedDestination, err := validation.ValidateAndNormalizeIATACode(destination)
	if err != nil {
		return nil, err
	}

	if normalizedOrigin == normalizedDestination {
		return nil, exceptions.ErrSameOriginAndDestination
	}

	flight := &models.Flight{
		ID:            uuid.New(),
		Number:        normalizedNumber,
		Origin:        normalizedOrigin,
		Destination:   normalizedDestination,
		DepartureTime: departure,
		ArrivalTime:   arrival,
		Status:        models.FlightStatusScheduled,
	}

	if err := s.Repo.CreateFlight(ctx, flight); err != nil {
		logger.Error("Failed to create flight in database", "flight_id", flight.ID, "err", err)
		return nil, err
	}

	if err := s.Cache.SetFlight(ctx, flight); err != nil {
		logger.Warn("Failed to cache flight", "flight_id", flight.ID, "err", err)
	}

	return flight, nil
}
