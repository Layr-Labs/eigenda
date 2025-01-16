package codecs

import (
	"fmt"
)

// ProtoBlob is data that is dispersed to eigenDA. TODO: write a good description of a proto blob
//
// A ProtoBlob will contain either an encodedPayload, or a coeffPoly. Whether the ProtoBlob contains the former or the
// latter is determined by how the dispersing client has been configured.
type ProtoBlob struct {
	encodedPayload *encodedPayload
	coeffPoly      *coeffPoly
}

// protoBlobFromEncodedPayload creates a ProtoBlob containing an encodedPayload
func protoBlobFromEncodedPayload(encodedPayload *encodedPayload) *ProtoBlob {
	return &ProtoBlob{encodedPayload: encodedPayload}
}

// blobFromCoeffPoly creates a ProtoBlob containing a coeffPoly
func protoBlobFromCoeffPoly(poly *coeffPoly) *ProtoBlob {
	return &ProtoBlob{coeffPoly: poly}
}

// NewProtoBlob initializes a ProtoBlob from raw bytes, and the expected BlobForm
//
// This function will return an error if the input bytes cannot be successfully interpreted as the claimed BlobForm
func NewProtoBlob(bytes []byte, blobForm BlobForm) (*ProtoBlob, error) {
	switch blobForm {
	case Eval:
		encodedPayload, err := newEncodedPayload(bytes)
		if err != nil {
			return nil, fmt.Errorf("new encoded payload: %v", err)
		}

		return protoBlobFromEncodedPayload(encodedPayload), nil
	case Coeff:
		coeffPoly, err := coeffPolyFromBytes(bytes)
		if err != nil {
			return nil, fmt.Errorf("new coeff poly: %v", err)
		}

		return protoBlobFromCoeffPoly(coeffPoly), nil
	default:
		return nil, fmt.Errorf("unsupported blob form type: %v", blobForm)
	}
}

// GetBytes gets the raw bytes of the Blob
func (b *ProtoBlob) GetBytes() []byte {
	if b.encodedPayload == nil {
		return b.encodedPayload.getBytes()
	} else {
		return b.coeffPoly.getBytes()
	}
}

// ToPayload converts the ProtoBlob into a Payload
func (b *ProtoBlob) ToPayload() (*Payload, error) {
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
		return nil, fmt.Errorf("proto blob has no contents")
	}

	payload, err := encodedPayload.decode()
	if err != nil {
		return nil, fmt.Errorf("decode encoded payload: %v", err)
	}

	return payload, nil
}
