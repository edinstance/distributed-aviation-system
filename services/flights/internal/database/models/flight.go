package models

import (
	"time"

	"github.com/google/uuid"
)

type Flight struct {
	ID             uuid.UUID    `db:"id" json:"id"`
	Number         string       `db:"number" json:"number"`
	Origin         string       `db:"origin" json:"origin"`
	Destination    string       `db:"destination" json:"destination"`
	DepartureTime  time.Time    `db:"departure_time" json:"departure_time"`
	ArrivalTime    time.Time    `db:"arrival_time" json:"arrival_time"`
	Status         FlightStatus `db:"status" json:"status"`
	AircraftID     uuid.UUID    `db:"aircraft_id" json:"aircraft_id"`
	CreatedBy      uuid.UUID    `db:"created_by" json:"-"`
	LastUpdatedBy  uuid.UUID    `db:"last_updated_by" json:"-"`
	OrganizationID uuid.UUID    `db:"organization_id" json:"-"`
	CreatedAt      time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time    `db:"updated_at" json:"updated_at"`
}

func (Flight) IsEntity() {}
