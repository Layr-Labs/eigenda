package encoding

import (
	"golang.org/x/exp/constraints"

	"math"
)

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetBlobLength(blobSize uint) uint {
	return RoundUpDivide(blobSize, BYTES_PER_SYMBOL)
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetBlobLengthPowerOf2(blobSize uint) uint {
	return NextPowerOf2(GetBlobLength(blobSize))
}

// GetBlobSize converts from blob length in symbols to blob size in bytes. This is not an exact conversion.
func GetBlobSize(blobLength uint) uint {
	return blobLength * BYTES_PER_SYMBOL
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetEncodedBlobLength(blobLength uint, quorumThreshold, advThreshold uint8) uint {
	return RoundUpDivide(blobLength*100, uint(quorumThreshold-advThreshold))
}

func RoundUpDivide[T constraints.Integer](a, b T) T {
	return (a + b - 1) / b
}

func NextPowerOf2[T constraints.Integer](d T) T {
	nextPower := math.Ceil(math.Log2(float64(d)))
	return T(math.Pow(2.0, nextPower))
}
