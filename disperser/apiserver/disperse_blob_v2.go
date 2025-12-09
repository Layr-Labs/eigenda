package apiserver

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/api/grpc/controller"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	blobstore "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *DispersalServerV2) DisperseBlob(
	ctx context.Context,
	req *pb.DisperseBlobRequest,
) (*pb.DisperseBlobReply, error) {
	reply, st := s.disperseBlob(ctx, req)
	api.LogResponseStatus(s.logger, st)
	if st != nil {
		// nolint:wrapcheck
		return reply, st.Err()
	}
	return reply, nil
}

func (s *DispersalServerV2) disperseBlob(
	ctx context.Context,
	req *pb.DisperseBlobRequest,
) (*pb.DisperseBlobReply, *status.Status) {
	start := time.Now()
	defer func() {
		s.metrics.reportDisperseBlobLatency(time.Since(start))
	}()

	// Validate the request
	onchainState := s.onchainState.Load()
	if onchainState == nil {
		return nil, status.New(codes.Internal, "onchain state is not available")
	}
	blobHeader, err := s.validateDispersalRequest(req, onchainState)
	if err != nil {
		return nil, status.Newf(codes.InvalidArgument, "failed to validate request: %s", err.Error())
	}

	if st := s.checkBlobExistence(ctx, blobHeader); st != nil && st.Code() != codes.OK {
		return nil, st
	}

	if s.useControllerMediatedPayments {
		// Use the new controller-based payment system
		authorizePaymentRequest := &controller.AuthorizePaymentRequest{
			BlobHeader:      req.GetBlobHeader(),
			ClientSignature: req.GetSignature(),
		}
		_, err := s.controllerClient.AuthorizePayment(ctx, authorizePaymentRequest)
		if err != nil {
			return nil, status.Convert(err)
		}
	} else {
		// Use the legacy payment metering system
		// Check against payment meter to make sure there is quota remaining
		if st := s.checkPaymentMeter(ctx, req, start); st != nil && st.Code() != codes.OK {
			return nil, st
		}
	}

	finishedValidation := time.Now()
	s.metrics.reportValidateDispersalRequestLatency(finishedValidation.Sub(start))

	blob := req.GetBlob()
	s.metrics.reportDisperseBlobSize(len(blob))
	s.logger.Debug(
		"received a new blob dispersal request",
		"blobSizeBytes",
		len(blob),
		"quorums",
		req.GetBlobHeader().GetQuorumNumbers(),
	)

	blobKey, st := s.StoreBlob(
		ctx, blob, blobHeader, req.GetSignature(), req.GetAnchorSignature(), time.Now(), onchainState.TTL)
	if st != nil && st.Code() != codes.OK {
		return nil, st
	}
	s.logger.Debug("stored blob", "blobKey", blobKey.Hex())

	// Update Account asynchronously after successful blob storage
	go func() {
		accountID := blobHeader.PaymentMetadata.AccountID
		timestamp := uint64(time.Now().Unix())

		// Use a timeout context for the async database operation
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.blobMetadataStore.UpdateAccount(ctx, accountID, timestamp); err != nil {
			s.logger.Warn("failed to update account", "accountID", accountID.Hex(), "error", err)
		}
	}()

	s.metrics.reportStoreBlobLatency(time.Since(finishedValidation))

	return &pb.DisperseBlobReply{
		Result:  dispv2.Queued.ToProfobuf(),
		BlobKey: blobKey[:],
	}, status.New(codes.OK, "blob dispersal request accepted")
}

func (s *DispersalServerV2) StoreBlob(
	ctx context.Context,
	data []byte,
	blobHeader *corev2.BlobHeader,
	signature []byte,
	anchorSignature []byte,
	requestedAt time.Time,
	ttl time.Duration,
) (corev2.BlobKey, *status.Status) {
	blobKey, err := blobHeader.BlobKey()
	if err != nil {
		return corev2.BlobKey{}, status.Newf(codes.InvalidArgument, "failed to get blob key: %v", err)
	}

	if err := s.blobStore.StoreBlob(ctx, blobKey, data); err != nil {
		s.logger.Warn("failed to store blob", "err", err, "blobKey", blobKey.Hex())
		if errors.Is(err, blobstore.ErrAlreadyExists) {
			return corev2.BlobKey{}, status.Newf(codes.AlreadyExists, "blob already exists: %s", blobKey.Hex())
		}

		return corev2.BlobKey{}, status.Newf(codes.Internal, "failed to store blob: %v", err)
	}

	s.logger.Debug("storing blob metadata", "blobHeader", blobHeader)
	blobMetadata := &dispv2.BlobMetadata{
		BlobHeader:      blobHeader,
		Signature:       signature,
		AnchorSignature: anchorSignature,
		BlobStatus:      dispv2.Queued,
		Expiry:          uint64(requestedAt.Add(ttl).Unix()),
		NumRetries:      0,
		BlobSize:        uint64(len(data)),
		RequestedAt:     uint64(requestedAt.UnixNano()),
		UpdatedAt:       uint64(requestedAt.UnixNano()),
	}
	err = s.blobMetadataStore.PutBlobMetadata(ctx, blobMetadata)
	if err != nil {
		s.logger.Warn("failed to store blob metadata", "err", err, "blobKey", blobKey.Hex())
		if errors.Is(err, blobstore.ErrAlreadyExists) {
			return corev2.BlobKey{}, status.Newf(codes.AlreadyExists, "blob metadata already exists: %s", blobKey.Hex())
		}

		return corev2.BlobKey{}, status.Newf(codes.Internal, "failed to store blob metadata: %v", err)
	}
	return blobKey, status.New(codes.OK, "blob stored successfully")
}

func (s *DispersalServerV2) checkPaymentMeter(
	ctx context.Context,
	req *pb.DisperseBlobRequest,
	receivedAt time.Time,
) *status.Status {
	blobHeaderProto := req.GetBlobHeader()
	blobHeader, err := corev2.BlobHeaderFromProtobuf(blobHeaderProto)
	if err != nil {
		return status.Newf(codes.InvalidArgument, "invalid blob header: %s", err.Error())
	}
	blobLength := encoding.GetBlobLengthPowerOf2(uint32(len(req.GetBlob())))

	// handle payments and check rate limits
	timestamp := blobHeaderProto.GetPaymentHeader().GetTimestamp()
	cumulativePayment := new(big.Int).SetBytes(blobHeaderProto.GetPaymentHeader().GetCumulativePayment())
	accountID := blobHeaderProto.GetPaymentHeader().GetAccountId()
	if !gethcommon.IsHexAddress(accountID) {
		return status.Newf(codes.InvalidArgument, "invalid account ID: %s", accountID)
	}

	paymentHeader := core.PaymentMetadata{
		AccountID:         gethcommon.HexToAddress(accountID),
		Timestamp:         timestamp,
		CumulativePayment: cumulativePayment,
	}

	symbolsCharged, err := s.meterer.MeterRequest(
		ctx,
		paymentHeader,
		uint64(blobLength),
		blobHeader.QuorumNumbers,
		receivedAt,
	)
	if err != nil {
		return status.New(codes.ResourceExhausted, err.Error())
	}
	s.metrics.reportDisperseMeteredBytes(int(symbolsCharged) * encoding.BYTES_PER_SYMBOL)

	return status.New(codes.OK, "payment meter check passed")
}

func (s *DispersalServerV2) validateDispersalRequest(
	req *pb.DisperseBlobRequest,
	onchainState *OnchainState) (*corev2.BlobHeader, error) {

	signature := req.GetSignature()
	if len(signature) != 65 {
		return nil, fmt.Errorf("signature is expected to be 65 bytes, but got %d bytes", len(signature))
	}

	blob := req.GetBlob()
	blobSize := uint32(len(blob))
	if blobSize == 0 {
		return nil, errors.New("blob size must be greater than 0")
	}
	blobLength := encoding.GetBlobLengthPowerOf2(blobSize)
	if blobLength > s.maxNumSymbolsPerBlob {
		return nil, errors.New("blob size too big")
	}

	blobHeaderProto := req.GetBlobHeader()
	if blobHeaderProto.GetCommitment() == nil {
		return nil, errors.New("blob header must contain commitments")
	}

	if blobHeaderProto.GetCommitment() == nil {
		return nil, errors.New("blob header must contain a commitment")
	}
	commitedBlobLength := blobHeaderProto.GetCommitment().GetLength()
	if commitedBlobLength == 0 || commitedBlobLength != math.NextPowOf2u32(commitedBlobLength) {
		return nil, errors.New("invalid commitment length, must be a power of 2")
	}
	lengthPowerOf2 := encoding.GetBlobLengthPowerOf2(blobSize)
	if lengthPowerOf2 > commitedBlobLength {
		return nil, fmt.Errorf("commitment length %d is less than blob length %d", commitedBlobLength, lengthPowerOf2)
	}

	blobHeader, err := corev2.BlobHeaderFromProtobuf(blobHeaderProto)
	if err != nil {
		return nil, fmt.Errorf("invalid blob header: %w", err)
	}

	if blobHeader.PaymentMetadata == (core.PaymentMetadata{}) {
		return nil, errors.New("payment metadata is required")
	}

	if s.ReservedOnly && blobHeader.PaymentMetadata.CumulativePayment.Sign() != 0 {
		return nil, errors.New("on-demand payments are not supported by reserved-only mode disperser")
	}

	timestampIsNegative := blobHeader.PaymentMetadata.Timestamp < 0
	paymentIsNegative := blobHeader.PaymentMetadata.CumulativePayment.Cmp(big.NewInt(0)) == -1
	timestampIsZeroAndPaymentIsZero := blobHeader.PaymentMetadata.Timestamp == 0 &&
		blobHeader.PaymentMetadata.CumulativePayment.Cmp(big.NewInt(0)) == 0
	if timestampIsNegative || paymentIsNegative || timestampIsZeroAndPaymentIsZero {
		return nil, errors.New("invalid payment metadata")
	}

	if err := s.validateDispersalTimestamp(blobHeader); err != nil {
		return nil, err
	}

	if len(blobHeaderProto.GetQuorumNumbers()) == 0 {
		return nil, errors.New("blob header must contain at least one quorum number")
	}

	if len(blobHeaderProto.GetQuorumNumbers()) > int(onchainState.QuorumCount) {
		return nil, fmt.Errorf("too many quorum numbers specified: maximum is %d", onchainState.QuorumCount)
	}

	for _, quorum := range blobHeaderProto.GetQuorumNumbers() {
		if quorum > corev2.MaxQuorumID || uint8(quorum) >= onchainState.QuorumCount {
			return nil, fmt.Errorf("invalid quorum number %d; maximum is %d", quorum, onchainState.QuorumCount)
		}
	}

	// validate every 32 bytes is a valid field element
	_, err = rs.ToFrArray(blob)
	if err != nil {
		s.logger.Error("failed to convert a 32bytes as a field element", "err", err)
		return nil, errors.New(
			"encountered an error to convert a 32-bytes into a valid field element, please use the correct format where every 32bytes(big-endian) is less than 21888242871839275222246405745257275088548364400416034343698204186575808495617",
		)
	}

	if _, ok := onchainState.BlobVersionParameters.Get(corev2.BlobVersion(blobHeaderProto.GetVersion())); !ok {
		return nil, fmt.Errorf(
			"invalid blob version %d; valid blob versions are: %v",
			blobHeaderProto.GetVersion(),
			onchainState.BlobVersionParameters.Keys(),
		)
	}

	if err = s.validateAnchorSignature(req, blobHeader); err != nil {
		return nil, fmt.Errorf("validate anchor signature: %w", err)
	}

	commitments, err := s.committer.GetCommitmentsForPaddedLength(blob)
	if err != nil {
		return nil, fmt.Errorf("failed to get commitments: %w", err)
	}
	// TODO(samlaf): should differentiate 400 from 500 errors here
	if err = commitments.Equal(&blobHeader.BlobCommitments); err != nil {
		return nil, fmt.Errorf("invalid blob commitment: %w", err)
	}

	return blobHeader, nil
}

// Validates the anchor signature included in the DisperseBlobRequest.
//
// If TolerateMissingAnchorSignature is true, then this method will pass validation even if no anchor signature is
// provided in the request.
//
// If an anchor signature is provided, it will be validated whether or not TolerateMissingAnchorSignature is true.
// While validating the anchor signature, this method will also verify that the disperser ID and chain ID in the request
// match the expected values.
func (s *DispersalServerV2) validateAnchorSignature(
	req *pb.DisperseBlobRequest,
	blobHeader *corev2.BlobHeader,
) error {
	anchorSignature := req.GetAnchorSignature()

	if len(anchorSignature) == 0 {
		if s.serverConfig.TolerateMissingAnchorSignature {
			return nil
		}

		return errors.New("anchor signature is required but not provided")
	}

	if len(anchorSignature) != 65 {
		return fmt.Errorf("anchor signature length is unexpected: %d", len(anchorSignature))
	}

	if req.GetDisperserId() != s.serverConfig.DisperserId {
		return fmt.Errorf(
			"disperser ID mismatch: request specifies %d but this disperser is %d",
			req.GetDisperserId(),
			s.serverConfig.DisperserId,
		)
	}

	reqChainId, err := common.ChainIdFromBytes(req.GetChainId())
	if err != nil {
		return fmt.Errorf("invalid chain ID: %w", err)
	}
	if s.serverConfig.ChainId.Cmp(reqChainId) != 0 {
		return fmt.Errorf(
			"chain ID mismatch: request specifies %s but this disperser is on chain %s",
			reqChainId.String(),
			s.serverConfig.ChainId.String(),
		)
	}

	blobKey, err := blobHeader.BlobKey()
	if err != nil {
		return fmt.Errorf("compute blob key: %w", err)
	}

	anchorHash, err := hashing.ComputeDispersalAnchorHash(reqChainId, req.GetDisperserId(), blobKey)
	if err != nil {
		return fmt.Errorf("compute anchor hash: %w", err)
	}

	anchorSigPubKey, err := crypto.SigToPub(anchorHash, anchorSignature)
	if err != nil {
		return fmt.Errorf("recover public key from anchor signature: %w", err)
	}

	if blobHeader.PaymentMetadata.AccountID.Cmp(crypto.PubkeyToAddress(*anchorSigPubKey)) != 0 {
		return errors.New("anchor signature doesn't match account ID")
	}

	return nil
}

// Validates that the dispersal timestamp in the blob header is neither too old, nor too far in the future.
func (s *DispersalServerV2) validateDispersalTimestamp(blobHeader *corev2.BlobHeader) error {
	dispersalTime := time.Unix(0, blobHeader.PaymentMetadata.Timestamp)
	dispersalAge := s.getNow().Sub(dispersalTime)
	driftSeconds := dispersalAge.Seconds()
	accountID := blobHeader.PaymentMetadata.AccountID.Hex()

	if dispersalAge > s.MaxDispersalAge {
		s.metrics.reportDispersalTimestampRejected("stale")
		s.metrics.reportDispersalTimestampDrift(driftSeconds, "rejected", accountID)
		return fmt.Errorf("potential clock drift detected: dispersal timestamp is too old. "+
			"age=%v, max_age=%v, timestamp_unix_nanos=%d, timestamp_utc=%s",
			dispersalAge,
			s.MaxDispersalAge,
			blobHeader.PaymentMetadata.Timestamp,
			dispersalTime.UTC().Format(time.RFC3339),
		)
	}

	// If dispersalAge is negative, the timestamp is in the future
	if dispersalAge < -s.MaxFutureDispersalTime {
		s.metrics.reportDispersalTimestampRejected("future")
		s.metrics.reportDispersalTimestampDrift(driftSeconds, "rejected", accountID)
		return fmt.Errorf("potential clock drift detected: dispersal timestamp is too far in the future. "+
			"future_offset=%v, max_future_offset=%v, timestamp_unix_nanos=%d, timestamp_utc=%s",
			-dispersalAge,
			s.MaxFutureDispersalTime,
			blobHeader.PaymentMetadata.Timestamp,
			dispersalTime.UTC().Format(time.RFC3339))
	}

	// Record accepted timestamp drift
	s.metrics.reportDispersalTimestampDrift(driftSeconds, "accepted", accountID)

	return nil
}

func (s *DispersalServerV2) checkBlobExistence(ctx context.Context, blobHeader *corev2.BlobHeader) *status.Status {
	blobKey, err := blobHeader.BlobKey()
	if err != nil {
		return status.Newf(codes.InvalidArgument, "failed to parse blob key: %v", err.Error())
	}

	// check if blob already exists
	exists, err := s.blobMetadataStore.CheckBlobExists(ctx, blobKey)
	if err != nil {
		return status.Newf(codes.Internal, "failed to check blob existence: %s", err.Error())
	}

	if exists {
		return status.Newf(codes.AlreadyExists, "blob already exists: %s", blobKey.Hex())
	}

	return status.New(codes.OK, "")
}
