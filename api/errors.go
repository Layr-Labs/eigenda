package api

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
