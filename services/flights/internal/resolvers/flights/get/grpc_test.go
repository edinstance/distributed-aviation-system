package get

import (
	"context"
	"errors"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFlightGrpcResolverGetFlight(t *testing.T) {
	id := uuid.New()
	number := "AA123"
	origin := "LAX"
	destination := "JFK"
	departureTime := time.Date(2024, 12, 15, 10, 0, 0, 0, time.UTC)
	arrivalTime := time.Date(2024, 12, 15, 15, 0, 0, 0, time.UTC)
	expectedFlight := &models.Flight{
		ID:            id,
		Number:        number,
		Origin:        origin,
		Destination:   destination,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		Status:        models.FlightStatusScheduled,
	}

	tests := []struct {
		name          string
		id            string
		serviceSetup  func(*MockFlightService)
		expectErr     bool
		expectedError string
		expectNil     bool
	}{
		{
			name: "success",
			id:   id.String(),
			serviceSetup: func(m *MockFlightService) {
				m.On("GetFlightByID",
					mock.Anything, id,
				).Return(expectedFlight, nil)
			},
			expectErr: false,
		},
		{
			name: "flight not found",
			id:   id.String(),
			serviceSetup: func(m *MockFlightService) {
				m.On("GetFlightByID",
					mock.Anything, id,
				).Return(nil, exceptions.ErrNotFound)
			},
			expectErr: false,
			expectNil: true,
		},
		{
			name:         "invalid id",
			id:           "fake uuid",
			serviceSetup: func(m *MockFlightService) {},
			expectErr:    true,
		},
		{
			name:          "service not configured",
			id:            id.String(),
			serviceSetup:  func(_ *MockFlightService) {},
			expectErr:     true,
			expectedError: "service not configured",
		},
		{
			name: "service returns error",
			id:   id.String(),
			serviceSetup: func(m *MockFlightService) {
				m.On("GetFlightByID", mock.Anything, id).
					Return(nil, errors.New("db error"))
			},
			expectErr:     true,
			expectedError: "db error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var resolver *FlightResolver

			if tc.name == "service not configured" {
				resolver = &FlightResolver{service: nil}
			} else {
				mockService := &MockFlightService{}
				tc.serviceSetup(mockService)
				resolver = &FlightResolver{service: mockService}
			}

			req := connect.NewRequest(&v1.GetFlightByIdRequest{Id: tc.id})
			resp, err := resolver.GetFlightByIdGRPC(context.Background(), req)

			if tc.expectErr {
				assert.Error(t, err)
				if tc.expectedError != "" {
					assert.Contains(t, err.Error(), tc.expectedError)
				}
				assert.Nil(t, resp)
				return
			}

			assert.NoError(t, err)
			if tc.expectNil {
				assert.Nil(t, resp.Msg.Flight)
			} else {
				got := resp.Msg.Flight
				assert.Equal(t, expectedFlight.ID.String(), got.Id)
				assert.Equal(t, expectedFlight.Number, got.Number)
				assert.Equal(t, expectedFlight.Origin, got.Origin)
				assert.Equal(t, expectedFlight.Destination, got.Destination)
				assert.True(t, got.DepartureTime.AsTime().Equal(expectedFlight.DepartureTime))
				assert.True(t, got.ArrivalTime.AsTime().Equal(expectedFlight.ArrivalTime))
				assert.Equal(t, timestamppb.New(expectedFlight.CreatedAt).AsTime().IsZero(), true) // created/updated might be zero in mock
			}
		})
	}
}
