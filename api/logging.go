package api

import (
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LogResponseStatus(logger logging.Logger, s *status.Status) {
	if s == nil {
		logger.Debug("gRPC response status nil")
		return
	}
	switch s.Code() {
	case codes.OK:
		logger.Debug("gRPC response status", "code", s.Code(), "message", s.Message())
	case codes.Unknown,
		codes.FailedPrecondition,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unimplemented,
		codes.Internal,
		codes.Unavailable,
		codes.DataLoss:
		logger.Error("gRPC response status", "code", s.Code(), "message", s.Message())
	case codes.Canceled,
		codes.InvalidArgument,
		codes.DeadlineExceeded,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.ResourceExhausted,
		codes.Unauthenticated:
		logger.Warn("gRPC response status", "code", s.Code(), "message", s.Message())
	}
}
