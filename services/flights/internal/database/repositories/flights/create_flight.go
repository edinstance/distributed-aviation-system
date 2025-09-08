package flights

import (
	"context"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
)

func (flightRepository *FlightRepository) CreateFlight(ctx context.Context, f *models.Flight) error {
	sql := `
        INSERT INTO flights (
            id, number, origin, destination,
            departure_time, arrival_time, status
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING created_at, updated_at
    `

	err := flightRepository.Pool.QueryRow(
		ctx,
		sql,
		f.ID,
		f.Number,
		f.Origin,
		f.Destination,
		f.DepartureTime,
		f.ArrivalTime,
		f.Status,
	).Scan(&f.CreatedAt, &f.UpdatedAt)

	return err
}
