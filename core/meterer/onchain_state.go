package meterer

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
)

// PaymentAccounts (For reservations and on-demand payments)

// OnchainPaymentState is an interface for getting information about the current chain state for payments.
type OnchainPayment interface {
	RefreshOnchainPaymentState(ctx context.Context, tx *eth.Reader) error
	GetActiveReservationByAccount(ctx context.Context, accountID string) (core.ActiveReservation, error)
	GetOnDemandPaymentByAccount(ctx context.Context, accountID string) (core.OnDemandPayment, error)
	GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error)
	GetGlobalSymbolsPerSecond() uint64
	GetMinNumSymbols() uint32
	GetPricePerSymbol() uint32
	GetReservationWindow() uint32
}

var _ OnchainPayment = (*OnchainPaymentState)(nil)

type OnchainPaymentState struct {
	tx *eth.Reader

	ActiveReservations map[string]core.ActiveReservation
	OnDemandPayments   map[string]core.OnDemandPayment

	ReservationsLock sync.RWMutex
	OnDemandLocks    sync.RWMutex

	PaymentVaultParams PaymentVaultParams
}

type PaymentVaultParams struct {
	GlobalSymbolsPerSecond uint64
	MinNumSymbols          uint32
	PricePerSymbol         uint32
	ReservationWindow      uint32
	OnDemandQuorumNumbers  []uint8
}

func NewOnchainPaymentState(ctx context.Context, tx *eth.Reader) (OnchainPaymentState, error) {
	paymentVaultParams, err := GetPaymentVaultParams(ctx, tx)
	if err != nil {
		return OnchainPaymentState{}, err
	}

	return OnchainPaymentState{
		tx:                 tx,
		ActiveReservations: make(map[string]core.ActiveReservation),
		OnDemandPayments:   make(map[string]core.OnDemandPayment),
		PaymentVaultParams: paymentVaultParams,
	}, nil
}

func GetPaymentVaultParams(ctx context.Context, tx *eth.Reader) (PaymentVaultParams, error) {
	blockNumber, err := tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return PaymentVaultParams{}, err
	}

	quorumNumbers, err := tx.GetRequiredQuorumNumbers(ctx, blockNumber)
	if err != nil {
		return PaymentVaultParams{}, err
	}

	globalSymbolsPerSecond, err := tx.GetGlobalSymbolsPerSecond(ctx)
	if err != nil {
		return PaymentVaultParams{}, err
	}

	minNumSymbols, err := tx.GetMinNumSymbols(ctx)
	if err != nil {
		return PaymentVaultParams{}, err
	}

	pricePerSymbol, err := tx.GetPricePerSymbol(ctx)
	if err != nil {
		return PaymentVaultParams{}, err
	}

	reservationWindow, err := tx.GetReservationWindow(ctx)
	if err != nil {
		return PaymentVaultParams{}, err
	}

	return PaymentVaultParams{
		OnDemandQuorumNumbers:  quorumNumbers,
		GlobalSymbolsPerSecond: globalSymbolsPerSecond,
		MinNumSymbols:          minNumSymbols,
		PricePerSymbol:         pricePerSymbol,
		ReservationWindow:      reservationWindow,
	}, nil
}

// RefreshOnchainPaymentState returns the current onchain payment state
func (pcs *OnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context, tx *eth.Reader) error {
	paymentVaultParams, err := GetPaymentVaultParams(ctx, tx)
	if err != nil {
		return err
	}
	// These parameters should be rarely updated, but we refresh them anyway
	pcs.PaymentVaultParams = paymentVaultParams
	return nil
}

// GetActiveReservationByAccount returns a pointer to the active reservation for the given account ID; no writes will be made to the reservation
func (pcs *OnchainPaymentState) GetActiveReservationByAccount(ctx context.Context, accountID string) (core.ActiveReservation, error) {
	if reservation, ok := pcs.ActiveReservations[accountID]; ok {
		return reservation, nil
	}
	res, err := pcs.GetActiveReservationByAccountOnChain(ctx, accountID)
	if err != nil {
		return core.ActiveReservation{}, err
	}

	pcs.ReservationsLock.Lock()
	pcs.ActiveReservations[accountID] = res
	pcs.ReservationsLock.Unlock()
	return res, nil
}

// GetActiveReservationByAccountOnChain returns on-chain reservation for the given account ID
func (pcs *OnchainPaymentState) GetActiveReservationByAccountOnChain(ctx context.Context, accountID string) (core.ActiveReservation, error) {
	blockNumber, err := pcs.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return core.ActiveReservation{}, err
	}
	res, err := pcs.tx.GetActiveReservationByAccount(ctx, blockNumber, accountID)
	if err != nil {
		return core.ActiveReservation{}, fmt.Errorf("reservation account not found on-chain: %w", err)
	}
	return res, nil
}

// GetOnDemandPaymentByAccount returns a pointer to the on-demand payment for the given account ID; no writes will be made to the payment
func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID string) (core.OnDemandPayment, error) {
	if payment, ok := pcs.OnDemandPayments[accountID]; ok {
		return payment, nil
	}
	res, err := pcs.GetOnDemandPaymentByAccountOnChain(ctx, accountID)
	if err != nil {
		return core.OnDemandPayment{}, err
	}

	pcs.OnDemandLocks.Lock()
	pcs.OnDemandPayments[accountID] = res
	pcs.OnDemandLocks.Unlock()
	return res, nil
}

func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccountOnChain(ctx context.Context, accountID string) (core.OnDemandPayment, error) {
	blockNumber, err := pcs.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return core.OnDemandPayment{}, err
	}
	res, err := pcs.tx.GetOnDemandPaymentByAccount(ctx, blockNumber, accountID)
	if err != nil {
		return core.OnDemandPayment{}, fmt.Errorf("on-demand not found on-chain: %w", err)
	}
	return res, nil
}

func (pcs *OnchainPaymentState) GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error) {
	blockNumber, err := pcs.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	return pcs.tx.GetRequiredQuorumNumbers(ctx, blockNumber)
}

func (pcs *OnchainPaymentState) GetGlobalSymbolsPerSecond() uint64 {
	return pcs.PaymentVaultParams.GlobalSymbolsPerSecond
}

func (pcs *OnchainPaymentState) GetMinNumSymbols() uint32 {
	return pcs.PaymentVaultParams.MinNumSymbols
}

func (pcs *OnchainPaymentState) GetPricePerSymbol() uint32 {
	return pcs.PaymentVaultParams.PricePerSymbol
}

func (pcs *OnchainPaymentState) GetReservationWindow() uint32 {
	return pcs.PaymentVaultParams.ReservationWindow
}
