package apiserver

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

func (s *DispersalServerV2) DisperseBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*pb.DisperseBlobReply, error) {
	start := time.Now()
	defer func() {
		s.metrics.reportDisperseBlobLatency(time.Since(start))
	}()

	// Validate the request
	onchainState := s.onchainState.Load()
	if onchainState == nil {
		return nil, api.NewErrorInternal("onchain state is nil")
	}
	if err := s.validateDispersalRequest(req, onchainState); err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to validate the request: %v", err))
	}

	// Check against payment meter to make sure there is quota remaining
	if err := s.checkPaymentMeter(ctx, req); err != nil {
		return nil, err
	}

	finishedValidation := time.Now()
	s.metrics.reportValidateDispersalRequestLatency(finishedValidation.Sub(start))

	blob := req.GetBlob()
	s.metrics.reportDisperseBlobSize(len(blob))
	blobHeader, err := corev2.BlobHeaderFromProtobuf(req.GetBlobHeader())
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to parse the blob header proto: %v", err))
	}
	s.logger.Debug("received a new blob dispersal request", "blobSizeBytes", len(blob), "quorums", req.GetBlobHeader().GetQuorumNumbers())

	blobKey, err := s.StoreBlob(ctx, blob, blobHeader, req.GetSignature(), time.Now(), onchainState.TTL)
	if err != nil {
		return nil, err
	}

	s.metrics.reportStoreBlobLatency(time.Since(finishedValidation))

	return &pb.DisperseBlobReply{
		Result:  dispv2.Queued.ToProfobuf(),
		BlobKey: blobKey[:],
	}, nil
}

func (s *DispersalServerV2) StoreBlob(ctx context.Context, data []byte, blobHeader *corev2.BlobHeader, signature []byte, requestedAt time.Time, ttl time.Duration) (corev2.BlobKey, error) {
	blobKey, err := blobHeader.BlobKey()
	if err != nil {
		return corev2.BlobKey{}, api.NewErrorInvalidArg(fmt.Sprintf("failed to get blob key: %v", err))
	}

	if err := s.blobStore.StoreBlob(ctx, blobKey, data); err != nil {
		s.logger.Warn("failed to store blob", "err", err, "blobKey", blobKey.Hex())
		if errors.Is(err, common.ErrAlreadyExists) {
			return corev2.BlobKey{}, api.NewErrorAlreadyExists(fmt.Sprintf("blob already exists: %s", blobKey.Hex()))
		}

		return corev2.BlobKey{}, api.NewErrorInternal(fmt.Sprintf("failed to store blob: %v", err))
	}

	blobMetadata := &dispv2.BlobMetadata{
		BlobHeader:  blobHeader,
		Signature:   signature,
		BlobStatus:  dispv2.Queued,
		Expiry:      uint64(requestedAt.Add(ttl).Unix()),
		NumRetries:  0,
		BlobSize:    uint64(len(data)),
		RequestedAt: uint64(requestedAt.UnixNano()),
		UpdatedAt:   uint64(requestedAt.UnixNano()),
	}
	err = s.blobMetadataStore.PutBlobMetadata(ctx, blobMetadata)
	if err != nil {
		s.logger.Warn("failed to store blob metadata", "err", err, "blobKey", blobKey.Hex())
		if errors.Is(err, common.ErrAlreadyExists) {
			return corev2.BlobKey{}, api.NewErrorAlreadyExists(fmt.Sprintf("blob metadata already exists: %s", blobKey.Hex()))
		}

		return corev2.BlobKey{}, api.NewErrorInternal(fmt.Sprintf("failed to store blob metadata: %v", err))
	}
	return blobKey, err
}

func (s *DispersalServerV2) checkPaymentMeter(ctx context.Context, req *pb.DisperseBlobRequest) error {
	blobHeaderProto := req.GetBlobHeader()
	blobHeader, err := corev2.BlobHeaderFromProtobuf(blobHeaderProto)
	if err != nil {
		return api.NewErrorInvalidArg(fmt.Sprintf("invalid blob header: %s", err.Error()))
	}
	blobLength := encoding.GetBlobLengthPowerOf2(uint(len(req.GetBlob())))

	// handle payments and check rate limits
	reservationPeriod := blobHeaderProto.GetPaymentHeader().GetReservationPeriod()
	cumulativePayment := new(big.Int).SetBytes(blobHeaderProto.GetPaymentHeader().GetCumulativePayment())
	accountID := blobHeaderProto.GetPaymentHeader().GetAccountId()

	paymentHeader := core.PaymentMetadata{
		AccountID:         accountID,
		ReservationPeriod: reservationPeriod,
		CumulativePayment: cumulativePayment,
	}

	symbolsCharged, err := s.meterer.MeterRequest(ctx, paymentHeader, blobLength, blobHeader.QuorumNumbers)
	if err != nil {
		return api.NewErrorResourceExhausted(err.Error())
	}
	s.metrics.reportDisperseMeteredBytes(int(symbolsCharged) * encoding.BYTES_PER_SYMBOL)

	return nil
}

func (s *DispersalServerV2) validateDispersalRequest(
	req *pb.DisperseBlobRequest,
	onchainState *OnchainState) error {

	signature := req.GetSignature()
	if len(signature) != 65 {
		return fmt.Errorf("signature is expected to be 65 bytes, but got %d bytes", len(signature))
	}
	blob := req.GetBlob()
	blobSize := len(blob)
	if blobSize == 0 {
		return errors.New("blob size must be greater than 0")
	}
	blobLength := encoding.GetBlobLengthPowerOf2(uint(blobSize))
	if blobLength > uint(s.maxNumSymbolsPerBlob) {
		return errors.New("blob size too big")
	}

	blobHeaderProto := req.GetBlobHeader()
	if blobHeaderProto.GetCommitment() == nil {
		return errors.New("blob header must contain commitments")
	}

	if blobHeaderProto.GetCommitment() == nil {
		return errors.New("blob header must contain a commitment")
	}
	commitedBlobLength := blobHeaderProto.GetCommitment().GetLength()
	if commitedBlobLength == 0 || commitedBlobLength != encoding.NextPowerOf2(commitedBlobLength) {
		return errors.New("invalid commitment length, must be a power of 2")
	}
	lengthPowerOf2 := encoding.GetBlobLengthPowerOf2(uint(blobSize))
	if lengthPowerOf2 > uint(commitedBlobLength) {
		return fmt.Errorf("commitment length %d is less than blob length %d", commitedBlobLength, lengthPowerOf2)
	}

	blobHeader, err := corev2.BlobHeaderFromProtobuf(blobHeaderProto)
	if err != nil {
		return fmt.Errorf("invalid blob header: %w", err)
	}

	if blobHeader.PaymentMetadata == (core.PaymentMetadata{}) {
		return errors.New("payment metadata is required")
	}

	if len(blobHeader.PaymentMetadata.AccountID) == 0 || (blobHeader.PaymentMetadata.ReservationPeriod == 0 && blobHeader.PaymentMetadata.CumulativePayment.Cmp(big.NewInt(0)) == 0) {
		return errors.New("invalid payment metadata")
	}

	if len(blobHeaderProto.GetQuorumNumbers()) == 0 {
		return errors.New("blob header must contain at least one quorum number")
	}

	if len(blobHeaderProto.GetQuorumNumbers()) > int(onchainState.QuorumCount) {
		return fmt.Errorf("too many quorum numbers specified: maximum is %d", onchainState.QuorumCount)
	}

	for _, quorum := range blobHeaderProto.GetQuorumNumbers() {
		if quorum > corev2.MaxQuorumID || uint8(quorum) >= onchainState.QuorumCount {
			return fmt.Errorf("invalid quorum number %d; maximum is %d", quorum, onchainState.QuorumCount)
		}
	}

	// validate every 32 bytes is a valid field element
	_, err = rs.ToFrArray(blob)
	if err != nil {
		s.logger.Error("failed to convert a 32bytes as a field element", "err", err)
		return errors.New("encountered an error to convert a 32-bytes into a valid field element, please use the correct format where every 32bytes(big-endian) is less than 21888242871839275222246405745257275088548364400416034343698204186575808495617")
	}

	if _, ok := onchainState.BlobVersionParameters.Get(corev2.BlobVersion(blobHeaderProto.GetVersion())); !ok {
		return fmt.Errorf("invalid blob version %d; valid blob versions are: %v", blobHeaderProto.GetVersion(), onchainState.BlobVersionParameters.Keys())
	}

	if err = s.authenticator.AuthenticateBlobRequest(blobHeader, signature); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	commitments, err := s.prover.GetCommitmentsForPaddedLength(blob)
	if err != nil {
		return fmt.Errorf("failed to get commitments: %w", err)
	}
	if !commitments.Equal(&blobHeader.BlobCommitments) {
		return errors.New("invalid blob commitment")
	}

	return nil
}
