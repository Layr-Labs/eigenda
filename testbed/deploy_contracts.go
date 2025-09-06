package testbed

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// KeyInfo represents information about a private key
type KeyInfo struct {
	PrivateKey string
	Password   string
	KeyFile    string
}

// PrivateKeyMaps holds the ECDSA and BLS key mappings
type PrivateKeyMaps struct {
	EcdsaMap map[string]KeyInfo
	BlsMap   map[string]KeyInfo
}

// LoadPrivateKeysInput contains all the inputs needed to load private keys
type LoadPrivateKeysInput struct {
	NumOperators int
	NumRelays    int
}

// GetAnvilDefaultKeys returns the default private keys from Anvil's test mnemonic
// These keys are from: "test test test test test test test test test junk"
func GetAnvilDefaultKeys() (defaultKey string, batcher0Key string) {
	// Account #0: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 (10,000 ETH)
	defaultKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	// Account #1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 (10,000 ETH)
	batcher0Key = "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

	return defaultKey, batcher0Key
}

// LoadPrivateKeys constructs a mapping between service names (e.g., 'deployer', 'dis0', 'opr1') and private keys
func LoadPrivateKeys(input LoadPrivateKeysInput) (*PrivateKeyMaps, error) {
	// Get funded Anvil keys for deployer and batcher
	deployerKey, batcherKey := GetAnvilDefaultKeys()

	// Use testbed secrets for other services
	keyPath := "secrets"

	// Build the list of service names
	names := make([]string, 0)

	// Add single deployer
	names = append(names, "deployer")

	addNames := func(prefix string, num int) {
		for i := 0; i < num; i++ {
			names = append(names, fmt.Sprintf("%v%v", prefix, i))
		}
	}
	addNames("dis", 2)
	addNames("opr", input.NumOperators)
	addNames("staker", input.NumOperators)
	addNames("retriever", 1)
	addNames("relay", input.NumRelays)

	// Read ECDSA private keys from secrets
	fileData, err := os.ReadFile(filepath.Join(keyPath, "ecdsa_keys/private_key_hex.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to read ECDSA private keys: %w", err)
	}
	ecdsaPks := strings.Split(string(fileData), "\n")

	// Read ECDSA passwords
	fileData, err = os.ReadFile(filepath.Join(keyPath, "ecdsa_keys/password.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to read ECDSA passwords: %w", err)
	}
	ecdsaPwds := strings.Split(string(fileData), "\n")

	// Read BLS private keys
	fileData, err = os.ReadFile(filepath.Join(keyPath, "bls_keys/private_key_hex.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to read BLS private keys: %w", err)
	}
	blsPks := strings.Split(string(fileData), "\n")

	// Read BLS passwords
	fileData, err = os.ReadFile(filepath.Join(keyPath, "bls_keys/password.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to read BLS passwords: %w", err)
	}
	blsPwds := strings.Split(string(fileData), "\n")

	if len(ecdsaPks) != len(blsPks) || len(blsPks) != len(ecdsaPwds) || len(ecdsaPwds) != len(blsPwds) {
		return nil, errors.New("the number of keys and passwords for ECDSA and BLS must be the same")
	}

	// Initialize maps
	result := &PrivateKeyMaps{
		EcdsaMap: make(map[string]KeyInfo),
		BlsMap:   make(map[string]KeyInfo),
	}

	// Add keys for each service name
	// Start at index 0 for reading from secrets (we'll skip indices for deployer and dis0)
	secretIndex := 0
	for _, name := range names {
		switch name {
		case "deployer":
			// Deployer uses Anvil account #0
			result.EcdsaMap[name] = KeyInfo{
				PrivateKey: deployerKey,
				Password:   "",
				KeyFile:    "",
			}
			// No BLS key for deployer
			result.BlsMap[name] = KeyInfo{
				PrivateKey: "",
				Password:   "",
				KeyFile:    "",
			}
		case "dis0":
			// First disperser (batcher) uses Anvil account #1
			result.EcdsaMap[name] = KeyInfo{
				PrivateKey: batcherKey,
				Password:   "",
				KeyFile:    "",
			}
			// No BLS key for batcher
			result.BlsMap[name] = KeyInfo{
				PrivateKey: "",
				Password:   "",
				KeyFile:    "",
			}
		default:
			// All other services use keys from secrets
			if secretIndex >= len(ecdsaPks) {
				return nil, errors.New("not enough keys in secrets")
			}

			result.EcdsaMap[name] = KeyInfo{
				PrivateKey: ecdsaPks[secretIndex],
				Password:   ecdsaPwds[secretIndex],
				KeyFile:    fmt.Sprintf("%s/ecdsa_keys/keys/%v.ecdsa.key.json", keyPath, secretIndex+1),
			}
			result.BlsMap[name] = KeyInfo{
				PrivateKey: blsPks[secretIndex],
				Password:   blsPwds[secretIndex],
				KeyFile:    fmt.Sprintf("%s/bls_keys/keys/%v.bls.key.json", keyPath, secretIndex+1),
			}

			secretIndex++
		}
	}

	return result, nil
}

// EigenDADeployConfig contains configuration for deploying EigenDA contracts
type EigenDADeployConfig struct {
	UseDefaults         bool       `json:"useDefaults"`
	NumStrategies       int        `json:"numStrategies"`
	MaxOperatorCount    int        `json:"maxOperatorCount"`
	StakerPrivateKeys   []string   `json:"stakerPrivateKeys"`
	ConfirmerPrivateKey string     `json:"confirmerPrivateKey"`
	StakerTokenAmounts  [][]string `json:"-"`
	OperatorPrivateKeys []string   `json:"-"`
}

// Custom JSON marshaling for EigenDADeployConfig
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

	operatorPrivateKeysStr := "["
	for i, key := range cfg.OperatorPrivateKeys {
		operatorPrivateKeysStr += "\"" + key + "\""
		if i < len(cfg.OperatorPrivateKeys)-1 {
			operatorPrivateKeysStr += ","
		}
	}
	operatorPrivateKeysStr += "]"

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
		return nil, fmt.Errorf("failed to marshal remaining fields: %w", err)
	}

	// Convert to map to add custom fields
	var result map[string]interface{}
	if err := json.Unmarshal(remainingJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to map: %w", err)
	}

	// Add the custom formatted fields as raw JSON
	result["stakerTokenAmounts"] = json.RawMessage(amountsStr)
	result["operatorPrivateKeys"] = json.RawMessage(operatorPrivateKeysStr)

	finalJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal final result: %w", err)
	}
	return finalJSON, nil
}

// V1CertVerifierDeployConfig contains configuration for deploying V1 CertVerifier
type V1CertVerifierDeployConfig struct {
	ServiceManager                string   `json:"eigenDAServiceManager"`
	RequiredQuorums               []uint32 `json:"requiredQuorums"`
	RequiredAdversarialThresholds []uint32 `json:"adversaryThresholds"`
	RequiredConfirmationQuorums   []uint32 `json:"confirmationThresholds"`
}

// EigenDAContract holds deployed EigenDA contract addresses
type EigenDAContract struct {
	EigenDADirectory       string `json:"eigenDADirectory"`
	ServiceManager         string `json:"eigenDAServiceManager"`
	OperatorStateRetriever string `json:"operatorStateRetriever"`
	BlsApkRegistry         string `json:"blsApkRegistry"`
	RegistryCoordinator    string `json:"registryCoordinator"`
	CertVerifierLegacy     string `json:"eigenDALegacyCertVerifier"`
	CertVerifier           string `json:"eigenDACertVerifier"`
	CertVerifierRouter     string `json:"eigenDACertVerifierRouter"`
}

// Stakes represents token staking configuration
type Stakes struct {
	Total        float32   `yaml:"total"`
	Distribution []float32 `yaml:"distribution"`
}

// DeploymentConfig holds all configuration for deploying contracts
type DeploymentConfig struct {
	AnvilRPCURL      string
	DeployerKey      string
	NumOperators     int
	NumRelays        int
	Stakes           []Stakes
	MaxOperatorCount int
	PrivateKeys      *PrivateKeyMaps
	Logger           logging.Logger
}

// DeploymentResult holds the results of contract deployment
type DeploymentResult struct {
	EigenDA               EigenDAContract
	EigenDAV1CertVerifier string
	EigenDAV2CertVerifier string
}

// DeployEigenDAContracts deploys EigenDA core system and along with Eigenlayer contracts on a local anvil chain.
// This calls the SetupEigenDA.s.sol forge script to initialize the deployment.
//
// TODO: SetupEigenDA.s.sol is pretty legacy and its primary function is to set up the EigenDA environment for the inabox environment.
// There exists a DeployEigenDA.s.sol script that has been used in production to deploy environments but it currently does not handle the
// Eigenlayer contracts. We should consider deprecating SetupEigenDA.s.sol in favor of DeployEigenDA.s.sol.
func DeployEigenDAContracts(config DeploymentConfig) (*DeploymentResult, error) {
	if config.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	config.Logger.Info("Deploy the EigenDA and EigenLayer contracts")

	result := &DeploymentResult{}

	// Save current directory and change to contracts
	origDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	defer func() {
		_ = os.Chdir(origDir)
	}()

	contractsDir := "../contracts"
	if err := os.Chdir(contractsDir); err != nil {
		return nil, fmt.Errorf("failed to change to contracts directory: %w", err)
	}

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil {
		config.Logger.Info("Successfully changed to absolute path", "path", cwd)
	}

	eigendaDeployConfig, err := generateEigenDADeployConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error generating eigenda deploy config: %w", err)
	}

	data, err := json.Marshal(&eigendaDeployConfig)
	if err != nil {
		return nil, fmt.Errorf("error marshaling eigenda deploy config: %w", err)
	}
	err = os.WriteFile("script/input/eigenda_deploy_config.json", data, 0644)
	if err != nil {
		return nil, fmt.Errorf("error writing eigenda deploy config: %w", err)
	}

	config.Logger.Info("Executing EigenDA deployer script", "script", "script/SetUpEigenDA.s.sol:SetupEigenDA")
	err = execForgeScript(
		"script/SetUpEigenDA.s.sol:SetupEigenDA",
		config.DeployerKey,
		config.AnvilRPCURL,
		nil,
		config.Logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute EigenDA deployer script: %w", err)
	}

	// Add relevant addresses to path
	data, err = os.ReadFile("script/output/eigenda_deploy_output.json")
	if err != nil {
		return nil, fmt.Errorf("error reading eigenda deploy output: %w", err)
	}

	err = json.Unmarshal(data, &result.EigenDA)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling eigenda deploy output: %w", err)
	}

	// Deploy V1 CertVerifier
	certVerifierV1DeployCfg := generateV1CertVerifierDeployConfig(result.EigenDA.ServiceManager)
	data, err = json.Marshal(&certVerifierV1DeployCfg)
	if err != nil {
		return nil, fmt.Errorf("error marshaling certverifier config: %w", err)
	}

	// NOTE: this is pretty janky and is a short-term solution until V1 contract usage
	//       can be deprecated.
	if err := os.WriteFile("script/deploy/certverifier/config/v1/inabox_deploy_config_v1.json", data, 0644); err != nil {
		return nil, fmt.Errorf("error writing certverifier config: %w", err)
	}

	config.Logger.Info("Executing CertVerifierDeployerV1 script")
	if err := execForgeScript("script/deploy/certverifier/CertVerifierDeployerV1.s.sol:CertVerifierDeployerV1",
		config.DeployerKey,
		config.AnvilRPCURL,
		[]string{"--sig", "run(string, string)", "inabox_deploy_config_v1.json", "inabox_v1_deploy.json"},
		config.Logger); err != nil {
		return nil, fmt.Errorf("failed to execute CertVerifierDeployerV1 script: %w", err)
	}

	data, err = os.ReadFile("script/deploy/certverifier/output/inabox_v1_deploy.json")
	if err != nil {
		return nil, fmt.Errorf("error reading certverifier output: %w", err)
	}

	var verifierAddress struct{ EigenDACertVerifier string }
	err = json.Unmarshal(data, &verifierAddress)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling verifier address: %w", err)
	}
	result.EigenDAV1CertVerifier = verifierAddress.EigenDACertVerifier

	config.Logger.Debug("Deployment results",
		"EigenDADirectory", result.EigenDA.EigenDADirectory,
		"ServiceManager", result.EigenDA.ServiceManager,
		"OperatorStateRetriever", result.EigenDA.OperatorStateRetriever,
		"BlsApkRegistry", result.EigenDA.BlsApkRegistry,
		"RegistryCoordinator", result.EigenDA.RegistryCoordinator,
		"CertVerifierLegacy", result.EigenDA.CertVerifierLegacy,
		"CertVerifier", result.EigenDA.CertVerifier,
		"CertVerifierRouter", result.EigenDA.CertVerifierRouter,
		"V1CertVerifier", result.EigenDAV1CertVerifier,
	)

	return result, nil
}

// execForgeScript executes a forge script with the given parameters
func execForgeScript(script, privateKey, rpcURL string, extraArgs []string, logger logging.Logger) error {
	args := []string{"script", script,
		"--rpc-url", rpcURL,
		"--private-key", privateKey,
		"--broadcast"}

	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	cmd := exec.Command("forge", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Info("Running forge command", "command", "forge "+strings.Join(args, " "))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("forge script failed: %w", err)
	}

	return nil
}

// generateEigenDADeployConfig generates input config fed into SetUpEigenDA.s.sol foundry script
func generateEigenDADeployConfig(config DeploymentConfig) (EigenDADeployConfig, error) {
	operators := make([]string, 0)
	stakers := make([]string, 0)
	maxOperatorCount := config.MaxOperatorCount
	if maxOperatorCount == 0 {
		maxOperatorCount = config.NumOperators
	}

	numStrategies := len(config.Stakes)
	if numStrategies == 0 {
		// Default to 2 strategies if not specified
		numStrategies = 2
		config.Stakes = []Stakes{
			{Total: 1e18, Distribution: make([]float32, config.NumOperators)},
			{Total: 1e18, Distribution: make([]float32, config.NumOperators)},
		}
		// Equal distribution
		for i := 0; i < config.NumOperators; i++ {
			config.Stakes[0].Distribution[i] = 1.0 / float32(config.NumOperators)
			config.Stakes[1].Distribution[i] = 1.0 / float32(config.NumOperators)
		}
	}

	total := make([]float32, numStrategies)
	stakes := make([][]string, numStrategies)

	for quorum, stake := range config.Stakes {
		for _, s := range stake.Distribution {
			total[quorum] += s
		}
	}

	for quorum := 0; quorum < numStrategies; quorum++ {
		stakes[quorum] = make([]string, len(config.Stakes[quorum].Distribution))
		for ind, stake := range config.Stakes[quorum].Distribution {
			stakes[quorum][ind] = strconv.FormatFloat(float64(stake/total[quorum]*config.Stakes[quorum].Total), 'f', 0, 32)
		}
	}

	for i := 0; i < config.NumOperators; i++ {
		stakerName := fmt.Sprintf("staker%d", i)
		operatorName := fmt.Sprintf("opr%d", i)

		// Get keys for staker and operator
		stakerKey, ok := config.PrivateKeys.EcdsaMap[stakerName]
		if !ok {
			return EigenDADeployConfig{}, fmt.Errorf("failed to get key for %s", stakerName)
		}

		operatorKey, ok := config.PrivateKeys.EcdsaMap[operatorName]
		if !ok {
			return EigenDADeployConfig{}, fmt.Errorf("failed to get key for %s", operatorName)
		}

		stakers = append(stakers, stakerKey.PrivateKey)
		operators = append(operators, operatorKey.PrivateKey)
	}

	// Use batcher0 key as the batch confirmer
	batcherKeyInfo, ok := config.PrivateKeys.EcdsaMap["batcher0"]
	if !ok {
		return EigenDADeployConfig{}, fmt.Errorf("failed to get key for batcher0")
	}
	batcherKey := batcherKeyInfo.PrivateKey

	deployConfig := EigenDADeployConfig{
		UseDefaults:         true,
		NumStrategies:       numStrategies,
		MaxOperatorCount:    maxOperatorCount,
		StakerPrivateKeys:   stakers,
		StakerTokenAmounts:  stakes,
		OperatorPrivateKeys: operators,
		ConfirmerPrivateKey: batcherKey,
	}

	return deployConfig, nil
}

// generateV1CertVerifierDeployConfig generates the input config used for deploying the V1 CertVerifier
// NOTE: this will be killed in the future with eventual deprecation of V1
func generateV1CertVerifierDeployConfig(serviceManager string) V1CertVerifierDeployConfig {
	config := V1CertVerifierDeployConfig{
		ServiceManager:                serviceManager,
		RequiredQuorums:               []uint32{0, 1},
		RequiredAdversarialThresholds: []uint32{33, 33},
		RequiredConfirmationQuorums:   []uint32{55, 55},
	}

	return config
}

// DeployToAnvil is a convenience function to deploy contracts to an Anvil instance
func DeployContractsToAnvil(anvilURL string, numOperators int, logger logging.Logger) (*DeploymentResult, error) {
	// Load private keys
	privateKeys, err := LoadPrivateKeys(LoadPrivateKeysInput{
		NumOperators: numOperators,
		NumRelays:    1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load private keys: %w", err)
	}

	// Use Anvil's default first account which has 10,000 ETH
	// This is needed because the deployment script needs to fund staker accounts
	anvilDefaultKey, _ := GetAnvilDefaultKeys()

	// Deploy contracts
	config := DeploymentConfig{
		AnvilRPCURL:      anvilURL,
		DeployerKey:      anvilDefaultKey,
		NumOperators:     numOperators,
		NumRelays:        1,
		MaxOperatorCount: numOperators,
		PrivateKeys:      privateKeys,
		Logger:           logger,
	}

	return DeployEigenDAContracts(config)
}
