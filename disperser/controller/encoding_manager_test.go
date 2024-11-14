package controller_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispcommon "github.com/Layr-Labs/eigenda/disperser/common"
	commonv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	dispmock "github.com/Layr-Labs/eigenda/disperser/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	blockNumber = uint64(100)
)

type testComponents struct {
	EncodingManager *controller.EncodingManager
	Pool            common.WorkerPool
	EncodingClient  *dispmock.MockEncoderClientV2
	ChainReader     *coremock.MockWriter
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
				assert.Error(t, err)
			} else {
				assert.NoError(t, tt.err)
				assert.Len(t, got, int(tt.numRelays))
				seen := make(map[corev2.RelayKey]struct{})
				for _, relay := range got {
					assert.Contains(t, tt.availableRelays, relay)
					seen[relay] = struct{}{}
				}
				assert.Equal(t, len(seen), len(got))
			}
		})
	}
}

func TestEncodingManagerHandleBatch(t *testing.T) {
	ctx := context.Background()
	blobHeader1 := &corev2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   []core.QuorumID{0},
		BlobCommitments: mockCommitment,
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x1234",
			BinIndex:          0,
			CumulativePayment: big.NewInt(532),
		},
	}
	blobKey1, err := blobHeader1.BlobKey()
	assert.NoError(t, err)
	now := time.Now()
	metadata1 := &commonv2.BlobMetadata{
		BlobHeader: blobHeader1,
		BlobStatus: commonv2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err = blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	assert.NoError(t, err)

	c := newTestComponents(t)
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(&encoding.FragmentInfo{
		TotalChunkSizeBytes: 100,
		FragmentSizeBytes:   1024 * 1024 * 4,
	}, nil)

	err = c.EncodingManager.HandleBatch(ctx)
	assert.NoError(t, err)
	c.Pool.StopWait()

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Equal(t, commonv2.Encoded, fetchedMetadata.BlobStatus)
	assert.Greater(t, fetchedMetadata.UpdatedAt, metadata1.UpdatedAt)

	fetchedCert, fetchedFragmentInfo, err := blobMetadataStore.GetBlobCertificate(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Equal(t, fetchedCert.BlobHeader, blobHeader1)
	for _, relayKey := range fetchedCert.RelayKeys {
		assert.Contains(t, c.EncodingManager.AvailableRelays, relayKey)
	}
	assert.Equal(t, fetchedFragmentInfo.TotalChunkSizeBytes, uint32(100))
	assert.Equal(t, fetchedFragmentInfo.FragmentSizeBytes, uint32(1024*1024*4))
}

func TestEncodingManagerHandleBatchNoBlobs(t *testing.T) {
	ctx := context.Background()
	c := newTestComponents(t)
	err := c.EncodingManager.HandleBatch(ctx)
	assert.ErrorContains(t, err, "no blobs to encode")
}

func TestEncodingManagerHandleBatchRetrySuccess(t *testing.T) {
	ctx := context.Background()
	blobHeader1 := &corev2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   []core.QuorumID{0},
		BlobCommitments: mockCommitment,
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x12345",
			BinIndex:          0,
			CumulativePayment: big.NewInt(532),
		},
	}
	blobKey1, err := blobHeader1.BlobKey()
	assert.NoError(t, err)
	now := time.Now()
	metadata1 := &commonv2.BlobMetadata{
		BlobHeader: blobHeader1,
		BlobStatus: commonv2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err = blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	assert.NoError(t, err)

	c := newTestComponents(t)
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(&encoding.FragmentInfo{
		TotalChunkSizeBytes: 100,
		FragmentSizeBytes:   1024 * 1024 * 4,
	}, nil)

	err = c.EncodingManager.HandleBatch(ctx)
	assert.NoError(t, err)
	c.Pool.StopWait()

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Equal(t, commonv2.Encoded, fetchedMetadata.BlobStatus)
	assert.Greater(t, fetchedMetadata.UpdatedAt, metadata1.UpdatedAt)

	fetchedCert, fetchedFragmentInfo, err := blobMetadataStore.GetBlobCertificate(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Equal(t, fetchedCert.BlobHeader, blobHeader1)
	for _, relayKey := range fetchedCert.RelayKeys {
		assert.Contains(t, c.EncodingManager.AvailableRelays, relayKey)
	}
	assert.Equal(t, fetchedFragmentInfo.TotalChunkSizeBytes, uint32(100))
	assert.Equal(t, fetchedFragmentInfo.FragmentSizeBytes, uint32(1024*1024*4))
	c.EncodingClient.AssertNumberOfCalls(t, "EncodeBlob", 2)
}

func TestEncodingManagerHandleBatchRetryFailure(t *testing.T) {
	ctx := context.Background()
	blobHeader1 := &corev2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   []core.QuorumID{0},
		BlobCommitments: mockCommitment,
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x123456",
			BinIndex:          0,
			CumulativePayment: big.NewInt(532),
		},
	}
	blobKey1, err := blobHeader1.BlobKey()
	assert.NoError(t, err)
	now := time.Now()
	metadata1 := &commonv2.BlobMetadata{
		BlobHeader: blobHeader1,
		BlobStatus: commonv2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err = blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	assert.NoError(t, err)

	c := newTestComponents(t)
	c.EncodingClient.On("EncodeBlob", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Twice()

	err = c.EncodingManager.HandleBatch(ctx)
	assert.NoError(t, err)
	c.Pool.StopWait()

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	// marked as failed
	assert.Equal(t, commonv2.Failed, fetchedMetadata.BlobStatus)
	assert.Greater(t, fetchedMetadata.UpdatedAt, metadata1.UpdatedAt)

	fetchedCert, fetchedFragmentInfo, err := blobMetadataStore.GetBlobCertificate(ctx, blobKey1)
	assert.ErrorIs(t, err, dispcommon.ErrMetadataNotFound)
	assert.Nil(t, fetchedCert)
	assert.Nil(t, fetchedFragmentInfo)
	c.EncodingClient.AssertNumberOfCalls(t, "EncodeBlob", 2)
}

func newTestComponents(t *testing.T) *testComponents {
	logger := logging.NewNoopLogger()
	// logger, err := common.NewLogger(common.DefaultLoggerConfig())
	// assert.NoError(t, err)
	pool := workerpool.New(5)
	encodingClient := dispmock.NewMockEncoderClientV2()
	chainReader := &coremock.MockWriter{}
	chainReader.On("GetCurrentBlockNumber").Return(blockNumber, nil)
	em, err := controller.NewEncodingManager(controller.EncodingManagerConfig{
		PullInterval:           1 * time.Second,
		EncodingRequestTimeout: 5 * time.Second,
		StoreTimeout:           5 * time.Second,
		NumEncodingRetries:     1,
		NumRelayAssignment:     2,
		AvailableRelays:        []corev2.RelayKey{0, 1, 2, 3},
	}, blobMetadataStore, pool, encodingClient, chainReader, logger)
	assert.NoError(t, err)
	return &testComponents{
		EncodingManager: em,
		Pool:            pool,
		EncodingClient:  encodingClient,
		ChainReader:     chainReader,
	}
}
