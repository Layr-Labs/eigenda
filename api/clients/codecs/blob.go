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
// A Blob is represented under the hood by an array of field elements, which represent a polynomial in coefficient form
type Blob struct {
	coeffPolynomial []fr.Element
	// blobLengthSymbols must be a power of 2, and should match the blobLength claimed in the BlobCommitment
	//
	// This value must be specified, rather than computed from the length of the coeffPolynomial, due to an edge case
	// illustrated by the following example: imagine a user disperses a very small blob, only 64 bytes, and the last 40
	// bytes are trailing zeros. When a different user fetches the blob from a relay, it's possible that the relay could
	// truncate the trailing zeros. If we were to say that blobLengthSymbols = nextPowerOf2(len(coeffPolynomial)), then the
	// user fetching and reconstructing this blob would determine that the blob length is 1 symbol, when it's actually 2.
	blobLengthSymbols uint32
}

// BlobFromBytes initializes a Blob from bytes
func BlobFromBytes(bytes []byte, blobLengthSymbols uint32) (*Blob, error) {
	coeffPolynomial, err := rs.ToFrArray(bytes)
	if err != nil {
		return nil, fmt.Errorf("bytes to field elements: %w", err)
	}

	return BlobFromPolynomial(coeffPolynomial, blobLengthSymbols)
}

// BlobFromPolynomial initializes a blob from a polynomial
func BlobFromPolynomial(coeffPolynomial []fr.Element, blobLengthSymbols uint32) (*Blob, error) {
	return &Blob{
		coeffPolynomial:   coeffPolynomial,
		blobLengthSymbols: blobLengthSymbols,
	}, nil
}

// GetBytes gets the raw bytes of the Blob
func (b *Blob) GetBytes() []byte {
	return rs.FieldElementsToBytes(b.coeffPolynomial)
}

// ToPayload converts the Blob into a Payload
//
// The payloadForm indicates how payloads are interpreted. The way that payloads are interpreted dictates what
// conversion, if any, must be performed when creating a payload from the blob.
func (b *Blob) ToPayload(payloadForm PolynomialForm) (*Payload, error) {
	encodedPayload, err := b.toEncodedPayload(payloadForm)
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
//
// The payloadForm indicates how payloads are interpreted. The way that payloads are interpreted dictates what
// conversion, if any, must be performed when creating an encoded payload from the blob.
func (b *Blob) toEncodedPayload(payloadForm PolynomialForm) (*encodedPayload, error) {
	maxPermissiblePayloadLength, err := codec.GetMaxPermissiblePayloadLength(b.blobLengthSymbols)
	if err != nil {
		return nil, fmt.Errorf("get max permissible payload length: %w", err)
	}

	var payloadElements []fr.Element
	switch payloadForm {
	case PolynomialFormCoeff:
		// the payload is interpreted as coefficients of the polynomial, so no conversion needs to be done, given that
		// eigenda also interprets blobs as coefficients
		payloadElements = b.coeffPolynomial
	case PolynomialFormEval:
		// the payload is interpreted as evaluations of the polynomial, so the coefficient representation contained
		// in the blob must be converted to the evaluation form
		payloadElements, err = b.computeEvalPoly()
		if err != nil {
			return nil, fmt.Errorf("compute eval poly: %w", err)
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
	fftSettings := fft.FFTSettingsFromBlobLengthSymbols(b.blobLengthSymbols)

	// the FFT method pads to the next power of 2, so we don't need to do that manually
	fftedElements, err := fftSettings.FFT(b.coeffPolynomial, false)
	if err != nil {
		return nil, fmt.Errorf("perform FFT: %w", err)
	}

	return fftedElements, nil
}
