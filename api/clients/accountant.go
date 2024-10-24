package clients

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
)

type IAccountant interface {
	AccountBlob(ctx context.Context, data []byte, quorums []uint8) (uint32, uint64, error)
}

type Accountant struct {
	// on-chain states
	reservation       core.ActiveReservation
	onDemand          core.OnDemandPayment
	reservationWindow uint32
	pricePerSymbol    uint32
	minNumSymbols     uint32

	// local accounting
	// contains 3 bins; 0 for current bin, 1 for next bin, 2 for overflowed bin
	binUsages         []uint64
	cumulativePayment *big.Int
	stopRotation      chan struct{}

	paymentSigner core.PaymentSigner
}

func NewAccountant(reservation core.ActiveReservation, onDemand core.OnDemandPayment, reservationWindow uint32, pricePerSymbol uint32, minNumSymbols uint32, paymentSigner core.PaymentSigner) Accountant {
	//TODO: client storage; currently every instance starts fresh but on-chain or a small store makes more sense
	// Also client is currently responsible for supplying network params, we need to add RPC in order to be automatic
	a := Accountant{
		reservation:       reservation,
		onDemand:          onDemand,
		reservationWindow: reservationWindow,
		pricePerSymbol:    pricePerSymbol,
		minNumSymbols:     minNumSymbols,
		binUsages:         []uint64{0, 0, 0},
		cumulativePayment: big.NewInt(0),
		stopRotation:      make(chan struct{}),
		paymentSigner:     paymentSigner,
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
func (a *Accountant) BlobPaymentInfo(ctx context.Context, dataLength uint64) (uint32, *big.Int, error) {
	//TODO: do we need to lock the binUsages here in case the blob rotation happens in the middle of the function?
	// binUsage := a.binUsages[0] + dataLength
	a.binUsages[0] += dataLength
	now := time.Now().Unix()
	currentBinIndex := meterer.GetBinIndex(uint64(now), a.reservationWindow)

	// first attempt to use the active reservation
	binLimit := a.reservation.SymbolsPerSec * uint64(a.reservationWindow)
	if a.binUsages[0] <= binLimit {
		return currentBinIndex, big.NewInt(0), nil
	}

	// Allow one overflow when the overflow bin is empty, the current usage and new length are both less than the limit
	if a.binUsages[2] == 0 && a.binUsages[0]-dataLength < binLimit && dataLength <= binLimit {
		a.binUsages[2] += a.binUsages[0] - binLimit
		return currentBinIndex, big.NewInt(0), nil
	}

	// reservation not available, attempt on-demand
	//todo: rollback if disperser respond with some type of rejection?
	a.binUsages[0] -= dataLength
	incrementRequired := big.NewInt(int64(a.PaymentCharged(uint32(dataLength))))
	a.cumulativePayment.Add(a.cumulativePayment, incrementRequired)
	if a.cumulativePayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		return 0, a.cumulativePayment, nil
	}
	fmt.Println("Accountant cannot approve payment for this blob")
	return 0, big.NewInt(0), errors.New("Accountant cannot approve payment for this blob")
}

// accountant provides and records payment information
func (a *Accountant) AccountBlob(ctx context.Context, dataLength uint64, quorums []uint8) (*commonpb.PaymentHeader, []byte, error) {
	binIndex, cumulativePayment, err := a.BlobPaymentInfo(ctx, dataLength)
	if err != nil {
		return nil, nil, err
	}

	accountID, err := a.paymentSigner.GetAccountID()
	if err != nil {
		return nil, nil, err
	}
	pm := &core.PaymentMetadata{
		AccountID:         accountID,
		BinIndex:          binIndex,
		CumulativePayment: cumulativePayment,
	}
	protoPaymentHeader := pm.ConvertToProtoPaymentHeader()

	signature, err := a.paymentSigner.SignBlobPayment(protoPaymentHeader)
	if err != nil {
		return nil, nil, err
	}

	return protoPaymentHeader, signature, nil
}

// PaymentCharged returns the chargeable price for a given data length
func (a *Accountant) PaymentCharged(dataLength uint32) uint64 {
	return uint64(core.RoundUpDivide(uint(a.BlobSizeCharged(dataLength)*a.pricePerSymbol), uint(a.minNumSymbols)))
}

// BlobSizeCharged returns the chargeable data length for a given data length
func (a *Accountant) BlobSizeCharged(dataLength uint32) uint32 {
	return max(dataLength, uint32(a.minNumSymbols))
}

func (a *Accountant) SetPaymentState(paymentState *disperser_rpc.GetPaymentStateReply) {
	quorumNumbers := make([]uint8, len(paymentState.Reservation.QuorumNumbers))
	for i, quorum := range paymentState.Reservation.QuorumNumbers {
		quorumNumbers[i] = uint8(quorum)
	}
	quorumSplit := make([]uint8, len(paymentState.Reservation.QuorumSplit))
	for i, quorum := range paymentState.Reservation.QuorumSplit {
		quorumSplit[i] = uint8(quorum)
	}
	a.onDemand.CumulativePayment = new(big.Int).SetBytes(paymentState.OnChainCumulativePayment)
	a.reservation.SymbolsPerSec = uint64(paymentState.PaymentGlobalParams.GlobalSymbolsPerSecond)
	a.reservation.StartTimestamp = uint64(paymentState.Reservation.StartTimestamp)
	a.reservation.EndTimestamp = uint64(paymentState.Reservation.EndTimestamp)
	a.reservation.QuorumNumbers = quorumNumbers
	a.reservation.QuorumSplit = quorumSplit
	a.reservationWindow = uint32(paymentState.PaymentGlobalParams.ReservationWindow)
	a.pricePerSymbol = uint32(paymentState.PaymentGlobalParams.PricePerSymbol)
	a.minNumSymbols = uint32(paymentState.PaymentGlobalParams.MinNumSymbols)
}
