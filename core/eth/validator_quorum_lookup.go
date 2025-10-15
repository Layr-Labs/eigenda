package eth

import (
	"context"
	"math"
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

	return nil, nil
}

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

	// This method returns a bitmap as a big.Int, where the first 255 bits represent membership in quorums 0-254.
	bigIntBitmap, err := v.registryCoordinator.GetCurrentQuorumBitmap(opts, validatorID)
	if err != nil {
		return nil, err
	}
	bitmap := bigIntBitmap.Bytes()

	quorumIDs := make([]core.QuorumID, 0)

	// Although technically 254 is the max quorum ID (due to an embarrassing off-by-one typo), it doesn't hurt
	// to check bit 255. It's possible that the typo will be fixed in the future, and this
	// code should still work if that happens.
	for i := 0; i <= math.MaxUint8; i++ {
		byteIndex := i / 8
		bitIndex := i % 8

		bit := (bitmap[byteIndex] >> bitIndex) & 1
		if bit == 1 {
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
	cache *lru.Cache[core.OperatorID, []core.QuorumID]
}

// Create a new cached ValidatorQuorumLookup with the given cache size.
func NewCachedValidatorQuorumLookup(
	base ValidatorQuorumLookup,
	cacheSize int,
) (ValidatorQuorumLookup, error) {

	cache, err := lru.New[core.OperatorID, []core.QuorumID](cacheSize)
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

	if quorums, ok := c.cache.Get(validatorAddress); ok {
		return quorums, nil
	}

	quorums, err := c.base.GetQuorumsForValidator(ctx, validatorAddress, referenceBlockNumber)
	if err != nil {
		return nil, err
	}

	c.cache.Add(validatorAddress, quorums)

	return quorums, nil
}
