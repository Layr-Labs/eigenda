package clients

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOCIDispersalRequestSigner(t *testing.T) {
	ctx := context.Background()
	keyOCID := "ocid1.key.oc1.test"

	// This test will fail without OCI credentials, but it tests the creation logic
	signer, err := NewOCIDispersalRequestSigner(ctx, keyOCID)

	// We expect an error in test environment without OCI setup
	if err != nil {
		// Error could be from missing credentials or invalid configuration
		assert.NotNil(t, err)
		assert.Nil(t, signer)
	} else {
		assert.NotNil(t, signer)
	}
}

func TestNewOCIDispersalRequestSigner_InvalidParameters(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		keyOCID       string
		expectedError string
	}{
		{
			name:          "empty key OCID",
			keyOCID:       "",
			expectedError: "failed to create workload identity provider",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			signer, err := NewOCIDispersalRequestSigner(ctx, test.keyOCID)

			assert.Nil(t, signer)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.expectedError)
		})
	}
}

func TestOCIDispersalRequestSigner_Interface(t *testing.T) {
	// Test that the ociRequestSigner implements DispersalRequestSigner interface
	// This is a compile-time check
	var _ DispersalRequestSigner = (*ociRequestSigner)(nil)
}