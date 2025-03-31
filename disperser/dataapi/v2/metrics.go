package v2

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// FetchMetricsSummary godoc
//
//	@Summary	Fetch metrics summary
//	@Tags		Metrics
//	@Produce	json
//	@Param		start	query		int	false	"Start unix timestamp [default: 1 hour ago]"
//	@Param		end		query		int	false	"End unix timestamp [default: unix time now]"
//	@Success	200		{object}	MetricSummary
//	@Failure	400		{object}	ErrorResponse	"error: Bad request"
//	@Failure	404		{object}	ErrorResponse	"error: Not found"
//	@Failure	500		{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/summary  [get]
func (s *ServerV2) FetchMetricsSummary(c *gin.Context) {
	handlerStart := time.Now()

	now := handlerStart
	start, err := strconv.ParseInt(c.DefaultQuery("start", "0"), 10, 64)
	if err != nil || start == 0 {
		start = now.Add(-time.Hour * 1).Unix()
	}

	end, err := strconv.ParseInt(c.DefaultQuery("end", "0"), 10, 64)
	if err != nil || end == 0 {
		end = now.Unix()
	}

	result, err := s.metricsHandler.GetCompleteBlobSize(c.Request.Context(), start, end)
	if err != nil || len(result.Values) == 0 {
		s.metrics.IncrementFailedRequestNum("FetchMetricsSummary")
		errorResponse(c, err)
		return
	}

	size := len(result.Values)
	totalBytes := result.Values[size-1].Value - result.Values[0].Value
	timeDuration := result.Values[size-1].Timestamp.Sub(result.Values[0].Timestamp).Seconds()
	metricSummary := &MetricSummary{
		TotalBytesPosted:      uint64(totalBytes),
		AverageBytesPerSecond: totalBytes / timeDuration,
		StartTimestampSec:     result.Values[0].Timestamp.Unix(),
		EndTimestampSec:       result.Values[size-1].Timestamp.Unix(),
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchMetricsSummary")
	s.metrics.ObserveLatency("FetchMetricsSummary", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxMetricAge))
	c.JSON(http.StatusOK, metricSummary)
}

// FetchMetricsThroughputTimeseries godoc
//
//	@Summary	Fetch throughput time series
//	@Tags		Metrics
//	@Produce	json
//	@Param		start	query		int	false	"Start unix timestamp [default: 1 hour ago]"
//	@Param		end		query		int	false	"End unix timestamp [default: unix time now]"
//	@Success	200		{object}	[]Throughput
//	@Failure	400		{object}	ErrorResponse	"error: Bad request"
//	@Failure	404		{object}	ErrorResponse	"error: Not found"
//	@Failure	500		{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/timeseries/throughput  [get]
func (s *ServerV2) FetchMetricsThroughputTimeseries(c *gin.Context) {
	handlerStart := time.Now()

	now := handlerStart
	start, err := strconv.ParseInt(c.DefaultQuery("start", "0"), 10, 64)
	if err != nil || start == 0 {
		start = now.Add(-time.Hour * 1).Unix()
	}

	end, err := strconv.ParseInt(c.DefaultQuery("end", "0"), 10, 64)
	if err != nil || end == 0 {
		end = now.Unix()
	}

	ths, err := s.metricsHandler.GetThroughputTimeseries(c.Request.Context(), start, end)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchMetricsThroughputTimeseries")
		errorResponse(c, err)
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchMetricsThroughputTimeseries")
	s.metrics.ObserveLatency("FetchMetricsThroughputTimeseries", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxThroughputAge))
	c.JSON(http.StatusOK, ths)
}
