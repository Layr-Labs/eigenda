package traffic

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	retrivereth "github.com/Layr-Labs/eigenda/retriever/eth"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/Layr-Labs/eigenda/tools/traffic/workers"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
)

// Generator simulates read/write traffic to the DA service.
//
//	┌------------┐                                       ┌------------┐
//	|   writer   |-┐             ┌------------┐          |   reader   |-┐
//	└------------┘ |-┐  -------> |  verifier  | -------> └------------┘ |-┐
//	  └------------┘ |           └------------┘            └------------┘ |
//	    └------------┘                                       └------------┘
//
// The traffic generator is built from three principal components: one or more writers
// that write blobs, a verifier that polls the dispenser service until blobs are confirmed,
// and one or more readers that read blobs.
//
// When a writer finishes writing a blob, it sends information about that blob to the verifier.
// When the verifier observes that a blob has been confirmed, it sends information about the blob
// to the readers. The readers only attempt to read blobs that have been confirmed by the verifier.
type Generator struct {
	ctx              *context.Context
	cancel           *context.CancelFunc
	waitGroup        *sync.WaitGroup
	generatorMetrics metrics.Metrics
	logger           *logging.Logger
	disperserClient  clients.DisperserClient
	eigenDAClient    *clients.EigenDAClient
	config           *config.Config

	writers  []*workers.BlobWriter
	verifier *workers.BlobVerifier
	readers  []*workers.BlobReader
}

func NewTrafficGenerator(config *config.Config, signer core.BlobRequestSigner) (*Generator, error) {
	logger, err := common.NewLogger(config.LoggingConfig)
	if err != nil {
		return nil, err
	}

	clientConfig := clients.EigenDAClientConfig{
		RPC:                 config.DisperserHostname + ":" + config.DisperserPort,
		DisableTLS:          config.DisableTlS,
		SignerPrivateKeyHex: config.SignerPrivateKey,
	}
	err = clientConfig.CheckAndSetDefaults()
	if err != nil {
		return nil, err
	}

	logger2 := log.NewLogger(log.NewTerminalHandler(os.Stderr, true))
	client, err := clients.NewEigenDAClient(logger2, clientConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}

	generatorMetrics := metrics.NewMetrics(config.MetricsHTTPPort, logger)

	blobTable := table.NewBlobTable()

	disperserConfig := clients.Config{
		Hostname:          config.DisperserHostname,
		Port:              config.DisperserPort,
		Timeout:           config.DisperserTimeout,
		UseSecureGrpcFlag: config.DisperserUseSecureGrpcFlag,
	}
	disperserClient := clients.NewDisperserClient(&disperserConfig, signer)
	statusVerifier := workers.NewBlobVerifier(
		&ctx,
		&waitGroup,
		logger,
		workers.NewTicker(config.WorkerConfig.VerifierInterval),
		&config.WorkerConfig,
		&blobTable,
		disperserClient,
		generatorMetrics)

	writers := make([]*workers.BlobWriter, 0)
	for i := 0; i < int(config.WorkerConfig.NumWriteInstances); i++ {
		writer := workers.NewBlobWriter(
			&ctx,
			&waitGroup,
			logger,
			workers.NewTicker(config.WorkerConfig.WriteRequestInterval),
			&config.WorkerConfig,
			disperserClient,
			&statusVerifier,
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
			workers.NewTicker(config.WorkerConfig.ReadRequestInterval),
			&config.WorkerConfig,
			retriever,
			chainClient,
			&blobTable,
			generatorMetrics)
		readers = append(readers, &reader)
	}

	return &Generator{
		ctx:              &ctx,
		cancel:           &cancel,
		waitGroup:        &waitGroup,
		generatorMetrics: generatorMetrics,
		logger:           &logger,
		disperserClient:  clients.NewDisperserClient(&disperserConfig, signer),
		eigenDAClient:    client,
		config:           config,
		writers:          writers,
		verifier:         &statusVerifier,
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

	ethClientConfig := geth.EthClientConfig{
		RPCURLs:    []string{config.EthClientHostname + ":" + config.EthClientPort},
		NumRetries: int(config.EthClientRetries),
	}
	gethClient, err := geth.NewMultiHomingClient(ethClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		panic(fmt.Sprintf("Unable to instantiate geth client: %s", err))
	}

	tx, err := eth.NewTransactor(
		logger,
		gethClient,
		config.BlsOperatorStateRetriever,
		config.EigenDAServiceManager)
	if err != nil {
		panic(fmt.Sprintf("Unable to instantiate transactor: %s", err))
	}

	cs := eth.NewChainState(tx, gethClient)

	// This is the indexer when config.UseGraph is true
	chainStateConfig := thegraph.Config{
		Endpoint:     config.TheGraphUrl,
		PullInterval: config.TheGraphPullInterval,
		MaxRetries:   int(config.TheGraphRetries),
	}
	chainState := thegraph.MakeIndexedChainState(chainStateConfig, cs, logger)

	var assignmentCoordinator core.AssignmentCoordinator = &core.StdAssignmentCoordinator{}

	nodeClient := clients.NewNodeClient(config.NodeClientTimeout)

	encoderConfig := kzg.KzgConfig{
		G1Path:          config.EncoderG1Path,
		G2Path:          config.EncoderG2Path,
		CacheDir:        config.EncoderCacheDir,
		SRSOrder:        config.EncoderSRSOrder,
		SRSNumberToLoad: config.EncoderSRSNumberToLoad,
		NumWorker:       config.EncoderNumWorkers,
	}
	v, err := verifier.NewVerifier(&encoderConfig, true)
	if err != nil {
		panic(fmt.Sprintf("Unable to build verifier: %s", err))
	}

	retriever, err := clients.NewRetrievalClient(
		logger,
		chainState,
		assignmentCoordinator,
		nodeClient,
		v,
		int(config.RetrieverNumConnections))

	if err != nil {
		panic(fmt.Sprintf("Unable to build retriever: %s", err))
	}

	chainClient := retrivereth.NewChainClient(gethClient, logger)

	return retriever, chainClient
}

// Start instantiates goroutines that generate read/write traffic, continues until a SIGTERM is observed.
func (generator *Generator) Start() error {

	generator.generatorMetrics.Start()
	generator.verifier.Start()

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
