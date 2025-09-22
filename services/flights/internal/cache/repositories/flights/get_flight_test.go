package flights

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFlightCacheGetFlight(t *testing.T) {
	ctx := context.Background()
	mockClient := new(MockRedisClient)
	cache := &flightCache{client: mockClient, ttl: time.Hour}

	flight := &models.Flight{
		ID:            uuid.New(),
		Number:        "AA123",
		Origin:        "LAX",
		Destination:   "JFK",
		DepartureTime: time.Now().UTC(),
		ArrivalTime:   time.Now().Add(5 * time.Hour).UTC(),
		Status:        models.FlightStatusScheduled,
	}
	flightJSON, _ := json.Marshal(flight)
	key := "flight:" + flight.ID.String()

	tests := []struct {
		name        string
		setupMock   func()
		expectErr   bool
		expectedNil bool
		expectedID  uuid.UUID
	}{
		{
			name: "cache hit",
			setupMock: func() {
				mockClient.On("Get", mock.Anything, key).
					Return(string(flightJSON), nil).Once()
			},
			expectedNil: false,
			expectedID:  flight.ID,
		},
		{
			name: "cache miss (redis.Nil)",
			setupMock: func() {
				mockClient.On("Get", mock.Anything, key).
					Return("", redis.Nil).Once()
			},
			expectedNil: true,
		},
		{
			name: "redis error",
			setupMock: func() {
				mockClient.On("Get", mock.Anything, key).
					Return("", errors.New("connection failed")).Once()
			},
			expectErr: true,
		},
		{
			name: "invalid json",
			setupMock: func() {
				mockClient.On("Get", mock.Anything, key).
					Return("not-json", nil).Once()
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			got, err := cache.GetFlight(ctx, flight.ID)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else if tc.expectedNil {
				assert.NoError(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tc.expectedID, got.ID)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
