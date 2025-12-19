package middleware

import (
	"errors"

	"google.golang.org/grpc/status"
)

// BlacklistableError marks an error as counting toward disperser blacklisting.
//
// It wraps an underlying error (typically one that maps to a gRPC status) so the caller can both:
// - return the original error to the client unchanged (including status/code), and
// - allow middleware to detect the "blacklist class" via errors.As and record a strike.
type BlacklistableError struct {
	Reason string
	Err    error
}

func (e *BlacklistableError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *BlacklistableError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// GRPCStatus delegates to the wrapped error if it supports it, preserving status codes.
func (e *BlacklistableError) GRPCStatus() *status.Status {
	if e == nil {
		return status.New(0, "")
	}
	type grpcStatusProvider interface {
		GRPCStatus() *status.Status
	}
	if e.Err == nil {
		return status.New(0, "")
	}
	if p, ok := e.Err.(grpcStatusProvider); ok {
		return p.GRPCStatus()
	}
	return status.Convert(e.Err)
}

// WrapBlacklistable wraps err as a BlacklistableError with the given reason.
func WrapBlacklistable(reason string, err error) error {
	if err == nil {
		return nil
	}
	return &BlacklistableError{Reason: reason, Err: err}
}

// AsBlacklistable returns the BlacklistableError if err is (or wraps) one.
func AsBlacklistable(err error) (*BlacklistableError, bool) {
	var be *BlacklistableError
	if errors.As(err, &be) && be != nil {
		return be, true
	}
	return nil, false
}
