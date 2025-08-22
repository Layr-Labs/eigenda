package testinfra

import (
	"context"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	caws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/testinfra/containers"
	"github.com/Layr-Labs/eigenda/common/testinfra/deployment"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
		
		// Populate AWS configuration for tests
		im.result.AWSConfig = &AWSTestConfig{
			EndpointURL:     localstack.Endpoint(),
			Region:          localstack.Region(),
			AccessKeyID:     "localstack",
			SecretAccessKey: "localstack",
		}
		
		// Deploy AWS resources if configured
		if im.config.LocalStack.DeployResources {
			fmt.Println("Deploying LocalStack AWS resources...")
			err = deployment.DeployLocalStackResources(ctx, localstack, im.config.LocalStack.Resources)
			if err != nil {
				return nil, fmt.Errorf("failed to deploy LocalStack resources: %w", err)
			}
			
			// Add deployed resource names to the result
			im.result.AWSConfig.BucketName = im.config.LocalStack.Resources.BucketName
			im.result.AWSConfig.MetadataTableName = im.config.LocalStack.Resources.MetadataTableName
			im.result.AWSConfig.BucketTableName = im.config.LocalStack.Resources.BucketTableName
			im.result.AWSConfig.V2MetadataTableName = im.config.LocalStack.Resources.V2MetadataTableName
			im.result.AWSConfig.V2PaymentPrefix = im.config.LocalStack.Resources.V2PaymentPrefix
			
			fmt.Println("Successfully deployed LocalStack AWS resources")
		}
		
		// Set AWS environment variables for tests to use LocalStack
		if err := im.result.AWSConfig.SetEnvironmentVariables(); err != nil {
			return nil, fmt.Errorf("failed to set AWS environment variables: %w", err)
		}
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

	// 4. Deploy EigenDA contracts if enabled (depends on Anvil)
	if im.config.EigenDA.Enabled && im.anvil != nil {
		// Set the RPC URL to use the Anvil instance we just started
		if im.config.EigenDA.Deployer.RPC == "" {
			im.config.EigenDA.Deployer.RPC = im.result.AnvilRPC
		}
		
		err := im.deployEigenDAContracts(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy EigenDA contracts: %w", err)
		}
	}

	success = true
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

// StartWithEigenDA starts infrastructure with EigenDA contract deployment enabled
func StartWithEigenDA(ctx context.Context, rootPath string) (*InfraManager, *InfraResult, error) {
	config := DefaultEigenDAConfig(rootPath)
	return StartCustom(ctx, config)
}

// generateDisperserKeypair generates a KMS keypair for the disperser
func (im *InfraManager) generateDisperserKeypair() (string, gethcommon.Address, error) {
	if im.localstack == nil {
		return "", gethcommon.Address{}, fmt.Errorf("LocalStack not available for KMS key generation")
	}

	// Get LocalStack endpoint
	endpoint := im.localstack.Endpoint()
	
	// Create KMS client
	keyManager := kms.New(kms.Options{
		Region:       "us-east-1",
		BaseEndpoint: aws.String(endpoint),
	})

	// Create the KMS key
	createKeyOutput, err := keyManager.CreateKey(context.Background(), &kms.CreateKeyInput{
		KeySpec:  types.KeySpecEccSecgP256k1,
		KeyUsage: types.KeyUsageTypeSignVerify,
	})
	if err != nil {
		return "", gethcommon.Address{}, fmt.Errorf("could not create KMS key: %w", err)
	}

	keyID := *createKeyOutput.KeyMetadata.KeyId

	// Load the public key and convert to address
	publicKey, err := caws.LoadPublicKeyKMS(context.Background(), keyManager, keyID)
	if err != nil {
		return "", gethcommon.Address{}, fmt.Errorf("could not load public key: %w", err)
	}

	address := crypto.PubkeyToAddress(*publicKey)
	
	return keyID, address, nil
}

// deployEigenDAContracts handles EigenDA contract deployment and registration
func (im *InfraManager) deployEigenDAContracts(_ context.Context) error {
	// Get private key for deployer account
	privateKey, err := im.anvil.GetPrivateKey(0)
	if err != nil {
		return fmt.Errorf("failed to get private key: %w", err)
	}

	// Clean the private key
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKey = strings.TrimPrefix(privateKey, "0X")

	// Create logger (assuming a default logger is available)
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Create eth client for contract operations
	ethClient, err := geth.NewMultiHomingClient(geth.EthClientConfig{
		RPCURLs:          []string{im.result.AnvilRPC},
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
		NumRetries:       3,
	}, gethcommon.Address{}, logger)
	if err != nil {
		return fmt.Errorf("failed to create eth client: %w", err)
	}

	// Create contract deployment manager using testinfra deployment utilities
	manager := deployment.NewContractDeploymentManager(im.config.EigenDA.RootPath, im.config.EigenDA.Deployer)

	// Deploy contracts if requested
	if im.config.EigenDA.DeployContracts {
		// Create deployment config matching the inabox testconfig-anvil.yaml structure
		numStrategies := 2 // Two strategies as per the stakes configuration
		maxOperatorCount := 3 // maxOperatorCount from services.counts
		
		// Stakes configuration from testconfig-anvil.yaml:
		// - total: 100e18, distribution: [1, 4, 6, 10] 
		// - total: 100e18, distribution: [1, 3, 8, 9]
		stakeDistribution := [][]float32{
			{1, 4, 6, 10},  // Strategy 0 distribution
			{1, 3, 8, 9},   // Strategy 1 distribution  
		}
		stakeTotals := []float32{100e18, 100e18} // 100e18 tokens per strategy

		// Private keys must be provided in the config
		if len(im.config.EigenDA.PrivateKeys) == 0 {
			return fmt.Errorf("private keys for operators and stakers must be provided in EigenDA config")
		}
		privateKeys := im.config.EigenDA.PrivateKeys

		deployConfig := deployment.GenerateEigenDADeployConfig(
			numStrategies,
			maxOperatorCount,
			stakeDistribution,
			stakeTotals,
			privateKeys,
		)

		err := manager.DeployEigenDAContracts(deployConfig)
		if err != nil {
			return fmt.Errorf("failed to deploy EigenDA contracts: %w", err)
		}
		
		// The deployment manager populates its own EigenDAContracts field
		// Copy it to our result
		im.result.EigenDAContracts = manager.EigenDAContracts
	}

	// Handle disperser keypair generation and registration
	if im.config.EigenDA.GenerateDisperserKeypair && im.localstack != nil {
		// Generate KMS keypair for disperser
		keyID, address, err := im.generateDisperserKeypair()
		if err != nil {
			return fmt.Errorf("failed to generate disperser keypair: %w", err)
		}
		
		// Store in config and result for output
		im.config.EigenDA.DisperserKMSKeyID = keyID
		im.config.EigenDA.DisperserAddress = address.Hex()
		im.result.DisperserKMSKeyID = keyID
		im.result.DisperserAddress = address
		
		// Register the disperser address if contracts are available
		if im.result.EigenDAContracts != nil {
			err = manager.RegisterDisperserAddress(ethClient, address)
			if err != nil {
				return fmt.Errorf("failed to register disperser address: %w", err)
			}
		}
	} else if im.config.EigenDA.RegisterDisperser && im.result.EigenDAContracts != nil {
		// Use provided disperser address if available
		if im.config.EigenDA.DisperserAddress != "" {
			disperserAddr := gethcommon.HexToAddress(im.config.EigenDA.DisperserAddress)
			err = manager.RegisterDisperserAddress(ethClient, disperserAddr)
			if err != nil {
				return fmt.Errorf("failed to register disperser address: %w", err)
			}
		} else {
			// Fallback: derive disperser address from deployer private key
			privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
			if err != nil {
				return fmt.Errorf("failed to parse private key: %w", err)
			}
			disperserAddress := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)

			err = manager.RegisterDisperserAddress(ethClient, disperserAddress)
			if err != nil {
				return fmt.Errorf("failed to register disperser keypair: %w", err)
			}
		}
	}

	// Register blob versions and relays if requested and contracts are available
	if im.config.EigenDA.RegisterBlobVersionAndRelays && im.result.EigenDAContracts != nil {
		// Set the service manager address in manager (used by RegisterBlobVersionsAndRelays)
		manager.EigenDAContracts.ServiceManager = im.result.EigenDAContracts.ServiceManager

		// Register relay addresses for v2
		// These are the standard relay ports used in inabox tests
		relayPorts := []string{"32035", "32037", "32039", "32041"}

		err = manager.RegisterBlobVersionsAndRelays(ethClient, im.config.EigenDA.BlobVersionParams, relayPorts)
		if err != nil {
			return fmt.Errorf("failed to register blob versions and relays: %w", err)
		}
	}

	return nil
}
