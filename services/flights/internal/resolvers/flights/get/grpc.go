package get

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models/converters"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *FlightResolver) GetFlightByIdGRPC(
	ctx context.Context,
	req *connect.Request[v1.GetFlightByIdRequest],
) (*connect.Response[v1.GetFlightByIdResponse], error) {

	if r.service == nil {
		logger.Error("GetFlight service not configured")
		return nil, errors.New("service not configured")
	}

	id := req.Msg.GetId()

	logger.Debug("GetFlight GRPC request", "id", id)

	flightId, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid flight ID format", "id", id, "err", err)
		return nil, errors.New("invalid flight ID format")
	}

	flight, err := r.service.GetFlightByID(ctx, flightId)
	if err != nil {
		logger.Error("Failed to get flight", "id", id, "err", err)
		return nil, err
	}

	logger.Debug("GetFlight GRPC response retrieved", "id", flight.ID)

	resp := &v1.GetFlightByIdResponse{
		Flight: &v1.Flight{
			Id:            flight.ID.String(),
			Number:        flight.Number,
			Origin:        flight.Origin,
			Destination:   flight.Destination,
			DepartureTime: timestamppb.New(flight.DepartureTime),
			ArrivalTime:   timestamppb.New(flight.ArrivalTime),
			Status:        converters.ToProtoStatus(flight.Status),
			CreatedAt:     timestamppb.New(flight.CreatedAt),
			UpdatedAt:     timestamppb.New(flight.UpdatedAt),
		},
	}
	return connect.NewResponse(resp), nil
}
