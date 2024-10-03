package meterer

import (
	"context"
	"errors"

	"github.com/Layr-Labs/eigenda/core/eth"
)

/* HEAVILY MOCKED */
// PaymentAccounts (For reservations and on-demand payments)

type TokenAmount uint64 // TODO: change to uint128

// OperatorInfo contains information about an operator which is stored on the blockchain state,
// corresponding to a particular quorum
type ActiveReservation struct {
	DataRate   uint64 // Bandwidth per reservation bin
	StartEpoch uint32
	EndEpoch   uint32

	QuorumNumbers []uint8
	QuorumSplit   []byte // ordered mapping of quorum number to payment split; on-chain validation should ensure split <= 100
}

type OnDemandPayment struct {
	CumulativePayment TokenAmount // Total amount deposited by the user
}

// ActiveReservations contains information about the current state of active reservations
// map account ID to the ActiveReservation for that account.
type ActiveReservations struct {
	Reservations map[string]*ActiveReservation
}

// OnDemandPayments contains information about the current state of on-demand payments
// Map from account ID to the OnDemandPayment for that account.
type OnDemandPayments struct {
	Payments map[string]*OnDemandPayment
}

// OnchainPaymentState is an interface for getting information about the current chain state for payments.
type OnchainPaymentState struct {
	tx *eth.Transactor

	ActiveReservations *ActiveReservations
	OnDemandPayments   *OnDemandPayments
	// FUNCTIONS IF THIS STRUCT WAS AN INTERFACE?
	// GetActiveReservations(ctx context.Context, blockNumber uint) (map[string]*ActiveReservations, error)
	// GetActiveReservationByAccount(ctx context.Context, blockNumber uint, accountID string) (*ActiveReservation, error)
	// GetOnDemandPayments(ctx context.Context, blockNumber uint) (map[string]*OnDemandPayments, error)
	// GetOnDemandPaymentByAccount(ctx context.Context, blockNumber uint, accountID string) (*OnDemandPayment, error)
}

func NewOnchainPaymentState() *OnchainPaymentState {
	return &OnchainPaymentState{
		ActiveReservations: &ActiveReservations{},
		OnDemandPayments:   &OnDemandPayments{},
	}
}

// Mock data initialization method (mocked structs)
func (pcs *OnchainPaymentState) InitializeOnchainPaymentState() {
	// update with a pull from chain (write pulling functions in/core/eth/tx.go)
	pcs.ActiveReservations.Reservations = map[string]*ActiveReservation{}
	pcs.OnDemandPayments.Payments = map[string]*OnDemandPayment{}
}

func (pcs *OnchainPaymentState) GetActiveReservations(ctx context.Context, blockNumber uint) (*ActiveReservations, error) {
	return pcs.ActiveReservations, nil
}

func (pcs *OnchainPaymentState) GetActiveReservationByAccount(ctx context.Context, blockNumber uint, accountID string) (*ActiveReservation, error) {
	if reservation, ok := pcs.ActiveReservations.Reservations[accountID]; ok {
		return reservation, nil
	}
	return nil, errors.New("reservation not found")
}

func (pcs *OnchainPaymentState) GetOnDemandPayments(ctx context.Context, blockNumber uint) (*OnDemandPayments, error) {
	return pcs.OnDemandPayments, nil
}

func (pcs *OnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, blockNumber uint, accountID string) (*OnDemandPayment, error) {
	if payment, ok := pcs.OnDemandPayments.Payments[accountID]; ok {
		return payment, nil
	}
	return nil, errors.New("payment not found")
}
