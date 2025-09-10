package payments

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"

	pb "github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/api/hashing"
	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentAuthorizationHandler struct {
	keyID      string
	publicKey  *ecdsa.PublicKey
	keyManager *kms.Client
}

// NewPaymentAuthorizationHandler creates a new PaymentAuthorizationHandler.
func NewPaymentAuthorizationHandler(
	ctx context.Context,
	region string,
	endpoint string,
	keyID string,
) (*PaymentAuthorizationHandler, error) {

	// Load the AWS SDK configuration, which will automatically detect credentials
	// from environment variables, IAM roles, or AWS config files
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var keyManager *kms.Client
	if endpoint != "" {
		keyManager = kms.New(kms.Options{
			Region:       region,
			BaseEndpoint: aws.String(endpoint),
		})
	} else {
		keyManager = kms.NewFromConfig(cfg)
	}

	publicKey, err := aws2.LoadPublicKeyKMS(ctx, keyManager, keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ecdsa public key: %w", err)
	}

	return &PaymentAuthorizationHandler{
		keyID:      keyID,
		publicKey:  publicKey,
		keyManager: keyManager,
	}, nil
}

// Processes a payment authorization request. Verifies the signature, and checks whether the payment is valid.
func (h *PaymentAuthorizationHandler) AuthorizePayment(
	ctx context.Context,
	request *pb.AuthorizePaymentRequest,
) (*pb.AuthorizePaymentReply, error) {
	if err := h.verifyDisperserSignature(request); err != nil {
		return nil, err
	}

	// TODO(litt3): Implement actual payment authorization logic
	return nil, errors.New("Payment authorization not implemented")
}

// Verifies the disperser's signature on the payment authorization request.
func (h *PaymentAuthorizationHandler) verifyDisperserSignature(request *pb.AuthorizePaymentRequest) error {
	if len(request.DisperserSignature) == 0 {
		return status.Errorf(codes.Unauthenticated, "disperser signature is required")
	}

	requestHash, err := hashing.HashAuthorizePaymentRequest(request)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to hash request: %v", err)
	}

	if len(request.DisperserSignature) != 65 {
		return status.Errorf(codes.Unauthenticated, "invalid disperser signature length")
	}

	// Remove the recovery ID (last byte) for verification
	valid := crypto.VerifySignature(crypto.FromECDSAPub(h.publicKey), requestHash, request.GetDisperserSignature()[:64])
	if !valid {
		return status.Errorf(codes.Unauthenticated, "invalid disperser signature")
	}

	return nil
}
