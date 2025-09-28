package server

import (
	"context"

	"connectrpc.com/connect"
	cacheRepository "github.com/edinstance/distributed-aviation-system/services/flights/internal/cache/repositories/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/config"
	flightRepository "github.com/edinstance/distributed-aviation-system/services/flights/internal/database/repositories/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
	v1connect "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1/flightsv1connect"
	createFlightsResolver "github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/create"
	getFlightsResolver "github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/get"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// FlightsServer implements the FlightsServiceHandler interface
type GrpcFlightsServer struct {
	createFlightResolver *createFlightsResolver.FlightResolver
	getFlightsResolver   *getFlightsResolver.FlightResolver
}

func NewGrpcFlightsServer(pool *pgxpool.Pool, client *redis.Client) *GrpcFlightsServer {
	logger.Debug("Creating new FlightsServer")
	dbRepo := flightRepository.NewFlightRepository(pool)
	cacheRepo := cacheRepository.NewRedisFlightRepository(client, config.App.CacheTTL)

	flightService := flights.NewFlightsService(dbRepo, cacheRepo)

	return &GrpcFlightsServer{
		createFlightResolver: createFlightsResolver.NewCreateFlightResolver(flightService),
		getFlightsResolver:   getFlightsResolver.NewGetFlightResolver(flightService),
	}
}

// Ensure FlightsServer implements the interface
var _ v1connect.FlightsServiceHandler = (*GrpcFlightsServer)(nil)

func (s *GrpcFlightsServer) CreateFlight(
	ctx context.Context,
	req *connect.Request[v1.CreateFlightRequest],
) (*connect.Response[v1.CreateFlightResponse], error) {
	return s.createFlightResolver.CreateFlightGRPC(ctx, req)
}

func (s *GrpcFlightsServer) GetFlightById(
	ctx context.Context,
	c *connect.Request[v1.GetFlightByIdRequest],
) (*connect.Response[v1.GetFlightByIdResponse], error) {
	return s.getFlightsResolver.GetFlightByIdGRPC(ctx, c)
}
