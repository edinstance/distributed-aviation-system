package models

import (
	"time"

	"github.com/google/uuid"
)

type Flight struct {
	ID            uuid.UUID
	Number        string
	Origin        string
	Destination   string
	DepartureTime time.Time
	ArrivalTime   time.Time
	Status        string
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
