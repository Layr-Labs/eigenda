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
	quorumCount uint8

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

	blob := getBlobFromRequest(request.DisperseRequest)

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

	blob := getBlobFromRequest(req)

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

	err = s.validateBlobRequest(ctx, blob)
	if err != nil {
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
		quorumInfos := confirmationInfo.BlobQuorumInfos
		slices.SortStableFunc[[]*core.BlobQuorumInfo](quorumInfos, func(a, b *core.BlobQuorumInfo) int {
			return int(a.QuorumID) - int(b.QuorumID)
		})
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
