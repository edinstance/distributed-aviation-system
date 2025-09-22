package flights

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type FlightRepository struct {
	pool DB
}

// NewFlightRepository returns a new FlightRepository backed by the provided *pgxpool.Pool.
// The pool is stored via the package DB interface to enable dependency injection and testing.
func NewFlightRepository(pool *pgxpool.Pool) *FlightRepository {
	return &FlightRepository{pool: pool}
}
