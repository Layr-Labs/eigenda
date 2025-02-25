package relay

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	relayRegistryBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARelayRegistry"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// DefaultRelayUrlProvider provides relay URL strings, based on relay key.
type DefaultRelayUrlProvider struct {
	relayRegistryCaller *relayRegistryBindings.ContractEigenDARelayRegistryCaller
}

var _ RelayUrlProvider = &DefaultRelayUrlProvider{}

// NewDefaultRelayUrlProvider constructs a DefaultRelayUrlProvider
func NewDefaultRelayUrlProvider(
	ethClient common.EthClient,
	relayRegistryAddress gethcommon.Address,
) (*DefaultRelayUrlProvider, error) {
	relayRegistryContractCaller, err := relayRegistryBindings.NewContractEigenDARelayRegistryCaller(
		relayRegistryAddress,
		ethClient)
	if err != nil {
		return nil, fmt.Errorf("NewContractEigenDARelayRegistryCaller: %w", err)
	}

	return &DefaultRelayUrlProvider{
		relayRegistryCaller: relayRegistryContractCaller,
	}, nil
}

// GetRelayUrl gets the URL string for a given relayKey
func (rup *DefaultRelayUrlProvider) GetRelayUrl(ctx context.Context, relayKey v2.RelayKey) (string, error) {
	relayUrl, err := rup.relayRegistryCaller.RelayKeyToUrl(&bind.CallOpts{Context: ctx}, relayKey)
	if err != nil {
		return "", fmt.Errorf("fetch relay key (%d) URL from EigenDARelayRegistry contract: %w", relayKey, err)
	}

	return relayUrl, nil
}

// GetRelayCount gets the number of relays that exist in the registry
func (rup *DefaultRelayUrlProvider) GetRelayCount(ctx context.Context) (uint32, error) {
	relayCount, err := rup.relayRegistryCaller.NextRelayKey(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("get next relay key from EigenDARelayRegistry contract: %w", err)
	}

	return relayCount, nil
}
