package controller

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestIsInvalidDisperserIDError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		assert.False(t, IsInvalidDisperserIDError(nil))
	})

	t.Run("non-gRPC error", func(t *testing.T) {
		err := errors.New("some error")
		assert.False(t, IsInvalidDisperserIDError(err))
	})

	t.Run("wrong gRPC status code", func(t *testing.T) {
		err := status.Error(codes.Internal, "failed to authenticate request: failed to get disperser key")
		assert.False(t, IsInvalidDisperserIDError(err))
	})

	t.Run("InvalidArgument but missing authentication pattern", func(t *testing.T) {
		err := status.Error(codes.InvalidArgument, "some other error")
		assert.False(t, IsInvalidDisperserIDError(err))
	})

	t.Run("InvalidArgument with authentication but missing disperser key/address", func(t *testing.T) {
		err := status.Error(codes.InvalidArgument, "failed to authenticate request: wrong signature")
		assert.False(t, IsInvalidDisperserIDError(err))
	})

	t.Run("valid invalid disperser ID error - disperser key", func(t *testing.T) {
		err := status.Error(codes.InvalidArgument, "failed to authenticate request: failed to get disperser key")
		assert.True(t, IsInvalidDisperserIDError(err))
	})

	t.Run("valid invalid disperser ID error - disperser address", func(t *testing.T) {
		err := status.Error(codes.InvalidArgument, "failed to authenticate request: failed to get disperser address")
		assert.True(t, IsInvalidDisperserIDError(err))
	})

	t.Run("valid error with nested message", func(t *testing.T) {
		err := status.Error(codes.InvalidArgument, "failed to authenticate request: failed to verify request: failed to get disperser key: key not found")
		assert.True(t, IsInvalidDisperserIDError(err))
	})

	t.Run("case sensitivity", func(t *testing.T) {
		// Should not match if case is different
		err := status.Error(codes.InvalidArgument, "Failed to Authenticate Request: Failed to Get Disperser Key")
		assert.False(t, IsInvalidDisperserIDError(err))
	})
}
