package clients

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/common/oci"
	"github.com/oracle/oci-go-sdk/v65/keymanagement"
)

var _ DispersalRequestSigner = &ociRequestSigner{}

type ociRequestSigner struct {
	keyOCID          string
	publicKey        *ecdsa.PublicKey
	cryptoClient     keymanagement.KmsCryptoClient
	managementClient keymanagement.KmsManagementClient
}

func (s *ociRequestSigner) SignStoreChunksRequest(
	ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {
	hash, err := hashing.HashStoreChunksRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to hash request: %w", err)
	}

	signature, err := oci.SignKMS(ctx, s.cryptoClient, s.keyOCID, s.publicKey, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	return signature, nil
}
