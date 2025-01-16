package codecs

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// evalPoly is a polynomial in evaluation form.
//
// The underlying bytes represent 32 byte field elements, and the field elements represent the evaluation at the
// polynomial's roots of unity
type evalPoly struct {
	fieldElements []fr.Element
}

// evalPolyFromBytes creates a new evalPoly from bytes. This function performs the necessary checks to guarantee that the
// bytes are well-formed, and returns a new object if they are
func evalPolyFromBytes(bytes []byte) (*evalPoly, error) {
	paddedBytes := encoding.PadToPowerOfTwo(bytes)

	fieldElements, err := rs.BytesToFieldElements(paddedBytes)
	if err != nil {
		return nil, fmt.Errorf("deserialize field elements: %w", err)
	}

	return &evalPoly{fieldElements: fieldElements}, nil
}

// evalPolyFromElements creates a new evalPoly from field elements.
func evalPolyFromElements(fieldElements []fr.Element) (*evalPoly, error) {
	return &evalPoly{fieldElements: fieldElements}, nil
}

// toCoeffPoly converts an evalPoly to a coeffPoly, using the IFFT operation
func (ep *evalPoly) toCoeffPoly() (*coeffPoly, error) {
	maxScale := uint8(math.Log2(float64(len(ep.fieldElements))))
	ifftedElements, err := fft.NewFFTSettings(maxScale).FFT(ep.fieldElements, true)
	if err != nil {
		return nil, fmt.Errorf("perform IFFT: %w", err)
	}

	coeffPoly, err := coeffPolyFromElements(ifftedElements)
	if err != nil {
		return nil, fmt.Errorf("construct coeff poly: %w", err)
	}

	return coeffPoly, nil
}

// toEncodedPayload converts an evalPoly into an encoded payload
//
// This conversion entails removing the power-of-2 padding which is added to an encodedPayload when originally creating
// an evalPoly.
func (ep *evalPoly) toEncodedPayload() (*encodedPayload, error) {
	polynomialBytes := rs.FieldElementsToBytes(ep.fieldElements)

	payloadLength := binary.BigEndian.Uint32(polynomialBytes[2:6])

	// add 32 to the padded data length, since the encoded payload includes a payload header
	encodedPayloadLength := codec.GetPaddedDataLength(payloadLength) + 32

	if uint32(len(polynomialBytes)) < payloadLength {
		return nil, fmt.Errorf(
			"polynomial contains fewer bytes (%d) than expected encoded payload (%d), as determined by claimed length in payload header (%d)",
			len(polynomialBytes), encodedPayloadLength, payloadLength)
	}

	encodedPayloadBytes := make([]byte, encodedPayloadLength)
	copy(encodedPayloadBytes, polynomialBytes[:encodedPayloadLength])

	encodedPayload, err := newEncodedPayload(encodedPayloadBytes)
	if err != nil {
		return nil, fmt.Errorf("construct encoded payload: %w", err)
	}

	return encodedPayload, nil
}
