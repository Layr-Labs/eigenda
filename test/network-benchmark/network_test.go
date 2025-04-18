package network_benchmark

import (
	"fmt"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/relay"
	testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/docker/go-units"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const bufSize = units.MiB

var dataPerTransfer = int64(1 * units.MiB)
var totalDataToTransfer = int64(100 * units.GiB)
var parallelism = 8

// TODO claude: respect these addresses in the test
var serverAddress = "localhost:50051"
var clientAddress = "localhost:50052"

func worker(
	t *testing.T,
	client TestClient,
	dataSize int64,
	seed int64,
	randomData *reusableRandomness,
	totalDataToTransfer int64,
	totalLatency *atomic.Uint64,
	totalDataTransferred *atomic.Uint64,
	transferCount *atomic.Uint64,
	finishedChan chan struct{},
) {
	rand := testrandom.NewTestRandom(seed)
	var dataTransferred int64

	for dataTransferred < totalDataToTransfer {
		// Generate a random seed for deterministic but varying data
		seed := rand.Int63()

		// Measure latency for this request
		requestStart := time.Now()
		data, err := client.GetData(int64(dataSize), seed)
		if err != nil {
			t.Fatalf("Failed to get data: %v", err)
		}
		requestLatency := time.Since(requestStart)

		// Regenerate the data using the same seed and verify it matches
		expectedData := randomData.getData(int64(dataSize), seed)
		if len(data) != len(expectedData) {
			t.Fatalf("Data length mismatch: got %d, expected %d", len(data), len(expectedData))
		}

		// Compare the data
		for i := 0; i < len(data); i++ {
			if data[i] != expectedData[i] {
				t.Fatalf("Data mismatch at index %d: got %d, expected %d", i, data[i], expectedData[i])
			}
		}

		// Update metrics
		totalLatency.Add(uint64(requestLatency.Nanoseconds()))
		dataSize := int64(len(data))
		dataTransferred += dataSize
		totalDataTransferred.Add(uint64(dataSize))
		transferCount.Add(1) // Count each successful transfer
	}

	finishedChan <- struct{}{}
}

func throughputTest(t *testing.T, server TestServer, clients []TestClient) {
	rand := testrandom.NewTestRandom()
	randomData := newReusableRandomness(units.GiB, rand.Int63())
	server.SetRandomData(randomData)

	start := time.Now()

	totalLatency := &atomic.Uint64{}
	totalDataTransferred := &atomic.Uint64{}
	transferCount := &atomic.Uint64{}

	finishedChan := make(chan struct{}, len(clients))

	for i := 0; i < len(clients); i++ {
		go worker(
			t,
			clients[i],
			dataPerTransfer,
			rand.Int63(),
			randomData,
			totalDataToTransfer/int64(len(clients)), // Divide total between workers
			totalLatency,
			totalDataTransferred,
			transferCount,
			finishedChan)
	}

	// Periodically print status updates
	statusTicker := time.NewTicker(1 * time.Second)
	defer statusTicker.Stop()

	// Use a separate goroutine to display stats while the workers are running
	done := false
	go func() {
		for !done {
			select {
			case <-statusTicker.C:
				elapsedSoFar := time.Since(start)
				currentLatencyNs := totalLatency.Load()
				currentTotal := totalDataTransferred.Load()
				currentCount := transferCount.Load()

				// Calculate current metrics
				throughputSoFar := float64(currentTotal) / elapsedSoFar.Seconds()
				var avgLatencySoFar time.Duration
				if currentCount > 0 {
					avgLatencySoFar = time.Duration(currentLatencyNs / currentCount)
				}

				fmt.Printf("[%s] Workers: %d, Transferred: %s (%.2f%%), Avg Latency: %s, Throughput: %s/s\n",
					elapsedSoFar.Round(time.Second),
					len(clients),
					units.BytesSize(float64(currentTotal)),
					100.0*float64(currentTotal)/float64(totalDataToTransfer),
					avgLatencySoFar,
					units.BytesSize(throughputSoFar))
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Wait for all workers to finish
	for i := 0; i < len(clients); i++ {
		<-finishedChan
	}

	elapsed := time.Since(start)
	done = true // Signal status printer to exit

	// Calculate final metrics from atomic values
	finalTotalNs := totalLatency.Load()
	finalTotal := totalDataTransferred.Load()
	finalCount := transferCount.Load()

	var avgLatency time.Duration
	if finalCount > 0 {
		avgLatency = time.Duration(finalTotalNs / finalCount)
	}
	throughput := float64(finalTotal) / elapsed.Seconds()

	// Print the benchmark results
	fmt.Printf("\n--- Benchmark Results ---\n")
	fmt.Printf("Parallelism: %d workers\n", len(clients))
	fmt.Printf("Total time: %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("Average latency: %s\n", avgLatency)
	fmt.Printf("Average throughput: %s/s\n", units.BytesSize(throughput))
	fmt.Printf("Total data transferred: %s\n", units.BytesSize(float64(finalTotal)))
	fmt.Printf("Number of transfers: %d\n", finalCount)
}

func TestProtobufThroughput(t *testing.T) {
	// Use real network addresses instead of in-memory bufconn
	// Start gRPC server with the real server address
	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	protobufServer := NewProtobufServer()
	relay.RegisterThroughputTestServer(server, protobufServer)

	go func() {
		if err := server.Serve(lis); err != nil {
			t.Fatalf("Failed to serve: %v", err)
		}
	}()
	defer server.Stop()

	// Set up clients that connect to the real server address
	clients := make([]TestClient, parallelism)
	for i := 0; i < parallelism; i++ {
		// Create a new connection for each client
		conn, err := grpc.Dial(
			serverAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			t.Fatalf("Failed to dial server: %v", err)
		}
		defer conn.Close()

		grpcClient := relay.NewThroughputTestClient(conn)
		clients[i] = newProtobufClient(grpcClient)
	}

	// Run the benchmark test
	throughputTest(t, protobufServer.(TestServer), clients)
}
