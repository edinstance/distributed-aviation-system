package aircraft_client

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	aircraftv1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/aircraft/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type mockAircraftServiceClient struct {
	resp *aircraftv1.GetAircraftByIdResponse
	err  error
}

func (m *mockAircraftServiceClient) GetAircraftById(ctx context.Context, in *aircraftv1.GetAircraftByIdRequest, opts ...grpc.CallOption) (*aircraftv1.GetAircraftByIdResponse, error) {
	return m.resp, m.err
}

// ensure mock satisfies interface at compile time
var _ aircraftv1.AircraftServiceClient = (*mockAircraftServiceClient)(nil)

func TestValidateAircraftExistsSuccess(t *testing.T) {
	id := uuid.New()
	c := &AircraftClient{client: &mockAircraftServiceClient{resp: &aircraftv1.GetAircraftByIdResponse{Aircraft: &aircraftv1.Aircraft{Id: id.String()}}}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := c.ValidateAircraftExists(ctx, id); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateAircraftExistsNotFound(t *testing.T) {
	id := uuid.New()
	c := &AircraftClient{client: &mockAircraftServiceClient{resp: &aircraftv1.GetAircraftByIdResponse{Aircraft: nil}}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := c.ValidateAircraftExists(ctx, id)
	if !errors.Is(err, exceptions.ErrAircraftNotFound) {
		t.Fatalf("expected ErrAircraftNotFound, got %v", err)
	}
}

func TestValidateAircraftExistsDownstreamError(t *testing.T) {
	id := uuid.New()
	c := &AircraftClient{client: &mockAircraftServiceClient{err: errors.New("downstream error")}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := c.ValidateAircraftExists(ctx, id)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, exceptions.ErrDownstreamClientDown) {
		t.Fatalf("expected ErrDownstreamClientDown, got %v", err)
	}
}
