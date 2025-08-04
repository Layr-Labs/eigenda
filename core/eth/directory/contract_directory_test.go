package directory

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: unit test that compares this against /Users/cody/ws/master-eigenda/contracts/src/core/libraries/v3/address-directory/AddressDirectoryConstants.sol

// TODO fix path
const solConstantsPath = "../../../contracts/src/core/libraries/v3/address-directory/AddressDirectoryConstants.sol"

func parseAddressDirectoryConstants(t *testing.T) map[string]struct{} {
	contractNames := make(map[string]struct{})

	solString, err := os.ReadFile(solConstantsPath)
	require.NoError(t, err)

	lines := strings.Split(string(solString), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)

		// The lines we want to parse have the following format:
		//        string internal constant REGISTRY_COORDINATOR_NAME = "REGISTRY_COORDINATOR";

		if strings.HasPrefix(lines[i], "//") ||
			lines[i] == "" ||
			!strings.HasPrefix(line, "string internal constant") ||
			strings.Count(line, "=") != 1 ||
			strings.Count(line, "\"") != 2 {
			continue
		}

		firstQuoteIndex := strings.Index(line, "\"")
		lastQuoteIndex := strings.LastIndex(line, "\"")
		if firstQuoteIndex == -1 || lastQuoteIndex == -1 || firstQuoteIndex >= lastQuoteIndex {
			continue
		}

		// Extract the contract name from the line
		contractName := line[firstQuoteIndex+1 : lastQuoteIndex]
		if contractName == "" {
			continue
		}

		contractNames[contractName] = struct{}{}
	}

	return contractNames
}

func TestContractNameList(t *testing.T) {
	parsedContractSet := parseAddressDirectoryConstants(t)

	knownContractSet := make(map[string]struct{}, len(knownContracts))
	for _, contractName := range knownContracts {
		knownContractSet[string(contractName)] = struct{}{}
	}

	for contractName := range parsedContractSet {
		_, exists := knownContractSet[contractName]
		require.Truef(t, exists,
			"Contract %s is defined in the Solidity constants but not in the known contracts list", contractName)
	}
	for contractName := range knownContractSet {
		_, exists := parsedContractSet[contractName]
		require.Truef(t, exists,
			"Contract %s is defined in the known contracts list but not in the Solidity constants", contractName)
	}
}
