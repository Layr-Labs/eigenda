package hashing

import (
	"fmt"

	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"golang.org/x/crypto/sha3"
)

// This file contains code for hashing gRPC messages that are sent to the relay.

// RelayGetChunksRequestDomain is the domain for hashing StoreChunksRequest messages (i.e. this string
// is added to the digest before hashing the message). This makes it difficult for an attacker to create a
// different type of object that has the same hash as a StoreChunksRequest.
const RelayGetChunksRequestDomain = "relay.GetChunksRequest"

// HashGetChunksRequest hashes the given GetChunksRequest.
func HashGetChunksRequest(request *pb.GetChunksRequest) ([]byte, error) {
	hasher := sha3.NewLegacyKeccak256()

	hasher.Write([]byte(RelayGetChunksRequestDomain))

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
