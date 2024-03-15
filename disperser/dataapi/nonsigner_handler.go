package dataapi

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
)

func (s *server) getOperatorNonsigningRate(ctx context.Context, intervalSeconds int64) (*OperatorsNonsigningPercentage, error) {
	batches, err := s.subgraphClient.QueryBatchNonSigningInfoInInterval(ctx, intervalSeconds)
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
	for i := range nonsigners {
		nonsignerAddressToId[nonsignerAddresses[i].Hex()] = nonsigners[i]
	}

	// Create operators' quorum intervals.
	operatorQuorumIntervals, err := s.createOperatorQuorumIntervals(ctx, nonsigners, nonsignerAddressToId, startBlock, endBlock)
	if err != nil {
		return nil, err
	}

	// Compute num batches failed, where numFailed[op][q] is the number of batches
	// failed to sign for operator "op" and quorum "q".
	numFailed := computeNumFailed(batches, operatorQuorumIntervals)

	// Compute num batches responsible, where numResponsible[op][q] is the number of batches
	// that operator "op" and quorum "q" are responsible for.
	numResponsible := computeNumResponsible(batches, operatorQuorumIntervals)

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
				nonsignerMetric := OperatorNonsigningPercentageMetrics{
					OperatorId:           fmt.Sprintf("0x%s", op),
					QuorumId:             q,
					TotalUnsignedBatches: unsignedCount,
					TotalBatches:         totalCount,
					Percentage:           pf,
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
		Meta: Meta{
			Size: len(nonsignerMetrics),
		},
		Data: nonsignerMetrics,
	}, nil
}

func (s *server) createOperatorQuorumIntervals(ctx context.Context, nonsigners []core.OperatorID, nonsignerAddressToId map[string]core.OperatorID, startBlock, endBlock uint32) (OperatorQuorumIntervals, error) {
	// Get operators' initial quorums (at startBlock).
	bitmaps, err := s.transactor.GetQuorumBitmapForOperatorsAtBlockNumber(ctx, nonsigners, startBlock)
	if err != nil {
		return nil, err
	}
	operatorInitialQuorum := make(map[string][]uint8)
	for i := range bitmaps {
		operatorInitialQuorum[nonsigners[i].Hex()] = eth.BitmapToQuorumIds(bitmaps[i])
	}

	// Get operators' quorum change events from [startBlock+1, endBlock].
	addedToQuorum, removedFromQuorum, err := s.getOperatorQuorumEvents(ctx, startBlock, endBlock, nonsignerAddressToId)
	if err != nil {
		return nil, err
	}

	// Create operators' quorum intervals.
	operatorQuorumIntervals, err := CreateOperatorQuorumIntervals(startBlock, endBlock, operatorInitialQuorum, addedToQuorum, removedFromQuorum)
	if err != nil {
		return nil, err
	}

	return operatorQuorumIntervals, nil
}

func (s *server) getOperatorQuorumEvents(ctx context.Context, startBlock, endBlock uint32, nonsignerAddressToId map[string]core.OperatorID) (map[string][]*OperatorQuorum, map[string][]*OperatorQuorum, error) {
	addedToQuorum := make(map[string][]*OperatorQuorum)
	removedFromQuorum := make(map[string][]*OperatorQuorum)
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

func getNonSigners(batches []*BatchNonSigningInfo) ([]core.OperatorID, error) {
	nonsignerSet := map[string]struct{}{}
	for _, b := range batches {
		for _, op := range b.NonSigners {
			nonsignerSet[op] = struct{}{}
		}
	}
	nonsigners := make([]core.OperatorID, 0)
	for op := range nonsignerSet {
		hexstr := strings.TrimPrefix(op, "0x")
		b, err := hex.DecodeString(hexstr)
		if err != nil {
			return nil, err
		}
		nonsigners = append(nonsigners, core.OperatorID(b))
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

func computeNumFailed(batches []*BatchNonSigningInfo, operatorQuorumIntervals OperatorQuorumIntervals) map[string]map[uint8]int {
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

func computeNumResponsible(batches []*BatchNonSigningInfo, operatorQuorumIntervals OperatorQuorumIntervals) map[string]map[uint8]int {
	// Create quorumBatches, where quorumBatches[q].AccuBatches is the total number of
	// batches in block interval [startBlock, b] for quorum "q".
	quorumBatches := CreatQuorumBatches(batches)

	numResponsible := make(map[string]map[uint8]int)
	for op, val := range operatorQuorumIntervals {
		for q, intervals := range val {
			numBatches := 0
			for _, interval := range intervals {
				numBatches = numBatches + ComputeNumBatches(quorumBatches[q], interval.StartBlock, interval.EndBlock)
			}
			if _, ok := numResponsible[op]; !ok {
				numResponsible[op] = make(map[uint8]int)
			}
			numResponsible[op][q] = numBatches
		}
	}
	return numResponsible
}
