package apiserver

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"slices"
	"strings"
	"sync"
	"time"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/common"
	healthcheck "github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var errSystemBlobRateLimit = fmt.Errorf("request ratelimited: system blob limit")
var errSystemThroughputRateLimit = fmt.Errorf("request ratelimited: system throughput limit")
var errAccountBlobRateLimit = fmt.Errorf("request ratelimited: account blob limit")
var errAccountThroughputRateLimit = fmt.Errorf("request ratelimited: account throughput limit")

const systemAccountKey = "system"

const maxBlobSize = 2 * 1024 * 1024 // 2 MiB

type DispersalServer struct {
	pb.UnimplementedDisperserServer
	mu *sync.RWMutex

	serverConfig disperser.ServerConfig
	rateConfig   RateConfig

	blobStore    disperser.BlobStore
	tx           core.Transactor
	quorumConfig QuorumConfig

	ratelimiter   common.RateLimiter
	authenticator core.BlobRequestAuthenticator

	metrics *disperser.Metrics

	logger common.Logger
}

type QuorumConfig struct {
	RequiredQuorums []core.QuorumID
	QuorumCount     uint8
	SecurityParams  []*core.SecurityParam
}

// NewServer creates a new Server struct with the provided parameters.
//
// Note: The Server's chunks store will be created at config.DbPath+"/chunk".
func NewDispersalServer(
	serverConfig disperser.ServerConfig,
	store disperser.BlobStore,
	tx core.Transactor,
	logger common.Logger,
	metrics *disperser.Metrics,
	ratelimiter common.RateLimiter,
	rateConfig RateConfig,
) *DispersalServer {
	for ip, rateInfoByQuorum := range rateConfig.Allowlist {
		for quorumID, rateInfo := range rateInfoByQuorum {
			logger.Info("[Allowlist]", "ip", ip, "quorumID", quorumID, "throughput", rateInfo.Throughput, "blobRate", rateInfo.BlobRate)
		}
	}

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
	}
}

func (s *DispersalServer) DisperseBlobAuthenticated(stream pb.Disperser_DisperseBlobAuthenticatedServer) error {

	// Process disperse_request
	in, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("error receiving next message: %v", err)
	}

	request, ok := in.Payload.(*pb.AuthenticatedRequest_DisperseRequest)
	if !ok {
		return errors.New("expected DisperseBlobRequest")
	}

	blob, err := s.validateRequestAndGetBlob(stream.Context(), request.DisperseRequest)
	if err != nil {
		for _, quorumID := range request.DisperseRequest.CustomQuorumNumbers {
			s.metrics.HandleFailedRequest(fmt.Sprint(quorumID), len(request.DisperseRequest.GetData()), "DisperseBlob")
		}
		return err
	}

	// Get the ethereum address associated with the public key. This is just for convenience so we can put addresses instead of public keys in the allowlist.
	// Decode public key
	publicKeyBytes, err := hexutil.Decode(blob.RequestHeader.AccountID)
	if err != nil {
		return fmt.Errorf("failed to decode public key (%v): %v", blob.RequestHeader.AccountID, err)
	}

	pubKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to decode public key (%v): %v", blob.RequestHeader.AccountID, err)
	}

	authenticatedAddress := crypto.PubkeyToAddress(*pubKey).String()

	// Send back challenge to client
	challenge := rand.Uint32()
	err = stream.Send(&pb.AuthenticatedReply{Payload: &pb.AuthenticatedReply_BlobAuthHeader{
		BlobAuthHeader: &pb.BlobAuthHeader{
			ChallengeParameter: challenge,
		},
	}})
	if err != nil {
		return err
	}

	// Recieve challenge_reply
	in, err = stream.Recv()
	if err != nil {
		return fmt.Errorf("error receiving next message: %v", err)
	}

	challengeReply, ok := in.Payload.(*pb.AuthenticatedRequest_AuthenticationData)
	if !ok {
		return errors.New("expected AuthenticationData")
	}

	blob.RequestHeader.Nonce = challenge
	blob.RequestHeader.AuthenticationData = challengeReply.AuthenticationData.AuthenticationData

	err = s.authenticator.AuthenticateBlobRequest(blob.RequestHeader.BlobAuthHeader)
	if err != nil {
		return fmt.Errorf("failed to authenticate blob request: %v", err)
	}

	// Disperse the blob
	reply, err := s.disperseBlob(stream.Context(), blob, authenticatedAddress)
	if err != nil {
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

	return nil

}

func (s *DispersalServer) DisperseBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*pb.DisperseBlobReply, error) {

	blob, err := s.validateRequestAndGetBlob(ctx, req)
	if err != nil {
		for _, quorumID := range req.CustomQuorumNumbers {
			s.metrics.HandleFailedRequest(fmt.Sprint(quorumID), len(req.GetData()), "DisperseBlob")
		}
		return nil, err
	}

	reply, err := s.disperseBlob(ctx, blob, "")
	if err != nil {
		s.logger.Info("failed to disperse blob", "err", err)
	}
	return reply, err
}

func (s *DispersalServer) disperseBlob(ctx context.Context, blob *core.Blob, authenticatedAddress string) (*pb.DisperseBlobReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("DisperseBlob", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	securityParams := blob.RequestHeader.SecurityParams

	// securityParams := blob.RequestHeader.SecurityParams
	securityParamsStrings := make([]string, len(securityParams))
	for i, sp := range securityParams {
		securityParamsStrings[i] = sp.String()
	}

	blobSize := len(blob.Data)

	origin, err := common.GetClientAddress(ctx, s.rateConfig.ClientIPHeader, 2, true)
	if err != nil {
		for _, param := range securityParams {
			quorumId := string(param.QuorumID)
			s.metrics.HandleFailedRequest(quorumId, blobSize, "DisperseBlob")
		}
		return nil, err
	}

	s.logger.Debug("received a new blob request", "origin", origin, "securityParams", strings.Join(securityParamsStrings, ", "))

	if s.ratelimiter != nil {
		err := s.checkRateLimitsAndAddRates(ctx, blob, origin, authenticatedAddress)
		if err != nil {
			for _, param := range securityParams {
				quorumId := string(param.QuorumID)
				if errors.Is(err, errSystemBlobRateLimit) || errors.Is(err, errSystemThroughputRateLimit) {
					s.metrics.HandleSystemRateLimitedRequest(quorumId, blobSize, "DisperseBlob")
				} else if errors.Is(err, errAccountBlobRateLimit) || errors.Is(err, errAccountThroughputRateLimit) {
					s.metrics.HandleAccountRateLimitedRequest(quorumId, blobSize, "DisperseBlob")
				} else {
					s.metrics.HandleFailedRequest(quorumId, blobSize, "DisperseBlob")
				}
			}
			return nil, err
		}
	}

	requestedAt := uint64(time.Now().UnixNano())
	metadataKey, err := s.blobStore.StoreBlob(ctx, blob, requestedAt)
	if err != nil {
		for _, param := range securityParams {
			quorumId := string(param.QuorumID)
			s.metrics.HandleBlobStoreFailedRequest(quorumId, blobSize, "DisperseBlob")
		}
		s.logger.Error("failed to store blob", "err", err)
		return nil, fmt.Errorf("failed to store blob, please try again later")
	}

	for _, param := range securityParams {
		quorumId := string(param.QuorumID)
		s.metrics.HandleSuccessfulRequest(quorumId, blobSize, "DisperseBlob")
	}

	s.logger.Info("successfully received a new blob: ", "key", metadataKey.String())
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
				return rates, key, nil
			}
		}
	}

	// Check if the origin is in the allowlist

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

		break
	}

	return rates, key, nil

}

func (s *DispersalServer) checkRateLimitsAndAddRates(ctx context.Context, blob *core.Blob, origin, authenticatedAddress string) error {

	// TODO(robert): Remove these locks once we have resolved ratelimiting approach
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, param := range blob.RequestHeader.SecurityParams {

		rates, ok := s.rateConfig.QuorumRateInfos[param.QuorumID]
		if !ok {
			return fmt.Errorf("no configured rate exists for quorum %d", param.QuorumID)
		}
		accountRates, accountKey, err := s.getAccountRate(origin, authenticatedAddress, param.QuorumID)
		if err != nil {
			return err
		}

		// Get the encoded blob size from the blob header. Calculation is done in a way that nodes can replicate
		blobSize := len(blob.Data)
		length := encoding.GetBlobLength(uint(blobSize))
		encodedLength := encoding.GetEncodedBlobLength(length, uint8(param.ConfirmationThreshold), uint8(param.AdversaryThreshold))
		encodedSize := encoding.GetBlobSize(encodedLength)

		s.logger.Debug("checking rate limits", "origin", origin, "address", authenticatedAddress, "quorum", param.QuorumID, "encodedSize", encodedSize, "blobSize", blobSize,
			"accountThroughput", accountRates.Throughput, "accountBlobRate", accountRates.BlobRate, "accountKey", accountKey)

		// Check System Ratelimit
		systemQuorumKey := fmt.Sprintf("%s:%d", systemAccountKey, param.QuorumID)
		allowed, err := s.ratelimiter.AllowRequest(ctx, systemQuorumKey, encodedSize, rates.TotalUnauthThroughput)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("system byte ratelimit exceeded", "systemQuorumKey", systemQuorumKey, "rate", rates.TotalUnauthThroughput)
			return errSystemThroughputRateLimit
		}

		systemQuorumKey = fmt.Sprintf("%s:%d-blobrate", systemAccountKey, param.QuorumID)
		allowed, err = s.ratelimiter.AllowRequest(ctx, systemQuorumKey, blobRateMultiplier, rates.TotalUnauthBlobRate)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("system blob ratelimit exceeded", "systemQuorumKey", systemQuorumKey, "rate", float32(rates.TotalUnauthBlobRate)/blobRateMultiplier)
			return errSystemBlobRateLimit
		}

		// Check Account Ratelimit

		accountQuorumKey := fmt.Sprintf("%s:%d", accountKey, param.QuorumID)
		allowed, err = s.ratelimiter.AllowRequest(ctx, accountQuorumKey, encodedSize, accountRates.Throughput)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("account byte ratelimit exceeded", "accountQuorumKey", accountQuorumKey, "rate", accountRates.Throughput)
			return errAccountThroughputRateLimit
		}

		accountQuorumKey = fmt.Sprintf("%s:%d-blobrate", accountKey, param.QuorumID)
		allowed, err = s.ratelimiter.AllowRequest(ctx, accountQuorumKey, blobRateMultiplier, accountRates.BlobRate)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("account blob ratelimit exceeded", "accountQuorumKey", accountQuorumKey, "rate", float32(accountRates.BlobRate)/blobRateMultiplier)
			return errAccountBlobRateLimit
		}

		// Update the quorum rate
		blob.RequestHeader.SecurityParams[i].QuorumRate = accountRates.Throughput
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
		return nil, fmt.Errorf("invalid request: request_id must not be empty")
	}

	s.logger.Info("received a new blob status request", "requestID", string(requestID))
	metadataKey, err := disperser.ParseBlobKey(string(requestID))
	if err != nil {
		return nil, err
	}

	s.logger.Debug("metadataKey", "metadataKey", metadataKey.String())
	metadata, err := s.blobStore.GetBlobMetadata(ctx, metadataKey)
	if err != nil {
		return nil, err
	}

	isConfirmed, err := metadata.IsConfirmed()
	if err != nil {
		return nil, err
	}

	s.logger.Debug("isConfirmed", "metadata", metadata, "isConfirmed", isConfirmed)
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

	s.logger.Info("received a new blob retrieval request", "batchHeaderHash", req.BatchHeaderHash, "blobIndex", req.BlobIndex)

	batchHeaderHash := req.GetBatchHeaderHash()
	// Convert to [32]byte
	var batchHeaderHash32 [32]byte
	copy(batchHeaderHash32[:], batchHeaderHash)

	blobIndex := req.GetBlobIndex()

	blobMetadata, err := s.blobStore.GetMetadataInBatch(ctx, batchHeaderHash32, blobIndex)
	if err != nil {
		s.logger.Error("Failed to retrieve blob metadata", "err", err)
		s.metrics.IncrementFailedBlobRequestNum("", "RetrieveBlob")

		return nil, err
	}

	data, err := s.blobStore.GetBlobContent(ctx, blobMetadata.BlobHash)
	if err != nil {
		s.logger.Error("Failed to retrieve blob", "err", err)
		s.metrics.HandleFailedRequest("", len(data), "RetrieveBlob")

		return nil, err
	}

	s.metrics.HandleSuccessfulRequest("", len(data), "RetrieveBlob")

	return &pb.RetrieveBlobReply{
		Data: data,
	}, nil
}

func (s *DispersalServer) Start(ctx context.Context) error {
	s.logger.Trace("Entering Start function...")
	defer s.logger.Trace("Exiting Start function...")

	// Serve grpc requests
	addr := fmt.Sprintf("%s:%s", disperser.Localhost, s.serverConfig.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not start tcp listener")
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB
	gs := grpc.NewServer(opt)
	reflection.Register(gs)
	pb.RegisterDisperserServer(gs, s)

	// Register Server for Health Checks
	name := pb.Disperser_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, gs)

	s.logger.Info("port", s.serverConfig.GrpcPort, "address", listener.Addr().String(), "GRPC Listening")
	if err := gs.Serve(listener); err != nil {
		return fmt.Errorf("could not start GRPC server")
	}

	return nil
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
		newConfig.SecurityParams = securityParams
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
	case disperser.Processing:
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
	if blobSize > maxBlobSize {
		return nil, fmt.Errorf("blob size cannot exceed 2 MiB")
	}
	if blobSize == 0 {
		return nil, fmt.Errorf("blob size must be greater than 0")
	}

	quorumConfig, err := s.updateQuorumConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get quorum config: %w", err)
	}

	seenQuorums := make(map[uint8]struct{})
	// The quorum ID must be in range [0, 254]. It'll actually be converted
	// to uint8, so it cannot be greater than 254.
	for i := range req.GetCustomQuorumNumbers() {

		if req.GetCustomQuorumNumbers()[i] > 254 {
			return nil, fmt.Errorf("invalid request: quorum_numbers must be in range [0, 254], but found %d", req.GetCustomQuorumNumbers()[i])
		}

		quorumID := uint8(req.GetCustomQuorumNumbers()[i])
		if _, ok := seenQuorums[quorumID]; ok {
			return nil, fmt.Errorf("invalid request: quorum_numbers must not contain duplicates")
		}
		seenQuorums[quorumID] = struct{}{}

		if quorumID >= quorumConfig.QuorumCount {
			if err != nil {
				return nil, fmt.Errorf("invalid request: the quorum_numbers must be in range [0, %d], but found %d", s.quorumConfig.QuorumCount-1, quorumID)
			}
		}
	}

	// Add the required quorums to the list of quorums to check
	for _, quorumID := range quorumConfig.RequiredQuorums {
		if _, ok := seenQuorums[quorumID]; ok {
			return nil, fmt.Errorf("invalid request: quorum_numbers should not include the required quorums, but required quorum %d was found", quorumID)
		}
		seenQuorums[quorumID] = struct{}{}
	}

	if len(seenQuorums) == 0 {
		return nil, fmt.Errorf("invalid request: the blob must be sent to at least one quorum")
	}

	params := make([]*core.SecurityParam, len(seenQuorums))
	i := 0
	for quorumID := range seenQuorums {
		params[i] = &core.SecurityParam{
			QuorumID:              core.QuorumID(quorumID),
			AdversaryThreshold:    quorumConfig.SecurityParams[i].AdversaryThreshold,
			ConfirmationThreshold: quorumConfig.SecurityParams[i].ConfirmationThreshold,
		}
		i++
	}

	header := core.BlobRequestHeader{
		BlobAuthHeader: core.BlobAuthHeader{
			AccountID: req.AccountId,
		},
		SecurityParams: params,
	}

	if err := header.Validate(); err != nil {
		s.logger.Warn("invalid header", "err", err)
		return nil, err
	}

	blob := &core.Blob{
		RequestHeader: header,
		Data:          data,
	}

	return blob, nil
}
