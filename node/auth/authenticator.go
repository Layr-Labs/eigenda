package auth

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

// RequestAuthenticator authenticates requests to the DA node. This object is thread safe.
//
// This class has largely been future-proofed for decentralized dispersers, with the exception of the
// preloadCache method, which will need to be updated to handle decentralized dispersers.
type RequestAuthenticator interface {
	// AuthenticateStoreChunksRequest authenticates a StoreChunksRequest, returning an error if the request is invalid.
	// Returns the hash of the request and an error if the request is invalid.
	AuthenticateStoreChunksRequest(
		ctx context.Context,
		request *grpc.StoreChunksRequest,
		now time.Time) ([]byte, error)
}

// keyWithTimeout is a key with that key's expiration time. After a key "expires", it should be reloaded
// from the chain state in case the key has been changed.
type keyWithTimeout struct {
	key        gethcommon.Address
	expiration time.Time
}

var _ RequestAuthenticator = &requestAuthenticator{}

type requestAuthenticator struct {
	// chainReader is used to read the chain state.
	chainReader core.Reader

	// keyCache is used to cache the public keys of dispersers. The uint32 map keys are disperser IDs. Disperser
	// IDs are serial numbers, with the original EigenDA disperser assigned ID 0. The map values contain
	// the public key of the disperser and the time when the local cache of the key will expire.
	keyCache *lru.Cache[uint32 /* disperser ID */, *keyWithTimeout]

	// keyTimeoutDuration is the duration for which a key is cached. After this duration, the key should be
	// reloaded from the chain state in case the key has been changed.
	keyTimeoutDuration time.Duration

	// disperserBlacklist is a list of disperser ID that has been blacklisted by the validator node
	disperserBlacklist []uint32
}

// NewRequestAuthenticator creates a new RequestAuthenticator.
func NewRequestAuthenticator(
	ctx context.Context,
	chainReader core.Reader,
	keyCacheSize int,
	keyTimeoutDuration time.Duration,
	disperserBlacklist []uint32,
	now time.Time) (RequestAuthenticator, error) {

	keyCache, err := lru.New[uint32, *keyWithTimeout](keyCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create key cache: %w", err)
	}

	authenticator := &requestAuthenticator{
		chainReader:        chainReader,
		keyCache:           keyCache,
		keyTimeoutDuration: keyTimeoutDuration,
		disperserBlacklist: disperserBlacklist,
	}

	err = authenticator.preloadCache(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to preload cache: %w", err)
	}

	return authenticator, nil
}

func (a *requestAuthenticator) preloadCache(ctx context.Context, now time.Time) error {
	// Populate the cache with the EigenLabs disperser ID; other disperser IDs will be added as they are encountered.
	_, err := a.getDisperserKey(ctx, now, api.EigenLabsDisperserID)
	if err != nil {
		return fmt.Errorf("failed to get operator key: %w", err)
	}

	return nil
}

func (a *requestAuthenticator) AuthenticateStoreChunksRequest(
	ctx context.Context,
	request *grpc.StoreChunksRequest,
	now time.Time) ([]byte, error) {

	key, err := a.getDisperserKey(ctx, now, request.DisperserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator key: %w", err)
	}

	hash, err := VerifyStoreChunksRequest(*key, request)
	if err != nil {
		return nil, fmt.Errorf("failed to verify request: %w", err)
	}

	return hash, nil
}

// getDisperserKey returns the public key of the operator with the given ID, caching the result.
func (a *requestAuthenticator) getDisperserKey(
	ctx context.Context,
	now time.Time,
	disperserID uint32) (*gethcommon.Address, error) {

	// Check if the disperser ID is blacklisted
	if slices.Contains(a.disperserBlacklist, disperserID) {
		return nil, fmt.Errorf("blacklisted disperser ID: %d", disperserID)
	}

	key, ok := a.keyCache.Get(disperserID)
	if ok {
		expirationTime := key.expiration
		if now.Before(expirationTime) {
			return &key.key, nil
		}
	}

	address, err := a.chainReader.GetDisperserAddress(ctx, disperserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get disperser address: %w", err)
		// if the IP address of the request consistently sends key that are not found, add it to blacklist
	}

	a.keyCache.Add(disperserID, &keyWithTimeout{
		key:        address,
		expiration: now.Add(a.keyTimeoutDuration),
	})

	return &address, nil
}

// // refreshCache refreshes all existing disperser keys in the cache.
// // It gathers all keys first, then makes a single multi-read call to the chain reader,
// // and finally repopulates the cache with the fresh values.
// func (a *requestAuthenticator) refreshCache(ctx context.Context, now time.Time) error {
// 	// Get all existing disperser IDs from the cache
// 	disperserIDs := make([]uint32, 0)
// 	for _, id := range a.keyCache.Keys() {
// 		disperserIDs = append(disperserIDs, id)
// 	}

// 	if len(disperserIDs) == 0 {
// 		return nil
// 	}

// 	// Get all keys in a single call
// 	keys, err := a.chainReader.GetDisperserAddresses(ctx, disperserIDs)
// 	if err != nil {
// 		return fmt.Errorf("failed to get disperser keys: %w", err)
// 	}

// 	// Repopulate the cache with fresh values
// 	for i, id := range disperserIDs {
// 		a.keyCache.Add(id, &keyWithTimeout{
// 			key:        keys[i],
// 			expiration: now.Add(a.keyTimeoutDuration),
// 		})
// 	}

// 	return nil
// }
