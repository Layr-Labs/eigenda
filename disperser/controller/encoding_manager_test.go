package controller_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	commonmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	dispcommon "github.com/Layr-Labs/eigenda/disperser/common"
	commonv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	dispmock "github.com/Layr-Labs/eigenda/disperser/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/gammazero/workerpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	blockNumber = uint64(100)
)

type HeartbeatMessage struct {
	Component string    // e.g., "encodingManager" or "dispatcher"
	Timestamp time.Time // when the heartbeat was sent
}

type testComponents struct {
	EncodingManager *controller.EncodingManager
	Pool            common.WorkerPool
	EncodingClient  *dispmock.MockEncoderClientV2
	ChainReader     *coremock.MockWriter
	MockPool        *commonmock.MockWorkerpool
	BlobSet         *controller.MockBlobSet
	GetHeartbeats   func() []HeartbeatMessage
}

func TestGetRelayKeys(t *testing.T) {
	// Test cases for GetRelayKeys function
	tests := []struct {
		name            string
		numRelays       uint16
		availableRelays []corev2.RelayKey
		err             error
	}{
		{
			name:            "Single relay",
			numRelays:       1,
			availableRelays: []corev2.RelayKey{0},
			err:             nil,
		},
		{
			name:            "Choose more than whats available",
			numRelays:       2,
			availableRelays: []corev2.RelayKey{0},
			err:             nil,
		},
		{
			name:            "All relays",
			numRelays:       2,
			availableRelays: []corev2.RelayKey{0, 1},
			err:             nil,
		},
		{
			name:            "Choose 1 from multiple relays",
			numRelays:       3,
			availableRelays: []corev2.RelayKey{0, 1, 2, 3},
			err:             nil,
		},
		{
			name:            "Choose 2 from multiple relays",
			numRelays:       2,
			availableRelays: []corev2.RelayKey{0, 1, 2, 3},
			err:             nil,
		},
		{
			name:            "No relays",
			numRelays:       0,
			availableRelays: []corev2.RelayKey{},
			err:             nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			availableRelaysCopy := make([]corev2.RelayKey, len(tt.availableRelays))
			copy(availableRelaysCopy, tt.availableRelays)

			got, err := controller.GetRelayKeys(tt.numRelays, tt.availableRelays)
			if err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, tt.err)
				require.Len(t, got, int(tt.numRelays))
				seen := make(map[corev2.RelayKey]struct{})
				for _, relay := range got {
					require.Contains(t, tt.availableRelays, relay)
					seen[relay] = struct{}{}
				}
				require.Equal(t, len(seen), len(got))
				// GetRelayKeys should not modify the original list of available relays.
				require.Equal(t, availableRelaysCopy, tt.availableRelays)
			}
		})
	}
}

func TestEncodingManagerHandleBatch(t *testing.T) {
	ctx := context.Background()
	blobKey1, blobHeader1 := newBlob(t, []core.QuorumID{0, 1})
	now := time.Now()
	metadata1 := &commonv2.BlobMetadata{
		BlobHeader: blobHeader1,
		BlobStatus: commonv2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err := blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	require.NoError(t, err)

	c := newTestComponents(t, false)
	c.BlobSet.On("Contains", mock.Anything).Return(false)
	c.BlobSet.On("AddBlob", mock.Anything).Return(nil)
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(&encoding.FragmentInfo{
		TotalChunkSizeBytes: 100,
		FragmentSizeBytes:   1024 * 1024 * 4,
	}, nil)

	defer func() {
		heartbeats := c.GetHeartbeats()
		require.NotEmpty(t, heartbeats, "Expected heartbeats, but none were received")

		// Verify that the heartbeat messages have the correct component identifier.
		for _, hb := range heartbeats {
			require.Equal(t, "encodingManager", hb.Component, "Expected heartbeat from encodingManager")
		}

		// Check that timestamps are increasing
		for i := 1; i < len(heartbeats); i++ {
			require.True(t, heartbeats[i].Timestamp.After(heartbeats[i-1].Timestamp) ||
				heartbeats[i].Timestamp.Equal(heartbeats[i-1].Timestamp),
				"Heartbeat timestamps should be increasing")
		}
	}()

	err = c.EncodingManager.HandleBatch(ctx)
	require.NoError(t, err)
	c.Pool.StopWait()
	c.BlobSet.AssertCalled(t, "Contains", blobKey1)
	c.BlobSet.AssertCalled(t, "AddBlob", blobKey1)

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	require.NoError(t, err)
	require.Equal(t, commonv2.Encoded, fetchedMetadata.BlobStatus)
	require.Greater(t, fetchedMetadata.UpdatedAt, metadata1.UpdatedAt)

	fetchedCert, fetchedFragmentInfo, err := blobMetadataStore.GetBlobCertificate(ctx, blobKey1)
	require.NoError(t, err)
	require.Equal(t, fetchedCert.BlobHeader, blobHeader1)
	for _, relayKey := range fetchedCert.RelayKeys {
		require.Contains(t, c.EncodingManager.AvailableRelays, relayKey)
	}
	require.Equal(t, fetchedFragmentInfo.TotalChunkSizeBytes, uint32(100))
	require.Equal(t, fetchedFragmentInfo.FragmentSizeBytes, uint32(1024*1024*4))

	deleteBlobs(t, blobMetadataStore, []corev2.BlobKey{blobKey1}, nil)
}

func TestEncodingManagerHandleBatchDedup(t *testing.T) {
	ctx := context.Background()
	blobKey1, blobHeader1 := newBlob(t, []core.QuorumID{0, 1})
	now := time.Now()
	metadata1 := &commonv2.BlobMetadata{
		BlobHeader: blobHeader1,
		BlobStatus: commonv2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err := blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	require.NoError(t, err)

	c := newTestComponents(t, false)
	c.BlobSet.On("Contains", blobKey1).Return(true).Once()
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(&encoding.FragmentInfo{
		TotalChunkSizeBytes: 100,
		FragmentSizeBytes:   1024 * 1024 * 4,
	}, nil)

	defer func() {
		heartbeats := c.GetHeartbeats()
		require.NotEmpty(t, heartbeats, "Expected heartbeats, but none were received")

		// Verify that the heartbeat messages have the correct component identifier.
		for _, hb := range heartbeats {
			require.Equal(t, "encodingManager", hb.Component, "Expected heartbeat from encodingManager")
		}

		// Check that timestamps are increasing
		for i := 1; i < len(heartbeats); i++ {
			require.True(t, heartbeats[i].Timestamp.After(heartbeats[i-1].Timestamp) ||
				heartbeats[i].Timestamp.Equal(heartbeats[i-1].Timestamp),
				"Heartbeat timestamps should be increasing")
		}
	}()

	err = c.EncodingManager.HandleBatch(ctx)
	require.ErrorContains(t, err, "no blobs to encode")
	c.Pool.StopWait()
	c.BlobSet.AssertCalled(t, "Contains", blobKey1)
	c.BlobSet.AssertNotCalled(t, "AddBlob", blobKey1)
	deleteBlobs(t, blobMetadataStore, []corev2.BlobKey{blobKey1}, nil)
}

func TestEncodingManagerHandleManyBatches(t *testing.T) {
	ctx := context.Background()
	numBlobs := 12
	keys := make([]corev2.BlobKey, 0)
	headers := make([]*corev2.BlobHeader, 0)
	metadata := make([]*commonv2.BlobMetadata, 0)
	for i := 0; i < numBlobs; i++ {
		k, h := newBlob(t, []core.QuorumID{0, 1})
		keys = append(keys, k)
		headers = append(headers, h)
		now := time.Now()
		metadata = append(metadata, &commonv2.BlobMetadata{
			BlobHeader: headers[i],
			BlobStatus: commonv2.Queued,
			Expiry:     uint64(now.Add(time.Hour).Unix()),
			NumRetries: 0,
			UpdatedAt:  uint64(now.UnixNano()),
		})
		err := blobMetadataStore.PutBlobMetadata(ctx, metadata[i])
		require.NoError(t, err)
	}

	c := newTestComponents(t, true)
	c.BlobSet.On("Contains", mock.Anything).Return(false)
	c.BlobSet.On("AddBlob", mock.Anything).Return(nil)
	numIterations := (numBlobs + int(c.EncodingManager.MaxNumBlobsPerIteration) - 1) / int(c.EncodingManager.MaxNumBlobsPerIteration)
	c.MockPool.On("Submit", mock.Anything).Return(nil).Times(numBlobs + numIterations)

	expectedNumTasks := 0
	for i := 0; i < numIterations; i++ {
		err := c.EncodingManager.HandleBatch(ctx)
		require.NoError(t, err)
		if i < numIterations-1 {
			expectedNumTasks += int(c.EncodingManager.MaxNumBlobsPerIteration)
			c.MockPool.AssertNumberOfCalls(t, "Submit", expectedNumTasks)

			// add blobs to the queue with UpdatedAt in the past
			// these should be skipped in this loop
			key, header := newBlob(t, []core.QuorumID{0, 1})
			keys = append(keys, key)
			now := time.Now()
			meta := &commonv2.BlobMetadata{
				BlobHeader: header,
				BlobStatus: commonv2.Queued,
				Expiry:     uint64(now.Add(time.Hour).Unix()),
				NumRetries: 0,
				UpdatedAt:  uint64(now.Add(-time.Hour).UnixNano()),
			}
			err := blobMetadataStore.PutBlobMetadata(ctx, meta)
			require.NoError(t, err)
		} else {
			expectedNumTasks += numBlobs % int(c.EncodingManager.MaxNumBlobsPerIteration)
			c.MockPool.AssertNumberOfCalls(t, "Submit", expectedNumTasks)
		}
	}

	for i := 0; i < numBlobs; i++ {
		err := blobMetadataStore.UpdateBlobStatus(ctx, keys[i], commonv2.Encoded)
		require.NoError(t, err)
	}

	// should handle blobs with UpdatedAt in the past
	err := c.EncodingManager.HandleBatch(ctx)
	require.NoError(t, err)
	c.MockPool.AssertNumberOfCalls(t, "Submit", expectedNumTasks+numIterations-1)

	for i := 0; i < numIterations-1; i++ {
		err := blobMetadataStore.UpdateBlobStatus(ctx, keys[numBlobs+i], commonv2.Encoded)
		require.NoError(t, err)
	}

	defer func() {
		heartbeats := c.GetHeartbeats()
		require.NotEmpty(t, heartbeats, "Expected heartbeats, but none were received")

		// Verify that the heartbeat messages have the correct component identifier.
		for _, hb := range heartbeats {
			require.Equal(t, "encodingManager", hb.Component, "Expected heartbeat from encodingManager")
		}

		// Check that timestamps are increasing
		for i := 1; i < len(heartbeats); i++ {
			require.True(t, heartbeats[i].Timestamp.After(heartbeats[i-1].Timestamp) ||
				heartbeats[i].Timestamp.Equal(heartbeats[i-1].Timestamp),
				"Heartbeat timestamps should be increasing")
		}
	}()

	// no more blobs to encode
	err = c.EncodingManager.HandleBatch(ctx)
	require.ErrorContains(t, err, "no blobs to encode")

	deleteBlobs(t, blobMetadataStore, keys, nil)
}

func TestEncodingManagerHandleBatchNoBlobs(t *testing.T) {
	ctx := context.Background()
	c := newTestComponents(t, false)
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	defer func() {
		heartbeats := c.GetHeartbeats()
		require.NotEmpty(t, heartbeats, "Expected heartbeats, but none were received")

		// Verify that the heartbeat messages have the correct component identifier.
		for _, hb := range heartbeats {
			require.Equal(t, "encodingManager", hb.Component, "Expected heartbeat from encodingManager")
		}

		// Check that timestamps are increasing
		for i := 1; i < len(heartbeats); i++ {
			require.True(t, heartbeats[i].Timestamp.After(heartbeats[i-1].Timestamp) ||
				heartbeats[i].Timestamp.Equal(heartbeats[i-1].Timestamp),
				"Heartbeat timestamps should be increasing")
		}
	}()

	err := c.EncodingManager.HandleBatch(ctx)
	require.ErrorContains(t, err, "no blobs to encode")
}

func TestEncodingManagerHandleBatchRetrySuccess(t *testing.T) {
	ctx := context.Background()
	blobKey1, blobHeader1 := newBlob(t, []core.QuorumID{0, 1})
	now := time.Now()
	metadata1 := &commonv2.BlobMetadata{
		BlobHeader: blobHeader1,
		BlobStatus: commonv2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err := blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	require.NoError(t, err)

	c := newTestComponents(t, false)
	c.BlobSet.On("Contains", mock.Anything).Return(false)
	c.BlobSet.On("AddBlob", mock.Anything).Return(nil)
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(&encoding.FragmentInfo{
		TotalChunkSizeBytes: 100,
		FragmentSizeBytes:   1024 * 1024 * 4,
	}, nil)

	defer func() {
		heartbeats := c.GetHeartbeats()
		require.NotEmpty(t, heartbeats, "Expected heartbeats, but none were received")

		// Verify that the heartbeat messages have the correct component identifier.
		for _, hb := range heartbeats {
			require.Equal(t, "encodingManager", hb.Component, "Expected heartbeat from encodingManager")
		}

		// Check that timestamps are increasing
		for i := 1; i < len(heartbeats); i++ {
			require.True(t, heartbeats[i].Timestamp.After(heartbeats[i-1].Timestamp) ||
				heartbeats[i].Timestamp.Equal(heartbeats[i-1].Timestamp),
				"Heartbeat timestamps should be increasing")
		}
	}()

	err = c.EncodingManager.HandleBatch(ctx)
	require.NoError(t, err)
	c.Pool.StopWait()

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	require.NoError(t, err)
	require.Equal(t, commonv2.Encoded, fetchedMetadata.BlobStatus)
	require.Greater(t, fetchedMetadata.UpdatedAt, metadata1.UpdatedAt)

	fetchedCert, fetchedFragmentInfo, err := blobMetadataStore.GetBlobCertificate(ctx, blobKey1)
	require.NoError(t, err)
	require.Equal(t, fetchedCert.BlobHeader, blobHeader1)
	for _, relayKey := range fetchedCert.RelayKeys {
		require.Contains(t, c.EncodingManager.AvailableRelays, relayKey)
	}
	require.Equal(t, fetchedFragmentInfo.TotalChunkSizeBytes, uint32(100))
	require.Equal(t, fetchedFragmentInfo.FragmentSizeBytes, uint32(1024*1024*4))
	c.EncodingClient.AssertNumberOfCalls(t, "EncodeBlob", 2)

	deleteBlobs(t, blobMetadataStore, []corev2.BlobKey{blobKey1}, nil)
}

func TestEncodingManagerHandleBatchRetryFailure(t *testing.T) {
	ctx := context.Background()
	blobKey1, blobHeader1 := newBlob(t, []core.QuorumID{0, 1})
	now := time.Now()
	metadata1 := &commonv2.BlobMetadata{
		BlobHeader: blobHeader1,
		BlobStatus: commonv2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err := blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	require.NoError(t, err)

	c := newTestComponents(t, false)
	c.BlobSet.On("Contains", mock.Anything).Return(false)
	c.BlobSet.On("AddBlob", mock.Anything).Return(nil)
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Twice()

	defer func() {
		heartbeats := c.GetHeartbeats()
		require.NotEmpty(t, heartbeats, "Expected heartbeats, but none were received")

		// Verify that the heartbeat messages have the correct component identifier.
		for _, hb := range heartbeats {
			require.Equal(t, "encodingManager", hb.Component, "Expected heartbeat from encodingManager")
		}

		// Check that timestamps are increasing
		for i := 1; i < len(heartbeats); i++ {
			require.True(t, heartbeats[i].Timestamp.After(heartbeats[i-1].Timestamp) ||
				heartbeats[i].Timestamp.Equal(heartbeats[i-1].Timestamp),
				"Heartbeat timestamps should be increasing")
		}
	}()

	err = c.EncodingManager.HandleBatch(ctx)
	require.NoError(t, err)
	c.Pool.StopWait()

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	require.NoError(t, err)
	// marked as failed
	require.Equal(t, commonv2.Failed, fetchedMetadata.BlobStatus)
	require.Greater(t, fetchedMetadata.UpdatedAt, metadata1.UpdatedAt)

	fetchedCert, fetchedFragmentInfo, err := blobMetadataStore.GetBlobCertificate(ctx, blobKey1)
	require.ErrorIs(t, err, dispcommon.ErrMetadataNotFound)
	require.Nil(t, fetchedCert)
	require.Nil(t, fetchedFragmentInfo)
	c.EncodingClient.AssertNumberOfCalls(t, "EncodeBlob", 2)

	deleteBlobs(t, blobMetadataStore, []corev2.BlobKey{blobKey1}, nil)
}

func newTestComponents(t *testing.T, mockPool bool) *testComponents {
	logger := testutils.GetLogger()
	// logger, err := common.NewLogger(common.DefaultLoggerConfig())
	// require.NoError(t, err)
	var pool common.WorkerPool
	var mockP *commonmock.MockWorkerpool
	if mockPool {
		mockP = &commonmock.MockWorkerpool{}
		pool = mockP
	} else {
		pool = workerpool.New(5)
	}
	encodingClient := dispmock.NewMockEncoderClientV2()
	chainReader := &coremock.MockWriter{}
	chainReader.On("GetCurrentBlockNumber").Return(blockNumber, nil)
	chainReader.On("GetAllVersionedBlobParams", mock.Anything).Return(map[v2.BlobVersion]*core.BlobVersionParameters{
		0: {
			NumChunks:       8192,
			CodingRate:      8,
			MaxNumOperators: 3537,
		},
	}, nil)
	onchainRefreshInterval := 1 * time.Millisecond
	blobSet := &controller.MockBlobSet{}
	blobSet.On("Size", mock.Anything).Return(0)

	// Heartbeat tracking variables
	var mu sync.Mutex
	var heartbeatsReceived []HeartbeatMessage
	doneListening := make(chan bool)

	// Mocked signalHeartbeat function
	mockSignalHeartbeat := func() {
		mu.Lock()
		defer mu.Unlock()
		heartbeatsReceived = append(heartbeatsReceived, HeartbeatMessage{
			Component: "encodingManager",
			Timestamp: time.Now(),
		})
	}

	em, err := controller.NewEncodingManager(&controller.EncodingManagerConfig{
		PullInterval:                1 * time.Second,
		EncodingRequestTimeout:      5 * time.Second,
		StoreTimeout:                5 * time.Second,
		NumEncodingRetries:          1,
		NumRelayAssignment:          2,
		AvailableRelays:             []corev2.RelayKey{0, 1, 2, 3},
		MaxNumBlobsPerIteration:     5,
		OnchainStateRefreshInterval: onchainRefreshInterval,
	}, blobMetadataStore, pool, encodingClient, chainReader, logger, prometheus.NewRegistry(), blobSet, mockSignalHeartbeat)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*onchainRefreshInterval)
	defer cancel()
	// Start the encoding manager to fetch the onchain state
	_ = em.Start(ctx)
	return &testComponents{
		EncodingManager: em,
		Pool:            pool,
		EncodingClient:  encodingClient,
		ChainReader:     chainReader,
		MockPool:        mockP,
		BlobSet:         blobSet,
		GetHeartbeats: func() []HeartbeatMessage {
			close(doneListening) // Stop tracking
			mu.Lock()
			defer mu.Unlock()
			return heartbeatsReceived
		},
	}
}
