package middleware

import (
	"context"
	"fmt"
	"time"

	validatorpb "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/node/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ctxKey int

const ctxKeyAuthenticatedDisperserID ctxKey = iota

func authenticatedDisperserIDFromContext(ctx context.Context) (uint32, bool) {
	v := ctx.Value(ctxKeyAuthenticatedDisperserID)
	if v == nil {
		return 0, false
	}
	id, ok := v.(uint32)
	return id, ok
}

// AuthenticatedDisperserIDFromContext returns the authenticated disperser ID (if present).
//
// This is set by StoreChunksDisperserAuthAndBlacklistInterceptor().
func AuthenticatedDisperserIDFromContext(ctx context.Context) (uint32, bool) {
	return authenticatedDisperserIDFromContext(ctx)
}

// StoreChunksDisperserAuthAndRateLimitInterceptor authenticates StoreChunks requests and rejects any requests from
// rate-limited dispersers before entering the handler.
//
// IMPORTANT: rate limiting is only enforced after request authentication. This prevents an attacker from spoofing
// a disperser ID and causing an honest disperser to be rate limited.
func StoreChunksDisperserAuthAndRateLimitInterceptor(
	rateLimiter *DisperserRateLimiter,
	requestAuthenticator auth.RequestAuthenticator,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info == nil || info.FullMethod != validatorpb.Dispersal_StoreChunks_FullMethodName {
			return handler(ctx, req)
		}

		storeReq, ok := req.(*validatorpb.StoreChunksRequest)
		if !ok {
			return nil, status.Errorf(codes.Internal, "unexpected request type for %s: %T", info.FullMethod, req)
		}

		now := time.Now()
		_, err := requestAuthenticator.AuthenticateStoreChunksRequest(ctx, storeReq, now)
		if err != nil {
			// Do NOT rate limit here; the disperser identity is not proven if auth fails.
			return nil, status.Errorf(codes.InvalidArgument, "failed to authenticate request: %v", err)
		}

		disperserID := storeReq.GetDisperserID()
		if rateLimiter != nil && !rateLimiter.Allow(disperserID, now) {
			return nil, status.Error(codes.ResourceExhausted, fmt.Sprintf("disperser %d is rate limited", disperserID))
		}

		ctx = context.WithValue(ctx, ctxKeyAuthenticatedDisperserID, disperserID)

		res, handlerErr := handler(ctx, req)

		return res, handlerErr
	}
}
