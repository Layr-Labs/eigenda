package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountID(t *testing.T) {
	// Test case with known private key and expected account ID
	privateKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
	expectedAccountID := "0x1aa8226f6d354380dDE75eE6B634875c4203e522"

	// Create signer instance
	signer := NewLocalBlobRequestSigner(privateKey)

	// Get account ID
	accountID, err := signer.GetAccountID()
	assert.NoError(t, err)
	assert.Equal(t, expectedAccountID, accountID)
}
