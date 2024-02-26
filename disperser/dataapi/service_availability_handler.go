package dataapi

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var (
	mutex sync.Mutex
)

func (s *server) getServiceAvailability(ctx context.Context, services []string) ([]*ServiceAvailability, error) {

	if services == nil {
		return nil, fmt.Errorf("services cannot be nil")
	}

	availabilityStatuses := make([]*ServiceAvailability, len(services))

	for i, serviceName := range services {
		pool, ok := s.getClientPool(serviceName)
		if !ok {
			return nil, fmt.Errorf("Invalid ServiceName: %s", serviceName)
		}
		conn, err := getClientConn(pool)
		if err != nil {
			return nil, fmt.Errorf("Error getting client connection: %v", err)
		}
		var availabilityStatus *ServiceAvailability
		client := grpc_health_v1.NewHealthClient(conn)
		response, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})

		if err != nil {
			availabilityStatus = &ServiceAvailability{
				ServiceName:   serviceName,
				ServiceStatus: grpc_health_v1.HealthCheckResponse_NOT_SERVING.String(),
			}
			availabilityStatuses[i] = availabilityStatus
		} else {

			availabilityStatus = &ServiceAvailability{
				ServiceName:   serviceName,
				ServiceStatus: response.Status.String(),
			}
			availabilityStatuses[i] = availabilityStatus
		}

		if availabilityStatuses[i].ServiceStatus == grpc_health_v1.HealthCheckResponse_SERVING.String() {
			pool.clients <- conn
		} else {
			conn.Close()
		}
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

	// Register Disperser and Churner client pools
	// These are the only public services that the dataapi server will connect to
	s.clientPools["Disperser"], err = newClientPool(poolSize, s.disperserHostName)
	if err != nil {
		return err
	}
	s.clientPools["Churner"], err = newClientPool(poolSize, s.churnerHostName)
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
		// Context with timeout for each dial attempt
		// Attempt to create a connection to the server only for certain error codes
		conn, err := createConnWithRetry(serverAddr)
		if err != nil {
			return nil, err
		}
		pool.clients <- conn
	}
	return pool, nil
}

// getClientPool retrieves a client pool for a given service hostname
func (s *server) getClientPool(serviceName string) (*ClientPool, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	pool, ok := s.clientPools[serviceName]
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

func createConnWithRetry(serverAddr string) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 30 * time.Second // Maximum time for retrying.

	// The operation to retry: gRPC dial with a timeout
	operation := func() error {
		var err error
		// Context with timeout for each dial attempt
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err = grpc.DialContext(ctx, serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				// Setup Retry for Transient Errors
				switch st.Code() {
				case codes.DeadlineExceeded: // No data transmitted before deadline expires
					log.Printf("Retrying due to deadline exceeded: %v", err)
					return err
				case codes.Unavailable: // Connection break after some data is transmitted
					log.Printf("Retrying due to service unavailability: %v", err)
					return err
				}
			}

			// Permanent error, don't retry
			return backoff.Permanent(err)
		}
		return nil
	}

	// Execute the operation with retry
	err := backoff.Retry(operation, b)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s with retry: %w", serverAddr, err)
	}
	return conn, nil
}
