package clients

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common/geth"
	registry_binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARelayRegistry"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/ethereum/go-ethereum/common"
)

// Provides accessor methods for read interfacing with
// relay registry contract (core/EigenDARelayRegistry.sol)
type RelayRegistryCaller interface {
	// Reads the current sockets or key->url mapping
	// from the on-chain registry
	GetSockets() (map[corev2.RelayKey]string, error)
}

// NewRelayRegistryClient constructs a RelayRegistryClient
func NewRelayRegistryClient(
	ethClient geth.EthClient,
	relayRegistryAddr string,
) (RelayRegistryClient, error) {

	registryCaller, err := registry_binding.NewContractEigenDARelayRegistryCaller(
		common.HexToAddress(relayRegistryAddr),
		ethClient)

	if err != nil {
		return nil, fmt.Errorf("bind to relay registry contract at %s: %w", relayRegistryAddr, err)
	}

	return &relayRegistryClient{
		caller: registryCaller,
	}, nil
}

// TODO: Mitigate the risk of contract mutability where the on-chain registry
// can be subject to new assertions or relay mapping changes.
type relayRegistryClient struct {
	caller *registry_binding.ContractEigenDARelayRegistryCaller
}

func (rrc *relayRegistryClient) GetSockets() (map[corev2.RelayKey]string, error) {
	// read the # of relays by processing next key position
	key, err := rrc.caller.NextRelayKey(nil)
	if err != nil {
		return nil, fmt.Errorf("get next relay key: %+w", err)
	}

	// iterate over each relay key index to construct registry state mapping

	m := make(map[corev2.RelayKey]string)
	for i := uint32(0); i < key; i++ {
		url, err := rrc.caller.RelayKeyToUrl(nil, i)
		if err != nil {
			return nil, fmt.Errorf("fetch url for relay #%d: %+w", i, err)
		}

		m[i] = url
	}

	return m, nil
}