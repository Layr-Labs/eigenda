package aws

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"math/big"
)

// This file contains utility methods for working with AWS KMS using ecdsa on the KeySpecEccSecgP256k1 curve.
// This code was adapted from code in https://github.com/Layr-Labs/eigensdk-go/tree/dev/signerv2

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

// LoadPublicKeyKMS loads the public key from AWS KMS.
func LoadPublicKeyKMS(
	ctx context.Context,
	client *kms.Client,
	keyId string) (*ecdsa.PublicKey, error) {

	getPubKeyOutput, err := client.GetPublicKey(ctx, &kms.GetPublicKeyInput{
		KeyId: aws.String(keyId),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get public key for KeyId=%s: %w", keyId, err)
	}

	key, err := ParsePublicKeyKMS(getPubKeyOutput.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key for KeyId=%s: %w", keyId, err)
	}

	return key, nil
}

// ParsePublicKeyKMS parses the public key from AWS KMS format into an ecdsa.PublicKey.
func ParsePublicKeyKMS(bytes []byte) (*ecdsa.PublicKey, error) {
	var asn1pubk asn1EcPublicKey
	_, err := asn1.Unmarshal(bytes, &asn1pubk)
	if err != nil {
		return nil, fmt.Errorf("asn1.Uunmarshal failed: %w", err)
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

// SignKMS signs a hash using the provided public using AWS KMS.
// The signature is returned in the 65-byte format used by Ethereum.
func SignKMS(
	ctx context.Context,
	client *kms.Client,
	keyId string,
	publicKey *ecdsa.PublicKey,
	hash []byte) ([]byte, error) {

	signOutput, err := client.Sign(ctx, &kms.SignInput{
		KeyId:            aws.String(keyId),
		SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
		MessageType:      types.MessageTypeDigest,
		Message:          hash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sign hash: %w", err)
	}

	signature, err := ParseSignatureKMS(publicKey, hash, signOutput.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signature: %w", err)
	}

	return signature, nil
}

// ParseSignatureKMS parses a signature (KeySpecEccSecgP256k1) in the format returned by amazon KMS into the
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
		return nil, err
	}

	if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(publicKeyBytes) {
		signature = append(rsSignature, []byte{1}...)
		recoveredPublicKeyBytes, err = crypto.Ecrecover(hash, signature)
		if err != nil {
			return nil, err
		}

		if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(publicKeyBytes) {
			return nil, errors.New("can not reconstruct public key from sig")
		}
	}

	return signature, nil
}
