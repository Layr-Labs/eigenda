package v2

import (
	"crypto/sha256"
	"math/big"
	"testing"

	corev1 "github.com/Layr-Labs/eigenda/core"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestSignBlobRequest(t *testing.T) {
	privateKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
	signer := NewLocalBlobRequestSigner(privateKey)
	accountID, err := signer.GetAccountID()
	require.NoError(t, err)
	require.Equal(t, "0x1aa8226f6d354380dDE75eE6B634875c4203e522", accountID)

	var commitX, commitY fp.Element
	_, err = commitX.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
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

	header := &core.BlobHeader{
		BlobCommitments: encoding.BlobCommitments{
			Commitment:       commitment,
			LengthCommitment: &lengthCommitment,
			LengthProof:      &lengthProof,
			Length:           48,
		},
		BlobVersion:   1,
		QuorumNumbers: []corev1.QuorumID{1, 2},
		PaymentMetadata: corev1.PaymentMetadata{
			AccountID:         accountID,
			CumulativePayment: big.NewInt(100),
			ReservationPeriod: 100,
		},
	}

	// Sign the blob request
	signature, err := signer.SignBlobRequest(header)
	require.NoError(t, err)
	require.NotNil(t, signature)

	// Verify the signature
	blobKey, err := header.BlobKey()
	require.NoError(t, err)

	// Recover the public key from the signature
	pubKey, err := crypto.SigToPub(blobKey[:], signature)
	require.NoError(t, err)

	// Verify that the recovered address matches the signer's address
	recoveredAddr := crypto.PubkeyToAddress(*pubKey).Hex()
	assert.Equal(t, accountID, recoveredAddr)
}

func TestSignPaymentStateRequest(t *testing.T) {
	privateKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
	signer := NewLocalBlobRequestSigner(privateKey)
	expectedAddr := "0x1aa8226f6d354380dDE75eE6B634875c4203e522"
	accountID, err := signer.GetAccountID()
	require.NoError(t, err)
	hash := sha256.Sum256([]byte(accountID))

	// Sign payment state request
	signature, err := signer.SignPaymentStateRequest()
	require.NoError(t, err)
	require.NotNil(t, signature)

	// Recover the public key from the signature
	pubKey, err := crypto.SigToPub(hash[:], signature)
	require.NoError(t, err)

	// Verify that the recovered address matches the signer's address
	recoveredAddr := crypto.PubkeyToAddress(*pubKey).Hex()
	assert.Equal(t, expectedAddr, recoveredAddr)
}

func TestNoopSigner(t *testing.T) {
	signer := NewLocalNoopSigner()

	t.Run("SignBlobRequest", func(t *testing.T) {
		sig, err := signer.SignBlobRequest(nil)
		assert.Error(t, err)
		assert.Nil(t, sig)
		assert.Equal(t, "noop signer cannot sign blob request", err.Error())
	})

	t.Run("SignPaymentStateRequest", func(t *testing.T) {
		sig, err := signer.SignPaymentStateRequest()
		assert.Error(t, err)
		assert.Nil(t, sig)
		assert.Equal(t, "noop signer cannot sign payment state request", err.Error())
	})

	t.Run("GetAccountID", func(t *testing.T) {
		accountID, err := signer.GetAccountID()
		assert.Error(t, err)
		assert.Empty(t, accountID)
		assert.Equal(t, "noop signer cannot get accountID", err.Error())
	})
}
