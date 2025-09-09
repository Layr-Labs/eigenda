package payments

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	pb "github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/api/hashing"
	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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

	// TODO: Implement actual payment authorization logic
	// This should include:
	// 1. Validate the client signature against the account's public key
	// 2. Check account balance and reservation status
	// 3. Calculate the cost of the requested dispersal
	// 4. Verify the account has sufficient funds
	// 5. Record the payment authorization for audit purposes

	// Example: Simulate insufficient balance error with structured metadata
	// In real implementation, this would be based on actual balance checks
	simulateInsufficientBalance := false

	if simulateInsufficientBalance {
		accountID := request.GetBlobHeader().GetPaymentHeader().GetAccountId()
		currentBalance := uint64(100) // In production, get actual balance
		requiredCost := uint64(150)   // In production, calculate actual cost

		return nil, h.newInsufficientBalanceError(accountID, currentBalance, requiredCost)
	}

	return &pb.AuthorizePaymentReply{}, nil
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
	valid := crypto.VerifySignature(crypto.FromECDSAPub(h.publicKey), requestHash, request.DisperserSignature[:64])
	if !valid {
		return status.Errorf(codes.Unauthenticated, "invalid disperser signature")
	}

	return nil
}

// newInsufficientBalanceError creates a structured gRPC error for insufficient balance
// with detailed metadata about the account balance and required cost.
func (h *PaymentAuthorizationHandler) newInsufficientBalanceError(
	accountID string,
	currentBalance uint64,
	requiredCost uint64,
) error {
	deficit := uint64(0)
	if requiredCost > currentBalance {
		deficit = requiredCost - currentBalance
	}

	st := status.New(codes.FailedPrecondition, "insufficient balance for blob dispersal")

	// Add structured error details with metadata
	// TODO: make this match the same metadata that is returned from structured payments error.
	// Consider creating a method on the structured payments error that wraps it as a gRPC error.
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Reason: "INSUFFICIENT_BALANCE",
		Domain: "payment",
		Metadata: map[string]string{
			"account_id":      accountID,
			"current_balance": fmt.Sprintf("%d", currentBalance),
			"required_cost":   fmt.Sprintf("%d", requiredCost),
			"deficit":         fmt.Sprintf("%d", deficit),
		},
	})

	if err != nil {
		// If we can't add details, return the basic error
		return st.Err()
	}

	return st.Err()
}
