package churner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedChurnerServer

	config  *Config
	churner *churner
	// the signature with the lastest expiry
	latestExpiry                int64
	lastRequestTimeByOperatorID map[core.OperatorID]time.Time

	logger  logging.Logger
	metrics *Metrics
}

func NewServer(
	config *Config,
	churner *churner,
	logger logging.Logger,
	metrics *Metrics,
) *Server {
	return &Server{
		config:                      config,
		churner:                     churner,
		latestExpiry:                int64(0),
		lastRequestTimeByOperatorID: make(map[core.OperatorID]time.Time),
		logger:                      logger.With("component", "ChurnerServer"),
		metrics:                     metrics,
	}
}

func (s *Server) Start(metricsConfig MetricsConfig) error {
	// Enable Metrics Block
	if metricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", metricsConfig.HTTPPort)
		s.metrics.Start(context.Background())
		s.logger.Info("Enabled metrics for Churner", "socket", httpSocket)
	}
	return nil
}

func (s *Server) Churn(ctx context.Context, req *pb.ChurnRequest) (*pb.ChurnReply, error) {
	s.logger.Info("Churn request received", "operator", req.GetOperatorAddress(), "quorumIds", req.GetQuorumIds())
	
	err := s.validateChurnRequest(ctx, req)
	if err != nil {
		s.logger.Error("Churn request validation failed", "error", err, "operator", req.GetOperatorAddress())
		s.metrics.IncrementFailedRequestNum("Churn", FailReasonInvalidRequest)
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("invalid request: %s", err.Error()))
	}
	s.logger.Info("Churn request validated successfully", "operator", req.GetOperatorAddress())

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("Churn", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()
	s.logger.Info("Received request: ", "QuorumIds", req.GetQuorumIds())

	now := time.Now()
	// Global rate limiting: check that we are after the previous approval's expiry
	s.logger.Info("Checking global rate limit", "now", now.Unix(), "latestExpiry", s.latestExpiry)
	if now.Unix() < s.latestExpiry {
		s.logger.Warn("Previous approval not expired", "retryIn", s.latestExpiry-now.Unix())
		s.metrics.IncrementFailedRequestNum("Churn", FailReasonPrevApprovalNotExpired)
		return nil, api.NewErrorResourceExhausted(fmt.Sprintf("previous approval not expired, retry in %d seconds", s.latestExpiry-now.Unix()))
	}

	s.logger.Info("Creating churn request object")
	request, err := createChurnRequest(req)
	if err != nil {
		s.logger.Error("Failed to create churn request", "error", err)
		s.metrics.IncrementFailedRequestNum("Churn", FailReasonInvalidRequest)
		return nil, api.NewErrorInvalidArg(err.Error())
	}

	s.logger.Info("Verifying request signature")
	operatorToRegisterAddress, err := s.churner.VerifyRequestSignature(ctx, request)
	if err != nil {
		s.logger.Error("Failed to verify request signature", "error", err)
		s.metrics.IncrementFailedRequestNum("Churn", FailReasonInvalidSignature)
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to verify request signature: %s", err.Error()))
	}
	s.logger.Info("Request signature verified", "operatorAddress", operatorToRegisterAddress.Hex())

	// Per-operator rate limiting: check if the request should be rate limited
	s.logger.Info("Checking per-operator rate limit")
	err = s.checkShouldBeRateLimited(now, *request)
	if err != nil {
		s.logger.Warn("Rate limit exceeded", "error", err)
		s.metrics.IncrementFailedRequestNum("Churn", FailReasonRateLimitExceeded)
		return nil, api.NewErrorResourceExhausted(fmt.Sprintf("rate limiter error: %s", err.Error()))
	}

	s.logger.Info("Processing churn request", "operator", operatorToRegisterAddress.Hex())
	response, err := s.churner.ProcessChurnRequest(ctx, operatorToRegisterAddress, request)
	if err != nil {
		s.logger.Error("Failed to process churn request", "error", err)
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		s.metrics.IncrementFailedRequestNum("Churn", FailReasonProcessChurnRequestFailed)
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to process churn request: %s", err.Error()))
	}

	// update the latest expiry
	s.latestExpiry = response.SignatureWithSaltAndExpiry.Expiry.Int64()

	operatorsToChurn := convertToOperatorsToChurnGrpc(response.OperatorsToChurn)

	s.metrics.IncrementSuccessfulRequestNum("Churn")
	return &pb.ChurnReply{
		SignatureWithSaltAndExpiry: &pb.SignatureWithSaltAndExpiry{
			Signature: response.SignatureWithSaltAndExpiry.Signature,
			Salt:      response.SignatureWithSaltAndExpiry.Salt[:],
			Expiry:    response.SignatureWithSaltAndExpiry.Expiry.Int64(),
		},
		OperatorsToChurn: operatorsToChurn,
	}, nil
}

func (s *Server) checkShouldBeRateLimited(now time.Time, request ChurnRequest) error {
	operatorToRegisterId := request.OperatorToRegisterPubkeyG1.GetOperatorID()
	lastRequestTimestamp := s.lastRequestTimeByOperatorID[operatorToRegisterId]
	if now.Unix() < lastRequestTimestamp.Add(s.config.PerPublicKeyRateLimit).Unix() {
		return fmt.Errorf("operatorID Rate Limit Exceeded: %d", operatorToRegisterId)
	}
	s.lastRequestTimeByOperatorID[operatorToRegisterId] = now
	return nil
}

func (s *Server) validateChurnRequest(ctx context.Context, req *pb.ChurnRequest) error {

	if len(req.GetOperatorRequestSignature()) != 64 {
		return errors.New("invalid signature length")
	}

	if len(req.GetOperatorToRegisterPubkeyG1()) != 64 {
		return errors.New("invalid operatorToRegisterPubkeyG1 length")
	}

	if len(req.GetOperatorToRegisterPubkeyG2()) != 128 {
		return errors.New("invalid operatorToRegisterPubkeyG2 length")
	}

	if len(req.GetSalt()) != 32 {
		return errors.New("invalid salt length")
	}

	// TODO: ensure that all quorumIDs are valid
	if len(req.GetQuorumIds()) == 0 || len(req.GetQuorumIds()) > 255 {
		return fmt.Errorf("invalid quorumIds length %d", len(req.GetQuorumIds()))
	}

	seenQuorums := make(map[int]struct{})
	for quorumID := range req.GetQuorumIds() {
		// make sure there are no duplicate quorum IDs
		if _, ok := seenQuorums[quorumID]; ok {
			return errors.New("invalid request: security_params must not contain duplicate quorum_id")
		}
		seenQuorums[quorumID] = struct{}{}

		if quorumID >= int(s.churner.QuorumCount) {
			err := s.churner.UpdateQuorumCount(ctx)
			if err != nil {
				return fmt.Errorf("failed to get onchain quorum count: %w", err)
			}

			if quorumID >= int(s.churner.QuorumCount) {
				return fmt.Errorf("invalid request: the quorum_id must be in range [0, %d], but found %d", s.churner.QuorumCount-1, quorumID)
			}
		}
	}

	return nil

}

func createChurnRequest(req *pb.ChurnRequest) (*ChurnRequest, error) {

	sigPoint, err := new(core.G1Point).Deserialize(req.GetOperatorRequestSignature())
	if err != nil {
		return nil, err
	}
	signature := &core.Signature{G1Point: sigPoint}

	address := gethcommon.HexToAddress(req.GetOperatorAddress())

	salt := [32]byte{}
	copy(salt[:], req.GetSalt())

	quorumIDs := make([]core.QuorumID, len(req.GetQuorumIds()))
	for i, id := range req.GetQuorumIds() {
		quorumIDs[i] = core.QuorumID(id)
	}

	pubkeyG1, err := new(core.G1Point).Deserialize(req.GetOperatorToRegisterPubkeyG1())
	if err != nil {
		return nil, err
	}
	pubkeyG2, err := new(core.G2Point).Deserialize(req.GetOperatorToRegisterPubkeyG2())
	if err != nil {
		return nil, err
	}

	return &ChurnRequest{
		OperatorAddress:            address,
		OperatorToRegisterPubkeyG1: pubkeyG1,
		OperatorToRegisterPubkeyG2: pubkeyG2,
		OperatorRequestSignature:   signature,
		Salt:                       salt,
		QuorumIDs:                  quorumIDs,
	}, nil
}

func convertToOperatorsToChurnGrpc(operatorsToChurn []core.OperatorToChurn) []*pb.OperatorToChurn {
	operatorsToChurnGRPC := make([]*pb.OperatorToChurn, len(operatorsToChurn))
	for i, operator := range operatorsToChurn {
		var pubkey []byte
		if operator.Pubkey != nil {
			pubkey = operator.Pubkey.Serialize()
		}
		operatorsToChurnGRPC[i] = &pb.OperatorToChurn{
			Operator: operator.Operator.Bytes(),
			QuorumId: uint32(operator.QuorumId),
			Pubkey:   pubkey,
		}
	}
	return operatorsToChurnGRPC
}
