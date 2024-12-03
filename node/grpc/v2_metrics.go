package grpc

import (
	"github.com/Layr-Labs/eigenda/common/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
	"os"
	"path/filepath"
	"time"
)

// V2Metrics encapsulates metrics for the v2 DA node.
type V2Metrics struct {
	metricsServer    metrics.Metrics
	grpcServerOption grpc.ServerOption

	StoreChunksLatency  metrics.LatencyMetric
	StoreChunksDataSize metrics.GaugeMetric

	GetChunksLatency  metrics.LatencyMetric
	GetChunksDataSize metrics.GaugeMetric
}

// NewV2Metrics creates a new V2Metrics instance. dbSizePollPeriod is the period at which the database size is polled.
// If set to 0, the database size is not polled.
func NewV2Metrics(
	logger logging.Logger,
	port int,
	dbDir string,
	dbSizePollPeriod time.Duration) (*V2Metrics, error) {

	server := metrics.NewMetrics(logger, "eigenda_node", port)

	grpcMetrics := grpcprom.NewServerMetrics()
	server.RegisterExternalMetrics(grpcMetrics)
	grpcServerOption := grpc.UnaryInterceptor(
		grpcMetrics.UnaryServerInterceptor(),
	)

	storeChunksLatency, err := server.NewLatencyMetric(
		"store_chunks_latency",
		"The latency of a StoreChunks() RPC call.",
		nil,
		metrics.NewQuantile(0.5),
		metrics.NewQuantile(0.9),
		metrics.NewQuantile(0.99))
	if err != nil {
		return nil, err
	}

	storeChunksDataSize, err := server.NewGaugeMetric(
		"store_chunks_data_size",
		"bytes",
		"The size of the data requested to be stored by StoreChunks() RPC calls.",
		nil)
	if err != nil {
		return nil, err
	}

	getChunksLatency, err := server.NewLatencyMetric(
		"get_chunks_latency",
		"The latency of a GetChunks() RPC call.",
		nil,
		metrics.NewQuantile(0.5),
		metrics.NewQuantile(0.9),
		metrics.NewQuantile(0.99))
	if err != nil {
		return nil, err
	}

	getChunksDataSize, err := server.NewGaugeMetric(
		"get_chunks_data_size",
		"bytes",
		"The size of the data requested to be retrieved by GetChunks() RPC calls.",
		nil)
	if err != nil {
		return nil, err
	}

	if dbSizePollPeriod.Nanoseconds() > 0 {
		err = server.NewAutoGauge(
			"db_size",
			"bytes",
			"The size of the leveldb database.",
			dbSizePollPeriod,
			func() float64 {
				var size int64
				err = filepath.Walk(dbDir, func(_ string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() {
						size += info.Size()
					}
					return err
				})
				if err != nil {
					logger.Errorf("failed to get database size (for metrics reporting): %v", err)
					return -1.0
				}
				return float64(size)
			})
		if err != nil {
			return nil, err
		}
	}

	return &V2Metrics{
		metricsServer:       server,
		grpcServerOption:    grpcServerOption,
		StoreChunksLatency:  storeChunksLatency,
		StoreChunksDataSize: storeChunksDataSize,
		GetChunksLatency:    getChunksLatency,
		GetChunksDataSize:   getChunksDataSize,
	}, nil
}

// Start starts the metrics server.
func (m *V2Metrics) Start() error {
	return m.metricsServer.Start()
}

// Stop stops the metrics server.
func (m *V2Metrics) Stop() error {
	return m.metricsServer.Stop()
}

// GetGRPCServerOption returns the gRPC server option that enables automatic GRPC metrics collection.
func (m *V2Metrics) GetGRPCServerOption() grpc.ServerOption {
	return m.grpcServerOption
}
