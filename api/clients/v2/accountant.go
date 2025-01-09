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
	reservation       *core.ReservedPayment
	onDemand          *core.OnDemandPayment
	reservationWindow uint32
	pricePerSymbol    uint32
	minNumSymbols     uint32

	// local accounting
	// contains 3 bins; circular wrapping of indices
	periodRecords     []PeriodRecord
	usageLock         sync.Mutex
	cumulativePayment *big.Int

	// number of bins in the circular accounting, restricted by minNumBins which is 3
	numBins uint32
}

type PeriodRecord struct {
	Index uint32
	Usage uint64
}

func NewAccountant(accountID string, reservation *core.ReservedPayment, onDemand *core.OnDemandPayment, reservationWindow uint32, pricePerSymbol uint32, minNumSymbols uint32, numBins uint32) *Accountant {
	periodRecords := make([]PeriodRecord, numBins)
	for i := range periodRecords {
		periodRecords[i] = PeriodRecord{Index: uint32(i), Usage: 0}
	}
	a := Accountant{
		accountID:         accountID,
		reservation:       reservation,
		onDemand:          onDemand,
		reservationWindow: reservationWindow,
		pricePerSymbol:    pricePerSymbol,
		minNumSymbols:     minNumSymbols,
		periodRecords:     periodRecords,
		cumulativePayment: big.NewInt(0),
		numBins:           max(numBins, uint32(meterer.MinNumBins)),
	}
	// TODO: add a routine to refresh the on-chain state occasionally?
	return &a
}

// BlobPaymentInfo calculates and records payment information. The accountant
// will attempt to use the active reservation first and check for quorum settings,
// then on-demand if the reservation is not available. The returned values are
// reservation period for reservation payments and cumulative payment for on-demand payments,
// and both fields are used to create the payment header and signature
func (a *Accountant) BlobPaymentInfo(ctx context.Context, numSymbols uint32, quorumNumbers []uint8) (uint32, *big.Int, error) {
	now := time.Now().Unix()
	currentReservationPeriod := meterer.GetReservationPeriod(uint64(now), a.reservationWindow)
	symbolUsage := uint64(a.SymbolsCharged(numSymbols))

	a.usageLock.Lock()
	defer a.usageLock.Unlock()
	relativePeriodRecord := a.GetRelativePeriodRecord(currentReservationPeriod)
	relativePeriodRecord.Usage += symbolUsage

	// first attempt to use the active reservation
	binLimit := a.reservation.SymbolsPerSecond * uint64(a.reservationWindow)
	if relativePeriodRecord.Usage <= binLimit {
		if err := QuorumCheck(quorumNumbers, a.reservation.QuorumNumbers); err != nil {
			return 0, big.NewInt(0), err
		}
		return currentReservationPeriod, big.NewInt(0), nil
	}

	overflowPeriodRecord := a.GetRelativePeriodRecord(currentReservationPeriod + 2)
	// Allow one overflow when the overflow bin is empty, the current usage and new length are both less than the limit
	if overflowPeriodRecord.Usage == 0 && relativePeriodRecord.Usage-symbolUsage < binLimit && symbolUsage <= binLimit {
		overflowPeriodRecord.Usage += relativePeriodRecord.Usage - binLimit
		if err := QuorumCheck(quorumNumbers, a.reservation.QuorumNumbers); err != nil {
			return 0, big.NewInt(0), err
		}
		return currentReservationPeriod, big.NewInt(0), nil
	}

	// reservation not available, rollback reservation records, attempt on-demand
	//todo: rollback on-demand if disperser respond with some type of rejection?
	relativePeriodRecord.Usage -= symbolUsage
	incrementRequired := big.NewInt(int64(a.PaymentCharged(numSymbols)))
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
func (a *Accountant) AccountBlob(ctx context.Context, numSymbols uint32, quorums []uint8, salt uint32) (*core.PaymentMetadata, error) {
	reservationPeriod, cumulativePayment, err := a.BlobPaymentInfo(ctx, numSymbols, quorums)
	if err != nil {
		return nil, err
	}

	pm := &core.PaymentMetadata{
		AccountID:         a.accountID,
		ReservationPeriod: reservationPeriod,
		CumulativePayment: cumulativePayment,
		Salt:              salt,
	}

	return pm, nil
}

// TODO: PaymentCharged and SymbolsCharged copied from meterer, should be refactored
// PaymentCharged returns the chargeable price for a given data length
func (a *Accountant) PaymentCharged(numSymbols uint32) uint64 {
	return uint64(a.SymbolsCharged(numSymbols)) * uint64(a.pricePerSymbol)
}

// SymbolsCharged returns the number of symbols charged for a given data length
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func (a *Accountant) SymbolsCharged(numSymbols uint32) uint32 {
	if numSymbols <= a.minNumSymbols {
		return a.minNumSymbols
	}
	// Round up to the nearest multiple of MinNumSymbols
	return uint32(core.RoundUpDivide(uint(numSymbols), uint(a.minNumSymbols))) * a.minNumSymbols
}

func (a *Accountant) GetRelativePeriodRecord(index uint32) *PeriodRecord {
	relativeIndex := index % a.numBins
	if a.periodRecords[relativeIndex].Index != uint32(index) {
		a.periodRecords[relativeIndex] = PeriodRecord{
			Index: uint32(index),
			Usage: 0,
		}
	}

	return &a.periodRecords[relativeIndex]
}

// SetPaymentState sets the accountant's state from the disperser's response
// We require disperser to return a valid set of global parameters, but optional
// account level on/off-chain state. If on-chain fields are not present, we use
// dummy values that disable accountant from using the corresponding payment method.
// If off-chain fields are not present, we assume the account has no payment history
// and set accoutant state to use initial values.
func (a *Accountant) SetPaymentState(paymentState *disperser_rpc.GetPaymentStateReply) error {
	if paymentState == nil {
		return fmt.Errorf("payment state cannot be nil")
	} else if paymentState.GetPaymentGlobalParams() == nil {
		return fmt.Errorf("payment global params cannot be nil")
	}

	a.minNumSymbols = uint32(paymentState.GetPaymentGlobalParams().GetMinNumSymbols())
	a.pricePerSymbol = uint32(paymentState.GetPaymentGlobalParams().GetPricePerSymbol())
	a.reservationWindow = uint32(paymentState.GetPaymentGlobalParams().GetReservationWindow())

	if paymentState.GetOnchainCumulativePayment() == nil {
		a.onDemand = &core.OnDemandPayment{
			CumulativePayment: big.NewInt(0),
		}
	} else {
		a.onDemand = &core.OnDemandPayment{
			CumulativePayment: new(big.Int).SetBytes(paymentState.GetOnchainCumulativePayment()),
		}
	}

	if paymentState.GetCumulativePayment() == nil {
		a.cumulativePayment = big.NewInt(0)
	} else {
		a.cumulativePayment = new(big.Int).SetBytes(paymentState.GetCumulativePayment())
	}

	if paymentState.GetReservation() == nil {
		a.reservation = &core.ReservedPayment{
			SymbolsPerSecond: 0,
			StartTimestamp:   0,
			EndTimestamp:     0,
			QuorumNumbers:    []uint8{},
			QuorumSplits:     []byte{},
		}
	} else {
		quorumNumbers := make([]uint8, len(paymentState.GetReservation().GetQuorumNumbers()))
		for i, quorum := range paymentState.GetReservation().GetQuorumNumbers() {
			quorumNumbers[i] = uint8(quorum)
		}
		quorumSplits := make([]uint8, len(paymentState.GetReservation().GetQuorumSplits()))
		for i, quorum := range paymentState.GetReservation().GetQuorumSplits() {
			quorumSplits[i] = uint8(quorum)
		}
		a.reservation = &core.ReservedPayment{
			SymbolsPerSecond: uint64(paymentState.GetReservation().GetSymbolsPerSecond()),
			StartTimestamp:   uint64(paymentState.GetReservation().GetStartTimestamp()),
			EndTimestamp:     uint64(paymentState.GetReservation().GetEndTimestamp()),
			QuorumNumbers:    quorumNumbers,
			QuorumSplits:     quorumSplits,
		}
	}

	periodRecords := make([]PeriodRecord, len(paymentState.GetPeriodRecords()))
	for i, record := range paymentState.GetPeriodRecords() {
		if record == nil {
			periodRecords[i] = PeriodRecord{Index: 0, Usage: 0}
		} else {
			periodRecords[i] = PeriodRecord{
				Index: record.Index,
				Usage: record.Usage,
			}
		}
	}
	a.periodRecords = periodRecords
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
