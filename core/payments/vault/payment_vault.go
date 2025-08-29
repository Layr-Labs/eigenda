package vault

import (
	"context"
	"errors"
	"fmt"
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

func (pv *paymentVault) GetTotalDeposit(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	if pv.paymentVaultBinding == nil {
		return nil, errors.New("payment vault not deployed")
	}
	onDemandPayment, err := pv.paymentVaultBinding.GetOnDemandTotalDeposit(&bind.CallOpts{
		Context: ctx,
	}, accountID)
	if err != nil {
		return nil, err
	}
	if onDemandPayment.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("ondemand payment does not exist for given account")
	}
	return onDemandPayment, nil
}

// GetGlobalSymbolsPerSecond retrieves the global symbols per second parameter
func (pv *paymentVault) GetGlobalSymbolsPerSecond(ctx context.Context) (uint64, error) {
	if pv.paymentVaultBinding == nil {
		return 0, errors.New("payment vault not deployed")
	}
	globalSymbolsPerSecond, err := pv.paymentVaultBinding.GlobalSymbolsPerPeriod(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, err
	}
	return globalSymbolsPerSecond, nil
}

// GetGlobalRatePeriodInterval retrieves the global rate period interval parameter
func (pv *paymentVault) GetGlobalRatePeriodInterval(ctx context.Context) (uint64, error) {
	if pv.paymentVaultBinding == nil {
		return 0, errors.New("payment vault not deployed")
	}
	globalRateBinInterval, err := pv.paymentVaultBinding.GlobalRatePeriodInterval(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, err
	}
	return globalRateBinInterval, nil
}

// GetMinNumSymbols retrieves the minimum number of symbols parameter
func (pv *paymentVault) GetMinNumSymbols(ctx context.Context) (uint64, error) {
	if pv.paymentVaultBinding == nil {
		return 0, errors.New("payment vault not deployed")
	}
	minNumSymbols, err := pv.paymentVaultBinding.MinNumSymbols(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, err
	}
	return minNumSymbols, nil
}

// GetPricePerSymbol retrieves the price per symbol parameter
func (pv *paymentVault) GetPricePerSymbol(ctx context.Context) (uint64, error) {
	if pv.paymentVaultBinding == nil {
		return 0, errors.New("payment vault not deployed")
	}
	pricePerSymbol, err := pv.paymentVaultBinding.PricePerSymbol(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, err
	}
	return pricePerSymbol, nil
}

// GetReservationWindow retrieves the reservation window parameter
func (pv *paymentVault) GetReservationWindow(ctx context.Context) (uint64, error) {
	if pv.paymentVaultBinding == nil {
		return 0, errors.New("payment vault not deployed")
	}
	reservationWindow, err := pv.paymentVaultBinding.ReservationPeriodInterval(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, err
	}
	return reservationWindow, nil
}
