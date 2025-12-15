package hashing

import (
	"fmt"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing/v2/serialize"
	"golang.org/x/crypto/sha3"
)

func HashStoreChunksRequest_Canonical(request *grpc.StoreChunksRequest) ([]byte, error) {
	canonicalRequest, err := serialize.SerializeStoreChunksRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize store chunks request: %w", err)
	}
	hasher := sha3.New256()
	_, _ = hasher.Write(canonicalRequest)
	return hasher.Sum(nil), nil
}

func HashStoreChunksRequest_V2_Canonical(request *grpc.StoreChunksRequest) ([]byte, error) {
	canonicalRequest, err := serialize.SerializeStoreChunksRequestV2(request)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize store chunks request: %w", err)
	}
	hasher := sha3.New256()
	_, _ = hasher.Write(canonicalRequest)
	return hasher.Sum(nil), nil
}
