package grpc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/mem"
	"google.golang.org/grpc/peer"
)

// ServerV2 implements the Node v2 proto APIs.
type ServerV2 struct {
	pb.UnimplementedDispersalServer
	pb.UnimplementedRetrievalServer

	config        *node.Config
	node          *node.Node
	ratelimiter   common.RateLimiter
	logger        logging.Logger
	metrics       *MetricsV2
	authenticator auth.RequestAuthenticator
}

// NewServerV2 creates a new Server instance with the provided parameters.
func NewServerV2(
	ctx context.Context,
	config *node.Config,
	node *node.Node,
	logger logging.Logger,
	ratelimiter common.RateLimiter,
	registry *prometheus.Registry,
	reader core.Reader) (*ServerV2, error) {

	metrics, err := NewV2Metrics(logger, registry)
	if err != nil {
		return nil, err
	}

	var authenticator auth.RequestAuthenticator
	if !config.DisableDispersalAuthentication {
		authenticator, err = auth.NewRequestAuthenticator(
			ctx,
			reader,
			config.DispersalAuthenticationKeyCacheSize,
			config.DisperserKeyTimeout,
			config.DispersalAuthenticationTimeout,
			func(id uint32) bool {
				return id == api.EigenLabsDisperserID
			},
			time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to create authenticator: %w", err)
		}
	}

	return &ServerV2{
		config:        config,
		node:          node,
		ratelimiter:   ratelimiter,
		logger:        logger,
		metrics:       metrics,
		authenticator: authenticator,
	}, nil
}

func (s *ServerV2) GetNodeInfo(ctx context.Context, in *pb.GetNodeInfoRequest) (*pb.GetNodeInfoReply, error) {
	if s.config.DisableNodeInfoResources {
		return &pb.GetNodeInfoReply{Semver: node.SemVer}, nil
	}

	memBytes := uint64(0)
	v, err := mem.VirtualMemory()
	if err == nil {
		memBytes = v.Total
	}

	return &pb.GetNodeInfoReply{
		Semver:   node.SemVer,
		Os:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		NumCpu:   uint32(runtime.GOMAXPROCS(0)),
		MemBytes: memBytes,
	}, nil
}

func (s *ServerV2) StoreChunks(ctx context.Context, in *pb.StoreChunksRequest) (*pb.StoreChunksReply, error) {
	start := time.Now()

	if !s.config.EnableV2 {
		return nil, api.NewErrorInvalidArg("v2 API is disabled")
	}

	if s.node.StoreV2 == nil {
		return nil, api.NewErrorInternal("v2 store not initialized")
	}

	if s.node.BLSSigner == nil {
		return nil, api.NewErrorInternal("missing bls signer")
	}

	batch, err := s.validateStoreChunksRequest(in)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to validate store chunk request: %v", err))
	}

	batchHeaderHash, err := batch.BatchHeader.Hash()
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to serialize batch header hash: %v", err))
	}

	if s.authenticator != nil {
		disperserPeer, ok := peer.FromContext(ctx)
		if !ok {
			return nil, api.NewErrorInvalidArg("could not get peer information from request context")
		}
		disperserAddress := disperserPeer.Addr.String()

		err := s.authenticator.AuthenticateStoreChunksRequest(ctx, disperserAddress, in, time.Now())
		if err != nil {
			return nil, api.WrapAsUnauthenticated(err, "failed to authenticate request")
		}
	}

	s.logger.Info("new StoreChunks request", "batchHeaderHash", hex.EncodeToString(batchHeaderHash[:]), "numBlobs", len(batch.BlobCertificates), "referenceBlockNumber", batch.BatchHeader.ReferenceBlockNumber)
	operatorState, err := s.node.ChainState.GetOperatorStateByOperator(ctx, uint(batch.BatchHeader.ReferenceBlockNumber), s.node.Config.ID)
	if err != nil {
		return nil, api.WrapAsInternal(err, "failed to get the operator state")
	}

	blobShards, rawBundles, err := s.node.DownloadBundles(ctx, batch, operatorState)
	if err != nil {
		return nil, api.WrapAsInternal(err, "failed to download batch")
	}

	type storeResult struct {
		keys []kvstore.Key
		err  error
	}
	storeChan := make(chan storeResult)
	go func() {
		keys, size, err := s.node.StoreV2.StoreBatch(batch, rawBundles)
		if err != nil {
			storeChan <- storeResult{
				keys: nil,
				err:  err,
			}
			return
		}

		s.metrics.ReportStoreChunksRequestSize(size)

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
		return nil, api.WrapAsInternal(err, "failed to validate batch")
	}

	res := <-storeChan
	if res.err != nil {
		return nil, api.WrapAsInternal(err, "failed to store batch")
	}

	sig, err := s.node.BLSSigner.Sign(ctx, batchHeaderHash[:])
	if err != nil {
		return nil, api.WrapAsInternal(err, "failed to sign batch")
	}

	s.metrics.ReportStoreChunksLatency(time.Since(start))

	return &pb.StoreChunksReply{
		Signature: sig,
	}, nil
}

// validateStoreChunksRequest validates the StoreChunksRequest and returns deserialized batch in the request
func (s *ServerV2) validateStoreChunksRequest(req *pb.StoreChunksRequest) (*corev2.Batch, error) {
	if req.GetDisperserID() != api.EigenLabsDisperserID {
		return nil, fmt.Errorf("disperserID is invalid: %d", req.GetDisperserID())
	}

	if len(req.GetSignature()) != 64 {
		return nil, errors.New("signature must be 64 bytes")
	}

	if req.GetBatch() == nil {
		return nil, errors.New("missing batch in request")
	}

	// BatchFromProtobuf internally validates the Batch while deserializing
	batch, err := corev2.BatchFromProtobuf(req.GetBatch())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize batch: %v", err)
	}

	return batch, nil
}

func (s *ServerV2) GetChunks(ctx context.Context, in *pb.GetChunksRequest) (*pb.GetChunksReply, error) {
	start := time.Now()

	if !s.config.EnableV2 {
		return nil, api.NewErrorInvalidArg("v2 API is disabled")
	}

	if s.node.StoreV2 == nil {
		return nil, api.NewErrorInternal("v2 store not initialized")
	}

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

	size := 0
	if len(chunks) > 0 {
		size = len(chunks[0]) * len(chunks)
	}
	s.metrics.ReportGetChunksDataSize(size)

	s.metrics.ReportGetChunksLatency(time.Since(start))

	return &pb.GetChunksReply{
		Chunks: chunks,
	}, nil
}
