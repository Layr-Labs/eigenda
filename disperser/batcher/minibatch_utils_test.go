package batcher_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDispersedBlobsSuccess(t *testing.T) {
	ctx := context.Background()
	operators, err := mockChainState.GetIndexedOperators(ctx, uint(0))
	assert.NoError(t, err)
	batchID := uuid.New()
	// Two minibatches
	dispersals := make([]*batcher.MinibatchDispersal, 0)
	for opID := range operators {
		dispersals = append(dispersals, &batcher.MinibatchDispersal{
			BatchID:        batchID,
			MinibatchIndex: 0,
			OperatorID:     opID,
			RequestedAt:    time.Now(),
			DispersalResponse: batcher.DispersalResponse{
				Signatures:  [][32]byte{{0}},
				RespondedAt: time.Now(),
				Error:       nil,
			},
		})
		dispersals = append(dispersals, &batcher.MinibatchDispersal{
			BatchID:        batchID,
			MinibatchIndex: 1,
			OperatorID:     opID,
			RequestedAt:    time.Now(),
			DispersalResponse: batcher.DispersalResponse{
				Signatures:  [][32]byte{{0}},
				RespondedAt: time.Now(),
				Error:       nil,
			},
		})
	}

	// 2 blobs in each minibatch
	m0 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash0",
			MetadataHash: "metadataHash0",
		},
		MinibatchIndex: 0,
		BlobIndex:      0,
	}
	m1 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash1",
			MetadataHash: "metadataHash1",
		},
		MinibatchIndex: 0,
		BlobIndex:      1,
	}
	m2 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash2",
			MetadataHash: "metadataHash2",
		},
		MinibatchIndex: 1,
		BlobIndex:      0,
	}
	m3 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash3",
			MetadataHash: "metadataHash3",
		},
		MinibatchIndex: 1,
		BlobIndex:      1,
	}
	successful, failed, err := batcher.DispersedBlobs(operators, dispersals, []*batcher.BlobMinibatchMapping{m0, m1, m2, m3})
	assert.NoError(t, err)
	assert.Len(t, successful, 4)
	assert.Len(t, failed, 0)
	assert.ElementsMatch(t, successful, []*batcher.BlobMinibatchMapping{m0, m1, m2, m3})
}

func TestDispersedBlobsFailedDispersal(t *testing.T) {
	ctx := context.Background()
	operators, err := mockChainState.GetIndexedOperators(ctx, uint(0))
	assert.NoError(t, err)
	batchID := uuid.New()
	// Two minibatches
	// One minibatch has failed dispersal
	dispersals := make([]*batcher.MinibatchDispersal, 0)
	for opID := range operators {
		dispersals = append(dispersals, &batcher.MinibatchDispersal{
			BatchID:        batchID,
			MinibatchIndex: 0,
			OperatorID:     opID,
			RequestedAt:    time.Now(),
			DispersalResponse: batcher.DispersalResponse{
				Signatures:  [][32]byte{{0}},
				RespondedAt: time.Now(),
				Error:       nil,
			},
		})
		if opID == opId0 {
			dispersals = append(dispersals, &batcher.MinibatchDispersal{
				BatchID:        batchID,
				MinibatchIndex: 1,
				OperatorID:     opID,
				RequestedAt:    time.Now(),
				DispersalResponse: batcher.DispersalResponse{
					Signatures:  nil,
					RespondedAt: time.Now(),
					Error:       errors.New("unknown error"),
				},
			})
		} else {
			dispersals = append(dispersals, &batcher.MinibatchDispersal{
				BatchID:        batchID,
				MinibatchIndex: 1,
				OperatorID:     opID,
				RequestedAt:    time.Now(),
				DispersalResponse: batcher.DispersalResponse{
					Signatures:  [][32]byte{{0}},
					RespondedAt: time.Now(),
					Error:       nil,
				},
			})
		}
	}

	// 2 blobs in each minibatch
	m0 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash0",
			MetadataHash: "metadataHash0",
		},
		MinibatchIndex: 0,
		BlobIndex:      0,
	}
	m1 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash1",
			MetadataHash: "metadataHash1",
		},
		MinibatchIndex: 0,
		BlobIndex:      1,
	}
	m2 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash2",
			MetadataHash: "metadataHash2",
		},
		MinibatchIndex: 1,
		BlobIndex:      0,
	}
	m3 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash3",
			MetadataHash: "metadataHash3",
		},
		MinibatchIndex: 1,
		BlobIndex:      1,
	}
	successful, failed, err := batcher.DispersedBlobs(operators, dispersals, []*batcher.BlobMinibatchMapping{m0, m1, m2, m3})
	assert.NoError(t, err)
	assert.Len(t, successful, 2)
	assert.Len(t, failed, 2)
	assert.ElementsMatch(t, successful, []*batcher.BlobMinibatchMapping{m0, m1})
	assert.ElementsMatch(t, failed, []*batcher.BlobMinibatchMapping{m2, m3})
}

func TestDispersedBlobsMissingDispersal(t *testing.T) {
	ctx := context.Background()
	operators, err := mockChainState.GetIndexedOperators(ctx, uint(0))
	assert.NoError(t, err)
	batchID := uuid.New()
	// Two minibatches
	// One minibatch has missing dispersal
	dispersals := make([]*batcher.MinibatchDispersal, 0)
	for opID := range operators {
		dispersals = append(dispersals, &batcher.MinibatchDispersal{
			BatchID:        batchID,
			MinibatchIndex: 0,
			OperatorID:     opID,
			RequestedAt:    time.Now(),
			DispersalResponse: batcher.DispersalResponse{
				Signatures:  [][32]byte{{0}},
				RespondedAt: time.Now(),
				Error:       nil,
			},
		})
		if opID != opId0 {
			dispersals = append(dispersals, &batcher.MinibatchDispersal{
				BatchID:        batchID,
				MinibatchIndex: 1,
				OperatorID:     opID,
				RequestedAt:    time.Now(),
				DispersalResponse: batcher.DispersalResponse{
					Signatures:  [][32]byte{{0}},
					RespondedAt: time.Now(),
					Error:       nil,
				},
			})
		}
	}

	// 2 blobs in each minibatch
	m0 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash0",
			MetadataHash: "metadataHash0",
		},
		MinibatchIndex: 0,
		BlobIndex:      0,
	}
	m1 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash1",
			MetadataHash: "metadataHash1",
		},
		MinibatchIndex: 0,
		BlobIndex:      1,
	}
	m2 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash2",
			MetadataHash: "metadataHash2",
		},
		MinibatchIndex: 1,
		BlobIndex:      0,
	}
	m3 := &batcher.BlobMinibatchMapping{
		BatchID: batchID,
		BlobKey: &disperser.BlobKey{
			BlobHash:     "blobHash3",
			MetadataHash: "metadataHash3",
		},
		MinibatchIndex: 1,
		BlobIndex:      1,
	}
	successful, failed, err := batcher.DispersedBlobs(operators, dispersals, []*batcher.BlobMinibatchMapping{m0, m1, m2, m3})
	assert.NoError(t, err)
	assert.Len(t, successful, 2)
	assert.Len(t, failed, 2)
	assert.ElementsMatch(t, successful, []*batcher.BlobMinibatchMapping{m0, m1})
	assert.ElementsMatch(t, failed, []*batcher.BlobMinibatchMapping{m2, m3})
}
