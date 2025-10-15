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

// A utility for converting back and forth between validator IDs and Ethereum addresses.
type ValidatorIDToAddressConverter interface {
	// Given a validator ID, find the validator's corresponding Ethereum address.
	ValidatorIDToAddress(ctx context.Context, validatorID core.OperatorID) (geth.Address, error)

	// Given a validator's Ethereum address, find the corresponding validator ID.
	ValidatorAddressToID(ctx context.Context, validatorAddress geth.Address) (core.OperatorID, error)
}

var _ ValidatorIDToAddressConverter = (*validatorIDToAddressConverter)(nil)

// A standard implementation of the ValidatorIDToAddressConverter interface.
type validatorIDToAddressConverter struct {
	registryCoordinator *regcoordinator.ContractEigenDARegistryCoordinator
}

func NewValidatorIDToAddressConverter(
	contractBackend bind.ContractBackend,
	registryCoordinatorAddress geth.Address,
) (ValidatorIDToAddressConverter, error) {

	registryCoordinator, err := regcoordinator.NewContractEigenDARegistryCoordinator(
		registryCoordinatorAddress,
		contractBackend)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry coordinator client: %w", err)
	}

	return &validatorIDToAddressConverter{
		registryCoordinator: registryCoordinator,
	}, nil

}

func (v *validatorIDToAddressConverter) ValidatorAddressToID(
	ctx context.Context,
	validatorAddress geth.Address,
) (core.OperatorID, error) {

	operatorInfo, err := v.registryCoordinator.GetOperator(&bind.CallOpts{Context: ctx}, validatorAddress)
	if err != nil {
		return core.OperatorID{}, fmt.Errorf("failed to get operator ID from address: %w", err)
	}
	validatorID := operatorInfo.OperatorId

	if validatorID == (core.OperatorID{}) {
		return core.OperatorID{}, fmt.Errorf("no operator found with address %s", validatorAddress.Hex())
	}

	return validatorID, nil

}

func (v *validatorIDToAddressConverter) ValidatorIDToAddress(
	ctx context.Context,
	validatorID core.OperatorID,
) (geth.Address, error) {

	address, err := v.registryCoordinator.GetOperatorFromId(&bind.CallOpts{Context: ctx}, validatorID)
	if err != nil {
		var zero geth.Address
		return zero, fmt.Errorf("failed to get operator address from ID: %w", err)
	}

	if address == (geth.Address{}) {
		return geth.Address{}, fmt.Errorf("no operator found with ID 0x%s", validatorID.Hex())
	}

	return address, nil
}

var _ ValidatorIDToAddressConverter = (*cachedValidatorIDToAddressConverter)(nil)

// A cached version of ValidatorIDToAddressConverter.
type cachedValidatorIDToAddressConverter struct {
	base             ValidatorIDToAddressConverter
	idToAddressCache *lru.Cache[core.OperatorID, geth.Address]
	addressToIDCache *lru.Cache[geth.Address, core.OperatorID]
}

func NewCachedValidatorIDToAddressConverter(
	base ValidatorIDToAddressConverter,
	cacheSize int,
) (ValidatorIDToAddressConverter, error) {

	idToAddressCache, err := lru.New[core.OperatorID, geth.Address](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create ID to address cache: %w", err)
	}

	addressToIDCache, err := lru.New[geth.Address, core.OperatorID](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create address to ID cache: %w", err)
	}

	return &cachedValidatorIDToAddressConverter{
		base:             base,
		idToAddressCache: idToAddressCache,
		addressToIDCache: addressToIDCache,
	}, nil
}

func (c *cachedValidatorIDToAddressConverter) ValidatorAddressToID(
	ctx context.Context,
	validatorAddress geth.Address,
) (core.OperatorID, error) {

	if id, ok := c.addressToIDCache.Get(validatorAddress); ok {
		return id, nil
	}

	id, err := c.base.ValidatorAddressToID(ctx, validatorAddress)
	if err != nil {
		return core.OperatorID{}, fmt.Errorf("failed to get validator ID from address: %w", err)
	}

	c.addressToIDCache.Add(validatorAddress, id)
	c.idToAddressCache.Add(id, validatorAddress)

	return id, nil
}

func (c *cachedValidatorIDToAddressConverter) ValidatorIDToAddress(
	ctx context.Context,
	validatorID core.OperatorID,
) (geth.Address, error) {

	if address, ok := c.idToAddressCache.Get(validatorID); ok {
		return address, nil
	}

	address, err := c.base.ValidatorIDToAddress(ctx, validatorID)
	if err != nil {
		return geth.Address{}, fmt.Errorf("failed to get validator address from ID: %w", err)
	}

	c.idToAddressCache.Add(validatorID, address)
	c.addressToIDCache.Add(address, validatorID)

	return address, nil
}
