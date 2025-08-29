package vault

import (
	"context"
	"math/big"

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

	// Error injection for testing error paths
	getTotalDepositsErr       error
	getTotalDepositErr        error
	getGlobalSymbolsPerSecErr error
	getMinNumSymbolsErr       error
	getPricePerSymbolErr      error
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

func (t *TestPaymentVault) SetDeposit(account gethcommon.Address, amount *big.Int) {
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

func (t *TestPaymentVault) SetGetTotalDepositsErr(err error) {
	t.getTotalDepositsErr = err
}

func (t *TestPaymentVault) SetGetTotalDepositErr(err error) {
	t.getTotalDepositErr = err
}

func (t *TestPaymentVault) SetGetGlobalSymbolsPerSecErr(err error) {
	t.getGlobalSymbolsPerSecErr = err
}

func (t *TestPaymentVault) SetGetMinNumSymbolsErr(err error) {
	t.getMinNumSymbolsErr = err
}

func (t *TestPaymentVault) SetGetPricePerSymbolErr(err error) {
	t.getPricePerSymbolErr = err
}

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

func (t *TestPaymentVault) GetTotalDeposit(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	if t.getTotalDepositErr != nil {
		return nil, t.getTotalDepositErr
	}

	if deposit, exists := t.totalDeposits[accountID]; exists {
		return new(big.Int).Set(deposit), nil
	}
	return big.NewInt(0), nil
}

func (t *TestPaymentVault) GetGlobalSymbolsPerSecond(ctx context.Context) (uint64, error) {
	if t.getGlobalSymbolsPerSecErr != nil {
		return 0, t.getGlobalSymbolsPerSecErr
	}
	return t.globalSymbolsPerSecond, nil
}

func (t *TestPaymentVault) GetMinNumSymbols(ctx context.Context) (uint64, error) {
	if t.getMinNumSymbolsErr != nil {
		return 0, t.getMinNumSymbolsErr
	}
	return t.minNumSymbols, nil
}

func (t *TestPaymentVault) GetPricePerSymbol(ctx context.Context) (uint64, error) {
	if t.getPricePerSymbolErr != nil {
		return 0, t.getPricePerSymbolErr
	}
	return t.PricePerSymbol, nil
}
