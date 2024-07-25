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

// TODO use consistent snake case on metrics

// TrafficGenerator simulates read/write traffic to the DA service.
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
type TrafficGenerator struct {
	ctx             *context.Context
	cancel          *context.CancelFunc
	waitGroup       *sync.WaitGroup
	metrics         *Metrics
	logger          *logging.Logger
	disperserClient clients.DisperserClient
	eigenDAClient   *clients.EigenDAClient
	config          *Config

	writers  []*BlobWriter
	verifier *BlobVerifier
	readers  []*BlobReader
}

func NewTrafficGenerator(config *Config, signer core.BlobRequestSigner) (*TrafficGenerator, error) {
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Signer private key: '%s', length: %d\n", config.SignerPrivateKey, len(config.SignerPrivateKey))

	clientConfig := clients.EigenDAClientConfig{
		RPC:                 config.Config.Hostname + ":" + config.Config.Port,
		DisableTLS:          config.DisableTlS,
		SignerPrivateKeyHex: config.SignerPrivateKey,
	}
	err = clientConfig.CheckAndSetDefaults()
	if err != nil {
		return nil, err
	}

	logger2 := log.NewLogger(log.NewTerminalHandler(os.Stderr, true)) // TODO
	client, err := clients.NewEigenDAClient(logger2, clientConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}

	metrics := NewMetrics(config.MetricsHTTPPort, logger)

	table := NewBlobTable()

	disperserClient := clients.NewDisperserClient(&config.Config, signer)
	statusVerifier := NewStatusVerifier(
		&ctx,
		&waitGroup,
		logger,
		config,
		&table,
		&disperserClient,
		metrics)

	writers := make([]*BlobWriter, 0)
	for i := 0; i < int(config.NumWriteInstances); i++ {
		writer := NewBlobWriter(
			&ctx,
			&waitGroup,
			logger,
			config,
			&disperserClient,
			&statusVerifier,
			metrics)
		writers = append(writers, &writer)
	}

	retriever, chainClient := buildRetriever(config)

	readers := make([]*BlobReader, 0)
	for i := 0; i < int(config.NumReadInstances); i++ {
		reader := NewBlobReader(
			&ctx,
			&waitGroup,
			logger,
			config,
			retriever,
			chainClient,
			&table,
			metrics)
		readers = append(readers, &reader)
	}

	return &TrafficGenerator{
		ctx:             &ctx,
		cancel:          &cancel,
		waitGroup:       &waitGroup,
		metrics:         metrics,
		logger:          &logger,
		disperserClient: clients.NewDisperserClient(&config.Config, signer),
		eigenDAClient:   client,
		config:          config,
		writers:         writers,
		verifier:        &statusVerifier,
		readers:         readers,
	}, nil
}

// buildRetriever creates a retriever client for the traffic generator.
func buildRetriever(config *Config) (clients.RetrievalClient, retrivereth.ChainClient) {
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

	tx, err := eth.NewTransactor(
		logger,
		gethClient,
		config.BlsOperatorStateRetriever,
		config.EigenDAServiceManager)

	cs := eth.NewChainState(tx, gethClient)

	// This is the indexer when config.UseGraph is true
	chainStateConfig := thegraph.Config{
		Endpoint:     config.TheGraphUrl,
		PullInterval: config.TheGraphPullInterval,
		MaxRetries:   int(config.TheGraphRetries),
	}
	chainState := thegraph.MakeIndexedChainState(chainStateConfig, cs, logger)

	var assignmentCoordinator core.AssignmentCoordinator = &core.StdAssignmentCoordinator{}

	nodeClient := clients.NewNodeClient(10 * time.Second)

	encoderConfig := kzg.KzgConfig{
		G1Path:          "../../inabox/resources/kzg/g1.point",
		G2Path:          "../../inabox/resources/kzg/g2.point",
		CacheDir:        "../../inabox/resources/kzg/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       12,
	}
	v, err := verifier.NewVerifier(&encoderConfig, true)
	if err != nil {
		panic(fmt.Sprintf("Unable to build verifier: %s", err))
	}

	numConnections := 20

	retriever, err := clients.NewRetrievalClient(
		logger,
		chainState,
		assignmentCoordinator,
		nodeClient,
		v,
		numConnections)

	if err != nil {
		panic(fmt.Sprintf("Unable to build retriever: %s", err))
	}

	chainClient := retrivereth.NewChainClient(gethClient, logger)

	return retriever, chainClient
}

// Start instantiates goroutines that generate read/write traffic, continues until a SIGTERM is observed.
func (generator *TrafficGenerator) Start() error {

	// TODO add configuration

	generator.metrics.Start(*generator.ctx) // TODO put context into metrics constructor
	generator.verifier.Start(time.Second)

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
