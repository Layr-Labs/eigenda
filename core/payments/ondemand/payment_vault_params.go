package ondemand

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type PaymentVaultParams struct {
	GlobalSymbolsPerSecond uint64
	MinNumSymbols          uint64
	PricePerSymbol         uint64
}

// Gets the global PaymentVault parameters, that govern
func GetPaymentVaultParams(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	logger logging.Logger,
) (*PaymentVaultParams, error) {
	globalSymbolsPerSecond, err := paymentVault.GetGlobalSymbolsPerSecond(ctx)
	if err != nil {
		return nil, fmt.Errorf("get global symbols per second: %w", err)
	}

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("get min num symbols: %w", err)
	}

	pricePerSymbol, err := paymentVault.GetPricePerSymbol(ctx)
	if err != nil {
		return nil, fmt.Errorf("get price per symbol: %w", err)
	}

	return &PaymentVaultParams{
		GlobalSymbolsPerSecond: globalSymbolsPerSecond,
		MinNumSymbols:          minNumSymbols,
		PricePerSymbol:         pricePerSymbol,
	}, nil
}
