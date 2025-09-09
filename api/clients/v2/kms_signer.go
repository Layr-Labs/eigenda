package clients

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// KMSSigner is a base struct that provides KMS signing functionality.
// It can be embedded by other signers to reuse the common KMS setup and signing logic.
type KMSSigner struct {
	keyID      string
	publicKey  *ecdsa.PublicKey
	keyManager *kms.Client
}

// NewKMSSigner creates a new KMSSigner with the specified KMS configuration.
func NewKMSSigner(
	ctx context.Context,
	region string,
	endpoint string,
	keyID string) (*KMSSigner, error) {

	// Load the AWS SDK configuration, which will automatically detect credentials
	// from environment variables, IAM roles, or AWS config files
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	var keyManager *kms.Client
	if endpoint != "" {
		keyManager = kms.New(kms.Options{
			Region:       region,
			BaseEndpoint: aws.String(endpoint),
		})
	} else {
		keyManager = kms.NewFromConfig(cfg)
	}

	key, err := aws2.LoadPublicKeyKMS(ctx, keyManager, keyID)
	if err != nil {
		return nil, fmt.Errorf("get ecdsa public key: %w", err)
	}

	return &KMSSigner{
		keyID:      keyID,
		publicKey:  key,
		keyManager: keyManager,
	}, nil
}

// SignHash signs a hash using the configured KMS key.
func (s *KMSSigner) SignHash(ctx context.Context, hash []byte) ([]byte, error) {
	signature, err := aws2.SignKMS(ctx, s.keyManager, s.keyID, s.publicKey, hash)
	if err != nil {
		return nil, fmt.Errorf("sign hash: %w", err)
	}

	return signature, nil
}
