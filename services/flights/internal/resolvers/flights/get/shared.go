package get

import (
	"context"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
)

type FlightGetter interface {
	GetFlightByID(ctx context.Context, id uuid.UUID) (*models.Flight, error)
}

type FlightResolver struct {
	service FlightGetter
}

// NewGetFlightResolver creates a new FlightResolver configured with the given FlightGetter.
// NewGetFlightResolver returns a FlightResolver that delegates flight retrieval to the provided FlightGetter.
func NewGetFlightResolver(service FlightGetter) *FlightResolver {
	return &FlightResolver{service: service}
}
