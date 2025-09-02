package vault

import (
	"context"
	"math/big"

	bindings "github.com/Layr-Labs/eigenda/contracts/bindings/PaymentVault"
	"github.com/Layr-Labs/eigenda/core/payments"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// TestPaymentVault is a test implementation of the PaymentVault interface
type TestPaymentVault struct {
	// Storage for individual account deposits
	totalDeposits map[gethcommon.Address]*big.Int

	// Storage for individual account reservations
	reservations map[gethcommon.Address]*bindings.IPaymentVaultReservation

	// Global parameters
	globalSymbolsPerSecond uint64
	minNumSymbols          uint64
	PricePerSymbol         uint64

	// Error injection for testing error paths
	getTotalDepositsErr       error
	getTotalDepositErr        error
	getGlobalSymbolsPerSecErr error
	getMinNumSymbolsErr       error
	getPricePerSymbolErr      error
}

var _ payments.PaymentVault = &TestPaymentVault{}

// NewTestPaymentVault creates a new test payment vault with default values
func NewTestPaymentVault() *TestPaymentVault {
	return &TestPaymentVault{
		totalDeposits:          make(map[gethcommon.Address]*big.Int),
		reservations:           make(map[gethcommon.Address]*bindings.IPaymentVaultReservation),
		globalSymbolsPerSecond: 1000,
		minNumSymbols:          1,
		PricePerSymbol:         100,
	}
}

// SetDeposit sets the deposit amount for a specific account
func (t *TestPaymentVault) SetDeposit(account gethcommon.Address, amount *big.Int) {
	if amount == nil {
		delete(t.totalDeposits, account)
	} else {
		t.totalDeposits[account] = new(big.Int).Set(amount)
	}
}

// SetGlobalSymbolsPerSecond sets the global symbols per second parameter
func (t *TestPaymentVault) SetGlobalSymbolsPerSecond(value uint64) {
	t.globalSymbolsPerSecond = value
}

// SetMinNumSymbols sets the minimum number of symbols parameter
func (t *TestPaymentVault) SetMinNumSymbols(value uint64) {
	t.minNumSymbols = value
}

// SetPricePerSymbol sets the price per symbol parameter
func (t *TestPaymentVault) SetPricePerSymbol(value uint64) {
	t.PricePerSymbol = value
}

// SetGetTotalDepositsErr sets the error to return from GetTotalDeposits
func (t *TestPaymentVault) SetGetTotalDepositsErr(err error) {
	t.getTotalDepositsErr = err
}

// SetGetTotalDepositErr sets the error to return from GetTotalDeposit
func (t *TestPaymentVault) SetGetTotalDepositErr(err error) {
	t.getTotalDepositErr = err
}

// SetGetGlobalSymbolsPerSecErr sets the error to return from GetGlobalSymbolsPerSecond
func (t *TestPaymentVault) SetGetGlobalSymbolsPerSecErr(err error) {
	t.getGlobalSymbolsPerSecErr = err
}

// SetGetMinNumSymbolsErr sets the error to return from GetMinNumSymbols
func (t *TestPaymentVault) SetGetMinNumSymbolsErr(err error) {
	t.getMinNumSymbolsErr = err
}

// SetGetPricePerSymbolErr sets the error to return from GetPricePerSymbol
func (t *TestPaymentVault) SetGetPricePerSymbolErr(err error) {
	t.getPricePerSymbolErr = err
}

// GetTotalDeposits retrieves on-demand payment information for multiple accounts
func (t *TestPaymentVault) GetTotalDeposits(ctx context.Context, accountIDs []gethcommon.Address) ([]*big.Int, error) {
	if t.getTotalDepositsErr != nil {
		return nil, t.getTotalDepositsErr
	}

	result := make([]*big.Int, len(accountIDs))
	for i, accountID := range accountIDs {
		if deposit, exists := t.totalDeposits[accountID]; exists {
			result[i] = new(big.Int).Set(deposit)
		} else {
			result[i] = big.NewInt(0)
		}
	}
	return result, nil
}

// GetTotalDeposit retrieves on-demand payment information for a single account
func (t *TestPaymentVault) GetTotalDeposit(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	if t.getTotalDepositErr != nil {
		return nil, t.getTotalDepositErr
	}

	if deposit, exists := t.totalDeposits[accountID]; exists {
		return new(big.Int).Set(deposit), nil
	}
	return big.NewInt(0), nil
}

// GetGlobalSymbolsPerSecond retrieves the global symbols per second parameter
func (t *TestPaymentVault) GetGlobalSymbolsPerSecond(ctx context.Context) (uint64, error) {
	if t.getGlobalSymbolsPerSecErr != nil {
		return 0, t.getGlobalSymbolsPerSecErr
	}
	return t.globalSymbolsPerSecond, nil
}

// GetMinNumSymbols retrieves the minimum number of symbols parameter
func (t *TestPaymentVault) GetMinNumSymbols(ctx context.Context) (uint64, error) {
	if t.getMinNumSymbolsErr != nil {
		return 0, t.getMinNumSymbolsErr
	}
	return t.minNumSymbols, nil
}

// GetPricePerSymbol retrieves the price per symbol parameter
func (t *TestPaymentVault) GetPricePerSymbol(ctx context.Context) (uint64, error) {
	if t.getPricePerSymbolErr != nil {
		return 0, t.getPricePerSymbolErr
	}
	return t.PricePerSymbol, nil
}

func (t *TestPaymentVault) SetReservation(account gethcommon.Address, reservation *bindings.IPaymentVaultReservation) {
	if reservation == nil {
		delete(t.reservations, account)
	} else {
		t.reservations[account] = reservation
	}
}

func (t *TestPaymentVault) GetReservations(
	ctx context.Context,
	accountIDs []gethcommon.Address,
) ([]*bindings.IPaymentVaultReservation, error) {
	result := make([]*bindings.IPaymentVaultReservation, len(accountIDs))
	for i, accountID := range accountIDs {
		if reservation, exists := t.reservations[accountID]; exists {
			result[i] = reservation
		} else {
			result[i] = nil
		}
	}
	return result, nil
}

func (t *TestPaymentVault) GetReservation(ctx context.Context, accountID gethcommon.Address) (*bindings.IPaymentVaultReservation, error) {
	if res, exists := t.reservations[accountID]; exists {
		return res, nil
	}
	return nil, nil
}
