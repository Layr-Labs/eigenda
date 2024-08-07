package batcher_test

import (
	"context"
	"math/big"
	"testing"
	"time"

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
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gammazero/workerpool"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	asgn := &core.StdAssignmentCoordinator{}
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
	p, err := makeTestProver()
	assert.NoError(t, err)
	encoderClient := disperser.NewLocalEncoderClient(p)
	metrics := bat.NewMetrics("9100", logger)
	trigger := bat.NewEncodedSizeNotifier(
		make(chan struct{}, 1),
		10*1024*1024,
	)
	encodingStreamer, err := bat.NewEncodingStreamer(streamerConfig, blobStore, mockChainState, encoderClient, asgn, trigger, encodingWorkerPool, metrics.EncodingStreamerMetrics, logger)
	assert.NoError(t, err)
	pool := workerpool.New(int(10))
	minibatcher, err := bat.NewMinibatcher(bat.MinibatcherConfig{
		PullInterval:              100 * time.Millisecond,
		MaxNumConnections:         10,
		MaxNumRetriesPerBlob:      2,
		MaxNumRetriesPerDispersal: 1,
	}, blobStore, minibatchStore, dispatcher, mockChainState, asgn, encodingStreamer, ethClient, pool, logger)
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
	b, err := bat.NewBatchConfirmer(config, blobStore, minibatchStore, dispatcher, mockChainState, asgn, encodingStreamer, agg, ethClient, transactor, txnManager, minibatcher, logger)
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

func TestBatchConfirmerIteration(t *testing.T) {
	components := makeBatchConfirmer(t)
	b := components.batchConfirmer
	batchID, err := uuid.NewV7()
	assert.NoError(t, err)
	batch := &bat.BatchRecord{
		ID:                   batchID,
		CreatedAt:            time.Now(),
		ReferenceBlockNumber: uint(initialBlock),
		Status:               bat.BatchStatusFormed,
		HeaderHash:           [32]byte{},
		AggregatePubKey:      &core.G2Point{},
		AggregateSignature:   &core.Signature{},
	}
	err = b.MinibatchStore.PutBatch(context.Background(), batch)
	assert.NoError(t, err)
	err = b.MinibatchStore.PutMinibatch(context.Background(), &bat.MinibatchRecord{
		BatchID:              batchID,
		MinibatchIndex:       0,
		BlobHeaderHashes:     [][32]byte{{1}, {2}},
		BatchSize:            0,
		ReferenceBlockNumber: uint(initialBlock),
	})
	assert.NoError(t, err)
	err = b.MinibatchStore.PutMinibatch(context.Background(), &bat.MinibatchRecord{
		BatchID:              batchID,
		MinibatchIndex:       1,
		BlobHeaderHashes:     [][32]byte{{3}, {4}},
		BatchSize:            0,
		ReferenceBlockNumber: uint(initialBlock),
	})
	assert.NoError(t, err)
	operatorState := components.chainData.GetTotalOperatorState(context.Background(), 0)

	for opID, opInfo := range operatorState.PrivateOperators {
		req0 := &bat.DispersalRequest{
			BatchID:        batchID,
			MinibatchIndex: 0,
			OperatorID:     opID,
			Socket:         opInfo.DispersalPort,
			NumBlobs:       2,
			RequestedAt:    time.Now(),
			BlobHash:       "0",
			MetadataHash:   "0",
		}
		err = b.MinibatchStore.PutDispersalRequest(context.Background(), req0)
		assert.NoError(t, err)
		err = b.MinibatchStore.PutDispersalResponse(context.Background(), &bat.DispersalResponse{
			DispersalRequest: *req0,
			Signatures:       []*core.Signature{opInfo.KeyPair.SignMessage([32]byte{0})},
			RespondedAt:      time.Now(),
			Error:            nil,
		})
		assert.NoError(t, err)

		req1 := &bat.DispersalRequest{
			BatchID:        batchID,
			MinibatchIndex: 1,
			OperatorID:     opID,
			Socket:         opInfo.DispersalPort,
			NumBlobs:       2,
			RequestedAt:    time.Now(),
			BlobHash:       "1",
			MetadataHash:   "1",
		}
		err = b.MinibatchStore.PutDispersalRequest(context.Background(), req1)
		assert.NoError(t, err)
		err = b.MinibatchStore.PutDispersalResponse(context.Background(), &bat.DispersalResponse{
			DispersalRequest: *req1,
			Signatures:       []*core.Signature{opInfo.KeyPair.SignMessage([32]byte{1})},
			RespondedAt:      time.Now(),
			Error:            nil,
		})
		assert.NoError(t, err)
	}

	b.Minibatcher.Batches[batchID] = &bat.BatchState{
		BatchID:              batchID,
		ReferenceBlockNumber: uint(initialBlock),
		BlobHeaders: []*core.BlobHeader{
			{
				AccountID: "0",
				QuorumInfos: []*core.BlobQuorumInfo{
					{
						SecurityParam: core.SecurityParam{
							QuorumID:              0,
							AdversaryThreshold:    30,
							ConfirmationThreshold: 80,
						},
					},
					{
						SecurityParam: core.SecurityParam{
							QuorumID:              1,
							AdversaryThreshold:    30,
							ConfirmationThreshold: 80,
						},
					},
				},
			},
			{
				AccountID: "1",
				QuorumInfos: []*core.BlobQuorumInfo{
					{
						SecurityParam: core.SecurityParam{
							QuorumID:              0,
							AdversaryThreshold:    30,
							ConfirmationThreshold: 80,
						},
					},
					{
						SecurityParam: core.SecurityParam{
							QuorumID:              1,
							AdversaryThreshold:    30,
							ConfirmationThreshold: 80,
						},
					},
				},
			},
		},
		BlobMetadata:   []*disperser.BlobMetadata{},
		OperatorState:  operatorState.IndexedOperatorState,
		NumMinibatches: 2,
	}

	signChan := make(chan core.SigningMessage, 4)
	batchHeaderHash := [32]byte{93, 156, 41, 17, 3, 78, 159, 243, 222, 111, 54, 107, 237, 48, 243, 176, 224, 151, 96, 151, 159, 99, 118, 186, 53, 192, 72, 59, 160, 73, 7, 213}
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
		},
	).Return(txn, nil)
	components.txnManager.On("ProcessTransaction").Return(nil)
	err = b.HandleSingleBatch(context.Background())
	assert.NoError(t, err)
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
		HeaderHash:           [32]byte{},
		AggregatePubKey:      &core.G2Point{},
		AggregateSignature:   &core.Signature{},
	}
	err = b.MinibatchStore.PutBatch(context.Background(), batch)
	assert.NoError(t, err)
	err = b.MinibatchStore.PutMinibatch(context.Background(), &bat.MinibatchRecord{
		BatchID:              batchID,
		MinibatchIndex:       0,
		BlobHeaderHashes:     [][32]byte{{1}, {2}},
		BatchSize:            0,
		ReferenceBlockNumber: uint(initialBlock),
	})
	assert.NoError(t, err)
	err = b.MinibatchStore.PutMinibatch(context.Background(), &bat.MinibatchRecord{
		BatchID:              batchID,
		MinibatchIndex:       1,
		BlobHeaderHashes:     [][32]byte{{3}, {4}},
		BatchSize:            0,
		ReferenceBlockNumber: uint(initialBlock),
	})
	assert.NoError(t, err)
	operatorState := components.chainData.GetTotalOperatorState(context.Background(), 0)

	for opID, opInfo := range operatorState.PrivateOperators {
		req0 := &bat.DispersalRequest{
			BatchID:        batchID,
			MinibatchIndex: 0,
			OperatorID:     opID,
			Socket:         opInfo.DispersalPort,
			NumBlobs:       2,
			RequestedAt:    time.Now(),
			BlobHash:       "0",
			MetadataHash:   "0",
		}
		err = b.MinibatchStore.PutDispersalRequest(context.Background(), req0)
		assert.NoError(t, err)
		err = b.MinibatchStore.PutDispersalResponse(context.Background(), &bat.DispersalResponse{
			DispersalRequest: *req0,
			Signatures:       []*core.Signature{opInfo.KeyPair.SignMessage([32]byte{0})},
			RespondedAt:      time.Now(),
			Error:            nil,
		})
		assert.NoError(t, err)

		req1 := &bat.DispersalRequest{
			BatchID:        batchID,
			MinibatchIndex: 1,
			OperatorID:     opID,
			Socket:         opInfo.DispersalPort,
			NumBlobs:       2,
			RequestedAt:    time.Now(),
			BlobHash:       "1",
			MetadataHash:   "1",
		}
		err = b.MinibatchStore.PutDispersalRequest(context.Background(), req1)
		assert.NoError(t, err)
		// Missing RespondedAt
		err = b.MinibatchStore.PutDispersalResponse(context.Background(), &bat.DispersalResponse{
			DispersalRequest: *req1,
		})
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
		HeaderHash:           [32]byte{},
		AggregatePubKey:      &core.G2Point{},
		AggregateSignature:   &core.Signature{},
	}
	err = b.MinibatchStore.PutBatch(context.Background(), batch)
	assert.NoError(t, err)
	err = b.MinibatchStore.PutMinibatch(context.Background(), &bat.MinibatchRecord{
		BatchID:              batchID,
		MinibatchIndex:       0,
		BlobHeaderHashes:     [][32]byte{{1}, {2}},
		BatchSize:            0,
		ReferenceBlockNumber: uint(initialBlock),
	})
	assert.NoError(t, err)
	err = b.MinibatchStore.PutMinibatch(context.Background(), &bat.MinibatchRecord{
		BatchID:              batchID,
		MinibatchIndex:       1,
		BlobHeaderHashes:     [][32]byte{{3}, {4}},
		BatchSize:            0,
		ReferenceBlockNumber: uint(initialBlock),
	})
	assert.NoError(t, err)
	operatorState := components.chainData.GetTotalOperatorState(context.Background(), 0)

	for opID, opInfo := range operatorState.PrivateOperators {
		req0 := &bat.DispersalRequest{
			BatchID:        batchID,
			MinibatchIndex: 0,
			OperatorID:     opID,
			Socket:         opInfo.DispersalPort,
			NumBlobs:       2,
			RequestedAt:    time.Now(),
			BlobHash:       "0",
			MetadataHash:   "0",
		}
		err = b.MinibatchStore.PutDispersalRequest(context.Background(), req0)
		assert.NoError(t, err)
		err = b.MinibatchStore.PutDispersalResponse(context.Background(), &bat.DispersalResponse{
			DispersalRequest: *req0,
			Signatures:       []*core.Signature{opInfo.KeyPair.SignMessage([32]byte{0})},
			RespondedAt:      time.Now(),
			Error:            nil,
		})
		assert.NoError(t, err)

		req1 := &bat.DispersalRequest{
			BatchID:        batchID,
			MinibatchIndex: 1,
			OperatorID:     opID,
			Socket:         opInfo.DispersalPort,
			NumBlobs:       2,
			RequestedAt:    time.Now(),
			BlobHash:       "1",
			MetadataHash:   "1",
		}
		err = b.MinibatchStore.PutDispersalRequest(context.Background(), req1)
		assert.NoError(t, err)
		err = b.MinibatchStore.PutDispersalResponse(context.Background(), &bat.DispersalResponse{
			DispersalRequest: *req1,
			Signatures:       []*core.Signature{opInfo.KeyPair.SignMessage([32]byte{1})},
			RespondedAt:      time.Now(),
			Error:            nil,
		})
		assert.NoError(t, err)
	}

	b.Minibatcher.Batches[batchID] = &bat.BatchState{
		BatchID:              batchID,
		ReferenceBlockNumber: uint(initialBlock),
		BlobHeaders: []*core.BlobHeader{
			{
				AccountID: "0",
				QuorumInfos: []*core.BlobQuorumInfo{
					{
						SecurityParam: core.SecurityParam{
							QuorumID:              0,
							AdversaryThreshold:    30,
							ConfirmationThreshold: 80,
						},
					},
					{
						SecurityParam: core.SecurityParam{
							QuorumID:              1,
							AdversaryThreshold:    30,
							ConfirmationThreshold: 80,
						},
					},
				},
			},
			{
				AccountID: "1",
				QuorumInfos: []*core.BlobQuorumInfo{
					{
						SecurityParam: core.SecurityParam{
							QuorumID:              0,
							AdversaryThreshold:    30,
							ConfirmationThreshold: 80,
						},
					},
					{
						SecurityParam: core.SecurityParam{
							QuorumID:              1,
							AdversaryThreshold:    30,
							ConfirmationThreshold: 80,
						},
					},
				},
			},
		},
		BlobMetadata:   []*disperser.BlobMetadata{},
		OperatorState:  operatorState.IndexedOperatorState,
		NumMinibatches: 2,
	}

	signChan := make(chan core.SigningMessage, 4)
	batchHeaderHash := [32]byte{93, 156, 41, 17, 3, 78, 159, 243, 222, 111, 54, 107, 237, 48, 243, 176, 224, 151, 96, 151, 159, 99, 118, 186, 53, 192, 72, 59, 160, 73, 7, 213}
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
}
