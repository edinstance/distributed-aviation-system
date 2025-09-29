package aircraft_client

import (
	"context"
	"fmt"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	aircraftv1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/aircraft/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AircraftValidator interface {
	ValidateAircraftExists(ctx context.Context, aircraftID uuid.UUID) error
}

type AircraftClient struct {
	conn   *grpc.ClientConn
	client aircraftv1.AircraftServiceClient
}

func NewAircraftClient(address string) (*AircraftClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("%w: %v", exceptions.ErrDownstreamClientDown, err)
	}

	client := aircraftv1.NewAircraftServiceClient(conn)

	return &AircraftClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *AircraftClient) Close() error {
	return c.conn.Close()
}
