package hashing

import (
	"errors"

	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"golang.org/x/crypto/sha3"
)

// This file contains code for hashing gRPC messages that are sent to the relay.

var (
	iByte = []byte{0x69}
	rByte = []byte{0x72}
)

// HashGetChunksRequest hashes the given GetChunksRequest.
func HashGetChunksRequest(request *pb.GetChunksRequest) ([]byte, error) {
	hasher := sha3.NewLegacyKeccak256()

	hasher.Write(request.GetOperatorId())
	for _, chunkRequest := range request.GetChunkRequests() {
		if chunkRequest.GetByIndex() != nil {
			getByIndex := chunkRequest.GetByIndex()
			hasher.Write(iByte)
			hasher.Write(getByIndex.GetBlobKey())
			for _, index := range getByIndex.ChunkIndices {
				hashUint32(hasher, index)
			}
		} else if chunkRequest.GetByRange() != nil {
			getByRange := chunkRequest.GetByRange()
			hasher.Write(rByte)
			hasher.Write(getByRange.GetBlobKey())
			hashUint32(hasher, getByRange.GetStartIndex())
			hashUint32(hasher, getByRange.GetEndIndex())
		} else {
			return nil, errors.New("invalid chunk request: must be either by index or by range")
		}
	}

	return hasher.Sum(nil), nil
}
