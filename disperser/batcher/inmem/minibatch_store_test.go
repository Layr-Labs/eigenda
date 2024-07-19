package inmem_test

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/batcher/inmem"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newMinibatchStore() batcher.MinibatchStore {
	return inmem.NewMinibatchStore(nil)
}

func TestPutBatch(t *testing.T) {
	s := newMinibatchStore()
	id, err := uuid.NewV7()
	assert.NoError(t, err)

	batch := &batcher.BatchRecord{
		ID:                   id,
		CreatedAt:            time.Now().UTC(),
		ReferenceBlockNumber: 1,
		HeaderHash:           [32]byte{1},
		AggregatePubKey:      nil,
		AggregateSignature:   nil,
	}
	err = s.PutBatch(batch)
	assert.NoError(t, err)
	b, err := s.GetBatch(batch.ID)
	assert.NoError(t, err)
	assert.Equal(t, batch, b)
}

func TestPutMiniBatch(t *testing.T) {
	s := newMinibatchStore()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	minibatch := &batcher.MinibatchRecord{
		BatchID:              id,
		MinibatchIndex:       12,
		BlobHeaderHashes:     [][32]byte{{1}},
		BatchSize:            1,
		ReferenceBlockNumber: 1,
	}
	err = s.PutMiniBatch(minibatch)
	assert.NoError(t, err)
	m, err := s.GetMiniBatch(minibatch.BatchID, minibatch.MinibatchIndex)
	assert.NoError(t, err)
	assert.Equal(t, minibatch, m)
}

func TestPutDispersalRequest(t *testing.T) {
	s := newMinibatchStore()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	request := &batcher.DispersalRequest{
		BatchID:         id,
		MinibatchIndex:  0,
		OperatorID:      core.OperatorID([32]byte{1}),
		OperatorAddress: gcommon.HexToAddress("0x0"),
		NumBlobs:        1,
		RequestedAt:     time.Now().UTC(),
	}
	err = s.PutDispersalRequest(request)
	assert.NoError(t, err)
	r, err := s.GetDispersalRequest(request.BatchID, request.MinibatchIndex)
	assert.NoError(t, err)
	assert.Equal(t, request, r)
}

func TestPutDispersalResponse(t *testing.T) {
	s := newMinibatchStore()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	response := &batcher.DispersalResponse{
		DispersalRequest: batcher.DispersalRequest{
			BatchID:         id,
			MinibatchIndex:  0,
			OperatorID:      core.OperatorID([32]byte{1}),
			OperatorAddress: gcommon.HexToAddress("0x0"),
			NumBlobs:        1,
			RequestedAt:     time.Now().UTC(),
		},
		Signatures:  nil,
		RespondedAt: time.Now().UTC(),
		Error:       nil,
	}
	err = s.PutDispersalResponse(response)
	assert.NoError(t, err)
	r, err := s.GetDispersalResponse(response.BatchID, response.MinibatchIndex)
	assert.NoError(t, err)
	assert.Equal(t, response, r)
}
