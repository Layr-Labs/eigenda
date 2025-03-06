package v2

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/gin-gonic/gin"
)

// FetchBatchFeed godoc
//
//	@Summary	Fetch batch feed in specified direction
//	@Tags		Batches
//	@Produce	json
//	@Param		direction	query		string	false	"Direction to fetch: 'forward' (oldest to newest, ASC order) or 'backward' (newest to oldest, DESC order) [default: forward]"
//	@Param		before		query		string	false	"Fetch batches before this time, exclusive (ISO 8601 format, example: 2006-01-02T15:04:05Z) [default: now]"
//	@Param		after		query		string	false	"Fetch batches after this time, exclusive (ISO 8601 format, example: 2006-01-02T15:04:05Z); must be smaller than 'before' [default: before-1h]"
//	@Param		limit		query		int		false	"Maximum number of batches to return; if limit <= 0 or >1000, it's treated as 1000 [default: 20; max: 1000]"
//	@Success	200			{object}	BatchFeedResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/batches/feed [get]
func (s *ServerV2) FetchBatchFeed(c *gin.Context) {
	handlerStart := time.Now()
	var err error

	// Validate direction
	direction := "forward"
	if dirStr := c.Query("direction"); dirStr != "" {
		if dirStr != "forward" && dirStr != "backward" {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("`direction` must be either \"forward\" or \"backward\", found: %q", dirStr))
			return
		}
		direction = dirStr
	}

	now := handlerStart
	oldestTime := now.Add(-maxBlobAge)

	// Handle before parameter
	beforeTime := now
	if c.Query("before") != "" {
		beforeTime, err = time.Parse("2006-01-02T15:04:05Z", c.Query("before"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse `before` param: %w", err))
			return
		}
		if beforeTime.Before(oldestTime) {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("`before` time cannot be more than 14 days in the past, found: `%s`", c.Query("before")))
			return
		}
		if now.Before(beforeTime) {
			beforeTime = now
		}
	}

	// Handle after parameter
	afterTime := beforeTime.Add(-time.Hour)
	if c.Query("after") != "" {
		afterTime, err = time.Parse("2006-01-02T15:04:05Z", c.Query("after"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse `after` param: %w", err))
			return
		}
		if now.Before(afterTime) {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("`after` must be before current time, found: `%s`", c.Query("after")))
			return
		}
		if afterTime.Before(oldestTime) {
			afterTime = oldestTime
		}
	}

	// Validate time range
	if !afterTime.Before(beforeTime) {
		s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
		invalidParamsErrorResponse(c, fmt.Errorf("`after` timestamp (%s) must be earlier than `before` timestamp (%s)",
			afterTime.Format(time.RFC3339), beforeTime.Format(time.RFC3339)))
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
		invalidParamsErrorResponse(c, fmt.Errorf("failed to parse `limit` param: %w", err))
		return
	}
	if limit <= 0 || limit > maxNumBatchesPerBatchFeedResponse {
		limit = maxNumBatchesPerBatchFeedResponse
	}

	var attestations []*corev2.Attestation

	if direction == "forward" {
		attestations, err = s.batchFeedCache.Get(
			c.Request.Context(),
			afterTime.Add(time.Nanosecond), // +1ns to make it exclusive
			beforeTime,
			Ascending,
			limit,
		)
	} else {
		attestations, err = s.batchFeedCache.Get(
			c.Request.Context(),
			afterTime.Add(time.Nanosecond), // +1ns to make it exclusive
			beforeTime,
			Descending,
			limit,
		)
	}

	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBatchFeed")
		errorResponse(c, fmt.Errorf("failed to fetch feed from blob metadata store: %w", err))
		return
	}

	batches := make([]*BatchInfo, len(attestations))
	for i, at := range attestations {
		batchHeaderHash, err := at.BatchHeader.Hash()
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchBatchFeed")
			errorResponse(c, fmt.Errorf("failed to compute batch header hash from batch header: %w", err))
			return
		}
		batches[i] = &BatchInfo{
			BatchHeaderHash:         hex.EncodeToString(batchHeaderHash[:]),
			BatchHeader:             at.BatchHeader,
			AttestedAt:              at.AttestedAt,
			AggregatedSignature:     at.Sigma,
			QuorumNumbers:           at.QuorumNumbers,
			QuorumSignedPercentages: at.QuorumResults,
		}
	}

	response := &BatchFeedResponse{
		Batches: batches,
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchBatchFeed")
	s.metrics.ObserveLatency("FetchBatchFeed", time.Since(handlerStart))
	c.JSON(http.StatusOK, response)
}

// FetchBatch godoc
//
//	@Summary	Fetch batch by the batch header hash
//	@Tags		Batches
//	@Produce	json
//	@Param		batch_header_hash	path		string	true	"Batch header hash in hex string"
//	@Success	200					{object}	BatchResponse
//	@Failure	400					{object}	ErrorResponse	"error: Bad request"
//	@Failure	404					{object}	ErrorResponse	"error: Not found"
//	@Failure	500					{object}	ErrorResponse	"error: Server error"
//	@Router		/batches/{batch_header_hash} [get]
func (s *ServerV2) FetchBatch(c *gin.Context) {
	handlerStart := time.Now()

	batchHeaderHashHex := c.Param("batch_header_hash")
	batchHeaderHash, err := dataapi.ConvertHexadecimalToBytes([]byte(batchHeaderHashHex))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBatch")
		errorResponse(c, errors.New("invalid batch header hash"))
		return
	}
	batchHeader, attestation, err := s.blobMetadataStore.GetSignedBatch(c.Request.Context(), batchHeaderHash)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBatch")
		errorResponse(c, err)
		return
	}
	// TODO: support fetch of blob inclusion info
	batchResponse := &BatchResponse{
		BatchHeaderHash: batchHeaderHashHex,
		SignedBatch: &SignedBatch{
			BatchHeader: batchHeader,
			Attestation: attestation,
		},
	}
	s.metrics.IncrementSuccessfulRequestNum("FetchBatch")
	s.metrics.ObserveLatency("FetchBatch", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	c.JSON(http.StatusOK, batchResponse)
}
