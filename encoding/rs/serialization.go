package rs

import (
	"errors"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

const BYTES_PER_SYMBOL = 32 // fr.Element serializes to 32 bytes

// Serialize serializes Frame to a compact byte format
// Each fr.Element is serialized consecutively
func Serialize(f *Frame) ([]byte, error) {
	// Pre-allocate buffer for all coefficients
	coded := make([]byte, 0, BYTES_PER_SYMBOL*len(f.Coeffs))

	// Append each coefficient's bytes
	for _, coeff := range f.Coeffs {
		coded = append(coded, coeff.Marshal()...)
	}

	return coded, nil
}

// Deserialize reconstructs Frame from compact byte format
func Deserialize(data []byte) (*Frame, error) {
	if len(data)%BYTES_PER_SYMBOL != 0 {
		return nil, errors.New("invalid data length")
	}

	numCoeffs := len(data) / BYTES_PER_SYMBOL
	frame := &Frame{
		Coeffs: make([]fr.Element, numCoeffs),
	}

	// Read coefficients
	for i := 0; i < numCoeffs; i++ {
		start := i * BYTES_PER_SYMBOL
		frame.Coeffs[i].Unmarshal(data[start : start+BYTES_PER_SYMBOL])
	}

	return frame, nil
}
