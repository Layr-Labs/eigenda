package clients

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/disperser/meterer"
)

type IAccountant interface {
	AccountBlob(ctx context.Context, data []byte, quorums []uint8) (uint32, uint64, error)
}

type Accountant struct {
	// on-chain states
	reservation        meterer.ActiveReservation
	onDemand           meterer.OnDemandPayment
	reservationWindow  uint32
	pricePerChargeable uint32
	minChargeableSize  uint32

	// local accounting
	// contains 3 bins; 0 for current bin, 1 for next bin, 2 for overflowed bin
	binUsages         []uint64
	cumulativePayment uint64
	stopRotation      chan struct{}
}

func NewAccountant(reservation meterer.ActiveReservation, onDemand meterer.OnDemandPayment, reservationWindow uint32, pricePerChargeable uint32, minChargeableSize uint32) Accountant {
	a := Accountant{
		reservation:        reservation,
		onDemand:           onDemand,
		reservationWindow:  reservationWindow,
		pricePerChargeable: pricePerChargeable,
		minChargeableSize:  minChargeableSize,
		binUsages:          []uint64{0, 0, 0},
		cumulativePayment:  0,
		stopRotation:       make(chan struct{}),
	}
	go a.startBinRotation()
	return a
}

func (a *Accountant) startBinRotation() {
	ticker := time.NewTicker(time.Duration(a.reservationWindow) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.rotateBins()
		case <-a.stopRotation:
			return
		}
	}
}

func (a *Accountant) rotateBins() {
	// Shift bins: bin_i to bin_{i-1}, add 0 to bin2
	a.binUsages[0] = a.binUsages[1]
	a.binUsages[1] = a.binUsages[2]
	a.binUsages[2] = 0
}

func (a *Accountant) Stop() {
	close(a.stopRotation)
}

// accountant provides and records payment information
func (a *Accountant) AccountBlob(ctx context.Context, dataLength uint64, quorums []uint8) (uint32, uint64, error) {
	//TODO: do we need to lock the binUsages here in case the rotation happens in the middle of the function?
	currentBinUsage := a.binUsages[0]
	currentBinIndex := meterer.GetCurrentBinIndex(a.reservationWindow)

	// first attempt to use the active reservation
	if currentBinUsage+dataLength <= a.reservation.DataRate {
		a.binUsages[0] += dataLength
		return currentBinIndex, 0, nil
	}

	// Allow one overflow when the overflow bin is empty, the current usage and new length are both less than the limit
	if a.binUsages[2] == 0 && currentBinUsage < a.reservation.DataRate && dataLength <= a.reservation.DataRate {
		fmt.Println("in overflow:", currentBinUsage, dataLength, a.reservation.DataRate)
		a.binUsages[0] += dataLength
		a.binUsages[2] += currentBinUsage + dataLength - a.reservation.DataRate
		return currentBinIndex, 0, nil
	}

	fmt.Println("in ondemand:", currentBinUsage, dataLength, a.reservation.DataRate)
	// reservation not available, attempt on-demand
	//todo: rollback if disperser respond with some type of rejection?
	incrementRequired := uint64(max(uint32(dataLength), a.minChargeableSize)) * uint64(a.pricePerChargeable) / uint64(a.minChargeableSize)
	a.cumulativePayment += incrementRequired
	if a.cumulativePayment <= uint64(a.onDemand.CumulativePayment) {
		return 0, a.cumulativePayment, nil
	}
	return 0, 0, errors.New("Accountant cannot approve payment for this blob")
}
