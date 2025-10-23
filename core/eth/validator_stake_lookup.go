package eth

import (
	"context"
	"fmt"
	"math/big"

	contractStakeRegistry "github.com/Layr-Labs/eigenda/contracts/bindings/StakeRegistry"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

// A utility for looking up a validator's stake.
type ValidatorStakeLookup interface {

	// Get a validator's stake in a specific quorum at a specific reference block number.
	GetValidatorStake(
		ctx context.Context,
		quorumID core.QuorumID,
		validatorID core.OperatorID,
		referenceBlockNumber uint64,
	) (*big.Int, error)

	// Get the total stake of all validators in a specific quorum at a specific reference block number.
	GetTotalQuorumStake(
		ctx context.Context,
		quorumID core.QuorumID,
		referenceBlockNumber uint64,
	) (*big.Int, error)

	// Get a validator's stake fraction (i.e., their stake divided by the total stake) in a specific quorum.
	// Returns a number between 0.0 and 1.0.
	GetValidatorStakeFraction(
		ctx context.Context,
		quorumID core.QuorumID,
		validatorID core.OperatorID,
		referenceBlockNumber uint64,
	) (float64, error)
}

var _ ValidatorStakeLookup = (*validatorStakeLookup)(nil)

// A standard implementation of the ValidatorStakeLookup interface.
type validatorStakeLookup struct {
	stakeRegistry *contractStakeRegistry.ContractStakeRegistry
}

// Create a new ValidatorStakeLookup instance.
func NewValidatorStakeLookup(
	backend bind.ContractBackend,
	stakeRegistryAddress gethcommon.Address,
) (ValidatorStakeLookup, error) {

	stakeRegistry, err := contractStakeRegistry.NewContractStakeRegistry(
		stakeRegistryAddress,
		backend,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stake registry contract instance: %w", err)
	}

	return &validatorStakeLookup{
		stakeRegistry: stakeRegistry,
	}, nil
}

func (v *validatorStakeLookup) GetTotalQuorumStake(
	ctx context.Context,
	quorumID core.QuorumID,
	referenceBlockNumber uint64,
) (*big.Int, error) {

	opts := &bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(referenceBlockNumber)),
	}

	stake, err := v.stakeRegistry.GetCurrentTotalStake(opts, quorumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total quorum stake: %w", err)
	}
	return stake, nil
}

func (v *validatorStakeLookup) GetValidatorStake(
	ctx context.Context,
	quorumID core.QuorumID,
	validatorID core.OperatorID,
	referenceBlockNumber uint64,
) (*big.Int, error) {

	opts := &bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(referenceBlockNumber)),
	}

	stake, err := v.stakeRegistry.GetCurrentStake(opts, validatorID, quorumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get validator stake: %w", err)
	}
	return stake, nil
}

func (v *validatorStakeLookup) GetValidatorStakeFraction(
	ctx context.Context,
	quorumID core.QuorumID,
	validatorID core.OperatorID,
	referenceBlockNumber uint64,
) (float64, error) {

	validatorStake, err := v.GetValidatorStake(ctx, quorumID, validatorID, referenceBlockNumber)
	if err != nil {
		return 0.0, fmt.Errorf("failed to get validator stake: %w", err)
	}

	totalStake, err := v.GetTotalQuorumStake(ctx, quorumID, referenceBlockNumber)
	if err != nil {
		return 0.0, fmt.Errorf("failed to get total quorum stake: %w", err)
	}

	if totalStake.Cmp(big.NewInt(0)) == 0 {
		return 0.0, nil // Avoid division by zero; if total stake is zero, return 0.0 fraction.
	}

	fraction := new(big.Rat).SetFrac(validatorStake, totalStake)
	floatFraction, _ := fraction.Float64()

	return floatFraction, nil
}

var _ ValidatorStakeLookup = (*cachedValidatorStakeLookup)(nil)

// A cached implementation of the ValidatorStakeLookup interface.
type cachedValidatorStakeLookup struct {
	base                ValidatorStakeLookup
	totalStakeCache     *lru.Cache[validatorStakeLookupTotalStakeCacheKey, *big.Int]
	validatorStakeCache *lru.Cache[validatorStakeLookupValidatorStakeCacheKey, *big.Int]
}

func NewCachedValidatorStakeLookup(
	base ValidatorStakeLookup,
	cacheSize int,
) (ValidatorStakeLookup, error) {

	totalStakeCache, err := lru.New[validatorStakeLookupTotalStakeCacheKey, *big.Int](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create total stake cache: %w", err)
	}

	validatorStakeCache, err := lru.New[validatorStakeLookupValidatorStakeCacheKey, *big.Int](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator stake cache: %w", err)
	}

	return &cachedValidatorStakeLookup{
		base:                base,
		totalStakeCache:     totalStakeCache,
		validatorStakeCache: validatorStakeCache,
	}, nil
}

type validatorStakeLookupTotalStakeCacheKey struct {
	quorumID             core.QuorumID
	referenceBlockNumber uint64
}

type validatorStakeLookupValidatorStakeCacheKey struct {
	quorumID             core.QuorumID
	validatorID          core.OperatorID
	referenceBlockNumber uint64
}

func (c *cachedValidatorStakeLookup) GetTotalQuorumStake(
	ctx context.Context,
	quorumID core.QuorumID,
	referenceBlockNumber uint64,
) (*big.Int, error) {

	key := validatorStakeLookupTotalStakeCacheKey{
		quorumID:             quorumID,
		referenceBlockNumber: referenceBlockNumber,
	}

	if stake, ok := c.totalStakeCache.Get(key); ok {
		return stake, nil
	}

	stake, err := c.base.GetTotalQuorumStake(ctx, quorumID, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get total quorum stake: %w", err)
	}

	c.totalStakeCache.Add(key, stake)

	return stake, nil
}

func (c *cachedValidatorStakeLookup) GetValidatorStake(
	ctx context.Context,
	quorumID core.QuorumID,
	validatorID core.OperatorID,
	referenceBlockNumber uint64,
) (*big.Int, error) {
	key := validatorStakeLookupValidatorStakeCacheKey{
		quorumID:             quorumID,
		validatorID:          validatorID,
		referenceBlockNumber: referenceBlockNumber,
	}

	if stake, ok := c.validatorStakeCache.Get(key); ok {
		return stake, nil
	}

	stake, err := c.base.GetValidatorStake(ctx, quorumID, validatorID, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get validator stake: %w", err)
	}

	c.validatorStakeCache.Add(key, stake)

	return stake, nil
}

func (c *cachedValidatorStakeLookup) GetValidatorStakeFraction(
	ctx context.Context,
	quorumID core.QuorumID,
	validatorID core.OperatorID,
	referenceBlockNumber uint64,
) (float64, error) {
	validatorStake, err := c.GetValidatorStake(ctx, quorumID, validatorID, referenceBlockNumber)
	if err != nil {
		return 0.0, fmt.Errorf("failed to get validator stake: %w", err)
	}

	totalStake, err := c.GetTotalQuorumStake(ctx, quorumID, referenceBlockNumber)
	if err != nil {
		return 0.0, fmt.Errorf("failed to get total quorum stake: %w", err)
	}

	if totalStake.Cmp(big.NewInt(0)) == 0 {
		return 0.0, nil // Avoid division by zero; if total stake is zero, return 0.0 fraction.
	}

	fraction := new(big.Rat).SetFrac(validatorStake, totalStake)
	floatFraction, _ := fraction.Float64()

	return floatFraction, nil
}
