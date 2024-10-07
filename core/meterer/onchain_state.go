package meterer

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/Layr-Labs/eigenda/core/eth"
)

/* HEAVILY MOCKED */

var (
	DummyReservationBytesLimit    = uint64(1024)
	DummyPaymentLimit             = TokenAmount(512)
	DummyMinimumChargeableSize    = uint32(128)
	DummyMinimumChargeablePayment = uint32(128)

	DummyReservation     = ActiveReservation{DataRate: DummyReservationBytesLimit, StartTimestamp: 0, EndTimestamp: math.MaxUint32, QuorumSplit: []byte{50, 50}, QuorumNumbers: []uint8{0, 1}}
	DummyOnDemandPayment = OnDemandPayment{CumulativePayment: DummyPaymentLimit}
)

// PaymentAccounts (For reservations and on-demand payments)

type TokenAmount uint64 // TODO: change to uint128

// OperatorInfo contains information about an operator which is stored on the blockchain state,
// corresponding to a particular quorum
type ActiveReservation struct {
	DataRate       uint64 // Bandwidth per reservation bin
	StartTimestamp uint64 // Unix timestamp that's valid for basically eternity
	EndTimestamp   uint64

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
	// TODO: update with pulling from chain; currently use a dummy
	pbk := "0x04cd9ba0357d1e5b929554e932cccdd6cf2d6e41d9d67907365b3e46cf005d5afd92b4f7bb3b829520be1a1b88641691973c98dfe68b07ee3613e270406285dfe8"
	pcs.ActiveReservations.Reservations = map[string]*ActiveReservation{
		pbk: &DummyReservation,
	}
	pcs.OnDemandPayments.Payments = map[string]*OnDemandPayment{
		pbk: &DummyOnDemandPayment,
	}

	fmt.Println("Initialized payment state with dummy reservation and on-demand payments")
}

func (pcs *OnchainPaymentState) GetCurrentBlockNumber(ctx context.Context) (uint, error) {
	return 0, nil
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
