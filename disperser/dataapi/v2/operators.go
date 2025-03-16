package v2

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// FetchOperatorSigningInfo godoc
//
//	@Summary	Fetch operators signing info
//	@Tags		Operators
//	@Produce	json
//	@Param		end				query		string	false	"Fetch operators signing info up to the end time (ISO 8601 format: 2006-01-02T15:04:05Z) [default: now]"
//	@Param		interval		query		int		false	"Fetch operators signing info starting from an interval (in seconds) before the end time [default: 3600]"
//	@Param		quorums			query		string	false	"Comma separated list of quorum IDs to fetch signing info for [default: 0,1]"
//	@Param		nonsigner_only	query		boolean	false	"Whether to only return operators with signing rate less than 100% [default: false]"
//	@Success	200				{object}	OperatorsSigningInfoResponse
//	@Failure	400				{object}	ErrorResponse	"error: Bad request"
//	@Failure	404				{object}	ErrorResponse	"error: Not found"
//	@Failure	500				{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/signing-info [get]
func (s *ServerV2) FetchOperatorSigningInfo(c *gin.Context) {
	handlerStart := time.Now()
	var err error

	now := handlerStart
	oldestTime := now.Add(-maxBlobAge)

	endTime := now
	if c.Query("end") != "" {
		endTime, err = time.Parse("2006-01-02T15:04:05Z", c.Query("end"))
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse end param: %w", err))
			return
		}
		if endTime.Before(oldestTime) {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
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
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse interval param: %w", err))
			return
		}
		if interval <= 0 {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
			invalidParamsErrorResponse(c, fmt.Errorf("interval must be greater than 0, found: %d", interval))
			return
		}
	}

	quorumStr := "0,1"
	if c.Query("quorums") != "" {
		quorumStr = c.Query("quorums")
	}
	quorums := strings.Split(quorumStr, ",")
	quorumsSeen := make(map[uint8]struct{}, 0)
	for _, idStr := range quorums {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
			invalidParamsErrorResponse(c, fmt.Errorf("failed to parse the provided quorum: %s", quorumStr))
			return
		}
		if id < 0 || id > maxQuorumIDAllowed {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorSigningInfo")
			invalidParamsErrorResponse(
				c, fmt.Errorf("the quorumID must be in range [0, %d], found: %d", maxQuorumIDAllowed, id),
			)
			return
		}
		quorumsSeen[uint8(id)] = struct{}{}
	}
	quorumIds := make([]uint8, 0, len(quorumsSeen))
	for q := range quorumsSeen {
		quorumIds = append(quorumIds, q)
	}

	nonsignerOnly := false
	if c.Query("nonsigner_only") != "" {
		nonsignerOnlyStr := c.Query("nonsigner_only")
		nonsignerOnly, err = strconv.ParseBool(nonsignerOnlyStr)
		if err != nil {
			invalidParamsErrorResponse(c, errors.New("the nonsigner_only param must be \"true\" or \"false\""))
			return
		}
	}

	startTime := endTime.Add(-time.Duration(interval) * time.Second)
	if startTime.Before(oldestTime) {
		startTime = oldestTime
	}

	attestations, err := s.batchFeedCache.Get(
		c.Request.Context(), startTime.Add(time.Nanosecond), endTime, Ascending, -1,
	)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchOperatorSigningInfo")
		errorResponse(c, fmt.Errorf("failed to fetch attestation feed from blob metadata store: %w", err))
		return
	}

	signingInfo, err := s.computeOperatorsSigningInfo(c.Request.Context(), attestations, quorumIds, nonsignerOnly)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchOperatorSigningInfo")
		errorResponse(c, fmt.Errorf("failed to compute the operators signing info: %w", err))
		return
	}
	startBlock, endBlock := computeBlockRange(attestations)
	response := OperatorsSigningInfoResponse{
		StartBlock:          startBlock,
		EndBlock:            endBlock,
		StartTimeUnixSec:    startTime.Unix(),
		EndTimeUnixSec:      endTime.Unix(),
		OperatorSigningInfo: signingInfo,
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchOperatorSigningInfo")
	s.metrics.ObserveLatency("FetchOperatorSigningInfo", time.Since(handlerStart))
	c.JSON(http.StatusOK, response)
}

// FetchOperatorsStake godoc
//
//	@Summary	Operator stake distribution query
//	@Tags		Operators
//	@Produce	json
//	@Param		operator_id	query		string	false	"Operator ID in hex string [default: all operators if unspecified]"
//	@Success	200			{object}	OperatorsStakeResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/stake [get]
func (s *ServerV2) FetchOperatorsStake(c *gin.Context) {
	handlerStart := time.Now()
	ctx := c.Request.Context()

	operatorId := c.DefaultQuery("operator_id", "")
	s.logger.Info("getting operators stake distribution", "operatorId", operatorId)

	currentBlock, err := s.indexedChainState.GetCurrentBlockNumber(c.Request.Context())
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchOperatorsStake")
		errorResponse(c, fmt.Errorf("failed to get current block number: %w", err))
		return
	}
	operatorsStakeResponse, err := s.operatorHandler.GetOperatorsStakeAtBlock(ctx, operatorId, uint32(currentBlock))
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchOperatorsStake")
		errorResponse(c, fmt.Errorf("failed to get operator stake: %w", err))
		return
	}
	operatorsStakeResponse.CurrentBlock = uint32(currentBlock)

	// Get operators' addresses in batch
	operatorsSeen := make(map[string]struct{}, 0)
	for _, ops := range operatorsStakeResponse.StakeRankedOperators {
		for _, op := range ops {
			operatorsSeen[op.OperatorId] = struct{}{}
		}
	}
	operatorIDs := make([]core.OperatorID, 0)
	for id := range operatorsSeen {
		opId, err := core.OperatorIDFromHex(id)
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchOperatorsStake")
			errorResponse(c, fmt.Errorf("malformed operator ID: %w", err))
			return
		}
		operatorIDs = append(operatorIDs, opId)
	}
	// Get the address for the operators.
	// operatorAddresses[i] is the address for operatorIDs[i].
	operatorAddresses, err := s.chainReader.BatchOperatorIDToAddress(ctx, operatorIDs)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("FetchOperatorsStake")
		errorResponse(c, fmt.Errorf("failed to get operator addresses from IDs: %w", err))
		return
	}
	idToAddress := make(map[string]string, 0)
	for i := range operatorIDs {
		idToAddress[operatorIDs[i].Hex()] = operatorAddresses[i].Hex()
	}
	for _, ops := range operatorsStakeResponse.StakeRankedOperators {
		for _, op := range ops {
			op.OperatorAddress = idToAddress[op.OperatorId]
		}
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchOperatorsStake")
	s.metrics.ObserveLatency("FetchOperatorsStake", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxOperatorsStakeAge))
	c.JSON(http.StatusOK, operatorsStakeResponse)
}

// FetchOperatorsNodeInfo godoc
//
//	@Summary	Active operator semver
//	@Tags		Operators
//	@Produce	json
//	@Success	200	{object}	SemverReportResponse
//	@Failure	500	{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/node-info [get]
func (s *ServerV2) FetchOperatorsNodeInfo(c *gin.Context) {
	handlerStart := time.Now()

	report, err := s.operatorHandler.ScanOperatorsHostInfo(c.Request.Context())
	if err != nil {
		s.logger.Error("failed to scan operators host info", "error", err)
		s.metrics.IncrementFailedRequestNum("FetchOperatorsNodeInfo")
		errorResponse(c, err)
	}

	s.metrics.IncrementSuccessfulRequestNum("FetchOperatorsNodeInfo")
	s.metrics.ObserveLatency("FetchOperatorsNodeInfo", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxOperatorPortCheckAge))
	c.JSON(http.StatusOK, report)
}

// FetchOperatorsResponses godoc
//
//	@Summary	Fetch operator attestation response for a batch
//	@Tags		Operators
//	@Produce	json
//	@Param		batch_header_hash	path		string	true	"Batch header hash in hex string"
//	@Param		operator_id			query		string	false	"Operator ID in hex string [default: all operators if unspecified]"
//	@Success	200					{object}	OperatorDispersalResponses
//	@Failure	400					{object}	ErrorResponse	"error: Bad request"
//	@Failure	404					{object}	ErrorResponse	"error: Not found"
//	@Failure	500					{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/{batch_header_hash} [get]
func (s *ServerV2) FetchOperatorsResponses(c *gin.Context) {
	handlerStart := time.Now()

	batchHeaderHashHex := c.Param("batch_header_hash")
	batchHeaderHash, err := dataapi.ConvertHexadecimalToBytes([]byte(batchHeaderHashHex))
	if err != nil {
		s.metrics.IncrementInvalidArgRequestNum("FetchOperatorsResponses")
		errorResponse(c, errors.New("invalid batch header hash"))
		return
	}
	operatorIdStr := c.DefaultQuery("operator_id", "")

	operatorResponses := make([]*corev2.DispersalResponse, 0)
	if operatorIdStr == "" {
		res, err := s.blobMetadataStore.GetDispersalResponses(c.Request.Context(), batchHeaderHash)
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchOperatorsResponses")
			errorResponse(c, err)
			return
		}
		operatorResponses = append(operatorResponses, res...)
	} else {
		operatorId, err := core.OperatorIDFromHex(operatorIdStr)
		if err != nil {
			s.metrics.IncrementInvalidArgRequestNum("FetchOperatorsResponses")
			errorResponse(c, errors.New("invalid operatorId"))
			return
		}

		res, err := s.blobMetadataStore.GetDispersalResponse(c.Request.Context(), batchHeaderHash, operatorId)
		if err != nil {
			s.metrics.IncrementFailedRequestNum("FetchOperatorsResponses")
			errorResponse(c, err)
			return
		}
		operatorResponses = append(operatorResponses, res)
	}
	response := &OperatorDispersalResponses{
		Responses: operatorResponses,
	}
	s.metrics.IncrementSuccessfulRequestNum("FetchOperatorsResponses")
	s.metrics.ObserveLatency("FetchOperatorsResponses", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxOperatorResponseAge))
	c.JSON(http.StatusOK, response)
}

// CheckOperatorsLiveness godoc
//
//	@Summary	Check operator v2 node liveness
//	@Tags		Operators
//	@Produce	json
//	@Param		operator_id	query		string	false	"Operator ID in hex string [default: all operators if unspecified]"
//	@Success	200			{object}	OperatorLivenessResponse
//	@Failure	400			{object}	ErrorResponse	"error: Bad request"
//	@Failure	404			{object}	ErrorResponse	"error: Not found"
//	@Failure	500			{object}	ErrorResponse	"error: Server error"
//	@Router		/operators/liveness [get]
func (s *ServerV2) CheckOperatorsLiveness(c *gin.Context) {
	handlerStart := time.Now()

	operatorId := c.DefaultQuery("operator_id", "")
	s.logger.Info("checking operator ports", "operatorId", operatorId)
	portCheckResponse, err := s.operatorHandler.ProbeV2OperatorPorts(c.Request.Context(), operatorId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			err = errNotFound
			s.logger.Warn("operator not found", "operatorId", operatorId)
			s.metrics.IncrementNotFoundRequestNum("CheckOperatorsLiveness")
		} else {
			s.logger.Error("operator port check failed", "error", err)
			s.metrics.IncrementFailedRequestNum("CheckOperatorsLiveness")
		}
		errorResponse(c, err)
		return
	}

	s.metrics.IncrementSuccessfulRequestNum("CheckOperatorsLiveness")
	s.metrics.ObserveLatency("CheckOperatorsLiveness", time.Since(handlerStart))
	c.Writer.Header().Set(cacheControlParam, fmt.Sprintf("max-age=%d", maxOperatorPortCheckAge))
	c.JSON(http.StatusOK, portCheckResponse)
}

func (s *ServerV2) computeOperatorsSigningInfo(
	ctx context.Context,
	attestations []*corev2.Attestation,
	quorumIDs []uint8,
	nonsignerOnly bool,
) ([]*OperatorSigningInfo, error) {
	if len(attestations) == 0 {
		return nil, errors.New("no attestations to compute signing info")
	}

	// Compute the block number range [startBlock, endBlock] (both inclusive) when the
	// attestations have happened.
	startBlock, endBlock := computeBlockRange(attestations)

	// Get quorum change events in range [startBlock+1, endBlock].
	// We don't need the events at startBlock because we'll fetch all active operators and
	// quorums at startBlock.
	operatorQuorumEvents, err := s.subgraphClient.QueryOperatorQuorumEvent(ctx, startBlock+1, endBlock)
	if err != nil {
		return nil, err
	}

	// Get operators of interest to compute signing info, which includes:
	// - operators that were active at startBlock
	// - operators that joined after startBlock
	operatorList, err := s.getOperatorsOfInterest(
		ctx, startBlock, endBlock, quorumIDs, operatorQuorumEvents,
	)
	if err != nil {
		return nil, err
	}

	// Create operators' quorum intervals: OperatorQuorumIntervals[op][q] is a sequence of
	// increasing and non-overlapping block intervals during which the operator "op" is
	// registered in quorum "q".
	operatorQuorumIntervals, _, err := s.operatorHandler.CreateOperatorQuorumIntervals(
		ctx, operatorList, operatorQuorumEvents, startBlock, endBlock,
	)
	if err != nil {
		return nil, err
	}

	// Compute num batches failed, where numFailed[op][q] is the number of batches
	// failed to sign for quorum "q" by operator "op".
	numFailed := computeNumFailed(attestations, operatorQuorumIntervals)

	// Compute num batches responsible, where numResponsible[op][q] is the number of batches
	// that operator "op" are responsible for in quorum "q".
	numResponsible := computeNumResponsible(attestations, operatorQuorumIntervals)

	totalNumBatchesPerQuorum := computeTotalNumBatchesPerQuorum(attestations)

	state, err := s.chainState.GetOperatorState(ctx, uint(endBlock), quorumIDs)
	if err != nil {
		return nil, err
	}
	signingInfo := make([]*OperatorSigningInfo, 0)
	for _, op := range operatorList.GetOperatorIds() {
		for _, q := range quorumIDs {
			operatorId := op.Hex()

			numShouldHaveSigned := 0
			if num, exist := safeAccess(numResponsible, operatorId, q); exist {
				numShouldHaveSigned = num
			}
			// The operator op received no batch that it should sign.
			if numShouldHaveSigned == 0 {
				continue
			}

			numFailedToSign := 0
			if num, exist := safeAccess(numFailed, operatorId, q); exist {
				numFailedToSign = num
			}

			if nonsignerOnly && numFailedToSign == 0 {
				continue
			}

			operatorAddress, ok := operatorList.GetAddress(operatorId)
			if !ok {
				// This should never happen (becuase OperatorList ensures the 1:1 mapping
				// between ID and address), but we don't fail the entire request, just
				// mark internal error for the address field to signal the issue.
				operatorAddress = "Unexpected internal error"
				s.logger.Error("Internal error: failed to find address for operatorId", "operatorId", operatorId)
			}

			// Signing percentage with 8 decimal (e.g. 95.75000000, which means 95.75%).
			// We need 8 decimal because if there is one attestation per second, then we
			// need to have resolution 1/(3600*24*14), which is 8.26719577e-7. At this
			// resolution we can capture the signing rate difference caused by 1 unsigned
			// batch.
			signingPercentage := math.Round(
				(float64(numShouldHaveSigned-numFailedToSign)/float64(numShouldHaveSigned))*100*1e8,
			) / 1e8

			stakePercentage := float64(0)
			if stake, ok := state.Operators[q][op]; ok {
				totalStake := new(big.Float).SetInt(state.Totals[q].Stake)
				stakeRatio := new(big.Float).Quo(
					new(big.Float).SetInt(stake.Stake),
					totalStake,
				)
				stakeRatio.Mul(stakeRatio, big.NewFloat(100))
				stakePercentage, _ = stakeRatio.Float64()
			}

			si := &OperatorSigningInfo{
				OperatorId:              operatorId,
				OperatorAddress:         operatorAddress,
				QuorumId:                q,
				TotalUnsignedBatches:    numFailedToSign,
				TotalResponsibleBatches: numShouldHaveSigned,
				TotalBatches:            totalNumBatchesPerQuorum[q],
				SigningPercentage:       signingPercentage,
				StakePercentage:         stakePercentage,
			}
			signingInfo = append(signingInfo, si)
		}
	}

	// Sort by descending order of signing rate and then ascending order of <quorumId, operatorId>.
	sort.Slice(signingInfo, func(i, j int) bool {
		if signingInfo[i].SigningPercentage == signingInfo[j].SigningPercentage {
			if signingInfo[i].OperatorId == signingInfo[j].OperatorId {
				return signingInfo[i].QuorumId < signingInfo[j].QuorumId
			}
			return signingInfo[i].OperatorId < signingInfo[j].OperatorId
		}
		return signingInfo[i].SigningPercentage > signingInfo[j].SigningPercentage
	})

	return signingInfo, nil
}

// getOperatorsOfInterest returns operators that we want to compute signing info for.
//
// This contains two parts:
// - the operators that were active at the startBlock
// - the operators that joined after startBlock
func (s *ServerV2) getOperatorsOfInterest(
	ctx context.Context,
	startBlock, endBlock uint32,
	quorumIDs []uint8,
	operatorQuorumEvents *dataapi.OperatorQuorumEvents,
) (*dataapi.OperatorList, error) {
	operatorList := dataapi.NewOperatorList()

	// The first part: active operators at startBlock
	operatorsByQuorum, err := s.chainReader.GetOperatorStakesForQuorums(ctx, quorumIDs, startBlock)
	if err != nil {
		return nil, err
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
	operatorAddresses, err := s.chainReader.BatchOperatorIDToAddress(ctx, operatorIDs)
	if err != nil {
		return nil, err
	}
	for i := range operatorIDs {
		operatorList.Add(operatorIDs[i], operatorAddresses[i].Hex())
	}

	// The second part: new operators after startBlock.
	newAddresses := make(map[string]struct{}, 0)
	for op := range operatorQuorumEvents.AddedToQuorum {
		if _, exist := operatorList.GetID(op); !exist {
			newAddresses[op] = struct{}{}
		}
	}
	for op := range operatorQuorumEvents.RemovedFromQuorum {
		if _, exist := operatorList.GetID(op); !exist {
			newAddresses[op] = struct{}{}
		}
	}
	addresses := make([]gethcommon.Address, 0, len(newAddresses))
	for addr := range newAddresses {
		addresses = append(addresses, gethcommon.HexToAddress(addr))
	}
	operatorIds, err := s.chainReader.BatchOperatorAddressToID(ctx, addresses)
	if err != nil {
		return nil, err
	}
	// We merge the new operators observed in AddedToQuorum and RemovedFromQuorum
	// into the operator set.
	for i := 0; i < len(operatorIds); i++ {
		operatorList.Add(operatorIds[i], addresses[i].Hex())
	}

	return operatorList, nil
}

func computeNumFailed(
	attestations []*corev2.Attestation,
	operatorQuorumIntervals dataapi.OperatorQuorumIntervals,
) map[string]map[uint8]int {
	numFailed := make(map[string]map[uint8]int)
	for _, at := range attestations {
		for _, pubkey := range at.NonSignerPubKeys {
			op := pubkey.GetOperatorID().Hex()
			// Note: avg number of quorums per operator is a small number, so use brute
			// force here (otherwise, we can create a map to make it more efficient)
			for _, operatorQuorum := range operatorQuorumIntervals.GetQuorums(
				op,
				uint32(at.ReferenceBlockNumber),
			) {
				for _, batchQuorum := range at.QuorumNumbers {
					if operatorQuorum == batchQuorum {
						if _, ok := numFailed[op]; !ok {
							numFailed[op] = make(map[uint8]int)
						}
						numFailed[op][operatorQuorum]++
						break
					}
				}
			}
		}
	}
	return numFailed
}

func computeNumResponsible(
	attestations []*corev2.Attestation,
	operatorQuorumIntervals dataapi.OperatorQuorumIntervals,
) map[string]map[uint8]int {
	// Create quorumBatches, where quorumBatches[q].AccuBatches is the total number of
	// batches in block interval [startBlock, b] for quorum "q".
	quorumBatches := dataapi.CreatQuorumBatches(dataapi.CreateQuorumBatchMapV2(attestations))

	numResponsible := make(map[string]map[uint8]int)
	for op, val := range operatorQuorumIntervals {
		if _, ok := numResponsible[op]; !ok {
			numResponsible[op] = make(map[uint8]int)
		}
		for q, intervals := range val {
			numBatches := 0
			if _, ok := quorumBatches[q]; ok {
				for _, interval := range intervals {
					numBatches += dataapi.ComputeNumBatches(
						quorumBatches[q], interval.StartBlock, interval.EndBlock,
					)
				}
			}
			numResponsible[op][q] = numBatches
		}
	}

	return numResponsible
}

func computeTotalNumBatchesPerQuorum(attestations []*corev2.Attestation) map[uint8]int {
	numBatchesPerQuorum := make(map[uint8]int)
	for _, at := range attestations {
		for _, q := range at.QuorumNumbers {
			numBatchesPerQuorum[q]++
		}
	}
	return numBatchesPerQuorum
}

func computeBlockRange(attestations []*corev2.Attestation) (uint32, uint32) {
	if len(attestations) == 0 {
		return 0, 0
	}
	startBlock := attestations[0].ReferenceBlockNumber
	endBlock := attestations[0].ReferenceBlockNumber
	for i := range attestations {
		if startBlock > attestations[i].ReferenceBlockNumber {
			startBlock = attestations[i].ReferenceBlockNumber
		}
		if endBlock < attestations[i].ReferenceBlockNumber {
			endBlock = attestations[i].ReferenceBlockNumber
		}
	}
	return uint32(startBlock), uint32(endBlock)
}
