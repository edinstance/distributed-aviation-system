package flights

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func defaultTestDeps() (*FakeRepo, *FakeFlightsCache, *FakeAircraftClient) {
	return &FakeRepo{}, &FakeFlightsCache{}, &FakeAircraftClient{}
}

func TestCreateFlight(t *testing.T) {
	now := time.Now()
	dep := now.Add(1 * time.Hour)
	arr := dep.Add(2 * time.Hour)

	repoErr := errors.New("db failure")
	aircraftErr := errors.New("aircraft not found")

	tests := []struct {
		name        string
		number      string
		origin      string
		dest        string
		departure   time.Time
		arrival     time.Time
		setup       func(r *FakeRepo, c *FakeFlightsCache, a *FakeAircraftClient)
		expectError error
	}{
		{
			name:      "valid flight",
			number:    "AA123",
			origin:    "JFK",
			dest:      "LHR",
			departure: dep,
			arrival:   arr,
			setup:     func(r *FakeRepo, c *FakeFlightsCache, a *FakeAircraftClient) {},
		},
		{
			name:        "arrival before departure",
			number:      "AA123",
			origin:      "JFK",
			dest:        "LHR",
			departure:   arr,
			arrival:     dep,
			setup:       func(r *FakeRepo, c *FakeFlightsCache, a *FakeAircraftClient) {},
			expectError: exceptions.ErrInvalidTimes,
		},
		{
			name:        "same origin and destination",
			number:      "AA123",
			origin:      "JFK",
			dest:        "JFK",
			departure:   dep,
			arrival:     arr,
			setup:       func(r *FakeRepo, c *FakeFlightsCache, a *FakeAircraftClient) {},
			expectError: exceptions.ErrSameOriginAndDestination,
		},
		{
			name:      "repo error",
			number:    "AA123",
			origin:    "JFK",
			dest:      "LHR",
			departure: dep,
			arrival:   arr,
			setup: func(r *FakeRepo, _ *FakeFlightsCache, _ *FakeAircraftClient) {
				r.CreateFlightFn = func(ctx context.Context, f *models.Flight) error {
					return repoErr
				}
			},
			expectError: repoErr,
		},
		{
			name:        "invalid flight number",
			number:      "123",
			origin:      "JFK",
			dest:        "LHR",
			departure:   dep,
			arrival:     arr,
			setup:       func(r *FakeRepo, c *FakeFlightsCache, a *FakeAircraftClient) {},
			expectError: exceptions.ErrInvalidFlightNumber,
		},
		{
			name:        "invalid origin",
			number:      "BA121",
			origin:      "JFK132",
			dest:        "LHR",
			departure:   dep,
			arrival:     arr,
			setup:       func(r *FakeRepo, c *FakeFlightsCache, a *FakeAircraftClient) {},
			expectError: exceptions.ErrInvalidIATACode,
		},
		{
			name:        "invalid destination",
			number:      "BA121",
			origin:      "LHR",
			dest:        "LHR1232",
			departure:   dep,
			arrival:     arr,
			setup:       func(r *FakeRepo, c *FakeFlightsCache, a *FakeAircraftClient) {},
			expectError: exceptions.ErrInvalidIATACode,
		},
		{
			name:      "aircraft validation error",
			number:    "BA121",
			origin:    "LHR",
			dest:      "LGW",
			departure: dep,
			arrival:   arr,
			setup: func(_ *FakeRepo, _ *FakeFlightsCache, r *FakeAircraftClient) {
				r.ValidateAircraftExistsFn = func(ctx context.Context, id uuid.UUID) error {
					return aircraftErr
				}
			},
			expectError: aircraftErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, cache, aircraft := defaultTestDeps()
			tt.setup(repo, cache, aircraft)

			svc := NewFlightsService(repo, cache, aircraft)

			flight, err := svc.CreateFlight(
				context.Background(),
				tt.number,
				tt.origin,
				tt.dest,
				tt.departure,
				tt.arrival,
				uuid.New(),
			)

			if tt.expectError != nil {
				assert.Nil(t, flight)
				assert.ErrorIs(t, err, tt.expectError)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, flight)
			assert.Equal(t, tt.number, flight.Number)
			assert.Equal(t, tt.origin, flight.Origin)
			assert.Equal(t, tt.dest, flight.Destination)
			assert.Equal(t, models.FlightStatusScheduled, flight.Status)
			assert.NotEqual(t, uuid.Nil, flight.ID)
		})
	}
}
