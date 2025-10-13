package clients

import (
	"context"
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/oci"
	oraclecommon "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
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
		if kmsConfig.KeyID == "" {
			return nil, fmt.Errorf("OCI KMS key OCID is required")
		}
		return NewOCIDispersalRequestSigner(ctx, kmsConfig.KeyID)

	default:
		return nil, fmt.Errorf("unsupported KMS provider: %s (supported: aws, oci)", kmsConfig.Provider)
	}
}

// NewOCIDispersalRequestSigner creates an OCI KMS-based DispersalRequestSigner
func NewOCIDispersalRequestSigner(
	ctx context.Context,
	keyOCID string) (DispersalRequestSigner, error) {

	// Create OCI configuration provider using workload identity
	configProvider, err := auth.OkeWorkloadIdentityConfigurationProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to create workload identity provider: %w", err)
	}

	// Get region and compartment ID from environment
	region := os.Getenv("OCI_REGION")
	if region == "" {
		region = os.Getenv("OCI_RESOURCE_PRINCIPAL_REGION")
	}
	if region == "" {
		return nil, fmt.Errorf("OCI_REGION or OCI_RESOURCE_PRINCIPAL_REGION environment variable is required")
	}

	compartmentID := os.Getenv("OCI_COMPARTMENT_ID")
	if compartmentID == "" {
		return nil, fmt.Errorf("OCI_COMPARTMENT_ID environment variable is required")
	}

	// Discover vault endpoints by finding the vault that contains the key
	cryptoEndpoint, managementEndpoint, err := findVaultEndpointsForKey(ctx, configProvider, keyOCID, compartmentID, region)
	if err != nil {
		return nil, fmt.Errorf("failed to discover vault endpoints for key %s: %w", keyOCID, err)
	}

	// Create OCI KMS clients with discovered endpoints
	cryptoClient, err := keymanagement.NewKmsCryptoClientWithConfigurationProvider(configProvider, cryptoEndpoint)
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

// findVaultEndpointsForKey discovers the vault endpoints for a given key OCID
// by listing all vaults and searching for the key across them
func findVaultEndpointsForKey(ctx context.Context, provider oraclecommon.ConfigurationProvider, keyOCID, compartmentID, region string) (cryptoEndpoint, managementEndpoint string, err error) {
	// Create vault client to list vaults
	vaultClient, err := keymanagement.NewKmsVaultClientWithConfigurationProvider(provider)
	if err != nil {
		return "", "", fmt.Errorf("failed to create KMS vault client: %w", err)
	}

	vaultClient.SetRegion(region)

	// List all vaults
	vaultReq := keymanagement.ListVaultsRequest{
		CompartmentId: oraclecommon.String(compartmentID),
	}

	vaultResp, err := vaultClient.ListVaults(ctx, vaultReq)
	if err != nil {
		return "", "", fmt.Errorf("failed to list vaults: %w", err)
	}

	// Search for the key in each vault
	for _, vault := range vaultResp.Items {
		if vault.LifecycleState != keymanagement.VaultSummaryLifecycleStateActive {
			continue
		}

		// Create management client for this vault
		kmsClient, err := keymanagement.NewKmsManagementClientWithConfigurationProvider(provider, *vault.ManagementEndpoint)
		if err != nil {
			continue // Skip this vault if we can't create the client
		}

		// Check if the key exists in this vault
		// We can try to get the key directly by OCID
		getKeyReq := keymanagement.GetKeyRequest{
			KeyId: oraclecommon.String(keyOCID),
		}

		_, err = kmsClient.GetKey(ctx, getKeyReq)
		if err == nil {
			// Key found in this vault
			return *vault.CryptoEndpoint, *vault.ManagementEndpoint, nil
		}

		// If direct lookup fails, the key might not be in this vault
		// Continue to next vault
	}

	return "", "", fmt.Errorf("key %s not found in any vault in compartment %s", keyOCID, compartmentID)
}
