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

func TestCreateDebitSlipBasic(t *testing.T) {
	accountID := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	timestamp := time.Now()
	transactionID := "test-transaction-123"
	sourceID := "disperser-1"
	amount := big.NewInt(1000000000000000000) // 1 ETH in wei
	symbolsCharged := uint64(1024)
	quorumNumbers := []core.QuorumID{0, 1}
	
	paymentHeader := PaymentHeader{
		AccountID:     accountID,
		Timestamp:     timestamp,
		TransactionID: transactionID,
		UsageType:     OnDemandPayment,
	}
	
	debitSlip, err := CreateDebitSlip(
		paymentHeader,
		OnDemandPayment,
		sourceID,
		amount,
		symbolsCharged,
		quorumNumbers,
	)
	
	require.NoError(t, err)
	require.NotNil(t, debitSlip)
	
	assert.Equal(t, paymentHeader.AccountID, debitSlip.PaymentHeader.AccountID)
	assert.Equal(t, OnDemandPayment, debitSlip.PaymentUsageType)
	assert.Equal(t, sourceID, debitSlip.SourceID)
	assert.Equal(t, amount, debitSlip.Amount)
	assert.Equal(t, symbolsCharged, debitSlip.SymbolsCharged)
	assert.Equal(t, quorumNumbers, debitSlip.QuorumNumbers)
}

func TestCreateReservationPaymentHeaderBasic(t *testing.T) {
	accountID := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	timestamp := time.Now()
	transactionID := "reservation-transaction-123"
	
	header, err := CreateReservationPaymentHeader(
		accountID,
		timestamp,
		transactionID,
	)
	
	require.NoError(t, err)
	require.NotNil(t, header)
	
	assert.Equal(t, accountID, header.AccountID)
	assert.Equal(t, timestamp, header.Timestamp)
	assert.Equal(t, transactionID, header.TransactionID)
	assert.Equal(t, ReservationPayment, header.UsageType)
}

func TestPaymentUsageTypeValues(t *testing.T) {
	assert.Equal(t, 0, int(ReservationPayment))
	assert.Equal(t, 1, int(OnDemandPayment))
}