package apiserver

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/common"
	healthcheck "github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/disperser"
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
	mu *sync.Mutex

	config disperser.ServerConfig

	blobStore   disperser.BlobStore
	tx          core.Transactor
	quorumCount uint16

	rateConfig    RateConfig
	ratelimiter   common.RateLimiter
	authenticator core.BlobRequestAuthenticator

	metrics *disperser.Metrics

	logger common.Logger
}

// NewServer creates a new Server struct with the provided parameters.
//
// Note: The Server's chunks store will be created at config.DbPath+"/chunk".
func NewDispersalServer(
	config disperser.ServerConfig,
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
		config:        config,
		blobStore:     store,
		tx:            tx,
		quorumCount:   0,
		metrics:       metrics,
		logger:        logger,
		ratelimiter:   ratelimiter,
		authenticator: authenticator,
		rateConfig:    rateConfig,
		mu:            &sync.Mutex{},
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

	blob := getBlobFromRequest(request.DisperseRequest)

	blob.RequestHeader.Nonce = challenge
	blob.RequestHeader.AuthenticationData = challengeReply.AuthenticationData.AuthenticationData

	err = s.authenticator.AuthenticateBlobRequest(blob.RequestHeader.BlobAuthHeader)
	if err != nil {
		return fmt.Errorf("failed to authenticate blob request: %v", err)
	}

	// Get the ethereum address associated with the public key. This is just for convenience so we can put addresses instead of public keys in the allowlist.
	authenticatedAddress := ""
	// Decode public key
	publicKeyBytes, err := hexutil.Decode(blob.RequestHeader.AccountID)
	if err != nil {
		s.logger.Warn("failed to decode public key (%v): %v", blob.RequestHeader.AccountID, err)
	} else {
		pubKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
		if err != nil {
			s.logger.Warn("failed to decode public key (%v): %v", blob.RequestHeader.AccountID, err)
		} else {
			// Get the address
			authenticatedAddress = crypto.PubkeyToAddress(*pubKey).String()
		}
	}

	// Disperse the blob
	reply, err := s.disperseBlob(stream.Context(), blob, authenticatedAddress)
	if err != nil {
		return err
	}

	// Send back disperse_reply
	err = stream.Send(&pb.AuthenticatedReply{Payload: &pb.AuthenticatedReply_DisperseReply{
		DisperseReply: reply,
	}})
	if err != nil {
		return err
	}

	return nil

}

func (s *DispersalServer) DisperseBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*pb.DisperseBlobReply, error) {

	blob := getBlobFromRequest(req)

	return s.disperseBlob(ctx, blob, "")

}

func (s *DispersalServer) disperseBlob(ctx context.Context, blob *core.Blob, authenticatedAddress string) (*pb.DisperseBlobReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("DisperseBlob", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	securityParams := blob.RequestHeader.SecurityParams
	if len(securityParams) == 0 {
		return nil, fmt.Errorf("invalid request: security_params must not be empty")
	}
	if len(securityParams) > 256 {
		return nil, fmt.Errorf("invalid request: security_params must not exceed 256")
	}

	seenQuorums := make(map[uint8]struct{})
	// The quorum ID must be in range [0, 255]. It'll actually be converted
	// to uint8, so it cannot be greater than 255.
	for _, param := range securityParams {
		if _, ok := seenQuorums[param.QuorumID]; ok {
			return nil, fmt.Errorf("invalid request: security_params must not contain duplicate quorum_id")
		}
		seenQuorums[param.QuorumID] = struct{}{}

		if uint16(param.QuorumID) >= s.quorumCount {
			err := s.updateQuorumCount(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get onchain quorum count: %w", err)
			}

			if uint16(param.QuorumID) >= s.quorumCount {
				return nil, fmt.Errorf("invalid request: the quorum_id must be in range [0, %d], but found %d", s.quorumCount-1, param.QuorumID)
			}
		}
	}

	blobSize := len(blob.Data)
	// The blob size in bytes must be in range [1, maxBlobSize].
	if blobSize > maxBlobSize {
		return nil, fmt.Errorf("blob size cannot exceed 2 MiB")
	}
	if blobSize == 0 {
		return nil, fmt.Errorf("blob size must be greater than 0")
	}

	origin, err := common.GetClientAddress(ctx, s.rateConfig.ClientIPHeader, 2, true)
	if err != nil {
		for _, param := range securityParams {
			quorumId := string(param.QuorumID)
			s.metrics.HandleFailedRequest(quorumId, blobSize, "DisperseBlob")
		}
		return nil, err
	}

	s.logger.Debug("received a new blob request", "origin", origin, "securityParams", securityParams)

	if err := blob.RequestHeader.Validate(); err != nil {
		s.logger.Warn("invalid header", "err", err)
		for _, param := range securityParams {
			quorumId := string(param.QuorumID)
			s.metrics.HandleFailedRequest(quorumId, blobSize, "DisperseBlob")
		}
		return nil, err
	}

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
			s.metrics.HandleFailedRequest(quorumId, blobSize, "DisperseBlob")
		}
		return nil, err
	}

	for _, param := range securityParams {
		quorumId := string(param.QuorumID)
		s.metrics.HandleSuccessfulRequest(quorumId, blobSize, "DisperseBlob")
	}

	s.logger.Info("received a new blob: ", "key", metadataKey.String())
	return &pb.DisperseBlobReply{
		Result:    pb.BlobStatus_PROCESSING,
		RequestId: []byte(metadataKey.String()),
	}, nil
}

type RateKeys struct {
	ThroughputKey string
	BlobRateKey   string
}

func (s *DispersalServer) getAccountRate(origin, authenticatedAddress string, quorumID core.QuorumID) (*PerUserRateInfo, *RateKeys, error) {
	unauthRates, ok := s.rateConfig.QuorumRateInfos[quorumID]
	if !ok {
		return nil, nil, fmt.Errorf("no configured rate exists for quorum %d", quorumID)
	}

	// Check if the origin is in the allowlist
	rates := &PerUserRateInfo{
		Throughput: unauthRates.PerUserUnauthThroughput,
		BlobRate:   unauthRates.PerUserUnauthBlobRate,
	}

	keys := &RateKeys{
		ThroughputKey: "ip:" + origin,
		BlobRateKey:   "ip:" + origin,
	}

	for account, rateInfoByQuorum := range s.rateConfig.Allowlist {
		if !strings.Contains(origin, account) {
			continue
		}

		rateInfo, ok := rateInfoByQuorum[quorumID]
		if !ok {
			break
		}

		if rateInfo.Throughput > unauthRates.PerUserUnauthThroughput {
			rates.Throughput = rateInfo.Throughput
		}

		if rateInfo.BlobRate > unauthRates.PerUserUnauthBlobRate {
			rates.BlobRate = rateInfo.BlobRate
		}

		break
	}

	// Check if the address is in the allowlist
	quorumRates, ok := s.rateConfig.Allowlist[authenticatedAddress]
	if ok {
		rateInfo, ok := quorumRates[quorumID]
		if ok {

			if rateInfo.Throughput > rates.Throughput {
				rates.Throughput = rateInfo.Throughput
				keys.ThroughputKey = "address:" + authenticatedAddress
			}

			if rateInfo.BlobRate > rates.BlobRate {
				rates.BlobRate = rateInfo.BlobRate
				keys.BlobRateKey = "address:" + authenticatedAddress
			}

		}
	}

	return rates, keys, nil

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
		accountRates, keys, err := s.getAccountRate(origin, authenticatedAddress, param.QuorumID)
		if err != nil {
			return err
		}

		// Get the encoded blob size from the blob header. Calculation is done in a way that nodes can replicate
		blobSize := len(blob.Data)
		length := core.GetBlobLength(uint(blobSize))
		encodedLength := core.GetEncodedBlobLength(length, uint8(param.QuorumThreshold), uint8(param.AdversaryThreshold))
		encodedSize := core.GetBlobSize(encodedLength)

		s.logger.Debug("checking rate limits", "origin", origin, "address", authenticatedAddress, "quorum", param.QuorumID, "encodedSize", encodedSize, "blobSize", blobSize,
			"accountThroughput", accountRates.Throughput, "accountBlobRate", accountRates.BlobRate,
			"throughputKey", keys.ThroughputKey, "blobRateKey", keys.BlobRateKey)

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

		userQuorumKey := fmt.Sprintf("%s:%d", keys.ThroughputKey, param.QuorumID)
		allowed, err = s.ratelimiter.AllowRequest(ctx, userQuorumKey, encodedSize, accountRates.Throughput)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("account byte ratelimit exceeded", "userQuorumKey", userQuorumKey, "rate", accountRates.Throughput)
			return errAccountThroughputRateLimit
		}

		userQuorumKey = fmt.Sprintf("%s:%d-blobrate", keys.BlobRateKey, param.QuorumID)
		allowed, err = s.ratelimiter.AllowRequest(ctx, userQuorumKey, blobRateMultiplier, accountRates.BlobRate)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("account blob ratelimit exceeded", "userQuorumKey", userQuorumKey, "rate", float32(accountRates.BlobRate)/blobRateMultiplier)
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
		commit, err := confirmationInfo.BlobCommitment.Commitment.Serialize()
		if err != nil {
			return nil, err
		}

		dataLength := uint32(confirmationInfo.BlobCommitment.Length)
		quorumInfos := confirmationInfo.BlobQuorumInfos
		blobQuorumParams := make([]*pb.BlobQuorumParam, len(quorumInfos))
		quorumNumbers := make([]byte, len(quorumInfos))
		quorumPercentSigned := make([]byte, len(quorumInfos))
		quorumIndexes := make([]byte, len(quorumInfos))
		for i, quorumInfo := range quorumInfos {
			blobQuorumParams[i] = &pb.BlobQuorumParam{
				QuorumNumber:                 uint32(quorumInfo.QuorumID),
				AdversaryThresholdPercentage: uint32(quorumInfo.AdversaryThreshold),
				QuorumThresholdPercentage:    uint32(quorumInfo.QuorumThreshold),
				ChunkLength:                  uint32(quorumInfo.ChunkLength),
			}
			quorumNumbers[i] = quorumInfo.QuorumID
			quorumPercentSigned[i] = confirmationInfo.QuorumResults[quorumInfo.QuorumID].PercentSigned
			quorumIndexes[i] = byte(i)
		}

		return &pb.BlobStatusReply{
			Status: getResponseStatus(metadata.BlobStatus),
			Info: &pb.BlobInfo{
				BlobHeader: &pb.BlobHeader{
					Commitment:       commit,
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
	addr := fmt.Sprintf("%s:%s", disperser.Localhost, s.config.GrpcPort)
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

	s.logger.Info("port", s.config.GrpcPort, "address", listener.Addr().String(), "GRPC Listening")
	if err := gs.Serve(listener); err != nil {
		return fmt.Errorf("could not start GRPC server")
	}

	return nil
}

func (s *DispersalServer) updateQuorumCount(ctx context.Context) error {
	currentBlock, err := s.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return err
	}
	count, err := s.tx.GetQuorumCount(ctx, currentBlock)
	if err != nil {
		return err
	}

	s.logger.Debug("updating quorum count", "currentBlock", currentBlock, "count", count)
	s.mu.Lock()
	s.quorumCount = count
	s.mu.Unlock()
	return nil
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

func getBlobFromRequest(req *pb.DisperseBlobRequest) *core.Blob {
	params := make([]*core.SecurityParam, len(req.SecurityParams))

	for i, param := range req.GetSecurityParams() {
		params[i] = &core.SecurityParam{
			QuorumID:           core.QuorumID(param.QuorumId),
			AdversaryThreshold: uint8(param.AdversaryThreshold),
			QuorumThreshold:    uint8(param.QuorumThreshold),
		}
	}

	data := req.GetData()

	blob := &core.Blob{
		RequestHeader: core.BlobRequestHeader{
			BlobAuthHeader: core.BlobAuthHeader{
				AccountID: req.AccountId,
			},
			SecurityParams: params,
		},
		Data: data,
	}

	return blob
}
