package flights

import (
	"context"
	"errors"
	"fmt"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (flightRepository *FlightRepository) GetFlightByID(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	tracer := otel.Tracer("flights-service")
	ctx, span := tracer.Start(ctx, "db.get_flight_by_id")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.operation", "select"),
		attribute.String("db.table", "flights"),
		attribute.String("flight.id", id.String()),
	)

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
		span.RecordError(err)
		if errors.Is(err, pgx.ErrNoRows) {
			span.SetAttributes(attribute.String("db.result", "not_found"))
			return nil, exceptions.ErrNotFound
		}
		span.SetAttributes(attribute.String("db.result", "error"))
		return nil, fmt.Errorf("get flight %s: %w", id, err)
	}

	span.SetAttributes(
		attribute.String("db.result", "success"),
		attribute.String("flight.number", flight.Number),
	)

	return &flight, nil
}
