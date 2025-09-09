package flights

import (
	"context"
	"fmt"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
)

func (flightRepository *FlightRepository) CreateFlight(ctx context.Context, f *models.Flight) error {
	const query = `
        INSERT INTO flights (
            id, number, origin, destination,
            departure_time, arrival_time, status
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING created_at, updated_at
    `

	err := flightRepository.pool.QueryRow(
		ctx,
		query,
		f.ID,
		f.Number,
		f.Origin,
		f.Destination,
		f.DepartureTime,
		f.ArrivalTime,
		f.Status,
	).Scan(&f.CreatedAt, &f.UpdatedAt)

	if err != nil {
		logger.Error("Error saving flight to db", "id", f.ID, "error", err)
		return fmt.Errorf("create flight %s: %w", f.ID, err)
	}
	return nil
}
