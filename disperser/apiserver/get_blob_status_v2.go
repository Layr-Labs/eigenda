package apiserver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispcommon "github.com/Layr-Labs/eigenda/disperser/common"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
)

func (s *DispersalServerV2) GetBlobStatus(ctx context.Context, req *pb.BlobStatusRequest) (*pb.BlobStatusReply, error) {
	start := time.Now()
	defer func() {
		s.metrics.reportGetBlobStatusLatency(time.Since(start))
	}()

	if req.GetBlobKey() == nil || len(req.GetBlobKey()) != 32 {
		return nil, api.NewErrorInvalidArg("invalid blob key")
	}

	blobKey, err := corev2.BytesToBlobKey(req.GetBlobKey())
	if err != nil {
		return nil, api.NewErrorInvalidArg("invalid blob key")
	}

	metadata, err := s.blobMetadataStore.GetBlobMetadata(ctx, blobKey)
	if err != nil {
		s.logger.Error("failed to get blob metadata", "err", err, "blobKey", blobKey.Hex())
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get blob metadata: %s", err.Error()))
	}

	if metadata.BlobStatus != dispv2.Certified {
		return &pb.BlobStatusReply{
			Status: metadata.BlobStatus.ToProfobuf(),
		}, nil
	}

	cert, _, err := s.blobMetadataStore.GetBlobCertificate(ctx, blobKey)
	if err != nil {
		s.logger.Error("failed to get blob certificate", "err", err, "blobKey", blobKey.Hex())
		if errors.Is(err, dispcommon.ErrMetadataNotFound) {
			return nil, api.NewErrorNotFound("no such blob certificate found")
		}
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get blob certificate: %s", err.Error()))
	}

	// For certified blobs, include signed batch and blob inclusion info
	blobInclusionInfos, err := s.blobMetadataStore.GetBlobInclusionInfos(ctx, blobKey)
	if err != nil {
		s.logger.Error("failed to get blob inclusion info", "err", err, "blobKey", blobKey.Hex())
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to get blob inclusion info: %s", err.Error()))
	}

	if len(blobInclusionInfos) == 0 {
		s.logger.Error("no inclusion info found for certified blob", "blobKey", blobKey.Hex())
		return nil, api.NewErrorInternal("no inclusion info found")
	}

	if len(blobInclusionInfos) > 1 {
		s.logger.Warn("multiple inclusion info found for certified blob", "blobKey", blobKey.Hex())
	}

	for _, inclusionInfo := range blobInclusionInfos {
		// get the signed batch from this inclusion info
		batchHeaderHash, err := inclusionInfo.BatchHeader.Hash()
		if err != nil {
			s.logger.Error("failed to get batch header hash", "err", err, "blobKey", blobKey.Hex())
			continue
		}
		batchHeader, attestation, err := s.blobMetadataStore.GetSignedBatch(ctx, batchHeaderHash)
		if err != nil {
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
		}, nil
	}

	s.logger.Error("no signed batch found for certified blob", "blobKey", blobKey.Hex())
	return nil, api.NewErrorInternal("no signed batch found")
}
