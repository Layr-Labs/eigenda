package deploy

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Staker struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private"`
	Stake      string `json:"stake"`
}

// Docker compose
type testbed struct {
	Services map[string]map[string]interface{} `yaml:"services"`
}

type Service struct {
	Image         string   `yaml:"image"`
	Volumes       []string `yaml:"volumes"`
	Ports         []string `yaml:"ports"`
	EnvFile       []string `yaml:"env_file"`
	Command       []string `yaml:"command"`
	ContainerName string   `yaml:"container_name"`
}

type EnvList map[string]string

type ContractDeployer struct {
	Name            string `yaml:"name"`
	RPC             string `yaml:"rpc"`
	VerifierURL     string `yaml:"verifierUrl"`
	VerifyContracts bool   `yaml:"verifyContracts"`
	Slow            bool   `yaml:"slow"`
	DeploySubgraphs bool   `yaml:"deploySubgraphs"`
	// PrivateKey string `yaml:"private_key"`
}

type TelemetryConfig struct {
	IsNeeded   bool     `yaml:"isNeeded"`
	ConfigPath string   `yaml:"configPath"`
	DockerSd   []string `yaml:"dockerSd"`
}

type EigenDADeployConfig struct {
	UseDefaults         bool       `json:"useDefaults"`
	NumStrategies       int        `json:"numStrategies"`
	MaxOperatorCount    int        `json:"maxOperatorCount"`
	StakerPrivateKeys   []string   `json:"stakerPrivateKeys"`
	ConfirmerPrivateKey string     `json:"confirmerPrivateKey"`
	StakerTokenAmounts  [][]string `json:"-"`
	OperatorPrivateKeys []string   `json:"-"`
}

func (cfg *EigenDADeployConfig) MarshalJSON() ([]byte, error) {
	// Convert StakerTokenAmounts to custom string format without quotes
	amountsStr := "["
	for i, subAmounts := range cfg.StakerTokenAmounts {
		amountsStr += "[" + strings.Join(subAmounts, ",") + "]"
		if i < len(cfg.StakerTokenAmounts)-1 {
			amountsStr += ","
		}
	}
	amountsStr += "]"

	operatorPrivateKyesStr := "["
	for i, key := range cfg.OperatorPrivateKeys {
		operatorPrivateKyesStr += "\"" + key + "\""
		if i < len(cfg.OperatorPrivateKeys)-1 {
			operatorPrivateKyesStr += ","
		}
	}
	operatorPrivateKyesStr += "]"

	// Marshal the remaining fields
	remainingFields := map[string]interface{}{
		"useDefaults":         cfg.UseDefaults,
		"numStrategies":       cfg.NumStrategies,
		"maxOperatorCount":    cfg.MaxOperatorCount,
		"stakerPrivateKeys":   cfg.StakerPrivateKeys,
		"confirmerPrivateKey": cfg.ConfirmerPrivateKey,
	}

	remainingJSON, err := json.Marshal(remainingFields)
	if err != nil {
		return nil, err
	}

	// Remove the trailing } from the remaining JSON and append the custom StakerTokenAmounts
	customJSON := string(remainingJSON)[:len(remainingJSON)-1] + `,"stakerTokenAmounts":` + amountsStr + `,"operatorPrivateKeys":` + operatorPrivateKyesStr + "}"
	return []byte(customJSON), nil
}

type EigenDAContract struct {
	Deployer               string `yaml:"deployer"`
	ServiceManager         string `json:"eigenDAServiceManager"`
	OperatorStateRetreiver string `json:"operatorStateRetriever"`
	BlsApkRegistry         string `json:"blsApkRegistry"`
	RegistryCoordinator    string `json:"registryCoordinator"`
	BlobVerifier           string `json:"blobVerifier"`
}

type Stakes struct {
	Total        float32   `yaml:"total"`
	Distribution []float32 `yaml:"distribution"`
}

type ServicesSpec struct {
	Counts struct {
		NumOpr              int `yaml:"operators"`
		NumMaxOperatorCount int `yaml:"maxOperatorCount"`
		NumRelays           int `yaml:"relays"`
	} `yaml:"counts"`
	Stakes    []Stakes  `yaml:"stakes"`
	BasePort  int       `yaml:"basePort"`
	Variables Variables `yaml:"variables"`
}

type Variables map[string]map[string]string

type KeyInfo struct {
	// The private key (e.g. ECDSA or BLS) in string.
	PrivateKey string `yaml:"privateKey"`
	// The password used to encrypt the private key.
	Password string `yaml:"password"`
	// The file path to the encrypted private key.
	KeyFile string `yaml:"keyFile"`
}

type BlobVersionParam struct {
	CodingRate      uint32 `yaml:"codingRate"`
	MaxNumOperators uint32 `yaml:"maxNumOperators"`
	NumChunks       uint32 `yaml:"numChunks"`
}

type PkConfig struct {
	EcdsaMap map[string]KeyInfo `yaml:"ecdsaMap"`
	BlsMap   map[string]KeyInfo `yaml:"blsMap"`
}

type Environment struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func (e Environment) IsLocal() bool {
	return e.Type == "local"
}

type Config struct {
	rootPath string

	Path     string
	TestName string

	Environment Environment `yaml:"environment"`

	Deployers []*ContractDeployer `yaml:"deployers"`

	EigenDA           EigenDAContract     `yaml:"eigenda"`
	BlobVersionParams []*BlobVersionParam `yaml:"blobVersions"`
	MockRollup        string              `yaml:"mockRollup" json:"mockRollup"`

	Pks *PkConfig `yaml:"privateKeys"`

	Services ServicesSpec `yaml:"services"`

	Telemetry TelemetryConfig `yaml:"telemetry"`

	Churner    ChurnerVars
	Dispersers []DisperserVars
	Batcher    []BatcherVars
	Encoder    []EncoderVars
	Operators  []OperatorVars
	Stakers    []Staker
	Retriever  RetrieverVars
	Controller ControllerVars
	Relays     []RelayVars
}

func (c Config) IsEigenDADeployed() bool {
	return c.EigenDA.ServiceManager != ""
}

func NewTestConfig(testName, rootPath string) (testEnv *Config) {

	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		log.Panicf("Error %s:", err.Error())
	}

	testPath := filepath.Join(rootPath, "inabox/testdata/"+testName)

	configPath := testPath + "/config.lock.yaml"
	if _, err := os.Stat(configPath); err != nil {
		configPath = testPath + "/config.yaml"

	}
	data := readFile(configPath)

	err = yaml.Unmarshal(data, &testEnv)
	if err != nil {
		log.Panicf("Error %s:", err.Error())
	}
	testEnv.TestName = testName
	testEnv.Path = testPath
	testEnv.rootPath = rootPath

	return
}

func (env *Config) SaveTestConfig() {
	obj, _ := yaml.Marshal(env)
	writeFile(env.Path+"/config.lock.yaml", obj)
}
