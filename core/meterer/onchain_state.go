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
	CurrentOnchainPaymentState(ctx context.Context, tx *eth.Transactor) (OnchainPaymentState, error)
	GetActiveReservations(ctx context.Context) (map[string]core.ActiveReservation, error)
	GetActiveReservationByAccount(ctx context.Context, accountID string) (core.ActiveReservation, error)
	GetOnDemandPayments(ctx context.Context) (map[string]core.OnDemandPayment, error)
	GetOnDemandPaymentByAccount(ctx context.Context, accountID string) (core.OnDemandPayment, error)
	GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error)
}

type OnchainPaymentState struct {
	tx *eth.Transactor

	ActiveReservations    map[string]core.ActiveReservation
	OnDemandPayments      map[string]core.OnDemandPayment
	OnDemandQuorumNumbers []uint8
}

func NewOnchainPaymentState(ctx context.Context, tx *eth.Transactor) (OnchainPaymentState, error) {
	blockNumber, err := tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return OnchainPaymentState{}, err
	}

	quorumNumbers, err := tx.GetRequiredQuorumNumbers(ctx, blockNumber)
	if err != nil {
		return OnchainPaymentState{}, err
	}

	return OnchainPaymentState{
		tx:                    tx,
		ActiveReservations:    make(map[string]core.ActiveReservation),
		OnDemandPayments:      make(map[string]core.OnDemandPayment),
		OnDemandQuorumNumbers: quorumNumbers,
	}, nil
}

// CurrentOnchainPaymentState returns the current onchain payment state (TODO: can optimize based on contract interface)
func (pcs *OnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context, tx *eth.Transactor) error {
	blockNumber, err := tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return err
	}

	accountIDs := make([]string, 0, len(pcs.ActiveReservations))
	for accountID := range pcs.ActiveReservations {
		accountIDs = append(accountIDs, accountID)
	}

	activeReservations, err := tx.GetActiveReservations(ctx, blockNumber, accountIDs)
	if err != nil {
		return err
	}
	pcs.ActiveReservations = activeReservations

	accountIDs = make([]string, 0, len(pcs.OnDemandPayments))
	for accountID := range pcs.OnDemandPayments {
		accountIDs = append(accountIDs, accountID)
	}

	onDemandPayments, err := tx.GetOnDemandPayments(ctx, blockNumber, accountIDs)
	if err != nil {
		return err
	}
	pcs.OnDemandPayments = onDemandPayments

	return nil
}

func (pcs *OnchainPaymentState) GetActiveReservations(ctx context.Context, blockNumber uint) (map[string]core.ActiveReservation, error) {
	return pcs.ActiveReservations, nil
}

// GetActiveReservationByAccount returns a pointer to the active reservation for the given account ID; no writes will be made to the reservation
func (pcs *OnchainPaymentState) GetActiveReservationByAccount(ctx context.Context, blockNumber uint32, accountID string) (core.ActiveReservation, error) {
	if reservation, ok := pcs.ActiveReservations[accountID]; ok {
		return reservation, nil
	}
	// pulls the chain state
	res, err := pcs.tx.GetActiveReservationByAccount(ctx, blockNumber, accountID)
	if err != nil {
		return core.ActiveReservation{}, errors.New("payment not found")
	}

	pcs.ActiveReservations[accountID] = res
	return res, nil
}

func (pcs *OnchainPaymentState) GetOnDemandPayments(ctx context.Context, blockNumber uint) (map[string]core.OnDemandPayment, error) {
	return pcs.OnDemandPayments, nil
}

// GetOnDemandPaymentByAccount returns a pointer to the on-demand payment for the given account ID; no writes will be made to the payment
func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, blockNumber uint32, accountID string) (core.OnDemandPayment, error) {
	if payment, ok := pcs.OnDemandPayments[accountID]; ok {
		return payment, nil
	}
	// pulls the chain state
	res, err := pcs.tx.GetOnDemandPaymentByAccount(ctx, blockNumber, accountID)
	if err != nil {
		return core.OnDemandPayment{}, errors.New("payment not found")
	}

	pcs.OnDemandPayments[accountID] = res
	return res, nil
}

func (pcs *OnchainPaymentState) GetOnDemandQuorumNumbers(ctx context.Context, blockNumber uint32) ([]uint8, error) {
	return pcs.tx.GetRequiredQuorumNumbers(ctx, blockNumber)
}
