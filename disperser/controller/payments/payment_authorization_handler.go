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
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/controller/metrics"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Handles payment authorization requests received from API servers.
type PaymentAuthorizationHandler struct {
	onDemandMeterer      *meterer.OnDemandMeterer
	onDemandValidator    *ondemand.OnDemandPaymentValidator
	reservationValidator *reservation.ReservationPaymentValidator
	metrics              *metrics.PaymentAuthorizationMetrics
}

// Panics if construction fails: we cannot operate if we cannot handle payments
func NewPaymentAuthorizationHandler(
	onDemandMeterer *meterer.OnDemandMeterer,
	onDemandValidator *ondemand.OnDemandPaymentValidator,
	reservationValidator *reservation.ReservationPaymentValidator,
	metrics *metrics.PaymentAuthorizationMetrics,
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

	handler := &PaymentAuthorizationHandler{
		onDemandMeterer:      onDemandMeterer,
		onDemandValidator:    onDemandValidator,
		reservationValidator: reservationValidator,
		metrics:              metrics,
	}

	metrics.RegisterReservationCacheSize(func() int {
		return reservationValidator.GetCacheSize()
	})
	metrics.RegisterOnDemandCacheSize(func() int {
		return onDemandValidator.GetCacheSize()
	})

	return handler
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
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("signature length is unexpected: %d, signature: %s",
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
		h.metrics.IncrementOnDemandGlobalMeterExhausted()
		return api.NewErrorResourceExhausted(fmt.Sprintf("global rate limit exceeded: %v", err))
	}

	probe.SetStage("on_demand_validation")
	err = h.onDemandValidator.Debit(ctx, accountID, symbolCount, quorumNumbers)
	if err == nil {
		h.metrics.RecordOnDemandPaymentSuccess(symbolCount)
		return nil
	}

	h.onDemandMeterer.CancelDispersal(reservation)

	var insufficientFundsErr *ondemand.InsufficientFundsError
	if errors.As(err, &insufficientFundsErr) {
		h.metrics.IncrementOnDemandInsufficientFunds()
		return api.NewErrorPermissionDenied(err.Error())
	}
	var quorumNotSupportedErr *ondemand.QuorumNotSupportedError
	if errors.As(err, &quorumNotSupportedErr) {
		h.metrics.IncrementOnDemandQuorumNotSupported()
		return api.NewErrorInvalidArg(err.Error())
	}

	h.metrics.IncrementOnDemandUnexpectedErrors()
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
		h.metrics.RecordReservationPaymentSuccess(symbolCount)
		return nil
	}
	if err == nil {
		h.metrics.IncrementReservationInsufficientFunds()
		return api.NewErrorPermissionDenied(fmt.Sprintf(
			"reservation payment validation failed for account %s: insufficient bandwidth for %d symbols, time: %s",
			accountID.Hex(), symbolCount, dispersalTime.Format(time.RFC3339)))
	}

	var quorumNotPermittedErr *reservation.QuorumNotPermittedError
	if errors.As(err, &quorumNotPermittedErr) {
		h.metrics.IncrementReservationQuorumNotPermitted()
		return api.NewErrorInvalidArg(err.Error())
	}
	var timeOutOfRangeErr *reservation.TimeOutOfRangeError
	if errors.As(err, &timeOutOfRangeErr) {
		h.metrics.IncrementReservationTimeOutOfRange()
		return api.NewErrorInvalidArg(err.Error())
	}
	var timeMovedBackwardErr *reservation.TimeMovedBackwardError
	if errors.As(err, &timeMovedBackwardErr) {
		h.metrics.IncrementReservationTimeMovedBackward()
		return api.NewErrorInternal(err.Error())
	}

	h.metrics.IncrementReservationUnexpectedErrors()
	return api.NewErrorInternal(fmt.Sprintf(
		"reservation payment validation failed for account %s, symbolCount: %d, quorums: %v, time: %s: %v",
		accountID.Hex(), symbolCount, quorumNumbers, dispersalTime.Format(time.RFC3339), err))
}
