package ethinitscript

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os/exec"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/inabox/config"
	"github.com/Layr-Labs/eigenda/inabox/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func HasDeployedEigenDAContracts(rpcUrl string, eigendaContractConfig *config.EigenDAContract) bool {
	// Replace with your Ethereum node URL
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// Replace with the address you're checking
	address := common.HexToAddress(eigendaContractConfig.ServiceManager)

	// Context for the network call
	ctx := context.Background()

	// Get the bytecode at the address
	bytecode, err := client.CodeAt(ctx, address, nil) // nil is latest block
	if err != nil {
		log.Fatalf("Failed to get bytecode: %v", err)
	}

	return len(bytecode) > 0
}

func SetupEigenDA(cfg *config.ConfigLock, contractsDir string) *exec.Cmd {
	log.Print("Deploying EigenDA...")
	deployer, ok := cfg.Config.GetDeployer(cfg.Config.EigenDA.Deployer)
	if !ok {
		log.Panicf("Deployer improperly configured")
	}

	var rpcUrl string
	var cmd *exec.Cmd
	if deployer.LocalAnvil {
		log.Print("Starting anvil...")
		var err error
		cmd, err = utils.StartCommand("anvil", "--host", "0.0.0.0")
		if err != nil {
			log.Panicf("Error starting anvil: %v", err)
		}
		err = utils.WaitForServer("localhost:8545", 10*time.Second)
		if err != nil {
			log.Panicf("Error starting anvil: %v", err)
		}
		log.Print("Anvil started.")
		rpcUrl = "http://localhost:8545"
		// TODO: Handle graceful termination of anvil as part of the graceful termination of the script
	} else {
		log.Print("Skipping anvil.")
		rpcUrl = deployer.RPC
	}

	hasDeployed := HasDeployedEigenDAContracts(rpcUrl, &cfg.Config.EigenDA)
	if !hasDeployed {
		deployerKey, _ := cfg.GetKey(deployer.Name)
		eigendaDeployConfig := GenerateEigenDADeployConfig(cfg)
		serviceManagerAddr := cfg.Config.EigenDA.ServiceManager
		DeployEigenDAContracts(deployerKey, rpcUrl, eigendaDeployConfig, serviceManagerAddr, contractsDir)

		log.Print("Test environment has succesfully deployed!")
	} else {
		log.Print("Test environment already was succesfully deployed!")
	}
	return cmd
}

func DeployEigenDAContracts(deployerKey string, rpcUrl string, eigendaDeployConfig *config.EigenDADeployConfig, serviceManagerAddr string, contractsDir string) {
	log.Print("Deploy the EigenDA and EigenLayer contracts")

	dir := utils.MustGetwd()
	defer utils.MustChdir(dir)
	utils.MustChdir(contractsDir)

	data, err := json.Marshal(eigendaDeployConfig)
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}
	utils.MustWriteFile("script/eigenda_deploy_config.json", data)

	execForgeScript("script/SetUpEigenDA.s.sol:SetupEigenDA", deployerKey, rpcUrl, nil)

	blobHeader := &core.BlobHeader{
		QuorumInfos: []*core.BlobQuorumInfo{
			{
				SecurityParam: core.SecurityParam{
					QuorumID:           0,
					QuorumThreshold:    100,
					AdversaryThreshold: 80,
				},
				QuantizationFactor: 1,
			},
		},
	}
	hash, err := blobHeader.GetQuorumBlobParamsHash()
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}
	hashStr := fmt.Sprintf("%x", hash)

	execForgeScript(
		"script/MockRollupDeployer.s.sol:MockRollupDeployer",
		deployerKey,
		rpcUrl,
		[]string{"--sig", "run(address,bytes32,uint256)", serviceManagerAddr, hashStr, big.NewInt(1e18).String()},
	)

	//add rollup address to path
	data = utils.MustReadFile("script/output/mock_rollup_deploy_output.json")
	var rollupAddr struct{ MockRollup string }
	err = json.Unmarshal(data, &rollupAddr)
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}

	// TODO: We don't do anything with this mock rollup right now
}

func GenerateEigenDADeployConfig(lock *config.ConfigLock) *config.EigenDADeployConfig {
	operators := make([]string, 0)
	stakers := make([]string, 0)
	maxOperatorCount := lock.Config.Services.Counts.NumMaxOperatorCount

	total := float32(0)
	stakes := [][]string{make([]string, len(lock.Config.Services.Stakes.Distribution))}

	for _, stake := range lock.Config.Services.Stakes.Distribution {
		total += stake
	}

	for ind, stake := range lock.Config.Services.Stakes.Distribution {
		stakes[0][ind] = strconv.FormatFloat(float64(stake/total*lock.Config.Services.Stakes.Total), 'f', 0, 32)
	}

	for i := 0; i < len(lock.Config.Services.Stakes.Distribution); i++ {
		stakerName := fmt.Sprintf("staker%d", i)
		operatorName := fmt.Sprintf("opr%d", i)

		stakers = append(stakers, lock.GetKeyString(stakerName))
		operators = append(operators, lock.GetKeyString(operatorName))
	}

	deployConfig := &config.EigenDADeployConfig{
		UseDefaults:         true,
		NumStrategies:       1,
		MaxOperatorCount:    maxOperatorCount,
		StakerPrivateKeys:   stakers,
		StakerTokenAmounts:  stakes,
		OperatorPrivateKeys: operators,
	}

	return deployConfig
}
