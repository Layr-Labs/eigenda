package meterer

import (
	"context"
	"crypto/ecdsa"
	"errors"

	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/ethereum/go-ethereum/crypto"
)

/* HEAVILY MOCKED */
// PaymentAccounts (For reservations and on-demand payments)

type TokenAmount uint64 // TODO: change to uint128

// OperatorInfo contains information about an operator which is stored on the blockchain state,
// corresponding to a particular quorum
type ActiveReservation struct {
	DataRate    uint32 // Bandwidth per reservation bin
	StartEpoch  uint64
	EndEpoch    uint64
	QuorumSplit []byte
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

	// GetActiveReservations(ctx context.Context, blockNumber uint) (map[string]*ActiveReservations, error)
	// GetActiveReservationByAccount(ctx context.Context, blockNumber uint, accountID string) (*ActiveReservation, error)
	// GetOnDemandPayments(ctx context.Context, blockNumber uint) (map[string]*OnDemandPayments, error)
	// GetOnDemandPaymentByAccount(ctx context.Context, blockNumber uint, accountID string) (*OnDemandPayment, error)
}

type MockedOnchainPaymentState struct {
	MockActiveReservations *ActiveReservations
	MockOnDemandPayments   *OnDemandPayments
}

func NewMockedOnchainPaymentState() *MockedOnchainPaymentState {
	return &MockedOnchainPaymentState{
		MockActiveReservations: &ActiveReservations{},
		MockOnDemandPayments:   &OnDemandPayments{},
	}
}

// Mock data initialization method
func (pcs *MockedOnchainPaymentState) InitializeMockData(privateKey1 *ecdsa.PrivateKey, privateKey2 *ecdsa.PrivateKey) {
	// Initialize mock active reservations
	binIndex := GetCurrentBinIndex()
	pcs.MockActiveReservations.Reservations = map[string]*ActiveReservation{
		crypto.PubkeyToAddress(privateKey1.PublicKey).Hex(): {DataRate: 100, StartEpoch: binIndex + 2, EndEpoch: binIndex + 5, QuorumSplit: []byte{50, 50}},
		crypto.PubkeyToAddress(privateKey2.PublicKey).Hex(): {DataRate: 200, StartEpoch: binIndex - 2, EndEpoch: binIndex + 10, QuorumSplit: []byte{30, 70}},
	}
	pcs.MockOnDemandPayments.Payments = map[string]*OnDemandPayment{
		crypto.PubkeyToAddress(privateKey1.PublicKey).Hex(): {CumulativePayment: 1000},
		crypto.PubkeyToAddress(privateKey2.PublicKey).Hex(): {CumulativePayment: 3000},
	}
}

func (pcs *MockedOnchainPaymentState) MockedGetActiveReservations(ctx context.Context, blockNumber uint) (*ActiveReservations, error) {
	return pcs.MockActiveReservations, nil
}

func (pcs *MockedOnchainPaymentState) MockedGetActiveReservationByAccount(ctx context.Context, blockNumber uint, accountID string) (*ActiveReservation, error) {
	if reservation, ok := pcs.MockActiveReservations.Reservations[accountID]; ok {
		return reservation, nil
	}
	return nil, errors.New("reservation not found")
}

func (pcs *MockedOnchainPaymentState) MockedGetOnDemandPayments(ctx context.Context, blockNumber uint) (*OnDemandPayments, error) {
	return pcs.MockOnDemandPayments, nil
}

func (pcs *MockedOnchainPaymentState) MockedGetOnDemandPaymentByAccount(ctx context.Context, blockNumber uint, accountID string) (*OnDemandPayment, error) {
	if payment, ok := pcs.MockOnDemandPayments.Payments[accountID]; ok {
		return payment, nil
	}
	return nil, errors.New("payment not found")
}
