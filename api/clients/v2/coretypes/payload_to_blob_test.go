// nolint: lll // Example contains long lines to print output
package coretypes

import (
	"encoding/hex"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/encoding"
)

// Example demonstrating the conversion process from a payload, to an encodedPayload interpreted as
// evaluations of a polynomial, which is then IFFT'd to produce a Blob in coefficient form.
// This example demonstrates the process that [Payload.ToBlob] performs internally.
func Example_payloadToBlobConversion() {
	// We create a payload of 2 symbols (64 bytes), which with an EncodedPayloadHeader of 1 symbol (32 bytes),
	// will result in an encoded payload of 3 symbols (96 bytes). Because blobs have to be powers of 2,
	// the blob length will be 4 symbols (128 bytes).
	numSymbols := uint64(2)
	payloadBytesPerSymbols := uint64(encoding.BYTES_PER_SYMBOL - 1)
	payloadBytes := make([]byte, numSymbols*payloadBytesPerSymbols)
	for i := range numSymbols {
		payloadBytes[i*payloadBytesPerSymbols] = byte(i + 1)
	}
	payload := Payload(payloadBytes)
	fmt.Printf("Payload bytes (len %d):\n%s\n", len(payload), hex.EncodeToString(payload))

	encodedPayload := payload.ToEncodedPayload()
	fmt.Printf("Encoded Payload bytes (len %d):\n%s\n", len(encodedPayload.Serialize()), hex.EncodeToString(encodedPayload.Serialize()))

	// Change to [codecs.PolynomialFormCoeff] to see that blob will have the same bytes as the encoded payload.
	blob, err := encodedPayload.ToBlob(codecs.PolynomialFormEval)
	if err != nil {
		panic(err)
	}
	// Now we have a Blob that can be serialized and dispersed on eigenDA.
	blobBytes := blob.Serialize()
	fmt.Printf("Blob bytes (len %d):\n%s\n", len(blobBytes), hex.EncodeToString(blobBytes))

	// Output:
	// Payload bytes (len 62):
	// 0100000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000
	// Encoded Payload bytes (len 128):
	// 00000000003e0000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
	// Blob bytes (len 128):
	// 0000c000000f80000000000000000000000000000000000000000000000000000b51701f1769982df83a9dbe76a1a7ac21abbab2ec7461a00b07d4200db2ec4900004000000f80000000000000000000000000000000000000000000000000002511de53c9e707fbc015a7f80adfb0b106882d958d450ef138da2173e24d13b8
}
