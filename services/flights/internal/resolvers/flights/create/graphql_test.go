package create

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFlightService struct {
	mock.Mock
}

func (m *MockFlightService) CreateFlight(
	ctx context.Context,
	number, origin, dest string,
	dep, arr time.Time,
) (*models.Flight, error) {
	args := m.Called(ctx, number, origin, dest, dep, arr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Flight), args.Error(1)
}

func TestFlightResolver_CreateFlight(t *testing.T) {
	number := "AA123"
	origin := "LAX"
	destination := "JFK"
	departureTime := time.Date(2024, 12, 15, 10, 0, 0, 0, time.UTC)
	arrivalTime := time.Date(2024, 12, 15, 15, 0, 0, 0, time.UTC)
	expectedFlight := &models.Flight{
		ID:            uuid.New(),
		Number:        number,
		Origin:        origin,
		Destination:   destination,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		Status:        models.FlightStatusScheduled,
	}

	tests := []struct {
		name           string
		serviceSetup   func(*MockFlightService)
		expectErr      bool
		expectedError  string
		expectedFlight *models.Flight
	}{
		{
			name: "success",
			serviceSetup: func(m *MockFlightService) {
				m.On("CreateFlight",
					mock.Anything, number, origin, destination, departureTime, arrivalTime,
				).Return(expectedFlight, nil)
			},
			expectErr:      false,
			expectedFlight: expectedFlight,
		},
		{
			name: "service not configured",
			serviceSetup: func(_ *MockFlightService) {
				// No service setup, will use nil service in resolver
			},
			expectErr:     true,
			expectedError: "service not configured",
		},
		{
			name: "service returns error",
			serviceSetup: func(m *MockFlightService) {
				m.On("CreateFlight",
					mock.Anything, number, origin, destination, departureTime, arrivalTime,
				).Return(nil, errors.New("db error"))
			},
			expectErr:     true,
			expectedError: "db error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var resolver *FlightResolver

			if tc.name == "service not configured" {
				// Resolver with nil service
				resolver = &FlightResolver{}
			} else {
				mockService := &MockFlightService{}
				tc.serviceSetup(mockService)
				resolver = &FlightResolver{service: mockService}
			}

			flight, err := resolver.CreateFlight(
				context.Background(),
				number,
				origin,
				destination,
				departureTime,
				arrivalTime,
			)

			if tc.expectErr {
				assert.Error(t, err)
				if tc.expectedError != "" {
					assert.Contains(t, err.Error(), tc.expectedError)
				}
				assert.Nil(t, flight)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedFlight, flight)
			}
		})
	}
}
