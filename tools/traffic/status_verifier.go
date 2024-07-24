package traffic

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"time"
)

// StatusVerifier periodically polls the disperser service to verify the status of blobs that were recently written.
// When blobs become confirmed, the status verifier updates the blob table accordingly.
// This is a thread safe data structure.
type StatusVerifier struct {

	// A table of confirmed blobs. Blobs are added here when they are confirmed by the disperser service.
	table *BlobTable

	// The maximum number of reads permitted against an individual blob, or -1 if unlimited.
	blobReadLimit int32

	// The disperser client used to monitor the disperser service.
	dispenser *clients.DisperserClient

	// The keys of blobs that have not yet been confirmed by the disperser service.
	unconfirmedKeys []*[]byte

	// Newly added keys that require verification.
	keyChannel chan *[]byte

	getStatusLatencyMetric            LatencyMetric
	getStatusErrorCountMetric         CountMetric
	unknownCountMetric                CountMetric
	processingCountMetric             CountMetric
	dispersingCountMetric             CountMetric
	failedCountMetric                 CountMetric
	insufficientSignaturesCountMetric CountMetric
	confirmedCountMetric              CountMetric
	finalizedCountMetric              CountMetric

	// TODO metric for average time for blob to become confirmed
}

// NewStatusVerifier creates a new StatusVerifier instance.
func NewStatusVerifier(
	table *BlobTable,
	disperser *clients.DisperserClient,
	blobReadLimit int32,
	metrics *Metrics) StatusVerifier {

	return StatusVerifier{
		table:                             table,
		blobReadLimit:                     blobReadLimit,
		dispenser:                         disperser,
		unconfirmedKeys:                   make([]*[]byte, 0),
		keyChannel:                        make(chan *[]byte),
		getStatusLatencyMetric:            metrics.NewLatencyMetric("getStatus"),
		getStatusErrorCountMetric:         metrics.NewCountMetric("getStatusError"),
		unknownCountMetric:                metrics.NewCountMetric("getStatusUnknown"),
		processingCountMetric:             metrics.NewCountMetric("getStatusProcessing"),
		dispersingCountMetric:             metrics.NewCountMetric("getStatusDispersing"),
		failedCountMetric:                 metrics.NewCountMetric("getStatusFailed"),
		insufficientSignaturesCountMetric: metrics.NewCountMetric("getStatusInsufficientSignatures"),
		confirmedCountMetric:              metrics.NewCountMetric("getStatusConfirmed"),
		finalizedCountMetric:              metrics.NewCountMetric("getStatusFinalized"),
	}
}

// AddUnconfirmedKey adds a key to the list of unconfirmed keys.
func (verifier *StatusVerifier) AddUnconfirmedKey(key *[]byte) {
	verifier.keyChannel <- key
}

// Start begins the status goroutine, which periodically polls
// the disperser service to verify the status of blobs.
func (verifier *StatusVerifier) Start(ctx context.Context, period time.Duration) {
	go verifier.monitor(ctx, period)
}

// monitor periodically polls the disperser service to verify the status of blobs.
func (verifier *StatusVerifier) monitor(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			return
		case key := <-verifier.keyChannel:
			verifier.unconfirmedKeys = append(verifier.unconfirmedKeys, key)
		case <-ticker.C:
			verifier.poll(ctx)
		}
	}
}

// poll checks all unconfirmed keys to see if they have been confirmed by the disperser service.
// If a key is confirmed, it is added to the blob table and removed from the list of unconfirmed keys.
func (verifier *StatusVerifier) poll(ctx context.Context) {

	// TODO If the number of unconfirmed blobs is high and the time to confirm his high, this is not efficient.

	unconfirmedKeys := make([]*[]byte, 0)
	for _, key := range verifier.unconfirmedKeys {
		confirmed := verifier.checkStatusForBlob(ctx, key)
		if !confirmed {
			unconfirmedKeys = append(unconfirmedKeys, key)
		}
	}
	verifier.unconfirmedKeys = unconfirmedKeys
}

// checkStatusForBlob checks the status of a blob. Returns true if the final blob status is known, false otherwise.
func (verifier *StatusVerifier) checkStatusForBlob(ctx context.Context, key *[]byte) bool {
	status, err := (*verifier.dispenser).GetBlobStatus(ctx, *key)

	if err != nil {
		fmt.Println("Error getting blob status:", err) // TODO is this proper?
		verifier.getStatusErrorCountMetric.Increment()
		return false
	}

	switch status.GetStatus() {

	case disperser.BlobStatus_UNKNOWN:
		verifier.unknownCountMetric.Increment()
		return false
	case disperser.BlobStatus_PROCESSING:
		verifier.processingCountMetric.Increment()
		return false
	case disperser.BlobStatus_DISPERSING:
		verifier.dispersingCountMetric.Increment()
		return false

	case disperser.BlobStatus_FAILED:
		verifier.failedCountMetric.Increment()
		fmt.Println("Blob dispersal failed:", status.GetStatus()) // TODO use logger
		return true
	case disperser.BlobStatus_INSUFFICIENT_SIGNATURES:
		verifier.insufficientSignaturesCountMetric.Increment()
		fmt.Println("Blob dispersal failed:", status.GetStatus()) // TODO use logger
		return true

	case disperser.BlobStatus_CONFIRMED:
		verifier.confirmedCountMetric.Increment()
		return true
	case disperser.BlobStatus_FINALIZED:
		verifier.finalizedCountMetric.Increment()
		batchHeaderHash := status.GetInfo().BlobVerificationProof.BatchMetadata.BatchHeaderHash
		blobIndex := status.GetInfo().BlobVerificationProof.GetBlobIndex()

		blobMetadata := NewBlobMetadata(key, &batchHeaderHash, blobIndex, -1) // TODO permits
		verifier.table.Add(blobMetadata)
		fmt.Printf("Confirmed blob, batch header hash: %s, blobIndex %d\n", base64.StdEncoding.EncodeToString(batchHeaderHash), blobIndex)
		//fmt.Println("Confirmed blob")

		return true

	default:
		fmt.Println("Unknown blob status:", status.GetStatus()) // TODO use logger
		return true
	}
}
