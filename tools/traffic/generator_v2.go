package traffic

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	retrivereth "github.com/Layr-Labs/eigenda/retriever/eth"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/Layr-Labs/eigenda/tools/traffic/workers"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
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
// that write blobs, a statusTracker that polls the dispenser service until blobs are confirmed,
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
	eigenDAClient    *clients.EigenDAClient
	config           *config.Config

	writers       []*workers.BlobWriter
	statusTracker *workers.BlobStatusTracker
	readers       []*workers.BlobReader
}

func NewTrafficGeneratorV2(config *config.Config) (*Generator, error) {
	logger, err := common.NewLogger(config.LoggingConfig)
	if err != nil {
		return nil, err
	}

	var signer core.BlobRequestSigner
	if config.EigenDAClientConfig.SignerPrivateKeyHex != "" {
		signer = auth.NewLocalBlobRequestSigner(config.EigenDAClientConfig.SignerPrivateKeyHex)
	}

	logger2 := log.NewLogger(log.NewTerminalHandler(os.Stderr, true))
	client, err := clients.NewEigenDAClient(logger2, *config.EigenDAClientConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}

	generatorMetrics := metrics.NewMetrics(
		config.MetricsHTTPPort,
		logger,
		config.WorkerConfig.MetricsBlacklist,
		config.WorkerConfig.MetricsFuzzyBlacklist)

	blobTable := table.NewBlobStore()

	unconfirmedKeyChannel := make(chan *workers.UnconfirmedKey, 100)

	// TODO: create a dedicated reservation for traffic generator
	disperserClient, err := clients.NewDisperserClient(config.DisperserClientConfig, signer)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("new disperser-client: %w", err)
	}
	statusVerifier := workers.NewBlobStatusTracker(
		&ctx,
		&waitGroup,
		logger,
		&config.WorkerConfig,
		unconfirmedKeyChannel,
		blobTable,
		disperserClient,
		generatorMetrics)

	writers := make([]*workers.BlobWriter, 0)
	for i := 0; i < int(config.WorkerConfig.NumWriteInstances); i++ {
		writer := workers.NewBlobWriter(
			&ctx,
			&waitGroup,
			logger,
			&config.WorkerConfig,
			disperserClient,
			unconfirmedKeyChannel,
			generatorMetrics)
		writers = append(writers, &writer)
	}

	retriever, chainClient := buildRetriever(config)

	readers := make([]*workers.BlobReader, 0)
	for i := 0; i < int(config.WorkerConfig.NumReadInstances); i++ {
		reader := workers.NewBlobReader(
			&ctx,
			&waitGroup,
			logger,
			&config.WorkerConfig,
			retriever,
			chainClient,
			blobTable,
			generatorMetrics)
		readers = append(readers, &reader)
	}

	return &Generator{
		ctx:              &ctx,
		cancel:           &cancel,
		waitGroup:        &waitGroup,
		generatorMetrics: generatorMetrics,
		logger:           &logger,
		disperserClient:  disperserClient,
		eigenDAClient:    client,
		config:           config,
		writers:          writers,
		statusTracker:    &statusVerifier,
		readers:          readers,
	}, nil
}

// buildRetriever creates a retriever client for the traffic generator.
func buildRetriever(config *config.Config) (clients.RetrievalClient, retrivereth.ChainClient) {
	loggerConfig := common.DefaultLoggerConfig()

	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		panic(fmt.Sprintf("Unable to instantiate logger: %s", err))
	}

	gethClient, err := geth.NewMultiHomingClient(config.RetrievalClientConfig.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		panic(fmt.Sprintf("Unable to instantiate geth client: %s", err))
	}

	tx, err := eth.NewReader(
		logger,
		gethClient,
		config.RetrievalClientConfig.BLSOperatorStateRetrieverAddr,
		config.RetrievalClientConfig.EigenDAServiceManagerAddr)
	if err != nil {
		panic(fmt.Sprintf("Unable to instantiate transactor: %s", err))
	}

	cs := eth.NewChainState(tx, gethClient)

	chainState := thegraph.MakeIndexedChainState(*config.TheGraphConfig, cs, logger)

	var assignmentCoordinator core.AssignmentCoordinator = &core.StdAssignmentCoordinator{}

	nodeClient := clients.NewNodeClient(config.NodeClientTimeout)

	config.RetrievalClientConfig.EncoderConfig.LoadG2Points = true
	v, err := verifier.NewVerifier(&config.RetrievalClientConfig.EncoderConfig, nil)
	if err != nil {
		panic(fmt.Sprintf("Unable to build statusTracker: %s", err))
	}

	retriever, err := clients.NewRetrievalClient(
		logger,
		chainState,
		assignmentCoordinator,
		nodeClient,
		v,
		config.RetrievalClientConfig.NumConnections)

	if err != nil {
		panic(fmt.Sprintf("Unable to build retriever: %s", err))
	}

	chainClient := retrivereth.NewChainClient(gethClient, logger)

	return retriever, chainClient
}

// Start instantiates goroutines that generate read/write traffic, continues until a SIGTERM is observed.
func (generator *Generator) Start() error {

	generator.generatorMetrics.Start()
	generator.statusTracker.Start()

	for _, writer := range generator.writers {
		writer.Start()
		time.Sleep(generator.config.InstanceLaunchInterval)
	}

	for _, reader := range generator.readers {
		reader.Start()
		time.Sleep(generator.config.InstanceLaunchInterval)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	(*generator.cancel)()
	generator.waitGroup.Wait()
	return nil
}
