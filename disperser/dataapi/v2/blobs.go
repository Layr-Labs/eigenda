package v2

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/gin-gonic/gin"
)

// FetchBlobFeed godoc
//
//	@Summary	Fetch blob feed in specified direction
//	@Tags		Blobs
//	@Produce	json
//	@Param		direction	query		string	false	"Direction to fetch: 'forward' (oldest to newest, ASC order) or 'backward' (newest to oldest, DESC order) [default: forward]"
//	@Param		before		query		string	false	"Fetch blobs before this time, exclusive (ISO 8601 format, example: 2006-01-02T15:04:05Z) [default: now]"
//	@Param		after		query		string	false	"Fetch blobs after this time, exclusive (ISO 8601 format, example: 2006-01-02T15:04:05Z); must be smaller than `before` [default: before-1h]"
//	@Param		cursor		query		string	false	"Pagination cursor (opaque string from previous response); for 'forward' direction, overrides `after` and fetches blobs from `cursor` to `before`; for 'backward' direction, overrides `before` and fetches blobs from `cursor` to `after` (all bounds exclusive) [default: empty]"
//	@Param		limit		query		int		false	"Maximum number of blobs to return; if limit <= 0 or >1000, it's treated as 1000 [default: 20; max: 1000]"
//	@Success	200			{object}	BlobFeedResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/blobs/feed [get]
func (s *ServerV2) FetchBlobFeed(c *gin.Context) {
	handlerStart := time.Now()
	var err error

	// Validate direction
	direction := "forward"
	if dirStr := c.Query("direction"); dirStr != "" {
		if dirStr != "forward" && dirStr != "backward" {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
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
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse `before` param: %w", err))
			return
		}
		if beforeTime.Before(oldestTime) {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("`before` time cannot be more than 14 days in the past, found: %q", c.Query("before")))
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
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse `after` param: %w", err))
			return
		}
		if now.Before(afterTime) {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("`after` must be before current time, found: %q", c.Query("after")))
			return
		}
		if afterTime.Before(oldestTime) {
			afterTime = oldestTime
		}
	}

	// Validate time range
	if !afterTime.Before(beforeTime) {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
		invalidParamsErrorResponse(c, fmt.Errorf("`after` timestamp (%q) must be earlier than `before` timestamp (%q)",
			afterTime.Format(time.RFC3339), beforeTime.Format(time.RFC3339)))
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
		invalidParamsErrorResponse(c, fmt.Errorf("failed to parse `limit` param: %w", err))
		return
	}
	if limit <= 0 || limit > maxNumBlobsPerBlobFeedResponse {
		limit = maxNumBlobsPerBlobFeedResponse
	}

	// Convert times to cursors
	afterCursor := blobstore.BlobFeedCursor{
		RequestedAt: uint64(afterTime.UnixNano()),
	}
	beforeCursor := blobstore.BlobFeedCursor{
		RequestedAt: uint64(beforeTime.UnixNano()),
	}

	current := blobstore.BlobFeedCursor{
		RequestedAt: 0,
	}
	// Handle cursor if provided
	if cursorStr := c.Query("cursor"); cursorStr != "" {
		cursor, err := new(blobstore.BlobFeedCursor).FromCursorKey(cursorStr)
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchBlobFeed")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse the `cursor`: %w", err))
			return
		}
		current = *cursor
	}

	var blobs []*v2.BlobMetadata
	var nextCursor *blobstore.BlobFeedCursor

	if direction == "forward" {
		startCursor := afterCursor
		// The presence of `cursor` param will override the `after` param
		if current.RequestedAt > 0 {
			startCursor = current
		}
		blobs, nextCursor, err = s.blobMetadataStore.GetBlobMetadataByRequestedAtForward(
			c.Request.Context(),
			startCursor,
			beforeCursor,
			limit,
		)
	} else {
		endCursor := beforeCursor
		// The presence of `cursor` param will override the `before` param
		if current.RequestedAt > 0 {
			endCursor = current
		}
		blobs, nextCursor, err = s.blobMetadataStore.GetBlobMetadataByRequestedAtBackward(
			c.Request.Context(),
			endCursor,
			afterCursor,
			limit,
		)
	}

	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchBlobFeed")
		errorResponse(c, fmt.Errorf("failed to fetch feed from blob metadata store: %w", err))
		return
	}

	s.sendBlobFeedResponse(c, blobs, nextCursor, handlerStart)
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
	metadata, cached := s.blobMetadataCache.Get(blobKey.Hex())
	if !cached {
		metadata, err = s.blobMetadataStore.GetBlobMetadata(c.Request.Context(), blobKey)
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchBlob")
			errorResponse(c, err)
			return
		}
		s.blobMetadataCache.Add(blobKey.Hex(), metadata)
	} else {
		s.metrics.IncrementCacheHit("FetchBlob")
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
//	@Summary	Fetch blob certificate by blob key
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
	cert, cached := s.blobCertificateCache.Get(blobKey.Hex())
	if !cached {
		cert, _, err = s.blobMetadataStore.GetBlobCertificate(c.Request.Context(), blobKey)
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchBlobCertificate")
			errorResponse(c, err)
			return
		}
		s.blobCertificateCache.Add(blobKey.Hex(), cert)
	} else {
		s.metrics.IncrementCacheHit("FetchBlobCertificate")
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

	ctx := c.Request.Context()
	blobKey, err := corev2.HexToBlobKey(c.Param("blob_key"))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchBlobAttestationInfo")
		invalidParamsErrorResponse(c, fmt.Errorf("failed to parse blob_key param: %w", err))
		return
	}

	response, cached := s.blobAttestationInfoResponseCache.Get(blobKey.Hex())
	if !cached {
		response, err = s.getBlobAttestationInfoResponse(ctx, blobKey)
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchBlobAttestationInfo")
			errorResponse(c, err)
			return
		}
		s.blobAttestationInfoResponseCache.Add(blobKey.Hex(), response)
	} else {
		s.metrics.IncrementCacheHit("FetchBlobAttestationInfo")
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchBlobAttestationInfo")
	s.metrics.ObserveLatency("FetchBlobAttestationInfo", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	c.JSON(http.StatusOK, response)
}

func (s *ServerV2) getBlobAttestationInfoResponse(ctx context.Context, blobKey corev2.BlobKey) (*BlobAttestationInfoResponse, error) {
	var err error
	attestationInfo, cached := s.blobAttestationInfoCache.Get(blobKey.Hex())
	if !cached {
		attestationInfo, err = s.blobMetadataStore.GetBlobAttestationInfo(ctx, blobKey)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch blob attestation info: %w", err)
		}
		s.blobAttestationInfoCache.Add(blobKey.Hex(), attestationInfo)
	}

	batchHeaderHash, err := attestationInfo.InclusionInfo.BatchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to get batch header hash from blob inclusion info:        %w", err)
	}

	// Get quorums that this blob was dispersed to
	metadata, cached := s.blobMetadataCache.Get(blobKey.Hex())
	if !cached {
		metadata, err = s.blobMetadataStore.GetBlobMetadata(ctx, blobKey)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch blob metadata: %w", err)
		}
		s.blobMetadataCache.Add(blobKey.Hex(), metadata)
	}

	blobQuorums := make(map[uint8]struct{}, 0)
	for _, q := range metadata.BlobHeader.QuorumNumbers {
		blobQuorums[q] = struct{}{}
	}

	// Get all operators for the attestation
	operatorList, operatorsByQuorum, err := s.getAllOperatorsForAttestation(ctx, attestationInfo.Attestation)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch operators at reference block number: %w", err)
	}

	// Get all nonsigners (of the batch that this blob is part of)
	nonsigners := make(map[core.OperatorID]struct{}, 0)
	for i := 0; i < len(attestationInfo.Attestation.NonSignerPubKeys); i++ {
		opId := attestationInfo.Attestation.NonSignerPubKeys[i].GetOperatorID()
		nonsigners[opId] = struct{}{}
	}

	// Compute the signers and nonsigners for the blob, for each quorum that the blob was dispersed to
	blobSigners := make(map[uint8][]OperatorIdentity, 0)
	blobNonsigners := make(map[uint8][]OperatorIdentity, 0)
	for q, innerMap := range operatorsByQuorum {
		// Make sure the blob was dispersed to the quorum
		if _, exist := blobQuorums[q]; !exist {
			continue
		}
		for _, op := range innerMap {
			id := op.OperatorID.Hex()
			addr, exist := operatorList.GetAddress(id)
			// This should never happen becuase OperatorList ensures the 1:1 mapping
			if !exist {
				addr = "Unexpected internal error"
				s.logger.Error("Internal error: failed to find address for operatorId", "operatorId", op.OperatorID.Hex())
			}
			if _, exist := nonsigners[op.OperatorID]; exist {
				blobNonsigners[q] = append(blobNonsigners[q], OperatorIdentity{
					OperatorId:      id,
					OperatorAddress: addr,
				})
			} else {
				blobSigners[q] = append(blobSigners[q], OperatorIdentity{
					OperatorId:      id,
					OperatorAddress: addr,
				})
			}
		}
	}

	return &BlobAttestationInfoResponse{
		BlobKey:         blobKey.Hex(),
		BatchHeaderHash: hex.EncodeToString(batchHeaderHash[:]),
		InclusionInfo:   attestationInfo.InclusionInfo,
		AttestationInfo: &AttestationInfo{
			Attestation: attestationInfo.Attestation,
			Signers:     blobSigners,
			Nonsigners:  blobNonsigners,
		},
	}, nil
}

func (s *ServerV2) getAllOperatorsForAttestation(ctx context.Context, attestation *corev2.Attestation) (*dataapi.OperatorList, core.OperatorStakes, error) {
	rbn := attestation.ReferenceBlockNumber
	operatorsByQuorum, err := s.chainReader.GetOperatorStakesForQuorums(ctx, attestation.QuorumNumbers, uint32(rbn))
	if err != nil {
		return nil, nil, err
	}

	operatorsSeen := make(map[core.OperatorID]struct{}, 0)
	for _, ops := range operatorsByQuorum {
		for _, op := range ops {
			operatorsSeen[op.OperatorID] = struct{}{}
		}
	}
	operatorIDs := make([]core.OperatorID, 0)
	for id := range operatorsSeen {
		operatorIDs = append(operatorIDs, id)
	}

	// Get the address for the operators.
	// operatorAddresses[i] is the address for operatorIDs[i].
	operatorList := dataapi.NewOperatorList()
	operatorAddresses, err := s.chainReader.BatchOperatorIDToAddress(ctx, operatorIDs)
	if err != nil {
		return nil, nil, err
	}
	for i := range operatorIDs {
		operatorList.Add(operatorIDs[i], operatorAddresses[i].Hex())
	}

	return operatorList, operatorsByQuorum, nil
}

func (s *ServerV2) sendBlobFeedResponse(
	c *gin.Context,
	blobs []*v2.BlobMetadata,
	nextCursor *blobstore.BlobFeedCursor,
	handlerStart time.Time,
) {
	cursorStr := ""
	if nextCursor != nil {
		cursorStr = nextCursor.ToCursorKey()
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
		Blobs:  blobInfo,
		Cursor: cursorStr,
	}
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxFeedBlobAge))
	s.metrics.IncrementSuccessfulRequestNum("FetchBlobFeed")
	s.metrics.ObserveLatency("FetchBlobFeed", time.Since(handlerStart))
	c.JSON(http.StatusOK, response)
}
