package grpc

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
)

const namespace = "eigenda_node"

// MetricsV2 encapsulates metrics for the v2 DA node.
type MetricsV2 struct {
	logger logging.Logger

	registry         *prometheus.Registry
	grpcServerOption grpc.ServerOption

	storeChunksRequestSize *prometheus.GaugeVec

	getChunksLatency  *prometheus.SummaryVec
	getChunksDataSize *prometheus.GaugeVec

	storeChunksStageTimer *common.StageTimer
}

// NewV2Metrics creates a new MetricsV2 instance. dbSizePollPeriod is the period at which the database size is polled.
// If set to 0, the database size is not polled.
func NewV2Metrics(logger logging.Logger, registry *prometheus.Registry) (*MetricsV2, error) {

	// These should be re-enabled once the legacy v1 metrics are removed.
	//registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	//registry.MustRegister(collectors.NewGoCollector())

	grpcMetrics := grpcprom.NewServerMetrics()
	registry.MustRegister(grpcMetrics)
	grpcServerOption := grpc.UnaryInterceptor(
		grpcMetrics.UnaryServerInterceptor(),
	)

	storeChunksRequestSize := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "store_chunks_request_size_bytes",
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

	storeChunksStageTimer := common.NewStageTimer(registry, namespace, "store_chunks", false)

	return &MetricsV2{
		logger:                 logger,
		registry:               registry,
		grpcServerOption:       grpcServerOption,
		storeChunksRequestSize: storeChunksRequestSize,
		getChunksLatency:       getChunksLatency,
		getChunksDataSize:      getChunksDataSize,
		storeChunksStageTimer:  storeChunksStageTimer,
	}, nil
}

// GetGRPCServerOption returns the gRPC server option that enables automatic GRPC metrics collection.
func (m *MetricsV2) GetGRPCServerOption() grpc.ServerOption {
	return m.grpcServerOption
}

// GetStoreChunksProbe returns a probe for measuring the latency of the StoreChunks() RPC call.
func (m *MetricsV2) GetStoreChunksProbe() *common.SequenceProbe {
	return m.storeChunksStageTimer.NewSequence()
}

func (m *MetricsV2) ReportStoreChunksRequestSize(size uint64) {
	m.storeChunksRequestSize.WithLabelValues().Set(float64(size))
}

func (m *MetricsV2) ReportGetChunksLatency(latency time.Duration) {
	m.getChunksLatency.WithLabelValues().Observe(common.ToMilliseconds(latency))
}

func (m *MetricsV2) ReportGetChunksDataSize(size int) {
	m.getChunksDataSize.WithLabelValues().Set(float64(size))
}
