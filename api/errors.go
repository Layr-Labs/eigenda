package api

import (
	"fmt"
	"strconv"

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

// ==================================================================
// API CLIENT ERRORS
// Note: These errors are currently used by api/clients.
// Eventually it might be useful to use them across the entire codebase.
// ==================================================================

// Code below is adapted from https://github.com/aws/smithy-go/blob/main/errors.go

// ErrorAPI is the most generic API and protocol agnostic error interface that we eventually
// want every api clients returned errors to implement.
//
// This way, consumers of the api clients can tell what kind of error they are dealing with,
// most broadly whether it is a client or server fault (see the ErrorFault type).
type ErrorAPI interface {
	error

	// ErrorCode returns the error code for the API exception.
	ErrorCode() ErrorCode
	// ErrorFault returns the fault for the API exception.
	ErrorFault() ErrorFault
}

// ErrorAPIGeneric provides a generic concrete API error type that implements APIError
// and clients can use when they don't have a more concrete error type to use.
// It would typically be
type ErrorAPIGeneric struct {
	Err   error
	Code  ErrorCode
	Fault ErrorFault
}

func NewErrorAPIGeneric(code ErrorCode, err error) *ErrorAPIGeneric {
	errGeneric := &ErrorAPIGeneric{
		Err:   err,
		Code:  code,
		Fault: ErrorFaultUnknown,
	}
	if code >= 400 && code < 500 {
		errGeneric.Fault = ErrorFaultClient
	} else if code >= 500 && code < 600 {
		errGeneric.Fault = ErrorFaultServer
	}
	return errGeneric
}

// ErrorCode returns the error code for the API exception.
func (e *ErrorAPIGeneric) ErrorCode() ErrorCode { return e.Code }

// ErrorFault returns the fault for the API exception.
func (e *ErrorAPIGeneric) ErrorFault() ErrorFault { return e.Fault }

func (e *ErrorAPIGeneric) Error() string {
	return fmt.Sprintf("api error %d: %s", e.Code, e.Error())
}

// We implement Unwrap so that errors.Is and errors.As work as expected.
func (e *ErrorAPIGeneric) Unwrap() error { return e.Err }

var _ ErrorAPI = (*ErrorAPIGeneric)(nil)

// ErrorCode is a subset of HTTP error codes that are relevant to the API.
//
// We might eventually switch to the more precise grpc error codes
// https://cloud.google.com/apis/design/errors#handling_errors
type ErrorCode uint16

const (
	ErrorCodeUnknown    ErrorCode = 0
	ErrorCodeBadRequest ErrorCode = 400
	ErrorCodeInternal   ErrorCode = 500
	// 503 is used to signify that eigenda is temporarily unavailable,
	// and suggest to the caller (most likely some rollup batcher via the eigenda-proxy)
	// to fallback to ethda for some amount of time.
	// See https://github.com/ethereum-optimism/specs/issues/434 for more details.
	ErrorCodeUnavailable ErrorCode = 503
)

func (f ErrorCode) String() string {
	return strconv.Itoa(int(f))
}

// ErrorFault provides the broadest categorization of an error (client, server, or unknown).
type ErrorFault int

// ErrorFault enumeration values
const (
	ErrorFaultUnknown ErrorFault = iota
	ErrorFaultServer
	ErrorFaultClient
)

func (f ErrorFault) String() string {
	switch f {
	case ErrorFaultServer:
		return "server"
	case ErrorFaultClient:
		return "client"
	default:
		return "unknown"
	}
}
