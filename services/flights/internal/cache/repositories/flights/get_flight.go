package flights

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (r *flightCache) GetFlight(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	key := fmt.Sprintf("flight:%s", id.String())

	val, err := r.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error getting data from the cache: %w", err)
	}

	var flight models.Flight

	if err := json.Unmarshal([]byte(val), &flight); err != nil {
		return nil, fmt.Errorf("error converting cache data: %w", err)
	}

	return &flight, nil
}
