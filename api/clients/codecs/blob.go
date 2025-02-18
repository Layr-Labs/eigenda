package codecs

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Blob is data that is dispersed on eigenDA.
//
// A Blob is represented under the hood by a coeff polynomial
type Blob struct {
	coeffPolynomial []fr.Element
	// blobLength must be a power of 2, and should match the blobLength claimed in the BlobCommitment
	// This is the blob length IN SYMBOLS, not in bytes
	blobLength uint32
}

// BlobFromBytes initializes a Blob from bytes
//
// blobLength is the length of the blob IN SYMBOLS
func BlobFromBytes(bytes []byte, blobLength uint32) (*Blob, error) {
	coeffPolynomial, err := rs.ToFrArray(bytes)
	if err != nil {
		return nil, fmt.Errorf("bytes to field elements: %w", err)
	}

	return BlobFromPolynomial(coeffPolynomial, blobLength)
}

// BlobFromPolynomial initializes a blob from a polynomial
//
// blobLength is the length of the blob IN SYMBOLS
func BlobFromPolynomial(coeffPolynomial []fr.Element, blobLength uint32) (*Blob, error) {
	return &Blob{
		coeffPolynomial: coeffPolynomial,
		blobLength:      blobLength,
	}, nil
}

// GetBytes gets the raw bytes of the Blob
func (b *Blob) GetBytes() []byte {
	return rs.FieldElementsToBytes(b.coeffPolynomial)
}

// ToPayload converts the Blob into a Payload
//
// The payloadStartingForm indicates how payloads are constructed by the dispersing client. Based on the starting form
// of the payload, we can determine what operations must be done to the blob in order to reconstruct the original payload
func (b *Blob) ToPayload(payloadStartingForm PolynomialForm) (*Payload, error) {
	encodedPayload, err := b.toEncodedPayload(payloadStartingForm)
	if err != nil {
		return nil, fmt.Errorf("to encoded payload: %w", err)
	}

	payload, err := encodedPayload.decode()
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}

	return payload, nil
}

// toEncodedPayload creates an encodedPayload from the blob
func (b *Blob) toEncodedPayload(payloadStartingForm PolynomialForm) (*encodedPayload, error) {
	maxPermissiblePayloadLength, err := codec.GetMaxPermissiblePayloadLength(b.blobLength)
	if err != nil {
		return nil, fmt.Errorf("get max permissible payload length: %w", err)
	}

	var payloadElements []fr.Element
	switch payloadStartingForm {
	case PolynomialFormCoeff:
		// the payload started off in coefficient form, so no conversion needs to be done
		payloadElements = b.coeffPolynomial
	case PolynomialFormEval:
		// the payload started off in evaluation form, so we first need to convert the blob's coeff poly into an eval poly
		payloadElements, err = b.computeEvalPoly()
		if err != nil {
			return nil, fmt.Errorf("coeff poly to eval poly: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid polynomial form")
	}

	encodedPayload, err := encodedPayloadFromElements(payloadElements, maxPermissiblePayloadLength)
	if err != nil {
		return nil, fmt.Errorf("encoded payload from elements %w", err)
	}

	return encodedPayload, nil
}

// computeEvalPoly converts a blob's coeffPoly to an evalPoly, using the FFT operation
func (b *Blob) computeEvalPoly() ([]fr.Element, error) {
	// TODO (litt3): this could conceivably be optimized, so that multiple objects share an instance of FFTSettings,
	//  which has enough roots of unity for general use. If the following construction of FFTSettings ever proves
	//  to present a computational burden, consider making this change.
	fftSettings := fft.FFTSettingsFromBlobLength(b.blobLength)

	// the FFT method pads to the next power of 2, so we don't need to do that manually
	fftedElements, err := fftSettings.FFT(b.coeffPolynomial, false)
	if err != nil {
		return nil, fmt.Errorf("perform FFT: %w", err)
	}

	return fftedElements, nil
}