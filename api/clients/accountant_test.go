package clients

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/disperser/meterer"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountant(t *testing.T) {
	reservation := meterer.ActiveReservation{
		DataRate:      1000,
		StartEpoch:    100,
		EndEpoch:      200,
		QuorumSplit:   []byte{50, 50},
		QuorumNumbers: []uint8{0, 1},
	}
	onDemand := meterer.OnDemandPayment{
		CumulativePayment: 500,
	}
	reservationWindow := uint32(60)
	pricePerChargeable := uint32(1)
	minChargeableSize := uint32(100)

	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerChargeable, minChargeableSize)

	assert.Equal(t, reservation, accountant.reservation)
	assert.Equal(t, onDemand, accountant.onDemand)
	assert.Equal(t, reservationWindow, accountant.reservationWindow)
	assert.Equal(t, pricePerChargeable, accountant.pricePerChargeable)
	assert.Equal(t, minChargeableSize, accountant.minChargeableSize)
	assert.Equal(t, []uint64{0, 0, 0}, accountant.binUsages)
	assert.Equal(t, uint64(0), accountant.cumulativePayment)
}

func TestAccountBlob_Reservation(t *testing.T) {
	reservation := meterer.ActiveReservation{
		DataRate:      1000,
		StartEpoch:    100,
		EndEpoch:      200,
		QuorumSplit:   []byte{50, 50},
		QuorumNumbers: []uint8{0, 1},
	}
	onDemand := meterer.OnDemandPayment{
		CumulativePayment: 500,
	}
	reservationWindow := uint32(60)
	pricePerChargeable := uint32(1)
	minChargeableSize := uint32(100)

	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerChargeable, minChargeableSize)

	ctx := context.Background()
	dataLength := uint64(500)
	quorums := []uint8{0, 1}

	binIndex, payment, err := accountant.AccountBlob(ctx, dataLength, quorums)

	assert.NoError(t, err)
	assert.Equal(t, meterer.GetCurrentBinIndex(reservationWindow), binIndex)
	assert.Equal(t, uint64(0), payment)
	assert.Equal(t, []uint64{500, 0, 0}, accountant.binUsages)

	dataLength = uint64(700)

	binIndex, payment, err = accountant.AccountBlob(ctx, dataLength, quorums)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, binIndex)
	assert.Equal(t, uint64(0), payment)
	assert.Equal(t, []uint64{1200, 0, 200}, accountant.binUsages)

	// Second call should use on-demand payment
	binIndex, payment, err = accountant.AccountBlob(ctx, 300, quorums)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), binIndex)
	assert.Equal(t, uint64(3), payment)
}

func TestAccountBlob_OnDemand(t *testing.T) {
	reservation := meterer.ActiveReservation{
		DataRate:      1000,
		StartEpoch:    100,
		EndEpoch:      200,
		QuorumSplit:   []byte{50, 50},
		QuorumNumbers: []uint8{0, 1},
	}
	onDemand := meterer.OnDemandPayment{
		CumulativePayment: 500,
	}
	reservationWindow := uint32(60)
	pricePerChargeable := uint32(1)
	minChargeableSize := uint32(100)

	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerChargeable, minChargeableSize)

	ctx := context.Background()
	dataLength := uint64(1500)
	quorums := []uint8{0, 1}

	binIndex, payment, err := accountant.AccountBlob(ctx, dataLength, quorums)

	expectedPayment := uint64(dataLength * uint64(pricePerChargeable) / uint64(minChargeableSize))
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), binIndex)
	assert.Equal(t, expectedPayment, payment)
	assert.Equal(t, []uint64{0, 0, 0}, accountant.binUsages)
	assert.Equal(t, expectedPayment, accountant.cumulativePayment)
}

func TestAccountBlob_InsufficientOnDemand(t *testing.T) {
	reservation := meterer.ActiveReservation{}
	onDemand := meterer.OnDemandPayment{
		CumulativePayment: 500,
	}
	reservationWindow := uint32(60)
	pricePerChargeable := uint32(100)
	minChargeableSize := uint32(100)

	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerChargeable, minChargeableSize)

	ctx := context.Background()
	dataLength := uint64(2000)
	quorums := []uint8{0, 1}

	_, _, err := accountant.AccountBlob(ctx, dataLength, quorums)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Accountant cannot approve payment for this blob")
}

func TestAccountBlobCallSeries(t *testing.T) {
	reservation := meterer.ActiveReservation{
		DataRate:      1000,
		StartEpoch:    100,
		EndEpoch:      200,
		QuorumSplit:   []byte{50, 50},
		QuorumNumbers: []uint8{0, 1},
	}
	onDemand := meterer.OnDemandPayment{
		CumulativePayment: 1000,
	}
	reservationWindow := uint32(60)
	pricePerChargeable := uint32(100)
	minChargeableSize := uint32(100)

	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerChargeable, minChargeableSize)

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// First call: Use reservation
	binIndex, payment, err := accountant.AccountBlob(ctx, 800, quorums)
	assert.NoError(t, err)
	assert.Equal(t, meterer.GetCurrentBinIndex(reservationWindow), binIndex)
	assert.Equal(t, uint64(0), payment)

	// Second call: Use remaining reservation + overflow
	binIndex, payment, err = accountant.AccountBlob(ctx, 300, quorums)
	assert.NoError(t, err)
	assert.Equal(t, meterer.GetCurrentBinIndex(reservationWindow), binIndex)
	assert.Equal(t, uint64(0), payment)

	// Third call: Use on-demand
	binIndex, payment, err = accountant.AccountBlob(ctx, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), binIndex)
	assert.Equal(t, uint64(500), payment)

	// Fourth call: Insufficient on-demand
	_, _, err = accountant.AccountBlob(ctx, 600, quorums)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Accountant cannot approve payment for this blob")
}

func TestAccountBlob_BinRotation(t *testing.T) {
	reservation := meterer.ActiveReservation{
		DataRate:      1000,
		StartEpoch:    100,
		EndEpoch:      200,
		QuorumSplit:   []byte{50, 50},
		QuorumNumbers: []uint8{0, 1},
	}
	onDemand := meterer.OnDemandPayment{
		CumulativePayment: 1000,
	}
	reservationWindow := uint32(1) // Set to 1 second for testing
	pricePerChargeable := uint32(1)
	minChargeableSize := uint32(100)

	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerChargeable, minChargeableSize)
	defer accountant.Stop()

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// First call
	_, _, err := accountant.AccountBlob(ctx, 800, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{800, 0, 0}, accountant.binUsages)

	// Wait for bin rotation
	time.Sleep(2 * time.Second)

	// Second call after bin rotation
	_, _, err = accountant.AccountBlob(ctx, 300, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{300, 0, 0}, accountant.binUsages)

	// Third call
	_, _, err = accountant.AccountBlob(ctx, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{800, 0, 0}, accountant.binUsages)
}

func TestBinRotation(t *testing.T) {
	reservation := meterer.ActiveReservation{
		DataRate:      1000,
		StartEpoch:    100,
		EndEpoch:      200,
		QuorumSplit:   []byte{50, 50},
		QuorumNumbers: []uint8{0, 1},
	}
	onDemand := meterer.OnDemandPayment{
		CumulativePayment: 1000,
	}
	reservationWindow := uint32(1) // Set to 1 second for testing
	pricePerChargeable := uint32(1)
	minChargeableSize := uint32(100)

	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerChargeable, minChargeableSize)
	defer accountant.Stop()

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// First call
	_, _, err := accountant.AccountBlob(ctx, 800, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{800, 0, 0}, accountant.binUsages)

	// Second call for overflow
	_, _, err = accountant.AccountBlob(ctx, 800, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{1600, 0, 600}, accountant.binUsages)

	// Wait for bin rotation
	time.Sleep(1200 * time.Millisecond)

	_, _, err = accountant.AccountBlob(ctx, 300, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{300, 600, 0}, accountant.binUsages)

	// another bin rotation
	time.Sleep(1200 * time.Millisecond)

	_, _, err = accountant.AccountBlob(ctx, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{1100, 0, 100}, accountant.binUsages)
}

func TestConcurrentBinRotationAndAccountBlob(t *testing.T) {
	reservation := meterer.ActiveReservation{
		DataRate:      1000,
		StartEpoch:    100,
		EndEpoch:      200,
		QuorumSplit:   []byte{50, 50},
		QuorumNumbers: []uint8{0, 1},
	}
	onDemand := meterer.OnDemandPayment{
		CumulativePayment: 1000,
	}
	reservationWindow := uint32(1) // Set to 1 second for testing
	pricePerChargeable := uint32(1)
	minChargeableSize := uint32(100)

	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerChargeable, minChargeableSize)
	defer accountant.Stop()

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// Start concurrent AccountBlob calls
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				_, _, err := accountant.AccountBlob(ctx, 100, quorums)
				assert.NoError(t, err)
				time.Sleep(500 * time.Millisecond)
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check final state
	assert.Equal(t, uint64(1000), accountant.binUsages[0]+accountant.binUsages[1]+accountant.binUsages[2])
}

func TestAccountBlob_ReservationWithOneOverflow(t *testing.T) {
	reservation := meterer.ActiveReservation{
		DataRate:      1000,
		StartEpoch:    100,
		EndEpoch:      200,
		QuorumSplit:   []byte{50, 50},
		QuorumNumbers: []uint8{0, 1},
	}
	onDemand := meterer.OnDemandPayment{
		CumulativePayment: 1000,
	}
	reservationWindow := uint32(60)
	pricePerChargeable := uint32(1)
	minChargeableSize := uint32(100)

	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerChargeable, minChargeableSize)

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// Okay reservation
	binIndex, payment, err := accountant.AccountBlob(ctx, 800, quorums)
	assert.NoError(t, err)
	assert.Equal(t, meterer.GetCurrentBinIndex(reservationWindow), binIndex)
	assert.Equal(t, uint64(0), payment)
	assert.Equal(t, []uint64{800, 0, 0}, accountant.binUsages)

	// Second call: Allow one overflow
	_, payment, err = accountant.AccountBlob(ctx, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), payment)
	assert.Equal(t, []uint64{1300, 0, 300}, accountant.binUsages)

	// Third call: Should use on-demand payment
	binIndex, payment, err = accountant.AccountBlob(ctx, 200, quorums)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), binIndex)
	assert.Equal(t, uint64(2), payment)
	assert.Equal(t, []uint64{1300, 0, 300}, accountant.binUsages)
}

func TestAccountBlob_ReservationOverflowReset(t *testing.T) {
	reservation := meterer.ActiveReservation{
		DataRate:      1000,
		StartEpoch:    100,
		EndEpoch:      200,
		QuorumSplit:   []byte{50, 50},
		QuorumNumbers: []uint8{0, 1},
	}
	onDemand := meterer.OnDemandPayment{
		CumulativePayment: 1000,
	}
	reservationWindow := uint32(1) // Set to 1 second for testing
	pricePerChargeable := uint32(1)
	minChargeableSize := uint32(100)

	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerChargeable, minChargeableSize)
	defer accountant.Stop()

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// full reservation
	_, _, err := accountant.AccountBlob(ctx, 1000, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{1000, 0, 0}, accountant.binUsages)

	// no overflow
	_, payment, err := accountant.AccountBlob(ctx, 500, quorums)
	assert.Equal(t, []uint64{1000, 0, 0}, accountant.binUsages)
	assert.Equal(t, uint64(5), payment)

	// Wait for bin rotation
	time.Sleep(1500 * time.Millisecond)

	// Third call: Should use new bin and allow overflow again
	_, _, err = accountant.AccountBlob(ctx, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{500, 0, 0}, accountant.binUsages)
}
