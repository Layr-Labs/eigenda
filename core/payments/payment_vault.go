package payments

import (
	"context"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Defines the interface for payment vault contract operations
type PaymentVault interface {
	// GetTotalDeposits retrieves on-demand payment information for multiple accounts
	// Returns deposits in same order as accountIDs
	GetTotalDeposits(ctx context.Context, accountIDs []gethcommon.Address) ([]*big.Int, error)

	// GetTotalDeposit retrieves on-demand payment information for a single account
	GetTotalDeposit(ctx context.Context, accountID gethcommon.Address) (*big.Int, error)

	// GetGlobalSymbolsPerSecond retrieves the global symbols per second parameter
	GetGlobalSymbolsPerSecond(ctx context.Context) (uint64, error)

	// GetGlobalRatePeriodInterval retrieves the global rate period interval parameter
	GetGlobalRatePeriodInterval(ctx context.Context) (uint64, error)

	// GetMinNumSymbols retrieves the minimum number of symbols parameter
	GetMinNumSymbols(ctx context.Context) (uint64, error)

	// GetPricePerSymbol retrieves the price per symbol parameter
	GetPricePerSymbol(ctx context.Context) (uint64, error)

	// GetReservationWindow retrieves the reservation window parameter
	GetReservationWindow(ctx context.Context) (uint64, error)
}
