package auth

import (
	"testing"

	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
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
						BlobKey:      random.RandomBytes(32),
						ChunkIndices: indices,
					},
				},
			})
		} else {
			requestedChunks = append(requestedChunks, &pb.ChunkRequest{
				Request: &pb.ChunkRequest_ByRange{
					ByRange: &pb.ChunkRequestByRange{
						BlobKey:    random.RandomBytes(32),
						StartIndex: rand.Uint32(),
						EndIndex:   rand.Uint32(),
					},
				},
			})
		}
	}
	return &pb.GetChunksRequest{
		OperatorId:    random.RandomBytes(32),
		ChunkRequests: requestedChunks,
	}
}

func TestHashGetChunksRequest(t *testing.T) {
	random.InitializeRandom()

	requestA := randomGetChunksRequest()
	requestB := randomGetChunksRequest()

	// Hashing the same request twice should yield the same hash
	hashA, err := hashing.HashGetChunksRequest(requestA)
	require.NoError(t, err)
	hashAA, err := hashing.HashGetChunksRequest(requestA)
	require.NoError(t, err)
	require.Equal(t, hashA, hashAA)

	// Hashing different requests should yield different hashes
	hashB, err := hashing.HashGetChunksRequest(requestB)
	require.NoError(t, err)
	require.NotEqual(t, hashA, hashB)

	// Adding a signature should not affect the hash
	requestA.OperatorSignature = random.RandomBytes(32)
	hashAA, err = hashing.HashGetChunksRequest(requestA)
	require.NoError(t, err)
	require.Equal(t, hashA, hashAA)

	// Changing the requester ID should change the hash
	requestA.OperatorId = random.RandomBytes(32)
	hashAA, err = hashing.HashGetChunksRequest(requestA)
	require.NoError(t, err)
	require.NotEqual(t, hashA, hashAA)
}
