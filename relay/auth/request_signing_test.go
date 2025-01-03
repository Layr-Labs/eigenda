package auth

import (
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/api/hashing"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
	"testing"
)

func randomGetChunksRequest() *pb.GetChunksRequest {
	requestedChunks := make([]*pb.ChunkRequest, 0)
	requestCount := rand.Intn(10) + 1
	for i := 0; i < requestCount; i++ {

		if rand.Intn(2) == 0 {
			indices := make([]uint32, rand.Intn(10)+1)
			for j := 0; j < len(indices); j++ {
				indices[j] = rand.Uint32()
			}
			requestedChunks = append(requestedChunks, &pb.ChunkRequest{
				Request: &pb.ChunkRequest_ByIndex{
					ByIndex: &pb.ChunkRequestByIndex{
						BlobKey:      tu.RandomBytes(32),
						ChunkIndices: indices,
					},
				},
			})
		} else {
			requestedChunks = append(requestedChunks, &pb.ChunkRequest{
				Request: &pb.ChunkRequest_ByRange{
					ByRange: &pb.ChunkRequestByRange{
						BlobKey:    tu.RandomBytes(32),
						StartIndex: rand.Uint32(),
						EndIndex:   rand.Uint32(),
					},
				},
			})
		}
	}
	return &pb.GetChunksRequest{
		OperatorId:    tu.RandomBytes(32),
		ChunkRequests: requestedChunks,
	}
}

func TestHashGetChunksRequest(t *testing.T) {
	tu.InitializeRandom()

	requestA := randomGetChunksRequest()
	requestB := randomGetChunksRequest()

	// Hashing the same request twice should yield the same hash
	hashA := hashing.HashGetChunksRequest(requestA)
	hashAA := hashing.HashGetChunksRequest(requestA)
	require.Equal(t, hashA, hashAA)

	// Hashing different requests should yield different hashes
	hashB := hashing.HashGetChunksRequest(requestB)
	require.NotEqual(t, hashA, hashB)

	// Adding a signature should not affect the hash
	requestA.OperatorSignature = tu.RandomBytes(32)
	hashAA = hashing.HashGetChunksRequest(requestA)
	require.Equal(t, hashA, hashAA)

	// Changing the requester ID should change the hash
	requestA.OperatorId = tu.RandomBytes(32)
	hashAA = hashing.HashGetChunksRequest(requestA)
	require.NotEqual(t, hashA, hashAA)
}
