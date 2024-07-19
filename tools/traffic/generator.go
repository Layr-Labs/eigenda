package traffic

import (
	"context"
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
//		┌------------┐                                       ┌------------┐
//		|  writer    |-┐             ┌------------┐          |  reader    |-┐
//	 └------------┘ |-┐  -------> |  verifier  | -------> └------------┘ |-┐
//	   └------------┘ |           └------------┘            └------------┘ |
//	     └------------┘                                       └------------┘
//
// The traffic generator is built from three principal components: one or more writers
// that write blobs, a verifier that polls the dispenser service until blobs are confirmed,
// and one or more readers that read blobs.
//
// When a writer finishes writing a blob, it
// sends information about that blob to the verifier. When the verifier observes that a blob
// has been confirmed, it sends information about the blob to the readers. The readers
// only attempt to read blobs that have been confirmed by the verifier.
type TrafficGenerator struct {
	Logger          logging.Logger
	DisperserClient clients.DisperserClient
	Config          *Config
}

func NewTrafficGenerator(config *Config, signer core.BlobRequestSigner) (*TrafficGenerator, error) {
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return nil, err
	}

	return &TrafficGenerator{
		Logger:          logger,
		DisperserClient: clients.NewDisperserClient(&config.Config, signer),
		Config:          config,
	}, nil
}

// Run instantiates goroutines that generate read/write traffic, continues until a SIGTERM is observed.
func (g *TrafficGenerator) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	// TODO add configuration
	table := NewBlobTable()
	verifier := NewStatusVerifier(&table, &g.DisperserClient, -1)
	verifier.Start(ctx, time.Second)

	var wg sync.WaitGroup

	for i := 0; i < int(g.Config.NumWriteInstances); i++ {
		writer := NewBlobWriter(&ctx, &wg, g, &verifier)
		writer.Start()
		time.Sleep(g.Config.InstanceLaunchInterval)
	}

	// TODO start multiple readers
	reader := NewBlobReader(&ctx, &wg, g, &table)
	reader.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	cancel()
	wg.Wait()
	return nil
}
