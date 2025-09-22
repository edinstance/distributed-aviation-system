package converters

import (
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/database/models"
	v1 "github.com/edinstance/distributed-aviation-system/services/flights/internal/protobuf/flights/v1"
)

// ToProtoStatus converts a models.FlightStatus to its protobuf equivalent.
func ToProtoStatus(s models.FlightStatus) v1.FlightStatus {
	switch s {
	case models.FlightStatusScheduled:
		return v1.FlightStatus_FLIGHT_STATUS_SCHEDULED
	case models.FlightStatusDelayed:
		return v1.FlightStatus_FLIGHT_STATUS_DELAYED
	case models.FlightStatusDeparted:
		return v1.FlightStatus_FLIGHT_STATUS_DEPARTED
	case models.FlightStatusInProgress:
		return v1.FlightStatus_FLIGHT_STATUS_IN_PROGRESS
	case models.FlightStatusArrived:
		return v1.FlightStatus_FLIGHT_STATUS_ARRIVED
	case models.FlightStatusCancelled:
		return v1.FlightStatus_FLIGHT_STATUS_CANCELLED
	default:
		return v1.FlightStatus_FLIGHT_STATUS_UNSPECIFIED
	}
}

// FromProtoStatus converts a protobuf FlightStatus to the models.FlightStatus type.
func FromProtoStatus(p v1.FlightStatus) models.FlightStatus {
	switch p {
	case v1.FlightStatus_FLIGHT_STATUS_SCHEDULED:
		return models.FlightStatusScheduled
	case v1.FlightStatus_FLIGHT_STATUS_DELAYED:
		return models.FlightStatusDelayed
	case v1.FlightStatus_FLIGHT_STATUS_DEPARTED:
		return models.FlightStatusDeparted
	case v1.FlightStatus_FLIGHT_STATUS_IN_PROGRESS:
		return models.FlightStatusInProgress
	case v1.FlightStatus_FLIGHT_STATUS_ARRIVED:
		return models.FlightStatusArrived
	case v1.FlightStatus_FLIGHT_STATUS_CANCELLED:
		return models.FlightStatusCancelled
	default:
		return models.FlightStatusUnspecified
	}
}
