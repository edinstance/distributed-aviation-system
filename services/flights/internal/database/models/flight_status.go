package models

type FlightStatus string

const (
	FlightStatusScheduled   FlightStatus = "scheduled"
	FlightStatusDelayed     FlightStatus = "delayed"
	FlightStatusCancelled   FlightStatus = "cancelled"
	FlightStatusDeparted    FlightStatus = "departed"
	FlightStatusInProgress  FlightStatus = "in_progress"
	FlightStatusArrived     FlightStatus = "arrived"
	FlightStatusUnspecified FlightStatus = "unspecified"
)
