package node

import (
	"context"
	"math/big"
	"sort"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	eigenmetrics "github.com/Layr-Labs/eigensdk-go/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	Namespace = "node"
)

type Metrics struct {
	logger logging.Logger

	// The quorums the node is registered.
	RegisteredQuorums *prometheus.GaugeVec
	// Accumulated number of RPC requests received.
	AccNumRequests *prometheus.CounterVec
	// The latency (in ms) to process the request.
	RequestLatency *prometheus.SummaryVec
	// Accumulated number and size of batches processed by their statuses.
	AccuBatches *prometheus.CounterVec
	// Accumulated number and size of batches that have been removed from the Node.
	AccuRemovedBatches *prometheus.CounterVec
	// Accumulated number and size of blobs processed by quorums.
	AccuBlobs *prometheus.CounterVec
	// Total number of changes in the node's socket address.
	AccuSocketUpdates prometheus.Counter
	// avs node spec eigen_ metrics: https://eigen.nethermind.io/docs/spec/metrics/metrics-prom-spec
	EigenMetrics eigenmetrics.Metrics

	registry *prometheus.Registry
	// socketAddr is the address at which the metrics server will be listening.
	// should be in format ip:port
	socketAddr             string
	operatorId             core.OperatorID
	onchainMetricsInterval int64
	tx                     core.Transactor
	chainState             core.ChainState
}

func NewMetrics(eigenMetrics eigenmetrics.Metrics, reg *prometheus.Registry, logger logging.Logger, socketAddr string, operatorId core.OperatorID, onchainMetricsInterval int64, tx core.Transactor, chainState core.ChainState) *Metrics {

	// Add Go module collectors
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	metrics := &Metrics{
		// The "type" label have values: stake_share, rank. The "stake_share" is stake share (in basis point),
		// and the "rank" is operator's ranking (the operator with highest amount of stake ranked as 1) by stake share in the quorum.
		RegisteredQuorums: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "registered",
				Help:      "the quorums the DA node is registered",
			},
			[]string{"quorum", "type"},
		),
		// The "status" label has values: success, failure.
		AccNumRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_rpc_requests_total",
				Help:      "the total number of requests processed by the DA node",
			},
			[]string{"method", "status"},
		),
		RequestLatency: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  Namespace,
				Name:       "request_latency_ms",
				Help:       "latency summary in milliseconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			},
			[]string{"method", "stage"},
		),
		AccuBlobs: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_blobs_total",
				Help:      "the total number and size of blobs processed by the DA node",
			},
			[]string{"type", "quorum"},
		),
		// The "status" label has values: received, validated, stored, signed.
		// These are the lifecycle of a batch at the DA Node.
		AccuBatches: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_processed_batches_total",
				Help:      "the total number and size of batches processed by the DA node",
			},
			[]string{"type", "status"},
		),
		AccuRemovedBatches: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_removed_batches_total",
				Help:      "the total number and size of batches that have been removed by the DA node",
			},
			[]string{"type"},
		),
		AccuSocketUpdates: promauto.With(reg).NewCounter(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_node_socket_updates_total",
				Help:      "the total number of node's socket address updates",
			},
		),
		EigenMetrics:           eigenMetrics,
		logger:                 logger.With("component", "NodeMetrics"),
		registry:               reg,
		socketAddr:             socketAddr,
		operatorId:             operatorId,
		onchainMetricsInterval: onchainMetricsInterval,
		tx:                     tx,
		chainState:             chainState,
	}

	return metrics
}

func (g *Metrics) Start() {
	_ = g.EigenMetrics.Start(context.Background(), g.registry)

	if g.onchainMetricsInterval > 0 {
		go g.collectOnchainMetrics()
	}
}

func (g *Metrics) RecordRPCRequest(method string, status string) {
	g.AccNumRequests.WithLabelValues(method, status).Inc()
}

func (g *Metrics) RecordSocketAddressChange() {
	g.AccuSocketUpdates.Inc()
}

func (g *Metrics) ObserveLatency(method, stage string, latencyMs float64) {
	g.RequestLatency.WithLabelValues(method, stage).Observe(latencyMs)
}

func (g *Metrics) RemoveNCurrentBatch(numBatches int, totalBatchSize int64) {
	for i := 0; i < numBatches; i++ {
		g.AccuRemovedBatches.WithLabelValues("number").Inc()
	}
	g.AccuRemovedBatches.WithLabelValues("size").Add(float64(totalBatchSize))
}

func (g *Metrics) AcceptBlobs(quorumId core.QuorumID, blobSize uint64) {
	quorum := strconv.Itoa(int(quorumId))
	g.AccuBlobs.WithLabelValues("number", quorum).Inc()
	g.AccuBlobs.WithLabelValues("size", quorum).Add(float64(blobSize))
}

func (g *Metrics) AcceptBatches(status string, batchSize uint64) {
	g.AccuBatches.WithLabelValues("number", status).Inc()
	g.AccuBatches.WithLabelValues("size", status).Add(float64(batchSize))
}

func (g *Metrics) collectOnchainMetrics() {
	ticker := time.NewTicker(time.Duration(uint64(g.onchainMetricsInterval)))
	defer ticker.Stop()

	// 3 chain RPC calls in each cycle.
	for range ticker.C {
		ctx := context.Background()
		blockNum, err := g.tx.GetCurrentBlockNumber(ctx)
		if err != nil {
			g.logger.Error("Failed to query chain RPC for current block number", "err", err)
			continue
		}
		bitmaps, err := g.tx.GetQuorumBitmapForOperatorsAtBlockNumber(ctx, []core.OperatorID{g.operatorId}, blockNum)
		if err != nil {
			g.logger.Error("Failed to query chain RPC for quorum bitmap", "blockNumber", blockNum, "err", err)
			continue
		}
		quorumIds := eth.BitmapToQuorumIds(bitmaps[0])
		if len(quorumIds) == 0 {
			g.logger.Warn("This node is currently not in any quorum", "blockNumber", blockNum, "operatorId", g.operatorId.Hex())
			continue
		}
		state, err := g.chainState.GetOperatorState(ctx, uint(blockNum), quorumIds)
		if err != nil {
			g.logger.Error("Failed to query chain RPC for operator state", "blockNumber", blockNum, "quorumIds", quorumIds, "err", err)
			continue
		}
		type OperatorStakeShare struct {
			operatorId core.OperatorID
			stakeShare float64
		}
		for q, operators := range state.Operators {
			operatorStakeShares := make([]*OperatorStakeShare, 0)
			for opId, opInfo := range operators {
				share, _ := new(big.Int).Div(new(big.Int).Mul(opInfo.Stake, big.NewInt(10000)), state.Totals[q].Stake).Float64()
				operatorStakeShares = append(operatorStakeShares, &OperatorStakeShare{operatorId: opId, stakeShare: share})
			}
			// Descending order by stake share in the quorum.
			sort.Slice(operatorStakeShares, func(i, j int) bool {
				if operatorStakeShares[i].stakeShare == operatorStakeShares[j].stakeShare {
					return operatorStakeShares[i].operatorId.Hex() < operatorStakeShares[j].operatorId.Hex()
				}
				return operatorStakeShares[i].stakeShare > operatorStakeShares[j].stakeShare
			})
			for i, op := range operatorStakeShares {
				if op.operatorId == g.operatorId {
					g.RegisteredQuorums.WithLabelValues(string(q), "stake_share").Set(op.stakeShare)
					g.RegisteredQuorums.WithLabelValues(string(q), "rank").Set(float64(i + 1))
					g.logger.Info("Current operator registration onchain", "operatorId", g.operatorId.Hex(), "blockNumber", blockNum, "quorumId", q, "stakeShare (basis point)", op.stakeShare, "rank", i+1)
					break
				}
			}
		}
	}
}
