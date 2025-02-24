package clients

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/common"
	relayRegistryBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARelayRegistry"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// RelayUrlProvider provides relay URL strings, based on relay key.
//
// Contains an internal cache, so that a given URL doesn't need to be fetched multiple times.
type RelayUrlProvider struct {
	logger logging.Logger
	relayRegistryCaller *relayRegistryBindings.ContractEigenDARelayRegistryCaller
	relayUrlCache sync.Map
}

// NewRelayUrlProvider constructs a RelayUrlProvider.
//
// This method initializes the provider's internal cache with the URLs of all relays that exist at the time of construction.
func NewRelayUrlProvider(
	ctx context.Context,
	logger logging.Logger,
	ethClient common.EthClient,
	relayRegistryAddress string,
) (*RelayUrlProvider, error) {
	relayRegistryContractCaller, err := relayRegistryBindings.NewContractEigenDARelayRegistryCaller(
		gethcommon.HexToAddress(relayRegistryAddress),
		ethClient)
	if err != nil {
		return nil, fmt.Errorf("NewContractEigenDARelayRegistryCaller: %w", err)
	}

	relayUrlProvider := &RelayUrlProvider{
		logger: logger,
		relayRegistryCaller: relayRegistryContractCaller,
	}

	err = relayUrlProvider.initializeCache(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialize relay URL cache: %w", err)
	}

	return relayUrlProvider, nil
}

// GetRelayUrl gets the URL string for a given relayKey
//
// If the internal cache already knows the URL for the relayKey, the known value is returned immediately.
// If the internal cache doesn't already know the URL for the relayKey, it attempts to fetch the URL with a call to the
// EigenDARelayRegistry contract. It returns the fetched value after the contract call succeeds, or returns an error if
// the call fails.
func (rup *RelayUrlProvider) GetRelayUrl(ctx context.Context, relayKey uint32) (string, error) {
	// the current contract doesn't allow updating the URL for a given relayKey, so if the value exists in the cache,
	// it's guaranteed to be correct.
	existingRelayUrl, valueFound := rup.relayUrlCache.Load(relayKey)
	if valueFound {
		return existingRelayUrl.(string), nil
	}

	fetchedUrl, err := rup.relayRegistryCaller.RelayKeyToUrl(&bind.CallOpts{Context: ctx}, relayKey)
	if err != nil {
		return "", fmt.Errorf("fetch relay key URL from EigenDARelayRegistry contract: %w", err)
	}

	rup.relayUrlCache.Store(relayKey, fetchedUrl)

	return fetchedUrl, nil
}

// initializeCache fetches the URL for all relays that exist at the time the RelayUrlProvider is created
//
// Returns an error if unable to fetch the number of relays that exist. If any given URL fetch fails during
// initialization, no error is returned: this method will do its best to initialize all relay URLs, even
// if some fetches fail.
func (rup *RelayUrlProvider) initializeCache(ctx context.Context) error {
	relayCount, err := rup.relayRegistryCaller.NextRelayKey(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("get next relay key from EigenDARelayRegistry contract: %w", err)
	}

	for relayKey := uint32(0); relayKey < relayCount; relayKey++ {
		// getting the url causes it to be saved in the cache
		_, err := rup.GetRelayUrl(ctx, relayKey)
		if err != nil {
			rup.logger.Errorf("failed to get URL for relay key %d: %v", relayKey, err)
		}
	}

	return nil
}
