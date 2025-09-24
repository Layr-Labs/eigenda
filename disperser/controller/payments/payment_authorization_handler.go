//nolint:wrapcheck // Directly returning errors from the api package is the correct pattern
package payments

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	grpccommon "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	"github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand/ondemandvalidation"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/Layr-Labs/eigenda/core/payments/reservation/reservationvalidation"
	core "github.com/Layr-Labs/eigenda/core/v2"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Handles payment authorization requests received from API servers.
type PaymentAuthorizationHandler struct {
	onDemandMeterer      *meterer.OnDemandMeterer
	onDemandValidator    *ondemandvalidation.OnDemandPaymentValidator
	reservationValidator *reservationvalidation.ReservationPaymentValidator
}

// Panics if construction fails: we cannot operate if we cannot handle payments
func NewPaymentAuthorizationHandler(
	onDemandMeterer *meterer.OnDemandMeterer,
	onDemandValidator *ondemandvalidation.OnDemandPaymentValidator,
	reservationValidator *reservationvalidation.ReservationPaymentValidator,
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
	blobHeader *grpccommon.BlobHeader,
	clientSignature []byte,
	probe *common.SequenceProbe,
) (*controller.AuthorizePaymentResponse, error) {
	probe.SetStage("request_validation")

	if len(clientSignature) != 65 {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("signature length %d is unexpected, signature: %s",
			len(clientSignature), hex.EncodeToString(clientSignature)))
	}

	coreHeader, err := core.BlobHeaderFromProtobuf(blobHeader)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf(
			"invalid blob header: %v, blobHeader: %s", err, blobHeader.String()))
	}

	blobKey, err := coreHeader.BlobKey()
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf(
			"failed to compute blob key: %v, blobHeader: %s", err, blobHeader.String()))
	}

	probe.SetStage("client_signature_verification")

	signerPubkey, err := crypto.SigToPub(blobKey[:], clientSignature)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf(
			"failed to recover public key from signature: %v, accountID: %s, signature: %s, blobKey: %s",
			err, coreHeader.PaymentMetadata.AccountID.Hex(),
			hex.EncodeToString(clientSignature), hex.EncodeToString(blobKey[:])))
	}

	accountID := coreHeader.PaymentMetadata.AccountID
	signerAddress := crypto.PubkeyToAddress(*signerPubkey)

	if accountID.Cmp(signerAddress) != 0 {
		return nil, api.NewErrorUnauthenticated(fmt.Sprintf(
			"signature %s doesn't match provided account, signerAddress: %s, accountID: %s",
			hex.EncodeToString(clientSignature), signerAddress.Hex(), accountID.Hex()))
	}

	symbolCount := uint32(coreHeader.BlobCommitments.Length)

	if coreHeader.PaymentMetadata.IsOnDemand() {
		err = h.authorizeOnDemandPayment(
			ctx, coreHeader.PaymentMetadata.AccountID, symbolCount, coreHeader.QuorumNumbers, probe)
	} else {
		dispersalTime := time.Unix(0, coreHeader.PaymentMetadata.Timestamp)
		err = h.authorizeReservationPayment(
			ctx, coreHeader.PaymentMetadata.AccountID, symbolCount, coreHeader.QuorumNumbers, dispersalTime, probe)
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
	probe *common.SequenceProbe,
) error {
	probe.SetStage("global_meter_check")
	reservation, err := h.onDemandMeterer.MeterDispersal(symbolCount)
	if err != nil {
		return api.NewErrorResourceExhausted(fmt.Sprintf("global rate limit exceeded: %v", err))
	}

	probe.SetStage("on_demand_validation")
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

	return api.NewErrorInternal(fmt.Sprintf(
		"on-demand payment validation failed for account %s, symbolCount: %d, quorums: %v: %v",
		accountID.Hex(), symbolCount, quorumNumbers, err))

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
	probe *common.SequenceProbe,
) error {
	probe.SetStage("reservation_validation")

	success, err := h.reservationValidator.Debit(ctx, accountID, symbolCount, quorumNumbers, dispersalTime)
	if success {
		return nil
	}
	if err == nil {
		return api.NewErrorPermissionDenied(fmt.Sprintf(
			"reservation payment validation failed for account %s: insufficient bandwidth for %d symbols, time: %s",
			accountID.Hex(), symbolCount, dispersalTime.Format(time.RFC3339)))
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

	return api.NewErrorInternal(fmt.Sprintf(
		"reservation payment validation failed for account %s, symbolCount: %d, quorums: %v, time: %s: %v",
		accountID.Hex(), symbolCount, quorumNumbers, dispersalTime.Format(time.RFC3339), err))
}
