package grpc_test

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/dispatcher"
	"github.com/stretchr/testify/assert"
)

func makeBatch(t *testing.T, blobSize int, numBlobs int, advThreshold, quorumThreshold int, refBlockNumber uint) (*core.BatchHeader, map[core.OperatorID][]*core.BlobMessage) {
	encoder, err := makeTestEncoder()
	assert.NoError(t, err)
	asn := &core.StdAssignmentCoordinator{}

	blobHeaders := make([]*core.BlobHeader, numBlobs)
	blobChunks := make([][]*core.Chunk, numBlobs)
	blobMessagesByOp := make(map[core.OperatorID][]*core.BlobMessage)
	for i := 0; i < numBlobs; i++ {
		// create data
		data := make([]byte, blobSize)
		_, err := rand.Read(data)
		assert.NoError(t, err)

		// encode data
		operatorState, err := chainState.GetOperatorState(context.Background(), 0, []core.QuorumID{0})
		assert.NoError(t, err)
		assignments, info, err := asn.GetAssignments(operatorState, 0, batcher.QuantizationFactor)
		assert.NoError(t, err)
		quorumInfo := batcher.QuorumInfo{
			Assignments:        assignments,
			Info:               info,
			QuantizationFactor: batcher.QuantizationFactor,
		}
		blobLength := core.GetBlobLength(uint(blobSize))
		numOperators := uint(len(quorumInfo.Assignments))
		chunkLength, err := asn.GetMinimumChunkLength(numOperators, blobLength, quorumInfo.QuantizationFactor, uint8(quorumThreshold), uint8(advThreshold))
		assert.NoError(t, err)
		params, err := core.GetEncodingParams(chunkLength, quorumInfo.Info.TotalChunks)
		assert.NoError(t, err)
		t.Logf("Encoding params: ChunkLength: %d, NumChunks: %d", params.ChunkLength, params.NumChunks)
		commits, chunks, err := encoder.Encode(data, params)
		assert.NoError(t, err)
		blobChunks[i] = chunks

		// populate blob header
		blobHeaders[i] = &core.BlobHeader{
			BlobCommitments: commits,
			QuorumInfos: []*core.BlobQuorumInfo{
				{
					SecurityParam: core.SecurityParam{
						QuorumID:           0,
						AdversaryThreshold: uint8(advThreshold),
						QuorumThreshold:    uint8(quorumThreshold),
					},
					QuantizationFactor: quorumInfo.QuantizationFactor,
					EncodedBlobLength:  params.ChunkLength * quorumInfo.QuantizationFactor * numOperators,
				},
			},
		}

		// populate blob messages
		for opID, assignment := range quorumInfo.Assignments {
			blobMessagesByOp[opID] = append(blobMessagesByOp[opID], &core.BlobMessage{
				BlobHeader: blobHeaders[i],
				Bundles:    make(core.Bundles),
			})
			blobMessagesByOp[opID][i].Bundles[0] = append(blobMessagesByOp[opID][i].Bundles[0], chunks[assignment.StartIndex:assignment.StartIndex+assignment.NumChunks]...)
		}
	}

	batchHeader := &core.BatchHeader{
		ReferenceBlockNumber: refBlockNumber,
		BatchRoot:            [32]byte{},
	}
	_, err = batchHeader.SetBatchRoot(blobHeaders)
	assert.NoError(t, err)
	return batchHeader, blobMessagesByOp
}

func TestStoreChunks(t *testing.T) {
	t.Skip("Skipping TestStoreChunks as it's meant to be tested manually to measure performance")

	server := newTestServer(t, false)
	// 50 X 200 KiB blobs
	batchHeader, blobMessagesByOp := makeBatch(t, 200*1024, 50, 80, 100, 0)
	numTotalChunks := 0
	for i := range blobMessagesByOp[opID] {
		numTotalChunks += len(blobMessagesByOp[opID][i].Bundles[0])
	}
	t.Logf("Batch numTotalChunks: %d", numTotalChunks)
	req, totalSize, err := dispatcher.GetStoreChunksRequest(blobMessagesByOp[opID], batchHeader)
	assert.NoError(t, err)
	assert.Equal(t, 50790400, totalSize)

	timer := time.Now()
	reply, err := server.StoreChunks(context.Background(), req)
	t.Log("StoreChunks took", time.Since(timer))
	assert.NoError(t, err)
	assert.NotNil(t, reply.GetSignature())
}
