package v2

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
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

// FetchNetworkSigningRate godoc
//
//	@Summary	Fetch network signing rate time series in the specified time range
//	@Tags		Metrics
//	@Produce	json
//	@Param		end			query		string	false	"Fetch network signing rate up to the end time (ISO 8601 format: 2006-01-02T15:04:05Z) [default: now]"
//	@Param		interval	query		int		false	"Fetch network signing rate starting from an interval (in seconds) before the end time [default: 3600]"
//	@Param		quorums		query		string	false	"Comma-separated list of quorum IDs to filter (e.g., 0,1) [default: 0,1]"
//	@Success	200			{object}	NetworkSigningRateResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/metrics/timeseries/network-signing-rate [get]
func (s *ServerV2) FetchNetworkSigningRate(c *gin.Context) {
	handlerStart := time.Now()
	var err error

	now := handlerStart
	oldestTime := now.Add(-maxBlobAge)

	endTime := now
	if c.Query("end") != "" {
		endTime, err = time.Parse("2006-01-02T15:04:05Z", c.Query("end"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchNetworkSigningRate")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse end param: %w", err))
			return
		}
		if endTime.Before(oldestTime) {
			s.metrics.IncrementInvalidArgRequestNum("FetchNetworkSigningRate")
			invalidParamsErrorResponse(
				c, fmt.Errorf("end time cannot be more than 14 days in the past, found: %s", c.Query("end")),
			)
			return
		}
	}

	interval := 3600
	if c.Query("interval") != "" {
		interval, err = strconv.Atoi(c.Query("interval"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchNetworkSigningRate")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse interval param: %w", err))
			return
		}
		if interval <= 0 {
			s.metrics.IncrementInvalidArgRequestNum("FetchNetworkSigningRate")
			invalidParamsErrorResponse(c, fmt.Errorf("interval must be greater than 0, found: %d", interval))
			return
		}
		if maxInterval := int(maxBlobAge / time.Second); interval > maxInterval {
			interval = maxInterval
		}
	}

	quorums := []uint8{0, 1}
	if quorumStr := c.Query("quorums"); quorumStr != "" {
		quorumStrs := strings.Split(quorumStr, ",")
		for _, qStr := range quorumStrs {
			q, err := strconv.ParseUint(qStr, 10, 8)
			if err != nil || q > maxQuorumIDAllowed {
				s.metrics.IncrementInvalidArgRequestNum("FetchNetworkSigningRate")
				if err != nil {
					invalidParamsErrorResponse(c, fmt.Errorf("failed to parse quorums param: %w", err))
				} else {
					invalidParamsErrorResponse(c, fmt.Errorf("the quorum ID must be in range [0, %d], found: %d", maxQuorumIDAllowed, q))
				}
				return
			}
			quorums = append(quorums, uint8(q))
		}
	}

	response := NetworkSigningRateResponse{
		QuorumSigningRates: make([]QuorumSigningRateData, 0, len(quorums)),
	}

	startTime := endTime.Add(-time.Duration(interval) * time.Second)
	for _, quorum := range quorums {
		result, err := s.metricsHandler.GetQuorumSigningRateTimeseries(c.Request.Context(), startTime, endTime, quorum)
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchNetworkSigningRate")
			errorResponse(c, err)
			return
		}
		if len(result.Values) > 0 {
			dataPoints := make([]SigningRateDataPoint, len(result.Values))
			for i, point := range result.Values {
				dataPoints[i] = SigningRateDataPoint{
					SigningRate: point.Value,
					Timestamp:   uint64(point.Timestamp.Unix()),
				}
			}
			data := QuorumSigningRateData{
				QuorumId:   fmt.Sprintf("%d", quorum),
				DataPoints: dataPoints,
			}
			response.QuorumSigningRates = append(response.QuorumSigningRates, data)
		}
	}

	// Sort the quorums by ID for consistent output
	sort.Slice(response.QuorumSigningRates, func(i, j int) bool {
		return response.QuorumSigningRates[i].QuorumId < response.QuorumSigningRates[j].QuorumId
	})

	s.metrics.IncrementSuccessfulRequestNum("FetchNetworkSigningRate")
	s.metrics.ObserveLatency("FetchNetworkSigningRate", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxSigningInfoAge))
	c.JSON(http.StatusOK, response)
}
