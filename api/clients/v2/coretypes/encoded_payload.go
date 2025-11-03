package coretypes

import (
	"encoding/binary"
	"fmt"
	gomath "math"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// EncodedPayload represents a payload that has had an encoding applied to it.
//
// It is an intermediary state between [Payload] and [Blob]. Most users should not need to interact with it directly,
// and should instead use [Payload.ToBlob] directly. EncodedPayloads are only exposed because secure rollup integrations
// with EigenDA need to decode them inside a fraud proof vm, in order to be able to discard wrongly encoded payloads.
// In such cases, a blob fetched from EigenDA can be transformed using [Blob.ToEncodedPayloadUnchecked]
// (note that this cannot error) and then sent to the fraud proof vm for verification. See
// https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#decode-blob-failed for more details.
//
// Example encoding:
//   - [Encoded Payload header (32 bytes total)] + [Encoded Payload Data (len is multiple of 32)]
//   - [0x00, version byte, big-endian uint32 len of payload, 0x00, ...] + [0x00, 31 bytes of data, 0x00, 31 bytes of data,...]
//
// An EncodedPayload can be interpreted as a polynomial, with each 32 byte chunk
// representing either a coefficient or an evaluation. Interpreting as coefficients has the advantage
// that the EncodedPayload already represents a [Blob]. Interpreting as evaluations has the advantage that
// point openings can be made (useful for interactive fraud proofs).
type EncodedPayload struct {
	// These bytes are kept private in order to force encapsulation, in case we decide in the future
	// to change the EncodedPayload's representation (eg. store [fr.Element]s directly instead).
	// Use [DeserializeEncodedPayloadUnchecked] to reconstruct a serialized EncodedPayload.
	//
	// The bytes should contain a power of 2 field elements, each 32 bytes long (see [EncodedPayload.checkLenInvariant]),
	// meaning valid lengths are [32, 64, 128, 256, ...]
	// The first 32 bytes represent the header (see [EncodedPayload.decodeHeader]),
	// and the body (bytes after header) contain serialized bn254 field elements.
	bytes []byte
}

// DeserializeEncodedPayloadUnchecked constructs an [EncodedPayload] from bytes array.
//
// It does not validate the bytes, to mimic the [Blob.ToEncodedPayloadUnchecked] process.
// The length, header, and body invariants are checked when calling [EncodedPayload.Decode].
func DeserializeEncodedPayloadUnchecked(bytes []byte) *EncodedPayload {
	return &EncodedPayload{bytes: bytes}
}

// Serialize returns the raw bytes of the encoded payload.
func (ep *EncodedPayload) Serialize() []byte {
	return ep.bytes
}

// LenSymbols returns the number of symbols in the encoded payload
func (ep *EncodedPayload) LenSymbols() uint32 {
	return uint32(len(ep.bytes)) / encoding.BYTES_PER_SYMBOL
}

// Decode applies the inverse of PayloadEncodingVersion0 to an EncodedPayload, and returns the decoded Payload
func (ep *EncodedPayload) Decode() (Payload, error) {
	err := ep.checkLenInvariant()
	if err != nil {
		return nil, fmt.Errorf("check length invariant: %w", err)
	}
	payloadLenInHeader, err := ep.decodeHeader()
	if err != nil {
		return nil, fmt.Errorf("decodeHeader: %w", err)
	}
	payload, err := ep.decodePayload(payloadLenInHeader)
	if err != nil {
		return nil, fmt.Errorf("decodePayload: %w", err)
	}
	return payload, nil
}

// ToBlob converts the EncodedPayload into a Blob
func (ep *EncodedPayload) ToBlob(payloadForm codecs.PolynomialForm) (*Blob, error) {
	if err := ep.checkLenInvariant(); err != nil {
		return nil, fmt.Errorf("check length invariant: %w", err)
	}
	fieldElements, err := rs.ToFrArray(ep.bytes)
	if err != nil {
		return nil, fmt.Errorf("encoded payload to field elements: %w", err)
	}

	var coeffPolynomial []fr.Element
	switch payloadForm {
	case codecs.PolynomialFormCoeff:
		// the payload is already in coefficient form. no conversion needs to take place, since blobs are also in
		// coefficient form
		coeffPolynomial = fieldElements
	case codecs.PolynomialFormEval:
		// the payload is in evaluation form, so we need to convert it to coeff form, since blobs are in coefficient form
		coeffPolynomial = evalToCoeffPoly(fieldElements)
	default:
		return nil, fmt.Errorf("unknown polynomial form: %v", payloadForm)
	}

	return blobFromCoefficients(coeffPolynomial)
}

// decodeHeader validates the header (first field element = 32 bytes) of the encoded payload,
// and returns the claimed length of the payload if the header is valid.
func (ep *EncodedPayload) decodeHeader() (uint32, error) {
	if len(ep.bytes) < codec.EncodedPayloadHeaderLenBytes {
		return 0, fmt.Errorf("encoded payload must be at least %d bytes long to contain a header, but got %d bytes",
			codec.EncodedPayloadHeaderLenBytes, len(ep.bytes))
	}
	if ep.bytes[0] != 0x00 {
		return 0, fmt.Errorf("encoded payload header first byte must be 0x00, but got %x", ep.bytes[0])
	}
	var payloadLength uint32
	switch ep.bytes[1] {
	case byte(codecs.PayloadEncodingVersion0):
		payloadLength = binary.BigEndian.Uint32(ep.bytes[2:6])
	default:
		return 0, fmt.Errorf("unknown encoded payload header version: %x", ep.bytes[1])
	}

	for _, b := range ep.bytes[6:codec.EncodedPayloadHeaderLenBytes] {
		if b != 0x00 {
			return 0, fmt.Errorf("padding in encoded payload header must be 0x00: %x", b)
		}
	}

	return payloadLength, nil
}

func (ep *EncodedPayload) decodePayload(payloadLen uint32) ([]byte, error) {
	body := ep.bytes[codec.EncodedPayloadHeaderLenBytes:]
	// Decode the body by removing internal 0 byte padding (0x00 initial byte for every 32 byte chunk)
	// The decodedBody should contain the payload bytes + potentially some external padding bytes.
	decodedPayload, err := codec.CheckAndRemoveInternalPadding(body, payloadLen)
	if err != nil {
		return nil, fmt.Errorf("remove internal padding: %w", err)
	}

	return Payload(decodedPayload), nil
}

// checkLenInvariant checks whether the encoded payload satisfies its length invariant.
// EncodedPayloads must contain a power of 2 number of Field Elements, each of length 32.
// This means the only valid encoded payloads have byte lengths of 32, 64, 128, 256, etc.
//
// Note that this function only checks the length invariant, meaning that it doesn't check that
// the 32 byte chunks are valid bn254 elements.
func (ep *EncodedPayload) checkLenInvariant() error {
	// this check is redundant since 0 is not a valid power of 32, but we keep it for clarity.
	if len(ep.bytes) < codec.EncodedPayloadHeaderLenBytes {
		return fmt.Errorf("encoded payload must be at least %d bytes long to contain a valid header, "+
			"but got %d bytes", codec.EncodedPayloadHeaderLenBytes, len(ep.bytes))
	}
	if len(ep.bytes)%encoding.BYTES_PER_SYMBOL != 0 {
		return fmt.Errorf("encoded payload must be a multiple of %d bytes (bn254 field element), "+
			"but got %d bytes", encoding.BYTES_PER_SYMBOL, len(ep.bytes))
	}
	// We could equivalently check that len(ep.bytes) is a power of 2 given that we've already
	// checked that it's a multiple of 32, but this invariant is closer to the representation of
	// the encoded payload as a polynomial, and is also more meaningful given
	// that the length in [encoding.BlobCommitments.Length] is in field elements.
	numFieldElements := len(ep.bytes) / encoding.BYTES_PER_SYMBOL
	if !math.IsPowerOfTwo(numFieldElements) {
		return fmt.Errorf("encoded payload must be a power of 2 field elements (32 bytes chunks), "+
			"but got %d field elements", numFieldElements)
	}
	return nil
}

// evalToCoeffPoly converts an evalPoly to a coeffPoly, using the IFFT operation
func evalToCoeffPoly(evalPoly []fr.Element) []fr.Element {
	// TODO (litt3): this could conceivably be optimized, so that multiple objects share an instance of FFTSettings,
	//  which has enough roots of unity for general use. If the following construction of FFTSettings ever proves
	//  to present a computational burden, consider making this change.
	fftSettings := fftSettingsFromBlobLengthSymbols(uint32(len(evalPoly)))

	// the FFT method pads to the next power of 2, so we don't need to do that manually
	ifftedElements, err := fftSettings.FFT(evalPoly, true)
	if err != nil {
		panic("bug: FFT only returns an error if we don't have enough roots of unity, " +
			"which is impossible because we already checked it above")
	}

	return ifftedElements
}

// fftSettingsFromBlobLengthSymbols accepts a blob length, and returns a new instance of FFT settings.
// blobLengthSymbols should be a power of 2, and the function will panic if it is not.
func fftSettingsFromBlobLengthSymbols(blobLengthSymbols uint32) *fft.FFTSettings {
	if !math.IsPowerOfTwo(blobLengthSymbols) {
		panic(fmt.Sprintf("blob length symbols %d is not a power of 2", blobLengthSymbols))
	}
	maxScale := uint8(gomath.Log2(float64(blobLengthSymbols)))
	return fft.NewFFTSettings(maxScale)
}
