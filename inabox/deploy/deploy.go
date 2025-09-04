package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	caws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	relayreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARelayRegistry"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	thresholdreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAThresholdRegistry"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gcommon "github.com/ethereum/go-ethereum/common"
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
func (env *Config) getKeyString(name string) (string, error) {
	key, _, err := env.getKey(name)
	if err != nil {
		return "", fmt.Errorf("could not get key for %s: %w", name, err)
	}

	keyInt, ok := new(big.Int).SetString(key, 0)
	if !ok {
		return "", fmt.Errorf("could not parse key %s", key)
	}

	return keyInt.String(), nil
}

// generateV1CertVerifierDeployConfig generates the input config used for deploying the V1 CertVerifier
// NOTE: this will be killed in the future with eventual deprecation of V1
func (env *Config) generateV1CertVerifierDeployConfig(ethClient common.EthClient) V1CertVerifierDeployConfig {
	config := V1CertVerifierDeployConfig{
		ServiceManager:                env.EigenDA.ServiceManager,
		RequiredQuorums:               []uint32{0, 1},
		RequiredAdversarialThresholds: []uint32{33, 33},
		RequiredConfirmationQuorums:   []uint32{55, 55},
	}

	return config
}

// generateEigenDADeployConfig generates input config fed into SetUpEigenDA.s.sol foundry script
func (env *Config) generateEigenDADeployConfig() (EigenDADeployConfig, error) {

	operators := make([]string, 0)
	stakers := make([]string, 0)
	maxOperatorCount := env.Services.Counts.NumMaxOperatorCount

	numStrategies := len(env.Services.Stakes)
	total := make([]float32, numStrategies)
	stakes := make([][]string, numStrategies)

	for quorum, stake := range env.Services.Stakes {
		for _, s := range stake.Distribution {
			total[quorum] += s
		}
	}

	for quorum := 0; quorum < numStrategies; quorum++ {
		stakes[quorum] = make([]string, len(env.Services.Stakes[quorum].Distribution))
		for ind, stake := range env.Services.Stakes[quorum].Distribution {
			stakes[quorum][ind] = strconv.FormatFloat(float64(stake/total[quorum]*env.Services.Stakes[quorum].Total), 'f', 0, 32)
		}
	}

	for i := 0; i < len(env.Services.Stakes[0].Distribution); i++ {
		stakerName := fmt.Sprintf("staker%d", i)
		operatorName := fmt.Sprintf("opr%d", i)

		// Get keys for staker and operator
		stakerKey, err := env.getKeyString(stakerName)
		if err != nil {
			return EigenDADeployConfig{}, fmt.Errorf("failed to get key for %s: %w", stakerName, err)
		}

		operatorKey, err := env.getKeyString(operatorName)
		if err != nil {
			return EigenDADeployConfig{}, fmt.Errorf("failed to get key for %s: %w", operatorName, err)
		}

		stakers = append(stakers, stakerKey)
		operators = append(operators, operatorKey)
	}

	// Get batcher0 key
	batcherKey, err := env.getKeyString("batcher0")
	if err != nil {
		return EigenDADeployConfig{}, fmt.Errorf("failed to get key for batcher0: %w", err)
	}

	config := EigenDADeployConfig{
		UseDefaults:         true,
		NumStrategies:       numStrategies,
		MaxOperatorCount:    maxOperatorCount,
		StakerPrivateKeys:   stakers,
		StakerTokenAmounts:  stakes,
		OperatorPrivateKeys: operators,
		ConfirmerPrivateKey: batcherKey,
	}

	return config, nil
}

// deployEigenDAContracts deploys EigenDA core system and peripheral contracts on local anvil chain
func (env *Config) deployEigenDAContracts() error {
	env.logger.Info("Deploy the EigenDA and EigenLayer contracts")

	// get deployer
	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	if !ok {
		return fmt.Errorf("deployer improperly configured")
	}

	if err := changeDirectory(filepath.Join(env.rootPath, "contracts")); err != nil {
		return fmt.Errorf("failed to change directories: %w", err)
	}

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil {
		env.logger.Info("Successfully changed to absolute path", "path", cwd)
	}

	eigendaDeployConfig, err := env.generateEigenDADeployConfig()
	if err != nil {
		return fmt.Errorf("error generating eigenda deploy config: %w", err)
	}

	data, err := json.Marshal(&eigendaDeployConfig)
	if err != nil {
		return fmt.Errorf("error marshaling eigenda deploy config: %w", err)
	}
	err = writeFile("script/input/eigenda_deploy_config.json", data)
	if err != nil {
		return fmt.Errorf("error writing eigenda deploy config: %w", err)
	}

	env.logger.Info("Executing EigenDA deployer script", "script", "script/SetUpEigenDA.s.sol:SetupEigenDA", "rpc", deployer.RPC, "deployer", deployer.Name, "privateKey", env.Pks.EcdsaMap[deployer.Name].PrivateKey)
	err = execForgeScript(
		"script/SetUpEigenDA.s.sol:SetupEigenDA",
		env.Pks.EcdsaMap[deployer.Name].PrivateKey,
		deployer,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to execute EigenDA deployer script: %w", err)
	}

	// Add relevant addresses to path
	data, err = readFile("script/output/eigenda_deploy_output.json")
	if err != nil {
		return fmt.Errorf("error reading eigenda deploy output: %w", err)
	}

	err = json.Unmarshal(data, &env.EigenDA)
	if err != nil {
		return fmt.Errorf("error unmarshaling eigenda deploy output: %w", err)
	}

	ethClient, err := geth.NewClient(geth.EthClientConfig{
		RPCURLs:          []string{deployer.RPC},
		PrivateKeyString: env.Pks.EcdsaMap[deployer.Name].PrivateKey[2:],
		NumConfirmations: 0,
		NumRetries:       0,
	}, gcommon.Address{}, 0, env.logger)
	if err != nil {
		return fmt.Errorf("error creating eth client: %w", err)
	}

	certVerifierV1DeployCfg := env.generateV1CertVerifierDeployConfig(ethClient)
	data, err = json.Marshal(&certVerifierV1DeployCfg)
	if err != nil {
		return fmt.Errorf("error marshaling certverifier config: %w", err)
	}

	// NOTE: this is pretty janky and is a short-term solution until V1 contract usage
	//       can be deprecated.
	if err := writeFile("script/deploy/certverifier/config/v1/inabox_deploy_config_v1.json", data); err != nil {
		return fmt.Errorf("error writing certverifier config: %w", err)
	}

	env.logger.Info("Executing CertVerifierDeployerV1 script", "script", "script/deploy/certverifier/CertVerifierDeployerV1.s.sol:CertVerifierDeployerV1", "rpc", deployer.RPC, "deployer", deployer.Name, "privateKey", env.Pks.EcdsaMap[deployer.Name].PrivateKey)
	if err := execForgeScript("script/deploy/certverifier/CertVerifierDeployerV1.s.sol:CertVerifierDeployerV1", env.Pks.EcdsaMap[deployer.Name].PrivateKey, deployer, []string{"--sig", "run(string, string)", "inabox_deploy_config_v1.json", "inabox_v1_deploy.json"}); err != nil {
		return fmt.Errorf("failed to execute CertVerifierDeployerV1 script: %w", err)
	}

	data, err = readFile("script/deploy/certverifier/output/inabox_v1_deploy.json")
	if err != nil {
		return fmt.Errorf("error reading certverifier output: %w", err)
	}

	var verifierAddress struct{ EigenDACertVerifier string }
	err = json.Unmarshal(data, &verifierAddress)
	if err != nil {
		return fmt.Errorf("error unmarshaling verifier address: %w", err)
	}
	env.EigenDAV1CertVerifier = verifierAddress.EigenDACertVerifier

	return nil
}

// Deploys a EigenDA experiment
// TODO: Figure out what necessitates experiment nomenclature
func (env *Config) DeployExperiment() error {
	if err := changeDirectory(filepath.Join(env.rootPath, "inabox")); err != nil {
		return fmt.Errorf("error changing directories: %w", err)
	}

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil {
		env.logger.Info("Successfully changed to absolute path", "path", cwd)
	}

	defer env.SaveTestConfig()

	env.logger.Info("Deploying experiment...")

	// Create a new experiment and deploy the contracts

	err := env.loadPrivateKeys()
	if err != nil {
		return fmt.Errorf("could not load private keys: %w", err)
	}

	if env.EigenDA.Deployer != "" && !env.IsEigenDADeployed() {
		env.logger.Info("Deploying EigenDA")
		err = env.deployEigenDAContracts()
		if err != nil {
			return fmt.Errorf("error deploying EigenDA contracts: %w", err)
		}
	}

	if deployer, ok := env.GetDeployer(env.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
		startBlock, err := GetLatestBlockNumber(env.logger, env.Deployers[0].RPC)
		if err != nil {
			return fmt.Errorf("error getting latest block number: %w", err)
		}

		err = env.deploySubgraphs(startBlock)
		if err != nil {
			return fmt.Errorf("error deploying subgraphs: %w", err)
		}
	}

	// Ideally these should be set in GenerateAllVariables, but they need to be used in GenerateDisperserKeypair
	// which is called before GenerateAllVariables
	env.localstackEndpoint = "http://localhost:4570"
	env.localstackRegion = "us-east-1"

	env.logger.Info("Generating disperser keypair")
	err = env.GenerateDisperserKeypair()
	if err != nil {
		env.logger.Errorf("could not generate disperser keypair: %v", err)
		panic(err)
	}

	env.logger.Info("Generating variables")
	env.GenerateAllVariables()

	// Register blob versions, relays, and disperser keypair
	if env.EigenDA.Deployer != "" && env.IsEigenDADeployed() {
		ethClient, err := geth.NewMultiHomingClient(geth.EthClientConfig{
			RPCURLs:          []string{env.Deployers[0].RPC},
			PrivateKeyString: env.Pks.EcdsaMap[env.EigenDA.Deployer].PrivateKey[2:],
			NumConfirmations: 0,
			NumRetries:       3,
		}, gcommon.Address{}, env.logger)
		if err != nil {
			env.logger.Errorf("could not create eth client for registration: %v", err)
		} else {
			env.logger.Info("Registering blob versions and relays")
			env.RegisterBlobVersionAndRelays(ethClient)

			env.logger.Info("Registering disperser keypair")
			err = env.RegisterDisperserKeypair(ethClient)
			if err != nil {
				env.logger.Errorf("could not register disperser keypair: %v", err)
			}
		}
	}

	env.logger.Info("Test environment has successfully deployed!")
	return nil
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
			env.logger.Warnf("Unable to reach local stack, skipping disperser keypair generation. Error: %v", err)
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
	env.logger.Infof("Generated disperser keypair: key ID: %s, address: %s",
		env.DisperserKMSKeyID, env.DisperserAddress.Hex())

	return nil
}

// RegisterDisperserKeypair registers the disperser's public key on-chain.
func (env *Config) RegisterDisperserKeypair(ethClient common.EthClient) error {
	// Write the disperser's public key to on-chain storage
	writer, err := eth.NewWriter(
		env.logger,
		ethClient,
		env.EigenDA.OperatorStateRetriever,
		env.EigenDA.ServiceManager,
	)
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}

	err = writer.SetDisperserAddress(context.Background(), env.DisperserAddress)
	if err != nil {
		return fmt.Errorf("could not set disperser address: %v", err)
	}

	// Read the disperser's public key from on-chain storage to verify it was written correctly

	retryTimeout := time.Now().Add(1 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)

	for time.Now().Before(retryTimeout) {
		address, err := writer.GetDisperserAddress(context.Background(), 0)
		if err != nil {
			env.logger.Warnf("could not get disperser address: %v", err)
		} else {
			if address != env.DisperserAddress {
				return fmt.Errorf("expected disperser address %s, got %s", env.DisperserAddress, address)
			}
			return nil
		}

		<-ticker.C
	}

	return fmt.Errorf("timed out waiting for disperser address to be set")
}

// RegisterBlobVersionAndRelays initializes blob versions in ThresholdRegistry contract
// and relays in RelayRegistry contract
func (env *Config) RegisterBlobVersionAndRelays(ethClient common.EthClient) {
	dasmAddr := gcommon.HexToAddress(env.EigenDA.ServiceManager)
	contractEigenDAServiceManager, err := eigendasrvmg.NewContractEigenDAServiceManager(dasmAddr, ethClient)
	if err != nil {
		env.logger.Fatal("Error creating EigenDAServiceManager contract", "error", err)
	}
	thresholdRegistryAddr, err := contractEigenDAServiceManager.EigenDAThresholdRegistry(&bind.CallOpts{})
	if err != nil {
		env.logger.Fatal("Error getting threshold registry address", "error", err)
	}
	contractThresholdRegistry, err := thresholdreg.NewContractEigenDAThresholdRegistry(thresholdRegistryAddr, ethClient)
	if err != nil {
		env.logger.Fatal("Error creating threshold registry contract", "error", err)
	}
	opts, err := ethClient.GetNoSendTransactOpts()
	if err != nil {
		env.logger.Fatal("Error getting transaction opts", "error", err)
	}
	for _, blobVersionParam := range env.BlobVersionParams {
		txn, err := contractThresholdRegistry.AddVersionedBlobParams(opts, thresholdreg.EigenDATypesV1VersionedBlobParams{
			MaxNumOperators: blobVersionParam.MaxNumOperators,
			NumChunks:       blobVersionParam.NumChunks,
			CodingRate:      uint8(blobVersionParam.CodingRate),
		})
		if err != nil {
			env.logger.Fatal("Error adding versioned blob params", "error", err)
		}
		err = ethClient.SendTransaction(context.Background(), txn)
		if err != nil {
			env.logger.Fatal("Error sending blob version transaction", "error", err)
		}
	}

	relayAddr, err := contractEigenDAServiceManager.EigenDARelayRegistry(&bind.CallOpts{})
	if err != nil {
		env.logger.Fatal("Error getting relay registry address", "error", err)
	}
	contractRelayRegistry, err := relayreg.NewContractEigenDARelayRegistry(relayAddr, ethClient)
	if err != nil {
		env.logger.Fatal("Error creating relay registry contract", "error", err)
	}

	ethAddr := ethClient.GetAccountAddress()
	for _, relayVars := range env.Relays {
		url := fmt.Sprintf("0.0.0.0:%s", relayVars.RELAY_GRPC_PORT)
		txn, err := contractRelayRegistry.AddRelayInfo(opts, relayreg.EigenDATypesV2RelayInfo{
			RelayAddress: ethAddr,
			RelayURL:     url,
		})
		if err != nil {
			env.logger.Fatal("Error adding relay info", "error", err)
		}
		err = ethClient.SendTransaction(context.Background(), txn)
		if err != nil {
			env.logger.Fatal("Error sending relay transaction", "error", err)
		}
	}
}

// TODO: Supply the test path to the runner utility
func (env *Config) StartBinaries() {
	if err := changeDirectory(filepath.Join(env.rootPath, "inabox")); err != nil {
		env.logger.Fatal("Error changing directories", "error", err)
	}

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil {
		env.logger.Info("Successfully changed to absolute path", "path", cwd)
	}

	env.logger.Info("Starting binaries")
	err := execCmd("./bin.sh", []string{"start-detached"}, []string{}, true)
	if err != nil {
		env.logger.Fatal("Failed to start binaries, check testdata directory for more information", "error", err)
	}

	env.logger.Info("Binaries started successfully!")
}

// TODO: Supply the test path to the runner utility
func (env *Config) StopBinaries() {
	if err := changeDirectory(filepath.Join(env.rootPath, "inabox")); err != nil {
		env.logger.Fatal("Error changing directories", "error", err)
	}

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil {
		env.logger.Info("Successfully changed to absolute path", "path", cwd)
	}

	err := execCmd("./bin.sh", []string{"stop"}, []string{}, true)
	if err != nil {
		env.logger.Fatal("Failed to stop binaries", "error", err)
	}
}

func (env *Config) RunNodePluginBinary(operation string, operator OperatorVars) {
	if err := changeDirectory(filepath.Join(env.rootPath, "inabox")); err != nil {
		env.logger.Fatal("Error changing directories", "error", err)
	}

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil {
		env.logger.Info("Successfully changed to absolute path", "path", cwd)
	}

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
		env.logger.Fatal("Failed to run node plugin", "error", err)
	}
}
