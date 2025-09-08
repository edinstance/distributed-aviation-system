package flights

import (
	"context"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	flightrepo "github.com/edinstance/distributed-aviation-system/services/flights/internal/database/repositories/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/google/uuid"
)

// Service holds business logic for Flights
type Service struct {
	repo *flightrepo.FlightRepository
}

func NewFlightsService(repo *flightrepo.FlightRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateFlight(
	ctx context.Context,
	number string,
	origin string,
	destination string,
	departure time.Time,
	arrival time.Time,
) (*models.Flight, error) {

	if !arrival.After(departure) {
		return nil, exceptions.ErrInvalidTimes
	}

	flight := &models.Flight{
		ID:            uuid.New(),
		Number:        number,
		Origin:        origin,
		Destination:   destination,
		DepartureTime: departure,
		ArrivalTime:   arrival,
		Status:        "SCHEDULED",
	}

	if err := s.repo.CreateFlight(ctx, flight); err != nil {
		return nil, err
	}

	return flight, nil
}
