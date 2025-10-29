package auth

import (
	"context"
	"fmt"
	"time"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
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

// keyWithTimeout contains a single key with its expiration time.
type keyWithTimeout struct {
	key        gethcommon.Address
	expiration time.Time
}

var _ RequestAuthenticator = &requestAuthenticator{}

type requestAuthenticator struct {
	// chainReader is used to read the chain state.
	chainReader core.Reader

	// keyCache is used to cache the public keys from the disperser registry. The uint32 map keys are the
	// registry indices (0, 1, 2, ...) and the values contain the address at that index with expiration time.
	keyCache *lru.Cache[uint32 /* registry index */, *keyWithTimeout]

	// keyTimeoutDuration is the duration for which a key is cached. After this duration, the key should be
	// reloaded from the chain state in case the key has been changed.
	keyTimeoutDuration time.Duration

	// keyCacheCapacity is the maximum number of keys to scan from the registry
	keyCacheCapacity int

	// logger for debug output
	logger logging.Logger
}

// NewRequestAuthenticator creates a new RequestAuthenticator.
func NewRequestAuthenticator(
	ctx context.Context,
	chainReader core.Reader,
	keyCacheSize int,
	keyTimeoutDuration time.Duration,
	logger logging.Logger,
	now time.Time) (RequestAuthenticator, error) {

	keyCache, err := lru.New[uint32, *keyWithTimeout](keyCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create key cache: %w", err)
	}

	authenticator := &requestAuthenticator{
		chainReader:        chainReader,
		keyCache:           keyCache,
		keyTimeoutDuration: keyTimeoutDuration,
		keyCacheCapacity:   keyCacheSize,
		logger:             logger,
	}

	err = authenticator.preloadCache(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to preload cache: %w", err)
	}

	return authenticator, nil
}

func (a *requestAuthenticator) preloadCache(ctx context.Context, now time.Time) error {
	_, err := a.getAllDisperserKeys(ctx, now)
	if err != nil {
		return fmt.Errorf("failed to get disperser keys: %w", err)
	}

	return nil
}

func (a *requestAuthenticator) AuthenticateStoreChunksRequest(
	ctx context.Context,
	request *grpc.StoreChunksRequest,
	now time.Time) ([]byte, error) {

	keys, err := a.getAllDisperserKeys(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get disperser keys: %w", err)
	}

	hash, err := VerifyStoreChunksRequestWithKeys(keys, request)
	if err != nil {
		return nil, fmt.Errorf("failed to verify request: %w", err)
	}

	return hash, nil
}

// getAllDisperserKeys returns all public keys from the disperser registry, caching each address individually.
func (a *requestAuthenticator) getAllDisperserKeys(
	ctx context.Context,
	now time.Time) ([]gethcommon.Address, error) {

	maxKeys := uint32(a.keyCacheCapacity)

	// Collect all valid cached addresses
	var cachedAddresses []gethcommon.Address
	needRefresh := false

	// Check each index in the cache
	for i := uint32(0); i < maxKeys; i++ {
		key, ok := a.keyCache.Get(i)
		if !ok || now.After(key.expiration) {
			needRefresh = true
			break
		}
		cachedAddresses = append(cachedAddresses, key.key)
	}

	if !needRefresh && len(cachedAddresses) > 0 {
		return cachedAddresses, nil
	}

	// Fetch fresh data from chain
	addresses, err := a.chainReader.GetAllDisperserAddresses(ctx, maxKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to get disperser addresses: %w", err)
	}

	// Cache each address individually by its index (zero addresses are not stored)
	for i, addr := range addresses {
		a.logger.Debug("Adding disperser key",
			"index", i,
			"address", addr.Hex())

		a.keyCache.Add(uint32(i), &keyWithTimeout{
			key:        addr,
			expiration: now.Add(a.keyTimeoutDuration),
		})
	}

	return addresses, nil
}
