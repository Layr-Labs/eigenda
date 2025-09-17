package eth

import (
	"context"
	"fmt"

	regcoordinator "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARegistryCoordinator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	geth "github.com/ethereum/go-ethereum/common"
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
		return geth.Address{}, fmt.Errorf("no operator found with ID %d", validatorID)
	}

	return address, nil
}
