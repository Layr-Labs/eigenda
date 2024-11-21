package apiserver

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

func (s *DispersalServerV2) DisperseBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*pb.DisperseBlobReply, error) {
	if err := s.validateDispersalRequest(req); err != nil {
		return nil, err
	}

	data := req.GetData()
	blobHeader, err := corev2.BlobHeaderFromProtobuf(req.GetBlobHeader())
	if err != nil {
		return nil, api.NewErrorInternal(err.Error())
	}
	s.logger.Debug("received a new blob dispersal request", "blobSizeBytes", len(data), "quorums", req.GetBlobHeader().GetQuorumNumbers())

	// TODO(ian-shim): handle payments and check rate limits

	blobKey, err := s.StoreBlob(ctx, data, blobHeader, time.Now())
	if err != nil {
		return nil, api.NewErrorInternal(err.Error())
	}

	return &pb.DisperseBlobReply{
		Result:  dispv2.Queued.ToProfobuf(),
		BlobKey: blobKey[:],
	}, nil
}

func (s *DispersalServerV2) StoreBlob(ctx context.Context, data []byte, blobHeader *corev2.BlobHeader, requestedAt time.Time) (corev2.BlobKey, error) {
	blobKey, err := blobHeader.BlobKey()
	if err != nil {
		return corev2.BlobKey{}, err
	}

	if err := s.blobStore.StoreBlob(ctx, blobKey, data); err != nil {
		return corev2.BlobKey{}, err
	}

	blobMetadata := &dispv2.BlobMetadata{
		BlobHeader:  blobHeader,
		BlobStatus:  dispv2.Queued,
		Expiry:      uint64(requestedAt.Add(s.onchainState.TTL).Unix()),
		NumRetries:  0,
		BlobSize:    uint64(len(data)),
		RequestedAt: uint64(requestedAt.UnixNano()),
		UpdatedAt:   uint64(requestedAt.UnixNano()),
	}
	err = s.blobMetadataStore.PutBlobMetadata(ctx, blobMetadata)
	return blobKey, err
}

func (s *DispersalServerV2) validateDispersalRequest(req *pb.DisperseBlobRequest) error {
	data := req.GetData()
	blobSize := len(data)
	if blobSize == 0 {
		return api.NewErrorInvalidArg("blob size must be greater than 0")
	}
	blobLength := encoding.GetBlobLengthPowerOf2(uint(blobSize))
	if blobLength > uint(s.maxNumSymbolsPerBlob) {
		return api.NewErrorInvalidArg("blob size too big")
	}

	blobHeaderProto := req.GetBlobHeader()
	if blobHeaderProto.GetCommitment() == nil {
		return api.NewErrorInvalidArg("blob header must contain commitments")
	}

	if len(blobHeaderProto.GetQuorumNumbers()) == 0 {
		return api.NewErrorInvalidArg("blob header must contain at least one quorum number")
	}

	if len(blobHeaderProto.GetQuorumNumbers()) > int(s.onchainState.QuorumCount) {
		return api.NewErrorInvalidArg(fmt.Sprintf("too many quorum numbers specified: maximum is %d", s.onchainState.QuorumCount))
	}

	for _, quorum := range blobHeaderProto.GetQuorumNumbers() {
		if quorum > corev2.MaxQuorumID || uint8(quorum) >= s.onchainState.QuorumCount {
			return api.NewErrorInvalidArg(fmt.Sprintf("invalid quorum number %d; maximum is %d", quorum, s.onchainState.QuorumCount))
		}
	}

	// validate every 32 bytes is a valid field element
	_, err := rs.ToFrArray(data)
	if err != nil {
		s.logger.Error("failed to convert a 32bytes as a field element", "err", err)
		return api.NewErrorInvalidArg("encountered an error to convert a 32-bytes into a valid field element, please use the correct format where every 32bytes(big-endian) is less than 21888242871839275222246405745257275088548364400416034343698204186575808495617")
	}

	if _, ok := s.onchainState.BlobVersionParameters[corev2.BlobVersion(blobHeaderProto.GetVersion())]; !ok {
		validVersions := make([]int32, 0, len(s.onchainState.BlobVersionParameters))
		for version := range s.onchainState.BlobVersionParameters {
			validVersions = append(validVersions, int32(version))
		}
		return api.NewErrorInvalidArg(fmt.Sprintf("invalid blob version %d; valid blob versions are: %v", blobHeaderProto.GetVersion(), validVersions))
	}

	blobHeader, err := corev2.BlobHeaderFromProtobuf(blobHeaderProto)
	if err != nil {
		return api.NewErrorInvalidArg(fmt.Sprintf("invalid blob header: %s", err.Error()))
	}
	if err = s.authenticator.AuthenticateBlobRequest(blobHeader); err != nil {
		return api.NewErrorInvalidArg(fmt.Sprintf("authentication failed: %s", err.Error()))
	}

	if len(blobHeader.PaymentMetadata.AccountID) == 0 || blobHeader.PaymentMetadata.BinIndex == 0 || blobHeader.PaymentMetadata.CumulativePayment == nil {
		return api.NewErrorInvalidArg("invalid payment metadata")
	}

	commitments, err := s.prover.GetCommitmentsForPaddedLength(data)
	if err != nil {
		return api.NewErrorInternal(fmt.Sprintf("failed to get commitments: %v", err))
	}
	if !commitments.Equal(&blobHeader.BlobCommitments) {
		return api.NewErrorInvalidArg("invalid blob commitment")
	}

	return nil
}
