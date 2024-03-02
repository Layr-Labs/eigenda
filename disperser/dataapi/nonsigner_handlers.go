package dataapi

import "fmt"

type OperatorQuorum struct {
	Operator      string
	QuorumNumbers []byte
	BlockNumber   uint32
}

// Representing an interval [StartBlock, EndBlock] (inclusive).
type BlockInterval struct {
	StartBlock uint32
	EndBlock   uint32
}

// OperatorQuorumIntervals[op][q] is a sequence of increasing and non-overlapping
// intervals during which the operator "op" is registered in quorum "q".
type OperatorQuorumIntervals map[string]map[uint8][]BlockInterval

// GetQuorums returns the quorums the operator is registered in at the given block number.
func (oqi OperatorQuorumIntervals) GetQuorums(operatorId string, blockNum uint32) []uint8 {
	quorums := make([]uint8, 0)
	for q, intervals := range oqi[operatorId] {
		// Note: if len(intervals) is large, we can perform binary search here.
		// In practice it shouldn't be quite small given that the quorum change is
		// not frequent, so search it with brute force here.
		live := false
		for _, interval := range intervals {
			if interval.StartBlock > blockNum {
				break
			}
			if blockNum <= interval.EndBlock {
				live = true
				break
			}
		}
		if live {
			quorums = append(quorums, q)
		}
	}
	return quorums
}

// CreateOperatorQuorumIntervals creates OperatorQuorumIntervals that are within the
// the block interval [startBlock, endBlock] for operators.
//
// The initial quorums at startBlock for all operators are provided by
// "operatorInitialQuorum", and the quorum change events during [startBlock + 1, endBlock]
// are provided by "addedToQuorum" and "removedFromQuorum".
func CreateOperatorQuorumIntervals(
	startBlock uint32,
	endBlock uint32,
	operatorInitialQuorum map[string][]uint8,
	addedToQuorum map[string][]*OperatorQuorum,
	removedFromQuorum map[string][]*OperatorQuorum,
) (OperatorQuorumIntervals, error) {
	if startBlock > endBlock {
		return nil, fmt.Errorf("the startBlock must be no less than endBlock, but found startBlock: %d, endBlock: %d", startBlock, endBlock)
	}
	operatorQuorumIntervals := make(OperatorQuorumIntervals)
	for op, initialQuorums := range operatorInitialQuorum {
		if len(initialQuorums) == 0 {
			return nil, fmt.Errorf("the operator: %s must be in at least one quorum at block: %d", op, startBlock)
		}
		operatorQuorumIntervals[op] = make(map[uint8][]BlockInterval)
		openQuorum := make(map[uint8]uint32)
		for _, q := range initialQuorums {
			openQuorum[q] = startBlock
		}
		added := addedToQuorum[op]
		removed := removedFromQuorum[op]
		i, j := 0, 0
		for i < len(added) && j < len(removed) {
			// Skip the block if it's not after startBlock, because the operatorInitialQuorum
			// already gives the state at startBlock.
			if added[i].BlockNumber <= startBlock {
				i++
				continue
			}
			if removed[j].BlockNumber <= startBlock {
				j++
				continue
			}
			// TODO: Having quorum addition and removal in the same block is a valid case.
			// Will come up a followup fix to handle this special case.
			if added[i].BlockNumber == removed[j].BlockNumber {
				return nil, fmt.Errorf("Not yet supported: the operator: %s was adding and removing quorums at the same block number: %d", op, added[i].BlockNumber)
			}
			if added[i].BlockNumber < removed[j].BlockNumber {
				for _, q := range added[i].QuorumNumbers {
					start, ok := openQuorum[q]
					if ok {
						return nil, fmt.Errorf("cannot add operator: %s to quorum: %d at block number: %d, because it is already in the quorum since block number: %d", op, q, start, added[i].BlockNumber)
					}
					openQuorum[q] = added[i].BlockNumber
				}
				i++
			} else {
				err := removeQuorums(removed[j], openQuorum, operatorQuorumIntervals)
				if err != nil {
					return nil, err
				}
				j++
			}
		}
		for ; i < len(added); i++ {
			for _, q := range added[i].QuorumNumbers {
				start, ok := openQuorum[q]
				if ok {
					return nil, fmt.Errorf("cannot add operator: %s to quorum: %d at block number: %d, because it is already in the quorum since block number: %d", op, q, start, added[i].BlockNumber)
				}
				openQuorum[q] = added[i].BlockNumber
			}
		}
		for ; j < len(removed); j++ {
			err := removeQuorums(removed[j], openQuorum, operatorQuorumIntervals)
			if err != nil {
				return nil, err
			}
		}
		for q, start := range openQuorum {
			interval := BlockInterval{
				StartBlock: start,
				EndBlock:   endBlock,
			}
			_, ok := operatorQuorumIntervals[op][q]
			if !ok {
				operatorQuorumIntervals[op][q] = make([]BlockInterval, 0)
			}
			operatorQuorumIntervals[op][q] = append(operatorQuorumIntervals[op][q], interval)
		}
	}

	return operatorQuorumIntervals, nil
}

// removeQuorums handles a quorum removal event, which marks the end of membership in a
// quorum, so it'll form a block interval.
func removeQuorums(operatorQuorum *OperatorQuorum, openQuorum map[uint8]uint32, result OperatorQuorumIntervals) error {
	op := operatorQuorum.Operator
	for _, q := range operatorQuorum.QuorumNumbers {
		start, ok := openQuorum[q]
		if !ok {
			return fmt.Errorf("cannot remove a quorum: %d, because the operator: %s is not in the quorum at block number: %d", q, op, operatorQuorum.BlockNumber)
		}
		if start >= operatorQuorum.BlockNumber {
			return fmt.Errorf("deregistration block number: %d must be strictly greater than its registration block number: %d, for operator: %s, quorum: %d", operatorQuorum.BlockNumber, start, op, q)
		}
		interval := BlockInterval{
			StartBlock: start,
			// The operator is NOT live at the block it's deregistered.
			EndBlock: operatorQuorum.BlockNumber - 1,
		}
		_, ok = result[op][q]
		if !ok {
			result[op][q] = make([]BlockInterval, 0)
		}
		result[op][q] = append(result[op][q], interval)
		delete(openQuorum, q)
	}
	return nil
}
