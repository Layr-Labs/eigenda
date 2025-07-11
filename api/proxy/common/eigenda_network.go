package common

import "fmt"

type EigenDANetwork string

const (
	MainnetEigenDANetwork        EigenDANetwork = "mainnet"
	HoleskyTestnetEigenDANetwork EigenDANetwork = "holesky_testnet"
	HoleskyPreprodEigenDANetwork EigenDANetwork = "holesky_preprod"
	SepoliaTestnetEigenDANetwork EigenDANetwork = "sepolia_testnet"
)

// GetEigenDADirectory returns, as a string, the address of the EigenDADirectory contract for the network.
func (n EigenDANetwork) GetEigenDADirectory() (string, error) {
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


// GetDisperserAddress gets a string representing the address of the disperser for the network.
// The format of the returned address is "<hostname>:<port>"
func (n EigenDANetwork) GetDisperserAddress() (string, error) {
	switch n {
	case MainnetEigenDANetwork:
		return "disperser-mainnet.eigenda.xyz:443", nil
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
func EigenDANetworksFromChainID(chainID string) ([]EigenDANetwork, error) {
	networks, ok := chainIDToNetworkMap[chainID]
	if !ok {
		return nil, fmt.Errorf("unknown chain ID: %s", chainID)
	}
	return networks, nil
}

func EigenDANetworkFromString(inputString string) (EigenDANetwork, error) {
	network := EigenDANetwork(inputString)

	switch network {
	case MainnetEigenDANetwork, HoleskyTestnetEigenDANetwork, HoleskyPreprodEigenDANetwork, SepoliaTestnetEigenDANetwork:
		return network, nil
	default:
		return "", fmt.Errorf("unknown network type: %s", inputString)
	}
}
