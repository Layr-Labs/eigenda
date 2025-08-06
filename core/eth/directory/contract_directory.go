package directory

import (
	"context"
	"fmt"
	"sync"

	contractIEigenDADirectory "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDADirectory"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type ContractName string

// EigenDA uses many different contracts. It used to be the case that each contract address had to be provided via
// configuration, which was hard to maintain and error-prone. Now, contract addresses are registered onchain in the
// "EigenDA directory" contract. This struct is a convenience wrapper for interacting with the directory contract.
//
// Originally, the contract directory was just referred to as "the directory" or "the EigenDA directory". The term
// "directory" is extremely overloaded and is poorly descriptive, and the prefix "EigenDA" doesn't help since everything
// in this repo qualifies for that prefix. Unfortunately, the name of the contract is hard to change now. As a general
// rule of thumb, we should use "contract directory" when referring to this service, and "contract directory contract"
// when referring specifically to the solidity contract.
type ContractDirectory struct {
	logger logging.Logger

	// Only look up each address once. Most of our code only looks this stuff up at startup, so there isn't much
	// point in checking a particular contract address multiple times.
	addressCache map[ContractName]gethcommon.Address

	// a handle for calling the EigenDA directory contract.
	caller *contractIEigenDADirectory.ContractIEigenDADirectoryCaller

	// A set of all known contract addresses. Used to prevent magic strings from sneaking into the codebase.
	legalContractSet map[ContractName]struct{}

	// Used to make this utility thread safe.
	lock sync.Mutex
}

// Create a new ContractDirectory instance.
//
// The requireCompleteness flag can be used to enforce that the offchain contract list matches the onchain contract
// list. In general, this should be enabled in test code to raise the alarm when the offchain contract list drifts
// out of sync with the onchain contract list. In production code, this should be disabled just in case the offending
// contracts are not important for the code to function correctly.
func NewContractDirectory(
	ctx context.Context,
	logger logging.Logger,
	client bind.ContractBackend,
	contractDirectoryAddress gethcommon.Address,
) (*ContractDirectory, error) {

	caller, err := contractIEigenDADirectory.NewContractIEigenDADirectoryCaller(contractDirectoryAddress, client)
	if err != nil {
		return nil, fmt.Errorf("NewContractDirectory: %w", err)
	}

	legalAddressSet := make(map[ContractName]struct{})
	for _, contractName := range knownContracts {
		legalAddressSet[contractName] = struct{}{}
	}

	d := &ContractDirectory{
		logger:           logger,
		addressCache:     make(map[ContractName]gethcommon.Address),
		caller:           caller,
		legalContractSet: legalAddressSet,
	}

	err = d.verifyContractList(ctx)
	if err != nil {
		return nil, fmt.Errorf("verifyContractList: %w", err)
	}

	return d, nil
}

// GetContractAddress returns the address of a contract by its name. Only contracts defined in the const
// block above should be used. It is sharply discouraged to use this function with a magic string.
func (d *ContractDirectory) GetContractAddress(
	ctx context.Context,
	contractName ContractName,
) (gethcommon.Address, error) {
	if contractName == "" {
		return gethcommon.Address{}, fmt.Errorf("contract name cannot be empty")
	}

	// This is not very granular. But since this is uniquely to be a performance hotspot, we can do the simple thing.
	d.lock.Lock()
	defer d.lock.Unlock()

	address, ok := d.addressCache[contractName]
	if ok {
		return address, nil
	}

	// Before we look up the address, make sure it's in our list of known contracts.
	if _, exists := d.legalContractSet[contractName]; !exists {
		return gethcommon.Address{}, fmt.Errorf("contract %s is not a known contract", contractName)
	}

	address, err := d.caller.GetAddress0(&bind.CallOpts{Context: ctx}, (string)(contractName))
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("GetAddress0: %w", err)
	}

	if address == (gethcommon.Address{}) {
		return gethcommon.Address{}, fmt.Errorf("constract %s is not registered onchain", contractName)
	}

	d.addressCache[contractName] = address

	d.logger.Debugf("fetched address for contract %s: %s", contractName, address.Hex())
	return address, nil
}

// Checks to see if the list of contracts tracked by this ContractDirectory matches the contracts currently registered
// in the EigenDA directory contract. Creates some noisy logs if there are any discrepancies.
func (d *ContractDirectory) verifyContractList(ctx context.Context) error {
	registeredContracts, err := d.caller.GetAllNames(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("GetAllNames: %w", err)
	}

	complete := true

	registeredContractSet := make(map[string]struct{}, len(registeredContracts))
	for _, name := range registeredContracts {
		registeredContractSet[name] = struct{}{}
	}

	for _, contractName := range knownContracts {
		_, exists := registeredContractSet[string(contractName)]
		if !exists {
			d.logger.Errorf(
				"Contract %s is known to offchain code but not registered in the "+
					"onchain EigenDA contract directory", contractName)
			complete = false
		}
	}

	if complete {
		d.logger.Infof("Onchain contract list matches offchain contract list")
	} else {
		d.logger.Warnf("Onchain contract list does not match offchain contract list")
	}

	return nil
}
