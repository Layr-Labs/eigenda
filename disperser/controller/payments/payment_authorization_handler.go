package payments

import (
	"context"
	"fmt"

	clients "github.com/Layr-Labs/eigenda/api/clients/v2"
	pb "github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentAuthorizationHandler struct {
	logger    logging.Logger
	kmsSigner *clients.KMSSigner
}

// NewPaymentAuthorizationHandler creates a new PaymentAuthorizationHandler.
func NewPaymentAuthorizationHandler(
	logger logging.Logger,
	kmsSigner *clients.KMSSigner,
) *PaymentAuthorizationHandler {
	return &PaymentAuthorizationHandler{
		logger:    logger,
		kmsSigner: kmsSigner,
	}
}

// AuthorizePayment processes a payment authorization request.
func (h *PaymentAuthorizationHandler) AuthorizePayment(
	ctx context.Context,
	request *pb.AuthorizePaymentRequest,
) (*pb.AuthorizePaymentReply, error) {
	// Sign the request with the disperser's KMS key to prove this controller authorized it
	if h.kmsSigner != nil {
		// Hash the request first
		hash, err := hashing.HashAuthorizePaymentRequest(request)
		if err != nil {
			h.logger.Error("Failed to hash AuthorizePaymentRequest", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to hash request: %v", err)
		}

		// Sign the hash
		signature, err := h.kmsSigner.SignHash(ctx, hash)
		if err != nil {
			h.logger.Error("Failed to sign AuthorizePaymentRequest", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to sign request: %v", err)
		}

		// Store the signature in the request for audit/logging purposes
		request.DisperserSignature = signature
		h.logger.Debug("Signed AuthorizePaymentRequest", "signatureLength", len(signature), "hashLength", len(hash))
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

// newInsufficientBalanceError creates a structured gRPC error for insufficient balance
// with detailed metadata about the account balance and required cost.
func (h *PaymentAuthorizationHandler) newInsufficientBalanceError(accountID string, currentBalance, requiredCost uint64) error {
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
		h.logger.Error("failed to add error details", "error", err)
		return st.Err()
	}

	return st.Err()
}
