package payments

import (
	"context"

	common "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	pb "github.com/Layr-Labs/eigenda/api/grpc/controller"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handles payment authorization requests received from API servers.
type PaymentAuthorizationHandler struct {
}

func NewPaymentAuthorizationHandler() *PaymentAuthorizationHandler {
	return &PaymentAuthorizationHandler{}
}

// Checks whether the payment is valid.
func (h *PaymentAuthorizationHandler) AuthorizePayment(
	ctx context.Context,
	blobHeader *common.BlobHeader,
) (*pb.AuthorizePaymentReply, error) {
	// TODO(litt3): Implement actual payment authorization logic
	return nil, status.Errorf(codes.Internal, "Payment authorization not implemented")
}
