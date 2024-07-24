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
	"github.com/Layr-Labs/eigensdk-go/logging"
)

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
	Logger          logging.Logger
	DisperserClient clients.DisperserClient
	EigenDAClient   *clients.EigenDAClient
	Config          *Config
}

func NewTrafficGenerator(config *Config, signer core.BlobRequestSigner) (*TrafficGenerator, error) {
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Signer private key: '%s', length: %d\n", config.SignerPrivateKey, len(config.SignerPrivateKey))

	clientConfig := clients.EigenDAClientConfig{
		RPC:                 "localhost:32003", // TODO make this configurable
		DisableTLS:          true,              // TODO config
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

	return &TrafficGenerator{
		Logger:          logger,
		DisperserClient: clients.NewDisperserClient(&config.Config, signer),
		EigenDAClient:   client,
		Config:          config,
	}, nil
}

// buildRetriever creates a retriever client for the traffic generator.
func (g *TrafficGenerator) buildRetriever() (clients.RetrievalClient, retrivereth.ChainClient) {

	//loggerConfig := common.LoggerConfig{
	//	Format: "text",
	// }

	loggerConfig := common.DefaultLoggerConfig()

	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		panic(err) // TODO
	}

	ethClientConfig := geth.EthClientConfig{
		RPCURLs:    []string{"http://localhost:8545"},
		NumRetries: 2,
	}
	gethClient, err := geth.NewMultiHomingClient(ethClientConfig, gethcommon.Address{}, logger)

	tx, err := eth.NewTransactor(
		logger,
		gethClient,
		"0x5f3f1dBD7B74C6B46e8c44f98792A1dAf8d69154",
		"0x851356ae760d987E095750cCeb3bC6014560891C")

	cs := eth.NewChainState(tx, gethClient)

	// -------------

	// This is the indexer when config.UseGraph is true
	chainStateConfig := thegraph.Config{
		Endpoint:     "http://localhost:8000/subgraphs/name/Layr-Labs/eigenda-operator-state",
		PullInterval: 100 * time.Millisecond,
		MaxRetries:   5,
	}
	chainState := thegraph.MakeIndexedChainState(chainStateConfig, cs, logger)

	// This is the indexer when config.UseGraph is false.
	//rpcClient, err := rpc.Dial("http://localhost:8545")
	//indexerConfig := indexer.Config{
	//	PullInterval: time.Second,
	//}
	//indexer, err := coreindexer.CreateNewIndexer(
	//	&indexerConfig,
	//	gethClient,
	//	rpcClient,
	//	"0x851356ae760d987E095750cCeb3bC6014560891C", // eigenDaServeManagerAddr
	//	logger,
	//)
	//if err != nil {
	//	panic(err) // TODO
	//}
	//chainState, err := coreindexer.NewIndexedChainState(cs, indexer)
	//if err != nil {
	//	panic(err) // TODO
	//}

	// -------------

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
		panic(err) // TODO
	}

	numConnections := 20

	//var retriever *clients.RetrievalClient
	retriever, err := clients.NewRetrievalClient(
		logger,
		chainState,
		assignmentCoordinator,
		nodeClient,
		v,
		numConnections)

	if err != nil {
		panic(err) // TODO
	}

	chainClient := retrivereth.NewChainClient(gethClient, logger)

	return retriever, chainClient
}

// Run instantiates goroutines that generate read/write traffic, continues until a SIGTERM is observed.
func (g *TrafficGenerator) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	metrics := NewMetrics("9101", g.Logger) // TODO config
	metrics.Start(ctx)

	// TODO add configuration
	table := NewBlobTable()
	verifier := NewStatusVerifier(&table, &g.DisperserClient, -1, metrics)
	verifier.Start(ctx, time.Second)

	var wg sync.WaitGroup

	for i := 0; i < int(g.Config.NumWriteInstances); i++ {
		writer := NewBlobWriter(&ctx, &wg, g, &verifier, metrics)
		writer.Start()
		time.Sleep(g.Config.InstanceLaunchInterval)
	}

	retriever, chainClient := g.buildRetriever()

	// TODO start multiple readers
	reader := NewBlobReader(&ctx, &wg, retriever, chainClient, &table, metrics)
	reader.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	cancel()
	wg.Wait()
	return nil
}
