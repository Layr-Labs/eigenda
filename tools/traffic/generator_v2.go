package traffic

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/workers"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Generator simulates read/write traffic to the DA service.
//
//	┌------------┐                                           ┌------------┐
//	|   writer   |-┐             ┌----------------┐          |   reader   |-┐
//	└------------┘ |-┐  -------> | status tracker | -------> └------------┘ |-┐
//	  └------------┘ |           └----------------┘            └------------┘ |
//	    └------------┘                                           └------------┘
//
// The traffic generator is built from three principal components: one or more writers
// that write blobs, a statusTracker that polls the disperser service until blobs are confirmed,
// and one or more readers that read blobs.
//
// When a writer finishes writing a blob, it sends information about that blob to the statusTracker.
// When the statusTracker observes that a blob has been confirmed, it sends information about the blob
// to the readers. The readers only attempt to read blobs that have been confirmed by the statusTracker.
type Generator struct {
	ctx              *context.Context
	cancel           *context.CancelFunc
	waitGroup        *sync.WaitGroup
	generatorMetrics metrics.Metrics
	logger           logging.Logger
	disperserClient  clients.DisperserClient
	// eigenDAClient    *clients.EigenDAClient #TODO: Add this back in when the client is implemented
	config *config.Config

	writers []*workers.BlobWriter
}

func NewTrafficGeneratorV2(config *config.Config) (*Generator, error) {
	logger, err := common.NewLogger(config.LoggingConfig)
	if err != nil {
		return nil, err
	}

	var signer *auth.LocalBlobRequestSigner
	if config.SignerPrivateKey != "" {
		signer = auth.NewLocalBlobRequestSigner(config.SignerPrivateKey)
	} else {
		logger.Error("signer private key is required")
		return nil, fmt.Errorf("signer private key is required")
	}

	signerAccountId, err := signer.GetAccountID()
	if err != nil {
		return nil, fmt.Errorf("error getting account ID: %w", err)
	}
	accountId := gethcommon.HexToAddress(signerAccountId)
	logger.Info("Initializing traffic generator", "accountId", accountId)

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}

	generatorMetrics := metrics.NewMetrics(
		config.MetricsHTTPPort,
		logger,
	)

	disperserClient, err := clients.NewDisperserClient(config.DisperserClientConfig, signer, nil, nil)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("new disperser-client: %w", err)
	}

	writers := make([]*workers.BlobWriter, 0)
	for i := 0; i < int(config.BlobWriterConfig.NumWriteInstances); i++ {
		writer := workers.NewBlobWriter(
			&ctx,
			&config.BlobWriterConfig,
			&waitGroup,
			logger,
			disperserClient,
			generatorMetrics)
		writers = append(writers, &writer)
	}

	return &Generator{
		ctx:              &ctx,
		cancel:           &cancel,
		waitGroup:        &waitGroup,
		generatorMetrics: generatorMetrics,
		logger:           logger,
		disperserClient:  disperserClient,
		config:           config,
		writers:          writers,
	}, nil
}

// Start instantiates goroutines that generate read/write traffic.
func (generator *Generator) Start() error {
	// Start metrics server
	generator.generatorMetrics.Start()

	// Start writers
	generator.logger.Info("Starting writers")
	for _, writer := range generator.writers {
		generator.logger.Info("Starting writer", "writer", writer)
		writer.Start()
		time.Sleep(generator.config.InstanceLaunchInterval)
	}

	// Wait for context cancellation to keep the process running
	<-(*generator.ctx).Done()
	generator.logger.Info("Generator received stop signal")
	return nil
}

func (generator *Generator) Stop() error {
	// Cancel context to stop all workers
	(*generator.cancel)()

	// Set a timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown metrics server
	if err := generator.generatorMetrics.Shutdown(); err != nil {
		generator.logger.Error("Failed to shutdown metrics server", "err", err)
	}

	// Wait for all workers with timeout
	done := make(chan struct{})
	go func() {
		generator.waitGroup.Wait()
		close(done)
	}()

	select {
	case <-done:
		generator.logger.Info("All workers shut down gracefully")
		return nil
	case <-shutdownCtx.Done():
		generator.logger.Warn("Shutdown timed out, forcing exit")
		return fmt.Errorf("shutdown timed out after 10 seconds")
	}
}
