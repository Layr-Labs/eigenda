package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/disperser/controller/metadata"
	controllerpayments "github.com/Layr-Labs/eigenda/disperser/controller/payments"
	"github.com/Layr-Labs/eigenda/disperser/controller/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gammazero/workerpool"
)

var (
	version   string
	gitCommit string
	gitDate   string
)

func main() {

	cfg, err := config.Bootstrap(controller.DefaultControllerConfig)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	err = RunController(context.Background(), cfg)
	if err != nil {
		panic(fmt.Sprintf("controller setup failed: %v", err))
	}
}

func RunController(ctx context.Context, cfg *controller.ControllerConfig) error {

	loggerConfig := common.DefaultLoggerConfig()
	loggerConfig.Format = common.LogFormat(cfg.LogFormat)
	loggerConfig.HandlerOpts.Level = slog.Level(cfg.LogLevel)

	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Reset readiness probe upon start-up
	if err := os.Remove(cfg.ControllerReadinessProbePath); err != nil {
		logger.Warn("Failed to clean up readiness file", "error", err, "path", cfg.ControllerReadinessProbePath)
	}

	dynamoClient, err := dynamodb.NewClient(cfg.AwsClientConfig, logger)
	if err != nil {
		return fmt.Errorf("failed to create DynamoDB client: %w", err)
	}
	gethClient, err := geth.NewMultiHomingClient(cfg.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return fmt.Errorf("failed to create geth client: %w", err)
	}

	contractDirectory, err := directory.NewContractDirectory(
		ctx,
		logger,
		gethClient,
		gethcommon.HexToAddress(cfg.EigenDAContractDirectoryAddress))
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

	logger.Infof("Starting metrics server at port %d", cfg.MetricsPort)
	addr := fmt.Sprintf(":%d", cfg.MetricsPort)
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
		cfg.DynamoDBTableName,
	)
	blobMetadataStore := blobstore.NewInstrumentedMetadataStore(baseBlobMetadataStore, blobstore.InstrumentedMetadataStoreConfig{
		ServiceName: "controller",
		Registry:    metricsRegistry,
		Backend:     blobstore.BackendDynamoDB,
	})

	controllerLivenessChan := make(chan healthcheck.HeartbeatMessage, 10)

	encoderClient, err := encoder.NewEncoderClientV2(cfg.EncodingManagerConfig.EncoderAddress)
	if err != nil {
		return fmt.Errorf("failed to create encoder client: %v", err)
	}
	encodingPool := workerpool.New(cfg.EncodingManagerConfig.NumConcurrentRequests)
	encodingManagerBlobSet := controller.NewBlobSet()
	encodingManager, err := controller.NewEncodingManager(
		&cfg.EncodingManagerConfig,
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
	dispatcherPool := workerpool.New(cfg.NumConcurrentRequests)
	chainState := eth.NewChainState(chainReader, gethClient)
	var ics core.IndexedChainState
	if cfg.UseGraph {
		logger.Info("Using graph node")

		logger.Info("Connecting to subgraph", "url", cfg.ChainStateConfig.Endpoint)
		ics = thegraph.MakeIndexedChainState(cfg.ChainStateConfig, chainState, logger)
	} else {
		return fmt.Errorf("built-in indexer is deprecated and will be removed soon, please use UseGraph=true")
	}

	var requestSigner clients.DispersalRequestSigner
	if cfg.DisperserStoreChunksSigningDisabled {
		logger.Warn("StoreChunks() signing is disabled")
	} else {
		requestSigner, err = clients.NewDispersalRequestSigner(
			ctx,
			*cfg.DispersalRequestSignerConfig,
		)
		if err != nil {
			return fmt.Errorf("failed to create request signer: %v", err)
		}
	}

	nodeClientManager, err := controller.NewNodeClientManager(
		cfg.NodeClientCacheSize,
		requestSigner,
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
		cfg.BatchMetadataUpdatePeriod,
		cfg.FinalizationBlockDelay,
	)
	if err != nil {
		return fmt.Errorf("failed to create batch metadata manager: %w", err)
	}

	dispatcher, err := controller.NewDispatcher(
		cfg,
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

	// nolint:nestif
	if cfg.EnableGrpcServer {
		logger.Info("Controller gRPC server ENABLED", "port", cfg.GrpcPort)
		var paymentAuthorizationHandler *controllerpayments.PaymentAuthorizationHandler
		if cfg.EnablePaymentAuthentication {
			logger.Info("Payment authentication ENABLED - building payment authorization handler")
			paymentAuthorizationHandler, err = controller.BuildPaymentAuthorizationHandler(
				ctx,
				logger,
				cfg.PaymentAuthorizationConfig,
				contractDirectory,
				gethClient,
				dynamoClient.GetAwsClient(),
				metricsRegistry,
			)
			if err != nil {
				return fmt.Errorf("build payment authorization handler: %w", err)
			}
		} else {
			logger.Warn("Payment authentication DISABLED - payment requests will fail")
		}

		// Create listener for the gRPC server
		addr := fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to create listener: %w", err)
		}

		grpcServer, err := server.NewServer(
			ctx,
			cfg,
			logger,
			metricsRegistry,
			paymentAuthorizationHandler,
			listener)
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
	if _, err := os.Create(cfg.ControllerReadinessProbePath); err != nil {
		logger.Warn("Failed to create readiness file", "error", err, "path", cfg.ControllerReadinessProbePath)
	}

	// Start heartbeat monitor
	go func() {
		err := healthcheck.NewHeartbeatMonitor(
			logger,
			controllerLivenessChan,
			cfg.HeartbeatMonitorConfig,
		)
		if err != nil {
			logger.Warn("Heartbeat monitor failed", "err", err)
		}
	}()

	return nil
}
