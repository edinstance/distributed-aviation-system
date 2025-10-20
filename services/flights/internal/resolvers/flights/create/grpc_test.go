package create

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
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Helper to create a request with required user context headers
func newRequestWithUserContext(req *v1.CreateFlightRequest) *connect.Request[v1.CreateFlightRequest] {
	connectReq := connect.NewRequest(req)
	connectReq.Header().Set("x-user-sub", "123e4567-e89b-12d3-a456-426614174000")
	connectReq.Header().Set("x-org-id", "987fcdeb-51a2-43d1-9f87-123456789abc")
	connectReq.Header().Set("x-user-roles", "user")
	return connectReq
}

type fakeService struct {
	createFn func(ctx context.Context, number, origin, dest string, dep, arr time.Time, aircraftId uuid.UUID) (*models.Flight, error)
}

func (f *fakeService) CreateFlight(
	ctx context.Context, number, origin, dest string, dep, arr time.Time, aircraftId uuid.UUID,
) (*models.Flight, error) {
	return f.createFn(ctx, number, origin, dest, dep, arr, aircraftId)
}

func TestCreateFlightGRPCValidation(testingHelper *testing.T) {
	dep := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	arr := dep.Add(2 * time.Hour)

	f := &fakeService{
		createFn: func(ctx context.Context, number, origin, dest string, dep, arr time.Time, aircraftId uuid.UUID) (*models.Flight, error) {
			switch {
			case number == "":
				return nil, exceptions.ErrInvalidInput
			case origin == "":
				return nil, exceptions.ErrInvalidInput
			case dest == "":
				return nil, exceptions.ErrInvalidInput
			default:
				return nil, nil
			}
		},
	}

	resolver := NewCreateFlightResolver(f)

	testCases := []struct {
		name    string
		req     *v1.CreateFlightRequest
		wantErr bool
		code    connect.Code
	}{
		{
			name: "missing number",
			req: &v1.CreateFlightRequest{
				Origin:        "LHR",
				Destination:   "LGW",
				DepartureTime: timestamppb.New(dep),
				ArrivalTime:   timestamppb.New(arr),
			},
			wantErr: true,
			code:    connect.CodeInvalidArgument,
		},
		{
			name: "missing origin",
			req: &v1.CreateFlightRequest{
				Number:        "AB123",
				Destination:   "LGW",
				DepartureTime: timestamppb.New(dep),
				ArrivalTime:   timestamppb.New(arr),
			},
			wantErr: true,
			code:    connect.CodeInvalidArgument,
		},
		{
			name: "missing destination",
			req: &v1.CreateFlightRequest{
				Number:        "AB123",
				Origin:        "LHR",
				DepartureTime: timestamppb.New(dep),
				ArrivalTime:   timestamppb.New(arr),
			},
			wantErr: true,
			code:    connect.CodeInvalidArgument,
		},
		{
			name: "missing departure time",
			req: &v1.CreateFlightRequest{
				Number:      "AB123",
				Origin:      "LHR",
				Destination: "LGW",
				ArrivalTime: timestamppb.New(arr),
			},
			wantErr: true,
			code:    connect.CodeInvalidArgument,
		},
		{
			name: "missing arrival time",
			req: &v1.CreateFlightRequest{
				Number:        "AB123",
				Origin:        "LHR",
				Destination:   "LGW",
				DepartureTime: timestamppb.New(dep),
			},
			wantErr: true,
			code:    connect.CodeInvalidArgument,
		},
		{
			name: "invalid arrival time",
			req: &v1.CreateFlightRequest{
				Number:        "AB123",
				Origin:        "LHR",
				Destination:   "LGW",
				DepartureTime: timestamppb.New(dep),
				ArrivalTime:   &timestamppb.Timestamp{Seconds: 253402300800, Nanos: 0},
			},
			wantErr: true,
			code:    connect.CodeInvalidArgument,
		},
		{
			name: "invalid departure time",
			req: &v1.CreateFlightRequest{
				Number:        "AB123",
				Origin:        "LHR",
				Destination:   "LGW",
				ArrivalTime:   timestamppb.New(dep),
				DepartureTime: &timestamppb.Timestamp{Seconds: 253402300800, Nanos: 0},
			},
			wantErr: true,
			code:    connect.CodeInvalidArgument,
		},
	}

	for _, testCase := range testCases {
		testingHelper.Run(testCase.name, func(testingHelper *testing.T) {
			ctx := context.Background()
			req := newRequestWithUserContext(testCase.req)
			resp, err := resolver.CreateFlightGRPC(ctx, req)

			if !testCase.wantErr {
				if err != nil {
					testingHelper.Fatalf("unexpected error: %v", err)
				}
				if resp == nil {
					testingHelper.Fatal("expected non-nil response")
				}
				return
			}

			if err == nil {
				testingHelper.Fatalf("expected error, got resp=%v", resp)
			}
			var connectionError *connect.Error
			if !errors.As(err, &connectionError) {
				testingHelper.Fatalf("expected connect.Error, got %T", err)
			}
			if connectionError.Code() != testCase.code {
				testingHelper.Fatalf("expected code %v, got %v", testCase.code, connectionError.Code())
			}
		})
	}
}

func TestCreateFlightGRPCServiceNotConfigured(testingHelper *testing.T) {
	resolver := NewCreateFlightResolver(nil)

	req := newRequestWithUserContext(&v1.CreateFlightRequest{
		Number:        "AB123",
		Origin:        "LHR",
		Destination:   "LGW",
		DepartureTime: timestamppb.Now(),
		ArrivalTime:   timestamppb.Now(),
	})

	_, err := resolver.CreateFlightGRPC(context.Background(), req)
	if err == nil {
		testingHelper.Fatal("expected error, got nil")
	}
	var connectionError *connect.Error
	if !errors.As(err, &connectionError) {
		testingHelper.Fatalf("expected connect.Error, got %T", err)
	}
	if connectionError.Code() != connect.CodeInternal {
		testingHelper.Errorf("expected CodeInternal, got %v", connectionError.Code())
	}
}

func TestCreateFlightGRPCSuccess(testingHelper *testing.T) {
	dep := time.Date(2025, 2, 1, 10, 0, 0, 0, time.UTC)
	arr := dep.Add(3 * time.Hour)

	flightId := uuid.New()

	f := &fakeService{
		createFn: func(ctx context.Context, number, origin, dest string, depTime, arrTime time.Time, aircraftId uuid.UUID) (*models.Flight, error) {
			return &models.Flight{
				ID:            flightId,
				Number:        number,
				Origin:        origin,
				Destination:   dest,
				DepartureTime: depTime,
				ArrivalTime:   arrTime,
				Status:        models.FlightStatusScheduled,
				AircraftID:    aircraftId,
			}, nil
		},
	}

	resolver := NewCreateFlightResolver(f)

	req := newRequestWithUserContext(&v1.CreateFlightRequest{
		Number:        "XY789",
		Origin:        "LHR",
		Destination:   "LGW",
		DepartureTime: timestamppb.New(dep),
		ArrivalTime:   timestamppb.New(arr),
		AircraftId:    uuid.NewString(),
	})

	resp, err := resolver.CreateFlightGRPC(context.Background(), req)
	if err != nil {
		testingHelper.Fatalf("unexpected error: %v", err)
	}
	if resp.Msg.Flight == nil {
		testingHelper.Fatal("expected flight in response, got nil")
	}

	got := resp.Msg.Flight
	if got.Id != flightId.String() {
		testingHelper.Errorf("expected id %s, got %s", flightId, got.Id)
	}
	if got.Number != "XY789" || got.Origin != "LHR" || got.Destination != "LGW" {
		testingHelper.Errorf("unexpected flight details: %+v", got)
	}
	if got.Status != v1.FlightStatus_FLIGHT_STATUS_SCHEDULED {
		testingHelper.Errorf("expected status scheduled, got %v", got.Status)
	}
}

func TestCreateFlightGRPCServiceError(testingHelper *testing.T) {
	dep := time.Now()
	arr := dep.Add(1 * time.Hour)

	f := &fakeService{
		createFn: func(ctx context.Context, number, origin, dest string, depTime, arrTime time.Time, aircraftId uuid.UUID) (*models.Flight, error) {
			return nil, errors.New("db failure")
		},
	}
	resolver := NewCreateFlightResolver(f)

	req := newRequestWithUserContext(&v1.CreateFlightRequest{
		Number:        "CD456",
		Origin:        "LHR",
		Destination:   "LGW",
		DepartureTime: timestamppb.New(dep),
		ArrivalTime:   timestamppb.New(arr),
		AircraftId:    uuid.NewString(),
	})

	resp, err := resolver.CreateFlightGRPC(context.Background(), req)
	if err == nil {
		testingHelper.Fatal("expected error, got nil")
	}
	if resp != nil {
		testingHelper.Errorf("expected nil response, got %+v", resp)
	}
	var connectionError *connect.Error
	if !errors.As(err, &connectionError) {
		testingHelper.Fatalf("expected connect.Error, got %T", err)
	}
	if connectionError.Code() != connect.CodeInternal {
		testingHelper.Errorf("expected CodeInternal, got %v", connectionError.Code())
	}
}
