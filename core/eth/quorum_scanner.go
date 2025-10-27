package eth

import (
	"context"
	"fmt"
	"math/big"

	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARegistryCoordinator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

// A utility that is capable of producing a list of all registered quorums.
type QuorumScanner interface {

	// Get all quorums registered at the given reference block number. Quorums are returned
	// sorted from least to greatest.
	GetQuorums(ctx context.Context, referenceBlockNumber uint64) ([]core.QuorumID, error)
}

var _ QuorumScanner = (*quorumScanner)(nil)

// A standard implementation of the QuorumScanner.
type quorumScanner struct {
	// A handle for communicating with the registry coordinator contract.
	registryCoordinator *regcoordinator.ContractEigenDARegistryCoordinator
}

// Create a new QuorumScanner instance. This instance is thread safe but not cached.
func NewQuorumScanner(
	contractBackend bind.ContractBackend,
	registryCoordinatorAddress gethcommon.Address,
) (QuorumScanner, error) {

	registryCoordinator, err := regcoordinator.NewContractEigenDARegistryCoordinator(
		registryCoordinatorAddress,
		contractBackend)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry coordinator client: %w", err)
	}

	return &quorumScanner{
		registryCoordinator: registryCoordinator,
	}, nil
}

func (q *quorumScanner) GetQuorums(ctx context.Context, referenceBlockNumber uint64) ([]core.QuorumID, error) {
	// Quorums are assigned starting at 0, and then sequentially without gaps. If we
	// know the number of quorums, we can generate a list of quorum IDs.

	quorumCount, err := q.registryCoordinator.QuorumCount(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: new(big.Int).SetUint64(referenceBlockNumber),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get quorum count: %w", err)
	}

	quorums := make([]core.QuorumID, quorumCount)
	for i := uint8(0); i < quorumCount; i++ {
		quorums[i] = i
	}

	return quorums, nil
}

var _ QuorumScanner = (*cachedQuorumScanner)(nil)

// A cached QuorumScanner implementation.
type cachedQuorumScanner struct {
	base  QuorumScanner
	cache *lru.Cache[uint64, []core.QuorumID]
}

// Create a new cached QuorumScanner that wraps the given base QuorumScanner. This implementation is thread safe.
func NewCachedQuorumScanner(base QuorumScanner, cacheSize int) (QuorumScanner, error) {
	cache, err := lru.New[uint64, []core.QuorumID](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}
	return &cachedQuorumScanner{
		base:  base,
		cache: cache,
	}, nil
}

func (c *cachedQuorumScanner) GetQuorums(ctx context.Context, referenceBlockNumber uint64) ([]core.QuorumID, error) {
	if quorums, ok := c.cache.Get(referenceBlockNumber); ok {
		return quorums, nil
	}

	quorums, err := c.base.GetQuorums(ctx, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get quorums: %w", err)
	}

	c.cache.Add(referenceBlockNumber, quorums)
	return quorums, nil
}

// Convert a list of quorums to a byte slice, where each byte is the ID of a quorum.
// This is the format expected by many smart contract functions.
func QuorumListToBytes(quorums []core.QuorumID) []byte {
	result := make([]byte, len(quorums))
	copy(result, quorums)
	return result
}
