package clients

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// RequestSigner encapsulates the logic for signing GetChunks requests.
type RequestSigner interface {
	// SignStoreChunksRequest signs a StoreChunksRequest. Does not modify the request
	// (i.e. it does not insert the signature).
	SignStoreChunksRequest(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error)
}

var _ RequestSigner = &requestSigner{}

type requestSigner struct {
	keyID      string
	publicKey  *ecdsa.PublicKey
	keyManager *kms.Client
}

// NewRequestSigner creates a new RequestSigner.
func NewRequestSigner(
	ctx context.Context,
	region string,
	endpoint string,
	keyID string) (RequestSigner, error) {

	keyManager := kms.New(kms.Options{
		Region:       region,
		BaseEndpoint: aws.String(endpoint),
	})

	key, err := common.LoadPublicKeyKMS(ctx, keyManager, keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ecdsa public key: %w", err)
	}

	return &requestSigner{
		keyID:      keyID,
		publicKey:  key,
		keyManager: keyManager,
	}, nil
}

func (s *requestSigner) SignStoreChunksRequest(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {
	hash := auth.HashStoreChunksRequest(request)

	signature, err := common.SignKMS(ctx, s.keyManager, s.keyID, s.publicKey, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	return signature, nil
}
