package integration

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"runtime"
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
	"github.com/Layr-Labs/eigenda/encoding/kzg"
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
	Logger     logging.Logger
}

// SetupOperatorHarness creates and initializes the operator harness
func SetupOperatorHarness(
	ctx context.Context,
	config *OperatorHarnessConfig,
	chainHarness *ChainHarness,
) (*OperatorHarness, error) {
	harness := &OperatorHarness{
		Servers:   make([]*grpc.Server, 0),
		ServersV2: make([]*grpc.ServerV2, 0),
	}

	// Store references we'll need
	harness.testConfig = config.TestConfig
	harness.testName = config.TestName
	harness.logger = config.Logger
	harness.chainHarness = chainHarness

	// Start all operators
	if err := harness.StartOperators(ctx); err != nil {
		return nil, err
	}

	return harness, nil
}

// OperatorHarness manages operator instances for integration tests
type OperatorHarness struct {
	Servers   []*grpc.Server
	ServersV2 []*grpc.ServerV2

	// Internal fields for operator management
	testConfig   *deploy.Config
	testName     string
	logger       logging.Logger
	chainHarness *ChainHarness
	srsG1Path    string
	srsG2Path    string
}

// getSRSPaths returns the correct paths to SRS files based on the source file location.
// This uses runtime.Caller to determine where this file is located and calculates
// the relative path to the resources/srs directory from there.
func getSRSPaths() (g1Path, g2Path string, err error) {
	// Get the path of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", "", fmt.Errorf("failed to get caller information")
	}

	// We need to go up 2 directories from tests/ to get to inabox/, then up one more to get to the project root
	// From project root, resources/srs is the target
	testDir := filepath.Dir(filename)
	inaboxDir := filepath.Dir(testDir)
	projectRoot := filepath.Dir(inaboxDir)

	g1Path = filepath.Join(projectRoot, "resources", "srs", "g1.point")
	g2Path = filepath.Join(projectRoot, "resources", "srs", "g2.point")

	return g1Path, g2Path, nil
}

// operatorListeners holds the four network listeners for a single operator
type operatorListeners struct {
	v1Dispersal  net.Listener
	v1Retrieval  net.Listener
	v2Dispersal  net.Listener
	v2Retrieval  net.Listener
}

// StartOperators starts all operator nodes configured in the test config
func (oh *OperatorHarness) StartOperators(ctx context.Context) error {
	// Get SRS paths first - fail early if we can't find them
	g1Path, g2Path, err := getSRSPaths()
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

	oh.logger.Info("Pre-creating listeners for operator goroutines", "count", operatorCount)

	// Pre-create all listeners with port 0 (OS assigns ports)
	allListeners := make([]operatorListeners, operatorCount)
	for i := range operatorCount {
		listeners, err := oh.createOperatorListeners(i)
		if err != nil {
			// Clean up any listeners we created before failing
			oh.closeAllListeners(allListeners[:i])
			return err
		}
		allListeners[i] = listeners

		oh.logger.Info("Created listeners for operator", "index", i,
			"v1DispersalPort", listeners.v1Dispersal.Addr().(*net.TCPAddr).Port,
			"v1RetrievalPort", listeners.v1Retrieval.Addr().(*net.TCPAddr).Port,
			"v2DispersalPort", listeners.v2Dispersal.Addr().(*net.TCPAddr).Port,
			"v2RetrievalPort", listeners.v2Retrieval.Addr().(*net.TCPAddr).Port)
	}

	// Now start each operator with its pre-created listeners
	for i, listeners := range allListeners {
		server, serverV2, err := oh.startOperator(ctx, i, listeners)
		if err != nil {
			// Clean up any operators we started and all remaining listeners
			oh.Cleanup(context.Background(), oh.logger)
			oh.closeAllListeners(allListeners[i:])
			return fmt.Errorf("failed to start operator %d: %w", i, err)
		}
		oh.Servers = append(oh.Servers, server)
		oh.ServersV2 = append(oh.ServersV2, serverV2)
		oh.logger.Info("Started operator", "index", i,
			"v1DispersalPort", server.GetDispersalPort(),
			"v1RetrievalPort", server.GetRetrievalPort(),
			"v2DispersalPort", serverV2.GetDispersalPort(),
			"v2RetrievalPort", serverV2.GetRetrievalPort())
	}

	return nil
}

// createOperatorListeners creates all four listeners for an operator
func (oh *OperatorHarness) createOperatorListeners(operatorIndex int) (operatorListeners, error) {
	var listeners operatorListeners

	v1DispersalListener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return listeners, fmt.Errorf("failed to create v1 dispersal listener for operator %d: %w", operatorIndex, err)
	}

	v1RetrievalListener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		_ = v1DispersalListener.Close()
		return listeners, fmt.Errorf("failed to create v1 retrieval listener for operator %d: %w", operatorIndex, err)
	}

	v2DispersalListener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		_ = v1DispersalListener.Close()
		_ = v1RetrievalListener.Close()
		return listeners, fmt.Errorf("failed to create v2 dispersal listener for operator %d: %w", operatorIndex, err)
	}

	v2RetrievalListener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		_ = v1DispersalListener.Close()
		_ = v1RetrievalListener.Close()
		_ = v2DispersalListener.Close()
		return listeners, fmt.Errorf("failed to create v2 retrieval listener for operator %d: %w", operatorIndex, err)
	}

	listeners.v1Dispersal = v1DispersalListener
	listeners.v1Retrieval = v1RetrievalListener
	listeners.v2Dispersal = v2DispersalListener
	listeners.v2Retrieval = v2RetrievalListener

	return listeners, nil
}

// closeAllListeners closes all listeners in the slice
func (oh *OperatorHarness) closeAllListeners(allListeners []operatorListeners) {
	for i, listeners := range allListeners {
		if listeners.v1Dispersal != nil {
			if err := listeners.v1Dispersal.Close(); err != nil {
				oh.logger.Warn("Failed to close v1 dispersal listener", "operatorIndex", i, "error", err)
			}
		}
		if listeners.v1Retrieval != nil {
			if err := listeners.v1Retrieval.Close(); err != nil {
				oh.logger.Warn("Failed to close v1 retrieval listener", "operatorIndex", i, "error", err)
			}
		}
		if listeners.v2Dispersal != nil {
			if err := listeners.v2Dispersal.Close(); err != nil {
				oh.logger.Warn("Failed to close v2 dispersal listener", "operatorIndex", i, "error", err)
			}
		}
		if listeners.v2Retrieval != nil {
			if err := listeners.v2Retrieval.Close(); err != nil {
				oh.logger.Warn("Failed to close v2 retrieval listener", "operatorIndex", i, "error", err)
			}
		}
	}
}

// startOperator starts a single operator with the given index and pre-created listeners
func (oh *OperatorHarness) startOperator(ctx context.Context, operatorIndex int, listeners operatorListeners) (*grpc.Server, *grpc.ServerV2, error) {
	// Get operator's private key
	operatorName := fmt.Sprintf("opr%d", operatorIndex)

	// Check if operator exists in test config
	if oh.testConfig.Pks == nil || oh.testConfig.Pks.EcdsaMap == nil {
		return nil, nil, fmt.Errorf("no private keys configured")
	}

	operatorKey, ok := oh.testConfig.Pks.EcdsaMap[operatorName]
	if !ok {
		return nil, nil, fmt.Errorf("operator %s not found in config", operatorName)
	}

	// Get BLS key configuration
	blsKey, blsOk := oh.testConfig.Pks.BlsMap[operatorName]
	if !blsOk {
		return nil, nil, fmt.Errorf("BLS key for %s not found in config", operatorName)
	}

	// Create logs directory
	// TODO(dmanc): If possible we should have a centralized place for creating loggers and injecting them into the config.
	logsDir := fmt.Sprintf("testdata/%s/logs", oh.testName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/operator_%d.log", logsDir, operatorIndex)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open operator log file: %w", err)
	}

	// Extract actual ports assigned by OS from the pre-created listeners
	dispersalPort := fmt.Sprintf("%d", listeners.v1Dispersal.Addr().(*net.TCPAddr).Port)
	retrievalPort := fmt.Sprintf("%d", listeners.v1Retrieval.Addr().(*net.TCPAddr).Port)
	v2DispersalPort := fmt.Sprintf("%d", listeners.v2Dispersal.Addr().(*net.TCPAddr).Port)
	v2RetrievalPort := fmt.Sprintf("%d", listeners.v2Retrieval.Addr().(*net.TCPAddr).Port)
	nodeApiPort := fmt.Sprintf("3710%d", operatorIndex)
	metricsPort := 3800 + operatorIndex

	// TODO(dmanc): The node config is quite a beast. This is a configuration that
	// passed the tests after a bunch of trial and error.
	// We really need better validation on the node constructor.

	// TODO(dmanc): In addition to loggers, we should have a centralized place for creating
	// configuration and injecting it into the harness config.
	nodeConfig := &node.Config{
		Hostname:                       "localhost",
		RetrievalPort:                  retrievalPort,
		DispersalPort:                  dispersalPort,
		V2RetrievalPort:                v2RetrievalPort,
		V2DispersalPort:                v2DispersalPort,
		InternalRetrievalPort:          retrievalPort,
		InternalDispersalPort:          dispersalPort,
		InternalV2RetrievalPort:        v2RetrievalPort,
		InternalV2DispersalPort:        v2DispersalPort,
		EnableNodeApi:                  true,
		NodeApiPort:                    nodeApiPort,
		EnableMetrics:                  true,
		MetricsPort:                    metricsPort,
		Timeout:                        30 * time.Second,
		RegisterNodeAtStart:            true,
		ExpirationPollIntervalSec:      10,
		EnableV1:                       true,
		EnableV2:                       true,
		DbPath:                         fmt.Sprintf("testdata/%s/db/operator_%d", oh.testName, operatorIndex),
		LogPath:                        logFilePath,
		ChurnerUrl:                     oh.chainHarness.Churner.URL,
		EnableTestMode:                 true,
		NumBatchValidators:             1,
		QuorumIDList:                   []core.QuorumID{0, 1},
		EigenDADirectory:               oh.testConfig.EigenDA.EigenDADirectory,
		DisableDispersalAuthentication: true, // TODO: set to false
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
	}

	// Create operator logger
	operatorLogger, err := common.NewLogger(&nodeConfig.LoggerConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create operator logger: %w", err)
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
		return nil, nil, fmt.Errorf("failed to create bucket store: %w", err)
	}
	ratelimiter := ratelimit.NewRateLimiter(reg, globalParams, bucketStore, operatorLogger)

	// Create RPC calls collector
	rpcCallsCollector := rpccalls.NewCollector(node.AppName, reg)

	// Create geth client
	gethClient, err := geth.NewInstrumentedEthClient(nodeConfig.EthClientConfig, rpcCallsCollector, operatorLogger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create geth client: %w", err)
	}

	// Create contract directory
	contractDirectory, err := directory.NewContractDirectory(
		ctx,
		operatorLogger,
		gethClient,
		gethcommon.HexToAddress(nodeConfig.EigenDADirectory))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create contract directory: %w", err)
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
		return nil, nil, fmt.Errorf("failed to create operator node: %w", err)
	}

	// Create operator gRPC server
	operatorServer := grpc.NewServer(
		nodeConfig,
		operatorNode,
		operatorLogger,
		ratelimiter,
		softwareVersion,
		listeners.v1Dispersal,
		listeners.v1Retrieval,
	)

	// Create v2 server if enabled
	var serverV2 *grpc.ServerV2
	if nodeConfig.EnableV2 {
		// Get operator state retriever and service manager addresses
		operatorStateRetrieverAddress, err := contractDirectory.GetContractAddress(
			ctx, directory.OperatorStateRetriever)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get OperatorStateRetriever address: %w", err)
		}

		eigenDAServiceManagerAddress, err := contractDirectory.GetContractAddress(
			ctx, directory.ServiceManager)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get ServiceManager address: %w", err)
		}

		// Create eth reader for v2 server
		reader, err := coreeth.NewReader(
			operatorLogger,
			gethClient,
			operatorStateRetrieverAddress.Hex(),
			eigenDAServiceManagerAddress.Hex())
		if err != nil {
			return nil, nil, fmt.Errorf("cannot create eth.Reader: %w", err)
		}

		// Create v2 server
		serverV2, err = grpc.NewServerV2(
			ctx,
			nodeConfig,
			operatorNode,
			operatorLogger,
			ratelimiter,
			reg,
			reader,
			softwareVersion,
			listeners.v2Dispersal,
			listeners.v2Retrieval)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create server v2: %w", err)
		}
	}

	// Start all gRPC servers using the new RunServers function
	err = grpc.RunServers(operatorServer, serverV2, nodeConfig, operatorLogger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start gRPC servers: %w", err)
	}

	// Wait for servers to be ready
	time.Sleep(100 * time.Millisecond)
	operatorLogger.Info("Operator servers started successfully",
		"v1DispersalPort", listeners.v1Dispersal.Addr().(*net.TCPAddr).Port,
		"v1RetrievalPort", listeners.v1Retrieval.Addr().(*net.TCPAddr).Port,
		"v2DispersalPort", listeners.v2Dispersal.Addr().(*net.TCPAddr).Port,
		"v2RetrievalPort", listeners.v2Retrieval.Addr().(*net.TCPAddr).Port,
		"operatorIndex", operatorIndex,
		"logFile", logFilePath)

	return operatorServer, serverV2, nil
}

// Cleanup releases resources held by the OperatorHarness
func (oh *OperatorHarness) Cleanup(ctx context.Context, logger logging.Logger) {
	if len(oh.Servers) == 0 {
		return
	}
	logger.Info("Stopping all operator servers", "count", len(oh.Servers))
	for i, server := range oh.Servers {
		if server == nil {
			continue
		}
		logger.Info("Stopping operator", "index", i)
		server.Stop()
	}
	oh.Servers = nil
	oh.ServersV2 = nil
}
