package validator

import (
	"context"
	"fmt"

	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO factory?

// A ValidatorGRPCManager is responsible for maintaining gRPC client connections with the validator nodes.
type ValidatorGRPCManager interface {

	// DownloadChunks downloads chunks from a validator node.
	DownloadChunks(
		ctx context.Context,
		key v2.BlobKey,
		operatorID core.OperatorID,
		quorumID core.QuorumID,
	) (*grpcnode.GetChunksReply, error)
}

var _ ValidatorGRPCManager = (*validatorGRPCManager)(nil)

// validatorGRPCManager is a standalone implementation of the ValidatorGRPCManager interface.
type validatorGRPCManager struct {
	logger logging.Logger

	// Information about the operators for each quorum.
	operatorInfo map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo
}

// NewValidatorGRPCManager creates a new ValidatorGRPCManager instance.
func NewValidatorGRPCManager(
	logger logging.Logger,
	operatorInfo map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo,
) ValidatorGRPCManager {
	return &validatorGRPCManager{
		logger:       logger,
		operatorInfo: operatorInfo,
	}
}

func (m *validatorGRPCManager) DownloadChunks(
	ctx context.Context,
	key v2.BlobKey,
	operatorID core.OperatorID,
	quorumID core.QuorumID,
) (*grpcnode.GetChunksReply, error) {

	// TODO(cody.littley) we can get a tighter bound?
	maxBlobSize := 16 * units.MiB // maximum size of the original blob
	encodingRate := 8             // worst case scenario if one validator has 100% stake
	fudgeFactor := units.MiB      // to allow for some overhead from things like protobuf encoding
	maxMessageSize := maxBlobSize*encodingRate + fudgeFactor

	quorumOpInfo, ok := m.operatorInfo[quorumID]
	if !ok {
		return nil, fmt.Errorf("quorum %d not found", quorumID)
	}
	opInfo, ok := quorumOpInfo[operatorID]
	if !ok {
		return nil, fmt.Errorf("operator %s not found in quorum %d", operatorID.Hex(), quorumID)
	}

	conn, err := grpc.NewClient(
		opInfo.Socket.GetV2RetrievalSocket(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize)),
	)
	defer func() {
		err := conn.Close()
		if err != nil {
			m.logger.Error("validator retriever failed to close connection", "err", err)
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to operator %s: %w", operatorID.Hex(), err)
	}

	client := grpcnode.NewRetrievalClient(conn)
	request := &grpcnode.GetChunksRequest{
		BlobKey:  key[:],
		QuorumId: uint32(quorumID),
	}

	reply, err := client.GetChunks(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunks from operator %s: %w", operatorID.Hex(), err)
	}

	return reply, nil
}
