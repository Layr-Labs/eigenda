package integration_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
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
	grpclib "google.golang.org/grpc"
)

// OperatorInstance holds the state for a single operator
type OperatorInstance struct {
	Node                *node.Node
	Server              *grpc.Server
	ServerV2            *grpc.ServerV2
	DispersalServer     *grpclib.Server
	RetrievalServer     *grpclib.Server
	V2DispersalServer   *grpclib.Server
	V2RetrievalServer   *grpclib.Server
	DispersalListener   net.Listener
	RetrievalListener   net.Listener
	V2DispersalListener net.Listener
	V2RetrievalListener net.Listener
	DispersalPort       string
	RetrievalPort       string
	V2DispersalPort     string
	V2RetrievalPort     string
	Logger              logging.Logger
}

// StartOperatorForInfrastructure starts an operator node server as part of the global infrastructure.
// This should be called after Anvil and other containers are started.
func StartOperatorForInfrastructure(infra *InfrastructureHarness, operatorIndex int) (*OperatorInstance, error) {
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

	operatorConfig := &node.Config{
		Hostname:                       "localhost",
		RetrievalPort:                  retrievalPort,
		DispersalPort:                  dispersalPort,
		V2RetrievalPort:                v2RetrievalPort,
		V2DispersalPort:                v2DispersalPort,
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
		ChurnerUrl:                     "localhost:32002",
		EnableTestMode:                 true,
		NumBatchValidators:             1,
		QuorumIDList:                   []core.QuorumID{0, 1}, // Default to quorums 0 and 1
		EigenDADirectory:               infra.TestConfig.EigenDA.EigenDADirectory,
		DisableDispersalAuthentication: true, // TODO: enable
		EthClientConfig: geth.EthClientConfig{
			RPCURLs:          []string{"http://localhost:8545"},
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
		StoreChunksBufferSizeBytes:          2 * units.GiB, // 2GB buffer for storing chunks
		GetChunksHotCacheReadLimitMB:        10 * units.GiB / units.MiB, // 10 GB/s for tests
		GetChunksHotBurstLimitMB:            10 * units.GiB / units.MiB, // 10 GB burst
		GetChunksColdCacheReadLimitMB:       1 * units.GiB / units.MiB,  // 1 GB/s for tests
		GetChunksColdBurstLimitMB:           1 * units.GiB / units.MiB,  // 1 GB burst
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
		context.Background(),
		operatorLogger,
		gethClient,
		gethcommon.HexToAddress(operatorConfig.EigenDADirectory))
	if err != nil {
		return nil, fmt.Errorf("failed to create contract directory: %w", err)
	}

	// Create public IP provider
	pubIPProvider := &mockPublicIPProvider{ip: "127.0.0.1"}

	// Create version info
	softwareVersion := &version.Semver{}

	// Create node instance
	operatorNode, err := node.NewNode(
		context.Background(),
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
			context.Background(), directory.OperatorStateRetriever)
		if err != nil {
			return nil, fmt.Errorf("failed to get OperatorStateRetriever address: %w", err)
		}

		eigenDAServiceManagerAddress, err := contractDirectory.GetContractAddress(
			context.Background(), directory.ServiceManager)
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
			context.Background(),
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

	// Create and start dispersal server
	var dispersalServer *grpclib.Server
	var dispersalListener net.Listener
	if operatorConfig.EnableV1 {
		dispersalListener, err = net.Listen("tcp", fmt.Sprintf(":%s", dispersalPort))
		if err != nil {
			return nil, fmt.Errorf("failed to listen on dispersal port %s: %w", dispersalPort, err)
		}

		dispersalServer = grpclib.NewServer(grpclib.MaxRecvMsgSize(60 * 1024 * 1024 * 1024)) // 60 GiB
		pb.RegisterDispersalServer(dispersalServer, operatorServer)
		healthcheck.RegisterHealthServer("node.Dispersal", dispersalServer)

		go func() {
			operatorLogger.Info("Starting dispersal server", "port", dispersalPort)
			if err := dispersalServer.Serve(dispersalListener); err != nil {
				operatorLogger.Info("Dispersal server stopped", "error", err)
			}
		}()
	}

	// Create and start retrieval server
	var retrievalServer *grpclib.Server
	var retrievalListener net.Listener
	if operatorConfig.EnableV1 {
		retrievalListener, err = net.Listen("tcp", fmt.Sprintf(":%s", retrievalPort))
		if err != nil {
			return nil, fmt.Errorf("failed to listen on retrieval port %s: %w", retrievalPort, err)
		}

		retrievalServer = grpclib.NewServer(grpclib.MaxRecvMsgSize(1024 * 1024 * 300)) // 300 MiB
		pb.RegisterRetrievalServer(retrievalServer, operatorServer)
		healthcheck.RegisterHealthServer("node.Retrieval", retrievalServer)

		go func() {
			operatorLogger.Info("Starting retrieval server", "port", retrievalPort)
			if err := retrievalServer.Serve(retrievalListener); err != nil {
				operatorLogger.Info("Retrieval server stopped", "error", err)
			}
		}()
	}

	// Create and start V2 dispersal server
	var v2DispersalServer *grpclib.Server
	var v2DispersalListener net.Listener
	if operatorConfig.EnableV2 {
		v2DispersalListener, err = net.Listen("tcp", fmt.Sprintf(":%s", v2DispersalPort))
		if err != nil {
			return nil, fmt.Errorf("failed to listen on v2 dispersal port %s: %w", v2DispersalPort, err)
		}

		v2DispersalServer = grpclib.NewServer(grpclib.MaxRecvMsgSize(1024 * 1024 * 300)) // 300 MiB
		validator.RegisterDispersalServer(v2DispersalServer, serverV2)
		healthcheck.RegisterHealthServer("node.v2.Dispersal", v2DispersalServer)

		go func() {
			operatorLogger.Info("Starting v2 dispersal server", "port", v2DispersalPort)
			if err := v2DispersalServer.Serve(v2DispersalListener); err != nil {
				operatorLogger.Info("V2 Dispersal server stopped", "error", err)
			}
		}()
	}

	// Create and start V2 retrieval server
	var v2RetrievalServer *grpclib.Server
	var v2RetrievalListener net.Listener
	if operatorConfig.EnableV2 {
		v2RetrievalListener, err = net.Listen("tcp", fmt.Sprintf(":%s", v2RetrievalPort))
		if err != nil {
			return nil, fmt.Errorf("failed to listen on v2 retrieval port %s: %w", v2RetrievalPort, err)
		}

		v2RetrievalServer = grpclib.NewServer(grpclib.MaxRecvMsgSize(1024 * 1024 * 300)) // 300 MiB
		validator.RegisterRetrievalServer(v2RetrievalServer, serverV2)
		healthcheck.RegisterHealthServer("node.v2.Retrieval", v2RetrievalServer)

		go func() {
			operatorLogger.Info("Starting v2 retrieval server", "port", v2RetrievalPort)
			if err := v2RetrievalServer.Serve(v2RetrievalListener); err != nil {
				operatorLogger.Info("V2 Retrieval server stopped", "error", err)
			}
		}()
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
		Node:                operatorNode,
		Server:              operatorServer,
		ServerV2:            serverV2,
		DispersalServer:     dispersalServer,
		RetrievalServer:     retrievalServer,
		V2DispersalServer:   v2DispersalServer,
		V2RetrievalServer:   v2RetrievalServer,
		DispersalListener:   dispersalListener,
		RetrievalListener:   retrievalListener,
		V2DispersalListener: v2DispersalListener,
		V2RetrievalListener: v2RetrievalListener,
		DispersalPort:       dispersalPort,
		RetrievalPort:       retrievalPort,
		V2DispersalPort:     v2DispersalPort,
		V2RetrievalPort:     v2RetrievalPort,
		Logger:              operatorLogger,
	}, nil
}

// StartOperatorsForInfrastructure starts all operator nodes configured in the test config
func StartOperatorsForInfrastructure(infra *InfrastructureHarness) error {
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
		instance, err := StartOperatorForInfrastructure(infra, i)
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

	// Stop the dispersal server
	if instance.DispersalServer != nil {
		instance.Logger.Info("Stopping dispersal server")
		instance.DispersalServer.GracefulStop()
	}

	// Stop the retrieval server
	if instance.RetrievalServer != nil {
		instance.Logger.Info("Stopping retrieval server")
		instance.RetrievalServer.GracefulStop()
	}

	// Stop the v2 dispersal server
	if instance.V2DispersalServer != nil {
		instance.Logger.Info("Stopping v2 dispersal server")
		instance.V2DispersalServer.GracefulStop()
	}

	// Stop the v2 retrieval server
	if instance.V2RetrievalServer != nil {
		instance.Logger.Info("Stopping v2 retrieval server")
		instance.V2RetrievalServer.GracefulStop()
	}

	// Close the dispersal listener
	if instance.DispersalListener != nil {
		_ = instance.DispersalListener.Close()
	}

	// Close the retrieval listener
	if instance.RetrievalListener != nil {
		_ = instance.RetrievalListener.Close()
	}

	// Close the v2 dispersal listener
	if instance.V2DispersalListener != nil {
		_ = instance.V2DispersalListener.Close()
	}

	// Close the v2 retrieval listener
	if instance.V2RetrievalListener != nil {
		_ = instance.V2RetrievalListener.Close()
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

// mockPublicIPProvider is a simple mock implementation for testing
type mockPublicIPProvider struct {
	ip string
}

func (m *mockPublicIPProvider) PublicIPAddress(_ context.Context) (string, error) {
	return m.ip, nil
}

func (m *mockPublicIPProvider) Name() string {
	return "mock"
}

// Ensure mockPublicIPProvider implements pubip.Provider interface
var _ interface {
	PublicIPAddress(context.Context) (string, error)
	Name() string
} = (*mockPublicIPProvider)(nil)
