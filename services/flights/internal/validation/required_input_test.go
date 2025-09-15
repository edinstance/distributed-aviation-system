package validation

import (
	"fmt"
	"testing"
)

func TestValidateRequiredInput(testHelper *testing.T) {
	testCases := []struct {
		fields          map[string]any
		expectedError   error
		expectedSuccess bool
	}{
		{
			fields:          map[string]any{"username": "alice", "password": "secret"},
			expectedError:   nil,
			expectedSuccess: true,
		},
		{
			fields:          map[string]any{"username": "", "password": "secret"},
			expectedError:   fmt.Errorf("missing required field(s): username"),
			expectedSuccess: false,
		},
		{
			fields:          map[string]any{"username": nil, "password": "secret"},
			expectedError:   fmt.Errorf("missing required field(s): username"),
			expectedSuccess: false,
		},
	}

	for _, testCase := range testCases {
		validationError := ValidateRequiredInput(testCase.fields)

		if testCase.expectedSuccess {
			if validationError != nil {
				testHelper.Errorf("expected no error, got %v", validationError)
			}
			continue
		}

		if validationError == nil {
			testHelper.Errorf("expected error %v, got nil", testCase.expectedError)
			continue
		}

		if validationError.Error() != testCase.expectedError.Error() {
			testHelper.Errorf("expected error %q, got %q",
				testCase.expectedError.Error(), validationError.Error())
		}
	}
}
