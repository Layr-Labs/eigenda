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
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

var requiredQuorums = []core.QuorumID{0, 1}

type Accountant struct {
	// on-chain states
	accountID   gethcommon.Address
	reservation map[core.QuorumID]*core.ReservedPayment
	onDemand    *core.OnDemandPayment

	// per-quorum payment configurations
	quorumPaymentConfigs  map[core.QuorumID]*core.PaymentQuorumConfig
	quorumProtocolConfigs map[core.QuorumID]*core.PaymentQuorumProtocolConfig

	// local accounting
	// contains 3 bins; circular wrapping of indices
	periodRecords     map[core.QuorumID][]PeriodRecord
	usageLock         sync.Mutex
	cumulativePayment *big.Int

	// number of bins in the circular accounting, restricted by minNumBins which is 3
	numBins uint32

	logger logging.Logger
}

type PeriodRecord struct {
	Index uint32
	Usage uint64
}

func NewAccountant(accountID gethcommon.Address, reservation map[uint8]*core.ReservedPayment, onDemand *core.OnDemandPayment, numBins uint32, logger logging.Logger) *Accountant {
	periodRecords := CreateEmptyReservationUsage(reservation, numBins)
	a := Accountant{
		accountID:             accountID,
		reservation:           reservation,
		onDemand:              onDemand,
		quorumPaymentConfigs:  make(map[core.QuorumID]*core.PaymentQuorumConfig),
		quorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
		periodRecords:         periodRecords,
		cumulativePayment:     big.NewInt(0),
		numBins:               max(numBins, uint32(meterer.MinNumBins)),
		logger:                logger,
	}
	// TODO: add a routine to refresh the on-chain state occasionally?
	return &a
}

// GetMinNumSymbols returns the minimum number of symbols for a given quorum
func (a *Accountant) GetMinNumSymbols(quorumID core.QuorumID) uint64 {
	if config, exists := a.quorumProtocolConfigs[quorumID]; exists {
		return config.MinNumSymbols
	}
	// Fallback to quorum 0 if the specific quorum config doesn't exist
	if config, exists := a.quorumProtocolConfigs[0]; exists {
		return config.MinNumSymbols
	}
	// Last resort fallback
	return 1
}

// GetPricePerSymbol returns the price per symbol for a given quorum
func (a *Accountant) GetPricePerSymbol(quorumID core.QuorumID) uint64 {
	if config, exists := a.quorumPaymentConfigs[quorumID]; exists {
		return config.OnDemandPricePerSymbol
	}
	// Fallback to quorum 0 if the specific quorum config doesn't exist
	if config, exists := a.quorumPaymentConfigs[0]; exists {
		return config.OnDemandPricePerSymbol
	}
	// Last resort fallback
	return 1
}

// GetReservationWindow returns the reservation window for a given quorum
func (a *Accountant) GetReservationWindow(quorumID core.QuorumID) uint64 {
	if config, exists := a.quorumProtocolConfigs[quorumID]; exists {
		return config.ReservationRateLimitWindow
	}
	// Fallback to quorum 0 if the specific quorum config doesn't exist
	if config, exists := a.quorumProtocolConfigs[0]; exists {
		return config.ReservationRateLimitWindow
	}
	// Last resort fallback
	return 1
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

	// first attempt to use the active reservation for each quorum
	// get all the quorum as part of a.reservation
	quorumWithReservation := make([]uint8, 0, len(a.reservation))
	for quorumNumber := range a.reservation {
		quorumWithReservation = append(quorumWithReservation, quorumNumber)
	}
	useReservation := true
	// check the input quorumNumbers is a subset of quorumWithReservation
	for _, quorumNumber := range quorumNumbers {
		if !slices.Contains(quorumWithReservation, quorumNumber) {
			useReservation = false
			a.logger.Warn("no reservation found for quorum", "quorum", quorumNumber)
		}
	}

	// check reservation feasibility before ondemand usage
	if useReservation {
		// Use the first quorum's reservation window for period calculation
		// (this should be consistent across quorums)
		firstQuorumID := core.QuorumID(quorumNumbers[0])
		reservationWindow := a.GetReservationWindow(firstQuorumID)
		currentReservationPeriod := meterer.GetReservationPeriodByNanosecond(timestamp, reservationWindow)
		symbolUsage := a.SymbolsCharged(numSymbols, firstQuorumID)
		a.usageLock.Lock()
		defer a.usageLock.Unlock()

		for quorumNumber, res := range a.reservation {
			relativePeriodRecord := a.GetRelativePeriodRecord(currentReservationPeriod, quorumNumber)
			relativePeriodRecord.Usage += symbolUsage

			quorumReservationWindow := a.GetReservationWindow(core.QuorumID(quorumNumber))
			binLimit := res.SymbolsPerSecond * uint64(quorumReservationWindow)
			if relativePeriodRecord.Usage <= binLimit {
				a.logger.Info("using reservation", "quorum", quorumNumber, "period", currentReservationPeriod, "usage", relativePeriodRecord.Usage, "binLimit", binLimit)
				continue
			}

			overflowPeriodRecord := a.GetRelativePeriodRecord(currentReservationPeriod+2, quorumNumber)
			// Allow one overflow when the overflow bin is empty, the current usage and new length are both less than the limit
			if overflowPeriodRecord.Usage == 0 && relativePeriodRecord.Usage-symbolUsage < binLimit && symbolUsage <= binLimit {
				overflowPeriodRecord.Usage += relativePeriodRecord.Usage - binLimit
				relativePeriodRecord.Usage = binLimit
				a.logger.Info("reservation bin overflowed, using overflow bin", "quorum", quorumNumber, "overflowPeriod", currentReservationPeriod+2, "overflowUsage", overflowPeriodRecord.Usage)
				continue
			}

			// Rollback usage for this quorum since we couldn't use it
			useReservation = false
			relativePeriodRecord.Usage -= symbolUsage
			a.logger.Warn("reservation bin full, rolling back usage", "quorum", quorumNumber, "period", currentReservationPeriod)
		}
	}

	if useReservation {
		a.logger.Info("reservation payment successfully generated for requested quorums", "quorums", quorumNumbers, "symbols", numSymbols)
		return big.NewInt(0), nil
	}

	// reservation not available for any quorums, attempt on-demand
	// on-demand can be applied to required quorums only, but on-chain record is only on quorum 0
	//todo: rollback on-demand if disperser respond with some ratelimit rejection
	incrementRequired := big.NewInt(int64(a.PaymentCharged(numSymbols, core.QuorumID(0))))
	resultingPayment := big.NewInt(0)
	resultingPayment.Add(a.cumulativePayment, incrementRequired)
	if resultingPayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		if err := QuorumCheck(quorumNumbers, requiredQuorums); err != nil {
			a.logger.Error("quorum check failed for on-demand payment", "err", err)
			return big.NewInt(0), err
		}
		a.logger.Info("using on-demand payment", "increment", incrementRequired, "cumulative", a.cumulativePayment)
		a.cumulativePayment.Add(a.cumulativePayment, incrementRequired)
		return a.cumulativePayment, nil
	}
	a.logger.Error("no bandwidth reservation and insufficient on-demand payment", "account", a.accountID.Hex(), "required", incrementRequired, "cumulative", a.cumulativePayment, "onDemand", a.onDemand.CumulativePayment)
	return big.NewInt(0), fmt.Errorf(
		"no bandwidth reservation found for account %s, and current cumulativePayment balance insufficient "+
			"to make an on-demand dispersal. Consider depositing more eth to the PaymentVault contract.", a.accountID.Hex())
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
// PaymentCharged returns the chargeable price for a given data length for a specific quorum
func (a *Accountant) PaymentCharged(numSymbols uint64, quorumID core.QuorumID) uint64 {
	return a.SymbolsCharged(numSymbols, quorumID) * a.GetPricePerSymbol(quorumID)
}

// SymbolsCharged returns the number of symbols charged for a given data length for a specific quorum
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func (a *Accountant) SymbolsCharged(numSymbols uint64, quorumID core.QuorumID) uint64 {
	minNumSymbols := a.GetMinNumSymbols(quorumID)
	if numSymbols <= minNumSymbols {
		return minNumSymbols
	}
	// Round up to the nearest multiple of MinNumSymbols
	return core.RoundUpDivide(numSymbols, minNumSymbols) * minNumSymbols
}

func (a *Accountant) GetRelativePeriodRecord(index uint64, quorumNumber uint8) *PeriodRecord {
	relativeIndex := uint32(index % uint64(a.numBins))
	if a.periodRecords[quorumNumber][relativeIndex].Index != uint32(index) {
		a.periodRecords[quorumNumber][relativeIndex] = PeriodRecord{
			Index: uint32(index),
			Usage: 0,
		}
	}

	return &a.periodRecords[quorumNumber][relativeIndex]
}

// SetPaymentState sets the accountant's state from the disperser's response
// We require disperser to return a valid set of global parameters, but optional
// account level on/off-chain state. If on-chain fields are not present, we use
// dummy values that disable accountant from using the corresponding payment method.
// If off-chain fields are not present, we assume the account has no payment history
// and set accoutant state to use initial values.
func (a *Accountant) SetPaymentState(paymentState *disperser_rpc.GetQuorumSpecificPaymentStateReply) error {
	if paymentState == nil {
		a.logger.Error("payment state cannot be nil")
		return fmt.Errorf("payment state cannot be nil")
	} else if paymentState.GetPaymentVaultParams() == nil {
		a.logger.Error("payment vault params cannot be nil")
		return fmt.Errorf("payment vault params cannot be nil")
	}

	vaultParams := paymentState.GetPaymentVaultParams()
	
	if vaultParams.GetQuorumPaymentConfigs() == nil {
		a.logger.Error("payment quorum configs cannot be nil")
		return fmt.Errorf("payment quorum configs cannot be nil")
	}
	
	if vaultParams.GetQuorumProtocolConfigs() == nil {
		a.logger.Error("payment quorum protocol configs cannot be nil")
		return fmt.Errorf("payment quorum protocol configs cannot be nil")
	}

	// Initialize the per-quorum configuration maps
	a.quorumPaymentConfigs = make(map[core.QuorumID]*core.PaymentQuorumConfig)
	a.quorumProtocolConfigs = make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig)

	// Convert protobuf configs to core types
	for quorumID, pbPaymentConfig := range vaultParams.GetQuorumPaymentConfigs() {
		a.quorumPaymentConfigs[core.QuorumID(quorumID)] = &core.PaymentQuorumConfig{
			ReservationSymbolsPerSecond: pbPaymentConfig.GetReservationSymbolsPerSecond(),
			OnDemandSymbolsPerSecond:    pbPaymentConfig.GetOnDemandSymbolsPerSecond(),
			OnDemandPricePerSymbol:      pbPaymentConfig.GetOnDemandPricePerSymbol(),
		}
	}

	for quorumID, pbProtocolConfig := range vaultParams.GetQuorumProtocolConfigs() {
		a.quorumProtocolConfigs[core.QuorumID(quorumID)] = &core.PaymentQuorumProtocolConfig{
			MinNumSymbols:              pbProtocolConfig.GetMinNumSymbols(),
			ReservationAdvanceWindow:   pbProtocolConfig.GetReservationAdvanceWindow(),
			ReservationRateLimitWindow: pbProtocolConfig.GetReservationRateLimitWindow(),
			OnDemandRateLimitWindow:    pbProtocolConfig.GetOnDemandRateLimitWindow(),
			OnDemandEnabled:            pbProtocolConfig.GetOnDemandEnabled(),
		}
	}

	a.logger.Info("updated payment state with per-quorum configurations",
		"numPaymentConfigs", len(a.quorumPaymentConfigs),
		"numProtocolConfigs", len(a.quorumProtocolConfigs))

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

	if paymentState.GetReservations() == nil {
		a.reservation = make(map[core.QuorumID]*core.ReservedPayment)
	} else {
		a.reservation = make(map[core.QuorumID]*core.ReservedPayment)
		for _, reservation := range paymentState.GetReservations() {
			a.reservation[core.QuorumID(reservation.QuorumNumber)] = &core.ReservedPayment{
				SymbolsPerSecond: uint64(reservation.GetSymbolsPerSecond()),
				StartTimestamp:   uint64(reservation.GetStartTimestamp()),
				EndTimestamp:     uint64(reservation.GetEndTimestamp()),
				QuorumNumbers:    []core.QuorumID{},
				QuorumSplits:     []byte{},
			}
		}
	}

	// periodRecords should be a map of quorumNumbers (the quorum numbers same as reservations)
	// and the value should be a slice of PeriodRecord, which is a circular array of length numBins

	periodRecords := make(map[core.QuorumID][]PeriodRecord)
	for quorumNumber, _ := range a.reservation {
		periodRecords[quorumNumber] = make([]PeriodRecord, a.numBins)
		for i := uint32(0); i < a.numBins; i++ {
			periodRecords[quorumNumber][i] = PeriodRecord{
				Index: i,
				Usage: 0,
			}
		}
	}

	for _, record := range paymentState.GetPeriodRecords() {
		quorumNumber := uint8(record.QuorumNumber)
		if _, exists := periodRecords[quorumNumber]; !exists {
			periodRecords[quorumNumber] = make([]PeriodRecord, a.numBins)
			for i := uint32(0); i < a.numBins; i++ {
				periodRecords[quorumNumber][i] = PeriodRecord{
					Index: i,
					Usage: 0,
				}
			}
		}
		idx := record.Index % a.numBins
		periodRecords[quorumNumber][idx] = PeriodRecord{
			Index: record.Index,
			Usage: record.Usage,
		}
	}
	a.periodRecords = periodRecords

	a.logger.Info("payment state updated", "reservations", a.reservation, "periodRecords", a.periodRecords, "onchain cumulative deposit", a.onDemand, "used amount", a.cumulativePayment)

	return nil
}

// CreateEmptyReservationUsage creates empty reservation usage records for the provided quorum numbers
func CreateEmptyReservationUsage(quorumNumbers map[core.QuorumID]*core.ReservedPayment, numBins uint32) map[core.QuorumID][]PeriodRecord {
	reservationUsage := make(map[core.QuorumID][]PeriodRecord)
	for quorumNumber := range quorumNumbers {
		reservationUsage[quorumNumber] = make([]PeriodRecord, numBins)
		for i := range reservationUsage[quorumNumber] {
			reservationUsage[quorumNumber][i] = PeriodRecord{Index: uint32(i), Usage: 0}
		}
	}
	return reservationUsage
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
