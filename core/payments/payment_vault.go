package payments

import (
	"context"
	"math/big"

	bindings "github.com/Layr-Labs/eigenda/contracts/bindings/PaymentVault"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Defines the interface for payment vault contract operations
type PaymentVault interface {
	// Retrieves total on-demand deposits (in wei) for multiple accounts.
	// Returns deposits in same order as accountIDs. Zero returned for accounts with no deposits.
	GetTotalDeposits(ctx context.Context, accountIDs []gethcommon.Address) ([]*big.Int, error)

	// Retrieves total on-demand deposits (in wei) for a single account.
	// Returns zero if the account has no deposits.
	GetTotalDeposit(ctx context.Context, accountID gethcommon.Address) (*big.Int, error)

	// Retrieves the global rate limit (symbols per second) for on-demand dispersals.
	GetGlobalSymbolsPerSecond(ctx context.Context) (uint64, error)

	// Retrieves the minimum billable size for all dispersals.
	// Dispersals are rounded up to the nearest multiple of this value for accounting.
	GetMinNumSymbols(ctx context.Context) (uint64, error)

	// Retrieves the price per symbol (in wei) for on-demand payments.
	GetPricePerSymbol(ctx context.Context) (uint64, error)

	// Retrieves reservation information for multiple accounts.
	// Returns reservations in same order as accountIDs. Returns nil for accounts with no reservation.
	GetReservations(ctx context.Context, accountIDs []gethcommon.Address) ([]*bindings.IPaymentVaultReservation, error)

	// Retrieves reservation information for a single account.
	// Returns nil if the account has no reservation.
	GetReservation(ctx context.Context, accountID gethcommon.Address) (*bindings.IPaymentVaultReservation, error)
}
