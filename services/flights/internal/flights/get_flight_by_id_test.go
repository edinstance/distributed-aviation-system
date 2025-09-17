package flights

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_GetFlightByID_WithFakeRepo(t *testing.T) {
	validID := uuid.New()

	expectedFlight := &models.Flight{
		ID:            validID,
		Number:        "AA123",
		Origin:        "LAX",
		Destination:   "JFK",
		DepartureTime: time.Now().Add(1 * time.Hour),
		ArrivalTime:   time.Now().Add(6 * time.Hour),
		Status:        models.FlightStatusScheduled,
	}

	tests := []struct {
		name           string
		fakeRepo       *FakeRepo
		inputID        uuid.UUID
		expectErr      bool
		expectedErrMsg string
		expectedFlight *models.Flight
	}{
		{
			name: "success",
			fakeRepo: &FakeRepo{
				GetFlightFn: func(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
					return expectedFlight, nil
				},
			},
			inputID:        validID,
			expectErr:      false,
			expectedFlight: expectedFlight,
		},
		{
			name: "not found error",
			fakeRepo: &FakeRepo{
				GetFlightFn: func(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
					return nil, errors.New("not found")
				},
			},
			inputID:        validID,
			expectErr:      true,
			expectedErrMsg: "not found",
		},
		{
			name: "nil flight but no error",
			fakeRepo: &FakeRepo{
				GetFlightFn: func(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
					return nil, nil
				},
			},
			inputID:        validID,
			expectErr:      false, // service will happily forward nil,nil
			expectedFlight: nil,   // might want to guard against this in Service
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{Repo: tc.fakeRepo}

			flight, err := service.GetFlightByID(context.Background(), tc.inputID)

			if tc.expectErr {
				assert.Error(t, err)
				if tc.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tc.expectedErrMsg)
				}
				assert.Nil(t, flight)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedFlight, flight)
			}
		})
	}
}
