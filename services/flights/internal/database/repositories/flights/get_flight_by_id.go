package flights

import (
	"context"
	"errors"
	"fmt"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (flightRepository *FlightRepository) GetFlightByID(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	const query = `
        SELECT id, number, origin, destination, departure_time, arrival_time, status, aircraft_id, created_at, updated_at
        FROM flights
        WHERE id = $1
    `

	var flight models.Flight
	err := flightRepository.pool.QueryRow(ctx, query, id).Scan(
		&flight.ID,
		&flight.Number,
		&flight.Origin,
		&flight.Destination,
		&flight.DepartureTime,
		&flight.ArrivalTime,
		&flight.Status,
		&flight.AircraftID,
		&flight.CreatedAt,
		&flight.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exceptions.ErrNotFound
		}
		return nil, fmt.Errorf("get flight %s: %w", id, err)
	}

	return &flight, nil
}
