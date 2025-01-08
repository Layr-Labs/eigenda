package controller_test

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	commonmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	dispcommon "github.com/Layr-Labs/eigenda/disperser/common"
	commonv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	dispmock "github.com/Layr-Labs/eigenda/disperser/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gammazero/workerpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	blockNumber = uint64(100)
)

type testComponents struct {
	EncodingManager *controller.EncodingManager
	Pool            common.WorkerPool
	EncodingClient  *dispmock.MockEncoderClientV2
	ChainReader     *coremock.MockWriter
	MockPool        *commonmock.MockWorkerpool
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
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(&encoding.FragmentInfo{
		TotalChunkSizeBytes: 100,
		FragmentSizeBytes:   1024 * 1024 * 4,
	}, nil)

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

	deleteBlobs(t, blobMetadataStore, []corev2.BlobKey{blobKey1}, nil)
}

func TestEncodingManagerHandleManyBatches(t *testing.T) {
	ctx := context.Background()
	numBlobs := 12
	keys := make([]corev2.BlobKey, numBlobs)
	headers := make([]*corev2.BlobHeader, numBlobs)
	metadata := make([]*commonv2.BlobMetadata, numBlobs)
	for i := 0; i < numBlobs; i++ {
		keys[i], headers[i] = newBlob(t, []core.QuorumID{0, 1})
		now := time.Now()
		metadata[i] = &commonv2.BlobMetadata{
			BlobHeader: headers[i],
			BlobStatus: commonv2.Queued,
			Expiry:     uint64(now.Add(time.Hour).Unix()),
			NumRetries: 0,
			UpdatedAt:  uint64(now.UnixNano()),
		}
		err := blobMetadataStore.PutBlobMetadata(ctx, metadata[i])
		require.NoError(t, err)
	}

	c := newTestComponents(t, true)
	c.MockPool.On("Submit", mock.Anything).Return(nil).Times(numBlobs + 1)

	numIterations := (numBlobs + int(c.EncodingManager.MaxNumBlobsPerIteration) - 1) / int(c.EncodingManager.MaxNumBlobsPerIteration)
	expectedNumTasks := 0
	for i := 0; i < numIterations; i++ {
		err := c.EncodingManager.HandleBatch(ctx)
		require.NoError(t, err)
		if i < numIterations-1 {
			expectedNumTasks += int(c.EncodingManager.MaxNumBlobsPerIteration)
			c.MockPool.AssertNumberOfCalls(t, "Submit", expectedNumTasks)
		} else {
			expectedNumTasks += numBlobs % int(c.EncodingManager.MaxNumBlobsPerIteration)
			c.MockPool.AssertNumberOfCalls(t, "Submit", expectedNumTasks)
		}
	}
	err := c.EncodingManager.HandleBatch(ctx)
	require.ErrorContains(t, err, "no blobs to encode")

	// new record
	key, header := newBlob(t, []core.QuorumID{0, 1})
	now := time.Now()
	meta := &commonv2.BlobMetadata{
		BlobHeader: header,
		BlobStatus: commonv2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err = blobMetadataStore.PutBlobMetadata(ctx, meta)
	require.NoError(t, err)
	err = c.EncodingManager.HandleBatch(ctx)
	require.NoError(t, err)
	c.MockPool.AssertNumberOfCalls(t, "Submit", expectedNumTasks+1)

	deleteBlobs(t, blobMetadataStore, keys, nil)
	deleteBlobs(t, blobMetadataStore, []corev2.BlobKey{key}, nil)
}

func TestEncodingManagerHandleBatchNoBlobs(t *testing.T) {
	ctx := context.Background()
	c := newTestComponents(t, false)
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
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
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(&encoding.FragmentInfo{
		TotalChunkSizeBytes: 100,
		FragmentSizeBytes:   1024 * 1024 * 4,
	}, nil)

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
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Twice()

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
	logger := logging.NewNoopLogger()
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

	em, err := controller.NewEncodingManager(&controller.EncodingManagerConfig{
		PullInterval:                1 * time.Second,
		EncodingRequestTimeout:      5 * time.Second,
		StoreTimeout:                5 * time.Second,
		NumEncodingRetries:          1,
		NumRelayAssignment:          2,
		AvailableRelays:             []corev2.RelayKey{0, 1, 2, 3},
		MaxNumBlobsPerIteration:     5,
		OnchainStateRefreshInterval: onchainRefreshInterval,
	}, blobMetadataStore, pool, encodingClient, chainReader, logger, prometheus.NewRegistry())
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
	}
}
