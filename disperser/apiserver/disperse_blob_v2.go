package apiserver

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

func (s *DispersalServerV2) DisperseBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*pb.DisperseBlobReply, error) {
	onchainState := s.onchainState.Load()
	if onchainState == nil {
		return nil, api.NewErrorInternal("onchain state is nil")
	}

	if err := s.validateDispersalRequest(ctx, req, onchainState); err != nil {
		return nil, err
	}

	data := req.GetData()
	blobHeader, err := corev2.BlobHeaderFromProtobuf(req.GetBlobHeader())
	if err != nil {
		return nil, api.NewErrorInternal(err.Error())
	}
	s.logger.Debug("received a new blob dispersal request", "blobSizeBytes", len(data), "quorums", req.GetBlobHeader().GetQuorumNumbers())

	blobKey, err := s.StoreBlob(ctx, data, blobHeader, time.Now(), onchainState.TTL)
	if err != nil {
		return nil, err
	}

	return &pb.DisperseBlobReply{
		Result:  dispv2.Queued.ToProfobuf(),
		BlobKey: blobKey[:],
	}, nil
}

func (s *DispersalServerV2) StoreBlob(ctx context.Context, data []byte, blobHeader *corev2.BlobHeader, requestedAt time.Time, ttl time.Duration) (corev2.BlobKey, error) {
	blobKey, err := blobHeader.BlobKey()
	if err != nil {
		return corev2.BlobKey{}, api.NewErrorInvalidArg(fmt.Sprintf("failed to get blob key: %v", err))
	}

	if err := s.blobStore.StoreBlob(ctx, blobKey, data); err != nil {
		return corev2.BlobKey{}, api.NewErrorInternal(fmt.Sprintf("failed to store blob: %v", err))
	}

	blobMetadata := &dispv2.BlobMetadata{
		BlobHeader:  blobHeader,
		BlobStatus:  dispv2.Queued,
		Expiry:      uint64(requestedAt.Add(ttl).Unix()),
		NumRetries:  0,
		BlobSize:    uint64(len(data)),
		RequestedAt: uint64(requestedAt.UnixNano()),
		UpdatedAt:   uint64(requestedAt.UnixNano()),
	}
	err = s.blobMetadataStore.PutBlobMetadata(ctx, blobMetadata)
	return blobKey, err
}

func (s *DispersalServerV2) validateDispersalRequest(ctx context.Context, req *pb.DisperseBlobRequest, onchainState *OnchainState) error {
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

	if len(blobHeaderProto.GetQuorumNumbers()) > int(onchainState.QuorumCount) {
		return api.NewErrorInvalidArg(fmt.Sprintf("too many quorum numbers specified: maximum is %d", onchainState.QuorumCount))
	}

	for _, quorum := range blobHeaderProto.GetQuorumNumbers() {
		if quorum > corev2.MaxQuorumID || uint8(quorum) >= onchainState.QuorumCount {
			return api.NewErrorInvalidArg(fmt.Sprintf("invalid quorum number %d; maximum is %d", quorum, onchainState.QuorumCount))
		}
	}

	// validate every 32 bytes is a valid field element
	_, err := rs.ToFrArray(data)
	if err != nil {
		s.logger.Error("failed to convert a 32bytes as a field element", "err", err)
		return api.NewErrorInvalidArg("encountered an error to convert a 32-bytes into a valid field element, please use the correct format where every 32bytes(big-endian) is less than 21888242871839275222246405745257275088548364400416034343698204186575808495617")
	}

	if _, ok := onchainState.BlobVersionParameters.Get(corev2.BlobVersion(blobHeaderProto.GetVersion())); !ok {
		return api.NewErrorInvalidArg(fmt.Sprintf("invalid blob version %d; valid blob versions are: %v", blobHeaderProto.GetVersion(), onchainState.BlobVersionParameters.Keys()))
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

	// handle payments and check rate limits
	if blobHeaderProto.GetPaymentHeader() != nil {
		binIndex := blobHeaderProto.GetPaymentHeader().GetBinIndex()
		cumulativePayment := new(big.Int).SetBytes(blobHeaderProto.GetPaymentHeader().GetCumulativePayment())
		accountID := blobHeaderProto.GetPaymentHeader().GetAccountId()

		paymentHeader := core.PaymentMetadata{
			AccountID:         accountID,
			BinIndex:          binIndex,
			CumulativePayment: cumulativePayment,
		}

		err := s.meterer.MeterRequest(ctx, paymentHeader, blobLength, blobHeader.QuorumNumbers)
		if err != nil {
			return api.NewErrorResourceExhausted(err.Error())
		}
	} else {
		return api.NewErrorInvalidArg("payment header is required")
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
