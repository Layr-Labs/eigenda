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
	accountID    gethcommon.Address
	reservations map[core.QuorumID]*core.QuorumReservation
	onDemand     *core.OnDemandPayment

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
}

type PeriodRecord struct {
	Index uint32
	Usage uint64
}

// NewAccountant initializes an accountant with the given account ID, reservations, on-demand payment, and number of bins
// TODO: Consider making this initialization take all the fields as arguments or entirely empty as clients are typically
// syncing the onchain configurations and offchain usage with the disperser server
func NewAccountant(accountID gethcommon.Address, reservations map[uint8]*core.QuorumReservation, onDemand *core.OnDemandPayment, numBins uint32) *Accountant {
	periodRecords := CreateEmptyReservationUsage(reservations, numBins)
	a := Accountant{
		accountID:             accountID,
		reservations:          reservations,
		onDemand:              onDemand,
		quorumPaymentConfigs:  make(map[core.QuorumID]*core.PaymentQuorumConfig),
		quorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
		periodRecords:         periodRecords,
		cumulativePayment:     big.NewInt(0),
		numBins:               max(numBins, uint32(meterer.MinNumBins)),
	}
	// TODO: add a routine to refresh the on-chain state occasionally?
	return &a
}

// validateQuorumReservations checks if all requested quorums have valid reservations
func (a *Accountant) validateQuorumReservations(quorumNumbers []uint8, timestamp int64) error {
	if len(quorumNumbers) == 0 {
		return fmt.Errorf("no quorum numbers provided")
	}

	for _, quorumNumber := range quorumNumbers {
		// check if the quorum number is in the reservations
		if reservation, exists := a.reservations[core.QuorumID(quorumNumber)]; !exists {
			return fmt.Errorf("No reservation found on quorum %d", quorumNumber)
		} else {
			// check if the reservation is active
			if !meterer.IsWithinTimeRange(uint64(reservation.StartTimestamp), uint64(reservation.EndTimestamp), timestamp) {
				return fmt.Errorf("reservation not active")
			}
		}
	}
	return nil
}

// calculateReservationUsage calculates the symbol usage for a given number of symbols and quorum
func (a *Accountant) calculateReservationUsage(numSymbols uint64, quorumNumber uint8) (uint64, error) {
	_, err := a.GetReservationWindow(quorumNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to get reservation window for quorum %d: %w", quorumNumber, err)
	}
	symbolsCharged, err := a.SymbolsCharged(numSymbols, quorumNumber)
	if err != nil {
		return 0, err
	}
	return symbolsCharged, nil
}

// updateReservationUsage updates the usage records for a quorum's reservation
func (a *Accountant) updateReservationUsage(quorumNumber uint8, currentPeriod uint64, symbolUsage uint64) error {
	res, exists := a.reservations[quorumNumber]
	if !exists {
		return fmt.Errorf("reservation not found for quorum %d", quorumNumber)
	}

	return a.processQuorumReservation(quorumNumber, res, currentPeriod, symbolUsage)
}

// BlobPaymentInfo calculates and records payment information. The accountant
// will attempt to use the active reservation first and check for quorum settings,
// then on-demand if the reservation is not available. It takes in a timestamp at
// the current UNIX time in nanoseconds, and returns a cumulative payment for on-
// demand payments in units of wei. Both timestamp and cumulative payment are used
// to create the payment header and signature, with non-zero cumulative payment
// indicating on-demand payment.
func (a *Accountant) BlobPaymentInfo(
	ctx context.Context,
	numSymbols uint64,
	quorumNumbers []uint8,
	timestamp int64) (*big.Int, error) {

	// Always try to use reservation first
	payment, err := a.ReservationUsage(numSymbols, quorumNumbers, timestamp)
	if err == nil {
		return payment, nil
	}

	// Fall back to on-demand payment if reservation fails
	return a.OnDemandUsage(numSymbols, quorumNumbers)
}

// ReservationUsage attempts to use the reservation for the requested quorums
// Returns (0, nil) if successful, or (nil, error) if reservation cannot be used
func (a *Accountant) ReservationUsage(numSymbols uint64, quorumNumbers []uint8, timestamp int64) (*big.Int, error) {
	// Validate quorum reservations
	if err := a.validateQuorumReservations(quorumNumbers, timestamp); err != nil {
		return nil, err
	}

	// Lock for updating usage records
	a.usageLock.Lock()
	defer a.usageLock.Unlock()

	// Try to use reservation for each quorum
	for _, quorumNumber := range quorumNumbers {
		// Calculate current period and symbol usage
		reservationWindow, err := a.GetReservationWindow(quorumNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get reservation window for quorum %d: %w", quorumNumber, err)
		}

		currentReservationPeriod := meterer.GetReservationPeriodByNanosecond(timestamp, reservationWindow)
		symbolUsage, err := a.calculateReservationUsage(numSymbols, quorumNumber)
		if err != nil {
			return nil, err
		}

		if err := a.updateReservationUsage(quorumNumber, currentReservationPeriod, symbolUsage); err != nil {
			// Rollback usage for this quorum
			relativePeriodRecord := a.GetRelativePeriodRecord(currentReservationPeriod, quorumNumber)
			relativePeriodRecord.Usage -= symbolUsage
			return nil, err
		}
	}

	return big.NewInt(0), nil
}

// processQuorumReservation handles the reservation usage for a single quorum
func (a *Accountant) processQuorumReservation(
	quorumNumber uint8,
	res *core.QuorumReservation,
	currentReservationPeriod uint64,
	symbolUsage uint64) error {

	relativePeriodRecord := a.GetRelativePeriodRecord(currentReservationPeriod, quorumNumber)
	relativePeriodRecord.Usage += symbolUsage

	quorumReservationWindow, err := a.GetReservationWindow(core.QuorumID(quorumNumber))
	if err != nil {
		return fmt.Errorf("failed to get reservation window for quorum %d: %w", quorumNumber, err)
	}

	binLimit := res.SymbolsPerSecond * uint64(quorumReservationWindow)

	// Check if we're within the bin limit
	if relativePeriodRecord.Usage <= binLimit {
		return nil
	}

	// Try to use overflow bin if we're over the limit
	overflowPeriodRecord := a.GetRelativePeriodRecord(currentReservationPeriod+2, quorumNumber)
	canUseOverflow := overflowPeriodRecord.Usage == 0 &&
		relativePeriodRecord.Usage-symbolUsage < binLimit &&
		symbolUsage <= binLimit

	if canUseOverflow {
		overflowPeriodRecord.Usage += relativePeriodRecord.Usage - binLimit
		relativePeriodRecord.Usage = binLimit
		return nil
	}

	return fmt.Errorf("reservation limit exceeded for quorum %d", quorumNumber)
}

// OnDemandUsage handles the on-demand payment calculation and validation
func (a *Accountant) OnDemandUsage(numSymbols uint64, quorumNumbers []uint8) (*big.Int, error) {
	// Verify quorum requirements
	if err := QuorumCheck(quorumNumbers, requiredQuorums); err != nil {
		return nil, err
	}

	// Calculate payment needed for the number of symbols
	paymentCharged, err := a.PaymentCharged(numSymbols, meterer.OnDemandQuorumID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate payment charged: %w", err)
	}

	// Calculate the increment required to add to the cumulative payment
	incrementRequired := new(big.Int).SetUint64(paymentCharged)
	resultingPayment := big.NewInt(0)
	resultingPayment.Add(a.cumulativePayment, incrementRequired)

	// Check if we have sufficient balance
	if resultingPayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		a.cumulativePayment.Add(a.cumulativePayment, incrementRequired)
		return a.cumulativePayment, nil
	}

	return nil, fmt.Errorf(
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

// PaymentCharged returns the chargeable price for a given data length for a specific quorum
func (a *Accountant) PaymentCharged(numSymbols uint64, quorumID core.QuorumID) (uint64, error) {
	pricePerSymbol, err := a.GetPricePerSymbol(quorumID)
	if err != nil {
		return 0, err
	}
	symbolsCharged, err := a.SymbolsCharged(numSymbols, quorumID)
	if err != nil {
		return 0, err
	}
	return symbolsCharged * pricePerSymbol, nil
}

// SymbolsCharged returns the number of symbols charged for a given data length for a specific quorum
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func (a *Accountant) SymbolsCharged(numSymbols uint64, quorumID core.QuorumID) (uint64, error) {
	minNumSymbols, err := a.GetMinNumSymbols(quorumID)
	if err != nil {
		return 0, err
	}
	if numSymbols <= minNumSymbols {
		return minNumSymbols, nil
	}
	// Round up to the nearest multiple of MinNumSymbols
	return core.RoundUpDivide(numSymbols, minNumSymbols) * minNumSymbols, nil
}

// GetMinNumSymbols returns the minimum number of symbols for a given quorum
func (a *Accountant) GetMinNumSymbols(quorumID core.QuorumID) (uint64, error) {
	if config, exists := a.quorumProtocolConfigs[quorumID]; exists {
		return config.MinNumSymbols, nil
	}
	return 0, fmt.Errorf("quorum ID %d not found in protocol configs", quorumID)
}

// GetPricePerSymbol returns the price per symbol for a given quorum
func (a *Accountant) GetPricePerSymbol(quorumID core.QuorumID) (uint64, error) {
	if config, exists := a.quorumPaymentConfigs[quorumID]; exists {
		return config.OnDemandPricePerSymbol, nil
	}
	return 0, fmt.Errorf("quorum ID %d not found in payment configs", quorumID)
}

// GetReservationWindow returns the reservation window for a given quorum
func (a *Accountant) GetReservationWindow(quorumID core.QuorumID) (uint64, error) {
	if config, exists := a.quorumProtocolConfigs[quorumID]; exists {
		return config.ReservationRateLimitWindow, nil
	}
	return 0, fmt.Errorf("quorum ID %d not found in protocol configs", quorumID)
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
func (a *Accountant) SetPaymentState(paymentState *disperser_rpc.GetPaymentStateForAllQuorumsReply) error {
	if paymentState == nil {
		return fmt.Errorf("payment state cannot be nil")
	} else if paymentState.GetPaymentVaultParams() == nil {
		return fmt.Errorf("payment vault params cannot be nil")
	}

	vaultParams := paymentState.GetPaymentVaultParams()

	if vaultParams.GetQuorumPaymentConfigs() == nil {
		return fmt.Errorf("payment quorum configs cannot be nil")
	}

	if vaultParams.GetQuorumProtocolConfigs() == nil {
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
		a.reservations = make(map[core.QuorumID]*core.QuorumReservation)
	} else {
		a.reservations = make(map[core.QuorumID]*core.QuorumReservation)
		for quorumNumber, reservation := range paymentState.GetReservations() {
			a.reservations[core.QuorumID(quorumNumber)] = reservation
		}
		a.periodRecords = CreateEmptyReservationUsage(a.reservations, a.numBins)
		for quorumNumber, periodRecords := range paymentState.GetPeriodRecords() {
			if periodRecords != nil {
				for _, record := range periodRecords.GetRecords() {
					idx := record.Index % a.numBins
					a.periodRecords[core.QuorumID(quorumNumber)][idx] = PeriodRecord{
						Index: record.Index,
						Usage: record.Usage,
					}
				}
			}
		}
	}

	return nil
}

// CreateEmptyReservationUsage creates empty reservation usage records for the provided quorum numbers
func CreateEmptyReservationUsage(quorumNumbers map[core.QuorumID]*core.QuorumReservation, numBins uint32) map[core.QuorumID][]PeriodRecord {
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
			return fmt.Errorf("provided quorum number %v not allowed; allowed quorum numbers: %v", quorum, allowedNumbers)
		}
	}
	return nil
}
