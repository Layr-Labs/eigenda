package dataapi

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus"
	"github.com/prometheus/common/model"
)

const (
	// maxNumOfDataPoints is the maximum number of data points that can be queried from Prometheus based on latency that this API can provide
	maxNumOfDataPoints = 3500
)

type (
	PrometheusClient interface {
		QueryDisperserBlobSizeBytesPerSecond(ctx context.Context, start time.Time, end time.Time) (*PrometheusResult, error)
		QueryDisperserAvgThroughputBlobSizeBytes(ctx context.Context, start time.Time, end time.Time, windowSizeInSec uint16) (*PrometheusResult, error)
	}

	PrometheusResultValues struct {
		Timestamp time.Time
		Value     float64
	}

	PrometheusResult struct {
		Values []*PrometheusResultValues
	}

	prometheusClient struct {
		api     prometheus.Api
		cluster string
	}
)

var _ PrometheusClient = (*prometheusClient)(nil)

func NewPrometheusClient(api prometheus.Api, cluster string) *prometheusClient {
	return &prometheusClient{api: api, cluster: cluster}
}

func (pc *prometheusClient) QueryDisperserBlobSizeBytesPerSecond(ctx context.Context, start time.Time, end time.Time) (*PrometheusResult, error) {
	query := fmt.Sprintf("eigenda_batcher_blobs_total{state=\"confirmed\",data=\"size\",cluster=\"%s\"}", pc.cluster)
	return pc.queryRange(ctx, query, start, end)
}

func (pc *prometheusClient) QueryDisperserAvgThroughputBlobSizeBytes(ctx context.Context, start time.Time, end time.Time, throughputRateSecs uint16) (*PrometheusResult, error) {
	query := fmt.Sprintf("avg_over_time( sum by (job) (rate(eigenda_batcher_blobs_total{state=\"confirmed\",data=\"size\",cluster=\"%s\"}[%ds])) [9m:])", pc.cluster, throughputRateSecs)
	return pc.queryRange(ctx, query, start, end)
}

func (pc *prometheusClient) queryRange(ctx context.Context, query string, start time.Time, end time.Time) (*PrometheusResult, error) {
	numSecondsInTimeRange := end.Sub(start).Seconds()
	step := uint64(numSecondsInTimeRange / maxNumOfDataPoints)
	if step < 1 {
		step = 1
	}

	v, _, err := pc.api.QueryRange(ctx, query, start, end, time.Duration(step)*time.Second)
	if err != nil {
		return nil, err
	}

	values := make([]*PrometheusResultValues, 0)
	if len(v.(model.Matrix)) == 0 {
		return &PrometheusResult{
			Values: values,
		}, nil
	}

	for _, v := range v.(model.Matrix)[0].Values {
		values = append(values, &PrometheusResultValues{
			Timestamp: v.Timestamp.Time(),
			Value:     float64(v.Value),
		})
	}

	return &PrometheusResult{
		Values: values,
	}, nil
}
