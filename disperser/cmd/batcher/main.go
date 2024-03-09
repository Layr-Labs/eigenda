package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shurcooL/graphql"

	"github.com/Layr-Labs/eigenda/common"
	coreindexer "github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigenda/core/thegraph"

	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	dispatcher "github.com/Layr-Labs/eigenda/disperser/batcher/grpc"
	"github.com/Layr-Labs/eigenda/disperser/cmd/batcher/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/fireblocks"
	walletsdk "github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli"
)

var (
	// version is the version of the binary.
	version   string
	gitCommit string
	gitDate   string
	// Note: Changing these paths will require updating the k8s deployment
	readinessProbePath      string        = "/tmp/ready"
	healthProbePath         string        = "/tmp/health"
	maxStallDuration        time.Duration = 240 * time.Second
	handleBatchLivenessChan               = make(chan time.Time, 1)
)

func main() {
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)
	app.Name = "batcher"
	app.Usage = "EigenDA Batcher"
	app.Description = "Service for creating a batch from queued blobs, distributing coded chunks to nodes, and confirming onchain"

	app.Action = RunBatcher
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	if _, err := os.Create(healthProbePath); err != nil {
		log.Printf("Failed to create healthProbe file: %v", err)
	}

	// Start HeartBeat Monitor
	go heartbeatMonitor(healthProbePath, maxStallDuration)

	select {}
}

func RunBatcher(ctx *cli.Context) error {

	// Clean up readiness file
	if err := os.Remove(readinessProbePath); err != nil {
		log.Printf("Failed to clean up readiness file: %v at path %v \n", err, readinessProbePath)
	}

	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	bucketName := config.BlobstoreConfig.BucketName
	s3Client, err := s3.NewClient(context.Background(), config.AwsClientConfig, logger)
	if err != nil {
		return err
	}
	logger.Info("Initialized S3 client", "bucket", bucketName)

	dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
	if err != nil {
		return err
	}

	dispatcher := dispatcher.NewDispatcher(&dispatcher.Config{
		Timeout: config.TimeoutConfig.AttestationTimeout,
	}, logger)
	asgn := &core.StdAssignmentCoordinator{}

	client, err := geth.NewClient(config.EthClientConfig, logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return err
	}
	rpcClient, err := rpc.Dial(config.EthClientConfig.RPCURL)
	if err != nil {
		return err
	}
	tx, err := coreeth.NewTransactor(logger, client, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return err
	}
	agg, err := core.NewStdSignatureAggregator(logger, tx)
	if err != nil {
		return err
	}
	blockStaleMeasure, err := tx.GetBlockStaleMeasure(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get BLOCK_STALE_MEASURE: %w", err)
	}
	storeDurationBlocks, err := tx.GetStoreDurationBlocks(context.Background())
	if err != nil || storeDurationBlocks == 0 {
		return fmt.Errorf("failed to get STORE_DURATION_BLOCKS: %w", err)
	}
	blobMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, config.BlobstoreConfig.TableName, time.Duration((storeDurationBlocks+blockStaleMeasure)*12)*time.Second)
	queue := blobstore.NewSharedStorage(bucketName, s3Client, blobMetadataStore, logger)

	cs := coreeth.NewChainState(tx, client)

	var ics core.IndexedChainState
	if config.UseGraph {
		logger.Info("Using graph node")
		querier := graphql.NewClient(config.GraphUrl, nil)
		logger.Info("Connecting to subgraph", "url", config.GraphUrl)
		ics = thegraph.NewIndexedChainState(cs, querier, logger)
	} else {
		logger.Info("Using built-in indexer")

		indexer, err := coreindexer.CreateNewIndexer(
			&config.IndexerConfig,
			client,
			rpcClient,
			config.EigenDAServiceManagerAddr,
			logger,
		)
		if err != nil {
			return err
		}
		ics, err = coreindexer.NewIndexedChainState(cs, indexer)
		if err != nil {
			return err
		}
	}

	metrics := batcher.NewMetrics(config.MetricsConfig.HTTPPort, logger)

	if len(config.BatcherConfig.EncoderSocket) == 0 {
		return errors.New("encoder socket must be specified")
	}
	encoderClient, err := encoder.NewEncoderClient(config.BatcherConfig.EncoderSocket, config.TimeoutConfig.EncodingTimeout)
	if err != nil {
		return err
	}
	finalizer := batcher.NewFinalizer(config.TimeoutConfig.ChainReadTimeout, config.BatcherConfig.FinalizerInterval, queue, client, rpcClient, config.BatcherConfig.MaxNumRetriesPerBlob, 1000, config.BatcherConfig.FinalizerPoolSize, logger, metrics.FinalizerMetrics)
	var wallet walletsdk.Wallet
	if len(config.FireblocksConfig.APIKey) > 0 &&
		len(config.FireblocksConfig.SecretKeyPath) > 0 &&
		len(config.FireblocksConfig.BaseURL) > 0 &&
		len(config.FireblocksConfig.VaultAccountName) > 0 {
		secretKey, err := os.ReadFile(config.FireblocksConfig.SecretKeyPath)
		if err != nil {
			return fmt.Errorf("Cannot read fireblocks secret from %s: %w", config.FireblocksConfig.SecretKeyPath, err)
		}
		fireblocksClient, err := fireblocks.NewClient(
			config.FireblocksConfig.APIKey,
			secretKey,
			config.FireblocksConfig.BaseURL,
			config.TimeoutConfig.ChainReadTimeout,
			logger,
		)
		if err != nil {
			return err
		}
		wallet, err = walletsdk.NewFireblocksWallet(fireblocksClient, client, config.FireblocksConfig.VaultAccountName, logger)
		if err != nil {
			return err
		}
	} else if len(config.EthClientConfig.PrivateKeyString) > 0 {
		privateKey, err := crypto.HexToECDSA(config.EthClientConfig.PrivateKeyString)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}
		chainID, err := client.ChainID(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get chain ID: %w", err)
		}
		signerV2, address, err := signerv2.SignerFromConfig(signerv2.Config{PrivateKey: privateKey}, chainID)
		if err != nil {
			return err
		}
		wallet, err = walletsdk.NewPrivateKeyWallet(client, signerV2, address, logger)
		if err != nil {
			return err
		}
	} else {
		return errors.New("no wallet is configured. Either Fireblocks or PrivateKey wallet should be configured")
	}

	txnManager := batcher.NewTxnManager(client, wallet, config.EthClientConfig.NumConfirmations, 20, config.TimeoutConfig.ChainWriteTimeout, logger, metrics.TxnManagerMetrics)
	batcher, err := batcher.NewBatcher(config.BatcherConfig, config.TimeoutConfig, queue, dispatcher, ics, asgn, encoderClient, agg, client, finalizer, tx, txnManager, logger, metrics, handleBatchLivenessChan)
	if err != nil {
		return err
	}

	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Batcher", "socket", httpSocket)
	}

	err = batcher.Start(context.Background())
	if err != nil {
		return err
	}

	// Signal readiness
	if _, err := os.Create(readinessProbePath); err != nil {
		log.Printf("Failed to create readiness file: %v at path %v \n", err, readinessProbePath)
	}

	return nil

}

// process liveness signal from handleBatch Go Routine
func heartbeatMonitor(filePath string, maxStallDuration time.Duration) {
	var lastHeartbeat time.Time
	stallTimer := time.NewTimer(maxStallDuration)

	for {
		select {
		// HeartBeat from Goroutine on Batcher Pull Interval
		case heartbeat, ok := <-handleBatchLivenessChan:
			if !ok {
				log.Println("handleBatchLivenessChan closed, stopping health probe")
				return
			}
			log.Printf("Received heartbeat from HandleBatch GoRoutine: %v\n", heartbeat)
			lastHeartbeat = heartbeat
			if err := os.WriteFile(filePath, []byte(lastHeartbeat.String()), 0666); err != nil {
				log.Printf("Failed to update heartbeat file: %v", err)
			} else {
				log.Printf("Updated heartbeat file: %v with time %v\n", filePath, lastHeartbeat)
			}
			stallTimer.Reset(maxStallDuration) // Reset timer on new heartbeat

		case <-stallTimer.C:
			log.Println("No heartbeat received within max stall duration, stopping health probe")
			return
		}
	}
}
