package codecs

import (
	"fmt"
)

// Blob is data that is dispersed on eigenDA.
//
// A Blob will contain either an encodedPayload, or a coeffPoly. Whether the Blob contains the former or the latter
// is determined by how the dispersing client has been configured.
type Blob struct {
	encodedPayload *encodedPayload
	coeffPoly      *coeffPoly
}

// BlobFromEncodedPayload creates a Blob containing an encodedPayload
func blobFromEncodedPayload(encodedPayload *encodedPayload) *Blob {
	return &Blob{encodedPayload: encodedPayload}
}

// blobFromCoeffPoly creates a Blob containing a coeffPoly
func blobFromCoeffPoly(poly *coeffPoly) *Blob {
	return &Blob{coeffPoly: poly}
}

// NewBlob initializes a Blob from raw bytes, and the expected BlobForm
//
// This function will return an error if the input bytes cannot be successfully interpreted as the claimed BlobForm
func NewBlob(bytes []byte, blobForm BlobForm) (*Blob, error) {
	switch blobForm {
	case Eval:
		encodedPayload, err := newEncodedPayload(bytes)
		if err != nil {
			return nil, fmt.Errorf("new encoded payload: %v", err)
		}

		return blobFromEncodedPayload(encodedPayload), nil
	case Coeff:
		coeffPoly, err := coeffPolyFromBytes(bytes)
		if err != nil {
			return nil, fmt.Errorf("new coeff poly: %v", err)
		}

		return blobFromCoeffPoly(coeffPoly), nil
	default:
		return nil, fmt.Errorf("unsupported blob form type: %v", blobForm)
	}
}

// GetBytes gets the raw bytes of the Blob
func (b *Blob) GetBytes() []byte {
	if b.encodedPayload == nil {
		return b.encodedPayload.getBytes()
	} else {
		return b.coeffPoly.getBytes()
	}
}

// ToPayload converts the Blob into a Payload
func (b *Blob) ToPayload() (*Payload, error) {
	var encodedPayload *encodedPayload
	var err error
	if b.encodedPayload != nil {
		encodedPayload = b.encodedPayload
	} else if b.coeffPoly != nil {
		evalPoly, err := b.coeffPoly.toEvalPoly()
		if err != nil {
			return nil, fmt.Errorf("coeff poly to eval poly: %v", err)
		}

		encodedPayload, err = evalPoly.toEncodedPayload()
		if err != nil {
			return nil, fmt.Errorf("eval poly to encoded payload: %v", err)
		}
	} else {
		return nil, fmt.Errorf("blob has no contents")
	}

	payload, err := encodedPayload.decode()
	if err != nil {
		return nil, fmt.Errorf("decode encoded payload: %v", err)
	}

	return payload, nil
}
