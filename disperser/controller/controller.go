package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller/metadata"
	controllerpayments "github.com/Layr-Labs/eigenda/disperser/controller/payments"
	"github.com/Layr-Labs/eigenda/disperser/controller/server"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gammazero/workerpool"
	"github.com/prometheus/client_golang/prometheus"
)

// Controller holds the running controller components along with the server.
// Server is optional and can be nil.
type Controller struct {
	EncodingManager *EncodingManager
	Dispatcher      *Dispatcher
	LivenessChan    chan healthcheck.HeartbeatMessage
	Server          *server.Server

	logger        logging.Logger
	metadataStore blobstore.MetadataStore
}

// NewController creates a new Controller with the given dependencies.
// This constructs all components but does not start them.
func NewController(
	logger logging.Logger,
	metadataStore blobstore.MetadataStore,
	chainReader core.Reader,
	encoderClient disperser.EncoderClientV2,
	indexedChainState core.IndexedChainState,
	batchMetadataManager metadata.BatchMetadataManager,
	signatureAggregator core.SignatureAggregator,
	nodeClientManager NodeClientManager,
	metricsRegistry *prometheus.Registry,
	encodingManagerConfig *EncodingManagerConfig,
	dispatcherConfig *DispatcherConfig,
	// Optional server dependencies
	serverConfig *server.Config,
	paymentAuthorizationHandler *controllerpayments.PaymentAuthorizationHandler,
) (*Controller, error) {
	// Create liveness channel shared between encoding manager and dispatcher
	controllerLivenessChan := make(chan healthcheck.HeartbeatMessage, 10)

	// Create worker pools
	encodingPool := workerpool.New(encodingManagerConfig.NumConcurrentRequests)
	dispatcherPool := workerpool.New(dispatcherConfig.NumConcurrentRequests)

	// Create blob sets
	encodingManagerBlobSet := NewBlobSet()
	dispatcherBlobSet := NewBlobSet()

	// Create encoding manager
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

	// Create beforeDispatch callback to coordinate between encoding manager and dispatcher
	beforeDispatch := func(blobKey corev2.BlobKey) error {
		encodingManagerBlobSet.RemoveBlob(blobKey)
		return nil
	}

	// Create dispatcher
	dispatcher, err := NewDispatcher(
		dispatcherConfig,
		metadataStore,
		dispatcherPool,
		indexedChainState,
		batchMetadataManager,
		signatureAggregator,
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

	// Create gRPC server if enabled
	var grpcServer *server.Server
	if serverConfig != nil && serverConfig.EnableServer {
		logger.Info("Controller gRPC server ENABLED", "port", serverConfig.GrpcPort)
		if serverConfig.EnablePaymentAuthentication {
			if paymentAuthorizationHandler == nil {
				return nil, fmt.Errorf("payment authentication enabled but payment authorization handler is nil")
			}
			logger.Info("Payment authentication ENABLED")
		} else {
			logger.Warn("Payment authentication DISABLED - payment requests will fail")
		}

		var err error
		grpcServer, err = server.NewServer(
			context.Background(),
			*serverConfig,
			logger,
			metricsRegistry,
			paymentAuthorizationHandler)
		if err != nil {
			return nil, fmt.Errorf("create gRPC server: %w", err)
		}
	} else if serverConfig != nil && !serverConfig.EnableServer {
		logger.Info("Controller gRPC server disabled")
	}

	return &Controller{
		EncodingManager: encodingManager,
		Dispatcher:      dispatcher,
		LivenessChan:    controllerLivenessChan,
		Server:          grpcServer,
		logger:          logger,
		metadataStore:   metadataStore,
	}, nil
}

// StartOptions contains optional dependencies needed to start the Controller.
type StartOptions struct {
	MetricsServer          *http.Server
	ReadinessProbePath     string
	HeartbeatMonitorConfig healthcheck.HeartbeatMonitorConfig
}

// Start recovers state and starts all controller components.
func (c *Controller) Start(ctx context.Context, opts StartOptions) error {
	// Recover state
	if err := RecoverState(ctx, c.metadataStore, c.logger); err != nil {
		return fmt.Errorf("failed to recover state: %w", err)
	}

	// Start encoding manager
	if err := c.EncodingManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start encoding manager: %w", err)
	}

	// Start dispatcher
	if err := c.Dispatcher.Start(ctx); err != nil {
		return fmt.Errorf("failed to start dispatcher: %w", err)
	}

	// Start heartbeat monitor if configured
	if opts.HeartbeatMonitorConfig.FilePath != "" {
		go func() {
			err := healthcheck.NewHeartbeatMonitor(
				c.logger,
				c.LivenessChan,
				healthcheck.HeartbeatMonitorConfig{
					FilePath:         opts.HeartbeatMonitorConfig.FilePath,
					MaxStallDuration: opts.HeartbeatMonitorConfig.MaxStallDuration,
				},
			)
			if err != nil {
				c.logger.Warn("Heartbeat monitor exited with error", "err", err)
			}
		}()
	}

	// Start metrics server if provided
	if opts.MetricsServer != nil {
		go func() {
			err := opts.MetricsServer.ListenAndServe()
			if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
				c.logger.Errorf("metrics server error: %v", err)
			}
		}()
	}

	// Create readiness probe file once the controller starts successfully
	if opts.ReadinessProbePath != "" {
		if _, err := os.Create(opts.ReadinessProbePath); err != nil {
			c.logger.Warn("Failed to create readiness file", "error", err, "path", opts.ReadinessProbePath)
		}
	}

	// Create health probe file
	if opts.HeartbeatMonitorConfig.FilePath != "" {
		if _, err := os.Create(opts.HeartbeatMonitorConfig.FilePath); err != nil {
			c.logger.Warn("Failed to create health probe file", "error", err, "path", opts.HeartbeatMonitorConfig.FilePath)
		}
	}

	// Start gRPC server if enabled
	if c.Server != nil {
		go func() {
			c.logger.Info("Starting controller gRPC server")
			if err := c.Server.Start(); err != nil {
				panic(fmt.Sprintf("gRPC server failed: %v", err))
			}
		}()
	}

	return nil
}
