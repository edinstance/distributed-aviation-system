package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFlightJSONSerialization(testHelper *testing.T) {
	flight := Flight{
		ID:            uuid.New(),
		Number:        "AA123",
		Origin:        "JFK",
		Destination:   "LAX",
		DepartureTime: time.Now(),
		ArrivalTime:   time.Now().Add(5 * time.Hour),
		Status:        FlightStatusScheduled,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	data, err := json.Marshal(flight)
	assert.NoError(testHelper, err)

	var decoded Flight
	err = json.Unmarshal(data, &decoded)
	assert.NoError(testHelper, err)
	assert.Equal(testHelper, flight.Number, decoded.Number)
}
