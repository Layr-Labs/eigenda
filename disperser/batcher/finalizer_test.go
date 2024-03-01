package batcher_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/common/inmem"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"

	m "github.com/stretchr/testify/mock"
)

const timeout = 5 * time.Second
const loopInterval = 6 * time.Minute

type MockBlobStore struct {
	m.Mock
}

func TestFinalizedBlob(t *testing.T) {
	queue := inmem.NewBlobStore()
	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	assert.NoError(t, err)
	ethClient := &mock.MockEthClient{}
	rpcClient := &mock.MockRPCEthClient{}

	latestFinalBlock := int64(1_000_010)
	rpcClient.On("CallContext", m.Anything, m.Anything, "eth_getBlockByNumber", "finalized", false).
		Run(func(args m.Arguments) {
			args[1].(*types.Header).Number = big.NewInt(latestFinalBlock)
		}).Return(nil).Once()
	ethClient.On("TransactionReceipt", m.Anything, m.Anything).Return(&types.Receipt{
		BlockNumber: new(big.Int).SetUint64(1_000_000),
	}, nil)

	metrics := batcher.NewMetrics("9100", logger)
	finalizer := batcher.NewFinalizer(timeout, loopInterval, queue, ethClient, rpcClient, 1, 1, 1, logger, metrics.FinalizerMetrics)

	requestedAt := uint64(time.Now().UnixNano())
	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:           0,
		AdversaryThreshold: 80,
	}})
	ctx := context.Background()
	metadataKey1, err := queue.StoreBlob(ctx, &blob, requestedAt)
	assert.NoError(t, err)
	metadataKey2, err := queue.StoreBlob(ctx, &blob, requestedAt+1)
	assert.NoError(t, err)
	batchHeaderHash := [32]byte{1, 2, 3}
	blobIndex := uint32(10)
	sigRecordHash := [32]byte{0}
	inclusionProof := []byte{1, 2, 3, 4, 5}
	confirmationInfo := &disperser.ConfirmationInfo{
		BatchHeaderHash:         batchHeaderHash,
		BlobIndex:               blobIndex,
		SignatoryRecordHash:     sigRecordHash,
		ReferenceBlockNumber:    132,
		BatchRoot:               []byte("hello"),
		BlobInclusionProof:      inclusionProof,
		BlobCommitment:          &encoding.BlobCommitments{},
		BatchID:                 99,
		ConfirmationTxnHash:     common.HexToHash("0x123"),
		ConfirmationBlockNumber: uint32(150),
		Fee:                     []byte{0},
	}
	metadata1 := &disperser.BlobMetadata{
		BlobHash:     metadataKey1.BlobHash,
		MetadataHash: metadataKey1.MetadataHash,
		BlobStatus:   disperser.Processing,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: blob.RequestHeader.SecurityParams,
			},
			RequestedAt: requestedAt,
		},
	}
	metadata2 := &disperser.BlobMetadata{
		BlobHash:     metadataKey2.BlobHash,
		MetadataHash: metadataKey2.MetadataHash,
		BlobStatus:   disperser.Processing,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: blob.RequestHeader.SecurityParams,
			},
			RequestedAt: requestedAt + 1,
		},
	}
	m, err := queue.MarkBlobConfirmed(ctx, metadata1, confirmationInfo)
	assert.Equal(t, disperser.Confirmed, m.BlobStatus)
	assert.NoError(t, err)
	m, err = queue.MarkBlobConfirmed(ctx, metadata2, confirmationInfo)
	assert.Equal(t, disperser.Confirmed, m.BlobStatus)
	assert.NoError(t, err)

	err = finalizer.FinalizeBlobs(context.Background())
	assert.NoError(t, err)

	metadatas, err := queue.GetBlobMetadataByStatus(ctx, disperser.Confirmed)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 0)

	metadatas, err = queue.GetBlobMetadataByStatus(ctx, disperser.Finalized)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 2)

	assert.ElementsMatch(t, []string{metadatas[0].BlobHash, metadatas[1].BlobHash}, []string{metadataKey1.BlobHash, metadataKey2.BlobHash})
	assert.Equal(t, metadatas[0].BlobStatus, disperser.Finalized)
	assert.Equal(t, metadatas[1].BlobStatus, disperser.Finalized)
	assert.ElementsMatch(t, []uint64{metadatas[0].RequestMetadata.RequestedAt, metadatas[1].RequestMetadata.RequestedAt}, []uint64{requestedAt, requestedAt + 1})
	assert.Equal(t, metadatas[0].RequestMetadata.SecurityParams, blob.RequestHeader.SecurityParams)
	assert.Equal(t, metadatas[1].RequestMetadata.SecurityParams, blob.RequestHeader.SecurityParams)
}

func TestUnfinalizedBlob(t *testing.T) {
	ctx := context.Background()
	queue := inmem.NewBlobStore()
	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	assert.NoError(t, err)
	ethClient := &mock.MockEthClient{}
	rpcClient := &mock.MockRPCEthClient{}

	latestFinalBlock := int64(1_000_010)
	rpcClient.On("CallContext", m.Anything, m.Anything, "eth_getBlockByNumber", "finalized", false).
		Run(func(args m.Arguments) {
			args[1].(*types.Header).Number = big.NewInt(latestFinalBlock)
		}).Return(nil).Once()
	ethClient.On("TransactionReceipt", m.Anything, m.Anything).Return(&types.Receipt{
		BlockNumber: new(big.Int).SetUint64(1_000_100),
	}, nil)

	metrics := batcher.NewMetrics("9100", logger)
	finalizer := batcher.NewFinalizer(timeout, loopInterval, queue, ethClient, rpcClient, 1, 1, 1, logger, metrics.FinalizerMetrics)

	requestedAt := uint64(time.Now().UnixNano())
	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:           0,
		AdversaryThreshold: 80,
	}})
	metadataKey, err := queue.StoreBlob(ctx, &blob, requestedAt)
	assert.NoError(t, err)
	batchHeaderHash := [32]byte{1, 2, 3}
	blobIndex := uint32(10)
	sigRecordHash := [32]byte{0}
	inclusionProof := []byte{1, 2, 3, 4, 5}
	confirmationInfo := &disperser.ConfirmationInfo{
		BatchHeaderHash:         batchHeaderHash,
		BlobIndex:               blobIndex,
		SignatoryRecordHash:     sigRecordHash,
		ReferenceBlockNumber:    132,
		BatchRoot:               []byte("hello"),
		BlobInclusionProof:      inclusionProof,
		BlobCommitment:          &encoding.BlobCommitments{},
		BatchID:                 99,
		ConfirmationTxnHash:     common.HexToHash("0x123"),
		ConfirmationBlockNumber: uint32(150),
		Fee:                     []byte{0},
	}
	metadata := &disperser.BlobMetadata{
		BlobHash:     metadataKey.BlobHash,
		MetadataHash: metadataKey.MetadataHash,
		BlobStatus:   disperser.Processing,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: blob.RequestHeader.SecurityParams,
			},
			BlobSize:    uint(len(blob.Data)),
			RequestedAt: requestedAt,
		},
	}
	m, err := queue.MarkBlobConfirmed(ctx, metadata, confirmationInfo)
	assert.NoError(t, err)
	assert.Equal(t, disperser.Confirmed, m.BlobStatus)
	err = finalizer.FinalizeBlobs(context.Background())
	assert.NoError(t, err)

	metadatas, err := queue.GetBlobMetadataByStatus(ctx, disperser.Confirmed)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 1)

	metadatas, err = queue.GetBlobMetadataByStatus(ctx, disperser.Finalized)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 0)
}

func TestNoReceipt(t *testing.T) {
	ctx := context.Background()
	queue := inmem.NewBlobStore()
	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	assert.NoError(t, err)
	ethClient := &mock.MockEthClient{}
	rpcClient := &mock.MockRPCEthClient{}

	latestFinalBlock := int64(1_000_010)
	rpcClient.On("CallContext", m.Anything, m.Anything, "eth_getBlockByNumber", "finalized", false).
		Run(func(args m.Arguments) {
			args[1].(*types.Header).Number = big.NewInt(latestFinalBlock)
		}).Return(nil)
	ethClient.On("TransactionReceipt", m.Anything, m.Anything).Return(nil, ethereum.NotFound)

	metrics := batcher.NewMetrics("9100", logger)
	finalizer := batcher.NewFinalizer(timeout, loopInterval, queue, ethClient, rpcClient, 1, 1, 1, logger, metrics.FinalizerMetrics)

	requestedAt := uint64(time.Now().UnixNano())
	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:           0,
		AdversaryThreshold: 80,
	}})
	metadataKey, err := queue.StoreBlob(ctx, &blob, requestedAt)
	assert.NoError(t, err)
	batchHeaderHash := [32]byte{1, 2, 3}
	blobIndex := uint32(10)
	sigRecordHash := [32]byte{0}
	inclusionProof := []byte{1, 2, 3, 4, 5}
	confirmationInfo := &disperser.ConfirmationInfo{
		BatchHeaderHash:         batchHeaderHash,
		BlobIndex:               blobIndex,
		SignatoryRecordHash:     sigRecordHash,
		ReferenceBlockNumber:    132,
		BatchRoot:               []byte("hello"),
		BlobInclusionProof:      inclusionProof,
		BlobCommitment:          &encoding.BlobCommitments{},
		BatchID:                 99,
		ConfirmationTxnHash:     common.HexToHash("0x123"),
		ConfirmationBlockNumber: uint32(150),
		Fee:                     []byte{0},
	}
	metadata := &disperser.BlobMetadata{
		BlobHash:     metadataKey.BlobHash,
		MetadataHash: metadataKey.MetadataHash,
		BlobStatus:   disperser.Processing,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: blob.RequestHeader.SecurityParams,
			},
			BlobSize:    uint(len(blob.Data)),
			RequestedAt: requestedAt,
		},
	}
	m, err := queue.MarkBlobConfirmed(ctx, metadata, confirmationInfo)
	assert.NoError(t, err)
	assert.Equal(t, disperser.Confirmed, m.BlobStatus)
	err = finalizer.FinalizeBlobs(context.Background())
	assert.NoError(t, err)

	// status should be kept at confirmed
	metadatas, err := queue.GetBlobMetadataByStatus(ctx, disperser.Finalized)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 0)
	metadatas, err = queue.GetBlobMetadataByStatus(ctx, disperser.Failed)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 0)
	metadatas, err = queue.GetBlobMetadataByStatus(ctx, disperser.Confirmed)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 1)
	// num retries should be incremented
	assert.Equal(t, metadatas[0].NumRetries, uint(1))

	// try again
	err = finalizer.FinalizeBlobs(context.Background())
	assert.NoError(t, err)

	// status should be transitioned to failed
	metadatas, err = queue.GetBlobMetadataByStatus(ctx, disperser.Finalized)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 0)
	metadatas, err = queue.GetBlobMetadataByStatus(ctx, disperser.Confirmed)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 0)
	metadatas, err = queue.GetBlobMetadataByStatus(ctx, disperser.Failed)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 1)
	// num retries should be the same
	assert.Equal(t, metadatas[0].NumRetries, uint(1))
}

func TestFinalizeBlobs_NilMetadatas(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockBlobStore := new(MockBlobStore)
	// Mock setup to return nil metadatas and simulate a nil reference scenario
	var nilMetadatas []*disperser.BlobMetadata
	mockBlobStore.On(
		"GetBlobMetadataByStatusWithPagination",
		ctx,
		disperser.Confirmed,
		m.AnythingOfType("int32"),
		m.AnythingOfType("*disperser.BlobStoreExclusiveStartKey"),
	).Return(nilMetadatas, (*disperser.BlobStoreExclusiveStartKey)(nil), nil)

	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	assert.NoError(t, err)
	ethClient := &mock.MockEthClient{}
	rpcClient := &mock.MockRPCEthClient{}

	latestFinalBlock := int64(1_000_010)
	rpcClient.On("CallContext", m.Anything, m.Anything, "eth_getBlockByNumber", "finalized", false).
		Run(func(args m.Arguments) {
			args[1].(*types.Header).Number = big.NewInt(latestFinalBlock)
		}).Return(nil).Once()
	ethClient.On("TransactionReceipt", m.Anything, m.Anything).Return(&types.Receipt{
		BlockNumber: new(big.Int).SetUint64(1_000_100),
	}, nil)

	metrics := batcher.NewMetrics("9100", logger)

	finalizer := batcher.NewFinalizer(timeout, loopInterval, mockBlobStore, ethClient, rpcClient, 1, 1, 1, logger, metrics.FinalizerMetrics)

	// Act
	err = finalizer.FinalizeBlobs(ctx)

	// Assert
	assert.NoError(t, err, "FinalizeBlobs should not return an error even if metadatas is nil")
	mockBlobStore.AssertExpectations(t)
}

func TestFinalizeBlobsWithMockBlobStore(t *testing.T) {
	// Setup
	ctx := context.Background()
	securityParams := []*core.SecurityParam{{
		QuorumID:           1,
		AdversaryThreshold: 80,
		QuorumRate:         32000,
	}}

	blob := &core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: securityParams,
		},
		Data: []byte("test"),
	}

	blobHash := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	blobSize := uint(len(blob.Data))

	mockBlobStore := new(MockBlobStore)
	blobKey1 := disperser.BlobKey{
		BlobHash:     blobHash,
		MetadataHash: "hash",
	}

	metadata1 := &disperser.BlobMetadata{
		MetadataHash: blobKey1.MetadataHash,
		BlobHash:     blobHash,
		BlobStatus:   disperser.Confirmed,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          blobSize,
			RequestedAt:       123,
		},
	}

	// Mock setup to return nil metadatas and simulate a nil reference scenario
	var metadatas []*disperser.BlobMetadata = []*disperser.BlobMetadata{metadata1}
	mockBlobStore.On(
		"GetBlobMetadataByStatusWithPagination",
		ctx,
		disperser.Confirmed,
		m.AnythingOfType("int32"),
		m.AnythingOfType("*disperser.BlobStoreExclusiveStartKey"),
	).Return(metadatas, (*disperser.BlobStoreExclusiveStartKey)(nil), nil)

	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	assert.NoError(t, err)
	ethClient := &mock.MockEthClient{}
	rpcClient := &mock.MockRPCEthClient{}

	latestFinalBlock := int64(1_000_010)
	rpcClient.On("CallContext", m.Anything, m.Anything, "eth_getBlockByNumber", "finalized", false).
		Run(func(args m.Arguments) {
			args[1].(*types.Header).Number = big.NewInt(latestFinalBlock)
		}).Return(nil).Once()
	ethClient.On("TransactionReceipt", m.Anything, m.Anything).Return(&types.Receipt{
		BlockNumber: new(big.Int).SetUint64(1_000_100),
	}, nil)

	metrics := batcher.NewMetrics("9100", logger)

	finalizer := batcher.NewFinalizer(timeout, loopInterval, mockBlobStore, ethClient, rpcClient, 1, 1, 1, logger, metrics.FinalizerMetrics)

	// Act
	err = finalizer.FinalizeBlobs(ctx)

	// Assert
	assert.NoError(t, err, "FinalizeBlobs should not return an error")
	mockBlobStore.AssertExpectations(t)
}

func (m *MockBlobStore) StoreBlob(ctx context.Context, blob *core.Blob, requestedAt uint64) (disperser.BlobKey, error) {
	args := m.Called(ctx, blob, requestedAt)
	return args.Get(0).(disperser.BlobKey), args.Error(1)
}

func (m *MockBlobStore) GetBlobContent(ctx context.Context, blobHash disperser.BlobHash) ([]byte, error) {
	args := m.Called(ctx, blobHash)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockBlobStore) MarkBlobConfirmed(ctx context.Context, existingMetadata *disperser.BlobMetadata, confirmationInfo *disperser.ConfirmationInfo) (*disperser.BlobMetadata, error) {
	args := m.Called(ctx, existingMetadata, confirmationInfo)
	return args.Get(0).(*disperser.BlobMetadata), args.Error(1)
}

func (m *MockBlobStore) MarkBlobInsufficientSignatures(ctx context.Context, existingMetadata *disperser.BlobMetadata, confirmationInfo *disperser.ConfirmationInfo) (*disperser.BlobMetadata, error) {
	args := m.Called(ctx, existingMetadata, confirmationInfo)
	return args.Get(0).(*disperser.BlobMetadata), args.Error(1)
}

func (m *MockBlobStore) MarkBlobFinalized(ctx context.Context, blobKey disperser.BlobKey) error {
	args := m.Called(ctx, blobKey)
	return args.Error(0)
}

func (m *MockBlobStore) MarkBlobProcessing(ctx context.Context, blobKey disperser.BlobKey) error {
	args := m.Called(ctx, blobKey)
	return args.Error(0)
}

func (m *MockBlobStore) MarkBlobFailed(ctx context.Context, blobKey disperser.BlobKey) error {
	args := m.Called(ctx, blobKey)
	return args.Error(0)
}

func (m *MockBlobStore) IncrementBlobRetryCount(ctx context.Context, existingMetadata *disperser.BlobMetadata) error {
	args := m.Called(ctx, existingMetadata)
	return args.Error(0)
}

func (m *MockBlobStore) GetBlobsByMetadata(ctx context.Context, metadata []*disperser.BlobMetadata) (map[disperser.BlobKey]*core.Blob, error) {
	args := m.Called(ctx, metadata)
	return args.Get(0).(map[disperser.BlobKey]*core.Blob), args.Error(1)
}

func (m *MockBlobStore) GetBlobMetadataByStatus(ctx context.Context, blobStatus disperser.BlobStatus) ([]*disperser.BlobMetadata, error) {
	args := m.Called(ctx, blobStatus)
	return args.Get(0).([]*disperser.BlobMetadata), args.Error(1)
}

func (m *MockBlobStore) GetMetadataInBatch(ctx context.Context, batchHeaderHash [32]byte, blobIndex uint32) (*disperser.BlobMetadata, error) {
	args := m.Called(ctx, batchHeaderHash, blobIndex)
	return args.Get(0).(*disperser.BlobMetadata), args.Error(1)
}

func (m *MockBlobStore) GetBlobMetadataByStatusWithPagination(ctx context.Context, blobStatus disperser.BlobStatus, limit int32, exclusiveStartKey *disperser.BlobStoreExclusiveStartKey) ([]*disperser.BlobMetadata, *disperser.BlobStoreExclusiveStartKey, error) {
	args := m.Called(ctx, blobStatus, limit, exclusiveStartKey)
	return args.Get(0).([]*disperser.BlobMetadata), args.Get(1).(*disperser.BlobStoreExclusiveStartKey), args.Error(2)
}

func (m *MockBlobStore) GetAllBlobMetadataByBatch(ctx context.Context, batchHeaderHash [32]byte) ([]*disperser.BlobMetadata, error) {
	args := m.Called(ctx, batchHeaderHash)
	return args.Get(0).([]*disperser.BlobMetadata), args.Error(1)
}

func (m *MockBlobStore) GetBlobMetadata(ctx context.Context, blobKey disperser.BlobKey) (*disperser.BlobMetadata, error) {
	args := m.Called(ctx, blobKey)
	return args.Get(0).(*disperser.BlobMetadata), args.Error(1)
}

func (m *MockBlobStore) HandleBlobFailure(ctx context.Context, metadata *disperser.BlobMetadata, maxRetry uint) error {
	args := m.Called(ctx, metadata, maxRetry)
	return args.Error(0)
}
