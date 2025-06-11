package clients

import (
	"encoding/hex"
	"fmt"
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

type newAccountantTest struct {
	name              string
	accountId         gethcommon.Address
	reservation       *core.ReservedPayment
	onDemand          *core.OnDemandPayment
	reservationWindow uint64
	pricePerSymbol    uint64
	minNumSymbols     uint64
	expectedRecords   []PeriodRecord
}

func TestNewAccountant(t *testing.T) {
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))

	tests := []newAccountantTest{
		{
			name:      "basic accountant initialization",
			accountId: accountId,
			reservation: &core.ReservedPayment{
				SymbolsPerSecond: 100,
				StartTimestamp:   100,
				EndTimestamp:     200,
				QuorumSplits:     []byte{50, 50},
				QuorumNumbers:    []uint8{0, 1},
			},
			onDemand: &core.OnDemandPayment{
				CumulativePayment: big.NewInt(500),
			},
			reservationWindow: 6,
			pricePerSymbol:    1,
			minNumSymbols:     100,
			expectedRecords: []PeriodRecord{
				{Index: 0, Usage: 0},
				{Index: 1, Usage: 0},
				{Index: 2, Usage: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountant := NewAccountant(tt.accountId, tt.reservation, tt.onDemand, tt.reservationWindow, tt.pricePerSymbol, tt.minNumSymbols, numBins)

			assert.NotNil(t, accountant)
			assert.Equal(t, tt.reservation, accountant.reservation)
			assert.Equal(t, tt.onDemand, accountant.onDemand)
			assert.Equal(t, tt.reservationWindow, accountant.reservationWindow)
			assert.Equal(t, tt.pricePerSymbol, accountant.pricePerSymbol)
			assert.Equal(t, tt.minNumSymbols, accountant.minNumSymbols)
			assert.Equal(t, tt.expectedRecords, accountant.periodRecords)
			assert.Equal(t, big.NewInt(0), accountant.cumulativePayment)
		})
	}
}

type accountBlobReservationTest struct {
	name           string
	symbolLength   uint64
	expectedHeader *core.PaymentMetadata
	expectedState  periodRecordState
}

func TestAccountBlob_Reservation(t *testing.T) {
	baseTime := time.Now()
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 200,
		StartTimestamp:   uint64(baseTime.Add(-10 * time.Second).Unix()),
		EndTimestamp:     uint64(baseTime.Add(10 * time.Second).Unix()),
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

	quorums := []uint8{0, 1}

	tests := []accountBlobReservationTest{
		{
			name:         "First call - use reservation",
			symbolLength: 500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 500,
				NextPeriodUsage:    0,
				OverflowUsage:      0,
			},
		},
		{
			name:         "Second call - use reservation with overflow",
			symbolLength: 700,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 1200,
				NextPeriodUsage:    0,
				OverflowUsage:      200,
			},
		},
		{
			name:         "Third call - use on-demand payment",
			symbolLength: 300,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(300),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 1200,
				NextPeriodUsage:    0,
				OverflowUsage:      200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := baseTime.UnixNano()
			header, err := accountant.AccountBlob(now, tt.symbolLength, quorums)

			assert.NoError(t, err)
			assert.NotEqual(t, uint64(0), header.Timestamp)
			assert.Equal(t, tt.expectedHeader.AccountID, header.AccountID)
			assert.Equal(t, tt.expectedHeader.CumulativePayment, header.CumulativePayment)

			period := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
			assert.Equal(t, tt.expectedState.CurrentPeriodUsage, getRelativePeriodRecord(period, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.NextPeriodUsage, getRelativePeriodRecord(period+reservationWindow, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.OverflowUsage, getRelativePeriodRecord(meterer.GetOverflowPeriod(period, reservationWindow), reservationWindow, accountant.periodRecords).Usage)
		})
	}
}

type accountBlobOnDemandTest struct {
	name           string
	symbolLength   uint64
	expectedHeader *core.PaymentMetadata
	expectedState  periodRecordState
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

	quorums := []uint8{0, 1}
	baseTime := time.Now().UnixNano()

	tests := []accountBlobOnDemandTest{
		{
			name:         "Use on-demand payment",
			symbolLength: 1500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(1500),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 0,
				NextPeriodUsage:    0,
				OverflowUsage:      0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := baseTime
			header, err := accountant.AccountBlob(now, tt.symbolLength, quorums)

			assert.NoError(t, err)
			assert.NotEqual(t, uint64(0), header.Timestamp)
			assert.Equal(t, tt.expectedHeader.AccountID, header.AccountID)
			assert.Equal(t, tt.expectedHeader.CumulativePayment, header.CumulativePayment)

			period := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
			assert.Equal(t, tt.expectedState.CurrentPeriodUsage, getRelativePeriodRecord(period, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.NextPeriodUsage, getRelativePeriodRecord(period+reservationWindow, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.OverflowUsage, getRelativePeriodRecord(meterer.GetOverflowPeriod(period, reservationWindow), reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedHeader.CumulativePayment, accountant.cumulativePayment)
		})
	}
}

type accountBlobInsufficientOnDemandTest struct {
	name         string
	symbolLength uint64
	expectError  bool
	errorMessage string
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
	baseTime := time.Now()
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 200,
		StartTimestamp:   uint64(baseTime.Add(-10 * time.Second).Unix()),
		EndTimestamp:     uint64(baseTime.Add(10 * time.Second).Unix()),
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
	baseTime := time.Now()
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 1000,
		StartTimestamp:   uint64(baseTime.Add(-10 * time.Second).Unix()),
		EndTimestamp:     uint64(baseTime.Add(10 * time.Second).Unix()),
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

	quorums := []uint8{0, 1}

	tests := []accountBlobBinRotationTest{
		{
			name:         "First call - Initial usage",
			timeOffset:   0,
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
			name:         "Second call - After window",
			timeOffset:   time.Duration(reservationWindow) * time.Second,
			symbolLength: 300,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 300,
				NextPeriodUsage:    0,
				OverflowUsage:      800, // previous bin not reset yet
			},
		},
		{
			name:         "Third call - Same window",
			timeOffset:   time.Duration(reservationWindow) * time.Second,
			symbolLength: 500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 800,
				NextPeriodUsage:    0,
				OverflowUsage:      800,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := baseTime.UnixNano() + tt.timeOffset.Nanoseconds()
			header, err := accountant.AccountBlob(now, tt.symbolLength, quorums)

			assert.NoError(t, err)
			assert.NotEqual(t, uint64(0), header.Timestamp)
			assert.Equal(t, tt.expectedHeader.AccountID, header.AccountID)
			assert.Equal(t, tt.expectedHeader.CumulativePayment, header.CumulativePayment)

			period := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
			assert.Equal(t, tt.expectedState.CurrentPeriodUsage, getRelativePeriodRecord(period, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.NextPeriodUsage, getRelativePeriodRecord(period+reservationWindow, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.OverflowUsage, getRelativePeriodRecord(meterer.GetOverflowPeriod(period, reservationWindow), reservationWindow, accountant.periodRecords).Usage)
		})
	}
}

type concurrentBinRotationTest struct {
	name          string
	numGoroutines int
	symbolLength  uint64
	expectedTotal uint64
}

func TestConcurrentBinRotationAndAccountBlob(t *testing.T) {
	baseTime := time.Now()
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 1000,
		StartTimestamp:   uint64(baseTime.Add(-10 * time.Second).Unix()),
		EndTimestamp:     uint64(baseTime.Add(10 * time.Second).Unix()),
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

	quorums := []uint8{0, 1}

	tests := []concurrentBinRotationTest{
		{
			name:          "Concurrent calls with 10 goroutines",
			numGoroutines: 10,
			symbolLength:  100,
			expectedTotal: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			for i := 0; i < tt.numGoroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					now := time.Now().UnixNano()
					_, err := accountant.AccountBlob(now, tt.symbolLength, quorums)
					assert.NoError(t, err)
				}()
			}

			// Wait for all goroutines to finish
			wg.Wait()

			// Check final state
			now := time.Now().UnixNano()
			period := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
			totalUsage := getRelativePeriodRecord(period, reservationWindow, accountant.periodRecords).Usage +
				getRelativePeriodRecord(period+reservationWindow, reservationWindow, accountant.periodRecords).Usage +
				getRelativePeriodRecord(meterer.GetOverflowPeriod(period, reservationWindow), reservationWindow, accountant.periodRecords).Usage
			assert.Equal(t, tt.expectedTotal, totalUsage)
		})
	}
}

type accountBlobReservationWithOneOverflowTest struct {
	name           string
	symbolLength   uint64
	expectedHeader *core.PaymentMetadata
	expectedState  periodRecordState
}

func TestAccountBlob_ReservationWithOneOverflow(t *testing.T) {
	baseTime := time.Now()
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 200,
		StartTimestamp:   uint64(baseTime.Add(-10 * time.Second).Unix()),
		EndTimestamp:     uint64(baseTime.Add(10 * time.Second).Unix()),
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

	quorums := []uint8{0, 1}

	tests := []accountBlobReservationWithOneOverflowTest{
		{
			name:         "First call - Okay reservation",
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
			name:         "Second call - Allow one overflow",
			symbolLength: 500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 1300,
				NextPeriodUsage:    0,
				OverflowUsage:      300,
			},
		},
		{
			name:         "Third call - Use on-demand payment",
			symbolLength: 200,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(200),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 1300,
				NextPeriodUsage:    0,
				OverflowUsage:      300,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := baseTime.UnixNano()
			header, err := accountant.AccountBlob(now, tt.symbolLength, quorums)

			assert.NoError(t, err)
			assert.NotEqual(t, uint64(0), header.Timestamp)
			assert.Equal(t, tt.expectedHeader.AccountID, header.AccountID)
			assert.Equal(t, tt.expectedHeader.CumulativePayment, header.CumulativePayment)

			period := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
			assert.Equal(t, tt.expectedState.CurrentPeriodUsage, getRelativePeriodRecord(period, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.NextPeriodUsage, getRelativePeriodRecord(period+reservationWindow, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.OverflowUsage, getRelativePeriodRecord(meterer.GetOverflowPeriod(period, reservationWindow), reservationWindow, accountant.periodRecords).Usage)
		})
	}
}

type accountBlobReservationOverflowResetTest struct {
	name           string
	timeOffset     time.Duration
	symbolLength   uint64
	expectedHeader *core.PaymentMetadata
	expectedState  periodRecordState
}

func TestAccountBlob_ReservationOverflowReset(t *testing.T) {
	baseTime := time.Now()
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 1000,
		StartTimestamp:   uint64(baseTime.Add(-10 * time.Second).Unix()),
		EndTimestamp:     uint64(baseTime.Add(10 * time.Second).Unix()),
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

	quorums := []uint8{0, 1}

	tests := []accountBlobReservationOverflowResetTest{
		{
			name:         "First call - Full reservation",
			timeOffset:   0,
			symbolLength: 1000,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 1000,
				NextPeriodUsage:    0,
				OverflowUsage:      0,
			},
		},
		{
			name:         "Second call - No overflow",
			timeOffset:   0,
			symbolLength: 500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(500),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 1000,
				NextPeriodUsage:    0,
				OverflowUsage:      0,
			},
		},
		{
			name:         "Third call - New window",
			timeOffset:   time.Duration(reservationWindow) * time.Second,
			symbolLength: 500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 500,
				NextPeriodUsage:    0,
				OverflowUsage:      1000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := baseTime.UnixNano() + tt.timeOffset.Nanoseconds()
			header, err := accountant.AccountBlob(now, tt.symbolLength, quorums)

			assert.NoError(t, err)
			assert.NotEqual(t, uint64(0), header.Timestamp)
			assert.Equal(t, tt.expectedHeader.AccountID, header.AccountID)
			assert.Equal(t, tt.expectedHeader.CumulativePayment, header.CumulativePayment)

			period := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
			assert.Equal(t, tt.expectedState.CurrentPeriodUsage, getRelativePeriodRecord(period, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.NextPeriodUsage, getRelativePeriodRecord(period+reservationWindow, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.OverflowUsage, getRelativePeriodRecord(meterer.GetOverflowPeriod(period, reservationWindow), reservationWindow, accountant.periodRecords).Usage)
		})
	}
}

type periodRecordState struct {
	CurrentPeriodUsage uint64
	NextPeriodUsage    uint64
	OverflowUsage      uint64
}

type testScenario struct {
	name           string
	timeOffset     time.Duration // Offset from base time
	symbolLength   uint64
	expectedHeader *core.PaymentMetadata
	expectedState  periodRecordState
}

func TestAccountBlob_ReservationOverflowWithWindow(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 1000,
		StartTimestamp:   uint64(baseTime.Add(-10 * time.Second).Unix()),
		EndTimestamp:     uint64(baseTime.Add(10 * time.Second).Unix()),
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

	quorums := []uint8{0, 1}

	windowSize := time.Duration(reservationWindow) * time.Second

	scenarios := []testScenario{
		{
			name:         "Current bin under limit -> Simple increment",
			timeOffset:   0,
			symbolLength: 1000,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 1000,
				NextPeriodUsage:    0,
				OverflowUsage:      0,
			},
		},
		{
			name:         "Current bin over limit but overflow bin empty -> can use overflow period",
			timeOffset:   windowSize / 2,
			symbolLength: 1500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 2500,
				NextPeriodUsage:    0,
				OverflowUsage:      500,
			},
		},
		{
			name:         "Current bin over limit and overflow bin used -> reject, use on-demand",
			timeOffset:   windowSize/2 + 100*time.Nanosecond,
			symbolLength: 500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(500),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 2500,
				NextPeriodUsage:    0,
				OverflowUsage:      500,
			},
		},
		{
			name:         "New window - request cannot fit into a bin -> reject, use on-demand",
			timeOffset:   windowSize,
			symbolLength: 2500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(3000),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 0,
				NextPeriodUsage:    500,
				OverflowUsage:      0,
			},
		},
		{
			name:         "New window - request within bin limit -> Simple increment",
			timeOffset:   windowSize + 100*time.Nanosecond,
			symbolLength: 1000,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 1000,
				NextPeriodUsage:    500,
				OverflowUsage:      0,
			},
		},
		{
			name:         "New window - current bin over limit but can use overflow",
			timeOffset:   windowSize + windowSize/2,
			symbolLength: 1500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 2500,
				NextPeriodUsage:    500,
				OverflowUsage:      500,
			},
		},
		{
			name:         "New window 2 - Exact bin limit usage",
			timeOffset:   2 * windowSize,
			symbolLength: 1500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(0),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 2000,
				NextPeriodUsage:    500,
				OverflowUsage:      2500,
			},
		},
		{
			name:         "New window 2 - current bin at limit -> use on-demand",
			timeOffset:   2*windowSize + 100*time.Nanosecond,
			symbolLength: 500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(3500),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 2000,
				NextPeriodUsage:    500,
				OverflowUsage:      0,
			},
		},
		{
			name:         "New window 2 - current bin at limit, on-demand used up, cannot serve",
			timeOffset:   2*windowSize + 100*time.Nanosecond,
			symbolLength: 500,
			expectedHeader: &core.PaymentMetadata{
				AccountID:         accountId,
				CumulativePayment: big.NewInt(3500),
			},
			expectedState: periodRecordState{
				CurrentPeriodUsage: 2000,
				NextPeriodUsage:    500,
				OverflowUsage:      0,
			},
		},
	}

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			now := baseTime.UnixNano() + tt.timeOffset.Nanoseconds()
			header, err := accountant.AccountBlob(now, tt.symbolLength, quorums)

			if tt.name == "New window 2 - current bin at limit, on-demand used up, cannot serve" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "insufficient ondemand payment")
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uint64(0), header.Timestamp)
				assert.Equal(t, tt.expectedHeader.AccountID, header.AccountID)
				assert.Equal(t, tt.expectedHeader.CumulativePayment, header.CumulativePayment)
			}

			period := meterer.GetReservationPeriodByNanosecond(now, reservationWindow)
			assert.Equal(t, tt.expectedState.CurrentPeriodUsage, getRelativePeriodRecord(period, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.NextPeriodUsage, getRelativePeriodRecord(period+reservationWindow, reservationWindow, accountant.periodRecords).Usage)
			assert.Equal(t, tt.expectedState.OverflowUsage, getRelativePeriodRecord(meterer.GetOverflowPeriod(period, reservationWindow), reservationWindow, accountant.periodRecords).Usage)
		})
	}
}

type quorumCheckTest struct {
	name           string
	quorumNumbers  []uint8
	allowedNumbers []uint8
	expectError    bool
	errorMessage   string
}

func TestQuorumCheck(t *testing.T) {
	tests := []quorumCheckTest{
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

type setPaymentStateTest struct {
	name          string
	state         *disperser_rpc.GetPaymentStateReply
	expectError   bool
	errorMessage  string
	expectedState *Accountant
}

func TestSetPaymentState(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey.D.Bytes()))

	emptyReservation := &core.ReservedPayment{
		SymbolsPerSecond: 0,
		StartTimestamp:   0,
		EndTimestamp:     0,
		QuorumNumbers:    []uint8{},
		QuorumSplits:     []byte{},
	}

	emptyOnDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(0),
	}

	tests := []setPaymentStateTest{
		{
			name:         "nil payment state",
			state:        nil,
			expectError:  true,
			errorMessage: "payment state cannot be nil",
		},
		{
			name:         "nil payment global params",
			state:        &disperser_rpc.GetPaymentStateReply{},
			expectError:  true,
			errorMessage: "payment global params cannot be nil",
		},
		{
			name: "successful set payment state with all fields",
			state: &disperser_rpc.GetPaymentStateReply{
				PaymentGlobalParams: &disperser_rpc.PaymentGlobalParams{
					MinNumSymbols:     100,
					PricePerSymbol:    50,
					ReservationWindow: 60,
				},
				OnchainCumulativePayment: big.NewInt(1000).Bytes(),
				CumulativePayment:        big.NewInt(500).Bytes(),
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
			},
			expectError: false,
			expectedState: &Accountant{
				accountID: accountId,
				reservation: &core.ReservedPayment{
					SymbolsPerSecond: 300,
					StartTimestamp:   100,
					EndTimestamp:     200,
					QuorumNumbers:    []uint8{0},
					QuorumSplits:     []byte{100},
				},
				onDemand: &core.OnDemandPayment{
					CumulativePayment: big.NewInt(1000),
				},
				reservationWindow: 60,
				pricePerSymbol:    50,
				minNumSymbols:     100,
				cumulativePayment: big.NewInt(500),
				periodRecords: []PeriodRecord{
					{Index: 1, Usage: 150},
					{Index: 0, Usage: 0},
					{Index: 0, Usage: 0},
				},
			},
		},
		{
			name: "successful set payment state with minimal fields",
			state: &disperser_rpc.GetPaymentStateReply{
				PaymentGlobalParams: &disperser_rpc.PaymentGlobalParams{
					MinNumSymbols:     50,
					PricePerSymbol:    25,
					ReservationWindow: 30,
				},
			},
			expectError: false,
			expectedState: &Accountant{
				accountID:         accountId,
				reservation:       emptyReservation,
				onDemand:          emptyOnDemand,
				reservationWindow: 30,
				pricePerSymbol:    25,
				minNumSymbols:     50,
				cumulativePayment: big.NewInt(0),
				periodRecords:     []PeriodRecord{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountant := NewAccountant(accountId, emptyReservation, emptyOnDemand, 0, 0, 0, numBins)
			err := accountant.SetPaymentState(tt.state)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedState.pricePerSymbol, accountant.pricePerSymbol)
				assert.Equal(t, tt.expectedState.reservationWindow, accountant.reservationWindow)
				assert.Equal(t, tt.expectedState.minNumSymbols, accountant.minNumSymbols)
				assert.Equal(t, tt.expectedState.onDemand.CumulativePayment, accountant.onDemand.CumulativePayment)
				assert.Equal(t, tt.expectedState.cumulativePayment, accountant.cumulativePayment)
				assert.Equal(t, tt.expectedState.reservation.SymbolsPerSecond, accountant.reservation.SymbolsPerSecond)
				assert.Equal(t, tt.expectedState.reservation.StartTimestamp, accountant.reservation.StartTimestamp)
				assert.Equal(t, tt.expectedState.reservation.EndTimestamp, accountant.reservation.EndTimestamp)
				assert.Equal(t, tt.expectedState.reservation.QuorumNumbers, accountant.reservation.QuorumNumbers)

				// Check period records
				for i := range tt.expectedState.periodRecords {
					assert.Equal(t, tt.expectedState.periodRecords[i].Index, accountant.periodRecords[i].Index)
					assert.Equal(t, tt.expectedState.periodRecords[i].Usage, accountant.periodRecords[i].Usage)
				}
			}
		})
	}
}

// getRelativePeriodRecord returns the period record for the given index
func getRelativePeriodRecord(index uint64, reservationWindow uint64, periodRecords []PeriodRecord) *PeriodRecord {
	relativeIndex := uint32((index / reservationWindow) % uint64(len(periodRecords)))
	// Return empty record if the index is greater than the number of bins (should never happen by accountant initialization)
	if relativeIndex >= uint32(len(periodRecords)) {
		panic(fmt.Sprintf("relativeIndex %d is greater than the number of bins %d cached", relativeIndex, len(periodRecords)))
	}
	return &periodRecords[relativeIndex]
}
