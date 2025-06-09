package clients

import (
	"context"
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

	symbolUsage := meterer.SymbolsCharged(numSymbols, a.minNumSymbols)

	// Always try to use reservation first
	payment, err := a.ReservationUsage(symbolUsage, quorumNumbers, timestamp)
	if err == nil {
		return payment, nil
	}

	// Fall back to on-demand payment if reservation fails
	return a.OnDemandUsage(symbolUsage, quorumNumbers)
}

// ReservationUsage attempts to use the reservation for the given request.
// Returns (0, nil) if successful, or (nil, error) if reservation cannot be used.
func (a *Accountant) ReservationUsage(
	symbolUsage uint64,
	quorumNumbers []uint8,
	timestamp int64) (*big.Int, error) {

	currentReservationPeriod := meterer.GetReservationPeriodByNanosecond(timestamp, a.reservationWindow)

	a.usageLock.Lock()
	defer a.usageLock.Unlock()

	relativePeriodRecord := a.GetRelativePeriodRecord(currentReservationPeriod)
	relativePeriodRecord.Usage += symbolUsage

	// Check if we can use the reservation within the bin limit
	binLimit := a.reservation.SymbolsPerSecond * uint64(a.reservationWindow)
	if relativePeriodRecord.Usage <= binLimit {
		if err := meterer.ValidateQuorum(quorumNumbers, a.reservation.QuorumNumbers); err != nil {
			return nil, err
		}
		return big.NewInt(0), nil
	}

	// Try to use overflow bin if available
	overflowPeriodRecord := a.GetRelativePeriodRecord(currentReservationPeriod + 2)
	if overflowPeriodRecord.Usage == 0 && relativePeriodRecord.Usage-symbolUsage < binLimit && symbolUsage <= binLimit {
		if err := meterer.ValidateQuorum(quorumNumbers, a.reservation.QuorumNumbers); err != nil {
			return nil, err
		}
		overflowPeriodRecord.Usage += relativePeriodRecord.Usage - binLimit
		return big.NewInt(0), nil
	}

	// Reservation not sufficient for the request, rollback the usage
	relativePeriodRecord.Usage -= symbolUsage
	return nil, fmt.Errorf("insufficient reservation")
}

// OnDemandUsage attempts to use on-demand payment for the given request.
// Returns the cumulative payment if successful, or an error if on-demand cannot be used.
func (a *Accountant) OnDemandUsage(
	symbolUsage uint64,
	quorumNumbers []uint8) (*big.Int, error) {

	incrementRequired := meterer.PaymentCharged(symbolUsage, a.pricePerSymbol)
	resultingPayment := big.NewInt(0)
	resultingPayment.Add(a.cumulativePayment, incrementRequired)

	if resultingPayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		if err := meterer.ValidateQuorum(quorumNumbers, requiredQuorums); err != nil {
			return nil, err
		}
		a.cumulativePayment.Add(a.cumulativePayment, incrementRequired)
		return a.cumulativePayment, nil
	}

	return nil, fmt.Errorf("insufficient ondemand payment")
}

// AccountBlob accountant provides and records payment information
func (a *Accountant) AccountBlob(
	ctx context.Context,
	timestamp int64,
	numSymbols uint64,
	quorums []uint8) (*core.PaymentMetadata, error) {

	cumulativePayment, err := a.BlobPaymentInfo(ctx, numSymbols, quorums, timestamp)
	if err != nil {
		return nil, fmt.Errorf("cannot create payment infomation for reservation or on-demand. Consider depositing more eth to the PaymentVault contract for your account. For more details, see https://docs.eigenda.xyz/core-concepts/payments#disperser-client-requirements. Account: %s, Error: %s", a.accountID.Hex(), err.Error())
	}

	pm := &core.PaymentMetadata{
		AccountID:         a.accountID,
		Timestamp:         timestamp,
		CumulativePayment: cumulativePayment,
	}

	return pm, nil
}

func (a *Accountant) GetRelativePeriodRecord(index uint64) *PeriodRecord {
	relativeIndex := uint32(index % uint64(a.numBins))
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
