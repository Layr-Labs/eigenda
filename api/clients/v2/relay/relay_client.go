package relay

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	relaygrpc "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/hashicorp/go-multierror"
)

// MessageSigner is a function that signs a message with a private BLS key.
type MessageSigner func(ctx context.Context, data [32]byte) (*core.Signature, error)

type RelayClientConfig struct {
	UseSecureGrpcFlag  bool
	MaxGRPCMessageSize uint
	OperatorID         *core.OperatorID
	MessageSigner      MessageSigner
	// The number of parallel connections open to each relay.
	ConnectionPoolSize int // TODO make sure this is configured with flags
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
	// relayLockProvider provides locks that correspond to individual relay keys
	relayLockProvider *KeyLock[corev2.RelayKey]
	// connectionPoolSize is the number of parallel connections open to each relay.
	connectionPoolSize int
	// relayInitializationStatus maps relay key to a bool `map[corev2.RelayKey]bool`
	// the boolean value indicates whether the connection to that relay has been initialized
	relayInitializationStatus sync.Map
	// For each relay, we maintain a pool of gRPC clients that can be used to make requests to that relay. The key
	// in this map is the relay key, and the value is a pool of gRPC clients.
	relayClientPools sync.Map
	// relayUrlProvider knows how to retrieve the relay URLs
	relayUrlProvider RelayUrlProvider
}

var _ RelayClient = (*relayClient)(nil)

// NewRelayClient creates a new RelayClient. It keeps a connection to each relay and reuses it for subsequent requests,
// and the connection is lazily instantiated.
func NewRelayClient(
	config *RelayClientConfig,
	logger logging.Logger,
	relayUrlProvider RelayUrlProvider,
) (RelayClient, error) {

	if config == nil {
		return nil, errors.New("nil config")
	}

	if config.MaxGRPCMessageSize == 0 {
		return nil, errors.New("max gRPC message size must be greater than 0")
	}

	connectionPoolSize := config.ConnectionPoolSize
	if connectionPoolSize <= 0 {
		connectionPoolSize = 1
	}

	logger.Info("creating relay client")

	return &relayClient{
		config:             config,
		logger:             logger.With("component", "RelayClient"),
		relayLockProvider:  NewKeyLock[corev2.RelayKey](),
		relayUrlProvider:   relayUrlProvider,
		connectionPoolSize: connectionPoolSize,
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

	hash, err := hashing.HashGetChunksRequest(request)
	if err != nil {
		return fmt.Errorf("failed to hash get chunks request: %v", err)
	}
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
		Timestamp:     uint32(time.Now().Unix()),
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
		Timestamp:     uint32(time.Now().Unix()),
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
	maybeClientPool, ok := c.relayClientPools.Load(key)
	if !ok {
		return nil, fmt.Errorf("no grpc client pool for relay key: %v", key)
	}
	clientPool, ok := maybeClientPool.(*common.GRPCClientPool[relaygrpc.RelayClient])
	if !ok {
		return nil, fmt.Errorf("invalid grpc client for relay key: %v", key)
	}
	return clientPool.GetClient(), nil
}

// initOnceGrpcConnection initializes the GRPC connection for a given relay, and is guaranteed to only be completed
// once per relay. If initialization fails, it will be retried by the next caller.
func (c *relayClient) initOnceGrpcConnection(ctx context.Context, key corev2.RelayKey) error {
	_, alreadyInitialized := c.relayInitializationStatus.Load(key)
	if alreadyInitialized {
		// this is the standard case, where the grpc connection has already been initialized
		return nil
	}

	// In cases were the value hasn't already been initialized, we must acquire a conceptual lock on the relay in
	// question. This allows us to guarantee that a connection with a given relay is only initialized a single time
	releaseKeyLock := c.relayLockProvider.AcquireKeyLock(key)
	defer releaseKeyLock()

	_, alreadyInitialized = c.relayInitializationStatus.Load(key)
	if alreadyInitialized {
		// If we find that the connection was initialized in the time it took to acquire a conceptual lock on the relay,
		// that means that a different caller did the necessary work already
		return nil
	}

	relayUrl, err := c.relayUrlProvider.GetRelayUrl(ctx, key)
	if err != nil {
		return fmt.Errorf("get relay url for key %d: %w", key, err)
	}

	dialOptions := clients.GetGrpcDialOptions(c.config.UseSecureGrpcFlag, c.config.MaxGRPCMessageSize)

	pool, err := common.NewGRPClientPool(
		c.logger,
		relaygrpc.NewRelayClient,
		c.config.ConnectionPoolSize,
		relayUrl,
		dialOptions...)
	if err != nil {
		return fmt.Errorf("failed to create gRPC client pool for relay %d: %w", key, err)
	}

	c.relayClientPools.Store(key, pool)

	// only set the initialization status to true if everything was successful.
	c.relayInitializationStatus.Store(key, true)

	return nil
}

func (c *relayClient) Close() error {
	var errList *multierror.Error

	c.relayClientPools.Range(
		func(k, v interface{}) bool {
			pool, ok := v.(*common.GRPCClientPool[relaygrpc.RelayClient])
			if !ok {
				errList = multierror.Append(errList, fmt.Errorf("invalid connection for relay key: %v", k))
				return true
			}

			if pool != nil {
				err := pool.Close()
				c.relayClientPools.Delete(k)
				if err != nil {
					c.logger.Error("failed to close connection", "err", err)
					errList = multierror.Append(errList, err)
				}
			}
			return true
		})
	return errList.ErrorOrNil()
}
