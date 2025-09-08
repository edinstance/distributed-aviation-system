package flights

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type FlightRepository struct {
	Pool *pgxpool.Pool
}

func NewFlightRepository(pool *pgxpool.Pool) *FlightRepository {
	return &FlightRepository{Pool: pool}
}
