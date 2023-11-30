package node_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	"github.com/Layr-Labs/eigensdk-go/metrics"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

// Creates a batch and returns its header and blobs.
func CreateBatch(t *testing.T) (*core.BatchHeader, []*core.BlobMessage, []*pb.Blob) {
	var commitX, commitY, lengthX, lengthY fp.Element
	_, err := commitX.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	assert.NoError(t, err)
	_, err = commitY.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	assert.NoError(t, err)
	_, err = lengthX.SetString("18730744272503541936633286178165146673834730535090946570310418711896464442549")
	assert.NoError(t, err)
	_, err = lengthY.SetString("15356431458378126778840641829778151778222945686256112821552210070627093656047")
	assert.NoError(t, err)

	commitment := bn254.G1Point{
		X: commitX,
		Y: commitY,
	}
	lengthProof := bn254.G1Point{
		X: lengthX,
		Y: lengthY,
	}

	numOperators := uint(10)
	adversaryThreshold := uint8(90)

	quorumHeader := &core.BlobQuorumInfo{
		SecurityParam: core.SecurityParam{
			QuorumID:           0,
			AdversaryThreshold: adversaryThreshold,
		},
		QuantizationFactor: 1,
		EncodedBlobLength:  32 * 1 * numOperators,
	}
	chunk1 := &core.Chunk{
		Proof:  commitment,
		Coeffs: []core.Symbol{bn254.ONE},
	}

	blobMessage := []*core.BlobMessage{
		{
			BlobHeader: &core.BlobHeader{
				BlobCommitments: core.BlobCommitments{
					Commitment:  &core.Commitment{G1Point: &commitment},
					LengthProof: &core.Commitment{G1Point: &lengthProof},
					Length:      48,
				},
				QuorumInfos: []*core.BlobQuorumInfo{quorumHeader},
			},
			Bundles: core.Bundles{
				core.QuorumID(0): []*core.Chunk{
					chunk1,
				},
			},
		},
		{
			BlobHeader: &core.BlobHeader{
				BlobCommitments: core.BlobCommitments{
					Commitment: &core.Commitment{G1Point: &commitment},
					Length:     50,
				},
				QuorumInfos: []*core.BlobQuorumInfo{quorumHeader},
			},
			Bundles: core.Bundles{
				core.QuorumID(0): []*core.Chunk{
					chunk1,
				},
			},
		},
	}

	batchHeader := core.BatchHeader{
		BatchRoot:            [32]byte{0},
		ReferenceBlockNumber: 0,
	}

	serializedCommitment0, err := core.Commitment{G1Point: &commitment}.Serialize()
	assert.NoError(t, err)
	serializedLengthProof0, err := core.Commitment{G1Point: &lengthProof}.Serialize()
	assert.NoError(t, err)

	quorumHeaderProto := &pb.BlobQuorumInfo{
		QuorumId:           uint32(quorumHeader.QuorumID),
		AdversaryThreshold: uint32(quorumHeader.AdversaryThreshold), QuantizationFactor: uint32(quorumHeader.QuantizationFactor),
		EncodedBlobLength: uint32(quorumHeader.EncodedBlobLength),
	}

	blobHeaderProto0 := &pb.BlobHeader{
		Commitment:    serializedCommitment0,
		LengthProof:   serializedLengthProof0,
		Length:        uint32(48),
		QuorumHeaders: []*pb.BlobQuorumInfo{quorumHeaderProto},
	}

	blobHeaderProto1 := &pb.BlobHeader{
		Commitment:    serializedCommitment0,
		LengthProof:   serializedLengthProof0,
		Length:        uint32(50),
		QuorumHeaders: []*pb.BlobQuorumInfo{quorumHeaderProto},
	}
	blobs := []*pb.Blob{
		{
			Header: blobHeaderProto0,
		},
		{
			Header: blobHeaderProto1,
		},
	}
	return &batchHeader, blobMessage, blobs
}

func TestStoringBlob(t *testing.T) {
	staleMeasure := uint32(1)
	storeDuration := uint32(1)
	noopMetrics := metrics.NewNoopMetrics()
	reg := prometheus.NewRegistry()
	s, _ := node.NewLevelDBStore(t.TempDir(), &mock.Logger{}, node.NewMetrics(noopMetrics, reg, &mock.Logger{}, ":9090"), staleMeasure, storeDuration)
	ctx := context.Background()

	// Empty store
	blobKey := []byte{1, 2}
	assert.False(t, s.HasKey(ctx, blobKey))

	// Prepare data to store.
	batchHeader, blobs, blobsProto := CreateBatch(t)
	batchHeaderBytes, _ := batchHeader.Serialize()

	// Store a batch.
	_, err := s.StoreBatch(ctx, batchHeader, blobs, blobsProto)
	assert.Nil(t, err)

	// Check existence: batch header.
	batchHeaderHash, err := batchHeader.GetBatchHeaderHash()
	assert.Nil(t, err)
	batchHeaderKey := node.EncodeBatchHeaderKey(batchHeaderHash)
	assert.True(t, s.HasKey(ctx, batchHeaderKey))
	header, err := s.GetBatchHeader(ctx, batchHeaderHash)
	assert.Nil(t, err)
	assert.True(t, bytes.Equal(header, batchHeaderBytes))

	// Check existence: blob headers.
	blobHeaderKey1, err := node.EncodeBlobHeaderKey(batchHeaderHash, 0)
	assert.Nil(t, err)
	assert.True(t, s.HasKey(ctx, blobHeaderKey1))
	blobHeaderKey2, err := node.EncodeBlobHeaderKey(batchHeaderHash, 1)
	assert.Nil(t, err)
	assert.True(t, s.HasKey(ctx, blobHeaderKey2))
	blobHeaderBytes1, err := s.GetBlobHeader(ctx, batchHeaderHash, 0)
	assert.Nil(t, err)
	expected, err := proto.Marshal(blobsProto[0].GetHeader())
	assert.Nil(t, err)
	assert.True(t, bytes.Equal(blobHeaderBytes1, expected))

	// Check existence: blob chunks.
	blobKey1, err := node.EncodeBlobKey(batchHeaderHash, 0, 0)
	assert.Nil(t, err)
	assert.True(t, s.HasKey(ctx, blobKey1))
	blobKey2, err := node.EncodeBlobKey(batchHeaderHash, 1, 0)
	assert.Nil(t, err)
	assert.True(t, s.HasKey(ctx, blobKey2))

	// Store the batch again it should be no-op.
	_, err = s.StoreBatch(ctx, batchHeader, blobs, blobsProto)
	assert.NotNil(t, err)
	assert.Equal(t, err, node.ErrBatchAlreadyExist)

	// Expire the batches.
	curTime := time.Now().Unix() + int64(staleMeasure+storeDuration)*12
	// Try to expire at a time before expiry, so nothing will be expired.
	numDeleted, err := s.DeleteExpiredEntries(curTime-10, 1)
	assert.Nil(t, err)
	assert.Equal(t, numDeleted, 0)
	assert.True(t, s.HasKey(ctx, batchHeaderKey))
	// Then expire it at a time post expiry, so the batch will get purged.
	numDeleted, err = s.DeleteExpiredEntries(curTime+10, 1)
	assert.Nil(t, err)
	assert.Equal(t, numDeleted, 1)
	assert.False(t, s.HasKey(ctx, batchHeaderKey))
	assert.False(t, s.HasKey(ctx, blobHeaderKey1))
	assert.False(t, s.HasKey(ctx, blobHeaderKey2))
	assert.False(t, s.HasKey(ctx, blobKey1))
	assert.False(t, s.HasKey(ctx, blobKey2))
}
