package eth

import (
	"fmt"
	"slices"

	"github.com/Layr-Labs/eigenda/common"
	eigendadirectory "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDADirectory"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// AddressDirectoryReader wraps the address directory contract and provides
// safe getters for contract addresses with zero address validation
type AddressDirectoryReader struct {
	contract      *eigendadirectory.ContractIEigenDADirectory
	contractNames []string
}

// NewAddressDirectoryReader creates a new AddressDirectoryReader
func NewAddressDirectoryReader(addressDirectoryHexAddr string, client common.EthClient) (*AddressDirectoryReader, error) {
	addressDirectoryAddr := gethcommon.HexToAddress(addressDirectoryHexAddr)
	contract, err := eigendadirectory.NewContractIEigenDADirectory(addressDirectoryAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create EigenDADirectory contract: %w", err)
	}

	contractNames, err := contract.GetAllNames(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all contract names: %w, addressDirectoryHexAddr: %s", err, addressDirectoryHexAddr)
	}

	return &AddressDirectoryReader{
		contract:      contract,
		contractNames: contractNames,
	}, nil
}

// getAddressWithValidation reads the directory to get an address by the contract name
// and validates it's not zero
func (r *AddressDirectoryReader) getAddressWithValidation(contractName string) (gethcommon.Address, error) {
	if !slices.Contains(r.contractNames, contractName) {
		return gethcommon.Address{}, fmt.Errorf("contract %s not found in address directory", contractName)
	}

	addr, err := r.contract.GetAddress0(&bind.CallOpts{}, contractName)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("failed to get %s address: %w", contractName, err)
	}
	if addr == (gethcommon.Address{}) {
		return gethcommon.Address{}, fmt.Errorf("%s address is zero", contractName)
	}
	return addr, nil
}

// GetAllContractNames returns the names of all contracts in the address directory
func (r *AddressDirectoryReader) GetAllContractNames() ([]string, error) {
	names, err := r.contract.GetAllNames(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all contract names: %w", err)
	}
	r.contractNames = names
	return names, nil
}

// GetOperatorStateRetrieverAddress returns the operator state retriever address with validation
func (r *AddressDirectoryReader) GetOperatorStateRetrieverAddress() (gethcommon.Address, error) {
	return r.getAddressWithValidation(ContractNames.OperatorStateRetriever)
}

// GetServiceManagerAddress returns the service manager address with validation
func (r *AddressDirectoryReader) GetServiceManagerAddress() (gethcommon.Address, error) {
	return r.getAddressWithValidation(ContractNames.ServiceManager)
}

// TODO: add other getters for other contracts; they are not needed for the current usage
