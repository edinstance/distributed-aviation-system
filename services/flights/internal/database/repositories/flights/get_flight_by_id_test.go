package flights

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlightRepositoryGetFlightByID(t *testing.T) {
	cases := []struct {
		name         string
		mockErr      error
		returnRows   bool
		expectErr    bool
		assertChecks func(t *testing.T, flight *models.Flight, err error, createdAt, updatedAt time.Time, expectedID uuid.UUID)
	}{
		{
			name:       "Success",
			mockErr:    nil,
			returnRows: true,
			expectErr:  false,
			assertChecks: func(t *testing.T, flight *models.Flight, err error, createdAt, updatedAt time.Time, expectedID uuid.UUID) {
				require.NoError(t, err)
				assert.Equal(t, expectedID, flight.ID)
				assert.Equal(t, "AA123", flight.Number)
				assert.Equal(t, createdAt, flight.CreatedAt)
				assert.Equal(t, updatedAt, flight.UpdatedAt)
			},
		},
		{
			name:       "Not Found",
			mockErr:    pgx.ErrNoRows,
			returnRows: false,
			expectErr:  true,
			assertChecks: func(t *testing.T, flight *models.Flight, err error, createdAt, updatedAt time.Time, expectedID uuid.UUID) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
				assert.Nil(t, flight)
			},
		},
		{
			name:       "Database Error",
			mockErr:    errors.New("connection reset"),
			returnRows: false,
			expectErr:  true,
			assertChecks: func(t *testing.T, flight *models.Flight, err error, createdAt, updatedAt time.Time, expectedID uuid.UUID) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "get flight")
				assert.Contains(t, err.Error(), expectedID.String())
				assert.Nil(t, flight)
			},
		},
		{
			name:       "Context Cancelled",
			mockErr:    context.Canceled,
			returnRows: false,
			expectErr:  true,
			assertChecks: func(t *testing.T, flight *models.Flight, err error, createdAt, updatedAt time.Time, expectedID uuid.UUID) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "get flight")
				assert.Nil(t, flight)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			flightID := uuid.New()
			expectedSQL := `
				SELECT id, number, origin, destination, departure_time, arrival_time, status, created_at, updated_at
				FROM flights
				WHERE id = $1
			`
			createdAt := time.Date(2024, 12, 15, 9, 0, 0, 0, time.UTC)
			updatedAt := createdAt

			expect := mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
				WithArgs(flightID)

			if tc.returnRows {
				expect.WillReturnRows(
					pgxmock.NewRows([]string{
						"id", "number", "origin", "destination", "departure_time", "arrival_time", "status", "created_at", "updated_at",
					}).AddRow(
						flightID,
						"AA123",
						"LAX",
						"JFK",
						time.Date(2024, 12, 15, 10, 0, 0, 0, time.UTC),
						time.Date(2024, 12, 15, 15, 0, 0, 0, time.UTC),
						models.FlightStatusScheduled,
						createdAt,
						updatedAt,
					),
				)
			} else {
				expect.WillReturnError(tc.mockErr)
			}

			repo := &FlightRepository{pool: mock}
			ctx := context.Background()
			if tc.name == "ContextCancelled" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			flight, err := repo.GetFlightByID(ctx, flightID)
			tc.assertChecks(t, flight, err, createdAt, updatedAt, flightID)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
