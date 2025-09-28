package integration_test

import (
	"fmt"
	"io"
	"log/slog"
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
	"github.com/Layr-Labs/eigenda/encoding/kzg"
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
	ServerRunner    *grpc.ServerRunner
	DispersalPort   string
	RetrievalPort   string
	V2DispersalPort string
	V2RetrievalPort string
	Logger          logging.Logger
}

// StartOperatorForInfrastructure starts an operator node server as part of the global infrastructure.
// This should be called after Anvil and the Churner are started.
func StartOperatorForInfrastructure(infra *InfrastructureHarness, operatorIndex int, anvilRPC string, churnerRPC string) (*OperatorInstance, error) {
	// Get operator's private key
	var privateKey string
	operatorName := fmt.Sprintf("opr%d", operatorIndex)

	// Check if operator exists in test config
	if infra.TestConfig.Pks == nil || infra.TestConfig.Pks.EcdsaMap == nil {
		return nil, fmt.Errorf("no private keys configured")
	}

	operatorKey, ok := infra.TestConfig.Pks.EcdsaMap[operatorName]
	if !ok {
		return nil, fmt.Errorf("operator %s not found in config", operatorName)
	}
	privateKey = strings.TrimPrefix(operatorKey.PrivateKey, "0x")

	// Create logs directory
	logsDir := fmt.Sprintf("testdata/%s/logs", infra.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/operator_%d.log", logsDir, operatorIndex)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open operator log file: %w", err)
	}

	// Get BLS key configuration - use the same operator name as for ECDSA key
	blsKey, blsOk := infra.TestConfig.Pks.BlsMap[operatorName]
	if !blsOk {
		return nil, fmt.Errorf("BLS key for %s not found in config", operatorName)
	}

	// Create operator configuration
	retrievalPort := fmt.Sprintf("3410%d", operatorIndex)
	dispersalPort := fmt.Sprintf("3310%d", operatorIndex)
	v2RetrievalPort := fmt.Sprintf("3510%d", operatorIndex)
	v2DispersalPort := fmt.Sprintf("3610%d", operatorIndex)
	nodeApiPort := fmt.Sprintf("3710%d", operatorIndex)
	metricsPort := 3800 + operatorIndex

	// TODO(dmanc): The node config is quite a beast. This is a configuration that passed the tests after a bunch of trial and error.
	// We really need better validation on the node constructor.
	operatorConfig := &node.Config{
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
		DbPath:                         fmt.Sprintf("testdata/%s/db/operator_%d", infra.TestName, operatorIndex),
		LogPath:                        logFilePath,
		ChurnerUrl:                     churnerRPC,
		EnableTestMode:                 true,
		NumBatchValidators:             1,
		QuorumIDList:                   []core.QuorumID{0, 1},
		EigenDADirectory:               infra.TestConfig.EigenDA.EigenDADirectory,
		DisableDispersalAuthentication: true, // TODO: set to false
		EthClientConfig: geth.EthClientConfig{
			RPCURLs:          []string{anvilRPC},
			PrivateKeyString: privateKey,
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
			G1Path:          "../resources/srs/g1.point",
			G2Path:          "../resources/srs/g2.point",
			CacheDir:        fmt.Sprintf("testdata/%s/cache/operator_%d", infra.TestName, operatorIndex),
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
	operatorLogger, err := common.NewLogger(&operatorConfig.LoggerConfig)
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
	gethClient, err := geth.NewInstrumentedEthClient(operatorConfig.EthClientConfig, rpcCallsCollector, operatorLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create geth client: %w", err)
	}

	// Create contract directory
	contractDirectory, err := directory.NewContractDirectory(
		infra.Ctx,
		operatorLogger,
		gethClient,
		gethcommon.HexToAddress(operatorConfig.EigenDADirectory))
	if err != nil {
		return nil, fmt.Errorf("failed to create contract directory: %w", err)
	}

	// Create version info
	softwareVersion := &version.Semver{}

	// Create mock IP provider for testing (returns "localhost")
	pubIPProvider := pubip.ProviderOrDefault(operatorLogger, "mockip")

	// Create node instance
	operatorNode, err := node.NewNode(
		infra.Ctx,
		reg,
		operatorConfig,
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
		operatorConfig,
		operatorNode,
		operatorLogger,
		ratelimiter,
		softwareVersion,
	)

	// Create v2 server if enabled
	var serverV2 *grpc.ServerV2
	if operatorConfig.EnableV2 {
		// Get operator state retriever and service manager addresses
		operatorStateRetrieverAddress, err := contractDirectory.GetContractAddress(
			infra.Ctx, directory.OperatorStateRetriever)
		if err != nil {
			return nil, fmt.Errorf("failed to get OperatorStateRetriever address: %w", err)
		}

		eigenDAServiceManagerAddress, err := contractDirectory.GetContractAddress(
			infra.Ctx, directory.ServiceManager)
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
			infra.Ctx,
			operatorConfig,
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
	runner, err := grpc.RunServers(operatorServer, serverV2, operatorConfig, operatorLogger)
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
		ServerRunner:    runner,
		DispersalPort:   dispersalPort,
		RetrievalPort:   retrievalPort,
		V2DispersalPort: v2DispersalPort,
		V2RetrievalPort: v2RetrievalPort,
		Logger:          operatorLogger,
	}, nil
}

// StartOperatorsForInfrastructure starts all operator nodes configured in the test config
func StartOperatorsForInfrastructure(infra *InfrastructureHarness, anvilRPC string, churnerRPC string) error {
	// Count how many operator configs exist
	operatorCount := 0
	for {
		operatorName := fmt.Sprintf("opr%d", operatorCount)
		if _, ok := infra.TestConfig.Pks.EcdsaMap[operatorName]; !ok {
			break
		}
		operatorCount++
	}

	if operatorCount == 0 {
		return fmt.Errorf("no operators found in config")
	}

	infra.Logger.Info("Starting operator goroutines", "count", operatorCount)

	// Start each operator
	for i := 0; i < operatorCount; i++ {
		instance, err := StartOperatorForInfrastructure(infra, i, anvilRPC, churnerRPC)
		if err != nil {
			// Clean up any operators we started before failing
			for _, inst := range infra.OperatorInstances {
				StopOperator(inst)
			}
			return fmt.Errorf("failed to start operator %d: %w", i, err)
		}
		infra.OperatorInstances = append(infra.OperatorInstances, instance)
		infra.Logger.Info("Started operator", "index", i, "dispersalPort", instance.DispersalPort, "retrievalPort", instance.RetrievalPort)
	}

	return nil
}

// StopOperator gracefully stops an operator instance
func StopOperator(instance *OperatorInstance) {
	if instance == nil {
		return
	}

	instance.Logger.Info("Stopping operator")

	// Use the ServerRunner to gracefully stop all servers
	if instance.ServerRunner != nil {
		instance.ServerRunner.Stop()
	}

	instance.Logger.Info("Operator stopped")
}

// StopAllOperators stops all operator instances in the infrastructure
func StopAllOperators(infra *InfrastructureHarness) {
	if infra == nil || len(infra.OperatorInstances) == 0 {
		return
	}

	infra.Logger.Info("Stopping all operator goroutines")
	for i, instance := range infra.OperatorInstances {
		infra.Logger.Info("Stopping operator", "index", i)
		StopOperator(instance)
	}
	infra.OperatorInstances = nil
}
