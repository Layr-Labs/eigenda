package deploy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	caws "github.com/Layr-Labs/eigenda/common/aws"
	relayreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDARelayRegistry"
	eigendasrvmg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	thresholdreg "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAThresholdRegistry"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// convertToTestbedPrivateKeys converts the current PkConfig to testbed.PrivateKeyMaps
func (env *Config) convertToTestbedPrivateKeys() *testbed.PrivateKeyMaps {
	if env.Pks == nil {
		return nil
	}

	result := &testbed.PrivateKeyMaps{
		EcdsaMap: make(map[string]testbed.KeyInfo),
		BlsMap:   make(map[string]testbed.KeyInfo),
	}

	for name, keyInfo := range env.Pks.EcdsaMap {
		result.EcdsaMap[name] = testbed.KeyInfo{
			PrivateKey: keyInfo.PrivateKey,
			Password:   keyInfo.Password,
			KeyFile:    keyInfo.KeyFile,
		}
	}

	for name, keyInfo := range env.Pks.BlsMap {
		result.BlsMap[name] = testbed.KeyInfo{
			PrivateKey: keyInfo.PrivateKey,
			Password:   keyInfo.Password,
			KeyFile:    keyInfo.KeyFile,
		}
	}

	return result
}

// deployEigenDAContracts deploys EigenDA core system and peripheral contracts on local anvil chain
func (env *Config) deployEigenDAContracts() error {
	logger.Info("Deploy the EigenDA and EigenLayer contracts using testbed")

	// get deployer
	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	if !ok {
		return fmt.Errorf("deployer improperly configured")
	}

	// Convert Stakes to testbed format
	stakes := make([]testbed.Stakes, len(env.Services.Stakes))
	for i, stake := range env.Services.Stakes {
		stakes[i] = testbed.Stakes{
			Total:        stake.Total,
			Distribution: stake.Distribution,
		}
	}

	// Create deployment config for testbed
	deployConfig := testbed.DeploymentConfig{
		AnvilRPCURL:      deployer.RPC,
		DeployerKey:      env.Pks.EcdsaMap[deployer.Name].PrivateKey,
		NumOperators:     env.Services.Counts.NumOpr,
		NumRelays:        env.Services.Counts.NumRelays,
		Stakes:           stakes,
		MaxOperatorCount: env.Services.Counts.NumMaxOperatorCount,
		PrivateKeys:      env.convertToTestbedPrivateKeys(),
		Logger:           logger,
	}

	// Deploy contracts using testbed
	result, err := testbed.DeployEigenDAContracts(deployConfig)
	if err != nil {
		return fmt.Errorf("failed to deploy EigenDA contracts: %w", err)
	}

	// Copy results to env
	env.EigenDA = EigenDAContract{
		Deployer:               env.EigenDA.Deployer,
		EigenDADirectory:       result.EigenDA.EigenDADirectory,
		ServiceManager:         result.EigenDA.ServiceManager,
		OperatorStateRetriever: result.EigenDA.OperatorStateRetriever,
		BlsApkRegistry:         result.EigenDA.BlsApkRegistry,
		RegistryCoordinator:    result.EigenDA.RegistryCoordinator,
		CertVerifierLegacyV2:   result.EigenDA.CertVerifierLegacyV2,
		CertVerifierLegacyV3:   result.EigenDA.CertVerifierLegacyV3,
		CertVerifier:           result.EigenDA.CertVerifier,
		CertVerifierRouter:     result.EigenDA.CertVerifierRouter,
	}

	return nil
}

// Deploys a EigenDA experiment
func (env *Config) DeployExperiment() error {
	if err := changeDirectory(filepath.Join(env.rootPath, "inabox")); err != nil {
		return fmt.Errorf("error changing directories: %w", err)
	}

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil {
		logger.Info("Successfully changed to absolute path", "path", cwd)
	}

	defer env.SaveTestConfig()

	logger.Info("Deploying experiment...")

	// Create a new experiment and deploy the contracts

	err := env.loadPrivateKeys()
	if err != nil {
		return fmt.Errorf("could not load private keys: %w", err)
	}

	if env.EigenDA.Deployer != "" && !env.IsEigenDADeployed() {
		logger.Info("Deploying EigenDA")
		err = env.deployEigenDAContracts()
		if err != nil {
			return fmt.Errorf("error deploying EigenDA contracts: %w", err)
		}
	}

	if deployer, ok := env.GetDeployer(env.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
		startBlock, err := GetLatestBlockNumber(env.Deployers[0].RPC)
		if err != nil {
			return fmt.Errorf("error getting latest block number: %w", err)
		}

		config := testbed.SubgraphDeploymentConfig{
			RootPath:            env.rootPath,
			RegistryCoordinator: env.EigenDA.RegistryCoordinator,
			BlsApkRegistry:      env.EigenDA.BlsApkRegistry,
			ServiceManager:      env.EigenDA.ServiceManager,
			Logger:              logger,
		}

		err = testbed.DeploySubgraphs(config, startBlock)
		if err != nil {
			return fmt.Errorf("error deploying subgraphs: %w", err)
		}
	}

	// Ideally these should be set in GenerateAllVariables, but they need to be used in GenerateDisperserKeypair
	// which is called before GenerateAllVariables
	env.localstackEndpoint = "http://localhost:4570"
	env.localstackRegion = "us-east-1"

	logger.Info("Test environment has successfully deployed!")
	return nil
}

// GenerateDisperserKeypair generates a disperser keypair using AWS KMS.
func (env *Config) GenerateDisperserKeypair() error {
	// Skip if we already have a disperser key
	if env.DisperserKMSKeyID != "" {
		logger.Info("Disperser keypair already exists, skipping generation")
		return nil
	}

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
			logger.Warnf("Unable to reach local stack, skipping disperser keypair generation. Error: %v", err)
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
	logger.Infof("Generated disperser keypair: key ID: %s, address: %s",
		env.DisperserKMSKeyID, env.DisperserAddress.Hex())

	return nil
}

// PerformDisperserRegistrations registers the disperser keypair onchain.
func (env *Config) PerformDisperserRegistrations(ethClient common.EthClient) {
	// Only register disperser keypair if we have a valid address
	if env.DisperserAddress != (gcommon.Address{}) {
		logger.Info("Registering disperser keypair")
		err := env.registerDisperserKeypair(ethClient)
		if err != nil {
			logger.Errorf("could not register disperser keypair: %v", err)
		}
	} else {
		logger.Info("Skipping disperser keypair registration")
	}
}

// RegisterDisperserKeypair registers the disperser's public key on-chain.
func (env *Config) registerDisperserKeypair(ethClient common.EthClient) error {
	// Write the disperser's public key to on-chain storage
	writer, err := eth.NewWriter(
		logger,
		ethClient,
		env.EigenDA.OperatorStateRetriever,
		env.EigenDA.ServiceManager,
	)
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}

	err = writer.SetDisperserAddress(context.Background(), 0, env.DisperserAddress)
	if err != nil {
		return fmt.Errorf("could not set disperser address: %v", err)
	}

	// Read the disperser's public key from on-chain storage to verify it was written correctly

	retryTimeout := time.Now().Add(1 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)

	for time.Now().Before(retryTimeout) {
		address, err := writer.GetDisperserAddress(context.Background(), 0)
		if err != nil {
			logger.Warnf("could not get disperser address: %v", err)
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

// RegisterBlobVersions initializes blob versions in ThresholdRegistry contract
func (env *Config) RegisterBlobVersions(ethClient common.EthClient) {
	dasmAddr := gcommon.HexToAddress(env.EigenDA.ServiceManager)
	if (dasmAddr == gcommon.Address{}) {
		logger.Fatal("Service Manager address is nil")
	}
	contractEigenDAServiceManager, err := eigendasrvmg.NewContractEigenDAServiceManager(dasmAddr, ethClient)
	if err != nil {
		logger.Fatal("Error creating EigenDAServiceManager contract", "error", err)
	}
	thresholdRegistryAddr, err := contractEigenDAServiceManager.EigenDAThresholdRegistry(&bind.CallOpts{})
	if err != nil {
		logger.Fatal("Error getting threshold registry address", "error", err)
	}
	contractThresholdRegistry, err := thresholdreg.NewContractEigenDAThresholdRegistry(thresholdRegistryAddr, ethClient)
	if err != nil {
		logger.Fatal("Error creating threshold registry contract", "error", err)
	}
	opts, err := ethClient.GetNoSendTransactOpts()
	if err != nil {
		logger.Fatal("Error getting transaction opts", "error", err)
	}
	for _, blobVersionParam := range env.BlobVersionParams {
		txn, err := contractThresholdRegistry.AddVersionedBlobParams(opts, thresholdreg.EigenDATypesV1VersionedBlobParams{
			MaxNumOperators: blobVersionParam.MaxNumOperators,
			NumChunks:       blobVersionParam.NumChunks,
			CodingRate:      uint8(blobVersionParam.CodingRate),
		})
		if err != nil {
			logger.Fatal("Error adding versioned blob params", "error", err)
		}
		err = ethClient.SendTransaction(context.Background(), txn)
		if err != nil {
			logger.Fatal("Error sending blob version transaction", "error", err)
		}
	}
}

// RegisterRelays initializes relays in RelayRegistry contract
func (env *Config) RegisterRelays(ethClient common.EthClient, relayURLs []string, relayAddress gcommon.Address) {
	dasmAddr := gcommon.HexToAddress(env.EigenDA.ServiceManager)
	if (dasmAddr == gcommon.Address{}) {
		logger.Fatal("Service Manager address is nil")
	}
	contractEigenDAServiceManager, err := eigendasrvmg.NewContractEigenDAServiceManager(dasmAddr, ethClient)
	if err != nil {
		logger.Fatal("Error creating EigenDAServiceManager contract", "error", err)
	}
	relayAddr, err := contractEigenDAServiceManager.EigenDARelayRegistry(&bind.CallOpts{})
	if err != nil {
		logger.Fatal("Error getting relay registry address", "error", err)
	}
	contractRelayRegistry, err := relayreg.NewContractEigenDARelayRegistry(relayAddr, ethClient)
	if err != nil {
		logger.Fatal("Error creating relay registry contract", "error", err)
	}
	opts, err := ethClient.GetNoSendTransactOpts()
	if err != nil {
		logger.Fatal("Error getting transaction opts", "error", err)
	}
	for _, url := range relayURLs {
		txn, err := contractRelayRegistry.AddRelayInfo(opts, relayreg.EigenDATypesV2RelayInfo{
			RelayAddress: relayAddress,
			RelayURL:     url,
		})
		if err != nil {
			logger.Fatal("Error adding relay info", "error", err)
		}
		err = ethClient.SendTransaction(context.Background(), txn)
		if err != nil {
			logger.Fatal("Error sending relay transaction", "error", err)
		}
	}
}
