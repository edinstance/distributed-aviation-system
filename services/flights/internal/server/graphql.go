package server

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	flightRepository "github.com/edinstance/distributed-aviation-system/services/flights/internal/database/repositories/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/graphql"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/graphql/resolvers"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/create"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/get"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

func newGraphQLHandler(pool *pgxpool.Pool) http.Handler {
	logger.Info("Setting up GraphQL Handler")
	flightService := flights.NewFlightsService(flightRepository.NewFlightRepository(pool))
	graphqlCreateFlightResolver := create.NewCreateFlightResolver(flightService)
	graphqlGetFlightResolver := get.NewGetFlightResolver(flightService)

	resolver := &resolvers.Resolver{
		CreateFlightResolver: graphqlCreateFlightResolver,
		GetFlightResolver:    graphqlGetFlightResolver,
	}

	srv := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	logger.Info("GraphQL Handler setup")
	return srv
}
