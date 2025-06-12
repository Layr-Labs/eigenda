package clients

import (
	"errors"
	"fmt"
	"math/big"
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
	// contains an array of period records, with length of max(MinNumBins, numBins)
	// numBins can be arbitrarily bigger than MinNumBins if the client wants to track more history in the cache
	periodRecords     []PeriodRecord
	cumulativePayment *big.Int

	// locks for concurrent access to period records and on-demand payment
	periodRecordsLock sync.Mutex
	onDemandLock      sync.Mutex
}

// PeriodRecord contains the index of the reservation period and the usage of the period
type PeriodRecord struct {
	// Index is start timestamp of the period in seconds; it is always a multiple of the reservation window
	Index uint32
	// Usage is the usage of the period in symbols
	Usage uint64
}

func NewAccountant(accountID gethcommon.Address, reservation *core.ReservedPayment, onDemand *core.OnDemandPayment, reservationWindow uint64, pricePerSymbol uint64, minNumSymbols uint64, numBins uint32) *Accountant {
	periodRecords := make([]PeriodRecord, max(numBins, uint32(meterer.MinNumBins)))
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
	}
	// TODO: add a routine to refresh the on-chain state occasionally?
	return &a
}

// reservationUsage attempts to use the reservation for the given request.
func (a *Accountant) reservationUsage(
	symbolUsage uint64,
	quorumNumbers []uint8,
	timestamp int64) error {
	if err := meterer.ValidateQuorum(quorumNumbers, a.reservation.QuorumNumbers); err != nil {
		return err
	}
	if !a.reservation.IsActiveByNanosecond(timestamp) {
		return fmt.Errorf("reservation is not active at timestamp %d", timestamp)
	}

	reservationWindow := a.reservationWindow
	currentReservationPeriod := meterer.GetReservationPeriodByNanosecond(timestamp, reservationWindow)

	a.periodRecordsLock.Lock()
	defer a.periodRecordsLock.Unlock()
	relativePeriodRecord := a.getOrRefreshRelativePeriodRecord(currentReservationPeriod, reservationWindow)
	relativePeriodRecord.Usage += symbolUsage

	// Check if we can use the reservation within the bin limit
	binLimit := meterer.GetReservationBinLimit(a.reservation, a.reservationWindow)
	if relativePeriodRecord.Usage <= binLimit {
		return nil
	}

	overflowPeriodRecord := a.getOrRefreshRelativePeriodRecord(meterer.GetOverflowPeriod(currentReservationPeriod, reservationWindow), reservationWindow)
	// Allow one overflow when the overflow bin is empty, the current usage and new length are both less than the limit
	if overflowPeriodRecord.Usage == 0 && relativePeriodRecord.Usage-symbolUsage < binLimit && symbolUsage <= binLimit {
		overflowPeriodRecord.Usage += relativePeriodRecord.Usage - binLimit
		return nil
	}

	// Reservation not sufficient for the request, rollback the usage
	relativePeriodRecord.Usage -= symbolUsage
	return errors.New("insufficient reservation")
}

// onDemandUsage attempts to use on-demand payment for the given request.
// Returns the cumulative payment if successful, or an error if on-demand cannot be used.
func (a *Accountant) onDemandUsage(symbolUsage uint64, quorumNumbers []uint8) (*big.Int, error) {
	if err := meterer.ValidateQuorum(quorumNumbers, requiredQuorums); err != nil {
		return nil, err
	}

	a.onDemandLock.Lock()
	defer a.onDemandLock.Unlock()

	incrementRequired := meterer.PaymentCharged(symbolUsage, a.pricePerSymbol)
	resultingPayment := new(big.Int).Add(a.cumulativePayment, incrementRequired)

	if resultingPayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		a.cumulativePayment.Add(a.cumulativePayment, incrementRequired)
		return a.cumulativePayment, nil
	}

	return nil, errors.New("insufficient ondemand payment")
}

// AccountBlob accountant generates payment information for a request. The accountant
// takes in a timestamp at the current UNIX time in nanoseconds, number of symbols of the request,
// and the quorums to disperse the request to. It will attempt to use the active reservation first
// and then on-demand if the reservation is not available or insufficient for the request.
// It returns a payment metadata object that will be used to create the payment header and signature,
// as specified in api/proto/common/v2/common_v2.proto
func (a *Accountant) AccountBlob(
	timestamp int64,
	numSymbols uint64,
	quorums []uint8) (*core.PaymentMetadata, error) {

	symbolUsage := meterer.SymbolsCharged(numSymbols, a.minNumSymbols)

	// Always try to use reservation first
	err := a.reservationUsage(symbolUsage, quorums, timestamp)
	if err == nil {
		return &core.PaymentMetadata{
			AccountID:         a.accountID,
			Timestamp:         timestamp,
			CumulativePayment: big.NewInt(0),
		}, nil
	}

	// Fall back to on-demand payment if reservation fails
	cumulativePayment, err := a.onDemandUsage(symbolUsage, quorums)
	if err != nil {
		return nil, fmt.Errorf("cannot create payment information for reservation or on-demand. Consider depositing more eth to the PaymentVault contract for your account. For more details, see https://docs.eigenda.xyz/core-concepts/payments#disperser-client-requirements. Account: %s, Error: %w", a.accountID.Hex(), err)
	}

	pm := &core.PaymentMetadata{
		AccountID:         a.accountID,
		Timestamp:         timestamp,
		CumulativePayment: cumulativePayment,
	}

	return pm, nil
}

// getOrRefreshRelativePeriodRecord returns the period record for the given index (which is in seconds and is the multiple of the reservation window),
// wrapping around the circular buffer and clearing the record if the index is greater than the number of bins
func (a *Accountant) getOrRefreshRelativePeriodRecord(index uint64, reservationWindow uint64) *PeriodRecord {
	relativeIndex := uint32((index / reservationWindow) % uint64(len(a.periodRecords)))
	if relativeIndex >= uint32(len(a.periodRecords)) {
		panic(fmt.Sprintf("relativeIndex %d is greater than the number of bins %d cached", relativeIndex, len(a.periodRecords)))
	}
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
// and set accountant state to use initial values.
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
