package clients

import (
	"context"
	"fmt"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
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
	keyManager *kms.Client
}

// NewRequestSigner creates a new RequestSigner.
func NewRequestSigner(
	region string,
	endpoint string,
	keyID string) RequestSigner {

	keyManager := kms.New(kms.Options{
		Region:       region,
		BaseEndpoint: aws.String(endpoint),
	})

	return &requestSigner{
		keyID:      keyID,
		keyManager: keyManager,
	}
}

func (s *requestSigner) SignStoreChunksRequest(ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {
	hash := auth.HashStoreChunksRequest(request)

	signOutput, err := s.keyManager.Sign(
		ctx,
		&kms.SignInput{
			KeyId:            aws.String(s.keyID),
			Message:          hash,
			SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
			MessageType:      types.MessageTypeDigest,
		})

	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	return signOutput.Signature, nil
}
