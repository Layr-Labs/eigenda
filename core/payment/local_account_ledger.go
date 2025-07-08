package payment

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/meterer/payment_logic"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

var _ AccountLedger = (*LocalAccountLedger)(nil)

// LocalAccountLedger implements AccountLedger for client-side payment tracking.
type LocalAccountLedger struct {
	// on-chain per-quorum reservation configs
	reservations map[core.QuorumID]*core.ReservedPayment
	// on-chain on-demand payment balance
	onDemand *core.OnDemandPayment
	// off-chain usage tracking for reservations
	periodRecords meterer.QuorumPeriodRecords
	// off-chain local consumed amount
	cumulativePayment *big.Int
	// concurrent access protection
	mutex sync.RWMutex
}

// NewLocalAccountLedger creates LocalAccountLedger from disperser RPC response data.
func NewLocalAccountLedger(
	reservations map[uint32]*disperser_rpc.QuorumReservation,
	periodRecords map[uint32]*disperser_rpc.PeriodRecords,
	onchainCumulativePayment []byte,
	cumulativePayment []byte,
) (*LocalAccountLedger, error) {
	onchainReservations := meterer.ReservationsFromProtobuf(reservations)
	localPeriodRecords := meterer.FromProtoRecords(periodRecords)

	onchainOnDemand := &core.OnDemandPayment{
		CumulativePayment: new(big.Int).SetBytes(onchainCumulativePayment),
	}
	localOnDemand := new(big.Int).SetBytes(cumulativePayment)

	return &LocalAccountLedger{
		reservations:      onchainReservations,
		onDemand:          onchainOnDemand,
		periodRecords:     localPeriodRecords,
		cumulativePayment: localOnDemand,
	}, nil
}

// CreatePaymentHeader determines payment method and creates PaymentMetadata.
// Logic: try reservations first (CumulativePayment=0), fallback to on-demand (CumulativePayment=new total).
// No changes made to the AccountLedger state.
func (lal *LocalAccountLedger) CreatePaymentHeader(
	accountID gethcommon.Address,
	timestampNs int64,
	numSymbols uint64,
	quorumNumbers []core.QuorumID,
	params *meterer.PaymentVaultParams,
	receivedAtNs int64,
) (core.PaymentMetadata, error) {
	if len(quorumNumbers) == 0 {
		return core.PaymentMetadata{}, fmt.Errorf("no quorums provided")
	}
	if numSymbols == 0 {
		return core.PaymentMetadata{}, fmt.Errorf("zero symbols requested")
	}

	lal.mutex.RLock()
	defer lal.mutex.RUnlock()

	// Try reservation first
	reservationValidationErr := payment_logic.ValidateReservations(lal.reservations, params.QuorumProtocolConfigs, quorumNumbers, timestampNs, receivedAtNs)
	if reservationValidationErr == nil {
		// Check usage limits with deep copy; no changes made to the original periodRecords
		periodRecordsCopy := lal.periodRecords.DeepCopy()
		var reservationUsageErr error
		for _, quorumNumber := range quorumNumbers {
			reservation := lal.reservations[quorumNumber]
			_, protocolConfig, err := params.GetQuorumConfigs(quorumNumber)
			if err != nil {
				reservationUsageErr = err
				break
			}
			if err := periodRecordsCopy.UpdateUsage(quorumNumber, timestampNs, numSymbols, reservation, protocolConfig); err != nil {
				reservationUsageErr = err
				break
			}
		}

		if reservationUsageErr == nil {
			return core.PaymentMetadata{
				AccountID:         accountID,
				Timestamp:         timestampNs,
				CumulativePayment: big.NewInt(0),
			}, nil
		}
	}

	// Fall back to on-demand payment
	onDemandErr := payment_logic.ValidateQuorum(quorumNumbers, params.OnDemandQuorumNumbers)
	if onDemandErr != nil {
		return core.PaymentMetadata{}, fmt.Errorf("invalid requested quorum for on-demand: %v", onDemandErr)
	}

	paymentQuorumConfig, protocolConfig, err := params.GetQuorumConfigs(meterer.OnDemandQuorumID)
	if err != nil {
		return core.PaymentMetadata{}, fmt.Errorf("invalid quorum config for on-demand quorum: %v", err)
	}

	symbolsCharged := payment_logic.SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
	paymentCharged := payment_logic.PaymentCharged(symbolsCharged, paymentQuorumConfig.OnDemandPricePerSymbol)

	resultingPayment := new(big.Int).Add(lal.cumulativePayment, paymentCharged)
	if resultingPayment.Cmp(lal.onDemand.CumulativePayment) <= 0 {
		return core.PaymentMetadata{
			AccountID:         accountID,
			Timestamp:         timestampNs,
			CumulativePayment: resultingPayment,
		}, nil
	}

	return core.PaymentMetadata{}, fmt.Errorf("insufficient ondemand payment: required %v, available %v", resultingPayment, lal.onDemand.CumulativePayment)
}

// Debit applies payment state updates based on the DebitSlip payment metadata.
// Uses CumulativePayment=0 for reservation updates, non-zero for on-demand updates.
func (lal *LocalAccountLedger) Debit(
	ctx context.Context,
	slip *DebitSlip,
	params *meterer.PaymentVaultParams,
) (*big.Int, error) {
	if err := lal.validateDebitInputs(slip, params); err != nil {
		return nil, err
	}

	lal.mutex.Lock()
	defer lal.mutex.Unlock()

	if slip.PaymentMetadata.CumulativePayment.Cmp(big.NewInt(0)) == 0 {
		return lal.processReservationDebit(slip, params)
	}
	return lal.processOnDemandDebit(slip.PaymentMetadata.CumulativePayment)
}

func (lal *LocalAccountLedger) validateDebitInputs(slip *DebitSlip, params *meterer.PaymentVaultParams) error {
	if slip == nil {
		return fmt.Errorf("debit slip cannot be nil")
	}
	if params == nil {
		return fmt.Errorf("payment vault params cannot be nil")
	}
	if len(slip.QuorumNumbers) == 0 {
		return ErrNoQuorums
	}
	if slip.NumSymbols == 0 {
		return ErrZeroSymbols
	}
	return nil
}

// processReservationDebit handles reservation-based payment debits
func (lal *LocalAccountLedger) processReservationDebit(
	slip *DebitSlip,
	params *meterer.PaymentVaultParams,
) (*big.Int, error) {
	periodRecordsCopy := lal.periodRecords.DeepCopy()

	for _, quorumNumber := range slip.QuorumNumbers {
		reservation := lal.reservations[quorumNumber]
		if reservation == nil {
			return nil, fmt.Errorf("no reservation found for quorum %d", quorumNumber)
		}

		_, protocolConfig, err := params.GetQuorumConfigs(quorumNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get config for quorum %d: %w", quorumNumber, err)
		}

		if err := periodRecordsCopy.UpdateUsage(quorumNumber, slip.GetTimestamp(), slip.NumSymbols, reservation, protocolConfig); err != nil {
			return nil, fmt.Errorf("reservation usage failed for quorum %d: %w", quorumNumber, err)
		}
	}

	lal.periodRecords = periodRecordsCopy
	return nil, nil
}

// processOnDemandDebit handles on-demand payment debits
func (lal *LocalAccountLedger) processOnDemandDebit(paymentAmount *big.Int) (*big.Int, error) {
	// Validate payment amount
	if paymentAmount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid on-demand payment amount: %v", paymentAmount)
	}

	// Ensure we don't go backwards in cumulative payment
	if paymentAmount.Cmp(lal.cumulativePayment) < 0 {
		return nil, fmt.Errorf("payment amount %v is less than current cumulative payment %v",
			paymentAmount, lal.cumulativePayment)
	}

	// Update cumulative payment
	lal.cumulativePayment.Set(paymentAmount)
	return new(big.Int).Set(paymentAmount), nil
}

// RevertDebit undoes previous debit operations using a DebitSlip.
// Reverts reservation usage (nil payment) or on-demand payment (non-nil payment).
// Note: current Accountant doesn't support reverting; this function will not be used in the near future.
func (lal *LocalAccountLedger) RevertDebit(
	ctx context.Context,
	slip *DebitSlip,
	params *meterer.PaymentVaultParams,
	previousCumulativePayment *big.Int,
) error {
	lal.mutex.Lock()
	defer lal.mutex.Unlock()

	if previousCumulativePayment == nil {
		// Revert reservation usage
		if params == nil {
			return errors.New("payment vault params cannot be nil")
		}

		periodRecordsCopy := lal.periodRecords.DeepCopy()

		for _, quorumNumber := range slip.QuorumNumbers {
			_, protocolConfig, err := params.GetQuorumConfigs(quorumNumber)
			if err != nil {
				return fmt.Errorf("failed to get config for quorum %d: %w", quorumNumber, err)
			}
			reservationPeriod := payment_logic.GetReservationPeriodByNanosecond(slip.GetTimestamp(), protocolConfig.ReservationRateLimitWindow)

			records := periodRecordsCopy[quorumNumber]
			for _, record := range records {
				if record != nil && record.Index == uint32(reservationPeriod) {
					record.Usage -= slip.NumSymbols
					break
				}
			}
		}

		lal.periodRecords = periodRecordsCopy
		return nil
	}

	// Revert on-demand payment
	lal.cumulativePayment = new(big.Int).Set(previousCumulativePayment)
	return nil
}

// GetAccountStateProtobuf returns account state in protobuf format for RPC transmission.
// Note: this function is not used in the current implementation as disperser client will not be sharing the state through protobuf.
func (lal *LocalAccountLedger) GetAccountStateProtobuf() (
	reservations map[uint32]*disperser_rpc.QuorumReservation,
	periodRecords map[uint32]*disperser_rpc.PeriodRecords,
	onchainCumulativePayment []byte,
	cumulativePayment []byte,
) {
	lal.mutex.RLock()
	defer lal.mutex.RUnlock()

	protoReservations := make(map[uint32]*disperser_rpc.QuorumReservation)
	for quorumID, reservation := range lal.reservations {
		protoReservations[uint32(quorumID)] = &disperser_rpc.QuorumReservation{
			SymbolsPerSecond: reservation.SymbolsPerSecond,
			StartTimestamp:   uint32(reservation.StartTimestamp),
			EndTimestamp:     uint32(reservation.EndTimestamp),
		}
	}

	protoPeriodRecords := make(map[uint32]*disperser_rpc.PeriodRecords)
	for quorumID, records := range lal.periodRecords {
		if len(records) > 0 {
			protoRecords := make([]*disperser_rpc.PeriodRecord, 0, len(records))
			for _, record := range records {
				if record != nil {
					protoRecords = append(protoRecords, &disperser_rpc.PeriodRecord{
						Index: record.Index,
						Usage: record.Usage,
					})
				}
			}
			protoPeriodRecords[uint32(quorumID)] = &disperser_rpc.PeriodRecords{
				Records: protoRecords,
			}
		}
	}

	var onchainCumulativePaymentBytes []byte
	if lal.onDemand != nil && lal.onDemand.CumulativePayment != nil {
		onchainCumulativePaymentBytes = lal.onDemand.CumulativePayment.Bytes()
	}

	var cumulativePaymentBytes []byte
	if lal.cumulativePayment != nil {
		cumulativePaymentBytes = lal.cumulativePayment.Bytes()
	}

	return protoReservations, protoPeriodRecords, onchainCumulativePaymentBytes, cumulativePaymentBytes
}
