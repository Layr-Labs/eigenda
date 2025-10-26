package auth

import (
	"context"
	"fmt"
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

// keysWithTimeout contains all keys for a disperser with their expiration time.
type keysWithTimeout struct {
	keys       []gethcommon.Address
	expiration time.Time
}

var _ RequestAuthenticator = &requestAuthenticator{}

type requestAuthenticator struct {
	// chainReader is used to read the chain state.
	chainReader core.Reader

	// keyCache is used to cache the public keys of dispersers. The uint32 map keys are disperser IDs. Disperser
	// IDs are serial numbers, with the original EigenDA disperser assigned ID 0. The map values contain
	// all public keys of the disperser and the time when the local cache of the keys will expire.
	keyCache *lru.Cache[uint32 /* disperser ID */, *keysWithTimeout]

	// keyTimeoutDuration is the duration for which a key is cached. After this duration, the key should be
	// reloaded from the chain state in case the key has been changed.
	keyTimeoutDuration time.Duration

	// keyLimit is the maximum number of keys to check per disperser.
	keyLimit int

	// disperserIDFilter is a function that returns true if the given disperser ID is valid.
	disperserIDFilter func(uint32) bool
}

// NewRequestAuthenticator creates a new RequestAuthenticator.
func NewRequestAuthenticator(
	ctx context.Context,
	chainReader core.Reader,
	keyCacheSize int,
	keyTimeoutDuration time.Duration,
	keyLimit int,
	disperserIDFilter func(uint32) bool,
	now time.Time) (RequestAuthenticator, error) {

	keyCache, err := lru.New[uint32, *keysWithTimeout](keyCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create key cache: %w", err)
	}

	authenticator := &requestAuthenticator{
		chainReader:        chainReader,
		keyCache:           keyCache,
		keyTimeoutDuration: keyTimeoutDuration,
		keyLimit:           keyLimit,
		disperserIDFilter:  disperserIDFilter,
	}

	err = authenticator.preloadCache(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to preload cache: %w", err)
	}

	return authenticator, nil
}

func (a *requestAuthenticator) preloadCache(ctx context.Context, now time.Time) error {
	// this will need to be updated for decentralized dispersers
	_, err := a.getDisperserKeys(ctx, now, api.EigenLabsDisperserID)
	if err != nil {
		return fmt.Errorf("failed to get disperser keys: %w", err)
	}

	return nil
}

func (a *requestAuthenticator) AuthenticateStoreChunksRequest(
	ctx context.Context,
	request *grpc.StoreChunksRequest,
	now time.Time) ([]byte, error) {

	keys, err := a.getDisperserKeys(ctx, now, request.GetDisperserID())
	if err != nil {
		return nil, fmt.Errorf("failed to get disperser keys: %w", err)
	}

	hash, err := VerifyStoreChunksRequestWithKeys(keys, request)
	if err != nil {
		return nil, fmt.Errorf("failed to verify request: %w", err)
	}

	return hash, nil
}

// getDisperserKeys returns all public keys of the disperser with the given ID, caching the result.
func (a *requestAuthenticator) getDisperserKeys(
	ctx context.Context,
	now time.Time,
	disperserID uint32) ([]gethcommon.Address, error) {

	if !a.disperserIDFilter(disperserID) {
		return nil, fmt.Errorf("invalid disperser ID: %d", disperserID)
	}

	keys, ok := a.keyCache.Get(disperserID)
	if ok {
		expirationTime := keys.expiration
		if now.Before(expirationTime) {
			return keys.keys, nil
		}
	}

	addresses, err := a.chainReader.GetAllDisperserAddresses(ctx, disperserID, uint32(a.keyLimit))
	if err != nil {
		return nil, fmt.Errorf("failed to get disperser addresses: %w", err)
	}

	a.keyCache.Add(disperserID, &keysWithTimeout{
		keys:       addresses,
		expiration: now.Add(a.keyTimeoutDuration),
	})

	return addresses, nil
}
