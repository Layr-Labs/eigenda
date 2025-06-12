package clients

import (
	"encoding/hex"
	"math"
	"math/big"
	"sync"
	"testing"
	"time"

	v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create standard PaymentVaultParams for testing
func createTestPaymentVaultParams(reservationWindow, pricePerSymbol, minNumSymbols uint64) *v2.PaymentVaultParams {
	return &v2.PaymentVaultParams{
		QuorumPaymentConfigs: map[uint32]*v2.PaymentQuorumConfig{
			0: {
				ReservationSymbolsPerSecond: 2000,
				OnDemandSymbolsPerSecond:    1000,
				OnDemandPricePerSymbol:      pricePerSymbol,
			},
			1: {
				ReservationSymbolsPerSecond: 2000,
				OnDemandSymbolsPerSecond:    1000,
				OnDemandPricePerSymbol:      pricePerSymbol,
			},
		},
		QuorumProtocolConfigs: map[uint32]*v2.PaymentQuorumProtocolConfig{
			0: {
				MinNumSymbols:              minNumSymbols,
				ReservationAdvanceWindow:   10,
				ReservationRateLimitWindow: reservationWindow,
				OnDemandRateLimitWindow:    30,
				OnDemandEnabled:            true,
			},
			1: {
				MinNumSymbols:              minNumSymbols,
				ReservationAdvanceWindow:   10,
				ReservationRateLimitWindow: reservationWindow,
				OnDemandRateLimitWindow:    30,
				OnDemandEnabled:            true,
			},
		},
		OnDemandQuorumNumbers: []uint32{0, 1},
	}
}

func TestNewAccountant(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	reservations := map[uint8]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 100,
			StartTimestamp:   uint64(time.Now().Unix()),
			EndTimestamp:     uint64(time.Now().Add(time.Hour).Unix()),
		},
		1: {
			SymbolsPerSecond: 200,
			StartTimestamp:   uint64(time.Now().Unix()),
			EndTimestamp:     uint64(time.Now().Add(time.Hour).Unix()),
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}

	acc := NewAccountant(accountID)
	now := time.Now()
	err := acc.SetPaymentState(&v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(10, 1, 100),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: map[uint32]*v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint32(now.Unix()),
				EndTimestamp:     uint32(now.Add(time.Hour).Unix()),
			},
			1: {
				SymbolsPerSecond: 200,
				StartTimestamp:   uint32(now.Unix()),
				EndTimestamp:     uint32(now.Add(time.Hour).Unix()),
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, acc)
	assert.Equal(t, accountID, acc.accountID)
	assert.Equal(t, reservations, acc.reservations)
	assert.Equal(t, onDemand, acc.onDemand)
	assert.NotNil(t, acc.periodRecords)

	// Should return minimum symbols for zero symbols
	header, err := acc.AccountBlob(now.UnixNano(), 0, []uint8{0})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "zero symbols requested")
	assert.Nil(t, header)

	// Should return zero payment for empty quorum list
	header, err = acc.AccountBlob(now.UnixNano(), 50, []uint8{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no quorums provided")
	assert.Nil(t, header)

	// Try to use maximum possible symbols
	maxSymbols := uint64(math.MaxUint64)
	header, err = acc.AccountBlob(now.UnixNano(), maxSymbols, []uint8{0})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "current cumulativePayment balance insufficient")
	assert.Nil(t, header)
}

func TestAccountant_Reservation(t *testing.T) {
	reservation := &v2.QuorumReservation{
		SymbolsPerSecond: 200,
		StartTimestamp:   uint32(time.Now().Unix()),
		EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
	}
	reservationWindow := uint64(5)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)

	// Create payment state with test configurations
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(500).Bytes(),
		Reservations:             map[uint32]*v2.QuorumReservation{0: reservation, 1: reservation},
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	symbolLength := uint64(500)
	quorums := []uint8{0, 1}
	now := time.Now().UnixNano()

	header, err := accountant.AccountBlob(now, symbolLength, []uint8{0, 1, 2})
	assert.Nil(t, header)
	assert.Error(t, err, "no reservation found for quorum")

	header, err = accountant.AccountBlob(now, symbolLength, quorums)
	assert.NoError(t, err)
	window := accountant.paymentVaultParams.QuorumProtocolConfigs[0].ReservationRateLimitWindow
	assert.Equal(t, meterer.GetReservationPeriod(time.Now().Unix(), window), meterer.GetReservationPeriodByNanosecond(header.Timestamp, window))
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Check that usage for quorum 0 was updated properly
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, window)
	record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, 0)
	assert.Equal(t, symbolLength, record.Usage)

	symbolLength = uint64(700)

	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(now, symbolLength, quorums)

	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// With overflow, usage should be at the limit in the current period
	binLimit := meterer.GetReservationBinLimit(accountant.reservations[0], window)
	expectedUsage := binLimit
	assert.Equal(t, expectedUsage, record.Usage)

	// Check overflow usage
	overflowIndex := meterer.GetOverflowPeriod(currentPeriod, window)
	expectedOverflow := symbolLength - (binLimit - 500) // 700 - (1000 - 500) = 200
	relativeRecord := accountant.periodRecords.GetRelativePeriodRecord(overflowIndex, 0)
	assert.Equal(t, expectedOverflow, relativeRecord.Usage)

	// Second call should use on-demand payment
	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(now, 300, quorums)

	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(300), header.CumulativePayment)
}

func TestAccountant_OnDemand(t *testing.T) {
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	numSymbols := uint64(1500)
	quorums := []uint8{0, 1}
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)
	err = accountant.SetPaymentState(&v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(5, 1, 100),
		OnchainCumulativePayment: big.NewInt(1500).Bytes(),
		Reservations:             map[uint32]*v2.QuorumReservation{},
	})
	assert.NoError(t, err)

	now := time.Now().UnixNano()

	// valid payment
	header, err := accountant.AccountBlob(now, numSymbols, quorums)
	assert.NoError(t, err)

	pricePerSymbol := accountant.paymentVaultParams.QuorumPaymentConfigs[0].OnDemandPricePerSymbol
	assert.NoError(t, err)
	expectedPayment := big.NewInt(int64(numSymbols * pricePerSymbol))
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, expectedPayment, header.CumulativePayment)

	// Check that no reservation usage was recorded
	for quorumNumber, records := range accountant.periodRecords {
		for _, record := range records {
			assert.Equal(t, uint64(0), record.Usage, "Usage should be 0 for quorum %d", quorumNumber)
		}
	}

	assert.Equal(t, expectedPayment, accountant.cumulativePayment)

}

type accountBlobInsufficientOnDemandTest struct {
	name         string
	symbolLength uint64
	expectError  bool
	errorMessage string
}

func TestAccountant_InsufficientOnDemand(t *testing.T) {
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)

	// Create payment state with test configurations
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(60, 100, 100),
		OnchainCumulativePayment: big.NewInt(500).Bytes(),
		Reservations:             map[uint32]*v2.QuorumReservation{},
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	quorums := []uint8{0, 1}
	baseTime := time.Now().UnixNano()

	tests := []accountBlobInsufficientOnDemandTest{
		{
			name:         "Insufficient on-demand payment",
			symbolLength: 2000,
			expectError:  true,
			errorMessage: "balance insufficient to make an on-demand dispersal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := baseTime
			header, err := accountant.AccountBlob(now, tt.symbolLength, quorums)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
				assert.Nil(t, header)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, header)
			}
		})
	}
}

type accountBlobCallSeriesTest struct {
	name           string
	symbolLength   uint64
	expectedHeader *core.PaymentMetadata
	expectedState  meterer.PeriodRecord
	expectError    bool
	errorMessage   string
}

func TestAccountant_AccountBlobCallSeries(t *testing.T) {
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	reservationWindow := uint64(5)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)

	// Create payment state with test configurations
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: map[uint32]*v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 200,
				StartTimestamp:   uint32(time.Now().Unix()),
				EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
			},
			1: {
				SymbolsPerSecond: 200,
				StartTimestamp:   uint32(time.Now().Unix()),
				EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
			},
		},
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	quorums := []uint8{0, 1}

	tests := []accountBlobCallSeriesTest{
		{
			name:         "First call - Use reservation",
			symbolLength: 800,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: meterer.PeriodRecord{
				Index: 800,
				Usage: 800,
			},
		},
		{
			name:         "Second call - Use remaining reservation + overflow",
			symbolLength: 300,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: meterer.PeriodRecord{
				Index: 1100,
				Usage: 100,
			},
		},
		{
			name:         "Third call - Use on-demand",
			symbolLength: 500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(500),
			},
			expectedState: meterer.PeriodRecord{
				Index: 1100,
				Usage: 100,
			},
		},
		{
			name:         "Fourth call - Insufficient on-demand",
			symbolLength: 600,
			expectError:  true,
			errorMessage: "balance insufficient to make an on-demand dispersal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now().UnixNano()
			header, err := accountant.AccountBlob(now, tt.symbolLength, quorums)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uint64(0), header.Timestamp)
				assert.Equal(t, tt.expectedHeader.AccountID, header.AccountID)
				assert.Equal(t, tt.expectedHeader.CumulativePayment, header.CumulativePayment)
			}
		})
	}
}

func TestAccountBlob_BinRotation(t *testing.T) {
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	reservationWindow := uint64(1)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)

	// Create payment state with test configurations
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: map[uint32]*v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 1000,
				StartTimestamp:   uint32(time.Now().Unix()),
				EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
			},
			1: {
				SymbolsPerSecond: 1000,
				StartTimestamp:   uint32(time.Now().Unix()),
				EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
			},
		},
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	quorums := []uint8{0, 1}

	// First call
	now := time.Now().UnixNano()
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(now, 800, quorums)
	assert.NoError(t, err)

	// Check bin 0 has usage 800 for charged quorums
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(800), record.Usage)
	}

	// Second call
	now += int64(reservationWindow) * time.Second.Nanoseconds()
	nextPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(now, 300, quorums)
	assert.NoError(t, err)

	// Check bin 1 has usage 300
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(800), record.Usage)
		record = accountant.periodRecords.GetRelativePeriodRecord(nextPeriod, quorumNumber)
		assert.Equal(t, uint64(300), record.Usage)
	}

	// Third call
	_, err = accountant.AccountBlob(now, 500, quorums)
	assert.NoError(t, err)

	// Check bin 1 now has usage 800
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(800), record.Usage)
		record = accountant.periodRecords.GetRelativePeriodRecord(nextPeriod, quorumNumber)
		assert.Equal(t, uint64(800), record.Usage)
	}
}

func TestAccountant_Concurrent(t *testing.T) {
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	reservationWindow := uint64(1)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)

	// Create payment state with test configurations
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: map[uint32]*v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 1000,
				StartTimestamp:   uint32(time.Now().Unix()),
				EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
			},
			1: {
				SymbolsPerSecond: 1000,
				StartTimestamp:   uint32(time.Now().Unix()),
				EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
			},
		},
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	quorums := []uint8{0, 1}

	// Start concurrent AccountBlob calls
	now := time.Now().UnixNano()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := accountant.AccountBlob(now, 100, quorums)
			assert.NoError(t, err)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check final state
	for _, quorumNumber := range quorums {
		currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(1000), record.Usage)
	}
}

func TestAccountBlob_ReservationWithOneOverflow(t *testing.T) {
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	reservationWindow := uint64(5)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)

	// Create payment state with test configurations
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: map[uint32]*v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 200,
				StartTimestamp:   uint32(time.Now().Unix()),
				EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
			},
			1: {
				SymbolsPerSecond: 200,
				StartTimestamp:   uint32(time.Now().Unix()),
				EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
			},
		},
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	quorums := []uint8{0, 1}
	now := time.Now().UnixNano()
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	header, err := accountant.AccountBlob(now, 800, quorums)
	assert.NoError(t, err)
	timestamp := (time.Duration(header.Timestamp) * time.Nanosecond).Seconds()
	assert.Equal(t, uint64(meterer.GetReservationPeriodByNanosecond(now, reservationWindow)), uint64(meterer.GetReservationPeriod(int64(timestamp), reservationWindow)))
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Check current period usage
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(800), record.Usage)
	}
	// Second call: Allow one overflow
	header, err = accountant.AccountBlob(now, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Check current period is at limit
	binLimit := accountant.reservations[0].SymbolsPerSecond * uint64(reservationWindow) // 1000
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, binLimit, record.Usage)
	}

	// Check overflow period has the overflow
	for _, quorumNumber := range quorums {
		overflowRecord := accountant.periodRecords.GetRelativePeriodRecord(meterer.GetOverflowPeriod(currentPeriod, reservationWindow), quorumNumber)
		assert.Equal(t, uint64(300), overflowRecord.Usage) // 800 + 500 - 1000 = 300
	}

	// Third call: Should use on-demand payment
	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(now, 200, quorums)
	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(200), header.CumulativePayment)
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(binLimit), record.Usage)
		record = accountant.periodRecords.GetRelativePeriodRecord(meterer.GetOverflowPeriod(currentPeriod, reservationWindow), quorumNumber)
		assert.Equal(t, uint64(300), record.Usage)
	}
}

func TestAccountBlob_ReservationOverflowReset(t *testing.T) {
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	reservationWindow := uint64(1)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)

	// Create payment state with test configurations
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: map[uint32]*v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 1000,
				StartTimestamp:   uint32(time.Now().Unix()),
				EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
			},
			1: {
				SymbolsPerSecond: 1000,
				StartTimestamp:   uint32(time.Now().Unix()),
				EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
			},
		},
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	quorums := []uint8{0, 1}

	// full reservation
	now := time.Now().UnixNano()
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(now, 1000, quorums)
	assert.NoError(t, err)

	// Check current period is at limit
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(1000), record.Usage)
	}

	// no overflow
	now = time.Now().UnixNano()
	header, err := accountant.AccountBlob(now, 500, quorums)
	assert.NoError(t, err)
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(1000), record.Usage)
	}
	assert.Equal(t, big.NewInt(500), header.CumulativePayment)

	// Wait for next reservation duration
	time.Sleep(time.Duration(reservationWindow) * time.Second)

	// Third call: Should use new bin and allow overflow again
	now = time.Now().UnixNano()
	nextPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	assert.Equal(t, currentPeriod+1, nextPeriod, "Should be next period")

	_, err = accountant.AccountBlob(now, 500, quorums)
	assert.NoError(t, err)

	// Check next period has usage 500
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(nextPeriod, quorumNumber)
		assert.Equal(t, uint64(500), record.Usage)
	}
}

func TestAccountant_SetPaymentState(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	acc := NewAccountant(accountID)

	tests := []struct {
		name    string
		state   *v2.GetPaymentStateForAllQuorumsReply
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil payment state",
			state:   nil,
			wantErr: true,
			errMsg:  "payment state cannot be nil",
		},
		{
			name:    "nil payment vault params",
			state:   &v2.GetPaymentStateForAllQuorumsReply{},
			wantErr: true,
			errMsg:  "payment vault params cannot be nil",
		},
		{
			name: "successful state update",
			state: &v2.GetPaymentStateForAllQuorumsReply{
				PaymentVaultParams: &v2.PaymentVaultParams{
					QuorumPaymentConfigs: map[uint32]*v2.PaymentQuorumConfig{
						0: {
							ReservationSymbolsPerSecond: 100,
							OnDemandSymbolsPerSecond:    200,
							OnDemandPricePerSymbol:      10,
						},
					},
					QuorumProtocolConfigs: map[uint32]*v2.PaymentQuorumProtocolConfig{
						0: {
							MinNumSymbols:              1,
							ReservationAdvanceWindow:   2,
							ReservationRateLimitWindow: 3,
							OnDemandRateLimitWindow:    4,
							OnDemandEnabled:            true,
						},
					},
				},
				Reservations: map[uint32]*v2.QuorumReservation{
					0: {
						SymbolsPerSecond: 100,
						StartTimestamp:   1000,
						EndTimestamp:     2000,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := acc.SetPaymentState(tt.state)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, acc.paymentVaultParams.QuorumPaymentConfigs[0])
				assert.NotNil(t, acc.paymentVaultParams.QuorumProtocolConfigs[0])
				assert.NotNil(t, acc.reservations[0])
			}
		})
	}
}

func TestAccountant_UseReservation(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now().UnixNano()
	acc := NewAccountant(accountID)
	err := acc.SetPaymentState(&v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams: createTestPaymentVaultParams(6, 1, 100),
		Reservations: map[uint32]*v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint32(now/1e9 - 1000),
				EndTimestamp:     uint32(now/1e9 + 1000),
			},
			1: {
				SymbolsPerSecond: 200,
				StartTimestamp:   uint32(now/1e9 - 1000),
				EndTimestamp:     uint32(now/1e9 + 1000),
			},
		},
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
	})
	require.NoError(t, err)
	// Set up protocol configs
	binInterval := uint64(1)
	acc.paymentVaultParams.QuorumProtocolConfigs[0] = &core.PaymentQuorumProtocolConfig{
		ReservationRateLimitWindow: binInterval,
		MinNumSymbols:              1,
	}
	acc.paymentVaultParams.QuorumProtocolConfigs[1] = &core.PaymentQuorumProtocolConfig{
		ReservationRateLimitWindow: binInterval,
		MinNumSymbols:              1,
	}

	tests := []struct {
		name          string
		numSymbols    uint64
		quorumNumbers []uint8
		timestamp     int64
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "successful reservation usage",
			numSymbols:    50,
			quorumNumbers: []uint8{0, 1},
			timestamp:     now,
			wantErr:       false,
		},
		{
			name:          "quorum without reservation",
			numSymbols:    50,
			quorumNumbers: []uint8{0, 2},
			timestamp:     now,
			wantErr:       true,
			errMsg:        "quorum number mismatch",
		},
		{
			name:          "reservation limit exceeded",
			numSymbols:    200,
			quorumNumbers: []uint8{0},
			timestamp:     now,
			wantErr:       true,
			errMsg:        "exceeds bin limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := acc.reservationUsage(tt.numSymbols, tt.quorumNumbers, tt.timestamp)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccountant_UseOnDemand(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now().UnixNano()
	acc := NewAccountant(accountID)
	err := acc.SetPaymentState(&v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams: createTestPaymentVaultParams(6, 1, 100),
		Reservations: map[uint32]*v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint32(now/1e9 - 1000),
				EndTimestamp:     uint32(now/1e9 - 500),
			},
		},
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
	})
	require.NoError(t, err)
	// Set up payment configs
	acc.paymentVaultParams.QuorumPaymentConfigs[0] = &core.PaymentQuorumConfig{
		OnDemandPricePerSymbol: 10,
	}
	acc.paymentVaultParams.QuorumProtocolConfigs[0] = &core.PaymentQuorumProtocolConfig{
		MinNumSymbols: 1,
	}
	acc.paymentVaultParams.OnDemandQuorumNumbers = []uint8{0, 1}

	tests := []struct {
		name          string
		numSymbols    uint64
		quorumNumbers []uint8
		timestamp     int64
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "successful on-demand payment",
			numSymbols:    50,
			quorumNumbers: []uint8{0, 1},
			timestamp:     now,
			wantErr:       false,
		},
		{
			name:          "insufficient balance",
			numSymbols:    200,
			quorumNumbers: []uint8{0, 1},
			timestamp:     now,
			wantErr:       true,
			errMsg:        "current cumulativePayment balance insufficient",
		},
		{
			name:          "invalid quorum",
			numSymbols:    50,
			quorumNumbers: []uint8{2},
			timestamp:     now,
			wantErr:       true,
			errMsg:        "quorum number mismatch: 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, err := acc.onDemandUsage(tt.numSymbols, tt.quorumNumbers)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, payment)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, payment)
				assert.True(t, payment.Cmp(big.NewInt(0)) > 0)
			}
		})
	}
}

func TestAccountant_MultipleOverflows(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now()
	acc := NewAccountant(accountID)
	err := acc.SetPaymentState(&v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams: createTestPaymentVaultParams(6, 1, 100),
		Reservations: map[uint32]*v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint32(now.Add(time.Second * -1000).Unix()),
				EndTimestamp:     uint32(now.Add(time.Second * 1000).Unix()),
			},
		},
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
	})
	require.NoError(t, err)
	// Set up protocol configs
	binInterval := uint64(1)
	acc.paymentVaultParams.QuorumProtocolConfigs[0] = &core.PaymentQuorumProtocolConfig{
		ReservationRateLimitWindow: binInterval,
		OnDemandRateLimitWindow:    binInterval,
		MinNumSymbols:              1,
	}

	acc.paymentVaultParams.QuorumPaymentConfigs[0] = &core.PaymentQuorumConfig{
		ReservationSymbolsPerSecond: 0,
		OnDemandSymbolsPerSecond:    100,
		OnDemandPricePerSymbol:      1,
	}
	acc.paymentVaultParams.OnDemandQuorumNumbers = []uint8{0}

	// First call: Use current bin
	_, err = acc.AccountBlob(now.UnixNano(), 50, []uint8{0})
	require.NoError(t, err)
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now.UnixNano(), 1)
	record := acc.periodRecords.GetRelativePeriodRecord(currentPeriod, 0)
	assert.Equal(t, uint64(50), record.Usage)
	record = acc.periodRecords.GetRelativePeriodRecord(currentPeriod+1*binInterval, 0)
	assert.Equal(t, uint64(0), record.Usage)
	record = acc.periodRecords.GetRelativePeriodRecord(meterer.GetOverflowPeriod(currentPeriod, binInterval), 0)
	assert.Equal(t, uint64(0), record.Usage)

	// Second call: Overflow
	_, err = acc.AccountBlob(now.UnixNano(), 100, []uint8{0})
	require.NoError(t, err)
	record = acc.periodRecords.GetRelativePeriodRecord(currentPeriod, 0)
	assert.Equal(t, uint64(100), record.Usage)
	record = acc.periodRecords.GetRelativePeriodRecord(currentPeriod+1*binInterval, 0)
	assert.Equal(t, uint64(0), record.Usage)
	record = acc.periodRecords.GetRelativePeriodRecord(meterer.GetOverflowPeriod(currentPeriod, binInterval), 0)
	assert.Equal(t, uint64(50), record.Usage)

	// Third call: Cannot overflow again
	_, err = acc.AccountBlob(now.UnixNano(), 100, []uint8{0})
	require.NoError(t, err)
	record = acc.periodRecords.GetRelativePeriodRecord(currentPeriod, 0)
	assert.Equal(t, uint64(100), record.Usage)
	record = acc.periodRecords.GetRelativePeriodRecord(currentPeriod+1*binInterval, 0)
	assert.Equal(t, uint64(0), record.Usage)
	record = acc.periodRecords.GetRelativePeriodRecord(meterer.GetOverflowPeriod(currentPeriod, binInterval), 0)
	assert.Equal(t, uint64(50), record.Usage)
}

func TestAccountant_MixedReservationStates(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now()
	acc := NewAccountant(accountID)

	// Set up protocol configs
	vaultParams := &v2.PaymentVaultParams{
		QuorumProtocolConfigs: make(map[uint32]*v2.PaymentQuorumProtocolConfig),
		QuorumPaymentConfigs:  make(map[uint32]*v2.PaymentQuorumConfig),
	}
	for i := uint32(0); i < 3; i++ {
		// Set up protocol configs
		vaultParams.QuorumProtocolConfigs[i] = &v2.PaymentQuorumProtocolConfig{
			ReservationRateLimitWindow: 1,
			OnDemandRateLimitWindow:    1,
			MinNumSymbols:              1,
		}

		vaultParams.QuorumPaymentConfigs[i] = &v2.PaymentQuorumConfig{
			ReservationSymbolsPerSecond: 0,
			OnDemandSymbolsPerSecond:    100,
			OnDemandPricePerSymbol:      1,
		}
	}
	err := acc.SetPaymentState(&v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams: vaultParams,
		Reservations: map[uint32]*v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint32(now.Add(time.Hour * 1).Unix()),
				EndTimestamp:     uint32(now.Add(time.Hour * 2).Unix()), // Future
			},
			1: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint32(now.Add(time.Hour * -2).Unix()),
				EndTimestamp:     uint32(now.Add(time.Hour * -1).Unix()), // Expired
			},
			2: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint32(now.Add(time.Hour * -1).Unix()),
				EndTimestamp:     uint32(now.Add(time.Hour * 1).Unix()), // Active
			},
		},
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
	})
	require.NoError(t, err)
	acc.paymentVaultParams.OnDemandQuorumNumbers = []uint8{0, 1}

	// Reservations and OnDemand are not sufficient for all three quorums
	payment, err := acc.AccountBlob(now.UnixNano(), 50, []uint8{0, 1, 2})
	assert.Nil(t, payment)
	assert.Contains(t, err.Error(), "cannot create payment information")

	// Separate reservation dispersal is sufficient for quorum 2
	payment, err = acc.AccountBlob(now.UnixNano(), 50, []uint8{2})
	// 1749697512 1749701112 1749693912.770014000
	require.NoError(t, err)
	assert.NotNil(t, payment)
	assert.True(t, payment.CumulativePayment.Cmp(big.NewInt(0)) == 0)

	// Alternatively use ondemand for quorum 0 or/and 1
	payment, err = acc.AccountBlob(now.UnixNano(), 50, []uint8{0, 1})
	require.NoError(t, err)
	assert.NotNil(t, payment)
	assert.True(t, payment.CumulativePayment.Cmp(big.NewInt(0)) > 0)
}

func TestAccountant_ReservationRollback(t *testing.T) {
	reservation := &v2.QuorumReservation{
		SymbolsPerSecond: 50,
		StartTimestamp:   uint32(time.Now().Unix()),
		EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
	}
	reservationWindow := uint64(2)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)

	// Create payment state with test configurations
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams: &v2.PaymentVaultParams{
			QuorumPaymentConfigs: map[uint32]*v2.PaymentQuorumConfig{
				0: {
					ReservationSymbolsPerSecond: 100,
					OnDemandSymbolsPerSecond:    100,
					OnDemandPricePerSymbol:      1,
				},
				1: {
					ReservationSymbolsPerSecond: 100,
					OnDemandSymbolsPerSecond:    100,
					OnDemandPricePerSymbol:      1,
				},
			},
			QuorumProtocolConfigs: map[uint32]*v2.PaymentQuorumProtocolConfig{
				0: {
					MinNumSymbols:              1,
					ReservationAdvanceWindow:   10,
					ReservationRateLimitWindow: reservationWindow,
					OnDemandRateLimitWindow:    30,
					OnDemandEnabled:            true,
				},
				1: {
					MinNumSymbols:              1,
					ReservationAdvanceWindow:   10,
					ReservationRateLimitWindow: reservationWindow,
					OnDemandRateLimitWindow:    30,
					OnDemandEnabled:            true,
				},
			},
			OnDemandQuorumNumbers: []uint32{0, 1},
		},
		OnchainCumulativePayment: big.NewInt(0).Bytes(),
		Reservations:             map[uint32]*v2.QuorumReservation{0: reservation, 1: reservation},
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	// Test rollback when a later quorum fails
	now := time.Now().UnixNano()
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)

	// First update should succeed
	moreUsedQuorum := uint8(1)
	lessUsedQuorum := uint8(0)
	_, err = accountant.AccountBlob(now, 50, []uint8{moreUsedQuorum})
	assert.NoError(t, err)

	// Verify first quorum was updated
	record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, moreUsedQuorum)
	assert.Equal(t, uint64(50), record.Usage)

	// Use both quorums, more used quorum overflows
	_, err = accountant.AccountBlob(now, 60, []uint8{moreUsedQuorum, lessUsedQuorum})
	assert.NoError(t, err)
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, moreUsedQuorum)
	assert.Equal(t, uint64(100), record.Usage)
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, lessUsedQuorum)
	assert.Equal(t, uint64(60), record.Usage)
	record = accountant.periodRecords.GetRelativePeriodRecord(meterer.GetOverflowPeriod(currentPeriod, reservationWindow), moreUsedQuorum)
	assert.Equal(t, uint64(10), record.Usage)

	// Use both quorums, more used quorum cannot overflow again
	_, err = accountant.AccountBlob(now, 60, []uint8{moreUsedQuorum, lessUsedQuorum})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reservation limit exceeded")

	// No reservation updates were made
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, moreUsedQuorum)
	assert.Equal(t, uint64(100), record.Usage)
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, lessUsedQuorum)
	assert.Equal(t, uint64(60), record.Usage)
	record = accountant.periodRecords.GetRelativePeriodRecord(meterer.GetOverflowPeriod(currentPeriod, reservationWindow), moreUsedQuorum)
	assert.Equal(t, uint64(10), record.Usage)

	// Test rollback when a quorum doesn't exist
	_, err = accountant.AccountBlob(now, 50, []uint8{lessUsedQuorum, 2})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quorum number mismatch")
	// quorum usage rolled back
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, lessUsedQuorum)
	assert.Equal(t, uint64(60), record.Usage)

	// Test rollback when config is missing
	// Remove config for quorum 1
	delete(accountant.paymentVaultParams.QuorumProtocolConfigs, 1)
	_, err = accountant.AccountBlob(now, 50, []uint8{0, 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quorum config not found")
	// quorum usage stays the same
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, lessUsedQuorum)
	assert.Equal(t, uint64(60), record.Usage)
}
