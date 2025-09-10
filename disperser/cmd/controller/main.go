package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/disperser/controller/metadata"
	"github.com/Layr-Labs/eigenda/disperser/controller/metrics"
	"github.com/Layr-Labs/eigenda/disperser/controller/payments"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/cmd/controller/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/disperser/controller/grpcserver"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gammazero/workerpool"
	"github.com/urfave/cli"
)

var (
	version   string
	gitCommit string
	gitDate   string

	controllerMaxStallDuration = 240 * time.Second
)

func main() {
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)
	app.Name = "controller"
	app.Usage = "EigenDA Controller"
	app.Description = "EigenDA control plane for encoding and dispatching blobs"

	app.Action = RunController
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}
	select {}
}

func RunController(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(&config.LoggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Reset readiness probe upon start-up
	if err := os.Remove(config.ControllerReadinessProbePath); err != nil {
		logger.Warn("Failed to clean up readiness file", "error", err, "path", config.ControllerReadinessProbePath)
	}

	dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
	if err != nil {
		return fmt.Errorf("failed to create DynamoDB client: %w", err)
	}
	gethClient, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return fmt.Errorf("failed to create geth client: %w", err)
	}

	contractDirectory, err := directory.NewContractDirectory(
		context.Background(),
		logger,
		gethClient,
		gethcommon.HexToAddress(config.EigenDAContractDirectoryAddress))
	if err != nil {
		return fmt.Errorf("failed to create contract directory: %w", err)
	}

	operatorStateRetrieverAddress, err :=
		contractDirectory.GetContractAddress(context.Background(), directory.OperatorStateRetriever)
	if err != nil {
		return fmt.Errorf("failed to get OperatorStateRetriever address: %w", err)
	}
	serviceManagerAddress, err :=
		contractDirectory.GetContractAddress(context.Background(), directory.ServiceManager)
	if err != nil {
		return fmt.Errorf("failed to get ServiceManager address: %w", err)
	}
	registryCoordinatorAddress, err :=
		contractDirectory.GetContractAddress(context.Background(), directory.RegistryCoordinator)
	if err != nil {
		return fmt.Errorf("failed to get registry coordinator address: %w", err)
	}

	chainReader, err := eth.NewReader(
		logger,
		gethClient,
		operatorStateRetrieverAddress.Hex(),
		serviceManagerAddress.Hex())
	if err != nil {
		return fmt.Errorf("failed to create chain reader: %w", err)
	}

	metricsRegistry := prometheus.NewRegistry()
	metricsRegistry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	metricsRegistry.MustRegister(collectors.NewGoCollector())

	logger.Infof("Starting metrics server at port %d", config.MetricsPort)
	addr := fmt.Sprintf(":%d", config.MetricsPort)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		metricsRegistry,
		promhttp.HandlerOpts{},
	))
	metricsServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	baseBlobMetadataStore := blobstore.NewBlobMetadataStore(
		dynamoClient,
		logger,
		config.DynamoDBTableName,
	)
	blobMetadataStore := blobstore.NewInstrumentedMetadataStore(baseBlobMetadataStore, blobstore.InstrumentedMetadataStoreConfig{
		ServiceName: "controller",
		Registry:    metricsRegistry,
		Backend:     blobstore.BackendDynamoDB,
	})

	controllerLivenessChan := make(chan healthcheck.HeartbeatMessage, 10)

	encoderClient, err := encoder.NewEncoderClientV2(config.EncodingManagerConfig.EncoderAddress)
	if err != nil {
		return fmt.Errorf("failed to create encoder client: %v", err)
	}
	encodingPool := workerpool.New(config.NumConcurrentEncodingRequests)
	encodingManagerBlobSet := controller.NewBlobSet()
	encodingManager, err := controller.NewEncodingManager(
		&config.EncodingManagerConfig,
		blobMetadataStore,
		encodingPool,
		encoderClient,
		chainReader,
		logger,
		metricsRegistry,
		encodingManagerBlobSet,
		controllerLivenessChan,
	)
	if err != nil {
		return fmt.Errorf("failed to create encoding manager: %v", err)
	}

	sigAgg, err := core.NewStdSignatureAggregator(logger, chainReader)
	if err != nil {
		return fmt.Errorf("failed to create signature aggregator: %v", err)
	}
	dispatcherPool := workerpool.New(config.NumConcurrentDispersalRequests)
	chainState := eth.NewChainState(chainReader, gethClient)
	var ics core.IndexedChainState
	if config.UseGraph {
		logger.Info("Using graph node")

		logger.Info("Connecting to subgraph", "url", config.ChainStateConfig.Endpoint)
		ics = thegraph.MakeIndexedChainState(config.ChainStateConfig, chainState, logger)
	} else {
		logger.Info("Using built-in indexer")
		rpcClient, err := rpc.Dial(config.EthClientConfig.RPCURLs[0])
		if err != nil {
			return err
		}
		idx, err := indexer.CreateNewIndexer(
			&config.IndexerConfig,
			gethClient,
			rpcClient,
			serviceManagerAddress.Hex(),
			logger,
		)
		if err != nil {
			return fmt.Errorf("failed to create indexer: %w", err)
		}
		ics, err = indexer.NewIndexedChainState(chainState, idx)
		if err != nil {
			return fmt.Errorf("failed to create indexed chain state: %w", err)
		}
	}

	var requestSigner clients.DispersalRequestSigner
	if config.DisperserStoreChunksSigningDisabled {
		logger.Warn("StoreChunks() signing is disabled")
	} else {
		requestSigner, err = clients.NewDispersalRequestSigner(
			context.Background(),
			config.AwsClientConfig.Region,
			config.AwsClientConfig.EndpointURL,
			config.DisperserKMSKeyID)
		if err != nil {
			return fmt.Errorf("failed to create request signer: %v", err)
		}
	}

	nodeClientManager, err := controller.NewNodeClientManager(config.NodeClientCacheSize, requestSigner, logger)
	if err != nil {
		return fmt.Errorf("failed to create node client manager: %v", err)
	}
	beforeDispatch := func(blobKey corev2.BlobKey) error {
		encodingManagerBlobSet.RemoveBlob(blobKey)
		return nil
	}
	dispatcherBlobSet := controller.NewBlobSet()

	batchMetadataManager, err := metadata.NewBatchMetadataManager(
		context.Background(),
		logger,
		gethClient,
		ics,
		registryCoordinatorAddress,
		config.DispatcherConfig.BatchMetadataUpdatePeriod,
		config.DispatcherConfig.FinalizationBlockDelay,
	)
	if err != nil {
		return fmt.Errorf("failed to create batch metadata manager: %w", err)
	}

	dispatcher, err := controller.NewDispatcher(
		&config.DispatcherConfig,
		blobMetadataStore,
		dispatcherPool,
		ics,
		batchMetadataManager,
		sigAgg,
		nodeClientManager,
		logger,
		metricsRegistry,
		beforeDispatch,
		dispatcherBlobSet,
		controllerLivenessChan,
	)
	if err != nil {
		return fmt.Errorf("failed to create dispatcher: %v", err)
	}

	c := context.Background()

	err = controller.RecoverState(c, blobMetadataStore, logger)
	if err != nil {
		return fmt.Errorf("failed to recover state: %v", err)
	}

	err = encodingManager.Start(c)
	if err != nil {
		return fmt.Errorf("failed to start encoding manager: %v", err)
	}

	err = dispatcher.Start(c)
	if err != nil {
		return fmt.Errorf("failed to start dispatcher: %v", err)
	}

	if config.ServerConfig.EnableServer {
		serverMetrics := metrics.NewServerMetrics(metricsRegistry, logger)

		var paymentAuthorizationHandler *payments.PaymentAuthorizationHandler
		if config.ServerConfig.EnablePaymentAuthentication {
			paymentAuthorizationHandler, err = payments.NewPaymentAuthorizationHandler(
				context.Background(),
				serverMetrics,
				config.AwsClientConfig.Region,
				config.AwsClientConfig.EndpointURL,
				config.DisperserKMSKeyID)
			if err != nil {
				return fmt.Errorf("create payment authorization handler: %w", err)
			}
		}

		grpcServer, err := grpcserver.NewServer(
			config.ServerConfig,
			logger,
			serverMetrics,
			paymentAuthorizationHandler)
		if err != nil {
			return fmt.Errorf("create gRPC server: %w", err)
		}

		go func() {
			logger.Info("Starting controller gRPC server", "port", config.ServerConfig.GrpcPort)
			if err := grpcServer.Start(); err != nil {
				logger.Error("gRPC server failed", "error", err)
			}
		}()
	} else {
		logger.Info("Controller gRPC server disabled")
	}

	go func() {
		err := metricsServer.ListenAndServe()
		if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
			logger.Errorf("metrics metricsServer error: %v", err)
		}
	}()

	// Create readiness probe file once the controller starts successfully
	if _, err := os.Create(config.ControllerReadinessProbePath); err != nil {
		logger.Warn("Failed to create readiness file", "error", err, "path", config.ControllerReadinessProbePath)
	}

	if _, err := os.Create(config.ControllerHealthProbePath); err != nil {
		logger.Warn("Failed to create healthProbe file: %v", err)
	}

	// Start heartbeat monitor
	go func() {
		err := healthcheck.HeartbeatMonitor(
			config.ControllerHealthProbePath,
			controllerMaxStallDuration,
			controllerLivenessChan,
			logger,
		)
		if err != nil {
			logger.Warn("Heartbeat monitor exited with error", "err", err)
		}
	}()

	return nil
}
