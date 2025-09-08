package exceptions

import "fmt"

var (
	ErrInvalidIATACode    = fmt.Errorf("IATA code must be exactly 3 uppercase letters A-Z")
	ErrInvalidFlightNumber = fmt.Errorf("flight number must contain airline code (2-3 letters) followed by digits, max 10 characters")
)