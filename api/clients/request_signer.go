package clients

import (
	"context"
	"crypto/ecdsa"
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

	getPublicKeyOutput, err := keyManager.GetPublicKey(ctx, &kms.GetPublicKeyInput{
		KeyId: aws.String(keyID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	publicKey, err := auth.ParseKMSPublicKey(getPublicKeyOutput.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &requestSigner{
		keyID:      keyID,
		publicKey:  publicKey,
		keyManager: keyManager,
	}, nil
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

	signature, err := auth.ParseKMSSignature(s.publicKey, hash, signOutput.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signature: %w", err)
	}

	return signature, nil
}
