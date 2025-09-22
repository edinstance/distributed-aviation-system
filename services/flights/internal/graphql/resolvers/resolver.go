package resolvers

import (
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/create"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/resolvers/flights/get"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CreateFlightResolver *create.FlightResolver
	GetFlightResolver    *get.FlightResolver
}
