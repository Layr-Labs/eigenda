package common

import (
	"fmt"
	"sync/atomic"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
)

// TODO unit test this

// A function that builds a gRPC client of type T.
type GRPCClientBuilder[T any] func(grpc.ClientConnInterface) T

// GRPCClientPool manages a pool of one or more gRPC clients.
type GRPCClientPool[T any] struct {
	// clients is a slice of gRPC clients of type T.
	clients []T

	// connections is a slice of gRPC client connections. We need to track this in order to be able to close the
	// connections when the pool is no longer needed.
	connections []*grpc.ClientConn

	// Incremented once per call to GetClient().
	callCount atomic.Uint64
}

// Creates a new GRPCClientPool with the specified client builder and size.
func NewGRPClientPool[T any](
	logger logging.Logger,
	clientBuilder GRPCClientBuilder[T],
	poolSize int,
	url string,
	dialOptions ...grpc.DialOption,
) (*GRPCClientPool[T], error) {

	if poolSize <= 0 {
		poolSize = 1
	}

	// Create the clients up front.
	connections := make([]*grpc.ClientConn, 0, poolSize)
	clients := make([]T, 0, poolSize)
	for i := 0; i < poolSize; i++ {
		conn, err := grpc.NewClient(url, dialOptions...)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC client connection to %s: %w", url, err)
		}
		connections = append(connections, conn)

		client := clientBuilder(conn)
		clients = append(clients, client)
	}

	clientType := fmt.Sprintf("%T", clients[0])
	logger.Infof("Creating gRPC client pool of size %d for %s with URL %s", poolSize, clientType, url)

	return &GRPCClientPool[T]{
		callCount:   atomic.Uint64{},
		connections: connections,
		clients:     clients,
	}, nil
}

// GetClient returns a gRPC client of type T. If this client manager maintains a pool of clients, then it will choose
// one from the pool to return.
func (m *GRPCClientPool[T]) GetClient() T {
	var client T
	if len(m.clients) == 1 {
		client = m.clients[0]
	} else {
		index := m.callCount.Add(1)
		client = m.clients[index%uint64(len(m.clients))]
	}

	return client
}

// Close closes all gRPC client connections in the pool and releases resources.
func (m *GRPCClientPool[T]) Close() error {
	var err error
	for _, conn := range m.connections {
		if closeErr := conn.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close gRPC client connection: %w", closeErr)
		}
	}

	m.connections = nil
	m.clients = nil

	return err
}
