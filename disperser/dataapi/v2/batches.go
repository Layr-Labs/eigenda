package v2

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/gin-gonic/gin"
)

// FetchBatchFeed godoc
//
//	@Summary	Fetch batch feed
//	@Tags		Batches
//	@Produce	json
//	@Param		end			query		string	false	"Fetch batches up to the end time (ISO 8601 format: 2006-01-02T15:04:05Z) [default: now]"
//	@Param		interval	query		int		false	"Fetch batches starting from an interval (in seconds) before the end time [default: 3600]"
//	@Param		limit		query		int		false	"The maximum number of batches to fetch. System max (1000) if limit <= 0 [default: 20; max: 1000]"
//	@Success	200			{object}	BatchFeedResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/batches/feed [get]
func (s *ServerV2) FetchBatchFeed(c *gin.Context) {
	handlerStart := time.Now()
	var err error

	now := handlerStart
	oldestTime := now.Add(-maxBlobAge)

	endTime := now
	if c.Query("end") != "" {
		endTime, err = time.Parse("2006-01-02T15:04:05Z", c.Query("end"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse end param: %w", err))
			return
		}
		if endTime.Before(oldestTime) {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("end time cannot be more than 14 days in the past, found: %s", c.Query("end")))
			return
		}
	}

	interval := 3600
	if c.Query("interval") != "" {
		interval, err = strconv.Atoi(c.Query("interval"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse interval param: %w", err))
			return
		}
		if interval <= 0 {
			s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("interval must be greater than 0, found: %d", interval))
			return
		}
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBatchFeed")
		invalidParamsErrorResponse(c, fmt.Errorf("failed to parse limit param: %w", err))
		return
	}
	if limit <= 0 || limit > maxNumBatchesPerBatchFeedResponse {
		limit = maxNumBatchesPerBatchFeedResponse
	}

	startTime := endTime.Add(-time.Duration(interval) * time.Second)
	attestations, err := s.blobMetadataStore.GetAttestationByAttestedAt(c.Request.Context(), uint64(startTime.UnixNano())+1, uint64(endTime.UnixNano()), limit)
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
	s.metrics.ObserveLatency("FetchBatchFeed", time.Since(now))
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
