package aws

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// Encapsulates AWS KMS signing operations for ECDSA signatures.
type KMSSigner struct {
	keyID     string
	publicKey *ecdsa.PublicKey
	kmsClient *kms.Client
}

func NewKMSSigner(ctx context.Context, kmsRegion string, keyID string, kmsEndpoint string) (*KMSSigner, error) {
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

	return &KMSSigner{
		keyID:     keyID,
		publicKey: publicKey,
		kmsClient: kmsClient,
	}, nil
}

// Signs a hash using the configured KMS key.
// The signature is returned in the 65-byte format used by Ethereum.
func (s *KMSSigner) Sign(ctx context.Context, hash []byte) ([]byte, error) {
	return SignKMS(ctx, s.kmsClient, s.keyID, s.publicKey, hash)
}
