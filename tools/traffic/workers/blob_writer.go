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
	// Config contains the configuration for the generator.
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
}

// NewBlobWriter creates a new BlobWriter instance.
func NewBlobWriter(
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
	writer.logger.Info("Starting blob writer")
	writer.waitGroup.Add(1)
	ticker := time.NewTicker(writer.config.WriteRequestInterval)

	go func() {
		defer writer.waitGroup.Done()
		defer ticker.Stop()

		for {
			select {
			case <-(*writer.ctx).Done():
				writer.logger.Info("context cancelled, stopping blob writer")
				return
			case <-ticker.C:
				if err := writer.writeNextBlob(); err != nil {
					writer.logger.Error("failed to write blob", "err", err)
				}
			}
		}
	}()
}

// writeNextBlob attempts to send a random blob to the disperser.
func (writer *BlobWriter) writeNextBlob() error {
	data, err := writer.getRandomData()
	if err != nil {
		writer.logger.Error("failed to get random data", "err", err)
		return err
	}
	start := time.Now()
	_, err = writer.sendRequest(data)
	if err != nil {
		writer.writeFailureMetric.Increment()
		writer.logger.Error("failed to send blob request", "err", err)
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
	ctxTimeout, cancel := context.WithTimeout(*writer.ctx, writer.config.WriteTimeout)
	defer cancel()

	writer.logger.Info("sending blob request", "size", len(data))
	status, key, err := writer.disperser.DisperseBlob(
		ctxTimeout,
		data,
		0,
		writer.config.CustomQuorums,
		0,
	)
	if err != nil {
		writer.logger.Error("failed to send blob request", "err", err)
		return
	}

	writer.logger.Info("blob request sent", "key", key.Hex(), "status", status.String())

	return
}
