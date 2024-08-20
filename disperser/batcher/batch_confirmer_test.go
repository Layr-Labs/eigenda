package batcher_test

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	cmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser"
	bat "github.com/Layr-Labs/eigenda/disperser/batcher"
	batcherinmem "github.com/Layr-Labs/eigenda/disperser/batcher/inmem"
	batchermock "github.com/Layr-Labs/eigenda/disperser/batcher/mock"
	batmock "github.com/Layr-Labs/eigenda/disperser/batcher/mock"
	"github.com/Layr-Labs/eigenda/disperser/common/inmem"
	dmock "github.com/Layr-Labs/eigenda/disperser/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gammazero/workerpool"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	assignmentCoordinator                 = &core.StdAssignmentCoordinator{}
	encoderProver         encoding.Prover = nil
)

type batchConfirmerComponents struct {
	batchConfirmer   *bat.BatchConfirmer
	blobStore        disperser.BlobStore
	minibatchStore   bat.MinibatchStore
	minibatcher      *bat.Minibatcher
	dispatcher       *dmock.Dispatcher
	chainData        *coremock.ChainDataMock
	transactor       *coremock.MockTransactor
	txnManager       *batchermock.MockTxnManager
	encodingStreamer *bat.EncodingStreamer
	ethClient        *cmock.MockEthClient
}

func makeBatchConfirmer(t *testing.T) *batchConfirmerComponents {
	logger := logging.NewNoopLogger()
	// logger, err := common.NewLogger(common.DefaultLoggerConfig())
	// assert.NoError(t, err)
	transactor := &coremock.MockTransactor{}
	agg, err := core.NewStdSignatureAggregator(logger, transactor)
	assert.NoError(t, err)
	state := mockChainState.GetTotalOperatorState(context.Background(), 0)
	dispatcher := dmock.NewDispatcher(state)
	blobStore := inmem.NewBlobStore()
	ethClient := &cmock.MockEthClient{}
	txnManager := batmock.NewTxnManager()
	minibatchStore := batcherinmem.NewMinibatchStore(logger)
	encodingWorkerPool := workerpool.New(10)
	encoderProver, err = makeTestProver()
	assert.NoError(t, err)
	encoderClient := disperser.NewLocalEncoderClient(encoderProver)
	metrics := bat.NewMetrics("9100", logger)
	trigger := bat.NewEncodedSizeNotifier(
		make(chan struct{}, 1),
		10*1024*1024,
	)
	encodingStreamer, err := bat.NewEncodingStreamer(streamerConfig, blobStore, mockChainState, encoderClient, assignmentCoordinator, trigger, encodingWorkerPool, metrics.EncodingStreamerMetrics, logger)
	assert.NoError(t, err)
	pool := workerpool.New(int(10))
	minibatcher, err := bat.NewMinibatcher(bat.MinibatcherConfig{
		PullInterval:              100 * time.Millisecond,
		MaxNumConnections:         10,
		MaxNumRetriesPerBlob:      2,
		MaxNumRetriesPerDispersal: 1,
	}, blobStore, minibatchStore, dispatcher, mockChainState, assignmentCoordinator, encodingStreamer, ethClient, pool, logger)
	assert.NoError(t, err)

	config := bat.BatchConfirmerConfig{
		PullInterval:                 100 * time.Millisecond,
		DispersalTimeout:             1 * time.Second,
		DispersalStatusCheckInterval: 100 * time.Millisecond,
		AttestationTimeout:           1 * time.Second,
		NumConnections:               1,
		SRSOrder:                     3000,
		MaxNumRetriesPerBlob:         2,
	}
	b, err := bat.NewBatchConfirmer(config, blobStore, minibatchStore, dispatcher, mockChainState, assignmentCoordinator, encodingStreamer, agg, ethClient, transactor, txnManager, minibatcher, logger)
	assert.NoError(t, err)

	ethClient.On("BlockNumber").Return(uint64(initialBlock), nil)

	return &batchConfirmerComponents{
		batchConfirmer:   b,
		blobStore:        blobStore,
		minibatchStore:   minibatchStore,
		minibatcher:      minibatcher,
		dispatcher:       dispatcher,
		chainData:        mockChainState,
		transactor:       transactor,
		txnManager:       txnManager,
		encodingStreamer: b.EncodingStreamer,
		ethClient:        ethClient,
	}
}

func generateBlobAndHeader(t *testing.T, operatorState *core.OperatorState, securityParams []*core.SecurityParam) (*core.Blob, *core.BlobHeader) {
	assert.NotNil(t, operatorState)
	assert.Greater(t, len(securityParams), 0)
	assert.NotNil(t, assignmentCoordinator)
	assert.NotNil(t, encoderProver)

	blob := makeTestBlob(securityParams)
	blobLength := encoding.GetBlobLength(uint(len(blob.Data)))
	blobHeader := &core.BlobHeader{}
	blobQuorumInfo := &core.BlobQuorumInfo{}
	chunkLength := uint(0)
	var err error
	for _, sp := range securityParams {
		chunkLength, err = assignmentCoordinator.CalculateChunkLength(operatorState, blobLength, streamerConfig.TargetNumChunks, sp)
		assert.NoError(t, err)
		blobQuorumInfo = &core.BlobQuorumInfo{
			SecurityParam: *sp,
			ChunkLength:   chunkLength,
		}
		blobHeader.QuorumInfos = append(blobHeader.QuorumInfos, blobQuorumInfo)
	}

	// use the latest security param to encode the blob
	_, info, err := assignmentCoordinator.GetAssignments(operatorState, blobLength, blobQuorumInfo)
	assert.NoError(t, err)
	params := encoding.ParamsFromMins(chunkLength, info.TotalChunks)
	commits, _, err := encoderProver.EncodeAndProve(blob.Data, params)
	assert.NoError(t, err)
	blobHeader.BlobCommitments = commits
	return &blob, blobHeader
}

func TestBatchConfirmerIteration(t *testing.T) {
	components := makeBatchConfirmer(t)
	ctx := context.Background()
	operatorState := components.chainData.GetTotalOperatorState(ctx, initialBlock)
	blob1, blobHeader1 := generateBlobAndHeader(t, operatorState.OperatorState, []*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}})
	blob2, blobHeader2 := generateBlobAndHeader(t, operatorState.OperatorState, []*core.SecurityParam{{
		QuorumID:              1,
		AdversaryThreshold:    70,
		ConfirmationThreshold: 100,
	}})
	b := components.batchConfirmer
	batchID, err := uuid.NewV7()
	assert.NoError(t, err)
	batch := &bat.BatchRecord{
		ID:                   batchID,
		CreatedAt:            time.Now(),
		ReferenceBlockNumber: uint(initialBlock),
		Status:               bat.BatchStatusFormed,
		NumMinibatches:       2,
	}

	// Set up batch
	err = b.MinibatchStore.PutBatch(context.Background(), batch)
	assert.NoError(t, err)
	requestedAt1, blobKey1 := queueBlob(t, ctx, blob1, components.blobStore)
	_, blobKey2 := queueBlob(t, ctx, blob2, components.blobStore)
	meta1, err := components.blobStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	meta2, err := components.blobStore.GetBlobMetadata(ctx, blobKey2)
	assert.NoError(t, err)
	batchHeader1 := &core.BatchHeader{
		ReferenceBlockNumber: initialBlock,
		BatchRoot:            [32]byte{},
	}
	_, err = batchHeader1.SetBatchRoot([]*core.BlobHeader{blobHeader1})
	assert.NoError(t, err)
	batchHeaderHash1, err := batchHeader1.GetBatchHeaderHash()
	assert.NoError(t, err)
	batchHeader2 := &core.BatchHeader{
		ReferenceBlockNumber: initialBlock,
		BatchRoot:            [32]byte{},
	}
	_, err = batchHeader2.SetBatchRoot([]*core.BlobHeader{blobHeader2})
	assert.NoError(t, err)
	batchHeaderHash2, err := batchHeader2.GetBatchHeaderHash()
	assert.NoError(t, err)
	// Set up dispersals
	for opID, opInfo := range operatorState.PrivateOperators {
		req0 := &bat.MinibatchDispersal{
			BatchID:        batchID,
			MinibatchIndex: 0,
			OperatorID:     opID,
			Socket:         opInfo.DispersalPort,
			NumBlobs:       1,
			RequestedAt:    time.Now(),
			DispersalResponse: bat.DispersalResponse{
				Signatures:  []*core.Signature{opInfo.KeyPair.SignMessage(batchHeaderHash1)},
				RespondedAt: time.Now(),
				Error:       nil,
			},
		}
		err = b.MinibatchStore.PutDispersal(context.Background(), req0)
		assert.NoError(t, err)

		req1 := &bat.MinibatchDispersal{
			BatchID:        batchID,
			MinibatchIndex: 1,
			OperatorID:     opID,
			Socket:         opInfo.DispersalPort,
			NumBlobs:       1,
			RequestedAt:    time.Now(),
			DispersalResponse: bat.DispersalResponse{
				Signatures:  []*core.Signature{opInfo.KeyPair.SignMessage(batchHeaderHash2)},
				RespondedAt: time.Now(),
				Error:       nil,
			},
		}
		err = b.MinibatchStore.PutDispersal(context.Background(), req1)
		assert.NoError(t, err)
	}

	// Set up batch state
	b.Minibatcher.Batches[batchID] = &bat.BatchState{
		BatchID:              batchID,
		ReferenceBlockNumber: uint(initialBlock),
		BlobHeaders: []*core.BlobHeader{
			blobHeader1,
			blobHeader2,
		},
		BlobMetadata:  []*disperser.BlobMetadata{meta1, meta2},
		OperatorState: operatorState.IndexedOperatorState,
	}

	// Receive signatures
	signChan := make(chan core.SigningMessage, 4)
	batchHeaderHash := [32]byte{138, 1, 226, 93, 51, 120, 236, 124, 91, 206, 100, 187, 237, 1, 193, 151, 137, 131, 30, 218, 139, 24, 221, 105, 141, 253, 242, 13, 239, 199, 179, 42}
	for opID, opInfo := range operatorState.PrivateOperators {
		signChan <- core.SigningMessage{
			Signature:       opInfo.KeyPair.SignMessage(batchHeaderHash),
			Operator:        opID,
			BatchHeaderHash: batchHeaderHash,
			Err:             nil,
		}
	}
	components.dispatcher.On("AttestBatch", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(signChan, nil)
	txn := types.NewTransaction(0, gethcommon.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	components.transactor.On("OperatorIDToAddress").Return(gethcommon.Address{}, nil)
	components.transactor.On("BuildConfirmBatchTxn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(
		func(args mock.Arguments) {
			assert.Equal(t, args.Get(1).(*core.BatchHeader).ReferenceBlockNumber, uint(initialBlock))
			assert.NotNil(t, args.Get(1).(*core.BatchHeader).BatchRoot)
			assert.Equal(t, args.Get(2).(map[uint8]*core.QuorumResult)[0].PercentSigned, uint8(100))
			assert.Equal(t, args.Get(2).(map[uint8]*core.QuorumResult)[1].PercentSigned, uint8(100))

			aggSig := args[3].(*core.SignatureAggregation)
			assert.Empty(t, aggSig.NonSigners)
			assert.Len(t, aggSig.QuorumAggPubKeys, 2)
			assert.Contains(t, aggSig.QuorumAggPubKeys, core.QuorumID(0))
			assert.Contains(t, aggSig.QuorumAggPubKeys, core.QuorumID(1))
			assert.Equal(t, aggSig.QuorumResults, map[core.QuorumID]*core.QuorumResult{
				core.QuorumID(0): {
					QuorumID:      core.QuorumID(0),
					PercentSigned: uint8(100),
				},
				core.QuorumID(1): {
					QuorumID:      core.QuorumID(1),
					PercentSigned: uint8(100),
				},
			})
		},
	).Return(txn, nil)
	components.txnManager.On("ProcessTransaction").Return(nil)
	err = b.HandleSingleBatch(ctx)
	assert.NoError(t, err)

	// Validate batch confirmation
	assert.Greater(t, len(components.txnManager.Requests), 0)
	// logData should be encoding 3 and 0
	logData, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	txHash := gethcommon.HexToHash("0x1234")
	blockNumber := big.NewInt(123)
	err = b.ProcessConfirmedBatch(ctx, &bat.ReceiptOrErr{
		Receipt: &types.Receipt{
			Logs: []*types.Log{
				{
					Topics: []gethcommon.Hash{common.BatchConfirmedEventSigHash, gethcommon.HexToHash("1234")},
					Data:   logData,
				},
			},
			BlockNumber: blockNumber,
			TxHash:      txHash,
		},
		Err:      nil,
		Metadata: components.txnManager.Requests[len(components.txnManager.Requests)-1].Metadata,
	})
	assert.NoError(t, err)
	// Check that the blob was processed
	meta1, err = components.blobStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Equal(t, blobKey1, meta1.GetBlobKey())
	assert.Equal(t, requestedAt1, meta1.RequestMetadata.RequestedAt)
	assert.Equal(t, disperser.Confirmed, meta1.BlobStatus)
	assert.Equal(t, meta1.ConfirmationInfo.BatchID, uint32(3))
	assert.Equal(t, meta1.ConfirmationInfo.ConfirmationTxnHash, txHash)
	assert.Equal(t, meta1.ConfirmationInfo.ConfirmationBlockNumber, uint32(blockNumber.Int64()))

	meta2, err = components.blobStore.GetBlobMetadata(ctx, blobKey2)
	assert.NoError(t, err)
	assert.Equal(t, blobKey2, meta2.GetBlobKey())
	assert.Equal(t, disperser.Confirmed, meta2.BlobStatus)

	res, err := components.encodingStreamer.EncodedBlobstore.GetEncodingResult(meta1.GetBlobKey(), 0)
	assert.ErrorContains(t, err, "no such key")
	assert.Nil(t, res)
	res, err = components.encodingStreamer.EncodedBlobstore.GetEncodingResult(meta2.GetBlobKey(), 1)
	assert.ErrorContains(t, err, "no such key")
	assert.Nil(t, res)
	count, size := components.encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, 0, count)
	assert.Equal(t, uint64(0), size)

	batch, err = components.minibatchStore.GetBatch(context.Background(), batch.ID)
	assert.NoError(t, err)
	assert.Equal(t, batch.Status, bat.BatchStatusAttested)
	batches, err := components.minibatchStore.GetBatchesByStatus(context.Background(), bat.BatchStatusAttested)
	assert.NoError(t, err)
	assert.Equal(t, len(batches), 1)
	assert.Equal(t, batches[0].ID, batch.ID)
}

func TestBatchConfirmerIterationFailure(t *testing.T) {
	// If dispersal responses are not received for all dispersal requests, the batch should not be confirmed
	components := makeBatchConfirmer(t)
	b := components.batchConfirmer
	batchID, err := uuid.NewV7()
	assert.NoError(t, err)
	batch := &bat.BatchRecord{
		ID:                   batchID,
		CreatedAt:            time.Now(),
		ReferenceBlockNumber: uint(initialBlock),
		Status:               bat.BatchStatusFormed,
	}
	err = b.MinibatchStore.PutBatch(context.Background(), batch)
	assert.NoError(t, err)
	operatorState := components.chainData.GetTotalOperatorState(context.Background(), 0)

	for opID, opInfo := range operatorState.PrivateOperators {
		req0 := &bat.MinibatchDispersal{
			BatchID:        batchID,
			MinibatchIndex: 0,
			OperatorID:     opID,
			Socket:         opInfo.DispersalPort,
			NumBlobs:       2,
			RequestedAt:    time.Now(),
			DispersalResponse: bat.DispersalResponse{
				Signatures:  []*core.Signature{opInfo.KeyPair.SignMessage([32]byte{0})},
				RespondedAt: time.Now(),
				Error:       nil,
			},
		}
		err = b.MinibatchStore.PutDispersal(context.Background(), req0)
		assert.NoError(t, err)
		req1 := &bat.MinibatchDispersal{
			BatchID:           batchID,
			MinibatchIndex:    1,
			OperatorID:        opID,
			Socket:            opInfo.DispersalPort,
			NumBlobs:          2,
			RequestedAt:       time.Now(),
			DispersalResponse: bat.DispersalResponse{
				// Missing RespondedAt
			},
		}
		err = b.MinibatchStore.PutDispersal(context.Background(), req1)
		assert.NoError(t, err)
	}

	err = b.HandleSingleBatch(context.Background())
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestBatchConfirmerInsufficientSignatures(t *testing.T) {
	components := makeBatchConfirmer(t)
	b := components.batchConfirmer
	batchID, err := uuid.NewV7()
	assert.NoError(t, err)
	batch := &bat.BatchRecord{
		ID:                   batchID,
		CreatedAt:            time.Now(),
		ReferenceBlockNumber: uint(initialBlock),
		Status:               bat.BatchStatusFormed,
		NumMinibatches:       2,
	}
	ctx := context.Background()
	err = b.MinibatchStore.PutBatch(ctx, batch)
	assert.NoError(t, err)

	operatorState := components.chainData.GetTotalOperatorState(context.Background(), 0)
	blob1, blobHeader1 := generateBlobAndHeader(t, operatorState.OperatorState, []*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}})
	blob2, blobHeader2 := generateBlobAndHeader(t, operatorState.OperatorState, []*core.SecurityParam{
		{
			QuorumID:              0,
			AdversaryThreshold:    70,
			ConfirmationThreshold: 100,
		},
		{
			QuorumID:              1,
			AdversaryThreshold:    70,
			ConfirmationThreshold: 100,
		}})
	_, blobKey1 := queueBlob(t, ctx, blob1, components.blobStore)
	_, blobKey2 := queueBlob(t, ctx, blob2, components.blobStore)
	meta1, err := components.blobStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	meta2, err := components.blobStore.GetBlobMetadata(ctx, blobKey2)
	assert.NoError(t, err)

	for opID, opInfo := range operatorState.PrivateOperators {
		req0 := &bat.MinibatchDispersal{
			BatchID:        batchID,
			MinibatchIndex: 0,
			OperatorID:     opID,
			Socket:         opInfo.DispersalPort,
			NumBlobs:       2,
			RequestedAt:    time.Now(),
			DispersalResponse: bat.DispersalResponse{
				Signatures:  []*core.Signature{opInfo.KeyPair.SignMessage([32]byte{0})},
				RespondedAt: time.Now(),
				Error:       nil,
			},
		}
		err = b.MinibatchStore.PutDispersal(context.Background(), req0)
		assert.NoError(t, err)
		req1 := &bat.MinibatchDispersal{
			BatchID:        batchID,
			MinibatchIndex: 1,
			OperatorID:     opID,
			Socket:         opInfo.DispersalPort,
			NumBlobs:       2,
			RequestedAt:    time.Now(),
			DispersalResponse: bat.DispersalResponse{
				Signatures:  []*core.Signature{opInfo.KeyPair.SignMessage([32]byte{1})},
				RespondedAt: time.Now(),
				Error:       nil,
			},
		}
		err = b.MinibatchStore.PutDispersal(context.Background(), req1)
		assert.NoError(t, err)
	}

	b.Minibatcher.Batches[batchID] = &bat.BatchState{
		BatchID:              batchID,
		ReferenceBlockNumber: uint(initialBlock),
		BlobHeaders: []*core.BlobHeader{
			blobHeader1,
			blobHeader2,
		},
		BlobMetadata:  []*disperser.BlobMetadata{meta1, meta2},
		OperatorState: operatorState.IndexedOperatorState,
	}

	signChan := make(chan core.SigningMessage, 4)
	batchHeader := &core.BatchHeader{
		ReferenceBlockNumber: uint(initialBlock),
		BatchRoot:            [32]byte{},
	}
	bhh1, err := blobHeader1.GetBlobHeaderHash()
	assert.NoError(t, err)
	bhh2, err := blobHeader2.GetBlobHeaderHash()
	assert.NoError(t, err)
	_, err = batchHeader.SetBatchRootFromBlobHeaderHashes([][32]byte{bhh1, bhh2})
	assert.NoError(t, err)
	batchHeaderHash, err := batchHeader.GetBatchHeaderHash()
	assert.NoError(t, err)
	for opID, opInfo := range operatorState.PrivateOperators {
		if opID == opId0 {
			signChan <- core.SigningMessage{
				Signature:       opInfo.KeyPair.SignMessage(batchHeaderHash),
				Operator:        opID,
				BatchHeaderHash: batchHeaderHash,
				Err:             nil,
			}
		} else {
			signChan <- core.SigningMessage{
				Signature:       nil,
				Operator:        opID,
				BatchHeaderHash: batchHeaderHash,
				Err:             context.DeadlineExceeded,
			}
		}
	}
	components.dispatcher.On("AttestBatch", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(signChan, nil)
	components.transactor.On("OperatorIDToAddress").Return(gethcommon.Address{}, nil)
	err = b.HandleSingleBatch(context.Background())
	assert.ErrorContains(t, err, "no blobs received sufficient signatures")
	batch, err = components.minibatchStore.GetBatch(context.Background(), batch.ID)
	assert.NoError(t, err)
	assert.Equal(t, batch.Status, bat.BatchStatusFailed)
}
