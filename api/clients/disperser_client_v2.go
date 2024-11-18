package clients

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/Layr-Labs/eigenda/api"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"google.golang.org/grpc"
)

type DisperserClientV2Config struct {
	Hostname          string
	Port              string
	UseSecureGrpcFlag bool
}

type DisperserClientV2 interface {
	Close() error
	DisperseBlob(ctx context.Context, data []byte, blobVersion corev2.BlobVersion, quorums []core.QuorumID) (*dispv2.BlobStatus, corev2.BlobKey, error)
	GetBlobStatus(ctx context.Context, blobKey corev2.BlobKey) (*disperser_rpc.BlobStatusReply, error)
	GetBlobCommitment(ctx context.Context, data []byte) (*disperser_rpc.BlobCommitmentReply, error)
}

type disperserClientV2 struct {
	config   *DisperserClientV2Config
	signer   corev2.BlobRequestSigner
	initOnce sync.Once
	conn     *grpc.ClientConn
	client   disperser_rpc.DisperserClient
	prover   encoding.Prover
}

var _ DisperserClientV2 = &disperserClientV2{}

// DisperserClientV2 maintains a single underlying grpc connection to the disperser server,
// through which it sends requests to disperse blobs and get blob status.
// The connection is established lazily on the first method call. Don't forget to call Close(),
// which is safe to call even if the connection was never established.
//
// DisperserClientV2 is safe to be used concurrently by multiple goroutines.
//
// Example usage:
//
//	client := NewDisperserClientV2(config, signer)
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
func NewDisperserClientV2(config *DisperserClientV2Config, signer corev2.BlobRequestSigner, prover encoding.Prover) (*disperserClientV2, error) {
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

	return &disperserClientV2{
		config: config,
		signer: signer,
		prover: prover,
		// conn and client are initialized lazily
	}, nil
}

// Close closes the grpc connection to the disperser server.
// It is thread safe and can be called multiple times.
func (c *disperserClientV2) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		c.client = nil
		return err
	}
	return nil
}

func (c *disperserClientV2) DisperseBlob(
	ctx context.Context,
	data []byte,
	blobVersion corev2.BlobVersion,
	quorums []core.QuorumID,
) (*dispv2.BlobStatus, corev2.BlobKey, error) {
	err := c.initOnceGrpcConnection()
	if err != nil {
		return nil, [32]byte{}, api.NewErrorFailover(err)
	}

	if c.signer == nil {
		return nil, [32]byte{}, api.NewErrorInternal("uninitialized signer for authenticated dispersal")
	}

	var payment core.PaymentMetadata
	accountId, err := c.signer.GetAccountID()
	if err != nil {
		return nil, [32]byte{}, api.NewErrorInvalidArg(fmt.Sprintf("please configure signer key if you want to use authenticated endpoint %v", err))
	}
	payment.AccountID = accountId
	// TODO: add payment metadata
	payment.BinIndex = 0
	payment.CumulativePayment = big.NewInt(0)

	if len(quorums) == 0 {
		return nil, [32]byte{}, api.NewErrorInvalidArg("quorum numbers must be provided")
	}

	for _, q := range quorums {
		if q > corev2.MaxQuorumID {
			return nil, [32]byte{}, api.NewErrorInvalidArg("quorum number must be less than 256")
		}
	}

	// check every 32 bytes of data are within the valid range for a bn254 field element
	_, err = rs.ToFrArray(data)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("encountered an error to convert a 32-bytes into a valid field element, please use the correct format where every 32bytes(big-endian) is less than 21888242871839275222246405745257275088548364400416034343698204186575808495617 %w", err)
	}

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
	} else {
		// if prover is configured, get commitments from prover
		blobCommitments, err = c.prover.GetCommitments(data)
		if err != nil {
			return nil, [32]byte{}, fmt.Errorf("error getting blob commitments: %w", err)
		}
	}

	blobHeader := &corev2.BlobHeader{
		BlobVersion:     blobVersion,
		BlobCommitments: blobCommitments,
		QuorumNumbers:   quorums,
		PaymentMetadata: payment,
	}
	sig, err := c.signer.SignBlobRequest(blobHeader)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("error signing blob request: %w", err)
	}
	blobHeader.Signature = sig
	blobHeaderProto, err := blobHeader.ToProtobuf()
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("error converting blob header to protobuf: %w", err)
	}
	request := &disperser_rpc.DisperseBlobRequest{
		Data:       data,
		BlobHeader: blobHeaderProto,
	}

	reply, err := c.client.DisperseBlob(ctx, request)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("error while calling DisperseBlob: %w", err)
	}

	blobStatus, err := dispv2.BlobStatusFromProtobuf(reply.GetResult())
	if err != nil {
		return nil, [32]byte{}, err
	}

	return &blobStatus, corev2.BlobKey(reply.GetBlobKey()), nil
}

// GetBlobStatus returns the status of a blob with the given blob key.
func (c *disperserClientV2) GetBlobStatus(ctx context.Context, blobKey corev2.BlobKey) (*disperser_rpc.BlobStatusReply, error) {
	err := c.initOnceGrpcConnection()
	if err != nil {
		return nil, api.NewErrorInternal(err.Error())
	}

	request := &disperser_rpc.BlobStatusRequest{
		BlobKey: blobKey[:],
	}
	return c.client.GetBlobStatus(ctx, request)
}

// GetBlobCommitment is a utility method that calculates commitment for a blob payload.
// While the blob commitment can be calculated by anyone, it requires SRS points to
// be loaded. For service that does not have access to SRS points, this method can be
// used to calculate the blob commitment in blob header, which is required for dispersal.
func (c *disperserClientV2) GetBlobCommitment(ctx context.Context, data []byte) (*disperser_rpc.BlobCommitmentReply, error) {
	err := c.initOnceGrpcConnection()
	if err != nil {
		return nil, api.NewErrorInternal(err.Error())
	}

	request := &disperser_rpc.BlobCommitmentRequest{
		Data: data,
	}
	return c.client.GetBlobCommitment(ctx, request)
}

// initOnceGrpcConnection initializes the grpc connection and client if they are not already initialized.
// If initialization fails, it caches the error and will return it on every subsequent call.
func (c *disperserClientV2) initOnceGrpcConnection() error {
	var initErr error
	c.initOnce.Do(func() {
		addr := fmt.Sprintf("%v:%v", c.config.Hostname, c.config.Port)
		dialOptions := getGrpcDialOptions(c.config.UseSecureGrpcFlag)
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
