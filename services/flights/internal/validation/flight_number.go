package validation

import (
	"regexp"
	"strings"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
)

// flightNumberRegex matches a normalised flight code.
var flightNumberRegex = regexp.MustCompile(`^[A-Z]{2,3}[0-9]{1,6}$`)

// ValidateAndNormalizeFlightNumber trims, upper-cases, and validates a flight code.
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
