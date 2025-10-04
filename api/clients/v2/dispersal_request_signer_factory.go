package clients

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/oci"
	oraclecommon "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/keymanagement"
)

// NewDispersalRequestSignerFromKMSConfig creates a DispersalRequestSigner based on KMS configuration
func NewDispersalRequestSignerFromKMSConfig(
	ctx context.Context,
	kmsConfig common.KMSKeyConfig,
	awsRegion, awsEndpoint string) (DispersalRequestSigner, error) {

	switch kmsConfig.Provider {
	case "aws":
		// Use existing AWS implementation with backward compatibility
		keyID := kmsConfig.KeyID
		if keyID == "" {
			return nil, fmt.Errorf("AWS KMS key ID is required")
		}
		region := kmsConfig.Region
		if region == "" {
			region = awsRegion // Fall back to AWS client config region for backward compatibility
		}
		return NewDispersalRequestSigner(ctx, region, awsEndpoint, keyID)

	case "oci":
		if kmsConfig.KeyOCID == "" || kmsConfig.KMSEndpoint == "" || kmsConfig.ManagementEndpoint == "" {
			return nil, fmt.Errorf("OCI KMS key OCID, KMS endpoint, and management endpoint are required")
		}
		return NewOCIDispersalRequestSigner(ctx, kmsConfig.KeyOCID, kmsConfig.KMSEndpoint, kmsConfig.ManagementEndpoint)

	default:
		return nil, fmt.Errorf("unsupported KMS provider: %s (supported: aws, oci)", kmsConfig.Provider)
	}
}

// NewOCIDispersalRequestSigner creates an OCI KMS-based DispersalRequestSigner
func NewOCIDispersalRequestSigner(
	ctx context.Context,
	keyOCID, kmsEndpoint, managementEndpoint string) (DispersalRequestSigner, error) {

	// Create OCI configuration provider
	configProvider := oraclecommon.DefaultConfigProvider()

	// Create OCI KMS clients
	cryptoClient, err := keymanagement.NewKmsCryptoClientWithConfigurationProvider(configProvider, kmsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create OCI KMS crypto client: %w", err)
	}

	managementClient, err := keymanagement.NewKmsManagementClientWithConfigurationProvider(
		configProvider, managementEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create OCI KMS management client: %w", err)
	}

	// Load the public key
	publicKey, err := oci.LoadPublicKeyKMS(ctx, managementClient, keyOCID)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key from OCI KMS: %w", err)
	}

	return &ociRequestSigner{
		keyOCID:          keyOCID,
		publicKey:        publicKey,
		cryptoClient:     cryptoClient,
		managementClient: managementClient,
	}, nil
}
