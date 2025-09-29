package create

import (
	"context"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
)

type FlightCreator interface {
	CreateFlight(ctx context.Context, number, origin, dest string, dep, arr time.Time, aircraftId uuid.UUID) (*models.Flight, error)
}

type FlightResolver struct {
	service FlightCreator
}

// NewCreateFlightResolver creates a new FlightResolver configured with the given FlightCreator.
// NewCreateFlightResolver returns a new FlightResolver configured with the given FlightCreator.
// The returned resolver delegates flight creation to the provided service.
func NewCreateFlightResolver(service FlightCreator) *FlightResolver {
	return &FlightResolver{service: service}
}
