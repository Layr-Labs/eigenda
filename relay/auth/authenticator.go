package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/hashing"

	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/core"
	lru "github.com/hashicorp/golang-lru/v2"
)

// RequestAuthenticator authenticates requests to the relay service. This object is thread safe.
type RequestAuthenticator interface {
	// AuthenticateGetChunksRequest authenticates a GetChunksRequest, returning an error if the request is invalid.
	AuthenticateGetChunksRequest(
		ctx context.Context,
		request *pb.GetChunksRequest,
		now time.Time) error
}

var _ RequestAuthenticator = &requestAuthenticator{}

type requestAuthenticator struct {
	ics core.IndexedChainState

	// keyCache is used to cache the public keys of operators. Operator keys are assumed to never change.
	keyCache *lru.Cache[core.OperatorID, *core.G2Point]
}

// NewRequestAuthenticator creates a new RequestAuthenticator.
func NewRequestAuthenticator(
	ctx context.Context,
	ics core.IndexedChainState,
	keyCacheSize int) (RequestAuthenticator, error) {

	keyCache, err := lru.New[core.OperatorID, *core.G2Point](keyCacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create key cache: %w", err)
	}

	authenticator := &requestAuthenticator{
		ics:      ics,
		keyCache: keyCache,
	}

	err = authenticator.preloadCache(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to preload cache: %w", err)
	}

	return authenticator, nil
}

func (a *requestAuthenticator) preloadCache(ctx context.Context) error {
	blockNumber, err := a.ics.GetCurrentBlockNumber(ctx)
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
	request *pb.GetChunksRequest,
	now time.Time) error {

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

	hash, err := hashing.HashGetChunksRequest(request)
	if err != nil {
		return fmt.Errorf("failed to hash request: %w", err)
	}
	isValid := signature.Verify(key, ([32]byte)(hash))

	if !isValid {
		return errors.New("signature verification failed")
	}

	return nil
}

// getOperatorKey returns the public key of the operator with the given ID, caching the result.
func (a *requestAuthenticator) getOperatorKey(ctx context.Context, operatorID core.OperatorID) (*core.G2Point, error) {
	key, ok := a.keyCache.Get(operatorID)
	if ok {
		return key, nil
	}

	blockNumber, err := a.ics.GetCurrentBlockNumber(ctx)
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
