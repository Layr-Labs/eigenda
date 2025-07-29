package common

import (
	"fmt"
	"sync/atomic"

	"google.golang.org/grpc"
)

// TODO unit test this

// GRPClientPool manages a pool of one or more gRPC clients.
type GRPClientPool[T any] struct {
	// clients is a slice of gRPC clients of type T.
	clients []T

	// Incremented once per call to GetClient().
	callCount atomic.Uint64
}

// A function that builds a gRPC client of type T.
type GRPCClientBuilder[T any] func(grpc.ClientConnInterface) T

// NewGRPClientManager creates a new GRPClientPool with the specified client builder and size.
func NewGRPClientManager[T any](
	clientBuilder GRPCClientBuilder[T],
	poolSize int,
	url string,
	dialOptions ...grpc.DialOption,
) (*GRPClientPool[T], error) {

	// Create the clients up front.
	clients := make([]T, 0, poolSize)
	for i := 0; i < poolSize; i++ {
		conn, err := grpc.NewClient(url, dialOptions...)
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC client connection to %s: %w", url, err)
		}

		client := clientBuilder(conn)
		if client == nil {
			return nil, fmt.Errorf("client builder returned nil for gRPC client")
		}

		clients = append(clients, client)
	}

	return &GRPClientPool[T]{
		callCount: atomic.Uint64{},
		clients:   make([]T, poolSize),
	}, nil
}

// GetClient returns a gRPC client of type T. If this client manager maintains a pool of clients, then it will choose
// one from the pool to return.
func (m *GRPClientPool[T]) GetClient() (T, error) {
	var client T
	if len(m.clients) == 0 {
		client = m.clients[0]
	} else {
		index := m.callCount.Add(1)
		client = m.clients[index%uint64(len(m.clients))]
	}

	return client, nil
}
