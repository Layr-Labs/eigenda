package integration_test

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/dispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/proxy/clients/standard_client"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/config"
	enabled_apis "github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	proxy_metrics "github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/arbitrum_altda"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/rest"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/eigenda/verify"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	integration "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/Layr-Labs/eigensdk-go/logging"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

// TestProxyAPIsEnabledRestALTDA tests to ensure that the enabled APIs expression is
// getting respected by the REST ALTDA Server when wiring up a proxy application instance
// with just `op-generic` mode enabled.
//
// This test has been migrated from api/proxy/test/e2e/server_rest_test.go to use inabox infrastructure.
func TestProxyAPIsEnabledRestALTDA(t *testing.T) {
	// Create fresh test harness from global infrastructure
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	ctx := context.Background()

	// Create proxy config with only op-generic API enabled
	proxyConfig, err := createProxyConfig(
		&enabled_apis.RestApisEnabled{
			OpGenericCommitment: true,  // only op-generic enabled
			StandardCommitment:  false, // standard disabled
			OpKeccakCommitment:  false, // keccak disabled
		},
	)
	require.NoError(t, err)

	// Start proxy REST server
	restServer, restURL, cleanup, err := startProxyServer(ctx, globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	t.Logf("Proxy server started at %s", restURL)

	// Test that standard commitment mode is disabled (should return 403)
	standardClient := standard_client.New(&standard_client.Config{
		URL: restURL,
	})

	testBlob := []byte("hello world")
	t.Log("Attempting to set data using standard commitment (should fail with 403)...")
	_, err = standardClient.SetData(ctx, testBlob)
	require.Error(t, err)
	require.ErrorContains(t, err, "403")

	// Test that op-generic mode works (should succeed)
	opGenericClient := altda.NewDAClient(restURL, false, false)

	t.Log("Setting data using op-generic commitment (should succeed)...")
	daCommit, err := opGenericClient.SetInput(ctx, testBlob)
	require.NoError(t, err)

	t.Log("Getting data using op-generic commitment (should succeed)...")
	preimage, err := opGenericClient.GetInput(ctx, daCommit, 0)
	require.NoError(t, err)
	require.Equal(t, testBlob, preimage)

	t.Log("TestProxyAPIsEnabledRestALTDA completed successfully")

	// Verify the server is still running
	require.NotNil(t, restServer)
}

// createProxyConfig creates a proxy configuration that connects to the inabox disperser
func createProxyConfig(
	enabledAPIs *enabled_apis.RestApisEnabled,
) (config.AppConfig, error) {
	payloadClientConfig := clientsv2.PayloadClientConfig{
		PayloadPolynomialForm: codecs.PolynomialFormEval,
		BlobVersion:           0,
	}

	// Get the Ethereum RPC URL from the test config
	ethRPCURL := globalInfra.TestConfig.Deployers[0].RPC

	// Get the disperser API server address from the infrastructure
	disperserHostname, disperserPort, err := net.SplitHostPort(globalInfra.DisperserHarness.APIServerAddress)
	if err != nil {
		return config.AppConfig{}, fmt.Errorf("invalid disperser API server address: %w", err)
	}

	// Get SRS paths using the utility function
	g1Path, g2Path, g2TrailingPath, err := integration.GetSRSPaths()
	if err != nil {
		return config.AppConfig{}, fmt.Errorf("failed to determine SRS file paths: %w", err)
	}

	// Construct cache directory path from g1Path
	srsDir := filepath.Dir(g1Path)
	cacheDir := filepath.Join(srsDir, "SRSTables")

	// Define max blob length
	maxBlobLengthBytes := uint64(16 * 1024 * 1024) // 16 MiB

	builderConfig := builder.Config{
		StoreConfig: store.Config{
			AsyncPutWorkers:  0,
			BackendsToEnable: []common.EigenDABackend{common.V2EigenDABackend},
			DispersalBackend: common.V2EigenDABackend,
		},
		ClientConfigV2: common.ClientConfigV2{
			DisperserClientCfg: dispersal.DisperserClientConfig{
				Hostname:          disperserHostname,
				Port:              disperserPort,
				UseSecureGrpcFlag: false, // inabox uses insecure gRPC
			},
			PayloadDisperserCfg: dispersal.PayloadDisperserConfig{
				PayloadClientConfig:    payloadClientConfig,
				DisperseBlobTimeout:    5 * time.Minute,
				BlobCompleteTimeout:    5 * time.Minute,
				BlobStatusPollInterval: 1 * time.Second,
				ContractCallTimeout:    5 * time.Second,
			},
			RelayPayloadRetrieverCfg: payloadretrieval.RelayPayloadRetrieverConfig{
				PayloadClientConfig: payloadClientConfig,
				RelayTimeout:        5 * time.Second,
			},
			PutTries:                           3,
			MaxBlobSizeBytes:                   maxBlobLengthBytes,
			EigenDACertVerifierOrRouterAddress: globalInfra.TestConfig.EigenDA.CertVerifierRouter,
			EigenDADirectory:                   globalInfra.TestConfig.EigenDA.EigenDADirectory,
			EigenDANetwork:                     "", // Empty for inabox (custom network)
			RetrieversToEnable:                 []common.RetrieverType{common.RelayRetrieverType, common.ValidatorRetrieverType},
			ClientLedgerMode:                   clientledger.ClientLedgerModeReservationOnly,
			VaultMonitorInterval:               30 * time.Second,
		},
		KzgConfig: kzg.KzgConfig{
			G1Path:          g1Path,
			G2Path:          g2Path,
			G2TrailingPath:  g2TrailingPath,
			CacheDir:        cacheDir,
			SRSOrder:        encoding.SRSOrder,
			NumWorker:       uint64(runtime.GOMAXPROCS(0)), // #nosec G115
			SRSNumberToLoad: maxBlobLengthBytes / 32,
			LoadG2Points:    false, // not needed for inabox tests
		},
		MemstoreEnabled: false, // Use actual disperser, not memstore
		VerifierConfigV1: verify.Config{
			VerifyCerts:          true,
			RPCURL:               ethRPCURL,
			EthConfirmationDepth: 1,
			WaitForFinalization:  false,
			MaxBlobSizeBytes:     maxBlobLengthBytes,
		},
	}

	secretConfig := common.SecretConfigV2{
		SignerPaymentKey: integration.GetDefaultTestPayloadDisperserConfig().PrivateKey,
		EthRPCURL:        ethRPCURL,
	}

	return config.AppConfig{
		StoreBuilderConfig: builderConfig,
		SecretConfig:       secretConfig,
		EnabledServersConfig: &enabled_apis.EnabledServersConfig{
			Metric:        false,
			ArbCustomDA:   false,
			RestAPIConfig: *enabledAPIs,
		},
		MetricsSvrConfig: proxy_metrics.Config{},
		RestSvrCfg: rest.Config{
			Host:        "127.0.0.1",
			Port:        0, // let OS assign a free port
			APIsEnabled: enabledAPIs,
		},
		ArbCustomDASvrCfg: arbitrum_altda.Config{
			Host: "127.0.0.1",
			Port: 0, // let OS assign a free port

		},
	}, nil
}

// startProxyServer starts a proxy REST server with the given configuration
func startProxyServer(
	ctx context.Context,
	logger logging.Logger,
	appConfig config.AppConfig,
) (*rest.Server, string, func(), error) {
	if err := appConfig.Check(); err != nil {
		return nil, "", nil, fmt.Errorf("invalid app config: %w", err)
	}

	metrics := proxy_metrics.NewEmulatedMetricer()

	// Build the eth client for contract interactions
	ethClient, _, err := common.BuildEthClient(
		ctx,
		logger,
		appConfig.SecretConfig.EthRPCURL,
		appConfig.StoreBuilderConfig.ClientConfigV2.EigenDANetwork,
	)
	if err != nil {
		return nil, "", nil, fmt.Errorf("build eth client: %w", err)
	}

	// Build storage managers
	certMgr, keccakMgr, err := builder.BuildManagers(
		ctx,
		logger,
		metrics,
		appConfig.StoreBuilderConfig,
		appConfig.SecretConfig,
		nil,
		ethClient,
	)
	if err != nil {
		return nil, "", nil, fmt.Errorf("build storage managers: %w", err)
	}

	// Create compatibility config
	compatibilityCfg, err := common.NewCompatibilityConfig(
		"test",
		"", // chainID not needed for inabox
		appConfig.StoreBuilderConfig.ClientConfigV2,
		false, // readOnlyMode
		appConfig.EnabledServersConfig.ToAPIStrings(),
	)
	if err != nil {
		return nil, "", nil, fmt.Errorf("new compatibility config: %w", err)
	}

	// Create and start REST server
	appConfig.RestSvrCfg.CompatibilityCfg = compatibilityCfg
	restServer := rest.NewServer(appConfig.RestSvrCfg, certMgr, keccakMgr, logger, metrics)
	router := mux.NewRouter()
	restServer.RegisterRoutes(router)

	if err := restServer.Start(router); err != nil {
		return nil, "", nil, fmt.Errorf("start proxy server: %w", err)
	}

	// Get the actual port assigned by the OS
	port := restServer.Port()
	restURL := fmt.Sprintf("http://127.0.0.1:%d", port)

	cleanup := func() {
		if err := restServer.Stop(); err != nil {
			logger.Error("Failed to stop proxy server", "err", err)
		}
	}

	return restServer, restURL, cleanup, nil
}
