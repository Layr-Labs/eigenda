package clients

import (
	"context"
	"fmt"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing"
)

// DispersalRequestSigner encapsulates the logic for signing GetChunks requests.
type DispersalRequestSigner interface {
	// SignStoreChunksRequest signs a StoreChunksRequest. Does not modify the request
	// (i.e. it does not insert the signature).
	SignStoreChunksRequest(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error)
}

var _ DispersalRequestSigner = &requestSigner{}

type requestSigner struct {
	*KMSSigner
}

// NewDispersalRequestSigner creates a new DispersalRequestSigner.
func NewDispersalRequestSigner(
	ctx context.Context,
	region string,
	endpoint string,
	keyID string) (DispersalRequestSigner, error) {

	kmsSigner, err := NewKMSSigner(ctx, region, endpoint, keyID)
	if err != nil {
		return nil, fmt.Errorf("create KMS signer: %w", err)
	}

	return &requestSigner{
		KMSSigner: kmsSigner,
	}, nil
}

func (s *requestSigner) SignStoreChunksRequest(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {
	hash, err := hashing.HashStoreChunksRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to hash request: %w", err)
	}

	return s.SignHash(ctx, hash)
}
