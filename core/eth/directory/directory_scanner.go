package directory

import (
	"context"
	"fmt"

	contractIEigenDADirectory "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDADirectory"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// A utility method that looks up all contracts in the EigenDA directory contract and returns a map from name
// to address.
func GetContractAddressMap(
	ctx context.Context,
	client bind.ContractBackend,
	directoryAddress gethcommon.Address) (map[string]gethcommon.Address, error) {

	caller, err := contractIEigenDADirectory.NewContractIEigenDADirectoryCaller(directoryAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create EigenDA directory contract caller: %w", err)
	}

	names, err := caller.GetAllNames(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("eth-call:get all contract names: %w", err)
	}

	addresses := make(map[string]gethcommon.Address)
	for _, name := range names {
		addr, err := caller.GetAddress0(&bind.CallOpts{Context: ctx}, name)
		if err != nil {
			return nil, fmt.Errorf("eth-call: get %s address: %w", name, err)
		}
		addresses[name] = addr
	}

	return addresses, nil
}
