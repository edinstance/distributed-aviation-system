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
// The returned Service's Repo field is set to the supplied repository.
func NewFlightsService(repo repository) *Service {
	return &Service{Repo: repo}
}
