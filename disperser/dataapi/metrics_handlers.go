package dataapi

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
)

func (s *server) getMetric(ctx context.Context, startTime int64, endTime int64) (*Metric, error) {
	blockNumber, err := s.transactor.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}
	quorumCount, err := s.transactor.GetQuorumCount(ctx, blockNumber)
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

	return &Metric{
		Throughput:          throughput,
		CostInGas:           costInGas,
		TotalStake:          totalStakePerQuorum[0],
		TotalStakePerQuorum: totalStakePerQuorum,
	}, nil
}

func (s *server) calculateTotalCostGasUsed(ctx context.Context) (float64, error) {
	batches, err := s.subgraphClient.QueryBatchesWithLimit(ctx, 1, 0)
	if err != nil {
		return 0, err
	}

	if len(batches) == 0 {
		return 0, nil
	}

	var (
		totalBlobSize uint
		totalGasUsed  float64
		batch         = batches[0]
	)

	if batch == nil {
		return 0, errors.New("error the latest batch is not valid")
	}

	batchHeaderHash, err := ConvertHexadecimalToBytes(batch.BatchHeaderHash)
	if err != nil {
		s.logger.Error("Failed to convert BatchHeaderHash to hex string: ", "batchHeaderHash", batch.BatchHeaderHash, "err", err)
		return 0, err
	}

	metadatas, err := s.blobstore.GetAllBlobMetadataByBatch(ctx, batchHeaderHash)
	if err != nil {
		s.logger.Error("Failed to get all blob metadata by batch: ", "batchHeaderHash", batchHeaderHash, "err", err)
		return 0, err
	}

	for _, metadata := range metadatas {
		totalBlobSize += metadata.RequestMetadata.BlobSize
	}

	if uint64(totalBlobSize) > 0 {
		totalGasUsed = float64(batch.GasFees.GasUsed) / float64(totalBlobSize)
	}
	return totalGasUsed, nil
}

func (s *server) getNonSigners(ctx context.Context, intervalSeconds int64) (*[]NonSigner, error) {
	nonSigners, err := s.subgraphClient.QueryBatchNonSigningOperatorIdsInInterval(ctx, intervalSeconds)
	if err != nil {
		return nil, err
	}

	nonSignersObj := make([]NonSigner, 0)
	for nonSigner, nonSigningAmount := range nonSigners {
		s.logger.Info("NonSigner", "nonSigner", nonSigner, "nonSigningAmount", nonSigningAmount)
		nonSignersObj = append(nonSignersObj, NonSigner{
			OperatorId: nonSigner,
			Count:      nonSigningAmount,
		})
	}

	return &nonSignersObj, nil
}
