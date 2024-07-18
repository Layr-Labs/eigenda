package traffic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
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
	var wg sync.WaitGroup
	for i := 0; i < int(g.Config.NumWriteInstances); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = g.StartWriteWorker(ctx)
		}()
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

// StartWriteWorker periodically sends (possibly) random blobs to a disperser at a configured rate.
func (g *TrafficGenerator) StartWriteWorker(ctx context.Context) error {
	data := make([]byte, g.Config.DataSize)
	_, err := rand.Read(data)
	if err != nil {
		return err
	}

	// TODO configuration for this stuff
	var table BlobTable = NewBlobTable()
	var verifier StatusVerifier = NewStatusVerifier(&table, &g.DisperserClient, -1)
	verifier.Start(ctx, time.Second)

	paddedData := codec.ConvertByPaddingEmptyByte(data)

	ticker := time.NewTicker(g.Config.WriteRequestInterval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			var key []byte
			if g.Config.RandomizeBlobs {
				_, err := rand.Read(data)
				if err != nil {
					return err
				}
				paddedData = codec.ConvertByPaddingEmptyByte(data)

				key, err = g.sendRequest(ctx, paddedData[:g.Config.DataSize])

				if err != nil {
					g.Logger.Error("failed to send blob request", "err:", err)
				}
				paddedData = nil
			} else {
				key, err = g.sendRequest(ctx, paddedData[:g.Config.DataSize])
				if err != nil {
					g.Logger.Error("failed to send blob request", "err:", err)
				}
			}

			fmt.Println("passing key to verifier") // TODO remove
			verifier.AddUnconfirmedKey(&key)
			fmt.Println("done passing key") // TODO remove
		}
	}
}

// sendRequest sends a blob to a disperser.
func (g *TrafficGenerator) sendRequest(ctx context.Context, data []byte) ([]byte /* key */, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, g.Config.Timeout)
	defer cancel()

	if g.Config.SignerPrivateKey != "" {
		blobStatus, key, err := g.DisperserClient.DisperseBlobAuthenticated(ctxTimeout, data, g.Config.CustomQuorums)
		if err != nil {
			return nil, err
		}

		g.Logger.Info("successfully dispersed new blob", "authenticated", true, "key", hex.EncodeToString(key), "status", blobStatus.String())
		return key, nil
	} else {
		blobStatus, key, err := g.DisperserClient.DisperseBlob(ctxTimeout, data, g.Config.CustomQuorums)
		if err != nil {
			return nil, err
		}

		g.Logger.Info("successfully dispersed new blob", "authenticated", false, "key", hex.EncodeToString(key), "status", blobStatus.String())
		return key, nil
	}
}
