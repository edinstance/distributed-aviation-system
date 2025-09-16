package server

import (
	"context"

	"connectrpc.com/connect"
	flightRepository "github.com/edinstance/distributed-aviation-system/services/flights/internal/database/repositories/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
	v1connect "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1/flightsv1connect"
	flightsResolver "github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/create"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FlightsServer implements the FlightsServiceHandler interface
type FlightsServer struct {
	createFlightResolver *flightsResolver.FlightResolver
}

// NewFlightsServer creates a new FlightsServer instance
func NewFlightsServer(pool *pgxpool.Pool) *FlightsServer {
	logger.Debug("Creating new FlightsServer")
	flightService := flights.NewFlightsService(flightRepository.NewFlightRepository(pool))
	return &FlightsServer{
		createFlightResolver: flightsResolver.NewCreateFlightResolver(flightService),
	}
}

// Ensure FlightsServer implements the interface
var _ v1connect.FlightsServiceHandler = (*FlightsServer)(nil)

func (s *FlightsServer) CreateFlight(
	ctx context.Context,
	req *connect.Request[v1.CreateFlightRequest],
) (*connect.Response[v1.CreateFlightResponse], error) {
	return s.createFlightResolver.CreateFlightGRPC(ctx, req)
}
