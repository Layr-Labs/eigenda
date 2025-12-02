package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/ejector"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	ctx := context.Background()

	err := run(ctx)
	if err != nil {
		panic(err)
	}

	// Block forever, the ejector runs in background goroutines.
	<-ctx.Done()
}

// Run the ejector. This method is split from main() so we only have to use panic() once.
func run(ctx context.Context) error {
	cfg, err := config.Bootstrap(ejector.DefaultRootEjectorConfig)
	if err != nil {
		return fmt.Errorf("failed to bootstrap config: %w", err)
	}
	secretConfig := cfg.Secret
	ejectorConfig := cfg.Config
	// Ensure we don't accidentally use cfg after this point.
	cfg = nil

	loggerConfig := common.DefaultLoggerConfig()
	loggerConfig.Format = common.LogFormat(ejectorConfig.LogOutputType)
	loggerConfig.HandlerOpts.NoColor = !ejectorConfig.LogColor

	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	var privateKey *ecdsa.PrivateKey
	var senderAddress gethcommon.Address

	if secretConfig.PrivateKey != "" {
		// Use private key if provided
		privateKey, err = crypto.HexToECDSA(secretConfig.PrivateKey)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}

		// Derive the public address from the private key
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			return fmt.Errorf("failed to get ECDSA public key")
		}
		senderAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
	} else {
		// Derive address from KMS public key
		senderAddress, err = getAddressFromKMS(
			ctx,
			ejectorConfig.KmsKeyId,
			ejectorConfig.KmsRegion,
			ejectorConfig.KmsEndpoint,
		)
		if err != nil {
			return fmt.Errorf("failed to get address from KMS: %w", err)
		}
	}

	gethClient, err := geth.NewMultiHomingClient(
		geth.EthClientConfig{
			RPCURLs:          secretConfig.EthRpcUrls,
			PrivateKeyString: secretConfig.PrivateKey,
			NumConfirmations: ejectorConfig.EthBlockConfirmations,
			NumRetries:       ejectorConfig.EthRpcRetryCount,
		},
		senderAddress,
		logger)
	if err != nil {
		return fmt.Errorf("failed to create geth client: %w", err)
	}

	contractDirectory, err := directory.NewContractDirectory(
		ctx,
		logger,
		gethClient,
		gethcommon.HexToAddress(ejectorConfig.ContractDirectoryAddress))
	if err != nil {
		return fmt.Errorf("failed to create contract directory: %w", err)
	}

	ejectionContractAddress, err := contractDirectory.GetContractAddress(ctx, directory.EigenDAEjectionManager)
	if err != nil {
		return fmt.Errorf("failed to get ejection manager address: %w", err)
	}

	registryCoordinatorAddress, err := contractDirectory.GetContractAddress(ctx, directory.RegistryCoordinator)
	if err != nil {
		return fmt.Errorf("failed to get registry coordinator address: %w", err)
	}

	chainID, err := gethClient.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	ejectionTransactor, err := ejector.NewEjectionTransactor(
		logger,
		gethClient,
		ejectionContractAddress,
		registryCoordinatorAddress,
		senderAddress,
		privateKey,
		chainID,
		ejectorConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to create ejection transactor: %w", err)
	}

	ejectionManager, err := ejector.NewEjectionManager(
		ctx,
		logger,
		ejectorConfig,
		time.Now,
		ejectionTransactor,
	)
	if err != nil {
		return fmt.Errorf("failed to create ejection manager: %w", err)
	}

	threadedEjectionManager := ejector.NewThreadedEjectionManager(ctx, logger, ejectionManager, ejectorConfig)

	// Currently used for both v1 and v2 signing rate lookups. Eventually, v2 will poll the controller for this info.
	dataApiSigningRateLookup := ejector.NewDataApiSigningRateLookup(
		logger,
		ejectorConfig.DataApiUrl,
		ejectorConfig.DataApiTimeout,
	)

	validatorIDToAddressConverter, err := eth.NewValidatorIDToAddressConverter(
		gethClient,
		registryCoordinatorAddress)
	if err != nil {
		return fmt.Errorf("failed to create validator ID to address converter: %w", err)
	}
	validatorIDToAddressConverter, err = eth.NewCachedValidatorIDToAddressConverter(validatorIDToAddressConverter, 1024)
	if err != nil {
		return fmt.Errorf("failed to create cached validator ID to address converter: %w", err)
	}

	referenceBlockProvider := eth.NewReferenceBlockProvider(
		logger,
		gethClient,
		ejectorConfig.ReferenceBlockNumberOffset,
	)

	validatorQuorumLookup, err := eth.NewValidatorQuorumLookup(
		gethClient,
		registryCoordinatorAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to create validator quorum lookup: %w", err)
	}
	validatorQuorumLookup, err = eth.NewCachedValidatorQuorumLookup(validatorQuorumLookup, 1024)
	if err != nil {
		return fmt.Errorf("failed to create cached validator quorum lookup: %w", err)
	}

	stakeRegistryAddress, err := contractDirectory.GetContractAddress(ctx, directory.StakeRegistry)
	if err != nil {
		return fmt.Errorf("failed to get stake registry address: %w", err)
	}

	validatorStakeLookup, err := eth.NewValidatorStakeLookup(gethClient, stakeRegistryAddress)
	if err != nil {
		return fmt.Errorf("failed to create validator stake lookup: %w", err)
	}
	validatorStakeLookup, err = eth.NewCachedValidatorStakeLookup(validatorStakeLookup, 1024)
	if err != nil {
		return fmt.Errorf("failed to create cached validator stake lookup: %w", err)
	}

	_ = ejector.NewEjector(
		ctx,
		logger,
		ejectorConfig,
		threadedEjectionManager,
		dataApiSigningRateLookup,
		dataApiSigningRateLookup,
		validatorIDToAddressConverter,
		referenceBlockProvider,
		validatorQuorumLookup,
		validatorStakeLookup,
	)
	return nil
}

// getAddressFromKMS retrieves the Ethereum address associated with a KMS key.
func getAddressFromKMS(
	ctx context.Context,
	kmsKeyId string,
	kmsRegion string,
	kmsEndpoint string,
) (gethcommon.Address, error) {
	// Load AWS config for the specified region
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(kmsRegion))
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create KMS client with optional custom endpoint
	kmsClient := awskms.NewFromConfig(awsCfg, func(o *awskms.Options) {
		if kmsEndpoint != "" {
			o.BaseEndpoint = &kmsEndpoint
		}
	})

	// Load the public key from KMS
	publicKey, err := aws.LoadPublicKeyKMS(ctx, kmsClient, kmsKeyId)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("failed to load public key from KMS: %w", err)
	}

	// Derive the Ethereum address from the public key
	address := crypto.PubkeyToAddress(*publicKey)
	return address, nil
}