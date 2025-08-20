package testinfra

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common/testinfra/containers"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

// InfraManager orchestrates the lifecycle of test infrastructure containers
type InfraManager struct {
	config     InfraConfig
	anvil      *containers.AnvilContainer
	localstack *containers.LocalStackContainer
	graphnode  *containers.GraphNodeContainer
	network    *testcontainers.DockerNetwork
	result     InfraResult
}

// NewInfraManager creates a new infrastructure manager with the given configuration
func NewInfraManager(config InfraConfig) *InfraManager {
	return &InfraManager{
		config: config,
	}
}

// Start initializes and starts all enabled infrastructure components
func (im *InfraManager) Start(ctx context.Context) (*InfraResult, error) {
	var success bool
	defer func() {
		if !success {
			im.cleanup(ctx)
		}
	}()

	// Create a shared network for all containers to communicate
	sharedNetwork, err := network.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create shared network: %w", err)
	}
	im.network = sharedNetwork

	// Start containers in dependency order

	// 1. Start Anvil blockchain if enabled
	if im.config.Anvil.Enabled {
		anvilConfig := containers.AnvilConfig{
			Enabled:   im.config.Anvil.Enabled,
			ChainID:   im.config.Anvil.ChainID,
			BlockTime: im.config.Anvil.BlockTime,
			GasLimit:  im.config.Anvil.GasLimit,
			GasPrice:  im.config.Anvil.GasPrice,
			Accounts:  im.config.Anvil.Accounts,
			Mnemonic:  im.config.Anvil.Mnemonic,
			Fork:      im.config.Anvil.Fork,
			ForkBlock: im.config.Anvil.ForkBlock,
		}
		anvil, err := containers.NewAnvilContainerWithNetwork(ctx, anvilConfig, sharedNetwork)
		if err != nil {
			return nil, fmt.Errorf("failed to start anvil: %w", err)
		}
		im.anvil = anvil
		im.result.AnvilRPC = anvil.RPCURL()
		im.result.AnvilChainID = anvil.ChainID()
	}

	// 2. Start LocalStack if enabled
	if im.config.LocalStack.Enabled {
		localstackConfig := containers.LocalStackConfig{
			Enabled:  im.config.LocalStack.Enabled,
			Services: im.config.LocalStack.Services,
			Region:   im.config.LocalStack.Region,
			Debug:    im.config.LocalStack.Debug,
		}
		localstack, err := containers.NewLocalStackContainerWithNetwork(ctx, localstackConfig, sharedNetwork)
		if err != nil {
			return nil, fmt.Errorf("failed to start localstack: %w", err)
		}
		im.localstack = localstack
		im.result.LocalStackURL = localstack.Endpoint()
	}

	// 3. Start Graph Node if enabled (depends on Anvil for Ethereum RPC)
	if im.config.GraphNode.Enabled {
		var ethereumRPC string

		// Use internal container network URL for Graph Node to reach Anvil
		if im.anvil != nil {
			internalRPC, err := im.anvil.InternalRPCURL(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get anvil internal RPC URL: %w", err)
			}
			ethereumRPC = internalRPC
		} else if im.config.GraphNode.EthereumRPC != "" {
			ethereumRPC = im.config.GraphNode.EthereumRPC
		} else {
			return nil, fmt.Errorf("graph node requires ethereum RPC but none provided and anvil not enabled")
		}

		graphnodeConfig := containers.GraphNodeConfig{
			Enabled:      im.config.GraphNode.Enabled,
			PostgresDB:   im.config.GraphNode.PostgresDB,
			PostgresUser: im.config.GraphNode.PostgresUser,
			PostgresPass: im.config.GraphNode.PostgresPass,
			EthereumRPC:  im.config.GraphNode.EthereumRPC,
			IPFSEndpoint: im.config.GraphNode.IPFSEndpoint,
		}
		graphnode, err := containers.NewGraphNodeContainerWithNetwork(ctx, graphnodeConfig, ethereumRPC, sharedNetwork)
		if err != nil {
			return nil, fmt.Errorf("failed to start graph node: %w", err)
		}
		im.graphnode = graphnode
		im.result.GraphNodeURL = graphnode.HTTPURL()
		im.result.GraphNodeAdminURL = graphnode.AdminURL()

		// Get IPFS URL if available
		if ipfsURL, err := graphnode.IPFSURL(ctx); err == nil {
			im.result.IPFSURL = ipfsURL
		}

		// Also expose the PostgreSQL URL for direct database access if needed
		if postgresContainer := graphnode.GetPostgres(); postgresContainer != nil {
			postgresHost, _ := postgresContainer.Host(ctx)
			postgresPort, _ := postgresContainer.MappedPort(ctx, "5432")
			if postgresHost != "" && postgresPort != "" {
				im.result.PostgresURL = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
					graphnodeConfig.PostgresUser, graphnodeConfig.PostgresPass, postgresHost, postgresPort.Port(), graphnodeConfig.PostgresDB)
			}
		}
	}

	return &im.result, nil
}

// Stop terminates all running containers
func (im *InfraManager) Stop(ctx context.Context) error {
	return im.cleanup(ctx)
}

// cleanup terminates all containers, collecting any errors
func (im *InfraManager) cleanup(ctx context.Context) error {
	var errs []error

	// Terminate in reverse dependency order
	if im.graphnode != nil {
		if err := im.graphnode.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate graph node: %w", err))
		}
		im.graphnode = nil
	}

	if im.localstack != nil {
		if err := im.localstack.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate localstack: %w", err))
		}
		im.localstack = nil
	}

	if im.anvil != nil {
		if err := im.anvil.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate anvil: %w", err))
		}
		im.anvil = nil
	}

	// Remove the shared network
	if im.network != nil {
		if err := im.network.Remove(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to remove shared network: %w", err))
		}
		im.network = nil
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during cleanup: %v", errs)
	}

	return nil
}

// GetAnvil returns the Anvil container if started
func (im *InfraManager) GetAnvil() *containers.AnvilContainer {
	return im.anvil
}

// GetLocalStack returns the LocalStack container if started
func (im *InfraManager) GetLocalStack() *containers.LocalStackContainer {
	return im.localstack
}

// GetGraphNode returns the Graph Node container if started
func (im *InfraManager) GetGraphNode() *containers.GraphNodeContainer {
	return im.graphnode
}

// GetResult returns the current infrastructure result
func (im *InfraManager) GetResult() *InfraResult {
	return &im.result
}

// StartMinimal starts only Anvil and LocalStack for basic testing
func StartMinimal(ctx context.Context) (*InfraManager, *InfraResult, error) {
	config := DefaultConfig()
	config.GraphNode.Enabled = false // Disable graph node for minimal setup

	manager := NewInfraManager(config)
	result, err := manager.Start(ctx)
	if err != nil {
		return nil, nil, err
	}

	return manager, result, nil
}

// StartFull starts all infrastructure components
func StartFull(ctx context.Context) (*InfraManager, *InfraResult, error) {
	config := DefaultConfig()
	config.GraphNode.Enabled = true // Enable all components

	manager := NewInfraManager(config)
	result, err := manager.Start(ctx)
	if err != nil {
		return nil, nil, err
	}

	return manager, result, nil
}

// StartCustom starts infrastructure with custom configuration
func StartCustom(ctx context.Context, config InfraConfig) (*InfraManager, *InfraResult, error) {
	manager := NewInfraManager(config)
	result, err := manager.Start(ctx)
	if err != nil {
		return nil, nil, err
	}

	return manager, result, nil
}
