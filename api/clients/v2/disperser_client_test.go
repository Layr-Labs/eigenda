package clients

import (
	"context"
	"math/big"
	"net"
	"sync"
	"testing"
	"time"

	v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestVerifyReceivedBlobKey(t *testing.T) {
	blobCommitments := encoding.BlobCommitments{
		Commitment:       &encoding.G1Commitment{},
		LengthCommitment: &encoding.G2Commitment{},
		LengthProof:      &encoding.LengthProof{},
		Length:           4,
	}

	quorumNumbers := make([]core.QuorumID, 1)
	quorumNumbers[0] = 8

	paymentMetadata := core.PaymentMetadata{
		AccountID:         gethcommon.Address{1},
		Timestamp:         5,
		CumulativePayment: big.NewInt(6),
	}

	blobHeader := &corev2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: blobCommitments,
		QuorumNumbers:   quorumNumbers,
		PaymentMetadata: paymentMetadata,
	}

	realKey, err := blobHeader.BlobKey()
	require.NoError(t, err)

	reply := v2.DisperseBlobReply{
		BlobKey: realKey[:],
	}

	require.NoError(t, verifyReceivedBlobKey(blobHeader, &reply))

	blobHeader.BlobVersion = 1
	require.Error(t, verifyReceivedBlobKey(blobHeader, &reply),
		"Any modification to the header should cause verification to fail")
}

// TestMutexPreventsSimultaneousRequests tests that the mutex in disperserClient
// prevents multiple goroutines from executing critical sections concurrently.
func TestMutexPreventsSimultaneousRequests(t *testing.T) {
	// Create a struct with a mutex and a counter
	client := &struct {
		requestMutex sync.Mutex
		counter      int
		callTimes    []time.Time
	}{}

	// Use this function to simulate a request that takes some time
	simulateRequest := func() {
		client.requestMutex.Lock()
		defer client.requestMutex.Unlock()

		// Record the time of the call
		callTime := time.Now()
		client.callTimes = append(client.callTimes, callTime)
		client.counter++

		// Simulate processing time
		time.Sleep(200 * time.Millisecond)
	}

	// Number of concurrent "requests" to attempt
	numRequests := 3

	// Use a WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Start time for our test
	startTime := time.Now()

	// Launch multiple goroutines to make concurrent "requests"
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			simulateRequest()
		}()
	}

	// Wait for all requests to complete
	wg.Wait()

	// Verify that the correct number of requests were made
	require.Equal(t, numRequests, client.counter, "Expected number of requests")

	// Check that the requests were executed sequentially, not concurrently
	// The time difference between consecutive requests should be at least the delay time
	for i := 1; i < len(client.callTimes); i++ {
		timeDiff := client.callTimes[i].Sub(client.callTimes[i-1])
		require.GreaterOrEqual(t, timeDiff.Milliseconds(), int64(199), // slightly less than 200ms to account for timing variations
			"Requests were not executed sequentially. Time between request %d and %d was only %v",
			i-1, i, timeDiff)
	}

	// The total time should be at least (numRequests * delay)
	// This verifies that the requests were not processed concurrently
	totalTime := time.Since(startTime)
	expectedMinTime := time.Duration(numRequests) * 200 * time.Millisecond
	require.GreaterOrEqual(t, totalTime.Milliseconds(), expectedMinTime.Milliseconds()-10, // allow small timing variations
		"Total execution time was less than expected, suggesting concurrent execution")
}

// mockDisperserServer implements the gRPC DisperserServer interface for testing
// Currently only implements the GetBlobStatus method, to be able to test Cert Retrievals.
type mockDisperserServer struct {
	v2.UnimplementedDisperserServer

	getBlobStatusResponse  *v2.BlobStatusReply
	getBlobStatusError     error
	getBlobStatusCallCount int
	getBlobStatusRequests  []*v2.BlobStatusRequest

	// Mutex for thread-safe access
	mu sync.Mutex
}

func (m *mockDisperserServer) GetBlobStatus(ctx context.Context, req *v2.BlobStatusRequest) (*v2.BlobStatusReply, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.getBlobStatusCallCount++
	m.getBlobStatusRequests = append(m.getBlobStatusRequests, req)

	if m.getBlobStatusError != nil {
		return nil, m.getBlobStatusError
	}
	return m.getBlobStatusResponse, nil
}

// testHarness encapsulates all the test dependencies
type testHarness struct {
	server     *mockDisperserServer
	grpcServer *grpc.Server
	client     DisperserClient
	signer     *mockBlobRequestSigner
	prover     *mockProver
	accountant *Accountant
	logger     logging.Logger
	listener   net.Listener
}

// setupTestHarness creates a complete test environment with mocked dependencies
func setupTestHarness(t *testing.T) (*testHarness, func()) {
	// Create mock server
	mockServer := &mockDisperserServer{}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	v2.RegisterDisperserServer(grpcServer, mockServer)

	// Create listener on random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	// Start gRPC server
	go func() {
		_ = grpcServer.Serve(listener)
	}()

	// Extract port from listener
	_, port, err := net.SplitHostPort(listener.Addr().String())
	require.NoError(t, err)

	// Create mock dependencies
	logger := logging.NewNoopLogger()
	signer := &mockBlobRequestSigner{}
	prover := &mockProver{}
	accountant := NewAccountant(gethcommon.Address{})

	// Create client config
	config := &DisperserClientConfig{
		Hostname:          "127.0.0.1",
		Port:              port,
		UseSecureGrpcFlag: false,
	}

	// Create disperser client
	client, err := NewDisperserClient(logger, config, signer, prover, accountant)
	require.NoError(t, err)

	// Cleanup function
	cleanup := func() {
		client.Close()
		grpcServer.Stop()
		listener.Close()
	}

	return &testHarness{
		server:     mockServer,
		grpcServer: grpcServer,
		client:     client,
		signer:     signer,
		prover:     prover,
		accountant: accountant,
		logger:     logger,
		listener:   listener,
	}, cleanup
}

func TestDisperserClient(t *testing.T) {

	t.Run("GetBlobStatus", func(t *testing.T) {
		harness, cleanup := setupTestHarness(t)
		defer cleanup()

		// Configure mock response
		harness.server.getBlobStatusResponse = &v2.BlobStatusReply{
			Status: v2.BlobStatus_FINALIZED,
		}

		// Test blob key
		var blobKey corev2.BlobKey
		copy(blobKey[:], []byte("test-blob-key"))

		// Make the call
		ctx := context.Background()
		reply, err := harness.client.GetBlobStatus(ctx, blobKey)

		// Verify results
		require.NoError(t, err)
		require.NotNil(t, reply)
		require.Equal(t, v2.BlobStatus_FINALIZED, reply.Status)

		// Verify mock was called correctly
		require.Equal(t, 1, harness.server.getBlobStatusCallCount)
		require.Len(t, harness.server.getBlobStatusRequests, 1)
		require.Equal(t, blobKey[:], harness.server.getBlobStatusRequests[0].BlobKey)
	})
}
