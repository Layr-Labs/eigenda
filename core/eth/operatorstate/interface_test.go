package operatorstate

import (
	"testing"

	"github.com/Layr-Labs/eigenda/core"
)

// TestIndexedChainStateInterfaceCompliance verifies that IndexedChainState implements
// the core.IndexedChainState interface completely
func TestIndexedChainStateInterfaceCompliance(t *testing.T) {
	// This test ensures that IndexedChainState properly implements the required interfaces
	// If the interface is not fully implemented, this will fail to compile

	var _ core.ChainState = (*IndexedChainState)(nil)
	var _ core.IndexedChainState = (*IndexedChainState)(nil)

	t.Log("IndexedChainState implements core.IndexedChainState interface")
}

// TestContractClientCompilation verifies ContractClient compiles correctly
func TestContractClientCompilation(t *testing.T) {
	// This test ensures ContractClient compiles without issues
	// If there are missing imports or syntax errors, this will fail to compile

	var _ = (*ContractClient)(nil)

	t.Log("ContractClient compiles successfully")
}
