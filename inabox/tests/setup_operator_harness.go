package integration

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/pubip"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/common/version"
	"github.com/Layr-Labs/eigenda/core"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/payments/reservation/reservationvalidation"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/grpc"
	"github.com/Layr-Labs/eigensdk-go/logging"
	rpccalls "github.com/Layr-Labs/eigensdk-go/metrics/collectors/rpc_calls"
	blssignerTypes "github.com/Layr-Labs/eigensdk-go/signer/bls/types"
	"github.com/docker/go-units"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
)

// OperatorHarnessConfig contains the configuration for setting up the operator harness
type OperatorHarnessConfig struct {
	TestConfig *deploy.Config
	TestName   string
}

// OperatorHarness manages operator instances for integration tests
type OperatorHarness struct {
	ServersV2 []*grpc.ServerV2

	// Internal fields for operator management
	testConfig   *deploy.Config
	testName     string
	chainHarness *ChainHarness
	srsG1Path    string
	srsG2Path    string
}

// SetupOperatorHarness creates and initializes the operator harness
func SetupOperatorHarness(
	ctx context.Context,
	logger logging.Logger,
	chainHarness *ChainHarness,
	config *OperatorHarnessConfig,
) (*OperatorHarness, error) {
	harness := &OperatorHarness{
		ServersV2: make([]*grpc.ServerV2, 0),
	}

	// Store references we'll need
	harness.testConfig = config.TestConfig
	harness.testName = config.TestName
	harness.chainHarness = chainHarness

	// Start all operators
	if err := harness.StartOperators(ctx, logger); err != nil {
		return nil, err
	}

	return harness, nil
}

// operatorListeners holds the network listeners for a single operator
type operatorListeners struct {
	v1 grpc.Listeners
	v2 grpc.Listeners
}

// StartOperators starts all operator nodes configured in the test config
func (oh *OperatorHarness) StartOperators(ctx context.Context, logger logging.Logger) error {
	// Get SRS paths first - fail early if we can't find them
	g1Path, g2Path, _, err := getSRSPaths()
	if err != nil {
		return fmt.Errorf("failed to determine SRS file paths: %w", err)
	}

	// Store them in the harness for use by startOperator
	oh.srsG1Path = g1Path
	oh.srsG2Path = g2Path

	// Check that chain dependencies are available
	if oh.chainHarness == nil || oh.chainHarness.Anvil == nil {
		return fmt.Errorf("AnvilContainer is not initialized")
	}

	if oh.chainHarness.Churner.URL == "" {
		return fmt.Errorf("churner has not been started (ChurnerURL is empty)")
	}

	// Count how many operator configs exist
	operatorCount := 0
	for {
		operatorName := fmt.Sprintf("opr%d", operatorCount)
		if _, ok := oh.testConfig.Pks.EcdsaMap[operatorName]; !ok {
			break
		}
		operatorCount++
	}
	if operatorCount == 0 {
		return fmt.Errorf("no operators found in config")
	}

	logger.Info("Starting operators", "count", operatorCount)

	// Create listeners and start each operator
	for i := range operatorCount {
		v1Listeners, err := grpc.CreateListeners("0", "0")
		if err != nil {
			return fmt.Errorf("failed to create v1 listeners for operator %d: %w", i, err)
		}

		v2Listeners, err := grpc.CreateListeners("0", "0")
		if err != nil {
			v1Listeners.Close()
			return fmt.Errorf("failed to create v2 listeners for opersator %d: %w", i, err)
		}

		listeners := operatorListeners{
			v1: v1Listeners,
			v2: v2Listeners,
		}

		// Note: on success, the server takes ownership of the listeners and they will be closed when
		// the infrastructure harness calls Cleanup().
		serverV2, err := oh.startOperator(ctx, logger, i, listeners)
		if err != nil {
			// Close the listeners we just created since startOperator failed
			listeners.v1.Close()
			listeners.v2.Close()

			// Clean up any operators we've already started
			oh.stopAllOperators(logger)
			return fmt.Errorf("failed to start operator %d: %w", i, err)
		}

		oh.ServersV2 = append(oh.ServersV2, serverV2)
		logger.Info("Started operator", "index", i,
			"v2DispersalPort", serverV2.GetDispersalPort(),
			"v2RetrievalPort", serverV2.GetRetrievalPort())
	}

	return nil
}

// startOperator starts a single operator with the given index and pre-created listeners
// On success, the returned server takes ownership of the listeners and will close them
// when Stop() is called. On failure, the caller retains ownership of the listeners.
func (oh *OperatorHarness) startOperator(
	ctx context.Context,
	logger logging.Logger,
	operatorIndex int,
	listeners operatorListeners,
) (*grpc.ServerV2, error) {
	// Get operator's private key
	operatorName := fmt.Sprintf("opr%d", operatorIndex)

	// Check if operator exists in test config
	if oh.testConfig.Pks == nil || oh.testConfig.Pks.EcdsaMap == nil {
		return nil, fmt.Errorf("no private keys configured")
	}

	operatorKey, ok := oh.testConfig.Pks.EcdsaMap[operatorName]
	if !ok {
		return nil, fmt.Errorf("operator %s not found in config", operatorName)
	}

	// Get BLS key configuration
	blsKey, blsOk := oh.testConfig.Pks.BlsMap[operatorName]
	if !blsOk {
		return nil, fmt.Errorf("BLS key for %s not found in config", operatorName)
	}

	// Create logs directory
	// TODO(dmanc): If possible we should have a centralized place for creating loggers and injecting them into the config.
	logsDir := fmt.Sprintf("testdata/%s/logs", oh.testName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/operator_%d.log", logsDir, operatorIndex)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open operator log file: %w", err)
	}

	// Extract actual ports assigned by OS from the pre-created listeners
	v2DispersalPort := fmt.Sprintf("%d", listeners.v2.Dispersal.Addr().(*net.TCPAddr).Port)
	v2RetrievalPort := fmt.Sprintf("%d", listeners.v2.Retrieval.Addr().(*net.TCPAddr).Port)
	nodeApiPort := fmt.Sprintf("3710%d", operatorIndex)
	metricsPort := 3800 + operatorIndex

	// TODO(dmanc): The node config is quite a beast. This is a configuration that
	// passed the tests after a bunch of trial and error.
	// We really need better validation on the node constructor.

	// TODO(dmanc): In addition to loggers, we should have a centralized place for creating
	// configuration and injecting it into the harness config.

	reservationLedgerCacheConfig, err := reservationvalidation.NewReservationLedgerCacheConfig(
		1024,
		120*time.Second,
		ratelimit.OverfillOncePermitted,
		1*time.Second, // Matches controller and API server update interval
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation ledger cache config: %w", err)
	}

	nodeConfig := &node.Config{
		Hostname:                       "localhost",
		V2RetrievalPort:                v2RetrievalPort,
		V2DispersalPort:                v2DispersalPort,
		InternalV2RetrievalPort:        v2RetrievalPort,
		InternalV2DispersalPort:        v2DispersalPort,
		EnableNodeApi:                  true,
		NodeApiPort:                    nodeApiPort,
		EnableMetrics:                  true,
		MetricsPort:                    metricsPort,
		Timeout:                        30 * time.Second,
		RegisterNodeAtStart:            true,
		ExpirationPollIntervalSec:      10,
		DbPath:                         fmt.Sprintf("testdata/%s/db/operator_%d", oh.testName, operatorIndex),
		LogPath:                        logFilePath,
		ChurnerUrl:                     oh.chainHarness.Churner.URL,
		EnableTestMode:                 true,
		NumBatchValidators:             1,
		QuorumIDList:                   []core.QuorumID{0, 1},
		EigenDADirectory:               oh.testConfig.EigenDA.EigenDADirectory,
		StoreChunksRequestMaxPastAge:   5 * time.Minute,
		StoreChunksRequestMaxFutureAge: 5 * time.Minute,
		EthClientConfig: geth.EthClientConfig{
			RPCURLs:          []string{oh.chainHarness.GetAnvilRPCUrl()},
			PrivateKeyString: strings.TrimPrefix(operatorKey.PrivateKey, "0x"),
		},
		LoggerConfig: common.LoggerConfig{
			Format:       common.TextLogFormat,
			OutputWriter: io.MultiWriter(os.Stdout, logFile),
			HandlerOpts: logging.SLoggerOptions{
				Level:     slog.LevelDebug,
				NoColor:   true,
				AddSource: true,
			},
		},
		BlsSignerConfig: blssignerTypes.SignerConfig{
			SignerType: blssignerTypes.PrivateKey,
			PrivateKey: strings.TrimPrefix(blsKey.PrivateKey, "0x"),
		},
		EncoderConfig: kzg.KzgConfig{
			G1Path:          oh.srsG1Path,
			G2Path:          oh.srsG2Path,
			CacheDir:        fmt.Sprintf("testdata/%s/cache/operator_%d", oh.testName, operatorIndex),
			SRSOrder:        10000,
			SRSNumberToLoad: 10000,
			NumWorker:       4,
		},
		OnchainStateRefreshInterval:         10 * time.Second,
		OperatorStateCacheSize:              64,
		ChunkDownloadTimeout:                10 * time.Second,
		DownloadPoolSize:                    10,
		DispersalAuthenticationKeyCacheSize: 100,
		DisperserKeyTimeout:                 10 * time.Minute,
		RelayMaxMessageSize:                 units.GiB,
		EjectionSentinelPeriod:              5 * time.Minute,
		StoreChunksBufferTimeout:            10 * time.Second,
		StoreChunksBufferSizeBytes:          2 * units.GiB,
		GetChunksHotCacheReadLimitMB:        10 * units.GiB / units.MiB,
		GetChunksHotBurstLimitMB:            10 * units.GiB / units.MiB,
		GetChunksColdCacheReadLimitMB:       1 * units.GiB / units.MiB,
		GetChunksColdBurstLimitMB:           1 * units.GiB / units.MiB,
		GRPCMsgSizeLimitV2:                  1024 * 1024 * 300,
		ReservationLedgerCacheConfig:        reservationLedgerCacheConfig,
		EnablePerAccountPaymentMetrics:      false,
	}

	// Create operator logger
	operatorLogger, err := common.NewLogger(&nodeConfig.LoggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create operator logger: %w", err)
	}

	// Create metrics registry
	reg := prometheus.NewRegistry()

	// Create rate limiter
	globalParams := common.GlobalRateParams{
		BucketSizes: []time.Duration{450 * time.Second},
		Multipliers: []float32{2},
		CountFailed: true,
	}
	bucketStore, err := store.NewLocalParamStore[common.RateBucketParams](10000)
	if err != nil {
		return nil, fmt.Errorf("failed to create bucket store: %w", err)
	}
	ratelimiter := ratelimit.NewRateLimiter(reg, globalParams, bucketStore, operatorLogger)

	// Create RPC calls collector
	rpcCallsCollector := rpccalls.NewCollector(node.AppName, reg)

	// Create geth client
	gethClient, err := geth.NewInstrumentedEthClient(nodeConfig.EthClientConfig, rpcCallsCollector, operatorLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create geth client: %w", err)
	}

	// Create contract directory
	contractDirectory, err := directory.NewContractDirectory(
		ctx,
		operatorLogger,
		gethClient,
		gethcommon.HexToAddress(nodeConfig.EigenDADirectory))
	if err != nil {
		return nil, fmt.Errorf("failed to create contract directory: %w", err)
	}

	// Create version info
	softwareVersion := &version.Semver{}

	// Create mock IP provider for testing (returns "localhost")
	pubIPProvider := pubip.ProviderOrDefault(operatorLogger, "mockip")

	// Create node instance
	operatorNode, err := node.NewNode(
		ctx,
		reg,
		nodeConfig,
		contractDirectory,
		pubIPProvider,
		gethClient,
		operatorLogger,
		softwareVersion,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create operator node: %w", err)
	}

	// Create v2 server
	// Get operator state retriever and service manager addresses
	operatorStateRetrieverAddress, err := contractDirectory.GetContractAddress(
		ctx, directory.OperatorStateRetriever)
	if err != nil {
		return nil, fmt.Errorf("failed to get OperatorStateRetriever address: %w", err)
	}

	eigenDAServiceManagerAddress, err := contractDirectory.GetContractAddress(
		ctx, directory.ServiceManager)
	if err != nil {
		return nil, fmt.Errorf("failed to get ServiceManager address: %w", err)
	}

	// Create eth reader for v2 server
	reader, err := coreeth.NewReader(
		operatorLogger,
		gethClient,
		operatorStateRetrieverAddress.Hex(),
		eigenDAServiceManagerAddress.Hex())
	if err != nil {
		return nil, fmt.Errorf("cannot create eth.Reader: %w", err)
	}

	// Create v2 server
	serverV2, err := grpc.NewServerV2(
		ctx,
		nodeConfig,
		operatorNode,
		operatorLogger,
		ratelimiter,
		reg,
		reader,
		softwareVersion,
		listeners.v2.Dispersal,
		listeners.v2.Retrieval)
	if err != nil {
		return nil, fmt.Errorf("failed to create server v2: %w", err)
	}

	// Start all gRPC servers using the RunServers function
	err = grpc.RunServers(serverV2, nodeConfig, operatorLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to start gRPC servers: %w", err)
	}

	// Wait for servers to be ready
	time.Sleep(100 * time.Millisecond)
	logger.Info("Operator servers started successfully",
		"v2DispersalPort", listeners.v2.Dispersal.Addr().(*net.TCPAddr).Port,
		"v2RetrievalPort", listeners.v2.Retrieval.Addr().(*net.TCPAddr).Port,
		"operatorIndex", operatorIndex,
		"logFile", logFilePath)

	return serverV2, nil
}

// stopAllOperators stops all running operator servers
func (oh *OperatorHarness) stopAllOperators(logger logging.Logger) {
	// Stop V2 servers
	for i, serverV2 := range oh.ServersV2 {
		if serverV2 != nil {
			logger.Info("Stopping operator v2", "index", i)
			serverV2.Stop()
		}
	}

	// Clear the slice
	oh.ServersV2 = nil
}

// Cleanup is a public method for external cleanup.
func (oh *OperatorHarness) Cleanup(logger logging.Logger) {
	oh.stopAllOperators(logger)
}
