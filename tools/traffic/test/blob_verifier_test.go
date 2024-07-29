package test

import (
	"context"
	"fmt"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/Layr-Labs/eigenda/tools/traffic/workers"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"sync"
	"testing"
	"time"
)

func getRandomStatus() disperser_rpc.BlobStatus {
	return disperser_rpc.BlobStatus(rand.Intn(7))
}

func isStatusTerminal(status disperser_rpc.BlobStatus) bool {
	switch status {
	case disperser_rpc.BlobStatus_UNKNOWN:
		return false
	case disperser_rpc.BlobStatus_PROCESSING:
		return false
	case disperser_rpc.BlobStatus_DISPERSING:
		return false

	case disperser_rpc.BlobStatus_INSUFFICIENT_SIGNATURES:
		return true
	case disperser_rpc.BlobStatus_FAILED:
		return true
	case disperser_rpc.BlobStatus_FINALIZED:
		return true
	case disperser_rpc.BlobStatus_CONFIRMED:
		return true
	default:
		panic("unknown status")
	}
}

func isStatusSuccess(status disperser_rpc.BlobStatus) bool {
	switch status {
	case disperser_rpc.BlobStatus_CONFIRMED:
		return true
	case disperser_rpc.BlobStatus_FINALIZED:
		return true
	default:
		return false
	}
}

func TestBlobVerifier(t *testing.T) {
	tu.InitializeRandom()

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.Nil(t, err)
	startTime := time.Unix(rand.Int63()%2_000_000_000, 0)
	ticker := newMockTicker(startTime)

	requiredDownloads := rand.Intn(10)
	config := &workers.Config{
		RequiredDownloads: float64(requiredDownloads),
	}

	blobTable := table.NewBlobTable()

	verifierMetrics := metrics.NewMockMetrics()

	lock := sync.Mutex{}

	disperserClient := newMockDisperserClient(t, &lock, true)

	verifier := workers.NewBlobVerifier(
		&ctx,
		&waitGroup,
		logger,
		ticker,
		config,
		&blobTable,
		disperserClient,
		verifierMetrics)

	verifier.Start()

	expectedGetStatusCount := 0
	statusCounts := make(map[disperser_rpc.BlobStatus]int)
	checksums := make(map[string][16]byte)
	sizes := make(map[string]uint)

	for i := 0; i < 100; i++ {

		// Add some new keys to track.
		newKeys := rand.Intn(10)
		for j := 0; j < newKeys; j++ {
			key := make([]byte, 16)
			checksum := [16]byte{}
			size := rand.Uint32()

			_, err = rand.Read(key)
			assert.Nil(t, err)
			_, err = rand.Read(checksum[:])
			assert.Nil(t, err)

			checksums[string(key)] = checksum
			sizes[string(key)] = uint(size)

			stringifiedKey := string(key)
			disperserClient.StatusMap[stringifiedKey] = disperser_rpc.BlobStatus_UNKNOWN

			verifier.AddUnconfirmedKey(&key, &checksum, uint(size))
		}

		// Choose some new statuses to be returned.
		// Count the number of status queries we expect to see in this iteration.
		for key, status := range disperserClient.StatusMap {
			if !isStatusTerminal(status) {
				// Blobs in a non-terminal status will be queried again.
				expectedGetStatusCount += 1
				// Set the next status to be returned.
				newStatus := getRandomStatus()
				disperserClient.StatusMap[key] = newStatus
				statusCounts[newStatus] += 1
			}
		}

		// Advance to the next cycle, allowing the verifier to process the new keys.
		ticker.Tick(time.Second)

		// Data is inserted asynchronously, so we may need to wait for it to be processed.
		tu.AssertEventuallyTrue(t, func() bool {
			lock.Lock()
			defer lock.Unlock()

			// Validate the number of calls made to the disperser client.
			if int(disperserClient.GetStatusCount) < expectedGetStatusCount {
				return false
			}

			// Read the data in the table into a map for quick lookup.
			tableData := make(map[string]*table.BlobMetadata)
			for i := uint(0); i < blobTable.Size(); i++ {
				metadata := blobTable.Get(i)
				tableData[string(*metadata.Key())] = metadata
			}

			blobsInFlight := 0
			for key, status := range disperserClient.StatusMap {
				metadata, present := tableData[key]

				if !isStatusTerminal(status) {
					blobsInFlight++
				}

				if isStatusSuccess(status) {
					// Successful blobs should be in the table.
					if !present {
						// Blob might not yet be in table due to timing.
						return false
					}
				} else {
					// Non-successful blobs should not be in the table.
					assert.False(t, present)
				}

				// Verify metadata.
				if present {
					assert.Equal(t, checksums[key], *metadata.Checksum())
					assert.Equal(t, sizes[key], metadata.Size())
					assert.Equal(t, requiredDownloads, metadata.RemainingReadPermits())
				}
			}

			// Verify metrics.
			for status, count := range statusCounts {
				metricName := fmt.Sprintf("get_status_%s", status.String())
				if float64(count) != verifierMetrics.GetCount(metricName) {
					return false
				}
			}
			if float64(blobsInFlight) != verifierMetrics.GetGaugeValue("blobs_in_flight") {
				fmt.Printf("expected blobs_in_flight to be %d, got %f\n", blobsInFlight, verifierMetrics.GetCount("blobs_in_flight"))
				return false
			}

			return true
		}, time.Second)
	}

	assert.Equal(t, expectedGetStatusCount, int(disperserClient.GetStatusCount))
	assert.Equal(t, 0, int(disperserClient.DisperseCount))

	cancel()
	tu.ExecuteWithTimeout(func() {
		waitGroup.Wait()
	}, time.Second)
}
