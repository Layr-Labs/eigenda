package meterer

import (
	"context"
	"math/big"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

const MinNumBins uint32 = 3

// MeteringStore defines the interface for storage backends
// used to track reservation and payment usage data
type MeteringStore interface {
	// IncrementBinUsages atomically increments the usage for a reservation bin and returns the new value
	// The key is AccountIDAndQuorum, formatted as {AccountID}:{quorumNumber}.
	IncrementBinUsages(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID, reservationPeriods map[core.QuorumID]uint64, sizes map[core.QuorumID]uint64) (map[core.QuorumID]uint64, error)

	// UpdateGlobalBin atomically increments the usage for a global bin and returns the new value
	UpdateGlobalBin(ctx context.Context, reservationPeriod uint64, size uint64) (uint64, error)

	// AddOnDemandPayment records a new on-demand payment and returns the previous payment amount if successful
	AddOnDemandPayment(ctx context.Context, paymentMetadata core.PaymentMetadata, paymentCharged *big.Int) (*big.Int, error)

	// RollbackOnDemandPayment rolls back a payment to the previous value
	RollbackOnDemandPayment(ctx context.Context, accountID gethcommon.Address, newPayment, oldPayment *big.Int) error

	// GetPeriodRecords fetches period records for the given account ID and reservation period across multiple quorums
	// Consider alternative replacement: singular quorumId and reservationPeriod as an input, require multiple calls
	GetPeriodRecords(ctx context.Context, accountID gethcommon.Address, quorumIds []core.QuorumID, reservationPeriods []uint64, numBins uint32) (map[core.QuorumID]*pb.PeriodRecords, error)

	// GetLargestCumulativePayment returns the largest cumulative payment for the given account
	GetLargestCumulativePayment(ctx context.Context, accountID gethcommon.Address) (*big.Int, error)

	// DecrementBinUsages atomically decrements the bin usage for each quorum in quorumNumbers for a specific account and reservation period.
	// The key is AccountIDAndQuorum, formatted as {AccountID}:{quorumNumber}.
	DecrementBinUsages(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID, reservationPeriods map[core.QuorumID]uint64, sizes map[core.QuorumID]uint64) error
}
