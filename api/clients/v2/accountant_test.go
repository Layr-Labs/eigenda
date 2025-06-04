package clients

import (
	"context"
	"encoding/hex"
	"math/big"
	"sync"
	"testing"
	"time"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

const numBins = uint32(3)

func TestNewAccountant(t *testing.T) {
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 100,
		StartTimestamp:   100,
		EndTimestamp:     200,
		QuorumSplits:     []byte{50, 50},
		QuorumNumbers:    []uint8{0, 1},
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
	assert.Equal(t, []PeriodRecord{{Index: 0, Usage: 0}, {Index: 1, Usage: 0}, {Index: 2, Usage: 0}}, accountant.periodRecords)
	assert.Equal(t, big.NewInt(0), accountant.cumulativePayment)
}

func TestAccountBlob_Reservation(t *testing.T) {
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 200,
		StartTimestamp:   100,
		EndTimestamp:     200,
		QuorumSplits:     []byte{50, 50},
		QuorumNumbers:    []uint8{0, 1},
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
	assert.Equal(t, isRotation([]uint64{500, 0, 0}, mapRecordUsage(accountant.periodRecords)), true)

	symbolLength = uint64(700)

	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(ctx, now, symbolLength, quorums)

	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)
	assert.Equal(t, isRotation([]uint64{1200, 0, 200}, mapRecordUsage(accountant.periodRecords)), true)

	// Second call should use on-demand payment
	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(ctx, now, 300, quorums)

	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(300), header.CumulativePayment)
}

func TestAccountBlob_OnDemand(t *testing.T) {
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 200,
		StartTimestamp:   100,
		EndTimestamp:     200,
		QuorumSplits:     []byte{50, 50},
		QuorumNumbers:    []uint8{0, 1},
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
	assert.Equal(t, isRotation([]uint64{0, 0, 0}, mapRecordUsage(accountant.periodRecords)), true)
	assert.Equal(t, expectedPayment, accountant.cumulativePayment)
}

func TestAccountBlob_InsufficientOnDemand(t *testing.T) {
	reservation := &core.ReservedPayment{}
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
	assert.Contains(t, err.Error(), "invalid payments")
}

func TestAccountBlobCallSeries(t *testing.T) {
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 200,
		StartTimestamp:   100,
		EndTimestamp:     200,
		QuorumSplits:     []byte{50, 50},
		QuorumNumbers:    []uint8{0, 1},
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
	assert.Contains(t, err.Error(), "invalid payments")
}

func TestAccountBlob_BinRotation(t *testing.T) {
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 1000,
		StartTimestamp:   100,
		EndTimestamp:     200,
		QuorumSplits:     []byte{50, 50},
		QuorumNumbers:    []uint8{0, 1},
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
	_, err = accountant.AccountBlob(ctx, now, 800, quorums)
	assert.NoError(t, err)
	assert.Equal(t, isRotation([]uint64{800, 0, 0}, mapRecordUsage(accountant.periodRecords)), true)

	// Second call
	now += int64(reservationWindow) * time.Second.Nanoseconds()
	_, err = accountant.AccountBlob(ctx, now, 300, quorums)
	assert.NoError(t, err)
	assert.Equal(t, isRotation([]uint64{800, 300, 0}, mapRecordUsage(accountant.periodRecords)), true)

	// Third call
	_, err = accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, isRotation([]uint64{800, 800, 0}, mapRecordUsage(accountant.periodRecords)), true)
}

func TestConcurrentBinRotationAndAccountBlob(t *testing.T) {
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 1000,
		StartTimestamp:   100,
		EndTimestamp:     200,
		QuorumSplits:     []byte{50, 50},
		QuorumNumbers:    []uint8{0, 1},
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

	// Check final state
	usages := mapRecordUsage(accountant.periodRecords)
	assert.Equal(t, uint64(1000), usages[0]+usages[1]+usages[2])
}

func TestAccountBlob_ReservationWithOneOverflow(t *testing.T) {
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 200,
		StartTimestamp:   100,
		EndTimestamp:     200,
		QuorumSplits:     []byte{50, 50},
		QuorumNumbers:    []uint8{0, 1},
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

	// Okay reservation
	header, err := accountant.AccountBlob(ctx, now, 800, quorums)
	assert.NoError(t, err)
	timestamp := (time.Duration(header.Timestamp) * time.Nanosecond).Seconds()
	assert.Equal(t, uint64(meterer.GetReservationPeriodByNanosecond(now, reservationWindow)), uint64(meterer.GetReservationPeriod(int64(timestamp), reservationWindow)))
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)
	assert.Equal(t, isRotation([]uint64{800, 0, 0}, mapRecordUsage(accountant.periodRecords)), true)

	// Second call: Allow one overflow
	header, err = accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)
	assert.Equal(t, isRotation([]uint64{1300, 0, 300}, mapRecordUsage(accountant.periodRecords)), true)

	// Third call: Should use on-demand payment
	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(ctx, now, 200, quorums)
	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(200), header.CumulativePayment)
	assert.Equal(t, isRotation([]uint64{1300, 0, 300}, mapRecordUsage(accountant.periodRecords)), true)
}

func TestAccountBlob_ReservationOverflowReset(t *testing.T) {
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 1000,
		StartTimestamp:   100,
		EndTimestamp:     200,
		QuorumSplits:     []byte{50, 50},
		QuorumNumbers:    []uint8{0, 1},
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
	_, err = accountant.AccountBlob(ctx, now, 1000, quorums)
	assert.NoError(t, err)
	assert.Equal(t, isRotation([]uint64{1000, 0, 0}, mapRecordUsage(accountant.periodRecords)), true)

	// no overflow
	now = time.Now().UnixNano()
	header, err := accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, isRotation([]uint64{1000, 0, 0}, mapRecordUsage(accountant.periodRecords)), true)
	assert.Equal(t, big.NewInt(500), header.CumulativePayment)

	// Wait for next reservation duration
	time.Sleep(time.Duration(reservationWindow) * time.Second)

	// Third call: Should use new bin and allow overflow again
	now = time.Now().UnixNano()
	_, err = accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, isRotation([]uint64{1000, 500, 0}, mapRecordUsage(accountant.periodRecords)), true)
}

// Test cases:
// Current bin under limit -> Simple increment
// Current bin over limit but overflow bin empty -> can use overflow for spillage
// Current bin over limit and overflow bin used -> reject, use on-demand
// New window - start fresh with current bin, request cannot fit into a bin -> reject, use on-demand
// New window - start fresh with current bin, within bin limit -> Simple increment
// New window - current bin over limit but can use overflow -> overwrite spillage in previous bin
// New window 2 - Exact bin limit usage -> Simple increment
// New window 2 - current bin at limit -> Simply reject and use on-demand; no spillage
func TestAccountBlob_ReservationOverflowWithWindow(t *testing.T) {
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 1000,
		StartTimestamp:   100,
		EndTimestamp:     200,
		QuorumSplits:     []byte{50, 50},
		QuorumNumbers:    []uint8{0, 1},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(3500),
	}
	reservationWindow := uint64(2) // Set to 2 seconds for testing
	pricePerSymbol := uint64(1)
	minNumSymbols := uint64(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, numBins)

	ctx := context.Background()
	quorums := []uint8{0, 1}

	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()
	windowSize := time.Duration(reservationWindow) * time.Second

	// Case 1: Current bin under limit -> Simple increment
	now := baseTime
	period := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(ctx, now, 1000, quorums)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1000), accountant.GetRelativePeriodRecord(period).Usage)
	assert.Equal(t, uint64(0), accountant.GetRelativePeriodRecord(period+reservationWindow).Usage)
	assert.Equal(t, uint64(0), accountant.GetRelativePeriodRecord(period+2*reservationWindow).Usage)

	// Case 2: Current bin over limit but overflow bin empty -> can use overflow period (current bin index + 2) for spillage
	now = baseTime + windowSize.Nanoseconds()/2
	period = meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(ctx, now, 1500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2500), accountant.GetRelativePeriodRecord(period).Usage)
	assert.Equal(t, uint64(0), accountant.GetRelativePeriodRecord(period+reservationWindow).Usage)
	assert.Equal(t, uint64(500), accountant.GetRelativePeriodRecord(period+2*reservationWindow).Usage)

	// Case 3: Current bin over limit and overflow bin used -> reject, use on-demand
	now = baseTime + windowSize.Nanoseconds()/2 + 100
	period = meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	header, err := accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(500), header.CumulativePayment)
	assert.Equal(t, uint64(2500), accountant.GetRelativePeriodRecord(period).Usage)
	assert.Equal(t, uint64(0), accountant.GetRelativePeriodRecord(period+reservationWindow).Usage)
	assert.Equal(t, uint64(500), accountant.GetRelativePeriodRecord(period+2*reservationWindow).Usage)

	// Case 4: New window - start fresh with the new current bin, request cannot fit into a bin -> reject, use on-demand
	now = baseTime + windowSize.Nanoseconds()
	period = meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	header, err = accountant.AccountBlob(ctx, now, 2500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(3000), header.CumulativePayment)
	assert.Equal(t, uint64(0), accountant.GetRelativePeriodRecord(period).Usage)
	assert.Equal(t, uint64(500), accountant.GetRelativePeriodRecord(period+reservationWindow).Usage)
	assert.Equal(t, uint64(0), accountant.GetRelativePeriodRecord(period+2*reservationWindow).Usage)

	// Case 5: New window - request within bin limit -> Simple increment
	now = baseTime + windowSize.Nanoseconds() + 100
	period = meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(ctx, now, 1000, quorums)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1000), accountant.GetRelativePeriodRecord(period).Usage)
	assert.Equal(t, uint64(500), accountant.GetRelativePeriodRecord(period+reservationWindow).Usage)
	assert.Equal(t, uint64(0), accountant.GetRelativePeriodRecord(period+2*reservationWindow).Usage)

	// Case 6: New window - current bin over limit but can use overflow -> overwrite spillage in previous bin
	// spillage = existing + new - bin limit = 1000 + 1500 - 2000 = 500
	now = baseTime + windowSize.Nanoseconds() + windowSize.Nanoseconds()/2
	period = meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(ctx, now, 1500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2500), accountant.GetRelativePeriodRecord(period).Usage)
	assert.Equal(t, uint64(500), accountant.GetRelativePeriodRecord(period+reservationWindow).Usage)
	assert.Equal(t, uint64(500), accountant.GetRelativePeriodRecord(period+2*reservationWindow).Usage)

	// Case 7: New window 2 - Exact bin limit usage -> Simple increment
	now = baseTime + 2*windowSize.Nanoseconds()
	period = meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(ctx, now, 1500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2000), accountant.GetRelativePeriodRecord(period).Usage)
	assert.Equal(t, uint64(500), accountant.GetRelativePeriodRecord(period+reservationWindow).Usage)
	assert.Equal(t, uint64(2500), accountant.GetRelativePeriodRecord(period+2*reservationWindow).Usage)

	// Case 8: New window 2 - current bin at limit -> Simply reject and use on-demand; no spillage
	// Past period record was cleaned up in the process
	now = baseTime + 2*windowSize.Nanoseconds() + 100
	period = meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	header, err = accountant.AccountBlob(ctx, now, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(3500), header.CumulativePayment)
	assert.Equal(t, uint64(2000), accountant.GetRelativePeriodRecord(period).Usage)
	assert.Equal(t, uint64(500), accountant.GetRelativePeriodRecord(period+reservationWindow).Usage)
	assert.Equal(t, uint64(0), accountant.GetRelativePeriodRecord(period+2*reservationWindow).Usage)

	// Case 9: New window 2 - current bin at limit, on-demand used up, cannot serve
	now = baseTime + 2*windowSize.Nanoseconds() + 100
	period = meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
	_, err = accountant.AccountBlob(ctx, now, 500, quorums)
	assert.Contains(t, err.Error(), "invalid payments")
	assert.Equal(t, uint64(2000), accountant.GetRelativePeriodRecord(period).Usage)
	assert.Equal(t, uint64(500), accountant.GetRelativePeriodRecord(period+reservationWindow).Usage)
	assert.Equal(t, uint64(0), accountant.GetRelativePeriodRecord(period+2*reservationWindow).Usage)
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

func mapRecordUsage(records []PeriodRecord) []uint64 {
	return []uint64{records[0].Usage, records[1].Usage, records[2].Usage}
}

func isRotation(arrA, arrB []uint64) bool {
	n := len(arrA)
	if n != len(arrB) {
		return false
	}

	doubleArrA := append(arrA, arrA...)
	// Check if arrB exists in doubleArrA as a subarray
	for i := 0; i < n; i++ {
		match := true
		for j := 0; j < n; j++ {
			if doubleArrA[i+j] != arrB[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func TestSetPaymentState(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey.D.Bytes()))

	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 0,
		StartTimestamp:   0,
		EndTimestamp:     0,
		QuorumNumbers:    []uint8{},
		QuorumSplits:     []byte{},
	}

	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(0),
	}

	accountant := NewAccountant(accountId, reservation, onDemand, 0, 0, 0, numBins)

	t.Run("nil payment state", func(t *testing.T) {
		err := accountant.SetPaymentState(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment state cannot be nil")
	})

	t.Run("nil payment global params", func(t *testing.T) {
		state := &disperser_rpc.GetPaymentStateReply{}
		err := accountant.SetPaymentState(state)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment global params cannot be nil")
	})

	t.Run("successful set payment state with all fields", func(t *testing.T) {
		onchainCumulativePayment := big.NewInt(1000).Bytes()
		cumulativePayment := big.NewInt(500).Bytes()

		state := &disperser_rpc.GetPaymentStateReply{
			PaymentGlobalParams: &disperser_rpc.PaymentGlobalParams{
				MinNumSymbols:     100,
				PricePerSymbol:    50,
				ReservationWindow: 60,
			},
			OnchainCumulativePayment: onchainCumulativePayment,
			CumulativePayment:        cumulativePayment,
			Reservation: &disperser_rpc.Reservation{
				SymbolsPerSecond: 300,
				StartTimestamp:   100,
				EndTimestamp:     200,
				QuorumNumbers:    []uint32{0},
				QuorumSplits:     []uint32{100},
			},
			PeriodRecords: []*disperser_rpc.PeriodRecord{
				{
					Index: 1,
					Usage: 150,
				},
				{
					Index: 0,
					Usage: 0,
				},
				{
					Index: 0,
					Usage: 0,
				},
			},
		}

		err := accountant.SetPaymentState(state)
		assert.NoError(t, err)

		// Verify the state was set correctly
		assert.Equal(t, uint64(50), accountant.pricePerSymbol)
		assert.Equal(t, uint64(60), accountant.reservationWindow)
		assert.Equal(t, uint64(100), accountant.minNumSymbols)
		assert.Equal(t, big.NewInt(1000), accountant.onDemand.CumulativePayment)
		assert.Equal(t, big.NewInt(500), accountant.cumulativePayment)
		assert.Equal(t, uint64(300), accountant.reservation.SymbolsPerSecond)
		assert.Equal(t, uint64(100), accountant.reservation.StartTimestamp)
		assert.Equal(t, uint64(200), accountant.reservation.EndTimestamp)
		assert.Equal(t, []uint8{0}, accountant.reservation.QuorumNumbers)

		// Check period records
		assert.Equal(t, uint32(1), accountant.periodRecords[0].Index)
		assert.Equal(t, uint64(150), accountant.periodRecords[0].Usage)
		assert.Equal(t, uint32(0), accountant.periodRecords[1].Index)
		assert.Equal(t, uint64(0), accountant.periodRecords[1].Usage)
		assert.Equal(t, uint32(0), accountant.periodRecords[2].Index)
		assert.Equal(t, uint64(0), accountant.periodRecords[2].Usage)
	})

	t.Run("successful set payment state with minimal fields", func(t *testing.T) {
		state := &disperser_rpc.GetPaymentStateReply{
			PaymentGlobalParams: &disperser_rpc.PaymentGlobalParams{
				MinNumSymbols:     50,
				PricePerSymbol:    25,
				ReservationWindow: 30,
			},
			// No OnchainCumulativePayment
			// No CumulativePayment
			// No Reservation
			// No PeriodRecords
		}

		err := accountant.SetPaymentState(state)
		assert.NoError(t, err)

		// Verify default values are set
		assert.Equal(t, uint64(25), accountant.pricePerSymbol)
		assert.Equal(t, uint64(30), accountant.reservationWindow)
		assert.Equal(t, uint64(50), accountant.minNumSymbols)
		assert.Equal(t, big.NewInt(0), accountant.onDemand.CumulativePayment)
		assert.Equal(t, big.NewInt(0), accountant.cumulativePayment)
		assert.Equal(t, uint64(0), accountant.reservation.SymbolsPerSecond)
		assert.Equal(t, uint64(0), accountant.reservation.StartTimestamp)
		assert.Equal(t, uint64(0), accountant.reservation.EndTimestamp)
		assert.Equal(t, []uint8{}, accountant.reservation.QuorumNumbers)

		// Verify period records are initialized but empty
		for i := range accountant.periodRecords {
			assert.Equal(t, uint32(i), accountant.periodRecords[i].Index)
			assert.Equal(t, uint64(0), accountant.periodRecords[i].Usage)
		}
	})
}
