package api

import (
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

// The canonical errors from the EigenDA gRPC API endpoints.
//
// Notes:
// - We start with a small (but sufficient) subset of google's error code convention, 
//   and expand when there is an important failure case to separate out. See:
//   https://cloud.google.com/apis/design/errors#handling_errors
// - Make sure that internally propagated errors are eventually wrapped in one of the 
//   user-facing errors defined here, since grpc otherwise returns an UNKNOWN error code,
//   which is harder to debug and understand for users.

func newErrorGRPC(code codes.Code, msg string) error {
	return status.Error(code, msg)
}

// HTTP Mapping: 400 Bad Request
func NewErrorInvalidArg(msg string) error {
	return newErrorGRPC(codes.InvalidArgument, msg)
}

// HTTP Mapping: 404 Not Found
func NewErrorNotFound(msg string) error {
	return newErrorGRPC(codes.NotFound, msg)
}

// HTTP Mapping: 429 Too Many Requests
func NewErrorResourceExhausted(msg string) error {
	return newErrorGRPC(codes.ResourceExhausted, msg)
}

// HTTP Mapping: 500 Internal Server Error
func NewErrorInternal(msg string) error {
	return newErrorGRPC(codes.Internal, msg)
}

// HTTP Mapping: 501 Not Implemented
func NewErrorUnimplemented() error {
	return newErrorGRPC(codes.Unimplemented, "not implemented")
}

// HTTP Mapping: 503 Service Unavailable
// Unavailable is used instead of 500 to indicate to the client that it can retry the operation.
// See the documentation for the FAILED_PRECONDITION error code in https://grpc.io/docs/guides/status-codes/
// which compares FAILED_PRECONDITION, UNAVAILABLE, and ABORTED.
func NewErrorUnavailable(msg string) error {
	return newErrorGRPC(codes.Unavailable, msg)
}

// HTTP Mapping: 503 Service Unavailable
// NewErrorUnavailableWithRetry is like NewUnavailableError, but allows the caller to specify a retry delay.
func NewErrorUnavailableWithRetry(msg string, delay time.Duration) error {
	st := status.New(codes.Unavailable, msg)

	retry := &errdetails.RetryInfo{
		RetryDelay: durationpb.New(delay),
	}

	statusWithDetails, err := st.WithDetails(retry)
	if err != nil {
		// If adding details failed, just return the status without details
		return st.Err()
	}

	return statusWithDetails.Err()
}
