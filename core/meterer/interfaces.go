package meterer

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// DebitSlip contains all information necessary to process a debit payment
type DebitSlip struct {
	// PaymentMetadata contains the core payment information (existing EigenDA type)
	PaymentMetadata *core.PaymentMetadata
	// NumSymbols is the number of symbols in the blob
	NumSymbols uint64
	// QuorumNumbers specifies which quorums this payment applies to
	QuorumNumbers []core.QuorumID
}

// AccountBalance represents the current balance state of an account
type AccountBalance struct {
	// AccountID is the account identifier
	AccountID gethcommon.Address
	// CumulativePayment is the total payment made by this account
	CumulativePayment *big.Int
	// ReservationUsage tracks usage across different reservation periods
	ReservationUsage map[uint32]*PeriodRecord
	// LastUpdateTime is when this balance was last updated
	LastUpdateTime time.Time
}

// AccountLedger defines the interface for tracking usage and payments for a single account
type AccountLedger interface {
	// CreateDebit processes a debit slip and updates the account ledger
	CreateDebit(debitSlip DebitSlip) error

	// CreateReservationPaymentMetadata creates a payment metadata for reservation-based payments
	CreateReservationPaymentMetadata(timestamp time.Time) (*core.PaymentMetadata, error)

	// CreateOnDemandPaymentMetadata creates a payment metadata for on-demand payments
	CreateOnDemandPaymentMetadata(timestamp time.Time, cumulativePayment *big.Int) (*core.PaymentMetadata, error)

	// GetAccountBalance retrieves the current balance state for this account
	GetAccountBalance() (*AccountBalance, error)

	// GetAccountID returns the account ID this ledger is tracking
	GetAccountID() gethcommon.Address
}

// CreateDebitSlip creates a new debit slip from payment components
func CreateDebitSlip(
	paymentMetadata *core.PaymentMetadata,
	numSymbols uint64,
	quorumNumbers []core.QuorumID,
) (*DebitSlip, error) {
	if paymentMetadata == nil {
		return nil, fmt.Errorf("payment metadata cannot be nil")
	}
	return &DebitSlip{
		PaymentMetadata: paymentMetadata,
		NumSymbols:      numSymbols,
		QuorumNumbers:   quorumNumbers,
	}, nil
}

// CreateReservationPaymentMetadata creates a payment metadata for reservation payments
func CreateReservationPaymentMetadata(
	accountID gethcommon.Address,
	timestamp time.Time,
) (*core.PaymentMetadata, error) {
	return &core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         timestamp.UnixNano(),
		CumulativePayment: big.NewInt(0), // Zero for reservation payments
	}, nil
}

// CreateOnDemandPaymentMetadata creates a payment metadata for on-demand payments
func CreateOnDemandPaymentMetadata(
	accountID gethcommon.Address,
	timestamp time.Time,
	cumulativePayment *big.Int,
) (*core.PaymentMetadata, error) {
	if cumulativePayment == nil || cumulativePayment.Sign() <= 0 {
		return nil, fmt.Errorf("cumulative payment must be positive for on-demand payments")
	}
	return &core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         timestamp.UnixNano(),
		CumulativePayment: new(big.Int).Set(cumulativePayment),
	}, nil
}

// IsOnDemandDebitSlip determines if a debit slip represents an on-demand payment
func IsOnDemandDebitSlip(debitSlip *DebitSlip) bool {
	if debitSlip == nil || debitSlip.PaymentMetadata == nil {
		return false
	}
	return IsOnDemandPayment(debitSlip.PaymentMetadata)
}


