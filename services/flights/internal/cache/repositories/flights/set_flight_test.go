package flights

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFlightCache_SetFlight(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockRedisClient)
	cache := &flightCache{client: mockClient, ttl: time.Hour}

	flight := &models.Flight{
		ID:          uuid.New(),
		Number:      "AA123",
		Origin:      "LAX",
		Destination: "JFK",
		Status:      models.FlightStatusScheduled,
	}

	key := "flight:" + flight.ID.String()

	serialized, _ := json.Marshal(flight)

	tests := []struct {
		name      string
		flight    *models.Flight
		setupMock func()
		expectErr bool
	}{
		{
			name:   "success",
			flight: flight,
			setupMock: func() {
				mockClient.On("Set", mock.Anything, key, serialized, time.Hour).
					Return(nil).Once()
			},
			expectErr: false,
		},
		{
			name:   "redis error",
			flight: flight,
			setupMock: func() {
				mockClient.On("Set", mock.Anything, key, serialized, time.Hour).
					Return(errors.New("redis down")).Once()
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := cache.SetFlight(ctx, tc.flight)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}
