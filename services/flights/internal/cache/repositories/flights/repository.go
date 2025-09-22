package flights

import (
	"context"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type FlightCacheRepository interface {
	GetFlight(ctx context.Context, id uuid.UUID) (*models.Flight, error)
	SetFlight(ctx context.Context, flight *models.Flight) error
}

type flightCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisFlightRepository(client *redis.Client, ttl time.Duration) FlightCacheRepository {
	return &flightCache{client: client, ttl: ttl}
}
