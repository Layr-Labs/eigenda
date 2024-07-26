package inmem_test

import (
	"context"
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
	ctx := context.Background()
	err = s.PutBatch(ctx, batch)
	assert.NoError(t, err)
	b, err := s.GetBatch(ctx, batch.ID)
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
	ctx := context.Background()
	err = s.PutMinibatch(ctx, minibatch)
	assert.NoError(t, err)
	m, err := s.GetMinibatch(ctx, minibatch.BatchID, minibatch.MinibatchIndex)
	assert.NoError(t, err)
	assert.Equal(t, minibatch, m)
}

func TestPutDispersalRequest(t *testing.T) {
	s := newMinibatchStore()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	minibatchIndex := uint(0)
	ctx := context.Background()
	req1 := &batcher.DispersalRequest{
		BatchID:         id,
		MinibatchIndex:  minibatchIndex,
		OperatorID:      core.OperatorID([32]byte{1}),
		OperatorAddress: gcommon.HexToAddress("0x0"),
		NumBlobs:        1,
		RequestedAt:     time.Now().UTC(),
	}
	err = s.PutDispersalRequest(ctx, req1)
	assert.NoError(t, err)
	req2 := &batcher.DispersalRequest{
		BatchID:         id,
		MinibatchIndex:  minibatchIndex,
		OperatorID:      core.OperatorID([32]byte{2}),
		OperatorAddress: gcommon.HexToAddress("0x0"),
		NumBlobs:        1,
		RequestedAt:     time.Now().UTC(),
	}
	err = s.PutDispersalRequest(ctx, req2)
	assert.NoError(t, err)

	r, err := s.GetMinibatchDispersalRequests(ctx, id, minibatchIndex)
	assert.NoError(t, err)
	assert.Len(t, r, 2)
	assert.Equal(t, req1, r[0])
	assert.Equal(t, req2, r[1])

	req, err := s.GetDispersalRequest(ctx, id, minibatchIndex, req1.OperatorID)
	assert.NoError(t, err)
	assert.Equal(t, req1, req)

	req, err = s.GetDispersalRequest(ctx, id, minibatchIndex, req2.OperatorID)
	assert.NoError(t, err)
	assert.Equal(t, req2, req)
}

func TestPutDispersalResponse(t *testing.T) {
	s := newMinibatchStore()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	ctx := context.Background()
	minibatchIndex := uint(0)
	resp1 := &batcher.DispersalResponse{
		DispersalRequest: batcher.DispersalRequest{
			BatchID:         id,
			MinibatchIndex:  minibatchIndex,
			OperatorID:      core.OperatorID([32]byte{1}),
			OperatorAddress: gcommon.HexToAddress("0x0"),
			NumBlobs:        1,
			RequestedAt:     time.Now().UTC(),
		},
		Signatures:  nil,
		RespondedAt: time.Now().UTC(),
		Error:       nil,
	}
	resp2 := &batcher.DispersalResponse{
		DispersalRequest: batcher.DispersalRequest{
			BatchID:         id,
			MinibatchIndex:  minibatchIndex,
			OperatorID:      core.OperatorID([32]byte{2}),
			OperatorAddress: gcommon.HexToAddress("0x0"),
			NumBlobs:        1,
			RequestedAt:     time.Now().UTC(),
		},
		Signatures:  nil,
		RespondedAt: time.Now().UTC(),
		Error:       nil,
	}
	err = s.PutDispersalResponse(ctx, resp1)
	assert.NoError(t, err)
	err = s.PutDispersalResponse(ctx, resp2)
	assert.NoError(t, err)

	r, err := s.GetMinibatchDispersalResponses(ctx, id, minibatchIndex)
	assert.NoError(t, err)
	assert.Len(t, r, 2)

	resp, err := s.GetDispersalResponse(ctx, id, minibatchIndex, resp1.OperatorID)
	assert.NoError(t, err)
	assert.Equal(t, resp1, resp)

	resp, err = s.GetDispersalResponse(ctx, id, minibatchIndex, resp2.OperatorID)
	assert.NoError(t, err)
	assert.Equal(t, resp2, resp)
}
