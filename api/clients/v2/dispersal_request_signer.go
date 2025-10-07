package clients

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing"
	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// DispersalRequestSigner encapsulates the logic for signing GetChunks requests.
type DispersalRequestSigner interface {
	// SignStoreChunksRequest signs a StoreChunksRequest. Does not modify the request
	// (i.e. it does not insert the signature).
	SignStoreChunksRequest(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error)
}

type DispersalRequestSignerConfig struct {
	Region   string
	Endpoint string
	KeyID    string
}

func DefaultDispersalRequestSignerConfig() DispersalRequestSignerConfig {
	return DispersalRequestSignerConfig{
		Region: "us-east-1",
	}
}

var _ DispersalRequestSigner = &requestSigner{}

type requestSigner struct {
	keyID      string
	publicKey  *ecdsa.PublicKey
	keyManager *kms.Client
}

// NewDispersalRequestSigner creates a new DispersalRequestSigner.
func NewDispersalRequestSigner(
	ctx context.Context,
	config DispersalRequestSignerConfig,
) (DispersalRequestSigner, error) {
	// Load the AWS SDK configuration, which will automatically detect credentials
	// from environment variables, IAM roles, or AWS config files
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(config.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var keyManager *kms.Client
	if config.Endpoint != "" {
		keyManager = kms.New(kms.Options{
			Region:       config.Region,
			BaseEndpoint: aws.String(config.Endpoint),
		})
	} else {
		keyManager = kms.NewFromConfig(cfg)
	}

	key, err := aws2.LoadPublicKeyKMS(ctx, keyManager, config.KeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ecdsa public key: %w", err)
	}

	return &requestSigner{
		keyID:      config.KeyID,
		publicKey:  key,
		keyManager: keyManager,
	}, nil
}

func (s *requestSigner) SignStoreChunksRequest(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {
	hash, err := hashing.HashStoreChunksRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to hash request: %w", err)
	}

	signature, err := aws2.SignKMS(ctx, s.keyManager, s.keyID, s.publicKey, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	return signature, nil
}
