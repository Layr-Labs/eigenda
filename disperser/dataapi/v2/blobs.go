package v2

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/gin-gonic/gin"
)

// FetchBlobFeed godoc
//
//	@Summary	Fetch blob feed
//	@Tags		Blobs
//	@Produce	json
//	@Param		end					query		string	false	"Fetch blobs up to the end time (ISO 8601 format: 2006-01-02T15:04:05Z) [default: now]"
//	@Param		interval			query		int		false	"Fetch blobs starting from an interval (in seconds) before the end time [default: 3600]"
//	@Param		pagination_token	query		string	false	"Fetch blobs starting from the pagination token (exclusively). Overrides the interval param if specified [default: empty]"
//	@Param		limit				query		int		false	"The maximum number of blobs to fetch. System max (1000) if limit <= 0 [default: 20; max: 1000]"
//	@Success	200					{object}	BlobFeedResponse
//	@Failure	400					{object}	ErrorResponse	"error: Bad request"
//	@Failure	404					{object}	ErrorResponse	"error: Not found"
//	@Failure	500					{object}	ErrorResponse	"error: Server error"
//	@Router		/blobs/feed [get]
func (s *ServerV2) FetchBlobFeed(c *gin.Context) {
	handlerStart := time.Now()
	var err error

	now := handlerStart
	oldestTime := now.Add(-maxBlobAge)

	endTime := now
	if c.Query("end") != "" {
		endTime, err = time.Parse("2006-01-02T15:04:05Z", c.Query("end"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse end param: %w", err))
			return
		}
		if endTime.Before(oldestTime) {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("end time cannot be more than 14 days in the past, found: %s", c.Query("end")))
			return
		}
	}

	interval := 3600
	if c.Query("interval") != "" {
		interval, err = strconv.Atoi(c.Query("interval"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse interval param: %w", err))
			return
		}
		if interval <= 0 {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("interval must be greater than 0, found: %d", interval))
			return
		}
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
		invalidParamsErrorResponse(c, fmt.Errorf("failed to parse limit param: %w", err))
		return
	}
	if limit <= 0 || limit > maxNumBlobsPerBlobFeedResponse {
		limit = maxNumBlobsPerBlobFeedResponse
	}

	paginationCursor := blobstore.BlobFeedCursor{
		RequestedAt: 0,
	}
	if c.Query("pagination_token") != "" {
		cursor, err := paginationCursor.FromCursorKey(c.Query("pagination_token"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse the pagination token: %w", err))
			return
		}
		paginationCursor = *cursor
	}

	startTime := endTime.Add(-time.Duration(interval) * time.Second)
	if startTime.Before(oldestTime) {
		startTime = oldestTime
	}
	startCursor := blobstore.BlobFeedCursor{
		RequestedAt: uint64(startTime.UnixNano()),
	}
	if startCursor.LessThan(&paginationCursor) {
		startCursor = paginationCursor
	}
	endCursor := blobstore.BlobFeedCursor{
		RequestedAt: uint64(endTime.UnixNano()),
	}

	blobs, paginationToken, err := s.blobMetadataStore.GetBlobMetadataByRequestedAt(c.Request.Context(), startCursor, endCursor, limit)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlobFeed")
		errorResponse(c, fmt.Errorf("failed to fetch feed from blob metadata store: %w", err))
		return
	}

	token := ""
	if paginationToken != nil {
		token = paginationToken.ToCursorKey()
	}
	blobInfo := make([]BlobInfo, len(blobs))
	for i := 0; i < len(blobs); i++ {
		bk, err := blobs[i].BlobHeader.BlobKey()
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchBlobFeed")
			errorResponse(c, fmt.Errorf("failed to serialize blob key: %w", err))
			return
		}
		blobInfo[i].BlobKey = bk.Hex()
		blobInfo[i].BlobMetadata = blobs[i]
	}
	response := &BlobFeedResponse{
		Blobs:           blobInfo,
		PaginationToken: token,
	}
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	s.metrics.IncrementSuccessfulRequestNum("FetchBlobFeed")
	s.metrics.ObserveLatency("FetchBlobFeed", time.Since(handlerStart))
	c.JSON(http.StatusOK, response)
}

// FetchBlob godoc
//
//	@Summary	Fetch blob metadata by blob key
//	@Tags		Blobs
//	@Produce	json
//	@Param		blob_key	path		string	true	"Blob key in hex string"
//	@Success	200			{object}	BlobResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/blobs/{blob_key} [get]
func (s *ServerV2) FetchBlob(c *gin.Context) {
	handlerStart := time.Now()

	blobKey, err := corev2.HexToBlobKey(c.Param("blob_key"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlob")
		errorResponse(c, err)
		return
	}
	metadata, err := s.blobMetadataStore.GetBlobMetadata(c.Request.Context(), blobKey)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlob")
		errorResponse(c, err)
		return
	}
	bk, err := metadata.BlobHeader.BlobKey()
	if err != nil || bk != blobKey {
		s.metrics.IncrementFailedRequestNum("FetchBlob")
		errorResponse(c, err)
		return
	}
	response := &BlobResponse{
		BlobKey:       bk.Hex(),
		BlobHeader:    metadata.BlobHeader,
		Status:        metadata.BlobStatus.String(),
		DispersedAt:   metadata.RequestedAt,
		BlobSizeBytes: metadata.BlobSize,
	}
	s.metrics.IncrementSuccessfulRequestNum("FetchBlob")
	s.metrics.ObserveLatency("FetchBlob", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	c.JSON(http.StatusOK, response)
}

// FetchBlobCertificate godoc
//
//	@Summary	Fetch blob certificate by blob key v2
//	@Tags		Blobs
//	@Produce	json
//	@Param		blob_key	path		string	true	"Blob key in hex string"
//	@Success	200			{object}	BlobCertificateResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/blobs/{blob_key}/certificate [get]
func (s *ServerV2) FetchBlobCertificate(c *gin.Context) {
	handlerStart := time.Now()

	blobKey, err := corev2.HexToBlobKey(c.Param("blob_key"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlobCertificate")
		errorResponse(c, err)
		return
	}
	cert, _, err := s.blobMetadataStore.GetBlobCertificate(c.Request.Context(), blobKey)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlobCertificate")
		errorResponse(c, err)
		return
	}
	response := &BlobCertificateResponse{
		Certificate: cert,
	}
	s.metrics.IncrementSuccessfulRequestNum("FetchBlobCertificate")
	s.metrics.ObserveLatency("FetchBlobCertificate", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	c.JSON(http.StatusOK, response)
}

// FetchBlobAttestationInfo godoc
//
//	@Summary	Fetch attestation info for a blob
//	@Tags		Blobs
//	@Produce	json
//	@Param		blob_key	path		string	true	"Blob key in hex string"
//	@Success	200			{object}	BlobAttestationInfoResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/blobs/{blob_key}/attestation-info [get]
func (s *ServerV2) FetchBlobAttestationInfo(c *gin.Context) {
	handlerStart := time.Now()
	blobKey, err := corev2.HexToBlobKey(c.Param("blob_key"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlobAttestationInfo")
		invalidParamsErrorResponse(c, fmt.Errorf("failed to parse blob_key param: %w", err))
		return
	}

	attestationInfo, err := s.blobMetadataStore.GetBlobAttestationInfo(c.Request.Context(), blobKey)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlobAttestationInfo")
		errorResponse(c, fmt.Errorf("failed to fetch blob attestation info: %w", err))
		return
	}

	batchHeaderHash, err := attestationInfo.InclusionInfo.BatchHeader.Hash()
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlobAttestationInfo")
		errorResponse(c, fmt.Errorf("failed to get batch header hash from blob inclusion info: %w", err))
		return
	}

	response := &BlobAttestationInfoResponse{
		BlobKey:         blobKey.Hex(),
		BatchHeaderHash: hex.EncodeToString(batchHeaderHash[:]),
		InclusionInfo:   attestationInfo.InclusionInfo,
		Attestation:     attestationInfo.Attestation,
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchBlobAttestationInfo")
	s.metrics.ObserveLatency("FetchBlobAttestationInfo", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	c.JSON(http.StatusOK, response)
}
