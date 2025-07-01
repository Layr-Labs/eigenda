package common

import "fmt"

type EigenDANetwork string

const (
	HoleskyTestnetEigenDANetwork EigenDANetwork = "holesky_testnet"
	HoleskyPreprodEigenDANetwork EigenDANetwork = "holesky_preprod"
	SepoliaTestnetEigenDANetwork EigenDANetwork = "sepolia_testnet"
)

// GetServiceManagerAddress returns, as a string, the address of the EigenDAServiceManager contract for the network.
func (n EigenDANetwork) GetServiceManagerAddress() (string, error) {
	switch n {
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
func (n EigenDANetwork) GetDisperserAddress() (string, error) {
	switch n {
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
func (n EigenDANetwork) GetBLSOperatorStateRetrieverAddress() (string, error) {
	switch n {
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
	case HoleskyTestnetEigenDANetwork, HoleskyPreprodEigenDANetwork, SepoliaTestnetEigenDANetwork:
		return network, nil
	default:
		return "", fmt.Errorf("unknown network type: %s", inputString)
	}
}
