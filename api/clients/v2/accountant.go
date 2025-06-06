package clients

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"sync"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

var requiredQuorums = []core.QuorumID{0, 1}

type Accountant struct {
	// on-chain states
	accountID         gethcommon.Address
	reservation       *core.ReservedPayment
	onDemand          *core.OnDemandPayment
	reservationWindow uint64
	pricePerSymbol    uint64
	minNumSymbols     uint64

	// local accounting
	// contains numBins of period records; circular wrapping on period record indices
	periodRecords     []PeriodRecord
	usageLock         sync.Mutex
	cumulativePayment *big.Int

	// number of bins in the circular accounting, restricted by minNumBins which is 3
	numBins uint32
}

// PeriodRecord contains the index of the reservation period and the usage of the period
type PeriodRecord struct {
	// Index is start timestamp of the period in seconds; it is always a multiple of the reservation window
	Index uint32
	// Usage is the usage of the period in symbols
	Usage uint64
}

func NewAccountant(accountID gethcommon.Address, reservation *core.ReservedPayment, onDemand *core.OnDemandPayment, reservationWindow uint64, pricePerSymbol uint64, minNumSymbols uint64, numBins uint32) *Accountant {
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
// then on-demand if the reservation is not available. It takes in a timestamp at
// the current UNIX time in nanoseconds, and returns a cumulative payment for on-
// demand payments in units of wei. Both timestamp and cumulative payment are used
// to create the payment header and signature, with non-zero cumulative payment
// indicating on-demand payment.
// These generated values are used to create the payment header and signature, as specified in
// api/proto/common/v2/common_v2.proto
func (a *Accountant) BlobPaymentInfo(
	ctx context.Context,
	numSymbols uint64,
	quorumNumbers []uint8,
	timestamp int64) (*big.Int, error) {
	reservationWindow := a.reservationWindow
	currentReservationPeriod := meterer.GetReservationPeriodByNanosecond(timestamp, reservationWindow)
	symbolUsage := a.SymbolsCharged(numSymbols)

	a.usageLock.Lock()
	defer a.usageLock.Unlock()
	relativePeriodRecord := a.GetOrRefreshRelativePeriodRecord(currentReservationPeriod, reservationWindow)
	relativePeriodRecord.Usage += symbolUsage

	// first attempt to use the active reservation
	binLimit := a.reservation.SymbolsPerSecond * uint64(a.reservationWindow)
	if relativePeriodRecord.Usage <= binLimit {
		if err := QuorumCheck(quorumNumbers, a.reservation.QuorumNumbers); err != nil {
			return big.NewInt(0), err
		}
		return big.NewInt(0), nil
	}

	overflowPeriodRecord := a.GetOrRefreshRelativePeriodRecord(currentReservationPeriod+2*reservationWindow, reservationWindow)
	// Allow one overflow when the overflow bin is empty, the current usage and new length are both less than the limit
	if overflowPeriodRecord.Usage == 0 && relativePeriodRecord.Usage-symbolUsage < binLimit && symbolUsage <= binLimit {
		if err := QuorumCheck(quorumNumbers, a.reservation.QuorumNumbers); err != nil {
			return big.NewInt(0), err
		}
		overflowPeriodRecord.Usage += relativePeriodRecord.Usage - binLimit
		return big.NewInt(0), nil
	}

	// reservation not available, rollback reservation records, attempt on-demand
	//todo: rollback on-demand if disperser respond with some type of rejection?
	relativePeriodRecord.Usage -= symbolUsage
	incrementRequired := big.NewInt(int64(a.PaymentCharged(numSymbols)))

	resultingPayment := big.NewInt(0)
	resultingPayment.Add(a.cumulativePayment, incrementRequired)
	if resultingPayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		if err := QuorumCheck(quorumNumbers, requiredQuorums); err != nil {
			return big.NewInt(0), err
		}
		a.cumulativePayment.Add(a.cumulativePayment, incrementRequired)
		return a.cumulativePayment, nil
	}
	return big.NewInt(0), fmt.Errorf(
		"invalid payments: no available bandwidth reservation found for account %s, and current cumulativePayment balance insufficient "+
			"to make an on-demand dispersal. Consider increasing reservation or cumulative payment on-chain. For more details, see https://docs.eigenda.xyz/core-concepts/payments#disperser-client-requirements", a.accountID.Hex())
}

// AccountBlob accountant provides and records payment information
func (a *Accountant) AccountBlob(
	ctx context.Context,
	timestamp int64,
	numSymbols uint64,
	quorums []uint8) (*core.PaymentMetadata, error) {

	cumulativePayment, err := a.BlobPaymentInfo(ctx, numSymbols, quorums, timestamp)
	if err != nil {
		return nil, err
	}

	pm := &core.PaymentMetadata{
		AccountID:         a.accountID,
		Timestamp:         timestamp,
		CumulativePayment: cumulativePayment,
	}

	return pm, nil
}

// TODO: PaymentCharged and SymbolsCharged copied from meterer, should be refactored
// PaymentCharged returns the chargeable price for a given data length
func (a *Accountant) PaymentCharged(numSymbols uint64) uint64 {
	return a.SymbolsCharged(numSymbols) * a.pricePerSymbol
}

// SymbolsCharged returns the number of symbols charged for a given data length
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func (a *Accountant) SymbolsCharged(numSymbols uint64) uint64 {
	if numSymbols <= a.minNumSymbols {
		return a.minNumSymbols
	}
	// Round up to the nearest multiple of MinNumSymbols
	return core.RoundUpDivide(numSymbols, a.minNumSymbols) * a.minNumSymbols
}

// GetRelativePeriodRecord returns the period record for the given index
// Returns an empty record if there is no record for the relative index
func (a *Accountant) GetRelativePeriodRecord(index uint64) *PeriodRecord {
	relativeIndex := uint32((index / a.reservationWindow) % uint64(a.numBins))
	// Return empty record if the index is greater than the number of bins (should never happen by accountant initialization)
	if relativeIndex >= uint32(a.numBins) {
		panic(fmt.Sprintf("relativeIndex %d is greater than the number of bins %d cached", relativeIndex, a.numBins))
	}
	return &a.periodRecords[relativeIndex]
}

// GetOrRefreshRelativePeriodRecord returns the period record for the given index (which is in seconds and is the multiple of the reservation window),
// wrapping around the circular buffer and clearing the record if the index is greater than the number of bins
func (a *Accountant) GetOrRefreshRelativePeriodRecord(index uint64, reservationWindow uint64) *PeriodRecord {
	relativeIndex := uint32((index / reservationWindow) % uint64(a.numBins))
	if a.periodRecords[relativeIndex].Index < uint32(index) {
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

	a.minNumSymbols = paymentState.GetPaymentGlobalParams().GetMinNumSymbols()
	a.pricePerSymbol = paymentState.GetPaymentGlobalParams().GetPricePerSymbol()
	a.reservationWindow = paymentState.GetPaymentGlobalParams().GetReservationWindow()

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
