package api

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The canonical errors from the EigenDA gRPC API endpoints.
//
// Notes:
// - Start the error space small (but sufficient), and expand when there is an important
//   failure case to separate out.
// - Avoid simply wrapping system-internal errors without checking if they are appropriate
//   in user-facing errors defined here. Consider map and convert system-internal errors
//   before return to users from APIs.

func NewGRPCError(code codes.Code, msg string) error {
	return status.Error(code, msg)
}

// HTTP Mapping: 400 Bad Request
func NewInvalidArgError(msg string) error {
	return NewGRPCError(codes.InvalidArgument, msg)
}

// HTTP Mapping: 404 Not Found
func NewNotFoundError(msg string) error {
	return NewGRPCError(codes.NotFound, msg)
}

// HTTP Mapping: 429 Too Many Requests
func NewResourceExhaustedError(msg string) error {
	return NewGRPCError(codes.ResourceExhausted, msg)
}

// HTTP Mapping: 500 Internal Server Error
func NewInternalError(msg string) error {
	return NewGRPCError(codes.Internal, msg)
}
