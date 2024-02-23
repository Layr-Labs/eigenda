package encoding

import (
	"github.com/Layr-Labs/eigenda/pkg/kzg"
	"golang.org/x/exp/constraints"

	"math"
)

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetBlobLength(blobSize uint) uint {
	symSize := uint(kzg.BYTES_PER_COEFFICIENT)
	return (blobSize + symSize - 1) / symSize
}

// GetBlobSize converts from blob length in symbols to blob size in bytes. This is not an exact conversion.
func GetBlobSize(blobLength uint) uint {
	return blobLength * kzg.BYTES_PER_COEFFICIENT
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetEncodedBlobLength(blobLength uint, quorumThreshold, advThreshold uint8) uint {
	return roundUpDivide(blobLength*100, uint(quorumThreshold-advThreshold))
}

func NextPowerOf2(d uint64) uint64 {
	nextPower := math.Ceil(math.Log2(float64(d)))
	return uint64(math.Pow(2.0, nextPower))
}

// func roundUpDivideBig(a, b *big.Int) *big.Int {

// 	one := new(big.Int).SetUint64(1)
// 	num := new(big.Int).Sub(new(big.Int).Add(a, b), one) // a + b - 1
// 	res := new(big.Int).Div(num, b)                      // (a + b - 1) / b
// 	return res

// }

func roundUpDivide[T constraints.Integer](a, b T) T {
	return (a + b - 1) / b

}
