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
	"github.com/Layr-Labs/eigenda/common/replay"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/mem"
)

// ServerV2 implements the Node v2 proto APIs.
type ServerV2 struct {
	pb.UnimplementedDispersalServer
	pb.UnimplementedRetrievalServer

	config         *node.Config
	node           *node.Node
	ratelimiter    common.RateLimiter
	logger         logging.Logger
	metrics        *MetricsV2
	authenticator  auth.RequestAuthenticator
	replayGuardian replay.ReplayGuardian
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
			func(id uint32) bool {
				return id == api.EigenLabsDisperserID
			},
			time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to create authenticator: %w", err)
		}
	}

	replayGuardian := replay.NewReplayGuardian(
		time.Now,
		config.StoreChunksRequestMaxPastAge,
		config.StoreChunksRequestMaxFutureAge)

	return &ServerV2{
		config:         config,
		node:           node,
		ratelimiter:    ratelimiter,
		logger:         logger,
		metrics:        metrics,
		authenticator:  authenticator,
		replayGuardian: replayGuardian,
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

	if s.node.BLSSigner == nil {
		return nil, api.NewErrorInternal("missing bls signer")
	}

	// Validate the request parameters (which is cheap) before starting any further
	// processing of the request.
	batch, err := s.validateStoreChunksRequest(in)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to validate store chunk request: %v", err))
	}

	batchHeaderHash, err := batch.BatchHeader.Hash()
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to serialize batch header hash: %v", err))
	}

	if s.authenticator != nil {
		hash, err := s.authenticator.AuthenticateStoreChunksRequest(ctx, in, time.Now())
		if err != nil {
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to authenticate request: %v", err))
		}

		timestamp := time.Unix(int64(in.Timestamp), 0)
		err = s.replayGuardian.VerifyRequest(hash, timestamp)
		if err != nil {
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to verify request: %v", err))
		}
	}

	s.logger.Info("new StoreChunks request", "batchHeaderHash", hex.EncodeToString(batchHeaderHash[:]), "numBlobs", len(batch.BlobCertificates), "referenceBlockNumber", batch.BatchHeader.ReferenceBlockNumber)
	operatorState, err := s.node.ChainState.GetOperatorStateByOperator(ctx, uint(batch.BatchHeader.ReferenceBlockNumber), s.node.Config.ID)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get the operator state: %v", err))
	}

	stageTimer := time.Now()
	blobShards, rawBundles, err := s.node.DownloadBundles(ctx, batch, operatorState)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get the operator state: %v", err))
	}
	s.metrics.ReportStoreChunksLatency("download", time.Since(stageTimer))

	err = s.validateAndStoreV2Chunks(ctx, batch, blobShards, rawBundles, operatorState, batchHeaderHash)
	if err != nil {
		return nil, err
	}

	sig, err := s.node.BLSSigner.Sign(ctx, batchHeaderHash[:])
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to sign batch: %v", err))
	}

	s.metrics.ReportStoreChunksLatency("total", time.Since(start))

	return &pb.StoreChunksReply{
		Signature: sig,
	}, nil
}

func (s *ServerV2) validateAndStoreV2Chunks(
	ctx context.Context,
	batch *corev2.Batch,
	blobShards []*corev2.BlobShard,
	rawBundles []*node.RawBundles,
	operatorState *core.OperatorState,
	batchHeaderHash [32]byte) error {

	batchData := make([]*node.BundleToStore, 0, len(rawBundles))
	for _, bundles := range rawBundles {
		blobKey, err := bundles.BlobCertificate.BlobHeader.BlobKey()
		if err != nil {
			return api.NewErrorInternal("failed to get blob key")
		}

		for quorum, bundle := range bundles.Bundles {
			bundleKey, err := node.BundleKey(blobKey, quorum)
			if err != nil {
				return api.NewErrorInternal("failed to get bundle key")
			}

			batchData = append(batchData, &node.BundleToStore{
				BundleKey:   bundleKey,
				BundleBytes: bundle,
			})
		}
	}

	if s.config.LittDBEnabled {
		return s.validateAndStoreV2ChunksLittDB(ctx, batch, blobShards, batchData, operatorState, batchHeaderHash)
	} else {
		return s.validateAndStoreV2ChunksLevelDB(ctx, batch, blobShards, batchData, operatorState, batchHeaderHash)
	}
}

func (s *ServerV2) validateAndStoreV2ChunksLevelDB(
	ctx context.Context,
	batch *corev2.Batch,
	blobShards []*corev2.BlobShard,
	batchData []*node.BundleToStore,
	operatorState *core.OperatorState,
	batchHeaderHash [32]byte) error {

	type storeResult struct {
		keys []kvstore.Key
		err  error
	}
	storeChan := make(chan storeResult)
	go func() {
		storageStart := time.Now()
		keys, size, err := s.node.StoreV2.StoreBatch(batchHeaderHash[:], batchData)
		if err != nil {
			storeChan <- storeResult{
				keys: nil,
				err:  err,
			}
			return
		}

		s.metrics.ReportStoreChunksRequestSize(size)
		s.metrics.ReportStoreChunksLatency("storage", time.Since(storageStart))
		storeChan <- storeResult{
			keys: keys,
			err:  nil,
		}
	}()

	stageTimer := time.Now()
	err := s.node.ValidateBatchV2(ctx, batch, blobShards, operatorState)
	if err != nil {
		res := <-storeChan
		if len(res.keys) > 0 {
			if deleteErr := s.node.StoreV2.DeleteKeys(res.keys); deleteErr != nil {
				s.logger.Error(
					"failed to delete keys",
					"err", deleteErr,
					"batchHeaderHash", hex.EncodeToString(batchHeaderHash[:]))
			}
		}
		return api.NewErrorInternal(fmt.Sprintf("failed to validate batch: %v", err))
	}
	s.metrics.ReportStoreChunksLatency("validation", time.Since(stageTimer))

	res := <-storeChan
	if res.err != nil {
		return api.NewErrorInternal(fmt.Sprintf("failed to store batch: %v", res.err))
	}

	return nil
}

func (s *ServerV2) validateAndStoreV2ChunksLittDB(
	ctx context.Context,
	batch *corev2.Batch,
	blobShards []*corev2.BlobShard,
	batchData []*node.BundleToStore,
	operatorState *core.OperatorState,
	batchHeaderHash [32]byte) error {

	stageTimer := time.Now()
	err := s.node.ValidateBatchV2(ctx, batch, blobShards, operatorState)
	if err != nil {
		return api.NewErrorInternal(fmt.Sprintf("failed to validate batch %s: %v", batchHeaderHash, err))
	}
	s.metrics.ReportStoreChunksLatency("validation", time.Since(stageTimer))

	storageStart := time.Now()
	_, size, err := s.node.StoreV2.StoreBatch(batchHeaderHash[:], batchData)
	if err != nil {
		return api.NewErrorInternal(fmt.Sprintf("failed to store batch %s: %v", batchHeaderHash, err))
	}

	s.metrics.ReportStoreChunksRequestSize(size)
	s.metrics.ReportStoreChunksLatency("storage", time.Since(storageStart))

	return nil
}

// validateStoreChunksRequest validates the StoreChunksRequest and returns deserialized batch in the request
func (s *ServerV2) validateStoreChunksRequest(req *pb.StoreChunksRequest) (*corev2.Batch, error) {
	// The signature is created by go-ethereum library, which contains 1 additional byte (for
	// recovering the public key from signature), so it's 65 bytes.
	if len(req.GetSignature()) != 65 {
		return nil, fmt.Errorf("signature must be 65 bytes, found %d bytes", len(req.GetSignature()))
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

	blobKey, err := corev2.BytesToBlobKey(in.GetBlobKey())
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("invalid blob key: %v", err))
	}

	if corev2.MaxQuorumID < in.GetQuorumId() {
		return nil, api.NewErrorInvalidArg("invalid quorum ID")
	}
	quorumID := core.QuorumID(in.GetQuorumId())

	bundleKey, err := node.BundleKey(blobKey, quorumID)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to get bundle key: %v", err))
	}

	bundleData, err := s.node.StoreV2.GetBundleData(bundleKey)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get chunks: %v", err))
	}

	chunks, _, err := node.DecodeChunks(bundleData)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to decode chunks: %v", err))
	}

	size := 0
	if len(chunks) > 0 {
		size = len(chunks[0]) * len(chunks)
	}
	s.metrics.ReportGetChunksDataSize(size)

	s.metrics.ReportGetChunksLatency(time.Since(start))

	return &pb.GetChunksReply{
		Chunks:              chunks,
		ChunkEncodingFormat: pb.ChunkEncodingFormat_GNARK,
	}, nil
}
