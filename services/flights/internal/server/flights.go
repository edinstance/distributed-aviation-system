package server

import (
	"context"

	"connectrpc.com/connect"
	flightRepository "github.com/edinstance/distributed-aviation-system/services/flights/internal/database/repositories/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
	v1connect "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1/flightsv1connect"
	createFlightsResolver "github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/create"
	getFlightsResolver "github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/get"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FlightsServer implements the FlightsServiceHandler interface
type FlightsServer struct {
	createFlightResolver *createFlightsResolver.FlightResolver
	getFlightsResolver   *getFlightsResolver.FlightResolver
}

// NewFlightsServer creates and returns a FlightsServer wired with resolvers for creating
// and retrieving flights. The provided Postgres connection pool is used to construct the
// NewFlightsServer creates a FlightsServer wired with a flight repository and service
// backed by the provided Postgres connection pool.
//
// The returned server contains create and get flight resolvers which share a single
// flights service built from a repository constructed using the given pgxpool.Pool.
func NewFlightsServer(pool *pgxpool.Pool) *FlightsServer {
	logger.Debug("Creating new FlightsServer")
	flightService := flights.NewFlightsService(flightRepository.NewFlightRepository(pool))
	return &FlightsServer{
		createFlightResolver: createFlightsResolver.NewCreateFlightResolver(flightService),
		getFlightsResolver:   getFlightsResolver.NewGetFlightResolver(flightService),
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

func (s *FlightsServer) GetFlightById(
	ctx context.Context,
	c *connect.Request[v1.GetFlightByIdRequest],
) (*connect.Response[v1.GetFlightByIdResponse], error) {
	return s.getFlightsResolver.GetFlightByIdGRPC(ctx, c)
}
