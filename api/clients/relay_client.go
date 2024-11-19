package clients

import (
	"context"
	"fmt"
	"sync"

	relaygrpc "github.com/Layr-Labs/eigenda/api/grpc/relay"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/grpc"
)

type RelayClientConfig struct {
	Sockets           map[corev2.RelayKey]string
	UseSecureGrpcFlag bool
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

type relayClient struct {
	config *RelayClientConfig

	initOnce map[corev2.RelayKey]*sync.Once
	conns    map[corev2.RelayKey]*grpc.ClientConn
	logger   logging.Logger

	grpcClients map[corev2.RelayKey]relaygrpc.RelayClient
}

var _ RelayClient = (*relayClient)(nil)

// NewRelayClient creates a new RelayClient that connects to the relays specified in the config.
// It keeps a connection to each relay and reuses it for subsequent requests, and the connection is lazily instantiated.
func NewRelayClient(config *RelayClientConfig, logger logging.Logger) (*relayClient, error) {
	if config == nil || len(config.Sockets) > 0 {
		return nil, fmt.Errorf("invalid config: %v", config)
	}

	initOnce := make(map[corev2.RelayKey]*sync.Once)
	conns := make(map[corev2.RelayKey]*grpc.ClientConn)
	grpcClients := make(map[corev2.RelayKey]relaygrpc.RelayClient)
	for key := range config.Sockets {
		initOnce[key] = &sync.Once{}
	}
	return &relayClient{
		config: config,

		initOnce: initOnce,
		conns:    conns,
		logger:   logger,

		grpcClients: grpcClients,
	}, nil
}

func (c *relayClient) GetBlob(ctx context.Context, relayKey corev2.RelayKey, blobKey corev2.BlobKey) ([]byte, error) {
	if err := c.initOnceGrpcConnection(relayKey); err != nil {
		return nil, err
	}

	client, ok := c.grpcClients[relayKey]
	if !ok {
		return nil, fmt.Errorf("no grpc client for relay key: %v", relayKey)
	}

	res, err := client.GetBlob(ctx, &relaygrpc.GetBlobRequest{
		BlobKey: blobKey[:],
	})
	if err != nil {
		return nil, err
	}

	return res.GetBlob(), nil
}

func (c *relayClient) GetChunksByRange(ctx context.Context, relayKey corev2.RelayKey, requests []*ChunkRequestByRange) ([][]byte, error) {
	if len(requests) == 0 {
		return nil, fmt.Errorf("no requests")
	}
	if err := c.initOnceGrpcConnection(relayKey); err != nil {
		return nil, err
	}

	client, ok := c.grpcClients[relayKey]
	if !ok {
		return nil, fmt.Errorf("no grpc client for relay key: %v", relayKey)
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
	res, err := client.GetChunks(ctx, &relaygrpc.GetChunksRequest{
		ChunkRequests: grpcRequests,
	})

	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func (c *relayClient) GetChunksByIndex(ctx context.Context, relayKey corev2.RelayKey, requests []*ChunkRequestByIndex) ([][]byte, error) {
	if len(requests) == 0 {
		return nil, fmt.Errorf("no requests")
	}
	if err := c.initOnceGrpcConnection(relayKey); err != nil {
		return nil, err
	}

	client, ok := c.grpcClients[relayKey]
	if !ok {
		return nil, fmt.Errorf("no grpc client for relay key: %v", relayKey)
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
	res, err := client.GetChunks(ctx, &relaygrpc.GetChunksRequest{
		ChunkRequests: grpcRequests,
	})

	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func (c *relayClient) initOnceGrpcConnection(key corev2.RelayKey) error {
	var initErr error
	c.initOnce[key].Do(func() {
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
		c.conns[key] = conn
		c.grpcClients[key] = relaygrpc.NewRelayClient(conn)
	})
	return initErr
}

func (c *relayClient) Close() error {
	var errList *multierror.Error
	for k, conn := range c.conns {
		if conn != nil {
			err := conn.Close()
			conn = nil
			c.grpcClients[k] = nil
			if err != nil {
				c.logger.Error("failed to close connection", "err", err)
				errList = multierror.Append(errList, err)
			}
		}
	}
	return errList.ErrorOrNil()
}
