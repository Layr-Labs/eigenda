package payment

import (
	"errors"
	"fmt"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
)

// QuorumPeriodRecords is a map from quorum number to a slice of period records
type QuorumPeriodRecords map[uint8][]*PeriodRecord

// PeriodRecord represents a single period record for a quorum
type PeriodRecord struct {
	// Index is start timestamp of the period in seconds; it is always a multiple of the reservation window
	Index uint32
	// Usage is the usage of the period in symbols
	Usage uint64
}

// updateRecord tracks a successful update for rollback purposes
type UpdateRecord struct {
	QuorumNumber uint8
	Period       uint64
	Usage        uint64
}

// GetRelativePeriodRecord returns the period record for the given index and quorum number; if the record does not exist, it is initialized to 0
func (pr QuorumPeriodRecords) GetRelativePeriodRecord(index uint64, quorumNumber uint8) *PeriodRecord {
	if _, exists := pr[quorumNumber]; !exists {
		pr[quorumNumber] = make([]*PeriodRecord, MinNumBins)
	}
	relativeIndex := uint32(index % uint64(MinNumBins))
	if pr[quorumNumber][relativeIndex] == nil {
		pr[quorumNumber][relativeIndex] = &PeriodRecord{
			Index: uint32(index),
			Usage: 0,
		}
	}
	return pr[quorumNumber][relativeIndex]
}

// UpdateUsage attempts to update the usage for a quorum's reservation period
// Returns error if the update would exceed the bin limit and cannot use overflow bin
//
// The function maintains a fixed-size circular buffer of numBins slots
// to track usage across an unbounded sequence of time periods by mapping each
// "absolute" period index onto a "relative" buffer index via modular arithmetic.
//
// Incoming timestamps are first bucketed into discrete reservation periods of
// length reservationWindow, yielding an integer period index. When a request
// for period p arrives, the system computes its buffer slot i = p mod numBins;
// if the stored period at slot i differs from p, the slot is reset (index
// updated, usage cleared) before accumulating usage.
//
// Controlled overflow allows unused capacity in a full bin to spill into future
// bins under strict conditions, and a sliding valid-period window ensures only
// recent periods are accepted.
func (pr QuorumPeriodRecords) UpdateUsage(
	quorumNumber uint8,
	timestamp int64,
	numSymbols uint64,
	reservation *ReservedPayment,
	protocolConfig *PaymentQuorumProtocolConfig,
) error {
	if reservation == nil {
		return errors.New("reservation cannot be nil")
	}
	if protocolConfig == nil {
		return errors.New("protocolConfig cannot be nil")
	}

	symbolUsage := SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
	binLimit := GetBinLimit(reservation.SymbolsPerSecond, protocolConfig.ReservationRateLimitWindow)

	if symbolUsage > binLimit {
		return errors.New("symbol usage exceeds bin limit")
	}

	currentPeriod := GetReservationPeriodByNanosecond(timestamp, protocolConfig.ReservationRateLimitWindow)
	relativePeriodRecord := pr.GetRelativePeriodRecord(currentPeriod, quorumNumber)
	oldUsage := relativePeriodRecord.Usage
	relativePeriodRecord.Usage += symbolUsage

	// within the bin limit
	if relativePeriodRecord.Usage <= binLimit {
		return nil
	}

	if oldUsage >= binLimit {
		return fmt.Errorf("reservation limit exceeded for quorum %d", quorumNumber)
	}

	// overflow bin if we're over the limit
	overflowPeriod := GetOverflowPeriod(currentPeriod, protocolConfig.ReservationRateLimitWindow)
	overflowPeriodRecord := pr.GetRelativePeriodRecord(overflowPeriod, quorumNumber)
	if overflowPeriodRecord.Usage == 0 {
		overflowPeriodRecord.Usage += relativePeriodRecord.Usage - binLimit
		relativePeriodRecord.Usage = binLimit
		return nil
	}

	return fmt.Errorf("reservation limit exceeded for quorum %d", quorumNumber)
}

// Make a deep copy of the period records
func (pr QuorumPeriodRecords) DeepCopy() QuorumPeriodRecords {
	copied := make(QuorumPeriodRecords)
	for quorumNumber, records := range pr {
		copied[quorumNumber] = make([]*PeriodRecord, len(records))
		for i, record := range records {
			if record != nil {
				// Create a new PeriodRecord with the same values
				copied[quorumNumber][i] = &PeriodRecord{
					Index: record.Index,
					Usage: record.Usage,
				}
			}
		}
	}
	return copied
}

// FromProtoRecords converts protobuf period records to native types
func FromProtoRecords(pbRecords map[uint32]*disperser_rpc.PeriodRecords) QuorumPeriodRecords {
	records := make(QuorumPeriodRecords)
	for quorumNumber, pbRecord := range pbRecords {
		if pbRecord == nil {
			continue
		}
		records[uint8(quorumNumber)] = make([]*PeriodRecord, MinNumBins)
		for i := range records[uint8(quorumNumber)] {
			records[uint8(quorumNumber)][i] = &PeriodRecord{
				Index: uint32(i),
				Usage: 0,
			}
		}
		// Populate with values from server
		for _, record := range pbRecord.GetRecords() {
			idx := record.Index % uint32(MinNumBins)
			records[uint8(quorumNumber)][idx] = &PeriodRecord{
				Index: record.Index,
				Usage: record.Usage,
			}
		}
	}
	return records
}
