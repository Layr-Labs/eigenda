package inmem_test

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/batcher/inmem"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
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
		Status:               1,
	}
	ctx := context.Background()
	err = s.PutBatch(ctx, batch)
	assert.NoError(t, err)
	b, err := s.GetBatch(ctx, batch.ID)
	assert.NoError(t, err)
	assert.Equal(t, batch, b)
}

func TestPutDispersal(t *testing.T) {
	s := newMinibatchStore()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	ctx := context.Background()
	minibatchIndex := uint(0)
	dispersal1 := &batcher.MinibatchDispersal{
		DispersalResponse: batcher.DispersalResponse{
			Signatures:  nil,
			RespondedAt: time.Now().UTC(),
			Error:       nil,
		},
		BatchID:         id,
		MinibatchIndex:  minibatchIndex,
		OperatorID:      core.OperatorID([32]byte{1}),
		OperatorAddress: gcommon.HexToAddress("0x0"),
		NumBlobs:        1,
		RequestedAt:     time.Now().UTC(),
	}
	dispersal2 := &batcher.MinibatchDispersal{
		DispersalResponse: batcher.DispersalResponse{
			Signatures:  nil,
			RespondedAt: time.Now().UTC(),
			Error:       nil,
		},
		BatchID:         id,
		MinibatchIndex:  minibatchIndex,
		OperatorID:      core.OperatorID([32]byte{2}),
		OperatorAddress: gcommon.HexToAddress("0x0"),
		NumBlobs:        1,
		RequestedAt:     time.Now().UTC(),
	}
	err = s.PutDispersal(ctx, dispersal1)
	assert.NoError(t, err)
	err = s.PutDispersal(ctx, dispersal2)
	assert.NoError(t, err)

	r, err := s.GetDispersalsByMinibatch(ctx, id, minibatchIndex)
	assert.NoError(t, err)
	assert.Len(t, r, 2)
	assert.Equal(t, dispersal1, r[0])
	assert.Equal(t, dispersal2, r[1])

	resp, err := s.GetDispersal(ctx, id, minibatchIndex, dispersal1.OperatorID)
	assert.NoError(t, err)
	assert.Equal(t, dispersal1, resp)

	resp, err = s.GetDispersal(ctx, id, minibatchIndex, dispersal2.OperatorID)
	assert.NoError(t, err)
	assert.Equal(t, dispersal2, resp)
}

func TestPutBlobMinibatchMappings(t *testing.T) {
	s := newMinibatchStore()
	ctx := context.Background()
	batchID, err := uuid.NewV7()
	assert.NoError(t, err)
	blobKey := disperser.BlobKey{
		BlobHash:     "blob-hash",
		MetadataHash: "metadata-hash",
	}
	var commitX, commitY, lengthX, lengthY fp.Element
	_, err = commitX.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	assert.NoError(t, err)
	_, err = commitY.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	assert.NoError(t, err)
	commitment := &encoding.G1Commitment{
		X: commitX,
		Y: commitY,
	}
	_, err = lengthX.SetString("18730744272503541936633286178165146673834730535090946570310418711896464442549")
	assert.NoError(t, err)
	_, err = lengthY.SetString("15356431458378126778840641829778151778222945686256112821552210070627093656047")
	assert.NoError(t, err)
	var lengthXA0, lengthXA1, lengthYA0, lengthYA1 fp.Element
	_, err = lengthXA0.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	assert.NoError(t, err)
	_, err = lengthXA1.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	assert.NoError(t, err)
	_, err = lengthYA0.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	assert.NoError(t, err)
	_, err = lengthYA1.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531")
	assert.NoError(t, err)

	var lengthProof, lengthCommitment bn254.G2Affine
	lengthProof.X.A0 = lengthXA0
	lengthProof.X.A1 = lengthXA1
	lengthProof.Y.A0 = lengthYA0
	lengthProof.Y.A1 = lengthYA1

	lengthCommitment = lengthProof
	expectedDataLength := 111
	expectedChunkLength := uint(222)
	err = s.PutBlobMinibatchMappings(ctx, []*batcher.BlobMinibatchMapping{
		{
			BlobKey:        &blobKey,
			BatchID:        batchID,
			MinibatchIndex: 11,
			BlobIndex:      22,
			BlobHeader: core.BlobHeader{
				BlobCommitments: encoding.BlobCommitments{
					Commitment:       commitment,
					LengthCommitment: (*encoding.G2Commitment)(&lengthCommitment),
					Length:           uint(expectedDataLength),
					LengthProof:      (*encoding.LengthProof)(&lengthProof),
				},
				QuorumInfos: []*core.BlobQuorumInfo{
					{
						ChunkLength: expectedChunkLength,
						SecurityParam: core.SecurityParam{
							QuorumID:              1,
							ConfirmationThreshold: 55,
							AdversaryThreshold:    33,
							QuorumRate:            123,
						},
					},
				},
				AccountID: "account-id",
			},
		},
	})
	assert.NoError(t, err)

	mapping, err := s.GetBlobMinibatchMappings(ctx, blobKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mapping))
	assert.Equal(t, &blobKey, mapping[0].BlobKey)
	assert.Equal(t, batchID, mapping[0].BatchID)
	assert.Equal(t, uint(11), mapping[0].MinibatchIndex)
	assert.Equal(t, uint(22), mapping[0].BlobIndex)
	assert.Equal(t, commitment, mapping[0].BlobCommitments.Commitment)
	lengthCommitmentBytes, err := mapping[0].BlobCommitments.LengthCommitment.Serialize()
	assert.NoError(t, err)
	expectedLengthCommitmentBytes := lengthCommitment.Bytes()
	assert.Equal(t, expectedLengthCommitmentBytes[:], lengthCommitmentBytes[:])
	assert.Equal(t, expectedDataLength, int(mapping[0].BlobCommitments.Length))
	lengthProofBytes, err := mapping[0].BlobCommitments.LengthProof.Serialize()
	assert.NoError(t, err)
	expectedLengthProofBytes := lengthProof.Bytes()
	assert.Equal(t, expectedLengthProofBytes[:], lengthProofBytes[:])
	assert.Len(t, mapping[0].QuorumInfos, 1)
	assert.Equal(t, expectedChunkLength, mapping[0].QuorumInfos[0].ChunkLength)
	assert.Equal(t, core.SecurityParam{
		QuorumID:              1,
		ConfirmationThreshold: 55,
		AdversaryThreshold:    33,
		QuorumRate:            123,
	}, mapping[0].QuorumInfos[0].SecurityParam)
}
