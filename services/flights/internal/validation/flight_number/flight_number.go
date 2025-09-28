package flight_number

import (
	"regexp"
	"strings"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
)

// flightNumberRegex matches a normalised flight code.
var flightNumberRegex = regexp.MustCompile(`^[A-Z]{2,3}[0-9]{1,6}$`)

// ValidateAndNormalizeFlightNumber normalises a flight code by trimming surrounding whitespace
// and converting it to upper-case, then validates it against the pattern `^[A-Z]{2,3}[0-9]{1,6}$`.
// On success it returns the normalised flight number. If the input is empty or does not match
// the pattern it returns an empty string and exceptions.ErrInvalidFlightNumber.
func ValidateAndNormalizeFlightNumber(number string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(number))
	if normalized == "" {
		return "", exceptions.ErrInvalidFlightNumber
	}
	if !flightNumberRegex.MatchString(normalized) {
		return "", exceptions.ErrInvalidFlightNumber
	}
	return normalized, nil
}
