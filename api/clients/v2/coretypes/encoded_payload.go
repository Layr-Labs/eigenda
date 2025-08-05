package coretypes

import (
	"encoding/binary"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// EncodedPayload represents a payload that has had an encoding applied to it,
// meaning that it is guaranteed to be a multiple of 32 bytes in length,
// each such 32 bytes representing a bn254 field element.
//
// Example encoding:
//   - [Encoded Payload header (32 bytes total)] + [Encoded Payload Data (len is multiple of 32)]
//   - [0x00, version byte, big-endian uint32 len of payload, 0x00, ...] + [0x00, 31 bytes of data, 0x00, 31 bytes of data,...]
type EncodedPayload struct {
	// the size of these bytes is guaranteed to be a multiple of 32
	bytes []byte
}

// Serialize returns the raw bytes of the encoded payload
func (ep *EncodedPayload) Serialize() []byte {
	return ep.bytes
}

// LenSymbols returns the number of symbols in the encoded payload
func (ep *EncodedPayload) LenSymbols() uint32 {
	return uint32(len(ep.bytes)+encoding.BYTES_PER_SYMBOL-1) / encoding.BYTES_PER_SYMBOL
}

// Decode applies the inverse of PayloadEncodingVersion0 to an EncodedPayload, and returns the decoded Payload
func (ep *EncodedPayload) Decode() (Payload, error) {
	if len(ep.bytes) < codec.PayloadHeaderSizeBytes {
		return nil, fmt.Errorf("encoded payload must be at least 32 bytes long, but got %d bytes", len(ep.bytes))
	}
	if ep.bytes[0] != 0x00 {
		return nil, fmt.Errorf("encoded payload header first byte must be 0x00, but got %x", ep.bytes[0])
	}
	if ep.bytes[1] != byte(codecs.PayloadEncodingVersion0) {
		return nil, fmt.Errorf("encoded payload header version byte must be %x, but got %x", codecs.PayloadEncodingVersion0, ep.bytes[1])
	}
	claimedLength := binary.BigEndian.Uint32(ep.bytes[2:6])

	// decode raw data modulo bn254
	unpaddedData, err := codec.RemoveInternalPadding(ep.bytes[32:])
	if err != nil {
		return nil, fmt.Errorf("remove internal padding: %w", err)
	}

	unpaddedDataLength := uint32(len(unpaddedData))

	// data length is checked when constructing an encoded payload. If this error is encountered, that means there
	// must be a flaw in the logic at construction time (or someone was bad and didn't use the proper construction methods)
	if unpaddedDataLength < claimedLength {
		return nil, fmt.Errorf(
			"length of unpadded data %d is less than length claimed in encoded payload header %d. this should never happen",
			unpaddedDataLength, claimedLength)
	}

	return Payload(unpaddedData[0:claimedLength]), nil
}

// ToBlob converts the EncodedPayload into a Blob
func (ep *EncodedPayload) ToBlob(payloadForm codecs.PolynomialForm) (*Blob, error) {
	fieldElements, err := rs.ToFrArray(ep.bytes)
	if err != nil {
		return nil, fmt.Errorf("encoded payload to field elements: %w", err)
	}

	blobLengthSymbols := uint32(encoding.NextPowerOf2(len(fieldElements)))

	var coeffPolynomial []fr.Element
	switch payloadForm {
	case codecs.PolynomialFormCoeff:
		// the payload is already in coefficient form. no conversion needs to take place, since blobs are also in
		// coefficient form
		coeffPolynomial = fieldElements
	case codecs.PolynomialFormEval:
		// the payload is in evaluation form, so we need to convert it to coeff form, since blobs are in coefficient form
		coeffPolynomial, err = evalToCoeffPoly(fieldElements, blobLengthSymbols)
		if err != nil {
			return nil, fmt.Errorf("eval poly to coeff poly: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown polynomial form: %v", payloadForm)
	}

	return BlobFromCoefficients(coeffPolynomial), nil
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
