package inputs

import (
	"fmt"
	"sort"
	"strings"
)

// ValidateRequiredInput checks that each entry in fields is present and non-empty.
// It treats string values that are empty after trimming whitespace, and nil values, as missing.
// If any fields are missing it returns an error listing their names (sorted, comma-separated);
// otherwise it returns nil.
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
		sort.Strings(missing)
		return fmt.Errorf("missing required field(s): %s", strings.Join(missing, ", "))
	}
	return nil
}
