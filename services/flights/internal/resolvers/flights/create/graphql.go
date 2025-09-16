package create

import (
	"context"
	"errors"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
)

type GraphQLCreateFlightResolver struct {
	service FlightCreator
}

func NewGraphQLCreateFlightResolver(service FlightCreator) *GraphQLCreateFlightResolver {
	return &GraphQLCreateFlightResolver{service: service}
}

func (r *GraphQLCreateFlightResolver) CreateFlight(
	ctx context.Context,
	number string,
	origin string,
	destination string,
	departureTime time.Time,
	arrivalTime time.Time,
) (*models.Flight, error) {
	logger.Debug("CreateFlight GraphQL request", "number", number)

	if r.service == nil {
		logger.Error("CreateFlight service not configured")
		return nil, errors.New("service not configured")
	}

	flight, err := r.service.CreateFlight(
		ctx,
		number,
		origin,
		destination,
		departureTime,
		arrivalTime,
	)
	if err != nil {
		logger.Error("Failed to create flight", "err", err)
		return nil, err
	}

	logger.Debug("CreateFlight GraphQL response created", "number", flight.Number, "id", flight.ID)
	return flight, nil
}
