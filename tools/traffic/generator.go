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
		return nil, fmt.Errorf("new logger: %w", err)
	}

	dispserserClient, err := clients.NewDisperserClient(&config.Config, signer)
	if err != nil {
		return nil, fmt.Errorf("new disperser-client: %w", err)
	}
	return &TrafficGenerator{
		Logger:          logger,
		DisperserClient: dispserserClient,
		Config:          config,
	}, nil
}

func (g *TrafficGenerator) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := 0; i < int(g.Config.NumInstances); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = g.StartTraffic(ctx)
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

func (g *TrafficGenerator) StartTraffic(ctx context.Context) error {
	data := make([]byte, g.Config.DataSize)
	_, err := rand.Read(data)
	if err != nil {
		return err
	}

	paddedData := codec.ConvertByPaddingEmptyByte(data)

	ticker := time.NewTicker(g.Config.RequestInterval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if g.Config.RandomizeBlobs {
				_, err := rand.Read(data)
				if err != nil {
					return err
				}
				paddedData = codec.ConvertByPaddingEmptyByte(data)

				err = g.sendRequest(ctx, paddedData[:g.Config.DataSize])
				if err != nil {
					g.Logger.Error("failed to send blob request", "err:", err)
				}
				paddedData = nil
			} else {
				err = g.sendRequest(ctx, paddedData[:g.Config.DataSize])
				if err != nil {
					g.Logger.Error("failed to send blob request", "err:", err)
				}
			}

		}
	}
}

func (g *TrafficGenerator) sendRequest(ctx context.Context, data []byte) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, g.Config.Timeout)
	defer cancel()

	if g.Config.SignerPrivateKey != "" {
		blobStatus, key, err := g.DisperserClient.DisperseBlobAuthenticated(ctxTimeout, data, g.Config.CustomQuorums)
		if err != nil {
			return err
		}

		g.Logger.Info("successfully dispersed new blob", "authenticated", true, "key", hex.EncodeToString(key), "status", blobStatus.String())
		return nil
	} else {
		blobStatus, key, err := g.DisperserClient.DisperseBlob(ctxTimeout, data, g.Config.CustomQuorums)
		if err != nil {
			return err
		}

		g.Logger.Info("successfully dispersed new blob", "authenticated", false, "key", hex.EncodeToString(key), "status", blobStatus.String())
		return nil
	}

}
