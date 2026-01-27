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
	}, logger)

	require.Error(t, err, "should fail to create signer with empty KeyID")
}

func TestKMSSignatureVerificationWithEmptyRegion(t *testing.T) {
	ctx := t.Context()

	// Try to create signer with empty Region - validation should catch it immediately
	_, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		Region:   "",
		Endpoint: localstackHost,
		KeyID:    "random_key_id",
	}, logger)

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
	}, logger)
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
		}, logger)
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
	}, logger)

	require.Error(t, err, "should fail to create signer with empty private key")
}

func TestLocalSignerWithInvalidPrivateKey(t *testing.T) {
	ctx := t.Context()

	_, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		PrivateKey: "invalid_hex",
	}, logger)

	require.Error(t, err, "should fail to create signer with invalid private key")
}

func TestLocalSignerPrivateKeyFormats(t *testing.T) {
	ctx := t.Context()

	// Test with 0x prefix - should fail
	_, err1 := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		PrivateKey: "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}, logger)
	require.Error(t, err1, "should fail with 0x prefix")

	// Test without 0x prefix - should succeed
	_, err2 := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		PrivateKey: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}, logger)
	require.NoError(t, err2, "should succeed without 0x prefix")
}

func TestLocalSignerWithBothKMSAndPrivateKey(t *testing.T) {
	ctx := t.Context()

	_, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		KeyID:      "some_key_id",
		PrivateKey: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		Region:     region,
	}, logger)

	require.Error(t, err, "should fail when both KeyID and PrivateKey are specified")
}

func TestNewKMSDispersalRequestSignerDirect(t *testing.T) {
	ctx := t.Context()
	_ = setupLocalStack(t)

	keyManager := kms.New(kms.Options{
		Region:       region,
		BaseEndpoint: aws.String(localstackHost),
	})

	// Create a test KMS key
	keyID, _ := createTestKMSKey(t, ctx, keyManager)

	// Test direct KMS factory function
	signer, err := NewKMSDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		Region:   region,
		Endpoint: localstackHost,
		KeyID:    keyID,
	}, logger)
	require.NoError(t, err, "failed to create KMS signer directly")
	require.NotNil(t, signer, "signer should not be nil")
}

func TestNewLocalDispersalRequestSignerDirect(t *testing.T) {
	// Generate a private key for testing
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate test private key")
	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))

	// Test direct local factory function
	signer, err := NewLocalDispersalRequestSigner(DispersalRequestSignerConfig{
		PrivateKey: privateKeyHex,
	})
	require.NoError(t, err, "failed to create local signer directly")
	require.NotNil(t, signer, "signer should not be nil")
}

func TestNewKMSDispersalRequestSignerErrors(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name        string
		config      DispersalRequestSignerConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "invalid_region_empty",
			config: DispersalRequestSignerConfig{
				KeyID:    "test-key",
				Region:   "",
				Endpoint: localstackHost,
			},
			expectError: true,
			errorMsg:    "should fail with empty region",
		},
		{
			name: "invalid_kms_endpoint",
			config: DispersalRequestSignerConfig{
				KeyID:    "non-existent-key",
				Region:   region,
				Endpoint: "http://invalid-endpoint:9999",
			},
			expectError: true,
			errorMsg:    "should fail with invalid endpoint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewKMSDispersalRequestSigner(ctx, tt.config, logger)
			if tt.expectError {
				require.Error(t, err, tt.errorMsg)
			} else {
				require.NoError(t, err, tt.errorMsg)
			}
		})
	}
}

func TestNewLocalDispersalRequestSignerErrors(t *testing.T) {
	tests := []struct {
		name        string
		config      DispersalRequestSignerConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "invalid_private_key_format",
			config: DispersalRequestSignerConfig{
				PrivateKey: "not-a-valid-hex-key",
			},
			expectError: true,
			errorMsg:    "should fail with invalid hex format",
		},
		{
			name: "empty_private_key",
			config: DispersalRequestSignerConfig{
				PrivateKey: "",
			},
			expectError: true,
			errorMsg:    "should fail with empty private key",
		},
		{
			name: "too_short_private_key",
			config: DispersalRequestSignerConfig{
				PrivateKey: "abc123",
			},
			expectError: true,
			errorMsg:    "should fail with too short private key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLocalDispersalRequestSigner(tt.config)
			if tt.expectError {
				require.Error(t, err, tt.errorMsg)
			} else {
				require.NoError(t, err, tt.errorMsg)
			}
		})
	}
}

func TestDefaultDispersalRequestSignerConfig(t *testing.T) {
	config := DefaultDispersalRequestSignerConfig()

	require.Equal(t, "", config.Endpoint, "default endpoint should be empty")
	require.Equal(t, "", config.KeyID, "default KeyID should be empty")
	require.Equal(t, "", config.PrivateKey, "default PrivateKey should be empty")
}

func TestDispersalRequestSignerConfigVerify(t *testing.T) {
	tests := []struct {
		name        string
		config      DispersalRequestSignerConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid_kms_config",
			config: DispersalRequestSignerConfig{
				KeyID:  "test-key",
				Region: "us-east-1",
			},
			expectError: false,
			errorMsg:    "valid KMS config should pass",
		},
		{
			name: "valid_local_config",
			config: DispersalRequestSignerConfig{
				PrivateKey: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			},
			expectError: false,
			errorMsg:    "valid local config should pass",
		},
		{
			name: "both_keyid_and_privatekey",
			config: DispersalRequestSignerConfig{
				KeyID:      "test-key",
				PrivateKey: "test-private-key",
				Region:     "us-east-1",
			},
			expectError: true,
			errorMsg:    "should fail when both KeyID and PrivateKey are provided",
		},
		{
			name: "neither_keyid_nor_privatekey",
			config: DispersalRequestSignerConfig{
				Region: "us-east-1",
			},
			expectError: true,
			errorMsg:    "should fail when neither KeyID nor PrivateKey is provided",
		},
		{
			name: "kms_without_region",
			config: DispersalRequestSignerConfig{
				KeyID: "test-key",
			},
			expectError: true,
			errorMsg:    "should fail when using KMS without region",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Verify()
			if tt.expectError {
				require.Error(t, err, tt.errorMsg)
			} else {
				require.NoError(t, err, tt.errorMsg)
			}
		})
	}
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
	}, logger)
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
		}, logger)
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

func TestKMSSignerEdgeCases(t *testing.T) {
	ctx := t.Context()
	_ = setupLocalStack(t)

	keyManager := kms.New(kms.Options{
		Region:       region,
		BaseEndpoint: aws.String(localstackHost),
	})

	// Create a test KMS key
	keyID, _ := createTestKMSKey(t, ctx, keyManager)

	signer, err := NewKMSDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		Region:   region,
		Endpoint: localstackHost,
		KeyID:    keyID,
	}, logger)
	require.NoError(t, err, "failed to create KMS signer")

	// Note: nil request test omitted as it would cause panic in hashing function,
	// which is expected behavior (caller should not pass nil)

	// Test with cancelled context
	t.Run("cancelled_context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		rand := random.NewTestRandom()
		request := auth.RandomStoreChunksRequest(rand)
		request.Signature = nil

		_, err := signer.SignStoreChunksRequest(cancelledCtx, request)
		require.Error(t, err, "should fail with cancelled context")
	})
}

func TestLocalSignerEdgeCases(t *testing.T) {
	ctx := t.Context()

	// Generate a private key for testing
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate test private key")
	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))

	signer, err := NewLocalDispersalRequestSigner(DispersalRequestSignerConfig{
		PrivateKey: privateKeyHex,
	})
	require.NoError(t, err, "failed to create local signer")

	// Note: nil request test omitted as it would cause panic in hashing function,
	// which is expected behavior (caller should not pass nil)

	// Test with cancelled context (should still work for local signing)
	t.Run("cancelled_context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		rand := random.NewTestRandom()
		request := auth.RandomStoreChunksRequest(rand)
		request.Signature = nil

		// Local signing should work even with cancelled context since it doesn't use network
		signature, err := signer.SignStoreChunksRequest(cancelledCtx, request)
		require.NoError(t, err, "local signing should work with cancelled context")
		require.NotNil(t, signature, "signature should not be nil")
		require.NotEmpty(t, signature, "signature should not be empty")
	})
}

func TestSignerTypeAssertion(t *testing.T) {
	ctx := t.Context()

	// Test that KMS factory returns KMS signer
	t.Run("kms_signer_type", func(t *testing.T) {
		_ = setupLocalStack(t)
		keyManager := kms.New(kms.Options{
			Region:       region,
			BaseEndpoint: aws.String(localstackHost),
		})
		keyID, _ := createTestKMSKey(t, ctx, keyManager)

		signer, err := NewKMSDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
			Region:   region,
			Endpoint: localstackHost,
			KeyID:    keyID,
		}, logger)
		require.NoError(t, err, "failed to create KMS signer")

		// Verify it's the correct concrete type
		_, ok := signer.(*kmsRequestSigner)
		require.True(t, ok, "should be kmsRequestSigner type")
	})

	// Test that local factory returns local signer
	t.Run("local_signer_type", func(t *testing.T) {
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err, "failed to generate test private key")
		privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))

		signer, err := NewLocalDispersalRequestSigner(DispersalRequestSignerConfig{
			PrivateKey: privateKeyHex,
		})
		require.NoError(t, err, "failed to create local signer")

		// Verify it's the correct concrete type
		_, ok := signer.(*localRequestSigner)
		require.True(t, ok, "should be localRequestSigner type")
	})
}

func TestNewDispersalRequestSignerRouting(t *testing.T) {
	ctx := t.Context()

	// Test routing to KMS signer
	t.Run("routes_to_kms", func(t *testing.T) {
		_ = setupLocalStack(t)
		keyManager := kms.New(kms.Options{
			Region:       region,
			BaseEndpoint: aws.String(localstackHost),
		})
		keyID, _ := createTestKMSKey(t, ctx, keyManager)

		signer, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
			Region:   region,
			Endpoint: localstackHost,
			KeyID:    keyID,
		}, logger)
		require.NoError(t, err, "should route to KMS signer")

		// Verify it routed to the correct type
		_, ok := signer.(*kmsRequestSigner)
		require.True(t, ok, "should have routed to kmsRequestSigner")
	})

	// Test routing to local signer
	t.Run("routes_to_local", func(t *testing.T) {
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err, "failed to generate test private key")
		privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))

		signer, err := NewDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
			PrivateKey: privateKeyHex,
		}, logger)
		require.NoError(t, err, "should route to local signer")

		// Verify it routed to the correct type
		_, ok := signer.(*localRequestSigner)
		require.True(t, ok, "should have routed to localRequestSigner")
	})
}

func TestKMSSignerWithDefaultConfig(t *testing.T) {
	ctx := t.Context()
	_ = setupLocalStack(t)

	keyManager := kms.New(kms.Options{
		Region:       region,
		BaseEndpoint: aws.String(localstackHost),
	})

	keyID, _ := createTestKMSKey(t, ctx, keyManager)

	// Test KMS signer without custom endpoint (uses default AWS config loading)
	_, err := NewKMSDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		Region: region,
		KeyID:  keyID,
		// No endpoint specified - should try to use default AWS config
	}, logger)
	// This will fail in test environment but we're testing the code path
	require.Error(t, err, "should fail to load default AWS config in test environment")
}

func TestKMSRegionOverride(t *testing.T) {
	ctx := t.Context()
	_ = setupLocalStack(t)

	keyManager := kms.New(kms.Options{
		Region:       region,
		BaseEndpoint: aws.String(localstackHost),
	})

	keyID, _ := createTestKMSKey(t, ctx, keyManager)

	// Test KMS region override - KMSRegion should be used instead of Region
	signer, err := NewKMSDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		Region:    "us-west-2", // This should be ignored
		KMSRegion: region,      // This should be used
		Endpoint:  localstackHost,
		KeyID:     keyID,
	}, logger)
	require.NoError(t, err, "should create KMS signer with region override")

	kmsSigner, ok := signer.(*kmsRequestSigner)
	require.True(t, ok, "should be kmsRequestSigner type")
	require.NotNil(t, kmsSigner.multiRegionSigner, "multi-region signer should be initialized")

	// Test signing with KMS region override
	rand := random.NewTestRandom()
	request := auth.RandomStoreChunksRequest(rand)
	request.Signature = nil

	signature, err := signer.SignStoreChunksRequest(ctx, request)
	require.NoError(t, err, "should sign successfully with KMS region override")
	require.NotNil(t, signature, "signature should not be nil")
	require.NotEmpty(t, signature, "signature should not be empty")
}

func TestSingleRegionBackwardCompatibility(t *testing.T) {
	ctx := t.Context()
	_ = setupLocalStack(t)

	keyManager := kms.New(kms.Options{
		Region:       region,
		BaseEndpoint: aws.String(localstackHost),
	})

	keyID, _ := createTestKMSKey(t, ctx, keyManager)

	// Test single-region mode (no fallback regions)
	signer, err := NewKMSDispersalRequestSigner(ctx, DispersalRequestSignerConfig{
		Region:   region,
		Endpoint: localstackHost,
		KeyID:    keyID,
	}, logger)
	require.NoError(t, err, "should create single-region KMS signer")

	kmsSigner, ok := signer.(*kmsRequestSigner)
	require.True(t, ok, "should be kmsRequestSigner type")
	require.NotNil(t, kmsSigner.multiRegionSigner, "multi-region signer should be set even for single-region mode")

	// Test signing with single-region setup
	rand := random.NewTestRandom()
	request := auth.RandomStoreChunksRequest(rand)
	request.Signature = nil

	signature, err := signer.SignStoreChunksRequest(ctx, request)
	require.NoError(t, err, "should sign successfully with single-region signer")
	require.NotNil(t, signature, "signature should not be nil")
	require.NotEmpty(t, signature, "signature should not be empty")
}

