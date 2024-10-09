package meterer

import (
	"context"
	"errors"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
)

// PaymentAccounts (For reservations and on-demand payments)

// OnchainPaymentState is an interface for getting information about the current chain state for payments.
type OnchainPayment interface {
	GetCurrentBlockNumber(ctx context.Context) (uint32, error)
	CurrentOnchainPaymentState(ctx context.Context, tx *eth.Transactor) (OnchainPaymentState, error)
	GetActiveReservations(ctx context.Context, blockNumber uint32) (map[string]core.ActiveReservation, error)
	GetActiveReservationsByAccount(ctx context.Context, blockNumber uint32, accountID string) (core.ActiveReservation, error)
	GetOnDemandPayments(ctx context.Context, blockNumber uint32) (map[string]core.OnDemandPayment, error)
	GetOnDemandPaymentByAccount(ctx context.Context, blockNumber uint32, accountID string) (core.OnDemandPayment, error)
}

type OnchainPaymentState struct {
	tx *eth.Transactor

	ActiveReservations map[string]core.ActiveReservation
	OnDemandPayments   map[string]core.OnDemandPayment
}

func NewOnchainPaymentState(ctx context.Context, tx *eth.Transactor) (OnchainPaymentState, error) {
	activeReservations, onDemandPayments, err := CurrentOnchainPaymentState(ctx, tx)
	if err != nil {
		return OnchainPaymentState{tx: tx}, err
	}

	return OnchainPaymentState{
		tx:                 tx,
		ActiveReservations: activeReservations,
		OnDemandPayments:   onDemandPayments,
	}, nil
}

// CurrentOnchainPaymentState returns the current onchain payment state (TODO: can optimize based on contract interface)
func CurrentOnchainPaymentState(ctx context.Context, tx *eth.Transactor) (map[string]core.ActiveReservation, map[string]core.OnDemandPayment, error) {
	blockNumber, err := tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, nil, err
	}

	activeReservations, err := tx.GetActiveReservations(ctx, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	onDemandPayments, err := tx.GetOnDemandPayments(ctx, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	return activeReservations, onDemandPayments, nil
}

func (pcs *OnchainPaymentState) GetCurrentBlockNumber(ctx context.Context) (uint32, error) {
	blockNumber, err := pcs.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (pcs *OnchainPaymentState) GetActiveReservations(ctx context.Context, blockNumber uint) (map[string]core.ActiveReservation, error) {
	return pcs.ActiveReservations, nil
}

// GetActiveReservationByAccount returns a pointer to the active reservation for the given account ID; no writes will be made to the reservation
func (pcs *OnchainPaymentState) GetActiveReservationByAccount(ctx context.Context, blockNumber uint, accountID string) (*core.ActiveReservation, error) {
	if reservation, ok := pcs.ActiveReservations[accountID]; ok {
		return &reservation, nil
	}
	return nil, errors.New("reservation not found")
}

func (pcs *OnchainPaymentState) GetOnDemandPayments(ctx context.Context, blockNumber uint) (map[string]core.OnDemandPayment, error) {
	return pcs.OnDemandPayments, nil
}

// GetOnDemandPaymentByAccount returns a pointer to the on-demand payment for the given account ID; no writes will be made to the payment
func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, blockNumber uint, accountID string) (*core.OnDemandPayment, error) {
	if payment, ok := pcs.OnDemandPayments[accountID]; ok {
		return &payment, nil
	}
	return nil, errors.New("payment not found")
}
