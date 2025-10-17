package flights

import (
	"context"
	"fmt"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/middleware"
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
		logger.ErrorContext(ctx, "Aircraft does not exist", "aircraft_id", aircraftId, "err", validationErr)
		return nil, validationErr
	}

	userContext := middleware.GetRequestUserContext(ctx)

	flight := &models.Flight{
		ID:             uuid.New(),
		Number:         normalizedNumber,
		Origin:         normalizedOrigin,
		Destination:    normalizedDestination,
		DepartureTime:  departure,
		ArrivalTime:    arrival,
		Status:         models.FlightStatusScheduled,
		AircraftID:     aircraftId,
		CreatedBy:      userContext.UserID,
		LastUpdatedBy:  userContext.UserID,
		OrganizationID: userContext.OrgID,
	}

	if err := service.Repo.CreateFlight(ctx, flight); err != nil {
		logger.ErrorContext(ctx, "Failed to create flight in database", "flight_id", flight.ID, "err", err)
		return nil, err
	}

	if err := service.Cache.SetFlight(ctx, flight); err != nil {
		logger.WarnContext(ctx, "Failed to cache flight", "flight_id", flight.ID, "err", err)
	}

	if err := service.KafkaPublisher.PublishFlightCreated(ctx, flight); err != nil {
		logger.WarnContext(ctx, "Failed to publish flight created event", "flight_id", flight.ID, "err", err)
	}

	logger.InfoContext(ctx, "Flight created", "flight_id", flight.ID, "number", flight.Number, "origin", flight.Origin, "destination", flight.Destination, "departure_time", flight.DepartureTime, "arrival_time", flight.ArrivalTime, "aircraft_id", flight.AircraftID)

	return flight, nil
}
