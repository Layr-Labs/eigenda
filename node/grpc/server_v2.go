package grpc

import (
	"context"
	"runtime"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/Layr-Labs/eigenda/common"
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
	return &pb.StoreChunksReply{}, api.NewErrorUnimplemented()
}

func (s *ServerV2) GetChunks(context.Context, *pb.GetChunksRequest) (*pb.GetChunksReply, error) {
	return &pb.GetChunksReply{}, api.NewErrorUnimplemented()
}
