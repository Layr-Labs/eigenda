package payment

import (
	"context"
	"math/big"

	disperser_v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core/meterer"
)

// AccountLedger defines the interface for tracking account payment state.
// Handles both reservation-based and on-demand payment modes with automatic fallback.
type AccountLedger interface {
	// Debit records symbol usage against the account using a DebitSlip.
	//
	// Returns cumulative payment for on-demand usage, nil for reservation usage.
	// Returns error if neither payment method can handle the request.
	Debit(
		ctx context.Context,
		slip *DebitSlip,
		params *meterer.PaymentVaultParams,
	) (*big.Int, error)

	// RevertDebit undoes a previous debit operation using a DebitSlip, restoring the account state.
	RevertDebit(
		ctx context.Context,
		slip *DebitSlip,
		params *meterer.PaymentVaultParams,
		previousCumulativePayment *big.Int,
	) error

	// GetAccountStateProtobuf returns account state as protobuf-compatible components
	// for serialization and wire transmission.
	GetAccountStateProtobuf() (
		reservations map[uint32]*disperser_v2.QuorumReservation,
		periodRecords map[uint32]*disperser_v2.PeriodRecords,
		onchainCumulativePayment []byte,
		cumulativePayment []byte,
	)
}
