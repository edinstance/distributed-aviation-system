package resolvers

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
	"github.com/google/uuid"
)

type CreateFlightResolver struct{}

func NewCreateFlightResolver() *CreateFlightResolver {
	return &CreateFlightResolver{}
}

func (r *CreateFlightResolver) CreateFlight(
	_ context.Context,
	req *connect.Request[v1.CreateFlightRequest],
) (*connect.Response[v1.CreateFlightResponse], error) {
	resp := &v1.CreateFlightResponse{
		Id:            uuid.NewString(),
		Number:        req.Msg.GetNumber(),
		Origin:        req.Msg.GetOrigin(),
		Destination:   req.Msg.GetDestination(),
		DepartureTime: req.Msg.GetDepartureTime(),
		ArrivalTime:   req.Msg.GetArrivalTime(),
		Status:        "SCHEDULED",
	}
	return connect.NewResponse(resp), nil
}
