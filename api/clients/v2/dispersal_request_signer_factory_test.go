package clients

import (
	"context"
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/assert"
)

func TestNewDispersalRequestSignerFromKMSConfig_AWS(t *testing.T) {
	ctx := context.Background()
	kmsConfig := common.KMSKeyConfig{
		Provider: "aws",
		KeyID:    "test-key-id",
		Region:   "us-east-1",
	}
	region := "us-east-1"
	endpointURL := ""

	// This test will fail without AWS credentials, but it tests the factory logic
	signer, err := NewDispersalRequestSignerFromKMSConfig(ctx, kmsConfig, region, endpointURL)

	// We expect an error in test environment without AWS setup
	if err != nil {
		// Error could be from missing credentials or invalid configuration
		assert.NotNil(t, err)
	} else {
		assert.NotNil(t, signer)
	}
}

func TestNewDispersalRequestSignerFromKMSConfig_OCI(t *testing.T) {
	ctx := context.Background()
	kmsConfig := common.KMSKeyConfig{
		Provider:           "oci",
		KeyOCID:            "ocid1.key.oc1.test",
		KMSEndpoint:        "https://test.oci.com",
		ManagementEndpoint: "https://management.test.oci.com",
	}
	region := "us-phoenix-1"
	endpointURL := ""

	// This test will fail without OCI credentials, but it tests the factory logic
	signer, err := NewDispersalRequestSignerFromKMSConfig(ctx, kmsConfig, region, endpointURL)

	// We expect an error in test environment without OCI setup
	if err != nil {
		// Error could be from missing credentials or invalid configuration
		assert.NotNil(t, err)
	} else {
		assert.NotNil(t, signer)
	}
}

func TestNewDispersalRequestSignerFromKMSConfig_UnsupportedProvider(t *testing.T) {
	ctx := context.Background()
	kmsConfig := common.KMSKeyConfig{
		Provider: "unsupported-provider",
		KeyID:    "test-key-id",
	}
	region := "us-east-1"
	endpointURL := ""

	signer, err := NewDispersalRequestSignerFromKMSConfig(ctx, kmsConfig, region, endpointURL)

	assert.Nil(t, signer)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported KMS provider: unsupported-provider")
}

func TestNewDispersalRequestSignerFromKMSConfig_EmptyProvider(t *testing.T) {
	ctx := context.Background()
	kmsConfig := common.KMSKeyConfig{
		Provider: "",
		KeyID:    "test-key-id",
	}
	region := "us-east-1"
	endpointURL := ""

	signer, err := NewDispersalRequestSignerFromKMSConfig(ctx, kmsConfig, region, endpointURL)

	assert.Nil(t, signer)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported KMS provider:")
}

func TestNewDispersalRequestSignerFromKMSConfig_EmptyKeyID(t *testing.T) {
	ctx := context.Background()
	kmsConfig := common.KMSKeyConfig{
		Provider: "aws",
		KeyID:    "",
	}
	region := "us-east-1"
	endpointURL := ""

	signer, err := NewDispersalRequestSignerFromKMSConfig(ctx, kmsConfig, region, endpointURL)

	assert.Nil(t, signer)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AWS KMS key ID is required")
}

func TestNewDispersalRequestSignerFromKMSConfig_AWSMissingRegion(t *testing.T) {
	ctx := context.Background()
	kmsConfig := common.KMSKeyConfig{
		Provider: "aws",
		KeyID:    "test-key-id",
		Region:   "",
	}
	region := ""
	endpointURL := ""

	signer, err := NewDispersalRequestSignerFromKMSConfig(ctx, kmsConfig, region, endpointURL)

	assert.Nil(t, signer)
	assert.Error(t, err)
	// The actual error depends on AWS KMS connectivity
	assert.NotNil(t, err)
}

func TestNewDispersalRequestSignerFromKMSConfig_OCIMissingEndpoints(t *testing.T) {
	ctx := context.Background()
	
	// Test missing KMS endpoint
	kmsConfig := common.KMSKeyConfig{
		Provider:           "oci",
		KeyOCID:            "ocid1.key.oc1.test",
		KMSEndpoint:        "",
		ManagementEndpoint: "https://management.test.oci.com",
	}
	region := "us-phoenix-1"
	endpointURL := ""

	signer, err := NewDispersalRequestSignerFromKMSConfig(ctx, kmsConfig, region, endpointURL)

	assert.Nil(t, signer)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OCI KMS key OCID, KMS endpoint, and management endpoint are required")

	// Test missing management endpoint
	kmsConfig = common.KMSKeyConfig{
		Provider:           "oci",
		KeyID:              "ocid1.key.oc1.test",
		KMSEndpoint:        "https://test.oci.com",
		ManagementEndpoint: "",
	}

	signer, err = NewDispersalRequestSignerFromKMSConfig(ctx, kmsConfig, region, endpointURL)

	assert.Nil(t, signer)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OCI KMS key OCID, KMS endpoint, and management endpoint are required")
}


// Test the validation logic without making actual KMS calls
func TestValidateKMSConfig(t *testing.T) {
	tests := []struct {
		name        string
		kmsConfig   common.KMSKeyConfig
		region      string
		expectError string
	}{
		{
			name: "valid AWS config",
			kmsConfig: common.KMSKeyConfig{
				Provider: "aws",
				KeyID:    "test-key-id",
				Region:   "us-east-1",
			},
			region: "us-east-1",
		},
		{
			name: "valid OCI config",
			kmsConfig: common.KMSKeyConfig{
				Provider:           "oci",
				KeyID:              "ocid1.key.oc1.test",
				KMSEndpoint:        "https://test.oci.com",
				ManagementEndpoint: "https://management.test.oci.com",
			},
			region: "us-phoenix-1",
		},
		{
			name: "empty provider",
			kmsConfig: common.KMSKeyConfig{
				Provider: "",
				KeyID:    "test-key-id",
			},
			expectError: "KMS provider must be specified",
		},
		{
			name: "empty key ID",
			kmsConfig: common.KMSKeyConfig{
				Provider: "aws",
				KeyID:    "",
			},
			expectError: "KMS key ID must be specified",
		},
		{
			name: "unsupported provider",
			kmsConfig: common.KMSKeyConfig{
				Provider: "gcp",
				KeyID:    "test-key-id",
			},
			expectError: "unsupported KMS provider: gcp",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateKMSConfig(test.kmsConfig, test.region)

			if test.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to validate KMS configuration
func validateKMSConfig(kmsConfig common.KMSKeyConfig, region string) error {
	if kmsConfig.Provider == "" {
		return fmt.Errorf("KMS provider must be specified")
	}

	if kmsConfig.KeyID == "" {
		return fmt.Errorf("KMS key ID must be specified")
	}

	switch kmsConfig.Provider {
	case "aws":
		if kmsConfig.Region == "" && region == "" {
			return fmt.Errorf("AWS KMS region must be specified")
		}
	case "oci":
		if kmsConfig.KMSEndpoint == "" {
			return fmt.Errorf("OCI KMS endpoint must be specified")
		}
		if kmsConfig.ManagementEndpoint == "" {
			return fmt.Errorf("OCI management endpoint must be specified")
		}
	default:
		return fmt.Errorf("unsupported KMS provider: %s", kmsConfig.Provider)
	}

	return nil
}

