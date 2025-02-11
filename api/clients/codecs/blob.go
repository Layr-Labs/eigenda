package codecs

import (
	"fmt"
)

// Blob is data that is dispersed on eigenDA.
//
// A Blob is represented under the hood by a coeff polynomial
type Blob struct {
	coeffPolynomial *coeffPoly
	// blobLength must be a power of 2, and should match the blobLength claimed in the BlobCommitment
	// This is the blob length in symbols, NOT in bytes
	blobLength uint32
}

// BlobFromBytes initializes a Blob from bytes, and a blobLength in symbols
func BlobFromBytes(bytes []byte, blobLength uint32) (*Blob, error) {
	poly, err := coeffPolyFromBytes(bytes)
	if err != nil {
		return nil, fmt.Errorf("polynomial from bytes: %w", err)
	}

	return BlobFromPolynomial(poly, blobLength)
}

// BlobFromPolynomial initializes a blob from a polynomial, and a blobLength in symbols
func BlobFromPolynomial(coeffPolynomial *coeffPoly, blobLength uint32) (*Blob, error) {
	return &Blob{
		coeffPolynomial: coeffPolynomial,
		blobLength:      blobLength}, nil
}

// GetBytes gets the raw bytes of the Blob
func (b *Blob) GetBytes() []byte {
	return b.coeffPolynomial.getBytes()
}

// ToPayload converts the Blob into a Payload
//
// The payloadStartingForm indicates how payloads are constructed by the dispersing client. Based on the starting form
// of the payload, we can determine what operations must be done to the blob in order to reconstruct the original payload
func (b *Blob) ToPayload(payloadStartingForm PolynomialForm) (*Payload, error) {
	var encodedPayload *encodedPayload
	var err error
	switch payloadStartingForm {
	case PolynomialFormCoeff:
		// the payload started off in coefficient form, so no conversion needs to be done
		encodedPayload, err = b.coeffPolynomial.toEncodedPayload(b.blobLength)
		if err != nil {
			return nil, fmt.Errorf("coeff poly to encoded payload: %w", err)
		}
	case PolynomialFormEval:
		// the payload started off in evaluation form, so we first need to convert the blob's coeff poly into an eval poly
		evalPoly, err := b.coeffPolynomial.toEvalPoly()
		if err != nil {
			return nil, fmt.Errorf("coeff poly to eval poly: %w", err)
		}

		encodedPayload, err = evalPoly.toEncodedPayload(b.blobLength)
		if err != nil {
			return nil, fmt.Errorf("eval poly to encoded payload: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid polynomial form")
	}

	payload, err := encodedPayload.decode()
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}

	return payload, nil
}
