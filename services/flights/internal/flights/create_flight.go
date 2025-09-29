package flights

import (
	"context"
	"fmt"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/validation/flight_number"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/validation/iata_codes"
	"github.com/google/uuid"
)

func (service *Service) CreateFlight(
	ctx context.Context,
	number string,
	origin string,
	destination string,
	departure time.Time,
	arrival time.Time,
	aircraftId uuid.UUID,
) (*models.Flight, error) {

	if !arrival.After(departure) {
		return nil, fmt.Errorf("%w: departure=%v, arrival=%v",
			exceptions.ErrInvalidTimes, departure, arrival)
	}

	normalizedNumber, err := flight_number.ValidateAndNormalizeFlightNumber(number)
	if err != nil {
		return nil, err
	}

	normalizedOrigin, err := iata_codes.ValidateAndNormalizeIATACode(origin)
	if err != nil {
		return nil, err
	}

	normalizedDestination, err := iata_codes.ValidateAndNormalizeIATACode(destination)
	if err != nil {
		return nil, err
	}

	if normalizedOrigin == normalizedDestination {
		return nil, exceptions.ErrSameOriginAndDestination
	}

	validationErr := service.AircraftClient.ValidateAircraftExists(ctx, aircraftId)
	if validationErr != nil {
		logger.Error("Aircraft does not exist", "aircraft_id", aircraftId, "err", validationErr)
		return nil, validationErr
	}

	flight := &models.Flight{
		ID:            uuid.New(),
		Number:        normalizedNumber,
		Origin:        normalizedOrigin,
		Destination:   normalizedDestination,
		DepartureTime: departure,
		ArrivalTime:   arrival,
		Status:        models.FlightStatusScheduled,
		AircraftID:    aircraftId,
	}

	if err := service.Repo.CreateFlight(ctx, flight); err != nil {
		logger.Error("Failed to create flight in database", "flight_id", flight.ID, "err", err)
		return nil, err
	}

	if err := service.Cache.SetFlight(ctx, flight); err != nil {
		logger.Warn("Failed to cache flight", "flight_id", flight.ID, "err", err)
	}

	return flight, nil
}
