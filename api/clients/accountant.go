package clients

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"sync"
	"time"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
)

var requiredQuorums = []uint8{0, 1}

type Accountant struct {
	// on-chain states
	accountID         string
	reservation       *core.ActiveReservation
	onDemand          *core.OnDemandPayment
	reservationWindow uint32
	pricePerSymbol    uint32
	minNumSymbols     uint32

	// local accounting
	// contains 3 bins; circular wrapping of indices
	binRecords        []BinRecord
	usageLock         sync.Mutex
	cumulativePayment *big.Int

	// number of bins in the circular accounting, restricted by minNumBins which is 3
	numBins uint32
}

type BinRecord struct {
	Index uint32
	Usage uint64
}

func NewAccountant(accountID string, reservation *core.ActiveReservation, onDemand *core.OnDemandPayment, reservationWindow uint32, pricePerSymbol uint32, minNumSymbols uint32, numBins uint32) *Accountant {
	//TODO: client storage; currently every instance starts fresh but on-chain or a small store makes more sense
	// Also client is currently responsible for supplying network params, we need to add RPC in order to be automatic
	// There's a subsequent PR that handles populating the accountant with on-chain state from the disperser
	binRecords := make([]BinRecord, numBins)
	for i := range binRecords {
		binRecords[i] = BinRecord{Index: uint32(i), Usage: 0}
	}
	a := Accountant{
		accountID:         accountID,
		reservation:       reservation,
		onDemand:          onDemand,
		reservationWindow: reservationWindow,
		pricePerSymbol:    pricePerSymbol,
		minNumSymbols:     minNumSymbols,
		binRecords:        binRecords,
		cumulativePayment: big.NewInt(0),
		numBins:           max(numBins, uint32(meterer.MinNumBins)),
	}
	// TODO: add a routine to refresh the on-chain state occasionally?
	return &a
}

// BlobPaymentInfo calculates and records payment information. The accountant
// will attempt to use the active reservation first and check for quorum settings,
// then on-demand if the reservation is not available. The returned values are
// bin index for reservation payments and cumulative payment for on-demand payments,
// and both fields are used to create the payment header and signature
func (a *Accountant) BlobPaymentInfo(ctx context.Context, numSymbols uint64, quorumNumbers []uint8) (uint32, *big.Int, error) {
	now := time.Now().Unix()
	currentBinIndex := meterer.GetBinIndex(uint64(now), a.reservationWindow)

	a.usageLock.Lock()
	defer a.usageLock.Unlock()
	relativeBinRecord := a.GetRelativeBinRecord(currentBinIndex)
	relativeBinRecord.Usage += numSymbols

	// first attempt to use the active reservation
	binLimit := a.reservation.SymbolsPerSec * uint64(a.reservationWindow)
	if relativeBinRecord.Usage <= binLimit {
		if err := QuorumCheck(quorumNumbers, a.reservation.QuorumNumbers); err != nil {
			return 0, big.NewInt(0), err
		}
		return currentBinIndex, big.NewInt(0), nil
	}

	overflowBinRecord := a.GetRelativeBinRecord(currentBinIndex + 2)
	// Allow one overflow when the overflow bin is empty, the current usage and new length are both less than the limit
	if overflowBinRecord.Usage == 0 && relativeBinRecord.Usage-numSymbols < binLimit && numSymbols <= binLimit {
		overflowBinRecord.Usage += relativeBinRecord.Usage - binLimit
		if err := QuorumCheck(quorumNumbers, a.reservation.QuorumNumbers); err != nil {
			return 0, big.NewInt(0), err
		}
		return currentBinIndex, big.NewInt(0), nil
	}

	// reservation not available, attempt on-demand
	//todo: rollback later if disperser respond with some type of rejection?
	relativeBinRecord.Usage -= numSymbols
	incrementRequired := big.NewInt(int64(a.PaymentCharged(uint(numSymbols))))
	a.cumulativePayment.Add(a.cumulativePayment, incrementRequired)
	if a.cumulativePayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		if err := QuorumCheck(quorumNumbers, requiredQuorums); err != nil {
			return 0, big.NewInt(0), err
		}
		return 0, a.cumulativePayment, nil
	}
	return 0, big.NewInt(0), fmt.Errorf("neither reservation nor on-demand payment is available")
}

// AccountBlob accountant provides and records payment information
func (a *Accountant) AccountBlob(ctx context.Context, numSymbols uint64, quorums []uint8) (*core.PaymentMetadata, error) {
	binIndex, cumulativePayment, err := a.BlobPaymentInfo(ctx, numSymbols, quorums)
	if err != nil {
		return nil, err
	}

	pm := &core.PaymentMetadata{
		AccountID:         a.accountID,
		BinIndex:          binIndex,
		CumulativePayment: cumulativePayment,
	}

	return pm, nil
}

// TODO: PaymentCharged and SymbolsCharged copied from meterer, should be refactored
// PaymentCharged returns the chargeable price for a given data length
func (a *Accountant) PaymentCharged(numSymbols uint) uint64 {
	return uint64(a.SymbolsCharged(numSymbols)) * uint64(a.pricePerSymbol)
}

// SymbolsCharged returns the number of symbols charged for a given data length
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func (a *Accountant) SymbolsCharged(numSymbols uint) uint32 {
	if numSymbols <= uint(a.minNumSymbols) {
		return a.minNumSymbols
	}
	// Round up to the nearest multiple of MinNumSymbols
	return uint32(core.RoundUpDivide(uint(numSymbols), uint(a.minNumSymbols))) * a.minNumSymbols
}

func (a *Accountant) GetRelativeBinRecord(index uint32) *BinRecord {
	relativeIndex := index % a.numBins
	if a.binRecords[relativeIndex].Index != uint32(index) {
		a.binRecords[relativeIndex] = BinRecord{
			Index: uint32(index),
			Usage: 0,
		}
	}

	return &a.binRecords[relativeIndex]
}

func (a *Accountant) SetPaymentState(paymentState *disperser_rpc.GetPaymentStateReply) error {
	if paymentState == nil {
		return fmt.Errorf("payment state cannot be nil")
	} else if paymentState.GetPaymentGlobalParams() == nil {
		return fmt.Errorf("payment global params cannot be nil")
	} else if paymentState.GetOnchainCumulativePayment() == nil {
		return fmt.Errorf("onchain cumulative payment cannot be nil")
	} else if paymentState.GetCumulativePayment() == nil {
		return fmt.Errorf("cumulative payment cannot be nil")
	} else if paymentState.GetReservation() == nil {
		return fmt.Errorf("reservation cannot be nil")
	} else if paymentState.GetReservation().GetQuorumNumbers() == nil {
		return fmt.Errorf("reservation quorum numbers cannot be nil")
	} else if paymentState.GetReservation().GetQuorumSplit() == nil {
		return fmt.Errorf("reservation quorum split cannot be nil")
	} else if paymentState.GetBinRecords() == nil {
		return fmt.Errorf("bin records cannot be nil")
	}

	a.minNumSymbols = uint32(paymentState.PaymentGlobalParams.MinNumSymbols)
	a.onDemand.CumulativePayment = new(big.Int).SetBytes(paymentState.OnchainCumulativePayment)
	a.cumulativePayment = new(big.Int).SetBytes(paymentState.CumulativePayment)
	a.pricePerSymbol = uint32(paymentState.PaymentGlobalParams.PricePerSymbol)

	a.reservation.SymbolsPerSec = uint64(paymentState.PaymentGlobalParams.GlobalSymbolsPerSecond)
	a.reservation.StartTimestamp = uint64(paymentState.Reservation.StartTimestamp)
	a.reservation.EndTimestamp = uint64(paymentState.Reservation.EndTimestamp)
	a.reservationWindow = uint32(paymentState.PaymentGlobalParams.ReservationWindow)

	quorumNumbers := make([]uint8, len(paymentState.Reservation.QuorumNumbers))
	for i, quorum := range paymentState.Reservation.QuorumNumbers {
		quorumNumbers[i] = uint8(quorum)
	}
	a.reservation.QuorumNumbers = quorumNumbers

	quorumSplit := make([]uint8, len(paymentState.Reservation.QuorumSplit))
	for i, quorum := range paymentState.Reservation.QuorumSplit {
		quorumSplit[i] = uint8(quorum)
	}
	a.reservation.QuorumSplit = quorumSplit

	binRecords := make([]BinRecord, len(paymentState.BinRecords))
	for i, record := range paymentState.BinRecords {
		binRecords[i] = BinRecord{
			Index: record.Index,
			Usage: record.Usage,
		}
	}
	a.binRecords = binRecords

	return nil
}

// QuorumCheck eagerly returns error if the check finds a quorum number not an element of the allowed quorum numbers
func QuorumCheck(quorumNumbers []uint8, allowedNumbers []uint8) error {
	if len(quorumNumbers) == 0 {
		return fmt.Errorf("no quorum numbers provided")
	}
	for _, quorum := range quorumNumbers {
		if !slices.Contains(allowedNumbers, quorum) {
			return fmt.Errorf("provided quorum number %v not allowed", quorum)
		}
	}
	return nil
}
