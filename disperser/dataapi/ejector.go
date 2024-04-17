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

type Ejector struct {
	Logger     logging.Logger
	Transactor core.Transactor
	Metrics    *Metrics

	// For serializing the ejection requests.
	mu sync.Mutex
}

func NewEjector(logger logging.Logger, tx core.Transactor, metrics *Metrics) *Ejector {
	return &Ejector{
		Logger:     logger.With("component", "Ejector"),
		Transactor: tx,
		Metrics:    metrics,
	}
}

func (e *Ejector) eject(ctx context.Context, nonsigningRate *OperatorsNonsigningPercentage, mode string) error {
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
			return computePerfScore(nonsigners[i]) < computePerfScore(nonsigners[j])
		}
		return nonsigners[i].QuorumId < nonsigners[j].QuorumId
	})

	operatorsByQuorum, err := e.convertOperators(nonsigners)
	if err != nil {
		return err
	}

	_, err = e.Transactor.EjectOperators(ctx, operatorsByQuorum)
	return err

	// TODO: get the txn response and update the metrics.
}

func (e *Ejector) convertOperators(nonsigners []*OperatorNonsigningPercentageMetrics) ([][]core.OperatorID, error) {
	var maxQuorumId uint8
	for _, metric := range nonsigners {
		if metric.QuorumId > maxQuorumId {
			maxQuorumId = metric.QuorumId
		}
	}

	result := make([][]core.OperatorID, maxQuorumId+1)
	for _, metric := range nonsigners {
		id, err := core.OperatorIDFromHex(metric.OperatorId)
		if err != nil {
			return nil, err
		}
		result[metric.QuorumId] = append(result[metric.QuorumId], id)
	}

	return result, nil
}
