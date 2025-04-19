package clients

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

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
	UseSecureGrpcFlag      bool
	MaxGRPCMessageSize     uint
	OperatorID             *core.OperatorID
	MessageSigner          MessageSigner
	NumConnectionsPerRelay int // Number of gRPC connections to open per relay for round-robin load balancing
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
	relayLockProvider *relay.KeyLock[corev2.RelayKey]
	// relayInitializationStatus maps relay key to a bool `map[corev2.RelayKey]bool`
	// the boolean value indicates whether the connection to that relay has been initialized
	relayInitializationStatus sync.Map
	// clientConnections maps relay key to a slice of gRPC connections: `map[corev2.RelayKey][]*grpc.ClientConn`
	// this map is maintained so that connections can be closed in Close
	clientConnections sync.Map
	// grpcRelayClients maps relay key to a slice of gRPC clients: `map[corev2.RelayKey][]relaygrpc.RelayClient`
	// these grpc relay clients are used to communicate with individual relays
	grpcRelayClients sync.Map
	// clientCounters maps relay key to an atomic counter for round-robin selection: `map[corev2.RelayKey]*atomic.Uint32`
	clientCounters sync.Map
	// relayUrlProvider knows how to retrieve the relay URLs
	relayUrlProvider relay.RelayUrlProvider
}

var _ RelayClient = (*relayClient)(nil)

// NewRelayClient creates a new RelayClient. It keeps connections to each relay and reuses them for subsequent requests,
// and the connections are lazily instantiated.
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

	// Default to 4 connections per relay if not specified
	if config.NumConnectionsPerRelay <= 0 {
		config.NumConnectionsPerRelay = 4
	}

	logger.Info("creating relay client", "connectionsPerRelay", config.NumConnectionsPerRelay)

	return &relayClient{
		config:            config,
		logger:            logger.With("component", "RelayClient"),
		relayLockProvider: relay.NewKeyLock[corev2.RelayKey](),
		relayUrlProvider:  relayUrlProvider,
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

	// Break large requests into batches of maxRequestsPerRPC
	const maxRequestsPerRPC = 64
	if len(requests) > maxRequestsPerRPC {
		c.logger.Debug("Breaking large request into batches", "relayKey", relayKey, "totalRequests", len(requests), "batchSize", maxRequestsPerRPC)

		var allData [][]byte
		for start := 0; start < len(requests); start += maxRequestsPerRPC {
			end := start + maxRequestsPerRPC
			if end > len(requests) {
				end = len(requests)
			}

			batchData, err := c.GetChunksByRange(ctx, relayKey, requests[start:end])
			if err != nil {
				return nil, err
			}
			allData = append(allData, batchData...)
		}
		return allData, nil
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

	c.logger.Debug("Processing chunk requests", "relayKey", relayKey, "count", len(grpcRequests))

	startTime := time.Now()
	res, err := client.GetChunks(ctx, request)
	duration := time.Since(startTime)

	if err != nil {
		c.logger.Error("Failed to get chunks by range", "relayKey", relayKey, "count", len(grpcRequests), "duration", duration, "error", err)
		return nil, err
	}

	c.logger.Debug("Completed GetChunks by range", "relayKey", relayKey, "count", len(grpcRequests), "duration", duration)

	return res.GetData(), nil
}

func (c *relayClient) GetChunksByIndex(
	ctx context.Context,
	relayKey corev2.RelayKey,
	requests []*ChunkRequestByIndex) ([][]byte, error) {

	if len(requests) == 0 {
		return nil, fmt.Errorf("no requests")
	}

	// Break large requests into batches of maxRequestsPerRPC
	const maxRequestsPerRPC = 64
	if len(requests) > maxRequestsPerRPC {
		c.logger.Debug("Breaking large request into batches", "relayKey", relayKey, "totalRequests", len(requests), "batchSize", maxRequestsPerRPC)

		var allData [][]byte
		for start := 0; start < len(requests); start += maxRequestsPerRPC {
			end := start + maxRequestsPerRPC
			if end > len(requests) {
				end = len(requests)
			}

			batchData, err := c.GetChunksByIndex(ctx, relayKey, requests[start:end])
			if err != nil {
				return nil, err
			}
			allData = append(allData, batchData...)
		}
		return allData, nil
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

	c.logger.Debug("Processing chunk requests by index", "relayKey", relayKey, "count", len(grpcRequests))

	startTime := time.Now()
	res, err := client.GetChunks(ctx, request)
	duration := time.Since(startTime)

	if err != nil {
		c.logger.Error("Failed to get chunks by index", "relayKey", relayKey, "count", len(grpcRequests), "duration", duration, "error", err)
		return nil, err
	}

	c.logger.Debug("Completed GetChunks by index", "relayKey", relayKey, "count", len(grpcRequests), "duration", duration)

	return res.GetData(), nil
}

// getClient gets a grpc relay client from the pool using round-robin selection
func (c *relayClient) getClient(ctx context.Context, key corev2.RelayKey) (relaygrpc.RelayClient, error) {
	if err := c.initOnceGrpcConnection(ctx, key); err != nil {
		return nil, fmt.Errorf("init grpc connection for key %d: %w", key, err)
	}

	maybeClientsSlice, ok := c.grpcRelayClients.Load(key)
	if !ok {
		return nil, fmt.Errorf("no grpc clients for relay key: %v", key)
	}

	clientsSlice, ok := maybeClientsSlice.([]relaygrpc.RelayClient)
	if !ok {
		return nil, fmt.Errorf("invalid grpc clients for relay key: %v", key)
	}

	if len(clientsSlice) == 0 {
		return nil, fmt.Errorf("empty grpc clients slice for relay key: %v", key)
	}

	// Get counter for round-robin selection
	maybeCounter, _ := c.clientCounters.Load(key)
	counter, ok := maybeCounter.(*atomic.Uint32)
	if !ok {
		return nil, fmt.Errorf("invalid counter for relay key: %v", key)
	}

	// Select client using round-robin
	idx := int(counter.Add(1) % uint32(len(clientsSlice)))
	return clientsSlice[idx], nil
}

// initOnceGrpcConnection initializes multiple GRPC connections for a given relay, and is guaranteed to only be completed
// once per relay. If initialization fails, it will be retried by the next caller.
func (c *relayClient) initOnceGrpcConnection(ctx context.Context, key corev2.RelayKey) error {
	_, alreadyInitialized := c.relayInitializationStatus.Load(key)
	if alreadyInitialized {
		// this is the standard case, where the grpc connections have already been initialized
		return nil
	}

	// In cases were the value hasn't already been initialized, we must acquire a conceptual lock on the relay in
	// question. This allows us to guarantee that connections with a given relay are only initialized a single time
	releaseKeyLock := c.relayLockProvider.AcquireKeyLock(key)
	defer releaseKeyLock()

	_, alreadyInitialized = c.relayInitializationStatus.Load(key)
	if alreadyInitialized {
		// If we find that the connections were initialized in the time it took to acquire a conceptual lock on the relay,
		// that means that a different caller did the necessary work already
		return nil
	}

	relayUrl, err := c.relayUrlProvider.GetRelayUrl(ctx, key)
	if err != nil {
		return fmt.Errorf("get relay url for key %d: %w", key, err)
	}

	// Create multiple connections and clients for this relay
	conns := make([]*grpc.ClientConn, c.config.NumConnectionsPerRelay)
	clients := make([]relaygrpc.RelayClient, c.config.NumConnectionsPerRelay)

	dialOptions := getGrpcDialOptions(c.config.UseSecureGrpcFlag, c.config.MaxGRPCMessageSize)

	for i := 0; i < c.config.NumConnectionsPerRelay; i++ {
		// Add unique user agent to each connection
		connDialOptions := append(dialOptions,
			grpc.WithUserAgent(fmt.Sprintf("relay-%d-conn-%d", key, i)))

		conn, err := grpc.NewClient(relayUrl, connDialOptions...)
		if err != nil {
			// Close any connections that were already created
			for j := 0; j < i; j++ {
				if conns[j] != nil {
					_ = conns[j].Close()
				}
			}
			return fmt.Errorf("create grpc client for key %d connection %d: %w", key, i, err)
		}

		conns[i] = conn
		clients[i] = relaygrpc.NewRelayClient(conn)
	}

	c.clientConnections.Store(key, conns)
	c.grpcRelayClients.Store(key, clients)
	c.clientCounters.Store(key, &atomic.Uint32{})

	// only set the initialization status to true if everything was successful.
	c.relayInitializationStatus.Store(key, true)

	c.logger.Info("Created multiple grpc connections", "relayKey", key, "numConnections", c.config.NumConnectionsPerRelay)
	return nil
}

func (c *relayClient) Close() error {
	var errList *multierror.Error
	c.clientConnections.Range(
		func(k, v interface{}) bool {
			conns, ok := v.([]*grpc.ClientConn)
			if !ok {
				errList = multierror.Append(errList, fmt.Errorf("invalid connections for relay key: %v", k))
				return true
			}

			for i, conn := range conns {
				if conn != nil {
					err := conn.Close()
					if err != nil {
						c.logger.Error("failed to close connection", "relayKey", k, "connIndex", i, "err", err)
						errList = multierror.Append(errList, err)
					}
				}
			}

			c.clientConnections.Delete(k)
			c.grpcRelayClients.Delete(k)
			c.clientCounters.Delete(k)

			return true
		})
	return errList.ErrorOrNil()
}
