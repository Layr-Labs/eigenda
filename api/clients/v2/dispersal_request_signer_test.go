package clients

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing"
	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

const (
	localstackPort = "4579"
	localstackHost = "http://0.0.0.0:4579"
	region         = "us-east-1"
)

var (
	logger = test.GetLogger()
)

// TODO: Good candidate to be extracted into test package as a utility
func setupLocalStack(t *testing.T) *testbed.LocalStackContainer {
	t.Helper()

	deployLocalStack := (os.Getenv("DEPLOY_LOCALSTACK") != "false")
	if !deployLocalStack {
		return nil
	}

	ctx := t.Context()
	localstackContainer, err := testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       localstackPort,
		Services:       []string{"kms"},
		Logger:         logger,
	})
	require.NoError(t, err, "failed to start localstack container")

	t.Cleanup(func() {
		logger.Info("Stopping localstack container")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = localstackContainer.Terminate(ctx)
	})

	return localstackContainer
}

func createTestKMSKey(
	t *testing.T, ctx context.Context, keyManager *kms.Client,
) (keyID string, publicAddress gethcommon.Address) {
	t.Helper()

	createKeyOutput, err := keyManager.CreateKey(ctx, &kms.CreateKeyInput{
		KeySpec:  types.KeySpecEccSecgP256k1,
		KeyUsage: types.KeyUsageTypeSignVerify,
	})
	require.NoError(t, err, "failed to create KMS key")

	keyID = *createKeyOutput.KeyMetadata.KeyId

	key, err := aws2.LoadPublicKeyKMS(ctx, keyManager, keyID)
	require.NoError(t, err, "failed to load public key from KMS")

	publicAddress = crypto.PubkeyToAddress(*key)
	return keyID, publicAddress
}

func TestKMSSignatureVerificationWithEmptyKeyID(t *testing.T) {
	ctx := t.Context()

	// Try to create signer with empty KeyID - validation should catch it immediately
	_, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		Region:   region,
		Endpoint: localstackHost,
		KeyID:    "",
	})

	require.Error(t, err, "should fail to create signer with empty KeyID")
}

func TestKMSSignatureVerificationWithEmptyRegion(t *testing.T) {
	ctx := t.Context()

	// Try to create signer with empty Region - validation should catch it immediately
	_, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		Region:   "",
		Endpoint: localstackHost,
		KeyID:    "random_key_id",
	})

	require.Error(t, err, "should fail to create signer with empty Region")
}

func TestKMSSignatureVerification(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()
	_ = setupLocalStack(t)

	keyManager := kms.New(kms.Options{
		Region:       region,
		BaseEndpoint: aws.String(localstackHost),
	})

	// Create a test KMS key
	keyID, publicAddress := createTestKMSKey(t, ctx, keyManager)

	// Create signer and request for all test scenarios
	signer, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		Region:   region,
		Endpoint: localstackHost,
		KeyID:    keyID,
	})
	require.NoError(t, err, "failed to create dispersal request signer")

	request := auth.RandomStoreChunksRequest(rand)
	request.Signature = nil

	// Sign the request
	validSignature, err := signer.SignStoreChunksRequest(ctx, request)
	require.NoError(t, err, "failed to sign store chunks request")

	// Table-driven test scenarios
	tests := []struct {
		name             string
		setupRequest     func() *grpc.StoreChunksRequest
		expectError      bool
		expectNilHash    bool
		errorDescription string
	}{
		{
			name: "valid_signature",
			setupRequest: func() *grpc.StoreChunksRequest {
				// Use the same request with valid signature
				req := &grpc.StoreChunksRequest{
					Batch:       request.GetBatch(),
					DisperserID: request.GetDisperserID(),
					Timestamp:   request.GetTimestamp(),
					Signature:   validSignature,
				}
				return req
			},
			expectError:      false,
			expectNilHash:    false,
			errorDescription: "valid signature should verify successfully",
		},
		{
			name: "corrupted_signature",
			setupRequest: func() *grpc.StoreChunksRequest {
				// Use the same request data with corrupted signature
				badSignature := make([]byte, len(validSignature))
				copy(badSignature, validSignature)
				badSignature[10] = badSignature[10] + 1
				req := &grpc.StoreChunksRequest{
					Batch:       request.GetBatch(),
					DisperserID: request.GetDisperserID(),
					Timestamp:   request.GetTimestamp(),
					Signature:   badSignature,
				}
				return req
			},
			expectError:      true,
			expectNilHash:    true,
			errorDescription: "corrupted signature should fail verification",
		},
		{
			name: "modified_request",
			setupRequest: func() *grpc.StoreChunksRequest {
				// Modify request data but use valid signature
				req := &grpc.StoreChunksRequest{
					Batch:       request.GetBatch(),
					DisperserID: request.GetDisperserID() + 1, // Modify disperser ID
					Timestamp:   request.GetTimestamp(),
					Signature:   validSignature,
				}
				return req
			},
			expectError:      true,
			expectNilHash:    true,
			errorDescription: "modified request should fail verification with valid signature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRequest := tt.setupRequest()

			hash, err := auth.VerifyStoreChunksRequest(publicAddress, testRequest)

			if tt.expectError {
				require.Error(t, err, tt.errorDescription)
			} else {
				require.NoError(t, err, tt.errorDescription)
			}

			if tt.expectNilHash {
				require.Nil(t, hash, "hash should be nil for failed verification")
			} else {
				require.NotNil(t, hash, "hash should not be nil for successful verification")
				// Verify hash matches expected
				expectedHash, err := hashing.HashStoreChunksRequest(testRequest)
				require.NoError(t, err, "failed to compute expected hash")
				require.Equal(t, expectedHash, hash, "computed hash should match expected hash")
			}
		})
	}

	// Test with a different KMS key to ensure multiple keys work
	t.Run("multiple_keys", func(t *testing.T) {
		keyID2, publicAddress2 := createTestKMSKey(t, ctx, keyManager)
		signer2, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
			Region:   region,
			Endpoint: localstackHost,
			KeyID:    keyID2,
		})
		require.NoError(t, err, "failed to create second dispersal request signer")

		request2 := auth.RandomStoreChunksRequest(rand)
		request2.Signature = nil

		signature2, err := signer2.SignStoreChunksRequest(ctx, request2)
		require.NoError(t, err, "failed to sign request with second key")

		request2.Signature = signature2
		hash, err := auth.VerifyStoreChunksRequest(publicAddress2, request2)
		require.NoError(t, err, "second key signature verification should succeed")
		require.NotNil(t, hash, "hash should not be nil for valid second key signature")
	})
}

func TestLocalSignerWithEmptyPrivateKey(t *testing.T) {
	ctx := t.Context()

	_, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		PrivateKey: "",
	})

	require.Error(t, err, "should fail to create signer with empty private key")
}

func TestLocalSignerWithInvalidPrivateKey(t *testing.T) {
	ctx := t.Context()

	_, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		PrivateKey: "invalid_hex",
	})

	require.Error(t, err, "should fail to create signer with invalid private key")
}

func TestLocalSignerPrivateKeyFormats(t *testing.T) {
	ctx := t.Context()

	// Test with 0x prefix - should fail
	_, err1 := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		PrivateKey: "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	})
	require.Error(t, err1, "should fail with 0x prefix")

	// Test without 0x prefix - should succeed
	_, err2 := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		PrivateKey: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	})
	require.NoError(t, err2, "should succeed without 0x prefix")
}

func TestLocalSignerWithBothKMSAndPrivateKey(t *testing.T) {
	ctx := t.Context()

	_, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		KeyID:      "some_key_id",
		PrivateKey: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		Region:     region,
	})

	require.Error(t, err, "should fail when both KeyID and PrivateKey are specified")
}

func TestLocalSignerSignatureVerification(t *testing.T) {
	ctx := t.Context()
	rand := random.NewTestRandom()

	// Generate a private key for testing
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate test private key")

	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))

	// Create signer with private key
	signer, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		PrivateKey: privateKeyHex,
	})
	require.NoError(t, err, "failed to create local dispersal request signer")

	request := auth.RandomStoreChunksRequest(rand)
	request.Signature = nil

	// Sign the request
	validSignature, err := signer.SignStoreChunksRequest(ctx, request)
	require.NoError(t, err, "failed to sign store chunks request")

	// Table-driven test scenarios
	tests := []struct {
		name             string
		setupRequest     func() *grpc.StoreChunksRequest
		expectError      bool
		expectNilHash    bool
		errorDescription string
	}{
		{
			name: "valid_signature",
			setupRequest: func() *grpc.StoreChunksRequest {
				req := &grpc.StoreChunksRequest{
					Batch:       request.GetBatch(),
					DisperserID: request.GetDisperserID(),
					Timestamp:   request.GetTimestamp(),
					Signature:   validSignature,
				}
				return req
			},
			expectError:      false,
			expectNilHash:    false,
			errorDescription: "valid signature should verify successfully",
		},
		{
			name: "corrupted_signature",
			setupRequest: func() *grpc.StoreChunksRequest {
				badSignature := make([]byte, len(validSignature))
				copy(badSignature, validSignature)
				badSignature[10] = badSignature[10] + 1
				req := &grpc.StoreChunksRequest{
					Batch:       request.GetBatch(),
					DisperserID: request.GetDisperserID(),
					Timestamp:   request.GetTimestamp(),
					Signature:   badSignature,
				}
				return req
			},
			expectError:      true,
			expectNilHash:    true,
			errorDescription: "corrupted signature should fail verification",
		},
		{
			name: "modified_request",
			setupRequest: func() *grpc.StoreChunksRequest {
				req := &grpc.StoreChunksRequest{
					Batch:       request.GetBatch(),
					DisperserID: request.GetDisperserID() + 1,
					Timestamp:   request.GetTimestamp(),
					Signature:   validSignature,
				}
				return req
			},
			expectError:      true,
			expectNilHash:    true,
			errorDescription: "modified request should fail verification with valid signature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRequest := tt.setupRequest()

			hash, err := auth.VerifyStoreChunksRequest(publicAddress, testRequest)

			if tt.expectError {
				require.Error(t, err, tt.errorDescription)
			} else {
				require.NoError(t, err, tt.errorDescription)
			}

			if tt.expectNilHash {
				require.Nil(t, hash, "hash should be nil for failed verification")
			} else {
				require.NotNil(t, hash, "hash should not be nil for successful verification")
				expectedHash, err := hashing.HashStoreChunksRequest(testRequest)
				require.NoError(t, err, "failed to compute expected hash")
				require.Equal(t, expectedHash, hash, "computed hash should match expected hash")
			}
		})
	}

	// Test with a different private key to ensure isolation
	t.Run("different_keys", func(t *testing.T) {
		privateKey2, err := crypto.GenerateKey()
		require.NoError(t, err, "failed to generate second test private key")

		publicAddress2 := crypto.PubkeyToAddress(privateKey2.PublicKey)

		signer2, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
			PrivateKey: fmt.Sprintf("%x", crypto.FromECDSA(privateKey2)),
		})
		require.NoError(t, err, "failed to create second local dispersal request signer")

		request2 := auth.RandomStoreChunksRequest(rand)
		request2.Signature = nil

		signature2, err := signer2.SignStoreChunksRequest(ctx, request2)
		require.NoError(t, err, "failed to sign request with second key")

		request2.Signature = signature2
		hash, err := auth.VerifyStoreChunksRequest(publicAddress2, request2)
		require.NoError(t, err, "second key signature verification should succeed")
		require.NotNil(t, hash, "hash should not be nil for valid second key signature")

		// Verify that first key cannot verify signature from second key
		_, err = auth.VerifyStoreChunksRequest(publicAddress, request2)
		require.Error(t, err, "first key should not verify signature from second key")
	})
}
