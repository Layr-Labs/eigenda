package testinfra

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

// InfraManager orchestrates the lifecycle of test infrastructure containers
type InfraManager struct {
	config     InfraConfig
	anvil      *containers.AnvilContainer
	localstack *containers.LocalStackContainer
	graphnode  *containers.GraphNodeContainer
	churner    *containers.ChurnerContainer
	dispersers []*containers.DisperserContainer
	encoders   []*containers.EncoderContainer
	batcher    *containers.BatcherContainer
	controller *containers.ControllerContainer
	operators  []*containers.OperatorContainer
	relays     []*containers.RelayContainer
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

		// Test subgraph connectivity
		subgraphURL := im.result.GraphNodeURL + "/subgraphs/name/Layr-Labs/eigenda-operator-state"
		fmt.Println("Testing deployed subgraph connectivity...")
		graphConnConfig := DefaultGraphConnectivityConfig()
		err = TestGraphNodeConnectivity(subgraphURL, graphConnConfig)
		if err != nil {
			fmt.Printf("📋 Graph Node Debug Info:\n")
			fmt.Printf("   - GraphQL URL: %s\n", subgraphURL)
			fmt.Printf("   - Admin URL: %s\n", im.result.GraphNodeAdminURL)
			fmt.Printf("   - IPFS URL: %s\n", im.result.IPFSURL)

			return nil, fmt.Errorf("subgraph connectivity test failed: %w", err)
		}
		fmt.Println("✅ Subgraph connectivity test passed")
	}

	// 6. Setup retrieval clients if enabled and EigenDA contracts are deployed
	if im.config.EigenDA.RetrievalClients.Enabled && im.result.EigenDAContracts != nil {
		fmt.Println("Setting up retrieval clients...")
		err := im.setupRetrievalClients(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to setup retrieval clients: %w", err)
		}
	}

	// 7. Start churner if enabled and EigenDA contracts are deployed
	if im.config.EigenDA.Churner.Enabled && im.result.EigenDAContracts != nil {
		fmt.Println("Starting churner service...")

		// Configure churner with deployed contract addresses
		churnerConfig := im.config.EigenDA.Churner.ChurnerConfig
		churnerConfig.EigenDADirectory = im.config.EigenDA.RootPath
		churnerConfig.OperatorStateRetriever = im.result.EigenDAContracts.OperatorStateRetriever
		churnerConfig.ServiceManager = im.result.EigenDAContracts.ServiceManager

		// Use internal Anvil URL for containers on the same network
		if im.network != nil && im.anvil != nil {
			churnerConfig.ChainRPC = "http://anvil:8545"
		} else {
			churnerConfig.ChainRPC = im.result.AnvilRPC
		}

		// Use deployer private key for churner
		deployerPrivateKey := im.config.EigenDA.Deployer.PrivateKey
		if deployerPrivateKey == "" {
			deployerPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		}
		churnerConfig.PrivateKey = deployerPrivateKey

		// Set Graph URL if Graph Node is enabled
		if im.config.GraphNode.Enabled && im.graphnode != nil {
			if im.network != nil {
				// Use internal URL for containers on the same network
				churnerConfig.GraphURL = "http://graph-node:8000/subgraphs/name/Layr-Labs/eigenda-operator-state"
			} else {
				churnerConfig.GraphURL = im.result.GraphNodeURL + "/subgraphs/name/Layr-Labs/eigenda-operator-state"
			}
		}

		// Start the churner container
		churner, err := containers.NewChurnerContainerWithNetwork(ctx, churnerConfig, im.network)
		if err != nil {
			return nil, fmt.Errorf("failed to start churner: %w", err)
		}
		im.churner = churner
		im.result.ChurnerURL = churner.URL()

		// Use internal URL for other containers
		if im.network != nil {
			im.result.ChurnerInternalURL = churner.InternalURL()
		}

		fmt.Printf("✅ Churner service started at %s\n", im.result.ChurnerURL)
	}

	// 8. Start operators if enabled and EigenDA contracts are deployed
	if im.config.EigenDA.Operators.Enabled && im.result.EigenDAContracts != nil {
		fmt.Printf("Starting %d operators...\n", im.config.EigenDA.Operators.Count)

		// Initialize operator addresses map
		im.result.OperatorAddresses = make(map[int]string)

		// Determine how many operators to actually run
		numToStart := im.config.EigenDA.Operators.Count

		// Start each operator
		for i := 0; i < numToStart; i++ {
			// Create operator config based on base config
			operatorConfig := containers.DefaultOperatorConfig(i)

			// Override with base config values if provided
			if im.config.EigenDA.Operators.BaseConfig.Image != "" {
				operatorConfig.Image = im.config.EigenDA.Operators.BaseConfig.Image
			}

			// Set contract addresses
			operatorConfig.BLSOperatorStateRetriever = im.result.EigenDAContracts.OperatorStateRetriever
			operatorConfig.EigenDAServiceManager = im.result.EigenDAContracts.ServiceManager

			// Debug log the contract addresses
			fmt.Printf("DEBUG: Setting operator %d contracts - BLSOperatorStateRetriever: %s, ServiceManager: %s\n",
				i, operatorConfig.BLSOperatorStateRetriever, operatorConfig.EigenDAServiceManager)

			// Use internal Anvil URL for containers on the same network
			if im.network != nil && im.anvil != nil {
				operatorConfig.ChainRPC = "http://anvil:8545"
			} else {
				operatorConfig.ChainRPC = im.result.AnvilRPC
			}

			// Set operator private keys from config
			operatorKey := fmt.Sprintf("opr%d", i)
			if privKey, ok := im.config.EigenDA.PrivateKeys[operatorKey]; ok {
				operatorConfig.PrivateKey = privKey
			}

			// Set BLS and ECDSA key file configuration
			// Use operator index + 1 to match the numbering in the key files (1-based)
			keyIndex := i + 1
			operatorConfig.BlsKeyFile = fmt.Sprintf("/app/secrets/bls_keys/keys/%d.bls.key.json", keyIndex)
			operatorConfig.EcdsaKeyFile = fmt.Sprintf("/app/secrets/ecdsa_keys/keys/%d.ecdsa.key.json", keyIndex)

			// Read BLS password from the password file
			blsPasswordFile := filepath.Join(im.config.EigenDA.RootPath, "common", "testinfra", "secrets", "bls_keys", "password.txt")
			blsPasswords, err := readPasswordFile(blsPasswordFile)
			if err != nil {
				fmt.Printf("Warning: Could not read BLS password file: %v\n", err)
				operatorConfig.BlsKeyFile = ""
			} else if keyIndex <= len(blsPasswords) {
				operatorConfig.BlsKeyPassword = blsPasswords[keyIndex-1]
				fmt.Printf("DEBUG: Loaded BLS password for operator %d\n", i)
			} else {
				fmt.Printf("Warning: No BLS password for operator %d\n", i)
				operatorConfig.BlsKeyFile = ""
			}

			// Read ECDSA password from the password file
			ecdsaPasswordFile := filepath.Join(im.config.EigenDA.RootPath, "common", "testinfra", "secrets", "ecdsa_keys", "password.txt")
			ecdsaPasswords, err := readPasswordFile(ecdsaPasswordFile)
			if err != nil {
				// Fallback to using private key directly if password file is not available
				fmt.Printf("Warning: Could not read ECDSA password file, falling back to private key: %v\n", err)
				operatorConfig.EcdsaPrivateKey = operatorConfig.PrivateKey
				operatorConfig.EcdsaKeyFile = ""
			} else if keyIndex <= len(ecdsaPasswords) {
				operatorConfig.EcdsaKeyPassword = ecdsaPasswords[keyIndex-1]

				// In test mode, we also need to provide the raw private key
				// Read the ECDSA private key from the private_key_hex.txt file
				ecdsaPrivKeyFile := filepath.Join(im.config.EigenDA.RootPath, "common", "testinfra", "secrets", "ecdsa_keys", "private_key_hex.txt")
				ecdsaKeys, err := readPasswordFile(ecdsaPrivKeyFile)
				if err == nil && keyIndex <= len(ecdsaKeys) {
					operatorConfig.EcdsaPrivateKey = ecdsaKeys[keyIndex-1]
					fmt.Printf("DEBUG: Loaded ECDSA private key for operator %d\n", i)
				} else {
					fmt.Printf("Warning: Could not load ECDSA private key for operator %d: %v\n", i, err)
				}
			} else {
				// Fallback if key index is beyond available passwords
				fmt.Printf("Warning: No ECDSA password for operator %d, falling back to private key\n", i)
				operatorConfig.EcdsaPrivateKey = operatorConfig.PrivateKey
				operatorConfig.EcdsaKeyFile = ""
			}

			// Set EigenDA directory for resource paths
			operatorConfig.EigenDADirectory = im.config.EigenDA.RootPath
			fmt.Printf("DEBUG: EigenDA root path for operator %d: %s\n", i, operatorConfig.EigenDADirectory)

			// Set the hostname for operator registration
			operatorConfig.Hostname = fmt.Sprintf("operator-%d.localtest.me", i)

			// Start the operator container
			operator, err := containers.NewOperatorContainerWithNetwork(ctx, operatorConfig, im.network)
			if err != nil {
				return nil, fmt.Errorf("failed to start operator %d: %w", i, err)
			}
			im.operators = append(im.operators, operator)

			// Store the operator's internal address for other services to use
			im.result.OperatorAddresses[i] = operator.GetInternalDispersalAddress()

			fmt.Printf("  ✅ Operator %d started: %s\n", i, operator.GetInternalDispersalAddress())
		}

		fmt.Printf("✅ %d operators started successfully\n", numToStart)
	}

	// 9. Start encoders if configured (needed by batcher)
	if len(im.config.EigenDA.Encoders) > 0 {
		fmt.Println("Starting encoder services...")

		// Initialize encoder URL maps
		im.result.EncoderURLs = make(map[string]string)
		im.result.EncoderInternalURLs = make(map[string]string)

		for _, encoderCfg := range im.config.EigenDA.Encoders {
			if !encoderCfg.Enabled {
				continue
			}

			// Configure encoder
			encoderConfig := encoderCfg.EncoderConfig

			// Configure AWS resources (required even for encoder v1)
			if im.config.LocalStack.Enabled && im.localstack != nil {
				// Use internal endpoint when containers are on the same network
				if im.network != nil {
					encoderConfig.AWSEndpointURL = im.localstack.InternalEndpoint()
				} else {
					encoderConfig.AWSEndpointURL = im.localstack.Endpoint()
				}
				encoderConfig.AWSRegion = im.localstack.Region()
				encoderConfig.AWSAccessKeyID = "localstack"
				encoderConfig.AWSSecretAccessKey = "localstack"

				// Use configured bucket name for encoder v2
				if encoderConfig.EncoderVersion == "2" && im.config.LocalStack.Resources.BucketName != "" {
					encoderConfig.S3BucketName = im.config.LocalStack.Resources.BucketName
				}
			} else {
				// Use default AWS region if LocalStack is not enabled
				if encoderConfig.AWSRegion == "" {
					encoderConfig.AWSRegion = "us-east-1"
				}
			}

			// Set KZG paths if not already configured - must be absolute paths for Docker
			// Use the KZG resources from root resources/srs directory
			if encoderConfig.G1Path == "" {
				// Use root resources/srs KZG files
				g1Path := filepath.Join(im.config.EigenDA.RootPath, "resources/srs/g1.point")
				absG1Path, err := filepath.Abs(g1Path)
				if err != nil {
					return nil, fmt.Errorf("failed to get absolute path for G1: %w", err)
				}
				encoderConfig.G1Path = absG1Path
			}
			if encoderConfig.G2Path == "" {
				// Use root resources/srs KZG files
				g2Path := filepath.Join(im.config.EigenDA.RootPath, "resources/srs/g2.point")
				absG2Path, err := filepath.Abs(g2Path)
				if err != nil {
					return nil, fmt.Errorf("failed to get absolute path for G2: %w", err)
				}
				encoderConfig.G2Path = absG2Path
			}
			if encoderConfig.G2PowerOf2Path == "" {
				// Use root resources/srs KZG files
				g2PowerOf2Path := filepath.Join(im.config.EigenDA.RootPath, "resources/srs/g2.point.powerOf2")
				absG2PowerOf2Path, err := filepath.Abs(g2PowerOf2Path)
				if err != nil {
					return nil, fmt.Errorf("failed to get absolute path for G2PowerOf2: %w", err)
				}
				encoderConfig.G2PowerOf2Path = absG2PowerOf2Path
			}
			if encoderConfig.CachePath == "" {
				// Use root resources/srs cache directory
				cachePath := filepath.Join(im.config.EigenDA.RootPath, "resources/srs/SRSTables")
				absCachePath, err := filepath.Abs(cachePath)
				if err != nil {
					return nil, fmt.Errorf("failed to get absolute path for cache: %w", err)
				}
				encoderConfig.CachePath = absCachePath
			}

			// Start the encoder container
			encoder, err := containers.NewEncoderContainerWithNetwork(ctx, encoderConfig, im.network)
			if err != nil {
				return nil, fmt.Errorf("failed to start encoder v%s: %w", encoderConfig.EncoderVersion, err)
			}
			im.encoders = append(im.encoders, encoder)

			// Store URLs by encoder version
			version := encoderConfig.EncoderVersion
			if version == "" {
				version = "1" // Default to v1
			}
			im.result.EncoderURLs[version] = encoder.URL()

			// Use internal URL for other containers
			if im.network != nil {
				im.result.EncoderInternalURLs[version] = encoder.InternalURL()
			}

			fmt.Printf("✅ Encoder v%s service started at %s\n", version, encoder.URL())
		}
	}

	// 10. Start dispersers if enabled and EigenDA contracts are deployed
	if len(im.config.EigenDA.Dispersers) > 0 && im.result.EigenDAContracts != nil {
		fmt.Printf("Starting %d dispersers...\n", len(im.config.EigenDA.Dispersers))

		// Initialize disperser URL maps
		im.result.DisperserURLs = make(map[string]string)
		im.result.DisperserInternalURLs = make(map[string]string)

		for _, disperserCfg := range im.config.EigenDA.Dispersers {
			if !disperserCfg.Enabled {
				continue
			}

			// Configure disperser with deployed contract addresses
			disperserConfig := disperserCfg.DisperserConfig
			disperserConfig.EigenDADirectory = im.config.EigenDA.RootPath
			disperserConfig.OperatorStateRetriever = im.result.EigenDAContracts.OperatorStateRetriever
			disperserConfig.ServiceManager = im.result.EigenDAContracts.ServiceManager

			// Use internal URLs for containers on the same network
			if im.network != nil {
				if im.anvil != nil {
					disperserConfig.ChainRPC = "http://anvil:8545"
				}
				if im.localstack != nil {
					disperserConfig.AWSEndpointURL = "http://localstack:4566"
				}
			} else {
				disperserConfig.ChainRPC = im.result.AnvilRPC
				if im.localstack != nil {
					disperserConfig.AWSEndpointURL = im.result.LocalStackURL
				}
			}

			// Set encoder address if encoders are available
			if len(im.encoders) > 0 {
				// For v1 disperser, use the first encoder
				// For v2 disperser, use the second encoder if available, otherwise first
				encoderIdx := 0
				if disperserConfig.Version == 2 && len(im.encoders) > 1 {
					encoderIdx = 1
				}
				if im.network != nil {
					// Use internal URL for containers on the same network
					disperserConfig.EncoderAddress = fmt.Sprintf("encoder-%d:%s", encoderIdx, im.encoders[encoderIdx].Config().GRPCPort)
				} else {
					disperserConfig.EncoderAddress = im.encoders[encoderIdx].URL()
				}
			}

			// Set KZG paths for v2 disperser (required)
			if disperserConfig.Version == 2 {
				// Use the same KZG paths as encoders if not already configured
				if disperserConfig.G1Path == "" && len(im.encoders) > 0 {
					// Copy KZG configuration from the encoder (they should have the same paths)
					encoderIdx := 0
					if len(im.encoders) > 1 {
						encoderIdx = 1 // Use v2 encoder's config if available
					}
					encoderConfig := im.encoders[encoderIdx].Config()
					disperserConfig.G1Path = encoderConfig.G1Path
					disperserConfig.G2Path = encoderConfig.G2Path
					disperserConfig.G2PowerOf2Path = encoderConfig.G2PowerOf2Path
					disperserConfig.CachePath = encoderConfig.CachePath
					disperserConfig.SRSOrder = encoderConfig.SRSOrder
					disperserConfig.SRSLoad = encoderConfig.SRSLoad
				}
			}

			// Start the disperser container
			disperser, err := containers.NewDisperserContainerWithNetwork(ctx, disperserConfig, im.network)
			if err != nil {
				return nil, fmt.Errorf("failed to start disperser v%d: %w", disperserConfig.Version, err)
			}
			im.dispersers = append(im.dispersers, disperser)

			// Store URLs by version
			versionStr := fmt.Sprintf("%d", disperserConfig.Version)
			im.result.DisperserURLs[versionStr] = disperser.URL()

			// Use internal URL for other containers
			if im.network != nil {
				im.result.DisperserInternalURLs[versionStr] = disperser.InternalURL()
			}

			fmt.Printf("✅ Disperser v%d service started at %s\n", disperserConfig.Version, disperser.URL())
		}
	}

	// 11. Start batcher if enabled and required components are available
	if im.config.EigenDA.Batcher.Enabled && im.result.EigenDAContracts != nil {
		fmt.Println("Starting batcher service...")

		// Configure batcher with deployed contract addresses and AWS resources
		batcherConfig := im.config.EigenDA.Batcher.BatcherConfig
		batcherConfig.EigenDADirectory = im.config.EigenDA.RootPath
		batcherConfig.OperatorStateRetriever = im.result.EigenDAContracts.OperatorStateRetriever
		batcherConfig.ServiceManager = im.result.EigenDAContracts.ServiceManager

		// Configure AWS resources if LocalStack is enabled
		if im.config.LocalStack.Enabled && im.localstack != nil {
			// Use internal endpoint when containers are on the same network
			if im.network != nil {
				batcherConfig.AWSEndpointURL = im.localstack.InternalEndpoint()
			} else {
				batcherConfig.AWSEndpointURL = im.localstack.Endpoint()
			}
			batcherConfig.AWSRegion = im.localstack.Region()
			batcherConfig.AWSAccessKeyID = "localstack"
			batcherConfig.AWSSecretAccessKey = "localstack"

			// Use configured bucket and table names
			if im.config.LocalStack.Resources.BucketName != "" {
				batcherConfig.S3BucketName = im.config.LocalStack.Resources.BucketName
			}
			if im.config.LocalStack.Resources.MetadataTableName != "" {
				batcherConfig.DynamoDBTableName = im.config.LocalStack.Resources.MetadataTableName
			}
		}

		// Use internal Anvil URL for containers on the same network
		if im.network != nil && im.anvil != nil {
			batcherConfig.ChainRPC = "http://anvil:8545"
		} else {
			batcherConfig.ChainRPC = im.result.AnvilRPC
		}

		// Use batcher0 private key
		batcherPrivateKey := im.config.EigenDA.PrivateKeys["batcher0"]
		if batcherPrivateKey == "" {
			// Use a default test private key if not provided
			batcherPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		}
		batcherConfig.PrivateKey = batcherPrivateKey

		// Set Graph URL if Graph Node is enabled
		if im.config.GraphNode.Enabled && im.graphnode != nil {
			if im.network != nil {
				// Use internal URL for containers on the same network
				batcherConfig.GraphURL = "http://graph-node:8000/subgraphs/name/Layr-Labs/eigenda-operator-state"
			} else {
				batcherConfig.GraphURL = im.result.GraphNodeURL + "/subgraphs/name/Layr-Labs/eigenda-operator-state"
			}
		}

		// Set encoder address if encoders are configured (batcher uses v1)
		if len(im.encoders) > 0 && im.result.EncoderURLs != nil {
			// Batcher needs encoder v1
			if im.network != nil {
				// Use internal URL for containers on the same network
				if url, ok := im.result.EncoderInternalURLs["1"]; ok {
					batcherConfig.EncoderAddress = url
				}
			} else {
				if url, ok := im.result.EncoderURLs["1"]; ok {
					batcherConfig.EncoderAddress = url
				}
			}
		}

		// Start the batcher container
		batcher, err := containers.NewBatcherContainerWithNetwork(ctx, batcherConfig, im.network)
		if err != nil {
			return nil, fmt.Errorf("failed to start batcher: %w", err)
		}
		im.batcher = batcher
		im.result.BatcherURL = batcher.URL()

		// Use internal URL for other containers
		if im.network != nil {
			im.result.BatcherInternalURL = batcher.InternalURL()
		}

		fmt.Printf("✅ Batcher service started (metrics at %s)\n", im.result.BatcherURL)
	}

	// 12. Start controller if enabled and encoders are available
	if im.config.EigenDA.Controller.Enabled && len(im.encoders) > 0 && im.result.EigenDAContracts != nil {
		fmt.Println("Starting controller service...")

		// Configure controller
		controllerConfig := im.config.EigenDA.Controller.ControllerConfig

		// Set encoder address (controller uses v2 encoder)
		if im.network != nil {
			if url, ok := im.result.EncoderInternalURLs["2"]; ok {
				controllerConfig.EncoderAddress = url
			}
		} else {
			if url, ok := im.result.EncoderURLs["2"]; ok {
				controllerConfig.EncoderAddress = url
			}
		}

		// Configure AWS resources
		if im.config.LocalStack.Enabled && im.localstack != nil {
			if im.network != nil {
				controllerConfig.AWSEndpointURL = im.localstack.InternalEndpoint()
			} else {
				controllerConfig.AWSEndpointURL = im.localstack.Endpoint()
			}
			controllerConfig.AWSRegion = im.localstack.Region()
			controllerConfig.AWSAccessKeyID = "localstack"
			controllerConfig.AWSSecretAccessKey = "localstack"
		}

		// Configure chain settings
		if im.anvil != nil {
			if im.network != nil {
				internalURL, err := im.anvil.InternalRPCURL(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to get anvil internal RPC URL: %w", err)
				}
				controllerConfig.ChainRPC = internalURL
			} else {
				controllerConfig.ChainRPC = im.anvil.RPCURL()
			}
		}

		// Set contract addresses from deployment
		controllerConfig.EigenDAServiceManager = im.result.EigenDAContracts.ServiceManager
		controllerConfig.BLSOperatorStateRetriever = im.result.EigenDAContracts.OperatorStateRetriever

		// Configure graph URL if available
		if im.config.GraphNode.Enabled && im.graphnode != nil {
			controllerConfig.UseGraph = true
			// GraphNode uses internal network name "graphnode" on port 8000
			if im.network != nil {
				controllerConfig.GraphURL = "http://graph-node:8000/subgraphs/name/Layr-Labs/eigenda-operator-state"
			} else {
				controllerConfig.GraphURL = im.result.GraphNodeURL + "/subgraphs/name/Layr-Labs/eigenda-operator-state"
			}
		}

		// Set disperser KMS key ID if available
		if im.result.DisperserKMSKeyID != "" {
			controllerConfig.DisperserKMSKeyID = im.result.DisperserKMSKeyID
		}

		// Start the controller container
		controller, err := containers.NewControllerContainerWithNetwork(ctx, controllerConfig, im.network)
		if err != nil {
			return nil, fmt.Errorf("failed to start controller: %w", err)
		}
		im.controller = controller
		im.result.ControllerMetricsURL = controller.MetricsURL()

		// Use internal URL for other containers
		if im.network != nil {
			im.result.ControllerInternalMetricsURL = controller.InternalMetricsURL()
		}

		fmt.Printf("✅ Controller service started (metrics at %s)\n", im.result.ControllerMetricsURL)
	}

	// 13. Start relays if enabled
	if im.config.EigenDA.Relays.Enabled && im.result.EigenDAContracts != nil {
		fmt.Printf("Starting %d relay services...\n", im.config.EigenDA.Relays.Count)

		// Initialize relay URL maps
		im.result.RelayURLs = make(map[int]string)
		im.result.RelayInternalURLs = make(map[int]string)

		for i := 0; i < im.config.EigenDA.Relays.Count; i++ {
			// Configure relay with base configuration
			relayConfig := im.config.EigenDA.Relays.BaseConfig
			relayConfig.ID = i

			// Override with default configuration if not set
			if relayConfig.Image == "" {
				relayConfig = containers.DefaultRelayConfig(i)
			} else {
				// Merge with defaults to ensure all required fields are set
				defaultConfig := containers.DefaultRelayConfig(i)
				relayConfig.ID = defaultConfig.ID
				if relayConfig.GRPCPort == "" {
					relayConfig.GRPCPort = defaultConfig.GRPCPort
				}
				if relayConfig.InternalGRPCPort == "" {
					relayConfig.InternalGRPCPort = defaultConfig.InternalGRPCPort
				}
				if relayConfig.MetricsPort == "" {
					relayConfig.MetricsPort = defaultConfig.MetricsPort
				}
				if relayConfig.Hostname == "" {
					relayConfig.Hostname = defaultConfig.Hostname
				}
			}

			// Set the hostname for relay registration
			relayConfig.Hostname = fmt.Sprintf("relay-%d.localtest.me", i)

			// Configure AWS resources if LocalStack is enabled
			if im.config.LocalStack.Enabled && im.localstack != nil {
				// Use internal endpoint when containers are on the same network
				if im.network != nil {
					relayConfig.AWSEndpointURL = im.localstack.InternalEndpoint()
				} else {
					relayConfig.AWSEndpointURL = im.localstack.Endpoint()
				}
				relayConfig.AWSRegion = im.localstack.Region()
				relayConfig.AWSAccessKeyID = "localstack"
				relayConfig.AWSSecretAccessKey = "localstack"

				// Use configured bucket and table names
				if im.config.LocalStack.Resources.BucketName != "" {
					relayConfig.BucketName = im.config.LocalStack.Resources.BucketName
				}
				if im.config.LocalStack.Resources.V2MetadataTableName != "" {
					relayConfig.MetadataTableName = im.config.LocalStack.Resources.V2MetadataTableName
				}
			}

			// Configure chain settings
			if im.anvil != nil {
				if im.network != nil {
					relayConfig.ChainRPC = "http://anvil:8545"
				} else {
					relayConfig.ChainRPC = im.anvil.RPCURL()
				}
			}

			// Configure graph URL if available
			if im.config.GraphNode.Enabled && im.graphnode != nil {
				// GraphNode uses internal network name "graphnode" on port 8000
				if im.network != nil {
					relayConfig.GraphURL = "http://graph-node:8000/subgraphs/name/Layr-Labs/eigenda-operator-state"
				} else {
					relayConfig.GraphURL = im.result.GraphNodeURL + "/subgraphs/name/Layr-Labs/eigenda-operator-state"
				}
			}

			// Set contract addresses from deployment
			relayConfig.BLSOperatorStateRetriever = im.result.EigenDAContracts.OperatorStateRetriever
			relayConfig.EigenDAServiceManager = im.result.EigenDAContracts.ServiceManager
			if im.result.EigenDAContracts.EigenDADirectory != "" {
				relayConfig.EigenDADirectory = im.result.EigenDAContracts.EigenDADirectory
			}

			// Set the relay keys - by default each relay serves its own index
			relayConfig.RelayKeys = []int{i}

			// Start the relay container
			relay, err := containers.NewRelayContainerWithNetwork(ctx, relayConfig, im.network)
			if err != nil {
				return nil, fmt.Errorf("failed to start relay %d: %w", i, err)
			}
			im.relays = append(im.relays, relay)

			// Store the relay's addresses
			im.result.RelayURLs[i] = relay.GetGRPCAddress()
			if im.network != nil {
				im.result.RelayInternalURLs[i] = relay.GetInternalGRPCAddress()
			}

			fmt.Printf("  ✅ Relay %d started: %s\n", i, relay.GetGRPCAddress())
		}

		fmt.Printf("✅ %d relays started successfully\n", im.config.EigenDA.Relays.Count)
	}

	// 14. Setup payload disperser if enabled and required components are available
	if im.config.EigenDA.PayloadDisperser.Enabled &&
		im.result.CertVerification != nil &&
		im.result.RetrievalClients != nil &&
		len(im.dispersers) > 0 {

		// Use the v2 disperser URL if available, otherwise v1
		disperserURL := ""
		if v2URL, ok := im.result.DisperserURLs["2"]; ok {
			disperserURL = v2URL
		} else if v1URL, ok := im.result.DisperserURLs["1"]; ok {
			disperserURL = v1URL
		}

		if disperserURL != "" {
			// Parse the URL to get hostname and port
			parts := strings.Split(disperserURL, ":")
			if len(parts) == 2 {
				// Override the config with actual disperser URL
				im.config.EigenDA.PayloadDisperser.DisperserHostname = parts[0]
				im.config.EigenDA.PayloadDisperser.DisperserPort = parts[1]
			}

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
	if im.batcher != nil {
		// Print log path for debugging
		if logPath := im.batcher.LogPath(); logPath != "" {
			fmt.Printf("Batcher logs available at: %s\n", logPath)
		}
		if err := im.batcher.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate batcher: %w", err))
		}
		im.batcher = nil
	}

	// Terminate controller
	if im.controller != nil {
		// Print log path for debugging
		if logPath := im.controller.LogPath(); logPath != "" {
			fmt.Printf("Controller logs available at: %s\n", logPath)
		}
		if err := im.controller.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate controller: %w", err))
		}
		im.controller = nil
	}

	// Terminate dispersers
	for _, disperser := range im.dispersers {
		if disperser != nil {
			// Print log path for debugging
			if logPath := disperser.LogPath(); logPath != "" {
				fmt.Printf("Disperser v%d logs available at: %s\n", disperser.Version(), logPath)
			}
			if err := disperser.Terminate(ctx); err != nil {
				errs = append(errs, fmt.Errorf("failed to terminate disperser v%d: %w", disperser.Version(), err))
			}
		}
	}
	im.dispersers = nil

	// Terminate encoders
	for _, encoder := range im.encoders {
		if encoder != nil {
			// Print log path for debugging
			if logPath := encoder.LogPath(); logPath != "" {
				fmt.Printf("Encoder logs available at: %s\n", logPath)
			}
			if err := encoder.Terminate(ctx); err != nil {
				errs = append(errs, fmt.Errorf("failed to terminate encoder: %w", err))
			}
		}
	}
	im.encoders = nil

	// Terminate operators
	for i, operator := range im.operators {
		if operator != nil {
			// Print log path for debugging
			if logPath := operator.LogPath(); logPath != "" {
				fmt.Printf("Operator %d logs available at: %s\n", i, logPath)
			}
			if err := operator.Terminate(ctx); err != nil {
				errs = append(errs, fmt.Errorf("failed to terminate operator %d: %w", i, err))
			}
		}
	}
	im.operators = nil

	// Terminate relays
	for i, relay := range im.relays {
		if relay != nil {
			// Print log path for debugging
			if logPath := relay.LogPath(); logPath != "" {
				fmt.Printf("Relay %d logs available at: %s\n", i, logPath)
			}
			if err := relay.Terminate(ctx); err != nil {
				errs = append(errs, fmt.Errorf("failed to terminate relay %d: %w", i, err))
			}
		}
	}
	im.relays = nil

	if im.churner != nil {
		// Print log path for debugging
		if logPath := im.churner.LogPath(); logPath != "" {
			fmt.Printf("Churner logs available at: %s\n", logPath)
		}
		if err := im.churner.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate churner: %w", err))
		}
		im.churner = nil
	}

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

// readPasswordFile reads a password file and returns a slice of passwords (one per line)
func readPasswordFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open password file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var passwords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		password := strings.TrimSpace(scanner.Text())
		if password != "" {
			passwords = append(passwords, password)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan password file: %w", err)
	}

	return passwords, nil
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
		// Create deployment config matching the test configuration structure
		numStrategies := 2    // Two strategies as per the stakes configuration
		maxOperatorCount := 3 // maxOperatorCount from services.counts

		// Stakes configuration from testconfig-anvil.yaml:
		// - total: 100e18, distribution: [1, 4, 6, 10]
		// - total: 100e18, distribution: [1, 3, 8, 9]
		stakeDistribution := [][]float32{
			{1, 4, 6, 10}, // Strategy 0 distribution
			{1, 3, 8, 9},  // Strategy 1 distribution
		}
		stakeTotals := []float32{100e18, 100e18} // 100e18 tokens per strategy

		// Generate test private keys if none provided
		privateKeys := im.config.EigenDA.PrivateKeys
		if len(privateKeys) == 0 {
			fmt.Println("No private keys provided, generating test keys for local testing...")
			privateKeys = generateTestPrivateKeys()
		}

		// Load operator ECDSA private keys from the secrets directory
		// These are the actual keys the operators will use for transactions
		ecdsaPrivKeyFile := filepath.Join(im.config.EigenDA.RootPath, "common", "testinfra", "secrets", "ecdsa_keys", "private_key_hex.txt")
		ecdsaKeys, err := readPasswordFile(ecdsaPrivKeyFile)
		if err == nil {
			// Add ECDSA keys to the privateKeys map with "opr{i}_ecdsa" naming
			for i, ecdsaKey := range ecdsaKeys {
				keyName := fmt.Sprintf("opr%d_ecdsa", i)
				privateKeys[keyName] = ecdsaKey
				fmt.Printf("Added operator %d ECDSA key for funding\n", i)
			}
		} else {
			fmt.Printf("Warning: Could not load operator ECDSA keys for funding, operators may not have ETH for gas: %v\n", err)
		}

		deployConfig := deployment.GenerateEigenDADeployConfig(
			numStrategies,
			maxOperatorCount,
			stakeDistribution,
			stakeTotals,
			privateKeys,
		)

		err = manager.DeployEigenDAContracts(deployConfig)
		if err != nil {
			return fmt.Errorf("failed to deploy EigenDA contracts: %w", err)
		}

		// The deployment manager populates its own EigenDAContracts field
		// Copy it to our result
		im.result.EigenDAContracts = manager.EigenDAContracts

		// Debug log the loaded contracts
		if im.result.EigenDAContracts != nil {
			fmt.Printf("DEBUG: Loaded EigenDA contracts - OperatorStateRetriever: %s, ServiceManager: %s\n",
				im.result.EigenDAContracts.OperatorStateRetriever,
				im.result.EigenDAContracts.ServiceManager)
		} else {
			fmt.Printf("WARNING: EigenDAContracts is nil after deployment!\n")
		}
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

		// Collect relay URLs from running relays if they exist
		var relayURLs []string
		if len(im.relays) > 0 {
			// Use actual relay URLs from running relays with proper hostnames
			for i, relay := range im.relays {
				// Use relay-N.localtest.me:PORT format for Docker network resolution
				url := fmt.Sprintf("relay-%d.localtest.me:%s", i, relay.Config().GRPCPort)
				relayURLs = append(relayURLs, url)
			}
		} else if im.config.EigenDA.Relays.Enabled {
			// If relays will be started later, use their expected URLs
			for i := 0; i < im.config.EigenDA.Relays.Count; i++ {
				basePort := 34000 + (i * 2) // Match DefaultRelayConfig port calculation
				url := fmt.Sprintf("relay-%d.localtest.me:%d", i, basePort)
				relayURLs = append(relayURLs, url)
			}
		} else {
			// Fallback to standard relay URLs used in tests
			relayURLs = []string{"localhost:32035", "localhost:32037", "localhost:32039", "localhost:32041"}
		}

		err = manager.RegisterBlobVersionsAndRelays(ethClient, im.config.EigenDA.BlobVersionParams, relayURLs)
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

	// No need to change directories - we'll use absolute paths to testinfra's KZG resources

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

	// Use default KZG paths from root resources/srs (official mainnet) if not specified
	g1Path := config.G1Path
	if g1Path == "" {
		g1Path = filepath.Join(im.config.EigenDA.RootPath, "resources/srs/g1.point")
	}
	g2Path := config.G2Path
	if g2Path == "" {
		g2Path = filepath.Join(im.config.EigenDA.RootPath, "resources/srs/g2.point")
	}
	g2PowerOf2Path := config.G2PowerOf2Path
	if g2PowerOf2Path == "" {
		g2PowerOf2Path = filepath.Join(im.config.EigenDA.RootPath, "resources/srs/g2.point.powerOf2")
	}
	cachePath := config.CachePath
	if cachePath == "" {
		cachePath = filepath.Join(im.config.EigenDA.RootPath, "resources/srs/cache")
	}

	// Get absolute paths
	g1Path, _ = filepath.Abs(g1Path)
	g2Path, _ = filepath.Abs(g2Path)
	g2PowerOf2Path, _ = filepath.Abs(g2PowerOf2Path)
	cachePath, _ = filepath.Abs(cachePath)

	// Print absolute paths, for debugging purposes
	fmt.Println("Using KZG resources from:")
	fmt.Printf(" - G1: %s\n", g1Path)
	fmt.Printf(" - G2: %s\n", g2Path)
	fmt.Printf(" - G2 (Power of 2): %s\n", g2PowerOf2Path)
	fmt.Printf(" - Cache: %s\n", cachePath)

	kzgConfig := &kzg.KzgConfig{
		G1Path:          g1Path,
		G2Path:          g2Path,
		G2PowerOf2Path:  g2PowerOf2Path,
		CacheDir:        cachePath,
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

// GraphConnectivityConfig contains configuration for testing Graph Node connectivity
type GraphConnectivityConfig struct {
	MaxRetries    int
	RetryInterval time.Duration
	Timeout       time.Duration
}

// DefaultGraphConnectivityConfig returns sensible default settings for graph connectivity testing
func DefaultGraphConnectivityConfig() GraphConnectivityConfig {
	return GraphConnectivityConfig{
		MaxRetries:    10,
		RetryInterval: 2 * time.Second,
		Timeout:       10 * time.Second,
	}
}

// TestGraphNodeConnectivity tests the connectivity to a Graph Node instance
// It queries the subgraph's _meta to verify it's accessible and synced
func TestGraphNodeConnectivity(graphURL string, config GraphConnectivityConfig) error {
	client := &http.Client{Timeout: config.Timeout}

	// Test the subgraph's actual data availability
	// This query checks if the subgraph is synced and has block data
	query := `{"query": "{_meta{block{number}}}"}`

	for i := 0; i < config.MaxRetries; i++ {
		fmt.Printf("Testing subgraph connectivity (attempt %d/%d)...\n", i+1, config.MaxRetries)

		req, err := http.NewRequest("POST", graphURL, strings.NewReader(query))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("  ❌ Connection failed: %v\n", err)
			if i < config.MaxRetries-1 {
				fmt.Printf("  ⏱️  Waiting %v before retry...\n", config.RetryInterval)
				time.Sleep(config.RetryInterval)
			}
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode == 200 {
			// Check if the response contains actual block data
			bodyStr := string(body)
			if strings.Contains(bodyStr, "block") && strings.Contains(bodyStr, "number") {
				fmt.Printf("  ✅ Subgraph is accessible and synced\n")
				return nil
			}
			fmt.Printf("  ⚠️  Subgraph responded but may not be fully synced: %s\n", bodyStr)
		} else {
			fmt.Printf("  ❌ HTTP %d: %s\n", resp.StatusCode, string(body))
		}

		if i < config.MaxRetries-1 {
			fmt.Printf("  ⏱️  Waiting %v before retry...\n", config.RetryInterval)
			time.Sleep(config.RetryInterval)
		}
	}

	return fmt.Errorf("graph node at %s is not accessible after %d attempts", graphURL, config.MaxRetries)
}

// generateTestPrivateKeys generates deterministic test private keys for local testing
func generateTestPrivateKeys() map[string]string {
	// These are hardcoded test keys for reproducibility in tests
	return map[string]string{
		"default":  "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", // anvil account 0 (deployer)
		"batcher0": "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d", // anvil account 1 (batcher/confirmer)
		"dis0":     "d1d51de8ce6bbaac0572e481268232898bfe46491766214c5738929dd557c552",
		"dis1":     "6374444d520f8ae51eee2683f4790644ee5f2d95ca4382fa78021e0460cb1663",
		"opr0":     "a2788f1c26c799b7e1ac32ababc0b598fc7e9c6fc3d319c461ae67ffb1ee57dd",
		"opr1":     "ea25637d76e7ddae9dab9bfac7467d76a1e3bf2d67941b267edc60f2b80d9413",
		"opr2":     "a9ab261a3f506a5e6402dbbaea7bee9496f12117dbe5fa24522e483c07bbe77c",
		"opr3":     "6f84250b1bffd06109bbfa46cc58fb3293008fd43e12a1a5d68d06ab25d060e8",
		"opr4":     "ff7a197fb9c52232f259c26f065c06968eeb982154abcd03d2d08d72641a362a",
		"opr5":     "e5d450c2ffdd19cbf55afbbde7b86e6b841e895546eea7813a9f7360fd38c2db",
		"opr6":     "a4c5553f2d13f96bac694272e94446bfe5e15ed853628c4bd9916e2b5509f956",
		"opr7":     "ef49de2f52c0552484214ebe8e5ba2b13a53dafda560584c1e2426e33dd699a3",
		"staker0":  "aa2b0489fc587a3d8ecac7d97ddea9fa4f2e23e53381ddd8f3b5356287706c28",
		"staker1":  "530f8ec291b5f48481809aa0d5d30f49e32d90620cddc7c178175c69229dbcfe",
		"staker2":  "253f81e5e1c027cf072a27184306b719f851b5b0f6338abe7e595e67ec7c6577",
		"staker3":  "56d6d5d6d7e808ee3cd70cbd44e6d23f1a736e3f94b376ff8a57f61d4fbccd39",
		"relay0":   "f820cde94ba36deefac7ba6a9d12f504b87bfb205c0c87f749008792bb8ba9c3",
		"relay1":   "efbd203977694c18ee6da3a2a42ba13dc95d769a9c814a9fc17e85f0e5eb8360",
		"relay2":   "a1bd1b667b2f37d4ce06d88a9d72e717943a8036ef2e10dc6df419698a77bb07",
		"relay3":   "824435bd114abbf405ad1f7b35fe9421346fb09b1b4cb9a67eea32fe68ff651c",
		"churner":  "b3fec0e8fa0461216ea04ea15faec83cc259e2b066561206f8f455171bdc6de3",
		"dataapi":  "40cc6882bb859e5ae339629f80c559c0c0a85ecca5eb2c58529dbde78a0a5ce4",
	}
}
