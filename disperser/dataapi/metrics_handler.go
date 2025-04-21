package dataapi

import (
	"context"
	"errors"
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
	version    DataApiVersion
}

func NewMetricsHandler(promClient PrometheusClient, version DataApiVersion) *MetricsHandler {
	return &MetricsHandler{
		promClient: promClient,
		version:    version,
	}
}

func (mh *MetricsHandler) GetCompleteBlobSize(ctx context.Context, startTime int64, endTime int64) (*PrometheusResult, error) {
	var result *PrometheusResult
	var err error
	if mh.version == V1 {
		result, err = mh.promClient.QueryDisperserBlobSizeBytesPerSecond(ctx, time.Unix(startTime, 0), time.Unix(endTime, 0))
	} else {
		result, err = mh.promClient.QueryDisperserBlobSizeBytesPerSecondV2(ctx, time.Unix(startTime, 0), time.Unix(endTime, 0))
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (mh *MetricsHandler) GetAvgThroughput(ctx context.Context, startTime int64, endTime int64) (float64, error) {
	var result *PrometheusResult
	var err error
	if mh.version == V1 {
		result, err = mh.promClient.QueryDisperserBlobSizeBytesPerSecond(ctx, time.Unix(startTime, 0), time.Unix(endTime, 0))
	} else {
		result, err = mh.promClient.QueryDisperserBlobSizeBytesPerSecondV2(ctx, time.Unix(startTime, 0), time.Unix(endTime, 0))
	}
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

func (mh *MetricsHandler) GetQuorumSigningRateTimeseries(ctx context.Context, startTime time.Time, endTime time.Time, quorumID uint8) (*PrometheusResult, error) {
	if mh.version != V2 {
		return nil, errors.New("only V2 signing rate fetch is supported")
	}

	result, err := mh.promClient.QueryQuorumNetworkSigningRateV2(ctx, startTime, endTime, quorumID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (mh *MetricsHandler) GetThroughputTimeseries(ctx context.Context, startTime int64, endTime int64) ([]*Throughput, error) {
	throughputRateSecs := uint16(defaultThroughputRateSecs)
	if endTime-startTime >= 7*24*60*60 {
		throughputRateSecs = uint16(sevenDayThroughputRateSecs)
	}

	var result *PrometheusResult
	var err error
	if mh.version == V1 {
		result, err = mh.promClient.QueryDisperserAvgThroughputBlobSizeBytes(ctx, time.Unix(startTime, 0), time.Unix(endTime, 0), throughputRateSecs)
	} else {
		result, err = mh.promClient.QueryDisperserAvgThroughputBlobSizeBytesV2(ctx, time.Unix(startTime, 0), time.Unix(endTime, 0), throughputRateSecs)
	}

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
