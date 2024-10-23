package auth_test

import (
	"encoding/hex"
	"testing"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentSigner(t *testing.T) {
	// Generate a new private key for testing
	privateKey, err := crypto.GenerateKey()
	// publicKey := &privateKey.PublicKey
	require.NoError(t, err)

	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
	signer := auth.NewPaymentSigner(privateKeyHex)

	t.Run("SignBlobPayment", func(t *testing.T) {
		header := &commonpb.PaymentHeader{
			BinIndex:          1,
			CumulativePayment: []byte{0x01, 0x02, 0x03},
			AccountId:         "",
		}

		signature, err := signer.SignBlobPayment(header)
		require.NoError(t, err)
		assert.NotEmpty(t, signature)

		// Verify the signature
		isValid := auth.VerifyPaymentSignature(header, signature)
		assert.True(t, isValid)
	})

	t.Run("VerifyPaymentSignature_InvalidSignature", func(t *testing.T) {
		header := &commonpb.PaymentHeader{
			BinIndex:          1,
			CumulativePayment: []byte{0x01, 0x02, 0x03},
			AccountId:         "",
		}

		// Create an invalid signature
		invalidSignature := make([]byte, 65)
		isValid := auth.VerifyPaymentSignature(header, invalidSignature)
		assert.False(t, isValid)
	})

	t.Run("VerifyPaymentSignature_ModifiedHeader", func(t *testing.T) {
		header := &commonpb.PaymentHeader{
			BinIndex:          1,
			CumulativePayment: []byte{0x01, 0x02, 0x03},
			AccountId:         "",
		}

		signature, err := signer.SignBlobPayment(header)
		require.NoError(t, err)

		// Modify the header after signing
		header.BinIndex = 2

		isValid := auth.VerifyPaymentSignature(header, signature)
		assert.False(t, isValid)
	})
}

func TestNoopPaymentSigner(t *testing.T) {
	signer := auth.NewNoopPaymentSigner()

	t.Run("SignBlobRequest", func(t *testing.T) {
		_, err := signer.SignBlobPayment(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "noop signer cannot sign blob payment header")
	})
}
