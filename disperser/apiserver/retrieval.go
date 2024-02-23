package apiserver

import (
	"context"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *DispersalServer) RetrieveBlob(ctx context.Context, req *pb.RetrieveBlobRequest) (*pb.RetrieveBlobReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("RetrieveBlob", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	s.logger.Info("received a new blob retrieval request", "batchHeaderHash", req.BatchHeaderHash, "blobIndex", req.BlobIndex)

	batchHeaderHash := req.GetBatchHeaderHash()
	// Convert to [32]byte
	var batchHeaderHash32 [32]byte
	copy(batchHeaderHash32[:], batchHeaderHash)

	blobIndex := req.GetBlobIndex()

	blobMetadata, err := s.blobStore.GetMetadataInBatch(ctx, batchHeaderHash32, blobIndex)
	if err != nil {
		s.logger.Error("Failed to retrieve blob metadata", "err", err)
		s.metrics.IncrementFailedBlobRequestNum("", "RetrieveBlob")

		return nil, err
	}

	data, err := s.blobStore.GetBlobContent(ctx, blobMetadata.BlobHash)
	if err != nil {
		s.logger.Error("Failed to retrieve blob", "err", err)
		s.metrics.HandleFailedRequest("", len(data), "RetrieveBlob")

		return nil, err
	}

	s.metrics.HandleSuccessfulRequest("", len(data), "RetrieveBlob")

	return &pb.RetrieveBlobReply{
		Data: data,
	}, nil
}
