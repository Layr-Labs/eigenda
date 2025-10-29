package common_test

import (
	"testing"

	clients_v2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/stretchr/testify/require"
)

func validClientConfigV2() common.ClientConfigV2 {
	return common.ClientConfigV2{
		DisperserClientCfg: clients_v2.DisperserClientConfig{
			Hostname: "localhost",
			Port:     "8080",
		},
		PayloadDisperserCfg:                payloaddispersal.PayloadDisperserConfig{},
		RelayPayloadRetrieverCfg:           payloadretrieval.RelayPayloadRetrieverConfig{},
		ValidatorPayloadRetrieverCfg:       payloadretrieval.ValidatorPayloadRetrieverConfig{},
		PutTries:                           3,
		MaxBlobSizeBytes:                   1024 * 1024, // 1 MiB
		EigenDACertVerifierOrRouterAddress: "0x1234567890abcdef1234567890abcdef12345678",
		RelayConnectionPoolSize:            10,
		RetrieversToEnable:                 []common.RetrieverType{common.RelayRetrieverType},
		EigenDADirectory:                   "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		RBNRecencyWindowSize:               100,
	}
}

func TestNewCompatibilityConfig(t *testing.T) {
	t.Parallel()

	clientConfig := validClientConfigV2()
	version := "1.2.3"
	chainID := "12345"
	APIsEnabled := []string{"put", "get"}
	readOnly := false

	result, err := common.NewCompatibilityConfig(
		version,
		chainID,
		clientConfig,
		readOnly,
		APIsEnabled,
	)

	require.NoError(t, err)
	require.Equal(t, version, result.Version)
	require.Equal(t, chainID, result.ChainID)
	require.Equal(t, clientConfig.EigenDADirectory, result.DirectoryAddress)
	require.Equal(t, clientConfig.EigenDACertVerifierOrRouterAddress, result.CertVerifierAddress)
	require.Equal(t, clientConfig.RBNRecencyWindowSize, result.RecencyWindowSize)
	require.Equal(t, APIsEnabled, result.APIsEnabled)
	require.Equal(t, readOnly, result.ReadOnlyMode)
	require.Greater(t, result.MaxPayloadSizeBytes, uint32(0))
}

func TestNewCompatibilityConfigVersionPrefixRemoval(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		inputVersion    string
		expectedVersion string
	}{
		{
			name:            "lowercase v prefix",
			inputVersion:    "v1.2.3",
			expectedVersion: "1.2.3",
		},
		{
			name:            "uppercase V prefix",
			inputVersion:    "V1.2.3",
			expectedVersion: "1.2.3",
		},
		{
			name:            "no prefix",
			inputVersion:    "1.2.3",
			expectedVersion: "1.2.3",
		},
		{
			name:            "empty version",
			inputVersion:    "",
			expectedVersion: "",
		},
		{
			name:            "version with metadata",
			inputVersion:    "v2.4.0-43-g3b4f9f40",
			expectedVersion: "2.4.0-43-g3b4f9f40",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clientConfig := validClientConfigV2()
			result, err := common.NewCompatibilityConfig(
				tc.inputVersion,
				"12345",
				clientConfig,
				false,
				[]string{"arb"},
			)

			require.NoError(t, err)
			require.Equal(t, tc.expectedVersion, result.Version)
		})
	}
}

func TestNewCompatibilityConfigReadOnlyMode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		readOnlyMode bool
	}{
		{
			name:         "read-only mode enabled",
			readOnlyMode: true,
		},
		{
			name:         "read-only mode disabled",
			readOnlyMode: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clientConfig := validClientConfigV2()
			result, err := common.NewCompatibilityConfig(
				"1.0.0",
				"12345",
				clientConfig,
				tc.readOnlyMode,
				[]string{"put"},
			)

			require.NoError(t, err)
			require.Equal(t, tc.readOnlyMode, result.ReadOnlyMode)
		})
	}
}

func TestNewCompatibilityConfigAPIsEnabled(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		APIsEnabled []string
	}{
		{
			name:        "single API",
			APIsEnabled: []string{"arb"},
		},
		{
			name:        "multiple APIs",
			APIsEnabled: []string{"arb", "op-generic", "standard"},
		},
		{
			name:        "no APIs",
			APIsEnabled: []string{},
		},
		{
			name:        "nil APIs",
			APIsEnabled: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clientConfig := validClientConfigV2()
			result, err := common.NewCompatibilityConfig(
				"1.0.0",
				"12345",
				clientConfig,
				false,
				tc.APIsEnabled,
			)

			require.NoError(t, err)
			require.Equal(t, tc.APIsEnabled, result.APIsEnabled)
		})
	}
}

func TestNewCompatibilityConfigChainID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		chainID string
	}{
		{
			name:    "numeric chain ID",
			chainID: "12345",
		},
		{
			name:    "empty chain ID (memstore)",
			chainID: "",
		},
		{
			name:    "mainnet chain ID",
			chainID: "1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clientConfig := validClientConfigV2()
			result, err := common.NewCompatibilityConfig(
				"1.0.0",
				tc.chainID,
				clientConfig,
				false,
				[]string{"arb"},
			)

			require.NoError(t, err)
			require.Equal(t, tc.chainID, result.ChainID)
		})
	}
}

func TestNewCompatibilityConfigMaxPayloadSize(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		maxBlobSizeBytes uint64
		wantErr          bool
	}{
		{
			name:             "valid blob size",
			maxBlobSizeBytes: 1024 * 1024, // 1 MiB
			wantErr:          false,
		},
		{
			name:             "larger blob size",
			maxBlobSizeBytes: 16 * 1024 * 1024, // 16 MiB
			wantErr:          false,
		},
		{
			name:             "zero blob size",
			maxBlobSizeBytes: 0,
			wantErr:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clientConfig := validClientConfigV2()
			clientConfig.MaxBlobSizeBytes = tc.maxBlobSizeBytes

			result, err := common.NewCompatibilityConfig(
				"1.0.0",
				"12345",
				clientConfig,
				false,
				[]string{"arb"},
			)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Greater(t, result.MaxPayloadSizeBytes, uint32(0))
				// The exact calculation is done by codec.BlobSymbolsToMaxPayloadSize
				// We just verify it's a reasonable value relative to input
				require.LessOrEqual(t, result.MaxPayloadSizeBytes, uint32(tc.maxBlobSizeBytes))
			}
		})
	}
}

func TestNewCompatibilityConfigContractAddresses(t *testing.T) {
	t.Parallel()

	directoryAddr := "0x1111111111111111111111111111111111111111"
	certVerifierAddr := "0x2222222222222222222222222222222222222222"

	clientConfig := validClientConfigV2()
	clientConfig.EigenDADirectory = directoryAddr
	clientConfig.EigenDACertVerifierOrRouterAddress = certVerifierAddr

	result, err := common.NewCompatibilityConfig(
		"1.0.0",
		"12345",
		clientConfig,
		false,
		[]string{"arb"},
	)

	require.NoError(t, err)
	require.Equal(t, directoryAddr, result.DirectoryAddress)
	require.Equal(t, certVerifierAddr, result.CertVerifierAddress)
}

func TestNewCompatibilityConfigRecencyWindow(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                  string
		recencyWindowSize     uint64
		expectedRecencyWindow uint64
	}{
		{
			name:                  "standard recency window",
			recencyWindowSize:     100,
			expectedRecencyWindow: 100,
		},
		{
			name:                  "zero recency window (disabled)",
			recencyWindowSize:     0,
			expectedRecencyWindow: 0,
		},
		{
			name:                  "large recency window",
			recencyWindowSize:     10000,
			expectedRecencyWindow: 10000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clientConfig := validClientConfigV2()
			clientConfig.RBNRecencyWindowSize = tc.recencyWindowSize

			result, err := common.NewCompatibilityConfig(
				"1.0.0",
				"12345",
				clientConfig,
				false,
				[]string{"arb"},
			)

			require.NoError(t, err)
			require.Equal(t, tc.expectedRecencyWindow, result.RecencyWindowSize)
		})
	}
}
