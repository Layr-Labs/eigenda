package clients

import (
	"context"
	"testing"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	grpccommon "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	commonv1 "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalDispersalRequestSigner(t *testing.T) {
	// Test private key (secp256k1)
	testPrivateKeyHex := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

	// Create signer from hex
	signer, err := NewLocalDispersalRequestSignerFromHex(testPrivateKeyHex)
	require.NoError(t, err)
	assert.NotNil(t, signer)

	// Create a test request
	request := &grpc.StoreChunksRequest{
		Batch: &grpccommon.Batch{
			Header: &grpccommon.BatchHeader{
				BatchRoot:            []byte("test-batch-root"),
				ReferenceBlockNumber: 12345,
			},
			BlobCertificates: []*grpccommon.BlobCertificate{
				{
					BlobHeader: &grpccommon.BlobHeader{
						Version:       1,
						QuorumNumbers: []uint32{0, 1},
						Commitment: &commonv1.BlobCommitment{
							Commitment:       []byte("test-commitment"),
							LengthCommitment: []byte("length-commitment"),
							LengthProof:      []byte("length-proof"),
							Length:           123,
						},
						PaymentHeader: &grpccommon.PaymentHeader{
							AccountId:         "test-account",
							Timestamp:         1234567890,
							CumulativePayment: []byte("cumulative-payment"),
						},
					},
					Signature:  []byte("test-signature"),
					RelayKeys:  []uint32{1, 2, 3},
				},
			},
		},
		DisperserID: 42,
		Timestamp:   1234567890,
	}

	// Test signing
	signature, err := signer.SignStoreChunksRequest(context.Background(), request)
	require.NoError(t, err)
	assert.NotNil(t, signature)

	// Signature should be 65 bytes for secp256k1 (Ethereum format)
	assert.Equal(t, 65, len(signature), "Signature should be 65 bytes for secp256k1")

	// Recovery ID should be valid (0 or 1)
	recoveryID := signature[64]
	assert.True(t, recoveryID == 0 || recoveryID == 1, "Recovery ID should be 0 or 1")
}

func TestLocalDispersalRequestSignerInvalidKey(t *testing.T) {
	// Test with invalid private key
	invalidKey := "invalid-key"

	signer, err := NewLocalDispersalRequestSignerFromHex(invalidKey)
	assert.Error(t, err)
	assert.Nil(t, signer)
}

func TestLocalDispersalRequestSignerFactory(t *testing.T) {
	testPrivateKeyHex := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

	kmsConfig := common.KMSKeyConfig{
		Provider:      "local",
		PrivateKeyHex: testPrivateKeyHex,
	}

	signer, err := NewDispersalRequestSignerFromKMSConfig(
		context.Background(),
		kmsConfig,
		"us-east-1", // Not used for local provider
		"",          // Not used for local provider  
	)

	require.NoError(t, err)
	assert.NotNil(t, signer)

	// Ensure it's the correct type
	localSigner, ok := signer.(*localRequestSigner)
	assert.True(t, ok, "Signer should be of type localRequestSigner")
	assert.NotNil(t, localSigner.privateKey)
}