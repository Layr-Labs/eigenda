package calldata_gas_estimator

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/tools/integration_utils/altdacommitment_parser"
	"github.com/Layr-Labs/eigenda/tools/integration_utils/flags"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli"
)

// CalldataGasInfo holds the breakdown of calldata gas calculations
type CalldataGasInfo struct {
	DataSize     int
	Zeros        uint64
	Nonzeros     uint64
	EIP2028Gas   uint64
	EIP7623Floor uint64
	FinalGas     uint64
}

func RunEstimator(ctx *cli.Context) {
	hexString := ctx.String(flags.CertHexFlag.Name)

	// Process the hex string to get binary data
	data, err := altdacommitment_parser.ProcessHexString(hexString)
	if err != nil {
		fmt.Printf("Gas Cost Estimation: Failed to process hex string: %v\n", err)
		return
	}

	// Calculate calldata gas
	info := CalculateCalldataGas(data)

	// display gas cost
	DisplayCalldataGasCost(info)
}

// CalculateCalldataGas processes data and returns detailed gas calculation breakdown
func CalculateCalldataGas(data []byte) CalldataGasInfo {
	info := CalldataGasInfo{
		DataSize: len(data),
	}

	if len(data) == 0 {
		return info
	}

	// Count zero and non-zero bytes
	for _, b := range data {
		if b == 0 {
			info.Zeros++
		} else {
			info.Nonzeros++
		}
	}

	// EIP-2028 "traditional" data gas pricing
	// 4 gas per zero byte, 16 gas per non-zero byte
	info.EIP2028Gas = info.Zeros*params.TxDataZeroGas + info.Nonzeros*params.TxDataNonZeroGasEIP2028

	// EIP-7623 floor pricing (tokens: 1 per zero byte, 4 per non-zero byte; 10 gas per token)
	// This creates a minimum floor to prevent cheap spam attacks
	tokens := info.Zeros + info.Nonzeros*params.TxTokenPerNonZeroByte
	info.EIP7623Floor = tokens * params.TxCostFloorPerToken

	// Return the higher of EIP-2028 traditional pricing or EIP-7623 floor pricing
	// This ensures we charge at least the floor price while maintaining backward compatibility
	if info.EIP7623Floor > info.EIP2028Gas {
		info.FinalGas = info.EIP7623Floor
	} else {
		info.FinalGas = info.EIP2028Gas
	}

	return info
}

// CalldataGas returns the gas charged for the calldata bytes alone (no 21k base, no access list).
// This function implements both EIP-2028 traditional data gas and EIP-7623 floor pricing,
// returning the higher of the two values.
func CalldataGas(data []byte) uint64 {
	return CalculateCalldataGas(data).FinalGas
}
