package meterer

import (
	"fmt"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
)

type QuorumPeriodRecords map[core.QuorumID][]*PeriodRecord

// PeriodRecord contains the index of the reservation period and the usage of the period
type PeriodRecord struct {
	// Index is start timestamp of the period in seconds; it is always a multiple of the reservation window
	Index uint32
	// Usage is the usage of the period in symbols
	Usage uint64
}

// updateRecord tracks a successful update for rollback purposes
type UpdateRecord struct {
	QuorumNumber core.QuorumID
	Period       uint64
	Usage        uint64
}

// GetRelativePeriodRecord returns the period record for the given index and quorum number; if the record does not exist, it is initialized to 0
func (pr QuorumPeriodRecords) GetRelativePeriodRecord(index uint64, quorumNumber core.QuorumID) *PeriodRecord {
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
func (pr QuorumPeriodRecords) UpdateUsage(
	quorumNumber core.QuorumID,
	currentPeriod uint64,
	overflowPeriod uint64,
	symbolUsage uint64,
	binLimit uint64,
) error {
	if symbolUsage > binLimit {
		return fmt.Errorf("symbol usage exceeds bin limit")
	}

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

// FromProtoRecords converts protobuf period records to QuorumPeriodRecords
func FromProtoRecords(protoRecords map[uint32]*disperser_rpc.PeriodRecords) QuorumPeriodRecords {
	records := make(QuorumPeriodRecords)
	for quorumNumber, protoRecord := range protoRecords {
		records[core.QuorumID(quorumNumber)] = make([]*PeriodRecord, MinNumBins)
		// Initialize all records to 0
		for i := range records[core.QuorumID(quorumNumber)] {
			records[core.QuorumID(quorumNumber)][i] = &PeriodRecord{
				Index: uint32(i),
				Usage: 0,
			}
		}
		// Populate with values from server
		for _, record := range protoRecord.GetRecords() {
			idx := record.Index % uint32(MinNumBins)
			records[core.QuorumID(quorumNumber)][idx] = &PeriodRecord{
				Index: record.Index,
				Usage: record.Usage,
			}
		}
	}
	return records
}
