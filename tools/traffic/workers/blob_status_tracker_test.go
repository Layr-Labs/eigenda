package workers

import (
	"context"
	"fmt"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	requiredDownloads := rand.Intn(10) + 1
	config := &config.WorkerConfig{
		RequiredDownloads: float64(requiredDownloads),
	}

	blobStore := table.NewBlobStore()

	verifierMetrics := metrics.NewMockMetrics()

	disperserClient := &MockDisperserClient{}

	verifier := NewBlobVerifier(
		&ctx,
		&waitGroup,
		logger,
		config,
		make(chan *UnconfirmedKey),
		blobStore,
		disperserClient,
		verifierMetrics)

	expectedGetStatusCount := 0
	statusCounts := make(map[disperser_rpc.BlobStatus]int)
	checksums := make(map[string][16]byte)
	sizes := make(map[string]uint)

	statusMap := make(map[string]disperser_rpc.BlobStatus)

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
			statusMap[stringifiedKey] = disperser_rpc.BlobStatus_UNKNOWN

			unconfirmedKey := &UnconfirmedKey{
				Key:            key,
				Checksum:       checksum,
				Size:           uint(size),
				SubmissionTime: time.Now(),
			}

			verifier.unconfirmedBlobs = append(verifier.unconfirmedBlobs, unconfirmedKey)
		}

		// Reset the mock disperser client.
		disperserClient.mock = mock.Mock{}
		expectedGetStatusCount = 0

		// Choose some new statuses to be returned.
		// Count the number of status queries we expect to see in this iteration.
		for key, status := range statusMap {
			var newStatus disperser_rpc.BlobStatus
			if isStatusTerminal(status) {
				newStatus = status
			} else {
				// Blobs in a non-terminal status will be queried again.
				expectedGetStatusCount += 1
				// Set the next status to be returned.
				newStatus = getRandomStatus()
				statusMap[key] = newStatus

				statusCounts[newStatus] += 1
			}
			disperserClient.mock.On("GetBlobStatus", []byte(key)).Return(
				&disperser_rpc.BlobStatusReply{
					Status: newStatus,
					Info: &disperser_rpc.BlobInfo{
						BlobVerificationProof: &disperser_rpc.BlobVerificationProof{
							BatchMetadata: &disperser_rpc.BatchMetadata{
								BatchHeaderHash: make([]byte, 0),
							},
						},
					},
				}, nil)
		}

		// Simulate advancement of time, allowing the verifier to process the new keys.
		verifier.poll()

		// Validate the number of calls made to the disperser client.
		disperserClient.mock.AssertNumberOfCalls(t, "GetBlobStatus", expectedGetStatusCount)

		// Read the data in the confirmedBlobs into a map for quick lookup.
		tableData := make(map[string]*table.BlobMetadata)
		for _, metadata := range blobStore.GetAll() {
			tableData[string(metadata.Key)] = metadata
		}

		blobsInFlight := 0
		for key, status := range statusMap {
			metadata, present := tableData[key]

			if !isStatusTerminal(status) {
				blobsInFlight++
			}

			if isStatusSuccess(status) {
				// Successful blobs should be in the confirmedBlobs.
				assert.True(t, present)
			} else {
				// Non-successful blobs should not be in the confirmedBlobs.
				assert.False(t, present)
			}

			// Verify metadata.
			if present {
				assert.Equal(t, checksums[key], metadata.Checksum)
				assert.Equal(t, sizes[key], metadata.Size)
				assert.Equal(t, requiredDownloads, metadata.RemainingReadPermits)
			}
		}

		// Verify metrics.
		for status, count := range statusCounts { // TODO
			metricName := fmt.Sprintf("get_status_%s", status.String())
			assert.Equal(t, float64(count), verifierMetrics.GetCount(metricName), "status: %s", status.String())
		}
		if float64(blobsInFlight) != verifierMetrics.GetGaugeValue("blobs_in_flight") {
			assert.Equal(t, float64(blobsInFlight), verifierMetrics.GetGaugeValue("blobs_in_flight"))
		}
	}

	cancel()
	tu.ExecuteWithTimeout(func() {
		waitGroup.Wait()
	}, time.Second)
}
