package vault

import (
	"context"
	"math/big"

	bindings "github.com/Layr-Labs/eigenda/contracts/bindings/v2/PaymentVault"
	"github.com/Layr-Labs/eigenda/core/payments"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// TestPaymentVault is a test implementation of the PaymentVault interface
type TestPaymentVault struct {
	// Storage for individual account deposits
	totalDeposits map[gethcommon.Address]*big.Int

	// Global parameters
	globalSymbolsPerSecond uint64
	minNumSymbols          uint64
	PricePerSymbol         uint64
}

var _ payments.PaymentVault = &TestPaymentVault{}

func NewTestPaymentVault() *TestPaymentVault {
	return &TestPaymentVault{
		totalDeposits:          make(map[gethcommon.Address]*big.Int),
		globalSymbolsPerSecond: 1000,
		minNumSymbols:          1,
		PricePerSymbol:         100,
	}
}

func (t *TestPaymentVault) SetTotalDeposit(account gethcommon.Address, amount *big.Int) {
	if amount == nil {
		delete(t.totalDeposits, account)
	} else {
		t.totalDeposits[account] = new(big.Int).Set(amount)
	}
}

func (t *TestPaymentVault) SetGlobalSymbolsPerSecond(value uint64) {
	t.globalSymbolsPerSecond = value
}

func (t *TestPaymentVault) SetMinNumSymbols(value uint64) {
	t.minNumSymbols = value
}

func (t *TestPaymentVault) SetPricePerSymbol(value uint64) {
	t.PricePerSymbol = value
}

func (t *TestPaymentVault) GetTotalDeposits(ctx context.Context, accountIDs []gethcommon.Address) ([]*big.Int, error) {
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

func (t *TestPaymentVault) GetTotalDeposit(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	if deposit, exists := t.totalDeposits[accountID]; exists {
		return new(big.Int).Set(deposit), nil
	}
	return big.NewInt(0), nil
}

func (t *TestPaymentVault) GetGlobalSymbolsPerSecond(ctx context.Context) (uint64, error) {
	return t.globalSymbolsPerSecond, nil
}

func (t *TestPaymentVault) GetMinNumSymbols(ctx context.Context) (uint64, error) {
	return t.minNumSymbols, nil
}

func (t *TestPaymentVault) GetPricePerSymbol(ctx context.Context) (uint64, error) {
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

func (t *TestPaymentVault) GetReservation(
	ctx context.Context,
	accountID gethcommon.Address,
) (*bindings.IPaymentVaultReservation, error) {
	if res, exists := t.reservations[accountID]; exists {
		return res, nil
	}
	return nil, nil
}
