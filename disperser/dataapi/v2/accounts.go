package v2

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	//"time"

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

// FetchAccountReservationUsage godoc
//
//	@Summary	Fetch reservation usage for an account
//	@Tags		Accounts
//	@Produce	json
//	@Param		account_id	path		string	true	"The account ID to fetch reservation usage for"
//	@Param		window		query		int		false	"Time window in hours to fetch reservation usage for [default: 24; max: 72]"
//	@Success	200			{object}	AccountReservationUsageResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/accounts/{account_id}/reservation/usage [get]
func (s *ServerV2) FetchAccountReservationUsage(c *gin.Context) {
	handlerStart := time.Now()
	var err error

	// Parse account ID
	accountStr := c.Param("account_id")
	if !gethcommon.IsHexAddress(accountStr) {
		s.metrics.IncrementInvalidArgRequestNum("FetchAccountReservationUsage")
		invalidParamsErrorResponse(c, errors.New("account id is not valid hex"))
		return
	}
	accountId := gethcommon.HexToAddress(accountStr)
	if accountId == (gethcommon.Address{}) {
		s.metrics.IncrementInvalidArgRequestNum("FetchAccountReservationUsage")
		invalidParamsErrorResponse(c, errors.New("zero account id is not valid"))
		return
	}

	// Parse window parameter
	window := 24 // default 24 hours
	if windowStr := c.Query("window"); windowStr != "" {
		parsedWindow, err := strconv.Atoi(windowStr)
		if err != nil || parsedWindow <= 0 || parsedWindow > 72 {
			s.metrics.IncrementInvalidArgRequestNum("FetchAccountReservationUsage")
			invalidParamsErrorResponse(c, errors.New("window must be between 1 and 72 hours"))
			return
		}
		window = parsedWindow
	}

	// Calculate reservation period
	now := time.Now()
	startTime := now.Add(-time.Duration(window) * time.Hour)

	// Get period records for the specified window (limit 1000)
	periodRecords, err := s.meterer.MeteringStore.GetPeriodRecords(c.Request.Context(), accountId, uint64(startTime.Unix()), 1000)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchAccountReservationUsage")
		errorResponse(c, fmt.Errorf("failed to fetch period records for account (%s): %w", accountId.Hex(), err))
		return
	}

	// Convert period records to response format
	records := make([]PeriodRecord, len(periodRecords))
	for i, record := range periodRecords {
		if record == nil {
			records[i] = PeriodRecord{
				ReservationPeriod: 0,
				Usage:             0,
			}
		} else {
			records[i] = PeriodRecord{
				ReservationPeriod: record.Index,
				Usage:             record.Usage,
			}
		}
	}

	response := &AccountReservationUsageResponse{
		AccountId: accountId.Hex(),
		Records:   records,
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchAccountReservationUsage")
	s.metrics.ObserveLatency("FetchAccountReservationUsage", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxBlobFeedAge))
	c.JSON(http.StatusOK, response)
}

// AccountReservationUsageResponse represents the response for account reservation usage
type AccountReservationUsageResponse struct {
	AccountId string         `json:"account_id"`
	Records   []PeriodRecord `json:"records"`
}

// PeriodRecord represents a single period's usage record
type PeriodRecord struct {
	ReservationPeriod uint32 `json:"reservation_period"`
	Usage             uint64 `json:"usage"`
}

type Reservation struct {
	SymbolsPerSecond uint64   `json:"symbols_per_second"`
	StartTimestamp   uint32   `json:"start_timestamp"`
	EndTimestamp     uint32   `json:"end_timestamp"`
	QuorumSplits     []uint32 `json:"quorum_splits"`
	QuorumNumbers    []uint32 `json:"quorum_numbers"`
}

type PaymentGlobalParams struct {
	GlobalSymbolsPerSecond uint64 `json:"global_symbols_per_second"`
	MinNumSymbols          uint64 `json:"min_num_symbols"`
	PricePerSymbol         uint64 `json:"price_per_symbol"`
	ReservationWindow      uint64 `json:"reservation_window"`
}
type AccountPaymentStateResponse struct {
	AccountId                string              `json:"account_id"`
	PaymentGlobalParams      PaymentGlobalParams `json:"payment_global_params"`
	PeriodRecords            []PeriodRecord      `json:"period_records"`
	Reservation              Reservation         `json:"reservation"`
	CumulativePayment        string              `json:"cumulative_payment"`
	OnchainCumulativePayment string              `json:"onchain_cumulative_payment"`
}

// FetchAccountPaymentState godoc
//
//	@Summary	Fetch payment state for an account
//	@Tags		Accounts
//	@Produce	json
//	@Param		account_id	path		string	true	"The account ID to fetch payment state for"
//	@Success	200			{object}	AccountPaymentStateResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/accounts/{account_id}/payment-state [get]
func (s *ServerV2) FetchAccountPaymentState(c *gin.Context) {
	handlerStart := time.Now()

	// Parse account ID
	accountStr := c.Param("account_id")
	if !gethcommon.IsHexAddress(accountStr) {
		s.metrics.IncrementInvalidArgRequestNum("FetchAccountPaymentState")
		invalidParamsErrorResponse(c, errors.New("account id is not valid hex"))
		return
	}
	accountId := gethcommon.HexToAddress(accountStr)
	if accountId == (gethcommon.Address{}) {
		s.metrics.IncrementInvalidArgRequestNum("FetchAccountPaymentState")
		invalidParamsErrorResponse(c, errors.New("zero account id is not valid"))
		return
	}

	if s.meterer == nil || s.meterer.ChainPaymentState == nil {
		errorResponse(c, fmt.Errorf("payment state is not available"))
		return
	}

	// on-chain global payment parameters
	globalSymbolsPerSecond := s.meterer.ChainPaymentState.GetGlobalSymbolsPerSecond()
	minNumSymbols := s.meterer.ChainPaymentState.GetMinNumSymbols()
	pricePerSymbol := s.meterer.ChainPaymentState.GetPricePerSymbol()
	reservationWindow := s.meterer.ChainPaymentState.GetReservationWindow()

	// off-chain account specific payment state
	now := time.Now()
	startTime := now.Add(-time.Duration(24) * time.Hour)

	periodRecords, err := s.meterer.MeteringStore.GetPeriodRecords(c.Request.Context(), accountId, uint64(startTime.Unix()), 1000)
	if err != nil {
		s.logger.Debug("failed to get reservation records, use placeholders", "err", err, "accountID", accountId)
	}
	var largestCumulativePaymentBytes []byte
	largestCumulativePayment, err := s.meterer.MeteringStore.GetLargestCumulativePayment(c.Request.Context(), accountId)
	if err != nil {
		s.logger.Debug("failed to get largest cumulative payment, use zero value", "err", err, "accountId", accountId)

	} else {
		largestCumulativePaymentBytes = largestCumulativePayment.Bytes()
	}

	// on-Chain account state
	var reservation *Reservation
	reservedPayment, err := s.meterer.ChainPaymentState.GetReservedPaymentByAccount(c.Request.Context(), accountId)
	if err != nil {
		s.logger.Debug("failed to get onchain reservation, use zero values", "err", err, "accountId", accountId)
	} else {
		quorumNumbers := make([]uint32, len(reservedPayment.QuorumNumbers))
		for i, v := range reservedPayment.QuorumNumbers {
			quorumNumbers[i] = uint32(v)
		}
		quorumSplits := make([]uint32, len(reservedPayment.QuorumSplits))
		for i, v := range reservedPayment.QuorumSplits {
			quorumSplits[i] = uint32(v)
		}

		reservation = &Reservation{
			SymbolsPerSecond: reservedPayment.SymbolsPerSecond,
			StartTimestamp:   uint32(reservedPayment.StartTimestamp),
			EndTimestamp:     uint32(reservedPayment.EndTimestamp),
			QuorumSplits:     quorumSplits,
			QuorumNumbers:    quorumNumbers,
		}
	}

	var onchainCumulativePaymentBytes []byte
	onDemandPayment, err := s.meterer.ChainPaymentState.GetOnDemandPaymentByAccount(c.Request.Context(), accountId)
	if err != nil {
		s.logger.Debug("failed to get ondemand payment, use zero value", "err", err, "accountId", accountId)
	} else {
		onchainCumulativePaymentBytes = onDemandPayment.CumulativePayment.Bytes()
	}

	paymentGlobalParams := PaymentGlobalParams{
		GlobalSymbolsPerSecond: globalSymbolsPerSecond,
		MinNumSymbols:          minNumSymbols,
		PricePerSymbol:         pricePerSymbol,
		ReservationWindow:      reservationWindow,
	}

	convertedRecords := make([]PeriodRecord, len(periodRecords))
	for i, pr := range periodRecords {
		if pr != nil {
			convertedRecords[i] = PeriodRecord{
				ReservationPeriod: pr.Index,
				Usage:             pr.Usage,
			}
		}
	}

	response := &AccountPaymentStateResponse{
		AccountId:                accountId.Hex(),
		PaymentGlobalParams:      paymentGlobalParams,
		PeriodRecords:            convertedRecords,
		Reservation:              *reservation,
		CumulativePayment:        string(largestCumulativePaymentBytes),
		OnchainCumulativePayment: string(onchainCumulativePaymentBytes),
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchAccountPaymentState")
	s.metrics.ObserveLatency("FetchAccountPaymentState", time.Since(handlerStart))
	c.JSON(http.StatusOK, response)
}
