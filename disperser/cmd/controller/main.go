package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand/ondemandvalidation"
	"github.com/Layr-Labs/eigenda/core/payments/reservation/reservationvalidation"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/disperser/cmd/controller/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/disperser/controller/metadata"
	"github.com/Layr-Labs/eigenda/disperser/controller/metrics"
	payments "github.com/Layr-Labs/eigenda/disperser/controller/payments"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gethcommon "github.com/ethereum/go-ethereum/common"
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
	encodingPool := workerpool.New(config.EncodingManagerConfig.NumConcurrentRequests)
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
	dispatcherPool := workerpool.New(config.DispatcherConfig.NumConcurrentRequests)
	chainState := eth.NewChainState(chainReader, gethClient)
	var ics core.IndexedChainState
	if config.UseGraph {
		logger.Info("Using graph node")

		logger.Info("Connecting to subgraph", "url", config.ChainStateConfig.Endpoint)
		ics = thegraph.MakeIndexedChainState(config.ChainStateConfig, chainState, logger)
	} else {
		return fmt.Errorf("built-in indexer is deprecated and will be removed soon, please use UseGraph=true")
	}
	logger.Info("Using graph node")
	logger.Info("Connecting to subgraph", "url", config.ChainStateConfig.Endpoint)

	var requestSigner clients.DispersalRequestSigner
	if config.DisperserStoreChunksSigningDisabled {
		logger.Warn("StoreChunks() signing is disabled")
	} else {
		requestSigner, err = clients.NewDispersalRequestSigner(
			context.Background(),
			config.DispersalRequestSignerConfig,
		)
		if err != nil {
			return fmt.Errorf("failed to create request signer: %v", err)
		}
	}

	nodeClientManager, err := controller.NewNodeClientManager(
		config.DispatcherConfig.NodeClientCacheSize,
		requestSigner,
		logger)
	if err != nil {
		return fmt.Errorf("failed to create signature aggregator: %w", err)
	}

	// Create indexed chain state
	chainState := eth.NewChainState(chainReader, gethClient)
	ics := thegraph.MakeIndexedChainState(config.ChainStateConfig, chainState, logger)

	// Create node client manager
	nodeClientManager, err := controller.NewNodeClientManager(
		config.DispatcherConfig.NodeClientCacheSize,
		requestSigner,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to create node client manager: %w", err)
	}

	// Create batch metadata manager
	batchMetadataManager, err := metadata.NewBatchMetadataManager(
		ctx,
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

	// Build payment authorization handler if needed
	var paymentAuthorizationHandler *payments.PaymentAuthorizationHandler
	if config.ServerConfig.EnableServer && config.ServerConfig.EnablePaymentAuthentication {
		logger.Info("Payment authentication ENABLED - building payment authorization handler")
		paymentAuthorizationHandler, err = buildPaymentAuthorizationHandler(
			ctx,
			logger,
			config.OnDemandConfig,
			config.ReservationConfig,
			contractDirectory,
			gethClient,
			dynamoClient.GetAwsClient(),
			metricsRegistry,
		)
		if err != nil {
			return fmt.Errorf("build payment authorization handler: %w", err)
		}
	}

	// Create controller
	c, err := controller.NewController(
		logger,
		metadataStore,
		chainReader,
		encoderClient,
		ics,
		batchMetadataManager,
		sigAgg,
		nodeClientManager,
		metricsRegistry,
		&config.EncodingManagerConfig,
		&config.DispatcherConfig,
		&config.ServerConfig,
		paymentAuthorizationHandler,
	)
	if err != nil {
		return fmt.Errorf("failed to create controller: %w", err)
	}

		go func() {
			logger.Info("Starting controller gRPC server", "port", config.ServerConfig.GrpcPort)
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
			config.HeartbeatMonitorConfig,
		)
		if err != nil {
			logger.Warn("Heartbeat monitor failed", "err", err)
		}
	}()

	return nil
}

func buildPaymentAuthorizationHandler(
	ctx context.Context,
	logger logging.Logger,
	onDemandConfig ondemandvalidation.OnDemandLedgerCacheConfig,
	reservationConfig reservationvalidation.ReservationLedgerCacheConfig,
	contractDirectory *directory.ContractDirectory,
	ethClient common.EthClient,
	awsDynamoClient *awsdynamodb.Client,
	metricsRegistry *prometheus.Registry,
) (*payments.PaymentAuthorizationHandler, error) {
	paymentVaultAddress, err := contractDirectory.GetContractAddress(ctx, directory.PaymentVault)
	if err != nil {
		return nil, fmt.Errorf("get PaymentVault address: %w", err)
	}

	paymentVault, err := vault.NewPaymentVault(logger, ethClient, paymentVaultAddress)
	if err != nil {
		return nil, fmt.Errorf("create payment vault: %w", err)
	}

	globalSymbolsPerSecond, err := paymentVault.GetGlobalSymbolsPerSecond(ctx)
	if err != nil {
		return nil, fmt.Errorf("get global symbols per second: %w", err)
	}

	globalRatePeriodInterval, err := paymentVault.GetGlobalRatePeriodInterval(ctx)
	if err != nil {
		return nil, fmt.Errorf("get global rate period interval: %w", err)
	}

	onDemandMeterer := meterer.NewOnDemandMeterer(
		globalSymbolsPerSecond,
		globalRatePeriodInterval,
		time.Now,
		meterer.NewOnDemandMetererMetrics(
			metricsRegistry,
			metrics.Namespace,
			metrics.AuthorizePaymentsSubsystem,
		),
	)

	onDemandValidator, err := ondemandvalidation.NewOnDemandPaymentValidator(
		ctx,
		logger,
		onDemandConfig,
		paymentVault,
		awsDynamoClient,
		ondemandvalidation.NewOnDemandValidatorMetrics(
			metricsRegistry,
			metrics.Namespace,
			metrics.AuthorizePaymentsSubsystem,
		),
		ondemandvalidation.NewOnDemandCacheMetrics(
			metricsRegistry,
			metrics.Namespace,
			metrics.AuthorizePaymentsSubsystem,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create on-demand payment validator: %w", err)
	}

	reservationValidator, err := reservationvalidation.NewReservationPaymentValidator(
		ctx,
		logger,
		reservationConfig,
		paymentVault,
		time.Now,
		reservationvalidation.NewReservationValidatorMetrics(
			metricsRegistry,
			metrics.Namespace,
			metrics.AuthorizePaymentsSubsystem,
		),
		reservationvalidation.NewReservationCacheMetrics(
			metricsRegistry,
			metrics.Namespace,
			metrics.AuthorizePaymentsSubsystem,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create reservation payment validator: %w", err)
	}

	return payments.NewPaymentAuthorizationHandler(
		onDemandMeterer,
		onDemandValidator,
		reservationValidator,
	), nil
}
