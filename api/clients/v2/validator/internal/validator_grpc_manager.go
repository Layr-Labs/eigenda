package internal

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

// A ValidatorGRPCManager is responsible for maintaining gRPC client connections with the validator nodes.
type ValidatorGRPCManager interface {

	// DownloadChunks downloads chunks from a validator node.
	DownloadChunks(
		ctx context.Context,
		key v2.BlobKey,
		operatorID core.OperatorID,
	) (*grpcnode.GetChunksReply, error)
}

// ValidatorGRPCManagerFactory is a function that creates a new ValidatorGRPCManager instance.
type ValidatorGRPCManagerFactory func(
	logger logging.Logger,
	socketMap map[core.OperatorID]core.OperatorSocket,
) ValidatorGRPCManager

var _ ValidatorGRPCManager = &validatorGRPCManager{}

// validatorGRPCManager is a standard implementation of the ValidatorGRPCManager interface.
type validatorGRPCManager struct {
	logger logging.Logger

	// Information about the operators for each quorum.
	socketMap map[core.OperatorID]core.OperatorSocket
}

var _ ValidatorGRPCManagerFactory = NewValidatorGRPCManager

// NewValidatorGRPCManager creates a new ValidatorGRPCManager instance.
func NewValidatorGRPCManager(
	logger logging.Logger,
	socketMap map[core.OperatorID]core.OperatorSocket,
) ValidatorGRPCManager {
	return &validatorGRPCManager{
		logger:    logger,
		socketMap: socketMap,
	}
}

func (m *validatorGRPCManager) DownloadChunks(
	ctx context.Context,
	key v2.BlobKey,
	operatorID core.OperatorID,
) (*grpcnode.GetChunksReply, error) {

	// TODO(cody.littley) we can get a tighter bound?
	maxBlobSize := 16 * units.MiB // maximum size of the original blob
	encodingRate := 8             // worst case scenario if one validator has 100% stake
	fudgeFactor := units.MiB      // to allow for some overhead from things like protobuf encoding
	maxMessageSize := maxBlobSize*encodingRate + fudgeFactor

	socket, ok := m.socketMap[operatorID]
	if !ok {
		return nil, fmt.Errorf("operator %s not found in socket map", operatorID.Hex())
	}

	conn, err := grpc.NewClient(
		socket.GetV2RetrievalSocket(),
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
		BlobKey: key[:],
	}

	reply, err := client.GetChunks(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunks from operator %s: %w", operatorID.Hex(), err)
	}

	return reply, nil
}
