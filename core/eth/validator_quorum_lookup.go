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

// A utility for looking up which quorums a given validator is a member of at a specific reference block number.
type ValidatorQuorumLookup interface {
	// Get the list of quorums that the given validator is a member of, at the specified reference block number.
	GetQuorumsForValidator(
		ctx context.Context,
		validatorAddress core.OperatorID,
		referenceBlockNumber uint64) ([]core.QuorumID, error)
}

var _ ValidatorQuorumLookup = (*validatorQuorumLookup)(nil)

// A standard implementation of the ValidatorQuorumLookup interface.
type validatorQuorumLookup struct {
	registryCoordinator *regcoordinator.ContractEigenDARegistryCoordinator
}

// Create a new ValidatorQuorumLookup instance.
func NewValidatorQuorumLookup(
	backend bind.ContractBackend,
	registryCoordinatorAddress gethcommon.Address,
) (ValidatorQuorumLookup, error) {

	registryCoordinator, err := regcoordinator.NewContractEigenDARegistryCoordinator(
		registryCoordinatorAddress,
		backend,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry coordinator contract instance: %w", err)
	}

	return &validatorQuorumLookup{
		registryCoordinator: registryCoordinator,
	}, nil
}

// TODO test this by hand to verify behavior.

func (v *validatorQuorumLookup) GetQuorumsForValidator(
	ctx context.Context,
	validatorID core.OperatorID,
	referenceBlockNumber uint64,
) ([]core.QuorumID, error) {

	blockNumber := big.NewInt(int64(referenceBlockNumber))

	opts := &bind.CallOpts{
		Context:     ctx,
		BlockNumber: blockNumber,
	}

	// This method returns a bitmap as a big.Int.
	bigIntBitmap, err := v.registryCoordinator.GetCurrentQuorumBitmap(opts, validatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quorum bitmap: %w", err)
	}

	quorumIDs := make([]core.QuorumID, 0)

	// An implementation detail of the solidity: the number returned by the contract is a bitmap backed by a
	// uint192, so we need to check each bit up to 192. If we check for higher bits, we will panic.
	for i := 0; i <= 192; i++ {
		present := bigIntBitmap.Bit(i)
		if present == 1 {
			quorumID := core.QuorumID(i)
			quorumIDs = append(quorumIDs, quorumID)
		}
	}

	return quorumIDs, nil
}

var _ ValidatorQuorumLookup = (*cachedValidatorQuorumLookup)(nil)

// A cached implementation of a ValidatorQuorumLookup.
type cachedValidatorQuorumLookup struct {
	base  ValidatorQuorumLookup
	cache *lru.Cache[validatorQuorumCacheKey, []core.QuorumID]
}

type validatorQuorumCacheKey struct {
	validatorID          core.OperatorID
	referenceBlockNumber uint64
}

// Create a new cached ValidatorQuorumLookup with the given cache size.
func NewCachedValidatorQuorumLookup(
	base ValidatorQuorumLookup,
	cacheSize int,
) (ValidatorQuorumLookup, error) {

	cache, err := lru.New[validatorQuorumCacheKey, []core.QuorumID](cacheSize)
	if err != nil {
		return nil, err
	}

	return &cachedValidatorQuorumLookup{
		base:  base,
		cache: cache,
	}, nil
}

// GetQuorumsForValidator implements ValidatorQuorumLookup.
func (c *cachedValidatorQuorumLookup) GetQuorumsForValidator(
	ctx context.Context,
	validatorAddress core.OperatorID,
	referenceBlockNumber uint64,
) ([]core.QuorumID, error) {

	key := validatorQuorumCacheKey{
		validatorID:          validatorAddress,
		referenceBlockNumber: referenceBlockNumber,
	}

	if quorums, ok := c.cache.Get(key); ok {
		return quorums, nil
	}

	quorums, err := c.base.GetQuorumsForValidator(ctx, validatorAddress, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get quorums for validator: %w", err)
	}

	c.cache.Add(key, quorums)

	return quorums, nil
}
