package mock

import (
	"context"
	"crypto/ecdsa"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/node/auth"
)

var _ clients.DispersalRequestSigner = &staticRequestSigner{}

// StaticRequestSigner is a DispersalRequestSigner that signs requests with a static key (i.e. it doesn't use AWS KMS).
// Useful for testing.
type staticRequestSigner struct {
	key *ecdsa.PrivateKey
}

func NewStaticRequestSigner(key *ecdsa.PrivateKey) clients.DispersalRequestSigner {
	return &staticRequestSigner{
		key: key,
	}
}

func (s *staticRequestSigner) SignStoreChunksRequest(
	ctx context.Context,
	request *validator.StoreChunksRequest) ([]byte, error) {

	return auth.SignStoreChunksRequest(s.key, request)
}
