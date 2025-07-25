package eth

import (
	"fmt"
	"slices"

	eigendadirectory "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDADirectory"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// EigenDADirectoryReader wraps the address directory contract and provides
// safe getters for contract addresses with zero address validation
type EigenDADirectoryReader struct {
	contract *eigendadirectory.ContractIEigenDADirectory
}

// NewEigenDADirectoryReader creates a new EigenDADirectoryReader
func NewEigenDADirectoryReader(eigendaDirectoryHexAddr string, client bind.ContractBackend) (*EigenDADirectoryReader, error) {
	if eigendaDirectoryHexAddr == "" || !gethcommon.IsHexAddress(eigendaDirectoryHexAddr) {
		return nil, fmt.Errorf("address directory must be a valid hex address: %s", eigendaDirectoryHexAddr)
	}

	eigendaDirectoryAddr := gethcommon.HexToAddress(eigendaDirectoryHexAddr)
	contract, err := eigendadirectory.NewContractIEigenDADirectory(eigendaDirectoryAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create EigenDADirectory contract: %w", err)
	}

	return &EigenDADirectoryReader{
		contract: contract,
	}, nil
}

// getAddressWithValidation reads the directory to get an address by the contract name
// and validates it's not zero
func (r *EigenDADirectoryReader) getAddressWithValidation(contractName string) (gethcommon.Address, error) {
	names, err := r.contract.GetAllNames(&bind.CallOpts{})
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("failed to get all contract names: %w", err)
	}
	if !slices.Contains(names, contractName) {
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

// GetOperatorStateRetrieverAddress returns the operator state retriever address with validation
func (r *EigenDADirectoryReader) GetOperatorStateRetrieverAddress() (gethcommon.Address, error) {
	return r.getAddressWithValidation(ContractNames.OperatorStateRetriever)
}

// GetServiceManagerAddress returns the service manager address with validation
func (r *EigenDADirectoryReader) GetServiceManagerAddress() (gethcommon.Address, error) {
	return r.getAddressWithValidation(ContractNames.ServiceManager)
}

// GetAllAddresses returns all contract addresses from the directory in a map keyed by contract name
func (r *EigenDADirectoryReader) GetAllAddresses() (map[string]gethcommon.Address, error) {
	names, err := r.contract.GetAllNames(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all contract names: %w", err)
	}

	addresses := make(map[string]gethcommon.Address)
	for _, name := range names {
		addr, err := r.contract.GetAddress0(&bind.CallOpts{}, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s address: %w", name, err)
		}
		addresses[name] = addr
	}

	return addresses, nil
}

// TODO: add other getters for other contracts; they are not needed for the current usage
