package blobqueue

import (
	"context"
	"time"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
)

// Responsible for polling for blobs that need to be put into batches for dispersal.
type BlobQueue interface {

	// Close the queue and free any resources.
	Close()

	// Get the channel that will receive blobs to disperse.
	GetBlobToDisperse() <-chan *BlobToDisperse

	// It is assumed that the blob source may return duplicates "for a while". After the blob source will no longer
	// return duplicates for a particular blob (i.e. when its status in DynamoDB is updated), this method will be
	// called. This is critical, as the BlobQueue must remember blobs it has already seen to prevent duplicate emission,
	// and this method allows it to prune its data structures without risk of re-emitting a blob.
	StopTracking(blobKey corev2.BlobKey)
}

var _ BlobQueue = (*blobQueue)(nil)

// A stab implementation of BlobQueue.
type blobQueue struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger logging.Logger

	// Provides blobs that need to be put into batches for dispersal.
	blobSource BlobSource

	// How often to poll for new blobs.
	pollInterval time.Duration

	// How long to wait when polling for new blobs before timing out.
	pollTimeout time.Duration

	// The channel with blobs to disperse.
	blobChan chan *BlobToDisperse

	// When StopTracking() is called, the blobKey is sent to this channel to be processed by the controlLoop goroutine.
	stopTrackingChan chan corev2.BlobKey

	// A set of blobs that have been previously observed, and that might be re-emitted by the BlobSource.
	// This is used to detect this re-emission and prevent duplicates from being sent to the blobQueue channel.
	observedBlobs map[corev2.BlobKey]struct{}

	// Encapsulates all metrics for the blob queue.
	metrics *blobQueueMetrics
}

// NewBlobQueue creates a new BlobQueue.
//
// If the metrics registry is nil, then no metrics will be registered.
func NewBlobQueue(
	ctx context.Context,
	logger logging.Logger,
	blobSource BlobSource,
	pollInterval time.Duration,
	pollTimeout time.Duration,
	queueSize uint64,
	registry *prometheus.Registry,
) BlobQueue {

	ctx, cancel := context.WithCancel(ctx)

	bq := &blobQueue{
		ctx:              ctx,
		cancel:           cancel,
		logger:           logger,
		blobSource:       blobSource,
		pollInterval:     pollInterval,
		pollTimeout:      pollTimeout,
		blobChan:         make(chan *BlobToDisperse, queueSize),
		observedBlobs:    make(map[corev2.BlobKey]struct{}, queueSize),
		stopTrackingChan: make(chan corev2.BlobKey, queueSize),
		metrics:          newBlobQueueMetrics(registry),
	}

	go bq.controlLoop()

	return bq
}

// Get the channel that will receive blobs to disperse.
func (q *blobQueue) GetBlobToDisperse() <-chan *BlobToDisperse {
	return q.blobChan
}

func (q *blobQueue) Close() {
	q.cancel()
}

func (q *blobQueue) StopTracking(blobKey corev2.BlobKey) {
	select {
	case <-q.ctx.Done():
		return
	case q.stopTrackingChan <- blobKey:
	}
}

// A goroutine that runs in the background and polls for blobs to disperse.
func (q *blobQueue) controlLoop() {
	ticker := time.NewTicker(q.pollInterval)
	defer ticker.Stop()
	defer close(q.blobChan)
	defer close(q.stopTrackingChan)

	for {
		select {
		case <-q.ctx.Done():
			return
		case blobKey := <-q.stopTrackingChan:
			q.stopTracking(blobKey)
		case <-ticker.C:
			q.pollForBlobs()
			q.metrics.reportQueueSize(uint64(len(q.blobChan)))
			q.metrics.reportDedupSetSize(uint64(len(q.observedBlobs)))
		}
	}
}

// Stops tracking the given blob key now that it is known that the BlobSource will no longer return it.
func (q *blobQueue) stopTracking(blobKey corev2.BlobKey) {
	delete(q.observedBlobs, blobKey)
}

// Attempts to get blobs from the BlobSource and put them into the queue.
func (q *blobQueue) pollForBlobs() {

	// Enforce a timeout for the sake of sanity.
	ctx, cancel := context.WithTimeout(q.ctx, q.pollTimeout)
	defer cancel()

	start := time.Now()

	// Get a batch of blobs to disperse.
	blobMetadata, err := q.blobSource.GetBlobsToDisperse(ctx)
	if err != nil {
		q.logger.Errorf("Error getting blobs to disperse: %v", err)
	}

	elapsed := time.Since(start)
	duplicateBlobCount := uint64(0)

	// For each blob, make sure it is unique and put it into the channel.
	for _, metadata := range blobMetadata {
		blobKey, err := metadata.BlobHeader.BlobKey()
		if err != nil {
			q.logger.Errorf("Error getting blob key: %v", err)
			continue
		}

		if q.isDuplicate(blobKey) {
			duplicateBlobCount++
			continue
		}

		select {
		case <-q.ctx.Done():
			return
		case q.blobChan <- newBlobToDisperse(q.metrics, metadata, blobKey):
		}
	}

	q.metrics.reportBlobSourcePoll(uint64(len(blobMetadata)), duplicateBlobCount, elapsed)
}

// Returns true if the blob has already been observed, and false otherwise. This method updates the internal
// state to mark the blob as observed, meaning that this method will return false exactly once for each unique blobKey
// and true for all subsequent calls with the same blobKey.
func (q *blobQueue) isDuplicate(blobKey corev2.BlobKey) bool {
	_, alreadyObserved := q.observedBlobs[blobKey]
	if alreadyObserved {
		return true
	}

	q.observedBlobs[blobKey] = struct{}{}
	return false
}
