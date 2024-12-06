package meterer

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// PaymentAccounts (For reservations and on-demand payments)

// OnchainPaymentState is an interface for getting information about the current chain state for payments.
type OnchainPayment interface {
	RefreshOnchainPaymentState(ctx context.Context, tx *eth.Reader) error
	GetActiveReservationByAccount(ctx context.Context, accountID gethcommon.Address) (*core.ActiveReservation, error)
	GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error)
	GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error)
	GetGlobalSymbolsPerSecond() uint64
	GetGlobalRateBinInterval() uint64
	GetMinNumSymbols() uint32
	GetPricePerSymbol() uint32
	GetReservationWindow() uint32
}

var _ OnchainPayment = (*OnchainPaymentState)(nil)

type OnchainPaymentState struct {
	tx *eth.Reader

	ActiveReservations map[gethcommon.Address]*core.ActiveReservation
	OnDemandPayments   map[gethcommon.Address]*core.OnDemandPayment

	ReservationsLock sync.RWMutex
	OnDemandLocks    sync.RWMutex

	PaymentVaultParams atomic.Pointer[PaymentVaultParams]
}

type PaymentVaultParams struct {
	GlobalSymbolsPerSecond uint64
	GlobalRateBinInterval  uint64
	MinNumSymbols          uint32
	PricePerSymbol         uint32
	ReservationWindow      uint32
	OnDemandQuorumNumbers  []uint8
}

func NewOnchainPaymentState(ctx context.Context, tx *eth.Reader) (*OnchainPaymentState, error) {
	paymentVaultParams, err := GetPaymentVaultParams(ctx, tx)
	if err != nil {
		return nil, err
	}

	state := OnchainPaymentState{
		tx:                 tx,
		ActiveReservations: make(map[gethcommon.Address]*core.ActiveReservation),
		OnDemandPayments:   make(map[gethcommon.Address]*core.OnDemandPayment),
		PaymentVaultParams: atomic.Pointer[PaymentVaultParams]{},
	}
	state.PaymentVaultParams.Store(paymentVaultParams)

	return &state, nil
}

func GetPaymentVaultParams(ctx context.Context, tx *eth.Reader) (*PaymentVaultParams, error) {
	blockNumber, err := tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	quorumNumbers, err := tx.GetRequiredQuorumNumbers(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	globalSymbolsPerSecond, err := tx.GetGlobalSymbolsPerSecond(ctx)
	if err != nil {
		return nil, err
	}

	minNumSymbols, err := tx.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, err
	}

	pricePerSymbol, err := tx.GetPricePerSymbol(ctx)
	if err != nil {
		return nil, err
	}

	reservationWindow, err := tx.GetReservationWindow(ctx)
	if err != nil {
		return nil, err
	}

	return &PaymentVaultParams{
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
	pcs.PaymentVaultParams.Store(paymentVaultParams)

	pcs.ReservationsLock.Lock()
	accountIDs := make([]gethcommon.Address, 0, len(pcs.ActiveReservations))
	for accountID := range pcs.ActiveReservations {
		accountIDs = append(accountIDs, accountID)
	}

	activeReservations, err := tx.GetActiveReservations(ctx, accountIDs)
	if err != nil {
		return err
	}
	pcs.ActiveReservations = activeReservations
	pcs.ReservationsLock.Unlock()

	pcs.OnDemandLocks.Lock()
	accountIDs = make([]gethcommon.Address, 0, len(pcs.OnDemandPayments))
	for accountID := range pcs.OnDemandPayments {
		accountIDs = append(accountIDs, accountID)
	}

	onDemandPayments, err := tx.GetOnDemandPayments(ctx, accountIDs)
	if err != nil {
		return err
	}
	pcs.OnDemandPayments = onDemandPayments
	pcs.OnDemandLocks.Unlock()

	return nil
}

// GetActiveReservationByAccount returns a pointer to the active reservation for the given account ID; no writes will be made to the reservation
func (pcs *OnchainPaymentState) GetActiveReservationByAccount(ctx context.Context, accountID gethcommon.Address) (*core.ActiveReservation, error) {
	pcs.ReservationsLock.RLock()
	defer pcs.ReservationsLock.RUnlock()
	if reservation, ok := (pcs.ActiveReservations)[accountID]; ok {
		return reservation, nil
	}

	// pulls the chain state
	res, err := pcs.tx.GetActiveReservationByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	pcs.ReservationsLock.Lock()
	(pcs.ActiveReservations)[accountID] = res
	pcs.ReservationsLock.Unlock()
	return res, nil
}

// GetActiveReservationByAccountOnChain returns on-chain reservation for the given account ID
func (pcs *OnchainPaymentState) GetActiveReservationByAccountOnChain(ctx context.Context, accountID gethcommon.Address) (*core.ActiveReservation, error) {
	res, err := pcs.tx.GetActiveReservationByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("reservation account not found on-chain: %w", err)
	}
	return res, nil
}

// GetOnDemandPaymentByAccount returns a pointer to the on-demand payment for the given account ID; no writes will be made to the payment
func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error) {
	pcs.OnDemandLocks.RLock()
	defer pcs.OnDemandLocks.RUnlock()
	if payment, ok := (pcs.OnDemandPayments)[accountID]; ok {
		return payment, nil
	}
	// pulls the chain state
	res, err := pcs.tx.GetOnDemandPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	pcs.OnDemandLocks.Lock()
	(pcs.OnDemandPayments)[accountID] = res
	pcs.OnDemandLocks.Unlock()
	return res, nil
}

func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccountOnChain(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error) {
	res, err := pcs.tx.GetOnDemandPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("on-demand not found on-chain: %w", err)
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
	return pcs.PaymentVaultParams.Load().GlobalSymbolsPerSecond
}

func (pcs *OnchainPaymentState) GetGlobalRateBinInterval() uint64 {
	return pcs.PaymentVaultParams.Load().GlobalRateBinInterval
}

func (pcs *OnchainPaymentState) GetMinNumSymbols() uint32 {
	return pcs.PaymentVaultParams.Load().MinNumSymbols
}

func (pcs *OnchainPaymentState) GetPricePerSymbol() uint32 {
	return pcs.PaymentVaultParams.Load().PricePerSymbol
}

func (pcs *OnchainPaymentState) GetReservationWindow() uint32 {
	return pcs.PaymentVaultParams.Load().ReservationWindow
}
