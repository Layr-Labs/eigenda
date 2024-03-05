package dataapi

import "fmt"

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
		// In practice it should be quite small given that the quorum change is
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
// The parameters:
//   - startBlock, endBlock: specifying the block interval of interest.
//     Requires: startBlock <= endBlock.
//   - operatorInitialQuorum: the initial quorums at startBlock that operators were
//     registered in.
//     Requires: operatorInitialQuorum[op] is non-empty for each operator "op".
//   - addedToQuorum, removedFromQuorum: a sequence of events that added/removed operators
//     to/from quorums.
//     Requires:
//     1) the block numbers for all events are in range [startBlock+1, endBlock];
//     2) the events are in ascending order by block number for each operator "op".
func CreateOperatorQuorumIntervals(
	startBlock uint32,
	endBlock uint32,
	operatorInitialQuorum map[string][]uint8,
	addedToQuorum map[string][]*OperatorQuorum,
	removedFromQuorum map[string][]*OperatorQuorum,
) (OperatorQuorumIntervals, error) {
	if startBlock > endBlock {
		msg := "the endBlock must be no less than startBlock, but found " +
			"startBlock: %d, endBlock: %d"
		return nil, fmt.Errorf(msg, startBlock, endBlock)
	}
	operatorQuorumIntervals := make(OperatorQuorumIntervals)
	addedToQuorumErr := "cannot add operator %s to quorum %d at block number %d, " +
		"the operator is already in the quorum since block number %d"
	for op, initialQuorums := range operatorInitialQuorum {
		if len(initialQuorums) == 0 {
			return nil, fmt.Errorf("operator %s must be in at least one quorum at block %d", op, startBlock)
		}
		operatorQuorumIntervals[op] = make(map[uint8][]BlockInterval)
		openQuorum := make(map[uint8]uint32)
		for _, q := range initialQuorums {
			openQuorum[q] = startBlock
		}
		added := addedToQuorum[op]
		removed := removedFromQuorum[op]
		if eventErr := validateQuorumEvents(added, removed, startBlock, endBlock); eventErr != nil {
			return nil, eventErr
		}
		i, j := 0, 0
		for i < len(added) && j < len(removed) {
			// TODO(jianoaix): Having quorum addition and removal in the same block is a valid case.
			// Come up a followup fix to handle this special case.
			if added[i].BlockNumber == removed[j].BlockNumber {
				msg := "Not yet supported: operator was adding and removing quorums at the " +
					"same block, operator: %s, block number: %d"
				return nil, fmt.Errorf(msg, op, added[i].BlockNumber)
			}
			if added[i].BlockNumber < removed[j].BlockNumber {
				for _, q := range added[i].QuorumNumbers {
					if start, ok := openQuorum[q]; ok {
						return nil, fmt.Errorf(addedToQuorumErr, op, q, added[i].BlockNumber, start)
					}
					openQuorum[q] = added[i].BlockNumber
				}
				i++
			} else {
				if err := removeQuorums(removed[j], openQuorum, operatorQuorumIntervals); err != nil {
					return nil, err
				}
				j++
			}
		}
		for ; i < len(added); i++ {
			for _, q := range added[i].QuorumNumbers {
				if start, ok := openQuorum[q]; ok {
					return nil, fmt.Errorf(addedToQuorumErr, op, q, added[i].BlockNumber, start)
				}
				openQuorum[q] = added[i].BlockNumber
			}
		}
		for ; j < len(removed); j++ {
			if err := removeQuorums(removed[j], openQuorum, operatorQuorumIntervals); err != nil {
				return nil, err
			}
		}
		for q, start := range openQuorum {
			interval := BlockInterval{
				StartBlock: start,
				EndBlock:   endBlock,
			}
			if _, ok := operatorQuorumIntervals[op][q]; !ok {
				operatorQuorumIntervals[op][q] = make([]BlockInterval, 0)
			}
			operatorQuorumIntervals[op][q] = append(operatorQuorumIntervals[op][q], interval)
		}
	}

	return operatorQuorumIntervals, nil
}

// removeQuorums handles a quorum removal event, which marks the end of membership in a quorum,
// so it'll form a block interval.
func removeQuorums(operatorQuorum *OperatorQuorum, openQuorum map[uint8]uint32, result OperatorQuorumIntervals) error {
	op := operatorQuorum.Operator
	for _, q := range operatorQuorum.QuorumNumbers {
		start, ok := openQuorum[q]
		if !ok {
			msg := "cannot remove a quorum %d, the operator %s is not yet in the quorum " +
				"at block number %d"
			return fmt.Errorf(msg, q, op, operatorQuorum.BlockNumber)
		}
		if start >= operatorQuorum.BlockNumber {
			msg := "deregistration block number %d must be strictly greater than its " +
				"registration block number %d, for operator %s, quorum %d"
			return fmt.Errorf(msg, operatorQuorum.BlockNumber, start, op, q)
		}
		interval := BlockInterval{
			StartBlock: start,
			// The operator is NOT live at the block it's deregistered.
			EndBlock: operatorQuorum.BlockNumber - 1,
		}
		if _, ok = result[op][q]; !ok {
			result[op][q] = make([]BlockInterval, 0)
		}
		result[op][q] = append(result[op][q], interval)
		delete(openQuorum, q)
	}
	return nil
}

// validateQuorumEvents validates the operator quorum events have the desired block numbers and are
// in ascending order by block numbers.
func validateQuorumEvents(added []*OperatorQuorum, removed []*OperatorQuorum, startBlock, endBlock uint32) error {
	validate := func(events []*OperatorQuorum) error {
		for i := range events {
			if events[i].BlockNumber <= startBlock || events[i].BlockNumber > endBlock {
				return fmt.Errorf("quorum events must be in range [%d, %d]", startBlock+1, endBlock)
			}
			if i > 0 && events[i].BlockNumber < events[i-1].BlockNumber {
				return fmt.Errorf("quorum events must be in ascending order by block number")
			}
		}
		return nil
	}
	if err := validate(added); err != nil {
		return err
	}
	return validate(removed)
}
