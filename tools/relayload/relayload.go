package relayload

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type DataAPIResponse struct {
	Blobs []struct {
		BlobKey string `json:"blob_key"`
	} `json:"blobs"`
}

type RelayLoad struct {
	relayUrl     string
	blobKeys     []string
	rangeSizes   []int
	requestSizes []int
	numThreads   int
	logger       logging.Logger
}

func fetchBlobKeysFromAPI(dataAPIUrl string, logger logging.Logger) ([]string, error) {
	requestUrl := fmt.Sprintf("%s/api/v2/blobs/feed?direction=backward&limit=1000", dataAPIUrl)
	logger.Info("Fetching blob keys from API", "url", requestUrl)
	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse DataAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}

	blobKeys := make([]string, len(apiResponse.Blobs))
	for i, blob := range apiResponse.Blobs {
		blobKeys[i] = blob.BlobKey
	}

	logger.Info("Fetched blob keys from API", "count", len(blobKeys))
	return blobKeys, nil
}

func NewRelayLoad(relayUrl string, dataAPIUrl string, rangeSizes []int, requestSizes []int, numThreads int, logger logging.Logger) *RelayLoad {
	if dataAPIUrl == "" {
		logger.Error("Data API URL is not set")
		os.Exit(1)
	}
	blobKeys, err := fetchBlobKeysFromAPI(dataAPIUrl, logger)
	if err != nil {
		logger.Error("Failed to fetch blob keys from API", "error", err)
		os.Exit(1)
	}
	return &RelayLoad{
		relayUrl:     relayUrl,
		blobKeys:     blobKeys,
		rangeSizes:   rangeSizes,
		requestSizes: requestSizes,
		numThreads:   numThreads,
		logger:       logger,
	}
}

func humanizeBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (c *RelayLoad) GetChunks(ctx context.Context, operatorId string) ([]byte, error) {
	operatorIdBytes, err := hex.DecodeString(operatorId)
	if err != nil {
		return nil, fmt.Errorf("invalid operator ID: %v", err)
	}

	requestSize := c.requestSizes[rand.Intn(len(c.requestSizes))]

	// Create chunk requests for all blob keys
	chunkRequests := make([]*relay.ChunkRequest, requestSize)
	for i := 0; i < requestSize; i++ {
		blobKey := c.blobKeys[rand.Intn(len(c.blobKeys))]
		blobKeyBytes, err := hex.DecodeString(blobKey)
		if err != nil {
			return nil, fmt.Errorf("invalid blob key: %v", err)
		}

		rangeSize := c.rangeSizes[rand.Intn(len(c.rangeSizes))]
		startIndex := rand.Intn(8064 - rangeSize)
		endIndex := startIndex + rangeSize

		chunkRequests[i] = &relay.ChunkRequest{
			Request: &relay.ChunkRequest_ByRange{
				ByRange: &relay.ChunkRequestByRange{
					BlobKey:    blobKeyBytes,
					StartIndex: uint32(startIndex),
					EndIndex:   uint32(endIndex),
				},
			},
		}
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"h2"},
		})),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1000 * 1024 * 1024)), // 1GB
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(1000 * 1024 * 1024)), // 1GB
	}

	conn, err := grpc.Dial(c.relayUrl, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	relayClient := relay.NewRelayClient(conn)

	request := &relay.GetChunksRequest{
		ChunkRequests: chunkRequests,
		OperatorId:    operatorIdBytes,
		Timestamp:     uint32(time.Now().Unix()),
	}

	reply, err := relayClient.GetChunks(ctx, request)
	if err != nil {
		return nil, err
	}

	totalSize := 0
	for _, data := range reply.Data {
		totalSize += len(data)
	}
	c.logger.Info("Received chunks", "request_size", requestSize, "total_size", humanizeBytes(totalSize))
	return reply.Data[0], nil
}

func (c *RelayLoad) RunParallel(ctx context.Context, operatorId string, numThreads int) error {
	var wg sync.WaitGroup
	errChan := make(chan error, numThreads)
	rateLimiter := time.NewTicker(100 * time.Millisecond) // Rate limit of 10 requests per second per thread
	defer rateLimiter.Stop()

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-rateLimiter.C:
					start := time.Now()
					_, err := c.GetChunks(ctx, operatorId)
					if err != nil {
						errChan <- fmt.Errorf("thread %d error: %v", threadID, err)
						return
					}
					elapsed := time.Since(start)
					c.logger.Info("Thread completed request",
						"thread_id", threadID,
						"duration", elapsed)
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete or an error to occur
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
