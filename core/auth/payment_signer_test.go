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
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
	signer, err := auth.NewPaymentSigner(privateKeyHex)
	require.NoError(t, err)

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
		err = auth.VerifyPaymentSignature(header, signature)
		assert.NoError(t, err)
	})

	t.Run("VerifyPaymentSignature_InvalidSignature", func(t *testing.T) {
		header := &commonpb.PaymentHeader{
			BinIndex:          1,
			CumulativePayment: []byte{0x01, 0x02, 0x03},
			AccountId:         "",
		}

		// Create an invalid signature
		invalidSignature := make([]byte, 65)
		err = auth.VerifyPaymentSignature(header, invalidSignature)
		assert.Error(t, err)
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

		err = auth.VerifyPaymentSignature(header, signature)
		assert.Error(t, err)
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
