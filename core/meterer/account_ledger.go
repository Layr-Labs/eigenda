package meterer

import (
	"context"
	"math/big"

	disperser_v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// AccountLedger defines the interface for tracking account payment state.
// Handles both reservation-based and on-demand payment modes with automatic fallback.
type AccountLedger interface {
	// Debit records symbol usage against the account.
	//
	// Returns cumulative payment for on-demand usage, nil for reservation usage.
	// Returns error if neither payment method can handle the request.
	Debit(
		ctx context.Context,
		accountID gethcommon.Address,
		timestampNs int64,
		numSymbols uint64,
		quorumNumbers []core.QuorumID,
		params *PaymentVaultParams,
	) (*big.Int, error)

	// RevertDebit undoes a previous debit operation, restoring the account state.
	RevertDebit(
		ctx context.Context,
		accountID gethcommon.Address,
		timestampNs int64,
		numSymbols uint64,
		quorumNumbers []core.QuorumID,
		params *PaymentVaultParams,
		payment *big.Int,
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
