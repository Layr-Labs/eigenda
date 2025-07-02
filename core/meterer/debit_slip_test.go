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
	tests := []struct {
		name            string
		paymentMetadata core.PaymentMetadata
		numSymbols      uint64
		quorumNumbers   []core.QuorumID
		expectedError   string
		checkFields     bool
	}{
		{
			name: "valid request",
			paymentMetadata: core.PaymentMetadata{
				AccountID:         gethcommon.HexToAddress("0x123"),
				Timestamp:         time.Now().UnixNano(),
				CumulativePayment: big.NewInt(100),
			},
			numSymbols:    1000,
			quorumNumbers: []core.QuorumID{0, 1, 2},
			expectedError: "",
			checkFields:   true,
		},
		{
			name: "no quorums",
			paymentMetadata: core.PaymentMetadata{
				AccountID:         gethcommon.HexToAddress("0x123"),
				Timestamp:         time.Now().UnixNano(),
				CumulativePayment: big.NewInt(100),
			},
			numSymbols:    100,
			quorumNumbers: []core.QuorumID{},
			expectedError: "no quorums provided",
		},
		{
			name: "zero symbols",
			paymentMetadata: core.PaymentMetadata{
				AccountID:         gethcommon.HexToAddress("0x123"),
				Timestamp:         time.Now().UnixNano(),
				CumulativePayment: big.NewInt(100),
			},
			numSymbols:    0,
			quorumNumbers: []core.QuorumID{0},
			expectedError: "zero symbols requested",
		},
		{
			name: "invalid account ID",
			paymentMetadata: core.PaymentMetadata{
				AccountID:         gethcommon.Address{},
				Timestamp:         time.Now().UnixNano(),
				CumulativePayment: big.NewInt(100),
			},
			numSymbols:    100,
			quorumNumbers: []core.QuorumID{0},
			expectedError: "invalid account ID",
		},
		{
			name: "invalid timestamp",
			paymentMetadata: core.PaymentMetadata{
				AccountID:         gethcommon.HexToAddress("0x123"),
				Timestamp:         0,
				CumulativePayment: big.NewInt(100),
			},
			numSymbols:    100,
			quorumNumbers: []core.QuorumID{0},
			expectedError: "invalid timestamp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := NewDebitSlip(tt.paymentMetadata, tt.numSymbols, tt.quorumNumbers)
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, request)

				if tt.checkFields {
					assert.Equal(t, tt.paymentMetadata, request.PaymentMetadata)
					assert.Equal(t, tt.numSymbols, request.NumSymbols)
					assert.Equal(t, tt.quorumNumbers, request.QuorumNumbers)
					assert.Equal(t, tt.paymentMetadata.AccountID, request.GetAccountID())
					assert.Equal(t, tt.paymentMetadata.Timestamp, request.GetTimestamp())
					assert.NotZero(t, request.ReceivedAt)
				}
			} else {
				require.Error(t, err)
				assert.Nil(t, request)
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
