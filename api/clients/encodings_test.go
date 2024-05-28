package clients

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"
)

// Helper function to generate a random byte slice of a given length
func randomByteSlice(length int64) []byte {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// TestDefaultBlobEncodingCodec tests the encoding and decoding of random byte streams
func TestDefaultBlobEncodingCodec(t *testing.T) {
	// Create an instance of the DefaultBlobEncodingCodec
	codec := DefaultBlobEncodingCodec{}

	// Number of test iterations
	const iterations = 100

	for i := 0; i < iterations; i++ {
		// Generate a random length for the byte slice
		length, err := rand.Int(rand.Reader, big.NewInt(1024)) // Random length between 0 and 1023
		if err != nil {
			panic(err)
		}
		originalData := randomByteSlice(length.Int64() + 1) // ensure it's not length 0

		// Encode the original data
		encodedData, err := codec.EncodeBlob(originalData)
		if err != nil {
			t.Fatalf("Iteration %d: failed to encode blob: %v", i, err)
		}

		// Decode the encoded data
		decodedData, err := codec.DecodeBlob(encodedData)
		if err != nil {
			t.Fatalf("Iteration %d: failed to decode blob: %v", i, err)
		}

		// Compare the original data with the decoded data
		if !bytes.Equal(originalData, decodedData) {
			t.Fatalf("Iteration %d: original and decoded data do not match\nOriginal: %v\nDecoded: %v", i, originalData, decodedData)
		}
	}
}

func TestIFFTBlobEncodingCodec(t *testing.T) {
	codec := IFFTBlobEncodingCodec{}

	// Test data
	originalData := []byte("test data for IFFT encoding")

	// Encode the data using IFFTBlobEncodingCodec
	encodedData, err := codec.EncodeBlob(originalData)
	if err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}

	// Decode the encoded data to verify it matches the original
	decodedData, err := codec.DecodeBlob(encodedData)
	if err != nil {
		t.Fatalf("Failed to decode data: %v", err)
	}

	// Check if the original and decoded data match
	if !bytes.Equal(originalData, decodedData) {
		t.Errorf("Original data and decoded data do not match. Original: %v, Decoded: %v", originalData, decodedData)
	}
}
