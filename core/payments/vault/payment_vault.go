package vault

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
	bindings "github.com/Layr-Labs/eigenda/contracts/bindings/PaymentVault"
	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Provides access to PaymentVault contract
type paymentVault struct {
	logger              logging.Logger
	paymentVaultBinding *bindings.ContractPaymentVault
}

var _ payments.PaymentVault = &paymentVault{}

// NewPaymentVault creates a new PaymentVault instance
func NewPaymentVault(
	logger logging.Logger,
	ethClient common.EthClient,
	paymentVaultAddress gethcommon.Address,
) (payments.PaymentVault, error) {
	if ethClient == nil {
		return nil, errors.New("ethClient cannot be nil")
	}

	paymentVaultBinding, err := bindings.NewContractPaymentVault(paymentVaultAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("new contract payment vault: %w", err)
	}

	return &paymentVault{
		logger:              logger,
		paymentVaultBinding: paymentVaultBinding,
	}, nil
}

// Retrieves total deposit information for multiple accounts
func (pv *paymentVault) GetTotalDeposits(
	ctx context.Context,
	accountIDs []gethcommon.Address,
) ([]*big.Int, error) {
	totalDeposits, err := pv.paymentVaultBinding.GetOnDemandTotalDeposits(&bind.CallOpts{Context: ctx}, accountIDs)
	if err != nil {
		return nil, fmt.Errorf("get on demand total deposits eth call: %w", err)
	}
	return totalDeposits, nil
}

// Retrieves total deposit information for a single account
func (pv *paymentVault) GetTotalDeposit(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	onDemandPayment, err := pv.paymentVaultBinding.GetOnDemandTotalDeposit(&bind.CallOpts{Context: ctx}, accountID)
	if err != nil {
		return nil, fmt.Errorf("get on demand total deposit for account %v eth call: %w", accountID.Hex(), err)
	}
	return onDemandPayment, nil
}

// Retrieves the global symbols per second parameter
func (pv *paymentVault) GetGlobalSymbolsPerSecond(ctx context.Context) (uint64, error) {
	globalSymbolsPerSecond, err := pv.paymentVaultBinding.GlobalSymbolsPerPeriod(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("global symbols per period eth call: %w", err)
	}
	return globalSymbolsPerSecond, nil
}

// Retrieves the minimum number of symbols parameter
func (pv *paymentVault) GetMinNumSymbols(ctx context.Context) (uint32, error) {
	minNumSymbols, err := pv.paymentVaultBinding.MinNumSymbols(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("min num symbols eth call: %w", err)
	}

	if minNumSymbols > math.MaxUint32 {
		return 0, fmt.Errorf("min num symbols > math.MaxUint32: this is nonsensically large, and cannot be handled")
	}

	return uint32(minNumSymbols), nil
}

// GetPricePerSymbol retrieves the price per symbol parameter
func (pv *paymentVault) GetPricePerSymbol(ctx context.Context) (uint64, error) {
	pricePerSymbol, err := pv.paymentVaultBinding.PricePerSymbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, fmt.Errorf("price per symbol eth call: %w", err)
	}
	return pricePerSymbol, nil
}

// Retrieves reservation information for multiple accounts
func (pv *paymentVault) GetReservations(
	ctx context.Context,
	accountIDs []gethcommon.Address,
) ([]*bindings.IPaymentVaultReservation, error) {
	reservations, err := pv.paymentVaultBinding.GetReservations(&bind.CallOpts{Context: ctx}, accountIDs)
	if err != nil {
		return nil, fmt.Errorf("get reservations eth call: %w", err)
	}

	result := make([]*bindings.IPaymentVaultReservation, len(reservations))
	for i, reservation := range reservations {
		// symbolsPerSecond > 0 indicates an active reservation
		if reservation.SymbolsPerSecond == 0 {
			result[i] = nil
			continue
		}

		result[i] = &reservation
	}
	return result, nil
}

// Retrieves reservation information for a single account
func (pv *paymentVault) GetReservation(
	ctx context.Context,
	accountID gethcommon.Address,
) (*bindings.IPaymentVaultReservation, error) {
	reservation, err := pv.paymentVaultBinding.GetReservation(&bind.CallOpts{Context: ctx}, accountID)
	if err != nil {
		return nil, fmt.Errorf("get reservation for account %v eth call: %w", accountID.Hex(), err)
	}

	if reservation.SymbolsPerSecond == 0 {
		return nil, nil
	}

	return &reservation, nil
}
