package flights

import (
	"context"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
)

type repository interface {
	CreateFlight(ctx context.Context, f *models.Flight) error
	GetFlightByID(ctx context.Context, id uuid.UUID) (*models.Flight, error)
}

type Service struct {
	Repo repository
}

// NewFlightsService returns a new *Service that uses the provided repository for flight persistence.
// NewFlightsService returns a new Service with its Repo field set to the provided repository.
func NewFlightsService(repo repository) *Service {
	return &Service{Repo: repo}
}
