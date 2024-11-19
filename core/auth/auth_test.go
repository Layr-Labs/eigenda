package auth_test

import (
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {

	// Make the authenticator
	authenticator := auth.NewAuthenticator()

	// Make the signer
	privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

	accountId, err := signer.GetAccountID()
	assert.NoError(t, err)

	testHeader := core.BlobAuthHeader{
		BlobCommitments:    encoding.BlobCommitments{},
		AccountID:          accountId,
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
	authenticator := auth.NewAuthenticator()

	// Make the signer
	privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

	accountId, err := signer.GetAccountID()
	assert.NoError(t, err)

	testHeader := core.BlobAuthHeader{
		BlobCommitments:    encoding.BlobCommitments{},
		AccountID:          accountId,
		Nonce:              rand.Uint32(),
		AuthenticationData: []byte{},
	}

	privateKeyHex = "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
	signer = auth.NewLocalBlobRequestSigner(privateKeyHex)

	// Sign the header
	signature, err := signer.SignBlobRequest(testHeader)
	assert.NoError(t, err)

	testHeader.AuthenticationData = signature

	err = authenticator.AuthenticateBlobRequest(testHeader)
	assert.Error(t, err)

}

func TestNoopSignerFail(t *testing.T) {
	signer := auth.NewLocalNoopSigner()
	accountId, err := signer.GetAccountID()
	assert.EqualError(t, err, "noop signer cannot get accountID")

	testHeader := core.BlobAuthHeader{
		BlobCommitments:    encoding.BlobCommitments{},
		AccountID:          accountId,
		Nonce:              rand.Uint32(),
		AuthenticationData: []byte{},
	}
	_, err = signer.SignBlobRequest(testHeader)
	assert.EqualError(t, err, "noop signer cannot sign blob request")
}
