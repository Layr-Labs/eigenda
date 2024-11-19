package workers

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// BlobWriter sends blobs to a disperser at a configured rate.
type BlobWriter struct {
	// The context for the generator. All work should cease when this context is cancelled.
	ctx *context.Context

	// Tracks the number of active goroutines within the generator.
	waitGroup *sync.WaitGroup

	// All logs should be written using this logger.
	logger logging.Logger

	// Config contains the configuration for the generator.
	config *config.WorkerConfig

	// disperser is the client used to send blobs to the disperser.
	disperser clients.DisperserClient

	// Unconfirmed keys are sent here.
	unconfirmedKeyChannel chan *UnconfirmedKey

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
	waitGroup *sync.WaitGroup,
	logger logging.Logger,
	config *config.WorkerConfig,
	disperser clients.DisperserClient,
	unconfirmedKeyChannel chan *UnconfirmedKey,
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
		ctx:                   ctx,
		waitGroup:             waitGroup,
		logger:                logger,
		config:                config,
		disperser:             disperser,
		unconfirmedKeyChannel: unconfirmedKeyChannel,
		fixedRandomData:       fixedRandomData,
		writeLatencyMetric:    generatorMetrics.NewLatencyMetric("write"),
		writeSuccessMetric:    generatorMetrics.NewCountMetric("write_success"),
		writeFailureMetric:    generatorMetrics.NewCountMetric("write_failure"),
	}
}

// Start begins the blob writer goroutine.
func (writer *BlobWriter) Start() {
	writer.waitGroup.Add(1)
	ticker := time.NewTicker(writer.config.WriteRequestInterval)

	go func() {
		defer writer.waitGroup.Done()

		for {
			select {
			case <-(*writer.ctx).Done():
				return
			case <-ticker.C:
				writer.writeNextBlob()
			}
		}
	}()
}

// writeNextBlob attempts to send a random blob to the disperser.
func (writer *BlobWriter) writeNextBlob() {
	data, err := writer.getRandomData()
	if err != nil {
		writer.logger.Error("failed to get random data", "err", err)
		return
	}
	start := time.Now()
	key, err := writer.sendRequest(data)
	if err != nil {
		writer.writeFailureMetric.Increment()
		writer.logger.Error("failed to send blob request", "err", err)
		return
	} else {
		end := time.Now()
		duration := end.Sub(start)
		writer.writeLatencyMetric.ReportLatency(duration)
	}

	writer.writeSuccessMetric.Increment()

	checksum := md5.Sum(data)

	writer.unconfirmedKeyChannel <- &UnconfirmedKey{
		Key:            key,
		Checksum:       checksum,
		Size:           uint(len(data)),
		SubmissionTime: time.Now(),
	}
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
func (writer *BlobWriter) sendRequest(data []byte) (key []byte, err error) {
	ctxTimeout, cancel := context.WithTimeout(*writer.ctx, writer.config.WriteTimeout)
	defer cancel()

	if writer.config.SignerPrivateKey != "" {
		_, key, err = writer.disperser.DisperseBlobAuthenticated(
			ctxTimeout,
			data,
			writer.config.CustomQuorums)
	} else {
		_, key, err = writer.disperser.DisperseBlob(
			ctxTimeout,
			data,
			writer.config.CustomQuorums)
	}
	return
}
