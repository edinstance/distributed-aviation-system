package exceptions

import (
	"errors"

	"connectrpc.com/connect"
)

var (
	ErrInvalidIATACode          = errors.New("IATA code must be exactly 3 uppercase letters A-Z")
	ErrInvalidFlightNumber      = errors.New("flight number must contain airline code (2-3 letters) followed by digits, max 10 characters")
	ErrSameOriginAndDestination = errors.New("duplicate origin and destination code")
	ErrInvalidTimes             = errors.New("arrival must be after departure")
	ErrInvalidInput             = errors.New("invalid input")
)

var errorCodeMap = map[error]connect.Code{
	ErrInvalidInput:             connect.CodeInvalidArgument,
	ErrInvalidTimes:             connect.CodeInvalidArgument,
	ErrInvalidFlightNumber:      connect.CodeInvalidArgument,
	ErrInvalidIATACode:          connect.CodeInvalidArgument,
	ErrSameOriginAndDestination: connect.CodeInvalidArgument,
}

// MapErrorToGrpcCode returns the corresponding connect.Code for the provided error.
// It checks the error against the package's sentinel errors in errorCodeMap using
// errors.Is (so wrapped errors are matched). If no mapping is found it returns
// connect.CodeInternal.
func MapErrorToGrpcCode(err error) connect.Code {
	for e, code := range errorCodeMap {
		if errors.Is(err, e) {
			return code
		}
	}
	return connect.CodeInternal
}
