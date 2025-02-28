package grpc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/shirou/gopsutil/mem"
	"github.com/wealdtech/go-merkletree/v2"
	"github.com/wealdtech/go-merkletree/v2/keccak256"

	_ "go.uber.org/automaxprocs"

	"google.golang.org/protobuf/proto"
)

// Server implements the Node proto APIs.
type Server struct {
	pb.UnimplementedDispersalServer
	pb.UnimplementedRetrievalServer

	node   *node.Node
	config *node.Config
	logger logging.Logger

	ratelimiter common.RateLimiter

	mu *sync.Mutex
}

// NewServer creates a new Server instance with the provided parameters.
//
// Note: The Server's chunks store will be created at config.DbPath+"/chunk".
func NewServer(config *node.Config, node *node.Node, logger logging.Logger, ratelimiter common.RateLimiter) *Server {

	return &Server{
		config:      config,
		logger:      logger,
		node:        node,
		ratelimiter: ratelimiter,
		mu:          &sync.Mutex{},
	}
}

func (s *Server) NodeInfo(ctx context.Context, in *pb.NodeInfoRequest) (*pb.NodeInfoReply, error) {
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

func (s *Server) handleStoreChunksRequest(ctx context.Context, in *pb.StoreChunksRequest) (*pb.StoreChunksReply, error) {
	start := time.Now()

	// Get batch header hash
	batchHeader, err := core.BatchHeaderFromProtobuf(in.GetBatchHeader())
	if err != nil {
		return nil, api.NewErrorInvalidArg(err.Error())
	}

	blobs, err := node.GetBlobMessages(in.GetBlobs(), s.node.Config.NumBatchDeserializationWorkers)
	if err != nil {
		return nil, err
	}

	s.node.Metrics.ObserveLatency("StoreChunks", "deserialization", float64(time.Since(start).Milliseconds()))
	s.node.Logger.Info("StoreChunksRequest deserialized", "duration", time.Since(start))

	sig, err := s.node.ProcessBatch(ctx, batchHeader, blobs, in.GetBlobs())
	if err != nil {
		return nil, err
	}

	sigData := sig.Serialize()

	return &pb.StoreChunksReply{Signature: sigData[:]}, nil
}

func (s *Server) validateStoreChunkRequest(in *pb.StoreChunksRequest) error {
	if in.GetBatchHeader() == nil {
		return api.NewErrorInvalidArg("missing batch_header in request")
	}
	if in.GetBatchHeader().GetBatchRoot() == nil {
		return api.NewErrorInvalidArg("missing batch_root in request")
	}
	if in.GetBatchHeader().GetReferenceBlockNumber() == 0 {
		return api.NewErrorInvalidArg("missing reference_block_number in request")
	}

	if len(in.GetBlobs()) == 0 {
		return api.NewErrorInvalidArg("missing blobs in request")
	}
	for _, blob := range in.Blobs {
		if blob.GetHeader() == nil {
			return api.NewErrorInvalidArg("missing blob header in request")
		}
		if node.ValidatePointsFromBlobHeader(blob.GetHeader()) != nil {
			return api.NewErrorInvalidArg("invalid points contained in the blob header in request")
		}
		if len(blob.GetHeader().GetQuorumHeaders()) == 0 {
			return api.NewErrorInvalidArg("missing quorum headers in request")
		}
		if len(blob.GetHeader().GetQuorumHeaders()) != len(blob.GetBundles()) {
			return api.NewErrorInvalidArg("the number of quorums must be the same as the number of bundles")
		}
		for _, q := range blob.GetHeader().GetQuorumHeaders() {
			if q.GetQuorumId() > core.MaxQuorumID {
				return api.NewErrorInvalidArg(fmt.Sprintf("quorum ID must be in range [0, %d], but found %d", core.MaxQuorumID, q.GetQuorumId()))
			}
			if err := core.ValidateSecurityParam(q.GetConfirmationThreshold(), q.GetAdversaryThreshold()); err != nil {
				return err
			}
		}
	}
	return nil
}

// StoreChunks is called by dispersers to store data.
func (s *Server) StoreChunks(ctx context.Context, in *pb.StoreChunksRequest) (*pb.StoreChunksReply, error) {
	start := time.Now()

	blobHeadersSize := 0
	bundleSize := 0
	for _, blob := range in.Blobs {
		blobHeadersSize += proto.Size(blob.GetHeader())
		for _, bundle := range blob.GetBundles() {
			bundleSize += proto.Size(bundle)
		}
	}
	s.node.Logger.Info("StoreChunks RPC request received", "num of blobs", len(in.Blobs), "request message size", proto.Size(in), "total size of blob headers", blobHeadersSize, "total size of bundles", bundleSize)

	// Validate the request.
	if err := s.validateStoreChunkRequest(in); err != nil {
		return nil, err
	}

	// Process the request.
	reply, err := s.handleStoreChunksRequest(ctx, in)

	// Record metrics.
	if err != nil {
		s.node.Metrics.RecordRPCRequest("StoreChunks", "failure", time.Since(start))
		s.node.Logger.Error("StoreChunks RPC failed", "duration", time.Since(start), "err", err)
	} else {
		s.node.Metrics.RecordRPCRequest("StoreChunks", "success", time.Since(start))
		s.node.Logger.Info("StoreChunks RPC succeeded", "duration", time.Since(start))
	}

	return reply, err
}

func (s *Server) StoreBlobs(ctx context.Context, in *pb.StoreBlobsRequest) (*pb.StoreBlobsReply, error) {
	return &pb.StoreBlobsReply{}, api.NewErrorUnimplemented()
}

func (s *Server) AttestBatch(ctx context.Context, in *pb.AttestBatchRequest) (*pb.AttestBatchReply, error) {
	return &pb.AttestBatchReply{}, api.NewErrorUnimplemented()
}

func (s *Server) RetrieveChunks(ctx context.Context, in *pb.RetrieveChunksRequest) (*pb.RetrieveChunksReply, error) {
	start := time.Now()

	if in.GetQuorumId() > core.MaxQuorumID {
		return nil, fmt.Errorf("invalid request: quorum ID must be in range [0, %d], but found %d", core.MaxQuorumID, in.GetQuorumId())
	}

	var batchHeaderHash [32]byte
	copy(batchHeaderHash[:], in.GetBatchHeaderHash())

	blobHeader, _, err := s.getBlobHeader(ctx, batchHeaderHash, int(in.GetBlobIndex()))
	if err != nil {
		return nil, err
	}

	retrieverID, err := common.GetClientAddress(ctx, s.config.ClientIPHeader, 1, false)
	if err != nil {
		return nil, err
	}

	quorumInfo := blobHeader.GetQuorumInfo(uint8(in.GetQuorumId()))
	if quorumInfo == nil {
		return nil, fmt.Errorf("invalid request: quorum ID %d not found in blob header", in.GetQuorumId())
	}
	encodedBlobSize := encoding.GetBlobSize(encoding.GetEncodedBlobLength(blobHeader.Length, quorumInfo.ConfirmationThreshold, quorumInfo.AdversaryThreshold))
	rate := quorumInfo.QuorumRate

	params := []common.RequestParams{
		{
			RequesterID: retrieverID,
			BlobSize:    encodedBlobSize,
			Rate:        rate,
		},
	}
	s.mu.Lock()
	allow, _, err := s.ratelimiter.AllowRequest(ctx, params)
	s.mu.Unlock()
	if err != nil {
		return nil, err
	}

	if !allow {
		return nil, errors.New("request rate limited")
	}

	chunks, format, err := s.node.Store.GetChunks(ctx, batchHeaderHash, int(in.GetBlobIndex()), uint8(in.GetQuorumId()))
	if err != nil {
		s.node.Metrics.RecordRPCRequest("RetrieveChunks", "failure", time.Since(start))
		return nil, fmt.Errorf("could not find chunks for batchHeaderHash %v, blob index: %v, quorumID: %v", hex.EncodeToString(batchHeaderHash[:]), in.GetBlobIndex(), in.GetQuorumId())
	}
	if !s.config.EnableGnarkBundleEncoding && format == pb.ChunkEncodingFormat_GNARK {
		s.node.Logger.Info("Converting chunks from Gnark back to Gob", "batchHeaderHash", hex.EncodeToString(batchHeaderHash[:]), "blobIndex", in.GetBlobIndex(), "quorumId", in.GetQuorumId())
		format = pb.ChunkEncodingFormat_GOB
		gobChunks := make([][]byte, 0, len(chunks))
		for _, c := range chunks {
			if len(c) == 0 {
				gobChunks = append(gobChunks, c)
				continue
			}
			decoded, err := new(encoding.Frame).DeserializeGnark(c)
			if err != nil {
				return nil, fmt.Errorf("the chunks are in Gnark but cannot be decoded: %v", err)
			}
			encoded, err := decoded.Serialize()
			if err != nil {
				return nil, err
			}
			gobChunks = append(gobChunks, encoded)
		}
		chunks = gobChunks
	}
	s.node.Metrics.RecordRPCRequest("RetrieveChunks", "success", time.Since(start))
	return &pb.RetrieveChunksReply{Chunks: chunks, ChunkEncodingFormat: format}, nil
}

func (s *Server) GetBlobHeader(ctx context.Context, in *pb.GetBlobHeaderRequest) (*pb.GetBlobHeaderReply, error) {
	var batchHeaderHash [32]byte
	copy(batchHeaderHash[:], in.GetBatchHeaderHash())

	blobHeader, protoBlobHeader, err := s.getBlobHeader(ctx, batchHeaderHash, int(in.GetBlobIndex()))
	if err != nil {
		return nil, err
	}

	blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
	if err != nil {
		return nil, err
	}

	tree, err := s.rebuildMerkleTree(batchHeaderHash)
	if err != nil {
		return nil, err
	}

	proof, err := tree.GenerateProof(blobHeaderHash[:], 0)
	if err != nil {
		return nil, err
	}

	return &pb.GetBlobHeaderReply{
		BlobHeader: protoBlobHeader,
		Proof: &pb.MerkleProof{
			Hashes: proof.Hashes,
			Index:  uint32(proof.Index),
		},
	}, nil
}

// rebuildMerkleTree rebuilds the merkle tree from the blob headers and batch header.
func (s *Server) rebuildMerkleTree(batchHeaderHash [32]byte) (*merkletree.MerkleTree, error) {
	batchHeaderBytes, err := s.node.Store.GetBatchHeader(context.Background(), batchHeaderHash)
	if err != nil {
		return nil, err
	}

	batchHeader, err := new(core.BatchHeader).Deserialize(batchHeaderBytes)
	if err != nil {
		return nil, err
	}

	blobIndex := 0
	leafs := make([][]byte, 0)
	for {
		blobHeaderBytes, err := s.node.Store.GetBlobHeader(context.Background(), batchHeaderHash, blobIndex)
		if err != nil {
			if errors.Is(err, node.ErrKeyNotFound) {
				break
			}
			return nil, err
		}

		var protoBlobHeader pb.BlobHeader
		err = proto.Unmarshal(blobHeaderBytes, &protoBlobHeader)
		if err != nil {
			return nil, err
		}

		blobHeader, err := core.BlobHeaderFromProtobuf(&protoBlobHeader)
		if err != nil {
			return nil, err
		}

		blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
		if err != nil {
			return nil, err
		}
		leafs = append(leafs, blobHeaderHash[:])
		blobIndex++
	}

	if len(leafs) == 0 {
		return nil, errors.New("no blob header found")
	}

	tree, err := merkletree.NewTree(merkletree.WithData(leafs), merkletree.WithHashType(keccak256.New()))
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(tree.Root(), batchHeader.BatchRoot[:]) {
		return nil, errors.New("invalid batch header")
	}

	return tree, nil
}

func (s *Server) getBlobHeader(ctx context.Context, batchHeaderHash [32]byte, blobIndex int) (*core.BlobHeader, *pb.BlobHeader, error) {

	blobHeaderBytes, err := s.node.Store.GetBlobHeader(ctx, batchHeaderHash, blobIndex)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get the blob header from Store: %w", err)
	}

	var protoBlobHeader pb.BlobHeader
	err = proto.Unmarshal(blobHeaderBytes, &protoBlobHeader)
	if err != nil {
		return nil, nil, err
	}

	blobHeader, err := core.BlobHeaderFromProtobuf(&protoBlobHeader)
	if err != nil {
		return nil, nil, err
	}

	return blobHeader, &protoBlobHeader, nil

}
