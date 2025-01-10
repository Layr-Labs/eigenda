package workers

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// BlobWriter sends blobs to a disperser at a configured rate.
type BlobWriter struct {
	// Name of the writer group this writer belongs to
	name string

	// Config contains the configuration for the blob writer.
	config *config.BlobWriterConfig

	// The context for the generator. All work should cease when this context is cancelled.
	ctx *context.Context

	// Tracks the number of active goroutines within the generator.
	waitGroup *sync.WaitGroup

	// All logs should be written using this logger.
	logger logging.Logger

	// disperser is the client used to send blobs to the disperser.
	disperser clients.DisperserClient

	// fixedRandomData contains random data for blobs if RandomizeBlobs is false, and nil otherwise.
	fixedRandomData []byte

	// writeLatencyMetric is used to record latency for write requests.
	writeLatencyMetric metrics.LatencyMetric

	// writeSuccessMetric is used to record the number of successful write requests.
	writeSuccessMetric metrics.CountMetric

	// writeFailureMetric is used to record the number of failed write requests.
	writeFailureMetric metrics.CountMetric

	// Mutex to protect configuration updates
	configMutex sync.RWMutex

	// Ticker for controlling write intervals
	ticker *time.Ticker

	// cancel is used to cancel the context
	cancel *context.CancelFunc
}

// NewBlobWriter creates a new BlobWriter instance.
func NewBlobWriter(
	name string,
	ctx *context.Context,
	config *config.BlobWriterConfig,
	waitGroup *sync.WaitGroup,
	logger logging.Logger,
	disperser clients.DisperserClient,
	generatorMetrics metrics.Metrics) BlobWriter {

	var fixedRandomData []byte
	if config.RandomizeBlobs {
		// New random data will be generated for each blob.
		fixedRandomData = nil
	} else {
		// Use this random data for each blob.
		fixedRandomData = make([]byte, config.DataSize)
		_, err := rand.Read(fixedRandomData)
		if err != nil {
			panic(fmt.Sprintf("unable to read random data: %s", err))
		}
		fixedRandomData = codec.ConvertByPaddingEmptyByte(fixedRandomData)
	}

	return BlobWriter{
		name:               name,
		ctx:                ctx,
		waitGroup:          waitGroup,
		logger:             logger,
		config:             config,
		disperser:          disperser,
		fixedRandomData:    fixedRandomData,
		writeLatencyMetric: generatorMetrics.NewLatencyMetric("write"),
		writeSuccessMetric: generatorMetrics.NewCountMetric("write_success"),
		writeFailureMetric: generatorMetrics.NewCountMetric("write_failure"),
	}
}

// Start begins the blob writer goroutine.
func (writer *BlobWriter) Start() {
	writer.logger.Info("Starting blob writer", "name", writer.name)
	writer.waitGroup.Add(1)
	writer.configMutex.Lock()
	writer.ticker = time.NewTicker(writer.config.WriteRequestInterval)
	writer.configMutex.Unlock()

	go func() {
		defer writer.waitGroup.Done()
		defer writer.ticker.Stop()

		for {
			select {
			case <-(*writer.ctx).Done():
				writer.logger.Info("context cancelled, stopping blob writer", "name", writer.name)
				return
			case <-writer.ticker.C:
				if err := writer.writeNextBlob(); err != nil {
					writer.logger.Error("failed to write blob", "name", writer.name, "err", err)
				}
			}
		}
	}()
}

// UpdateConfig updates the writer's configuration
func (writer *BlobWriter) UpdateConfig(config *config.BlobWriterConfig) {
	writer.configMutex.Lock()
	defer writer.configMutex.Unlock()

	// Update the ticker if the interval changed
	if writer.config.WriteRequestInterval != config.WriteRequestInterval {
		writer.ticker.Reset(config.WriteRequestInterval)
	}

	// Update the fixed random data if needed
	if writer.config.RandomizeBlobs != config.RandomizeBlobs || writer.config.DataSize != config.DataSize {
		if config.RandomizeBlobs {
			writer.fixedRandomData = nil
		} else {
			writer.fixedRandomData = make([]byte, config.DataSize)
			_, err := rand.Read(writer.fixedRandomData)
			if err != nil {
				writer.logger.Error("failed to generate new fixed random data", "name", writer.name, "err", err)
				return
			}
			writer.fixedRandomData = codec.ConvertByPaddingEmptyByte(writer.fixedRandomData)
		}
	}

	writer.config = config
	writer.logger.Info("Updated blob writer configuration",
		"name", writer.name,
		"writeInterval", config.WriteRequestInterval,
		"dataSize", config.DataSize,
		"randomizeBlobs", config.RandomizeBlobs)
}

// writeNextBlob attempts to send a random blob to the disperser.
func (writer *BlobWriter) writeNextBlob() error {
	data, err := writer.getRandomData()
	if err != nil {
		writer.logger.Error("failed to get random data", "name", writer.name, "err", err)
		return err
	}
	start := time.Now()
	_, err = writer.sendRequest(data)
	if err != nil {
		writer.writeFailureMetric.Increment()
		writer.logger.Error("failed to send blob request", "name", writer.name, "err", err)
		return err
	}

	end := time.Now()
	duration := end.Sub(start)
	writer.writeLatencyMetric.ReportLatency(duration)
	writer.writeSuccessMetric.Increment()

	return nil
}

// getRandomData returns a slice of random data to be used for a blob.
func (writer *BlobWriter) getRandomData() ([]byte, error) {
	writer.configMutex.RLock()
	defer writer.configMutex.RUnlock()

	if writer.fixedRandomData != nil {
		return writer.fixedRandomData, nil
	}

	data := make([]byte, writer.config.DataSize)
	_, err := rand.Read(data)
	if err != nil {
		return nil, fmt.Errorf("unable to read random data: %w", err)
	}
	data = codec.ConvertByPaddingEmptyByte(data)

	return data, nil
}

// sendRequest sends a blob to a disperser.
func (writer *BlobWriter) sendRequest(data []byte) (key v2.BlobKey, err error) {
	writer.configMutex.RLock()
	writeTimeout := writer.config.WriteTimeout
	customQuorums := writer.config.CustomQuorums
	writer.configMutex.RUnlock()

	ctxTimeout, cancel := context.WithTimeout(*writer.ctx, writeTimeout)
	defer cancel()

	writer.logger.Info("sending blob request", "name", writer.name, "size", len(data))
	status, key, err := writer.disperser.DisperseBlob(
		ctxTimeout,
		data,
		0,
		customQuorums,
		0,
	)
	if err != nil {
		writer.logger.Error("failed to send blob request", "name", writer.name, "err", err)
		return
	}

	writer.logger.Info("blob request sent", "name", writer.name, "key", key.Hex(), "status", status.String())

	return
}
