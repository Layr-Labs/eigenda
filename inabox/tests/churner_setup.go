package integration_test

import (
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
	"github.com/Layr-Labs/eigenda/operators/churner"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
)

// StartChurnerForInfrastructure starts the churner server as part of the global infrastructure.
// This should be called after Anvil and other containers are started.
// Returns the churner RPC address that can be used by operators.
func StartChurnerForInfrastructure(infra *InfrastructureHarness, anvilRPC string) (string, error) {
	// Get deployer's private key
	var privateKey string
	deployer, ok := infra.TestConfig.GetDeployer(infra.TestConfig.EigenDA.Deployer)
	if ok && deployer.Name != "" {
		privateKey = strings.TrimPrefix(infra.TestConfig.Pks.EcdsaMap[deployer.Name].PrivateKey, "0x")
	}

	// Create logs directory
	logsDir := fmt.Sprintf("testdata/%s/logs", infra.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/churner.log", logsDir)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open churner log file: %w", err)
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
		OperatorStateRetrieverAddr: infra.TestConfig.EigenDA.OperatorStateRetriever,
		EigenDAServiceManagerAddr:  infra.TestConfig.EigenDA.ServiceManager,
		EigenDADirectory:           infra.TestConfig.EigenDA.EigenDADirectory,
		GRPCPort:                   "32002",
		ChurnApprovalInterval:      15 * time.Minute,
		PerPublicKeyRateLimit:      1 * time.Second,
	}

	// Set graph URL if graph node is enabled
	if deployer.DeploySubgraphs && infra.GraphNodeContainer != nil {
		churnerConfig.ChainStateConfig = thegraph.Config{
			Endpoint: "http://localhost:8000/subgraphs/name/Layr-Labs/eigenda-operator-state",
		}
	}

	// Create churner logger
	churnerLogger, err := common.NewLogger(&churnerConfig.LoggerConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create churner logger: %w", err)
	}

	// Create geth client
	gethClient, err := geth.NewMultiHomingClient(churnerConfig.EthClientConfig, gethcommon.Address{}, churnerLogger)
	if err != nil {
		return "", fmt.Errorf("failed to create geth client: %w", err)
	}

	// Create writer
	churnerTx, err := coreeth.NewWriter(
		churnerLogger,
		gethClient,
		churnerConfig.OperatorStateRetrieverAddr,
		churnerConfig.EigenDAServiceManagerAddr)
	if err != nil {
		return "", fmt.Errorf("failed to create writer: %w", err)
	}

	// Create indexer
	chainState := coreeth.NewChainState(churnerTx, gethClient)
	indexer := thegraph.MakeIndexedChainState(churnerConfig.ChainStateConfig, chainState, churnerLogger)

	// Create churner
	churnerMetrics := churner.NewMetrics(churnerConfig.MetricsConfig.HTTPPort, churnerLogger)
	churnerInstance, err := churner.NewChurner(churnerConfig, indexer, churnerTx, churnerLogger, churnerMetrics)
	if err != nil {
		return "", fmt.Errorf("failed to create churner: %w", err)
	}

	// Create churner server
	churnerSvr := churner.NewServer(churnerConfig, churnerInstance, churnerLogger, churnerMetrics)
	err = churnerSvr.Start(churnerConfig.MetricsConfig)
	if err != nil {
		return "", fmt.Errorf("failed to start churner server metrics: %w", err)
	}

	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", churnerConfig.GRPCPort))
	if err != nil {
		return "", fmt.Errorf("failed to listen on port %s: %w", churnerConfig.GRPCPort, err)
	}
	infra.ChurnerListener = listener

	// Create and start gRPC server
	infra.ChurnerServer = grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 300))
	pb.RegisterChurnerServer(infra.ChurnerServer, churnerSvr)
	healthcheck.RegisterHealthServer(pb.Churner_ServiceDesc.ServiceName, infra.ChurnerServer)

	// Start serving in goroutine
	go func() {
		churnerLogger.Info("Starting churner gRPC server", "port", churnerConfig.GRPCPort)
		if err := infra.ChurnerServer.Serve(infra.ChurnerListener); err != nil {
			churnerLogger.Info("Churner gRPC server stopped", "error", err)
		}
	}()

	// TODO: Replace with proper health check endpoint
	time.Sleep(100 * time.Millisecond)
	churnerLogger.Info("Churner server started successfully", "port", churnerConfig.GRPCPort, "logFile", logFilePath)

	// Return the churner RPC address for operators to use
	return fmt.Sprintf("localhost:%s", churnerConfig.GRPCPort), nil
}
