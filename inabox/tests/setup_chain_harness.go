package integration

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/testcontainers/testcontainers-go"
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
