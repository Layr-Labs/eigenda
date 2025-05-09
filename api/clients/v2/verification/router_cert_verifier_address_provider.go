package verification

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// RouterAddressProvider is a dynamic provider which fetches cert verifier addresses by making eth_calls
// against the EigenDACertVerifierRouter contract on the given reference block number.
type RouterAddressProvider struct {
	routerBinding *binding.ContractEigenDACertVerifierRouterCaller
}

// Ensure RouterAddressProvider implements clients.CertVerifierAddressProvider
var _ clients.CertVerifierAddressProvider = &RouterAddressProvider{}

// BuildRouterAddressProvider creates a new RouterAddressProvider instance
// that implements the clients.CertVerifierAddressProvider interface
func BuildRouterAddressProvider(routerAddr gethcommon.Address, ethClient common.EthClient) (*RouterAddressProvider, error) {
	routerBinding, err := binding.NewContractEigenDACertVerifierRouterCaller(routerAddr, ethClient)
	if err != nil {
		return nil, err
	}

	return &RouterAddressProvider{
		routerBinding: routerBinding,
	}, nil
}

// GetCertVerifierAddress returns the cert verifier address for the given reference block number
func (rap *RouterAddressProvider) GetCertVerifierAddress(ctx context.Context, referenceBlockNumber uint64) (gethcommon.Address, error) {
	return rap.routerBinding.GetCertVerifierAt(&bind.CallOpts{Context: ctx}, uint32(referenceBlockNumber))
}
