package flights

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
)

func TestFlightRepositoryCreateFlight(testHelper *testing.T) {
	cases := []struct {
		name         string
		mockErr      error
		returnRows   bool
		expectErr    bool
		assertChecks func(testHelper *testing.T, flight *models.Flight, err error, createdAt, updatedAt time.Time)
	}{
		{
			name:       "Success",
			mockErr:    nil,
			returnRows: true,
			expectErr:  false,
			assertChecks: func(testHelper *testing.T, flight *models.Flight, err error, createdAt, updatedAt time.Time) {
				require.NoError(testHelper, err)
				assert.Equal(testHelper, createdAt, flight.CreatedAt)
				assert.Equal(testHelper, updatedAt, flight.UpdatedAt)
			},
		},
		{
			name:       "DatabaseError",
			mockErr:    pgx.ErrNoRows,
			returnRows: false,
			expectErr:  true,
			assertChecks: func(testHelper *testing.T, flight *models.Flight, err error, createdAt, updatedAt time.Time) {
				require.Error(testHelper, err)
				assert.Contains(testHelper, err.Error(), "create flight")
				assert.Contains(testHelper, err.Error(), flight.ID.String())
				assert.Zero(testHelper, flight.CreatedAt)
				assert.Zero(testHelper, flight.UpdatedAt)
			},
		},
		{
			name: "UniqueConstraintViolation",
			mockErr: &pgconn.PgError{
				Code:           pgerrcode.UniqueViolation,
				Message:        "duplicate key value violates unique constraint",
				ConstraintName: "unique_flight_instance",
			},
			returnRows: false,
			expectErr:  true,
			assertChecks: func(testHelper *testing.T, flight *models.Flight, err error, createdAt, updatedAt time.Time) {
				require.Error(testHelper, err)
				assert.Contains(testHelper, err.Error(), "flight with number")
				assert.Contains(testHelper, err.Error(), "already exists")
			},
		},
		{
			name:       "ContextCancelled",
			mockErr:    context.Canceled,
			returnRows: false,
			expectErr:  true,
			assertChecks: func(testHelper *testing.T, flight *models.Flight, err error, createdAt, updatedAt time.Time) {
				require.Error(testHelper, err)
				assert.Contains(testHelper, err.Error(), "create flight")
			},
		},
	}

	for _, tc := range cases {
		testHelper.Run(tc.name, func(testHelper *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(testHelper, err)
			defer mock.Close()

			flight := &models.Flight{
				ID:            uuid.New(),
				Number:        "AA123",
				Origin:        "LAX",
				Destination:   "JFK",
				DepartureTime: time.Date(2024, 12, 15, 10, 0, 0, 0, time.UTC),
				ArrivalTime:   time.Date(2024, 12, 15, 15, 0, 0, 0, time.UTC),
				Status:        models.FlightStatusScheduled,
			}

			expectedSQL := `INSERT INTO flights ( id, number, origin, destination, departure_time, arrival_time, status ) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING created_at, updated_at`
			createdAt := time.Date(2024, 12, 15, 9, 0, 0, 0, time.UTC)
			updatedAt := createdAt

			expectedSQL = regexp.QuoteMeta(expectedSQL)
			expect := mock.ExpectQuery(expectedSQL).WithArgs(
				flight.ID,
				flight.Number,
				flight.Origin,
				flight.Destination,
				flight.DepartureTime,
				flight.ArrivalTime,
				flight.Status,
			)

			if tc.returnRows {
				expect.WillReturnRows(
					pgxmock.NewRows([]string{"created_at", "updated_at"}).
						AddRow(createdAt, updatedAt),
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

			err = repo.CreateFlight(ctx, flight)
			tc.assertChecks(testHelper, flight, err, createdAt, updatedAt)

			assert.NoError(testHelper, mock.ExpectationsWereMet())
		})
	}
}

func TestFlightRepositoryCreateFlightAllFlightStatuses(testHelper *testing.T) {
	statuses := []models.FlightStatus{
		models.FlightStatusScheduled,
		models.FlightStatusDelayed,
		models.FlightStatusCancelled,
		models.FlightStatusDeparted,
		models.FlightStatusInProgress,
		models.FlightStatusArrived,
		models.FlightStatusUnspecified,
	}

	for _, status := range statuses {
		testHelper.Run(string(status), func(testHelper *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(testHelper, err)
			defer mock.Close()

			flight := &models.Flight{
				ID:            uuid.New(),
				Number:        "TEST123",
				Origin:        "LAX",
				Destination:   "JFK",
				DepartureTime: time.Date(2024, 12, 15, 10, 0, 0, 0, time.UTC),
				ArrivalTime:   time.Date(2024, 12, 15, 15, 0, 0, 0, time.UTC),
				Status:        status,
			}

			expectedSQL := `INSERT INTO flights ( id, number, origin, destination, departure_time, arrival_time, status ) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING created_at, updated_at`
			createdAt := time.Date(2024, 12, 15, 9, 0, 0, 0, time.UTC)
			updatedAt := createdAt

			expectedSQL = regexp.QuoteMeta(expectedSQL)

			mock.ExpectQuery(expectedSQL).
				WithArgs(
					flight.ID,
					flight.Number,
					flight.Origin,
					flight.Destination,
					flight.DepartureTime,
					flight.ArrivalTime,
					status,
				).
				WillReturnRows(pgxmock.NewRows([]string{"created_at", "updated_at"}).
					AddRow(createdAt, updatedAt))

			repo := &FlightRepository{pool: mock}
			ctx := context.Background()

			err = repo.CreateFlight(ctx, flight)

			assert.NoError(testHelper, err)
			assert.NotZero(testHelper, flight.CreatedAt)
			assert.NotZero(testHelper, flight.UpdatedAt)

			err = mock.ExpectationsWereMet()
			assert.NoError(testHelper, err)
		})
	}
}
