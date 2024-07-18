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
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	cancel()
	wg.Wait()
	return nil
}

// TODO maybe split reader/writer into separate files

// StartReadWorker periodically requests to download random blobs at a configured rate.
func (g *TrafficGenerator) StartReadWorker(ctx context.Context) error {
	ticker := time.NewTicker(g.Config.WriteRequestInterval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			// TODO determine which blob to download
			g.readRequest() // TODO add parameters
		}
	}
}

// readRequest reads a blob.
func (g *TrafficGenerator) readRequest() {
	// TODO
}
