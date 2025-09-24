package testbed

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/operators/churner"
	churnerflags "github.com/Layr-Labs/eigenda/operators/churner/flags"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ChurnerConfig defines configuration for the churner container
//
//nolint:lll // struct field documentation
type ChurnerConfig struct {
	// Enable churner service
	Enabled bool `env:"-"` // Skip in env mapping

	// Log configuration
	LogFormat string `env:"CHURNER_LOG_FORMAT"`
	LogLevel  string `env:"CHURNER_LOG_LEVEL"`

	// Network configuration
	Hostname string `env:"CHURNER_HOSTNAME"`
	GRPCPort string `env:"CHURNER_GRPC_PORT"`

	// EigenDA contract addresses
	EigenDADirectory       string `env:"CHURNER_EIGENDA_DIRECTORY"`
	OperatorStateRetriever string `env:"CHURNER_BLS_OPERATOR_STATE_RETRIVER"` // Note: typo in original
	ServiceManager         string `env:"CHURNER_EIGENDA_SERVICE_MANAGER"`

	// Chain configuration
	ChainRPC   string `env:"CHURNER_CHAIN_RPC"`
	PrivateKey string `env:"CHURNER_PRIVATE_KEY"`

	// Graph configuration
	GraphURL string `env:"CHURNER_GRAPH_URL"`

	// Metrics configuration
	EnableMetrics   bool   `env:"CHURNER_ENABLE_METRICS"`
	MetricsHTTPPort string `env:"CHURNER_METRICS_HTTP_PORT"`

	// Churner specific configuration
	PerPublicKeyRateLimit time.Duration `env:"CHURNER_PER_PUBLIC_KEY_RATE_LIMIT"`
	ChurnApprovalInterval time.Duration `env:"CHURNER_CHURN_APPROVAL_INTERVAL"`

	// Container configuration (not exposed as env vars)
	Image          string        `env:"-"`
	StartupTimeout time.Duration `env:"-"`
	ExposeHostPort bool          `env:"-"` // If true, binds container port to host port
	HostPort       string        `env:"-"` // Custom host port to bind to (defaults to GRPCPort if empty and ExposeHostPort is true)

	LogPath string `env:"CHURNER_LOG_PATH"` // Additional env vars that don't have direct struct fields

	// Test data directory path for logs and env files (not exposed as env var)
	TestDataPath string `env:"-"` // Path to inabox/testdata/<timestamp> directory

	// MockIndexer for testing - if provided, will be used instead of creating a real graph indexer
	MockIndexer interface{} `env:"-"` // Will be cast to thegraph.IndexedChainState
}

type ChurnerGoroutine struct {
	server   *grpc.Server
	listener net.Listener
	url      string
	cancel   context.CancelFunc
	config   ChurnerConfig
}

func StartChurnerGoroutine(config ChurnerConfig, logger logging.Logger) (*ChurnerGoroutine, error) {
	// Start listener on the specified port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCPort))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %s: %w", config.GRPCPort, err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Create FlagSet and apply churner flags
	fs := flag.NewFlagSet("churner", flag.ContinueOnError)
	for _, f := range churnerflags.Flags {
		f.Apply(fs)
	}

	// Parse the flag values
	args := buildChurnerArgs(config)
	err = fs.Parse(args)
	if err != nil {
		cancel()
		_ = listener.Close()
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	// Create CLI context with the FlagSet
	app := cli.NewApp()
	cliCtx := cli.NewContext(app, fs, nil)

	// Create the gRPC server
	gs := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024 * 1024 * 300),
	)

	// Channel to communicate startup errors
	errChan := make(chan error, 1)
	readyChan := make(chan struct{})

	// Start churner in goroutine
	go func() {
		defer gs.Stop()

		// Parse config
		churnerConfig, err := churner.NewConfig(cliCtx)
		if err != nil {
			logger.Error("Failed to parse churner config", "error", err)
			errChan <- fmt.Errorf("failed to parse churner config: %w", err)
			return
		}

		// Create logger from churner config (which respects the log path)
		churnerConfig.LoggerConfig.Format = common.TextLogFormat
		churnerConfig.LoggerConfig.HandlerOpts.NoColor = true
		churnerLogger, err := common.NewLogger(&churnerConfig.LoggerConfig)
		if err != nil {
			logger.Error("Failed to create churner logger", "error", err)
			errChan <- fmt.Errorf("failed to create churner logger: %w", err)
			return
		}
		logger = churnerLogger

		// Create clients
		gethClient, err := geth.NewMultiHomingClient(churnerConfig.EthClientConfig, gethcommon.Address{}, logger)
		if err != nil {
			logger.Error("Failed to create geth client", "error", err)
			errChan <- fmt.Errorf("failed to create geth client: %w", err)
			return
		}

		tx, err := coreeth.NewWriter(
			logger,
			gethClient,
			churnerConfig.OperatorStateRetrieverAddr,
			churnerConfig.EigenDAServiceManagerAddr)
		if err != nil {
			logger.Error("Failed to create writer", "error", err)
			errChan <- fmt.Errorf("failed to create writer: %w", err)
			return
		}

		cs := coreeth.NewChainState(tx, gethClient)

		// Use mock indexer if provided, otherwise create real graph indexer
		var indexer thegraph.IndexedChainState
		if config.MockIndexer != nil {
			var ok bool
			indexer, ok = config.MockIndexer.(thegraph.IndexedChainState)
			if !ok {
				logger.Error("Failed to cast MockIndexer to IndexedChainState")
				errChan <- fmt.Errorf("failed to cast MockIndexer to IndexedChainState")
				return
			}
		} else {
			indexer = thegraph.MakeIndexedChainState(churnerConfig.ChainStateConfig, cs, logger)
		}

		churnerMetrics := churner.NewMetrics(churnerConfig.MetricsConfig.HTTPPort, logger)
		cn, err := churner.NewChurner(churnerConfig, indexer, tx, logger, churnerMetrics)
		if err != nil {
			logger.Error("Failed to create churner", "error", err)
			errChan <- fmt.Errorf("failed to create churner: %w", err)
			return
		}

		churnerServer := churner.NewServer(churnerConfig, cn, logger, churnerMetrics)
		if err := churnerServer.Start(churnerConfig.MetricsConfig); err != nil {
			logger.Error("Failed to start churner server metrics", "error", err)
		}

		// Register services
		reflection.Register(gs)
		pb.RegisterChurnerServer(gs, churnerServer)
		healthcheck.RegisterHealthServer(pb.Churner_ServiceDesc.ServiceName, gs)

		logger.Info("Churner server starting", "port", config.GRPCPort)

		// Signal that the server is ready
		close(readyChan)

		// Start serving
		if err := gs.Serve(listener); err != nil && ctx.Err() == nil {
			logger.Error("Failed to serve", "error", err)
			errChan <- fmt.Errorf("failed to serve: %w", err)
		}
	}()

	// Wait for server to be ready or fail
	select {
	case err := <-errChan:
		cancel()
		_ = listener.Close()
		return nil, err
	case <-readyChan:
		// Server started successfully
		logger.Info("Churner server started successfully", "port", config.GRPCPort)
	case <-time.After(5 * time.Second):
		cancel()
		_ = listener.Close()
		return nil, fmt.Errorf("churner server startup timeout")
	}

	url := fmt.Sprintf("localhost:%s", config.GRPCPort)
	return &ChurnerGoroutine{
		server:   gs,
		listener: listener,
		url:      url,
		cancel:   cancel,
		config:   config,
	}, nil
}

// DefaultChurnerConfig returns a default churner configuration suitable for testing
func DefaultChurnerConfig() ChurnerConfig {
	return ChurnerConfig{
		Enabled:                true,
		LogFormat:              "text",
		LogLevel:               "debug",
		Hostname:               "0.0.0.0",
		GRPCPort:               "32002",
		EigenDADirectory:       "", // Will be populated from contract deployment
		OperatorStateRetriever: "", // Will be populated from contract deployment
		ServiceManager:         "", // Will be populated from contract deployment
		ChainRPC:               "", // Will be populated from Anvil
		PrivateKey:             "", // Will be populated from deployer key
		GraphURL:               "", // Will be populated from GraphNode if enabled
		EnableMetrics:          true,
		MetricsHTTPPort:        "9095",
		PerPublicKeyRateLimit:  1 * time.Second, // Fast for testing
		ChurnApprovalInterval:  900 * time.Second,
		Image:                  "ghcr.io/layr-labs/eigenda/churner:dev",
		StartupTimeout:         30 * time.Second,
		ExposeHostPort:         false,
		HostPort:               "",
	}
}

// buildChurnerArgs builds the command line arguments from the ChurnerConfig
func buildChurnerArgs(config ChurnerConfig) []string {
	args := []string{
		"--churner.grpc-port", config.GRPCPort,
		"--chain.rpc", config.ChainRPC,
		"--chain.private-key", config.PrivateKey,
		"--churner.bls-operator-state-retriever", config.OperatorStateRetriever,
		"--churner.eigenda-service-manager", config.ServiceManager,
		"--churner.eigenda-directory", config.EigenDADirectory,
		"--churner.hostname", config.Hostname,
		"--churner.log.level", config.LogLevel,
	}

	// Add log path if specified
	if config.LogPath != "" {
		args = append(args, "--churner.log.path", config.LogPath)
		// Force text format to match test logger
		args = append(args, "--churner.log.format", "text")
	}

	// Add optional parameters
	if config.GraphURL != "" {
		args = append(args, "--thegraph.endpoint", config.GraphURL)
	}

	if config.EnableMetrics {
		args = append(args, "--churner.enable-metrics", strconv.FormatBool(config.EnableMetrics))
		if config.MetricsHTTPPort != "" {
			args = append(args, "--churner.metrics-http-port", config.MetricsHTTPPort)
		}
	}

	if config.PerPublicKeyRateLimit > 0 {
		args = append(args, "--churner.per-public-key-rate-limit", config.PerPublicKeyRateLimit.String())
	}

	if config.ChurnApprovalInterval > 0 {
		args = append(args, "--churner.churn-approval-interval", config.ChurnApprovalInterval.String())
	}

	return args
}

// URL returns the churner service URL
func (c *ChurnerGoroutine) URL() string {
	return c.url
}

// Config returns the churner configuration
func (c *ChurnerGoroutine) Config() ChurnerConfig {
	return c.config
}

// Stop gracefully stops the churner goroutine
func (c *ChurnerGoroutine) Stop(ctx context.Context) {
	c.cancel()
	c.server.GracefulStop()
}
