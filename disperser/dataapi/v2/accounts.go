package v2

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// FetchAccountBlobFeed godoc
//
//	@Summary	Fetch blobs posted by an account in a time window by specific direction
//	@Tags		Accounts
//	@Produce	json
//	@Param		account_id	path		string	true	"The account ID to fetch blob feed for"
//	@Param		direction	query		string	false	"Direction to fetch: 'forward' (oldest to newest, ASC order) or 'backward' (newest to oldest, DESC order) [default: forward]"
//	@Param		before		query		string	false	"Fetch blobs before this time, exclusive (ISO 8601 format, example: 2006-01-02T15:04:05Z) [default: now]"
//	@Param		after		query		string	false	"Fetch blobs after this time, exclusive (ISO 8601 format, example: 2006-01-02T15:04:05Z); must be smaller than `before` [default: `before`-1h]"
//	@Param		limit		query		int		false	"Maximum number of blobs to return; if limit <= 0 or >1000, it's treated as 1000 [default: 20; max: 1000]"
//	@Success	200			{object}	AccountBlobFeedResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/accounts/{account_id}/blobs [get]
func (s *ServerV2) FetchAccountBlobFeed(c *gin.Context) {
	handlerStart := time.Now()
	var err error

	// Parse account ID
	accountStr := c.Param("account_id")
	if !gethcommon.IsHexAddress(accountStr) {
		s.metrics.IncrementInvalidArgRequestNum("FetchAccountBlobFeed")
		invalidParamsErrorResponse(c, errors.New("account id is not valid hex"))
		return
	}
	accountId := gethcommon.HexToAddress(accountStr)
	if accountId == (gethcommon.Address{}) {
		s.metrics.IncrementInvalidArgRequestNum("FetchAccountBlobFeed")
		invalidParamsErrorResponse(c, errors.New("zero account id is not valid"))
		return
	}

	// Parse the feed params
	params, err := ParseFeedParams(c, s.metrics, "FetchAccountBlobFeed")
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchAccountBlobFeed")
		invalidParamsErrorResponse(c, err)
		return
	}

	var blobs []*v2.BlobMetadata

	if params.direction == "forward" {
		blobs, err = s.blobMetadataStore.GetBlobMetadataByAccountID(
			c.Request.Context(),
			accountId,
			uint64(params.afterTime.UnixNano()),
			uint64(params.beforeTime.UnixNano()),
			params.limit,
			true, // ascending=true
		)
	} else {
		blobs, err = s.blobMetadataStore.GetBlobMetadataByAccountID(
			c.Request.Context(),
			accountId,
			uint64(params.afterTime.UnixNano()),
			uint64(params.beforeTime.UnixNano()),
			params.limit,
			false, // ascending=false
		)
	}

	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchAccountBlobFeed")
		errorResponse(c, fmt.Errorf("failed to fetch blobs from blob metadata store for account (%s): %w", accountId.Hex(), err))
		return
	}

	blobInfo := make([]BlobInfo, len(blobs))
	for i := 0; i < len(blobs); i++ {
		bk, err := blobs[i].BlobHeader.BlobKey()
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchAccountBlobFeed")
			errorResponse(c, fmt.Errorf("blob metadata is malformed and failed to serialize blob key: %w", err))
			return
		}
		blobInfo[i].BlobKey = bk.Hex()
		blobInfo[i].BlobMetadata = createBlobMetadata(blobs[i])
	}

	response := &AccountBlobFeedResponse{
		AccountId: accountId.Hex(),
		Blobs:     blobInfo,
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchAccountBlobFeed")
	s.metrics.ObserveLatency("FetchAccountBlobFeed", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxBlobFeedAge))
	c.JSON(http.StatusOK, response)
}

// FetchAccountFeed godoc
//
//	@Summary	Fetch accounts within a time window (sorted by latest timestamp)
//	@Tags		Accounts
//	@Produce	json
//	@Param		lookback_hours	query		int	false	"Number of hours to look back [default: 24; max: 24000 (1000 days)]"
//	@Success	200				{object}	AccountFeedResponse
//	@Failure	400				{object}	ErrorResponse	"error: Bad request"
//	@Failure	500				{object}	ErrorResponse	"error: Server error"
//	@Router		/accounts [get]
func (s *ServerV2) FetchAccountFeed(c *gin.Context) {
	handlerStart := time.Now()

	// Parse lookback_hours parameter
	lookbackHoursStr := c.Query("lookback_hours")
	lookbackHours := 24 // default to 24 hours
	if lookbackHoursStr != "" {
		parsedHours, err := strconv.Atoi(lookbackHoursStr)
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchAccountFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("invalid lookback_hours parameter: %w", err))
			return
		}
		if parsedHours > 24000 { // max 1000 days
			lookbackHours = 24000
		} else if parsedHours > 0 {
			lookbackHours = parsedHours
		}
	}

	lookbackSeconds := uint64(lookbackHours * 3600) // convert hours to seconds

	// Check cache first
	cacheKey := fmt.Sprintf("account_feed:%d", lookbackHours)
	if cached, ok := s.accountCache.Get(cacheKey); ok {
		s.metrics.IncrementCacheHit("FetchAccountFeed")
		s.metrics.IncrementSuccessfulRequestNum("FetchAccountFeed")
		s.metrics.ObserveLatency("FetchAccountFeed", time.Since(handlerStart))
		c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxAccountAge))
		c.JSON(http.StatusOK, cached)
		return
	}

	// Query accounts within time window
	accounts, err := s.blobMetadataStore.GetAccounts(c.Request.Context(), lookbackSeconds)
	if err != nil {
		s.logger.Error("failed to fetch accounts", "error", err, "lookbackHours", lookbackHours)
		s.metrics.IncrementFailedRequestNum("FetchAccountFeed")
		errorResponse(c, err)
		return
	}

	// Convert to API response format
	accountResponses := make([]AccountResponse, len(accounts))
	for i, account := range accounts {
		// Safely convert uint64 to int64 with bounds checking
		var timestamp int64
		if account.UpdatedAt > math.MaxInt64 {
			timestamp = 0
		} else {
			timestamp = int64(account.UpdatedAt)
		}

		accountResponses[i] = AccountResponse{
			Address:     account.Address.Hex(),
			DispersedAt: time.Unix(timestamp, 0).UTC().Format(time.RFC3339),
		}
	}

	response := &AccountFeedResponse{
		Accounts: accountResponses,
	}

	// Cache the response
	s.accountCache.Add(cacheKey, response)

	s.metrics.IncrementSuccessfulRequestNum("FetchAccountFeed")
	s.metrics.ObserveLatency("FetchAccountFeed", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxAccountAge))
	c.JSON(http.StatusOK, response)
}
