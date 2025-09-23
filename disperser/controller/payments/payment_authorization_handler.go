//nolint:wrapcheck // Directly returning errors from the api package is the correct pattern
package payments

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	common "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	"github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	core "github.com/Layr-Labs/eigenda/core/v2"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Handles payment authorization requests received from API servers.
type PaymentAuthorizationHandler struct {
	onDemandMeterer      *meterer.OnDemandMeterer
	onDemandValidator    *ondemand.OnDemandPaymentValidator
	reservationValidator *reservation.ReservationPaymentValidator
}

// Panics if construction fails: we cannot operate if we cannot handle payments
func NewPaymentAuthorizationHandler(
	onDemandMeterer *meterer.OnDemandMeterer,
	onDemandValidator *ondemand.OnDemandPaymentValidator,
	reservationValidator *reservation.ReservationPaymentValidator,
) *PaymentAuthorizationHandler {
	if onDemandMeterer == nil {
		panic("onDemandMeterer cannot be nil")
	}
	if onDemandValidator == nil {
		panic("onDemandValidator cannot be nil")
	}
	if reservationValidator == nil {
		panic("reservationValidator cannot be nil")
	}

	return &PaymentAuthorizationHandler{
		onDemandMeterer:      onDemandMeterer,
		onDemandValidator:    onDemandValidator,
		reservationValidator: reservationValidator,
	}
}

// Checks whether the payment is valid.
//
// Verifies the following:
// - client signature
// - payment validity
// - global on-demand throughput meter
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

	symbolCount := uint32(coreHeader.BlobCommitments.Length)

	if coreHeader.PaymentMetadata.IsOnDemand() {
		err = h.authorizeOnDemandPayment(
			ctx, coreHeader.PaymentMetadata.AccountID, symbolCount, coreHeader.QuorumNumbers)
	} else {
		dispersalTime := time.Unix(0, coreHeader.PaymentMetadata.Timestamp)
		err = h.authorizeReservationPayment(
			ctx, coreHeader.PaymentMetadata.AccountID, symbolCount, coreHeader.QuorumNumbers, dispersalTime)
	}

	if err != nil {
		return nil, err
	}

	return &controller.AuthorizePaymentResponse{}, nil
}

// Validates an on-demand payment.
//
// Steps:
// 1. Check the actual symbol count against the global rate limiter to enforce global throughput limits
// 2. Validate the payment with the on-demand validator
// 3. If payment validation fails, refund the global meter to avoid counting failed dispersals
func (h *PaymentAuthorizationHandler) authorizeOnDemandPayment(
	ctx context.Context,
	accountID gethcommon.Address,
	symbolCount uint32,
	quorumNumbers []uint8,
) error {
	reservation, err := h.onDemandMeterer.MeterDispersal(symbolCount)
	if err != nil {
		return api.NewErrorResourceExhausted(fmt.Sprintf("global rate limit exceeded: %v", err))
	}

	err = h.onDemandValidator.Debit(ctx, accountID, symbolCount, quorumNumbers)
	if err == nil {
		return nil
	}

	h.onDemandMeterer.CancelDispersal(reservation)

	var insufficientFundsErr *ondemand.InsufficientFundsError
	if errors.As(err, &insufficientFundsErr) {
		return api.NewErrorPermissionDenied(err.Error())
	}
	var quorumNotSupportedErr *ondemand.QuorumNotSupportedError
	if errors.As(err, &quorumNotSupportedErr) {
		return api.NewErrorInvalidArg(err.Error())
	}

	return api.NewErrorInternal(fmt.Sprintf("on-demand payment validation failed: %v", err))

}

// Validates a reservation payment.
//
// Note: No global metering is required for reservations as they are metered individually
func (h *PaymentAuthorizationHandler) authorizeReservationPayment(
	ctx context.Context,
	accountID gethcommon.Address,
	symbolCount uint32,
	quorumNumbers []uint8,
	dispersalTime time.Time,
) error {
	success, err := h.reservationValidator.Debit(ctx, accountID, symbolCount, quorumNumbers, dispersalTime)
	if success {
		return nil
	}
	if err == nil {
		return api.NewErrorPermissionDenied("reservation payment validation failed: insufficient bandwidth")
	}

	var quorumNotPermittedErr *reservation.QuorumNotPermittedError
	if errors.As(err, &quorumNotPermittedErr) {
		return api.NewErrorInvalidArg(err.Error())
	}
	var timeOutOfRangeErr *reservation.TimeOutOfRangeError
	if errors.As(err, &timeOutOfRangeErr) {
		return api.NewErrorInvalidArg(err.Error())
	}
	var timeMovedBackwardErr *ratelimit.TimeMovedBackwardError
	if errors.As(err, &timeMovedBackwardErr) {
		return api.NewErrorInternal(err.Error())
	}

	return api.NewErrorInternal(fmt.Sprintf("reservation payment validation failed: %v", err))
}
