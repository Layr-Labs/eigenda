package controller

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsInvalidDisperserIDError checks if the error is an invalid disperser ID error.
// These errors occur when a validator rejects a request because it doesn't recognize
// the disperser ID used to sign the request.
//
// The error is identified by:
// - gRPC status code: InvalidArgument
// - Error message containing authentication failure patterns related to disperser key/address lookup
func IsInvalidDisperserIDError(err error) bool {
	if err == nil {
		return false
	}

	// Check if this is a gRPC error with InvalidArgument status
	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	if st.Code() != codes.InvalidArgument {
		return false
	}

	// Check for authentication failure patterns
	errMsg := st.Message()
	if !strings.Contains(errMsg, "failed to authenticate") {
		return false
	}

	// Check for disperser key/address lookup failures
	// These indicate the validator doesn't know about this disperser ID
	return strings.Contains(errMsg, "failed to get disperser key") ||
		strings.Contains(errMsg, "failed to get disperser address")
}
