package auth

import (
	"context"
	"fmt"
	"time"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
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

	// IsDisperserAuthorized returns true if the disperser is authorized to disperse the given batch.
	// Returns true if the batch contains only reservation payments, or if the batch contains on-demand payments
	// and the disperser is authorized to handle them. Returns false if the batch contains on-demand payments
	// and the disperser is not authorized.
	IsDisperserAuthorized(disperserID uint32, batch *corev2.Batch) bool
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

	// keyCacheCapacity is the maximum number of keys that can be cached.
	keyCacheCapacity int

	// keyTimeoutDuration is the duration for which a key is cached. After this duration, the key should be
	// reloaded from the chain state in case the key has been changed.
	keyTimeoutDuration time.Duration

	// Set of disperser IDs authorized to submit on-demand payments.
	authorizedOnDemandDispersers map[uint32]struct{}
}

// NewRequestAuthenticator creates a new RequestAuthenticator.
func NewRequestAuthenticator(
	ctx context.Context,
	chainReader core.Reader,
	keyCacheSize int,
	keyTimeoutDuration time.Duration,
	authorizedOnDemandDispersers []uint32,
	now time.Time,
) (RequestAuthenticator, error) {

	keyCache, err := lru.New[uint32, *keyWithTimeout](keyCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create key cache: %w", err)
	}

	authorizedSet := make(map[uint32]struct{}, len(authorizedOnDemandDispersers))
	for _, id := range authorizedOnDemandDispersers {
		authorizedSet[id] = struct{}{}
	}

	authenticator := &requestAuthenticator{
		chainReader:                  chainReader,
		keyCache:                     keyCache,
		keyCacheCapacity:             keyCacheSize,
		keyTimeoutDuration:           keyTimeoutDuration,
		authorizedOnDemandDispersers: authorizedSet,
	}

	err = authenticator.preloadCache(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("failed to preload cache: %w", err)
	}

	return authenticator, nil
}

func (a *requestAuthenticator) preloadCache(ctx context.Context, now time.Time) error {
	// Preload disperser keys starting from ID 0 until we hit cache limit or resolve a default address 0x0
	for disperserID := uint32(0); disperserID < uint32(a.keyCacheCapacity); disperserID++ {
		address, err := a.chainReader.GetDisperserAddress(ctx, disperserID)
		if err != nil {
			fmt.Printf("failed to preload disperser key for ID %d: %v\n", disperserID, err)
			continue
		}

		// If we get a zero address (0x0), stop preloading as this indicates no more valid dispersers
		if address == (gethcommon.Address{}) {
			break
		}

		// Cache the key with timeout
		a.keyCache.Add(disperserID, &keyWithTimeout{
			key:        address,
			expiration: now.Add(a.keyTimeoutDuration),
		})
	}

	return nil
}

func (a *requestAuthenticator) AuthenticateStoreChunksRequest(
	ctx context.Context,
	request *grpc.StoreChunksRequest,
	now time.Time) ([]byte, error) {

	key, err := a.getDisperserKey(ctx, now, request.GetDisperserID())
	if err != nil {
		return nil, fmt.Errorf("failed to get operator key: %w", err)
	}

	hash, err := VerifyStoreChunksRequest(*key, request)
	if err != nil {
		return nil, fmt.Errorf("failed to verify request: %w", err)
	}

	return hash, nil
}

func (a *requestAuthenticator) IsDisperserAuthorized(disperserID uint32, batch *corev2.Batch) bool {
	hasOnDemand := false
	for _, cert := range batch.BlobCertificates {
		if cert.BlobHeader.PaymentMetadata.IsOnDemand() {
			hasOnDemand = true
			break
		}
	}

	if !hasOnDemand {
		return true
	}

	_, authorized := a.authorizedOnDemandDispersers[disperserID]
	return authorized
}

// getDisperserKey returns the public key of the operator with the given ID, caching the result.
func (a *requestAuthenticator) getDisperserKey(
	ctx context.Context,
	now time.Time,
	disperserID uint32) (*gethcommon.Address, error) {

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
	}

	a.keyCache.Add(disperserID, &keyWithTimeout{
		key:        address,
		expiration: now.Add(a.keyTimeoutDuration),
	})

	return &address, nil
}
