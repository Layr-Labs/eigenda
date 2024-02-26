package dataapi

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	mutex sync.Mutex
)

func (s *server) getServiceAvailability(ctx context.Context, hosts []string) ([]*ServiceAvailability, error) {

	if hosts == nil {
		return nil, fmt.Errorf("hostnames cannot be nil")
	}

	availabilityStatuses := make([]*ServiceAvailability, len(hosts))

	for i, host := range hosts {
		pool, ok := s.getClientPool(host)
		if !ok {
			return nil, fmt.Errorf("Invalid hostname: %s", host)
		}
		conn, err := getClientConn(pool)
		if err != nil {
			return nil, fmt.Errorf("Error getting client connection: %v", err)
		}
		defer conn.Close()
		client := grpc_health_v1.NewHealthClient(conn)
		response, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			return nil, fmt.Errorf("Error checking health of service: %v", err)
		}

		availabilityStatus := &ServiceAvailability{
			ServiceName:   host,
			ServiceStatus: response.Status.String(),
		}
		availabilityStatuses[i] = availabilityStatus
		// Return connection back to pool
		putClientConn(conn, pool)

	}
	return availabilityStatuses, nil
}

// Initializes the client pools for the server
func (s *server) InitGRPCClientPools(poolSize int) error {
	mutex.Lock()
	defer mutex.Unlock()

	var err error
	if s.clientPools == nil {
		s.clientPools = make(map[string]*ClientPool)
	}
	s.clientPools[s.disperserHostName], err = newClientPool(poolSize, s.disperserHostName)
	if err != nil {
		return err
	}
	s.clientPools[s.churnerHostName], err = newClientPool(poolSize, s.churnerHostName)
	if err != nil {
		return err
	}

	return nil
}

// newClientPool creates a client pool with prewarmed connections
func newClientPool(size int, serverAddr string) (*ClientPool, error) {
	pool := &ClientPool{
		clients: make(chan *grpc.ClientConn, size),
	}
	for i := 0; i < size; i++ {
		conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			return nil, err
		}
		pool.clients <- conn
	}
	return pool, nil
}

// getClientPool retrieves a client pool for a given service hostname
func (s *server) getClientPool(serviceHostName string) (*ClientPool, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	pool, ok := s.clientPools[serviceHostName]
	return pool, ok
}

// Get retrieves a gRPC client connection from the pool.
func getClientConn(pool *ClientPool) (*grpc.ClientConn, error) {
	select {
	case conn := <-pool.clients:
		return conn, nil
	default:
		// Handle the scenario when no connections are available in the pool.
		return nil, fmt.Errorf("no available connections in the pool")
	}
}

// puts a gRPC client connection back into the pool.
func putClientConn(conn *grpc.ClientConn, pool *ClientPool) {
	pool.clients <- conn // It's a good idea to check if the connection is still healthy before returning.
}
