package grpc

import (
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

const namespace = "eigenda_node"

// V2Metrics encapsulates metrics for the v2 DA node.
type V2Metrics struct {
	logger logging.Logger

	registry         *prometheus.Registry
	server           *http.Server
	grpcServerOption grpc.ServerOption

	storeChunksLatency  *prometheus.SummaryVec
	storeChunksDataSize *prometheus.GaugeVec

	getChunksLatency  *prometheus.SummaryVec
	getChunksDataSize *prometheus.GaugeVec

	dbSize           *prometheus.GaugeVec
	dbSizePollPeriod time.Duration
	dbDir            string
	isAlive          *atomic.Bool
}

// NewV2Metrics creates a new V2Metrics instance. dbSizePollPeriod is the period at which the database size is polled.
// If set to 0, the database size is not polled.
func NewV2Metrics(
	logger logging.Logger,
	port int,
	dbDir string,
	dbSizePollPeriod time.Duration) (*V2Metrics, error) {

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	logger.Infof("Starting metrics server at port %d", port)
	addr := fmt.Sprintf(":%d", port)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{},
	))
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	grpcMetrics := grpcprom.NewServerMetrics()
	registry.MustRegister(grpcMetrics)
	grpcServerOption := grpc.UnaryInterceptor(
		grpcMetrics.UnaryServerInterceptor(),
	)

	storeChunksLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "store_chunks_latency_ms",
			Help:       "The latency of a StoreChunks() RPC call.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{},
	)

	storeChunksDataSize := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "store_chunks_data_size_bytes",
			Help:      "The size of the data requested to be stored by StoreChunks() RPC calls.",
		},
		[]string{},
	)

	getChunksLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_chunks_latency_ms",
			Help:       "The latency of a GetChunks() RPC call.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{},
	)

	getChunksDataSize := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "get_chunks_data_size_bytes",
			Help:      "The size of the data requested to be retrieved by GetChunks() RPC calls.",
		},
		[]string{},
	)

	dbSize := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_size_bytes",
			Help:      "The size of the leveldb database.",
		},
		[]string{},
	)
	isAlive := &atomic.Bool{}
	isAlive.Store(true)

	return &V2Metrics{
		logger:              logger,
		registry:            registry,
		server:              server,
		grpcServerOption:    grpcServerOption,
		storeChunksLatency:  storeChunksLatency,
		storeChunksDataSize: storeChunksDataSize,
		getChunksLatency:    getChunksLatency,
		getChunksDataSize:   getChunksDataSize,
		dbSize:              dbSize,
		dbSizePollPeriod:    dbSizePollPeriod,
		dbDir:               dbDir,
		isAlive:             isAlive,
	}, nil
}

// Start starts the metrics server.
func (m *V2Metrics) Start() {
	go func() {
		err := m.server.ListenAndServe()
		if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
			m.logger.Errorf("metrics server error: %v", err)
		}
	}()

	if m.dbSizePollPeriod.Nanoseconds() == 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(m.dbSizePollPeriod)

		for m.isAlive.Load() {
			var size int64
			err := filepath.Walk(m.dbDir, func(_ string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					size += info.Size()
				}
				return err
			})

			if err != nil {
				m.logger.Errorf("failed to get database size: %v", err)
			} else {
				m.dbSize.WithLabelValues().Set(float64(size))
			}
			<-ticker.C
		}
	}()

}

// Stop stops the metrics server.
func (m *V2Metrics) Stop() error {
	m.isAlive.Store(false)
	return m.server.Close()
}

// GetGRPCServerOption returns the gRPC server option that enables automatic GRPC metrics collection.
func (m *V2Metrics) GetGRPCServerOption() grpc.ServerOption {
	return m.grpcServerOption
}

func (m *V2Metrics) ReportStoreChunksLatency(latency time.Duration) {
	m.storeChunksLatency.WithLabelValues().Observe(
		float64(latency.Nanoseconds()) / float64(time.Millisecond))
}

func (m *V2Metrics) ReportStoreChunksDataSize(size uint64) {
	m.storeChunksDataSize.WithLabelValues().Set(float64(size))
}

func (m *V2Metrics) ReportGetChunksLatency(latency time.Duration) {
	m.getChunksLatency.WithLabelValues().Observe(
		float64(latency.Nanoseconds()) / float64(time.Millisecond))
}

func (m *V2Metrics) ReportGetChunksDataSize(size uint64) {
	m.getChunksDataSize.WithLabelValues().Set(float64(size))
}
