package dataapi

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/core"
)

const (
	gweiMultiplier             = 1_000_000_000
	avgThroughputWindowSize    = 120 // The time window (in seconds) to calculate the data throughput.
	maxWorkersGetOperatorState = 10  // The maximum number of workers to use when querying operator state.
)

func (s *server) getMetric(ctx context.Context, startTime int64, endTime int64, limit int) (*Metric, error) {
	blockNumber, err := s.transactor.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}
	operatorState, err := s.chainState.GetOperatorState(ctx, uint(blockNumber), []core.QuorumID{core.QuorumID(0)})
	if err != nil {
		return nil, err
	}
	if len(operatorState.Operators) != 1 {
		return nil, fmt.Errorf("Requesting for one quorum (quorumID=0), but got %v", operatorState.Operators)
	}
	totalStake := big.NewInt(0)
	for _, op := range operatorState.Operators[0] {
		totalStake.Add(totalStake, op.Stake)
	}

	result, err := s.promClient.QueryDisperserBlobSizeBytesPerSecond(ctx, time.Unix(startTime, 0), time.Unix(endTime, 0))
	if err != nil {
		return nil, err
	}

	var (
		totalBytes   float64
		timeDuration float64
		troughput    float64
		valuesSize   = len(result.Values)
	)
	if valuesSize > 1 {
		totalBytes = result.Values[valuesSize-1].Value - result.Values[0].Value
		timeDuration = result.Values[len(result.Values)-1].Timestamp.Sub(result.Values[0].Timestamp).Seconds()
		troughput = totalBytes / timeDuration
	}

	costInWei, err := s.calculateTotalCostGasUsedInWei(ctx)
	if err != nil {
		return nil, err
	}

	return &Metric{
		Throughput: troughput,
		CostInWei:  costInWei,
		TotalStake: totalStake.Uint64(),
	}, nil
}

func (s *server) getThroughput(ctx context.Context, start int64, end int64) ([]*Throughput, error) {
	result, err := s.promClient.QueryDisperserBlobSizeBytesPerSecond(ctx, time.Unix(start, 0), time.Unix(end, 0))
	if err != nil {
		return nil, err
	}

	if len(result.Values) <= 1 {
		return []*Throughput{}, nil
	}

	return calculateAverageThroughput(result.Values, avgThroughputWindowSize), nil
}

func (s *server) calculateTotalCostGasUsedInWei(ctx context.Context) (uint64, error) {
	batches, err := s.subgraphClient.QueryBatchesWithLimit(ctx, 1, 0)
	if err != nil {
		return 0, err
	}

	if len(batches) == 0 {
		return 0, nil
	}

	var (
		totalBlobSize  uint
		totalCostInWei uint64
		batch          = batches[0]
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
		cost := float64(batch.GasFees.GasUsed) / float64(totalBlobSize)
		totalCostInWei = uint64(cost * gweiMultiplier)
	}
	return totalCostInWei, nil
}

func calculateAverageThroughput(values []*PrometheusResultValues, windowSize int64) []*Throughput {
	throughputs := make([]*Throughput, 0)
	totalBytesTransferred := float64(0)
	start := 0
	for i := avgThroughputWindowSize; i < len(values); i++ {
		currentTime := values[i].Timestamp.Unix()

		// The total number of iterations for this loop will be O(N) in aggregate after
		// the outer loop completes, so the amortized cost here is just O(1).
		for start < i && currentTime-values[start+1].Timestamp.Unix() > windowSize {
			start++
		}
		duration := currentTime - values[start].Timestamp.Unix()
		totalBytesTransferred = values[i].Value - values[start].Value
		averageThroughput := totalBytesTransferred / float64(duration)
		throughputs = append(throughputs, &Throughput{
			Timestamp:  uint64(currentTime),
			Throughput: averageThroughput,
		})
	}
	return throughputs
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
