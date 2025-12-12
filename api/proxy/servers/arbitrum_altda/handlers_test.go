package arbitrum_altda

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	proxy_common "github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/test/mocks"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var testLogger = logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{})

// createMockCert creates a mock versioned certificate for testing
func createMockCert() *certs.VersionedCert {
	return &certs.VersionedCert{
		Version:        certs.V2VersionByte,
		SerializedCert: []byte("mock cert data"),
	}
}

// createSequencerMsg creates a valid sequencer message with the given DA Cert
// and an empty message header
func createSequencerMsg(cert *certs.VersionedCert) hexutil.Bytes {
	messageHeader := make([]byte, MessageHeaderOffset)
	arbCommit := commitments.NewArbCommitment(*cert)
	daCommit := arbCommit.Encode()
	fullMsg := append(messageHeader, daCommit...)
	return hexutil.Bytes(fullMsg)
}

// TestMethod_GetMaxMessageSize verifies that the handler returns the correct max message size
func TestMethod_GetMaxMessageSize(t *testing.T) {
	testMaxPayloadSize := uint32(500)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEigenDAManager := mocks.NewMockIEigenDAManager(ctrl)
	compatCfg := proxy_common.CompatibilityConfig{
		Version:             "1.0.0",
		MaxPayloadSizeBytes: testMaxPayloadSize,
	}
	handlers := NewHandlers(mockEigenDAManager, testLogger, false, compatCfg)

	result, err := handlers.GetMaxMessageSize(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, int(testMaxPayloadSize), result.MaxSize)

}

// TestMethod_GetSupportedHeaderBytes verifies that the handler returns the correct header bytes
func TestMethod_GetSupportedHeaderBytes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEigenDAManager := mocks.NewMockIEigenDAManager(ctrl)
	compatCfg := proxy_common.CompatibilityConfig{Version: "1.0.0", MaxPayloadSizeBytes: 100_000_000}
	handlers := NewHandlers(mockEigenDAManager, testLogger, false, compatCfg)

	result, err := handlers.GetSupportedHeaderBytes(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.HeaderBytes, 1)
	require.Len(t, result.HeaderBytes[0], 2)
	require.Equal(t, uint8(commitments.ArbCustomDAHeaderByte), result.HeaderBytes[0][0])
	require.Equal(t, commitments.EigenDALayerByte, result.HeaderBytes[0][1])
}

// TestMethod_Store verifies the Store handler behavior using table-driven tests
func TestMethod_Store(t *testing.T) {
	mockCert := createMockCert()

	tests := []struct {
		name             string
		payload          []byte
		timeout          hexutil.Uint64
		dispersalBackend proxy_common.EigenDABackend
		mockPutReturn    *certs.VersionedCert
		mockPutError     error
		expectPutCall    bool
		expectError      bool
		errorContains    string
		errorIs          error
		validateResult   func(t *testing.T, result *StoreResult)
	}{
		{
			name:             "Success",
			payload:          []byte("test payload data"),
			timeout:          hexutil.Uint64(60),
			dispersalBackend: proxy_common.V2EigenDABackend,
			mockPutReturn:    mockCert,
			mockPutError:     nil,
			expectPutCall:    true,
			expectError:      false,
			validateResult: func(t *testing.T, result *StoreResult) {
				require.NotNil(t, result)
				require.NotNil(t, result.SerializedDACert)
				daCommit := commitments.NewArbCommitment(*mockCert)
				expectedEncoding := daCommit.Encode()
				require.Equal(t, expectedEncoding, []byte(result.SerializedDACert))
			},
		},
		{
			name:             "Error - Empty Payload Provided by DA Client",
			payload:          []byte{},
			timeout:          hexutil.Uint64(60),
			dispersalBackend: proxy_common.V2EigenDABackend,
			expectPutCall:    false,
			expectError:      true,
			errorContains:    "empty rollup payload",
		},
		{
			name:             "Error - Wrong Backend Type Configured",
			payload:          []byte("test payload"),
			timeout:          hexutil.Uint64(60),
			dispersalBackend: proxy_common.V1EigenDABackend,
			expectPutCall:    false,
			expectError:      true,
			errorContains:    "expected EigenDAV2 backend",
		},
		{
			name:             "Error - Failover Requested by Client",
			payload:          []byte("test payload"),
			timeout:          hexutil.Uint64(60),
			dispersalBackend: proxy_common.V2EigenDABackend,
			mockPutError:     &api.ErrorFailover{},
			expectPutCall:    true,
			expectError:      true,
			errorIs:          ErrFallbackRequested,
		},
		{
			name:             "Error - Dispersal Failed",
			payload:          []byte("test payload"),
			timeout:          hexutil.Uint64(60),
			dispersalBackend: proxy_common.V2EigenDABackend,
			mockPutError:     errors.New("put failed"),
			expectPutCall:    true,
			expectError:      true,
			errorContains:    "put rollup payload",
		},
		{
			name:             "Error - Batch Too Large",
			payload:          []byte("test payload that exceeds 10 bytes"),
			timeout:          hexutil.Uint64(60),
			dispersalBackend: proxy_common.V2EigenDABackend,
			expectPutCall:    false,
			expectError:      true,
			errorIs:          ErrMessageTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEigenDAManager := mocks.NewMockIEigenDAManager(ctrl)
			// Set MaxPayloadSizeBytes to 10 for the "Batch Too Large" test, otherwise use a large value
			maxPayloadSize := uint32(1000)
			if tt.name == "Error - Batch Too Large" {
				maxPayloadSize = 10
			}
			compatCfg := proxy_common.CompatibilityConfig{Version: "1.0.0", MaxPayloadSizeBytes: maxPayloadSize}
			handlers := NewHandlers(mockEigenDAManager, testLogger, false, compatCfg)

			mockEigenDAManager.EXPECT().
				GetDispersalBackend().
				Return(tt.dispersalBackend)

			if tt.expectPutCall {
				mockEigenDAManager.EXPECT().
					Put(gomock.Any(), tt.payload, coretypes.CertSerializationABI).
					Return(tt.mockPutReturn, tt.mockPutError)
			}

			result, err := handlers.Store(context.Background(), tt.payload, tt.timeout)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
				if tt.errorIs != nil {
					require.True(t, errors.Is(err, tt.errorIs))
				}
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}
		})
	}
}

// TestRecoverPayload verifies the RecoverPayload handler behavior using table-driven tests
func TestRecoverPayload(t *testing.T) {
	mockCert := createMockCert()

	tests := []struct {
		name               string
		sequencerMsg       hexutil.Bytes
		mockGetReturn      []byte
		mockGetError       error
		processInvalidCert bool
		expectError        bool
		errorContains      string
		errorIs            error
		validateResult     func(t *testing.T, result *PayloadResult)
	}{
		{
			name:          "Success - Valid Certificate",
			sequencerMsg:  createSequencerMsg(mockCert),
			mockGetReturn: []byte("recovered payload"),
			mockGetError:  nil,
			expectError:   false,
			validateResult: func(t *testing.T, result *PayloadResult) {
				require.NotNil(t, result)
				require.Equal(t, []byte("recovered payload"), result.Payload)
			},
		},
		{
			name:          "Error - Sequencer Message Too Small",
			sequencerMsg:  hexutil.Bytes([]byte("too short")),
			expectError:   true,
			errorContains: "deserialize DA Cert",
		},
		{
			name: "Error - Wrong Custom DA Header Byte",
			sequencerMsg: func() hexutil.Bytes {
				messageHeader := make([]byte, MessageHeaderOffset)
				wrongHeaderCommit := []byte{0xFF, commitments.EigenDALayerByte}
				wrongHeaderCommit = append(wrongHeaderCommit, []byte("some cert data")...)
				return hexutil.Bytes(append(messageHeader, wrongHeaderCommit...))
			}(),
			expectError:   true,
			errorContains: "CustomDAHeader byte",
		},
		{
			name:          "Error - Get Failed",
			sequencerMsg:  createSequencerMsg(mockCert),
			mockGetError:  errors.New("get failed"),
			expectError:   true,
			errorContains: "get rollup payload",
		},
		{
			name:               "Error - Certificate Validation Error With ProcessInvalidCert",
			sequencerMsg:       createSequencerMsg(mockCert),
			mockGetError:       &coretypes.DerivationError{},
			processInvalidCert: true,
			expectError:        true,
			errorIs:            ErrCertValidationError,
		},
		{
			name:               "Error - Certificate Validation Without ProcessInvalidCert",
			sequencerMsg:       createSequencerMsg(mockCert),
			mockGetError:       &coretypes.DerivationError{},
			processInvalidCert: false,
			expectError:        true,
			errorContains:      "get rollup payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEigenDAManager := mocks.NewMockIEigenDAManager(ctrl)
			compatCfg := proxy_common.CompatibilityConfig{Version: "1.0.0"}
			handlers := NewHandlers(mockEigenDAManager, testLogger, tt.processInvalidCert, compatCfg)

			// Only expect Get call if sequencer message is valid
			if len(tt.sequencerMsg) > DACertOffset &&
				tt.sequencerMsg[MessageHeaderOffset] == commitments.ArbCustomDAHeaderByte {
				mockEigenDAManager.EXPECT().
					Get(gomock.Any(), gomock.Any(), coretypes.CertSerializationABI, gomock.Any()).
					Return(tt.mockGetReturn, tt.mockGetError)
			}

			batchNum := hexutil.Uint64(1)
			batchBlockHash := common.HexToHash("0x1234")

			result, err := handlers.RecoverPayload(context.Background(), batchNum, batchBlockHash, tt.sequencerMsg)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
				if tt.errorIs != nil {
					require.True(t, errors.Is(err, tt.errorIs))
				}
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}
		})
	}
}

// TestCollectPreimages verifies the CollectPreimages handler behavior using table-driven tests
func TestCollectPreimages(t *testing.T) {
	mockCert := createMockCert()

	tests := []struct {
		name           string
		sequencerMsg   hexutil.Bytes
		mockGetReturn  []byte
		mockGetError   error
		expectError    bool
		expectNil      bool
		errorContains  string
		validateResult func(t *testing.T, result *PreimagesResult, sequencerMsg hexutil.Bytes)
	}{
		{
			name:          "Success - Valid Preimages",
			sequencerMsg:  createSequencerMsg(mockCert),
			mockGetReturn: []byte("recovered payload"),
			mockGetError:  nil,
			expectError:   false,
			validateResult: func(t *testing.T, result *PreimagesResult, sequencerMsg hexutil.Bytes) {
				require.NotNil(t, result)
				require.NotNil(t, result.Preimages)

				// Verify preimage mapping
				certHash := crypto.Keccak256Hash(sequencerMsg[MessageHeaderOffset:])
				preimageMap, exists := result.Preimages[CustomDAPreimageType]
				require.True(t, exists)
				preimage, exists := preimageMap[certHash]
				require.True(t, exists)
				require.Equal(t, []byte("recovered payload"), preimage)
			},
		},
		{
			name:          "Error - Invalid Certificate",
			sequencerMsg:  hexutil.Bytes([]byte("too short")),
			expectError:   true,
			errorContains: "deserialize cert",
		},
		{
			name:         "Success - Derivation Error Returns Nil",
			sequencerMsg: createSequencerMsg(mockCert),
			mockGetError: &coretypes.DerivationError{},
			expectError:  false,
			expectNil:    true,
		},
		{
			name:          "Error - Get Failed With Non-Derivation Error",
			sequencerMsg:  createSequencerMsg(mockCert),
			mockGetError:  errors.New("generic error"),
			expectError:   true,
			errorContains: "get rollup payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEigenDAManager := mocks.NewMockIEigenDAManager(ctrl)
			compatCfg := proxy_common.CompatibilityConfig{Version: "1.0.0"}
			handlers := NewHandlers(mockEigenDAManager, testLogger, false, compatCfg)

			// Only expect Get call if sequencer message is valid
			if len(tt.sequencerMsg) > DACertOffset &&
				tt.sequencerMsg[MessageHeaderOffset] == commitments.ArbCustomDAHeaderByte {
				mockEigenDAManager.EXPECT().
					Get(gomock.Any(), gomock.Any(), coretypes.CertSerializationABI, gomock.Any()).
					Return(tt.mockGetReturn, tt.mockGetError)
			}

			batchNum := hexutil.Uint64(1)
			batchBlockHash := common.HexToHash("0x1234")

			result, err := handlers.CollectPreimages(context.Background(), batchNum, batchBlockHash, tt.sequencerMsg)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				if tt.expectNil {
					require.Nil(t, result)
				} else if tt.validateResult != nil {
					tt.validateResult(t, result, tt.sequencerMsg)
				}
			}
		})
	}
}

// TestGenerateCertificateValidityProof verifies the GenerateCertificateValidityProof handler
func TestGenerateCertificateValidityProof(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEigenDAManager := mocks.NewMockIEigenDAManager(ctrl)
	compatCfg := proxy_common.CompatibilityConfig{Version: "1.0.0"}
	handlers := NewHandlers(mockEigenDAManager, testLogger, false, compatCfg)

	certificate := hexutil.Bytes([]byte("some certificate"))

	result, err := handlers.GenerateCertificateValidityProof(context.Background(), certificate)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, hexutil.Bytes([]byte{}), result.Proof)
}

// TestCompatibilityConfig verifies the CompatibilityConfig handler
func TestCompatibilityConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEigenDAManager := mocks.NewMockIEigenDAManager(ctrl)

	expectedConfig := proxy_common.CompatibilityConfig{
		Version:             "1.2.3",
		ChainID:             "17000",
		DirectoryAddress:    "0x1234567890abcdef",
		CertVerifierAddress: "0xfedcba0987654321",
		MaxPayloadSizeBytes: 16777216,
		RecencyWindowSize:   100,
		APIsEnabled:         []string{"api1", "api2"},
	}

	handlers := NewHandlers(mockEigenDAManager, testLogger, false, expectedConfig)

	result, err := handlers.CompatibilityConfig(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, expectedConfig.Version, result.Version)
	require.Equal(t, expectedConfig.ChainID, result.ChainID)
	require.Equal(t, expectedConfig.DirectoryAddress, result.DirectoryAddress)
	require.Equal(t, expectedConfig.CertVerifierAddress, result.CertVerifierAddress)
	require.Equal(t, expectedConfig.MaxPayloadSizeBytes, result.MaxPayloadSizeBytes)
	require.Equal(t, expectedConfig.RecencyWindowSize, result.RecencyWindowSize)
	require.Equal(t, expectedConfig.APIsEnabled, result.APIsEnabled)
}

// TestDeserializeCertFromSequencerMsg tests the Sequencer Message -> DA Cert
// deserialization logic
func TestDeserializeCertFromSequencerMsg(t *testing.T) {
	mockCert := createMockCert()

	tests := []struct {
		name          string
		sequencerMsg  hexutil.Bytes
		expectError   bool
		errorContains string
		validateCert  func(t *testing.T, cert *certs.VersionedCert)
	}{
		{
			name:         "Success - Valid Message",
			sequencerMsg: createSequencerMsg(mockCert),
			expectError:  false,
			validateCert: func(t *testing.T, cert *certs.VersionedCert) {
				require.NotNil(t, cert)
				require.Equal(t, mockCert.Version, cert.Version)
			},
		},
		{
			name:          "Error - Message Too Short",
			sequencerMsg:  hexutil.Bytes(make([]byte, DACertOffset-1)),
			expectError:   true,
			errorContains: "expected to be",
		},
		{
			name: "Error - Wrong CustomDA Header Byte",
			sequencerMsg: func() hexutil.Bytes {
				messageHeader := make([]byte, MessageHeaderOffset)
				wrongCommit := []byte{0xFF, commitments.EigenDALayerByte, byte(certs.V2VersionByte)}
				wrongCommit = append(wrongCommit, []byte("cert data")...)
				return hexutil.Bytes(append(messageHeader, wrongCommit...))
			}(),
			expectError:   true,
			errorContains: "CustomDAHeader byte",
		},
		{
			name: "Error - Wrong EigenDA Layer Byte",
			sequencerMsg: func() hexutil.Bytes {
				messageHeader := make([]byte, MessageHeaderOffset)
				wrongCommit := []byte{commitments.ArbCustomDAHeaderByte, 0xFF, byte(certs.V2VersionByte)}
				wrongCommit = append(wrongCommit, []byte("cert data")...)
				return hexutil.Bytes(append(messageHeader, wrongCommit...))
			}(),
			expectError:   true,
			errorContains: "EigenDALayer byte",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEigenDAManager := mocks.NewMockIEigenDAManager(ctrl)
			compatCfg := proxy_common.CompatibilityConfig{Version: "1.0.0"}
			handlers := NewHandlers(mockEigenDAManager, testLogger, false, compatCfg).(*Handlers)

			cert, err := handlers.deserializeCertFromSequencerMsg(tt.sequencerMsg)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
				require.Nil(t, cert)
			} else {
				require.NoError(t, err)
				if tt.validateCert != nil {
					tt.validateCert(t, cert)
				}
			}
		})
	}
}
