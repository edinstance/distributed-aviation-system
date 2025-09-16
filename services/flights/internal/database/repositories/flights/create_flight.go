package flights

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
		// Check if it's a Postgres error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == "unique_flight_instance" {
				return fmt.Errorf(
					"flight with number %s at %s already exists",
					f.Number,
					f.DepartureTime.Format(time.RFC3339),
				)
			} else {
				logger.Error("Error saving flight to db", "id", f.ID, "code", pgErr.Code, "constraint", pgErr.ConstraintName, "error", err)
				return fmt.Errorf("postgres error [%s]: %w", pgErr.Code, err)
			}
		}

		logger.Error("Error saving flight to db", "id", f.ID, "error", err)
		return fmt.Errorf("create flight %s: %w", f.ID, err)
	}

	return nil
}
