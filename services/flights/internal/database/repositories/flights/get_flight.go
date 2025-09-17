package flights

import (
	"context"
	"errors"
	"fmt"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (flightRepository *FlightRepository) GetFlightByID(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	const query = `
        SELECT id, number, origin, destination, departure_time, arrival_time, status, created_at, updated_at
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
		&flight.CreatedAt,
		&flight.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("flight with id %s not found", id)
		}
		return nil, fmt.Errorf("get flight %s: %w", id, err)
	}

	return &flight, nil
}
