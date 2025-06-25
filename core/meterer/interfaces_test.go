package meterer

import (
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDebitSlip(t *testing.T) {
	accountID := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	paymentMetadata := &core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().UnixNano(),
		CumulativePayment: big.NewInt(1000),
	}

	debitSlip, err := CreateDebitSlip(paymentMetadata, 100, []core.QuorumID{0, 1})

	require.NoError(t, err)
	require.NotNil(t, debitSlip)
	assert.Equal(t, paymentMetadata, debitSlip.PaymentMetadata)
	assert.Equal(t, uint64(100), debitSlip.NumSymbols)
	assert.Equal(t, []core.QuorumID{0, 1}, debitSlip.QuorumNumbers)
}

func TestCreateDebitSlipWithNilPaymentMetadata(t *testing.T) {
	debitSlip, err := CreateDebitSlip(nil, 100, []core.QuorumID{0, 1})

	require.Error(t, err)
	require.Nil(t, debitSlip)
	assert.Contains(t, err.Error(), "payment metadata cannot be nil")
}

func TestCreateReservationPaymentMetadata(t *testing.T) {
	accountID := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	timestamp := time.Now()

	metadata, err := CreateReservationPaymentMetadata(accountID, timestamp)

	require.NoError(t, err)
	require.NotNil(t, metadata)
	assert.Equal(t, accountID, metadata.AccountID)
	assert.Equal(t, timestamp.UnixNano(), metadata.Timestamp)
	assert.Equal(t, big.NewInt(0), metadata.CumulativePayment)
}

func TestCreateOnDemandPaymentMetadata(t *testing.T) {
	accountID := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	timestamp := time.Now()
	payment := big.NewInt(1000)

	metadata, err := CreateOnDemandPaymentMetadata(accountID, timestamp, payment)

	require.NoError(t, err)
	require.NotNil(t, metadata)
	assert.Equal(t, accountID, metadata.AccountID)
	assert.Equal(t, timestamp.UnixNano(), metadata.Timestamp)
	assert.Equal(t, payment, metadata.CumulativePayment)
	// Ensure it's a copy, not the same reference
	assert.NotSame(t, payment, metadata.CumulativePayment)
}

func TestCreateOnDemandPaymentMetadataWithInvalidPayment(t *testing.T) {
	accountID := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	timestamp := time.Now()

	// Test nil payment
	metadata, err := CreateOnDemandPaymentMetadata(accountID, timestamp, nil)
	require.Error(t, err)
	require.Nil(t, metadata)
	assert.Contains(t, err.Error(), "cumulative payment must be positive")

	// Test zero payment
	metadata, err = CreateOnDemandPaymentMetadata(accountID, timestamp, big.NewInt(0))
	require.Error(t, err)
	require.Nil(t, metadata)
	assert.Contains(t, err.Error(), "cumulative payment must be positive")

	// Test negative payment
	metadata, err = CreateOnDemandPaymentMetadata(accountID, timestamp, big.NewInt(-1))
	require.Error(t, err)
	require.Nil(t, metadata)
	assert.Contains(t, err.Error(), "cumulative payment must be positive")
}

func TestIsOnDemandDebitSlip(t *testing.T) {
	accountID := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	timestamp := time.Now()

	// Test reservation payment (cumulative payment = 0)
	reservationMetadata := &core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         timestamp.UnixNano(),
		CumulativePayment: big.NewInt(0),
	}
	reservationDebitSlip := &DebitSlip{
		PaymentMetadata: reservationMetadata,
		NumSymbols:      100,
		QuorumNumbers:   []core.QuorumID{0},
	}
	assert.False(t, IsOnDemandDebitSlip(reservationDebitSlip))

	// Test on-demand payment (cumulative payment > 0)
	onDemandMetadata := &core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         timestamp.UnixNano(),
		CumulativePayment: big.NewInt(1000),
	}
	onDemandDebitSlip := &DebitSlip{
		PaymentMetadata: onDemandMetadata,
		NumSymbols:      100,
		QuorumNumbers:   []core.QuorumID{0},
	}
	assert.True(t, IsOnDemandDebitSlip(onDemandDebitSlip))

	// Test nil debit slip
	assert.False(t, IsOnDemandDebitSlip(nil))

	// Test debit slip with nil payment metadata
	nilMetadataDebitSlip := &DebitSlip{
		PaymentMetadata: nil,
		NumSymbols:      100,
		QuorumNumbers:   []core.QuorumID{0},
	}
	assert.False(t, IsOnDemandDebitSlip(nilMetadataDebitSlip))
}

func TestAccountLedgerInterface(t *testing.T) {
	// Compile-time test to ensure interface methods exist
	var ledger AccountLedger
	_ = ledger.CreateDebit
	_ = ledger.CreateReservationPaymentMetadata
	_ = ledger.CreateOnDemandPaymentMetadata
	_ = ledger.GetAccountBalance
	_ = ledger.GetAccountID
}

func TestDebitSlipStructure(t *testing.T) {
	accountID := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	paymentMetadata := &core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().UnixNano(),
		CumulativePayment: big.NewInt(1000),
	}

	debitSlip := DebitSlip{
		PaymentMetadata: paymentMetadata,
		NumSymbols:      100,
		QuorumNumbers:   []core.QuorumID{0, 1},
	}

	assert.Equal(t, paymentMetadata, debitSlip.PaymentMetadata)
	assert.Equal(t, uint64(100), debitSlip.NumSymbols)
	assert.Equal(t, []core.QuorumID{0, 1}, debitSlip.QuorumNumbers)
}

func TestAccountBalanceStructure(t *testing.T) {
	accountID := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	cumulativePayment := big.NewInt(1000)
	reservationUsage := make(map[uint32]*PeriodRecord)
	lastUpdateTime := time.Now()

	accountBalance := AccountBalance{
		AccountID:         accountID,
		CumulativePayment: cumulativePayment,
		ReservationUsage:  reservationUsage,
		LastUpdateTime:    lastUpdateTime,
	}

	assert.Equal(t, accountID, accountBalance.AccountID)
	assert.Equal(t, cumulativePayment, accountBalance.CumulativePayment)
	assert.Equal(t, reservationUsage, accountBalance.ReservationUsage)
	assert.Equal(t, lastUpdateTime, accountBalance.LastUpdateTime)
}