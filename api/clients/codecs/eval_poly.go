package codecs

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// evalPoly is a polynomial in evaluation form.
//
// The underlying bytes represent 32 byte field elements, and the field elements represent the polynomial evaluation
// at roots of unity.
type evalPoly struct {
	fieldElements []fr.Element
}

// evalPolyFromElements creates a new evalPoly from field elements.
func evalPolyFromElements(elements []fr.Element) *evalPoly {
	return &evalPoly{fieldElements: elements}
}

// toCoeffPoly converts an evalPoly to a coeffPoly, using the IFFT operation
//
// blobLength (in SYMBOLS) is required, to be able to choose the correct parameters when performing FFT
func (p *evalPoly) toCoeffPoly(blobLength uint32) (*coeffPoly, error) {
	// TODO (litt3): this could conceivably be optimized, so that multiple objects share an instance of FFTSettings,
	//  which has enough roots of unity for general use. If the following construction of FFTSettings ever proves
	//  to present a computational burden, consider making this change.
	fftSettings := fft.FFTSettingsFromBlobLength(blobLength)

	// the FFT method pads to the next power of 2, so we don't need to do that manually
	ifftedElements, err := fftSettings.FFT(p.fieldElements, true)
	if err != nil {
		return nil, fmt.Errorf("perform IFFT: %w", err)
	}

	return coeffPolyFromElements(ifftedElements), nil
}

// toEncodedPayload converts an evalPoly into an encoded payload
//
// blobLength is required, to be able to perform length checks on the encoded payload during construction
func (p *evalPoly) toEncodedPayload(blobLength uint32) (*encodedPayload, error) {
	return encodedPayloadFromElements(p.fieldElements, blobLength)
}
