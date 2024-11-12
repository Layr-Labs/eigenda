package grpc

import (
	"context"
	"encoding/hex"
	"fmt"
	"runtime"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/shirou/gopsutil/mem"
)

// ServerV2 implements the Node v2 proto APIs.
type ServerV2 struct {
	pb.UnimplementedDispersalServer
	pb.UnimplementedRetrievalServer

	config      *node.Config
	node        *node.Node
	ratelimiter common.RateLimiter
	logger      logging.Logger
}

// NewServerV2 creates a new Server instance with the provided parameters.
func NewServerV2(
	config *node.Config,
	node *node.Node,
	logger logging.Logger,
	ratelimiter common.RateLimiter,
) *ServerV2 {
	return &ServerV2{
		config:      config,
		node:        node,
		ratelimiter: ratelimiter,
		logger:      logger,
	}
}

func (s *ServerV2) NodeInfo(ctx context.Context, in *pb.NodeInfoRequest) (*pb.NodeInfoReply, error) {
	if s.config.DisableNodeInfoResources {
		return &pb.NodeInfoReply{Semver: node.SemVer}, nil
	}

	memBytes := uint64(0)
	v, err := mem.VirtualMemory()
	if err == nil {
		memBytes = v.Total
	}

	return &pb.NodeInfoReply{Semver: node.SemVer, Os: runtime.GOOS, Arch: runtime.GOARCH, NumCpu: uint32(runtime.GOMAXPROCS(0)), MemBytes: memBytes}, nil
}

func (s *ServerV2) StoreChunks(ctx context.Context, in *pb.StoreChunksRequest) (*pb.StoreChunksReply, error) {
	batch, err := s.validateStoreChunksRequest(in)
	if err != nil {
		return nil, err
	}

	batchHeaderHash, err := batch.BatchHeader.Hash()
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("invalid batch header: %v", err))
	}

	operatorState, err := s.node.ChainState.GetOperatorStateByOperator(ctx, uint(batch.BatchHeader.ReferenceBlockNumber), s.node.Config.ID)
	if err != nil {
		return nil, err
	}

	blobShards, rawBundles, err := s.node.DownloadBundles(ctx, batch, operatorState)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to download batch: %v", err))
	}

	type storeResult struct {
		keys []kvstore.Key
		err  error
	}
	storeChan := make(chan storeResult)
	go func() {
		keys, err := s.node.StoreV2.StoreBatch(batch, rawBundles)
		if err != nil {
			storeChan <- storeResult{
				keys: nil,
				err:  fmt.Errorf("failed to store batch: %v", err),
			}
			return
		}

		storeChan <- storeResult{
			keys: keys,
			err:  nil,
		}
	}()

	err = s.node.ValidateBatchV2(ctx, batch, blobShards, operatorState)
	if err != nil {
		res := <-storeChan
		if len(res.keys) > 0 {
			if deleteErr := s.node.StoreV2.DeleteKeys(res.keys); deleteErr != nil {
				s.logger.Error("failed to delete keys", "err", deleteErr, "batchHeaderHash", hex.EncodeToString(batchHeaderHash[:]))
			}
		}
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to validate batch: %v", err))
	}

	res := <-storeChan
	if res.err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to store batch: %v", res.err))
	}

	sig := s.node.KeyPair.SignMessage(batchHeaderHash).Bytes()
	return &pb.StoreChunksReply{
		Signature: sig[:],
	}, nil
}

// validateStoreChunksRequest validates the StoreChunksRequest and returns deserialized batch in the request
func (s *ServerV2) validateStoreChunksRequest(req *pb.StoreChunksRequest) (*corev2.Batch, error) {
	if req.GetBatch() == nil {
		return nil, api.NewErrorInvalidArg("missing batch in request")
	}

	batch, err := corev2.BatchFromProtobuf(req.GetBatch())
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to deserialize batch: %v", err))
	}

	return batch, nil
}

func (s *ServerV2) GetChunks(ctx context.Context, in *pb.GetChunksRequest) (*pb.GetChunksReply, error) {
	blobKey, err := corev2.BytesToBlobKey(in.GetBlobKey())
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("invalid blob key: %v", err))
	}

	if corev2.MaxQuorumID < in.GetQuorumId() {
		return nil, api.NewErrorInvalidArg("invalid quorum ID")
	}
	quorumID := core.QuorumID(in.GetQuorumId())
	chunks, err := s.node.StoreV2.GetChunks(blobKey, quorumID)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get chunks: %v", err))
	}

	return &pb.GetChunksReply{
		Chunks: chunks,
	}, nil
}
