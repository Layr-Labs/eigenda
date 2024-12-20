package traffic

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	logger           *logging.Logger
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
		config.WorkerConfig.MetricsBlacklist,
		config.WorkerConfig.MetricsFuzzyBlacklist)

	uncertifiedKeyChannel := make(chan *workers.UncertifiedKey, 100)

	// TODO: create a dedicated reservation for traffic generator
	disperserClient, err := clients.NewDisperserClient(config.DisperserClientConfig, signer, nil, nil)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("new disperser-client: %w", err)
	}

	writers := make([]*workers.BlobWriter, 0)
	for i := 0; i < int(config.WorkerConfig.NumWriteInstances); i++ {
		writer := workers.NewBlobWriter(
			&ctx,
			&waitGroup,
			logger,
			&config.WorkerConfig,
			disperserClient,
			uncertifiedKeyChannel,
			generatorMetrics)
		writers = append(writers, &writer)
	}

	return &Generator{
		ctx:              &ctx,
		cancel:           &cancel,
		waitGroup:        &waitGroup,
		generatorMetrics: generatorMetrics,
		logger:           &logger,
		disperserClient:  disperserClient,
		config:           config,
		writers:          writers,
	}, nil
}

// Start instantiates goroutines that generate read/write traffic, continues until a SIGTERM is observed.
func (generator *Generator) Start() error {

	generator.generatorMetrics.Start()

	// generator.statusTracker.Start()

	for _, writer := range generator.writers {
		writer.Start()
		time.Sleep(generator.config.InstanceLaunchInterval)
	}

	// for _, reader := range generator.readers {
	// 	reader.Start()
	// 	time.Sleep(generator.config.InstanceLaunchInterval)
	// }

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	(*generator.cancel)()
	generator.waitGroup.Wait()
	return nil
}
