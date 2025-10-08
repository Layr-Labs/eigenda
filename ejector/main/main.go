package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/ejector"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

const ejectorEnvVarPrefix = "EJECTOR"

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
	cfg, err := config.Bootstrap(ejector.DefaultEjectorConfig, ejectorEnvVarPrefix)
	if err != nil {
		return fmt.Errorf("failed to bootstrap config: %w", err)
	}

	logger, err := common.NewLogger(cfg.LoggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// TODO is this important?
	var senderAddress gethcommon.Address

	gethClient, err := geth.NewMultiHomingClient(
		geth.EthClientConfig{
			RPCURLs:          cfg.EthRpcUrls,
			PrivateKeyString: cfg.PrivateKey,
			NumConfirmations: cfg.EthBlockConfirmations,
			NumRetries:       cfg.EthRpcRetryCount,
		},
		senderAddress,
		logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return fmt.Errorf("failed to create geth client: %w", err)
	}

	contractDirectory, err := directory.NewContractDirectory(
		context.Background(),
		logger,
		gethClient,
		gethcommon.HexToAddress(cfg.ContractDirectoryAddress))
	if err != nil {
		return fmt.Errorf("failed to create contract directory: %w", err)
	}

	ejectionTransactor := ejector.NewEjectionTransactor()

	ejectionManager, err := ejector.NewEjectionManager(
		ctx,
		logger,
		cfg,
		time.Now,
		ejectionTransactor,
	)
	if err != nil {
		return fmt.Errorf("failed to create ejection manager: %w", err)
	}

	threadedEjectionManager := ejector.NewThreadedEjectionManager(ctx, logger, ejectionManager, cfg)

	// Currently used for both v1 and v2 signing rate lookups. Eventually, v2 will poll the controller for this info.
	dataApiSigningRateLookup := ejector.NewDataApiSigningRateLookup(
		logger,
		cfg.DataApiUrl,
		cfg.DataApiTimeout,
	)

	registryCoordinatorAddress, err :=
		contractDirectory.GetContractAddress(context.Background(), directory.RegistryCoordinator)
	if err != nil {
		return fmt.Errorf("failed to get RegistryCoordinator address: %w", err)
	}

	validatorIDCache, err := eth.NewValidatorIDToAddressCache(
		gethClient,
		registryCoordinatorAddress,
		1024)
	if err != nil {
		return fmt.Errorf("failed to create validator ID to address cache: %w", err)
	}

	_ = ejector.NewEjector(
		ctx,
		logger,
		cfg,
		threadedEjectionManager,
		dataApiSigningRateLookup,
		dataApiSigningRateLookup,
		validatorIDCache,
	)
	return nil
}
