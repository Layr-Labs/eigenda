package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"sync"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/emirpasic/gods/queues"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
	lru "github.com/hashicorp/golang-lru/v2"
)

// RequestAuthenticator authenticates requests to the relay service. This object is thread safe.
type RequestAuthenticator interface {
	// AuthenticateGetChunksRequest authenticates a GetChunksRequest, returning an error if the request is invalid.
	// The origin is the address of the peer that sent the request. This may be used to cache auth results
	// in order to save server resources.
	AuthenticateGetChunksRequest(
		ctx context.Context,
		origin string,
		request *pb.GetChunksRequest,
		now time.Time) error
}

// authenticationTimeout is used to track the expiration of an auth.
type authenticationTimeout struct {
	origin     string
	expiration time.Time
}

var _ RequestAuthenticator = &requestAuthenticator{}

type requestAuthenticator struct {
	ics core.IndexedChainState

	// authenticatedClients is a set of client IDs that have been recently authenticated.
	authenticatedClients map[string]struct{}

	// authenticationTimeouts is a list of authentications that have been performed, along with their expiration times.
	authenticationTimeouts queues.Queue

	// authenticationTimeoutDuration is the duration for which an auth is valid.
	// If this is zero, then auth saving is disabled, and each request will be authenticated independently.
	authenticationTimeoutDuration time.Duration

	// savedAuthLock is used for thread safe atomic modification of the authenticatedClients map and the
	// authenticationTimeouts queue.
	savedAuthLock sync.Mutex

	// keyCache is used to cache the public keys of operators. Operator keys are assumed to never change.
	keyCache *lru.Cache[core.OperatorID, *core.G2Point]
}

// NewRequestAuthenticator creates a new RequestAuthenticator.
func NewRequestAuthenticator(
	ctx context.Context,
	ics core.IndexedChainState,
	keyCacheSize int,
	authenticationTimeoutDuration time.Duration) (RequestAuthenticator, error) {

	keyCache, err := lru.New[core.OperatorID, *core.G2Point](keyCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create key cache: %w", err)
	}

	authenticator := &requestAuthenticator{
		ics:                           ics,
		authenticatedClients:          make(map[string]struct{}),
		authenticationTimeouts:        linkedlistqueue.New(),
		authenticationTimeoutDuration: authenticationTimeoutDuration,
		keyCache:                      keyCache,
	}

	err = authenticator.preloadCache(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to preload cache: %w", err)
	}

	return authenticator, nil
}

func (a *requestAuthenticator) preloadCache(ctx context.Context) error {
	blockNumber, err := a.ics.GetCurrentBlockNumber()
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}
	operators, err := a.ics.GetIndexedOperators(ctx, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to get operators: %w", err)
	}

	for operatorID, operator := range operators {
		a.keyCache.Add(operatorID, operator.PubkeyG2)
	}

	return nil
}

func (a *requestAuthenticator) AuthenticateGetChunksRequest(
	ctx context.Context,
	origin string,
	request *pb.GetChunksRequest,
	now time.Time) error {

	if a.isAuthenticationStillValid(now, origin) {
		// We've recently authenticated this client. Do not authenticate again for a while.
		return nil
	}

	if request.OperatorId == nil || len(request.OperatorId) != 32 {
		return errors.New("invalid operator ID")
	}

	key, err := a.getOperatorKey(ctx, core.OperatorID(request.OperatorId))
	if err != nil {
		return fmt.Errorf("failed to get operator key: %w", err)
	}

	g1Point, err := (&core.G1Point{}).Deserialize(request.OperatorSignature)
	if err != nil {
		return fmt.Errorf("failed to deserialize signature: %w", err)
	}

	signature := core.Signature{
		G1Point: g1Point,
	}

	hash := hashing.HashGetChunksRequest(request)
	isValid := signature.Verify(key, ([32]byte)(hash))

	if !isValid {
		return errors.New("signature verification failed")
	}

	a.saveAuthenticationResult(now, origin)
	return nil
}

// getOperatorKey returns the public key of the operator with the given ID, caching the result.
func (a *requestAuthenticator) getOperatorKey(ctx context.Context, operatorID core.OperatorID) (*core.G2Point, error) {
	key, ok := a.keyCache.Get(operatorID)
	if ok {
		return key, nil
	}

	blockNumber, err := a.ics.GetCurrentBlockNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}
	operators, err := a.ics.GetIndexedOperators(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get operators: %w", err)
	}

	operator, ok := operators[operatorID]
	if !ok {
		return nil, errors.New("operator not found")
	}
	key = operator.PubkeyG2

	a.keyCache.Add(operatorID, key)
	return key, nil
}

// saveAuthenticationResult saves the result of an auth.
func (a *requestAuthenticator) saveAuthenticationResult(now time.Time, origin string) {
	if a.authenticationTimeoutDuration == 0 {
		// Authentication saving is disabled.
		return
	}

	a.savedAuthLock.Lock()
	defer a.savedAuthLock.Unlock()

	a.authenticatedClients[origin] = struct{}{}
	a.authenticationTimeouts.Enqueue(
		&authenticationTimeout{
			origin:     origin,
			expiration: now.Add(a.authenticationTimeoutDuration),
		})
}

// isAuthenticationStillValid returns true if the client at the given address has been authenticated recently.
func (a *requestAuthenticator) isAuthenticationStillValid(now time.Time, address string) bool {
	if a.authenticationTimeoutDuration == 0 {
		// Authentication saving is disabled.
		return false
	}

	a.savedAuthLock.Lock()
	defer a.savedAuthLock.Unlock()

	a.removeOldAuthentications(now)
	_, ok := a.authenticatedClients[address]
	return ok
}

// removeOldAuthentications removes any authentications that have expired.
// This method is not thread safe and should be called with the savedAuthLock held.
func (a *requestAuthenticator) removeOldAuthentications(now time.Time) {
	for a.authenticationTimeouts.Size() > 0 {
		val, _ := a.authenticationTimeouts.Peek()
		next := val.(*authenticationTimeout)
		if next.expiration.After(now) {
			break
		}
		delete(a.authenticatedClients, next.origin)
		a.authenticationTimeouts.Dequeue()
	}
}
