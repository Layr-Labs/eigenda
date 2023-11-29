package apiserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/common"
	healthcheck "github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var errSystemRateLimit = fmt.Errorf("request ratelimited: system limit")
var errAccountRateLimit = fmt.Errorf("request ratelimited: account limit")

const systemAccountKey = "system"

const maxBlobSize = 1024 * 512 // 512 KiB

type DispersalServer struct {
	pb.UnimplementedDisperserServer
	mu *sync.Mutex

	config disperser.ServerConfig

	blobStore   disperser.BlobStore
	tx          core.Transactor
	quorumCount uint16

	rateConfig  RateConfig
	ratelimiter common.RateLimiter

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
	return &DispersalServer{
		config:      config,
		blobStore:   store,
		tx:          tx,
		quorumCount: 0,
		metrics:     metrics,
		logger:      logger,
		ratelimiter: ratelimiter,
		rateConfig:  rateConfig,
		mu:          &sync.Mutex{},
	}
}

func (s *DispersalServer) DisperseBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*pb.DisperseBlobReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("DisperseBlob", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	securityParams := req.GetSecurityParams()
	if len(securityParams) == 0 {
		return nil, fmt.Errorf("invalid request: security_params must not be empty")
	}
	if len(securityParams) > 256 {
		return nil, fmt.Errorf("invalid request: security_params must not exceed 256")
	}

	seenQuorums := make(map[uint32]struct{})
	// The quorum ID must be in range [0, 255]. It'll actually be converted
	// to uint8, so it cannot be greater than 255.
	for _, param := range securityParams {
		if _, ok := seenQuorums[param.QuorumId]; ok {
			return nil, fmt.Errorf("invalid request: security_params must not contain duplicate quorum_id")
		}
		seenQuorums[param.QuorumId] = struct{}{}

		if param.GetQuorumId() >= uint32(s.quorumCount) {
			err := s.updateQuorumCount(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get onchain quorum count: %w", err)
			}

			if param.GetQuorumId() >= uint32(s.quorumCount) {
				return nil, fmt.Errorf("invalid request: the quorum_id must be in range [0, %d], but found %d", s.quorumCount-1, param.GetQuorumId())
			}
		}
	}

	blobSize := len(req.GetData())
	// The blob size in bytes must be in range [1, maxBlobSize].
	if blobSize > maxBlobSize {
		return nil, fmt.Errorf("blob size cannot exceed 512 KiB")
	}
	if blobSize == 0 {
		return nil, fmt.Errorf("blob size must be greater than 0")
	}

	blob := getBlobFromRequest(req)

	origin, err := common.GetClientAddress(ctx, s.rateConfig.ClientIPHeader, 2, true)
	if err != nil {
		for _, param := range securityParams {
			quorumId := string(uint8(param.GetQuorumId()))
			s.metrics.HandleFailedRequest(quorumId, blobSize, "DisperseBlob")
		}
		return nil, err
	}

	s.logger.Debug("received a new blob request", "origin", origin, "securityParams", securityParams)

	if err := blob.RequestHeader.Validate(); err != nil {
		s.logger.Warn("invalid header", "err", err)
		for _, param := range securityParams {
			quorumId := string(uint8(param.GetQuorumId()))
			s.metrics.HandleFailedRequest(quorumId, blobSize, "DisperseBlob")
		}
		return nil, err
	}

	if s.ratelimiter != nil {
		err := s.checkRateLimitsAndAddRates(ctx, blob, origin)
		if err != nil {
			for _, param := range securityParams {
				quorumId := string(uint8(param.GetQuorumId()))
				if errors.Is(err, errSystemRateLimit) {
					s.metrics.HandleSystemRateLimitedRequest(quorumId, blobSize, "DisperseBlob")
				} else if errors.Is(err, errAccountRateLimit) {
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
			quorumId := string(uint8(param.GetQuorumId()))
			s.metrics.HandleFailedRequest(quorumId, blobSize, "DisperseBlob")
		}
		return nil, err
	}

	for _, param := range securityParams {
		quorumId := string(uint8(param.GetQuorumId()))
		s.metrics.HandleSuccessfulRequest(quorumId, blobSize, "DisperseBlob")
	}

	s.logger.Info("received a new blob: ", "key", metadataKey.String())
	return &pb.DisperseBlobReply{
		Result:    pb.BlobStatus_PROCESSING,
		RequestId: []byte(metadataKey.String()),
	}, nil
}

func (s *DispersalServer) checkRateLimitsAndAddRates(ctx context.Context, blob *core.Blob, origin string) error {

	for _, param := range blob.RequestHeader.SecurityParams {

		rates, ok := s.rateConfig.QuorumRateInfos[param.QuorumID]
		if !ok {
			return fmt.Errorf("no configured rate exists for quorum %d", param.QuorumID)
		}

		// Get the encoded blob size from the blob header. Calculation is done in a way that nodes can replicate
		blobSize := len(blob.Data)
		length := core.GetBlobLength(uint(blobSize))
		encodedLength := core.GetEncodedBlobLength(length, uint8(blob.RequestHeader.SecurityParams[param.QuorumID].QuorumThreshold), uint8(blob.RequestHeader.SecurityParams[param.QuorumID].AdversaryThreshold))
		encodedSize := core.GetBlobSize(encodedLength)

		s.logger.Debug("checking rate limits", "origin", origin, "quorum", param.QuorumID, "encodedSize", encodedSize, "blobSize", blobSize)

		// Check System Ratelimit
		systemQuorumKey := fmt.Sprintf("%s:%d", systemAccountKey, param.QuorumID)
		allowed, err := s.ratelimiter.AllowRequest(ctx, systemQuorumKey, encodedSize, rates.TotalUnauthThroughput)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("system ratelimit exceeded", "systemQuorumKey", systemQuorumKey, "rate", rates.TotalUnauthThroughput)
			return errSystemRateLimit
		}

		blob.RequestHeader.AccountID = "ip:" + origin

		userQuorumKey := fmt.Sprintf("%s:%d", blob.RequestHeader.AccountID, param.QuorumID)
		allowed, err = s.ratelimiter.AllowRequest(ctx, userQuorumKey, encodedSize, rates.PerUserUnauthThroughput)
		if err != nil {
			return fmt.Errorf("ratelimiter error: %v", err)
		}
		if !allowed {
			s.logger.Warn("account ratelimit exceeded", "userQuorumKey", userQuorumKey, "rate", rates.PerUserUnauthThroughput)
			return errAccountRateLimit
		}

		// Update the quorum rate
		blob.RequestHeader.SecurityParams[param.QuorumID].QuorumRate = rates.PerUserUnauthThroughput
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
				QuantizationParam:            uint32(quorumInfo.QuantizationFactor),
				EncodedLength:                uint64(quorumInfo.EncodedBlobLength),
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
	healthcheck.RegisterHealthServer(gs)

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
			SecurityParams: params,
		},
		Data: data,
	}

	return blob
}
