package node

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/operators"
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

	// Rank of the operator in a particular registered quorum.
	RegisteredQuorumsRank *prometheus.GaugeVec
	// Stake share of the operator in a particular registered quorum.
	RegisteredQuorumsStakeShare *prometheus.GaugeVec
	// Accumulated number of RPC requests received.
	AccNumRequests *prometheus.CounterVec
	// The latency (in ms) to process the request.
	RequestLatency *prometheus.SummaryVec
	// Accumulated number and size of batches processed by their statuses.
	AccuBatches *prometheus.CounterVec
	// Accumulated number and size of batches that have been removed from the Node.
	AccuRemovedBatches *prometheus.CounterVec
	// Accumulated number and size of blobs that have been removed from the Node.
	AccuRemovedBlobs *prometheus.CounterVec
	// Accumulated number and size of blobs processed by quorums.
	AccuBlobs *prometheus.CounterVec
	// Total number of changes in the node's socket address.
	AccuSocketUpdates prometheus.Counter
	// avs node spec eigen_ metrics: https://eigen.nethermind.io/docs/spec/metrics/metrics-prom-spec
	EigenMetrics eigenmetrics.Metrics
	// Reachability gauge to monitoring the reachability of the node's retrieval/dispersal sockets
	ReachabilityGauge *prometheus.GaugeVec
	// The throughput (bytes per second) at which the data is written to database.
	DBWriteThroughput prometheus.Gauge

	registry *prometheus.Registry
	// socketAddr is the address at which the metrics server will be listening.
	// should be in format ip:port
	socketAddr             string
	operatorId             core.OperatorID
	onchainMetricsInterval int64
	tx                     core.Reader
	chainState             core.ChainState
	allQuorumCache         map[core.QuorumID]bool
}

func NewMetrics(eigenMetrics eigenmetrics.Metrics, reg *prometheus.Registry, logger logging.Logger, socketAddr string, operatorId core.OperatorID, onchainMetricsInterval int64, tx core.Reader, chainState core.ChainState) *Metrics {

	// Add Go module collectors
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	metrics := &Metrics{
		RegisteredQuorumsRank: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "registered_quorums_rank",
				Help:      "the rank of operator by TVL in that quorum (1 being the highest)",
			},
			[]string{"quorum"},
		),
		RegisteredQuorumsStakeShare: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "registered_quorums_stake_share",
				Help:      "the stake share of operator in basis points in that quorum",
			},
			[]string{"quorum"},
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
		AccuRemovedBlobs: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_removed_blobs_total",
				Help:      "the total number and size of blobs that have been removed by the DA node",
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
		ReachabilityGauge: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "reachability_status",
				Help:      "the reachability status of the nodes retrievel/dispersal sockets",
			},
			[]string{"service"},
		),
		DBWriteThroughput: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "db_write_throughput_bytes_per_second",
				Help:      "the throughput (bytes per second) at which the data is written to database",
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
		allQuorumCache:         make(map[core.QuorumID]bool),
	}

	return metrics
}

func (g *Metrics) Start() {
	_ = g.EigenMetrics.Start(context.Background(), g.registry)

	if g.onchainMetricsInterval > 0 {
		go g.collectOnchainMetrics()
	}
}

func (g *Metrics) RecordRPCRequest(method string, status string, duration time.Duration) {
	g.AccNumRequests.WithLabelValues(method, status).Inc()
	g.ObserveLatency(method, "total", float64(duration.Milliseconds()))
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

func (g *Metrics) RemoveNBlobs(numBlobs int, totalSize int64) {
	for i := 0; i < numBlobs; i++ {
		g.AccuRemovedBlobs.WithLabelValues("number").Inc()
	}
	g.AccuRemovedBatches.WithLabelValues("size").Add(float64(totalSize))
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

func (g *Metrics) RecordStoreChunksStage(stage string, dataSize uint64, latency time.Duration) {
	g.AcceptBatches(stage, dataSize)
	g.ObserveLatency("StoreChunks", stage, float64(latency.Milliseconds()))
}

func (g *Metrics) collectOnchainMetrics() {
	ticker := time.NewTicker(time.Duration(g.onchainMetricsInterval) * time.Second)
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
			g.ResetQuorumMetrics(blockNum)
			g.logger.Warn("This node is currently not in any quorum", "blockNumber", blockNum, "operatorId", g.operatorId.Hex())
			continue
		}
		state, err := g.chainState.GetOperatorState(ctx, uint(blockNum), quorumIds)
		if err != nil {
			g.logger.Error("Failed to query chain RPC for operator state", "blockNumber", blockNum, "quorumIds", quorumIds, "err", err)
			continue
		}
		_, quorumRankedOperators := operators.GetRankedOperators(state)
		for q := range state.Operators {
			for i, op := range quorumRankedOperators[q] {
				if op.OperatorId == g.operatorId {
					g.allQuorumCache[q] = true
					g.RegisteredQuorumsStakeShare.WithLabelValues(fmt.Sprintf("%d", q)).Set(op.StakeShare)
					g.RegisteredQuorumsRank.WithLabelValues(fmt.Sprintf("%d", q)).Set(float64(i + 1))
					g.logger.Info("Current operator registration onchain", "operatorId", g.operatorId.Hex(), "blockNumber", blockNum, "quorumId", q, "stakeShare (basis point)", op.StakeShare, "rank", i+1)
					break
				}
			}
		}
		// Check if operator deregistered for an existing quorum, set the stake share and rank to 0
		g.ResetQuorumMetrics(blockNum)
	}
}

func (g *Metrics) ResetQuorumMetrics(blockNum uint32) {
	// Check if operator deregistered for an existing quorum, set the stake share and rank to 0
	for q := range g.allQuorumCache {
		// If this quorum was deregistered then set the stake share and rank to 0
		if !g.allQuorumCache[q] {
			g.RegisteredQuorumsStakeShare.WithLabelValues(fmt.Sprintf("%d", q)).Set(0)
			g.RegisteredQuorumsRank.WithLabelValues(fmt.Sprintf("%d", q)).Set(0)
			g.logger.Info("Current operator deregistration onchain", "operatorId", g.operatorId.Hex(), "blockNumber", blockNum, "quorumId", q)
		}
		// Reset the cache to false for all quorum for next cycle
		g.allQuorumCache[q] = false
	}
}
