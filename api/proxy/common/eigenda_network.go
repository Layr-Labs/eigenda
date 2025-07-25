package common

import (
	"fmt"
	"strings"
)

// TODO: this should be moved outside of proxy, since it could be used by other packages/tools.
// For example tools/discovery is currently making use of it.
type EigenDANetwork string

const (
	HoleskyTestnetEigenDANetwork EigenDANetwork = "holesky_testnet"
	HoleskyPreprodEigenDANetwork EigenDANetwork = "holesky_preprod"
	SepoliaTestnetEigenDANetwork EigenDANetwork = "sepolia_testnet"
	MainnetEigenDANetwork        EigenDANetwork = "mainnet"
)

// GetEigenDADirectory returns, as a string, the address of the EigenDADirectory contract for the network.
// For more information about networks and contract addresses, see https://docs.eigenlayer.xyz/eigenda/networks/
func (n EigenDANetwork) GetEigenDADirectory() (string, error) {
	// TODO: These hardcoded addresses should eventually be fetched from the EigenDADirectory contract
	// to reduce duplication and ensure consistency across the codebase
	switch n {
	case MainnetEigenDANetwork:
		return "0x64AB2e9A86FA2E183CB6f01B2D4050c1c2dFAad4", nil
	case HoleskyTestnetEigenDANetwork:
		return "0x90776Ea0E99E4c38aA1Efe575a61B3E40160A2FE", nil
	case HoleskyPreprodEigenDANetwork:
		return "0xfB676e909f376efFDbDee7F17342aCF55f6Ec502", nil
	case SepoliaTestnetEigenDANetwork:
		return "0x9620dC4B3564198554e4D2b06dEFB7A369D90257", nil
	default:
		return "", fmt.Errorf("unknown network type: %s", n)
	}
}

// GetServiceManagerAddress returns, as a string, the address of the EigenDAServiceManager contract for the network.
// For more information about networks and contract addresses, see https://docs.eigenlayer.xyz/eigenda/networks/
// TODO: these should be fetched from the EigenDADirectory contract instead.
func (n EigenDANetwork) GetServiceManagerAddress() (string, error) {
	// TODO: These hardcoded addresses should eventually be fetched from the EigenDADirectory contract
	// to reduce duplication and ensure consistency across the codebase
	switch n {
	case MainnetEigenDANetwork:
		return "0x870679E138bCdf293b7Ff14dD44b70FC97e12fc0", nil
	case HoleskyTestnetEigenDANetwork:
		return "0xD4A7E1Bd8015057293f0D0A557088c286942e84b", nil
	case HoleskyPreprodEigenDANetwork:
		return "0x54A03db2784E3D0aCC08344D05385d0b62d4F432", nil
	case SepoliaTestnetEigenDANetwork:
		return "0x3a5acf46ba6890B8536420F4900AC9BC45Df4764", nil
	default:
		return "", fmt.Errorf("unknown network type: %s", n)
	}
}

// GetDisperserAddress gets a string representing the address of the disperser for the network.
// The format of the returned address is "<hostname>:<port>"
// For more information about networks and disperser endpoints, see https://docs.eigenlayer.xyz/eigenda/networks/
func (n EigenDANetwork) GetDisperserAddress() (string, error) {
	// TODO: These hardcoded addresses should eventually be fetched from the EigenDADirectory contract
	// to reduce duplication and ensure consistency across the codebase
	switch n {
	case MainnetEigenDANetwork:
		return "disperser.eigenda.xyz:443", nil
	case HoleskyTestnetEigenDANetwork:
		return "disperser-testnet-holesky.eigenda.xyz:443", nil
	case HoleskyPreprodEigenDANetwork:
		return "disperser-preprod-holesky.eigenda.xyz:443", nil
	case SepoliaTestnetEigenDANetwork:
		return "disperser-testnet-sepolia.eigenda.xyz:443", nil
	default:
		return "", fmt.Errorf("unknown network type: %s", n)
	}
}

// GetBLSOperatorStateRetrieverAddress returns, as a string, the address of the OperatorStateRetriever contract for the
// network
// For more information about networks and contract addresses, see https://docs.eigenlayer.xyz/eigenda/networks/
// TODO: these should be fetched from the EigenDADirectory contract instead.
func (n EigenDANetwork) GetBLSOperatorStateRetrieverAddress() (string, error) {
	// TODO: These hardcoded addresses should eventually be fetched from the EigenDADirectory contract
	// to reduce duplication and ensure consistency across the codebase
	switch n {
	case MainnetEigenDANetwork:
		return "0xEC35aa6521d23479318104E10B4aA216DBBE63Ce", nil
	case HoleskyTestnetEigenDANetwork, HoleskyPreprodEigenDANetwork:
		return "0x003497Dd77E5B73C40e8aCbB562C8bb0410320E7", nil
	case SepoliaTestnetEigenDANetwork:
		return "0x22478d082E9edaDc2baE8443E4aC9473F6E047Ff", nil
	default:
		return "", fmt.Errorf("unknown network: %s", n)
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

// Useful when an rpc url is provided, but no network is specified.
// In this case, we can use this function to automatically choose a default network.
func DefaultEigenDANetworkFromChainID(chainID string) (EigenDANetwork, error) {
	networks, err := EigenDANetworksFromChainID(chainID)
	if err != nil {
		return "", err
	}
	if len(networks) == 0 {
		return "", fmt.Errorf("no EigenDA network found for chain ID: %s", chainID)
	}
	return networks[0], nil
}

// EigenDANetworkFromString parses an inputString to an EigenDANetwork value.
// The returned EigenDANetwork is guaranteed to be non-nil.
// If an invalid network is provided, an error is returned.
func EigenDANetworkFromString(inputString string) (EigenDANetwork, error) {
	network := EigenDANetwork(inputString)

	switch network {
	case HoleskyTestnetEigenDANetwork, HoleskyPreprodEigenDANetwork, SepoliaTestnetEigenDANetwork, MainnetEigenDANetwork:
		return network, nil
	default:
		allowedNetworks := []string{
			MainnetEigenDANetwork.String(),
			HoleskyTestnetEigenDANetwork.String(),
			HoleskyPreprodEigenDANetwork.String(),
			SepoliaTestnetEigenDANetwork.String(),
		}
		return "", fmt.Errorf("invalid network: %s. Must be one of: %s",
			inputString, strings.Join(allowedNetworks, ", "))
	}
}
