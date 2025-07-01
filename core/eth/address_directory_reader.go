package eth

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	eigendadirectory "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDADirectory"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// AddressDirectoryReader wraps the address directory contract and provides
// safe getters for contract addresses with zero address validation
type AddressDirectoryReader struct {
	contract *eigendadirectory.ContractIEigenDADirectory
}

// NewAddressDirectoryReader creates a new AddressDirectoryReader
func NewAddressDirectoryReader(addressDirectoryHexAddr string, client common.EthClient) (*AddressDirectoryReader, error) {
	addressDirectoryAddr := gethcommon.HexToAddress(addressDirectoryHexAddr)
	contract, err := eigendadirectory.NewContractIEigenDADirectory(addressDirectoryAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create EigenDADirectory contract: %w", err)
	}

	return &AddressDirectoryReader{
		contract: contract,
	}, nil
}

// GetOperatorStateRetrieverAddress returns the operator state retriever address with validation
func (r *AddressDirectoryReader) GetOperatorStateRetrieverAddress() (gethcommon.Address, error) {
	addr, err := r.contract.GetAddress0(&bind.CallOpts{}, ContractNames.OperatorStateRetriever)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("failed to get operator state retriever address: %w", err)
	}
	if addr == (gethcommon.Address{}) {
		return gethcommon.Address{}, fmt.Errorf("operator state retriever address is zero - not deployed or registered in address directory")
	}
	return addr, nil
}

// GetServiceManagerAddress returns the service manager address with validation
func (r *AddressDirectoryReader) GetServiceManagerAddress() (gethcommon.Address, error) {
	addr, err := r.contract.GetAddress0(&bind.CallOpts{}, ContractNames.ServiceManager)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("failed to get service manager address: %w", err)
	}
	if addr == (gethcommon.Address{}) {
		return gethcommon.Address{}, fmt.Errorf("service manager address is zero - not deployed or registered in address directory")
	}
	return addr, nil
}

// TODO: add other getters for other contracts; they are not needed for the current usage
