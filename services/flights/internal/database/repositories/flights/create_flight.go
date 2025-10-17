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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (flightRepository *FlightRepository) CreateFlight(ctx context.Context, f *models.Flight) error {
	tracer := otel.Tracer("flights-service")
	ctx, span := tracer.Start(ctx, "db.create_flight")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.operation", "insert"),
		attribute.String("db.table", "flights"),
		attribute.String("flight.id", f.ID.String()),
		attribute.String("flight.number", f.Number),
		attribute.String("flight.origin", f.Origin),
		attribute.String("flight.destination", f.Destination),
	)

	const query = `
        INSERT INTO flights (
            id, number, origin, destination,
            departure_time, arrival_time, status, aircraft_id,
            created_by, last_updated_by, organization_id
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
		f.AircraftID,
		f.CreatedBy,
		f.LastUpdatedBy,
		f.OrganizationID,
	).Scan(&f.CreatedAt, &f.UpdatedAt)

	if err != nil {
		span.RecordError(err)

		// Check if it's a Postgres error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			span.SetAttributes(
				attribute.String("db.error.code", pgErr.Code),
				attribute.String("db.error.constraint", pgErr.ConstraintName),
			)

			if pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == "unique_flight_instance" {
				span.SetAttributes(attribute.String("db.result", "duplicate"))
				return fmt.Errorf(
					"flight with number %s at %s already exists",
					f.Number,
					f.DepartureTime.Format(time.RFC3339),
				)
			} else {
				logger.Error("Error saving flight to db", "id", f.ID, "code", pgErr.Code, "constraint", pgErr.ConstraintName, "error", err)
				span.SetAttributes(attribute.String("db.result", "postgres_error"))
				return fmt.Errorf("postgres error [%s]: %w", pgErr.Code, err)
			}
		}

		logger.Error("Error saving flight to db", "id", f.ID, "error", err)
		span.SetAttributes(attribute.String("db.result", "error"))
		return fmt.Errorf("create flight %s: %w", f.ID, err)
	}

	span.SetAttributes(attribute.String("db.result", "success"))
	return nil
}
