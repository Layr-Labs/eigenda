package meterer

import (
	"context"
	"math/big"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

const MinNumBins int32 = 3

// OffchainStore defines the interface for storing and retrieving offchain data
type OffchainStore interface {
	// UpdateReservationBin updates the reservation bin for an account and period
	UpdateReservationBin(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64, size uint64) (uint64, error)

	// UpdateGlobalBin updates the global bin for a period
	UpdateGlobalBin(ctx context.Context, reservationPeriod uint64, size uint64) (uint64, error)

	// AddOnDemandPayment adds a new on-demand payment
	AddOnDemandPayment(ctx context.Context, paymentMetadata core.PaymentMetadata, paymentCharged *big.Int) (*big.Int, error)

	// RollbackOnDemandPayment rolls back an on-demand payment
	RollbackOnDemandPayment(ctx context.Context, accountID gethcommon.Address, newPayment, oldPayment *big.Int) error

	// GetPeriodRecords gets period records for an account
	GetPeriodRecords(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64) ([MinNumBins]*pb.PeriodRecord, error)

	// GetLargestCumulativePayment gets the largest cumulative payment for an account
	GetLargestCumulativePayment(ctx context.Context, accountID gethcommon.Address) (*big.Int, error)

	// GetGlobalBinUsage gets the global bin usage for a period
	GetGlobalBinUsage(ctx context.Context, reservationPeriod uint64) (uint64, error)

	// GetReservationBin gets the reservation bin for an account and period
	GetReservationBin(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64) (uint64, error)

	// GetOnDemandPayment gets the on-demand payment for an account
	GetOnDemandPayment(ctx context.Context, accountID gethcommon.Address) (*big.Int, error)

	// GetGlobalBin gets the global bin for a period
	GetGlobalBin(ctx context.Context, reservationPeriod uint64) (uint64, error)

	// Destroy shuts down and permanently deletes all data in the store
	Destroy() error
}
