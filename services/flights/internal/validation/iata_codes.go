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

// ValidateAndNormalizeIATACode normalises an IATA airport code by trimming
// surrounding whitespace and converting to upper-case, then validates it must
// consist of three ASCII letters. On success it returns the normalised code.
// If the code is not a valid 3‑letter IATA code it returns an empty string and
// exceptions.ErrInvalidIATACode.
func ValidateAndNormalizeIATACode(code string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	if !iataCodeRegex.MatchString(normalized) {
		return "", exceptions.ErrInvalidIATACode
	}
	return normalized, nil
}
