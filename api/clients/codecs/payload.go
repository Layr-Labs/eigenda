package codecs

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
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
func (p *Payload) ToBlob(form PolynomialForm) (*Blob, error) {
	encodedPayload, err := newEncodedPayload(p)
	if err != nil {
		return nil, fmt.Errorf("encoding payload: %w", err)
	}

	var coeffPolynomial *coeffPoly
	switch form {
	case PolynomialFormCoeff:
		// the payload is already in coefficient form. no conversion needs to take place, since blobs are also in
		// coefficient form
		coeffPolynomial, err = encodedPayload.toCoeffPoly()
		if err != nil {
			return nil, fmt.Errorf("coeff poly from elements: %w", err)
		}
	case PolynomialFormEval:
		// the payload is in evaluation form, so we need to convert it to coeff form
		evalPoly, err := encodedPayload.toEvalPoly()
		if err != nil {
			return nil, fmt.Errorf("eval poly from elements: %w", err)
		}

		coeffPolynomial, err = evalPoly.toCoeffPoly()
		if err != nil {
			return nil, fmt.Errorf("eval poly to coeff poly: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown polynomial form: %v", form)
	}

	// it's possible that the number of field elements might already be a power of 2
	// in that case, calling NextPowerOf2 will just return the input value
	blobLength := uint32(encoding.NextPowerOf2(len(coeffPolynomial.fieldElements)))

	return BlobFromPolynomial(coeffPolynomial, blobLength)
}

// GetBytes returns the bytes that underlie the payload, i.e. the unprocessed user data
func (p *Payload) GetBytes() []byte {
	return p.bytes
}
