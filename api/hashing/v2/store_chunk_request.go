package hashing

import (
	"fmt"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing/v2/serialize"
	"golang.org/x/crypto/sha3"
)

// HashStoreChunksRequest hashes the given StoreChunksRequest using the canonical serialization.
func HashStoreChunksRequest(request *grpc.StoreChunksRequest) ([]byte, error) {
	canonicalRequest, err := serialize.SerializeStoreChunksRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize store chunks request: %w", err)
	}
	hasher := sha3.New256()
	_, _ = hasher.Write(canonicalRequest)
	return hasher.Sum(nil), nil
}
