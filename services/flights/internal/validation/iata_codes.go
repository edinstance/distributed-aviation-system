package validation

import (
	"regexp"
	"strings"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
)

var (
	iataCodeRegex = regexp.MustCompile(`^[A-Z]{3}$`)
)

func ValidateAndNormalizeIATACode(code string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	if !iataCodeRegex.MatchString(normalized) {
		return "", exceptions.ErrInvalidIATACode
	}
	return normalized, nil
}
