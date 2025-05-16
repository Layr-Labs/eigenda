package clients

import (
	"context"
	"encoding/hex"
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
)

const numBins = uint32(3)

func TestNewAccountant(t *testing.T) {
	reservation := map[uint8]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 100,
			StartTimestamp:   100,
			EndTimestamp:     200,
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}
	reservationWindow := uint64(6)
	pricePerSymbol := uint64(1)
	minNumSymbols := uint64(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, numBins)

	assert.NotNil(t, accountant)
	assert.Equal(t, reservation, accountant.reservation)
	assert.Equal(t, onDemand, accountant.onDemand)
	assert.Equal(t, reservationWindow, accountant.reservationWindow)
	assert.Equal(t, pricePerSymbol, accountant.pricePerSymbol)
	assert.Equal(t, minNumSymbols, accountant.minNumSymbols)

	// Check initialization of periodRecords
	for quorumNumber, records := range accountant.periodRecords {
		assert.Equal(t, int(numBins), len(records), "Should have numBins records for each quorum")
		for i, record := range records {
			assert.Equal(t, uint32(i), record.Index)
			assert.Equal(t, uint64(0), record.Usage)
			assert.Equal(t, quorumNumber, record.QuorumNumber)
		}
	}

	assert.Equal(t, big.NewInt(0), accountant.cumulativePayment)
}

func TestAccountBlob_Reservation(t *testing.T) {
	reservation := map[uint8]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 200,
			StartTimestamp:   100,
			EndTimestamp:     200,
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}
	reservationWindow := uint64(5)
	pricePerSymbol := uint64(1)
	minNumSymbols := uint64(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, numBins)

	ctx := context.Background()
	symbolLength := uint64(500)
	quorums := []uint8{0, 1}
	now := time.Now().UnixNano()

	header, err := accountant.AccountBlob(ctx, now, symbolLength, quorums)

	assert.NoError(t, err)
	assert.Equal(t, meterer.GetReservationPeriod(time.Now().Unix(), reservationWindow), meterer.GetReservationPeriodByNanosecond(header.Timestamp, reservationWindow))
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Check that usage for quorum 0 was updated properly
	quorum0Records := accountant.periodRecords[0]
	currentPeriodIndex := uint32(meterer.GetReservationPeriodByNanosecond(now, reservationWindow) % uint64(numBins))
	assert.Equal(t, symbolLength, quorum0Records[currentPeriodIndex].Usage)

	symbolLength = uint64(700)

	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(ctx, now, symbolLength, quorums)

	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// With overflow, usage should be at the limit in the current period
	binLimit := reservation[0].SymbolsPerSecond * uint64(reservationWindow)
	expectedUsage := binLimit
	assert.Equal(t, expectedUsage, quorum0Records[currentPeriodIndex].Usage)

	// The overflow should be in the next bin
	overflowIndex := uint32((meterer.GetReservationPeriodByNanosecond(now, reservationWindow) + 2) % uint64(numBins))
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
	reservation := map[uint8]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 200,
			StartTimestamp:   100,
			EndTimestamp:     200,
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1500),
	}
	reservationWindow := uint64(5)
	pricePerSymbol := uint64(1)
	minNumSymbols := uint64(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, numBins)

	ctx := context.Background()
	numSymbols := uint64(1500)
	quorums := []uint8{0, 1}
	now := time.Now().UnixNano()
	header, err := accountant.AccountBlob(ctx, now, numSymbols, quorums)
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
	reservation := map[uint8]*core.ReservedPayment{} // Empty reservation map
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}
	reservationWindow := uint64(60)
	pricePerSymbol := uint64(100)
	minNumSymbols := uint64(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, numBins)

	ctx := context.Background()
	numSymbols := uint64(2000)
	quorums := []uint8{0, 1}
	now := time.Now().UnixNano()
	_, err = accountant.AccountBlob(ctx, now, numSymbols, quorums)
	assert.Contains(t, err.Error(), "no bandwidth reservation found for account")
}

func TestAccountBlobCallSeries(t *testing.T) {
	reservation := map[uint8]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 200,
			StartTimestamp:   100,
			EndTimestamp:     200,
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(5)
	pricePerSymbol := uint64(1)
	minNumSymbols := uint64(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, numBins)

	ctx := context.Background()
	quorums := []uint8{0, 1}

	now := time.Now().UnixNano()
	// First call: Use reservation
	header, err := accountant.AccountBlob(ctx, now, 800, quorums)
	assert.NoError(t, err)
	timestamp := (time.Duration(header.Timestamp) * time.Nanosecond).Seconds()
	assert.Equal(t, uint64(meterer.GetReservationPeriodByNanosecond(now, reservationWindow)), uint64(meterer.GetReservationPeriod(int64(timestamp), reservationWindow)))
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Second call: Use remaining reservation + overflow
	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(ctx, now, 300, quorums)
	assert.NoError(t, err)
	timestamp = (time.Duration(header.Timestamp) * time.Nanosecond).Seconds()
	assert.Equal(t, uint64(meterer.GetReservationPeriodByNanosecond(now, reservationWindow)), uint64(meterer.GetReservationPeriod(int64(timestamp), reservationWindow)))
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Third call: Use on-demand
	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(500), header.CumulativePayment)

	// Fourth call: Insufficient on-demand
	now = time.Now().UnixNano()
	_, err = accountant.AccountBlob(ctx, now, 600, quorums)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no bandwidth reservation found for account")
}

func TestAccountBlob_BinRotation(t *testing.T) {
	reservation := map[uint8]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   100,
			EndTimestamp:     200,
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(1) // Set to 1 second for testing
	pricePerSymbol := uint64(1)
	minNumSymbols := uint64(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, numBins)

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// First call
	now := time.Now().UnixNano()
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(ctx, now, 800, quorums)
	assert.NoError(t, err)

	// Check bin 0 has usage 800
	record := accountant.GetRelativePeriodRecord(currentPeriod, 0)
	assert.Equal(t, uint64(800), record.Usage)

	// Second call (next period)
	now += int64(reservationWindow) * time.Second.Nanoseconds()
	nextPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	assert.Equal(t, currentPeriod+1, nextPeriod, "Should be next period")

	_, err = accountant.AccountBlob(ctx, now, 300, quorums)
	assert.NoError(t, err)

	// Check bin 1 has usage 300
	record = accountant.GetRelativePeriodRecord(nextPeriod, 0)
	assert.Equal(t, uint64(300), record.Usage)

	// Third call in same period
	_, err = accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)

	// Check bin 1 now has usage 800
	record = accountant.GetRelativePeriodRecord(nextPeriod, 0)
	assert.Equal(t, uint64(800), record.Usage)
}

func TestConcurrentBinRotationAndAccountBlob(t *testing.T) {
	reservation := map[uint8]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   100,
			EndTimestamp:     200,
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(1) // Set to 1 second for testing
	pricePerSymbol := uint64(1)
	minNumSymbols := uint64(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, numBins)

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// Start concurrent AccountBlob calls
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			now := time.Now().UnixNano()
			_, err := accountant.AccountBlob(ctx, now, 100, quorums)
			assert.NoError(t, err)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check final state - total usage should be 1000 across all bins for quorum 0
	totalUsage := uint64(0)
	for _, record := range accountant.periodRecords[0] {
		totalUsage += record.Usage
	}
	assert.Equal(t, uint64(1000), totalUsage)
}

func TestAccountBlob_ReservationWithOneOverflow(t *testing.T) {
	reservation := map[uint8]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 200,
			StartTimestamp:   100,
			EndTimestamp:     200,
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(5)
	pricePerSymbol := uint64(1)
	minNumSymbols := uint64(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, numBins)

	ctx := context.Background()
	quorums := []uint8{0, 1}
	now := time.Now().UnixNano()
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)

	// Okay reservation
	header, err := accountant.AccountBlob(ctx, now, 800, quorums)
	assert.NoError(t, err)
	timestamp := (time.Duration(header.Timestamp) * time.Nanosecond).Seconds()
	assert.Equal(t, uint64(meterer.GetReservationPeriodByNanosecond(now, reservationWindow)), uint64(meterer.GetReservationPeriod(int64(timestamp), reservationWindow)))
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Check current period usage
	record := accountant.GetRelativePeriodRecord(currentPeriod, 0)
	assert.Equal(t, uint64(800), record.Usage)

	// Second call: Allow one overflow
	header, err = accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Check current period is at limit
	binLimit := reservation[0].SymbolsPerSecond * uint64(reservationWindow) // 1000
	record = accountant.GetRelativePeriodRecord(currentPeriod, 0)
	assert.Equal(t, binLimit, record.Usage)

	// Check overflow period has the overflow
	overflowRecord := accountant.GetRelativePeriodRecord(currentPeriod+2, 0)
	assert.Equal(t, uint64(300), overflowRecord.Usage) // 800 + 500 - 1000 = 300

	// Third call: Should use on-demand payment
	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(ctx, now, 200, quorums)
	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(200), header.CumulativePayment)
}

func TestAccountBlob_ReservationOverflowReset(t *testing.T) {
	reservation := map[uint8]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   100,
			EndTimestamp:     200,
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(1) // Set to 1 second for testing
	pricePerSymbol := uint64(1)
	minNumSymbols := uint64(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, numBins)

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// full reservation
	now := time.Now().UnixNano()
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(ctx, now, 1000, quorums)
	assert.NoError(t, err)

	// Check current period is at limit
	record := accountant.GetRelativePeriodRecord(currentPeriod, 0)
	assert.Equal(t, uint64(1000), record.Usage)

	// no overflow
	now = time.Now().UnixNano()
	header, err := accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)
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
	record = accountant.GetRelativePeriodRecord(nextPeriod, 0)
	assert.Equal(t, uint64(500), record.Usage)
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
			errorMessage:   "no quorum numbers provided",
		},
		{
			name:           "invalid quorum number",
			quorumNumbers:  []uint8{0, 2},
			allowedNumbers: []uint8{0, 1},
			expectError:    true,
			errorMessage:   "provided quorum number 2 not allowed",
		},
		{
			name:           "empty allowed numbers",
			quorumNumbers:  []uint8{0},
			allowedNumbers: []uint8{},
			expectError:    true,
			errorMessage:   "provided quorum number 0 not allowed",
		},
		{
			name:           "multiple invalid quorums",
			quorumNumbers:  []uint8{2, 3, 4},
			allowedNumbers: []uint8{0, 1},
			expectError:    true,
			errorMessage:   "provided quorum number 2 not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := QuorumCheck(tt.quorumNumbers, tt.allowedNumbers)
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
	// Create accountant with initial state
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey.D.Bytes()))

	// Create with empty state
	accountant := NewAccountant(accountId,
		map[uint8]*core.ReservedPayment{},
		&core.OnDemandPayment{CumulativePayment: big.NewInt(0)},
		10, 1, 100, numBins)

	// Create payment state reply with sample data
	paymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentGlobalParams: &v2.PaymentGlobalParams{
			GlobalSymbolsPerSecond: 2000,
			MinNumSymbols:          200,
			PricePerSymbol:         2,
			ReservationWindow:      20,
		},
		CumulativePayment:        big.NewInt(500).Bytes(),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: []*v2.QuorumReservation{
			{
				QuorumNumber:     0,
				SymbolsPerSecond: 300,
				StartTimestamp:   150,
				EndTimestamp:     250,
			},
			{
				QuorumNumber:     1,
				SymbolsPerSecond: 400,
				StartTimestamp:   160,
				EndTimestamp:     260,
			},
		},
		PeriodRecords: []*v2.QuorumPeriodRecord{
			{
				QuorumNumber: 0,
				Index:        123,
				Usage:        600,
			},
			{
				QuorumNumber: 1,
				Index:        123,
				Usage:        700,
			},
		},
	}

	// Set payment state
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	// Verify global params are updated
	assert.Equal(t, uint64(200), accountant.minNumSymbols)
	assert.Equal(t, uint64(2), accountant.pricePerSymbol)
	assert.Equal(t, uint64(20), accountant.reservationWindow)

	// Verify on-demand payment is updated
	assert.Equal(t, big.NewInt(1000), accountant.onDemand.CumulativePayment)

	// Verify cumulative payment is updated
	assert.Equal(t, big.NewInt(500), accountant.cumulativePayment)

	// Verify reservations are updated
	assert.Equal(t, 2, len(accountant.reservation))
	assert.Equal(t, uint64(300), accountant.reservation[0].SymbolsPerSecond)
	assert.Equal(t, uint64(150), accountant.reservation[0].StartTimestamp)
	assert.Equal(t, uint64(250), accountant.reservation[0].EndTimestamp)
	assert.Equal(t, uint64(400), accountant.reservation[1].SymbolsPerSecond)
	assert.Equal(t, uint64(160), accountant.reservation[1].StartTimestamp)
	assert.Equal(t, uint64(260), accountant.reservation[1].EndTimestamp)

	// Test with nil values
	nilPaymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentGlobalParams: &v2.PaymentGlobalParams{
			GlobalSymbolsPerSecond: 3000,
			MinNumSymbols:          300,
			PricePerSymbol:         3,
			ReservationWindow:      30,
		},
		// No CumulativePayment
		// No OnchainCumulativePayment
		// No Reservations
	}

	err = accountant.SetPaymentState(nilPaymentState)
	assert.NoError(t, err)

	// Verify global params are updated
	assert.Equal(t, uint64(300), accountant.minNumSymbols)
	assert.Equal(t, uint64(3), accountant.pricePerSymbol)
	assert.Equal(t, uint64(30), accountant.reservationWindow)

	// Verify defaults when fields are nil
	assert.Equal(t, big.NewInt(0), accountant.onDemand.CumulativePayment)
	assert.Equal(t, big.NewInt(0), accountant.cumulativePayment)
	assert.Equal(t, 0, len(accountant.reservation))

	// Test with nil payment state
	err = accountant.SetPaymentState(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payment state cannot be nil")

	// Test with nil payment global params
	nilGlobalParams := &v2.GetQuorumSpecificPaymentStateReply{
		// No PaymentGlobalParams
	}
	err = accountant.SetPaymentState(nilGlobalParams)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payment global params cannot be nil")
}
