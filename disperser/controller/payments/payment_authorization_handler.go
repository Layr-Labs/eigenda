//nolint:wrapcheck // Directly returning errors from the api package is the correct pattern
package payments

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/Layr-Labs/eigenda/api"
	common "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	"github.com/Layr-Labs/eigenda/api/grpc/controller"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/ethereum/go-ethereum/crypto"
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
) (*controller.AuthorizePaymentResponse, error) {
	if len(clientSignature) != 65 {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("signature length is unexpected: %d", len(clientSignature)))
	}

	coreHeader, err := core.BlobHeaderFromProtobuf(blobHeader)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("invalid blob header: %v", err))
	}

	blobKey, err := coreHeader.BlobKey()
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to compute blob key: %v", err))
	}

	signerPubkey, err := crypto.SigToPub(blobKey[:], clientSignature)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to recover public key from signature: %v", err))
	}

	accountID := coreHeader.PaymentMetadata.AccountID
	signerAddress := crypto.PubkeyToAddress(*signerPubkey)

	if accountID.Cmp(signerAddress) != 0 {
		return nil, api.NewErrorUnauthenticated(fmt.Sprintf("signature %s doesn't match with provided account %s",
			hex.EncodeToString(clientSignature), accountID.Hex()))
	}

	// TODO(litt3): Implement actual payment authorization logic
	return nil, api.NewErrorInternal("Payment authorization not implemented")
}
