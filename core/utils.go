package core

import (
	"math"
	"math/big"

	"golang.org/x/exp/constraints"
)

func RoundUpDivideBig(a, b *big.Int) *big.Int {

	one := new(big.Int).SetUint64(1)
	num := new(big.Int).Sub(new(big.Int).Add(a, b), one) // a + b - 1
	res := new(big.Int).Div(num, b)                      // (a + b - 1) / b
	return res
}

func RoundUpDivide[T constraints.Integer](a, b T) T {
	return (a + b - 1) / b
}

func NextPowerOf2[T constraints.Integer](d T) T {
	nextPower := math.Ceil(math.Log2(float64(d)))
	return T(math.Pow(2.0, nextPower))
}

// PadToPowerOf2 pads a byte slice to the nearest power of 2 length by appending zeros
func PadToPowerOf2(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return []byte{0}
	}

	// If length is already a power of 2, return original
	if length&(length-1) == 0 {
		return data
	}

	// Create a new slice with the power-of-2 length
	paddedData := make([]byte, NextPowerOf2(uint64(len(data))))

	// Copy original data
	copy(paddedData, data)

	// The remaining bytes will be zero by default
	return paddedData
}
