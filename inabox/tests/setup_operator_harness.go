package integration

import (
	"context"
	"fmt"
	"io"
	"log/slog"
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

// OperatorInstance holds the state for a single operator
type OperatorInstance struct {
	Node            *node.Node
	Server          *grpc.Server
	ServerV2        *grpc.ServerV2
	DispersalPort   string
	RetrievalPort   string
	V2DispersalPort string
	V2RetrievalPort string
	Logger          logging.Logger
}

// OperatorHarnessConfig contains the configuration for setting up the operator harness
type OperatorHarnessConfig struct {
	TestConfig   *deploy.Config
	TestName     string
	Logger       logging.Logger
	ChainHarness *ChainHarness // Access to chain infrastructure
	Ctx          context.Context
}

// SetupOperatorHarness creates and initializes the operator harness
func SetupOperatorHarness(ctx context.Context, config *OperatorHarnessConfig) (*OperatorHarness, error) {
	harness := &OperatorHarness{
		OperatorInstances: make([]*OperatorInstance, 0),
	}

	// Store references we'll need
	harness.testConfig = config.TestConfig
	harness.testName = config.TestName
	harness.logger = config.Logger
	harness.chainHarness = config.ChainHarness
	harness.ctx = config.Ctx

	// Start all operators
	if err := harness.StartOperators(); err != nil {
		return nil, err
	}

	return harness, nil
}

// Add fields to OperatorHarness to support the methods
type OperatorHarness struct {
	OperatorInstances []*OperatorInstance

	// Internal fields for operator management
	testConfig   *deploy.Config
	testName     string
	logger       logging.Logger
	chainHarness *ChainHarness
	ctx          context.Context
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

// StartOperators starts all operator nodes configured in the test config
func (oh *OperatorHarness) StartOperators() error {
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

	oh.logger.Info("Starting operator goroutines", "count", operatorCount)

	// Start each operator
	for i := 0; i < operatorCount; i++ {
		instance, err := oh.startOperator(i)
		if err != nil {
			// Clean up any operators we started before failing
			oh.Cleanup(context.Background(), oh.logger)
			return fmt.Errorf("failed to start operator %d: %w", i, err)
		}
		oh.OperatorInstances = append(oh.OperatorInstances, instance)
		oh.logger.Info("Started operator", "index", i,
			"dispersalPort", instance.DispersalPort, "retrievalPort", instance.RetrievalPort)
	}

	return nil
}

// startOperator starts a single operator with the given index
func (oh *OperatorHarness) startOperator(operatorIndex int) (*OperatorInstance, error) {
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
	logsDir := fmt.Sprintf("testdata/%s/logs", oh.testName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/operator_%d.log", logsDir, operatorIndex)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open operator log file: %w", err)
	}

	// Create operator configuration
	retrievalPort := fmt.Sprintf("3410%d", operatorIndex)
	dispersalPort := fmt.Sprintf("3310%d", operatorIndex)
	v2RetrievalPort := fmt.Sprintf("3510%d", operatorIndex)
	v2DispersalPort := fmt.Sprintf("3610%d", operatorIndex)
	nodeApiPort := fmt.Sprintf("3710%d", operatorIndex)
	metricsPort := 3800 + operatorIndex

	// TODO(dmanc): The node config is quite a beast. This is a configuration that
	// passed the tests after a bunch of trial and error.
	// We really need better validation on the node constructor.
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
				Level:   slog.LevelDebug,
				NoColor: true,
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
		oh.ctx,
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
		oh.ctx,
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

	// Create operator gRPC server
	operatorServer := grpc.NewServer(
		nodeConfig,
		operatorNode,
		operatorLogger,
		ratelimiter,
		softwareVersion,
	)

	// Create v2 server if enabled
	var serverV2 *grpc.ServerV2
	if nodeConfig.EnableV2 {
		// Get operator state retriever and service manager addresses
		operatorStateRetrieverAddress, err := contractDirectory.GetContractAddress(
			oh.ctx, directory.OperatorStateRetriever)
		if err != nil {
			return nil, fmt.Errorf("failed to get OperatorStateRetriever address: %w", err)
		}

		eigenDAServiceManagerAddress, err := contractDirectory.GetContractAddress(
			oh.ctx, directory.ServiceManager)
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
		serverV2, err = grpc.NewServerV2(
			oh.ctx,
			nodeConfig,
			operatorNode,
			operatorLogger,
			ratelimiter,
			reg,
			reader,
			softwareVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to create server v2: %w", err)
		}
	}

	// Start all gRPC servers using the new RunServers function
	err = grpc.RunServers(operatorServer, serverV2, nodeConfig, operatorLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to start gRPC servers: %w", err)
	}

	// Wait for servers to be ready
	time.Sleep(100 * time.Millisecond)
	operatorLogger.Info("Operator servers started successfully",
		"dispersalPort", dispersalPort,
		"retrievalPort", retrievalPort,
		"v2DispersalPort", v2DispersalPort,
		"v2RetrievalPort", v2RetrievalPort,
		"operatorIndex", operatorIndex,
		"logFile", logFilePath)

	return &OperatorInstance{
		Node:            operatorNode,
		Server:          operatorServer,
		ServerV2:        serverV2,
		DispersalPort:   dispersalPort,
		RetrievalPort:   retrievalPort,
		V2DispersalPort: v2DispersalPort,
		V2RetrievalPort: v2RetrievalPort,
		Logger:          operatorLogger,
	}, nil
}

// Cleanup releases resources held by the OperatorHarness
func (oh *OperatorHarness) Cleanup(ctx context.Context, logger logging.Logger) {
	if len(oh.OperatorInstances) == 0 {
		return
	}
	logger.Info("Stopping all operator goroutines")
	for i, instance := range oh.OperatorInstances {
		if instance == nil {
			continue
		}
		logger.Info("Stopping operator", "index", i)
		StopOperator(instance)
	}
	oh.OperatorInstances = nil
}

// StopOperator gracefully stops an operator instance
func StopOperator(instance *OperatorInstance) {
	if instance == nil {
		return
	}

	instance.Logger.Info("Stopping operator")

	// TODO: Add graceful shutdown of node once it's implemented

	instance.Logger.Info("Operator stopped")
}

// StartOperatorForInfrastructure is a compatibility wrapper for existing code.
// It creates a temporary OperatorHarness and starts a single operator.
// New code should use OperatorHarness.startOperator directly.
func StartOperatorForInfrastructure(
	infra *InfrastructureHarness, operatorIndex int,
) (*OperatorInstance, error) {
	// Create a temporary harness with the infrastructure references
	harness := &OperatorHarness{
		testConfig:   infra.TestConfig,
		testName:     infra.TestName,
		logger:       infra.Logger,
		chainHarness: &infra.ChainHarness,
		ctx:          infra.Ctx,
	}

	return harness.startOperator(operatorIndex)
}

// StartOperatorsForInfrastructure is a compatibility wrapper for existing code.
// It updates the infrastructure's OperatorHarness with all started operators.
// New code should use SetupOperatorHarness and OperatorHarness.StartOperators directly.
func StartOperatorsForInfrastructure(infra *InfrastructureHarness) error {
	// Set up the operator harness with references from infrastructure
	infra.OperatorHarness = OperatorHarness{
		testConfig:   infra.TestConfig,
		testName:     infra.TestName,
		logger:       infra.Logger,
		chainHarness: &infra.ChainHarness,
		ctx:          infra.Ctx,
	}

	// Start all operators
	return infra.OperatorHarness.StartOperators()
}
