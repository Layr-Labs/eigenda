package ethinitscript

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/inabox/strategies/containers/config"
)

func DeployEigenDAContracts(deployerKey string, deployer *config.ContractDeployer, eigendaDeployConfig *config.EigenDADeployConfig, serviceManagerAddr string) {
	log.Print("Deploy the EigenDA and EigenLayer contracts")
	changeDirectory("/contracts")

	data, err := json.Marshal(eigendaDeployConfig)
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}
	writeFile("script/eigenda_deploy_config.json", data)

	execForgeScript("script/SetUpEigenDA.s.sol:SetupEigenDA", deployerKey, deployer, nil)

	blobHeader := &core.BlobHeader{
		QuorumInfos: []*core.BlobQuorumInfo{
			{
				SecurityParam: core.SecurityParam{
					QuorumID:           0,
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
		deployer,
		[]string{"--sig", "run(address,bytes32,uint256)", serviceManagerAddr, hashStr, big.NewInt(1e18).String()},
	)

	//add rollup address to path
	data = readFile("script/output/mock_rollup_deploy_output.json")
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
