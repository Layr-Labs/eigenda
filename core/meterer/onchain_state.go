package meterer

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// PaymentAccounts (For reservations and on-demand payments)

// OnchainPaymentState is an interface for getting information about the current chain state for payments.
type OnchainPayment interface {
	RefreshOnchainPaymentState(ctx context.Context) error
	GetReservedPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.ReservedPayment, error)
	GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error)
	GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error)
	GetGlobalSymbolsPerSecond() uint64
	GetGlobalRatePeriodInterval() uint32
	GetMinNumSymbols() uint32
	GetPricePerSymbol() uint32
	GetReservationWindow() uint32
}

var _ OnchainPayment = (*OnchainPaymentState)(nil)

type OnchainPaymentState struct {
	tx *eth.Reader

	ReservedPayments map[gethcommon.Address]*core.ReservedPayment
	OnDemandPayments map[gethcommon.Address]*core.OnDemandPayment

	ReservationsLock sync.RWMutex
	OnDemandLocks    sync.RWMutex

	PaymentVaultParams atomic.Pointer[PaymentVaultParams]
}

type PaymentVaultParams struct {
	GlobalSymbolsPerSecond   uint64
	GlobalRatePeriodInterval uint32
	MinNumSymbols            uint32
	PricePerSymbol           uint32
	ReservationWindow        uint32
	OnDemandQuorumNumbers    []uint8
}

func NewOnchainPaymentState(ctx context.Context, tx *eth.Reader) (*OnchainPaymentState, error) {
	state := OnchainPaymentState{
		tx:                 tx,
		ReservedPayments:   make(map[gethcommon.Address]*core.ReservedPayment),
		OnDemandPayments:   make(map[gethcommon.Address]*core.OnDemandPayment),
		PaymentVaultParams: atomic.Pointer[PaymentVaultParams]{},
	}

	paymentVaultParams, err := state.GetPaymentVaultParams(ctx)
	if err != nil {
		return nil, err
	}

	state.PaymentVaultParams.Store(paymentVaultParams)

	return &state, nil
}

func (pcs *OnchainPaymentState) GetPaymentVaultParams(ctx context.Context) (*PaymentVaultParams, error) {
	quorumNumbers, err := pcs.GetOnDemandQuorumNumbers(ctx)
	if err != nil {
		return nil, err
	}

	globalSymbolsPerSecond, err := pcs.tx.GetGlobalSymbolsPerSecond(ctx)
	if err != nil {
		return nil, err
	}

	globalRatePeriodInterval, err := pcs.tx.GetGlobalRatePeriodInterval(ctx)
	if err != nil {
		return nil, err
	}

	minNumSymbols, err := pcs.tx.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, err
	}

	pricePerSymbol, err := pcs.tx.GetPricePerSymbol(ctx)
	if err != nil {
		return nil, err
	}

	reservationWindow, err := pcs.tx.GetReservationWindow(ctx)
	if err != nil {
		return nil, err
	}

	return &PaymentVaultParams{
		OnDemandQuorumNumbers:    quorumNumbers,
		GlobalSymbolsPerSecond:   globalSymbolsPerSecond,
		GlobalRatePeriodInterval: globalRatePeriodInterval,
		MinNumSymbols:            minNumSymbols,
		PricePerSymbol:           pricePerSymbol,
		ReservationWindow:        reservationWindow,
	}, nil
}

// RefreshOnchainPaymentState returns the current onchain payment state
func (pcs *OnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context) error {
	paymentVaultParams, err := pcs.GetPaymentVaultParams(ctx)
	if err != nil {
		return err
	}
	// These parameters should be rarely updated, but we refresh them anyway
	pcs.PaymentVaultParams.Store(paymentVaultParams)

	pcs.ReservationsLock.Lock()
	accountIDs := make([]gethcommon.Address, 0, len(pcs.ReservedPayments))
	for accountID := range pcs.ReservedPayments {
		accountIDs = append(accountIDs, accountID)
	}

	reservedPayments, err := pcs.tx.GetReservedPayments(ctx, accountIDs)
	if err != nil {
		return err
	}
	pcs.ReservedPayments = reservedPayments
	pcs.ReservationsLock.Unlock()

	pcs.OnDemandLocks.Lock()
	accountIDs = make([]gethcommon.Address, 0, len(pcs.OnDemandPayments))
	for accountID := range pcs.OnDemandPayments {
		accountIDs = append(accountIDs, accountID)
	}

	onDemandPayments, err := pcs.tx.GetOnDemandPayments(ctx, accountIDs)
	if err != nil {
		return err
	}
	pcs.OnDemandPayments = onDemandPayments
	pcs.OnDemandLocks.Unlock()

	return nil
}

// GetReservedPaymentByAccount returns a pointer to the active reservation for the given account ID; no writes will be made to the reservation
func (pcs *OnchainPaymentState) GetReservedPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.ReservedPayment, error) {
	pcs.ReservationsLock.RLock()
	if reservation, ok := (pcs.ReservedPayments)[accountID]; ok {
		pcs.ReservationsLock.RUnlock()
		return reservation, nil
	}
	pcs.ReservationsLock.RUnlock()

	// pulls the chain state
	res, err := pcs.tx.GetReservedPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	pcs.ReservationsLock.Lock()
	(pcs.ReservedPayments)[accountID] = res
	pcs.ReservationsLock.Unlock()

	return res, nil
}

// GetOnDemandPaymentByAccount returns a pointer to the on-demand payment for the given account ID; no writes will be made to the payment
func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error) {
	pcs.OnDemandLocks.RLock()
	if payment, ok := (pcs.OnDemandPayments)[accountID]; ok {
		pcs.OnDemandLocks.RUnlock()
		return payment, nil
	}
	pcs.OnDemandLocks.RUnlock()

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

func (pcs *OnchainPaymentState) GetGlobalRatePeriodInterval() uint32 {
	return pcs.PaymentVaultParams.Load().GlobalRatePeriodInterval
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
