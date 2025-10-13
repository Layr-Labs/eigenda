package oci

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/keymanagement"
)

// This file contains utility methods for working with OCI KMS using ecdsa on the secp256k1 curve.
// This code was adapted from common/aws/kms.go

var secp256k1N = crypto.S256().Params().N
var secp256k1HalfN = new(big.Int).Div(secp256k1N, big.NewInt(2))

type asn1EcPublicKey struct {
	EcPublicKeyInfo asn1EcPublicKeyInfo
	PublicKey       asn1.BitString
}

type asn1EcPublicKeyInfo struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.ObjectIdentifier
}

type asn1EcSig struct {
	R asn1.RawValue
	S asn1.RawValue
}

// LoadPublicKeyKMS loads the public key from OCI KMS.
func LoadPublicKeyKMS(
	ctx context.Context,
	managementClient keymanagement.KmsManagementClient,
	keyId string) (*ecdsa.PublicKey, error) {

	// Get the key details first
	getKeyRequest := keymanagement.GetKeyRequest{
		KeyId: common.String(keyId),
	}

	keyResponse, err := managementClient.GetKey(ctx, getKeyRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get key for KeyId=%s: %w", keyId, err)
	}

	// Get the key version to retrieve the public key
	getKeyVersionRequest := keymanagement.GetKeyVersionRequest{
		KeyId:        common.String(keyId),
		KeyVersionId: keyResponse.CurrentKeyVersion,
	}

	keyVersionResponse, err := managementClient.GetKeyVersion(ctx, getKeyVersionRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get key version for KeyId=%s: %w", keyId, err)
	}

	if keyVersionResponse.PublicKey == nil {
		return nil, fmt.Errorf("public key not available for KeyId=%s", keyId)
	}

	key, err := ParsePublicKeyKMS([]byte(*keyVersionResponse.PublicKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key for KeyId=%s: %w", keyId, err)
	}

	return key, nil
}

// ParsePublicKeyKMS parses the public key from OCI KMS format (PEM) into an ecdsa.PublicKey.
func ParsePublicKeyKMS(keyBytes []byte) (*ecdsa.PublicKey, error) {
	// First, try to decode as PEM (which is what OCI KMS typically returns)
	block, _ := pem.Decode(keyBytes)
	var derBytes []byte

	if block != nil {
		// Successfully decoded PEM, use the DER bytes
		derBytes = block.Bytes
	} else {
		// Not PEM format, assume raw DER bytes
		derBytes = keyBytes
	}

	// Parse the DER-encoded public key using ASN.1
	var asn1pubk asn1EcPublicKey
	_, err := asn1.Unmarshal(derBytes, &asn1pubk)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key from OCI format: %w", err)
	}

	key, err := crypto.UnmarshalPubkey(asn1pubk.PublicKey.Bytes)
	if err != nil {
		return nil, fmt.Errorf("crypto.UnmarshalPubkey failed: %w", err)
	}

	return key, nil
}

func adjustSignatureLength(buffer []byte) []byte {
	if len(buffer) > 32 {
		buffer = buffer[len(buffer)-32:] // Take last 32 bytes
	}

	buffer = bytes.TrimLeft(buffer, "\x00")
	for len(buffer) < 32 {
		zeroBuf := []byte{0}
		buffer = append(zeroBuf, buffer...)
	}
	return buffer
}

// SignKMS signs a hash using the provided public key using OCI KMS.
// The signature is returned in the 65-byte format used by Ethereum.
func SignKMS(
	ctx context.Context,
	cryptoClient keymanagement.KmsCryptoClient,
	keyId string,
	publicKey *ecdsa.PublicKey,
	hash []byte) ([]byte, error) {

	// OCI KMS expects the message to be base64 encoded
	messageBase64 := base64.StdEncoding.EncodeToString(hash)

	signRequest := keymanagement.SignRequest{
		SignDataDetails: keymanagement.SignDataDetails{
			KeyId:            common.String(keyId),
			Message:          common.String(messageBase64),
			SigningAlgorithm: keymanagement.SignDataDetailsSigningAlgorithmEcdsaSha256,
			MessageType:      keymanagement.SignDataDetailsMessageTypeDigest,
		},
	}

	signResponse, err := cryptoClient.Sign(ctx, signRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to sign hash: %w", err)
	}

	// Decode the base64 signature
	signatureBytes, err := base64.StdEncoding.DecodeString(*signResponse.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}

	signature, err := ParseSignatureKMS(publicKey, hash, signatureBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signature: %w", err)
	}

	return signature, nil
}

// ParseSignatureKMS parses a signature (secp256k1) in the format returned by OCI KMS into the
// 65-byte format used by Ethereum.
func ParseSignatureKMS(
	publicKey *ecdsa.PublicKey,
	hash []byte,
	bytes []byte) ([]byte, error) {

	if !secp256k1.S256().IsOnCurve(publicKey.X, publicKey.Y) {
		return nil, errors.New("public key is not on curve")
	}

	publicKeyBytes := secp256k1.S256().Marshal(publicKey.X, publicKey.Y)

	var sigAsn1 asn1EcSig
	_, err := asn1.Unmarshal(bytes, &sigAsn1)
	if err != nil {
		return nil, fmt.Errorf("asn1.Unmarshal failed: %w", err)
	}

	r := sigAsn1.R.Bytes
	s := sigAsn1.S.Bytes

	// Adjust S value from signature according to Ethereum standard
	sBigInt := new(big.Int).SetBytes(s)
	if sBigInt.Cmp(secp256k1HalfN) > 0 {
		s = new(big.Int).Sub(secp256k1N, sBigInt).Bytes()
	}

	rsSignature := append(adjustSignatureLength(r), adjustSignatureLength(s)...)
	signature := append(rsSignature, []byte{0}...)

	recoveredPublicKeyBytes, err := crypto.Ecrecover(hash, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to recover public key: %w", err)
	}

	if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(publicKeyBytes) {
		signature = append(rsSignature, []byte{1}...)
		recoveredPublicKeyBytes, err = crypto.Ecrecover(hash, signature)
		if err != nil {
			return nil, fmt.Errorf("failed to recover public key with recovery ID 1: %w", err)
		}

		if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(publicKeyBytes) {
			return nil, errors.New("can not reconstruct public key from sig")
		}
	}

	return signature, nil
}
