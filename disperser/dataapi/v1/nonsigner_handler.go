package v1

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
)

func (s *server) getOperatorNonsigningRate(ctx context.Context, startTime, endTime int64, liveOnly bool) (*OperatorsNonsigningPercentage, error) {
	batches, err := s.subgraphClient.QueryBatchNonSigningInfoInInterval(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}
	if len(batches) == 0 {
		return &OperatorsNonsigningPercentage{}, nil
	}

	// Get the block interval of interest [startBlock, endBlock].
	startBlock := batches[0].ReferenceBlockNumber
	endBlock := batches[0].ReferenceBlockNumber
	for i := range batches {
		if startBlock > batches[i].ReferenceBlockNumber {
			startBlock = batches[i].ReferenceBlockNumber
		}
		if endBlock < batches[i].ReferenceBlockNumber {
			endBlock = batches[i].ReferenceBlockNumber
		}
	}

	// Get the nonsigner (in operatorId) list.
	nonsigners, err := getNonSigners(batches)
	if err != nil {
		return nil, err
	}
	if len(nonsigners) == 0 {
		return &OperatorsNonsigningPercentage{}, nil
	}

	// Get the address for the nonsigners (from their operatorIDs).
	// nonsignerAddresses[i] is the address for nonsigners[i].
	nonsignerAddresses, err := s.transactor.BatchOperatorIDToAddress(ctx, nonsigners)
	if err != nil {
		return nil, err
	}

	// Create a mapping from address to operatorID.
	nonsignerAddressToId := make(map[string]core.OperatorID)
	nonsignerIdToAddress := make(map[string]string)
	for i := range nonsigners {
		addr := strings.ToLower(nonsignerAddresses[i].Hex())
		nonsignerAddressToId[addr] = nonsigners[i]
		nonsignerIdToAddress[nonsigners[i].Hex()] = addr
	}

	// Create operators' quorum intervals.
	operatorQuorumIntervals, quorumIDs, err := s.createOperatorQuorumIntervals(ctx, nonsigners, nonsignerAddressToId, startBlock, endBlock)
	if err != nil {
		return nil, err
	}

	// Compute num batches failed, where numFailed[op][q] is the number of batches
	// failed to sign for operator "op" and quorum "q".
	numFailed := computeNumFailed(batches, operatorQuorumIntervals)

	// Compute num batches responsible, where numResponsible[op][q] is the number of batches
	// that operator "op" and quorum "q" are responsible for.
	numResponsible := computeNumResponsible(batches, operatorQuorumIntervals)

	state, err := s.chainState.GetOperatorState(ctx, uint(endBlock), quorumIDs)
	if err != nil {
		return nil, err
	}

	// Compute the nonsigning rate for each <operator, quorum> pair.
	nonsignerMetrics := make([]*OperatorNonsigningPercentageMetrics, 0)
	for op, val := range numResponsible {
		for q, totalCount := range val {
			if totalCount == 0 {
				continue
			}
			if unsignedCount, ok := numFailed[op][q]; ok {
				ps := fmt.Sprintf("%.2f", (float64(unsignedCount)/float64(totalCount))*100)
				pf, err := strconv.ParseFloat(ps, 64)
				if err != nil {
					return nil, err
				}

				opID, err := core.OperatorIDFromHex(op)
				if err != nil {
					return nil, err
				}

				stakePercentage := float64(0)
				if stake, ok := state.Operators[q][opID]; ok {
					totalStake := new(big.Float).SetInt(state.Totals[q].Stake)
					stakePercentage, _ = new(big.Float).Quo(
						new(big.Float).SetInt(stake.Stake),
						totalStake).Float64()
				} else if liveOnly {
					// Operator "opID" isn't live at "endBlock", skip it.
					continue
				}

				nonsignerMetric := OperatorNonsigningPercentageMetrics{
					OperatorId:           fmt.Sprintf("0x%s", op),
					OperatorAddress:      nonsignerIdToAddress[op],
					QuorumId:             q,
					TotalUnsignedBatches: unsignedCount,
					TotalBatches:         totalCount,
					Percentage:           pf,
					StakePercentage:      100 * stakePercentage,
				}
				nonsignerMetrics = append(nonsignerMetrics, &nonsignerMetric)
			}
		}
	}

	// Sort by descending order of nonsigning rate.
	sort.Slice(nonsignerMetrics, func(i, j int) bool {
		if nonsignerMetrics[i].Percentage == nonsignerMetrics[j].Percentage {
			if nonsignerMetrics[i].OperatorId == nonsignerMetrics[j].OperatorId {
				return nonsignerMetrics[i].QuorumId < nonsignerMetrics[j].QuorumId
			}
			return nonsignerMetrics[i].OperatorId < nonsignerMetrics[j].OperatorId
		}
		return nonsignerMetrics[i].Percentage > nonsignerMetrics[j].Percentage
	})

	return &OperatorsNonsigningPercentage{
		Meta: dataapi.Meta{
			Size: len(nonsignerMetrics),
		},
		Data: nonsignerMetrics,
	}, nil
}

func (s *server) createOperatorQuorumIntervals(ctx context.Context, nonsigners []core.OperatorID, nonsignerAddressToId map[string]core.OperatorID, startBlock, endBlock uint32) (dataapi.OperatorQuorumIntervals, []uint8, error) {
	// Get operators' initial quorums (at startBlock).
	quorumSeen := make(map[uint8]struct{}, 0)

	bitmaps, err := s.transactor.GetQuorumBitmapForOperatorsAtBlockNumber(ctx, nonsigners, startBlock)
	if err != nil {
		return nil, nil, err
	}
	operatorInitialQuorum := make(map[string][]uint8)
	for i := range bitmaps {
		opQuorumIDs := eth.BitmapToQuorumIds(bitmaps[i])
		operatorInitialQuorum[nonsigners[i].Hex()] = opQuorumIDs
		for _, q := range opQuorumIDs {
			quorumSeen[q] = struct{}{}
		}
	}

	// Get all quorums.
	allQuorums := make([]uint8, 0)
	for q := range quorumSeen {
		allQuorums = append(allQuorums, q)
	}

	// Get operators' quorum change events from [startBlock+1, endBlock].
	addedToQuorum, removedFromQuorum, err := s.getOperatorQuorumEvents(ctx, startBlock, endBlock, nonsignerAddressToId)
	if err != nil {
		return nil, nil, err
	}

	// Create operators' quorum intervals.
	operatorQuorumIntervals, err := dataapi.CreateOperatorQuorumIntervals(startBlock, endBlock, operatorInitialQuorum, addedToQuorum, removedFromQuorum)
	if err != nil {
		return nil, nil, err
	}

	return operatorQuorumIntervals, allQuorums, nil
}

func (s *server) getOperatorQuorumEvents(ctx context.Context, startBlock, endBlock uint32, nonsignerAddressToId map[string]core.OperatorID) (map[string][]*dataapi.OperatorQuorum, map[string][]*dataapi.OperatorQuorum, error) {
	addedToQuorum := make(map[string][]*dataapi.OperatorQuorum)
	removedFromQuorum := make(map[string][]*dataapi.OperatorQuorum)
	if startBlock == endBlock {
		return addedToQuorum, removedFromQuorum, nil
	}
	operatorQuorumEvents, err := s.subgraphClient.QueryOperatorQuorumEvent(ctx, startBlock+1, endBlock)
	if err != nil {
		return nil, nil, err
	}
	// Make quorum events organize by operatorID (instead of address) and drop those who
	// are not nonsigners.
	for op, events := range operatorQuorumEvents.AddedToQuorum {
		if id, ok := nonsignerAddressToId[op]; ok {
			addedToQuorum[id.Hex()] = events
		}
	}
	for op, events := range operatorQuorumEvents.RemovedFromQuorum {
		if id, ok := nonsignerAddressToId[op]; ok {
			removedFromQuorum[id.Hex()] = events
		}
	}
	return addedToQuorum, removedFromQuorum, nil
}

func getNonSigners(batches []*dataapi.BatchNonSigningInfo) ([]core.OperatorID, error) {
	nonsignerSet := map[string]struct{}{}
	for _, b := range batches {
		for _, op := range b.NonSigners {
			nonsignerSet[op] = struct{}{}
		}
	}
	nonsigners := make([]core.OperatorID, 0)
	for op := range nonsignerSet {
		id, err := core.OperatorIDFromHex(op)
		if err != nil {
			return nil, err
		}
		nonsigners = append(nonsigners, id)
	}
	sort.Slice(nonsigners, func(i, j int) bool {
		for k := range nonsigners[i] {
			if nonsigners[i][k] != nonsigners[j][k] {
				return nonsigners[i][k] < nonsigners[j][k]
			}
		}
		return false
	})
	return nonsigners, nil
}

func computeNumFailed(batches []*dataapi.BatchNonSigningInfo, operatorQuorumIntervals dataapi.OperatorQuorumIntervals) map[string]map[uint8]int {
	numFailed := make(map[string]map[uint8]int)
	for _, b := range batches {
		for _, op := range b.NonSigners {
			op := op[2:]
			// Note: avg number of quorums per operator is a small number, so use brute
			// force here (otherwise, we can create a map to make it more efficient)
			for _, operatorQuorum := range operatorQuorumIntervals.GetQuorums(op, b.ReferenceBlockNumber) {
				for _, batchQuorum := range b.QuorumNumbers {
					if operatorQuorum == batchQuorum {
						if _, ok := numFailed[op]; !ok {
							numFailed[op] = make(map[uint8]int)
						}
						numFailed[op][operatorQuorum]++
						break
					}
				}
			}
		}
	}
	return numFailed
}

func computeNumResponsible(batches []*dataapi.BatchNonSigningInfo, operatorQuorumIntervals dataapi.OperatorQuorumIntervals) map[string]map[uint8]int {
	// Create quorumBatches, where quorumBatches[q].AccuBatches is the total number of
	// batches in block interval [startBlock, b] for quorum "q".
	quorumBatches := dataapi.CreatQuorumBatches(batches)

	numResponsible := make(map[string]map[uint8]int)
	for op, val := range operatorQuorumIntervals {
		for q, intervals := range val {
			numBatches := 0
			if _, ok := quorumBatches[q]; ok {
				for _, interval := range intervals {
					numBatches = numBatches + dataapi.ComputeNumBatches(quorumBatches[q], interval.StartBlock, interval.EndBlock)
				}
			}
			if _, ok := numResponsible[op]; !ok {
				numResponsible[op] = make(map[uint8]int)
			}
			numResponsible[op][q] = numBatches
		}
	}
	return numResponsible
}
