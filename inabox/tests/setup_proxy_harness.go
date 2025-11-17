package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/config"
	"github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	proxy_metrics "github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/arbitrum_altda"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/rest"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/eigenda/verify"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gorilla/mux"
)

// ProxyHarnessConfig contains the configuration for setting up the proxy harness
type ProxyHarnessConfig struct {
	// Disperser endpoint for V2 API
	DisperserHostname string
	DisperserPort     string

	// Ethereum RPC endpoint
	EthRPCURL string

	// LocalStack S3 configuration
	S3BucketName   string
	LocalStackPort string

	// EigenDA contract addresses
	EigenDADirectory      string
	CertVerifierAddress   string
	ServiceManagerAddress string

	// Max blob size
	MaxBlobSizeBytes uint64

	// Which backends to enable
	BackendsToEnable []common.EigenDABackend
	DispersalBackend common.EigenDABackend

	// Enabled APIs
	EnabledAPIs *enablement.RestApisEnabled

	// Enable Arbitrum AltDA server
	EnableArbCustomDA bool
}

// ProxyHarness contains the proxy server components
type ProxyHarness struct {
	// REST server for EigenDA proxy
	RestServer *rest.Server

	// Arbitrum AltDA server (optional)
	ArbServer *arbitrum_altda.Server

	// Server addresses
	RestAddress string
	ArbAddress  string

	// Metrics
	Metrics *proxy_metrics.EmulatedMetricer

	// Configuration used
	Config ProxyHarnessConfig
}

// SetupProxyHarness creates and starts the proxy servers
func SetupProxyHarness(
	ctx context.Context,
	logger logging.Logger,
	config ProxyHarnessConfig,
) (*ProxyHarness, error) {
	logger.Info("Setting up Proxy Harness")

	// Set defaults
	if config.EnabledAPIs == nil {
		config.EnabledAPIs = &enablement.RestApisEnabled{
			Admin:               false,
			OpGenericCommitment: true,
			OpKeccakCommitment:  true,
			StandardCommitment:  true,
		}
	}

	if config.MaxBlobSizeBytes == 0 {
		config.MaxBlobSizeBytes = 1 * 1024 * 1024 // 1 MiB default
	}

	if len(config.BackendsToEnable) == 0 {
		// Enable both V1 and V2 by default
		config.BackendsToEnable = []common.EigenDABackend{
			common.V1EigenDABackend,
			common.V2EigenDABackend,
		}
	}

	if config.DispersalBackend == 0 {
		// Default to V2
		config.DispersalBackend = common.V2EigenDABackend
	}

	// Build proxy configuration
	appConfig, err := buildProxyConfig(config)
	if err != nil {
		return nil, fmt.Errorf("build proxy config: %w", err)
	}

	// Validate configuration
	if err := appConfig.Check(); err != nil {
		return nil, fmt.Errorf("invalid proxy config: %w", err)
	}

	// Build storage managers
	metrics := proxy_metrics.NewEmulatedMetricer()

	// Build eth client
	ethClient, chainID, err := common.BuildEthClient(
		ctx,
		logger,
		config.EthRPCURL,
		"", // EigenDANetwork not needed for inabox
	)
	if err != nil {
		return nil, fmt.Errorf("build eth client: %w", err)
	}

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
		return nil, fmt.Errorf("build storage managers: %w", err)
	}

	// Build compatibility config
	compatibilityCfg, err := common.NewCompatibilityConfig(
		"inabox",
		chainID,
		appConfig.StoreBuilderConfig.ClientConfigV2,
		false, // readOnlyMode
		appConfig.EnabledServersConfig.ToAPIStrings(),
	)
	if err != nil {
		return nil, fmt.Errorf("new compatibility config: %w", err)
	}

	harness := &ProxyHarness{
		Metrics: metrics,
		Config:  config,
	}

	// Start REST server
	if appConfig.EnabledServersConfig.RestAPIConfig.DAEndpointEnabled() {
		appConfig.RestSvrCfg.CompatibilityCfg = compatibilityCfg
		restServer := rest.NewServer(appConfig.RestSvrCfg, certMgr, keccakMgr, logger, metrics)

		router := mux.NewRouter()
		restServer.RegisterRoutes(router)

		if err := restServer.Start(router); err != nil {
			return nil, fmt.Errorf("start REST server: %w", err)
		}

		harness.RestServer = restServer
		harness.RestAddress = fmt.Sprintf("http://127.0.0.1:%d", restServer.Port())
		logger.Info("Proxy REST server started", "address", harness.RestAddress)
	}

	// Start Arbitrum AltDA server (optional)
	if config.EnableArbCustomDA {
		arbHandlers := arbitrum_altda.NewHandlers(certMgr, logger, compatibilityCfg)
		arbServer, err := arbitrum_altda.NewServer(ctx, &appConfig.ArbCustomDASvrCfg, arbHandlers)
		if err != nil {
			// Clean up REST server if arb server fails
			if harness.RestServer != nil {
				_ = harness.RestServer.Stop()
			}
			return nil, fmt.Errorf("create Arbitrum server: %w", err)
		}

		if err := arbServer.Start(); err != nil {
			// Clean up REST server if arb server fails
			if harness.RestServer != nil {
				_ = harness.RestServer.Stop()
			}
			return nil, fmt.Errorf("start Arbitrum server: %w", err)
		}

		harness.ArbServer = arbServer
		harness.ArbAddress = fmt.Sprintf("http://127.0.0.1:%d", arbServer.Port())
		logger.Info("Proxy Arbitrum AltDA server started", "address", harness.ArbAddress)
	}

	logger.Info("Proxy Harness setup complete")
	return harness, nil
}

// buildProxyConfig creates the proxy AppConfig from the harness config
func buildProxyConfig(harnessConfig ProxyHarnessConfig) (config.AppConfig, error) {
	// Get SRS paths using the utility function
	g1Path, g2Path, g2TrailingPath, err := getSRSPaths()
	if err != nil {
		return config.AppConfig{}, fmt.Errorf("failed to determine SRS file paths: %w", err)
	}

	// Construct cache directory path from g1Path
	srsDir := filepath.Dir(g1Path)
	cacheDir := filepath.Join(srsDir, "SRSTables")

	builderConfig := builder.Config{
		StoreConfig: store.Config{
			BackendsToEnable: harnessConfig.BackendsToEnable,
			DispersalBackend: harnessConfig.DispersalBackend,
		},
		ClientConfigV1: common.ClientConfigV1{
			EdaClientCfg: clients.EigenDAClientConfig{
				RPC:                      harnessConfig.DisperserHostname + ":" + harnessConfig.DisperserPort,
				StatusQueryTimeout:       15 * time.Minute,
				StatusQueryRetryInterval: 1 * time.Second,
				DisableTLS:               true, // inabox doesn't use TLS
				EthRpcUrl:                harnessConfig.EthRPCURL,
				SvcManagerAddr:           harnessConfig.ServiceManagerAddress,
			},
			MaxBlobSizeBytes: harnessConfig.MaxBlobSizeBytes,
			PutTries:         3,
		},
		VerifierConfigV1: verify.Config{
			VerifyCerts:          true,
			RPCURL:               harnessConfig.EthRPCURL,
			SvcManagerAddr:       harnessConfig.ServiceManagerAddress,
			EthConfirmationDepth: 0, // No need to wait for confirmations in inabox
			WaitForFinalization:  false,
			MaxBlobSizeBytes:     harnessConfig.MaxBlobSizeBytes,
		},
		KzgConfig: kzg.KzgConfig{
			G1Path:          g1Path,
			G2Path:          g2Path,
			G2TrailingPath:  g2TrailingPath,
			CacheDir:        cacheDir,
			SRSOrder:        encoding.SRSOrder,
			SRSNumberToLoad: harnessConfig.MaxBlobSizeBytes / 32,
			NumWorker:       uint64(runtime.GOMAXPROCS(0)),
			LoadG2Points:    false, // Not needed for inabox
		},
		ClientConfigV2: common.ClientConfigV2{
			DisperserClientCfg: clientsv2.DisperserClientConfig{
				Hostname:          harnessConfig.DisperserHostname,
				Port:              harnessConfig.DisperserPort,
				UseSecureGrpcFlag: false, // inabox doesn't use TLS
			},
			PayloadDisperserCfg: payloaddispersal.PayloadDisperserConfig{
				PayloadClientConfig:    *clientsv2.GetDefaultPayloadClientConfig(),
				DisperseBlobTimeout:    5 * time.Minute,
				BlobCompleteTimeout:    5 * time.Minute,
				BlobStatusPollInterval: 1 * time.Second,
				ContractCallTimeout:    5 * time.Second,
			},
			RelayPayloadRetrieverCfg: payloadretrieval.RelayPayloadRetrieverConfig{
				PayloadClientConfig: *clientsv2.GetDefaultPayloadClientConfig(),
				RelayTimeout:        5 * time.Second,
			},
			PutTries:                           3,
			MaxBlobSizeBytes:                   harnessConfig.MaxBlobSizeBytes,
			EigenDACertVerifierOrRouterAddress: harnessConfig.CertVerifierAddress,
			EigenDADirectory:                   harnessConfig.EigenDADirectory,
			RetrieversToEnable: []common.RetrieverType{
				common.RelayRetrieverType,
				common.ValidatorRetrieverType,
			},
			ClientLedgerMode:     clientledger.ClientLedgerModeReservationOnly,
			VaultMonitorInterval: 30 * time.Second,
		},
		MemstoreEnabled: false, // Use real disperser in inabox
	}

	// Configure S3 if provided
	if harnessConfig.S3BucketName != "" {
		builderConfig.S3Config = s3.Config{
			Bucket:          harnessConfig.S3BucketName,
			Path:            "",
			Endpoint:        "localhost:" + harnessConfig.LocalStackPort,
			EnableTLS:       false,
			AccessKeySecret: "localstack",
			AccessKeyID:     "localstack",
			CredentialType:  s3.CredentialTypeStatic,
		}
	}

	secretConfig := common.SecretConfigV2{
		SignerPaymentKey: "", // Will be set by test if needed
		EthRPCURL:        harnessConfig.EthRPCURL,
	}

	return config.AppConfig{
		StoreBuilderConfig: builderConfig,
		SecretConfig:       secretConfig,
		EnabledServersConfig: &enablement.EnabledServersConfig{
			Metric:        false,
			ArbCustomDA:   harnessConfig.EnableArbCustomDA,
			RestAPIConfig: *harnessConfig.EnabledAPIs,
		},
		MetricsSvrConfig: proxy_metrics.Config{},
		RestSvrCfg: rest.Config{
			Host:        "127.0.0.1",
			Port:        0, // Auto-assign port
			APIsEnabled: harnessConfig.EnabledAPIs,
		},
		ArbCustomDASvrCfg: arbitrum_altda.Config{
			Host: "127.0.0.1",
			Port: 0, // Auto-assign port
		},
	}, nil
}

// Cleanup shuts down the proxy servers
func (h *ProxyHarness) Cleanup(ctx context.Context, logger logging.Logger) {
	logger.Info("Cleaning up Proxy Harness")

	if h.RestServer != nil {
		if err := h.RestServer.Stop(); err != nil {
			logger.Error("Failed to stop REST server", "err", err)
		}
	}

	if h.ArbServer != nil {
		if err := h.ArbServer.Stop(); err != nil {
			logger.Error("Failed to stop Arbitrum server", "err", err)
		}
	}

	logger.Info("Proxy Harness cleanup complete")
}
