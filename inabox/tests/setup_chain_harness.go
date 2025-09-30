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

	pb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/operators/churner"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/grpc"
)

// ChainHarnessConfig contains the configuration for setting up the chain harness
type ChainHarnessConfig struct {
	TestConfig *deploy.Config
	TestName   string
	Logger     logging.Logger
	Network    *testcontainers.DockerNetwork
}

type ChainHarness struct {
	Anvil     *testbed.AnvilContainer
	GraphNode *testbed.GraphNodeContainer // Optional, only when subgraphs are deployed
	Churner   struct {
		Server   *grpc.Server
		Listener net.Listener
		URL      string
	}
	EthClient *geth.MultiHomingClient
}

// SetupChainHarness creates and initializes the chain infrastructure (Anvil, Graph Node, contracts, and Churner)
func SetupChainHarness(ctx context.Context, config *ChainHarnessConfig) (*ChainHarness, error) {
	harness := &ChainHarness{}

	// Step 1: Setup Anvil
	config.Logger.Info("Starting anvil")
	anvilContainer, err := testbed.NewAnvilContainerWithOptions(
		ctx,
		testbed.AnvilOptions{
			ExposeHostPort: true,
			HostPort:       "8545",
			Logger:         config.Logger,
			Network:        config.Network,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to start anvil: %w", err)
	}
	harness.Anvil = anvilContainer

	// Create eth client for contract interactions (after Anvil is running)
	ethClient, err := geth.NewMultiHomingClient(geth.EthClientConfig{
		RPCURLs:          []string{config.TestConfig.Deployers[0].RPC},
		PrivateKeyString: config.TestConfig.Pks.EcdsaMap[config.TestConfig.EigenDA.Deployer].PrivateKey[2:],
		NumConfirmations: 0,
		NumRetries:       3,
	}, gethcommon.Address{}, config.Logger)
	if err != nil {
		return nil, fmt.Errorf("could not create eth client for registration: %w", err)
	}
	harness.EthClient = ethClient

	// Step 2: Setup Graph Node if needed
	deployer, ok := config.TestConfig.GetDeployer(config.TestConfig.EigenDA.Deployer)
	if ok && deployer.DeploySubgraphs {
		config.Logger.Info("Starting graph node")
		anvilInternalEndpoint := harness.GetAnvilInternalEndpoint()
		graphNodeContainer, err := testbed.NewGraphNodeContainerWithOptions(
			ctx,
			testbed.GraphNodeOptions{
				PostgresDB:     "graph-node",
				PostgresUser:   "graph-node",
				PostgresPass:   "let-me-in",
				EthereumRPC:    anvilInternalEndpoint,
				ExposeHostPort: true,
				HostHTTPPort:   "8000",
				HostWSPort:     "8001",
				HostAdminPort:  "8020",
				HostIPFSPort:   "5001",
				Logger:         config.Logger,
				Network:        config.Network,
			})
		if err != nil {
			return nil, fmt.Errorf("failed to start graph node: %w", err)
		}
		harness.GraphNode = graphNodeContainer
	}

	// Step 3: Deploy contracts
	config.Logger.Info("Deploying experiment")
	err = config.TestConfig.DeployExperiment()
	if err != nil {
		return nil, fmt.Errorf("failed to deploy experiment: %w", err)
	}

	// Register blob versions
	config.TestConfig.RegisterBlobVersions(harness.EthClient)

	// Register relay URLs
	relayURLs := []string{
		"localhost:32035",
		"localhost:32037",
		"localhost:32039",
		"localhost:32041",
	}
	config.TestConfig.RegisterRelays(harness.EthClient, relayURLs, harness.EthClient.GetAccountAddress())

	// Step 4: Start Churner (requires deployed contracts)
	config.Logger.Info("Starting churner server")
	err = startChurner(harness, config)
	if err != nil {
		return nil, fmt.Errorf("failed to start churner server: %w", err)
	}
	config.Logger.Info("Churner server started", "address", harness.Churner.URL)

	return harness, nil
}

// GetAnvilInternalEndpoint returns the internal Docker network endpoint for Anvil
func (ch *ChainHarness) GetAnvilInternalEndpoint() string {
	if ch.Anvil == nil {
		return ""
	}
	return ch.Anvil.InternalEndpoint()
}

// GetAnvilRPCUrl returns the external RPC URL for Anvil
func (ch *ChainHarness) GetAnvilRPCUrl() string {
	if ch.Anvil == nil {
		return ""
	}
	return ch.Anvil.RpcURL()
}

// Cleanup releases resources held by the ChainHarness (excluding shared network)
func (ch *ChainHarness) Cleanup(ctx context.Context, logger logging.Logger) {
	if ch.Churner.Server != nil {
		logger.Info("Stopping churner server")
		ch.Churner.Server.GracefulStop()
		if ch.Churner.Listener != nil {
			_ = ch.Churner.Listener.Close()
		}
	}

	if ch.GraphNode != nil {
		logger.Info("Stopping graph node")
		if err := ch.GraphNode.Terminate(ctx); err != nil {
			logger.Warn("Failed to terminate graph node container", "error", err)
		}
	}

	if ch.Anvil != nil {
		logger.Info("Stopping anvil")
		if err := ch.Anvil.Terminate(ctx); err != nil {
			logger.Warn("Failed to terminate anvil container", "error", err)
		}
	}
}

// startChurner starts the churner server
func startChurner(harness *ChainHarness, config *ChainHarnessConfig) error {
	// Get Anvil RPC URL using the getter method
	anvilRPC := harness.GetAnvilRPCUrl()

	// Get deployer's private key
	var privateKey string
	deployer, ok := config.TestConfig.GetDeployer(config.TestConfig.EigenDA.Deployer)
	if ok && deployer.Name != "" {
		privateKey = strings.TrimPrefix(config.TestConfig.Pks.EcdsaMap[deployer.Name].PrivateKey, "0x")
	}

	// Create logs directory
	logsDir := fmt.Sprintf("testdata/%s/logs", config.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/churner.log", logsDir)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open churner log file: %w", err)
	}

	// Create churner configuration
	churnerConfig := &churner.Config{
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
		MetricsConfig: churner.MetricsConfig{
			HTTPPort:      "9095",
			EnableMetrics: true,
		},
		OperatorStateRetrieverAddr: config.TestConfig.EigenDA.OperatorStateRetriever,
		EigenDAServiceManagerAddr:  config.TestConfig.EigenDA.ServiceManager,
		EigenDADirectory:           config.TestConfig.EigenDA.EigenDADirectory,
		GRPCPort:                   "32002",
		ChurnApprovalInterval:      15 * time.Minute,
		PerPublicKeyRateLimit:      1 * time.Second,
	}

	// Set graph URL if graph node is enabled
	if deployer.DeploySubgraphs && harness.GraphNode != nil {
		churnerConfig.ChainStateConfig = thegraph.Config{
			Endpoint: "http://localhost:8000/subgraphs/name/Layr-Labs/eigenda-operator-state",
		}
	}

	// Create churner logger
	churnerLogger, err := common.NewLogger(&churnerConfig.LoggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create churner logger: %w", err)
	}

	// Create geth client
	gethClient, err := geth.NewMultiHomingClient(churnerConfig.EthClientConfig, gethcommon.Address{}, churnerLogger)
	if err != nil {
		return fmt.Errorf("failed to create geth client: %w", err)
	}

	// Create writer
	churnerTx, err := coreeth.NewWriter(
		churnerLogger,
		gethClient,
		churnerConfig.OperatorStateRetrieverAddr,
		churnerConfig.EigenDAServiceManagerAddr)
	if err != nil {
		return fmt.Errorf("failed to create writer: %w", err)
	}

	// Create indexer
	chainState := coreeth.NewChainState(churnerTx, gethClient)
	indexer := thegraph.MakeIndexedChainState(churnerConfig.ChainStateConfig, chainState, churnerLogger)

	// Create churner
	churnerMetrics := churner.NewMetrics(churnerConfig.MetricsConfig.HTTPPort, churnerLogger)
	churnerInstance, err := churner.NewChurner(churnerConfig, indexer, churnerTx, churnerLogger, churnerMetrics)
	if err != nil {
		return fmt.Errorf("failed to create churner: %w", err)
	}

	// Create churner server
	churnerSvr := churner.NewServer(churnerConfig, churnerInstance, churnerLogger, churnerMetrics)
	err = churnerSvr.Start(churnerConfig.MetricsConfig)
	if err != nil {
		return fmt.Errorf("failed to start churner server metrics: %w", err)
	}

	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", churnerConfig.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", churnerConfig.GRPCPort, err)
	}
	harness.Churner.Listener = listener

	// Create and start gRPC server
	harness.Churner.Server = grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 300))
	pb.RegisterChurnerServer(harness.Churner.Server, churnerSvr)
	healthcheck.RegisterHealthServer(pb.Churner_ServiceDesc.ServiceName, harness.Churner.Server)

	// Start serving in goroutine
	go func() {
		churnerLogger.Info("Starting churner gRPC server", "port", churnerConfig.GRPCPort)
		if err := harness.Churner.Server.Serve(harness.Churner.Listener); err != nil {
			churnerLogger.Info("Churner gRPC server stopped", "error", err)
		}
	}()

	// TODO: Replace with proper health check endpoint
	time.Sleep(100 * time.Millisecond)
	churnerLogger.Info("Churner server started successfully", "port", churnerConfig.GRPCPort, "logFile", logFilePath)

	// Store the churner RPC address
	harness.Churner.URL = fmt.Sprintf("localhost:%s", churnerConfig.GRPCPort)
	return nil
}
