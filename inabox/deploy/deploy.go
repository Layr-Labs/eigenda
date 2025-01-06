package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	relayreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARelayRegistry"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	thresholdreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAThresholdRegistry"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gcommon "github.com/ethereum/go-ethereum/common"
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

func (env *Config) getKeyString(name string) string {
	key, _ := env.getKey(name)
	keyInt, ok := new(big.Int).SetString(key, 0)
	if !ok {
		log.Panicf("Error: could not parse key %s", key)
	}
	return keyInt.String()
}

func (env *Config) generateEigenDADeployConfig() EigenDADeployConfig {

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

		stakers = append(stakers, env.getKeyString(stakerName))
		operators = append(operators, env.getKeyString(operatorName))
	}

	config := EigenDADeployConfig{
		UseDefaults:         true,
		NumStrategies:       numStrategies,
		MaxOperatorCount:    maxOperatorCount,
		StakerPrivateKeys:   stakers,
		StakerTokenAmounts:  stakes,
		OperatorPrivateKeys: operators,
		ConfirmerPrivateKey: env.getKeyString("batcher0"),
	}

	return config

}

func (env *Config) deployEigenDAContracts() {
	log.Print("Deploy the EigenDA and EigenLayer contracts")

	// get deployer
	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	if !ok {
		log.Panicf("Deployer improperly configured")
	}

	changeDirectory(filepath.Join(env.rootPath, "contracts"))

	eigendaDeployConfig := env.generateEigenDADeployConfig()
	data, err := json.Marshal(&eigendaDeployConfig)
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}
	writeFile("script/input/eigenda_deploy_config.json", data)

	execForgeScript("script/SetUpEigenDA.s.sol:SetupEigenDA", env.Pks.EcdsaMap[deployer.Name].PrivateKey, deployer, nil)

	//add relevant addresses to path
	data = readFile("script/output/eigenda_deploy_output.json")
	err = json.Unmarshal(data, &env.EigenDA)
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}
	execForgeScript("script/MockRollupDeployer.s.sol:MockRollupDeployer", env.Pks.EcdsaMap[deployer.Name].PrivateKey, deployer, []string{"--sig", "run(address)", env.EigenDA.ServiceManager})

	//add rollup address to path
	data = readFile("script/output/mock_rollup_deploy_output.json")
	var rollupAddr struct{ MockRollup string }
	err = json.Unmarshal(data, &rollupAddr)
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}
	env.MockRollup = rollupAddr.MockRollup
}

// Deploys a EigenDA experiment
func (env *Config) DeployExperiment() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	defer env.SaveTestConfig()

	log.Print("Deploying experiment...")

	// Log to file
	f, err := os.OpenFile(filepath.Join(env.Path, "deploy.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	// Create a new experiment and deploy the contracts

	err = env.loadPrivateKeys()
	if err != nil {
		log.Panicf("could not load private keys: %v", err)
	}

	if env.EigenDA.Deployer != "" && !env.IsEigenDADeployed() {
		fmt.Println("Deploying EigenDA")
		env.deployEigenDAContracts()
	}

	log.Print("Generating disperser keypair")
	keyID, disperserAddress, err := generateDisperserKeypair()
	if err != nil {
		log.Panicf("could not generate disperser keypair: %v", err)
	}
	env.DisperserAddress = disperserAddress
	env.DisperserKMSKeyID = keyID

	log.Print("Updating disperser address")
	err = env.updateDisperserAddress()
	if err != nil {
		log.Panicf("could not update disperser address: %v", err)
	}

	if deployer, ok := env.GetDeployer(env.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
		startBlock := GetLatestBlockNumber(env.Deployers[0].RPC)
		env.deploySubgraphs(startBlock)
	}

	// TODO remove
	b := make([]byte, 2048) // adjust buffer size to be larger than expected stack
	n := runtime.Stack(b, false)
	s := string(b[:n])
	log.Printf("Stack trace:\n %s", s)

	fmt.Println("Generating variables")
	env.GenerateAllVariables()

	fmt.Println("Test environment has successfully deployed!")
}

// GenerateDisperserKeypair generates a disperser keypair using AWS KMS. Returns the key ID and the public address.
func generateDisperserKeypair() (string, gcommon.Address, error) {
	keyManager := kms.New(kms.Options{
		Region:       "us-east-1",
		BaseEndpoint: aws.String("http://localhost:4570"), // TODO don't hard code this
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
		return "", gcommon.Address{}, err
	}

	keyID := *createKeyOutput.KeyMetadata.KeyId

	key, err := common.LoadPublicKeyKMS(context.Background(), keyManager, keyID)
	if err != nil {
		return "", gcommon.Address{}, err
	}

	publicAddress := crypto.PubkeyToAddress(*key)

	log.Printf("Generated disperser keypair: key ID: %s, address: %s", keyID, publicAddress.Hex())

	return keyID, publicAddress, nil
}

// updateDisperserAddress updates the disperser address in the retriever contract
func (env *Config) updateDisperserAddress() error {
	pk := env.Pks.EcdsaMap["default"].PrivateKey
	pk = strings.TrimPrefix(pk, "0x")
	pk = strings.TrimPrefix(pk, "0X")

	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return fmt.Errorf("could not create logger: %v", err)
	}

	ethClient, err := geth.NewMultiHomingClient(geth.EthClientConfig{
		RPCURLs:          []string{env.Deployers[0].RPC},
		PrivateKeyString: pk,
		NumConfirmations: 0,
		NumRetries:       0,
	}, gcommon.Address{}, logger)
	if err != nil {
		return fmt.Errorf("could not create eth client: %v", err)
	}

	writer, err := eth.NewWriter(
		logger,
		ethClient,
		env.Retriever.RETRIEVER_BLS_OPERATOR_STATE_RETRIVER,
		env.Retriever.RETRIEVER_EIGENDA_SERVICE_MANAGER)
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}

	err = writer.SetDisperserAddress(context.Background(), env.DisperserAddress)
	if err != nil {
		return fmt.Errorf("could not set disperser address: %v", err)
	}

	address, err := writer.GetDisperserAddress(context.Background(), 0)
	if err != nil {
		return fmt.Errorf("could not get disperser address: %v", err)
	}

	if address != env.DisperserAddress {
		return fmt.Errorf("expected disperser address %s, got %s", env.DisperserAddress, address)
	}

	return nil
}

func (env *Config) RegisterBlobVersionAndRelays(ethClient common.EthClient) map[uint32]string {
	dasmAddr := gcommon.HexToAddress(env.EigenDA.ServiceManager)
	contractEigenDAServiceManager, err := eigendasrvmg.NewContractEigenDAServiceManager(dasmAddr, ethClient)
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	thresholdRegistryAddr, err := contractEigenDAServiceManager.EigenDAThresholdRegistry(&bind.CallOpts{})
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	contractThresholdRegistry, err := thresholdreg.NewContractEigenDAThresholdRegistry(thresholdRegistryAddr, ethClient)
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	opts, err := ethClient.GetNoSendTransactOpts()
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	for _, blobVersionParam := range env.BlobVersionParams {
		txn, err := contractThresholdRegistry.AddVersionedBlobParams(opts, thresholdreg.VersionedBlobParams{
			MaxNumOperators: blobVersionParam.MaxNumOperators,
			NumChunks:       blobVersionParam.NumChunks,
			CodingRate:      uint8(blobVersionParam.CodingRate),
		})
		if err != nil {
			log.Panicf("Error: %s", err)
		}
		err = ethClient.SendTransaction(context.Background(), txn)
		if err != nil {
			log.Panicf("Error: %s", err)
		}
	}

	relayAddr, err := contractEigenDAServiceManager.EigenDARelayRegistry(&bind.CallOpts{})
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	contractRelayRegistry, err := relayreg.NewContractEigenDARelayRegistry(relayAddr, ethClient)
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	relays := map[uint32]string{}
	for i, relayVars := range env.Relays {
		url := fmt.Sprintf("0.0.0.0:%s", relayVars.RELAY_GRPC_PORT)
		txn, err := contractRelayRegistry.AddRelayInfo(opts, relayreg.RelayInfo{
			RelayAddress: gcommon.Address{0},
			RelayURL:     url,
		})
		if err != nil {
			log.Panicf("Error: %s", err)
		}
		err = ethClient.SendTransaction(context.Background(), txn)
		if err != nil {
			log.Panicf("Error: %s", err)
		}
		relays[uint32(i)] = url
	}

	return relays
}

// TODO: Supply the test path to the runner utility
func (env *Config) StartBinaries() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	err := execCmd("./bin.sh", []string{"start-detached"}, []string{})

	if err != nil {
		log.Panicf("Failed to start binaries. Err: %s", err)
	}
}

// TODO: Supply the test path to the runner utility
func (env *Config) StopBinaries() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	err := execCmd("./bin.sh", []string{"stop"}, []string{})
	if err != nil {
		log.Panicf("Failed to stop binaries. Err: %s", err)
	}
}

func (env *Config) StartAnvil() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	err := execCmd("./bin.sh", []string{"start-anvil"}, []string{})
	if err != nil {
		log.Panicf("Failed to start anvil. Err: %s", err)
	}
}

func (env *Config) StopAnvil() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	err := execCmd("./bin.sh", []string{"stop-anvil"}, []string{})
	if err != nil {
		log.Panicf("Failed to stop anvil. Err: %s", err)
	}
}

func (env *Config) RunNodePluginBinary(operation string, operator OperatorVars) {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))

	socket := string(core.MakeOperatorSocket(operator.NODE_HOSTNAME, operator.NODE_DISPERSAL_PORT, operator.NODE_RETRIEVAL_PORT))

	envVars := []string{
		"NODE_OPERATION=" + operation,
		"NODE_ECDSA_KEY_FILE=" + operator.NODE_ECDSA_KEY_FILE,
		"NODE_BLS_KEY_FILE=" + operator.NODE_BLS_KEY_FILE,
		"NODE_ECDSA_KEY_PASSWORD=" + operator.NODE_ECDSA_KEY_PASSWORD,
		"NODE_BLS_KEY_PASSWORD=" + operator.NODE_BLS_KEY_PASSWORD,
		"NODE_SOCKET=" + socket,
		"NODE_QUORUM_ID_LIST=" + operator.NODE_QUORUM_ID_LIST,
		"NODE_CHAIN_RPC=" + operator.NODE_CHAIN_RPC,
		"NODE_BLS_OPERATOR_STATE_RETRIVER=" + operator.NODE_BLS_OPERATOR_STATE_RETRIVER,
		"NODE_EIGENDA_SERVICE_MANAGER=" + operator.NODE_EIGENDA_SERVICE_MANAGER,
		"NODE_CHURNER_URL=" + operator.NODE_CHURNER_URL,
		"NODE_NUM_CONFIRMATIONS=0",
	}

	err := execCmd("./node-plugin.sh", []string{}, envVars)

	if err != nil {
		log.Panicf("Failed to run node plugin. Err: %s", err)
	}
}
