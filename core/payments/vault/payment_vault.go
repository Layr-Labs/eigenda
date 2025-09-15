package vault

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
	bindings "github.com/Layr-Labs/eigenda/contracts/bindings/v2/PaymentVault"
	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Provides access to PaymentVault contract
type paymentVault struct {
	logger              logging.Logger
	ethClient           common.EthClient
	paymentVaultAddress gethcommon.Address
	paymentVaultBinding *bindings.ContractPaymentVault
}

var _ payments.PaymentVault = &paymentVault{}

func NewPaymentVault(
	logger logging.Logger,
	ethClient common.EthClient,
	paymentVaultAddress gethcommon.Address,
) (payments.PaymentVault, error) {
	if ethClient == nil {
		return nil, errors.New("ethClient cannot be nil")
	}

	return &paymentVault{
		logger:              logger,
		ethClient:           ethClient,
		paymentVaultAddress: paymentVaultAddress,
		paymentVaultBinding: bindings.NewContractPaymentVault(),
	}, nil
}

// Retrieves total deposit information for multiple accounts
func (pv *paymentVault) GetTotalDeposits(
	ctx context.Context,
	accountIDs []gethcommon.Address,
) ([]*big.Int, error) {
	callData, err := pv.paymentVaultBinding.TryPackGetOnDemandTotalDeposits(accountIDs)
	if err != nil {
		return nil, fmt.Errorf("pack GetOnDemandTotalDeposits call: %w", err)
	}

	returnData, err := pv.ethClient.CallContract(ctx, ethereum.CallMsg{
		To:   &pv.paymentVaultAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("get on demand total deposits eth call: %w", err)
	}

	totalDeposits, err := pv.paymentVaultBinding.UnpackGetOnDemandTotalDeposits(returnData)
	if err != nil {
		return nil, fmt.Errorf("unpack GetOnDemandTotalDeposits return data: %w", err)
	}

	return totalDeposits, nil
}

// Retrieves total deposit information for a single account
func (pv *paymentVault) GetTotalDeposit(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	callData, err := pv.paymentVaultBinding.TryPackGetOnDemandTotalDeposit(accountID)
	if err != nil {
		return nil, fmt.Errorf("pack GetOnDemandTotalDeposit call: %w", err)
	}

	returnData, err := pv.ethClient.CallContract(ctx, ethereum.CallMsg{
		To:   &pv.paymentVaultAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("get on demand total deposit for account %v eth call: %w", accountID.Hex(), err)
	}

	onDemandPayment, err := pv.paymentVaultBinding.UnpackGetOnDemandTotalDeposit(returnData)
	if err != nil {
		return nil, fmt.Errorf("unpack GetOnDemandTotalDeposit return data: %w", err)
	}

	return onDemandPayment, nil
}

// Retrieves the global symbols per second parameter
func (pv *paymentVault) GetGlobalSymbolsPerSecond(ctx context.Context) (uint64, error) {
	callData, err := pv.paymentVaultBinding.TryPackGlobalSymbolsPerPeriod()
	if err != nil {
		return 0, fmt.Errorf("pack GlobalSymbolsPerPeriod call: %w", err)
	}

	returnData, err := pv.ethClient.CallContract(ctx, ethereum.CallMsg{
		To:   &pv.paymentVaultAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return 0, fmt.Errorf("global symbols per period eth call: %w", err)
	}

	globalSymbolsPerSecond, err := pv.paymentVaultBinding.UnpackGlobalSymbolsPerPeriod(returnData)
	if err != nil {
		return 0, fmt.Errorf("unpack GlobalSymbolsPerPeriod return data: %w", err)
	}

	return globalSymbolsPerSecond, nil
}

// Retrieves the minimum number of symbols parameter
func (pv *paymentVault) GetMinNumSymbols(ctx context.Context) (uint32, error) {
	callData, err := pv.paymentVaultBinding.TryPackMinNumSymbols()
	if err != nil {
		return 0, fmt.Errorf("pack MinNumSymbols call: %w", err)
	}

	returnData, err := pv.ethClient.CallContract(ctx, ethereum.CallMsg{
		To:   &pv.paymentVaultAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return 0, fmt.Errorf("min num symbols eth call: %w", err)
	}

	minNumSymbols, err := pv.paymentVaultBinding.UnpackMinNumSymbols(returnData)
	if err != nil {
		return 0, fmt.Errorf("unpack MinNumSymbols return data: %w", err)
	}

	if minNumSymbols > math.MaxUint32 {
		return 0, fmt.Errorf("min num symbols > math.MaxUint32: this is nonsensically large, and cannot be handled")
	}

	return uint32(minNumSymbols), nil
}

// GetPricePerSymbol retrieves the price per symbol parameter
func (pv *paymentVault) GetPricePerSymbol(ctx context.Context) (uint64, error) {
	callData, err := pv.paymentVaultBinding.TryPackPricePerSymbol()
	if err != nil {
		return 0, fmt.Errorf("pack PricePerSymbol call: %w", err)
	}

	returnData, err := pv.ethClient.CallContract(ctx, ethereum.CallMsg{
		To:   &pv.paymentVaultAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return 0, fmt.Errorf("price per symbol eth call: %w", err)
	}

	pricePerSymbol, err := pv.paymentVaultBinding.UnpackPricePerSymbol(returnData)
	if err != nil {
		return 0, fmt.Errorf("unpack PricePerSymbol return data: %w", err)
	}

	return pricePerSymbol, nil
}

// Retrieves reservation information for multiple accounts
func (pv *paymentVault) GetReservations(
	ctx context.Context,
	accountIDs []gethcommon.Address,
) ([]*bindings.IPaymentVaultReservation, error) {
	callData, err := pv.paymentVaultBinding.TryPackGetReservations(accountIDs)
	if err != nil {
		return nil, fmt.Errorf("pack GetReservations call: %w", err)
	}

	returnData, err := pv.ethClient.CallContract(ctx, ethereum.CallMsg{
		To:   &pv.paymentVaultAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("get reservations eth call: %w", err)
	}

	reservations, err := pv.paymentVaultBinding.UnpackGetReservations(returnData)
	if err != nil {
		return nil, fmt.Errorf("unpack GetReservations return data: %w", err)
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
	callData, err := pv.paymentVaultBinding.TryPackGetReservation(accountID)
	if err != nil {
		return nil, fmt.Errorf("pack GetReservation call: %w", err)
	}

	returnData, err := pv.ethClient.CallContract(ctx, ethereum.CallMsg{
		To:   &pv.paymentVaultAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("get reservation for account %v eth call: %w", accountID.Hex(), err)
	}

	reservation, err := pv.paymentVaultBinding.UnpackGetReservation(returnData)
	if err != nil {
		return nil, fmt.Errorf("unpack GetReservation return data: %w", err)
	}

	if reservation.SymbolsPerSecond == 0 {
		return nil, nil
	}

	return &reservation, nil
}
