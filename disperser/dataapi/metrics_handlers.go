package dataapi

import (
	"context"
	"errors"
	"time"
)

const (
	gweiMultiplier             = 1_000_000_000
	avgThroughputWindowSize    = 120 // The time window (in seconds) to calculate the data throughput.
	maxWorkersGetOperatorState = 10  // The maximum number of workers to use when querying operator state.
)

func (s *server) getMetric(ctx context.Context, startTime int64, endTime int64, limit int) (*Metric, error) {
	// operators, err := s.subgraphClient.QueryOperatorsWithLimit(ctx, limit)
	// if err != nil {
	// 	return nil, err
	// }

	// blockNumber, err := s.transactor.GetCurrentBlockNumber(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get current block number: %w", err)
	// }
	// totalStake, err := s.calculateTotalStake(operators, blockNumber)
	// if err != nil {
	// 	return nil, err
	// }

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
		TotalStake: 0,
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

// func (s *server) calculateTotalStake(operators []*Operator, blockNumber uint32) (int64, error) {
// 	var (
// 		totalStakeByOperatorChan = make(chan *big.Int, len(operators))
// 		pool                     = workerpool.New(maxWorkersGetOperatorState)
// 	)
//
// 	s.logger.Debug("Number of operators to calculate stake:", "numOperators", len(operators), "blockNumber", blockNumber)
// 	for _, o := range operators {
// 		operatorId, err := ConvertHexadecimalToBytes(o.OperatorId)
// 		if err != nil {
// 			s.logger.Error("Failed to convert operator id to hex string: ", "operatorId", operatorId, "err", err)
// 			return 0, err
// 		}
//
// 		pool.Submit(func() {
// 			operatorState, err := s.chainState.GetOperatorStateByOperator(context.Background(), uint(blockNumber), operatorId)
// 			if err != nil {
// 				s.logger.Error("Failed to get operator state: ", "operatorId", operatorId, "blockNumber", blockNumber, "err", err)
// 				totalStakeByOperatorChan <- big.NewInt(-1)
// 				return
// 			}
// 			totalStake := big.NewInt(0)
// 			s.logger.Debug("Operator state:", "operatorId", operatorId, "num quorums", len(operatorState.Totals))
// 			for quorumId, total := range operatorState.Totals {
// 				s.logger.Debug("Operator stake:", "operatorId", operatorId, "quorum", quorumId, "stake", (*total.Stake).Int64())
// 				totalStake.Add(totalStake, total.Stake)
// 			}
// 			totalStakeByOperatorChan <- totalStake
// 		})
// 	}
//
// 	pool.StopWait()
// 	close(totalStakeByOperatorChan)
//
// 	totalStake := big.NewInt(0)
// 	for total := range totalStakeByOperatorChan {
// 		if total.Int64() == -1 {
// 			return 0, errors.New("error getting operator state")
// 		}
// 		totalStake.Add(totalStake, total)
// 	}
// 	return totalStake.Int64(), nil
// }

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
