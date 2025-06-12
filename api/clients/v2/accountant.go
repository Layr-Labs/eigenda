package clients

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type Accountant struct {
	// on-chain states
	accountID    gethcommon.Address
	reservations map[core.QuorumID]*core.ReservedPayment
	onDemand     *core.OnDemandPayment

	paymentVaultParams *meterer.PaymentVaultParams

	// local accounting
	// contains 3 bins; circular wrapping of indices
	periodRecords     map[core.QuorumID][]PeriodRecord
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

// NewAccountant initializes an accountant with the given account ID, reservations, on-demand payment, and number of bins
// TODO: Consider making this initialization take all the fields as arguments or entirely empty as clients are typically
// syncing the onchain configurations and offchain usage with the disperser server
func NewAccountant(accountID gethcommon.Address, reservations map[uint8]*core.ReservedPayment, onDemand *core.OnDemandPayment) *Accountant {
	periodRecords := CreateEmptyReservationUsage(reservations)
	a := Accountant{
		accountID:    accountID,
		reservations: reservations,
		onDemand:     onDemand,
		paymentVaultParams: &meterer.PaymentVaultParams{
			QuorumPaymentConfigs:  make(map[uint8]*core.PaymentQuorumConfig),
			QuorumProtocolConfigs: make(map[uint8]*core.PaymentQuorumProtocolConfig),
			OnDemandQuorumNumbers: make([]uint8, 0),
		},
		periodRecords:     periodRecords,
		cumulativePayment: big.NewInt(0),
	}
	// TODO: add a routine to refresh the on-chain state occasionally?
	return &a
}

// updateReservationUsage updates the usage records for a quorum's reservation
func (a *Accountant) updateReservationUsage(quorumNumber core.QuorumID, currentPeriod uint64, symbolUsage uint64) error {
	res, exists := a.reservations[quorumNumber]
	if !exists {
		return fmt.Errorf("reservation not found for quorum %d", quorumNumber)
	}

	relativePeriodRecord := a.GetRelativePeriodRecord(currentPeriod, quorumNumber)
	relativePeriodRecord.Usage += symbolUsage

	quorumReservationWindow, err := a.GetReservationWindow(quorumNumber)
	if err != nil {
		return fmt.Errorf("failed to get reservation window for quorum %d: %w", quorumNumber, err)
	}

	binLimit := res.SymbolsPerSecond * uint64(quorumReservationWindow)

	// Check if we're within the bin limit
	if relativePeriodRecord.Usage <= binLimit {
		return nil
	}

	// Try to use overflow bin if we're over the limit
	overflowPeriodRecord := a.GetRelativePeriodRecord(meterer.GetOverflowPeriod(currentPeriod, quorumReservationWindow), quorumNumber)
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

// ReservationUsage attempts to use the reservation for the requested quorums
// Returns (0, nil) if successful, or (nil, error) if reservation cannot be used
func (a *Accountant) reservationUsage(numSymbols uint64, quorumNumbers []core.QuorumID, timestamp int64) error {
	// The two timestamps are the same for the accountant client for validating the reservation period; for the server the second timestamp is the received at time
	if err := meterer.ValidateReservations(a.reservations, a.paymentVaultParams.QuorumProtocolConfigs, quorumNumbers, timestamp, time.Unix(0, timestamp)); err != nil {
		return err
	}

	// Lock for updating usage records
	a.periodRecordsLock.Lock()
	defer a.periodRecordsLock.Unlock()

	for _, quorumNumber := range quorumNumbers {
		// Calculate current period and symbol usage
		reservationWindow, err := a.GetReservationWindow(quorumNumber)
		if err != nil {
			return fmt.Errorf("failed to get reservation window for quorum %d: %w", quorumNumber, err)
		}
		currentReservationPeriod := meterer.GetReservationPeriodByNanosecond(timestamp, reservationWindow)
		minNumSymbols, err := a.GetMinNumSymbols(quorumNumber)
		if err != nil {
			return err
		}
		symbolUsage := meterer.SymbolsCharged(numSymbols, minNumSymbols)

		if err := a.updateReservationUsage(quorumNumber, currentReservationPeriod, symbolUsage); err != nil {
			// Rollback usage for this quorum
			relativePeriodRecord := a.GetRelativePeriodRecord(currentReservationPeriod, quorumNumber)
			relativePeriodRecord.Usage -= symbolUsage
			return err
		}
	}

	return nil
}

// OnDemandUsage handles the on-demand payment calculation and validation
func (a *Accountant) onDemandUsage(numSymbols uint64, quorumNumbers []core.QuorumID) (*big.Int, error) {
	// Verify quorum requirements
	if err := meterer.ValidateQuorum(quorumNumbers, a.paymentVaultParams.OnDemandQuorumNumbers); err != nil {
		return nil, err
	}

	// Calculate payment needed for the number of symbols
	pricePerSymbol, err := a.GetPricePerSymbol(meterer.OnDemandQuorumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get price per symbol for on-demand quorum: %w", err)
	}
	minNumSymbols, err := a.GetMinNumSymbols(meterer.OnDemandQuorumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get min num symbols for on-demand quorum: %w", err)
	}
	symbolsCharged := meterer.SymbolsCharged(numSymbols, minNumSymbols)
	paymentCharged := meterer.PaymentCharged(symbolsCharged, pricePerSymbol)

	// Calculate the increment required to add to the cumulative payment
	a.onDemandLock.Lock()
	defer a.onDemandLock.Unlock()
	resultingPayment := new(big.Int).Add(a.cumulativePayment, paymentCharged)
	if resultingPayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		a.cumulativePayment.Add(a.cumulativePayment, paymentCharged)
		return a.cumulativePayment, nil
	}

	return nil, fmt.Errorf(
		"no bandwidth reservation found for account %s, and current cumulativePayment balance insufficient "+
			"to make an on-demand dispersal. Consider depositing more eth to the PaymentVault contract.", a.accountID.Hex())
}

// AccountBlob accountant provides and records payment information
func (a *Accountant) AccountBlob(
	timestamp int64,
	numSymbols uint64,
	quorums []uint8) (*core.PaymentMetadata, error) {

	// Always try to use reservation first
	err := a.reservationUsage(numSymbols, quorums, timestamp)
	if err == nil {
		return &core.PaymentMetadata{
			AccountID:         a.accountID,
			Timestamp:         timestamp,
			CumulativePayment: big.NewInt(0),
		}, nil
	}

	// Fall back to on-demand payment if reservation fails
	cumulativePayment, err := a.onDemandUsage(numSymbols, quorums)
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

// GetMinNumSymbols returns the minimum number of symbols for a given quorum
func (a *Accountant) GetMinNumSymbols(quorumID core.QuorumID) (uint64, error) {
	if a.paymentVaultParams == nil {
		return 0, fmt.Errorf("payment vault params is nil")
	}
	return a.paymentVaultParams.GetMinNumSymbols(quorumID)
}

// GetPricePerSymbol returns the price per symbol for a given quorum
func (a *Accountant) GetPricePerSymbol(quorumID core.QuorumID) (uint64, error) {
	if a.paymentVaultParams == nil {
		return 0, fmt.Errorf("payment vault params is nil")
	}
	return a.paymentVaultParams.GetPricePerSymbol(quorumID)
}

// GetReservationWindow returns the reservation window for a given quorum
func (a *Accountant) GetReservationWindow(quorumID core.QuorumID) (uint64, error) {
	if a.paymentVaultParams == nil {
		return 0, fmt.Errorf("payment vault params is nil")
	}
	return a.paymentVaultParams.GetReservationWindow(quorumID)
}

func (a *Accountant) GetRelativePeriodRecord(index uint64, quorumNumber core.QuorumID) *PeriodRecord {
	relativeIndex := uint32(index % uint64(meterer.MinNumBins))
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

	// Convert protobuf configs to core types
	quorumPaymentConfigs := make(map[core.QuorumID]*core.PaymentQuorumConfig)
	quorumProtocolConfigs := make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig)

	for quorumID, pbPaymentConfig := range vaultParams.GetQuorumPaymentConfigs() {
		quorumPaymentConfigs[core.QuorumID(quorumID)] = &core.PaymentQuorumConfig{
			ReservationSymbolsPerSecond: pbPaymentConfig.GetReservationSymbolsPerSecond(),
			OnDemandSymbolsPerSecond:    pbPaymentConfig.GetOnDemandSymbolsPerSecond(),
			OnDemandPricePerSymbol:      pbPaymentConfig.GetOnDemandPricePerSymbol(),
		}
	}

	for quorumID, pbProtocolConfig := range vaultParams.GetQuorumProtocolConfigs() {
		quorumProtocolConfigs[core.QuorumID(quorumID)] = &core.PaymentQuorumProtocolConfig{
			MinNumSymbols:              pbProtocolConfig.GetMinNumSymbols(),
			ReservationAdvanceWindow:   pbProtocolConfig.GetReservationAdvanceWindow(),
			ReservationRateLimitWindow: pbProtocolConfig.GetReservationRateLimitWindow(),
			OnDemandRateLimitWindow:    pbProtocolConfig.GetOnDemandRateLimitWindow(),
			OnDemandEnabled:            pbProtocolConfig.GetOnDemandEnabled(),
		}
	}

	// Convert uint32 slice to uint8 slice
	onDemandQuorumNumbers := make([]uint8, len(vaultParams.GetOnDemandQuorumNumbers()))
	for i, num := range vaultParams.GetOnDemandQuorumNumbers() {
		onDemandQuorumNumbers[i] = uint8(num)
	}

	a.paymentVaultParams = &meterer.PaymentVaultParams{
		QuorumPaymentConfigs:  quorumPaymentConfigs,
		QuorumProtocolConfigs: quorumProtocolConfigs,
		OnDemandQuorumNumbers: onDemandQuorumNumbers,
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
		a.reservations = make(map[core.QuorumID]*core.ReservedPayment)
	} else {
		a.reservations = make(map[core.QuorumID]*core.ReservedPayment)
		for quorumNumber, reservation := range paymentState.GetReservations() {
			a.reservations[core.QuorumID(quorumNumber)] = &core.ReservedPayment{
				SymbolsPerSecond: reservation.GetSymbolsPerSecond(),
				StartTimestamp:   uint64(reservation.GetStartTimestamp()),
				EndTimestamp:     uint64(reservation.GetEndTimestamp()),
			}
		}
		a.periodRecords = CreateEmptyReservationUsage(a.reservations)
		fmt.Println("created empty periodRecords", a.periodRecords)
		for quorumNumber, periodRecords := range paymentState.GetPeriodRecords() {
			if periodRecords != nil {
				for _, record := range periodRecords.GetRecords() {
					idx := record.Index % uint32(meterer.MinNumBins)
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
func CreateEmptyReservationUsage(quorumNumbers map[core.QuorumID]*core.ReservedPayment) map[core.QuorumID][]PeriodRecord {
	reservationUsage := make(map[core.QuorumID][]PeriodRecord)
	for quorumNumber := range quorumNumbers {
		reservationUsage[quorumNumber] = make([]PeriodRecord, meterer.MinNumBins)
		for i := range reservationUsage[quorumNumber] {
			reservationUsage[quorumNumber][i] = PeriodRecord{Index: uint32(i), Usage: 0}
		}
	}
	return reservationUsage
}
