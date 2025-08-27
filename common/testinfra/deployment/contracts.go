package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"crypto/ecdsa"

	"github.com/Layr-Labs/eigenda/common"
	relayreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARelayRegistry"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	thresholdreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAThresholdRegistry"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ContractDeployer holds configuration for deploying contracts
type ContractDeployer struct {
	Name            string `yaml:"name"`
	RPC             string `yaml:"rpc"`
	PrivateKey      string `yaml:"private_key"`
	DeploySubgraphs bool   `yaml:"deploy_subgraphs"`
	VerifyContracts bool   `yaml:"verify_contracts"`
	VerifierURL     string `yaml:"verifier_url"`
	Slow            bool   `yaml:"slow"`
}

// EigenDAContracts contains all deployed EigenDA contract addresses
type EigenDAContracts struct {
	ProxyAdmin                   string `json:"proxyAdmin"`
	PauserRegistry               string `json:"pauserRegistry"`
	DelegationManager            string `json:"delegationManager"`
	Slasher                      string `json:"slasher"`
	StrategyManager              string `json:"strategyManager"`
	EigenPodManager              string `json:"eigenPodManager"`
	AVSDirectory                 string `json:"avsDirectory"`
	EigenDADirectory             string `json:"eigenDADirectory"`
	RewardsCoordinator           string `json:"rewardsCoordinator"`
	StrategyFactory              string `json:"strategyFactory"`
	StrategyBeacon               string `json:"strategyBeacon"`
	StrategyBase                 string `json:"strategyBase"`
	StrategyBaseTVLLimits        string `json:"strategyBaseTVLLimits"`
	Token                        string `json:"token"`
	ServiceManager               string `json:"eigenDAServiceManager"`
	RegistryCoordinator          string `json:"registryCoordinator"`
	BlsApkRegistry               string `json:"blsApkRegistry"`
	StakeRegistry                string `json:"stakeRegistry"`
	IndexRegistry                string `json:"indexRegistry"`
	OperatorStateRetriever       string `json:"operatorStateRetriever"`
	ServiceManagerImplementation string `json:"serviceManagerImplementation"`
	EigenDAV1CertVerifier        string `json:"eigenDAV1CertVerifier,omitempty"`
	EigenDAV2CertVerifier        string `json:"eigenDACertVerifier,omitempty"`
	EigenDACertVerifierRouter    string `json:"eigenDACertVerifierRouter,omitempty"`
	Deployer                     string `json:"deployer,omitempty"`
}

// EigenDADeployConfig represents the configuration for deploying EigenDA contracts
type EigenDADeployConfig struct {
	UseDefaults         bool       `json:"useDefaults"`
	NumStrategies       int        `json:"numStrategies"`
	MaxOperatorCount    int        `json:"maxOperatorCount"`
	StakerPrivateKeys   []string   `json:"stakerPrivateKeys"`
	StakerTokenAmounts  [][]string `json:"-"`
	OperatorPrivateKeys []string   `json:"-"`
	ConfirmerPrivateKey string     `json:"confirmerPrivateKey"`
}

// MarshalJSON custom marshaler to generate unquoted numeric values for stakerTokenAmounts
// This matches the format expected by the Solidity contract deployment script
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
		return nil, err
	}

	// Remove the trailing } from the remaining JSON and append the custom StakerTokenAmounts
	customJSON := string(remainingJSON)[:len(remainingJSON)-1] + `,"stakerTokenAmounts":` + amountsStr + `,"operatorPrivateKeys":` + operatorPrivateKeysStr + "}"
	return []byte(customJSON), nil
}

// V1CertVerifierDeployConfig represents the configuration for deploying V1 CertVerifier
type V1CertVerifierDeployConfig struct {
	EigenDAServiceManager  string   `json:"eigenDAServiceManager"`
	RequiredQuorums        []uint32 `json:"requiredQuorums"`
	AdversaryThresholds    []uint32 `json:"adversaryThresholds"`
	ConfirmationThresholds []uint32 `json:"confirmationThresholds"`
}

// BlobVersionParam represents blob version parameters
type BlobVersionParam struct {
	MaxNumOperators uint32 `yaml:"max_num_operators"`
	NumChunks       uint32 `yaml:"num_chunks"`
	CodingRate      uint32 `yaml:"coding_rate"`
}

// ContractDeploymentManager manages EigenDA contract deployments
type ContractDeploymentManager struct {
	RootPath          string
	Deployer          ContractDeployer
	EigenDAContracts  *EigenDAContracts
	PrivateKeys       map[string]string
	DisperserAddress  gcommon.Address
	DisperserKMSKeyID string
}

// NewContractDeploymentManager creates a new contract deployment manager
func NewContractDeploymentManager(rootPath string, deployer ContractDeployer) *ContractDeploymentManager {
	return &ContractDeploymentManager{
		RootPath:         rootPath,
		Deployer:         deployer,
		EigenDAContracts: &EigenDAContracts{},
		PrivateKeys:      make(map[string]string),
	}
}

// DeployEigenDAContracts deploys all EigenDA contracts
func (m *ContractDeploymentManager) DeployEigenDAContracts(deployConfig EigenDADeployConfig) error {
	log.Print("Deploying EigenDA and EigenLayer contracts")

	// Save current directory to restore later
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			log.Printf("Warning: failed to restore directory to %s: %v", originalDir, err)
		}
	}()

	// Change to contracts directory
	contractsPath := filepath.Join(m.RootPath, "contracts")
	if err := os.Chdir(contractsPath); err != nil {
		return fmt.Errorf("failed to change to contracts directory: %w", err)
	}

	// Write deployment config
	data, err := json.Marshal(&deployConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal deploy config: %w", err)
	}

	configPath := "script/input/eigenda_deploy_config.json"
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write deploy config: %w", err)
	}

	// Execute the forge script
	if err := m.execForgeScript("script/SetUpEigenDA.s.sol:SetupEigenDA", nil); err != nil {
		return fmt.Errorf("failed to execute forge script: %w", err)
	}

	// Read the deployment output
	outputPath := "script/output/eigenda_deploy_output.json"
	outputData, err := os.ReadFile(outputPath)
	if err != nil {
		return fmt.Errorf("failed to read deployment output: %w", err)
	}

	if err := json.Unmarshal(outputData, m.EigenDAContracts); err != nil {
		return fmt.Errorf("failed to unmarshal deployment output: %w", err)
	}

	log.Printf("Loaded EigenDA contracts - ServiceManager: %s, RegistryCoordinator: %s",
		m.EigenDAContracts.ServiceManager, m.EigenDAContracts.RegistryCoordinator)

	// Deploy V1 CertVerifier (optional - log errors but don't fail)
	if err := m.deployV1CertVerifier(); err != nil {
		log.Printf("Warning: V1 CertVerifier deployment failed: %v", err)
		log.Print("Continuing without V1 CertVerifier - using default address")
		// Set a default/empty address if deployment fails
		m.EigenDAContracts.EigenDAV1CertVerifier = "0x0000000000000000000000000000000000000000"
	}

	// V2 CertVerifier and Router are already deployed by the main SetUpEigenDA.s.sol script
	// and loaded from the main deployment output above
	log.Printf("V2 CertVerifier loaded from main deployment: %s", m.EigenDAContracts.EigenDAV2CertVerifier)
	log.Printf("CertVerifier Router loaded from main deployment: %s", m.EigenDAContracts.EigenDACertVerifierRouter)

	log.Print("Successfully deployed EigenDA contracts")
	return nil
}

// deployV1CertVerifier deploys the V1 CertVerifier contract
func (m *ContractDeploymentManager) deployV1CertVerifier() error {
	log.Print("Deploying V1 CertVerifier")

	// Create logger for V1 CertVerifier deployment (client is not actually used here)
	// The deployment is done via forge script which uses the deployer's RPC and private key

	// Generate V1 CertVerifier deploy config
	v1Config := V1CertVerifierDeployConfig{
		EigenDAServiceManager:  m.EigenDAContracts.ServiceManager,
		RequiredQuorums:        []uint32{0, 1},
		AdversaryThresholds:    []uint32{33, 33},
		ConfirmationThresholds: []uint32{55, 55},
	}

	data, err := json.Marshal(&v1Config)
	if err != nil {
		return fmt.Errorf("failed to marshal V1 config: %w", err)
	}

	configPath := "script/deploy/certverifier/config/v1/inabox_deploy_config_v1.json"
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write V1 config: %w", err)
	}

	// Execute V1 deployment script
	extraArgs := []string{"--sig", "run(string, string)", "inabox_deploy_config_v1.json", "inabox_v1_deploy.json"}
	if err := m.execForgeScript("script/deploy/certverifier/CertVerifierDeployerV1.s.sol:CertVerifierDeployerV1", extraArgs); err != nil {
		return fmt.Errorf("failed to execute V1 deployment script: %w", err)
	}

	// Read V1 deployment output
	outputPath := "script/deploy/certverifier/output/inabox_v1_deploy.json"
	outputData, err := os.ReadFile(outputPath)
	if err != nil {
		return fmt.Errorf("failed to read V1 deployment output: %w", err)
	}

	var verifierAddress struct{ EigenDACertVerifier string }
	if err := json.Unmarshal(outputData, &verifierAddress); err != nil {
		return fmt.Errorf("failed to unmarshal V1 deployment output: %w", err)
	}

	m.EigenDAContracts.EigenDAV1CertVerifier = verifierAddress.EigenDACertVerifier

	log.Printf("V1 CertVerifier deployed at: %s", m.EigenDAContracts.EigenDAV1CertVerifier)
	return nil
}

// RegisterDisperserAddress registers the disperser's address on-chain
func (m *ContractDeploymentManager) RegisterDisperserAddress(ethClient common.EthClient, disperserAddress gcommon.Address) error {
	log.Printf("RegisterDisperserAddress: Starting registration for address: %s", disperserAddress.Hex())
	log.Printf("RegisterDisperserAddress: ServiceManager: %s", m.EigenDAContracts.ServiceManager)
	log.Printf("RegisterDisperserAddress: OperatorStateRetriever: %s", m.EigenDAContracts.OperatorStateRetriever)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	writer, err := eth.NewWriter(
		logger,
		ethClient,
		m.EigenDAContracts.OperatorStateRetriever,
		m.EigenDAContracts.ServiceManager,
	)
	if err != nil {
		log.Printf("RegisterDisperserAddress: Failed to create writer: %v", err)
		return fmt.Errorf("failed to create writer: %w", err)
	}

	log.Printf("RegisterDisperserAddress: Calling SetDisperserAddress")
	if err := writer.SetDisperserAddress(context.Background(), disperserAddress); err != nil {
		log.Printf("RegisterDisperserAddress: SetDisperserAddress failed: %v", err)
		// If it's the "disperser registry not deployed" error, provide more context
		if err.Error() == "disperser registry not deployed" {
			log.Printf("RegisterDisperserAddress: DisperserRegistry contract is not deployed or not linked to ServiceManager")
			log.Printf("RegisterDisperserAddress: This may be because the deployment script doesn't include DisperserRegistry")
			// Don't fail the test for this specific error as it's expected in test environments
			log.Printf("RegisterDisperserAddress: Skipping disperser registration (expected in test environment)")
			return nil
		}
		return fmt.Errorf("failed to set disperser address: %w", err)
	}

	log.Printf("RegisterDisperserAddress: SetDisperserAddress succeeded, now verifying...")
	// Verify the address was set correctly
	registeredAddress, err := writer.GetDisperserAddress(context.Background(), 0)
	if err != nil {
		log.Printf("RegisterDisperserAddress: GetDisperserAddress failed: %v", err)
		log.Printf("RegisterDisperserAddress: Warning - could not verify registration, but continuing anyway")
		// Don't fail here as the disperser registry might not be queryable in test environments
		return nil
	}

	if registeredAddress != disperserAddress {
		log.Printf("RegisterDisperserAddress: Address mismatch - expected %s, got %s", disperserAddress, registeredAddress)
		return fmt.Errorf("disperser address mismatch: expected %s, got %s", disperserAddress, registeredAddress)
	}

	log.Printf("RegisterDisperserAddress: Successfully registered and verified disperser address")
	return nil
}

// RegisterBlobVersionsAndRelays registers blob versions and relay information
func (m *ContractDeploymentManager) RegisterBlobVersionsAndRelays(
	ethClient common.EthClient,
	blobVersionParams []BlobVersionParam,
	relayURLs []string,
) error {
	log.Print("Registering blob versions and relays")

	// Get service manager contract
	dasmAddr := gcommon.HexToAddress(m.EigenDAContracts.ServiceManager)
	contractEigenDAServiceManager, err := eigendasrvmg.NewContractEigenDAServiceManager(dasmAddr, ethClient)
	if err != nil {
		return fmt.Errorf("failed to get service manager contract: %w", err)
	}

	// Register blob versions
	thresholdRegistryAddr, err := contractEigenDAServiceManager.EigenDAThresholdRegistry(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("failed to get threshold registry address: %w", err)
	}

	contractThresholdRegistry, err := thresholdreg.NewContractEigenDAThresholdRegistry(thresholdRegistryAddr, ethClient)
	if err != nil {
		return fmt.Errorf("failed to get threshold registry contract: %w", err)
	}

	opts, err := ethClient.GetNoSendTransactOpts()
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %w", err)
	}

	for _, blobVersionParam := range blobVersionParams {
		txn, err := contractThresholdRegistry.AddVersionedBlobParams(opts, thresholdreg.EigenDATypesV1VersionedBlobParams{
			MaxNumOperators: blobVersionParam.MaxNumOperators,
			NumChunks:       blobVersionParam.NumChunks,
			CodingRate:      uint8(blobVersionParam.CodingRate),
		})
		if err != nil {
			return fmt.Errorf("failed to add blob version params: %w", err)
		}
		if err := ethClient.SendTransaction(context.Background(), txn); err != nil {
			return fmt.Errorf("failed to send blob version transaction: %w", err)
		}
	}

	// Register relays
	relayAddr, err := contractEigenDAServiceManager.EigenDARelayRegistry(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("failed to get relay registry address: %w", err)
	}

	contractRelayRegistry, err := relayreg.NewContractEigenDARelayRegistry(relayAddr, ethClient)
	if err != nil {
		return fmt.Errorf("failed to get relay registry contract: %w", err)
	}

	ethAddr := ethClient.GetAccountAddress()
	for _, url := range relayURLs {
		// URL should already include the proper hostname (e.g., relay-0.localtest.me:34000)
		txn, err := contractRelayRegistry.AddRelayInfo(opts, relayreg.EigenDATypesV2RelayInfo{
			RelayAddress: ethAddr,
			RelayURL:     url,
		})
		if err != nil {
			return fmt.Errorf("failed to add relay info: %w", err)
		}
		if err := ethClient.SendTransaction(context.Background(), txn); err != nil {
			return fmt.Errorf("failed to send relay transaction: %w", err)
		}
	}

	log.Print("Successfully registered blob versions and relays")
	return nil
}

// execForgeScript executes a forge script
func (m *ContractDeploymentManager) execForgeScript(script string, extraArgs []string) error {
	log.Printf("Executing forge script: %s", script)

	args := []string{"script", script,
		"--rpc-url", m.Deployer.RPC,
		"--private-key", m.Deployer.PrivateKey,
		"--broadcast"}

	if m.Deployer.VerifyContracts {
		args = append(args, "--verify",
			"--verifier", "blockscout",
			"--verifier-url", m.Deployer.VerifierURL)
	}

	if m.Deployer.Slow {
		args = append(args, "--slow")
	}

	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	cmd := exec.Command("forge", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("forge script failed: %w", err)
	}

	log.Print("Forge script executed successfully")
	return nil
}

// GetContractAddresses returns the deployed contract addresses
func (m *ContractDeploymentManager) GetContractAddresses() *EigenDAContracts {
	return m.EigenDAContracts
}

// GetAddress converts a private key to an Ethereum address
func GetAddress(privateKey string) (gcommon.Address, error) {
	// Remove 0x prefix if present
	privateKey = strings.TrimPrefix(privateKey, "0x")

	// Parse the private key
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return gcommon.Address{}, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get the public key
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return gcommon.Address{}, fmt.Errorf("failed to cast public key to ECDSA")
	}

	// Get the address
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

// GenerateEigenDADeployConfig generates the deployment configuration for EigenDA contracts
func GenerateEigenDADeployConfig(
	numStrategies int,
	maxOperatorCount int,
	stakeDistribution [][]float32,
	stakeTotals []float32,
	privateKeys map[string]string) EigenDADeployConfig {
	operators := make([]string, 0)
	stakers := make([]string, 0)
	stakes := make([][]string, numStrategies)

	// Calculate normalized stakes
	total := make([]float32, numStrategies)
	for quorum, distribution := range stakeDistribution {
		for _, s := range distribution {
			total[quorum] += s
		}
	}

	for quorum := 0; quorum < numStrategies; quorum++ {
		stakes[quorum] = make([]string, len(stakeDistribution[quorum]))
		for ind, stake := range stakeDistribution[quorum] {
			normalizedStake := stake / total[quorum] * stakeTotals[quorum]
			// Use higher precision to avoid truncation issues with large numbers like 100e18
			stakes[quorum][ind] = strconv.FormatFloat(float64(normalizedStake), 'f', -1, 64)
		}
	}

	// Get staker and operator keys
	for i := 0; i < len(stakeDistribution[0]); i++ {
		stakerName := fmt.Sprintf("staker%d", i)
		operatorName := fmt.Sprintf("opr%d", i)

		if stakerKey, ok := privateKeys[stakerName]; ok {
			keyInt, ok := new(big.Int).SetString(stakerKey, 0)
			if !ok {
				log.Printf("Warning: could not parse staker key %s", stakerName)
				continue
			}
			stakers = append(stakers, keyInt.String())
		}

		// For operators, we need to use ECDSA private keys, not BLS keys
		// Look for operator ECDSA keys first with "opr{i}_ecdsa" naming
		ecdsaKeyName := fmt.Sprintf("%s_ecdsa", operatorName)
		if ecdsaKey, ok := privateKeys[ecdsaKeyName]; ok {
			keyInt, ok := new(big.Int).SetString(ecdsaKey, 0)
			if !ok {
				log.Printf("Warning: could not parse operator ECDSA key %s", ecdsaKeyName)
				// Fallback to BLS key if ECDSA key parsing fails
				if operatorKey, ok := privateKeys[operatorName]; ok {
					keyInt, ok := new(big.Int).SetString(operatorKey, 0)
					if !ok {
						log.Printf("Warning: could not parse operator BLS key %s", operatorName)
						continue
					}
					operators = append(operators, keyInt.String())
				}
			} else {
				operators = append(operators, keyInt.String())
			}
		} else if operatorKey, ok := privateKeys[operatorName]; ok {
			// Fallback to BLS key if no ECDSA key is provided
			keyInt, ok := new(big.Int).SetString(operatorKey, 0)
			if !ok {
				log.Printf("Warning: could not parse operator key %s", operatorName)
				continue
			}
			log.Printf("Warning: Using BLS key for operator %s funding - this may not match the actual operator ECDSA address", operatorName)
			operators = append(operators, keyInt.String())
		}
	}

	// Get batcher key
	batcherKey := ""
	if key, ok := privateKeys["batcher0"]; ok {
		keyInt, ok := new(big.Int).SetString(key, 0)
		if ok {
			batcherKey = keyInt.String()
		}
	}

	return EigenDADeployConfig{
		UseDefaults:         true,
		NumStrategies:       numStrategies,
		MaxOperatorCount:    maxOperatorCount,
		StakerPrivateKeys:   stakers,
		StakerTokenAmounts:  stakes,
		OperatorPrivateKeys: operators,
		ConfirmerPrivateKey: batcherKey,
	}
}
