package validation

import (
	"errors"
	"testing"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/exceptions"
)

func TestValidateAndNormalizeIATACode(testHelper *testing.T) {
	testCases := []struct {
		code               string
		expectedNormalized string
		expectedError      error
	}{
		{"LHR", "LHR", nil},
		{"lhr", "LHR", nil},
		{"", "", exceptions.ErrInvalidIATACode},
		{"ERROR", "", exceptions.ErrInvalidIATACode},
		{"123", "", exceptions.ErrInvalidIATACode},
		{"123456", "", exceptions.ErrInvalidIATACode},
	}

	for _, testCase := range testCases {
		result, err := ValidateAndNormalizeIATACode(testCase.code)
		if result != testCase.expectedNormalized {
			testHelper.Errorf("Expected normalization of %s to %s, got %s instead", testCase.code, testCase.expectedNormalized, result)
		}
		if !errors.Is(err, testCase.expectedError) {
			testHelper.Errorf("Expected error for %s to be %v, got %v instead", testCase.code, testCase.expectedError, err)
		}
	}
}
