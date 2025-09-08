package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
	v1connect "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1/flightsv1connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights"
)

// FlightsServer implements the FlightsServiceHandler interface
type FlightsServer struct {
	createFlightResolver *resolvers.CreateFlightResolver
}

// NewFlightsServer creates a new FlightsServer instance
func NewFlightsServer() *FlightsServer {
	logger.Debug("Creating new FlightsServer")
	return &FlightsServer{
		createFlightResolver: resolvers.NewCreateFlightResolver(),
	}
}

// Ensure FlightsServer implements the interface
var _ v1connect.FlightsServiceHandler = (*FlightsServer)(nil)

func (s *FlightsServer) CreateFlight(
	ctx context.Context,
	req *connect.Request[v1.CreateFlightRequest],
) (*connect.Response[v1.CreateFlightResponse], error) {
	return s.createFlightResolver.CreateFlight(ctx, req)
}
