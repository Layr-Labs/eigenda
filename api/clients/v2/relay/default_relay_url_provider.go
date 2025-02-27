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

// relayUrlProvider provides relay URL strings, based on relay key.
type relayUrlProvider struct {
	relayRegistryCaller *relayRegistryBindings.ContractEigenDARelayRegistryCaller
}

var _ RelayUrlProvider = &relayUrlProvider{}

// NewRelayUrlProvider constructs a relayUrlProvider
func NewRelayUrlProvider(
	ethClient common.EthClient,
	relayRegistryAddress gethcommon.Address,
) (RelayUrlProvider, error) {
	relayRegistryContractCaller, err := relayRegistryBindings.NewContractEigenDARelayRegistryCaller(
		relayRegistryAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("NewContractEigenDARelayRegistryCaller: %w", err)
	}

	return &relayUrlProvider{
		relayRegistryCaller: relayRegistryContractCaller,
	}, nil
}

// GetRelayUrl gets the URL string for a given relayKey
func (rup *relayUrlProvider) GetRelayUrl(ctx context.Context, relayKey v2.RelayKey) (string, error) {
	relayUrl, err := rup.relayRegistryCaller.RelayKeyToUrl(&bind.CallOpts{Context: ctx}, relayKey)
	if err != nil {
		return "", fmt.Errorf("fetch relay key (%d) URL from EigenDARelayRegistry contract: %w", relayKey, err)
	}

	return relayUrl, nil
}

// GetRelayCount gets the number of relays that exist in the registry
func (rup *relayUrlProvider) GetRelayCount(ctx context.Context) (uint32, error) {
	// NextRelayKey initializes to 0, and is incremented each time a relay is added
	// current logic doesn't support removing relays, so NextRelayKey therefore corresponds directly to relay count
	relayCount, err := rup.relayRegistryCaller.NextRelayKey(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("get next relay key from EigenDARelayRegistry contract: %w", err)
	}

	return relayCount, nil
}
