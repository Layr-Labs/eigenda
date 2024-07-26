package workers

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"sync"
	"time"
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
	config *Config

	// disperser is the client used to send blobs to the disperser.
	disperser *clients.DisperserClient

	// Responsible for polling on the status of a recently written blob until it becomes confirmed.
	verifier *BlobVerifier

	// fixedRandomData contains random data for blobs if RandomizeBlobs is false, and nil otherwise.
	fixedRandomData *[]byte

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
	config *Config,
	disperser *clients.DisperserClient,
	verifier *BlobVerifier,
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
		verifier:           verifier,
		fixedRandomData:    &fixedRandomData,
		writeLatencyMetric: generatorMetrics.NewLatencyMetric("write"),
		writeSuccessMetric: generatorMetrics.NewCountMetric("write_success"),
		writeFailureMetric: generatorMetrics.NewCountMetric("write_failure"),
	}
}

// Start begins the blob writer goroutine.
func (writer *BlobWriter) Start() {
	writer.waitGroup.Add(1)
	go func() {
		writer.run()
		writer.waitGroup.Done()
	}()
}

// run sends blobs to a disperser at a configured rate.
// Continues and dues not return until the context is cancelled.
func (writer *BlobWriter) run() {
	ticker := time.NewTicker(writer.config.WriteRequestInterval)
	for {
		select {
		case <-(*writer.ctx).Done():
			return
		case <-ticker.C:
			data := writer.getRandomData()
			key, err := metrics.InvokeAndReportLatency(&writer.writeLatencyMetric, func() ([]byte, error) {
				return writer.sendRequest(*data)
			})
			if err != nil {
				writer.writeFailureMetric.Increment()
				writer.logger.Error("failed to send blob request", "err:", err)
				continue
			}

			writer.writeSuccessMetric.Increment()

			checksum := md5.Sum(*data)
			writer.verifier.AddUnconfirmedKey(&key, &checksum, uint(len(*data)))
		}
	}
}

// getRandomData returns a slice of random data to be used for a blob.
func (writer *BlobWriter) getRandomData() *[]byte {
	if *writer.fixedRandomData != nil {
		return writer.fixedRandomData
	}

	data := make([]byte, writer.config.DataSize)
	_, err := rand.Read(data)
	if err != nil {
		panic(fmt.Sprintf("unable to read random data: %s", err))
	}
	data = codec.ConvertByPaddingEmptyByte(data)

	return &data
}

// sendRequest sends a blob to a disperser.
func (writer *BlobWriter) sendRequest(data []byte) ([]byte /* key */, error) {

	ctxTimeout, cancel := context.WithTimeout(*writer.ctx, writer.config.WriteTimeout)
	defer cancel()

	var key []byte
	var err error
	if writer.config.SignerPrivateKey != "" {
		_, key, err = (*writer.disperser).DisperseBlobAuthenticated(
			ctxTimeout,
			data,
			writer.config.CustomQuorums)
	} else {
		_, key, err = (*writer.disperser).DisperseBlob(
			ctxTimeout,
			data,
			writer.config.CustomQuorums)
	}
	if err != nil {
		return nil, err
	}
	return key, nil
}
