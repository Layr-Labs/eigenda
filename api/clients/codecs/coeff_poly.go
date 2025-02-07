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

// coeffPolyFromBytes creates a new polynomial from bytes. This function performs the necessary checks to guarantee that the
// bytes are well-formed, and returns a new object if they are
func coeffPolyFromBytes(bytes []byte) (*coeffPoly, error) {
	fieldElements, err := rs.BytesToFieldElements(bytes)
	if err != nil {
		return nil, fmt.Errorf("deserialize field elements: %w", err)
	}

	return &coeffPoly{fieldElements: fieldElements}, nil
}

// coeffPolyFromElements creates a new coeffPoly from field elements.
func coeffPolyFromElements(elements []fr.Element) *coeffPoly {
	return &coeffPoly{fieldElements: elements}
}

// toEvalPoly converts a coeffPoly to an evalPoly, using the FFT operation
func (p *coeffPoly) toEvalPoly() (*evalPoly, error) {
	// we need to pad to the next power of 2, to be able to take the FFT
	paddedLength := encoding.NextPowerOf2(len(p.fieldElements))
	padding := make([]fr.Element, paddedLength-len(p.fieldElements))
	paddedElements := append(p.fieldElements, padding...)

	maxScale := uint8(math.Log2(float64(len(paddedElements))))
	fftedElements, err := fft.NewFFTSettings(maxScale).FFT(paddedElements, false)
	if err != nil {
		return nil, fmt.Errorf("perform FFT: %w", err)
	}

	return evalPolyFromElements(fftedElements), nil
}

// GetBytes returns the bytes that underlie the polynomial
func (p *coeffPoly) getBytes() []byte {
	return rs.FieldElementsToBytes(p.fieldElements)
}

// toEncodedPayload converts a coeffPoly into an encoded payload
//
// This conversion entails removing the power-of-2 padding which is added to an encodedPayload when originally creating
// an evalPoly.
func (p *coeffPoly) toEncodedPayload() (*encodedPayload, error) {
	return encodedPayloadFromElements(p.fieldElements)
}
