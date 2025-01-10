package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Layr-Labs/eigenda/common"
	relayreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARelayRegistry"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	thresholdreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAThresholdRegistry"

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

	if deployer, ok := env.GetDeployer(env.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
		startBlock := GetLatestBlockNumber(env.Deployers[0].RPC)
		env.deploySubgraphs(startBlock)
	}

	fmt.Println("Generating variables")
	env.GenerateAllVariables()

	fmt.Println("Test environment has successfully deployed!")
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
	ethAddr := ethClient.GetAccountAddress()
	for i, relayVars := range env.Relays {
		url := fmt.Sprintf("0.0.0.0:%s", relayVars.RELAY_GRPC_PORT)
		txn, err := contractRelayRegistry.AddRelayInfo(opts, relayreg.RelayInfo{
			RelayAddress: ethAddr,
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
