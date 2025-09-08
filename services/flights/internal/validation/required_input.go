package validation

import (
	"fmt"
	"strings"
)

func ValidateRequiredInput(fields map[string]any) error {
	var missing []string
	for name, val := range fields {
		if stringValue, ok := val.(string); ok {
			if strings.TrimSpace(stringValue) == "" {
				missing = append(missing, name)
			}
			continue
		}

		if val == nil {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required field(s): %s", strings.Join(missing, ", "))
	}
	return nil
}
