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
	// nonsignerAddresses[i] is the address for nonsigners[i].
	nonsignerAddresses, err := s.transactor.BatchOperatorIDToAddress(ctx, nonsigners)
	if err != nil {
		return nil, err
	}
	nonsignerIdToAddress := make(map[core.OperatorID]string)
	nonsignerAddressToId := make(map[string]core.OperatorID)
	for i := range nonsigners {
		nonsignerIdToAddress[nonsigners[i]] = nonsignerAddresses[i].Hex()
		nonsignerAddressToId[nonsignerAddresses[i].Hex()] = nonsigners[i]
	}

	// Get operators' quorums at startBlock.
	bitmaps, err := s.transactor.GetQuorumBitmapForOperatorsAtBlockNumber(ctx, nonsigners, startBlock)
	if err != nil {
		return nil, err
	}
	operatorInitialQuorum := make(map[string][]uint8)
	for i := range bitmaps {
		operatorInitialQuorum[nonsigners[i].Hex()] = eth.BitmapToQuorumIds(bitmaps[i])
	}

	// Get operators' quorum change events from [startBlock+1, endBlock].
	addedToQuorum := make(map[string][]*OperatorQuorum)
	removedFromQuorum := make(map[string][]*OperatorQuorum)
	if startBlock < endBlock {
		operatorQuorumEvents, err := s.subgraphClient.QueryOperatorQuorumEvent(ctx, startBlock+1, endBlock)
		if err != nil {
			return nil, err
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
	}

	// Create operators' quorum intervals.
	operatorQuorumIntervals, err := CreateOperatorQuorumIntervals(startBlock, endBlock, operatorInitialQuorum, addedToQuorum, removedFromQuorum)
	if err != nil {
		return nil, err
	}
	// Compute num batches failed to sign for each <operatorId, quorumId>
	numFailed := make(map[string]map[uint8]int)
	for _, b := range batches {
		for _, op := range b.NonSigners {
			op := op[2:]
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

	// Create quorumBatches, where quorumBatches[q].AccuBatches is the total number of
	// batches in block interval [startBlock, b] for quorum "q".
	quorumBatches := CreatQuorumBatches(batches)

	// Compute num batches responsible to sign for each <operatorId, quorumId>
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

	operators := make([]*OperatorNonsigningPercentageMetrics, 0)
	for op, val := range numResponsible {
		for q, num := range val {
			if numResponsible[op][q] > 0 {
				if unsignedCount, ok := numFailed[op][q]; ok {
					ps := fmt.Sprintf("%.2f", (float64(unsignedCount)/float64(num))*100)
					pf, err := strconv.ParseFloat(ps, 64)
					if err != nil {
						return nil, err
					}
					operatorMetric := OperatorNonsigningPercentageMetrics{
						OperatorId:           op,
						QuorumId:             q,
						TotalUnsignedBatches: unsignedCount,
						TotalBatches:         num,
						Percentage:           pf,
					}
					operators = append(operators, &operatorMetric)
				}
			}
		}
	}

	sort.Slice(operators, func(i, j int) bool {
		if operators[i].Percentage == operators[j].Percentage {
			return operators[i].OperatorId < operators[j].OperatorId
		}
		return operators[i].Percentage > operators[j].Percentage
	})

	return &OperatorsNonsigningPercentage{
		Meta: Meta{
			Size: len(operators),
		},
		Data: operators,
	}, nil
}
