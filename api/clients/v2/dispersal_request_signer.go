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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// DispersalRequestSigner encapsulates the logic for signing GetChunks requests.
type DispersalRequestSigner interface {
	// SignStoreChunksRequest signs a StoreChunksRequest. Does not modify the request
	// (i.e. it does not insert the signature).
	SignStoreChunksRequest(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error)
}

type DispersalRequestSignerConfig struct {
	// KeyID is the AWS KMS key identifier used for signing requests. Optional if PrivateKey is provided.
	KeyID string `docs:"required"`
	// PrivateKey is a hex-encoded private key for local signing (without 0x prefix). Optional if KeyID is provided.
	PrivateKey string `docs:"required"`
	// Region is the default AWS region (e.g., "us-east-1"). Required if using KMS.
	Region string `docs:"required"`
	// KMSRegion is an optional AWS region override for KMS operations. When specified, this region is used
	// for KMS operations instead of Region. If empty, Region is used. This allows KMS keys to be stored in
	// a different region than other AWS resources.
	KMSRegion string
	// Endpoint is an optional custom AWS KMS endpoint URL. If empty, the standard AWS KMS endpoint is used.
	// This is primarily useful for testing with LocalStack or other custom KMS implementations. Default is empty.
	Endpoint string
}

var _ config.VerifiableConfig = &DispersalRequestSignerConfig{}

func DefaultDispersalRequestSignerConfig() DispersalRequestSignerConfig {
	return DispersalRequestSignerConfig{}
}

// Verify checks that the configuration is valid, returning an error if it is not.
func (c *DispersalRequestSignerConfig) Verify() error {
	if c.KeyID == "" && c.PrivateKey == "" {
		return errors.New("either KeyID or PrivateKey is required")
	}
	if c.KeyID != "" && c.PrivateKey != "" {
		return errors.New("KeyID and PrivateKey cannot be specified together")
	}
	if c.KeyID != "" && c.Region == "" {
		return errors.New("Region is required when using KMS")
	}

	return nil
}

// kmsRequestSigner implements DispersalRequestSigner using AWS KMS.
type kmsRequestSigner struct {
	keyID     string
	publicKey *ecdsa.PublicKey
	kmsClient *kms.Client
}

var _ DispersalRequestSigner = &kmsRequestSigner{}

// localRequestSigner implements DispersalRequestSigner using a local private key.
type localRequestSigner struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

var _ DispersalRequestSigner = &localRequestSigner{}

// NewDispersalRequestSigner creates a new DispersalRequestSigner.
func NewDispersalRequestSigner(
	ctx context.Context,
	config DispersalRequestSignerConfig,
) (DispersalRequestSigner, error) {
	if err := config.Verify(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Use KMS if KeyID is provided
	if config.KeyID != "" {
		return NewKMSDispersalRequestSigner(ctx, config)
	}

	// Use local private key
	return NewLocalDispersalRequestSigner(config)
}

// NewKMSDispersalRequestSigner creates a new KMS-based DispersalRequestSigner.
func NewKMSDispersalRequestSigner(
	ctx context.Context,
	config DispersalRequestSignerConfig,
) (DispersalRequestSigner, error) {
	// Determine which region to use for KMS operations.
	// If KMSRegion is specified, use it; otherwise fall back to Region.
	kmsRegion := config.Region
	if config.KMSRegion != "" {
		kmsRegion = config.KMSRegion
	}

	var kmsClient *kms.Client
	if config.Endpoint != "" {
		kmsClient = kms.New(kms.Options{
			Region:       kmsRegion,
			BaseEndpoint: aws.String(config.Endpoint),
		})
	} else {
		// Load the AWS SDK configuration, which will automatically detect credentials
		// from environment variables, IAM roles, or AWS config files
		cfg, err := awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(kmsRegion),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS config for region %s: %w", kmsRegion, err)
		}
		kmsClient = kms.NewFromConfig(cfg)
	}

	key, err := aws2.LoadPublicKeyKMS(ctx, kmsClient, config.KeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ecdsa public key from KMS in region %s: %w", kmsRegion, err)
	}

	return &kmsRequestSigner{
		keyID:     config.KeyID,
		publicKey: key,
		kmsClient: kmsClient,
	}, nil
}

// NewLocalDispersalRequestSigner creates a new local private key-based DispersalRequestSigner.
func NewLocalDispersalRequestSigner(
	config DispersalRequestSignerConfig,
) (DispersalRequestSigner, error) {
	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &localRequestSigner{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

func (s *kmsRequestSigner) SignStoreChunksRequest(
	ctx context.Context,
	request *grpc.StoreChunksRequest,
) ([]byte, error) {
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

func (s *localRequestSigner) SignStoreChunksRequest(
	ctx context.Context,
	request *grpc.StoreChunksRequest,
) ([]byte, error) {
	hash, err := hashing.HashStoreChunksRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to hash request: %w", err)
	}

	signature, err := crypto.Sign(hash, s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}
	return signature, nil
}
