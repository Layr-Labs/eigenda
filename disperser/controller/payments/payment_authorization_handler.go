package payments

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/api/hashing"
	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/disperser/controller/metrics"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handles payment authorization requests received from API servers
type PaymentAuthorizationHandler struct {
	metrics    *metrics.ServerMetrics
	keyID      string
	publicKey  *ecdsa.PublicKey
	keyManager *kms.Client
}

func NewPaymentAuthorizationHandler(
	ctx context.Context,
	metrics *metrics.ServerMetrics,
	region string,
	endpoint string,
	keyID string,
) (*PaymentAuthorizationHandler, error) {
	var keyManager *kms.Client
	if endpoint != "" {
		keyManager = kms.New(kms.Options{
			Region:       region,
			BaseEndpoint: aws.String(endpoint),
		})
	} else {
		awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
		if err != nil {
			return nil, fmt.Errorf("load AWS config: %w", err)
		}

		keyManager = kms.NewFromConfig(awsConfig)
	}

	publicKey, err := aws2.LoadPublicKeyKMS(ctx, keyManager, keyID)
	if err != nil {
		return nil, fmt.Errorf("get ecdsa public key: %w", err)
	}

	return &PaymentAuthorizationHandler{
		metrics:    metrics,
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
	start := time.Now()

	if err := h.verifyDisperserSignature(request); err != nil {
		h.metrics.ReportAuthorizePaymentSignatureFailure()
		return nil, err
	}

	h.metrics.ReportAuthorizePaymentSignatureLatency(time.Since(start))

	// TODO(litt3): Implement actual payment authorization logic
	if true {
		h.metrics.ReportAuthorizePaymentAuthFailure()
		return nil, status.Errorf(codes.Internal, "Payment authorization not implemented")
	}

	h.metrics.ReportAuthorizePaymentLatency(time.Since(start))
	return &pb.AuthorizePaymentReply{}, nil
}

// Verifies the disperser's signature on the payment authorization request.
func (h *PaymentAuthorizationHandler) verifyDisperserSignature(request *pb.AuthorizePaymentRequest) error {
	requestHash, err := hashing.HashAuthorizePaymentRequest(request)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to hash request: %v", err)
	}

	if len(request.DisperserSignature) != 65 {
		return status.Errorf(
			codes.Unauthenticated, "invalid disperser signature length %d", len(request.DisperserSignature))
	}

	// Remove the recovery ID (last byte) for verification
	valid := crypto.VerifySignature(crypto.FromECDSAPub(h.publicKey), requestHash, request.GetDisperserSignature()[:64])
	if !valid {
		return status.Errorf(codes.Unauthenticated, "invalid disperser signature")
	}

	return nil
}
