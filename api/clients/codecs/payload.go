package codecs

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Payload represents arbitrary user data, without any processing.
type Payload struct {
	bytes []byte
}

// NewPayload wraps an arbitrary array of bytes into a Payload type.
func NewPayload(payloadBytes []byte) *Payload {
	return &Payload{
		bytes: payloadBytes,
	}
}

// ToBlob converts the Payload bytes into a Blob
//
// The payloadForm indicates how payloads are interpreted. The form of a payload dictates what conversion, if any, must
// be performed when creating a blob from the payload.
func (p *Payload) ToBlob(payloadForm PolynomialForm) (*Blob, error) {
	encodedPayload, err := newEncodedPayload(p)
	if err != nil {
		return nil, fmt.Errorf("encoding payload: %w", err)
	}

	fieldElements, err := encodedPayload.toFieldElements()
	if err != nil {
		return nil, fmt.Errorf("encoded payload to field elements: %w", err)
	}

	blobLengthSymbols := uint32(encoding.NextPowerOf2(len(fieldElements)))

	var coeffPolynomial []fr.Element
	switch payloadForm {
	case PolynomialFormCoeff:
		// the payload is already in coefficient form. no conversion needs to take place, since blobs are also in
		// coefficient form
		coeffPolynomial = fieldElements
	case PolynomialFormEval:
		// the payload is in evaluation form, so we need to convert it to coeff form, since blobs are in coefficient form
		coeffPolynomial, err = evalToCoeffPoly(fieldElements, blobLengthSymbols)
		if err != nil {
			return nil, fmt.Errorf("eval poly to coeff poly: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown polynomial form: %v", payloadForm)
	}

	return BlobFromPolynomial(coeffPolynomial, blobLengthSymbols)
}

// Serialize returns the bytes that underlie the payload, i.e. the unprocessed user data
func (p *Payload) Serialize() []byte {
	return p.bytes
}

// evalToCoeffPoly converts an evalPoly to a coeffPoly, using the IFFT operation
//
// blobLengthSymbols is required, to be able to choose the correct parameters when performing FFT
func evalToCoeffPoly(evalPoly []fr.Element, blobLengthSymbols uint32) ([]fr.Element, error) {
	// TODO (litt3): this could conceivably be optimized, so that multiple objects share an instance of FFTSettings,
	//  which has enough roots of unity for general use. If the following construction of FFTSettings ever proves
	//  to present a computational burden, consider making this change.
	fftSettings := fft.FFTSettingsFromBlobLengthSymbols(blobLengthSymbols)

	// the FFT method pads to the next power of 2, so we don't need to do that manually
	ifftedElements, err := fftSettings.FFT(evalPoly, true)
	if err != nil {
		return nil, fmt.Errorf("perform IFFT: %w", err)
	}

	return ifftedElements, nil
}
