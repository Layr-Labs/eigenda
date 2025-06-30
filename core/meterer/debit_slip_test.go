package meterer

import (
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebitSlip_NewDebitSlip(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	timestamp := time.Now().UnixNano()

	paymentMetadata := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         timestamp,
		CumulativePayment: big.NewInt(100),
	}

	numSymbols := uint64(1000)
	quorumNumbers := []core.QuorumID{0, 1, 2}

	request := NewDebitSlip(paymentMetadata, numSymbols, quorumNumbers)

	assert.Equal(t, paymentMetadata, request.PaymentMetadata)
	assert.Equal(t, numSymbols, request.NumSymbols)
	assert.Equal(t, quorumNumbers, request.QuorumNumbers)
	assert.Equal(t, accountID, request.GetAccountID())
	assert.Equal(t, timestamp, request.GetTimestamp())
	assert.NotZero(t, request.ReceivedAt)
}

func TestDebitSlip_Validate(t *testing.T) {
	tests := []struct {
		name          string
		request       *DebitSlip
		expectedError string
	}{
		{
			name: "valid request",
			request: &DebitSlip{
				PaymentMetadata: core.PaymentMetadata{
					AccountID: gethcommon.HexToAddress("0x123"),
					Timestamp: time.Now().UnixNano(),
				},
				NumSymbols:    100,
				QuorumNumbers: []core.QuorumID{0, 1},
			},
			expectedError: "",
		},
		{
			name: "no quorums",
			request: &DebitSlip{
				PaymentMetadata: core.PaymentMetadata{
					AccountID: gethcommon.HexToAddress("0x123"),
					Timestamp: time.Now().UnixNano(),
				},
				NumSymbols:    100,
				QuorumNumbers: []core.QuorumID{},
			},
			expectedError: "no quorums provided",
		},
		{
			name: "zero symbols",
			request: &DebitSlip{
				PaymentMetadata: core.PaymentMetadata{
					AccountID: gethcommon.HexToAddress("0x123"),
					Timestamp: time.Now().UnixNano(),
				},
				NumSymbols:    0,
				QuorumNumbers: []core.QuorumID{0},
			},
			expectedError: "zero symbols requested",
		},
		{
			name: "invalid account ID",
			request: &DebitSlip{
				PaymentMetadata: core.PaymentMetadata{
					AccountID: gethcommon.Address{},
					Timestamp: time.Now().UnixNano(),
				},
				NumSymbols:    100,
				QuorumNumbers: []core.QuorumID{0},
			},
			expectedError: "invalid account ID",
		},
		{
			name: "invalid timestamp",
			request: &DebitSlip{
				PaymentMetadata: core.PaymentMetadata{
					AccountID: gethcommon.HexToAddress("0x123"),
					Timestamp: 0,
				},
				NumSymbols:    100,
				QuorumNumbers: []core.QuorumID{0},
			},
			expectedError: "invalid timestamp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestDebitSlip_WithMethods(t *testing.T) {
	request := &DebitSlip{
		PaymentMetadata: core.PaymentMetadata{
			AccountID: gethcommon.HexToAddress("0x123"),
			Timestamp: time.Now().UnixNano(),
		},
		NumSymbols:    100,
		QuorumNumbers: []core.QuorumID{0},
	}

	requestID := "test-request-123"
	receivedAt := time.Now().Add(-time.Hour)

	// Test method chaining
	result := request.WithRequestID(requestID).WithReceivedAt(receivedAt)

	assert.Equal(t, requestID, result.RequestID)
	assert.Equal(t, receivedAt, result.ReceivedAt)
	assert.Same(t, request, result) // Should return same instance for chaining
}

func TestDebitSlip_String(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	timestamp := int64(1234567890)

	request := &DebitSlip{
		PaymentMetadata: core.PaymentMetadata{
			AccountID: accountID,
			Timestamp: timestamp,
		},
		NumSymbols:    100,
		QuorumNumbers: []core.QuorumID{0, 1},
		RequestID:     "test-123",
	}

	str := request.String()

	assert.Contains(t, str, accountID.Hex())
	assert.Contains(t, str, "100")
	assert.Contains(t, str, "[0 1]")
	assert.Contains(t, str, "1234567890")
	assert.Contains(t, str, "test-123")
}
