package clients

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DisperserClientConfig struct {
	Hostname          string
	Port              string
	UseSecureGrpcFlag bool
	NtpServer         string
	NtpSyncInterval   time.Duration
}

// DisperserClient manages communication with the disperser server.
type DisperserClient interface {
	// Close closes the grpc connection to the disperser server.
	Close() error
	// DisperseBlob disperses a blob with the given data, blob version, and quorums.
	DisperseBlob(
		ctx context.Context,
		data []byte,
		blobVersion corev2.BlobVersion,
		quorums []core.QuorumID) (*dispv2.BlobStatus, corev2.BlobKey, error)
	// DisperseBlobWithProbe is similar to DisperseBlob, but also takes a SequenceProbe to capture metrics.
	// If the probe is nil, no metrics are captured.
	DisperseBlobWithProbe(
		ctx context.Context,
		data []byte,
		blobVersion corev2.BlobVersion,
		quorums []core.QuorumID,
		probe *common.SequenceProbe) (*dispv2.BlobStatus, corev2.BlobKey, error)
	// GetBlobStatus returns the status of a blob with the given blob key.
	GetBlobStatus(ctx context.Context, blobKey corev2.BlobKey) (*disperser_rpc.BlobStatusReply, error)
	// GetBlobCommitment returns the blob commitment for a given blob payload.
	GetBlobCommitment(ctx context.Context, data []byte) (*disperser_rpc.BlobCommitmentReply, error)
}
type disperserClient struct {
	config             *DisperserClientConfig
	signer             corev2.BlobRequestSigner
	initOnceGrpc       sync.Once
	initOnceAccountant sync.Once
	conn               *grpc.ClientConn
	client             disperser_rpc.DisperserClient
	prover             encoding.Prover
	accountant         *Accountant
	accountantLock     sync.Mutex
	ntpClock           *core.NTPSyncedClock
}

var _ DisperserClient = &disperserClient{}

// DisperserClient maintains a single underlying grpc connection to the disperser server,
// through which it sends requests to disperse blobs and get blob status.
// The connection is established lazily on the first method call. Don't forget to call Close(),
// which is safe to call even if the connection was never established.
//
// DisperserClient is safe to be used concurrently by multiple goroutines.
//
// Example usage:
//
//	client := NewDisperserClient(config, signer)
//	defer client.Close()
//
//	// The connection will be established on the first call
//	status, blobKey, err := client.DisperseBlob(ctx, data, blobHeader)
//	if err != nil {
//	    // Handle error
//	}
//
//	// Subsequent calls will use the existing connection
//	status2, blobKey2, err := client.DisperseBlob(ctx, data, blobHeader)
func NewDisperserClient(logger logging.Logger, config *DisperserClientConfig, signer corev2.BlobRequestSigner, prover encoding.Prover, accountant *Accountant) (*disperserClient, error) {
	if config == nil {
		return nil, api.NewErrorInvalidArg("config must be provided")
	}
	if config.Hostname == "" {
		return nil, api.NewErrorInvalidArg("hostname must be provided")
	}
	if config.Port == "" {
		return nil, api.NewErrorInvalidArg("port must be provided")
	}
	if signer == nil {
		return nil, api.NewErrorInvalidArg("signer must be provided")
	}

	// Set default NTP config if not provided
	if config.NtpServer == "" {
		config.NtpServer = "pool.ntp.org"
	}
	if config.NtpSyncInterval == 0 {
		config.NtpSyncInterval = 5 * time.Minute
	}

	// Initialize NTP synced clock
	ntpClock, err := core.NewNTPSyncedClock(context.Background(), config.NtpServer, config.NtpSyncInterval, logger.With("component", "DisperserClient"))
	if err != nil {
		return nil, fmt.Errorf("failed to create NTP clock: %w", err)
	}

	return &disperserClient{
		config:     config,
		signer:     signer,
		prover:     prover,
		accountant: accountant,
		ntpClock:   ntpClock,
		// conn and client are initialized lazily
	}, nil
}

// PopulateAccountant populates the accountant with the payment state from the disperser.
func (c *disperserClient) PopulateAccountant(ctx context.Context) error {
	if c.accountant == nil {
		accountId, err := c.signer.GetAccountID()
		if err != nil {
			return fmt.Errorf("error getting account ID: %w", err)
		}
		c.accountant = NewAccountant(accountId)
	}

	paymentStateProto, err := c.GetPaymentState(ctx)
	if err != nil {
		return fmt.Errorf("error getting payment state for initializing accountant: %w", err)
	}

	// Convert protobuf types to native Go types using meterer conversion function
	paymentVaultParams, reservations, cumulativePayment, onchainCumulativePayment, periodRecords, err := meterer.ConvertPaymentStateFromProtobuf(paymentStateProto)
	if err != nil {
		return fmt.Errorf("error converting payment state from protobuf: %w", err)
	}

	err = c.accountant.SetPaymentState(paymentVaultParams, reservations, cumulativePayment, onchainCumulativePayment, periodRecords)
	if err != nil {
		return fmt.Errorf("error setting payment state for accountant: %w", err)
	}

	return nil
}

// Close closes the grpc connection to the disperser server.
// It is thread safe and can be called multiple times.
func (c *disperserClient) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		c.client = nil
		return err
	}
	return nil
}

func (c *disperserClient) DisperseBlob(
	ctx context.Context,
	data []byte,
	blobVersion corev2.BlobVersion,
	quorums []core.QuorumID,
) (*dispv2.BlobStatus, corev2.BlobKey, error) {
	return c.DisperseBlobWithProbe(ctx, data, blobVersion, quorums, nil)
}

// DisperseBlobWithProbe disperses a blob with the given data, blob version, and quorums. If sequenceProbe is not nil,
// the probe is used to capture metrics during the dispersal process.
func (c *disperserClient) DisperseBlobWithProbe(
	ctx context.Context,
	data []byte,
	blobVersion corev2.BlobVersion,
	quorums []core.QuorumID,
	probe *common.SequenceProbe,
) (*dispv2.BlobStatus, corev2.BlobKey, error) {

	if len(quorums) == 0 {
		return nil, [32]byte{}, api.NewErrorInvalidArg("quorum numbers must be provided")
	}
	if c.signer == nil {
		return nil, [32]byte{}, api.NewErrorInternal("uninitialized signer for authenticated dispersal")
	}
	for _, q := range quorums {
		if q > corev2.MaxQuorumID {
			return nil, [32]byte{}, api.NewErrorInvalidArg("quorum number must be less than 256")
		}
	}

	err := c.initOnceGrpcConnection()
	if err != nil {
		return nil, [32]byte{}, api.NewErrorFailover(err)
	}

	probe.SetStage("acquire_accountant_lock")
	c.accountantLock.Lock()

	probe.SetStage("accountant")

	err = c.initOncePopulateAccountant(ctx)
	if err != nil {
		c.accountantLock.Unlock()
		return nil, [32]byte{}, api.NewErrorFailover(err)
	}

	symbolLength := encoding.GetBlobLengthPowerOf2(uint(len(data)))
	payment, err := c.accountant.AccountBlob(c.ntpClock.Now().UnixNano(), uint64(symbolLength), quorums)
	if err != nil {
		c.accountantLock.Unlock()
		return nil, [32]byte{}, fmt.Errorf("error accounting blob: %w", err)
	}

	if payment.CumulativePayment == nil || payment.CumulativePayment.Sign() == 0 {
		// This request is using reserved bandwidth, no need to prevent parallel dispersal.
		c.accountantLock.Unlock()
	} else {
		// This request is using on-demand bandwidth, current implementation requires sequential dispersal.
		defer c.accountantLock.Unlock()
	}

	probe.SetStage("verify_field_element")

	// check every 32 bytes of data are within the valid range for a bn254 field element
	_, err = rs.ToFrArray(data)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf(
			"encountered an error to convert a 32-bytes into a valid field element, "+
				"please use the correct format where every 32bytes(big-endian) is less than "+
				"21888242871839275222246405745257275088548364400416034343698204186575808495617 %w", err)
	}

	probe.SetStage("get_commitments")

	var blobCommitments encoding.BlobCommitments
	if c.prover == nil {
		// if prover is not configured, get blob commitments from disperser
		commitments, err := c.GetBlobCommitment(ctx, data)
		if err != nil {
			return nil, [32]byte{}, fmt.Errorf("error getting blob commitments: %w", err)
		}
		deserialized, err := encoding.BlobCommitmentsFromProtobuf(commitments.GetBlobCommitment())
		if err != nil {
			return nil, [32]byte{}, fmt.Errorf("error deserializing blob commitments: %w", err)
		}
		blobCommitments = *deserialized

		// We need to check that the disperser used the correct length. Even once checking the commitment from the
		// disperser has been implemented, there is still an edge case where the disperser could truncate trailing 0s,
		// yielding the wrong blob length, but not causing commitment verification to fail. It is important that the
		// commitment doesn't report a blob length smaller than expected, since this could cause payload parsing to
		// fail, if the length claimed in the encoded payload header is larger than the blob length in the commitment.
		lengthFromCommitment := commitments.GetBlobCommitment().GetLength()
		if lengthFromCommitment != uint32(symbolLength) {
			return nil, [32]byte{}, fmt.Errorf(
				"blob commitment length (%d) from disperser doesn't match expected length (%d): %w",
				lengthFromCommitment, symbolLength, err)
		}
	} else {
		// if prover is configured, get commitments from prover

		blobCommitments, err = c.prover.GetCommitmentsForPaddedLength(data)
		if err != nil {
			return nil, [32]byte{}, fmt.Errorf("error getting blob commitments: %w", err)
		}
	}

	blobHeader := &corev2.BlobHeader{
		BlobVersion:     blobVersion,
		BlobCommitments: blobCommitments,
		QuorumNumbers:   quorums,
		PaymentMetadata: *payment,
	}

	probe.SetStage("sign_blob_request")

	sig, err := c.signer.SignBlobRequest(blobHeader)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("error signing blob request: %w", err)
	}
	blobHeaderProto, err := blobHeader.ToProtobuf()
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("error converting blob header to protobuf: %w", err)
	}
	request := &disperser_rpc.DisperseBlobRequest{
		Blob:       data,
		Signature:  sig,
		BlobHeader: blobHeaderProto,
	}

	probe.SetStage("send_to_disperser")

	reply, err := c.client.DisperseBlob(ctx, request)
	if err != nil {
		// TODO: rollback payment for the accountant if the blob fails to disperse
		// because ondemand request hits global ratelimit.
		return nil, [32]byte{}, fmt.Errorf("error while calling DisperseBlob: %w", err)
	}

	blobStatus, err := dispv2.BlobStatusFromProtobuf(reply.GetResult())
	if err != nil {
		return nil, [32]byte{}, err
	}

	probe.SetStage("verify_blob_key")

	if verifyReceivedBlobKey(blobHeader, reply) != nil {
		return nil, [32]byte{}, fmt.Errorf("verify received blob key: %w", err)
	}

	return &blobStatus, corev2.BlobKey(reply.GetBlobKey()), nil
}

// verifyReceivedBlobKey computes the BlobKey from the BlobHeader which was sent to the disperser, and compares it with
// the BlobKey which was returned by the disperser in the DisperseBlobReply
//
// A successful verification guarantees that the disperser didn't make any modifications to the BlobHeader that it
// received from this client.
//
// This function returns nil if the verification succeeds, and otherwise returns an error describing the failure
func verifyReceivedBlobKey(
	// the blob header which was constructed locally and sent to the disperser
	blobHeader *corev2.BlobHeader,
	// the reply received back from the disperser
	disperserReply *disperser_rpc.DisperseBlobReply,
) error {

	actualBlobKey, err := blobHeader.BlobKey()
	if err != nil {
		// this shouldn't be possible, since the blob key has already been used when signing dispersal
		return fmt.Errorf("computing blob key: %w", err)
	}

	blobKeyFromDisperser, err := corev2.BytesToBlobKey(disperserReply.GetBlobKey())
	if err != nil {
		return fmt.Errorf("converting returned bytes to blob key: %w", err)
	}

	if actualBlobKey != blobKeyFromDisperser {
		return fmt.Errorf(
			"blob key returned by disperser (%v) doesn't match blob which was dispersed (%v)",
			blobKeyFromDisperser, actualBlobKey)
	}

	return nil
}

// GetBlobStatus returns the status of a blob with the given blob key.
func (c *disperserClient) GetBlobStatus(ctx context.Context, blobKey corev2.BlobKey) (*disperser_rpc.BlobStatusReply, error) {
	err := c.initOnceGrpcConnection()
	if err != nil {
		return nil, api.NewErrorInternal(err.Error())
	}

	request := &disperser_rpc.BlobStatusRequest{
		BlobKey: blobKey[:],
	}
	return c.client.GetBlobStatus(ctx, request)
}

// GetPaymentState returns the payment state of the disperser client
func (c *disperserClient) GetPaymentState(ctx context.Context) (*disperser_rpc.GetPaymentStateForAllQuorumsReply, error) {
	err := c.initOnceGrpcConnection()
	if err != nil {
		return nil, api.NewErrorInternal(err.Error())
	}

	accountID, err := c.signer.GetAccountID()
	if err != nil {
		return nil, fmt.Errorf("error getting signer's account ID: %w", err)
	}

	timestamp := uint64(time.Now().UnixNano())

	signature, err := c.signer.SignPaymentStateRequest(timestamp)
	if err != nil {
		return nil, fmt.Errorf("error signing payment state request: %w", err)
	}

	request := &disperser_rpc.GetPaymentStateForAllQuorumsRequest{
		AccountId: accountID.Hex(),
		Signature: signature,
		Timestamp: timestamp,
	}
	allQuorumsReply, err := c.client.GetPaymentStateForAllQuorums(ctx, request)
	if err != nil {
		// Check if error is "method not found" or "unimplemented"
		if isMethodNotFoundError(err) {
			// Fall back to old method
			return c.getPaymentStateFromLegacyAPI(ctx, accountID, signature, timestamp)
		}
		return nil, err
	}

	return allQuorumsReply, nil
}

// this is true if we are targeting a disperser that hasn't upgraded to the new API yet.
func isMethodNotFoundError(err error) bool {
	if st, ok := status.FromError(err); ok {
		return st.Code() == codes.Unimplemented
	}
	return false
}

// getPaymentStateFromLegacyAPI retrieves the payment state from the legacy GetPaymentState grpc method.
// It is needed until we have upgraded all dispersers (testnet and mainnet) to the new API.
// Check those endpoints for GetPaymentStateForAllQuorums using:
// `grpcurl disperser-testnet-holesky.eigenda.xyz:443 list disperser.v2.Disperser`
// `grpcurl disperser.eigenda.xyz:443 list disperser.v2.Disperser`
func (c *disperserClient) getPaymentStateFromLegacyAPI(
	ctx context.Context, accountID gethcommon.Address, signature []byte, timestamp uint64,
) (*disperser_rpc.GetPaymentStateForAllQuorumsReply, error) {
	oldRequest := &disperser_rpc.GetPaymentStateRequest{
		AccountId: accountID.Hex(),
		Signature: signature,
		Timestamp: timestamp,
	}

	oldResult, err := c.client.GetPaymentState(ctx, oldRequest)
	if err != nil {
		return nil, err
	}

	return convertLegacyPaymentStateToNew(oldResult)
}

// convertLegacyPaymentStateToNew converts the old GetPaymentStateReply to the new GetPaymentStateForAllQuorumsReply format
func convertLegacyPaymentStateToNew(legacyReply *disperser_rpc.GetPaymentStateReply) (*disperser_rpc.GetPaymentStateForAllQuorumsReply, error) {

	if legacyReply.PaymentGlobalParams == nil {
		return nil, fmt.Errorf("legacy payment state received from disperser does not contain global params")
	}
	// Convert PaymentGlobalParams to PaymentVaultParams
	var paymentVaultParams *disperser_rpc.PaymentVaultParams
	{
		paymentVaultParams = &disperser_rpc.PaymentVaultParams{
			QuorumPaymentConfigs:  make(map[uint32]*disperser_rpc.PaymentQuorumConfig),
			QuorumProtocolConfigs: make(map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig),
			OnDemandQuorumNumbers: legacyReply.PaymentGlobalParams.OnDemandQuorumNumbers,
		}

		// Apply the global params to all quorums, both on-demand and reservation.
		onDemandQuorums := legacyReply.PaymentGlobalParams.OnDemandQuorumNumbers
		if len(onDemandQuorums) == 0 {
			return nil, fmt.Errorf("no on-demand quorums specified in legacy PaymentGlobalParams received from disperser")
		}
		reservationQuorums := legacyReply.Reservation.QuorumNumbers
		// There may be overlapping quorums but it doesn't matter since we will apply the same global params to all of them.
		allQuorums := append(reservationQuorums, onDemandQuorums...)

		for _, quorumID := range allQuorums {
			paymentVaultParams.QuorumPaymentConfigs[quorumID] = &disperser_rpc.PaymentQuorumConfig{
				ReservationSymbolsPerSecond: 0, // Not available in legacy format
				OnDemandSymbolsPerSecond:    legacyReply.PaymentGlobalParams.GlobalSymbolsPerSecond,
				OnDemandPricePerSymbol:      legacyReply.PaymentGlobalParams.PricePerSymbol,
			}

			paymentVaultParams.QuorumProtocolConfigs[quorumID] = &disperser_rpc.PaymentQuorumProtocolConfig{
				MinNumSymbols: legacyReply.PaymentGlobalParams.MinNumSymbols,
				// ReservationAdvanceWindow is not used offchain at the moment so it's okay to set to any value.
				ReservationAdvanceWindow:   0,
				ReservationRateLimitWindow: legacyReply.PaymentGlobalParams.ReservationWindow,
				OnDemandRateLimitWindow:    0, // Not available in legacy format
			}
		}

		for _, quorumID := range onDemandQuorums {
			paymentVaultParams.QuorumProtocolConfigs[quorumID].OnDemandEnabled = true
		}
	}

	// If no reservation is available, return early with only payment vault params and cumulative payment info.
	if legacyReply.Reservation == nil {
		return &disperser_rpc.GetPaymentStateForAllQuorumsReply{
			PaymentVaultParams:       paymentVaultParams,
			CumulativePayment:        legacyReply.CumulativePayment,
			OnchainCumulativePayment: legacyReply.OnchainCumulativePayment,
		}, nil
	}

	// Otherwise there is a reservation available, so we need to convert it to the per-quorum format.

	// We first make sure that the disperser returned valid data.
	if len(legacyReply.PeriodRecords) == 0 {
		return nil, fmt.Errorf("legacy payment state received from disperser does not contain period records")
	}
	if len(legacyReply.Reservation.QuorumNumbers) == 0 {
		return nil, fmt.Errorf("legacy payment state received from disperser does not contain reservation quorums")
	}

	reservations := make(map[uint32]*disperser_rpc.QuorumReservation)
	periodRecords := make(map[uint32]*disperser_rpc.PeriodRecords)

	// Apply the reservation to all reservationQuorums mentioned in the reservation
	for _, quorumID := range legacyReply.Reservation.QuorumNumbers {
		reservations[quorumID] = &disperser_rpc.QuorumReservation{
			SymbolsPerSecond: legacyReply.Reservation.SymbolsPerSecond,
			StartTimestamp:   legacyReply.Reservation.StartTimestamp,
			EndTimestamp:     legacyReply.Reservation.EndTimestamp,
		}
		periodRecords[quorumID] = &disperser_rpc.PeriodRecords{
			Records: legacyReply.PeriodRecords,
		}
	}

	return &disperser_rpc.GetPaymentStateForAllQuorumsReply{
		PaymentVaultParams:       paymentVaultParams,
		PeriodRecords:            periodRecords,
		Reservations:             reservations,
		CumulativePayment:        legacyReply.CumulativePayment,
		OnchainCumulativePayment: legacyReply.OnchainCumulativePayment,
	}, nil
}

// GetBlobCommitment is a utility method that calculates commitment for a blob payload.
// While the blob commitment can be calculated by anyone, it requires SRS points to
// be loaded. For service that does not have access to SRS points, this method can be
// used to calculate the blob commitment in blob header, which is required for dispersal.
func (c *disperserClient) GetBlobCommitment(ctx context.Context, data []byte) (*disperser_rpc.BlobCommitmentReply, error) {
	err := c.initOnceGrpcConnection()
	if err != nil {
		return nil, api.NewErrorInternal(err.Error())
	}

	request := &disperser_rpc.BlobCommitmentRequest{
		Blob: data,
	}
	return c.client.GetBlobCommitment(ctx, request)
}

// initOnceGrpcConnection initializes the grpc connection and client if they are not already initialized.
// If initialization fails, it caches the error and will return it on every subsequent call.
func (c *disperserClient) initOnceGrpcConnection() error {
	var initErr error
	c.initOnceGrpc.Do(func() {
		addr := fmt.Sprintf("%v:%v", c.config.Hostname, c.config.Port)
		dialOptions := GetGrpcDialOptions(c.config.UseSecureGrpcFlag, 4*units.MiB)
		conn, err := grpc.NewClient(addr, dialOptions...)
		if err != nil {
			initErr = err
			return
		}
		c.conn = conn
		c.client = disperser_rpc.NewDisperserClient(conn)
	})
	if initErr != nil {
		return fmt.Errorf("initializing grpc connection: %w", initErr)
	}
	return nil
}

// initOncePopulateAccountant initializes the accountant if it is not already initialized.
// If initialization fails, it caches the error and will return it on every subsequent call.
func (c *disperserClient) initOncePopulateAccountant(ctx context.Context) error {
	var initErr error
	c.initOnceAccountant.Do(func() {
		if c.accountant == nil {
			err := c.PopulateAccountant(ctx)
			if err != nil {
				initErr = err
				return
			}
		}
	})
	if initErr != nil {
		return fmt.Errorf("populating accountant: %w", initErr)
	}
	return nil
}
