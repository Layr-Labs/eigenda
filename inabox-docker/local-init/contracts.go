package localinit

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigenda/core"
)

func (env *Config) GetDeployer(name string) (*ContractDeployer, bool) {
	for _, deployer := range env.Deployers {
		if deployer.Name == name {
			return deployer, true
		}
	}
	return nil, false
}

func (env *Config) DeployEigenDAContracts() {
	log.Print("Deploy the EigenDA and EigenLayer contracts")
	changeDirectory("/contracts")

	// get deployer
	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	if !ok {
		log.Panicf("Deployer improperly configured")
	}

	eigendaDeployConfig := env.GenerateEigenDADeployConfig()
	data, err := json.Marshal(&eigendaDeployConfig)
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}
	writeFile("script/eigenda_deploy_config.json", data)

	execForgeScript("script/SetUpEigenDA.s.sol:SetupEigenDA", env.Pks.EcdsaMap[deployer.Name].PrivateKey, deployer, nil)

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
		env.Pks.EcdsaMap[deployer.Name].PrivateKey,
		deployer,
		[]string{"--sig", "run(address,bytes32)", env.EigenDA.ServiceManager, hashStr},
	)

	//add rollup address to path
	data = readFile("script/output/mock_rollup_deploy_output.json")
	var rollupAddr struct{ MockRollup string }
	err = json.Unmarshal(data, &rollupAddr)
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}

	env.MockRollup = rollupAddr.MockRollup
	changeDirectory("/")
}

func (env *Config) GenerateEigenDADeployConfig() EigenDADeployConfig {

	operators := make([]string, 0)
	stakers := make([]string, 0)
	maxOperatorCount := env.Services.Counts.NumMaxOperatorCount

	total := float32(0)
	stakes := [][]string{make([]string, len(env.Services.Stakes.Distribution))}

	for _, stake := range env.Services.Stakes.Distribution {
		total += stake
	}

	for ind, stake := range env.Services.Stakes.Distribution {
		stakes[0][ind] = strconv.FormatFloat(float64(stake/total*env.Services.Stakes.Total), 'f', 0, 32)
	}

	for i := 0; i < len(env.Services.Stakes.Distribution); i++ {
		stakerName := fmt.Sprintf("staker%d", i)
		operatorName := fmt.Sprintf("opr%d", i)

		stakers = append(stakers, env.getKeyString(stakerName))
		operators = append(operators, env.getKeyString(operatorName))
	}

	config := EigenDADeployConfig{
		UseDefaults:         true,
		NumStrategies:       1,
		MaxOperatorCount:    maxOperatorCount,
		StakerPrivateKeys:   stakers,
		StakerTokenAmounts:  stakes,
		OperatorPrivateKeys: operators,
	}

	return config
}

// Constructs a mapping between service names/deployer names (e.g., 'dis0', 'opr1') and private keys. Order of priority: Map, List, File
func (env *Config) LoadPrivateKeys() error {

	// construct full list of names
	// nTotal := env.Services.Counts.NumDis + env.Services.Counts.NumOpr + env.Services.Counts.NumRet + env.Services.Counts.NumSeq + env.Services.Counts.NumCha
	// names := make([]string, len(env.Deployers)+nTotal)
	names := make([]string, 0)
	for _, d := range env.Deployers {
		names = append(names, d.Name)
	}
	addNames := func(prefix string, num int) {
		for i := 0; i < num; i++ {
			names = append(names, fmt.Sprintf("%v%v", prefix, i))
		}
	}
	addNames("dis", env.Services.Counts.NumDis)
	addNames("opr", env.Services.Counts.NumOpr)
	addNames("staker", env.Services.Counts.NumOpr)

	log.Println("service names:", names)

	// Collect private keys from file
	keyPath := "/secrets"

	// Read ECDSA private keys
	fileData := readFile(filepath.Join(keyPath, "ecdsa_keys/private_key_hex.txt"))
	ecdsaPks := strings.Split(string(fileData), "\n")
	// Read ECDSA passwords
	fileData = readFile(filepath.Join(keyPath, "ecdsa_keys/password.txt"))
	ecdsaPwds := strings.Split(string(fileData), "\n")
	// Read BLS private keys
	fileData = readFile(filepath.Join(keyPath, "bls_keys/private_key_hex.txt"))
	blsPks := strings.Split(string(fileData), "\n")
	// Read BLS passwords
	fileData = readFile(filepath.Join(keyPath, "bls_keys/password.txt"))
	blsPwds := strings.Split(string(fileData), "\n")

	if len(ecdsaPks) != len(blsPks) || len(blsPks) != len(ecdsaPwds) || len(ecdsaPwds) != len(blsPwds) {
		return errors.New("the number of keys and passwords for ECDSA and BLS must be the same")
	}

	// Add missing items to map
	if env.Pks.EcdsaMap == nil {
		env.Pks.EcdsaMap = make(map[string]KeyInfo)
	}
	if env.Pks.BlsMap == nil {
		env.Pks.BlsMap = make(map[string]KeyInfo)
	}

	ind := 0
	for _, name := range names {
		_, exists := env.Pks.EcdsaMap[name]
		if !exists {

			if ind >= len(ecdsaPks) {
				return errors.New("not enough pks")
			}

			env.Pks.EcdsaMap[name] = KeyInfo{
				PrivateKey: ecdsaPks[ind],
				Password:   ecdsaPwds[ind],
				KeyFile:    fmt.Sprintf("%s/ecdsa_keys/keys/%v.ecdsa.key.json", keyPath, ind+1),
			}
			env.Pks.BlsMap[name] = KeyInfo{
				PrivateKey: blsPks[ind],
				Password:   blsPwds[ind],
				KeyFile:    fmt.Sprintf("%s/bls_keys/keys/%v.bls.key.json", keyPath, ind+1),
			}

			ind++
		}
	}

	return nil
}
