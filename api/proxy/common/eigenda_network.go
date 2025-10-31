package common

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	common_eigenda "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	geth_common "github.com/ethereum/go-ethereum/common"
)

// TODO: this should be moved outside of proxy, since it could be used by other packages/tools.
// For example tools/discovery is currently making use of it.
type EigenDANetwork string

const (
	HoleskyTestnetEigenDANetwork EigenDANetwork = "holesky_testnet"
	HoleskyPreprodEigenDANetwork EigenDANetwork = "holesky_preprod"
	SepoliaTestnetEigenDANetwork EigenDANetwork = "sepolia_testnet"
	HoodiTestnetEigenDANetwork   EigenDANetwork = "hoodi_testnet"
	MainnetEigenDANetwork        EigenDANetwork = "mainnet"
)

// GetEigenDADirectory returns, as a string, the address of the EigenDADirectory contract for the network.
// For more information about networks and contract addresses, see https://docs.eigenlayer.xyz/eigenda/networks/
func (n EigenDANetwork) GetEigenDADirectory() string {
	// TODO: These hardcoded addresses should eventually be fetched from the EigenDADirectory contract
	// to reduce duplication and ensure consistency across the codebase
	switch n {
	case MainnetEigenDANetwork:
		return "0x64AB2e9A86FA2E183CB6f01B2D4050c1c2dFAad4"
	case HoleskyTestnetEigenDANetwork:
		return "0x90776Ea0E99E4c38aA1Efe575a61B3E40160A2FE"
	case HoleskyPreprodEigenDANetwork:
		return "0xfB676e909f376efFDbDee7F17342aCF55f6Ec502"
	case SepoliaTestnetEigenDANetwork:
		return "0x9620dC4B3564198554e4D2b06dEFB7A369D90257"
	case HoodiTestnetEigenDANetwork:
		return "0x5a44e56e88abcf610c68340c6814ae7f5c4369fd"
	default:
		panic(fmt.Sprintf("unknown EigenDA network: %s", n))
	}
}

// GetDisperserAddress gets a string representing the address of the disperser for the network.
// The format of the returned address is "<hostname>:<port>"
// For more information about networks and disperser endpoints, see https://docs.eigenlayer.xyz/eigenda/networks/
func (n EigenDANetwork) GetDisperserAddress() string {
	// TODO: These hardcoded addresses should eventually be fetched from the EigenDADirectory contract
	// to reduce duplication and ensure consistency across the codebase
	switch n {
	case MainnetEigenDANetwork:
		return "disperser.eigenda.xyz:443"
	case HoleskyTestnetEigenDANetwork:
		return "disperser-testnet-holesky.eigenda.xyz:443"
	case HoleskyPreprodEigenDANetwork:
		return "disperser-preprod-holesky.eigenda.xyz:443"
	case SepoliaTestnetEigenDANetwork:
		return "disperser-testnet-sepolia.eigenda.xyz:443"
	case HoodiTestnetEigenDANetwork:
		return "disperser-testnet-hoodi.eigenda.xyz:443"
	default:
		panic(fmt.Sprintf("unknown EigenDA network: %s", n))
	}
}

func (n EigenDANetwork) String() string {
	return string(n)
}

// chainIDToNetworkMap maps chain IDs to EigenDA networks
var chainIDToNetworkMap = map[string][]EigenDANetwork{
	"1":        {MainnetEigenDANetwork},
	"17000":    {HoleskyTestnetEigenDANetwork, HoleskyPreprodEigenDANetwork},
	"11155111": {SepoliaTestnetEigenDANetwork},
	"560048":   {HoodiTestnetEigenDANetwork},
}

// EigenDANetworksFromChainID returns the EigenDA network(s) for a given chain ID
// If no error occurs, the returned slice will contain one or more EigenDANetwork values.
func EigenDANetworksFromChainID(chainID string) ([]EigenDANetwork, error) {
	networks, ok := chainIDToNetworkMap[chainID]
	if !ok {
		return nil, fmt.Errorf("unknown chain ID: %s", chainID)
	}
	return networks, nil
}

// EigenDANetworkFromString parses an inputString to an EigenDANetwork value.
// The returned EigenDANetwork is guaranteed to be non-nil.
// If an invalid network is provided, an error is returned.
func EigenDANetworkFromString(inputString string) (EigenDANetwork, error) {
	network := EigenDANetwork(inputString)

	switch network {
	case HoleskyTestnetEigenDANetwork, HoleskyPreprodEigenDANetwork, SepoliaTestnetEigenDANetwork,
		HoodiTestnetEigenDANetwork, MainnetEigenDANetwork:
		return network, nil
	default:
		allowedNetworks := []string{
			MainnetEigenDANetwork.String(),
			HoleskyTestnetEigenDANetwork.String(),
			HoleskyPreprodEigenDANetwork.String(),
			SepoliaTestnetEigenDANetwork.String(),
			HoodiTestnetEigenDANetwork.String(),
		}
		return "", fmt.Errorf("invalid network: %s. Must be one of: %s",
			inputString, strings.Join(allowedNetworks, ", "))
	}
}

// BuildEthClient creates an Ethereum client using the provided RPC URL and, if set, validates that the chain ID
// matches the expected EigenDA network. It returns an ethClient, it's ChainID, and an error.
func BuildEthClient(ctx context.Context, log logging.Logger, ethRpcUrl string,
	expectedNetwork EigenDANetwork) (common_eigenda.EthClient, string, error) {
	gethCfg := geth.EthClientConfig{
		RPCURLs: []string{ethRpcUrl},
	}

	ethClient, err := geth.NewClient(gethCfg, geth_common.Address{}, 0, log)
	if err != nil {
		return nil, "", fmt.Errorf("create geth client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get chain ID from ETH RPC: %w", err)
	}

	log.Infof("Using chain id: %d", chainID.Uint64())

	// Validate that the chain ID matches the expected network
	if expectedNetwork != "" {
		actualNetworks, err := EigenDANetworksFromChainID(chainID.String())
		if err != nil {
			return nil, "", fmt.Errorf("unknown chain ID %s: %w", chainID.String(), err)
		}
		if !slices.Contains(actualNetworks, expectedNetwork) {
			return nil, "", fmt.Errorf("network mismatch: expected %s (based on configuration), but ETH RPC "+
				"returned chain ID %s which corresponds to %s",
				expectedNetwork, chainID.String(), actualNetworks)
		}

		log.Infof("Detected EigenDA network: %s. Will use for reading network default values if overrides "+
			"aren't provided.", expectedNetwork.String())
	}

	return ethClient, chainID.String(), nil
}
