package dispersal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	clients "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
)

const maxNumberOfConnections = 32

type DisperserClientConfig struct {
	NetworkAddress    string
	UseSecureGrpcFlag bool
	// The number of grpc connections to the disperser server. A value of 0 is treated as 1.
	DisperserConnectionCount uint
}

// DisperserClient manages communication with the disperser server.
type DisperserClient struct {
	logger     logging.Logger
	config     *DisperserClientConfig
	signer     corev2.BlobRequestSigner
	clientPool *common.GRPCClientPool[disperser_rpc.DisperserClient]
	committer  *committer.Committer
	metrics    metrics.DispersalMetricer
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
	metrics metrics.DispersalMetricer,
) (*DisperserClient, error) {
	if config == nil {
		return nil, fmt.Errorf("config must be provided")
	}
	if strings.TrimSpace(config.NetworkAddress) == "" {
		return nil, fmt.Errorf("network address must be provided")
	}
	if signer == nil {
		return nil, fmt.Errorf("signer must be provided")
	}
	if committer == nil {
		return nil, fmt.Errorf("committer must be provided")
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

	dialOptions := clients.GetGrpcDialOptions(config.UseSecureGrpcFlag, 4*units.MiB)
	clientPool, err := common.NewGRPCClientPool(
		logger,
		disperser_rpc.NewDisperserClient,
		connectionCount,
		config.NetworkAddress,
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
		metrics:    metrics,
	}, nil
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
			return nil, nil, api.NewErrorInvalidArg(fmt.Sprintf("quorum number %d must be <= %d", q, corev2.MaxQuorumID))
		}
	}

	if paymentMetadata == nil {
		//nolint:wrapcheck
		return nil, nil, api.NewErrorInvalidArg("payment metadata must be provided")
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

	blobCommitments, err := c.committer.GetCommitmentsForPaddedLength(data)
	if err != nil {
		return nil, nil, fmt.Errorf("get commitments for padded length: %w", err)
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

	client, err := c.clientPool.GetClient()
	if err != nil {
		return nil, nil, fmt.Errorf("get client: %w", err)
	}

	reply, err := client.DisperseBlob(ctx, request)
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

	client, err := c.clientPool.GetClient()
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	reply, err := client.GetBlobStatus(ctx, request)
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

	client, err := c.clientPool.GetClient()
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	reply, err := client.GetPaymentState(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error while calling GetPaymentState: %w", err)
	}
	return reply, nil
}
