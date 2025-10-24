package relay

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func (s *Server) validateGetChunksRequest(request *pb.GetChunksRequest) *status.Status {
	if request == nil {
		return status.New(codes.InvalidArgument, "request is nil")
	}
	if len(request.GetChunkRequests()) == 0 {
		return status.New(codes.InvalidArgument, "no chunk requests provided")
	}
	if len(request.GetChunkRequests()) > s.config.MaxKeysPerGetChunksRequest {
		return status.Newf(codes.InvalidArgument,
			"too many chunk requests provided, max is %d", s.config.MaxKeysPerGetChunksRequest)
	}

	for _, chunkRequest := range request.GetChunkRequests() {
		if chunkRequest.GetByIndex() == nil && chunkRequest.GetByRange() == nil {
			return status.New(codes.InvalidArgument, "chunk request must be either by index or by range")
		}
	}

	return status.New(codes.OK, "")
}

// GetChunks retrieves chunks from blobs stored by the relay.
func (s *Server) GetChunks(ctx context.Context, request *pb.GetChunksRequest) (*pb.GetChunksReply, error) {
	reply, st := s.getChunks(ctx, request)
	api.LogResponseStatus(s.logger, st)
	if st != nil {
		// nolint:wrapcheck
		return reply, st.Err()
	}
	return reply, nil
}

func (s *Server) getChunks(ctx context.Context, request *pb.GetChunksRequest) (*pb.GetChunksReply, *status.Status) {
	start := time.Now()

	if s.config.Timeouts.GetChunksTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.config.Timeouts.GetChunksTimeout)
		defer cancel()
	}
	res := s.validateGetChunksRequest(request)
	if res != nil && res.Code() != codes.OK {
		return nil, res
	}

	s.metrics.ReportChunkKeyCount(len(request.GetChunkRequests()))

	if s.authenticator != nil {
		client, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.New(codes.InvalidArgument, "could not get peer information")
		}
		clientAddress := client.Addr.String()

		hash, err := s.authenticator.AuthenticateGetChunksRequest(ctx, request)
		if err != nil {
			s.metrics.ReportChunkAuthFailure()
			s.logger.Debug("rejected GetChunks request", "client", clientAddress)
			return nil, status.Newf(codes.InvalidArgument, "auth failed: %v", err)
		}

		timestamp := time.Unix(int64(request.GetTimestamp()), 0)
		err = s.replayGuardian.VerifyRequest(hash, timestamp)
		if err != nil {
			s.metrics.ReportChunkAuthFailure()
			return nil, status.Newf(codes.InvalidArgument, "failed to verify request: %v", err)
		}

		s.logger.Debug("received authenticated GetChunks request", "client", clientAddress)
	}

	finishedAuthenticating := time.Now()
	if s.authenticator != nil {
		s.metrics.ReportChunkAuthenticationLatency(finishedAuthenticating.Sub(start))
	}

	clientID := string(request.GetOperatorId())
	err := s.chunkRateLimiter.BeginGetChunkOperation(time.Now(), clientID)
	if err != nil {
		return nil, status.Newf(codes.ResourceExhausted, "rate limit exceeded: %v", err)
	}
	defer s.chunkRateLimiter.FinishGetChunkOperation(clientID)

	// keys might contain duplicate keys
	keys, err := getKeysFromChunkRequest(request)
	if err != nil {
		return nil, status.Newf(codes.InvalidArgument, "invalid request: %v", err)
	}

	mMap, err := s.metadataProvider.GetMetadataForBlobs(ctx, keys)
	if err != nil {
		if strings.Contains(err.Error(), blobstore.ErrMetadataNotFound.Error()) {
			// nolint:wrapcheck
			return nil, status.Newf(codes.NotFound,
				"blob not found, check if blob exists and is assigned to this relay:: %v", keys)
		}
		// nolint:wrapcheck
		return nil, status.Newf(codes.Internal, "error fetching metadata for blob: %v", err)
	}

	finishedFetchingMetadata := time.Now()
	s.metrics.ReportChunkMetadataLatency(finishedFetchingMetadata.Sub(finishedAuthenticating))

	requiredBandwidth, err := computeChunkRequestRequiredBandwidth(request, mMap)
	if err != nil {
		return nil, status.Newf(codes.Internal, "error computing required bandwidth: %v", err)
	}
	s.metrics.ReportGetChunksRequestedBandwidthUsage(requiredBandwidth)
	err = s.chunkRateLimiter.RequestGetChunkBandwidth(time.Now(), clientID, requiredBandwidth)
	if err != nil {
		if strings.Contains(err.Error(), "internal error") {
			return nil, status.New(codes.Internal, err.Error())
		}
		return nil, buildInsufficientGetChunksBandwidthError(request, requiredBandwidth, err)
	}
	s.metrics.ReportGetChunksBandwidthUsage(requiredBandwidth)

	frames, err := s.chunkProvider.GetFrames(ctx, mMap)
	if err != nil {
		return nil, status.Newf(codes.Internal, "error fetching frames: %v", err)
	}

	bytesToSend, err := gatherChunkDataToSend(frames, request)
	if err != nil {
		return nil, status.Newf(codes.Internal, "error gathering chunk data: %v", err)
	}

	s.metrics.ReportChunkDataLatency(time.Since(finishedFetchingMetadata))
	s.metrics.ReportChunkLatency(time.Since(start))

	return &pb.GetChunksReply{
		Data: bytesToSend,
	}, nil
}

// getKeysFromChunkRequest gathers a slice of blob keys from a GetChunks request.
func getKeysFromChunkRequest(request *pb.GetChunksRequest) ([]v2.BlobKey, error) {
	keys := make([]v2.BlobKey, 0, len(request.GetChunkRequests()))

	for _, chunkRequest := range request.GetChunkRequests() {
		var key v2.BlobKey
		if chunkRequest.GetByIndex() != nil {
			var err error
			key, err = v2.BytesToBlobKey(chunkRequest.GetByIndex().GetBlobKey())
			if err != nil {
				return nil, fmt.Errorf("invalid blob key: %w", err)
			}
		} else {
			var err error
			key, err = v2.BytesToBlobKey(chunkRequest.GetByRange().GetBlobKey())
			if err != nil {
				return nil, fmt.Errorf("invalid blob key: %w", err)
			}
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// gatherChunkDataToSend takes the chunk data and narrows it down to the data requested in the GetChunks request.
func gatherChunkDataToSend(
	frames map[v2.BlobKey]*core.ChunksData,
	request *pb.GetChunksRequest) ([][]byte, error) {

	bytesToSend := make([][]byte, 0, len(request.GetChunkRequests()))

	for _, chunkRequest := range request.GetChunkRequests() {
		var framesSubset *core.ChunksData
		var err error

		if chunkRequest.GetByIndex() != nil {
			framesSubset, err = selectFrameSubsetByIndex(chunkRequest.GetByIndex(), frames)
		} else {
			framesSubset, err = selectFrameSubsetByRange(chunkRequest.GetByRange(), frames)
		}

		if err != nil {
			return nil, fmt.Errorf("error selecting frame subset: %w", err)
		}

		subsetBytes, err := framesSubset.FlattenToBundle()
		if err != nil {
			return nil, fmt.Errorf("error serializing frame subset: %w", err)
		}

		bytesToSend = append(bytesToSend, subsetBytes)
	}

	return bytesToSend, nil
}

// selectFrameSubsetByRange selects a subset of frames from a BinaryFrames object based on a range
func selectFrameSubsetByRange(
	request *pb.ChunkRequestByRange,
	allFrames map[v2.BlobKey]*core.ChunksData) (*core.ChunksData, error) {

	key := v2.BlobKey(request.GetBlobKey())
	startIndex := request.GetStartIndex()
	endIndex := request.GetEndIndex()

	frames, ok := allFrames[key]
	if !ok {
		return nil, fmt.Errorf("frames not found for key %s", key.Hex())
	}

	if startIndex > endIndex {
		return nil, fmt.Errorf(
			"chunk range %d-%d is invalid for key %s, start index must be less than or equal to end index",
			startIndex, endIndex, key.Hex())
	}
	if endIndex > uint32(len(frames.Chunks)) {
		return nil, fmt.Errorf(
			"chunk range %d-%d is invald for key %s, chunk count %d",
			startIndex, endIndex, key, len(frames.Chunks))
	}

	framesSubset := &core.ChunksData{
		Chunks:   frames.Chunks[startIndex:endIndex],
		Format:   frames.Format,
		ChunkLen: frames.ChunkLen,
	}

	return framesSubset, nil
}

// selectFrameSubsetByIndex selects a subset of frames from a BinaryFrames object based on a list of indices
func selectFrameSubsetByIndex(
	request *pb.ChunkRequestByIndex,
	allFrames map[v2.BlobKey]*core.ChunksData) (*core.ChunksData, error) {

	key := v2.BlobKey(request.GetBlobKey())
	frames, ok := allFrames[key]
	if !ok {
		return nil, fmt.Errorf("frames not found for key %s", key.Hex())
	}

	if len(request.GetChunkIndices()) > len(frames.Chunks) {
		return nil, fmt.Errorf("too many requested chunks for key %s, chunk count %d",
			key.Hex(), len(frames.Chunks))
	}

	framesSubset := &core.ChunksData{
		Format:   frames.Format,
		ChunkLen: frames.ChunkLen,
		Chunks:   make([][]byte, 0, len(request.GetChunkIndices())),
	}

	for _, index := range request.GetChunkIndices() {
		if index >= uint32(len(frames.Chunks)) {
			return nil, fmt.Errorf(
				"chunk index %d out of range for key %s, chunk count %d",
				index, key.Hex(), len(frames.Chunks))
		}

		framesSubset.Chunks = append(framesSubset.Chunks, frames.Chunks[index])
	}

	return framesSubset, nil
}

// computeChunkRequestRequiredBandwidth computes the bandwidth required to fulfill a GetChunks request.
func computeChunkRequestRequiredBandwidth(request *pb.GetChunksRequest, mMap metadataMap) (uint32, error) {
	requiredBandwidth := uint32(0)
	for _, req := range request.GetChunkRequests() {
		var metadata *blobMetadata
		var key v2.BlobKey
		var requestedChunks uint32

		if req.GetByIndex() != nil {
			key = v2.BlobKey(req.GetByIndex().GetBlobKey())
			metadata = mMap[key]
			requestedChunks = uint32(len(req.GetByIndex().GetChunkIndices()))
		} else {
			key = v2.BlobKey(req.GetByRange().GetBlobKey())
			metadata = mMap[key]

			if req.GetByRange().GetEndIndex() < req.GetByRange().GetStartIndex() {
				return 0, fmt.Errorf(
					"chunk range %d-%d is invalid for key %s, start index must be less than or equal to end index",
					req.GetByRange().GetStartIndex(), req.GetByRange().GetEndIndex(), key.Hex())
			}

			requestedChunks = req.GetByRange().GetEndIndex() - req.GetByRange().GetStartIndex()
		}

		if metadata == nil {
			return 0, fmt.Errorf("metadata not found for key %s", key.Hex())
		}

		requiredBandwidth += requestedChunks * metadata.chunkSizeBytes
	}

	return requiredBandwidth, nil
}

// buildInsufficientBandwidthError builds an informative error message for when there is insufficient
// bandwidth to serve a GetChunks() request.
func buildInsufficientGetChunksBandwidthError(
	request *pb.GetChunksRequest,
	requiredBandwidth uint32,
	originalError error) *status.Status {

	chunkCount := 0
	for _, chunkRequest := range request.GetChunkRequests() {
		if chunkRequest.GetByIndex() != nil {
			chunkCount += len(chunkRequest.GetByIndex().GetChunkIndices())
		} else {
			chunkCount += int(chunkRequest.GetByRange().GetEndIndex() - chunkRequest.GetByRange().GetStartIndex())
		}
	}

	blobCount := len(request.GetChunkRequests())

	return status.Newf(codes.ResourceExhausted, "unable to serve data (%d blobs, %d chunks, %d bytes): %v",
		blobCount, chunkCount, requiredBandwidth, originalError)
}
