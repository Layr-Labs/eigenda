package eth

import (
	"context"
	"fmt"

	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARegistryCoordinator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	geth "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

// Given a validator ID, find the validator's corresponding Ethereum address.
func ValidatorIDToAddress(
	ctx context.Context,
	contractBackend bind.ContractBackend,
	registryCoordinatorAddress geth.Address,
	validatorID core.OperatorID,
) (geth.Address, error) {

	registryCoordinator, err := regcoordinator.NewContractEigenDARegistryCoordinator(
		registryCoordinatorAddress,
		contractBackend)
	if err != nil {
		var zero geth.Address
		return zero, fmt.Errorf("failed to create registry coordinator client: %w", err)
	}

	address, err := registryCoordinator.GetOperatorFromId(&bind.CallOpts{Context: ctx}, validatorID)
	if err != nil {
		var zero geth.Address
		return zero, fmt.Errorf("failed to get operator address from ID: %w", err)
	}

	if address == (geth.Address{}) {
		return geth.Address{}, fmt.Errorf("no operator found with ID 0x%s", validatorID.Hex())
	}

	return address, nil
}

// A cache for validator ID to address mappings. Thread safe, but concurrent calls for the same validator ID
// may result in multiple onchain calls.
type ValidatorIDToAddressCache struct {
	// The contract backend used for making calls to the blockchain.
	contractBackend bind.ContractBackend

	// The address of the RegistryCoordinator contract.
	registryCoordinatorAddress geth.Address

	// A cache of previously looked up validator ID to address mappings.
	cache *lru.Cache[core.OperatorID, geth.Address]
}

// NewValidatorIDToAddressCache creates a new ValidatorIDToAddressCache with the given cache size.
func NewValidatorIDToAddressCache(
	contractBackend bind.ContractBackend,
	registryCoordinatorAddress geth.Address,
	cacheSize int,
) (*ValidatorIDToAddressCache, error) {

	cache, err := lru.New[core.OperatorID, geth.Address](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator ID to address cache: %w", err)
	}

	return &ValidatorIDToAddressCache{
		contractBackend:            contractBackend,
		registryCoordinatorAddress: registryCoordinatorAddress,
		cache:                      cache,
	}, nil
}

// GetValidatorAddress looks up the Ethereum address for the given validator ID, using the cache if possible.
func (c *ValidatorIDToAddressCache) GetValidatorAddress(
	ctx context.Context,
	validatorID core.OperatorID,
) (geth.Address, error) {

	if address, ok := c.cache.Get(validatorID); ok {
		return address, nil
	}

	address, err := ValidatorIDToAddress(ctx, c.contractBackend, c.registryCoordinatorAddress, validatorID)
	if err != nil {
		return geth.Address{}, err
	}

	c.cache.Add(validatorID, address)

	return address, nil
}
