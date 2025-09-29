package flight_number

import (
	"errors"
	"testing"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
)

func TestValidateAndNormalizeFlightNumber(testHelper *testing.T) {
	testCases := []struct {
		number             string
		expectedNormalized string
		expectedError      error
	}{
		{"BA12", "BA12", nil},
		{"", "", exceptions.ErrInvalidFlightNumber},
		{"ERROR12", "", exceptions.ErrInvalidFlightNumber},
		{"12BA", "", exceptions.ErrInvalidFlightNumber},
		{"BA1212121212121", "", exceptions.ErrInvalidFlightNumber},
	}

	for _, testCase := range testCases {
		result, err := ValidateAndNormalizeFlightNumber(testCase.number)
		if result != testCase.expectedNormalized {
			testHelper.Errorf("Expected normalization of %s to %s, got %s instead", testCase.number, testCase.expectedNormalized, result)
		}
		if !errors.Is(err, testCase.expectedError) {
			testHelper.Errorf("Expected error for %s to be %v, got %v instead", testCase.number, testCase.expectedError, err)
		}
	}
}
