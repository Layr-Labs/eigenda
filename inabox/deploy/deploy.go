package deploy

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	caws "github.com/Layr-Labs/eigenda/common/aws"
	testdeployment "github.com/Layr-Labs/eigenda/common/testinfra/deployment"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	churnerImage   = "ghcr.io/layr-labs/eigenda/churner:local"
	disImage       = "ghcr.io/layr-labs/eigenda/disperser:local"
	encoderImage   = "ghcr.io/layr-labs/eigenda/encoder:local"
	batcherImage   = "ghcr.io/layr-labs/eigenda/batcher:local"
	nodeImage      = "ghcr.io/layr-labs/eigenda/node:local"
	retrieverImage = "ghcr.io/layr-labs/eigenda/retriever:local"
	relayImage     = "ghcr.io/layr-labs/eigenda/relay:local"
)

// getKeyString retrieves a ECDSA private key string for a given Ethereum account
func (env *Config) getKeyString(name string) string {
	key, _ := env.getKey(name)
	keyInt, ok := new(big.Int).SetString(key, 0)
	if !ok {
		log.Panicf("Error: could not parse key %s", key)
	}
	return keyInt.String()
}

// deployEigenDAContracts deploys EigenDA core system and peripheral contracts using testinfra
func (env *Config) deployEigenDAContracts() {
	log.Print("Deploy the EigenDA and EigenLayer contracts using testinfra")

	// Get deployer configuration
	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	if !ok {
		log.Panicf("Deployer improperly configured")
	}

	// Create testinfra contract deployment manager
	contractDeployer := testdeployment.ContractDeployer{
		Name:            deployer.Name,
		RPC:             deployer.RPC,
		PrivateKey:      env.Pks.EcdsaMap[deployer.Name].PrivateKey,
		DeploySubgraphs: deployer.DeploySubgraphs,
		VerifyContracts: deployer.VerifyContracts,
		VerifierURL:     deployer.VerifierURL,
		Slow:            deployer.Slow,
	}

	manager := testdeployment.NewContractDeploymentManager(env.rootPath, contractDeployer)

	// Prepare deployment configuration
	numStrategies := len(env.Services.Stakes)
	stakeDistribution := make([][]float32, numStrategies)
	stakeTotals := make([]float32, numStrategies)

	for quorum, stake := range env.Services.Stakes {
		stakeDistribution[quorum] = stake.Distribution
		stakeTotals[quorum] = stake.Total
	}

	// Build private keys map for deployment
	privateKeys := make(map[string]string)
	for i := 0; i < len(env.Services.Stakes[0].Distribution); i++ {
		stakerName := fmt.Sprintf("staker%d", i)
		operatorName := fmt.Sprintf("opr%d", i)
		privateKeys[stakerName] = env.getKeyString(stakerName)
		privateKeys[operatorName] = env.getKeyString(operatorName)
	}
	privateKeys["batcher0"] = env.getKeyString("batcher0")

	deployConfig := testdeployment.GenerateEigenDADeployConfig(
		numStrategies,
		env.Services.Counts.NumMaxOperatorCount,
		stakeDistribution,
		stakeTotals,
		privateKeys,
	)

	// Deploy contracts using testinfra
	err := manager.DeployEigenDAContracts(deployConfig)
	if err != nil {
		log.Panicf("Failed to deploy EigenDA contracts: %s", err.Error())
	}

	// Copy contract addresses back to env
	contracts := manager.GetContractAddresses()
	env.EigenDA.ProxyAdmin = contracts.ProxyAdmin
	env.EigenDA.PauserRegistry = contracts.PauserRegistry
	env.EigenDA.DelegationManager = contracts.DelegationManager
	env.EigenDA.Slasher = contracts.Slasher
	env.EigenDA.StrategyManager = contracts.StrategyManager
	env.EigenDA.EigenPodManager = contracts.EigenPodManager
	env.EigenDA.AVSDirectory = contracts.AVSDirectory
	env.EigenDA.EigenDADirectory = contracts.EigenDADirectory
	env.EigenDA.RewardsCoordinator = contracts.RewardsCoordinator
	env.EigenDA.StrategyFactory = contracts.StrategyFactory
	env.EigenDA.StrategyBeacon = contracts.StrategyBeacon
	env.EigenDA.StrategyBase = contracts.StrategyBase
	env.EigenDA.StrategyBaseTVLLimits = contracts.StrategyBaseTVLLimits
	env.EigenDA.Token = contracts.Token
	env.EigenDA.ServiceManager = contracts.ServiceManager
	env.EigenDA.RegistryCoordinator = contracts.RegistryCoordinator
	env.EigenDA.BlsApkRegistry = contracts.BlsApkRegistry
	env.EigenDA.StakeRegistry = contracts.StakeRegistry
	env.EigenDA.IndexRegistry = contracts.IndexRegistry
	env.EigenDA.OperatorStateRetriever = contracts.OperatorStateRetriever
	env.EigenDA.ServiceManagerImplementation = contracts.ServiceManagerImplementation
	env.EigenDAV1CertVerifier = contracts.EigenDAV1CertVerifier
	env.EigenDAV2CertVerifier = contracts.EigenDAV2CertVerifier
	env.EigenDA.CertVerifier = contracts.EigenDAV2CertVerifier
	env.EigenDA.CertVerifierRouter = contracts.EigenDACertVerifierRouter
}

// Deploys a EigenDA experiment
// TODO: Figure out what necessitates experiment nomenclature
func (env *Config) DeployExperiment() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	defer env.SaveTestConfig()

	log.Print("Deploying experiment...")

	// Log to file
	f, err := os.OpenFile(filepath.Join(env.Path, "deploy.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("error opening file: %v", err)
	}
	defer core.CloseLogOnError(f, f.Name(), nil)
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	// Create a new experiment and deploy the contracts

	err = env.LoadPrivateKeys()
	if err != nil {
		log.Panicf("could not load private keys: %v", err)
	}

	if env.EigenDA.Deployer != "" && !env.IsEigenDADeployed() {
		fmt.Println("Deploying EigenDA")
		env.deployEigenDAContracts()
	}

	if deployer, ok := env.GetDeployer(env.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
		startBlock := GetLatestBlockNumber(env.Deployers[0].RPC)

		// Use testinfra subgraph deployment when available, otherwise fallback to legacy
		fmt.Println("Using testinfra-based subgraph deployment")
		err := env.deploySubgraphsWithTestinfra(startBlock)
		if err != nil {
			log.Panicf("Testinfra subgraph deployment failed: %v", err)
		}
	}

	// Ideally these should be set in GenerateAllVariables, but they need to be used in GenerateDisperserKeypair
	// which is called before GenerateAllVariables

	// Check if AWS_ENDPOINT_URL is set in environment (for dynamic testcontainer ports)
	if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
		env.localstackEndpoint = endpoint
	} else {
		env.localstackEndpoint = "http://localhost:4570"
	}
	env.localstackRegion = "us-east-1"

	// Update dynamic URLs in globals config to use testcontainer ports
	// This ensures environment files get the correct URLs instead of static ones from YAML
	if env.Services.Variables["globals"] == nil {
		env.Services.Variables["globals"] = make(map[string]string)
	}

	// Update RPC URL to use dynamic port from testcontainers
	if len(env.Deployers) > 0 && env.Deployers[0].RPC != "" {
		env.Services.Variables["globals"]["CHAIN_RPC"] = env.Deployers[0].RPC
	}

	// Update AWS endpoint URL to use dynamic port from testcontainers
	if env.localstackEndpoint != "" {
		env.Services.Variables["globals"]["AWS_ENDPOINT_URL"] = env.localstackEndpoint
	}

	fmt.Println("Generating disperser keypair")
	err = env.GenerateDisperserKeypair()
	if err != nil {
		log.Panicf("could not generate disperser keypair: %v", err)
	}

	fmt.Println("Generating variables")
	env.GenerateAllVariables()

	fmt.Println("Test environment has successfully deployed!")
}

// GenerateDisperserKeypair generates a disperser keypair using AWS KMS.
func (env *Config) GenerateDisperserKeypair() error {

	// Generate a keypair in AWS KMS

	keyManager := kms.New(kms.Options{
		Region:       env.localstackRegion,
		BaseEndpoint: aws.String(env.localstackEndpoint),
	})

	createKeyOutput, err := keyManager.CreateKey(context.Background(), &kms.CreateKeyInput{
		KeySpec:  types.KeySpecEccSecgP256k1,
		KeyUsage: types.KeyUsageTypeSignVerify,
	})
	if err != nil {
		if strings.Contains(err.Error(), "connect: connection refused") {
			log.Printf("Unable to reach local stack, skipping disperser keypair generation. Error: %v", err)
			err = nil
		}
		return err
	}

	env.DisperserKMSKeyID = *createKeyOutput.KeyMetadata.KeyId

	// Load the public key and convert it to an Ethereum address

	key, err := caws.LoadPublicKeyKMS(context.Background(), keyManager, env.DisperserKMSKeyID)
	if err != nil {
		return fmt.Errorf("could not load public key: %v", err)
	}

	env.DisperserAddress = crypto.PubkeyToAddress(*key)
	log.Printf("Generated disperser keypair: key ID: %s, address: %s",
		env.DisperserKMSKeyID, env.DisperserAddress.Hex())

	return nil
}

// RegisterDisperserKeypair registers the disperser's public key on-chain using testinfra.
func (env *Config) RegisterDisperserKeypair(ethClient common.EthClient) error {
	log.Printf("RegisterDisperserKeypair: Starting disperser registration")
	log.Printf("RegisterDisperserKeypair: Disperser address to register: %s", env.DisperserAddress.Hex())
	
	// Get deployer configuration
	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	if !ok {
		return fmt.Errorf("deployer improperly configured")
	}

	// Create testinfra contract deployment manager
	contractDeployer := testdeployment.ContractDeployer{
		Name:            deployer.Name,
		RPC:             deployer.RPC,
		PrivateKey:      env.Pks.EcdsaMap[deployer.Name].PrivateKey,
		DeploySubgraphs: deployer.DeploySubgraphs,
		VerifyContracts: deployer.VerifyContracts,
		VerifierURL:     deployer.VerifierURL,
		Slow:            deployer.Slow,
	}

	manager := testdeployment.NewContractDeploymentManager(env.rootPath, contractDeployer)
	// Set the contract addresses in the manager
	manager.EigenDAContracts.OperatorStateRetriever = env.EigenDA.OperatorStateRetriever
	manager.EigenDAContracts.ServiceManager = env.EigenDA.ServiceManager
	
	log.Printf("RegisterDisperserKeypair: ServiceManager address: %s", env.EigenDA.ServiceManager)
	log.Printf("RegisterDisperserKeypair: OperatorStateRetriever address: %s", env.EigenDA.OperatorStateRetriever)

	err := manager.RegisterDisperserAddress(ethClient, env.DisperserAddress)
	if err != nil {
		log.Printf("RegisterDisperserKeypair: Failed to register disperser: %v", err)
	} else {
		log.Printf("RegisterDisperserKeypair: Successfully registered disperser")
	}
	return err
}

// RegisterBlobVersionAndRelays initializes blob versions in ThresholdRegistry contract
// and relays in RelayRegistry contract using testinfra
func (env *Config) RegisterBlobVersionAndRelays(ethClient common.EthClient) {
	// Get deployer configuration
	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	if !ok {
		log.Panicf("Deployer improperly configured")
	}

	// Create testinfra contract deployment manager
	contractDeployer := testdeployment.ContractDeployer{
		Name:            deployer.Name,
		RPC:             deployer.RPC,
		PrivateKey:      env.Pks.EcdsaMap[deployer.Name].PrivateKey,
		DeploySubgraphs: deployer.DeploySubgraphs,
		VerifyContracts: deployer.VerifyContracts,
		VerifierURL:     deployer.VerifierURL,
		Slow:            deployer.Slow,
	}

	manager := testdeployment.NewContractDeploymentManager(env.rootPath, contractDeployer)
	// Set the service manager address in the manager
	manager.EigenDAContracts.ServiceManager = env.EigenDA.ServiceManager

	// Convert BlobVersionParams to testinfra format
	blobVersionParams := make([]testdeployment.BlobVersionParam, len(env.BlobVersionParams))
	for i, param := range env.BlobVersionParams {
		blobVersionParams[i] = testdeployment.BlobVersionParam{
			MaxNumOperators: param.MaxNumOperators,
			NumChunks:       param.NumChunks,
			CodingRate:      param.CodingRate,
		}
	}

	// Collect relay ports
	relayPorts := make([]string, len(env.Relays))
	for i, relayVars := range env.Relays {
		relayPorts[i] = relayVars.RELAY_GRPC_PORT
	}

	err := manager.RegisterBlobVersionsAndRelays(ethClient, blobVersionParams, relayPorts)
	if err != nil {
		log.Panicf("Failed to register blob versions and relays: %s", err)
	}
}

// TODO: Supply the test path to the runner utility
func (env *Config) StartBinaries() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	log.Printf("Starting binaries with start-detached command...")
	
	err := execCmd("./bin.sh", []string{"start-detached"}, []string{}, true)

	if err != nil {
		log.Panicf("Failed to start binaries. Err: %s", err)
	}
	log.Printf("Binaries started successfully")
}

// TODO: Supply the test path to the runner utility
func (env *Config) StopBinaries() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	err := execCmd("./bin.sh", []string{"stop"}, []string{}, true)
	if err != nil {
		log.Panicf("Failed to stop binaries. Err: %s", err)
	}
}

func (env *Config) ForceStopBinaries() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	err := execCmd("./bin.sh", []string{"force-stop"}, []string{}, true)
	if err != nil {
		log.Printf("Force stop completed with some errors (this is expected): %s", err)
	}
}

func (env *Config) StartAnvil() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	err := execCmd("./bin.sh", []string{"start-anvil"}, []string{}, false) // printing output causes hang
	if err != nil {
		log.Panicf("Failed to start anvil. Err: %s", err)
	}
}

func (env *Config) StopAnvil() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	err := execCmd("./bin.sh", []string{"stop-anvil"}, []string{}, true)
	if err != nil {
		log.Panicf("Failed to stop anvil. Err: %s", err)
	}
}

func (env *Config) RunNodePluginBinary(operation string, operator OperatorVars) {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))

	socket := string(core.MakeOperatorSocket(operator.NODE_HOSTNAME, operator.NODE_DISPERSAL_PORT, operator.NODE_RETRIEVAL_PORT, operator.NODE_V2_DISPERSAL_PORT, operator.NODE_V2_RETRIEVAL_PORT))

	envVars := []string{
		"NODE_OPERATION=" + operation,
		"NODE_ECDSA_KEY_FILE=" + operator.NODE_ECDSA_KEY_FILE,
		"NODE_BLS_KEY_FILE=" + operator.NODE_BLS_KEY_FILE,
		"NODE_ECDSA_KEY_PASSWORD=" + operator.NODE_ECDSA_KEY_PASSWORD,
		"NODE_BLS_KEY_PASSWORD=" + operator.NODE_BLS_KEY_PASSWORD,
		"NODE_SOCKET=" + socket,
		"NODE_QUORUM_ID_LIST=" + operator.NODE_QUORUM_ID_LIST,
		"NODE_CHAIN_RPC=" + operator.NODE_CHAIN_RPC,
		"NODE_EIGENDA_DIRECTORY=" + operator.NODE_EIGENDA_DIRECTORY,
		"NODE_BLS_OPERATOR_STATE_RETRIVER=" + operator.NODE_BLS_OPERATOR_STATE_RETRIVER,
		"NODE_EIGENDA_SERVICE_MANAGER=" + operator.NODE_EIGENDA_SERVICE_MANAGER,
		"NODE_CHURNER_URL=" + operator.NODE_CHURNER_URL,
		"NODE_NUM_CONFIRMATIONS=0",
	}

	err := execCmd("./node-plugin.sh", []string{}, envVars, true)

	if err != nil {
		log.Panicf("Failed to run node plugin. Err: %s", err)
	}
}

// deploySubgraphsWithTestinfra deploys subgraphs using the testinfra package
func (env *Config) deploySubgraphsWithTestinfra(startBlock int) error {
	if !env.Environment.IsLocal() {
		return fmt.Errorf("testinfra subgraph deployment only supported for local environments")
	}

	// Log contract addresses before subgraph deployment
	log.Printf("Subgraph deployment - ServiceManager: %s, RegistryCoordinator: %s, BlsApkRegistry: %s",
		env.EigenDA.ServiceManager, env.EigenDA.RegistryCoordinator, env.EigenDA.BlsApkRegistry)
	
	// Prepare subgraph deployment configuration
	deployConfig := testdeployment.SubgraphDeploymentConfig{
		RootPath: env.rootPath,
		Subgraphs: []testdeployment.SubgraphConfig{
			{
				Name:    "eigenda-operator-state",
				Path:    "eigenda-operator-state",
				Enabled: true,
			},
			{
				Name:    "eigenda-batch-metadata",
				Path:    "eigenda-batch-metadata",
				Enabled: true,
			},
		},
		EigenDAConfig: testdeployment.EigenDAContractAddresses{
			RegistryCoordinator: env.EigenDA.RegistryCoordinator,
			BlsApkRegistry:      env.EigenDA.BlsApkRegistry,
			ServiceManager:      env.EigenDA.ServiceManager,
		},
	}

	// Use the standalone testinfra deployment function with the URLs we have
	return testdeployment.DeploySubgraphsWithURLs(deployConfig, env.GraphAdminURL, env.IPFSURL, startBlock)
}
