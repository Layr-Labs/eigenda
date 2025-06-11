package clients

import (
	"context"
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

const numBins = uint32(3)

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
	}
}

func TestNewAccountant(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	reservations := map[uint8]*core.QuorumReservation{
		0: {
			SymbolsPerSecond: 100,
			StartTimestamp:   uint32(time.Now().Unix()),
			EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
		},
		1: {
			SymbolsPerSecond: 200,
			StartTimestamp:   uint32(time.Now().Unix()),
			EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	numBins := uint32(3)

	acc := NewAccountant(accountID, reservations, onDemand, numBins)

	require.NotNil(t, acc)
	assert.Equal(t, accountID, acc.accountID)
	assert.Equal(t, reservations, acc.reservations)
	assert.Equal(t, onDemand, acc.onDemand)
	assert.Equal(t, numBins, acc.numBins)
	assert.NotNil(t, acc.periodRecords)
	assert.NotNil(t, acc.quorumPaymentConfigs)
	assert.NotNil(t, acc.quorumProtocolConfigs)
}

func TestAccountBlob_Reservation(t *testing.T) {
	reservation := &core.QuorumReservation{
		SymbolsPerSecond: 200,
		StartTimestamp:   uint32(time.Now().Unix()),
		EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}
	reservationWindow := uint64(5)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.QuorumReservation{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins)

	// Create payment state with test configurations
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(500).Bytes(),
		Reservations:             map[uint32]*core.QuorumReservation{0: reservation, 1: reservation},
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	ctx := context.Background()
	symbolLength := uint64(500)
	quorums := []uint8{0, 1}
	now := time.Now().UnixNano()

	header, err := accountant.AccountBlob(ctx, now, symbolLength, []uint8{0, 1, 2})
	assert.Nil(t, header)
	assert.Error(t, err, "no reservation found for quorum")

	header, err = accountant.AccountBlob(ctx, now, symbolLength, quorums)
	assert.NoError(t, err)
	window, err := accountant.GetReservationWindow(0)
	assert.NoError(t, err)
	assert.Equal(t, meterer.GetReservationPeriod(time.Now().Unix(), window), meterer.GetReservationPeriodByNanosecond(header.Timestamp, window))
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Check that usage for quorum 0 was updated properly
	quorum0Records := accountant.periodRecords[0]
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, window) % uint64(numBins)
	assert.Equal(t, symbolLength, quorum0Records[currentPeriod].Usage)

	symbolLength = uint64(700)

	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(ctx, now, symbolLength, quorums)

	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// With overflow, usage should be at the limit in the current period
	binLimit := reservations[0].SymbolsPerSecond * uint64(window)
	expectedUsage := binLimit
	assert.Equal(t, expectedUsage, quorum0Records[currentPeriod].Usage)

	// The overflow should be in the next bin
	overflowIndex := uint32((meterer.GetReservationPeriodByNanosecond(now, window) + 2) % uint64(numBins))
	expectedOverflow := symbolLength - (binLimit - 500) // 700 - (1000 - 500) = 200
	assert.Equal(t, expectedOverflow, quorum0Records[overflowIndex].Usage)

	// Second call should use on-demand payment
	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(ctx, now, 300, quorums)

	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(300), header.CumulativePayment)
}

func TestAccountBlob_OnDemand(t *testing.T) {
	numSymbols := uint64(1500)
	quorums := []uint8{0, 1}

	reservations := map[uint8]*core.QuorumReservation{
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
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1500),
	}
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservations, onDemand, numBins)

	// Create payment state with test configurations
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(5, 1, 100),
		OnchainCumulativePayment: big.NewInt(1500).Bytes(),
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	ctx := context.Background()
	now := time.Now().UnixNano()
	header, err := accountant.AccountBlob(ctx, now, numSymbols, quorums)
	assert.NoError(t, err)

	pricePerSymbol, err := accountant.GetPricePerSymbol(0)
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

func TestAccountBlob_InsufficientOnDemand(t *testing.T) {
	reservation := map[uint8]*core.QuorumReservation{} // Empty reservation map
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, numBins)

	// Create payment state with test configurations
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       createTestPaymentVaultParams(60, 100, 100),
		OnchainCumulativePayment: big.NewInt(500).Bytes(),
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	ctx := context.Background()
	numSymbols := uint64(2000)
	quorums := []uint8{0, 1}
	baseTime := time.Now().UnixNano()

	tests := []accountBlobInsufficientOnDemandTest{
		{
			name:         "Insufficient on-demand payment",
			symbolLength: 2000,
			expectError:  true,
			errorMessage: "insufficient ondemand payment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := baseTime
			_, err := accountant.AccountBlob(now, tt.symbolLength, quorums)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type accountBlobCallSeriesTest struct {
	name           string
	symbolLength   uint64
	expectedHeader *core.PaymentMetadata
	expectedState  periodRecordState
	expectError    bool
	errorMessage   string
}

func TestAccountBlobCallSeries(t *testing.T) {
	reservation := &core.QuorumReservation{
		SymbolsPerSecond: 200,
		StartTimestamp:   uint32(time.Now().Unix()),
		EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(5)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.QuorumReservation{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins)

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

	ctx := context.Background()
	quorums := []uint8{0, 1}

	tests := []accountBlobCallSeriesTest{
		{
			name:         "First call - Use reservation",
			symbolLength: 800,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 800,
				NextPeriodUsage:    0,
				OverflowUsage:      0,
			},
		},
		{
			name:         "Second call - Use remaining reservation + overflow",
			symbolLength: 300,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 1100,
				NextPeriodUsage:    0,
				OverflowUsage:      100,
			},
		},
		{
			name:         "Third call - Use on-demand",
			symbolLength: 500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(500),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 1100,
				NextPeriodUsage:    0,
				OverflowUsage:      100,
			},
		},
		{
			name:         "Fourth call - Insufficient on-demand",
			symbolLength: 600,
			expectError:  true,
			errorMessage: "insufficient ondemand payment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := baseTime.UnixNano()
			header, err := accountant.AccountBlob(now, tt.symbolLength, quorums)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uint64(0), header.Timestamp)
				assert.Equal(t, tt.expectedHeader.AccountID, header.AccountID)
				assert.Equal(t, tt.expectedHeader.CumulativePayment, header.CumulativePayment)

				period := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
				assert.Equal(t, tt.expectedState.CurrentPeriodUsage, getRelativePeriodRecord(period, reservationWindow, accountant.periodRecords).Usage)
				assert.Equal(t, tt.expectedState.NextPeriodUsage, getRelativePeriodRecord(period+reservationWindow, reservationWindow, accountant.periodRecords).Usage)
				assert.Equal(t, tt.expectedState.OverflowUsage, getRelativePeriodRecord(meterer.GetOverflowPeriod(period, reservationWindow), reservationWindow, accountant.periodRecords).Usage)
			}
		})
	}
}

type accountBlobBinRotationTest struct {
	name           string
	timeOffset     time.Duration
	symbolLength   uint64
	expectedHeader *core.PaymentMetadata
	expectedState  periodRecordState
}

func TestAccountBlob_BinRotation(t *testing.T) {
	reservation := &core.QuorumReservation{
		SymbolsPerSecond: 1000,
		StartTimestamp:   uint32(time.Now().Unix()),
		EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(1) // Set to 1 second for testing

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.QuorumReservation{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins)

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

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// First call
	now := time.Now().UnixNano()
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(ctx, now, 800, quorums)
	assert.NoError(t, err)

	// Check bin 0 has usage 800 for charged quorums
	for _, quorumNumber := range quorums {
		relativeIndex := uint32(currentPeriod % uint64(numBins))
		assert.Equal(t, uint64(800), accountant.periodRecords[quorumNumber][relativeIndex].Usage)
	}

	// Second call
	now += int64(reservationWindow) * time.Second.Nanoseconds()
	nextPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(ctx, now, 300, quorums)
	assert.NoError(t, err)

	// Check bin 1 has usage 300
	for _, quorumNumber := range quorums {
		relativeIndex := uint32(currentPeriod % uint64(numBins))
		assert.Equal(t, uint64(800), accountant.periodRecords[quorumNumber][relativeIndex].Usage)
		relativeIndex = uint32(nextPeriod % uint64(numBins))
		assert.Equal(t, uint64(300), accountant.periodRecords[quorumNumber][relativeIndex].Usage)
	}

	// Third call
	_, err = accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)

	// Check bin 1 now has usage 800
	for _, quorumNumber := range quorums {
		relativeIndex := uint32(currentPeriod % uint64(numBins))
		assert.Equal(t, uint64(800), accountant.periodRecords[quorumNumber][relativeIndex].Usage)
		relativeIndex = uint32(nextPeriod % uint64(numBins))
		assert.Equal(t, uint64(800), accountant.periodRecords[quorumNumber][relativeIndex].Usage)
	}
}

func TestConcurrentBinRotationAndAccountBlob(t *testing.T) {
	reservation := &core.QuorumReservation{
		SymbolsPerSecond: 1000,
		StartTimestamp:   uint32(time.Now().Unix()),
		EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(1) // Set to 1 second for testing

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.QuorumReservation{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins)

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

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// Start concurrent AccountBlob calls
	now := time.Now().UnixNano()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := accountant.AccountBlob(ctx, now, 100, quorums)
			assert.NoError(t, err)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check final state
	for _, quorumNumber := range quorums {
		currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
		relativeIndex := uint32(currentPeriod % uint64(numBins))
		assert.Equal(t, uint64(1000), accountant.periodRecords[quorumNumber][relativeIndex].Usage)
	}
}

func TestAccountBlob_ReservationWithOneOverflow(t *testing.T) {
	reservation := &core.QuorumReservation{
		SymbolsPerSecond: 200,
		StartTimestamp:   uint32(time.Now().Unix()),
		EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(5)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.QuorumReservation{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins)

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

	ctx := context.Background()
	quorums := []uint8{0, 1}
	now := time.Now().UnixNano()
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	header, err := accountant.AccountBlob(ctx, now, 800, quorums)
	assert.NoError(t, err)
	timestamp := (time.Duration(header.Timestamp) * time.Nanosecond).Seconds()
	assert.Equal(t, uint64(meterer.GetReservationPeriodByNanosecond(now, reservationWindow)), uint64(meterer.GetReservationPeriod(int64(timestamp), reservationWindow)))
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Check current period usage
	for _, quorumNumber := range quorums {
		record := accountant.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(800), record.Usage)
	}
	// Second call: Allow one overflow
	header, err = accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Check current period is at limit
	binLimit := reservations[0].SymbolsPerSecond * uint64(reservationWindow) // 1000
	for _, quorumNumber := range quorums {
		record := accountant.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, binLimit, record.Usage)
	}

	// Check overflow period has the overflow
	for _, quorumNumber := range quorums {
		overflowRecord := accountant.GetRelativePeriodRecord(currentPeriod+2, quorumNumber)
		assert.Equal(t, uint64(300), overflowRecord.Usage) // 800 + 500 - 1000 = 300
	}

	// Third call: Should use on-demand payment
	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(ctx, now, 200, quorums)
	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(200), header.CumulativePayment)
	for _, quorumNumber := range quorums {
		record := accountant.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(binLimit), record.Usage)
		record = accountant.GetRelativePeriodRecord(currentPeriod+2, quorumNumber)
		assert.Equal(t, uint64(300), record.Usage)
	}
}

func TestAccountBlob_ReservationOverflowReset(t *testing.T) {
	reservation := &core.QuorumReservation{
		SymbolsPerSecond: 1000,
		StartTimestamp:   uint32(time.Now().Unix()),
		EndTimestamp:     uint32(time.Now().Add(time.Hour).Unix()),
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(1) // Set to 1 second for testing

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.QuorumReservation{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins)

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

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// full reservation
	now := time.Now().UnixNano()
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(ctx, now, 1000, quorums)
	assert.NoError(t, err)

	// Check current period is at limit
	for _, quorumNumber := range quorums {
		record := accountant.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(1000), record.Usage)
	}

	// no overflow
	now = time.Now().UnixNano()
	header, err := accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)
	for _, quorumNumber := range quorums {
		record := accountant.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(1000), record.Usage)
	}
	assert.Equal(t, big.NewInt(500), header.CumulativePayment)

	// Wait for next reservation duration
	time.Sleep(time.Duration(reservationWindow) * time.Second)

	// Third call: Should use new bin and allow overflow again
	now = time.Now().UnixNano()
	nextPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	assert.Equal(t, currentPeriod+1, nextPeriod, "Should be next period")

	_, err = accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)

	// Check next period has usage 500
	for _, quorumNumber := range quorums {
		record := accountant.GetRelativePeriodRecord(nextPeriod, quorumNumber)
		assert.Equal(t, uint64(500), record.Usage)
	}
}

func TestQuorumCheck(t *testing.T) {
	tests := []struct {
		name           string
		quorumNumbers  []uint8
		allowedNumbers []uint8
		expectError    bool
		errorMessage   string
	}{
		{
			name:           "valid quorum numbers",
			quorumNumbers:  []uint8{0, 1},
			allowedNumbers: []uint8{0, 1, 2},
			expectError:    false,
		},
		{
			name:           "empty quorum numbers",
			quorumNumbers:  []uint8{},
			allowedNumbers: []uint8{0, 1},
			expectError:    true,
			errorMessage:   "no quorum numbers provided in the request",
		},
		{
			name:           "invalid quorum number",
			quorumNumbers:  []uint8{0, 2},
			allowedNumbers: []uint8{0, 1},
			expectError:    true,
			errorMessage:   "quorum number mismatch: 2",
		},
		{
			name:           "empty allowed numbers",
			quorumNumbers:  []uint8{0},
			allowedNumbers: []uint8{},
			expectError:    true,
			errorMessage:   "quorum number mismatch: 0",
		},
		{
			name:           "multiple invalid quorums",
			quorumNumbers:  []uint8{2, 3, 4},
			allowedNumbers: []uint8{0, 1},
			expectError:    true,
			errorMessage:   "quorum number mismatch: 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := meterer.ValidateQuorum(tt.quorumNumbers, tt.allowedNumbers)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetPaymentState(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	acc := NewAccountant(accountID, nil, nil, 3)

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
				// Verify state was set correctly
				assert.NotNil(t, acc.quorumPaymentConfigs[0])
				assert.NotNil(t, acc.quorumProtocolConfigs[0])
				assert.NotNil(t, acc.reservations[0])
			}
		})
	}
}

func TestGetMinNumSymbols(t *testing.T) {
	accountant := NewAccountant(gethcommon.Address{}, map[uint8]*core.QuorumReservation{}, &core.OnDemandPayment{}, numBins)

	// Test with non-existent quorum
	_, err := accountant.GetMinNumSymbols(99)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in protocol configs")

	// Test with existing quorum after setting payment state
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams: createTestPaymentVaultParams(6, 1, 100),
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	minSymbols, err := accountant.GetMinNumSymbols(0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), minSymbols)
}

func TestGetPricePerSymbol(t *testing.T) {
	accountant := NewAccountant(gethcommon.Address{}, map[uint8]*core.QuorumReservation{}, &core.OnDemandPayment{}, numBins)

	// Test with non-existent quorum
	_, err := accountant.GetPricePerSymbol(99)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in payment configs")

	// Test with existing quorum after setting payment state
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams: createTestPaymentVaultParams(6, 1, 100),
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	price, err := accountant.GetPricePerSymbol(0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), price)
}

func TestGetReservationWindow(t *testing.T) {
	accountant := NewAccountant(gethcommon.Address{}, map[uint8]*core.QuorumReservation{}, &core.OnDemandPayment{}, numBins)

	// Test with non-existent quorum
	_, err := accountant.GetReservationWindow(99)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in protocol configs")

	// Test with existing quorum after setting payment state
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams: createTestPaymentVaultParams(6, 1, 100),
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	window, err := accountant.GetReservationWindow(0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(6), window)
}

func TestPaymentCharged(t *testing.T) {
	accountant := NewAccountant(gethcommon.Address{}, map[uint8]*core.QuorumReservation{}, &core.OnDemandPayment{}, numBins)

	// Test with non-existent quorum
	_, err := accountant.PaymentCharged(100, 99)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in payment configs")

	// Test with existing quorum after setting payment state
	paymentState := &v2.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams: createTestPaymentVaultParams(6, 1, 100),
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	// Test with numSymbols less than minNumSymbols
	payment, err := accountant.PaymentCharged(50, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), payment) // Should use minNumSymbols (100) * pricePerSymbol (1)

	// Test with numSymbols greater than minNumSymbols
	payment, err = accountant.PaymentCharged(150, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(200), payment) // Should round up to 200 * pricePerSymbol (1)
}

func TestBlobPaymentInfo_UseReservation(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now().UnixNano()
	reservations := map[uint8]*core.QuorumReservation{
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
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	acc := NewAccountant(accountID, reservations, onDemand, 3)

	// Set up protocol configs
	acc.quorumProtocolConfigs[0] = &core.PaymentQuorumProtocolConfig{
		ReservationRateLimitWindow: 1,
		MinNumSymbols:              1,
	}
	acc.quorumProtocolConfigs[1] = &core.PaymentQuorumProtocolConfig{
		ReservationRateLimitWindow: 1,
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
			errMsg:        "No reservation found on quorum 2",
		},
		{
			name:          "reservation limit exceeded",
			numSymbols:    200,
			quorumNumbers: []uint8{0},
			timestamp:     now,
			wantErr:       true,
			errMsg:        "reservation limit exceeded for quorum 0",
		},
		{
			name:          "empty quorum numbers",
			numSymbols:    50,
			quorumNumbers: []uint8{},
			timestamp:     now,
			wantErr:       true,
			errMsg:        "no quorum numbers provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, err := acc.ReservationUsage(tt.numSymbols, tt.quorumNumbers, tt.timestamp)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, payment)
			} else {
				require.NoError(t, err)
				assert.Equal(t, big.NewInt(0), payment)
			}
		})
	}
}

func TestBlobPaymentInfo_UseOnDemand(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now().UnixNano()
	reservations := map[uint8]*core.QuorumReservation{
		0: {
			SymbolsPerSecond: 100,
			StartTimestamp:   uint32(now/1e9 - 1000),
			EndTimestamp:     uint32(now/1e9 - 500), // Expired reservation
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	acc := NewAccountant(accountID, reservations, onDemand, 3)

	// Set up payment configs
	acc.quorumPaymentConfigs[0] = &core.PaymentQuorumConfig{
		OnDemandPricePerSymbol: 10,
	}
	acc.quorumProtocolConfigs[0] = &core.PaymentQuorumProtocolConfig{
		MinNumSymbols: 1,
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
			errMsg:        "provided quorum number 2 not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, err := acc.OnDemandUsage(tt.numSymbols, tt.quorumNumbers)
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

func TestProcessQuorumReservation(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now().UnixNano()
	reservations := map[uint8]*core.QuorumReservation{
		0: {
			SymbolsPerSecond: 100,
			StartTimestamp:   uint32(now/1e9 - 1000),
			EndTimestamp:     uint32(now/1e9 + 1000),
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	acc := NewAccountant(accountID, reservations, onDemand, 3)

	// Set up protocol configs
	acc.quorumProtocolConfigs[0] = &core.PaymentQuorumProtocolConfig{
		ReservationRateLimitWindow: 1,
	}

	tests := []struct {
		name             string
		quorumNumber     uint8
		reservation      *core.QuorumReservation
		currentPeriod    uint64
		symbolUsage      uint64
		initialUsage     uint64
		wantErr          bool
		errMsg           string
		expectedUsage    uint64
		expectedOverflow uint64
	}{
		{
			name:          "within bin limit",
			quorumNumber:  0,
			reservation:   reservations[0],
			currentPeriod: 1,
			symbolUsage:   50,
			initialUsage:  0,
			wantErr:       false,
			expectedUsage: 50,
		},
		{
			name:          "exact bin limit",
			quorumNumber:  0,
			reservation:   reservations[0],
			currentPeriod: 1,
			symbolUsage:   100,
			initialUsage:  0,
			wantErr:       false,
			expectedUsage: 100,
		},
		{
			name:             "overflow bin usage",
			quorumNumber:     0,
			reservation:      reservations[0],
			currentPeriod:    1,
			symbolUsage:      50,
			initialUsage:     80,
			wantErr:          false,
			expectedUsage:    100,
			expectedOverflow: 30,
		},
		{
			name:          "exceeds limit",
			quorumNumber:  0,
			reservation:   reservations[0],
			currentPeriod: 1,
			symbolUsage:   150,
			initialUsage:  0,
			wantErr:       true,
			errMsg:        "reservation limit exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set initial usage
			periodRecord := acc.GetRelativePeriodRecord(tt.currentPeriod, tt.quorumNumber)
			periodRecord.Usage = tt.initialUsage

			err := acc.processQuorumReservation(tt.quorumNumber, tt.reservation, tt.currentPeriod, tt.symbolUsage)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedUsage, periodRecord.Usage)
				if tt.expectedOverflow > 0 {
					overflowRecord := acc.GetRelativePeriodRecord(tt.currentPeriod+2, tt.quorumNumber)
					assert.Equal(t, tt.expectedOverflow, overflowRecord.Usage)
				}
			}
		})
	}
}

func TestBlobPaymentInfo_FutureReservation(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now().UnixNano()
	reservations := map[uint8]*core.QuorumReservation{
		0: {
			SymbolsPerSecond: 100,
			StartTimestamp:   uint32(now/1e9 + 1000), // Future reservation
			EndTimestamp:     uint32(now/1e9 + 2000),
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	acc := NewAccountant(accountID, reservations, onDemand, 3)

	// Set up protocol configs
	acc.quorumProtocolConfigs[0] = &core.PaymentQuorumProtocolConfig{
		ReservationRateLimitWindow: 1,
		OnDemandRateLimitWindow:    1,
		MinNumSymbols:              1,
	}

	acc.quorumPaymentConfigs[0] = &core.PaymentQuorumConfig{
		ReservationSymbolsPerSecond: 0,
		OnDemandSymbolsPerSecond:    100,
		OnDemandPricePerSymbol:      1,
	}

	// Should use on-demand since reservation is in future
	payment, err := acc.BlobPaymentInfo(context.Background(), 50, []uint8{0}, now)
	require.NoError(t, err)
	assert.NotNil(t, payment)
	assert.True(t, payment.Cmp(big.NewInt(0)) > 0)
}

func TestBlobPaymentInfo_MultipleOverflows(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now().UnixNano()
	reservations := map[uint8]*core.QuorumReservation{
		0: {
			SymbolsPerSecond: 100,
			StartTimestamp:   uint32(now/1e9 - 1000),
			EndTimestamp:     uint32(now/1e9 + 1000),
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	acc := NewAccountant(accountID, reservations, onDemand, 3)

	// Set up protocol configs
	binInterval := uint64(1)
	acc.quorumProtocolConfigs[0] = &core.PaymentQuorumProtocolConfig{
		ReservationRateLimitWindow: binInterval,
		OnDemandRateLimitWindow:    binInterval,
		MinNumSymbols:              1,
	}

	acc.quorumPaymentConfigs[0] = &core.PaymentQuorumConfig{
		ReservationSymbolsPerSecond: 0,
		OnDemandSymbolsPerSecond:    100,
		OnDemandPricePerSymbol:      1,
	}

	// First call: Use current bin
	ctx := context.Background()
	_, err := acc.BlobPaymentInfo(ctx, 50, []uint8{0}, now)
	require.NoError(t, err)
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, 1)
	assert.Equal(t, uint64(50), acc.periodRecords[0][currentPeriod%3].Usage)
	assert.Equal(t, uint64(0), acc.periodRecords[0][(currentPeriod+1)%3].Usage)
	assert.Equal(t, uint64(0), acc.periodRecords[0][(currentPeriod+2)%3].Usage)

	// Second call: Overflow
	_, err = acc.BlobPaymentInfo(ctx, 100, []uint8{0}, now)
	require.NoError(t, err)
	assert.Equal(t, uint64(100), acc.periodRecords[0][currentPeriod%3].Usage)
	assert.Equal(t, uint64(0), acc.periodRecords[0][(currentPeriod+1)%3].Usage)
	assert.Equal(t, uint64(50), acc.periodRecords[0][(currentPeriod+2)%3].Usage)

	// Third call: Cannot overflow again
	_, err = acc.BlobPaymentInfo(ctx, 100, []uint8{0}, now)
	require.NoError(t, err)
	assert.Equal(t, uint64(100), acc.periodRecords[0][currentPeriod%3].Usage)
	assert.Equal(t, uint64(0), acc.periodRecords[0][(currentPeriod+1)%3].Usage)
	assert.Equal(t, uint64(50), acc.periodRecords[0][(currentPeriod+2)%3].Usage)
}

func TestBlobPaymentInfo_MixedReservationStates(t *testing.T) {
	ctx := context.Background()
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now()
	reservations := map[uint8]*core.QuorumReservation{
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
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	acc := NewAccountant(accountID, reservations, onDemand, 3)

	// Set up protocol configs
	for i := uint8(0); i < 3; i++ {
		// Set up protocol configs
		acc.quorumProtocolConfigs[i] = &core.PaymentQuorumProtocolConfig{
			ReservationRateLimitWindow: 1,
			OnDemandRateLimitWindow:    1,
			MinNumSymbols:              1,
		}

		acc.quorumPaymentConfigs[i] = &core.PaymentQuorumConfig{
			ReservationSymbolsPerSecond: 0,
			OnDemandSymbolsPerSecond:    100,
			OnDemandPricePerSymbol:      1,
		}
	}

	// Reservations and OnDemand are not sufficient for all three quorums
	payment, err := acc.BlobPaymentInfo(ctx, 50, []uint8{0, 1, 2}, now.UnixNano())
	assert.Nil(t, payment)
	assert.Contains(t, err.Error(), "not allowed")

	// Separate reservation dispersal is sufficient for quorum 2
	payment, err = acc.BlobPaymentInfo(ctx, 50, []uint8{2}, now.UnixNano())
	require.NoError(t, err)
	assert.NotNil(t, payment)
	assert.True(t, payment.Cmp(big.NewInt(0)) == 0)

	// Alternatively use ondemand for quorum 0 or/and 1
	payment, err = acc.BlobPaymentInfo(ctx, 50, []uint8{0, 1}, now.UnixNano())
	require.NoError(t, err)
	assert.NotNil(t, payment)
	assert.True(t, payment.Cmp(big.NewInt(0)) > 0)

}

func TestBlobPaymentInfo_ZeroPayment(t *testing.T) {
	ctx := context.Background()
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now().UnixNano()
	reservations := map[uint8]*core.QuorumReservation{
		0: {
			SymbolsPerSecond: 100,
			StartTimestamp:   uint32(now/1e9 - 1000),
			EndTimestamp:     uint32(now/1e9 + 1000),
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	acc := NewAccountant(accountID, reservations, onDemand, 3)

	// Set up protocol configs
	acc.quorumProtocolConfigs[0] = &core.PaymentQuorumProtocolConfig{
		ReservationRateLimitWindow: 1,
		OnDemandRateLimitWindow:    1,
		MinNumSymbols:              1,
	}

	acc.quorumPaymentConfigs[0] = &core.PaymentQuorumConfig{
		ReservationSymbolsPerSecond: 0,
		OnDemandSymbolsPerSecond:    100,
		OnDemandPricePerSymbol:      1,
	}

	// Should return zero payment for zero symbols
	payment, err := acc.BlobPaymentInfo(ctx, 0, []uint8{0}, now)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), payment)

	// Should return zero payment for empty quorum list
	payment, err = acc.BlobPaymentInfo(ctx, 50, []uint8{}, now)
	require.Error(t, err)
	assert.Nil(t, payment)
}

func TestBlobPaymentInfo_MaximumPayment(t *testing.T) {
	ctx := context.Background()
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now().UnixNano()
	reservations := map[uint8]*core.QuorumReservation{
		0: {
			SymbolsPerSecond: 100,
			StartTimestamp:   uint32(now/1e9 - 1000),
			EndTimestamp:     uint32(now/1e9 + 1000),
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	acc := NewAccountant(accountID, reservations, onDemand, 3)

	// Set up protocol configs
	acc.quorumProtocolConfigs[0] = &core.PaymentQuorumProtocolConfig{
		ReservationRateLimitWindow: 1,
		OnDemandRateLimitWindow:    1,
		MinNumSymbols:              1,
	}

	acc.quorumPaymentConfigs[0] = &core.PaymentQuorumConfig{
		ReservationSymbolsPerSecond: 0,
		OnDemandSymbolsPerSecond:    100,
		OnDemandPricePerSymbol:      1,
	}

	// Try to use maximum possible symbols
	maxSymbols := uint64(math.MaxUint64)
	payment, err := acc.BlobPaymentInfo(ctx, maxSymbols, []uint8{0}, now)
	require.Error(t, err)
	assert.Nil(t, payment)
	assert.Contains(t, err.Error(), "current cumulativePayment balance insufficient")
}
