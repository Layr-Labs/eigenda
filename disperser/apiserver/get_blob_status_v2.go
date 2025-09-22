package apiserver

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispcommon "github.com/Layr-Labs/eigenda/disperser/common"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	blobstore "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *DispersalServerV2) GetBlobStatus(ctx context.Context, req *pb.BlobStatusRequest) (*pb.BlobStatusReply, error) {
	reply, st := s.getBlobStatus(ctx, req)
	api.LogResponseStatus(s.logger, st)
	if st != nil {
		// nolint:wrapcheck
		return reply, st.Err()
	}
	return reply, nil
}

func (s *DispersalServerV2) getBlobStatus(
	ctx context.Context,
	req *pb.BlobStatusRequest,
) (*pb.BlobStatusReply, *status.Status) {
	start := time.Now()
	defer func() {
		s.metrics.reportGetBlobStatusLatency(time.Since(start))
	}()

	if req.GetBlobKey() == nil || len(req.GetBlobKey()) != 32 {
		return nil, status.New(
			codes.InvalidArgument,
			fmt.Sprintf("blob key must be 32 bytes, got %d bytes", len(req.GetBlobKey())),
		)
	}

	blobKey, err := corev2.BytesToBlobKey(req.GetBlobKey())
	if err != nil {
		return nil, status.Newf(codes.InvalidArgument, "invalid blob key: %s", req.GetBlobKey())
	}

	metadata, err := s.blobMetadataStore.GetBlobMetadata(ctx, blobKey)
	if err != nil {
		if strings.Contains(err.Error(), "metadata not found") {
			s.logger.Info("blob metadata not found", "err", err, "blobKey", blobKey.Hex())
			return nil, status.New(codes.NotFound, "no such blob found")
		}
		s.logger.Warn("failed to get blob metadata", "err", err, "blobKey", blobKey.Hex())
		return nil, status.Newf(codes.Internal, "failed to get blob metadata: %v", err)
	}

	// If the blob is not complete or gathering signatures, return the status without the signed batch
	if metadata.BlobStatus != dispv2.Complete && metadata.BlobStatus != dispv2.GatheringSignatures {
		return &pb.BlobStatusReply{
			Status: metadata.BlobStatus.ToProfobuf(),
		}, status.New(codes.OK, "")
	}

	cert, _, err := s.blobMetadataStore.GetBlobCertificate(ctx, blobKey)
	if err != nil {
		if errors.Is(err, dispcommon.ErrMetadataNotFound) {
			return nil, status.New(codes.NotFound, "no such blob certificate found")
		}
		return nil, status.Newf(codes.Internal, "failed to get blob certificate: %v", err)
	}

	// For blobs in GatheringSignatures/Complete status, include signed batch and blob inclusion info
	blobInclusionInfos, err := s.blobMetadataStore.GetBlobInclusionInfos(ctx, blobKey)
	if err != nil {
		return nil, status.Newf(codes.Internal, "failed to get blob inclusion info for blob %s: %v", blobKey.Hex(), err)
	}

	if len(blobInclusionInfos) == 0 {
		return nil, status.Newf(codes.Internal, "no blob inclusion info found for blob %s", blobKey.Hex())
	}

	if len(blobInclusionInfos) > 1 {
		s.logger.Warn("multiple inclusion info found for blob", "blobKey", blobKey.Hex())
	}

	for _, inclusionInfo := range blobInclusionInfos {
		// get the signed batch from this inclusion info
		batchHeaderHash, err := inclusionInfo.BatchHeader.Hash()
		if err != nil {
			s.logger.Error(
				"failed to get batch header hash from blob inclusion info",
				"err",
				err,
				"blobKey",
				blobKey.Hex(),
			)
			continue
		}
		batchHeader, attestation, err := s.blobMetadataStore.GetSignedBatch(ctx, batchHeaderHash)
		if err != nil {
			if errors.Is(err, blobstore.ErrAttestationNotFound) {
				// attestation may not exist yet if the blob has not been processed by the dispatcher
				s.logger.Info("attestation not found for signed batch", "err", err, "blobKey", blobKey.Hex())
				continue
			}
			s.logger.Error("failed to get signed batch", "err", err, "blobKey", blobKey.Hex())
			continue
		}

		blobInclusionInfoProto, err := inclusionInfo.ToProtobuf(cert)
		if err != nil {
			s.logger.Error("failed to convert blob inclusion info to protobuf", "err", err, "blobKey", blobKey.Hex())
			continue
		}

		attestationProto, err := attestation.ToProtobuf()
		if err != nil {
			s.logger.Error("failed to convert attestation to protobuf", "err", err, "blobKey", blobKey.Hex())
			continue
		}

		// return the first signed batch found
		return &pb.BlobStatusReply{
			Status: metadata.BlobStatus.ToProfobuf(),
			SignedBatch: &pb.SignedBatch{
				Header:      batchHeader.ToProtobuf(),
				Attestation: attestationProto,
			},
			BlobInclusionInfo: blobInclusionInfoProto,
		}, status.New(codes.OK, "")
	}

	return nil, status.Newf(codes.Internal, "no signed batch found for blob %s", blobKey.Hex())
}
