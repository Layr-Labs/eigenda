package v2_test

import (
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	privateKeyHex = "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
)

func TestAuthentication(t *testing.T) {
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)
	authenticator := auth.NewAuthenticator()

	accountId, err := signer.GetAccountID()
	assert.NoError(t, err)
	header := testHeader(t, accountId)

	// Sign the header
	signature, err := signer.SignBlobRequest(header)
	assert.NoError(t, err)

	header.Signature = signature

	err = authenticator.AuthenticateBlobRequest(header)
	assert.NoError(t, err)

}

func TestAuthenticationFail(t *testing.T) {
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)
	authenticator := auth.NewAuthenticator()

	accountId, err := signer.GetAccountID()
	assert.NoError(t, err)

	header := testHeader(t, accountId)

	wrongPrivateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
	signer = auth.NewLocalBlobRequestSigner(wrongPrivateKeyHex)

	// Sign the header
	signature, err := signer.SignBlobRequest(header)
	assert.NoError(t, err)

	header.Signature = signature

	err = authenticator.AuthenticateBlobRequest(header)
	assert.Error(t, err)
}

func TestNoopSignerFail(t *testing.T) {
	signer := auth.NewLocalNoopSigner()
	accountId, err := signer.GetAccountID()
	assert.EqualError(t, err, "noop signer cannot get accountID")

	header := testHeader(t, accountId)

	_, err = signer.SignBlobRequest(header)
	assert.EqualError(t, err, "noop signer cannot sign blob request")
}

func testHeader(t *testing.T, accountID string) *corev2.BlobHeader {
	var commitX, commitY fp.Element
	_, err := commitX.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	assert.NoError(t, err)
	_, err = commitY.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	assert.NoError(t, err)

	commitment := &encoding.G1Commitment{
		X: commitX,
		Y: commitY,
	}
	var lengthXA0, lengthXA1, lengthYA0, lengthYA1 fp.Element
	_, err = lengthXA0.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	assert.NoError(t, err)
	_, err = lengthXA1.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	assert.NoError(t, err)
	_, err = lengthYA0.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	assert.NoError(t, err)
	_, err = lengthYA1.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531")
	assert.NoError(t, err)

	var lengthProof, lengthCommitment encoding.G2Commitment
	lengthProof.X.A0 = lengthXA0
	lengthProof.X.A1 = lengthXA1
	lengthProof.Y.A0 = lengthYA0
	lengthProof.Y.A1 = lengthYA1

	lengthCommitment = lengthProof

	return &corev2.BlobHeader{
		BlobVersion: 0,
		BlobCommitments: encoding.BlobCommitments{
			Commitment:       commitment,
			LengthCommitment: &lengthCommitment,
			LengthProof:      &lengthProof,
			Length:           50,
		},
		QuorumNumbers: []core.QuorumID{0, 1},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         accountID,
			ReservationPeriod: 5,
			CumulativePayment: big.NewInt(100),
		},
		Signature: []byte{},
	}
}

func TestAuthenticatePaymentStateRequestValid(t *testing.T) {
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)
	authenticator := auth.NewAuthenticator()

	signature, err := signer.SignPaymentStateRequest()
	assert.NoError(t, err)

	accountId, err := signer.GetAccountID()
	assert.NoError(t, err)

	err = authenticator.AuthenticatePaymentStateRequest(signature, accountId)
	assert.NoError(t, err)
}

func TestAuthenticatePaymentStateRequestInvalidSignatureLength(t *testing.T) {
	authenticator := auth.NewAuthenticator()

	err := authenticator.AuthenticatePaymentStateRequest([]byte{1, 2, 3}, "0x123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature length is unexpected")
}

func TestAuthenticatePaymentStateRequestInvalidPublicKey(t *testing.T) {
	authenticator := auth.NewAuthenticator()

	err := authenticator.AuthenticatePaymentStateRequest(make([]byte, 65), "not-hex-encoded")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to recover public key from signature")
}

func TestAuthenticatePaymentStateRequestSignatureMismatch(t *testing.T) {
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)
	authenticator := auth.NewAuthenticator()

	// Create a different signer with wrong private key
	wrongSigner := auth.NewLocalBlobRequestSigner("0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded")

	// Sign with wrong key
	accountId, err := signer.GetAccountID()
	assert.NoError(t, err)

	signature, err := wrongSigner.SignPaymentStateRequest()
	assert.NoError(t, err)

	err = authenticator.AuthenticatePaymentStateRequest(signature, accountId)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature doesn't match with provided public key")
}

func TestAuthenticatePaymentStateRequestCorruptedSignature(t *testing.T) {
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)
	authenticator := auth.NewAuthenticator()

	accountId, err := signer.GetAccountID()
	assert.NoError(t, err)

	hash := sha256.Sum256([]byte(accountId))
	signature, err := crypto.Sign(hash[:], signer.PrivateKey)
	assert.NoError(t, err)

	// Corrupt the signature
	signature[0] ^= 0x01

	err = authenticator.AuthenticatePaymentStateRequest(signature, accountId)
	assert.Error(t, err)
}
