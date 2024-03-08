package churner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
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

	logger  common.Logger
	metrics *Metrics
}

func NewServer(
	config *Config,
	churner *churner,
	logger common.Logger,
	metrics *Metrics,
) *Server {
	return &Server{
		config:                      config,
		churner:                     churner,
		latestExpiry:                int64(0),
		lastRequestTimeByOperatorID: make(map[core.OperatorID]time.Time),
		logger:                      logger,
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
	err := s.validateChurnRequest(ctx, req)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("Churn", FailReasonInvalidRequest)
		return nil, api.NewInvalidArgError(fmt.Sprintf("invalid request: %s", err.Error()))
	}

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("Churn", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()
	s.logger.Info("Received request: ", "QuorumIds", req.GetQuorumIds())

	now := time.Now()
	// Global rate limiting: check that we are after the previous approval's expiry
	if now.Unix() < s.latestExpiry {
		s.metrics.IncrementFailedRequestNum("Churn", FailReasonPrevApprovalNotExpired)
		return nil, api.NewResourceExhaustedError(fmt.Sprintf("previous approval not expired, retry in %d", s.latestExpiry-now.Unix()))
	}

	request := createChurnRequest(req)

	operatorToRegisterAddress, err := s.churner.VerifyRequestSignature(ctx, request)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("Churn", FailReasonInvalidSignature)
		return nil, api.NewInvalidArgError(fmt.Sprintf("failed to verify request signature: %s", err.Error()))
	}

	// Per-operator rate limiting: check if the request should be rate limited
	err = s.checkShouldBeRateLimited(now, *request)
	if err != nil {
		s.metrics.IncrementFailedRequestNum("Churn", FailReasonRateLimitExceeded)
		return nil, api.NewResourceExhaustedError(fmt.Sprintf("rate limiter error: %s", err.Error()))
	}

	response, err := s.churner.ProcessChurnRequest(ctx, operatorToRegisterAddress, request)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		s.metrics.IncrementFailedRequestNum("Churn", FailReasonProcessChurnRequestFailed)
		return nil, api.NewInternalError(fmt.Sprintf("failed to process churn request: %s", err.Error()))
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

	if len(req.OperatorRequestSignature) != 64 {
		return errors.New("invalid signature length")
	}

	if len(req.OperatorToRegisterPubkeyG1) != 64 {
		return errors.New("invalid operatorToRegisterPubkeyG1 length")
	}

	if len(req.OperatorToRegisterPubkeyG2) != 128 {
		return errors.New("invalid operatorToRegisterPubkeyG2 length")
	}

	if len(req.Salt) != 32 {
		return errors.New("invalid salt length")
	}

	// TODO: ensure that all quorumIDs are valid
	if len(req.QuorumIds) == 0 {
		return errors.New("invalid quorumIds length")
	}

	for quorumID := range req.GetQuorumIds() {
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

func createChurnRequest(req *pb.ChurnRequest) *ChurnRequest {
	signature := &core.Signature{G1Point: new(core.G1Point).Deserialize(req.GetOperatorRequestSignature())}

	address := gethcommon.HexToAddress(req.GetOperatorAddress())

	salt := [32]byte{}
	copy(salt[:], req.GetSalt())

	quorumIDs := make([]core.QuorumID, len(req.QuorumIds))
	for i, id := range req.QuorumIds {
		quorumIDs[i] = core.QuorumID(id)
	}

	return &ChurnRequest{
		OperatorAddress:            address,
		OperatorToRegisterPubkeyG1: new(core.G1Point).Deserialize(req.GetOperatorToRegisterPubkeyG1()),
		OperatorToRegisterPubkeyG2: new(core.G2Point).Deserialize(req.GetOperatorToRegisterPubkeyG2()),
		OperatorRequestSignature:   signature,
		Salt:                       salt,
		QuorumIDs:                  quorumIDs,
	}
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
