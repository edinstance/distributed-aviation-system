package flights

import (
	"context"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
)

func (s *Service) GetFlightByID(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	return s.Repo.GetFlightByID(ctx, id)
}
