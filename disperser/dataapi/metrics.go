package dataapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
	"github.com/Layr-Labs/eigenda/operators"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsConfig struct {
	HTTPPort      string
	EnableMetrics bool
}

type Metrics struct {
	registry *prometheus.Registry

	NumRequests    *prometheus.CounterVec
	Latency        *prometheus.SummaryVec
	OperatorsStake *prometheus.GaugeVec

	Semvers                *prometheus.GaugeVec
	SemversStakePctQuorum0 *prometheus.GaugeVec
	SemversStakePctQuorum1 *prometheus.GaugeVec
	SemversStakePctQuorum2 *prometheus.GaugeVec

	httpPort string
	logger   logging.Logger
}

func NewMetrics(blobMetadataStore *blobstore.BlobMetadataStore, httpPort string, logger logging.Logger) *Metrics {
	namespace := "eigenda_dataapi"
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(NewDynamoDBCollector(blobMetadataStore, logger))
	metrics := &Metrics{
		NumRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "requests",
				Help:      "the number of requests",
			},
			[]string{"status", "method"},
		),
		Latency: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Name:       "latency_ms",
				Help:       "latency summary in milliseconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			},
			[]string{"method"},
		),
		Semvers: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "node_semvers",
				Help: "Node semver install base",
			},
			[]string{"semver"},
		),
		SemversStakePctQuorum0: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "node_semvers_stake_pct_quorum_0",
				Help: "Node semver stake percentage in quorum 0",
			},
			[]string{"semver_stake_pct_quorum_0"},
		),
		SemversStakePctQuorum1: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "node_semvers_stake_pct_quorum_1",
				Help: "Node semver stake percentage in quorum 1",
			},
			[]string{"semver_stake_pct_quorum_1"},
		),
		SemversStakePctQuorum2: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "node_semvers_stake_pct_quorum_2",
				Help: "Node semver stake percentage in quorum 2",
			},
			[]string{"semver_stake_pct_quorum_2"},
		),
		OperatorsStake: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "operators_stake",
				Help:      "the sum of stake percentages of top N operators",
			},
			// The "quorum" can be: total, 0, 1, ...
			// The "topn" can be: 1, 2, 3, 5, 8, 10
			[]string{"quorum", "topn"},
		),
		registry: reg,
		httpPort: httpPort,
		logger:   logger.With("component", "DataAPIMetrics"),
	}
	return metrics
}

// ObserveLatency observes the latency of a stage in 'stage
func (g *Metrics) ObserveLatency(method string, latencyMs float64) {
	g.Latency.WithLabelValues(method).Observe(latencyMs)
}

// IncrementSuccessfulRequestNum increments the number of successful requests
func (g *Metrics) IncrementSuccessfulRequestNum(method string) {
	g.NumRequests.With(prometheus.Labels{
		"status": "success",
		"method": method,
	}).Inc()
}

// IncrementFailedRequestNum increments the number of failed requests
func (g *Metrics) IncrementFailedRequestNum(method string) {
	g.NumRequests.With(prometheus.Labels{
		"status": "failed",
		"method": method,
	}).Inc()
}

// IncrementNotFoundRequestNum increments the number of not found requests
func (g *Metrics) IncrementNotFoundRequestNum(method string) {
	g.NumRequests.With(prometheus.Labels{
		"status": "not found",
		"method": method,
	}).Inc()
}

// UpdateSemverMetrics updates the semver metrics
func (g *Metrics) UpdateSemverCounts(semverData map[string]*semver.SemverMetrics) {
	for semver, metrics := range semverData {
		g.Semvers.WithLabelValues(semver).Set(float64(metrics.Operators))
		for quorum, stakePct := range metrics.QuorumStakePercentage {
			switch quorum {
			case 0:
				g.SemversStakePctQuorum0.WithLabelValues(semver).Set(stakePct)
			case 1:
				g.SemversStakePctQuorum1.WithLabelValues(semver).Set(stakePct)
			case 2:
				g.SemversStakePctQuorum2.WithLabelValues(semver).Set(stakePct)
			default:
				g.logger.Error("Unable to log semver quorum stake percentage for quorum", "semver", semver, "quorum", quorum, "stake", stakePct)
			}
		}
	}
}

func (g *Metrics) updateStakeMetrics(rankedOperators []*operators.OperatorStakeShare, label string) {
	indices := []int{0, 1, 2, 4, 7, 9}
	accuStake := float64(0)
	idx := 0
	for i, op := range rankedOperators {
		accuStake += op.StakeShare
		if idx < len(indices) && i == indices[idx] {
			g.OperatorsStake.WithLabelValues(label, fmt.Sprintf("%d", i+1)).Set(accuStake / 100)
			idx++
		}
	}
}

func (g *Metrics) UpdateOperatorsStake(totalRanked []*operators.OperatorStakeShare, quorumRanked map[uint8][]*operators.OperatorStakeShare) {
	g.updateStakeMetrics(totalRanked, "total")
	for q, operators := range quorumRanked {
		g.updateStakeMetrics(operators, fmt.Sprintf("%d", q))
	}
}

// Start starts the metrics server
func (g *Metrics) Start(ctx context.Context) {
	g.logger.Info("Starting metrics server at ", "port", g.httpPort)
	addr := fmt.Sprintf(":%s", g.httpPort)
	go func() {
		log := g.logger
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(
			g.registry,
			promhttp.HandlerOpts{},
		))
		err := http.ListenAndServe(addr, mux)
		log.Error("Prometheus server failed", "err", err)
	}()
}

type DynamoDBCollector struct {
	blobMetadataStore *blobstore.BlobMetadataStore
	blobStatusMetric  *prometheus.Desc
	logger            logging.Logger
}

func NewDynamoDBCollector(blobMetadataStore *blobstore.BlobMetadataStore, logger logging.Logger) *DynamoDBCollector {
	return &DynamoDBCollector{
		blobMetadataStore: blobMetadataStore,
		blobStatusMetric: prometheus.NewDesc("dynamodb_blob_metadata_status_count",
			"Number of blobs with specific status in DynamoDB",
			[]string{"status"},
			nil,
		),
		logger: logger,
	}
}

func (collector *DynamoDBCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.blobStatusMetric
}

func (collector *DynamoDBCollector) Collect(ch chan<- prometheus.Metric) {
	count, err := collector.blobMetadataStore.GetBlobMetadataCountByStatus(context.Background(), disperser.Processing)
	if err != nil {
		collector.logger.Error("failed to get count of blob metadata by status", "err", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		collector.blobStatusMetric,
		prometheus.GaugeValue,
		float64(count),
		disperser.Processing.String(),
	)
}
