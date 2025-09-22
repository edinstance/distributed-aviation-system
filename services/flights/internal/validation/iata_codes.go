package validation

import (
	"regexp"
	"strings"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
)

// iataCodeRegex matches a normalised 3-letter IATA airport code.
var (
	iataCodeRegex = regexp.MustCompile(`^[A-Z]{3}$`)
)

// ValidateAndNormalizeIATACode trims, upper-cases, and validates a 3-letter IATA code.
func ValidateAndNormalizeIATACode(code string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	if !iataCodeRegex.MatchString(normalized) {
		return "", exceptions.ErrInvalidIATACode
	}
	return normalized, nil
}
