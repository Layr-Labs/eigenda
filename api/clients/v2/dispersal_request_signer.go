package clients

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing"
	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/pkg/errors"
)

// DispersalRequestSigner encapsulates the logic for signing GetChunks requests.
type DispersalRequestSigner interface {
	// SignStoreChunksRequest signs a StoreChunksRequest. Does not modify the request
	// (i.e. it does not insert the signature).
	SignStoreChunksRequest(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error)
}

type DispersalRequestSignerConfig struct {
	// KeyID is the AWS KMS key identifier used for signing requests. Required.
	KeyID string
	// Region is the AWS region where the KMS key is located (e.g., "us-east-1"). Required.
	Region string
	// Endpoint is an optional custom AWS KMS endpoint URL. If empty, the standard AWS KMS endpoint is used.
	// This is primarily useful for testing with LocalStack or other custom KMS implementations. Default is empty.
	Endpoint string
}

var _ config.VerifiableConfig = &DispersalRequestSignerConfig{}

func DefaultDispersalRequestSignerConfig() DispersalRequestSignerConfig {
	return DispersalRequestSignerConfig{
		Region:   "us-east-1",
		Endpoint: "",
	}
}

// Verify checks that the configuration is valid, returning an error if it is not.
func (c *DispersalRequestSignerConfig) Verify() error {
	if c.KeyID == "" {
		return errors.New("KeyID is required")
	}
	if c.Region == "" {
		return errors.New("Region is required")
	}

	return nil
}

var _ DispersalRequestSigner = &requestSigner{}

type requestSigner struct {
	keyID     string
	publicKey *ecdsa.PublicKey
	kmsClient *kms.Client
}

// NewDispersalRequestSigner creates a new DispersalRequestSigner.
func NewDispersalRequestSigner(
	ctx context.Context,
	config DispersalRequestSignerConfig,
) (DispersalRequestSigner, error) {
	if err := config.Verify(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	var kmsClient *kms.Client
	if config.Endpoint != "" {
		kmsClient = kms.New(kms.Options{
			Region:       config.Region,
			BaseEndpoint: aws.String(config.Endpoint),
		})
	} else {
		// Load the AWS SDK configuration, which will automatically detect credentials
		// from environment variables, IAM roles, or AWS config files
		cfg, err := awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(config.Region),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS config: %w", err)
		}
		kmsClient = kms.NewFromConfig(cfg)
	}

	key, err := aws2.LoadPublicKeyKMS(ctx, kmsClient, config.KeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ecdsa public key: %w", err)
	}

	return &requestSigner{
		keyID:     config.KeyID,
		publicKey: key,
		kmsClient: kmsClient,
	}, nil
}

func (s *requestSigner) SignStoreChunksRequest(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {
	hash, err := hashing.HashStoreChunksRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to hash request: %w", err)
	}

	signature, err := aws2.SignKMS(ctx, s.kmsClient, s.keyID, s.publicKey, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	return signature, nil
}
