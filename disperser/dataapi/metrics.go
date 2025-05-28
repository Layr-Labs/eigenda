package dataapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
	commonv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
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
	CacheHitsTotal *prometheus.CounterVec
	Latency        *prometheus.SummaryVec
	OperatorsStake *prometheus.GaugeVec

	// Cache metrics in v2
	BatchFeedCacheMetrics *FeedCacheMetrics

	Semvers                *prometheus.GaugeVec
	SemversStakePctQuorum0 *prometheus.GaugeVec
	SemversStakePctQuorum1 *prometheus.GaugeVec
	SemversStakePctQuorum2 *prometheus.GaugeVec

	httpPort string
	logger   logging.Logger
}

func NewMetrics(serverVersion uint, reg *prometheus.Registry, blobMetadataStore interface{}, httpPort string, logger logging.Logger) *Metrics {
	namespace := "eigenda_dataapi"
	if reg == nil {
		reg = prometheus.NewRegistry()
	}

	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())
	if serverVersion == 1 {
		if store, ok := blobMetadataStore.(*blobstore.BlobMetadataStore); ok {
			reg.MustRegister(NewDynamoDBCollector(store, logger))
		}
	} else if serverVersion == 2 {
		if store, ok := blobMetadataStore.(*blobstorev2.BlobMetadataStore); ok {
			reg.MustRegister(NewBlobMetadataStoreV2Collector(store, logger))
		}
	}
	metrics := &Metrics{
		NumRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "requests",
				Help:      "the number of requests",
			},
			[]string{"status", "method"},
		),
		// Cache hit rate for an API is CacheHitsTotal["method_foo"] / NumRequests["success"]["method_foo"]
		CacheHitsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "cache_hits_total",
				Help:      "the number of requests that hit the cache",
			},
			[]string{"method"},
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
		BatchFeedCacheMetrics: NewFeedCacheMetrics("batch_feed", reg),
		registry:              reg,
		httpPort:              httpPort,
		logger:                logger.With("component", "DataAPIMetrics"),
	}
	return metrics
}

// ObserveLatency observes the latency of a stage in 'stage
func (g *Metrics) ObserveLatency(method string, duration time.Duration) {
	g.Latency.WithLabelValues(method).Observe(float64(duration.Milliseconds()))
}

// IncrementCacheHit increments the number of requests that hit cache
func (g *Metrics) IncrementCacheHit(method string) {
	g.CacheHitsTotal.With(prometheus.Labels{
		"method": method,
	}).Inc()
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

// IncrementInvalidArgdRequestNum increments the number of failed requests with invalid args
func (g *Metrics) IncrementInvalidArgRequestNum(method string) {
	g.NumRequests.With(prometheus.Labels{
		"status": "invalid_args",
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

// BlobStatusMetrics holds the metrics for a specific blob status
type BlobStatusMetrics struct {
	gauge        prometheus.Gauge
	currentValue float64
}

// BlobMetadataStoreV2Collector collects metrics from the blob metadata store.
type BlobMetadataStoreV2Collector struct {
	blobMetadataStore *blobstorev2.BlobMetadataStore
	statusMetrics     map[commonv2.BlobStatus]*BlobStatusMetrics
	logger            logging.Logger
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewBlobMetadataStoreV2Collector(blobMetadataStore *blobstorev2.BlobMetadataStore, logger logging.Logger) *BlobMetadataStoreV2Collector {
	ctx, cancel := context.WithCancel(context.Background())
	collector := &BlobMetadataStoreV2Collector{
		blobMetadataStore: blobMetadataStore,
		statusMetrics:     make(map[commonv2.BlobStatus]*BlobStatusMetrics),
		logger:            logger,
		ctx:               ctx,
		cancel:            cancel,
	}

	// Create a gauge for each possible status (that is not terminal)
	for _, status := range []commonv2.BlobStatus{
		commonv2.Queued,
		commonv2.Encoded,
		commonv2.GatheringSignatures,
	} {
		gauge := prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "eigenda_blob_metadata_v2_status_count",
				Help: "Current number of blobs in this status. In case of timeouts or failures when querying the blob metadata store (e.g. when there are too many blobs), the last known value will be returned as stale data.",
				ConstLabels: prometheus.Labels{
					"status": status.String(),
				},
			},
		)
		prometheus.MustRegister(gauge)
		collector.statusMetrics[status] = &BlobStatusMetrics{
			gauge:        gauge,
			currentValue: 0,
		}
	}

	// Do initial count
	collector.updateCounts(context.Background())

	return collector
}

// countBlobsWithStatus counts blobs for a specific status with pagination and timeout handling
func (collector *BlobMetadataStoreV2Collector) countBlobsWithStatus(ctx context.Context, status commonv2.BlobStatus) (int32, error) {
	var totalCount int32
	var cursor *blobstorev2.StatusIndexCursor

	for {
		select {
		case <-ctx.Done():
			return totalCount, ctx.Err()
		default:
			blobs, nextCursor, err := collector.blobMetadataStore.GetBlobMetadataByStatusPaginated(ctx, status, cursor, 100)
			if err != nil {
				return totalCount, err
			}

			count := int32(len(blobs))
			totalCount += count

			collector.logger.Debug("Got partial count for status",
				"status", status.String(),
				"partial_count", count,
				"running_total", totalCount,
				"has_more", nextCursor != nil,
			)

			if count == 0 || nextCursor == nil {
				return totalCount, nil
			}
			cursor = nextCursor
		}
	}
}

func (collector *BlobMetadataStoreV2Collector) updateCounts(ctx context.Context) {
	collector.logger.Debug("Starting blob status count update")
	startTime := time.Now()

	for status, metrics := range collector.statusMetrics {
		statusCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		totalCount, err := collector.countBlobsWithStatus(statusCtx, status)
		defer cancel()

		if err != nil {
			collector.logger.Error("Failed to get count of blob metadata by status - using stale data",
				"status", status,
				"err", err,
				"current_count", metrics.currentValue,
			)
			continue // Keep using the last known value
		}

		metrics.gauge.Set(float64(totalCount))
		metrics.currentValue = float64(totalCount)

		collector.logger.Debug("Updated blob status count",
			"status", status.String(),
			"count", totalCount,
		)
	}

	collector.logger.Debug("Completed blob status count update",
		"duration_ms", time.Since(startTime).Milliseconds(),
	)
}

func (collector *BlobMetadataStoreV2Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metrics := range collector.statusMetrics {
		ch <- metrics.gauge.Desc()
	}
}

func (collector *BlobMetadataStoreV2Collector) Collect(ch chan<- prometheus.Metric) {
	collector.logger.Debug("Prometheus scrape triggered, updating counts")
	startTime := time.Now()

	// Create a context with timeout for the entire collection.
	// The default scrape timeout is 10 seconds so we set it to 8 seconds to allow for some time for the collection.
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	// Try to get fresh counts
	collector.updateCounts(ctx)

	// Send current gauge values (either fresh or stale)
	for _, metrics := range collector.statusMetrics {
		ch <- metrics.gauge
	}

	collector.logger.Debug("Completed blob metadata store v2 collector scrape",
		"duration_ms", time.Since(startTime).Milliseconds(),
	)
}
