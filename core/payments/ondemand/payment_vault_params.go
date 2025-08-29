package ondemand

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/core/payments/vault"
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
	paymentVault *vault.PaymentVault,
	logger logging.Logger,
) (*PaymentVaultParams, error) {
	blockNumber, err := paymentVault.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current block number: %w", err)
	}

	globalSymbolsPerSecond, err := paymentVault.GetGlobalSymbolsPerSecond(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("get global symbols per second: %w", err)
	}

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("get min num symbols: %w", err)
	}

	pricePerSymbol, err := paymentVault.GetPricePerSymbol(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("get price per symbol: %w", err)
	}

	return &PaymentVaultParams{
		GlobalSymbolsPerSecond: globalSymbolsPerSecond,
		MinNumSymbols:          minNumSymbols,
		PricePerSymbol:         pricePerSymbol,
	}, nil
}
