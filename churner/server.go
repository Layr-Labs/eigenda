package churner

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
)

type Server struct {
	pb.UnimplementedChurnerServer

	config  *Config
	churner *churner
	// the signature with the lastest expiry
	latestExpiry                int64
	lastRequestTimeByOperatorID map[core.OperatorID]time.Time

	logger common.Logger
}

func NewServer(
	config *Config,
	churner *churner,
	logger common.Logger,
) *Server {
	return &Server{
		config:                      config,
		churner:                     churner,
		latestExpiry:                int64(0),
		lastRequestTimeByOperatorID: make(map[core.OperatorID]time.Time),
		logger:                      logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO: Start Metrics
	return nil
}

func (s *Server) Churn(ctx context.Context, req *pb.ChurnRequest) (*pb.ChurnReply, error) {
	s.logger.Info("Received request: ", "QuorumIds", req.GetQuorumIds())

	now := time.Now()
	// check that we are after the previous approval's expiry
	if now.Unix() < s.latestExpiry {
		return nil, fmt.Errorf("previous approval not expired, retry in %d", s.latestExpiry-now.Unix())
	}

	for quorumID := range req.GetQuorumIds() {
		if quorumID >= int(s.churner.QuorumCount) {
			err := s.churner.UpdateQuorumCount(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get onchain quorum count: %w", err)
			}

			if quorumID >= int(s.churner.QuorumCount) {
				return nil, fmt.Errorf("Invalid request: the quorum_id must be in range [0, %d], but found %d", s.churner.QuorumCount-1, quorumID)
			}
		}
	}

	request := createChurnRequest(req)

	operatorToRegisterAddress, err := s.churner.VerifyRequestSignature(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to verify request signature: %w", err)
	}

	// check if the request should be rate limited
	err = s.checkShouldBeRateLimited(now, *request)
	if err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	response, err := s.churner.ProcessChurnRequest(ctx, operatorToRegisterAddress, request)
	if err != nil {
		return nil, fmt.Errorf("failed to process churn request: %w", err)
	}

	// update the latest expiry
	s.latestExpiry = response.SignatureWithSaltAndExpiry.Expiry.Int64()

	operatorsToChurn := convertToOperatorsToChurnGrpc(response.OperatorsToChurn)

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

func createChurnRequest(req *pb.ChurnRequest) *ChurnRequest {
	signature := &core.Signature{G1Point: new(core.G1Point).Deserialize(req.GetOperatorRequestSignature())}

	salt := [32]byte{}
	copy(salt[:], req.GetSalt())

	quorumIDs := make([]core.QuorumID, len(req.QuorumIds))
	for i, id := range req.QuorumIds {
		quorumIDs[i] = core.QuorumID(id)
	}

	return &ChurnRequest{
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
		operatorsToChurnGRPC[i] = &pb.OperatorToChurn{
			Operator: operator.Operator[:],
			QuorumId: uint32(operator.QuorumId),
			Pubkey:   operator.Pubkey.Serialize(),
		}
	}
	return operatorsToChurnGRPC
}
