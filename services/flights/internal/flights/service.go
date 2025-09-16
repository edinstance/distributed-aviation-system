package flights

import (
	"context"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
)

type repository interface {
	CreateFlight(ctx context.Context, f *models.Flight) error
}

type Service struct {
	Repo repository
}

func NewFlightsService(repo repository) *Service {
	return &Service{Repo: repo}
}
