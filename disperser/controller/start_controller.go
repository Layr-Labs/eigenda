package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand/ondemandvalidation"
	"github.com/Layr-Labs/eigenda/core/payments/reservation/reservationvalidation"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller/metadata"
	"github.com/Layr-Labs/eigenda/disperser/controller/metrics"
	controllerpayments "github.com/Layr-Labs/eigenda/disperser/controller/payments"
	"github.com/Layr-Labs/eigenda/disperser/controller/server"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigensdk-go/logging"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gammazero/workerpool"
	"github.com/prometheus/client_golang/prometheus"
)

// Controller holds the running controller components along with the server.
// Server is optional and can be nil.
type Controller struct {
	EncodingManager *EncodingManager
	Dispatcher      *Dispatcher
	LivenessChan    chan healthcheck.HeartbeatMessage

	Server *server.Server
}

// StartController creates and starts the encoding manager and dispatcher.
// This is the shared logic between the main controller binary and test harnesses.
func StartController(
	ctx context.Context,
	logger logging.Logger,
	ethClient common.EthClient,
	dynamoClient dynamodb.Client,
	awsDynamoClient *awsdynamodb.Client,
	metadataTableName string,
	metricsRegistry *prometheus.Registry,
	requestSigner clients.DispersalRequestSigner,
	// Config
	operatorStateRetrieverAddress gethcommon.Address,
	serviceManagerAddress gethcommon.Address,
	registryCoordinatorAddress gethcommon.Address,
	operatorStateSubgraphURL string,
	encodingManagerConfig *EncodingManagerConfig,
	dispatcherConfig *DispatcherConfig,
	numConcurrentEncodingRequests int,
	numConcurrentDispersalRequests int,
	nodeClientCacheSize int,

	// Chain state config
	chainStateConfig thegraph.Config,

	// Optional components (can be nil/empty for test harness)
	metricsServer *http.Server,
	readinessProbePath string,
	heartbeatMonitorConfig healthcheck.HeartbeatMonitorConfig,

	// Server config (can be nil for test harness)
	serverConfig *server.Config,
	onDemandConfig *ondemandvalidation.OnDemandLedgerCacheConfig,
	reservationConfig *reservationvalidation.ReservationLedgerCacheConfig,
	contractDirectory *directory.ContractDirectory,
) (*Controller, error) {
	// Create metadata store
	baseMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, metadataTableName)
	metadataStore := blobstore.NewInstrumentedMetadataStore(baseMetadataStore, blobstore.InstrumentedMetadataStoreConfig{
		ServiceName: "controller",
		Registry:    metricsRegistry,
		Backend:     blobstore.BackendDynamoDB,
	})

	// Create chain reader
	chainReader, err := eth.NewReader(
		logger,
		ethClient,
		operatorStateRetrieverAddress.Hex(),
		serviceManagerAddress.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to create chain reader: %w", err)
	}

	// Create encoder client
	encoderClient, err := encoder.NewEncoderClientV2(encodingManagerConfig.EncoderAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder client: %w", err)
	}

	// Create encoding manager
	encodingPool := workerpool.New(numConcurrentEncodingRequests)
	encodingManagerBlobSet := NewBlobSet()

	// Heartbeat monitor
	controllerLivenessChan := make(chan healthcheck.HeartbeatMessage, 10)
	encodingManager, err := NewEncodingManager(
		encodingManagerConfig,
		metadataStore,
		encodingPool,
		encoderClient,
		chainReader,
		logger,
		metricsRegistry,
		encodingManagerBlobSet,
		controllerLivenessChan,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoding manager: %w", err)
	}

	// Create signature aggregator
	sigAgg, err := core.NewStdSignatureAggregator(logger, chainReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature aggregator: %w", err)
	}

	// Create dispatcher components
	dispatcherPool := workerpool.New(numConcurrentDispersalRequests)
	chainState := eth.NewChainState(chainReader, ethClient)
	ics := thegraph.MakeIndexedChainState(chainStateConfig, chainState, logger)

	nodeClientManager, err := NewNodeClientManager(nodeClientCacheSize, requestSigner, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create node client manager: %w", err)
	}

	beforeDispatch := func(blobKey corev2.BlobKey) error {
		encodingManagerBlobSet.RemoveBlob(blobKey)
		return nil
	}
	dispatcherBlobSet := NewBlobSet()

	batchMetadataManager, err := metadata.NewBatchMetadataManager(
		ctx,
		logger,
		ethClient,
		ics,
		registryCoordinatorAddress,
		dispatcherConfig.BatchMetadataUpdatePeriod,
		dispatcherConfig.FinalizationBlockDelay,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch metadata manager: %w", err)
	}

	dispatcher, err := NewDispatcher(
		dispatcherConfig,
		metadataStore,
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
		return nil, fmt.Errorf("failed to create dispatcher: %w", err)
	}

	// Recover state
	if err := RecoverState(ctx, metadataStore, logger); err != nil {
		return nil, fmt.Errorf("failed to recover state: %w", err)
	}

	// Start encoding manager
	if err := encodingManager.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start encoding manager: %w", err)
	}

	// Start dispatcher
	if err := dispatcher.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start dispatcher: %w", err)
	}

	// Start heartbeat monitor
	if heartbeatMonitorConfig.FilePath != "" {
		go func() {
			err := healthcheck.HeartbeatMonitor(
				logger,
				controllerLivenessChan,
				healthcheck.HeartbeatMonitorConfig{
					FilePath:         heartbeatMonitorConfig.FilePath,
					MaxStallDuration: heartbeatMonitorConfig.MaxStallDuration,
				},
			)
			if err != nil {
				logger.Warn("Heartbeat monitor exited with error", "err", err)
			}
		}()
	}

	// Start metrics server if provided
	if metricsServer != nil {
		go func() {
			err := metricsServer.ListenAndServe()
			if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
				logger.Errorf("metrics server error: %v", err)
			}
		}()
	}

	// Create readiness probe file once the controller starts successfully
	if readinessProbePath != "" {
		if _, err := os.Create(readinessProbePath); err != nil {
			logger.Warn("Failed to create readiness file", "error", err, "path", readinessProbePath)
		}
	}

	// Create health probe file
	if heartbeatMonitorConfig.FilePath != "" {
		if _, err := os.Create(heartbeatMonitorConfig.FilePath); err != nil {
			logger.Warn("Failed to create health probe file", "error", err, "path", heartbeatMonitorConfig.FilePath)
		}
	}

	// Start gRPC server if enabled
	var grpcServer *server.Server
	if serverConfig != nil && serverConfig.EnableServer {
		logger.Info("Controller gRPC server ENABLED", "port", serverConfig.GrpcPort)
		var paymentAuthorizationHandler *controllerpayments.PaymentAuthorizationHandler
		if serverConfig.EnablePaymentAuthentication {
			logger.Info("Payment authentication ENABLED - building payment authorization handler")
			var err error
			paymentAuthorizationHandler, err = buildPaymentAuthorizationHandler(
				ctx,
				logger,
				*onDemandConfig,
				*reservationConfig,
				contractDirectory,
				ethClient,
				awsDynamoClient,
				metricsRegistry,
			)
			if err != nil {
				return nil, fmt.Errorf("build payment authorization handler: %w", err)
			}
		} else {
			logger.Warn("Payment authentication DISABLED - payment requests will fail")
		}

		var err error
		grpcServer, err = server.NewServer(
			ctx,
			*serverConfig,
			logger,
			metricsRegistry,
			paymentAuthorizationHandler)
		if err != nil {
			return nil, fmt.Errorf("create gRPC server: %w", err)
		}

		go func() {
			logger.Info("Starting controller gRPC server", "port", serverConfig.GrpcPort)
			if err := grpcServer.Start(); err != nil {
				panic(fmt.Sprintf("gRPC server failed: %v", err))
			}
		}()
	} else if serverConfig != nil && !serverConfig.EnableServer {
		logger.Info("Controller gRPC server disabled")
	}

	return &Controller{
		EncodingManager: encodingManager,
		Dispatcher:      dispatcher,
		LivenessChan:    controllerLivenessChan,
		Server:          grpcServer,
	}, nil
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
) (*controllerpayments.PaymentAuthorizationHandler, error) {
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

	return controllerpayments.NewPaymentAuthorizationHandler(
		onDemandMeterer,
		onDemandValidator,
		reservationValidator,
	), nil
}
