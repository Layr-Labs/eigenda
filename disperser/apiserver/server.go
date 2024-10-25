package apiserver

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/common"
	healthcheck "github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

const systemAccountKey = "system"

type DispersalServer struct {
	pb.UnimplementedDisperserServer
	mu *sync.RWMutex

	serverConfig disperser.ServerConfig
	rateConfig   RateConfig

	blobStore    disperser.BlobStore
	tx           core.Reader
	quorumConfig QuorumConfig

	ratelimiter   common.RateLimiter
	authenticator core.BlobRequestAuthenticator

	metrics *disperser.Metrics

	maxBlobSize int

	logger logging.Logger
}

type QuorumConfig struct {
	RequiredQuorums []core.QuorumID
	QuorumCount     uint8
	SecurityParams  map[core.QuorumID]core.SecurityParam
}

// NewServer creates a new Server struct with the provided parameters.
//
// Note: The Server's chunks store will be created at config.DbPath+"/chunk".
func NewDispersalServer(
	serverConfig disperser.ServerConfig,
	store disperser.BlobStore,
	tx core.Reader,
	_logger logging.Logger,
	metrics *disperser.Metrics,
	ratelimiter common.RateLimiter,
	rateConfig RateConfig,
	maxBlobSize int,
) *DispersalServer {
	logger := _logger.With("component", "DispersalServer")
	for account, rateInfoByQuorum := range rateConfig.Allowlist {
		for quorumID, rateInfo := range rateInfoByQuorum {
			logger.Info("[Allowlist]", "account", account, "name", rateInfo.Name, "quorumID", quorumID, "throughput", rateInfo.Throughput, "blobRate", rateInfo.BlobRate)
		}
	}
	logger.Info("allowlist config", "file", rateConfig.AllowlistFile, "refreshInterval", rateConfig.AllowlistRefreshInterval.String())

	authenticator := auth.NewAuthenticator(auth.AuthConfig{})

	return &DispersalServer{
		serverConfig:  serverConfig,
		rateConfig:    rateConfig,
		blobStore:     store,
		tx:            tx,
		metrics:       metrics,
		logger:        logger,
		ratelimiter:   ratelimiter,
		authenticator: authenticator,
		mu:            &sync.RWMutex{},
		quorumConfig:  QuorumConfig{},
		maxBlobSize:   maxBlobSize,
	}
}

func (s *DispersalServer) DisperseBlobAuthenticated(stream pb.Disperser_DisperseBlobAuthenticatedServer) error {

	// This uses the existing deadline of stream.Context() if it is earlier.
	ctx, cancel := context.WithTimeout(stream.Context(), s.serverConfig.GrpcTimeout)
	defer cancel()

	// Process disperse_request
	in, err := stream.Recv()
	if err != nil {
		s.metrics.HandleInvalidArgRpcRequest("DisperseBlobAuthenticated")
		s.metrics.HandleInvalidArgRequest("DisperseBlobAuthenticated")
		return api.NewErrorInvalidArg(fmt.Sprintf("error receiving next message: %v", err))
	}

	request, ok := in.GetPayload().(*pb.AuthenticatedRequest_DisperseRequest)
	if !ok {
		s.metrics.HandleInvalidArgRpcRequest("DisperseBlobAuthenticated")
		s.metrics.HandleInvalidArgRequest("DisperseBlobAuthenticated")
		return api.NewErrorInvalidArg("missing DisperseBlobRequest")
	}

	blob, err := s.validateRequestAndGetBlob(ctx, request.DisperseRequest)
	if err != nil {
		for _, quorumID := range request.DisperseRequest.CustomQuorumNumbers {
			s.metrics.HandleFailedRequest(codes.InvalidArgument.String(), fmt.Sprint(quorumID), len(request.DisperseRequest.GetData()), "DisperseBlobAuthenticated")
		}
		s.metrics.HandleInvalidArgRpcRequest("DisperseBlobAuthenticated")
		return api.NewErrorInvalidArg(err.Error())
	}

	// Get the ethereum address associated with the public key. This is just for convenience so we can put addresses instead of public keys in the allowlist.
	// Decode public key
	publicKeyBytes, err := hexutil.Decode(blob.RequestHeader.AccountID)
	if err != nil {
		s.metrics.HandleInvalidArgRpcRequest("DisperseBlobAuthenticated")
		s.metrics.HandleInvalidArgRequest("DisperseBlobAuthenticated")
		return api.NewErrorInvalidArg(fmt.Sprintf("failed to decode account ID (%v): %v", blob.RequestHeader.AccountID, err))
	}

	pubKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		s.metrics.HandleInvalidArgRpcRequest("DisperseBlobAuthenticated")
		s.metrics.HandleInvalidArgRequest("DisperseBlobAuthenticated")
		return api.NewErrorInvalidArg(fmt.Sprintf("failed to decode public key (%v): %v", hexutil.Encode(publicKeyBytes), err))
	}

	authenticatedAddress := crypto.PubkeyToAddress(*pubKey).String()

	// Send back challenge to client
	challengeBytes := make([]byte, 32)
	_, err = rand.Read(challengeBytes)
	if err != nil {
		s.metrics.HandleInvalidArgRpcRequest("DisperseBlobAuthenticated")
		s.metrics.HandleInvalidArgRequest("DisperseBlobAuthenticated")
		return api.NewErrorInvalidArg(fmt.Sprintf("failed to generate challenge: %v", err))
	}
	challenge := binary.LittleEndian.Uint32(challengeBytes)
	err = stream.Send(&pb.AuthenticatedReply{Payload: &pb.AuthenticatedReply_BlobAuthHeader{
		BlobAuthHeader: &pb.BlobAuthHeader{
			ChallengeParameter: challenge,
		},
	}})
	if err != nil {
		return err
	}

	// Create a channel for the result of stream.Recv()
	resultCh := make(chan *pb.AuthenticatedRequest)
	errCh := make(chan error)

	// Run stream.Recv() in a goroutine
	go func() {
		in, err := stream.Recv()
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- in
	}()

	// Use select to wait on either the result of stream.Recv() or the context being done
	select {
	case in = <-resultCh:
	case err := <-errCh:
		s.metrics.HandleInvalidArgRpcRequest("DisperseBlobAuthenticated")
		s.metrics.HandleInvalidArgRequest("DisperseBlobAuthenticated")
		return api.NewErrorInvalidArg(fmt.Sprintf("error receiving next message: %v", err))
	case <-ctx.Done():
		s.metrics.HandleInvalidArgRpcRequest("DisperseBlobAuthenticated")
		s.metrics.HandleInvalidArgRequest("DisperseBlobAuthenticated")
		return api.NewErrorInvalidArg("context deadline exceeded")
	}

	challengeReply, ok := in.GetPayload().(*pb.AuthenticatedRequest_AuthenticationData)
	if !ok {
		s.metrics.HandleInvalidArgRpcRequest("DisperseBlobAuthenticated")
		s.metrics.HandleInvalidArgRequest("DisperseBlobAuthenticated")
		return api.NewErrorInvalidArg("expected AuthenticationData")
	}

	blob.RequestHeader.Nonce = challenge
	blob.RequestHeader.AuthenticationData = challengeReply.AuthenticationData.GetAuthenticationData()

	err = s.authenticator.AuthenticateBlobRequest(blob.RequestHeader.BlobAuthHeader)
	if err != nil {
		s.metrics.HandleInvalidArgRpcRequest("DisperseBlobAuthenticated")
		s.metrics.HandleInvalidArgRequest("DisperseBlobAuthenticated")
		return api.NewErrorInvalidArg(fmt.Sprintf("failed to authenticate blob request: %v", err))
	}

	// Disperse the blob
	reply, err := s.disperseBlob(ctx, blob, authenticatedAddress, "DisperseBlobAuthenticated")
	if err != nil {
		// Note the disperseBlob already updated metrics for this error.
		s.logger.Info("failed to disperse blob", "err", err)
		return err
	}

	// Send back disperse_reply
	err = stream.Send(&pb.AuthenticatedReply{Payload: &pb.AuthenticatedReply_DisperseReply{
		DisperseReply: reply,
	}})
	if err != nil {
		s.logger.Error("failed to stream back DisperseReply", "err", err)
		return err
	}

	s.metrics.HandleSuccessfulRpcRequest("DisperseBlobAuthenticated")

	return nil

}

func (s *DispersalServer) DisperseBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*pb.DisperseBlobReply, error) {
	blob, err := s.validateRequestAndGetBlob(ctx, req)
	if err != nil {
		for _, quorumID := range req.CustomQuorumNumbers {
			s.metrics.HandleFailedRequest(codes.InvalidArgument.String(), fmt.Sprint(quorumID), len(req.GetData()), "DisperseBlob")
		}
		s.metrics.HandleInvalidArgRpcRequest("DisperseBlob")
		return nil, api.NewErrorInvalidArg(err.Error())
	}

	reply, err := s.disperseBlob(ctx, blob, "", "DisperseBlob")
	if err != nil {
		// Note the disperseBlob already updated metrics for this error.
		s.logger.Info("failed to disperse blob", "err", err)
	} else {
		s.metrics.HandleSuccessfulRpcRequest("DisperseBlob")
	}
	return reply, err
}

func (s *DispersalServer) DispersePaidBlob(ctx context.Context, req *pb.DispersePaidBlobRequest) (*pb.DisperseBlobReply, error) {
	return nil, api.NewErrorInternal("not implemented")
}

// Note: disperseBlob will internally update metrics upon an error; the caller doesn't need
// to track the error again.
func (s *DispersalServer) disperseBlob(ctx context.Context, blob *core.Blob, authenticatedAddress string, apiMethodName string) (*pb.DisperseBlobReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("DisperseBlob", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	securityParams := blob.RequestHeader.SecurityParams
	securityParamsStrings := make([]string, len(securityParams))
	for i, sp := range securityParams {
		securityParamsStrings[i] = sp.String()
	}

	blobSize := len(blob.Data)

	origin, err := common.GetClientAddress(ctx, s.rateConfig.ClientIPHeader, 2, true)
	if err != nil {
		for _, param := range securityParams {
			s.metrics.HandleFailedRequest(codes.InvalidArgument.String(), fmt.Sprintf("%d", param.QuorumID), blobSize, apiMethodName)
		}
		s.metrics.HandleInvalidArgRpcRequest(apiMethodName)
		return nil, api.NewErrorInvalidArg(err.Error())
	}

	s.logger.Debug("received a new blob dispersal request", "authenticatedAddress", authenticatedAddress, "origin", origin, "blobSizeBytes", blobSize, "securityParams", strings.Join(securityParamsStrings, ", "))

	if s.ratelimiter != nil {
		err := s.checkRateLimitsAndAddRatesToHeader(ctx, blob, origin, authenticatedAddress, apiMethodName)
		if err != nil {
			// Note checkRateLimitsAndAddRatesToHeader already updated the metrics for this error.
			return nil, err
		}
	}

	requestedAt := uint64(time.Now().UnixNano())
	metadataKey, err := s.blobStore.StoreBlob(ctx, blob, requestedAt)
	if err != nil {
		for _, param := range securityParams {
			s.metrics.HandleBlobStoreFailedRequest(fmt.Sprintf("%d", param.QuorumID), blobSize, apiMethodName)
		}
		s.metrics.HandleStoreFailureRpcRequest(apiMethodName)
		s.logger.Error("failed to store blob", "err", err)
		return nil, api.NewErrorInternal(fmt.Sprintf("store blob: %v", err))
	}

	for _, param := range securityParams {
		s.metrics.HandleSuccessfulRequest(fmt.Sprintf("%d", param.QuorumID), blobSize, apiMethodName)
	}

	return &pb.DisperseBlobReply{
		Result:    pb.BlobStatus_PROCESSING,
		RequestId: []byte(metadataKey.String()),
	}, nil
}

func (s *DispersalServer) getAccountRate(origin, authenticatedAddress string, quorumID core.QuorumID) (*PerUserRateInfo, string, error) {
	unauthRates, ok := s.rateConfig.QuorumRateInfos[quorumID]
	if !ok {
		return nil, "", fmt.Errorf("no configured rate exists for quorum %d", quorumID)
	}

	rates := &PerUserRateInfo{
		Name:       "",
		Throughput: unauthRates.PerUserUnauthThroughput,
		BlobRate:   unauthRates.PerUserUnauthBlobRate,
	}

	// Check if the address is in the allowlist
	if len(authenticatedAddress) > 0 {
		quorumRates, ok := s.rateConfig.Allowlist[authenticatedAddress]
		if ok {
			rateInfo, ok := quorumRates[quorumID]
			if ok {
				key := "address:" + authenticatedAddress
				if rateInfo.Throughput > 0 {
					rates.Throughput = rateInfo.Throughput
				}
				if rateInfo.BlobRate > 0 {
					rates.BlobRate = rateInfo.BlobRate
				}
				rates.Name = rateInfo.Name
				return rates, key, nil
			}
		}
	}

	// Check if the origin is in the allowlist

	// If the origin is not in the allowlist, we use the origin as the account key since
	// it is a more limited resource than an ETH public key
	key := "ip:" + origin

	for account, rateInfoByQuorum := range s.rateConfig.Allowlist {
		if !strings.Contains(origin, account) {
			continue
		}

		rateInfo, ok := rateInfoByQuorum[quorumID]
		if !ok {
			break
		}

		if rateInfo.Throughput > 0 {
			rates.Throughput = rateInfo.Throughput
		}

		if rateInfo.BlobRate > 0 {
			rates.BlobRate = rateInfo.BlobRate
		}

		if len(rateInfo.Name) > 0 {
			rates.Name = rateInfo.Name
		}

		break
	}

	return rates, key, nil
}

// Enum of rateTypes for the limiterInfo struct
type RateType uint8

const (
	SystemThroughputType RateType = iota
	SystemBlobRateType
	AccountThroughputType
	AccountBlobRateType
	RetrievalThroughputType
	RetrievalBlobRateType
)

func (r RateType) String() string {
	switch r {
	case SystemThroughputType:
		return "System throughput rate limit"
	case SystemBlobRateType:
		return "System blob rate limit"
	case AccountThroughputType:
		return "Account throughput rate limit"
	case AccountBlobRateType:
		return "Account blob rate limit"
	case RetrievalThroughputType:
		return "Retrieval throughput rate limit"
	case RetrievalBlobRateType:
		return "Retrieval blob rate limit"
	default:
		return "Unknown rate type"
	}
}

func (r RateType) Plug() string {
	switch r {
	case SystemThroughputType:
		return "system_throughput"
	case SystemBlobRateType:
		return "system_blob_rate"
	case AccountThroughputType:
		return "account_throughput"
	case AccountBlobRateType:
		return "account_blob_rate"
	case RetrievalThroughputType:
		return "retrieval_throughput"
	case RetrievalBlobRateType:
		return "retrieval_blob_rate"
	default:
		return "unknown_rate_type"
	}
}

type limiterInfo struct {
	RateType RateType
	QuorumID core.QuorumID
}

// checkRateLimitsAndAddRatesToHeader checks the configured rate limits for all of the quorums in the blob's security params,
// including both system and account level rates, relative to both the blob rate and the data bandwidth rate.
// The function will check for whitelist entries for both the authenticated address (if authenticated) and the origin.
// If no whitelist entry is found for either the origin or the authenticated address, the origin will be used as the account key
// and unauthenticated rates will be used. If the rate limit is exceeded, the function will return a ResourceExhaustedError.
// checkRateLimitsAndAddRatesToHeader will also update the blob's security params with the throughput rate for each qourum.
//
// This information is currently passed to the DA nodes for their use is ratelimiting retrieval requests. This retrieval ratelimiting
// is a temporary measure until the DA nodes are able to determine rates by themselves and will be simplified or replaced in the future.
func (s *DispersalServer) checkRateLimitsAndAddRatesToHeader(ctx context.Context, blob *core.Blob, origin, authenticatedAddress string, apiMethodName string) error {

	requestParams := make([]common.RequestParams, 0)

	blobSize := len(blob.Data)
	length := encoding.GetBlobLength(uint(blobSize))
	requesterName := ""
	for i, param := range blob.RequestHeader.SecurityParams {

		globalRates, ok := s.rateConfig.QuorumRateInfos[param.QuorumID]
		if !ok {
			s.metrics.HandleInternalFailureRpcRequest(apiMethodName)
			return api.NewErrorInternal(fmt.Sprintf("no configured rate exists for quorum %d", param.QuorumID))
		}

		accountRates, accountKey, err := s.getAccountRate(origin, authenticatedAddress, param.QuorumID)
		if err != nil {
			s.metrics.HandleInternalFailureRpcRequest(apiMethodName)
			return api.NewErrorInternal(err.Error())
		}

		// Note: There's an implicit assumption that an empty name means the account
		// is not in the allow list.
		requesterName = accountRates.Name

		// Update the quorum rate
		blob.RequestHeader.SecurityParams[i].QuorumRate = accountRates.Throughput

		// Update AccountID to accountKey.
		// This will use the origin as the account key if the user does not provide
		// an authenticated address.
		blob.RequestHeader.BlobAuthHeader.AccountID = accountKey

		// Get the encoded blob size from the blob header. Calculation is done in a way that nodes can replicate
		encodedLength := encoding.GetEncodedBlobLength(length, uint8(param.ConfirmationThreshold), uint8(param.AdversaryThreshold))
		encodedSize := encoding.GetBlobSize(encodedLength)

		// System Level
		key := fmt.Sprintf("%s:%d-%s", systemAccountKey, param.QuorumID, SystemThroughputType.Plug())
		requestParams = append(requestParams, common.RequestParams{
			RequesterID:   key,
			RequesterName: systemAccountKey,
			BlobSize:      encodedSize,
			Rate:          globalRates.TotalUnauthThroughput,
			Info: limiterInfo{
				RateType: SystemThroughputType,
				QuorumID: param.QuorumID,
			},
		})

		key = fmt.Sprintf("%s:%d-%s", systemAccountKey, param.QuorumID, SystemBlobRateType.Plug())
		requestParams = append(requestParams, common.RequestParams{
			RequesterID:   key,
			RequesterName: systemAccountKey,
			BlobSize:      blobRateMultiplier,
			Rate:          globalRates.TotalUnauthBlobRate,
			Info: limiterInfo{
				RateType: SystemBlobRateType,
				QuorumID: param.QuorumID,
			},
		})

		// Account Level
		key = fmt.Sprintf("%s:%d-%s", accountKey, param.QuorumID, AccountThroughputType.Plug())
		requestParams = append(requestParams, common.RequestParams{
			RequesterID:   key,
			RequesterName: requesterName,
			BlobSize:      encodedSize,
			Rate:          accountRates.Throughput,
			Info: limiterInfo{
				RateType: AccountThroughputType,
				QuorumID: param.QuorumID,
			},
		})

		key = fmt.Sprintf("%s:%d-%s", accountKey, param.QuorumID, AccountBlobRateType.Plug())
		requestParams = append(requestParams, common.RequestParams{
			RequesterID:   key,
			RequesterName: requesterName,
			BlobSize:      blobRateMultiplier,
			Rate:          accountRates.BlobRate,
			Info: limiterInfo{
				RateType: AccountBlobRateType,
				QuorumID: param.QuorumID,
			},
		})

	}

	s.mu.Lock()
	defer s.mu.Unlock()

	allowed, params, err := s.ratelimiter.AllowRequest(ctx, requestParams)
	if err != nil {
		s.metrics.HandleInternalFailureRpcRequest(apiMethodName)
		s.metrics.HandleFailedRequest(codes.Internal.String(), "", blobSize, apiMethodName)
		return api.NewErrorInternal(err.Error())
	}

	if !allowed {
		info, ok := params.Info.(limiterInfo)
		if !ok {
			s.metrics.HandleInternalFailureRpcRequest(apiMethodName)
			return api.NewErrorInternal("failed to cast limiterInfo")
		}
		if info.RateType == SystemThroughputType || info.RateType == SystemBlobRateType {
			s.metrics.HandleSystemRateLimitedRpcRequest(apiMethodName)
			s.metrics.HandleSystemRateLimitedRequest(fmt.Sprint(info.QuorumID), blobSize, apiMethodName)
		} else if info.RateType == AccountThroughputType || info.RateType == AccountBlobRateType {
			s.metrics.HandleAccountRateLimitedRpcRequest(apiMethodName)
			s.metrics.HandleAccountRateLimitedRequest(fmt.Sprint(info.QuorumID), blobSize, apiMethodName)
			s.logger.Info("request ratelimited", "requesterName", requesterName, "requesterID", params.RequesterID, "rateType", info.RateType.String(), "quorum", info.QuorumID)
		}
		errorString := fmt.Sprintf("request ratelimited: %s for quorum %d", info.RateType.String(), info.QuorumID)
		return api.NewErrorResourceExhausted(errorString)
	}

	return nil

}

func (s *DispersalServer) GetBlobStatus(ctx context.Context, req *pb.BlobStatusRequest) (*pb.BlobStatusReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("GetBlobStatus", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	requestID := req.GetRequestId()
	if len(requestID) == 0 {
		s.metrics.HandleInvalidArgRpcRequest("GetBlobStatus")
		s.metrics.HandleInvalidArgRequest("GetBlobStatus")
		return nil, api.NewErrorInvalidArg("request_id must not be empty")
	}

	s.logger.Info("received a new blob status request", "requestID", string(requestID))
	metadataKey, err := disperser.ParseBlobKey(string(requestID))
	if err != nil {
		s.metrics.HandleInvalidArgRpcRequest("GetBlobStatus")
		s.metrics.HandleInvalidArgRequest("GetBlobStatus")
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to parse the requestID: %s", err.Error()))
	}

	s.logger.Debug("metadataKey", "metadataKey", metadataKey.String())
	metadata, err := s.blobStore.GetBlobMetadata(ctx, metadataKey)
	if err != nil {
		if errors.Is(err, disperser.ErrMetadataNotFound) {
			s.metrics.HandleNotFoundRpcRequest("GetBlobStatus")
			s.metrics.HandleNotFoundRequest("GetBlobStatus")
			return nil, api.NewErrorNotFound("no metadata found for the requestID")
		}
		s.metrics.HandleInternalFailureRpcRequest("GetBlobStatus")
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get blob metadata, blobkey: %s", metadataKey.String()))
	}

	isConfirmed, err := metadata.IsConfirmed()
	if err != nil {
		s.metrics.HandleInternalFailureRpcRequest("GetBlobStatus")
		return nil, api.NewErrorInternal(fmt.Sprintf("missing confirmation information: %s", err.Error()))
	}

	s.metrics.HandleSuccessfulRpcRequest("GetBlobStatus")

	s.logger.Debug("isConfirmed", "metadataKey", metadataKey, "isConfirmed", isConfirmed)
	if isConfirmed {
		confirmationInfo := metadata.ConfirmationInfo
		dataLength := uint32(confirmationInfo.BlobCommitment.Length)
		quorumResults := confirmationInfo.QuorumResults
		batchQuorumIDs := make([]uint8, 0, len(quorumResults))
		for quorumID := range quorumResults {
			batchQuorumIDs = append(batchQuorumIDs, quorumID)
		}
		slices.Sort(batchQuorumIDs)
		quorumNumbers := make([]byte, len(batchQuorumIDs))
		quorumPercentSigned := make([]byte, len(batchQuorumIDs))
		for i, quorumID := range batchQuorumIDs {
			quorumNumbers[i] = quorumID
			quorumPercentSigned[i] = confirmationInfo.QuorumResults[quorumID].PercentSigned
		}

		quorumInfos := confirmationInfo.BlobQuorumInfos
		slices.SortStableFunc[[]*core.BlobQuorumInfo](quorumInfos, func(a, b *core.BlobQuorumInfo) int {
			return int(a.QuorumID) - int(b.QuorumID)
		})
		blobQuorumParams := make([]*pb.BlobQuorumParam, len(quorumInfos))
		quorumIndexes := make([]byte, len(quorumInfos))
		for i, quorumInfo := range quorumInfos {
			blobQuorumParams[i] = &pb.BlobQuorumParam{
				QuorumNumber:                    uint32(quorumInfo.QuorumID),
				AdversaryThresholdPercentage:    uint32(quorumInfo.AdversaryThreshold),
				ConfirmationThresholdPercentage: uint32(quorumInfo.ConfirmationThreshold),
				ChunkLength:                     uint32(quorumInfo.ChunkLength),
			}
			quorumIndexes[i] = byte(slices.Index(quorumNumbers, quorumInfo.QuorumID))
		}

		return &pb.BlobStatusReply{
			Status: getResponseStatus(metadata.BlobStatus),
			Info: &pb.BlobInfo{
				BlobHeader: &pb.BlobHeader{
					Commitment: &commonpb.G1Commitment{
						X: confirmationInfo.BlobCommitment.Commitment.X.Marshal(),
						Y: confirmationInfo.BlobCommitment.Commitment.Y.Marshal(),
					},
					DataLength:       dataLength,
					BlobQuorumParams: blobQuorumParams,
				},
				BlobVerificationProof: &pb.BlobVerificationProof{
					BatchId:   confirmationInfo.BatchID,
					BlobIndex: confirmationInfo.BlobIndex,
					BatchMetadata: &pb.BatchMetadata{
						BatchHeader: &pb.BatchHeader{
							BatchRoot:               confirmationInfo.BatchRoot,
							QuorumNumbers:           quorumNumbers,
							QuorumSignedPercentages: quorumPercentSigned,
							ReferenceBlockNumber:    confirmationInfo.ReferenceBlockNumber,
						},
						SignatoryRecordHash:     confirmationInfo.SignatoryRecordHash[:],
						Fee:                     confirmationInfo.Fee,
						ConfirmationBlockNumber: confirmationInfo.ConfirmationBlockNumber,
						BatchHeaderHash:         confirmationInfo.BatchHeaderHash[:],
					},
					InclusionProof: confirmationInfo.BlobInclusionProof,
					// ref: api/proto/disperser/disperser.proto:BlobVerificationProof.quorum_indexes
					QuorumIndexes: quorumIndexes,
				},
			},
		}, nil
	}

	return &pb.BlobStatusReply{
		Status: getResponseStatus(metadata.BlobStatus),
		Info:   &pb.BlobInfo{},
	}, nil
}

func (s *DispersalServer) RetrieveBlob(ctx context.Context, req *pb.RetrieveBlobRequest) (*pb.RetrieveBlobReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("RetrieveBlob", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	origin, err := common.GetClientAddress(ctx, s.rateConfig.ClientIPHeader, 2, true)
	if err != nil {
		s.metrics.HandleInvalidArgRpcRequest("RetrieveBlob")
		s.metrics.HandleInvalidArgRequest("RetrieveBlob")
		return nil, api.NewErrorInvalidArg(err.Error())
	}

	stageTimer := time.Now()
	// Check blob rate limit
	if s.ratelimiter != nil {
		allowed, param, err := s.ratelimiter.AllowRequest(ctx, []common.RequestParams{
			{
				RequesterID: fmt.Sprintf("%s:%s", origin, RetrievalBlobRateType.Plug()),
				BlobSize:    blobRateMultiplier,
				Rate:        s.rateConfig.RetrievalBlobRate,
				Info:        RetrievalBlobRateType.String(),
			},
		})
		if err != nil {
			s.metrics.HandleInternalFailureRpcRequest("RetrieveBlob")
			return nil, api.NewErrorInternal(fmt.Sprintf("ratelimiter error: %v", err))
		}
		if !allowed {
			s.metrics.HandleRateLimitedRpcRequest("RetrieveBlob")
			s.metrics.HandleFailedRequest(codes.ResourceExhausted.String(), "", 0, "RetrieveBlob")
			errorString := "request ratelimited"
			info, ok := param.Info.(string)
			if ok {
				errorString += ": " + info
			}
			return nil, api.NewErrorResourceExhausted(errorString)
		}
	}
	s.logger.Debug("checked retrieval blob rate limiting", "requesterID", fmt.Sprintf("%s:%s", origin, RetrievalBlobRateType.Plug()), "duration", time.Since(stageTimer).String())
	s.logger.Info("received a new blob retrieval request", "batchHeaderHash", req.BatchHeaderHash, "blobIndex", req.BlobIndex)

	batchHeaderHash := req.GetBatchHeaderHash()
	// Convert to [32]byte
	var batchHeaderHash32 [32]byte
	copy(batchHeaderHash32[:], batchHeaderHash)

	blobIndex := req.GetBlobIndex()

	stageTimer = time.Now()
	blobMetadata, err := s.blobStore.GetMetadataInBatch(ctx, batchHeaderHash32, blobIndex)
	if err != nil {
		s.logger.Error("Failed to retrieve blob metadata", "err", err)
		if errors.Is(err, disperser.ErrMetadataNotFound) {
			s.metrics.HandleNotFoundRpcRequest("RetrieveBlob")
			s.metrics.HandleNotFoundRequest("RetrieveBlob")
			return nil, api.NewErrorNotFound("no metadata found for the given batch header hash and blob index")
		}
		s.metrics.HandleInternalFailureRpcRequest("RetrieveBlob")
		s.metrics.IncrementFailedBlobRequestNum(codes.Internal.String(), "", "RetrieveBlob")
		return nil, api.NewErrorInternal("failed to get blob metadata, please retry")
	}

	if blobMetadata.Expiry < uint64(time.Now().Unix()) {
		s.metrics.HandleNotFoundRpcRequest("RetrieveBlob")
		s.metrics.HandleNotFoundRequest("RetrieveBlob")
		return nil, api.NewErrorNotFound("no metadata found for the given batch header hash and blob index")
	}

	s.logger.Debug("fetched blob metadata", "batchHeaderHash", req.BatchHeaderHash, "blobIndex", req.BlobIndex, "duration", time.Since(stageTimer).String())

	stageTimer = time.Now()
	// Check throughout rate limit
	blobSize := encoding.GetBlobSize(blobMetadata.ConfirmationInfo.BlobCommitment.Length)

	if s.ratelimiter != nil {
		allowed, param, err := s.ratelimiter.AllowRequest(ctx, []common.RequestParams{
			{
				RequesterID: fmt.Sprintf("%s:%s", origin, RetrievalThroughputType.Plug()),
				BlobSize:    blobSize,
				Rate:        s.rateConfig.RetrievalThroughput,
				Info:        RetrievalThroughputType.String(),
			},
		})
		if err != nil {
			s.metrics.HandleInternalFailureRpcRequest("RetrieveBlob")
			return nil, api.NewErrorInternal(fmt.Sprintf("ratelimiter error: %v", err))
		}
		if !allowed {
			s.metrics.HandleRateLimitedRpcRequest("RetrieveBlob")
			s.metrics.HandleFailedRequest(codes.ResourceExhausted.String(), "", 0, "RetrieveBlob")
			errorString := "request ratelimited"
			info, ok := param.Info.(string)
			if ok {
				errorString += ": " + info
			}
			return nil, api.NewErrorResourceExhausted(errorString)
		}
	}
	s.logger.Debug("checked retrieval throughput rate limiting", "requesterID", fmt.Sprintf("%s:%s", origin, RetrievalThroughputType.Plug()), "duration (ms)", time.Since(stageTimer).String())

	stageTimer = time.Now()
	data, err := s.blobStore.GetBlobContent(ctx, blobMetadata.BlobHash)
	if err != nil {
		s.logger.Error("Failed to retrieve blob", "err", err)
		s.metrics.HandleInternalFailureRpcRequest("RetrieveBlob")
		s.metrics.HandleFailedRequest(codes.Internal.String(), "", len(data), "RetrieveBlob")
		return nil, api.NewErrorInternal("failed to get blob data, please retry")
	}
	s.metrics.HandleSuccessfulRpcRequest("RetrieveBlob")
	s.metrics.HandleSuccessfulRequest("", len(data), "RetrieveBlob")

	s.logger.Debug("fetched blob content", "batchHeaderHash", req.BatchHeaderHash, "blobIndex", req.BlobIndex, "data size (bytes)", len(data), "duration", time.Since(stageTimer).String())

	return &pb.RetrieveBlobReply{
		Data: data,
	}, nil
}

func (s *DispersalServer) GetRateConfig() *RateConfig {
	return &s.rateConfig
}

func (s *DispersalServer) Start(ctx context.Context) error {
	go func() {
		t := time.NewTicker(s.rateConfig.AllowlistRefreshInterval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				s.LoadAllowlist()
			}
		}
	}()
	// Serve grpc requests
	addr := fmt.Sprintf("%s:%s", disperser.Localhost, s.serverConfig.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.New("could not start tcp listener")
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB

	gs := grpc.NewServer(opt)
	reflection.Register(gs)
	pb.RegisterDisperserServer(gs, s)

	// Register Server for Health Checks
	name := pb.Disperser_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, gs)

	s.logger.Info("GRPC Listening", "port", s.serverConfig.GrpcPort, "address", listener.Addr().String(), "maxBlobSize", s.maxBlobSize)

	if err := gs.Serve(listener); err != nil {
		return errors.New("could not start GRPC server")
	}

	return nil
}

func (s *DispersalServer) LoadAllowlist() {
	al, err := ReadAllowlistFromFile(s.rateConfig.AllowlistFile)
	if err != nil {
		s.logger.Error("failed to load allowlist", "err", err)
		return
	}
	s.rateConfig.Allowlist = al
	for account, rateInfoByQuorum := range al {
		for quorumID, rateInfo := range rateInfoByQuorum {
			s.logger.Info("[Allowlist]", "account", account, "name", rateInfo.Name, "quorumID", quorumID, "throughput", rateInfo.Throughput, "blobRate", rateInfo.BlobRate)
		}
	}
}

// updateQuorumConfig updates the quorum config and returns the updated quorum config. If the update fails,
// it will fallback to the old quorumConfig if it is set. This is to improve the robustness of the disperser to
// RPC failures since the quorum config is rarely updated. In the event that quorumConfig is incorrect, this will
// not result in a safety failure since all parameters are separately validated on the smart contract.

func (s *DispersalServer) updateQuorumConfig(ctx context.Context) (QuorumConfig, error) {

	s.mu.RLock()
	newConfig := s.quorumConfig
	s.mu.RUnlock()

	// If the quorum count is set, we will fallback to the old quorumConfig if the RPC fails
	fallbackToCache := newConfig.QuorumCount != 0

	currentBlock, err := s.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		s.logger.Error("failed to get current block number", "err", err)
		if !fallbackToCache {
			return QuorumConfig{}, err
		}
		return newConfig, nil
	}

	count, err := s.tx.GetQuorumCount(ctx, currentBlock)
	if err != nil {
		s.logger.Error("failed to get quorum count", "err", err)
		if !fallbackToCache {
			return QuorumConfig{}, err
		}
	} else {
		newConfig.QuorumCount = count
	}

	securityParams, err := s.tx.GetQuorumSecurityParams(ctx, currentBlock)
	if err != nil {
		s.logger.Error("failed to get quorum security params", "err", err)
		if !fallbackToCache {
			return QuorumConfig{}, err
		}
	} else {
		newConfig.SecurityParams = make(map[core.QuorumID]core.SecurityParam)
		for _, sp := range securityParams {
			newConfig.SecurityParams[sp.QuorumID] = sp
		}
	}

	requiredQuorums, err := s.tx.GetRequiredQuorumNumbers(ctx, currentBlock)
	if err != nil {
		s.logger.Error("failed to get quorum security params", "err", err)
		if !fallbackToCache {
			return QuorumConfig{}, err
		}
	} else {
		newConfig.RequiredQuorums = requiredQuorums
	}

	s.logger.Debug("updating quorum count", "currentBlock", currentBlock, "count", count)
	s.mu.Lock()
	s.quorumConfig = newConfig
	s.mu.Unlock()
	return newConfig, nil
}

func getResponseStatus(status disperser.BlobStatus) pb.BlobStatus {
	switch status {
	case disperser.Dispersing, disperser.Processing:
		return pb.BlobStatus_PROCESSING
	case disperser.Confirmed:
		return pb.BlobStatus_CONFIRMED
	case disperser.Failed:
		return pb.BlobStatus_FAILED
	case disperser.Finalized:
		return pb.BlobStatus_FINALIZED
	case disperser.InsufficientSignatures:
		return pb.BlobStatus_INSUFFICIENT_SIGNATURES
	default:
		return pb.BlobStatus_UNKNOWN
	}
}

func (s *DispersalServer) validateRequestAndGetBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*core.Blob, error) {

	data := req.GetData()
	blobSize := len(data)
	// The blob size in bytes must be in range [1, maxBlobSize].
	if blobSize > s.maxBlobSize {
		return nil, fmt.Errorf("blob size cannot exceed %v Bytes", s.maxBlobSize)
	}
	if blobSize == 0 {
		return nil, fmt.Errorf("blob size must be greater than 0")
	}

	if len(req.GetCustomQuorumNumbers()) > 256 {
		return nil, errors.New("number of custom_quorum_numbers must not exceed 256")
	}

	// validate every 32 bytes is a valid field element
	_, err := rs.ToFrArray(data)
	if err != nil {
		s.logger.Error("failed to convert a 32bytes as a field element", "err", err)
		return nil, api.NewErrorInvalidArg("encountered an error to convert a 32-bytes into a valid field element, please use the correct format where every 32bytes(big-endian) is less than 21888242871839275222246405745257275088548364400416034343698204186575808495617")
	}

	quorumConfig, err := s.updateQuorumConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get quorum config: %w", err)
	}

	if len(req.GetCustomQuorumNumbers()) > int(quorumConfig.QuorumCount) {
		return nil, errors.New("number of custom_quorum_numbers must not exceed number of quorums")
	}

	seenQuorums := make(map[uint8]struct{})
	// The quorum ID must be in range [0, 254]. It'll actually be converted
	// to uint8, so it cannot be greater than 254.
	for i := range req.GetCustomQuorumNumbers() {

		if req.GetCustomQuorumNumbers()[i] > core.MaxQuorumID {
			return nil, fmt.Errorf("custom_quorum_numbers must be in range [0, 254], but found %d", req.GetCustomQuorumNumbers()[i])
		}

		quorumID := uint8(req.GetCustomQuorumNumbers()[i])
		if quorumID >= quorumConfig.QuorumCount {
			return nil, fmt.Errorf("custom_quorum_numbers must be in range [0, %d], but found %d", s.quorumConfig.QuorumCount-1, quorumID)
		}

		if _, ok := seenQuorums[quorumID]; ok {
			return nil, fmt.Errorf("custom_quorum_numbers must not contain duplicates")
		}
		seenQuorums[quorumID] = struct{}{}

	}

	// Add the required quorums to the list of quorums to check
	for _, quorumID := range quorumConfig.RequiredQuorums {
		// Note: no matter if dual quorum staking is enabled or not, custom_quorum_numbers cannot
		// use any required quorums that are defined onchain.
		if _, ok := seenQuorums[quorumID]; ok {
			return nil, fmt.Errorf("custom_quorum_numbers should not include the required quorums %v, but required quorum %d was found", quorumConfig.RequiredQuorums, quorumID)
		}

		seenQuorums[quorumID] = struct{}{}
	}

	if len(seenQuorums) == 0 {
		return nil, fmt.Errorf("the blob must be sent to at least one quorum")
	}

	params := make([]*core.SecurityParam, len(seenQuorums))
	i := 0
	for quorumID := range seenQuorums {
		params[i] = &core.SecurityParam{
			QuorumID:              core.QuorumID(quorumID),
			AdversaryThreshold:    quorumConfig.SecurityParams[quorumID].AdversaryThreshold,
			ConfirmationThreshold: quorumConfig.SecurityParams[quorumID].ConfirmationThreshold,
		}
		err = params[i].Validate()
		if err != nil {
			return nil, fmt.Errorf("invalid request: %w", err)
		}
		i++
	}

	header := core.BlobRequestHeader{
		BlobAuthHeader: core.BlobAuthHeader{
			AccountID: req.AccountId,
		},
		SecurityParams: params,
	}

	blob := &core.Blob{
		RequestHeader: header,
		Data:          data,
	}

	return blob, nil
}
