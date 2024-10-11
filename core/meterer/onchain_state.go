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
	RefreshOnchainPaymentState(ctx context.Context, tx *eth.Transactor) error
	GetActiveReservations(ctx context.Context) (map[string]core.ActiveReservation, error)
	GetActiveReservationByAccount(ctx context.Context, accountID string) (core.ActiveReservation, error)
	GetOnDemandPayments(ctx context.Context) (map[string]core.OnDemandPayment, error)
	GetOnDemandPaymentByAccount(ctx context.Context, accountID string) (core.OnDemandPayment, error)
	GetOnDemandQuorumNumbers(ctx context.Context) []uint8
}

type OnchainPaymentState struct {
	tx *eth.Transactor

	ActiveReservations    map[string]core.ActiveReservation
	OnDemandPayments      map[string]core.OnDemandPayment
	OnDemandQuorumNumbers []uint8
}

var _ OnchainPayment = (*OnchainPaymentState)(nil)

func NewOnchainPaymentState(ctx context.Context, tx *eth.Transactor) (OnchainPaymentState, error) {
	initState := OnchainPaymentState{tx: tx}
	err := initState.RefreshOnchainPaymentState(ctx, tx)
	if err != nil {
		return OnchainPaymentState{tx: tx}, err
	}

	return initState, nil
}

// RefreshOnchainPaymentState returns the current onchain payment state (TODO: can optimize based on contract interface)
func (pcs OnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context, tx *eth.Transactor) error {
	blockNumber, err := tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return err
	}

	activeReservations, err := tx.GetActiveReservations(ctx, blockNumber)
	if err != nil {
		return err
	}
	pcs.ActiveReservations = activeReservations

	onDemandPayments, err := tx.GetOnDemandPayments(ctx, blockNumber)
	if err != nil {
		return err
	}
	pcs.OnDemandPayments = onDemandPayments

	quorumNumbers, err := tx.GetRequiredQuorumNumbers(ctx, blockNumber)
	if err != nil {
		return err
	}
	pcs.OnDemandQuorumNumbers = quorumNumbers

	return nil
}

func (pcs OnchainPaymentState) GetCurrentBlockNumber(ctx context.Context) (uint32, error) {
	blockNumber, err := pcs.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (pcs OnchainPaymentState) GetActiveReservations(ctx context.Context) (map[string]core.ActiveReservation, error) {
	return pcs.ActiveReservations, nil
}

// GetActiveReservationByAccount returns a pointer to the active reservation for the given account ID; no writes will be made to the reservation
func (pcs OnchainPaymentState) GetActiveReservationByAccount(ctx context.Context, accountID string) (core.ActiveReservation, error) {
	if reservation, ok := pcs.ActiveReservations[accountID]; ok {
		return reservation, nil
	}
	return core.ActiveReservation{}, errors.New("reservation not found")
}

func (pcs OnchainPaymentState) GetOnDemandPayments(ctx context.Context) (map[string]core.OnDemandPayment, error) {
	return pcs.OnDemandPayments, nil
}

// GetOnDemandPaymentByAccount returns a pointer to the on-demand payment for the given account ID; no writes will be made to the payment
func (pcs OnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID string) (core.OnDemandPayment, error) {
	if payment, ok := pcs.OnDemandPayments[accountID]; ok {
		return payment, nil
	}
	return core.OnDemandPayment{}, errors.New("payment not found")
}

func (pcs OnchainPaymentState) GetOnDemandQuorumNumbers(ctx context.Context) []uint8 {
	return pcs.OnDemandQuorumNumbers
}
