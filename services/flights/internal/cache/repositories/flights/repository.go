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

type redisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type flightCache struct {
	client redisClient
	ttl    time.Duration
}

func NewRedisFlightRepository(client *redis.Client, ttl time.Duration) FlightCacheRepository {
	return &flightCache{client: client, ttl: ttl}
}
