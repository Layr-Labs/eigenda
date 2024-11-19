package clients

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Hostname string
	Port     string
	// BlobDispersal Timeouts for both authenticated and unauthenticated dispersals
	// GetBlobStatus and RetrieveBlob timeouts are hardcoded to 60seconds
	// TODO: do we want to add config timeouts for those separate requests?
	Timeout           time.Duration
	UseSecureGrpcFlag bool
	// MaxRetrieveBlobSizeBytes is the maximum size of a blob that can be retrieved by using
	// the RetrieveBlob method. This is used to set the max message size for the grpc client.
	// DisperserClient uses a single underlying grpc channel shared for all methods,
	// but all other methods use the default 4MiB max message size, whereas RetrieveBlob
	// potentially needs a larger size.
	//
	// If not set, default value is 100MiB for forward compatibility.
	// Check official documentation for current max blob size on mainnet.
	MaxRetrieveBlobSizeBytes int
}

// Deprecated: Use &Config{...} directly instead
func NewConfig(hostname, port string, timeout time.Duration, useSecureGrpcFlag bool) *Config {
	return &Config{
		Hostname:          hostname,
		Port:              port,
		Timeout:           timeout,
		UseSecureGrpcFlag: useSecureGrpcFlag,
	}
}

type DisperserClient interface {
	Close() error
	DisperseBlob(ctx context.Context, data []byte, customQuorums []uint8) (*disperser.BlobStatus, []byte, error)
	// DisperseBlobAuthenticated disperses a blob with an authenticated request.
	// The BlobStatus returned will always be PROCESSSING if error is nil.
	DisperseBlobAuthenticated(ctx context.Context, data []byte, customQuorums []uint8) (*disperser.BlobStatus, []byte, error)
	GetBlobStatus(ctx context.Context, key []byte) (*disperser_rpc.BlobStatusReply, error)
	RetrieveBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error)
}

// See the NewDisperserClient constructor's documentation for details and usage examples.
type disperserClient struct {
	config *Config
	signer core.BlobRequestSigner
	// conn and client are not initialized in the constructor, but are initialized lazily
	// whenever a method is called, using initOnce to make sure initialization happens only once
	// and is thread-safe
	initOnce sync.Once
	// We use a single grpc connection, which allows a max number of concurrent open streams (from DisperseBlobAuthenticated).
	// This should be fine in most cases, as each such request should take <1sec per 1MB blob.
	// The MaxConcurrentStreams parameter is set by the server. If not set, then it defaults to the stdlib's
	// http2 default of 100-1000: https://github.com/golang/net/blob/4783315416d92ff3d4664762748bd21776b42b98/http2/transport.go#L55
	// This means a conservative estimate of 100-1000MB/sec, which should be amply sufficient.
	// If we ever need to increase this, we could either consider asking the disperser to increase its limit,
	// or to use a pool of connections here.
	// TODO: we should refactor or make a new constructor which allows setting conn and/or client
	//       via dependency injection. This would allow for testing via https://pkg.go.dev/google.golang.org/grpc/test/bufconn
	//       instead of a real network connection for eg.
	conn   *grpc.ClientConn
	client disperser_rpc.DisperserClient
}

var _ DisperserClient = &disperserClient{}

// DisperserClient maintains a single underlying grpc connection to the disperser server,
// through which it sends requests to disperse blobs, get blob status, and retrieve blobs.
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
//	status, requestId, err := client.DisperseBlob(ctx, someData, someQuorums)
//	if err != nil {
//	    // Handle error
//	}
//
//	// Subsequent calls will use the existing connection
//	status2, requestId2, err := client.DisperseBlob(ctx, otherData, otherQuorums)
func NewDisperserClient(config *Config, signer core.BlobRequestSigner) (*disperserClient, error) {
	if err := checkConfigAndSetDefaults(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &disperserClient{
		config: config,
		signer: signer,
		// conn and client are initialized lazily
	}, nil
}

func checkConfigAndSetDefaults(c *Config) error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}
	if c.Hostname == "" {
		return fmt.Errorf("config.Hostname is empty")
	}
	if c.Port == "" {
		return fmt.Errorf("config.Port is empty")
	}
	if c.Timeout == 0 {
		return fmt.Errorf("config.Timeout is 0")
	}
	if c.MaxRetrieveBlobSizeBytes == 0 {
		// Set to 100MiB for forward compatibility.
		// Check official documentation for current max blob size on mainnet.
		c.MaxRetrieveBlobSizeBytes = 100 * 1024 * 1024
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

func (c *disperserClient) DisperseBlob(ctx context.Context, data []byte, quorums []uint8) (*disperser.BlobStatus, []byte, error) {
	err := c.initOnceGrpcConnection()
	if err != nil {
		return nil, nil, api.NewErrorFailover(err)
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	quorumNumbers := make([]uint32, len(quorums))
	for i, q := range quorums {
		quorumNumbers[i] = uint32(q)
	}

	// check every 32 bytes of data are within the valid range for a bn254 field element
	_, err = rs.ToFrArray(data)
	if err != nil {
		return nil, nil, fmt.Errorf("encountered an error to convert a 32-bytes into a valid field element, please use the correct format where every 32bytes(big-endian) is less than 21888242871839275222246405745257275088548364400416034343698204186575808495617 %w", err)
	}
	request := &disperser_rpc.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: quorumNumbers,
	}

	reply, err := c.client.DisperseBlob(ctxTimeout, request)
	if err != nil {
		return nil, nil, err
	}

	blobStatus, err := disperser.FromBlobStatusProto(reply.GetResult())
	if err != nil {
		return nil, nil, err
	}

	return blobStatus, reply.GetRequestId(), nil
}

func (c *disperserClient) DisperseBlobAuthenticated(ctx context.Context, data []byte, quorums []uint8) (*disperser.BlobStatus, []byte, error) {
	err := c.initOnceGrpcConnection()
	if err != nil {
		return nil, nil, api.NewErrorFailover(err)
	}

	if c.signer == nil {
		return nil, nil, api.NewErrorInternal("uninitialized signer for authenticated dispersal")
	}

	// first check if signer is valid
	accountId, err := c.signer.GetAccountID()
	if err != nil {
		return nil, nil, api.NewErrorInvalidArg(fmt.Sprintf("please configure signer key if you want to use authenticated endpoint %v", err))
	}

	quorumNumbers := make([]uint32, len(quorums))
	for i, q := range quorums {
		quorumNumbers[i] = uint32(q)
	}

	// check every 32 bytes of data are within the valid range for a bn254 field element
	_, err = rs.ToFrArray(data)
	if err != nil {
		return nil, nil, api.NewErrorInvalidArg(
			fmt.Sprintf("encountered an error to convert a 32-bytes into a valid field element, "+
				"please use the correct format where every 32bytes(big-endian) is less than "+
				"21888242871839275222246405745257275088548364400416034343698204186575808495617, %v", err))
	}

	request := &disperser_rpc.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: quorumNumbers,
		AccountId:           accountId,
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	stream, err := c.client.DisperseBlobAuthenticated(ctxTimeout)
	if err != nil {
		// grpc client errors return grpc errors, so we can just wrap the error in a normal wrapError,
		// no need to wrap in another grpc error as we do with other errors above.
		return nil, nil, fmt.Errorf("error while calling DisperseBlobAuthenticated: %w", err)
	}

	// Send the initial request
	err = stream.Send(&disperser_rpc.AuthenticatedRequest{Payload: &disperser_rpc.AuthenticatedRequest_DisperseRequest{
		DisperseRequest: request,
	}})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Get the Challenge
	reply, err := stream.Recv()
	if err != nil {
		return nil, nil, fmt.Errorf("error while receiving: %w", err)
	}
	authHeaderReply, ok := reply.Payload.(*disperser_rpc.AuthenticatedReply_BlobAuthHeader)
	if !ok {
		return nil, nil, api.NewErrorInternal(fmt.Sprintf("client expected challenge from disperser, instead received: %v", reply))
	}

	authHeader := core.BlobAuthHeader{
		BlobCommitments: encoding.BlobCommitments{},
		AccountID:       "",
		Nonce:           authHeaderReply.BlobAuthHeader.ChallengeParameter,
	}
	authData, err := c.signer.SignBlobRequest(authHeader)
	if err != nil {
		return nil, nil, api.NewErrorInternal(fmt.Sprintf("error signing blob request: %v", err))
	}

	// Process challenge and send back challenge_reply
	err = stream.Send(&disperser_rpc.AuthenticatedRequest{Payload: &disperser_rpc.AuthenticatedRequest_AuthenticationData{
		AuthenticationData: &disperser_rpc.AuthenticationData{
			AuthenticationData: authData,
		},
	}})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send challenge reply: %w", err)
	}

	reply, err = stream.Recv()
	if err != nil {
		return nil, nil, fmt.Errorf("error while receiving final reply: %w", err)
	}
	disperseReply, ok := reply.Payload.(*disperser_rpc.AuthenticatedReply_DisperseReply) // Process the final disperse_reply
	if !ok {
		return nil, nil, api.NewErrorInternal(fmt.Sprintf("client expected DisperseReply from disperser, instead received: %v", reply))
	}

	blobStatus, err := disperser.FromBlobStatusProto(disperseReply.DisperseReply.GetResult())
	if err != nil {
		return nil, nil, api.NewErrorInternal(fmt.Sprintf("parsing blob status: %v", err))
	}

	// Assert: only status that makes sense is processing. Anything else is a bug on disperser side.
	if *blobStatus != disperser.Processing {
		return nil, nil, api.NewErrorInternal(fmt.Sprintf("expected status to be Processing, got %v", *blobStatus))
	}

	return blobStatus, disperseReply.DisperseReply.GetRequestId(), nil
}

func (c *disperserClient) GetBlobStatus(ctx context.Context, requestID []byte) (*disperser_rpc.BlobStatusReply, error) {
	err := c.initOnceGrpcConnection()
	if err != nil {
		return nil, api.NewErrorInternal(err.Error())
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	request := &disperser_rpc.BlobStatusRequest{
		RequestId: requestID,
	}
	return c.client.GetBlobStatus(ctxTimeout, request)
}

func (c *disperserClient) RetrieveBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error) {
	err := c.initOnceGrpcConnection()
	if err != nil {
		return nil, api.NewErrorInternal(err.Error())
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	reply, err := c.client.RetrieveBlob(ctxTimeout,
		&disperser_rpc.RetrieveBlobRequest{
			BatchHeaderHash: batchHeaderHash,
			BlobIndex:       blobIndex,
		},
		grpc.MaxCallRecvMsgSize(c.config.MaxRetrieveBlobSizeBytes)) // for client
	if err != nil {
		return nil, err
	}
	return reply.Data, nil
}

// initOnceGrpcConnection initializes the grpc connection and client if they are not already initialized.
// If initialization fails, it caches the error and will return it on every subsequent call.
func (c *disperserClient) initOnceGrpcConnection() error {
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

func getGrpcDialOptions(useSecureGrpcFlag bool) []grpc.DialOption {
	options := []grpc.DialOption{}
	if useSecureGrpcFlag {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	return options
}
