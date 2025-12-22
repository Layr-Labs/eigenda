package middleware

import (
	"context"
	"testing"
	"time"

	validatorpb "github.com/Layr-Labs/eigenda/api/grpc/validator"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockRequestAuthenticator struct {
	t *testing.T

	authenticateFn func(ctx context.Context, request *validatorpb.StoreChunksRequest, now time.Time) ([]byte, error)
}

func (m *mockRequestAuthenticator) AuthenticateStoreChunksRequest(
	ctx context.Context,
	request *validatorpb.StoreChunksRequest,
	now time.Time,
) ([]byte, error) {
	require.NotNil(m.t, m.t)
	require.NotNil(m.t, m.authenticateFn, "authenticateFn must be set")
	return m.authenticateFn(ctx, request, now)
}

func (m *mockRequestAuthenticator) IsDisperserAuthorized(uint32, *corev2.Batch) bool {
	// Not used by the interceptor; included to satisfy interface changes if any.
	return true
}

func TestStoreChunksDisperserAuthAndRateLimitInterceptor_PassThroughForOtherMethods(t *testing.T) {
	t.Parallel()

	var authCalled bool
	auth := &mockRequestAuthenticator{
		t: t,
		authenticateFn: func(context.Context, *validatorpb.StoreChunksRequest, time.Time) ([]byte, error) {
			authCalled = true
			return nil, nil
		},
	}

	interceptor := StoreChunksDisperserAuthAndRateLimitInterceptor(nil, auth)

	handlerCalled := false
	_, err := interceptor(
		context.Background(),
		&validatorpb.StoreChunksRequest{DisperserID: 1},
		&grpc.UnaryServerInfo{FullMethod: validatorpb.Dispersal_GetNodeInfo_FullMethodName},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return "ok", nil
		},
	)
	require.NoError(t, err)
	require.True(t, handlerCalled)
	require.False(t, authCalled, "auth should not be called for other methods")
}

func TestStoreChunksDisperserAuthAndRateLimitInterceptor_RejectsWhenAuthFails(t *testing.T) {
	t.Parallel()

	auth := &mockRequestAuthenticator{
		t: t,
		authenticateFn: func(context.Context, *validatorpb.StoreChunksRequest, time.Time) ([]byte, error) {
			return nil, status.Error(codes.Unauthenticated, "bad sig")
		},
	}

	interceptor := StoreChunksDisperserAuthAndRateLimitInterceptor(nil, auth)

	handlerCalled := false
	_, err := interceptor(
		context.Background(),
		&validatorpb.StoreChunksRequest{DisperserID: 7},
		&grpc.UnaryServerInfo{FullMethod: validatorpb.Dispersal_StoreChunks_FullMethodName},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return "ok", nil
		},
	)
	require.Error(t, err)
	require.False(t, handlerCalled)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestStoreChunksDisperserAuthAndRateLimitInterceptor_RejectsWhenRateLimited(t *testing.T) {
	t.Parallel()

	auth := &mockRequestAuthenticator{
		t: t,
		authenticateFn: func(context.Context, *validatorpb.StoreChunksRequest, time.Time) ([]byte, error) {
			return nil, nil
		},
	}

	limiter := NewDisperserRateLimiter(nil, 1, 1) // burst 1
	now := time.Now()
	require.True(t, limiter.Allow(9, now))
	require.False(t, limiter.Allow(9, now)) // exhaust immediately

	interceptor := StoreChunksDisperserAuthAndRateLimitInterceptor(limiter, auth)

	handlerCalled := false
	_, err := interceptor(
		context.Background(),
		&validatorpb.StoreChunksRequest{DisperserID: 9},
		&grpc.UnaryServerInfo{FullMethod: validatorpb.Dispersal_StoreChunks_FullMethodName},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return "ok", nil
		},
	)
	require.Error(t, err)
	require.False(t, handlerCalled)
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestStoreChunksDisperserAuthAndRateLimitInterceptor_AllowsAndInjectsAuthenticatedDisperserID(t *testing.T) {
	t.Parallel()

	auth := &mockRequestAuthenticator{
		t: t,
		authenticateFn: func(context.Context, *validatorpb.StoreChunksRequest, time.Time) ([]byte, error) {
			return nil, nil
		},
	}

	interceptor := StoreChunksDisperserAuthAndRateLimitInterceptor(
		NewDisperserRateLimiter(nil, 10, 10),
		auth,
	)

	var gotID uint32
	var gotOk bool
	_, err := interceptor(
		context.Background(),
		&validatorpb.StoreChunksRequest{DisperserID: 11},
		&grpc.UnaryServerInfo{FullMethod: validatorpb.Dispersal_StoreChunks_FullMethodName},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			gotID, gotOk = AuthenticatedDisperserIDFromContext(ctx)
			return "ok", nil
		},
	)
	require.NoError(t, err)
	require.True(t, gotOk)
	require.Equal(t, uint32(11), gotID)
}
