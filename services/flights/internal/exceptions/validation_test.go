package exceptions

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
)

func TestMapErrorToGrpcCode(testHelper *testing.T) {
	testCases := []struct {
		error               error
		expectedConnectCode connect.Code
	}{
		{ErrInvalidIATACode, connect.CodeInvalidArgument},
		{ErrInvalidFlightNumber, connect.CodeInvalidArgument},
		{ErrInvalidTimes, connect.CodeInvalidArgument},
		{ErrInvalidFlightNumber, connect.CodeInvalidArgument},
		{ErrInvalidInput, connect.CodeInvalidArgument},
		{error: error(nil), expectedConnectCode: connect.CodeInternal},
	}

	for _, testCase := range testCases {
		result := MapErrorToGrpcCode(testCase.error)
		assert.Equal(testHelper, testCase.expectedConnectCode, result)
	}
}
