package auth_test

import (
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {

	// Make the authenticator
	authenticator := auth.NewAuthenticator(auth.AuthConfig{})

	// Make the signer
	privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	signer := auth.NewSigner(privateKeyHex)

	testHeader := core.BlobAuthHeader{
		BlobCommitments:    core.BlobCommitments{},
		AccountID:          signer.GetAccountID(),
		Nonce:              rand.Uint32(),
		AuthenticationData: []byte{},
	}

	// Sign the header
	signature, err := signer.SignBlobRequest(testHeader)
	assert.NoError(t, err)

	testHeader.AuthenticationData = signature

	err = authenticator.AuthenticateBlobRequest(testHeader)
	assert.NoError(t, err)

}

func TestAuthenticationFail(t *testing.T) {

	// Make the authenticator
	authenticator := auth.NewAuthenticator(auth.AuthConfig{})

	// Make the signer
	privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	signer := auth.NewSigner(privateKeyHex)

	testHeader := core.BlobAuthHeader{
		BlobCommitments:    core.BlobCommitments{},
		AccountID:          signer.GetAccountID(),
		Nonce:              rand.Uint32(),
		AuthenticationData: []byte{},
	}

	privateKeyHex = "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
	signer = auth.NewSigner(privateKeyHex)

	// Sign the header
	signature, err := signer.SignBlobRequest(testHeader)
	assert.NoError(t, err)

	testHeader.AuthenticationData = signature

	err = authenticator.AuthenticateBlobRequest(testHeader)
	assert.Error(t, err)

}
