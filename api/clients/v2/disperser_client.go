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
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
	"google.golang.org/grpc"
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
	logger             logging.Logger
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
func NewDisperserClient(config *DisperserClientConfig, signer corev2.BlobRequestSigner, prover encoding.Prover, accountant *Accountant) (*disperserClient, error) {
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
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	logger = logger.With("component", "DisperserClient")

	ntpClock, err := core.NewNTPSyncedClock(context.Background(), config.NtpServer, config.NtpSyncInterval, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create NTP clock: %w", err)
	}

	return &disperserClient{
		config:     config,
		signer:     signer,
		prover:     prover,
		accountant: accountant,
		ntpClock:   ntpClock,
		logger:     logger,
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
		c.accountant = NewAccountant(accountId, nil, nil, 0, 0, 0, 0, c.logger)
	}

	paymentState, err := c.GetPaymentState(ctx)
	if err != nil {
		return fmt.Errorf("error getting payment state for initializing accountant: %w", err)
	}

	err = c.accountant.SetPaymentState(paymentState)
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
		return nil, [32]byte{}, api.NewErrorFailover(err)
	}

	symbolLength := encoding.GetBlobLengthPowerOf2(uint(len(data)))
	payment, err := c.accountant.AccountBlob(ctx, c.ntpClock.Now().UnixNano(), uint64(symbolLength), quorums)
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
func (c *disperserClient) GetPaymentState(ctx context.Context) (*disperser_rpc.GetPaymentStateReply, error) {
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

	request := &disperser_rpc.GetPaymentStateRequest{
		AccountId: accountID.Hex(),
		Signature: signature,
		Timestamp: timestamp,
	}
	return c.client.GetPaymentState(ctx, request)
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
