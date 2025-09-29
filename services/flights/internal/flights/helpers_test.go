package flights

import (
	"context"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"

	"github.com/google/uuid"
)

type FakeRepo struct {
	CreateFlightFn func(ctx context.Context, f *models.Flight) error
	GetFlightFn    func(ctx context.Context, id uuid.UUID) (*models.Flight, error)
}

type FakeFlightsCache struct {
	SaveFlightFn func(ctx context.Context, f *models.Flight) error
	GetFlightFn  func(ctx context.Context, id uuid.UUID) (*models.Flight, error)
}

type FakeAircraftClient struct {
	ValidateAircraftExistsFn func(ctx context.Context, id uuid.UUID) error
}

func (f FakeFlightsCache) GetFlight(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	if f.GetFlightFn == nil {
		return nil, nil
	}
	return f.GetFlightFn(ctx, id)
}

func (f FakeFlightsCache) SetFlight(ctx context.Context, flight *models.Flight) error {
	if f.SaveFlightFn == nil {
		return nil
	}
	return f.SaveFlightFn(ctx, flight)
}

func (f *FakeRepo) GetFlightByID(ctx context.Context, id uuid.UUID) (*models.Flight, error) {
	if f.GetFlightFn == nil {
		return nil, nil
	}
	return f.GetFlightFn(ctx, id)
}

func (f *FakeRepo) CreateFlight(ctx context.Context, fl *models.Flight) error {
	if f.CreateFlightFn == nil {
		return nil
	}
	return f.CreateFlightFn(ctx, fl)
}

func (f *FakeAircraftClient) ValidateAircraftExists(ctx context.Context, id uuid.UUID) error {
	if f.ValidateAircraftExistsFn == nil {
		return nil
	}
	return f.ValidateAircraftExistsFn(ctx, id)
}
