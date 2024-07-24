package traffic

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"sync"
	"time"
)

// TODO document

// BlobWriter sends blobs to a disperser at a configured rate.
type BlobWriter struct {
	// The context for the generator. All work should cease when this context is cancelled.
	ctx *context.Context

	// Tracks the number of active goroutines within the generator.
	waitGroup *sync.WaitGroup

	// TODO this type should be refactored maybe
	generator *TrafficGenerator

	// Responsible for polling on the status of a recently written blob until it becomes confirmed.
	verifier *StatusVerifier

	// fixedRandomData contains random data for blobs if RandomizeBlobs is false, and nil otherwise.
	fixedRandomData *[]byte

	// writeLatencyMetric is used to record latency for write requests.
	writeLatencyMetric LatencyMetric

	// writeSuccessMetric is used to record the number of successful write requests.
	writeSuccessMetric CountMetric

	// writeFailureMetric is used to record the number of failed write requests.
	writeFailureMetric CountMetric
}

// NewBlobWriter creates a new BlobWriter instance.
func NewBlobWriter(
	ctx *context.Context,
	waitGroup *sync.WaitGroup,
	generator *TrafficGenerator,
	verifier *StatusVerifier,
	metrics *Metrics) BlobWriter {

	var fixedRandomData []byte
	if generator.Config.RandomizeBlobs {
		// New random data will be generated for each blob.
		fixedRandomData = nil
	} else {
		// Use this random data for each blob.
		fixedRandomData := make([]byte, generator.Config.DataSize)
		_, err := rand.Read(fixedRandomData)
		if err != nil {
			panic(err)
		}
		fixedRandomData = codec.ConvertByPaddingEmptyByte(fixedRandomData)
	}

	return BlobWriter{
		ctx:                ctx,
		waitGroup:          waitGroup,
		generator:          generator,
		verifier:           verifier,
		fixedRandomData:    &fixedRandomData,
		writeLatencyMetric: metrics.NewLatencyMetric("write"),
		writeSuccessMetric: metrics.NewCountMetric("write_success"),
		writeFailureMetric: metrics.NewCountMetric("write_failure"),
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
	ticker := time.NewTicker(writer.generator.Config.WriteRequestInterval)
	for {
		select {
		case <-(*writer.ctx).Done():
			return
		case <-ticker.C:
			data := writer.getRandomData()
			key, err := InvokeAndReportLatency(&writer.writeLatencyMetric, func() ([]byte, error) {
				return writer.sendRequest(*data)
			})
			if err != nil {
				writer.writeFailureMetric.Increment()
				writer.generator.Logger.Error("failed to send blob request", "err:", err)
				continue
			}

			writer.writeSuccessMetric.Increment()
			fmt.Println("Sent blob with length", len(*data)) // TODO remove

			writer.verifier.AddUnconfirmedKey(&key)
		}
	}
}

// getRandomData returns a slice of random data to be used for a blob.
func (writer *BlobWriter) getRandomData() *[]byte {
	// TODO
	//if *writer.fixedRandomData != nil {
	//	return writer.fixedRandomData
	//}

	data := make([]byte, writer.generator.Config.DataSize)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	// TODO: get explanation why this is necessary
	data = codec.ConvertByPaddingEmptyByte(data)

	// TODO remove
	//if !fft.IsPowerOfTwo(uint64(len(data))) {
	//	p := uint32(math.Log2(float64(len(data)))) + 1
	//	newSize := 1 << p
	//	bytesToAdd := newSize - len(data)
	//	for i := 0; i < bytesToAdd; i++ {
	//		data = append(data, 0)
	//	}
	//}

	return &data
}

// sendRequest sends a blob to a disperser.
func (writer *BlobWriter) sendRequest(data []byte) ([]byte /* key */, error) {
	// TODO remove
	//if !fft.IsPowerOfTwo(uint64(len(data))) {
	//	return nil, fmt.Errorf("data length must be a power of two, data size = %d", len(data))
	//}

	// TODO add timeout to other types of requests
	ctxTimeout, cancel := context.WithTimeout(*writer.ctx, writer.generator.Config.Timeout)
	defer cancel()

	if writer.generator.Config.SignerPrivateKey != "" {
		_, key, err :=
			writer.generator.DisperserClient.DisperseBlobAuthenticated(ctxTimeout, data, writer.generator.Config.CustomQuorums)
		if err != nil {
			return nil, err
		}

		// TODO
		//writer.generator.Logger.Info("successfully dispersed new blob", "authenticated", true, "key", hex.EncodeToString(key), "status", blobStatus.String())
		return key, nil
	} else {
		_, key, err := writer.generator.DisperserClient.DisperseBlob(ctxTimeout, data, writer.generator.Config.CustomQuorums)
		if err != nil {
			return nil, err
		}

		// TODO
		//writer.generator.Logger.Info("successfully dispersed new blob", "authenticated", false, "key", hex.EncodeToString(key), "status", blobStatus.String())
		return key, nil
	}
}
