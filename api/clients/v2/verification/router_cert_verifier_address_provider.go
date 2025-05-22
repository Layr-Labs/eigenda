package verification

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// RouterAddressProvider is a dynamic provider which fetches cert verifier addresses by making eth_calls
// against the EigenDACertVerifierRouter contract at the given reference block number.
type RouterAddressProvider struct {
	routerBinding      *binding.ContractEigenDACertVerifierRouterCaller
	blockNumberMonitor *BlockNumberMonitor
}

// Ensure RouterAddressProvider implements clients.CertVerifierAddressProvider
var _ clients.CertVerifierAddressProvider = &RouterAddressProvider{}

// BuildRouterAddressProvider creates a new RouterAddressProvider instance
// that implements the clients.CertVerifierAddressProvider interface
func BuildRouterAddressProvider(routerAddr gethcommon.Address, ethClient common.EthClient, logger logging.Logger) (*RouterAddressProvider, error) {
	routerBinding, err := binding.NewContractEigenDACertVerifierRouterCaller(routerAddr, ethClient)
	if err != nil {
		return nil, err
	}

	// Create the BlockNumberMonitor
	blockNumberMonitor, err := NewBlockNumberMonitor(logger, ethClient, time.Second*3)
	if err != nil {
		return nil, fmt.Errorf("create block number monitor: %w", err)
	}

	return &RouterAddressProvider{
		routerBinding:      routerBinding,
		blockNumberMonitor: blockNumberMonitor,
	}, nil
}

// GetCertVerifierAddress returns the cert verifier address for the given reference block number
func (rap *RouterAddressProvider) GetCertVerifierAddress(ctx context.Context, referenceBlockNumber uint64) (gethcommon.Address, error) {
	// Wait for the local client to reach the reference block number
	if err := rap.blockNumberMonitor.WaitForBlockNumber(ctx, referenceBlockNumber); err != nil {
		return gethcommon.Address{}, fmt.Errorf("wait for block number: %w", err)
	}

	return rap.routerBinding.GetCertVerifierAt(&bind.CallOpts{Context: ctx}, uint32(referenceBlockNumber))
}
