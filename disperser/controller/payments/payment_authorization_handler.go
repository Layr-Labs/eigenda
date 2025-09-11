package payments

import (
	"context"
	"encoding/hex"

	common "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	pb "github.com/Layr-Labs/eigenda/api/grpc/controller/v1"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/ethereum/go-ethereum/crypto"
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
//
// First verifies client signature, then verifies that payment is valid
func (h *PaymentAuthorizationHandler) AuthorizePayment(
	ctx context.Context,
	blobHeader *common.BlobHeader,
	clientSignature []byte,
) (*pb.AuthorizePaymentResponse, error) {
	if len(clientSignature) != 65 {
		return nil, status.Errorf(codes.InvalidArgument, "signature length is unexpected: %d", len(clientSignature))
	}

	coreHeader, err := core.BlobHeaderFromProtobuf(blobHeader)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid blob header: %v", err)
	}

	blobKey, err := coreHeader.BlobKey()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to compute blob key: %v", err)
	}

	signerPubkey, err := crypto.SigToPub(blobKey[:], clientSignature)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to recover public key from signature: %v", err)
	}

	accountID := coreHeader.PaymentMetadata.AccountID
	signerAddress := crypto.PubkeyToAddress(*signerPubkey)

	if accountID.Cmp(signerAddress) != 0 {
		return nil, status.Errorf(codes.Unauthenticated, "signature %s doesn't match with provided account %s",
			hex.EncodeToString(clientSignature), accountID.Hex())
	}

	// TODO(litt3): Implement actual payment authorization logic
	return nil, status.Errorf(codes.Internal, "Payment authorization not implemented")
}
