package flights

import (
	"context"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/cache/repositories/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/clients/aircraft_client"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
)

type repository interface {
	CreateFlight(ctx context.Context, f *models.Flight) error
	GetFlightByID(ctx context.Context, id uuid.UUID) (*models.Flight, error)
}

type Service struct {
	Repo           repository
	Cache          flights.FlightCacheRepository
	AircraftClient aircraft_client.AircraftValidator
}

// NewFlightsService returns a new *Service that uses the provided repository for flight persistence.
func NewFlightsService(repo repository, cache flights.FlightCacheRepository,
	aircraftClient aircraft_client.AircraftValidator) *Service {
	return &Service{Repo: repo, Cache: cache, AircraftClient: aircraftClient}
}
