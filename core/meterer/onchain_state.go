package meterer

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// PaymentAccounts (For reservations and on-demand payments)

// OnchainPaymentState is an interface for getting information about the current chain state for payments.
type OnchainPayment interface {
	RefreshOnchainPaymentState(ctx context.Context, tx *eth.Reader) error
	GetReservedPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.ReservedPayment, error)
	GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error)
	GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error)
	GetGlobalSymbolsPerSecond() uint64
	GetGlobalRateBinInterval() uint32
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
	GlobalSymbolsPerSecond uint64
	GlobalRateBinInterval  uint32
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
		ReservedPayments:   make(map[gethcommon.Address]*core.ReservedPayment),
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
	accountIDs := make([]gethcommon.Address, 0, len(pcs.ReservedPayments))
	for accountID := range pcs.ReservedPayments {
		accountIDs = append(accountIDs, accountID)
	}

	reservedPayments, err := tx.GetReservedPayments(ctx, accountIDs)
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

	onDemandPayments, err := tx.GetOnDemandPayments(ctx, accountIDs)
	if err != nil {
		return err
	}
	pcs.OnDemandPayments = onDemandPayments
	pcs.OnDemandLocks.Unlock()

	return nil
}

// GetReservedPaymentByAccount returns a pointer to the active reservation for the given account ID; no writes will be made to the reservation
func (pcs *OnchainPaymentState) GetReservedPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.ReservedPayment, error) {
	timestamp := uint64(time.Now().Unix())
	pcs.ReservationsLock.Lock()
	defer pcs.ReservationsLock.Unlock()

	if reservation, ok := (pcs.ReservedPayments)[accountID]; ok {
		if !reservation.IsActive(timestamp) {
			// if reservation is expired, remove it from the local state; if it is not activated, we leave the reservation in the local state
			if reservation.EndTimestamp < timestamp {
				delete(pcs.ReservedPayments, accountID)
			}
			return nil, fmt.Errorf("reservation not active")
		}
		return reservation, nil
	}

	// pulls the chain state
	res, err := pcs.tx.GetReservedPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if !res.IsActive(timestamp) {
		if res.StartTimestamp > timestamp {
			// if reservation is not activated yet, we add it to the local state to reduce future on-chain calls
			(pcs.ReservedPayments)[accountID] = res
		}
		return nil, fmt.Errorf("reservation not active")
	}
	(pcs.ReservedPayments)[accountID] = res

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

func (pcs *OnchainPaymentState) GetGlobalRateBinInterval() uint32 {
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
