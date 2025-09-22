package flights

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
)

func (r *flightCache) SetFlight(ctx context.Context, flight *models.Flight) error {
	key := fmt.Sprintf("flight:%s", flight.ID.String())

	data, err := json.Marshal(flight)

	if err != nil {
		return fmt.Errorf("error setting cache data: %w", err)
	}

	return r.client.Set(ctx, key, data, r.ttl).Err()
}
