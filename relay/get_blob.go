package relay

import (
	"context"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetBlob retrieves a blob stored by the relay.
func (s *Server) GetBlob(ctx context.Context, request *pb.GetBlobRequest) (*pb.GetBlobReply, error) {
	reply, st := s.getBlob(ctx, request)
	api.LogResponseStatus(s.logger, st)
	if st != nil {
		// nolint:wrapcheck
		return reply, st.Err()
	}
	return reply, nil
}

func (s *Server) getBlob(ctx context.Context, request *pb.GetBlobRequest) (*pb.GetBlobReply, *status.Status) {
	start := time.Now()

	if s.config.Timeouts.GetBlobTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.config.Timeouts.GetBlobTimeout)
		defer cancel()
	}

	// Validate the request params before any further processing (as validation is cheaper)
	key, err := v2.BytesToBlobKey(request.GetBlobKey())
	if err != nil {
		return nil, status.Newf(codes.InvalidArgument, "invalid blob key: %v", err)
	}
	s.logger.Debug("GetBlob request received", "key", key.Hex())

	err = s.blobRateLimiter.BeginGetBlobOperation(time.Now())
	if err != nil {
		return nil, status.Newf(codes.ResourceExhausted, "rate limit exceeded: %v", err)
	}
	defer s.blobRateLimiter.FinishGetBlobOperation()

	keys := []v2.BlobKey{key}
	mMap, err := s.metadataProvider.GetMetadataForBlobs(ctx, keys)
	if err != nil {
		if strings.Contains(err.Error(), blobstore.ErrMetadataNotFound.Error()) {
			// nolint:wrapcheck
			return nil, status.Newf(codes.NotFound,
				"blob %s not found, check if blob exists and is assigned to this relay", key.Hex())
		}
		// nolint:wrapcheck
		return nil, status.Newf(codes.Internal, "error fetching metadata for blob: %v", err)
	}
	metadata := mMap[v2.BlobKey(request.GetBlobKey())]
	if metadata == nil {
		return nil, status.New(codes.NotFound, "blob not found")
	}

	finishedFetchingMetadata := time.Now()
	s.metrics.ReportBlobMetadataLatency(finishedFetchingMetadata.Sub(start))

	s.metrics.ReportBlobRequestedBandwidthUsage(int(metadata.blobSizeBytes))
	err = s.blobRateLimiter.RequestGetBlobBandwidth(time.Now(), metadata.blobSizeBytes)
	if err != nil {
		return nil, status.Newf(codes.ResourceExhausted, "bandwidth limit exceeded: %v", err)
	}

	data, err := s.blobProvider.GetBlob(ctx, key)
	if err != nil {
		if strings.Contains(err.Error(), blobstore.ErrBlobNotFound.Error()) {
			return nil, status.Newf(codes.NotFound, "blob %s not found", key.Hex())
		} else {
			s.logger.Errorf("error fetching blob %s: %v", key.Hex(), err)
			return nil, status.Newf(codes.NotFound,
				"relay encountered errors while attempting to fetch blob %s", key.Hex())
		}
	}

	s.metrics.ReportBlobBandwidthUsage(len(data))
	s.metrics.ReportBlobDataLatency(time.Since(finishedFetchingMetadata))
	s.metrics.ReportBlobLatency(time.Since(start))

	reply := &pb.GetBlobReply{
		Blob: data,
	}
	return reply, status.New(codes.OK, "")
}
