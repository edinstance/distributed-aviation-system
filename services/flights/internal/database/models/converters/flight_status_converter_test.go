package converters

import (
	"testing"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
	"github.com/stretchr/testify/assert"
)

func TestToProtoStatus(testHelper *testing.T) {
	tests := []struct {
		name     string
		input    models.FlightStatus
		expected v1.FlightStatus
	}{
		{
			input:    models.FlightStatusScheduled,
			expected: v1.FlightStatus_FLIGHT_STATUS_SCHEDULED,
		},
		{
			input:    models.FlightStatusDelayed,
			expected: v1.FlightStatus_FLIGHT_STATUS_DELAYED,
		},
		{
			input:    models.FlightStatusDeparted,
			expected: v1.FlightStatus_FLIGHT_STATUS_DEPARTED,
		},
		{
			input:    models.FlightStatusInProgress,
			expected: v1.FlightStatus_FLIGHT_STATUS_IN_PROGRESS,
		},
		{
			input:    models.FlightStatusArrived,
			expected: v1.FlightStatus_FLIGHT_STATUS_ARRIVED,
		},
		{
			input:    models.FlightStatusCancelled,
			expected: v1.FlightStatus_FLIGHT_STATUS_CANCELLED,
		},
		{
			input:    "",
			expected: v1.FlightStatus_FLIGHT_STATUS_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		testHelper.Run(tt.name, func(testHelper *testing.T) {
			result := ToProtoStatus(tt.input)
			assert.Equal(testHelper, tt.expected, result)
		})
	}
}

func TestFromProtoStatus(testHelper *testing.T) {
	tests := []struct {
		name     string
		input    v1.FlightStatus
		expected models.FlightStatus
	}{
		{
			name:     "Scheduled status",
			input:    v1.FlightStatus_FLIGHT_STATUS_SCHEDULED,
			expected: models.FlightStatusScheduled,
		},
		{
			name:     "Delayed status",
			input:    v1.FlightStatus_FLIGHT_STATUS_DELAYED,
			expected: models.FlightStatusDelayed,
		},
		{
			name:     "Departed status",
			input:    v1.FlightStatus_FLIGHT_STATUS_DEPARTED,
			expected: models.FlightStatusDeparted,
		},
		{
			name:     "In Progress status",
			input:    v1.FlightStatus_FLIGHT_STATUS_IN_PROGRESS,
			expected: models.FlightStatusInProgress,
		},
		{
			name:     "Arrived status",
			input:    v1.FlightStatus_FLIGHT_STATUS_ARRIVED,
			expected: models.FlightStatusArrived,
		},
		{
			name:     "Cancelled status",
			input:    v1.FlightStatus_FLIGHT_STATUS_CANCELLED,
			expected: models.FlightStatusCancelled,
		},
		{
			name:     "Unspecified status",
			input:    v1.FlightStatus_FLIGHT_STATUS_UNSPECIFIED,
			expected: models.FlightStatusUnspecified,
		},
		{
			name:     "Unknown status defaults to unspecified",
			input:    v1.FlightStatus(999),
			expected: models.FlightStatusUnspecified,
		},
	}

	for _, testCase := range tests {
		testHelper.Run(testCase.name, func(subTest *testing.T) {
			result := FromProtoStatus(testCase.input)
			assert.Equal(subTest, testCase.expected, result)
		})
	}
}

func TestStatusConversionRoundTrip(testHelper *testing.T) {
	statuses := []models.FlightStatus{
		models.FlightStatusScheduled,
		models.FlightStatusDelayed,
		models.FlightStatusDeparted,
		models.FlightStatusInProgress,
		models.FlightStatusArrived,
		models.FlightStatusCancelled,
		models.FlightStatusUnspecified,
	}

	for _, status := range statuses {
		testHelper.Run(string(status), func(subTest *testing.T) {
			proto := ToProtoStatus(status)
			converted := FromProtoStatus(proto)
			assert.Equal(subTest, status, converted)
		})
	}
}
