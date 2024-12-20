package traffic

import (
	"context"
	"fmt"
	"sync"
	"time"

	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"

	// "github.com/Layr-Labs/eigenda/common/geth"
	// "github.com/Layr-Labs/eigenda/core/auth"
	// "github.com/Layr-Labs/eigenda/core/eth"
	// "github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	// retrivereth "github.com/Layr-Labs/eigenda/retriever/eth"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	trafficconfig "github.com/Layr-Labs/eigenda/tools/traffic/config"
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
type WriterGroup struct {
	name    string
	writers []*workers.BlobWriter
	cancels map[*workers.BlobWriter]context.CancelFunc
}

type Generator struct {
	ctx              *context.Context
	cancel           *context.CancelFunc
	waitGroup        *sync.WaitGroup
	generatorMetrics metrics.Metrics
	logger           logging.Logger
	disperserClient  clientsv2.DisperserClient
	config           *config.Config
	writerGroups     map[string]*WriterGroup
	configManager    *config.RuntimeConfigManager
	mu               sync.RWMutex
}

func NewTrafficGeneratorV2(config *config.Config) (*Generator, error) {
	logger, err := common.NewLogger(config.LoggingConfig)
	if err != nil {
		return nil, err
	}

	signer, err := auth.NewLocalBlobRequestSigner(config.SignerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("new local blob request signer: %w", err)
	}

	signerAccountId, err := signer.GetAccountID()
	if err != nil {
		return nil, fmt.Errorf("error getting account ID: %w", err)
	}
	accountId := gethcommon.HexToAddress(signerAccountId)
	logger.Info("Initializing traffic generator", "accountId", accountId)

	if config.RuntimeConfigPath == "" {
		return nil, fmt.Errorf("runtime config path is required")
	}

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}

	generatorMetrics := metrics.NewMetrics(
		config.MetricsHTTPPort,
		logger,
	)

	disperserClient, err := clientsv2.NewDisperserClient(config.DisperserClientConfig, signer, nil, nil)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("new disperser-client: %w", err)
	}

	generator := &Generator{
		ctx:              &ctx,
		cancel:           &cancel,
		waitGroup:        &waitGroup,
		generatorMetrics: generatorMetrics,
		logger:           logger,
		disperserClient:  disperserClient,
		config:           config,
		writerGroups:     make(map[string]*WriterGroup),
	}

	// Initialize runtime config manager
	configManager, err := trafficconfig.NewRuntimeConfigManager(config.RuntimeConfigPath, generator.handleConfigUpdate)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize runtime config manager: %w", err)
	}
	generator.configManager = configManager

	return generator, nil
}

// handleConfigUpdate is called when the runtime configuration changes
func (generator *Generator) handleConfigUpdate(runtimeConfig *trafficconfig.RuntimeConfig) {
	generator.mu.Lock()
	defer generator.mu.Unlock()

	generator.logger.Info("Received runtime configuration update")

	// Track existing groups to identify which ones to remove
	existingGroups := make(map[string]bool)
	for name := range generator.writerGroups {
		existingGroups[name] = true
	}

	// Update or create writer groups
	for _, groupConfig := range runtimeConfig.WriterGroups {
		delete(existingGroups, groupConfig.Name)

		writerConfig := &trafficconfig.BlobWriterConfig{
			NumWriteInstances:    groupConfig.NumWriteInstances,
			WriteRequestInterval: groupConfig.WriteRequestInterval,
			DataSize:             groupConfig.DataSize,
			RandomizeBlobs:       groupConfig.RandomizeBlobs,
			WriteTimeout:         groupConfig.WriteTimeout,
			CustomQuorums:        groupConfig.CustomQuorums,
		}

		group, exists := generator.writerGroups[groupConfig.Name]
		if !exists {
			group = &WriterGroup{
				name:    groupConfig.Name,
				writers: make([]*workers.BlobWriter, 0),
				cancels: make(map[*workers.BlobWriter]context.CancelFunc),
			}
			generator.writerGroups[groupConfig.Name] = group
		}

		// Update writer count
		currentWriters := len(group.writers)
		targetWriters := int(groupConfig.NumWriteInstances)

		// Scale down if needed
		if targetWriters < currentWriters {
			for i := targetWriters; i < currentWriters; i++ {
				if cancel, exists := group.cancels[group.writers[i]]; exists {
					cancel()
					delete(group.cancels, group.writers[i])
				}
			}
			group.writers = group.writers[:targetWriters]
		}

		// Scale up if needed
		if targetWriters > currentWriters {
			for i := currentWriters; i < targetWriters; i++ {
				writerCtx, writerCancel := context.WithCancel(*generator.ctx)
				writer := workers.NewBlobWriter(
					groupConfig.Name,
					&writerCtx,
					writerConfig,
					generator.waitGroup,
					generator.logger,
					generator.disperserClient,
					generator.generatorMetrics)
				group.writers = append(group.writers, &writer)
				group.cancels[&writer] = writerCancel
				writer.Start()
			}
		}

		// Update configuration for existing writers
		for _, writer := range group.writers[:min(currentWriters, targetWriters)] {
			writer.UpdateConfig(writerConfig)
		}
	}

	// Remove any groups that are no longer in the config
	for name := range existingGroups {
		group := generator.writerGroups[name]
		for _, writer := range group.writers {
			if cancel, exists := group.cancels[writer]; exists {
				cancel()
			}
		}
		delete(generator.writerGroups, name)
	}

	// cs := eth.NewChainState(tx, gethClient)

	// var assignmentCoordinator core.AssignmentCoordinator = &core.StdAssignmentCoordinator{}

	// nodeClient := clients.NewNodeClient(config.NodeClientTimeout)

	// config.RetrievalClientConfig.EncoderConfig.LoadG2Points = true
	// v, err := verifier.NewVerifier(&config.RetrievalClientConfig.EncoderConfig, nil)
	// if err != nil {
	// 	panic(fmt.Sprintf("Unable to build statusTracker: %s", err))
	// }

	// retriever, err := clients.NewRetrievalClient(
	// 	logger,
	// 	cs,
	// 	assignmentCoordinator,
	// 	nodeClient,
	// 	v,
	// 	config.RetrievalClientConfig.NumConnections)

	// if err != nil {
	// 	panic(fmt.Sprintf("Unable to build retriever: %s", err))
	// }

	// chainClient := retrivereth.NewChainClient(gethClient, logger)

	// return retriever, chainClient
}

// Start instantiates goroutines that generate read/write traffic.
func (generator *Generator) Start() error {
	// Start metrics server
	if err := generator.generatorMetrics.Start(); err != nil {
		return fmt.Errorf("failed to start metrics server: %w", err)
	}

	// Start runtime config watcher if configured
	if generator.configManager != nil {
		generator.configManager.StartWatching(*generator.ctx)
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
