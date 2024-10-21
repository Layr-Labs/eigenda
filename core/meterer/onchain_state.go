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
	RefreshOnchainPaymentState(ctx context.Context, tx *eth.Transactor) error
	GetActiveReservations(ctx context.Context) (map[string]core.ActiveReservation, error)
	GetActiveReservationByAccount(ctx context.Context, accountID string) (core.ActiveReservation, error)
	GetOnDemandPayments(ctx context.Context) (map[string]core.OnDemandPayment, error)
	GetOnDemandPaymentByAccount(ctx context.Context, accountID string) (core.OnDemandPayment, error)
	GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error)
}

var _ OnchainPayment = (*OnchainPaymentState)(nil)

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

// RefreshOnchainPaymentState returns the current onchain payment state (TODO: can optimize based on contract interface)
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

func (pcs *OnchainPaymentState) GetActiveReservations(ctx context.Context) (map[string]core.ActiveReservation, error) {
	return pcs.ActiveReservations, nil
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

	pcs.ActiveReservations[accountID] = res
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
		return core.ActiveReservation{}, errors.New("reservation account not found on-chain")
	}
	return res, nil
}

func (pcs OnchainPaymentState) GetOnDemandPayments(ctx context.Context) (map[string]core.OnDemandPayment, error) {
	return pcs.OnDemandPayments, nil
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
	pcs.OnDemandPayments[accountID] = res
	return res, nil
}

func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccountOnChain(ctx context.Context, accountID string) (core.OnDemandPayment, error) {
	// pulls the chain state
	blockNumber, err := pcs.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return core.OnDemandPayment{}, err
	}
	res, err := pcs.tx.GetOnDemandPaymentByAccount(ctx, blockNumber, accountID)
	if err != nil {
		return core.OnDemandPayment{}, errors.New("on-demand not found on-chain")
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
