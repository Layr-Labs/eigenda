package dataapi

import (
	"context"
	"time"
)

const (
	defaultThroughputRateSecs  = 240 // 4m rate is used for < 7d window to match $__rate_interval
	sevenDayThroughputRateSecs = 660 // 11m rate is used for >= 7d window to match $__rate_interval
)

// metricHandler handles operations to collect metrics about the Disperser.
type MetricsHandler struct {
	// For accessing metrics info
	promClient PrometheusClient
}

func NewMetricsHandler(promClient PrometheusClient) *MetricsHandler {
	return &MetricsHandler{
		promClient: promClient,
	}
}

func (mh *MetricsHandler) GetAvgThroughput(ctx context.Context, startTime int64, endTime int64) (float64, error) {
	result, err := mh.promClient.QueryDisperserBlobSizeBytesPerSecond(ctx, time.Unix(startTime, 0), time.Unix(endTime, 0))
	if err != nil {
		return 0, err
	}
	size := len(result.Values)
	if size == 0 {
		return 0, nil
	}
	totalBytes := result.Values[size-1].Value - result.Values[0].Value
	timeDuration := result.Values[size-1].Timestamp.Sub(result.Values[0].Timestamp).Seconds()
	return totalBytes / timeDuration, nil
}

func (mh *MetricsHandler) GetThroughputTimeseries(ctx context.Context, startTime int64, endTime int64) ([]*Throughput, error) {
	throughputRateSecs := uint16(defaultThroughputRateSecs)
	if endTime-startTime >= 7*24*60*60 {
		throughputRateSecs = uint16(sevenDayThroughputRateSecs)
	}
	result, err := mh.promClient.QueryDisperserAvgThroughputBlobSizeBytes(ctx, time.Unix(startTime, 0), time.Unix(endTime, 0), throughputRateSecs)
	if err != nil {
		return nil, err
	}

	if len(result.Values) <= 1 {
		return []*Throughput{}, nil
	}

	throughputs := make([]*Throughput, 0)
	for i := throughputRateSecs; i < uint16(len(result.Values)); i++ {
		v := result.Values[i]
		throughputs = append(throughputs, &Throughput{
			Timestamp:  uint64(v.Timestamp.Unix()),
			Throughput: v.Value,
		})
	}

	return throughputs, nil
}
