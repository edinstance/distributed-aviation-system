package flights_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/flights"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type FakeRepo struct {
	CreateFlightFn func(ctx context.Context, f *models.Flight) error
}

func (f *FakeRepo) CreateFlight(ctx context.Context, fl *models.Flight) error {
	return f.CreateFlightFn(ctx, fl)
}

func newService(repo *FakeRepo) *flights.Service {
	return flights.NewFlightsService(repo)
}

func TestCreateFlight(t *testing.T) {
	now := time.Now()
	dep := now.Add(1 * time.Hour)
	arr := dep.Add(2 * time.Hour)

	repoErr := errors.New("db failure")

	tests := []struct {
		name        string
		number      string
		origin      string
		dest        string
		departure   time.Time
		arrival     time.Time
		repo        *FakeRepo
		expectError error
	}{
		{
			name:      "valid flight",
			number:    "AA123",
			origin:    "JFK",
			dest:      "LHR",
			departure: dep,
			arrival:   arr,
			repo: &FakeRepo{
				CreateFlightFn: func(ctx context.Context, f *models.Flight) error {
					return nil
				},
			},
			expectError: nil,
		},
		{
			name:        "arrival before departure",
			number:      "AA123",
			origin:      "JFK",
			dest:        "LHR",
			departure:   arr,
			arrival:     dep,
			repo:        &FakeRepo{},
			expectError: exceptions.ErrInvalidTimes,
		},
		{
			name:        "same origin and destination",
			number:      "AA123",
			origin:      "JFK",
			dest:        "JFK",
			departure:   dep,
			arrival:     arr,
			repo:        &FakeRepo{},
			expectError: exceptions.ErrSameOriginAndDestination,
		},
		{
			name:      "repo error",
			number:    "AA123",
			origin:    "JFK",
			dest:      "LHR",
			departure: dep,
			arrival:   arr,
			repo: &FakeRepo{
				CreateFlightFn: func(ctx context.Context, f *models.Flight) error {
					return repoErr
				},
			},
			expectError: repoErr,
		},
		{
			name:      "invalid flight number",
			number:    "123",
			origin:    "JFK",
			dest:      "LHR",
			departure: dep,
			arrival:   arr,
			repo: &FakeRepo{
				CreateFlightFn: func(ctx context.Context, f *models.Flight) error {
					return repoErr
				},
			},
			expectError: exceptions.ErrInvalidFlightNumber,
		},
		{
			name:      "invalid origin",
			number:    "BA121",
			origin:    "JFK132",
			dest:      "LHR",
			departure: dep,
			arrival:   arr,
			repo: &FakeRepo{
				CreateFlightFn: func(ctx context.Context, f *models.Flight) error {
					return repoErr
				},
			},
			expectError: exceptions.ErrInvalidIATACode,
		},
		{
			name:      "invalid destination",
			number:    "BA121",
			origin:    "LHR",
			dest:      "LHR1232",
			departure: dep,
			arrival:   arr,
			repo: &FakeRepo{
				CreateFlightFn: func(ctx context.Context, f *models.Flight) error {
					return repoErr
				},
			},
			expectError: exceptions.ErrInvalidIATACode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newService(tt.repo)

			flight, err := svc.CreateFlight(
				context.Background(),
				tt.number,
				tt.origin,
				tt.dest,
				tt.departure,
				tt.arrival,
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
