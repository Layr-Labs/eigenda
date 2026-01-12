package hashing

import (
	"fmt"

	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"golang.org/x/crypto/sha3"
)

// This file contains code for hashing gRPC messages that are sent to the relay.

// RelayGetChunksRequestDomain is the domain for hashing GetChunksRequest messages (i.e. this string
// is added to the digest before hashing the message). This makes it difficult for an attacker to create a
// different type of object that has the same hash as a GetChunksRequest.
const RelayGetChunksRequestDomain = "relay.GetChunksRequest"

// RelayGetValidatorChunksRequestDomain is the domain for hashing GetValidatorChunksRequest messages.
const RelayGetValidatorChunksRequestDomain = "relay.GetValidatorChunksRequest"

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
			err = hashByteArray(hasher, getByIndex.GetBlobKey())
			if err != nil {
				return nil, fmt.Errorf("failed to hash blob key: %w", err)
			}
			err = hashUint32Array(hasher, getByIndex.GetChunkIndices())
			if err != nil {
				return nil, fmt.Errorf("failed to hash ChunkIndices: %w", err)
			}
		} else if chunkRequest.GetByRange() != nil {
			getByRange := chunkRequest.GetByRange()
			hashChar(hasher, 'r')
			err = hashByteArray(hasher, getByRange.GetBlobKey())
			if err != nil {
				return nil, fmt.Errorf("failed to hash blob key: %w", err)
			}
			hashUint32(hasher, getByRange.GetStartIndex())
			hashUint32(hasher, getByRange.GetEndIndex())
		}
	}

	return hasher.Sum(nil), nil
}

// Hashes the given GetValidatorChunksRequest.
func HashGetValidatorChunksRequest(request *pb.GetValidatorChunksRequest) ([]byte, error) {
	hasher := sha3.NewLegacyKeccak256()

	hasher.Write([]byte(RelayGetValidatorChunksRequestDomain))

	err := hashByteArray(hasher, request.GetValidatorId())
	if err != nil {
		return nil, fmt.Errorf("hash validator ID: %w", err)
	}
	err = hashByteArray(hasher, request.GetBlobKey())
	if err != nil {
		return nil, fmt.Errorf("hash blob key: %w", err)
	}
	hashUint32(hasher, request.GetTimestamp())

	return hasher.Sum(nil), nil
}
