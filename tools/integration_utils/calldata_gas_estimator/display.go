package calldata_gas_estimator

import (
	"fmt"
)

// DisplayCalldataGasCost displays the estimated gas cost for posting the certificate as calldata
func DisplayCalldataGasCost(gasInfo CalldataGasInfo) {
	// Calculate EIP-7623 tokens for display
	// see https://eips.ethereum.org/EIPS/eip-7623
	eip7623Tokens := gasInfo.Zeros + gasInfo.Nonzeros*4 // TxTokenPerNonZeroByte=4

	fmt.Printf("\nStatic Calldata Gas Rough Cost Estimation:\n")
	fmt.Printf("  Data Size: %d bytes (%d zero, %d non-zero)\n", gasInfo.DataSize, gasInfo.Zeros, gasInfo.Nonzeros)
	fmt.Printf("  EIP-2028 Cost: %d gas (4×%d + 16×%d)\n", gasInfo.EIP2028Gas, gasInfo.Zeros, gasInfo.Nonzeros)
	fmt.Printf("  EIP-7623 Floor: %d gas (%d tokens × 10 gas/token)\n", gasInfo.EIP7623Floor, eip7623Tokens)
	fmt.Printf("  Calldata Gas: %d gas (higher of the two)\n", gasInfo.FinalGas)
	fmt.Printf("  Total with 21k base: %d gas\n", gasInfo.FinalGas+21000)
}
