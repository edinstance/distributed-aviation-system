package create

import (
	"context"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
)

type FlightCreator interface {
	CreateFlight(ctx context.Context, number, origin, dest string, dep, arr time.Time) (*models.Flight, error)
}

type FlightResolver struct {
	service FlightCreator
}

// NewCreateFlightResolver creates a new FlightResolver configured with the given FlightCreator.
// NewCreateFlightResolver returns a new FlightResolver that delegates flight creation to the provided FlightCreator.
func NewCreateFlightResolver(service FlightCreator) *FlightResolver {
	return &FlightResolver{service: service}
}
