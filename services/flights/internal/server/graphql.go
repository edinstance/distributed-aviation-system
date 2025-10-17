package server

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	cacheRepository "github.com/edinstance/distributed-aviation-system/services/flights/internal/cache/repositories/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/clients/aircraft_client"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/config"
	flightRepository "github.com/edinstance/distributed-aviation-system/services/flights/internal/database/repositories/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/directives"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/flights"
	graphqlschema "github.com/edinstance/distributed-aviation-system/services/flights/internal/graphql"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/graphql/resolvers"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/metrics"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/create"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/get"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func newGraphQLHandler(pool *pgxpool.Pool, client *redis.Client) http.Handler {
	logger.Info("Setting up GraphQL Handler")
	dbRepo := flightRepository.NewFlightRepository(pool)
	cacheRepo := cacheRepository.NewRedisFlightRepository(client, config.App.CacheTTL)
	aircraftClient, aircraftClientErr := aircraft_client.NewAircraftClient(config.App.AircraftServiceGrpcUrl)

	if aircraftClientErr != nil {
		logger.Error("Failed to create aircraft client", "err", aircraftClientErr)
		return nil
	}

	flightService := flights.NewFlightsService(dbRepo, cacheRepo, aircraftClient)
	graphqlCreateFlightResolver := create.NewCreateFlightResolver(flightService)
	graphqlGetFlightResolver := get.NewGetFlightResolver(flightService)

	resolver := &resolvers.Resolver{
		CreateFlightResolver: graphqlCreateFlightResolver,
		GetFlightResolver:    graphqlGetFlightResolver,
	}

	srv := handler.New(
		graphqlschema.NewExecutableSchema(
			graphqlschema.Config{
				Resolvers: resolver,
				Directives: graphqlschema.DirectiveRoot{
					Authentication: directives.AuthenticationDirective,
				},
			},
		),
	)

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

	// Tracing
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		tracer := otel.Tracer("flights-service")

		operationCtx := graphql.GetOperationContext(ctx)
		if operationCtx == nil {
			return next(ctx)
		}

		operationName := "unknown"
		if operationCtx.OperationName != "" {
			operationName = operationCtx.OperationName
		} else if operationCtx.Doc != nil && len(operationCtx.Doc.Operations) > 0 {
			operationName = string(operationCtx.Doc.Operations[0].Operation)
		}

		ctx, span := tracer.Start(ctx, "graphql."+operationName)
		span.SetAttributes(
			attribute.String("graphql.operation.name", operationName),
			attribute.String("graphql.operation.type", string(operationCtx.Operation.Operation)),
		)

		if operationCtx.Variables != nil {
			span.SetAttributes(attribute.Int("graphql.variables.count", len(operationCtx.Variables)))
		}

		respHandler := next(ctx)

		return func(ctx context.Context) *graphql.Response {
			defer span.End()
			return respHandler(ctx)
		}
	})

	// Metrics
	srv.AroundOperations(metrics.GraphQLMetricsInterceptor)

	logger.Info("GraphQL Handler setup")
	return srv
}
