package models

type FlightStatus string

const (
	FlightStatusScheduled   FlightStatus = "SCHEDULED"
	FlightStatusDelayed     FlightStatus = "DELAYED"
	FlightStatusCancelled   FlightStatus = "CANCELLED"
	FlightStatusDeparted    FlightStatus = "DEPARTED"
	FlightStatusInProgress  FlightStatus = "IN_PROGRESS"
	FlightStatusArrived     FlightStatus = "ARRIVED"
	FlightStatusUnspecified FlightStatus = "UNSPECIFIED"
)
