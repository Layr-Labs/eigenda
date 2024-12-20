package v2

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
)

func (s *ServerV2) getMetric(ctx context.Context, startTime int64, endTime int64) (*dataapi.Metric, error) {
	blockNumber, err := s.chainReader.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}
	quorumCount, err := s.chainReader.GetQuorumCount(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get quorum count: %w", err)
	}
	// assume quorum IDs are consequent integers starting from 0
	quorumIDs := make([]core.QuorumID, quorumCount)
	for i := 0; i < int(quorumCount); i++ {
		quorumIDs[i] = core.QuorumID(i)
	}

	operatorState, err := s.chainState.GetOperatorState(ctx, uint(blockNumber), quorumIDs)
	if err != nil {
		return nil, err
	}
	if len(operatorState.Operators) != int(quorumCount) {
		return nil, fmt.Errorf("Requesting for %d quorums (quorumID=%v), but got %v", quorumCount, quorumIDs, operatorState.Operators)
	}
	totalStakePerQuorum := map[core.QuorumID]*big.Int{}
	for quorumID, opInfoByID := range operatorState.Operators {
		for _, opInfo := range opInfoByID {
			if s, ok := totalStakePerQuorum[quorumID]; !ok {
				totalStakePerQuorum[quorumID] = new(big.Int).Set(opInfo.Stake)
			} else {
				s.Add(s, opInfo.Stake)
			}
		}
	}

	throughput, err := s.metricsHandler.GetAvgThroughput(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}

	costInGas, err := s.calculateTotalCostGasUsed(ctx)
	if err != nil {
		return nil, err
	}

	return &dataapi.Metric{
		Throughput:          throughput,
		CostInGas:           costInGas,
		TotalStake:          totalStakePerQuorum[0],
		TotalStakePerQuorum: totalStakePerQuorum,
	}, nil
}

func (s *ServerV2) calculateTotalCostGasUsed(ctx context.Context) (float64, error) {
	return 0, nil
}

func (s *ServerV2) getNonSigners(ctx context.Context, intervalSeconds int64) (*[]dataapi.NonSigner, error) {
	nonSigners, err := s.subgraphClient.QueryBatchNonSigningOperatorIdsInInterval(ctx, intervalSeconds)
	if err != nil {
		return nil, err
	}

	nonSignersObj := make([]dataapi.NonSigner, 0)
	for nonSigner, nonSigningAmount := range nonSigners {
		s.logger.Info("NonSigner", "nonSigner", nonSigner, "nonSigningAmount", nonSigningAmount)
		nonSignersObj = append(nonSignersObj, dataapi.NonSigner{
			OperatorId: nonSigner,
			Count:      nonSigningAmount,
		})
	}

	return &nonSignersObj, nil
}
