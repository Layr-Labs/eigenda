package clients

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
)

const maxNumberOfConnections = 32

type DisperserClientConfig struct {
	Hostname          string
	Port              string
	UseSecureGrpcFlag bool
	// The number of grpc connections to the disperser server. A value of 0 is treated as 1.
	DisperserConnectionCount uint
}

// DisperserClient manages communication with the disperser server.
type DisperserClient struct {
	logger                  logging.Logger
	config                  *DisperserClientConfig
	signer                  corev2.BlobRequestSigner
	clientPool              *common.GRPCClientPool[disperser_rpc.DisperserClient]
	committer               *committer.Committer
	accountant              *Accountant
	accountantLock          sync.Mutex
	initOnceAccountant      sync.Once
	initOnceAccountantError error
	metrics                 metrics.DispersalMetricer
}

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
func NewDisperserClient(
	logger logging.Logger,
	config *DisperserClientConfig,
	signer corev2.BlobRequestSigner,
	committer *committer.Committer,
	accountant *Accountant,
	metrics metrics.DispersalMetricer,
) (*DisperserClient, error) {
	if config == nil {
		return nil, fmt.Errorf("config must be provided")
	}
	if strings.TrimSpace(config.Hostname) == "" {
		return nil, fmt.Errorf("hostname must be provided")
	}
	if strings.TrimSpace(config.Port) == "" {
		return nil, fmt.Errorf("port must be provided")
	}
	if signer == nil {
		return nil, fmt.Errorf("signer must be provided")
	}
	if metrics == nil {
		return nil, fmt.Errorf("metrics must be provided")
	}

	var connectionCount uint
	if config.DisperserConnectionCount == 0 {
		connectionCount = 1
	}
	if config.DisperserConnectionCount > maxNumberOfConnections {
		connectionCount = maxNumberOfConnections
	}

	addr := fmt.Sprintf("%v:%v", config.Hostname, config.Port)
	dialOptions := GetGrpcDialOptions(config.UseSecureGrpcFlag, 4*units.MiB)
	clientPool, err := common.NewGRPCClientPool(
		logger,
		disperser_rpc.NewDisperserClient,
		connectionCount,
		addr,
		dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("new grpc client pool: %w", err)
	}

	return &DisperserClient{
		logger:     logger,
		config:     config,
		signer:     signer,
		clientPool: clientPool,
		committer:  committer,
		accountant: accountant,
		metrics:    metrics,
	}, nil
}

// PopulateAccountant populates the accountant with the payment state from the disperser.
func (c *DisperserClient) PopulateAccountant(ctx context.Context) error {
	if c.accountant == nil {
		return fmt.Errorf("accountant is nil")
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
func (c *DisperserClient) Close() error {
	if c.clientPool != nil {
		err := c.clientPool.Close()
		if err != nil {
			return fmt.Errorf("error closing client pool: %w", err)
		}
	}
	return nil
}

// Disperses a blob with the given data, blob version, and quorums.
//
// Returns the BlobHeader of the blob that was dispersed, and the DisperseBlobReply that was received from the
// disperser, if the dispersal was successful. Otherwise returns an error
func (c *DisperserClient) DisperseBlob(
	ctx context.Context,
	data []byte,
	blobVersion corev2.BlobVersion,
	quorums []core.QuorumID,
	probe *common.SequenceProbe,
	// if this is nil, that indicates we will use the legacy payment system to create the paymentMetadata
	paymentMetadata *core.PaymentMetadata,
) (*corev2.BlobHeader, *disperser_rpc.DisperseBlobReply, error) {
	if len(quorums) == 0 {
		//nolint:wrapcheck
		return nil, nil, api.NewErrorInvalidArg("quorum numbers must be provided")
	}
	if c.signer == nil {
		//nolint:wrapcheck
		return nil, nil, api.NewErrorInternal("uninitialized signer for authenticated dispersal")
	}
	for _, q := range quorums {
		if q > corev2.MaxQuorumID {
			//nolint:wrapcheck
			return nil, nil, api.NewErrorInvalidArg("quorum number must be less than 256")
		}
	}

	symbolLength := encoding.GetBlobLengthPowerOf2(uint32(len(data)))

	if paymentMetadata == nil {
		// we are using the legacy payment system
		probe.SetStage("acquire_accountant_lock")
		c.accountantLock.Lock()

		probe.SetStage("accountant")

		err := c.initOncePopulateAccountant(ctx)
		if err != nil {
			c.accountantLock.Unlock()
			return nil, nil, fmt.Errorf("error initializing accountant: %w", err)
		}

		paymentMetadata, err = c.accountant.AccountBlob(time.Now().UnixNano(), uint64(symbolLength), quorums)
		if err != nil {
			c.accountantLock.Unlock()
			return nil, nil, fmt.Errorf("error accounting blob: %w", err)
		}

		if paymentMetadata.CumulativePayment == nil || paymentMetadata.CumulativePayment.Sign() == 0 {
			// This request is using reserved bandwidth, no need to prevent parallel dispersal.
			c.accountantLock.Unlock()
		} else {
			// This request is using on-demand bandwidth, current implementation requires sequential dispersal.
			defer c.accountantLock.Unlock()
		}
	}

	probe.SetStage("verify_field_element")

	// check every 32 bytes of data are within the valid range for a bn254 field element
	_, err := rs.ToFrArray(data)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"encountered an error to convert a 32-bytes into a valid field element, "+
				"please use the correct format where every 32bytes(big-endian) is less than "+
				"21888242871839275222246405745257275088548364400416034343698204186575808495617 %w", err)
	}

	probe.SetStage("get_commitments")

	var blobCommitments encoding.BlobCommitments
	if c.committer == nil {
		// if committer is not configured, get blob commitments from disperser
		commitments, err := c.GetBlobCommitment(ctx, data)
		if err != nil {
			// Failover worthy error because it means the disperser is not responsive.
			return nil, nil, api.NewErrorFailover(fmt.Errorf("GetBlobCommitment rpc: %w", err))
		}
		deserialized, err := encoding.BlobCommitmentsFromProtobuf(commitments.GetBlobCommitment())
		if err != nil {
			return nil, nil, fmt.Errorf("error deserializing blob commitments: %w", err)
		}
		blobCommitments = *deserialized

		// We need to check that the disperser used the correct length. Even once checking the commitment from the
		// disperser has been implemented, there is still an edge case where the disperser could truncate trailing 0s,
		// yielding the wrong blob length, but not causing commitment verification to fail. It is important that the
		// commitment doesn't report a blob length smaller than expected, since this could cause payload parsing to
		// fail, if the length claimed in the encoded payload header is larger than the blob length in the commitment.
		lengthFromCommitment := commitments.GetBlobCommitment().GetLength()
		if lengthFromCommitment != uint32(symbolLength) {
			return nil, nil, fmt.Errorf(
				"blob commitment length (%d) from disperser doesn't match expected length (%d): %w",
				lengthFromCommitment, symbolLength, err)
		}
	} else {
		// if committer is configured, get commitments from committer
		blobCommitments, err = c.committer.GetCommitmentsForPaddedLength(data)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting blob commitments: %w", err)
		}
	}

	blobHeader := &corev2.BlobHeader{
		BlobVersion:     blobVersion,
		BlobCommitments: blobCommitments,
		QuorumNumbers:   quorums,
		PaymentMetadata: *paymentMetadata,
	}

	probe.SetStage("sign_blob_request")

	sig, err := c.signer.SignBlobRequest(blobHeader)
	if err != nil {
		return nil, nil, fmt.Errorf("error signing blob request: %w", err)
	}
	blobHeaderProto, err := blobHeader.ToProtobuf()
	if err != nil {
		return nil, nil, fmt.Errorf("error converting blob header to protobuf: %w", err)
	}
	request := &disperser_rpc.DisperseBlobRequest{
		Blob:       data,
		Signature:  sig,
		BlobHeader: blobHeaderProto,
	}

	probe.SetStage("send_to_disperser")

	reply, err := c.clientPool.GetClient().DisperseBlob(ctx, request)
	if err != nil {
		return nil, nil, api.NewErrorFailover(fmt.Errorf("DisperseBlob rpc: %w", err))
	}

	c.metrics.RecordBlobSizeBytes(len(data))

	return blobHeader, reply, nil
}

// GetBlobStatus returns the status of a blob with the given blob key.
func (c *DisperserClient) GetBlobStatus(
	ctx context.Context,
	blobKey corev2.BlobKey,
) (*disperser_rpc.BlobStatusReply, error) {
	request := &disperser_rpc.BlobStatusRequest{
		BlobKey: blobKey[:],
	}
	reply, err := c.clientPool.GetClient().GetBlobStatus(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error while calling GetBlobStatus: %w", err)
	}
	return reply, nil
}

// GetPaymentState returns the payment state of the disperser client
func (c *DisperserClient) GetPaymentState(ctx context.Context) (*disperser_rpc.GetPaymentStateReply, error) {
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
	reply, err := c.clientPool.GetClient().GetPaymentState(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error while calling GetPaymentState: %w", err)
	}
	return reply, nil
}

// GetBlobCommitment is a utility method that calculates commitment for a blob payload.
// While the blob commitment can be calculated by anyone, it requires SRS points to
// be loaded. For service that does not have access to SRS points, this method can be
// used to calculate the blob commitment in blob header, which is required for dispersal.
func (c *DisperserClient) GetBlobCommitment(
	ctx context.Context,
	data []byte,
) (*disperser_rpc.BlobCommitmentReply, error) {
	request := &disperser_rpc.BlobCommitmentRequest{
		Blob: data,
	}
	reply, err := c.clientPool.GetClient().GetBlobCommitment(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error while calling GetBlobCommitment: %w", err)
	}
	return reply, nil
}

// initOncePopulateAccountant initializes the accountant if it is not already initialized.
// If initialization fails, it caches the error and will return it on every subsequent call.
func (c *DisperserClient) initOncePopulateAccountant(ctx context.Context) error {
	c.initOnceAccountant.Do(func() {
		err := c.PopulateAccountant(ctx)
		if err != nil {
			c.initOnceAccountantError = err
			return
		}
	})
	if c.initOnceAccountantError != nil {
		return fmt.Errorf("populating accountant: %w", c.initOnceAccountantError)
	}
	return nil
}
