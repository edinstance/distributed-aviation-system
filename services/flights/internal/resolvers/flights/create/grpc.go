package create

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	app "github.com/edinstance/distributed-aviation-system/services/flights/internal/flights"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateFlightResolver struct {
	service *app.Service
}

func NewCreateFlightResolver(service *app.Service) *CreateFlightResolver {
	return &CreateFlightResolver{service: service}
}

func (r *CreateFlightResolver) CreateFlightGRPC(
	ctx context.Context,
	req *connect.Request[v1.CreateFlightRequest],
) (*connect.Response[v1.CreateFlightResponse], error) {
	logger.Debug("CreateFlight request", "number", req.Msg.GetNumber())

	if req.Msg.GetDepartureTime() == nil || req.Msg.GetArrivalTime() == nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("departure_time and arrival_time are required"),
		)
	}

	flight, err := r.service.CreateFlight(
		ctx,
		req.Msg.GetNumber(),
		req.Msg.GetOrigin(),
		req.Msg.GetDestination(),
		req.Msg.GetDepartureTime().AsTime(),
		req.Msg.GetArrivalTime().AsTime(),
	)

	if err != nil {
		logger.Error("Failed to create flight", "err", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Convert string status to protobuf enum
	var statusMap = map[string]v1.FlightStatus{
		"SCHEDULED": v1.FlightStatus_FLIGHT_STATUS_SCHEDULED,
		"DELAYED":   v1.FlightStatus_FLIGHT_STATUS_DELAYED,
		"CANCELLED": v1.FlightStatus_FLIGHT_STATUS_CANCELLED,
	}
	status := statusMap[flight.Status]
	if status == v1.FlightStatus(0) {
		status = v1.FlightStatus_FLIGHT_STATUS_UNSPECIFIED
	}

	logger.Debug("CreateFlight response", "number", req.Msg.GetNumber(), "status", status)

	resp := &v1.CreateFlightResponse{
		Flight: &v1.Flight{
			Id:            flight.ID.String(),
			Number:        flight.Number,
			Origin:        flight.Origin,
			Destination:   flight.Destination,
			DepartureTime: timestamppb.New(flight.DepartureTime),
			ArrivalTime:   timestamppb.New(flight.ArrivalTime),
			Status:        status,
			CreatedAt:     timestamppb.New(flight.CreatedAt),
			UpdatedAt:     timestamppb.New(flight.UpdatedAt),
		},
	}

	return connect.NewResponse(resp), nil
}
