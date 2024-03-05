package dataapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
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

	NumRequests *prometheus.CounterVec
	Latency     *prometheus.SummaryVec

	httpPort string
	logger   common.Logger
}

func NewMetrics(blobMetadataStore *blobstore.BlobMetadataStore, httpPort string, logger common.Logger) *Metrics {
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
		registry: reg,
		httpPort: httpPort,
		logger:   logger,
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
	blobMetadataStore    *blobstore.BlobMetadataStore
	blobStatusMetric     *prometheus.Desc
	scrapeDurationMetric *prometheus.GaugeVec
	logger               common.Logger
}

func NewDynamoDBCollector(blobMetadataStore *blobstore.BlobMetadataStore, logger common.Logger) *DynamoDBCollector {
	if blobMetadataStore == nil {
		logger.Error("BlobMetadataStore is nil, metrics will not be collected")
	}

	collector := &DynamoDBCollector{
		blobMetadataStore: blobMetadataStore,
		blobStatusMetric: prometheus.NewDesc("dynamodb_blob_metadata_status_count",
			"Number of blobs with specific status in DynamoDB",
			[]string{"status"},
			nil,
		),
		scrapeDurationMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dynamodb_collector_scrape_duration_seconds",
			Help: "Gauge of scrape duration for DynamoDB collector",
		}, []string{}),
		logger: logger,
	}

	return collector
}

func (collector *DynamoDBCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.blobStatusMetric
	collector.scrapeDurationMetric.Describe(ch)
}

func (collector *DynamoDBCollector) Collect(ch chan<- prometheus.Metric) {
	// Record the start time of the scrape
	startTime := time.Now()

	for _, status := range []disperser.BlobStatus{
		disperser.Processing,
		disperser.Confirmed,
		disperser.Failed,
		disperser.InsufficientSignatures,
	} {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		count, err := collector.getBlobMetadataByStatus(ctx, status)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				collector.logger.Error("Fetching blob metadata by status took longer than 60 seconds", "status", status)
			} else {
				collector.logger.Error("Failed to get count of blob metadata by status", "status", status, "err", err)
			}
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			collector.blobStatusMetric,
			prometheus.GaugeValue,
			float64(count),
			status.String(),
		)
	}

	// Record the scrape duration
	duration := time.Since(startTime).Seconds()
	collector.scrapeDurationMetric.WithLabelValues().Set(duration)
	collector.scrapeDurationMetric.Collect(ch)
}

// getBlobMetadataByStatus fetches the count of blob metadata by status from DynamoDB.
// It uses pagination to fetch all the metadata by status and returns the total count.
func (collector *DynamoDBCollector) getBlobMetadataByStatus(ctx context.Context, status disperser.BlobStatus) (int, error) {
	totalMetadata := 0

	metadatas, exclusiveStartKey, err := collector.blobMetadataStore.GetBlobMetadataByStatusWithPagination(ctx, status, 1000, nil)
	if err != nil {
		collector.logger.Error("failed to get blob metadata by status with pagination", "status", status.String(), "err", err)
		return 0, err
	}
	totalMetadata += len(metadatas) // Count the first batch of metadata

	for exclusiveStartKey != nil {
		metadatas, exclusiveStartKey, err = collector.blobMetadataStore.GetBlobMetadataByStatusWithPagination(ctx, status, 1000, exclusiveStartKey)
		if err != nil {
			collector.logger.Error("failed to get blob metadata by status with pagination in loop", "status", status.String(), "err", err)
			return totalMetadata, err
		}

		totalMetadata += len(metadatas)
	}

	return totalMetadata, nil
}
