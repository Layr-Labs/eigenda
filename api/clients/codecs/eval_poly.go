package codecs

import (
	"fmt"
	"math"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// evalPoly is a polynomial in evaluation form.
//
// The underlying bytes represent 32 byte field elements, and the field elements represent the polynomial evaluation
// at roots of unity.
//
// The number of field elements is always a power of 2.
type evalPoly struct {
	fieldElements []fr.Element
}

// evalPolyFromElements creates a new evalPoly from field elements.
func evalPolyFromElements(elements []fr.Element) *evalPoly {
	return &evalPoly{fieldElements: elements}
}

// toCoeffPoly converts an evalPoly to a coeffPoly, using the IFFT operation
func (p *evalPoly) toCoeffPoly() (*coeffPoly, error) {
	// we need to pad to the next power of 2, to be able to take the FFT
	paddedLength := encoding.NextPowerOf2(len(p.fieldElements))
	padding := make([]fr.Element, paddedLength-len(p.fieldElements))
	paddedElements := append(p.fieldElements, padding...)

	maxScale := uint8(math.Log2(float64(len(paddedElements))))
	ifftedElements, err := fft.NewFFTSettings(maxScale).FFT(paddedElements, true)
	if err != nil {
		return nil, fmt.Errorf("perform IFFT: %w", err)
	}

	return coeffPolyFromElements(ifftedElements), nil
}

// toEncodedPayload converts an evalPoly into an encoded payload
//
// This conversion entails removing the power-of-2 padding which is added to an encodedPayload when originally creating
// an evalPoly.
func (p *evalPoly) toEncodedPayload() (*encodedPayload, error) {
	return encodedPayloadFromElements(p.fieldElements)
}
