package codecs

import (
	"fmt"
	"math"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// coeffPoly is a polynomial in coefficient form.
//
// The underlying bytes represent 32 byte field elements, and each field element represents a coefficient
type coeffPoly struct {
	fieldElements []fr.Element
}

// coeffPolyFromBytes creates a new coeffPoly from bytes. This function performs the necessary checks to guarantee that the
// bytes are well-formed, and returns a new object if they are
func coeffPolyFromBytes(bytes []byte) (*coeffPoly, error) {
	if !encoding.IsPowerOfTwo(len(bytes)) {
		return nil, fmt.Errorf("bytes have length %d, expected a power of 2", len(bytes))
	}

	fieldElements, err := rs.BytesToFieldElements(bytes)
	if err != nil {
		return nil, fmt.Errorf("deserialize field elements: %w", err)
	}

	return &coeffPoly{fieldElements: fieldElements}, nil
}

// coeffPolyFromElements creates a new coeffPoly from field elements.
func coeffPolyFromElements(elements []fr.Element) (*coeffPoly, error) {
	return &coeffPoly{fieldElements: elements}, nil
}

// toEvalPoly converts a coeffPoly to an evalPoly, using the FFT operation
func (cp *coeffPoly) toEvalPoly() (*evalPoly, error) {
	maxScale := uint8(math.Log2(float64(len(cp.fieldElements))))
	fftedElements, err := fft.NewFFTSettings(maxScale).FFT(cp.fieldElements, false)
	if err != nil {
		return nil, fmt.Errorf("perform FFT: %w", err)
	}

	evalPoly, err := evalPolyFromElements(fftedElements)
	if err != nil {
		return nil, fmt.Errorf("construct eval poly: %w", err)
	}

	return evalPoly, nil
}

// GetBytes returns the bytes that underlie the polynomial
func (cp *coeffPoly) getBytes() []byte {
	return rs.FieldElementsToBytes(cp.fieldElements)
}
