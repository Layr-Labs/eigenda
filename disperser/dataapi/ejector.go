package dataapi

import (
	"context"
	"sort"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// The caller should ensure "stakeShare" is in range (0, 1].
func stakeShareToSLA(stakeShare float64) float64 {
	switch {
	case stakeShare > 0.1:
		return 0.975
	case stakeShare > 0.05:
		return 0.95
	default:
		return 0.9
	}
}

// operatorPerfScore scores an operator based on its stake share and nonsigning rate. The
// performance score will be in range [0, 1], with higher score indicating better performance.
func operatorPerfScore(stakeShare float64, nonsigningRate float64) float64 {
	if nonsigningRate == 0 {
		return 1.0
	}
	sla := stakeShareToSLA(stakeShare)
	perf := (1 - sla) / nonsigningRate
	return perf / (1.0 + perf)
}

func computePerfScore(metric *OperatorNonsigningPercentageMetrics) float64 {
	return operatorPerfScore(metric.StakePercentage, metric.Percentage/100.0)
}

type ejector struct {
	logger     logging.Logger
	transactor core.Transactor
	metrics    *Metrics

	// For serializing the ejection requests.
	mu sync.Mutex
}

func newEjector(logger logging.Logger, tx core.Transactor, metrics *Metrics) *ejector {
	return &ejector{
		logger:     logger.With("component", "Ejector"),
		transactor: tx,
		metrics:    metrics,
	}
}

func (e *ejector) eject(ctx context.Context, nonsigningRate *OperatorsNonsigningPercentage, mode string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	nonsigners := make([]*OperatorNonsigningPercentageMetrics, 0)
	for _, metric := range nonsigningRate.Data {
		// Collect only the nonsigners who violate the SLA.
		if metric.Percentage/100.0 > 1-stakeShareToSLA(metric.StakePercentage) {
			nonsigners = append(nonsigners, metric)
		}
	}

	// Rank the operators for each quorum by the operator performance score.
	// The operators with lower perf score will get ejected with priority in case of
	// rate limiting.
	sort.Slice(nonsigners, func(i, j int) bool {
		if nonsigners[i].QuorumId == nonsigners[j].QuorumId {
			if computePerfScore(nonsigners[i]) == computePerfScore(nonsigners[j]) {
				return float64(nonsigners[i].TotalUnsignedBatches)*nonsigners[i].StakePercentage > float64(nonsigners[j].TotalUnsignedBatches)*nonsigners[j].StakePercentage
			}
			return computePerfScore(nonsigners[i]) < computePerfScore(nonsigners[j])
		}
		return nonsigners[i].QuorumId < nonsigners[j].QuorumId
	})

	operatorsByQuorum, err := e.convertOperators(nonsigners)
	if err != nil {
		return err
	}

	receipt, err := e.transactor.EjectOperators(ctx, operatorsByQuorum)
	if err != nil {
		e.logger.Error("Ejection transaction failed", "err", err)
		return err
	}
	e.logger.Info("Ejection transaction succeeded", "receipt", receipt)

	e.metrics.UpdateEjectionGasUsed(receipt.GasUsed)

	// TODO: get the txn response and update the metrics.

	return nil
}

func (e *ejector) convertOperators(nonsigners []*OperatorNonsigningPercentageMetrics) ([][]core.OperatorID, error) {
	var maxQuorumId uint8
	for _, metric := range nonsigners {
		if metric.QuorumId > maxQuorumId {
			maxQuorumId = metric.QuorumId
		}
	}

	numOperatorByQuorum := make(map[uint8]int)
	stakeShareByQuorum := make(map[uint8]float64)

	result := make([][]core.OperatorID, maxQuorumId+1)
	for _, metric := range nonsigners {
		id, err := core.OperatorIDFromHex(metric.OperatorId)
		if err != nil {
			return nil, err
		}
		result[metric.QuorumId] = append(result[metric.QuorumId], id)
		numOperatorByQuorum[metric.QuorumId]++
		stakeShareByQuorum[metric.QuorumId] += metric.StakePercentage
	}
	e.metrics.UpdateRequestedOperatorMetric(numOperatorByQuorum, stakeShareByQuorum)

	return result, nil
}
