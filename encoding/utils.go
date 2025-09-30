package encoding

import (
	"github.com/Layr-Labs/eigenda/common/math"
)

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetBlobLength(blobSize uint32) uint32 {
	return math.RoundUpDivide(blobSize, BYTES_PER_SYMBOL)
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetBlobLengthPowerOf2(blobSize uint32) uint32 {
	return math.NextPowOf2u32(GetBlobLength(blobSize))
}

// GetBlobSize converts from blob length in symbols to blob size in bytes. This is not an exact conversion.
func GetBlobSize(blobLength uint32) uint32 {
	return blobLength * BYTES_PER_SYMBOL
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetEncodedBlobLength(blobLength uint32, quorumThreshold, advThreshold uint8) uint32 {
	return math.RoundUpDivide(blobLength*100, uint32(quorumThreshold-advThreshold))
}
