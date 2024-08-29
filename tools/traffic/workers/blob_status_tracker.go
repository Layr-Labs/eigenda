package workers

import (
	"context"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"math/rand"
	"sync"
	"time"
)

// BlobStatusTracker periodically polls the disperser service to verify the status of blobs that were recently written.
// When blobs become confirmed, the status verifier updates the blob blobsToRead accordingly.
// This is a thread safe data structure.
type BlobStatusTracker struct {

	// The context for the generator. All work should cease when this context is cancelled.
	ctx *context.Context

	// Tracks the number of active goroutines within the generator.
	waitGroup *sync.WaitGroup

	// All logs should be written using this logger.
	logger logging.Logger

	// config contains the configuration for the generator.
	config *config.WorkerConfig

	// Contains confirmed blobs. Blobs are added here when they are confirmed by the disperser service.
	confirmedBlobs *table.BlobStore

	// The disperser client used to monitor the disperser service.
	disperser clients.DisperserClient

	// The keys of blobs that have not yet been confirmed by the disperser service.
	unconfirmedBlobs []*UnconfirmedKey

	// Newly added keys that require verification.
	keyChannel chan *UnconfirmedKey

	blobsInFlightMetric               metrics.GaugeMetric
	getStatusLatencyMetric            metrics.LatencyMetric
	confirmationLatencyMetric         metrics.LatencyMetric
	getStatusErrorCountMetric         metrics.CountMetric
	unknownCountMetric                metrics.CountMetric
	processingCountMetric             metrics.CountMetric
	dispersingCountMetric             metrics.CountMetric
	failedCountMetric                 metrics.CountMetric
	insufficientSignaturesCountMetric metrics.CountMetric
	confirmedCountMetric              metrics.CountMetric
	finalizedCountMetric              metrics.CountMetric
}

// NewBlobVerifier creates a new BlobStatusTracker instance.
func NewBlobVerifier(
	ctx *context.Context,
	waitGroup *sync.WaitGroup,
	logger logging.Logger,
	config *config.WorkerConfig,
	keyChannel chan *UnconfirmedKey,
	table *table.BlobStore,
	disperser clients.DisperserClient,
	generatorMetrics metrics.Metrics) BlobStatusTracker {

	return BlobStatusTracker{
		ctx:                               ctx,
		waitGroup:                         waitGroup,
		logger:                            logger,
		config:                            config,
		keyChannel:                        keyChannel,
		confirmedBlobs:                    table,
		disperser:                         disperser,
		unconfirmedBlobs:                  make([]*UnconfirmedKey, 0),
		blobsInFlightMetric:               generatorMetrics.NewGaugeMetric("blobs_in_flight"),
		getStatusLatencyMetric:            generatorMetrics.NewLatencyMetric("get_status"),
		confirmationLatencyMetric:         generatorMetrics.NewLatencyMetric("confirmation"),
		getStatusErrorCountMetric:         generatorMetrics.NewCountMetric("get_status_ERROR"),
		unknownCountMetric:                generatorMetrics.NewCountMetric("get_status_UNKNOWN"),
		processingCountMetric:             generatorMetrics.NewCountMetric("get_status_PROCESSING"),
		dispersingCountMetric:             generatorMetrics.NewCountMetric("get_status_DISPERSING"),
		failedCountMetric:                 generatorMetrics.NewCountMetric("get_status_FAILED"),
		insufficientSignaturesCountMetric: generatorMetrics.NewCountMetric("get_status_INSUFFICIENT_SIGNATURES"),
		confirmedCountMetric:              generatorMetrics.NewCountMetric("get_status_CONFIRMED"),
		finalizedCountMetric:              generatorMetrics.NewCountMetric("get_status_FINALIZED"),
	}
}

// Start begins the status goroutine, which periodically polls
// the disperser service to verify the status of blobs.
func (verifier *BlobStatusTracker) Start() {
	verifier.waitGroup.Add(1)
	go verifier.monitor()
}

// monitor periodically polls the disperser service to verify the status of blobs.
func (verifier *BlobStatusTracker) monitor() {
	ticker := time.NewTicker(verifier.config.VerifierInterval)
	for {
		select {
		case <-(*verifier.ctx).Done():
			verifier.waitGroup.Done()
			return
		case key := <-verifier.keyChannel:
			verifier.unconfirmedBlobs = append(verifier.unconfirmedBlobs, key)
		case <-ticker.C:
			verifier.poll()
		}
	}
}

// poll checks all unconfirmed keys to see if they have been confirmed by the disperser service.
// If a Key is confirmed, it is added to the blob confirmedBlobs and removed from the list of unconfirmed keys.
func (verifier *BlobStatusTracker) poll() {

	// FUTURE WORK If the number of unconfirmed blobs is high and the time to confirm is high, this is not efficient.
	// Revisit this method if there are performance problems.

	nonFinalBlobs := make([]*UnconfirmedKey, 0)
	for _, key := range verifier.unconfirmedBlobs {

		blobStatus, err := verifier.getBlobStatus(key)
		if err != nil {
			verifier.logger.Error("failed to get blob status: ", "err:", err)
			// There was an error getting status. Try again later.
			nonFinalBlobs = append(nonFinalBlobs, key)
			continue
		}

		verifier.updateStatusMetrics(blobStatus.Status)
		if isBlobStatusTerminal(blobStatus.Status) {
			if isBlobStatusConfirmed(blobStatus.Status) {
				verifier.forwardToReader(key, blobStatus)
			}
		} else {
			// try again later
			nonFinalBlobs = append(nonFinalBlobs, key)
		}
	}
	verifier.unconfirmedBlobs = nonFinalBlobs
	verifier.blobsInFlightMetric.Set(float64(len(verifier.unconfirmedBlobs)))
}

// isBlobStatusTerminal returns true if the status is a terminal status.
func isBlobStatusTerminal(status disperser.BlobStatus) bool {
	switch status {
	case disperser.BlobStatus_FAILED:
		return true
	case disperser.BlobStatus_INSUFFICIENT_SIGNATURES:
		return true
	case disperser.BlobStatus_CONFIRMED:
		// Technically this isn't terminal, as confirmed blobs eventually should become finalized.
		// But it is terminal from the status tracker's perspective, since we stop tracking the blob
		// once it becomes either confirmed or finalized.
		return true
	case disperser.BlobStatus_FINALIZED:
		return true
	default:
		return false
	}
}

// isBlobStatusConfirmed returns true if the status is a confirmed status.
func isBlobStatusConfirmed(status disperser.BlobStatus) bool {
	switch status {
	case disperser.BlobStatus_CONFIRMED:
		return true
	case disperser.BlobStatus_FINALIZED:
		return true
	default:
		return false
	}
}

// updateStatusMetrics updates the metrics for the reported status of a blob.
func (verifier *BlobStatusTracker) updateStatusMetrics(status disperser.BlobStatus) {
	switch status {
	case disperser.BlobStatus_UNKNOWN:
		verifier.unknownCountMetric.Increment()
	case disperser.BlobStatus_PROCESSING:
		verifier.processingCountMetric.Increment()
	case disperser.BlobStatus_DISPERSING:
		verifier.dispersingCountMetric.Increment()
	case disperser.BlobStatus_FAILED:
		verifier.failedCountMetric.Increment()
	case disperser.BlobStatus_INSUFFICIENT_SIGNATURES:
		verifier.insufficientSignaturesCountMetric.Increment()
	case disperser.BlobStatus_CONFIRMED:
		verifier.confirmedCountMetric.Increment()
	case disperser.BlobStatus_FINALIZED:
		verifier.finalizedCountMetric.Increment()
	default:
		verifier.logger.Error("unknown blob status", "status:", status)
	}
}

// getBlobStatus gets the status of a blob from the disperser service. Returns nil if there was an error
// getting the status.
func (verifier *BlobStatusTracker) getBlobStatus(key *UnconfirmedKey) (*disperser.BlobStatusReply, error) {
	ctxTimeout, cancel := context.WithTimeout(*verifier.ctx, verifier.config.GetBlobStatusTimeout)
	defer cancel()

	status, err := metrics.InvokeAndReportLatency[*disperser.BlobStatusReply](verifier.getStatusLatencyMetric,
		func() (*disperser.BlobStatusReply, error) {
			return verifier.disperser.GetBlobStatus(ctxTimeout, key.Key)
		})

	if err != nil {
		verifier.getStatusErrorCountMetric.Increment()
		return nil, err
	}

	return status, nil
}

// forwardToReader forwards a blob to the reader. Only called once the blob is ready to be read.
func (verifier *BlobStatusTracker) forwardToReader(key *UnconfirmedKey, status *disperser.BlobStatusReply) {
	batchHeaderHash := [32]byte(status.GetInfo().BlobVerificationProof.BatchMetadata.BatchHeaderHash)
	blobIndex := status.GetInfo().BlobVerificationProof.GetBlobIndex()

	confirmationTime := time.Now()
	confirmationLatency := confirmationTime.Sub(key.SubmissionTime)
	verifier.confirmationLatencyMetric.ReportLatency(confirmationLatency)

	requiredDownloads := verifier.config.RequiredDownloads
	var downloadCount int32
	if requiredDownloads < 0 {
		// Allow unlimited downloads.
		downloadCount = -1
	} else if requiredDownloads == 0 {
		// Do not download blob.
		return
	} else if requiredDownloads < 1 {
		// Download blob with probability equal to requiredDownloads.
		if rand.Float64() < requiredDownloads {
			// Download the blob once.
			downloadCount = 1
		} else {
			// Do not download blob.
			return
		}
	} else {
		// Download blob requiredDownloads times.
		downloadCount = int32(requiredDownloads)
	}

	blobMetadata, err := table.NewBlobMetadata(key.Key, key.Checksum, key.Size, uint(blobIndex), batchHeaderHash, int(downloadCount))
	if err != nil {
		verifier.logger.Error("failed to create blob metadata", "err:", err)
		return
	}
	verifier.confirmedBlobs.Add(blobMetadata)
}
