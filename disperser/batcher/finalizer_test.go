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
	"github.com/Layr-Labs/eigenda/disperser/inmem"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"

	m "github.com/stretchr/testify/mock"
)

const timeout = 5 * time.Second
const loopInterval = 6 * time.Minute

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

	finalizer := batcher.NewFinalizer(timeout, loopInterval, queue, ethClient, rpcClient, logger)

	requestedAt := uint64(time.Now().UnixNano())
	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:           0,
		AdversaryThreshold: 80,
	}})
	ctx := context.Background()
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
		BlobCommitment:          &core.BlobCommitments{},
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
			RequestedAt: requestedAt,
		},
	}
	m, err := queue.MarkBlobConfirmed(ctx, metadata, confirmationInfo)
	assert.Equal(t, disperser.Confirmed, m.BlobStatus)
	assert.NoError(t, err)

	err = finalizer.FinalizeBlobs(context.Background())
	assert.NoError(t, err)

	metadatas, err := queue.GetBlobMetadataByStatus(ctx, disperser.Confirmed)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 0)

	metadatas, err = queue.GetBlobMetadataByStatus(ctx, disperser.Finalized)
	assert.NoError(t, err)
	assert.Len(t, metadatas, 1)

	assert.Equal(t, metadatas[0].BlobHash, metadataKey.BlobHash)
	assert.Equal(t, metadatas[0].BlobStatus, disperser.Finalized)
	assert.Equal(t, metadatas[0].RequestMetadata.RequestedAt, requestedAt)
	assert.Equal(t, metadatas[0].RequestMetadata.SecurityParams, blob.RequestHeader.SecurityParams)
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

	finalizer := batcher.NewFinalizer(timeout, loopInterval, queue, ethClient, rpcClient, logger)

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
		BlobCommitment:          &core.BlobCommitments{},
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
