package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/signingrate"
	"github.com/Layr-Labs/eigenda/disperser/controller/metadata"
	controllerpayments "github.com/Layr-Labs/eigenda/disperser/controller/payments"
	"github.com/Layr-Labs/eigenda/disperser/controller/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/common/nameremapping"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/cmd/controller/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gammazero/workerpool"
	"github.com/urfave/cli"
)

var (
	version   string
	gitCommit string
	gitDate   string
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

func RunController(cliCtx *cli.Context) error {
	config, err := NewConfig(cliCtx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(&config.Logger)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Reset readiness probe upon start-up
	if err := os.Remove(config.ControllerReadinessProbePath); err != nil {
		logger.Warn("Failed to clean up readiness file", "error", err, "path", config.ControllerReadinessProbePath)
	}

	dynamoClient, err := dynamodb.NewClient(config.AwsClient, logger)
	if err != nil {
		return fmt.Errorf("failed to create DynamoDB client: %w", err)
	}
	gethClient, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return fmt.Errorf("failed to create geth client: %w", err)
	}

	ctx := context.Background()

	contractDirectory, err := directory.NewContractDirectory(
		ctx,
		logger,
		gethClient,
		gethcommon.HexToAddress(config.EigenDAContractDirectoryAddress))
	if err != nil {
		return fmt.Errorf("failed to create contract directory: %w", err)
	}

	operatorStateRetrieverAddress, err :=
		contractDirectory.GetContractAddress(ctx, directory.OperatorStateRetriever)
	if err != nil {
		return fmt.Errorf("failed to get OperatorStateRetriever address: %w", err)
	}
	serviceManagerAddress, err :=
		contractDirectory.GetContractAddress(ctx, directory.ServiceManager)
	if err != nil {
		return fmt.Errorf("failed to get ServiceManager address: %w", err)
	}
	registryCoordinatorAddress, err :=
		contractDirectory.GetContractAddress(ctx, directory.RegistryCoordinator)
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

	var userAccountRemapping map[string]string
	if config.UserAccountRemappingFilePath != "" {
		userAccountRemapping, err = nameremapping.LoadNameRemapping(config.UserAccountRemappingFilePath)
		if err != nil {
			logger.Error("Failed to load user account remapping", "error", err)
		} else {
			logger.Info("Loaded user account remapping",
				"count", len(userAccountRemapping),
				"mappings", nameremapping.FormatMappings(userAccountRemapping))
		}
	}

	var validatorIdRemapping map[string]string
	if config.ValidatorIdRemappingFilePath != "" {
		validatorIdRemapping, err = nameremapping.LoadNameRemapping(
			config.ValidatorIdRemappingFilePath)
		if err != nil {
			logger.Error("Failed to load validator ID remapping", "error", err)
		} else {
			logger.Info("Loaded validator ID remapping",
				"count", len(validatorIdRemapping),
				"mappings", nameremapping.FormatMappings(validatorIdRemapping))
		}
	}

	encoderClient, err := encoder.NewEncoderClientV2(config.EncodingManager.EncoderAddress)
	if err != nil {
		return fmt.Errorf("failed to create encoder client: %v", err)
	}
	encodingPool := workerpool.New(config.EncodingManager.NumConcurrentRequests)
	encodingManagerBlobSet := controller.NewBlobSet()
	encodingManager, err := controller.NewEncodingManager(
		&config.EncodingManager,
		time.Now,
		blobMetadataStore,
		encodingPool,
		encoderClient,
		chainReader,
		logger,
		metricsRegistry,
		encodingManagerBlobSet,
		controllerLivenessChan,
		userAccountRemapping,
	)
	if err != nil {
		return fmt.Errorf("failed to create encoding manager: %v", err)
	}

	sigAgg, err := core.NewStdSignatureAggregator(logger, chainReader)
	if err != nil {
		return fmt.Errorf("failed to create signature aggregator: %v", err)
	}
	dispatcherPool := workerpool.New(config.NumConcurrentRequests)
	chainState := eth.NewChainState(chainReader, gethClient)
	var ics core.IndexedChainState
	if config.UseGraph {
		logger.Info("Using graph node")

		logger.Info("Connecting to subgraph", "url", config.ChainState.Endpoint)
		ics = thegraph.MakeIndexedChainState(config.ChainState, chainState, logger)
	} else {
		return fmt.Errorf("built-in indexer is deprecated and will be removed soon, please use UseGraph=true")
	}

	var requestSigner clients.DispersalRequestSigner
	if config.DisperserStoreChunksSigningDisabled {
		logger.Warn("StoreChunks() signing is disabled")
	} else {
		requestSigner, err = clients.NewDispersalRequestSigner(
			ctx,
			config.DispersalRequestSigner,
		)
		if err != nil {
			return fmt.Errorf("failed to create request signer: %v", err)
		}
	}

	nodeClientManager, err := controller.NewNodeClientManager(
		config.NodeClientCacheSize,
		requestSigner,
		config.DisperserID,
		logger)
	if err != nil {
		return fmt.Errorf("failed to create node client manager: %v", err)
	}
	beforeDispatch := func(blobKey corev2.BlobKey) error {
		encodingManagerBlobSet.RemoveBlob(blobKey)
		return nil
	}
	dispatcherBlobSet := controller.NewBlobSet()

	batchMetadataManager, err := metadata.NewBatchMetadataManager(
		ctx,
		logger,
		gethClient,
		ics,
		registryCoordinatorAddress,
		config.BatchMetadataUpdatePeriod,
		config.FinalizationBlockDelay,
	)
	if err != nil {
		return fmt.Errorf("failed to create batch metadata manager: %w", err)
	}

	signingRateTracker, err := signingrate.NewSigningRateTracker(
		logger,
		config.SigningRateRetentionPeriod,
		config.SigningRateBucketSpan,
		time.Now)
	if err != nil {
		return fmt.Errorf("failed to create signing rate tracker: %w", err)
	}
	signingRateTracker = signingrate.NewThreadsafeSigningRateTracker(ctx, signingRateTracker)

	signingRateStorage, err := signingrate.NewDynamoSigningRateStorage(
		ctx,
		logger,
		dynamoClient.GetAwsClient(),
		config.SigningRateDynamoDbTableName)
	if err != nil {
		return fmt.Errorf("failed to create signing rate storage: %w", err)
	}

	// Load existing signing rate data from persistent storage.
	err = signingrate.LoadSigningRateDataFromStorage(
		ctx,
		logger,
		signingRateTracker,
		signingRateStorage,
		config.SigningRateRetentionPeriod,
	)
	if err != nil {
		return fmt.Errorf("failed to load signing rate data from storage: %w", err)
	}

	// Periodically flush signing rate data to persistent storage.
	go signingrate.SigningRateStorageFlusher(
		ctx,
		logger,
		signingRateTracker,
		signingRateStorage,
		config.SigningRateFlushPeriod,
	)

	dispatcher, err := controller.NewController(
		ctx,
		config,
		time.Now,
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
		signingRateTracker,
		userAccountRemapping,
		validatorIdRemapping,
	)
	if err != nil {
		return fmt.Errorf("failed to create dispatcher: %v", err)
	}

	err = controller.RecoverState(ctx, blobMetadataStore, logger)
	if err != nil {
		return fmt.Errorf("failed to recover state: %v", err)
	}

	err = encodingManager.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start encoding manager: %v", err)
	}

	err = dispatcher.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start dispatcher: %v", err)
	}

	// nolint: nestif
	if config.Server.EnableServer {
		logger.Info("Controller gRPC server ENABLED", "port", config.Server.GrpcPort)
		var paymentAuthorizationHandler *controllerpayments.PaymentAuthorizationHandler
		if config.Server.EnablePaymentAuthentication {
			logger.Info("Payment authentication ENABLED - building payment authorization handler")
			paymentAuthorizationHandler, err = controller.BuildPaymentAuthorizationHandler(
				ctx,
				logger,
				config.PaymentAuthorization,
				contractDirectory,
				gethClient,
				dynamoClient.GetAwsClient(),
				metricsRegistry,
				userAccountRemapping,
			)
			if err != nil {
				return fmt.Errorf("build payment authorization handler: %w", err)
			}
		} else {
			logger.Warn("Payment authentication DISABLED - payment requests will fail")
		}

		// Create listener for the gRPC server
		addr := fmt.Sprintf("0.0.0.0:%d", config.Server.GrpcPort)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to create listener: %w", err)
		}

		grpcServer, err := server.NewServer(
			ctx,
			config.Server,
			logger,
			metricsRegistry,
			paymentAuthorizationHandler,
			listener,
			signingRateTracker)
		if err != nil {
			return fmt.Errorf("create gRPC server: %w", err)
		}

		go func() {
			logger.Info("Starting controller gRPC server", "address", listener.Addr().String())
			if err := grpcServer.Start(); err != nil {
				panic(fmt.Sprintf("gRPC server failed: %v", err))
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

	// Start heartbeat monitor
	go func() {
		err := healthcheck.NewHeartbeatMonitor(
			logger,
			controllerLivenessChan,
			config.HeartbeatMonitor,
		)
		if err != nil {
			logger.Warn("Heartbeat monitor failed", "err", err)
		}
	}()

	return nil
}
