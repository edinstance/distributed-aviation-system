package create

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models/converters"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/middleware"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	userContext "github.com/edinstance/distributed-aviation-system/services/flights/internal/context"
)

func extractUserContextFromMetadata(ctx context.Context, req *connect.Request[v1.CreateFlightRequest]) (context.Context, error) {
	logger.Debug("All gRPC headers", "headers", req.Header())

	userSub := req.Header().Get("x-user-sub")
	orgID := req.Header().Get("x-org-id")
	roles := req.Header().Get("x-user-roles")

	logger.Debug("Extracting user context from gRPC metadata", "userSub", userSub, "orgID", orgID, "roles", roles)

	// Require both user sub and org ID
	if userSub == "" {
		logger.Warn("Missing required user sub in metadata")
		return ctx, errors.New("missing required user authentication")
	}

	if orgID == "" {
		logger.Warn("Missing required organization ID in metadata")
		return ctx, errors.New("missing required organization context")
	}

	var parsedUserID, parsedOrgID uuid.UUID
	var err error

	parsedUserID, err = uuid.Parse(userSub)
	if err != nil {
		logger.Warn("Invalid user ID in metadata", "userSub", userSub, "err", err)
		return ctx, errors.New("invalid user ID format")
	}

	parsedOrgID, err = uuid.Parse(orgID)
	if err != nil {
		logger.Warn("Invalid organization ID in metadata", "orgID", orgID, "err", err)
		return ctx, errors.New("invalid organization ID format")
	}

	userCtx := &userContext.UserContext{
		UserID: parsedUserID,
		OrgID:  parsedOrgID,
		Roles:  roles,
	}

	logger.Debug("Created user context", "userID", parsedUserID, "orgID", parsedOrgID, "roles", roles)

	return middleware.SetUserContextInContext(ctx, userCtx), nil
}

func (r *FlightResolver) CreateFlightGRPC(
	ctx context.Context,
	req *connect.Request[v1.CreateFlightRequest],
) (*connect.Response[v1.CreateFlightResponse], error) {
	logger.Debug("CreateFlight request", "number", req.Msg.GetNumber())

	// Extract user context from gRPC metadata
	ctx, err := extractUserContextFromMetadata(ctx, req)
	if err != nil {
		logger.Error("Failed to extract user context", "err", err)
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

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
