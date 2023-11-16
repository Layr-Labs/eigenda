package config

import (
	"encoding/json"
	"strings"
)

type Staker struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private"`
	Stake      string `json:"stake"`
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
		"useDefaults":       cfg.UseDefaults,
		"numStrategies":     cfg.NumStrategies,
		"maxOperatorCount":  cfg.MaxOperatorCount,
		"stakerPrivateKeys": cfg.StakerPrivateKeys,
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
	Deployer                       string `yaml:"deployer"`
	ServiceManager                 string `yaml:"eigenDAServiceManager"`
	OperatorStateRetreiver         string `yaml:"blsOperatorStateRetriever"`
	PubkeyRegistry                 string `yaml:"blsPubkeyRegistry"`
	PubkeyCompendium               string `yaml:"pubkeyCompendium"`
	RegistryCoordinatorWithIndices string `yaml:"blsRegistryCoordinatorWithIndices"`
}

type ServicesSpec struct {
	Counts struct {
		NumDis              int `yaml:"dispersers"`
		NumOpr              int `yaml:"operators"`
		NumMaxOperatorCount int `yaml:"maxOperatorCount"`
	} `yaml:"counts"`
	Stakes struct {
		Total        float32   `yaml:"total"`
		Distribution []float32 `yaml:"distribution"`
	} `yaml:"stakes"`
	BasePort  int       `yaml:"basePort"`
	Variables Variables `yaml:"variables"`
}

type Variables map[string]map[string]string

type Environment struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func (e Environment) IsLocal() bool {
	return e.Type == "local"
}

type Config struct {
	Environment Environment         `yaml:"environment"`
	Deployers   []*ContractDeployer `yaml:"deployers"`
	EigenDA     EigenDAContract     `yaml:"eigenda"`
	MockRollup  string              `yaml:"mockRollup" json:"mockRollup"`
	Services    ServicesSpec        `yaml:"services"`
	Telemetry   TelemetryConfig     `yaml:"telemetry"`
	Pks         PkConfig            `yaml:"privateKeys"`
}
