package encoding

import (
	"github.com/Layr-Labs/eigenda/common/math"
)

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetBlobLength(blobSize uint) uint {
	return math.RoundUpDivide(blobSize, BYTES_PER_SYMBOL)
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetBlobLengthPowerOf2(blobSize uint) uint {
	return uint(math.NextPowOf2u32(uint32(GetBlobLength(blobSize))))
}

// GetBlobSize converts from blob length in symbols to blob size in bytes. This is not an exact conversion.
func GetBlobSize(blobLength uint) uint {
	return blobLength * BYTES_PER_SYMBOL
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetEncodedBlobLength(blobLength uint, quorumThreshold, advThreshold uint8) uint {
	return math.RoundUpDivide(blobLength*100, uint(quorumThreshold-advThreshold))
}
