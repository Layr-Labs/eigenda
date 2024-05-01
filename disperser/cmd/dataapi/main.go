package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/aws/secretmanager"
	"github.com/Layr-Labs/eigenda/common/geth"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/disperser/cmd/dataapi/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/fireblocks"
	walletsdk "github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/signerv2"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"
)

var (
	// version is the version of the binary.
	version   string
	gitCommit string
	gitDate   string
)

// @title			EigenDA Data Access API
// @description	This is the EigenDA Data Access API server.
// @version		1
// @Schemes		https http
func main() {
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)
	app.Name = "data-access-api"
	app.Usage = "EigenDA Data Access API"
	app.Description = "Service that provides access to data blobs."

	app.Action = RunDataApi
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func RunDataApi(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	s3Client, err := s3.NewClient(context.Background(), config.AwsClientConfig, logger)
	if err != nil {
		return err
	}

	dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
	if err != nil {
		return err
	}

	promApi, err := prometheus.NewApi(config.PrometheusConfig)
	if err != nil {
		return err
	}

	client, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		return err
	}

	tx, err := coreeth.NewTransactor(logger, client, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return err
	}

	wallet, err := getWallet(config, client, logger)
	if err != nil {
		return err
	}
	var (
		promClient        = dataapi.NewPrometheusClient(promApi, config.PrometheusConfig.Cluster)
		blobMetadataStore = blobstore.NewBlobMetadataStore(dynamoClient, logger, config.BlobstoreConfig.TableName, 0)
		sharedStorage     = blobstore.NewSharedStorage(config.BlobstoreConfig.BucketName, s3Client, blobMetadataStore, logger)
		subgraphApi       = subgraph.NewApi(config.SubgraphApiBatchMetadataAddr, config.SubgraphApiOperatorStateAddr)
		subgraphClient    = dataapi.NewSubgraphClient(subgraphApi, logger)
		chainState        = coreeth.NewChainState(tx, client)
		metrics           = dataapi.NewMetrics(blobMetadataStore, config.MetricsConfig.HTTPPort, logger)
		server            = dataapi.NewServer(
			dataapi.Config{
				ServerMode:         config.ServerMode,
				SocketAddr:         config.SocketAddr,
				AllowOrigins:       config.AllowOrigins,
				EjectionToken:      config.EjectionToken,
				DisperserHostname:  config.DisperserHostname,
				ChurnerHostname:    config.ChurnerHostname,
				BatcherHealthEndpt: config.BatcherHealthEndpt,
			},
			sharedStorage,
			promClient,
			subgraphClient,
			tx,
			chainState,
			dataapi.NewEjector(wallet, client, logger, tx, metrics, config.TxnTimeout),
			logger,
			metrics,
			nil,
			nil,
			nil,
		)
	)

	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Data Access API", "socket", httpSocket)
	}

	// Setup channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	// catch SIGINT (Ctrl+C) and SIGTERM (e.g., from `kill`)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a separate goroutine so that it doesn't block.
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Block until a signal is received.
	<-quit
	logger.Info("Shutting down server...")
	err = server.Shutdown()

	if err != nil {
		logger.Errorf("Failed to shutdown server: %v", err)
	}

	return err
}

func getWallet(config Config, ethClient common.EthClient, logger logging.Logger) (walletsdk.Wallet, error) {
	var wallet walletsdk.Wallet
	if !config.FireblocksConfig.Disable {
		validConfigflag := len(config.FireblocksConfig.APIKeyName) > 0 &&
			len(config.FireblocksConfig.SecretKeyName) > 0 &&
			len(config.FireblocksConfig.BaseURL) > 0 &&
			len(config.FireblocksConfig.VaultAccountName) > 0 &&
			len(config.FireblocksConfig.WalletAddress) > 0 &&
			len(config.FireblocksConfig.Region) > 0
		if !validConfigflag {
			return nil, errors.New("fireblocks config is either invalid or incomplete")
		}
		apiKey, err := secretmanager.ReadStringFromSecretManager(context.Background(), config.FireblocksConfig.APIKeyName, config.FireblocksConfig.Region)
		if err != nil {
			return nil, fmt.Errorf("cannot read fireblocks api key %s from secret manager: %w", config.FireblocksConfig.APIKeyName, err)
		}
		secretKey, err := secretmanager.ReadStringFromSecretManager(context.Background(), config.FireblocksConfig.SecretKeyName, config.FireblocksConfig.Region)
		if err != nil {
			return nil, fmt.Errorf("cannot read fireblocks secret key %s from secret manager: %w", config.FireblocksConfig.SecretKeyName, err)
		}
		fireblocksClient, err := fireblocks.NewClient(
			apiKey,
			[]byte(secretKey),
			config.FireblocksConfig.BaseURL,
			config.FireblocksConfig.APITimeout,
			logger.With("component", "FireblocksClient"),
		)
		if err != nil {
			return nil, err
		}
		wallet, err = walletsdk.NewFireblocksWallet(fireblocksClient, ethClient, config.FireblocksConfig.VaultAccountName, logger.With("component", "FireblocksWallet"))
		if err != nil {
			return nil, err
		}
		sender, err := wallet.SenderAddress(context.Background())
		if err != nil {
			return nil, err
		}
		if sender.Cmp(gethcommon.HexToAddress(config.FireblocksConfig.WalletAddress)) != 0 {
			return nil, fmt.Errorf("configured wallet address %s does not match derived address %s", config.FireblocksConfig.WalletAddress, sender.Hex())
		}
		logger.Info("Initialized Fireblocks wallet", "vaultAccountName", config.FireblocksConfig.VaultAccountName, "address", sender.Hex())
	} else if len(config.EthClientConfig.PrivateKeyString) > 0 {
		privateKey, err := crypto.HexToECDSA(config.EthClientConfig.PrivateKeyString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		chainID, err := ethClient.ChainID(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get chain ID: %w", err)
		}
		signerV2, address, err := signerv2.SignerFromConfig(signerv2.Config{PrivateKey: privateKey}, chainID)
		if err != nil {
			return nil, err
		}
		wallet, err = walletsdk.NewPrivateKeyWallet(ethClient, signerV2, address, logger.With("component", "PrivateKeyWallet"))
		if err != nil {
			return nil, err
		}
		logger.Info("Initialized PrivateKey wallet", "address", address.Hex())
	} else {
		return nil, errors.New("no wallet is configured. Either Fireblocks or PrivateKey wallet should be configured")
	}

	return wallet, nil
}
