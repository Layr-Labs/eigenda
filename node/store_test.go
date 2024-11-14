package node_test

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"testing"
	"time"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	pb "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/metrics"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

const (
	staleMeasure  = uint32(1)
	storeDuration = uint32(1)
)

func CreateBatch(t *testing.T) (*core.BatchHeader, []*core.BlobMessage, []*pb.Blob) {
	return CreateBatchWith(t, false)
}

// Creates a batch and returns its header and blobs.
func CreateBatchWith(t *testing.T, encodeBundle bool) (*core.BatchHeader, []*core.BlobMessage, []*pb.Blob) {
	var commitX, commitY, lengthX, lengthY fp.Element
	_, err := commitX.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	assert.NoError(t, err)
	_, err = commitY.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	assert.NoError(t, err)
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

	commitment := bn254.G1Affine{
		X: commitX,
		Y: commitY,
	}

	adversaryThreshold := uint8(90)
	quorumThreshold := uint8(100)

	quorumHeader := &core.BlobQuorumInfo{
		SecurityParam: core.SecurityParam{
			QuorumID:              0,
			ConfirmationThreshold: quorumThreshold,
			AdversaryThreshold:    adversaryThreshold,
		},
		ChunkLength: 10,
	}
	chunk1 := &encoding.Frame{
		Proof:  commitment,
		Coeffs: []encoding.Symbol{encoding.ONE},
	}
	chunk1bytes, err := chunk1.Serialize()
	assert.Nil(t, err)
	bundle1 := core.Bundle{
		chunk1,
	}
	bundle1bytes, err := bundle1.Serialize()
	assert.Nil(t, err)

	blobMessage := []*core.BlobMessage{
		{
			BlobHeader: &core.BlobHeader{
				BlobCommitments: encoding.BlobCommitments{
					Commitment:       (*encoding.G1Commitment)(&commitment),
					LengthCommitment: (*encoding.G2Commitment)(&lengthCommitment),
					LengthProof:      (*encoding.LengthProof)(&lengthProof),
					Length:           48,
				},
				QuorumInfos: []*core.BlobQuorumInfo{quorumHeader},
			},
			Bundles: core.Bundles{
				core.QuorumID(0): []*encoding.Frame{
					chunk1,
				},
			},
		},
		{
			BlobHeader: &core.BlobHeader{
				BlobCommitments: encoding.BlobCommitments{
					Commitment:       (*encoding.G1Commitment)(&commitment),
					LengthCommitment: (*encoding.G2Commitment)(&lengthCommitment),
					LengthProof:      (*encoding.G2Commitment)(&lengthProof),
					Length:           50,
				},
				QuorumInfos: []*core.BlobQuorumInfo{quorumHeader},
			},
			Bundles: core.Bundles{
				core.QuorumID(0): []*encoding.Frame{
					chunk1,
				},
			},
		},
	}

	batchHeader := core.BatchHeader{
		BatchRoot:            [32]byte{0},
		ReferenceBlockNumber: 0,
	}

	quorumHeaderProto := &pb.BlobQuorumInfo{
		QuorumId:              uint32(quorumHeader.QuorumID),
		AdversaryThreshold:    uint32(quorumHeader.AdversaryThreshold),
		ConfirmationThreshold: uint32(quorumHeader.ConfirmationThreshold),
		ChunkLength:           uint32(quorumHeader.ChunkLength),
	}

	blobHeaderProto0 := &pb.BlobHeader{
		Commitment: &commonpb.G1Commitment{
			X: commitment.X.Marshal(),
			Y: commitment.Y.Marshal(),
		},
		LengthCommitment: &pb.G2Commitment{
			XA0: lengthCommitment.X.A0.Marshal(),
			XA1: lengthCommitment.X.A1.Marshal(),
			YA0: lengthCommitment.Y.A0.Marshal(),
			YA1: lengthCommitment.Y.A1.Marshal(),
		},
		LengthProof: &pb.G2Commitment{
			XA0: lengthProof.X.A0.Marshal(),
			XA1: lengthProof.X.A1.Marshal(),
			YA0: lengthProof.Y.A0.Marshal(),
			YA1: lengthProof.Y.A1.Marshal(),
		},
		Length:        uint32(48),
		QuorumHeaders: []*pb.BlobQuorumInfo{quorumHeaderProto},
	}

	blobHeaderProto1 := &pb.BlobHeader{
		Commitment: &commonpb.G1Commitment{
			X: commitment.X.Marshal(),
			Y: commitment.Y.Marshal(),
		},
		LengthCommitment: &pb.G2Commitment{
			XA0: lengthCommitment.X.A0.Marshal(),
			XA1: lengthCommitment.X.A1.Marshal(),
			YA0: lengthCommitment.Y.A0.Marshal(),
			YA1: lengthCommitment.Y.A1.Marshal(),
		},
		LengthProof: &pb.G2Commitment{
			XA0: lengthProof.X.A0.Marshal(),
			XA1: lengthProof.X.A1.Marshal(),
			YA0: lengthProof.Y.A0.Marshal(),
			YA1: lengthProof.Y.A1.Marshal(),
		},
		Length:        uint32(50),
		QuorumHeaders: []*pb.BlobQuorumInfo{quorumHeaderProto},
	}
	var bundles []*pb.Bundle
	if encodeBundle {
		bundles = []*pb.Bundle{
			{
				Bundle: bundle1bytes,
			},
		}
	} else {
		bundles = []*pb.Bundle{
			{
				Chunks: [][]byte{
					chunk1bytes,
				},
			},
		}

	}
	blobs := []*pb.Blob{
		{
			Header:  blobHeaderProto0,
			Bundles: bundles,
		},
		{
			Header:  blobHeaderProto1,
			Bundles: bundles,
		},
	}
	return &batchHeader, blobMessage, blobs
}

func createStore(t *testing.T) *node.Store {
	noopMetrics := metrics.NewNoopMetrics()
	reg := prometheus.NewRegistry()
	logger := logging.NewNoopLogger()
	operatorId := [32]byte(hexutil.MustDecode("0x3fbfefcdc76462d2cdb7d0cea75f27223829481b8b4aa6881c94cb2126a316ad"))
	tx := &coremock.MockWriter{}
	dat, _ := mock.MakeChainDataMock(map[uint8]int{
		0: 6,
		1: 3,
	})
	s, _ := node.NewLevelDBStore(t.TempDir(), logger, node.NewMetrics(noopMetrics, reg, logger, ":9090", operatorId, -1, tx, dat), staleMeasure, storeDuration)
	return s
}

func TestEncodeDecodeChunks(t *testing.T) {
	decoded, format, err := node.DecodeChunks([]byte{})
	assert.Nil(t, err)
	assert.Equal(t, pb.ChunkEncodingFormat_UNKNOWN, format)
	assert.Equal(t, 0, len(decoded))

	numSamples := 32
	numChunks := 10
	chunkSize := 2 * 1024
	for n := 0; n < numSamples; n++ {
		chunks := make([][]byte, numChunks)
		for i := 0; i < numChunks; i++ {
			chunk := make([]byte, chunkSize)
			_, _ = cryptorand.Read(chunk)
			chunks[i] = chunk
		}
		encoded, err := node.EncodeChunks(chunks)
		assert.Nil(t, err)
		decoded, format, err := node.DecodeChunks(encoded)
		assert.Nil(t, err)
		assert.Equal(t, pb.ChunkEncodingFormat_GOB, format)
		for i := 0; i < numChunks; i++ {
			assert.True(t, bytes.Equal(decoded[i], chunks[i]))
		}
	}
}

func TestStoreBatchInvalidBlob(t *testing.T) {
	s := createStore(t)
	ctx := context.Background()
	batchHeader, blobs, blobsProto := CreateBatchWith(t, true)
	blobsProto[0].Bundles[0].Chunks = [][]byte{[]byte{1}}
	_, err := s.StoreBatch(ctx, batchHeader, blobs, blobsProto)
	assert.EqualError(t, err, "chunks of a bundle are encoded together already")
}

func TestStoreBatchSuccess(t *testing.T) {
	s := createStore(t)
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

	blobHeaderBytes2, err := s.GetBlobHeader(ctx, batchHeaderHash, 1)
	assert.Nil(t, err)
	expected, err = proto.Marshal(blobsProto[1].GetHeader())
	assert.Nil(t, err)
	assert.True(t, bytes.Equal(blobHeaderBytes2, expected))

	blobHeaderBytes3, err := s.GetBlobHeader(ctx, batchHeaderHash, 2)
	assert.ErrorIs(t, err, node.ErrKeyNotFound)
	assert.Nil(t, blobHeaderBytes3)

	// Check existence: blob chunks.
	blobKey1, err := node.EncodeBlobKey(batchHeaderHash, 0, 0)
	assert.Nil(t, err)
	assert.True(t, s.HasKey(ctx, blobKey1))
	blobKey2, err := node.EncodeBlobKey(batchHeaderHash, 1, 0)
	assert.Nil(t, err)
	assert.True(t, s.HasKey(ctx, blobKey2))

	// Check the chunks.
	chunks, format, err := s.GetChunks(ctx, batchHeaderHash, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, pb.ChunkEncodingFormat_GOB, format)
	assert.Equal(t, chunks, blobsProto[0].Bundles[0].Chunks)
	chunks, format, err = s.GetChunks(ctx, batchHeaderHash, 1, 0)
	assert.Nil(t, err)
	assert.Equal(t, pb.ChunkEncodingFormat_GOB, format)
	assert.Equal(t, chunks, blobsProto[1].Bundles[0].Chunks)

	// Store the batch again it should be no-op.
	_, err = s.StoreBatch(ctx, batchHeader, blobs, blobsProto)
	assert.NotNil(t, err)
	assert.Equal(t, err, node.ErrBatchAlreadyExist)

	// Expire the batches.
	curTime := time.Now().Unix() + int64(staleMeasure+storeDuration)*12
	// Try to expire at a time before expiry, so nothing will be expired.
	numDeleted, _, _, err := s.DeleteExpiredEntries(curTime-10, 1)
	assert.Nil(t, err)
	assert.Equal(t, numDeleted, 0)
	assert.True(t, s.HasKey(ctx, batchHeaderKey))
	// Then expire it at a time post expiry, so the batch will get purged.
	numDeleted, _, _, err = s.DeleteExpiredEntries(curTime+10, 1)
	assert.Nil(t, err)
	assert.Equal(t, numDeleted, 1)
	assert.False(t, s.HasKey(ctx, batchHeaderKey))
	assert.False(t, s.HasKey(ctx, blobHeaderKey1))
	assert.False(t, s.HasKey(ctx, blobHeaderKey2))
	assert.False(t, s.HasKey(ctx, blobKey1))
	assert.False(t, s.HasKey(ctx, blobKey2))
}

func decodeChunks(t *testing.T, s *node.Store, batchHeaderHash [32]byte, blobIdx int, chunkEncoding pb.ChunkEncodingFormat) []*encoding.Frame {
	ctx := context.Background()
	chunks, format, err := s.GetChunks(ctx, batchHeaderHash, blobIdx, 0)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(chunks))
	assert.Equal(t, chunkEncoding, format)
	var f *encoding.Frame
	switch chunkEncoding {
	case pb.ChunkEncodingFormat_GOB:
		f, err = new(encoding.Frame).Deserialize(chunks[0])
		assert.Nil(t, err)
	case pb.ChunkEncodingFormat_GNARK:
		f, err = new(encoding.Frame).DeserializeGnark(chunks[0])
		assert.Nil(t, err)
	}
	return []*encoding.Frame{f}
}

func checkBundleEquivalence(t *testing.T, bundle1, bundle2 []*encoding.Frame) {
	assert.Equal(t, len(bundle1), len(bundle2))
	for i := 0; i < len(bundle1); i++ {
		assert.True(t, bundle1[i].Proof.Equal(&bundle2[i].Proof))
		assert.Equal(t, len(bundle1[i].Coeffs), len(bundle2[i].Coeffs))
		for j := 0; j < len(bundle1[i].Coeffs); j++ {
			assert.True(t, bundle1[i].Coeffs[j].Equal(&bundle2[i].Coeffs[j]))
		}
	}
}

func TestBundleEncodingEquivalence(t *testing.T) {
	ctx := context.Background()
	// Gnark chunks
	s1 := createStore(t)
	batchHeader1, blobs1, blobsProto1 := CreateBatchWith(t, true)
	_, err := s1.StoreBatch(ctx, batchHeader1, blobs1, blobsProto1)
	assert.Nil(t, err)
	// Gob chunks
	s2 := createStore(t)
	batchHeader2, blobs2, blobsProto2 := CreateBatchWith(t, false)
	_, err = s2.StoreBatch(ctx, batchHeader2, blobs2, blobsProto2)
	assert.Nil(t, err)

	// Check parity
	batchHeaderHash, err := batchHeader1.GetBatchHeaderHash()
	assert.Nil(t, err)
	// The first blob
	bundle1 := decodeChunks(t, s1, batchHeaderHash, 0, pb.ChunkEncodingFormat_GNARK)
	bundle2 := decodeChunks(t, s2, batchHeaderHash, 0, pb.ChunkEncodingFormat_GOB)
	checkBundleEquivalence(t, bundle1, bundle2)
	// The second blob
	bundle1 = decodeChunks(t, s1, batchHeaderHash, 1, pb.ChunkEncodingFormat_GNARK)
	bundle2 = decodeChunks(t, s2, batchHeaderHash, 1, pb.ChunkEncodingFormat_GOB)
	checkBundleEquivalence(t, bundle1, bundle2)
}

func BenchmarkEncodeChunks(b *testing.B) {
	numSamples := 32
	numChunks := 10
	chunkSize := 2 * 1024
	sampleChunks := make([][][]byte, numSamples)
	for n := 0; n < numSamples; n++ {
		chunks := make([][]byte, numChunks)
		for i := 0; i < numChunks; i++ {
			chunk := make([]byte, chunkSize)
			_, _ = cryptorand.Read(chunk)
			chunks[i] = chunk
		}
		sampleChunks[n] = chunks
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = node.EncodeChunks(sampleChunks[i%numSamples])
	}
}

func BenchmarkDecocodeChunks(b *testing.B) {
	numSamples := 32
	numChunks := 10
	chunkSize := 2 * 1024
	sampleChunks := make([][]byte, numSamples)
	for n := 0; n < numSamples; n++ {
		chunks := make([][]byte, numChunks)
		for i := 0; i < numChunks; i++ {
			chunk := make([]byte, chunkSize)
			_, _ = cryptorand.Read(chunk)
			chunks[i] = chunk
		}
		encoded, _ := node.EncodeChunks(chunks)
		sampleChunks[n] = encoded
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = node.DecodeChunks(sampleChunks[i%numSamples])
	}
}
