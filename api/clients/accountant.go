package clients

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/meterer"
	"github.com/ethereum/go-ethereum/common"
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

	domainSigner     *meterer.EIP712Signer
	accountantSigner *ecdsa.PrivateKey
}

func NewAccountant(reservation meterer.ActiveReservation, onDemand meterer.OnDemandPayment, reservationWindow uint32, pricePerChargeable uint32, minChargeableSize uint32, accountantSigner *ecdsa.PrivateKey) Accountant {
	//todo: refactor signing strategy; EIP712 -> serializer? private key -> Local blob request signer?
	chainID := big.NewInt(17000)
	verifyingContract := common.HexToAddress("0x1234000000000000000000000000000000000000")

	a := Accountant{
		reservation:        reservation,
		onDemand:           onDemand,
		reservationWindow:  reservationWindow,
		pricePerChargeable: pricePerChargeable,
		minChargeableSize:  minChargeableSize,
		binUsages:          []uint64{0, 0, 0},
		cumulativePayment:  0,
		stopRotation:       make(chan struct{}),
		domainSigner:       meterer.NewEIP712Signer(chainID, verifyingContract),
		accountantSigner:   accountantSigner,
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

// accountant calculates and records payment information
func (a *Accountant) BlobPaymentInfo(ctx context.Context, dataLength uint64) (uint32, uint64, error) {
	//TODO: do we need to lock the binUsages here in case the blob rotation happens in the middle of the function?
	binUsage := a.binUsages[0] + dataLength
	now := time.Now().Unix()
	currentBinIndex := meterer.GetBinIndex(uint64(now), a.reservationWindow)

	// first attempt to use the active reservation
	if binUsage <= a.reservation.DataRate {
		return currentBinIndex, 0, nil
	}

	// Allow one overflow when the overflow bin is empty, the current usage and new length are both less than the limit
	if a.binUsages[2] == 0 && binUsage-dataLength < a.reservation.DataRate && dataLength <= a.reservation.DataRate {
		a.binUsages[2] += binUsage - a.reservation.DataRate
		return currentBinIndex, 0, nil
	}

	// reservation not available, attempt on-demand
	//todo: rollback if disperser respond with some type of rejection?
	incrementRequired := a.PaymentCharged(uint32(dataLength))
	a.cumulativePayment += incrementRequired
	if a.cumulativePayment <= uint64(a.onDemand.CumulativePayment) {
		return 0, a.cumulativePayment, nil
	}
	return 0, 0, errors.New("Accountant cannot approve payment for this blob")
}

// accountant provides and records payment information
func (a *Accountant) AccountBlob(ctx context.Context, dataLength uint64, quorums []uint8) (*meterer.BlobHeader, error) {
	binIndex, cumulativePayment, err := a.BlobPaymentInfo(ctx, dataLength)
	if err != nil {
		return nil, err
	}

	// TODO: have client first generate a commitment, then pass it here
	commitment := core.NewG1Point(big.NewInt(0), big.NewInt(1))
	header, err := meterer.ConstructBlobHeader(a.domainSigner, binIndex, cumulativePayment, *commitment, uint32(dataLength), quorums, a.accountantSigner)
	if err != nil {
		return nil, err
	}
	return header, nil
}

// PaymentCharged returns the chargeable price for a given data length
func (a *Accountant) PaymentCharged(dataLength uint32) uint64 {
	return uint64(core.RoundUpDivide(uint(a.BlobSizeCharged(dataLength)*a.pricePerChargeable), uint(a.minChargeableSize)))
}

// BlobSizeCharged returns the chargeable data length for a given data length
func (a *Accountant) BlobSizeCharged(dataLength uint32) uint32 {
	return max(dataLength, uint32(a.minChargeableSize))
}
