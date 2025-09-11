package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Encapsulates AWS KMS signature verification operations for ECDSA signatures.
type KMSSignatureVerifier struct {
	address gethcommon.Address
}

func NewKMSSignatureVerifier(
	ctx context.Context,
	kmsRegion string,
	keyID string,
	kmsEndpoint string,
) (*KMSSignatureVerifier, error) {
	if kmsRegion == "" {
		return nil, fmt.Errorf("KMS region is required")
	}
	if keyID == "" {
		return nil, fmt.Errorf("KMS key ID is required")
	}

	var kmsClient *kms.Client
	if kmsEndpoint != "" {
		kmsClient = kms.New(kms.Options{
			Region:       kmsRegion,
			BaseEndpoint: &kmsEndpoint,
		})
	} else {
		awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(kmsRegion))
		if err != nil {
			return nil, fmt.Errorf("load AWS config: %w", err)
		}

		kmsClient = kms.NewFromConfig(awsConfig)
	}

	publicKey, err := LoadPublicKeyKMS(ctx, kmsClient, keyID)
	if err != nil {
		return nil, fmt.Errorf("load public key from KMS: %w", err)
	}

	address := crypto.PubkeyToAddress(*publicKey)

	return &KMSSignatureVerifier{address: address}, nil
}

// Verifies that a signature was created by the KMS key.
// Returns nil if the signature is valid, or an error if it's not.
func (v *KMSSignatureVerifier) VerifySignature(hash []byte, signature []byte) error {
	if len(signature) != 65 {
		return fmt.Errorf("invalid signature length %d, expected 65", len(signature))
	}

	signingPubkey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return fmt.Errorf("recover public key from signature: %w", err)
	}

	signingAddress := crypto.PubkeyToAddress(*signingPubkey)
	if signingAddress != v.address {
		return fmt.Errorf(
			"signature doesn't match expected address: got %s, expected %s", signingAddress.Hex(), v.address.Hex())
	}

	return nil
}
