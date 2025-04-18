package network_benchmark

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/relay"
	testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/docker/go-units"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var dataPerTransfer = int64(1 * units.MiB)
var totalDataToTransfer = int64(100 * units.GiB)
var parallelism = 8

// Server and client addresses for the benchmark tests
var serverAddress = "localhost:50051"
var clientAddress = "localhost:50052"

var seed = int64(1337)

// Command-line flags to control test behavior
var runServer bool
var runClient bool

// Initialize flags in init()
func init() {
	if os.Getenv("RUN_SERVER") == "true" {
		runServer = true
	}
	if os.Getenv("RUN_CLIENT") == "true" {
		runClient = true
	}

	if runServer {
		fmt.Println("Running server...")
	} else {
		fmt.Println("Not running server...")
	}
	if runClient {
		fmt.Println("Running client...")
	} else {
		fmt.Println("Not running client...")
	}
}

// waitForCtrlC blocks until the user presses Ctrl+C
func waitForCtrlC() {
	fmt.Println("Server running. Press Ctrl+C to stop...")

	// Set up channel to receive interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for signal
	<-sigChan
	fmt.Println("\nReceived interrupt signal, shutting down...")
}

func worker(
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
			panic(fmt.Sprintf("Failed to get data: %v", err))
		}
		requestLatency := time.Since(requestStart)

		// Regenerate the data using the same seed and verify it matches
		expectedData := randomData.getData(int64(dataSize), seed)
		if len(data) != len(expectedData) {
			panic(fmt.Sprintf("Data length mismatch: %d vs %d", len(data), len(expectedData)))
		}

		// Compare the data
		for i := 0; i < len(data); i++ {
			if data[i] != expectedData[i] {
				panic(fmt.Sprintf("Data mismatch: %d vs %d", data[i], expectedData[i]))
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

func throughputTest(server TestServer, clients []TestClient) {
	rand := testrandom.NewTestRandom(seed)
	randomData := newReusableRandomness(units.GiB, seed)
	server.SetRandomData(randomData)

	start := time.Now()

	totalLatency := &atomic.Uint64{}
	totalDataTransferred := &atomic.Uint64{}
	transferCount := &atomic.Uint64{}

	finishedChan := make(chan struct{}, len(clients))

	for i := 0; i < len(clients); i++ {
		go worker(
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
	// Parse flags if they haven't been parsed yet (this happens when running the test directly)
	if !flag.Parsed() {
		flag.Parse()
	}

	var server *grpc.Server
	var pbufServer TestServer

	// Start the server if runServer flag is true
	if runServer {
		fmt.Printf("Starting Protobuf server on %s\n", serverAddress)

		// Use real network addresses instead of in-memory bufconn
		lis, err := net.Listen("tcp", serverAddress)
		if err != nil {
			t.Fatalf("Failed to listen: %v", err)
		}

		server = grpc.NewServer()
		pbufServer = NewProtobufServer()

		randomData := newReusableRandomness(units.GiB, seed)
		pbufServer.(TestServer).SetRandomData(randomData)
		fmt.Printf("random data initialized\n")

		// Type assertion to register with gRPC
		grpcServer, ok := pbufServer.(*protobufServer)
		if !ok {
			t.Fatalf("Failed to cast to protobufServer type")
		}
		relay.RegisterThroughputTestServer(server, grpcServer)

		// If only running the server and not the client, run the server and block until Ctrl+C
		if !runClient {
			fmt.Printf("Running Protobuf server only mode on %s\n", serverAddress)

			// Start the server in a goroutine
			go func() {
				if err := server.Serve(lis); err != nil {
					if err != grpc.ErrServerStopped {
						t.Errorf("Failed to serve: %v", err)
					}
				}
			}()

			// Wait for Ctrl+C
			waitForCtrlC()

			// Stop the server before exiting
			server.Stop()
			return
		}

		// If also running the client, start server in background
		go func() {
			if err := server.Serve(lis); err != nil {
				if err != grpc.ErrServerStopped {
					t.Errorf("Failed to serve: %v", err)
				}
			}
		}()
		defer server.Stop()
	} else {
		// Create a mock server for the client-only case
		pbufServer = NewProtobufServer()
	}

	// Run the client part if runClient flag is true
	if runClient {
		fmt.Printf("Starting Protobuf clients connecting to %s\n", serverAddress)

		// Set up clients that connect to the server address
		clients := make([]TestClient, parallelism)
		for i := 0; i < parallelism; i++ {
			// Create a new connection for each client
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Using grpc.DialContext with nolint directive to suppress the deprecation warning
			// We're using this method instead of NewClient for backwards compatibility
			//nolint:staticcheck
			conn, err := grpc.DialContext(
				ctx,
				serverAddress,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				//nolint:staticcheck
				grpc.WithBlock(),
			)
			if err != nil {
				t.Fatalf("Failed to dial server: %v", err)
			}
			defer conn.Close()

			grpcClient := relay.NewThroughputTestClient(conn)
			clients[i] = newProtobufClient(grpcClient)
		}

		// Run the benchmark test
		throughputTest(pbufServer, clients)
	}
}

func TestSocketThroughput(t *testing.T) {
	// Parse flags if they haven't been parsed yet
	if !flag.Parsed() {
		flag.Parse()
	}

	var server TestServer
	var err error

	// Start the server if runServer flag is true
	if runServer {
		fmt.Printf("Starting Socket server on %s\n", clientAddress)

		// Create a socket server listening on the clientAddress
		server, err = NewSocketServer(clientAddress)
		if err != nil {
			t.Fatalf("Failed to create socket server: %v", err)
		}

		randomData := newReusableRandomness(units.GiB, seed)
		server.SetRandomData(randomData)
		fmt.Printf("random data initialized\n")

		// If only running the server and not the client, block until Ctrl+C
		if !runClient {
			fmt.Printf("Running Socket server only mode on %s\n", clientAddress)

			// Wait for Ctrl+C
			waitForCtrlC()

			// Cleanup before exiting
			socketServerImpl, ok := server.(*socketServer)
			if !ok {
				t.Errorf("Failed to cast to socketServer type")
			} else {
				err = socketServerImpl.Close()
				if err != nil {
					t.Errorf("Failed to close socket server: %v", err)
				}
			}
			return
		}

		// Ensure server resources are cleaned up after the test when both server and client are running
		defer func() {
			socketServerImpl, ok := server.(*socketServer)
			if !ok {
				t.Errorf("Failed to cast to socketServer type")
			} else {
				err = socketServerImpl.Close()
				if err != nil {
					t.Errorf("Failed to close socket server: %v", err)
				}
			}
		}()

		// Wait for the server to start
		time.Sleep(100 * time.Millisecond)
	} else {
		// Create a mock server for the client-only case
		mockServer := &protobufServer{
			randomData: newReusableRandomness(units.MiB, seed),
		}
		server = mockServer
	}

	// Run the client part if runClient flag is true
	if runClient {
		fmt.Printf("Starting Socket clients connecting to %s\n", clientAddress)

		// Create socket clients connecting to the server
		clients := make([]TestClient, parallelism)
		for i := 0; i < parallelism; i++ {
			client, err := NewSocketClient(clientAddress)
			if err != nil {
				t.Fatalf("Failed to create socket client %d: %v", i, err)
			}

			// Ensure client resources are cleaned up after the test
			defer func(c TestClient) {
				socketClient, ok := c.(*socketClient)
				if !ok {
					t.Errorf("Failed to cast to socketClient type")
				} else {
					err := socketClient.Close()
					if err != nil {
						t.Errorf("Failed to close socket client: %v", err)
					}
				}
			}(client)

			clients[i] = client
		}

		// Run the benchmark test with socket server and clients
		throughputTest(server, clients)
	}
}
