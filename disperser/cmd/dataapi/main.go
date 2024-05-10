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
	"github.com/Layr-Labs/eigenda/common/geth"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/disperser/cmd/dataapi/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
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

	client, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.HexToAddress(config.FireblocksConfig.WalletAddress), logger)
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
		return common.NewFireblocksWallet(&config.FireblocksConfig, ethClient, logger)
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
