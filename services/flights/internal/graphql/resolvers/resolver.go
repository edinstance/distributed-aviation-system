package resolvers

import "github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/create"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CreateFlightResolver *create.FlightResolver
}
