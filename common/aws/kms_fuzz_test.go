package aws

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

// ecdsaSignature defines the ASN.1 structure for ECDSA signatures.
type ecdsaSignature struct {
	R, S *big.Int
}

// generateValidSignature generates a valid ECDSA signature and returns the public key, hash, and DER signature.
func generateValidSignature() (*ecdsa.PublicKey, []byte, []byte, error) {
	// Generate a secp256k1 ECDSA key pair.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, nil, nil, err
	}
	publicKey := &privateKey.PublicKey

	// Define a message and compute its SHA-256 hash.
	message := "Test message for ECDSA signature"
	hash := sha256.Sum256([]byte(message))

	// Sign the hash using the private key.
	signatureBytes, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// Convert the signature to DER format.
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:64])

	// Marshal R and S into ASN.1 DER format.
	derSignature, err := asn1.Marshal(ecdsaSignature{R: r, S: s})
	if err != nil {
		return nil, nil, nil, err
	}

	return publicKey, hash[:], derSignature, nil
}

// defineEdgeCases returns a slice of tuples containing publicKeyBytes, hashBytes, derSignatureBytes
func defineEdgeCases() [][3][]byte {
	var edgeCases [][3][]byte

	// Helper: Generate a valid signature to obtain a public key.
	pubKeyValidBytes, hashValid, derSigValid, err := generateValidSignature()
	if err != nil {
		panic("Failed to generate valid signature for edge cases")
	}
	publicKeyValid := crypto.FromECDSAPub(pubKeyValidBytes)

	// 1. Malformed Public Keys

	// a. Incorrect length (too short)
	publicKeyShort := []byte{0x04, 0x01, 0x02}
	derSignatureValid := derSigValid
	edgeCases = append(edgeCases, [3][]byte{publicKeyShort, hashValid, derSignatureValid})

	// b. Incorrect prefix
	publicKeyBadPrefix := make([]byte, 65)
	publicKeyBadPrefix[0] = 0x05 // Invalid prefix
	copy(publicKeyBadPrefix[1:], bytes.Repeat([]byte{0x01}, 64))
	edgeCases = append(edgeCases, [3][]byte{publicKeyBadPrefix, hashValid, derSignatureValid})

	// c. Coordinates not on curve (invalid X, Y)
	publicKeyInvalidXY := make([]byte, 65)
	publicKeyInvalidXY[0] = 0x04
	// Set X and Y to values that are not on the curve
	copy(publicKeyInvalidXY[1:], bytes.Repeat([]byte{0xFF}, 64))
	edgeCases = append(edgeCases, [3][]byte{publicKeyInvalidXY, hashValid, derSignatureValid})

	// 2. Malformed Signatures

	// a. Invalid DER encoding (truncated)
	derSignatureInvalidDER := []byte{0x30, 0x00} // Incomplete DER
	edgeCases = append(edgeCases, [3][]byte{publicKeyValid, hashValid, derSignatureInvalidDER})

	// b. R too long (33 bytes with leading zero)
	derSignatureRTooLong := []byte{
		0x30, 0x46, // SEQUENCE, length 70
		0x02, 0x21, // INTEGER, length 33
		0x00, // Leading zero
	}
	derSignatureRTooLong = append(derSignatureRTooLong, bytes.Repeat([]byte{0x01}, 32)...) // R
	derSignatureRTooLong = append(derSignatureRTooLong, 0x02, 0x20)                        // S INTEGER, length 32
	derSignatureRTooLong = append(derSignatureRTooLong, bytes.Repeat([]byte{0x02}, 32)...) // S
	edgeCases = append(edgeCases, [3][]byte{publicKeyValid, hashValid, derSignatureRTooLong})

	// c. S too short (31 bytes)
	derSignatureSTooShort := []byte{
		0x30, 0x44, // SEQUENCE, length 68
		0x02, 0x20, // INTEGER, length 32
	}
	derSignatureSTooShort = append(derSignatureSTooShort, bytes.Repeat([]byte{0x03}, 32)...) // R
	derSignatureSTooShort = append(derSignatureSTooShort, 0x02, 0x1F)                        // S INTEGER, length 31
	derSignatureSTooShort = append(derSignatureSTooShort, bytes.Repeat([]byte{0x04}, 31)...) // S
	edgeCases = append(edgeCases, [3][]byte{publicKeyValid, hashValid, derSignatureSTooShort})

	// 3. Invalid Hashes

	// a. Incorrect hash length (too short)
	hashTooShort := make([]byte, 16)
	edgeCases = append(edgeCases, [3][]byte{publicKeyValid, hashTooShort, derSignatureValid})

	// b. Empty hash
	hashEmpty := []byte{}
	edgeCases = append(edgeCases, [3][]byte{publicKeyValid, hashEmpty, derSignatureValid})

	// 4. Random Data

	// a. Completely random bytes
	randomPublicKey := bytes.Repeat([]byte{0xAB}, 65)
	randomHash := bytes.Repeat([]byte{0xCD}, 32)
	randomSignature := bytes.Repeat([]byte{0xEF}, 70)
	edgeCases = append(edgeCases, [3][]byte{randomPublicKey, randomHash, randomSignature})

	// 5. Boundary Conditions

	// a. R equals zero
	derSignatureRZero, _ := asn1.Marshal(ecdsaSignature{R: big.NewInt(0), S: big.NewInt(1)})
	edgeCases = append(edgeCases, [3][]byte{publicKeyValid, hashValid, derSignatureRZero})

	// b. S equals N (curve order)
	secp256k1N := crypto.S256().Params().N
	derSignatureSEqualsN, _ := asn1.Marshal(ecdsaSignature{R: big.NewInt(1), S: new(big.Int).Set(secp256k1N)})
	edgeCases = append(edgeCases, [3][]byte{publicKeyValid, hashValid, derSignatureSEqualsN})

	// c. S just above N/2
	secp256k1HalfN := new(big.Int).Div(crypto.S256().Params().N, big.NewInt(2))
	sAboveHalfN := new(big.Int).Add(secp256k1HalfN, big.NewInt(1))
	derSignatureSAboveHalfN, _ := asn1.Marshal(ecdsaSignature{R: big.NewInt(1), S: sAboveHalfN})
	edgeCases = append(edgeCases, [3][]byte{publicKeyValid, hashValid, derSignatureSAboveHalfN})

	// d. S just below N/2
	sBelowHalfN := new(big.Int).Sub(secp256k1HalfN, big.NewInt(1))
	derSignatureSBelowHalfN, _ := asn1.Marshal(ecdsaSignature{R: big.NewInt(1), S: sBelowHalfN})
	edgeCases = append(edgeCases, [3][]byte{publicKeyValid, hashValid, derSignatureSBelowHalfN})

	// 6. Extra Data

	// a. Extra bytes appended to the signature
	derSignatureExtra := append(derSignatureValid, 0x00, 0x01, 0x02)
	edgeCases = append(edgeCases, [3][]byte{publicKeyValid, hashValid, derSignatureExtra})

	// b. Missing bytes in the signature
	if len(derSignatureValid) > 2 {
		derSignatureMissing := derSignatureValid[:len(derSignatureValid)-2]
		edgeCases = append(edgeCases, [3][]byte{publicKeyValid, hashValid, derSignatureMissing})
	}

	return edgeCases
}

// FuzzParseSignatureKMS tests the ParseSignatureKMS function with various inputs, including edge cases.
func FuzzParseSignatureKMS(f *testing.F) {
	// Generate multiple valid seed inputs
	for i := 0; i < 5; i++ {
		publicKey, hash, derSignature, err := generateValidSignature()
		if err != nil {
			f.Fatalf("Failed to generate valid signature: %v", err)
		}
		publicKeyBytes := crypto.FromECDSAPub(publicKey)
		f.Add(publicKeyBytes, hash, derSignature)
	}

	// Incorporate edge cases into the fuzz corpus
	edgeCases := defineEdgeCases()
	for _, ec := range edgeCases {
		f.Add(ec[0], ec[1], ec[2])
	}

	// Define the fuzzing function
	f.Fuzz(func(t *testing.T, publicKeyBytes []byte, hashBytes []byte, derSignatureBytes []byte) {
		// Skip iteration if publicKeyBytes is not the correct length
		if len(publicKeyBytes) != 65 {
			return
		}

		// Attempt to parse the public key
		pubKey, err := ParsePublicKeyKMS(publicKeyBytes)
		if err != nil {
			// Invalid public key; acceptable for fuzzing
			return
		}

		// Attempt to parse the signature
		signature, err := ParseSignatureKMS(pubKey, hashBytes, derSignatureBytes)
		if err != nil {
			// Parsing failed; acceptable for fuzzing
			return
		}

		// Validate that the signature is exactly 65 bytes
		if len(signature) != 65 {
			t.Errorf("Expected signature length 65 bytes, got %d bytes", len(signature))
		}

		// if the code made it this far, then the pubkey and signature are valid so recovery must work.
		recoveredPubBytes, err := crypto.Ecrecover(hashBytes, signature)
		if err != nil {
			t.Errorf("Ecrecover failed: %v", err)
			return
		}

		// Compare the recovered public key with the original
		if !bytes.Equal(recoveredPubBytes, publicKeyBytes) {
			// Attempt with the possible V values
			signatureCheck := false
			if signature[64] == 27 {
				recoveredPubBytes, err = crypto.Ecrecover(hashBytes, signature)
				if err != nil {
					t.Errorf("Ecrecover failed with V=27: %v", err)
				} else if !bytes.Equal(recoveredPubBytes, publicKeyBytes) {
					t.Errorf("Recovered public key does not match original")
				} else {
					signatureCheck = true
				}
			}

			if !signatureCheck {
				signature[64] = 28
				recoveredPubBytes, err = crypto.Ecrecover(hashBytes, signature)
				if err != nil {
					t.Errorf("Ecrecover failed with V=28: %v", err)
					return
				}

				if !bytes.Equal(recoveredPubBytes, publicKeyBytes) {
					t.Errorf("Recovered public key does not match original")
					return
				}
			}
		}
	})
}
