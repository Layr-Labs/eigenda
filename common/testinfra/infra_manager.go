package testinfra

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	validatorclientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/common"
	caws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/testinfra/containers"
	"github.com/Layr-Labs/eigenda/common/testinfra/deployment"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
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

		// Initialize cert verification components if contracts were deployed
		if im.result.EigenDAContracts != nil {
			logger, err := common.NewLogger(common.DefaultLoggerConfig())
			if err != nil {
				return nil, fmt.Errorf("failed to create logger: %w", err)
			}

			// Get private key for eth client
			privateKey, err := im.anvil.GetPrivateKey(0)
			if err != nil {
				return nil, fmt.Errorf("failed to get private key: %w", err)
			}
			privateKey = strings.TrimPrefix(privateKey, "0x")
			privateKey = strings.TrimPrefix(privateKey, "0X")

			// Create eth client for cert verification
			ethClient, err := geth.NewMultiHomingClient(geth.EthClientConfig{
				RPCURLs:          []string{im.result.AnvilRPC},
				PrivateKeyString: privateKey,
				NumConfirmations: 0,
				NumRetries:       3,
			}, gethcommon.Address{}, logger)
			if err != nil {
				return nil, fmt.Errorf("failed to create eth client: %w", err)
			}

			// Initialize cert verification components
			certComponents, err := deployment.InitializeCertVerification(ctx, ethClient, logger, im.result.EigenDAContracts)
			if err != nil {
				return nil, fmt.Errorf("failed to initialize cert verification: %w", err)
			}
			im.result.CertVerification = certComponents
		}
	}

	// 5. Deploy subgraphs if Graph Node is enabled and EigenDA contracts are deployed
	if im.config.GraphNode.Enabled && im.graphnode != nil && im.result.EigenDAContracts != nil {
		// Log current working directory for debugging
		cwd, _ := os.Getwd()
		fmt.Printf("Current working directory: %s\n", cwd)
		fmt.Printf("EigenDA RootPath from config: %s\n", im.config.EigenDA.RootPath)
		
		// Prepare subgraph deployment config
		subgraphConfig := deployment.SubgraphDeploymentConfig{
			RootPath: im.config.EigenDA.RootPath,
			Subgraphs: []deployment.SubgraphConfig{
				{
					Name:    "eigenda-operator-state",
					Path:    "eigenda-operator-state",
					Enabled: true,
				},
			},
			EigenDAConfig: deployment.EigenDAContractAddresses{
				RegistryCoordinator: im.result.EigenDAContracts.RegistryCoordinator,
				ServiceManager:      im.result.EigenDAContracts.ServiceManager,
				BlsApkRegistry:      im.result.EigenDAContracts.BlsApkRegistry,
			},
		}
		
		fmt.Println("Deploying subgraphs to Graph Node...")
		err = deployment.DeploySubgraphs(ctx, im.graphnode, subgraphConfig, 0) // Start at block 0 for simplicity
		if err != nil {
			return nil, fmt.Errorf("failed to deploy subgraphs: %w", err)
		}
		fmt.Println("✅ Subgraphs deployed successfully")
		
		// Wait a bit for the subgraph to sync
		fmt.Println("Waiting for subgraph to sync...")
		time.Sleep(5 * time.Second)
	}

	// 6. Setup retrieval clients if enabled and EigenDA contracts are deployed
	if im.config.EigenDA.RetrievalClients.Enabled && im.result.EigenDAContracts != nil {
		fmt.Println("Setting up retrieval clients...")
		err := im.setupRetrievalClients(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to setup retrieval clients: %w", err)
		}
	}

	// 7. Setup payload disperser if enabled and required components are available
	if im.config.EigenDA.PayloadDisperser.Enabled && 
		im.result.CertVerification != nil && 
		im.result.RetrievalClients != nil {
		// Get deployer private key
		deployerPrivateKey := im.config.EigenDA.Deployer.PrivateKey
		if deployerPrivateKey == "" {
			// Use default anvil account 0 private key
			deployerPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		}
		err := im.setupPayloadDisperser(ctx, deployerPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to setup payload disperser: %w", err)
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

// setupRetrievalClients sets up all the retrieval clients
func (im *InfraManager) setupRetrievalClients(ctx context.Context) error {
	config := im.config.EigenDA.RetrievalClients
	
	// Skip if not enabled
	if !config.Enabled {
		return nil
	}

	// Save current working directory to restore it at the end
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	// Change to the inabox directory where KZG resources are located
	// This makes the relative paths work correctly
	inaboxDir := strings.TrimSuffix(im.config.EigenDA.RootPath, "/") + "/inabox"
	if err := os.Chdir(inaboxDir); err != nil {
		return fmt.Errorf("failed to change to inabox directory %s: %w", inaboxDir, err)
	}

	// Use Anvil RPC if not specified
	rpcURL := config.RPC
	if rpcURL == "" && im.anvil != nil {
		rpcURL = im.result.AnvilRPC
	}
	if rpcURL == "" {
		return fmt.Errorf("no RPC URL available for retrieval clients")
	}

	// Use contract addresses from deployment if not specified
	operatorStateRetriever := config.OperatorStateRetriever
	serviceManager := config.ServiceManager
	if operatorStateRetriever == "" && im.result.EigenDAContracts != nil {
		operatorStateRetriever = im.result.EigenDAContracts.OperatorStateRetriever
	}
	if serviceManager == "" && im.result.EigenDAContracts != nil {
		serviceManager = im.result.EigenDAContracts.ServiceManager
	}
	if operatorStateRetriever == "" || serviceManager == "" {
		return fmt.Errorf("contract addresses not available for retrieval clients")
	}

	// Create logger
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Setup eth client
	ethClientConfig := geth.EthClientConfig{
		RPCURLs:          []string{rpcURL},
		PrivateKeyString: "351b8eca372e64f64d514f90f223c5c4f86a04ff3dcead5c27293c547daab4ca", // just random private key
		NumConfirmations: 3,
		NumRetries:       0,
	}
	
	ethClient, err := geth.NewMultiHomingClient(ethClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		return fmt.Errorf("failed to create eth client: %w", err)
	}
	
	rpcClient, err := ethrpc.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to create RPC client: %w", err)
	}
	
	tx, err := eth.NewWriter(logger, ethClient, operatorStateRetriever, serviceManager)
	if err != nil {
		return fmt.Errorf("failed to create eth writer: %w", err)
	}

	cs := eth.NewChainState(tx, ethClient)
	agn := &core.StdAssignmentCoordinator{}
	nodeClient := clients.NewNodeClient(20 * time.Second)
	
	srsOrder, err := strconv.Atoi(config.SRSOrder)
	if err != nil {
		return fmt.Errorf("failed to parse SRS order: %w", err)
	}
	
	kzgConfig := &kzg.KzgConfig{
		G1Path:          config.G1Path,
		G2Path:          config.G2Path,
		G2PowerOf2Path:  config.G2PowerOf2Path,
		CacheDir:        config.CachePath,
		SRSOrder:        uint64(srsOrder),
		SRSNumberToLoad: uint64(srsOrder),
		NumWorker:       1,
		PreloadEncoder:  false,
		LoadG2Points:    true,
	}

	kzgVerifier, err := verifier.NewVerifier(kzgConfig, nil)
	if err != nil {
		return fmt.Errorf("failed to create KZG verifier: %w", err)
	}

	retrievalClient, err := clients.NewRetrievalClient(logger, cs, agn, nodeClient, kzgVerifier, 10)
	if err != nil {
		return fmt.Errorf("failed to create retrieval client: %w", err)
	}
	
	chainReader, err := eth.NewReader(logger, ethClient, operatorStateRetriever, serviceManager)
	if err != nil {
		return fmt.Errorf("failed to create chain reader: %w", err)
	}

	clientConfig := validatorclientsv2.DefaultClientConfig()
	retrievalClientV2 := validatorclientsv2.NewValidatorClient(logger, chainReader, cs, kzgVerifier, clientConfig, nil)

	validatorPayloadRetrieverConfig := payloadretrieval.ValidatorPayloadRetrieverConfig{
		PayloadClientConfig: *clientsv2.GetDefaultPayloadClientConfig(),
		RetrievalTimeout:    1 * time.Minute,
	}

	validatorRetrievalClientV2, err := payloadretrieval.NewValidatorPayloadRetriever(
		logger,
		validatorPayloadRetrieverConfig,
		retrievalClientV2,
		kzgVerifier.Srs.G1)
	if err != nil {
		return fmt.Errorf("failed to create validator retrieval client: %w", err)
	}

	relayClientConfig := &relay.RelayClientConfig{
		MaxGRPCMessageSize: 100 * 1024 * 1024, // 100 MB message size limit
	}

	relayUrlProvider, err := relay.NewRelayUrlProvider(ethClient, chainReader.GetRelayRegistryAddress())
	if err != nil {
		return fmt.Errorf("failed to create relay URL provider: %w", err)
	}

	relayClient, err := relay.NewRelayClient(relayClientConfig, logger, relayUrlProvider)
	if err != nil {
		return fmt.Errorf("failed to create relay client: %w", err)
	}

	relayPayloadRetrieverConfig := payloadretrieval.RelayPayloadRetrieverConfig{
		PayloadClientConfig: *clientsv2.GetDefaultPayloadClientConfig(),
		RelayTimeout:        5 * time.Second,
	}

	// Use a new random source for each client instance
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	relayRetrievalClientV2, err := payloadretrieval.NewRelayPayloadRetriever(
		logger,
		randSource,
		relayPayloadRetrieverConfig,
		relayClient,
		kzgVerifier.Srs.G1)
	if err != nil {
		return fmt.Errorf("failed to create relay retrieval client: %w", err)
	}

	// Store all clients in the result
	im.result.RetrievalClients = &RetrievalClientsComponents{
		EthClient:                  ethClient,
		RPCClient:                  rpcClient,
		RetrievalClient:            retrievalClient,
		ChainReader:                chainReader,
		RelayRetrievalClientV2:     relayRetrievalClientV2,
		ValidatorRetrievalClientV2: validatorRetrievalClientV2,
	}

	fmt.Println("✅ Retrieval clients initialized successfully")
	return nil
}

// setupPayloadDisperser sets up the payload disperser components
func (im *InfraManager) setupPayloadDisperser(ctx context.Context, deployerPrivateKey string) error {
	config := im.config.EigenDA.PayloadDisperser
	
	// Skip if not enabled
	if !config.Enabled {
		return nil
	}

	// Ensure cert verification and retrieval clients are available
	if im.result.CertVerification == nil {
		return fmt.Errorf("cert verification components required for payload disperser")
	}
	if im.result.RetrievalClients == nil {
		return fmt.Errorf("retrieval clients required for payload disperser")
	}

	fmt.Println("Setting up payload disperser...")

	// Create logger if not available
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Get eth client from retrieval clients
	ethClient := im.result.RetrievalClients.EthClient
	if ethClient == nil {
		return fmt.Errorf("eth client not available from retrieval clients")
	}

	// Get chain ID
	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create deployer transactor opts
	deployerTransactorOpts, err := deployment.CreateTransactorOpts(deployerPrivateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create deployer transactor opts: %w", err)
	}

	// Setup the payload disperser
	params := deployment.PayloadDisperserParams{
		Logger:                 logger,
		EthClient:              ethClient,
		DisperserPrivateKey:    config.DisperserPrivateKey,
		DisperserHostname:      config.DisperserHostname,
		DisperserPort:          config.DisperserPort,
		DisperseBlobTimeout:    config.DisperseBlobTimeout,
		BlobCompleteTimeout:    config.BlobCompleteTimeout,
		BlobStatusPollInterval: config.BlobStatusPollInterval,
		ContractCallTimeout:    config.ContractCallTimeout,
		CertBuilder:            im.result.CertVerification.CertBuilder,
		RouterCertVerifier:     im.result.CertVerification.RouterCertVerifier,
	}
	payloadDisperser, err := deployment.SetupPayloadDisperser(params)
	if err != nil {
		return fmt.Errorf("failed to setup payload disperser: %w", err)
	}

	// Store components in the result
	im.result.PayloadDisperser = &PayloadDisperserComponents{
		PayloadDisperser:       payloadDisperser,
		DeployerTransactorOpts: deployerTransactorOpts,
	}

	fmt.Println("✅ Payload disperser initialized successfully")
	return nil
}
