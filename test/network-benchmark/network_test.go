package network_benchmark

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/relay"
	testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/docker/go-units"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024 // 1MB

func TestProtobufThroughput(t *testing.T) {
	dataPerTransfer := 1 * units.MiB
	totalDataToTransfer := 10 * units.GiB
	parallelism := 8 // TODO claude, run this many transfers in parallel

	rand := testrandom.NewTestRandom()
	randomData := newReusableRandomness(units.GiB, rand.Int63())

	// Set up server and client on localhost
	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer()
	protobufServer := NewProtobufServer(randomData)
	relay.RegisterThroughputTestServer(server, protobufServer)

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Fatalf("Failed to serve: %v", err)
		}
	}()
	defer server.Stop()

	// Set up client
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	grpcClient := relay.NewThroughputTestClient(conn)
	client := newProtobufClient(grpcClient)

	// Initialize variables to track metrics
	var totalLatency time.Duration
	var totalDataTransferred int64
	var transferCount int64
	lastStatusTime := time.Now()
	statusInterval := 1 * time.Second

	start := time.Now()

	// Run the benchmark
	for totalDataTransferred < int64(totalDataToTransfer) {
		// Generate a random seed for deterministic but varying data
		seed := rand.Int63()

		// Measure latency for this request
		requestStart := time.Now()
		data, err := client.getData(int64(dataPerTransfer), seed)
		if err != nil {
			t.Fatalf("Failed to get data: %v", err)
		}
		requestLatency := time.Since(requestStart)

		// Regenerate the data using the same seed and verify it matches
		expectedData := randomData.getData(int64(dataPerTransfer), seed)
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
		totalLatency += requestLatency
		dataSize := int64(len(data))
		totalDataTransferred += dataSize
		transferCount++

		// Print status periodically
		if time.Since(lastStatusTime) >= statusInterval {
			elapsedSoFar := time.Since(start)
			throughputSoFar := float64(totalDataTransferred) / elapsedSoFar.Seconds()

			fmt.Printf("[%s] Transferred: %s (%.2f%%), Throughput: %s/s\n",
				elapsedSoFar.Round(time.Second),
				units.BytesSize(float64(totalDataTransferred)),
				100.0*float64(totalDataTransferred)/float64(totalDataToTransfer),
				units.BytesSize(throughputSoFar))

			lastStatusTime = time.Now()
		}
	}

	elapsed := time.Since(start)

	// Calculate final metrics
	avgLatency := totalLatency / time.Duration(transferCount)
	throughput := float64(totalDataTransferred) / elapsed.Seconds()

	// Print the benchmark results
	fmt.Printf("\n--- Benchmark Results ---\n")
	fmt.Printf("Total time: %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("Average latency: %s\n", avgLatency.Round(time.Microsecond))
	fmt.Printf("Average throughput: %s/s\n", units.BytesSize(throughput))
	fmt.Printf("Total data transferred: %s\n", units.BytesSize(float64(totalDataTransferred)))
	fmt.Printf("Number of transfers: %d\n", transferCount)
}
