package coretypes

import (
	"encoding/binary"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
)

// Payload represents arbitrary user data, without any processing.
type Payload []byte

// ToEncodedPayload performs the [codecs.PayloadEncodingVersion0] encoding to create an encoded payload.
func (p Payload) ToEncodedPayload() *EncodedPayload {
	// Encode payload modulo bn254, and align to 32 bytes
	encodedData := codec.PadPayload(p)

	// Calculate the length of the EncodedPayload in symbols (including the header) which has to be a power of 2.
	encodedDataLenSymbols := uint32(len(encodedData)) / encoding.BYTES_PER_SYMBOL
	encodedHeaderAndDataLenSymbols := codec.EncodedPayloadHeaderLenSymbols + encodedDataLenSymbols
	encodedPayloadLenSymbols := encoding.NextPowerOf2(encodedHeaderAndDataLenSymbols)

	encodedPayloadBytes := make([]byte, encodedPayloadLenSymbols*encoding.BYTES_PER_SYMBOL)

	// Write the header
	encodedPayloadHeader := encodedPayloadBytes[:codec.EncodedPayloadHeaderLenBytes]
	// first byte is always 0 to ensure the payloadHeader is a valid bn254 element
	encodedPayloadHeader[1] = byte(codecs.PayloadEncodingVersion0) // encode version byte
	// encode payload length as uint32
	binary.BigEndian.PutUint32(
		encodedPayloadHeader[2:6],
		uint32(len(p))) // uint32 should be more than enough to store the length (approx 4gb)

	// Write the encoded data, starting after the header
	copy(encodedPayloadBytes[codec.EncodedPayloadHeaderLenBytes:], encodedData)

	encodedPayload := DeserializeEncodedPayloadUnchecked(encodedPayloadBytes)
	if err := encodedPayload.checkLenInvariant(); err != nil {
		panic("bug converting payload to encodedPayload: broken EncodedPayload invariants:" + err.Error())
	}
	return encodedPayload
}

// ToBlob converts the Payload bytes into a Blob
//
// The payloadForm indicates how payloads are interpreted. The form of a payload dictates what conversion, if any, must
// be performed when creating a blob from the payload.
func (p Payload) ToBlob(payloadForm codecs.PolynomialForm) (*Blob, error) {
	return p.ToEncodedPayload().ToBlob(payloadForm)
}
