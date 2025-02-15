package codecs

import (
	"fmt"

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
	fieldElements, err := rs.ToFrArray(bytes)
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
//
// blobLength (in SYMBOLS) is required, to be able to choose the correct parameters when performing FFT
func (p *coeffPoly) toEvalPoly(blobLength uint32) (*evalPoly, error) {
	// TODO (litt3): this could conceivably be optimized, so that multiple objects share an instance of FFTSettings,
	//  which has enough roots of unity for general use. If the following construction of FFTSettings ever proves
	//  to present a computational burden, consider making this change.
	fftSettings := fft.FFTSettingsFromBlobLength(blobLength)

	// the FFT method pads to the next power of 2, so we don't need to do that manually
	fftedElements, err := fftSettings.FFT(p.fieldElements, false)
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
// blobLength is required, to be able to perform length checks on the encoded payload during construction.
// blobLength is in symbols, NOT bytes
func (p *coeffPoly) toEncodedPayload(blobLength uint32) (*encodedPayload, error) {
	return encodedPayloadFromElements(p.fieldElements, blobLength)
}
