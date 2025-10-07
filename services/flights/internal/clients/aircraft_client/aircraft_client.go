package aircraft_client

import (
	"context"
	"fmt"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/metrics"
	aircraftv1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/aircraft/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithUnaryInterceptor(metrics.OutboundGrpcUnaryClientInterceptor()),
	)

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
