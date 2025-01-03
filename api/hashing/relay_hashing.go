package hashing

import (
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"golang.org/x/crypto/sha3"
)

// This file contains code for hashing gRPC messages that are sent to the relay.

var (
	iByte = []byte{0x69}
	rByte = []byte{0x72}
)

// HashGetChunksRequest hashes the given GetChunksRequest.
func HashGetChunksRequest(request *pb.GetChunksRequest) []byte {
	hasher := sha3.NewLegacyKeccak256()

	hasher.Write(request.GetOperatorId())
	for _, chunkRequest := range request.GetChunkRequests() {
		if chunkRequest.GetByIndex() != nil {
			getByIndex := chunkRequest.GetByIndex()
			hasher.Write(iByte)
			hasher.Write(getByIndex.BlobKey)
			for _, index := range getByIndex.ChunkIndices {
				hashUint32(hasher, index)
			}
		} else {
			getByRange := chunkRequest.GetByRange()
			hasher.Write(rByte)
			hasher.Write(getByRange.BlobKey)
			hashUint32(hasher, getByRange.StartIndex)
			hashUint32(hasher, getByRange.EndIndex)
		}
	}

	return hasher.Sum(nil)
}
