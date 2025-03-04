package hashing

import (
	"fmt"

	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"golang.org/x/crypto/sha3"
)

// This file contains code for hashing gRPC messages that are sent to the relay.

// HashGetChunksRequest hashes the given GetChunksRequest.
func HashGetChunksRequest(request *pb.GetChunksRequest) ([]byte, error) {
	hasher := sha3.NewLegacyKeccak256()

	err := hashByteArray(hasher, request.GetOperatorId())
	if err != nil {
		return nil, fmt.Errorf("failed to hash operator ID: %w", err)
	}
	err = hashLength(hasher, request.GetChunkRequests())
	if err != nil {
		return nil, fmt.Errorf("failed to hash GetChunkRequests length: %w", err)
	}
	for _, chunkRequest := range request.GetChunkRequests() {
		if chunkRequest.GetByIndex() != nil {
			getByIndex := chunkRequest.GetByIndex()
			hashChar(hasher, 'i')
			err = hashByteArray(hasher, getByIndex.BlobKey)
			if err != nil {
				return nil, fmt.Errorf("failed to hash blob key: %w", err)
			}
			err = hashUint32Array(hasher, getByIndex.ChunkIndices)
			if err != nil {
				return nil, fmt.Errorf("failed to hash ChunkIndices: %w", err)
			}
		} else {
			getByRange := chunkRequest.GetByRange()
			hashChar(hasher, 'r')
			err = hashByteArray(hasher, getByRange.BlobKey)
			if err != nil {
				return nil, fmt.Errorf("failed to hash blob key: %w", err)
			}
			hashUint32(hasher, getByRange.StartIndex)
			hashUint32(hasher, getByRange.EndIndex)
		}
	}

	return hasher.Sum(nil), nil
}
