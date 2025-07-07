package meterer

import (
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

var (
	ErrNoQuorums        = errors.New("no quorums provided")
	ErrZeroSymbols      = errors.New("zero symbols requested")
	ErrInvalidAccount   = errors.New("invalid account ID")
	ErrInvalidTimestamp = errors.New("invalid timestamp")
)

// DebitSlip encapsulates all parameters needed for blob processing requests.
type DebitSlip struct {
	// Payment information - contains account ID, timestamp, and cumulative payment
	PaymentMetadata core.PaymentMetadata

	// Blob characteristics - number of symbols in the blob
	NumSymbols uint64

	// Target quorums for blob dispersal
	QuorumNumbers []core.QuorumID

	// Request identification (for tracking/logging)
	RequestID string // Optional, for tracing and debugging
}

// NewDebitSlip creates a new DebitSlip with the essential parameters and validates them
func NewDebitSlip(
	paymentMetadata core.PaymentMetadata,
	numSymbols uint64,
	quorumNumbers []core.QuorumID,
) (*DebitSlip, error) {
	// Validate the parameters during construction
	if len(quorumNumbers) == 0 {
		return nil, ErrNoQuorums
	}
	if numSymbols == 0 {
		return nil, ErrZeroSymbols
	}
	if paymentMetadata.AccountID == (gethcommon.Address{}) {
		return nil, ErrInvalidAccount
	}
	if paymentMetadata.Timestamp <= 0 {
		return nil, ErrInvalidTimestamp
	}

	return &DebitSlip{
		PaymentMetadata: paymentMetadata,
		NumSymbols:      numSymbols,
		QuorumNumbers:   quorumNumbers,
	}, nil
}

// GetAccountID returns the account ID from the payment metadata
func (ds *DebitSlip) GetAccountID() gethcommon.Address {
	return ds.PaymentMetadata.AccountID
}

// GetTimestamp returns the timestamp from the payment metadata
func (ds *DebitSlip) GetTimestamp() int64 {
	return ds.PaymentMetadata.Timestamp
}

// GetQuorumNumbersAsUint8 returns quorum numbers as uint8 slice
func (ds *DebitSlip) GetQuorumNumbersAsUint8() []uint8 {
	result := make([]uint8, len(ds.QuorumNumbers))
	for i, qid := range ds.QuorumNumbers {
		result[i] = uint8(qid)
	}
	return result
}

// WithRequestID sets the request ID for tracking purposes
func (ds *DebitSlip) WithRequestID(requestID string) *DebitSlip {
	ds.RequestID = requestID
	return ds
}

// String returns a string representation for logging and debugging
func (ds *DebitSlip) String() string {
	return fmt.Sprintf("DebitSlip{Account: %s, Symbols: %d, Quorums: %v, Timestamp: %d, RequestID: %s}",
		ds.PaymentMetadata.AccountID.Hex(),
		ds.NumSymbols,
		ds.QuorumNumbers,
		ds.PaymentMetadata.Timestamp,
		ds.RequestID,
	)
}
