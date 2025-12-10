package common

import (
	"fmt"
	"math/big"
)

// Converts a chain ID to 32-byte big-endian representation compatible with EIP-155.
// Returns an empty byte slice if chainId is nil.
func ChainIdToBytes(chainId *big.Int) []byte {
	if chainId == nil {
		return nil
	}

	bytes := make([]byte, 32)
	chainId.FillBytes(bytes)
	return bytes
}

// Converts 32-byte big-endian bytes to a chain ID.
//
// Returns an error if the input is not 32 bytes.
func ChainIdFromBytes(bytes []byte) (*big.Int, error) {
	if len(bytes) != 32 {
		return nil, fmt.Errorf("chainID must be 32 bytes, got %d", len(bytes))
	}
	return new(big.Int).SetBytes(bytes), nil
}
