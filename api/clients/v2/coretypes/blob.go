package coretypes

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Blob is data that is dispersed on eigenDA.
//
// A Blob is represented under the hood by an array of field elements (symbols),
// which represent a polynomial in coefficient form.
// A Blob must have a length (in symbols) that is a power of two. In particular, blobs of length 0 are not allowed.
// A Blob's length must match the blobLength in the BlobHeader's [encoding.BlobCommitments.Length].
type Blob struct {
	coeffPolynomial []fr.Element
}

// DeserializeBlob initializes a Blob from bytes.
// blobLengthSymbols is the length of the blob, which is present in the BlobHeader's [encoding.BlobCommitments.Length].
// The bytes passed in will be appended with zeros to match the blobLengthSymbols if they are shorter than that length,
// or an error will be returned if they are longer than that length.
func DeserializeBlob(bytes []byte, blobLengthSymbols uint32) (*Blob, error) {
	blobLengthBytes := blobLengthSymbols * encoding.BYTES_PER_SYMBOL
	if uint32(len(bytes)) > blobLengthBytes {
		return nil, fmt.Errorf(
			"length (%d bytes) is greater than claimed blob length (%d bytes)",
			len(bytes),
			blobLengthBytes)
	}
	// We pad with 0s up to blobLengthSymbols in case the bytes being deserialized have had trailing 0s truncated, as
	// illustrated by the following example: imagine a user disperses a very small blob, only 64 bytes, and the last 40
	// bytes are trailing zeros. When a different user fetches the blob from a relay, it's possible that the relay could
	// truncate the trailing zeros since that doesn't affect the KZG commitment. If we were to say that
	// blobLengthSymbols = nextPowerOf2(len(bytes)), then the user fetching and reconstructing this blob would determine
	// that the blob length is 1 symbol, when it's actually 2.
	if uint32(len(bytes)) < blobLengthBytes {
		bytes = append(bytes, make([]byte, blobLengthBytes-uint32(len(bytes)))...)
	}

	coeffPolynomial, err := rs.ToFrArray(bytes)
	if err != nil {
		return nil, fmt.Errorf("bytes to field elements: %w", err)
	}

	return blobFromCoefficients(coeffPolynomial)
}

// LenSymbols returns the number of coefficient symbols in the Blob.
func (b *Blob) LenSymbols() uint32 {
	return uint32(len(b.coeffPolynomial))
}

// LenBytes returns the number of bytes in the Blob.
func (b *Blob) LenBytes() uint32 {
	return uint32(len(b.coeffPolynomial) * encoding.BYTES_PER_SYMBOL)
}

// Serialize gets the raw bytes of the Blob
func (b *Blob) Serialize() []byte {
	return rs.SerializeFieldElements(b.coeffPolynomial)
}

// ToPayload converts the Blob into a Payload
//
// The payloadForm indicates how payloads are interpreted. The way that payloads are interpreted dictates what
// conversion, if any, must be performed when creating a payload from the blob.
func (b *Blob) ToPayload(payloadForm codecs.PolynomialForm) (Payload, error) {
	encodedPayload := b.ToEncodedPayloadUnchecked(payloadForm)

	payload, err := encodedPayload.Decode()
	if err != nil {
		return Payload{}, fmt.Errorf("decode payload: %w", err)
	}
	return payload, nil
}

// ToEncodedPayloadUnchecked creates an EncodedPayload from the blob.
//
// This method does not perform any validation on the blob or the resulting EncodedPayload.
// Most users should call [Blob.ToPayload] directly instead, but this method is exposed
// since some secure integrations require decoding the Payload (and checking the invariants)
// inside a fraud proof VM.
//
// The payloadForm indicates how payloads are interpreted. The way that payloads are interpreted dictates what
// conversion, if any, must be performed when creating an encoded payload from the blob.
func (b *Blob) ToEncodedPayloadUnchecked(payloadForm codecs.PolynomialForm) *EncodedPayload {
	var encodedPayloadElements []fr.Element
	switch payloadForm {
	case codecs.PolynomialFormCoeff:
		// the payload is interpreted as coefficients of the polynomial, so no conversion needs to be done, given that
		// eigenda also interprets blobs as coefficients
		encodedPayloadElements = b.coeffPolynomial
	case codecs.PolynomialFormEval:
		// the payload is interpreted as evaluations of the polynomial, so the coefficient representation contained
		// in the blob must be converted to the evaluation form
		encodedPayloadElements = b.toEvalPoly()
	default:
		panic(fmt.Sprintf("invalid codecs.PolynomialForm enum value: %d", payloadForm))
	}

	return DeserializeEncodedPayloadUnchecked(rs.SerializeFieldElements(encodedPayloadElements))

}

// toEvalPoly converts a blob's coeffPoly to an evalPoly, using the FFT operation
func (b *Blob) toEvalPoly() []fr.Element {
	// TODO (litt3): this could conceivably be optimized, so that multiple objects share an instance of FFTSettings,
	//  which has enough roots of unity for general use. If the following construction of FFTSettings ever proves
	//  to present a computational burden, consider making this change.
	fftSettings := fftSettingsFromBlobLengthSymbols(uint32(len(b.coeffPolynomial)))

	// the FFT method pads to the next power of 2, so we don't need to do that manually
	fftedElements, err := fftSettings.FFT(b.coeffPolynomial, false)
	if err != nil {
		panic("bug: FFT only returns an error if we don't have enough roots of unity, " +
			"which is impossible because we already checked it above")
	}
	return fftedElements
}

// blobFromCoefficients creates a blob from the coefficients of a polynomial.
// A Blob must have a power of two coefficients. Thus:
//   - If the passed coefficients are a power of 2 in length, they will be used as is.
//   - If the coefficients are not a power of two in length, they will be copied to a new slice that is padded with zeros
//     to the next power of two in length.
func blobFromCoefficients(coefficients []fr.Element) (*Blob, error) {
	if !encoding.IsPowerOfTwo(len(coefficients)) {
		return nil, fmt.Errorf("blob must have a power of 2 coefficients, but got %d coefficients", len(coefficients))
	}
	return &Blob{
		coeffPolynomial: coefficients,
	}, nil
}
