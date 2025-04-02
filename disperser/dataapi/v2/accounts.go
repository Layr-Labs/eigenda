package v2

import (
	"errors"
	"fmt"
	"net/http"
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
//	@Router		/operators/{operator_id}/dispersals [get]
func (s *ServerV2) FetchAccountBlobFeed(c *gin.Context) {
	handlerStart := time.Now()
	var err error

	// Parse account ID
	accountStr := c.Param("account_id")
	if !gethcommon.IsHexAddress(accountStr) {
		s.metrics.IncrementInvalidArgRequestNum("FetchAccountBlobFeed")
		errorResponse(c, errors.New("account id is not valid hex"))
		return
	}
	accountId := gethcommon.HexToAddress(accountStr)
	if accountId == (gethcommon.Address{}) {
		s.metrics.IncrementInvalidArgRequestNum("FetchAccountBlobFeed")
		errorResponse(c, errors.New("zero account id is not valid"))
		return
	}

	// Parse the feed params
	params, err := ParseFeedParams(c, s.metrics, "FetchAccountBlobFeed")
	if err != nil {
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
			errorResponse(c, fmt.Errorf("failed to serialize blob key: %w", err))
			return
		}
		blobInfo[i].BlobKey = bk.Hex()
		blobInfo[i].BlobMetadata = blobs[i]
	}

	response := &AccountBlobFeedResponse{
		AccountId: accountId.Hex(),
		Blobs:     blobInfo,
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchAccountBlobFeed")
	s.metrics.ObserveLatency("FetchAccountBlobFeed", time.Since(handlerStart))
	c.JSON(http.StatusOK, response)
}
