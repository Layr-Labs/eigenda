package grpc_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	dispatcher "github.com/Layr-Labs/eigenda/disperser/batcher/grpc"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/assert"
)

func makeBatch(t *testing.T, blobSize int, numBlobs int, advThreshold, quorumThreshold int, refBlockNumber uint) (*core.BatchHeader, map[core.OperatorID][]*core.BlobMessage) {
	p, _, err := makeTestComponents()
	assert.NoError(t, err)
	asn := &core.StdAssignmentCoordinator{}

	blobHeaders := make([]*core.BlobHeader, numBlobs)
	blobChunks := make([][]*encoding.Frame, numBlobs)
	blobMessagesByOp := make(map[core.OperatorID][]*core.BlobMessage)
	for i := 0; i < numBlobs; i++ {
		// create data
		ranData := make([]byte, blobSize)
		_, err := rand.Read(ranData)
		assert.NoError(t, err)

		data := codec.ConvertByPaddingEmptyByte(ranData)
		data = data[:blobSize]

		operatorState, err := chainState.GetOperatorState(context.Background(), 0, []core.QuorumID{0})
		assert.NoError(t, err)

		chunkLength, err := asn.CalculateChunkLength(operatorState, encoding.GetBlobLength(uint(blobSize)), 0, &core.SecurityParam{
			QuorumID:              0,
			AdversaryThreshold:    uint8(advThreshold),
			ConfirmationThreshold: uint8(quorumThreshold),
		})
		assert.NoError(t, err)

		blobQuorumInfo := &core.BlobQuorumInfo{
			SecurityParam: core.SecurityParam{
				QuorumID:              0,
				AdversaryThreshold:    uint8(advThreshold),
				ConfirmationThreshold: uint8(quorumThreshold),
			},
			ChunkLength: chunkLength,
		}

		// encode data

		assignments, info, err := asn.GetAssignments(operatorState, encoding.GetBlobLength(uint(blobSize)), blobQuorumInfo)
		assert.NoError(t, err)
		quorumInfo := batcher.QuorumInfo{
			Assignments:        assignments,
			Info:               info,
			QuantizationFactor: batcher.QuantizationFactor,
		}

		params := encoding.ParamsFromMins(chunkLength, quorumInfo.Info.TotalChunks)
		t.Logf("Encoding params: ChunkLength: %d, NumChunks: %d", params.ChunkLength, params.NumChunks)
		commits, chunks, err := p.EncodeAndProve(data, params)
		assert.NoError(t, err)
		blobChunks[i] = chunks

		// populate blob header
		blobHeaders[i] = &core.BlobHeader{
			BlobCommitments: commits,
			QuorumInfos:     []*core.BlobQuorumInfo{blobQuorumInfo},
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
	batchHeader, blobMessagesByOp := makeBatch(t, 200*1024, 50, 80, 100, 1)
	numTotalChunks := 0
	for i := range blobMessagesByOp[opID] {
		numTotalChunks += len(blobMessagesByOp[opID][i].Bundles[0])
	}
	t.Logf("Batch numTotalChunks: %d", numTotalChunks)
	req, totalSize, err := dispatcher.GetStoreChunksRequest(blobMessagesByOp[opID], batchHeader, false)
	fmt.Println("totalSize", totalSize)
	assert.NoError(t, err)
	assert.Equal(t, int64(26214400), totalSize)

	timer := time.Now()
	reply, err := server.StoreChunks(context.Background(), req)
	t.Log("StoreChunks took", time.Since(timer))
	assert.NoError(t, err)
	assert.NotNil(t, reply.GetSignature())
}
