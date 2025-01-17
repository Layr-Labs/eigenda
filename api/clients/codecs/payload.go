package codecs

import (
	"encoding/binary"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
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

// Encode applies the PayloadEncodingVersion0 to the original payload bytes
//
// Example encoding:
//
//                  Payload header (32 bytes total)                                  Encoded Payload Data
// [0x00, version byte, big-endian uint32 len of payload, 0x00, ...] + [0x00, 31 bytes of data, 0x00, 31 bytes of data,...]
func (p *Payload) encode() (*encodedPayload, error) {
	payloadHeader := make([]byte, 32)
	// first byte is always 0 to ensure the payloadHeader is a valid bn254 element
	payloadHeader[1] = byte(PayloadEncodingVersion0) // encode version byte

	// encode payload length as uint32
	binary.BigEndian.PutUint32(
		payloadHeader[2:6],
		uint32(len(p.bytes))) // uint32 should be more than enough to store the length (approx 4gb)

	// encode payload modulo bn254, and align to 32 bytes
	encodedData := codec.PadPayload(p.bytes)

	encodedPayload, err := newEncodedPayload(append(payloadHeader, encodedData...))
	if err != nil {
		return nil, fmt.Errorf("encoding payload: %w", err)
	}

	return encodedPayload, nil
}

// ToBlob converts the Payload bytes into a Blob
func (p *Payload) ToBlob(form BlobForm) (*Blob, error) {
	encodedPayload, err := p.encode()
	if err != nil {
		return nil, fmt.Errorf("encoding payload: %w", err)
	}

	switch form {
	case Eval:
		return blobFromEncodedPayload(encodedPayload), nil
	case Coeff:
		evalPolynomial, err := encodedPayload.toEvalPoly()
		if err != nil {
			return nil, fmt.Errorf("encoded payload to eval poly: %w", err)
		}

		coeffPoly, err := evalPolynomial.toCoeffPoly()
		if err != nil {
			return nil, fmt.Errorf("eval poly to coeff poly: %w", err)
		}

		return blobFromCoeffPoly(coeffPoly), nil
	default:
		return nil, fmt.Errorf("unknown polynomial form: %v", form)
	}
}

// GetBytes returns the bytes that underlie the payload, i.e. the unprocessed user data
func (p *Payload) GetBytes() []byte {
	return p.bytes
}
