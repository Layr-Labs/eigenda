package clients

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
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
	UseSecureGrpcFlag  bool
	MaxGRPCMessageSize uint
	OperatorID         *core.OperatorID
	MessageSigner      MessageSigner
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
	Close() error
}

// relayClient is a client for the entire relay subsystem.
//
// It is a wrapper around a collection of grpc relay clients, which are used to interact with individual relays.
type relayClient struct {
	logger logging.Logger
	config *RelayClientConfig
	// initOnce is used to ensure that the connection to each relay is initialized only once
	initOnce map[corev2.RelayKey]*sync.Once
	// initOnceMutex protects access to the initOnce map
	initOnceMutex sync.Mutex
	// clientConnections maps relay key to the gRPC connection: `map[corev2.RelayKey]*grpc.ClientConn`
	// this map is maintained so that connections can be closed in Close
	clientConnections sync.Map
	// grpcRelayClients maps relay key to the gRPC client: `map[corev2.RelayKey]relaygrpc.RelayClient`
	// these grpc relay clients are used to communicate with individual relays
	grpcRelayClients sync.Map
	// relayUrlProvider knows how to retrieve the relay URLs, and maintains an internal URL cache
	relayUrlProvider relay.RelayUrlProvider
}

var _ RelayClient = (*relayClient)(nil)

// NewRelayClient creates a new RelayClient that connects to the relays specified in the config.
// It keeps a connection to each relay and reuses it for subsequent requests, and the connection is lazily instantiated.
func NewRelayClient(
	config *RelayClientConfig,
	logger logging.Logger,
	relayUrlProvider relay.RelayUrlProvider,
) (RelayClient, error) {

	if config == nil {
		return nil, errors.New("nil config")
	}

	if config.MaxGRPCMessageSize == 0 {
		return nil, errors.New("max gRPC message size must be greater than 0")
	}

	logger.Info("creating relay client")

	return &relayClient{
		config:           config,
		logger:           logger.With("component", "RelayClient"),
		relayUrlProvider: relayUrlProvider,
	}, nil
}

func (c *relayClient) GetBlob(ctx context.Context, relayKey corev2.RelayKey, blobKey corev2.BlobKey) ([]byte, error) {
	client, err := c.getClient(ctx, relayKey)
	if err != nil {
		return nil, fmt.Errorf("get grpc client for key %d: %w", relayKey, err)
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

	client, err := c.getClient(ctx, relayKey)
	if err != nil {
		return nil, fmt.Errorf("get grpc relay client for key %d: %w", relayKey, err)
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

	client, err := c.getClient(ctx, relayKey)
	if err != nil {
		return nil, fmt.Errorf("get grpc relay client for key %d: %w", relayKey, err)
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

// getClient gets the grpc relay client, which has a connection to a given relay
func (c *relayClient) getClient(ctx context.Context, key corev2.RelayKey) (relaygrpc.RelayClient, error) {
	if err := c.initOnceGrpcConnection(ctx, key); err != nil {
		return nil, fmt.Errorf("init grpc connection for key %d: %w", key, err)
	}
	maybeClient, ok := c.grpcRelayClients.Load(key)
	if !ok {
		return nil, fmt.Errorf("no grpc client for relay key: %v", key)
	}
	client, ok := maybeClient.(relaygrpc.RelayClient)
	if !ok {
		return nil, fmt.Errorf("invalid grpc client for relay key: %v", key)
	}
	return client, nil
}

// initOnceGrpcConnection initializes the GRPC connection for a given relay, and is guaranteed to only do perform
// the initialization once per relay.
func (c *relayClient) initOnceGrpcConnection(ctx context.Context, key corev2.RelayKey) error {
	// we must use a mutex here instead of a sync.Map, because this method could be called concurrently, and if
	// two concurrent calls tried to `LoadOrStore` from a sync.Map at the same time, it's possible they would
	// each create a unique sync.Once object, and perform duplicate initialization
	c.initOnceMutex.Lock()
	once, ok := c.initOnce[key]
	if !ok {
		once = &sync.Once{}
		c.initOnce[key] = once
	}
	c.initOnceMutex.Unlock()

	var initErr error
	once.Do(
		func() {
			relayUrl, err := c.relayUrlProvider.GetRelayUrl(ctx, key)
			if err != nil {
				initErr = fmt.Errorf("get relay url for key %d: %w", key, err)
				return
			}

			dialOptions := getGrpcDialOptions(c.config.UseSecureGrpcFlag, c.config.MaxGRPCMessageSize)
			conn, err := grpc.NewClient(relayUrl, dialOptions...)
			if err != nil {
				initErr = fmt.Errorf("create grpc client for key %d: %w", key, err)
				return
			}
			c.clientConnections.Store(key, conn)
			c.grpcRelayClients.Store(key, relaygrpc.NewRelayClient(conn))
		})

	return initErr
}

func (c *relayClient) Close() error {
	var errList *multierror.Error
	c.clientConnections.Range(
		func(k, v interface{}) bool {
			conn, ok := v.(*grpc.ClientConn)
			if !ok {
				errList = multierror.Append(errList, fmt.Errorf("invalid connection for relay key: %v", k))
				return true
			}

			if conn != nil {
				err := conn.Close()
				c.clientConnections.Delete(k)
				c.grpcRelayClients.Delete(k)
				if err != nil {
					c.logger.Error("failed to close connection", "err", err)
					errList = multierror.Append(errList, err)
				}
			}
			return true
		})
	return errList.ErrorOrNil()
}
