package clients

import (
	"context"
	"encoding/hex"
	"math/big"
	"sync"
	"testing"
	"time"

	commonpbv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

const numBins = uint32(3)

// Helper function to create standard PaymentVaultParams for testing
func createTestPaymentVaultParams(reservationWindow, pricePerSymbol, minNumSymbols uint64) *commonpbv2.PaymentVaultParams {
	return &commonpbv2.PaymentVaultParams{
		QuorumPaymentConfigs: map[uint32]*commonpbv2.PaymentQuorumConfig{
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
		QuorumProtocolConfigs: map[uint32]*commonpbv2.PaymentQuorumProtocolConfig{
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
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 100,
		StartTimestamp:   100,
		EndTimestamp:     200,
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.ReservedPayment{0: reservation, 1: reservation}

	accountant := NewAccountant(accountId, reservations, onDemand, numBins, testutils.GetLogger())

	// Create payment state with test configurations
	paymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentVaultParams:       createTestPaymentVaultParams(6, 1, 100),
		OnchainCumulativePayment: big.NewInt(500).Bytes(),
		Reservations: []*v2.QuorumReservation{
			{
				QuorumNumber:     0,
				SymbolsPerSecond: 100,
				StartTimestamp:   100,
				EndTimestamp:     200,
			},
			{
				QuorumNumber:     1,
				SymbolsPerSecond: 100,
				StartTimestamp:   100,
				EndTimestamp:     200,
			},
		},
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	assert.NotNil(t, accountant)
	// After SetPaymentState, reservations and onDemand should be updated from the payment state
	assert.Equal(t, 2, len(accountant.reservation))
	assert.Equal(t, uint64(100), accountant.reservation[0].SymbolsPerSecond)
	assert.Equal(t, uint64(100), accountant.reservation[1].SymbolsPerSecond)
	assert.Equal(t, big.NewInt(500), accountant.onDemand.CumulativePayment)

	// Check initialization of periodRecords
	for _, records := range accountant.periodRecords {
		assert.Equal(t, int(numBins), len(records), "Should have numBins records for each quorum")
		for i, record := range records {
			assert.Equal(t, uint32(i), record.Index)
			assert.Equal(t, uint64(0), record.Usage)
		}
	}

	assert.Equal(t, big.NewInt(0), accountant.cumulativePayment)
}

func TestAccountBlob_Reservation(t *testing.T) {
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 200,
		StartTimestamp:   100,
		EndTimestamp:     200,
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}
	reservationWindow := uint64(5)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.ReservedPayment{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins, testutils.GetLogger())

	// Create payment state with test configurations
	paymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(500).Bytes(),
		Reservations: []*v2.QuorumReservation{
			{
				QuorumNumber:     0,
				SymbolsPerSecond: 200,
				StartTimestamp:   100,
				EndTimestamp:     200,
			},
			{
				QuorumNumber:     1,
				SymbolsPerSecond: 200,
				StartTimestamp:   100,
				EndTimestamp:     200,
			},
		},
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
	assert.Equal(t, meterer.GetReservationPeriod(time.Now().Unix(), reservationWindow), meterer.GetReservationPeriodByNanosecond(header.Timestamp, reservationWindow))
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// Check that usage for quorum 0 was updated properly
	quorum0Records := accountant.periodRecords[0]
	currentPeriod := meterer.GetReservationPeriodByNanosecond(now, reservationWindow) % uint64(numBins)
	assert.Equal(t, symbolLength, quorum0Records[currentPeriod].Usage)

	symbolLength = uint64(700)

	now = time.Now().UnixNano()
	header, err = accountant.AccountBlob(ctx, now, symbolLength, quorums)

	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), header.Timestamp)
	assert.Equal(t, big.NewInt(0), header.CumulativePayment)

	// With overflow, usage should be at the limit in the current period
	binLimit := reservations[0].SymbolsPerSecond * uint64(reservationWindow)
	expectedUsage := binLimit
	assert.Equal(t, expectedUsage, quorum0Records[currentPeriod].Usage)

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

	numSymbols := uint64(1500)
	quorums := []uint8{0, 1}

	reservations := map[uint8]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 200,
			StartTimestamp:   100,
			EndTimestamp:     200,
		},
		1: {
			SymbolsPerSecond: 200,
			StartTimestamp:   100,
			EndTimestamp:     200,
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1500),
	}
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservations, onDemand, numBins, testutils.GetLogger())

	// Create payment state with test configurations
	paymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentVaultParams:       createTestPaymentVaultParams(5, 1, 100),
		OnchainCumulativePayment: big.NewInt(1500).Bytes(),
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	ctx := context.Background()
	now := time.Now().UnixNano()
	header, err := accountant.AccountBlob(ctx, now, numSymbols, quorums)
	assert.NoError(t, err)

	expectedPayment := big.NewInt(int64(numSymbols * 1)) // pricePerSymbol = 1
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
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId, reservation, onDemand, numBins, testutils.GetLogger())

	// Create payment state with test configurations
	paymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentVaultParams:       createTestPaymentVaultParams(60, 100, 100),
		OnchainCumulativePayment: big.NewInt(500).Bytes(),
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

	ctx := context.Background()
	numSymbols := uint64(2000)
	quorums := []uint8{0, 1}
	now := time.Now().UnixNano()
	_, err = accountant.AccountBlob(ctx, now, numSymbols, quorums)
	assert.Contains(t, err.Error(), "no bandwidth reservation found for account")
}

func TestAccountBlobCallSeries(t *testing.T) {
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 200,
		StartTimestamp:   100,
		EndTimestamp:     200,
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(5)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.ReservedPayment{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins, testutils.GetLogger())

	// Create payment state with test configurations
	paymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: []*v2.QuorumReservation{
			{
				QuorumNumber:     0,
				SymbolsPerSecond: 200,
				StartTimestamp:   100,
				EndTimestamp:     200,
			},
			{
				QuorumNumber:     1,
				SymbolsPerSecond: 200,
				StartTimestamp:   100,
				EndTimestamp:     200,
			},
		},
	}
	err = accountant.SetPaymentState(paymentState)
	assert.NoError(t, err)

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
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 1000,
		StartTimestamp:   100,
		EndTimestamp:     200,
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(1) // Set to 1 second for testing

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.ReservedPayment{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins, testutils.GetLogger())

	// Create payment state with test configurations
	paymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: []*v2.QuorumReservation{
			{
				QuorumNumber:     0,
				SymbolsPerSecond: 1000,
				StartTimestamp:   100,
				EndTimestamp:     200,
			},
			{
				QuorumNumber:     1,
				SymbolsPerSecond: 1000,
				StartTimestamp:   100,
				EndTimestamp:     200,
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
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 1000,
		StartTimestamp:   100,
		EndTimestamp:     200,
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(1) // Set to 1 second for testing

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.ReservedPayment{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins, testutils.GetLogger())

	// Create payment state with test configurations
	paymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: []*v2.QuorumReservation{
			{
				QuorumNumber:     0,
				SymbolsPerSecond: 1000,
				StartTimestamp:   100,
				EndTimestamp:     200,
			},
			{
				QuorumNumber:     1,
				SymbolsPerSecond: 1000,
				StartTimestamp:   100,
				EndTimestamp:     200,
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
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 200,
		StartTimestamp:   100,
		EndTimestamp:     200,
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(5)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.ReservedPayment{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins, testutils.GetLogger())

	// Create payment state with test configurations
	paymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: []*v2.QuorumReservation{
			{
				QuorumNumber:     0,
				SymbolsPerSecond: 200,
				StartTimestamp:   100,
				EndTimestamp:     200,
			},
			{
				QuorumNumber:     1,
				SymbolsPerSecond: 200,
				StartTimestamp:   100,
				EndTimestamp:     200,
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
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 1000,
		StartTimestamp:   100,
		EndTimestamp:     200,
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint64(1) // Set to 1 second for testing

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	reservations := map[uint8]*core.ReservedPayment{0: reservation, 1: reservation}
	accountant := NewAccountant(accountId, reservations, onDemand, numBins, testutils.GetLogger())

	// Create payment state with test configurations
	paymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentVaultParams:       createTestPaymentVaultParams(reservationWindow, 1, 100),
		OnchainCumulativePayment: big.NewInt(1000).Bytes(),
		Reservations: []*v2.QuorumReservation{
			{
				QuorumNumber:     0,
				SymbolsPerSecond: 1000,
				StartTimestamp:   100,
				EndTimestamp:     200,
			},
			{
				QuorumNumber:     1,
				SymbolsPerSecond: 1000,
				StartTimestamp:   100,
				EndTimestamp:     200,
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
		numBins, testutils.GetLogger())

	// Create payment state reply with sample data
	paymentState := &v2.GetQuorumSpecificPaymentStateReply{
		PaymentVaultParams: &commonpbv2.PaymentVaultParams{
			QuorumPaymentConfigs: map[uint32]*commonpbv2.PaymentQuorumConfig{
				0: {
					ReservationSymbolsPerSecond: 2000,
					OnDemandSymbolsPerSecond:    1000,
					OnDemandPricePerSymbol:      2,
				},
				1: {
					ReservationSymbolsPerSecond: 3000,
					OnDemandSymbolsPerSecond:    1500,
					OnDemandPricePerSymbol:      3,
				},
			},
			QuorumProtocolConfigs: map[uint32]*commonpbv2.PaymentQuorumProtocolConfig{
				0: {
					MinNumSymbols:              200,
					ReservationAdvanceWindow:   20,
					ReservationRateLimitWindow: 60,
					OnDemandRateLimitWindow:    30,
					OnDemandEnabled:            true,
				},
				1: {
					MinNumSymbols:              300,
					ReservationAdvanceWindow:   25,
					ReservationRateLimitWindow: 70,
					OnDemandRateLimitWindow:    35,
					OnDemandEnabled:            true,
				},
			},
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

	// Verify per-quorum params are updated
	assert.Equal(t, uint64(200), accountant.GetMinNumSymbols(core.QuorumID(0)))
	assert.Equal(t, uint64(2), accountant.GetPricePerSymbol(core.QuorumID(0)))
	assert.Equal(t, uint64(60), accountant.GetReservationWindow(core.QuorumID(0)))
	assert.Equal(t, uint64(300), accountant.GetMinNumSymbols(core.QuorumID(1)))
	assert.Equal(t, uint64(3), accountant.GetPricePerSymbol(core.QuorumID(1)))
	assert.Equal(t, uint64(70), accountant.GetReservationWindow(core.QuorumID(1)))

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
		PaymentVaultParams: &commonpbv2.PaymentVaultParams{
			QuorumPaymentConfigs: map[uint32]*commonpbv2.PaymentQuorumConfig{
				0: {
					ReservationSymbolsPerSecond: 4000,
					OnDemandSymbolsPerSecond:    2000,
					OnDemandPricePerSymbol:      3,
				},
			},
			QuorumProtocolConfigs: map[uint32]*commonpbv2.PaymentQuorumProtocolConfig{
				0: {
					MinNumSymbols:              300,
					ReservationAdvanceWindow:   30,
					ReservationRateLimitWindow: 90,
					OnDemandRateLimitWindow:    45,
					OnDemandEnabled:            true,
				},
			},
		},
		// No CumulativePayment
		// No OnchainCumulativePayment
		// No Reservations
	}

	err = accountant.SetPaymentState(nilPaymentState)
	assert.NoError(t, err)

	// Verify per-quorum params are updated
	assert.Equal(t, uint64(300), accountant.GetMinNumSymbols(core.QuorumID(0)))
	assert.Equal(t, uint64(3), accountant.GetPricePerSymbol(core.QuorumID(0)))
	assert.Equal(t, uint64(90), accountant.GetReservationWindow(core.QuorumID(0)))

	// Verify defaults when fields are nil
	assert.Equal(t, big.NewInt(0), accountant.onDemand.CumulativePayment)
	assert.Equal(t, big.NewInt(0), accountant.cumulativePayment)
	assert.Equal(t, 0, len(accountant.reservation))

	// Test with nil payment state
	err = accountant.SetPaymentState(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payment state cannot be nil")

	// Test with nil payment vault params
	nilVaultParams := &v2.GetQuorumSpecificPaymentStateReply{
		// No PaymentVaultParams
	}
	err = accountant.SetPaymentState(nilVaultParams)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payment vault params cannot be nil")
}
