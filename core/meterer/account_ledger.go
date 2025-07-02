package meterer

import (
	"context"
	"math/big"

	disperser_v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// AccountLedger defines the standard interface for tracking account payment state,
// including both on-chain settings and local usage tracking. This interface abstracts
// the differences between client-side in-memory tracking and server-side persistent storage.
//
// Implementation Patterns:
//   - LocalAccountLedger: In-memory tracking for clients
//   - Future: DatabaseAccountLedger for persistent server-side tracking
//   - Future: DistributedAccountLedger for consensus-based tracking
type AccountLedger interface {
	// Debit attempts to record symbol usage against the ledger, handling both
	// reservation and on-demand payment modes automatically based on available payment methods.
	//
	// The method performs the following operations:
	//   1. Attempts reservation usage first if reservations exist and are active
	//   2. Falls back to on-demand usage if reservation fails or doesn't exist
	//   3. For reservation mode: validates reservations, calculates time periods, checks rate limits
	//   4. For on-demand mode: validates quorums, calculates payment, checks balance
	//   5. Updates the appropriate tracking state atomically
	//
	// Parameters:
	//   ctx: Context for cancellation and timeouts (currently unused in LocalAccountLedger)
	//   accountID: The account identifier (used for logging and future database implementations)
	//   timestampNs: The timestamp in nanoseconds when the usage occurred (used for reservations)
	//   numSymbols: The number of symbols to be charged
	//   quorumNumbers: The list of quorums that will be charged for this usage
	//   params: Payment vault parameters containing quorum configurations
	//
	// Returns:
	//   (*big.Int, nil): The new cumulative payment amount if on-demand was used, or nil if reservation was used
	//   (nil, error): If the usage cannot be accommodated:
	//     - "reservation not found for quorum X": No reservation exists for quorum
	//     - "reservation limit exceeded": Usage would exceed rate limits
	//     - "reservation expired": Current time is outside reservation window
	//     - "quorum X not enabled for on-demand": Requested quorum doesn't support on-demand
	//     - "insufficient ondemand payment": Would exceed available balance
	//     - "quorum config not found for X": Missing payment configuration
	//     - "no payment method available": Neither reservation nor on-demand can handle the usage
	Debit(
		ctx context.Context,
		accountID gethcommon.Address,
		timestampNs int64,
		numSymbols uint64,
		quorumNumbers []core.QuorumID,
		params *PaymentVaultParams,
	) (*big.Int, error)

	// RevertDebit reverts the debit operation, returning the ledger to its previous state.
	// This is useful for handling failed debits or partial usage.
	RevertDebit(
		ctx context.Context,
		accountID gethcommon.Address,
		timestampNs int64,
		numSymbols uint64,
		quorumNumbers []core.QuorumID,
		params *PaymentVaultParams,
		payment *big.Int,
	) error

	// GetAccountStateProtobuf returns the account state as separate protobuf-compatible components
	// optimized for wire transmission. This method avoids creating intermediate Go objects,
	// providing direct serializable representations of each state component.
	//
	// The returned components mirror the structure used in GetPaymentStateForAllQuorums:
	//   1. reservations: Per-quorum reservation settings (on-chain)
	//   2. period_records: Usage history for rate limiting (off-chain)
	//   3. onchain_cumulative_payment: Available on-demand balance (on-chain)
	//   4. cumulative_payment: Consumed on-demand payment (off-chain)
	//
	// Use Cases:
	//   - Wire transmission: Sending account state from disperser to client
	//   - Cross-service communication: Passing state between microservices
	//   - Efficient serialization: Avoiding intermediate object creation
	//   - State synchronization: Periodic sync with minimal overhead
	//
	// Returns:
	//   reservations: map[uint32]*disperser_v2.QuorumReservation - Per-quorum reservation settings
	//   periodRecords: map[uint32]*disperser_v2.PeriodRecords - Usage history per quorum
	//   onchainCumulativePayment: []byte - On-chain payment balance as bytes
	//   cumulativePayment: []byte - Local cumulative payment as bytes
	GetAccountStateProtobuf() (
		reservations map[uint32]*disperser_v2.QuorumReservation,
		periodRecords map[uint32]*disperser_v2.PeriodRecords,
		onchainCumulativePayment []byte,
		cumulativePayment []byte,
	)
}
