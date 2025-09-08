package validation

import (
	"regexp"
	"strings"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
)

var flightNumberRegex = regexp.MustCompile(`^[A-Z]{2,3}[0-9]{1,6}$`)

func ValidateAndNormalizeFlightNumber(number string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(number))
	if len(normalized) == 0 || len(normalized) > 10 {
		return "", exceptions.ErrInvalidFlightNumber
	}
	if !flightNumberRegex.MatchString(normalized) {
		return "", exceptions.ErrInvalidFlightNumber
	}
	return normalized, nil
}
