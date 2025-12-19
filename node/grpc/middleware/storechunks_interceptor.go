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

// StoreChunksDisperserAuthAndBlacklistInterceptor authenticates StoreChunks requests and rejects any requests from
// blacklisted dispersers before entering the handler.
//
// IMPORTANT: blacklisting is only enforced after request authentication. This prevents an attacker from spoofing
// a disperser ID and causing an honest disperser to be blacklisted.
func StoreChunksDisperserAuthAndBlacklistInterceptor(
	blacklist *DisperserBlacklist,
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
			// Do NOT blacklist here; the disperser identity is not proven if auth fails.
			return nil, status.Errorf(codes.InvalidArgument, "failed to authenticate request: %v", err)
		}

		disperserID := storeReq.GetDisperserID()
		if blacklist != nil && blacklist.IsBlacklisted(disperserID, now) {
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("disperser %d is temporarily blacklisted", disperserID))
		}

		ctx = context.WithValue(ctx, ctxKeyAuthenticatedDisperserID, disperserID)
		return handler(ctx, req)
	}
}
