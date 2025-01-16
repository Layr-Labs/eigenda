package clients

import (
	"context"
	"errors"
	"fmt"
	"sync"

	relaygrpc "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/grpc"
)

// MessageSigner is a function that signs a message with a private BLS key.
type MessageSigner func(ctx context.Context, data [32]byte) (*core.Signature, error)

type RelayClientConfig struct {
	Sockets           map[corev2.RelayKey]string
	UseSecureGrpcFlag bool
	OperatorID        *core.OperatorID
	MessageSigner     MessageSigner
}

type ChunkRequestByRange struct {
	BlobKey corev2.BlobKey
	Start   uint32
	End     uint32
}

type ChunkRequestByIndex struct {
	BlobKey corev2.BlobKey
	Indices []uint32
}

type RelayClient interface {
	// GetBlob retrieves a blob from a relay
	GetBlob(ctx context.Context, relayKey corev2.RelayKey, blobKey corev2.BlobKey) ([]byte, error)
	// GetChunksByRange retrieves blob chunks from a relay by chunk index range
	// The returned slice has the same length and ordering as the input slice, and the i-th element is the bundle for the i-th request.
	// Each bundle is a sequence of frames in raw form (i.e., serialized core.Bundle bytearray).
	GetChunksByRange(ctx context.Context, relayKey corev2.RelayKey, requests []*ChunkRequestByRange) ([][]byte, error)
	// GetChunksByIndex retrieves blob chunks from a relay by index
	// The returned slice has the same length and ordering as the input slice, and the i-th element is the bundle for the i-th request.
	// Each bundle is a sequence of frames in raw form (i.e., serialized core.Bundle bytearray).
	GetChunksByIndex(ctx context.Context, relayKey corev2.RelayKey, requests []*ChunkRequestByIndex) ([][]byte, error)
	// GetSockets returns the relay sockets
	GetSockets() map[corev2.RelayKey]string
	Close() error
}

type relayClient struct {
	config *RelayClientConfig

	// initOnce is used to ensure that the connection to each relay is initialized only once.
	// It maps relay key to a sync.Once instance: `map[corev2.RelayKey]*sync.Once`
	initOnce *sync.Map
	// conns maps relay key to the gRPC connection: `map[corev2.RelayKey]*grpc.ClientConn`
	conns  sync.Map
	logger logging.Logger

	// grpcClients maps relay key to the gRPC client: `map[corev2.RelayKey]relaygrpc.RelayClient`
	grpcClients sync.Map
}

var _ RelayClient = (*relayClient)(nil)

// NewRelayClient creates a new RelayClient that connects to the relays specified in the config.
// It keeps a connection to each relay and reuses it for subsequent requests, and the connection is lazily instantiated.
func NewRelayClient(config *RelayClientConfig, logger logging.Logger) (RelayClient, error) {
	if config == nil || len(config.Sockets) <= 0 {
		return nil, fmt.Errorf("invalid config: %v", config)
	}

	logger.Info("creating relay client", "urls", config.Sockets)

	initOnce := sync.Map{}
	for key := range config.Sockets {
		initOnce.Store(key, &sync.Once{})
	}
	return &relayClient{
		config: config,

		initOnce: &initOnce,
		logger:   logger.With("component", "RelayClient"),
	}, nil
}

func (c *relayClient) GetBlob(ctx context.Context, relayKey corev2.RelayKey, blobKey corev2.BlobKey) ([]byte, error) {
	client, err := c.getClient(relayKey)
	if err != nil {
		return nil, err
	}

	res, err := client.GetBlob(ctx, &relaygrpc.GetBlobRequest{
		BlobKey: blobKey[:],
	})
	if err != nil {
		return nil, err
	}

	return res.GetBlob(), nil
}

// signGetChunksRequest signs the GetChunksRequest with the operator's private key
// and sets the signature in the request.
func (c *relayClient) signGetChunksRequest(ctx context.Context, request *relaygrpc.GetChunksRequest) error {
	if c.config.OperatorID == nil {
		return errors.New("no operator ID provided in config, cannot sign get chunks request")
	}
	if c.config.MessageSigner == nil {
		return errors.New("no message signer provided in config, cannot sign get chunks request")
	}

	hash := hashing.HashGetChunksRequest(request)
	hashArray := [32]byte{}
	copy(hashArray[:], hash)
	signature, err := c.config.MessageSigner(ctx, hashArray)
	if err != nil {
		return fmt.Errorf("failed to sign get chunks request: %v", err)
	}
	sig := signature.SerializeCompressed()
	request.OperatorSignature = sig[:]
	return nil
}

func (c *relayClient) GetChunksByRange(
	ctx context.Context,
	relayKey corev2.RelayKey,
	requests []*ChunkRequestByRange) ([][]byte, error) {

	if len(requests) == 0 {
		return nil, fmt.Errorf("no requests")
	}
	client, err := c.getClient(relayKey)
	if err != nil {
		return nil, err
	}

	grpcRequests := make([]*relaygrpc.ChunkRequest, len(requests))
	for i, req := range requests {
		grpcRequests[i] = &relaygrpc.ChunkRequest{
			Request: &relaygrpc.ChunkRequest_ByRange{
				ByRange: &relaygrpc.ChunkRequestByRange{
					BlobKey:    req.BlobKey[:],
					StartIndex: req.Start,
					EndIndex:   req.End,
				},
			},
		}
	}

	request := &relaygrpc.GetChunksRequest{
		ChunkRequests: grpcRequests,
		OperatorId:    c.config.OperatorID[:],
	}
	err = c.signGetChunksRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	res, err := client.GetChunks(ctx, request)
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func (c *relayClient) GetChunksByIndex(
	ctx context.Context,
	relayKey corev2.RelayKey,
	requests []*ChunkRequestByIndex) ([][]byte, error) {

	if len(requests) == 0 {
		return nil, fmt.Errorf("no requests")
	}

	client, err := c.getClient(relayKey)
	if err != nil {
		return nil, err
	}

	grpcRequests := make([]*relaygrpc.ChunkRequest, len(requests))
	for i, req := range requests {
		grpcRequests[i] = &relaygrpc.ChunkRequest{
			Request: &relaygrpc.ChunkRequest_ByIndex{
				ByIndex: &relaygrpc.ChunkRequestByIndex{
					BlobKey:      req.BlobKey[:],
					ChunkIndices: req.Indices,
				},
			},
		}
	}

	request := &relaygrpc.GetChunksRequest{
		ChunkRequests: grpcRequests,
		OperatorId:    c.config.OperatorID[:],
	}
	err = c.signGetChunksRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	res, err := client.GetChunks(ctx, request)

	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func (c *relayClient) getClient(key corev2.RelayKey) (relaygrpc.RelayClient, error) {
	if err := c.initOnceGrpcConnection(key); err != nil {
		return nil, err
	}
	maybeClient, ok := c.grpcClients.Load(key)
	if !ok {
		return nil, fmt.Errorf("no grpc client for relay key: %v", key)
	}
	client, ok := maybeClient.(relaygrpc.RelayClient)
	if !ok {
		return nil, fmt.Errorf("invalid grpc client for relay key: %v", key)
	}
	return client, nil
}

func (c *relayClient) initOnceGrpcConnection(key corev2.RelayKey) error {
	var initErr error
	once, ok := c.initOnce.Load(key)
	if !ok {
		return fmt.Errorf("unknown relay key: %v", key)
	}
	once.(*sync.Once).Do(func() {
		socket, ok := c.config.Sockets[key]
		if !ok {
			initErr = fmt.Errorf("unknown relay key: %v", key)
			return
		}
		dialOptions := getGrpcDialOptions(c.config.UseSecureGrpcFlag)
		conn, err := grpc.NewClient(socket, dialOptions...)
		if err != nil {
			initErr = err
			return
		}
		c.conns.Store(key, conn)
		c.grpcClients.Store(key, relaygrpc.NewRelayClient(conn))
	})
	return initErr
}

func (c *relayClient) GetSockets() map[corev2.RelayKey]string {
	return c.config.Sockets
}

func (c *relayClient) Close() error {
	var errList *multierror.Error
	c.conns.Range(func(k, v interface{}) bool {
		conn, ok := v.(*grpc.ClientConn)
		if !ok {
			errList = multierror.Append(errList, fmt.Errorf("invalid connection for relay key: %v", k))
			return true
		}

		if conn != nil {
			err := conn.Close()
			c.conns.Delete(k)
			c.grpcClients.Delete(k)
			if err != nil {
				c.logger.Error("failed to close connection", "err", err)
				errList = multierror.Append(errList, err)
			}
		}
		return true
	})
	return errList.ErrorOrNil()
}
