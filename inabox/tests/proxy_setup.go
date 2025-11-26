package integration

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/dispersal"
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
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigenda/api/proxy/test/testutils"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gorilla/mux"
)

type ProxyTestConfig struct {
	UseMemstore      bool
	EnabledRestAPIs  *enablement.RestApisEnabled
	Expiration       time.Duration
	MaxBlobLength    string
	WriteThreadCount int
	WriteOnCacheMiss bool
	// at most one of the below options should be true
	UseKeccak256ModeS3            bool
	UseS3Caching                  bool
	UseS3Fallback                 bool
	ErrorOnSecondaryInsertFailure bool

	ClientLedgerMode     clientledger.ClientLedgerMode
	VaultMonitorInterval time.Duration

	GlobalInfra *InfrastructureHarness
}

// NewProxyTestConfig returns a new ProxyTestConfig
func NewProxyTestConfig(globalInfra *InfrastructureHarness) ProxyTestConfig {
	return ProxyTestConfig{
		UseMemstore: false,
		EnabledRestAPIs: &enablement.RestApisEnabled{
			Admin:               false,
			OpGenericCommitment: true,
			OpKeccakCommitment:  true,
			StandardCommitment:  true,
		},
		Expiration:                    14 * 24 * time.Hour,
		UseKeccak256ModeS3:            false,
		UseS3Caching:                  false,
		UseS3Fallback:                 false,
		WriteThreadCount:              0,
		WriteOnCacheMiss:              false,
		ErrorOnSecondaryInsertFailure: false,
		ClientLedgerMode:              clientledger.ClientLedgerModeReservationOnly,
		VaultMonitorInterval:          30 * time.Second,

		GlobalInfra: globalInfra,
	}
}

// createProxyConfig creates a proxy configuration that connects to the inabox disperser
func CreateProxyConfig(
	testCfg ProxyTestConfig,
) (config.AppConfig, error) {
	payloadClientConfig := clientsv2.PayloadClientConfig{
		PayloadPolynomialForm: codecs.PolynomialFormEval,
		BlobVersion:           0,
	}

	// Get the Ethereum RPC URL from the test config
	ethRPCURL := testCfg.GlobalInfra.TestConfig.Deployers[0].RPC

	// Get the disperser API server address from the infrastructure
	disperserHostname, disperserPort, err := net.SplitHostPort(testCfg.GlobalInfra.DisperserHarness.APIServerAddress)
	if err != nil {
		return config.AppConfig{}, fmt.Errorf("invalid disperser API server address: %w", err)
	}

	// Get SRS paths using the utility function
	g1Path, g2Path, g2TrailingPath, err := getSRSPaths()
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
			AsyncPutWorkers:               testCfg.WriteThreadCount,
			BackendsToEnable:              []common.EigenDABackend{common.V2EigenDABackend},
			DispersalBackend:              common.V2EigenDABackend,
			WriteOnCacheMiss:              testCfg.WriteOnCacheMiss,
			ErrorOnSecondaryInsertFailure: testCfg.ErrorOnSecondaryInsertFailure,
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
			EigenDACertVerifierOrRouterAddress: testCfg.GlobalInfra.TestConfig.EigenDA.CertVerifierRouter,
			EigenDADirectory:                   testCfg.GlobalInfra.TestConfig.EigenDA.EigenDADirectory,
			EigenDANetwork:                     "", // Empty for inabox (custom network)
			RetrieversToEnable:                 []common.RetrieverType{common.RelayRetrieverType, common.ValidatorRetrieverType},
			ClientLedgerMode:                   testCfg.ClientLedgerMode,
			VaultMonitorInterval:               testCfg.VaultMonitorInterval,
		},
		KzgConfig: kzg.KzgConfig{
			G1Path:          g1Path,
			G2Path:          g2Path,
			G2TrailingPath:  g2TrailingPath,
			CacheDir:        cacheDir,
			SRSOrder:        encoding.SRSOrder,
			NumWorker:       uint64(runtime.GOMAXPROCS(0)), // #nosec G115
			SRSNumberToLoad: maxBlobLengthBytes / 32,
			LoadG2Points:    true,
		},
		MemstoreConfig: memconfig.NewSafeConfig(
			memconfig.Config{
				BlobExpiration:   testCfg.Expiration,
				MaxBlobSizeBytes: maxBlobLengthBytes,
			}),
		MemstoreEnabled: testCfg.UseMemstore,
		VerifierConfigV1: verify.Config{
			VerifyCerts:          true,
			RPCURL:               ethRPCURL,
			EthConfirmationDepth: 1,
			WaitForFinalization:  false,
			MaxBlobSizeBytes:     maxBlobLengthBytes,
		},
	}

	localstack := testCfg.GlobalInfra.DisperserHarness.LocalStack
	awsConfig := localstack.GetAWSClientConfig()
	awsEndpoint := strings.TrimPrefix(awsConfig.EndpointURL, "http://")
	s3Config := s3.Config{
		Bucket:          testCfg.GlobalInfra.DisperserHarness.S3Buckets.BlobStore,
		Path:            "",
		Endpoint:        awsEndpoint,
		EnableTLS:       false,
		AccessKeySecret: awsConfig.SecretAccessKey,
		AccessKeyID:     awsConfig.AccessKey,
		CredentialType:  s3.CredentialTypeStatic,
	}

	switch {
	case testCfg.UseKeccak256ModeS3:
		builderConfig.S3Config = s3Config
	case testCfg.UseS3Caching:
		builderConfig.StoreConfig.CacheTargets = []string{"S3"}
		builderConfig.S3Config = s3Config
	case testCfg.UseS3Fallback:
		builderConfig.StoreConfig.FallbackTargets = []string{"S3"}
		builderConfig.S3Config = s3Config
	}

	secretConfig := common.SecretConfigV2{
		SignerPaymentKey: GetDefaultTestPayloadDisperserConfig().PrivateKey,
		EthRPCURL:        ethRPCURL,
	}

	return config.AppConfig{
		StoreBuilderConfig: builderConfig,
		SecretConfig:       secretConfig,
		EnabledServersConfig: &enablement.EnabledServersConfig{
			Metric:        false,
			ArbCustomDA:   false,
			RestAPIConfig: *testCfg.EnabledRestAPIs,
		},
		MetricsSvrConfig: proxy_metrics.Config{
			Host: "127.0.0.1",
			Port: 0, // let OS assign a free port
		},
		RestSvrCfg: rest.Config{
			Host:        "127.0.0.1",
			Port:        0, // let OS assign a free port
			APIsEnabled: testCfg.EnabledRestAPIs,
		},
		ArbCustomDASvrCfg: arbitrum_altda.Config{
			Host: "127.0.0.1",
			Port: 0, // let OS assign a free port
		},
	}, nil
}

// startProxyServer starts a proxy REST server with the given configuration
func StartProxyServer(
	ctx context.Context,
	logger logging.Logger,
	appConfig config.AppConfig,
) (*testutils.TestSuite, func(), error) {
	if err := appConfig.Check(); err != nil {
		return nil, nil, fmt.Errorf("invalid app config: %w", err)
	}

	var (
		restServer *rest.Server
		arbServer  *arbitrum_altda.Server
		metrics    = proxy_metrics.NewEmulatedMetricer()
	)

	// Build the eth client for contract interactions
	ethClient, chainID, err := common.BuildEthClient(
		ctx,
		logger,
		appConfig.SecretConfig.EthRPCURL,
		appConfig.StoreBuilderConfig.ClientConfigV2.EigenDANetwork,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("build eth client: %w", err)
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
		return nil, nil, fmt.Errorf("build storage managers: %w", err)
	}

	// Create compatibility config
	compatibilityCfg, err := common.NewCompatibilityConfig(
		"test",
		chainID,
		appConfig.StoreBuilderConfig.ClientConfigV2,
		false, // readOnlyMode
		appConfig.EnabledServersConfig.ToAPIStrings(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("new compatibility config: %w", err)
	}

	if appConfig.EnabledServersConfig.RestAPIConfig.DAEndpointEnabled() {
		// Create and start REST server
		appConfig.RestSvrCfg.CompatibilityCfg = compatibilityCfg
		restServer = rest.NewServer(appConfig.RestSvrCfg, certMgr, keccakMgr, logger, metrics)
		router := mux.NewRouter()
		restServer.RegisterRoutes(router)
		if appConfig.StoreBuilderConfig.MemstoreEnabled {
			memconfig.NewHandlerHTTP(logger, appConfig.StoreBuilderConfig.MemstoreConfig).
				RegisterMemstoreConfigHandlers(router)
		}

		if err := restServer.Start(router); err != nil {
			return nil, nil, fmt.Errorf("start proxy server: %w", err)
		}
	}

	if appConfig.EnabledServersConfig.ArbCustomDA {
		arbHandlers := arbitrum_altda.NewHandlers(certMgr, logger, compatibilityCfg)
		arbServer, err = arbitrum_altda.NewServer(ctx, &appConfig.ArbCustomDASvrCfg, arbHandlers)
		if err != nil {
			return nil, nil, fmt.Errorf("create arbitrum server: %v", err.Error())
		}

		if err := arbServer.Start(); err != nil {
			return nil, nil, fmt.Errorf("start arbitrum server: %v", err.Error())
		}
	}

	cleanup := func() {
		if appConfig.EnabledServersConfig.RestAPIConfig.DAEndpointEnabled() {
			if err := restServer.Stop(); err != nil {
				logger.Error("failed to stop proxy server", "err", err)
			}
		}

		if appConfig.EnabledServersConfig.ArbCustomDA {
			if err := arbServer.Stop(); err != nil {
				logger.Error("failed to stop arb server", "err", err)
			}
		}
	}

	return &testutils.TestSuite{
		Ctx:        ctx,
		Log:        logger,
		RestServer: restServer,
		Metrics:    metrics,
		ArbServer:  arbServer,
	}, cleanup, nil
}
