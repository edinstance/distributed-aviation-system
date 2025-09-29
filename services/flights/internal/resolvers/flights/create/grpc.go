package create

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models/converters"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *FlightResolver) CreateFlightGRPC(
	ctx context.Context,
	req *connect.Request[v1.CreateFlightRequest],
) (*connect.Response[v1.CreateFlightResponse], error) {
	logger.Debug("CreateFlight request", "number", req.Msg.GetNumber())

	if r.service == nil {
		logger.Error("CreateFlight service not configured")
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("service not configured"),
		)
	}

	number := req.Msg.GetNumber()
	origin := req.Msg.GetOrigin()
	dest := req.Msg.GetDestination()

	departureTS := req.Msg.GetDepartureTime()
	arrivalTS := req.Msg.GetArrivalTime()

	aircraftId, idErr := uuid.Parse(req.Msg.GetAircraftId())

	if idErr != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid aircraft ID"))
	}

	if departureTS == nil || arrivalTS == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("missing required timestamp(s)"))
	}
	if err := arrivalTS.CheckValid(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid arrival_time"))
	}
	if err := departureTS.CheckValid(); err != nil {
		logger.Debug("Invalid departure timestamp", "err", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid departure_time"))
	}

	flight, err := r.service.CreateFlight(
		ctx,
		number,
		origin,
		dest,
		departureTS.AsTime(),
		arrivalTS.AsTime(),
		aircraftId,
	)

	if err != nil {
		logger.Error("Failed to create flight", "err", err)
		return nil, connect.NewError(exceptions.MapErrorToGrpcCode(err), err)
	}

	resp := &v1.CreateFlightResponse{
		Flight: &v1.Flight{
			Id:            flight.ID.String(),
			Number:        flight.Number,
			Origin:        flight.Origin,
			Destination:   flight.Destination,
			DepartureTime: timestamppb.New(flight.DepartureTime),
			ArrivalTime:   timestamppb.New(flight.ArrivalTime),
			Status:        converters.ToProtoStatus(flight.Status),
			AircraftId:    flight.AircraftID.String(),
		},
	}

	logger.Debug("CreateFlight response created", "number", flight.Number, "id", flight.ID)
	return connect.NewResponse(resp), nil
}
