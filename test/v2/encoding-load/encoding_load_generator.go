package encodingload

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/common/pprof"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/docker/go-units"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/metadata"
)

// EncodingLoadGenerator is a utility for generating encoding load for the target network.
type EncodingLoadGenerator struct {
	// The blob store for storing blobs.
	blobStore *blobstorev2.BlobStore

	ctx    context.Context
	cancel context.CancelFunc

	// The configuration for the encoding load generator.
	config *EncodingLoadGeneratorConfig
	// The logger to use for the load test.
	logger logging.Logger
	// The encoder client for encoding operations.
	encoderClient disperser.EncoderClientV2
	// The time between starting each blob encoding.
	encodingPeriod time.Duration
	// The channel to limit the number of parallel blob encodings.
	parallelismLimiter chan struct{}
	// if true, the encoding load generator is running.
	alive atomic.Bool
	// The channel to signal when the encoding load generator is finished.
	finishedChan chan struct{}
	// The metrics for the encoding load generator.
	metrics *encodingLoadGeneratorMetrics
}

// ReadConfigFile loads a EncodingLoadGeneratorConfig from a file.
func ReadConfigFile(filePath string) (*EncodingLoadGeneratorConfig, error) {
	// Resolve tilde in path if present
	if len(filePath) > 0 && filePath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		filePath = homeDir + filePath[1:]
	}

	configFileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &EncodingLoadGeneratorConfig{}
	err = json.Unmarshal(configFileBytes, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return config, nil
}

// NewEncodingLoadGenerator creates a new EncodingLoadGenerator.
func NewEncodingLoadGenerator(
	config *EncodingLoadGeneratorConfig,
	logger logging.Logger,
	metricsRegistry *prometheus.Registry,
	encoderClient disperser.EncoderClientV2,
	blobStore *blobstorev2.BlobStore) *EncodingLoadGenerator {

	bytesPerSecond := config.MBPerSecond * units.MiB
	averageBlobSize := config.AverageBlobSizeMB * units.MiB

	encodingFrequency := bytesPerSecond / averageBlobSize
	encodingPeriod := 1 / encodingFrequency
	encodingPeriodAsDuration := time.Duration(encodingPeriod * float64(time.Second))

	parallelismLimiter := make(chan struct{}, config.MaxParallelism)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	metrics := newEncodingLoadGeneratorMetrics(metricsRegistry)

	if config.EnablePprof {
		pprofProfiler := pprof.NewPprofProfiler(fmt.Sprintf("%d", config.PprofHttpPort), logger)
		go pprofProfiler.Start()
		logger.Info("Enabled pprof", "port", config.PprofHttpPort)
	}

	return &EncodingLoadGenerator{
		blobStore:          blobStore,
		ctx:                ctx,
		cancel:             cancel,
		config:             config,
		logger:             logger,
		encoderClient:      encoderClient,
		encodingPeriod:     encodingPeriodAsDuration,
		parallelismLimiter: parallelismLimiter,
		alive:              atomic.Bool{},
		finishedChan:       make(chan struct{}),
		metrics:            metrics,
	}
}

// Start starts the encoding load generator. If block is true, this function will block until Stop() or
// the encoding load generator crashes. If block is false, this function will return immediately.
func (l *EncodingLoadGenerator) Start(block bool) {
	l.alive.Store(true)
	go l.run()

	if block {
		<-l.finishedChan
	}
}

// Stop stops the encoding load generator.
func (l *EncodingLoadGenerator) Stop() {
	l.alive.Store(false)
	l.cancel()
	<-l.finishedChan
}

// run is the main loop of the encoding load generator.
func (l *EncodingLoadGenerator) run() {
	defer close(l.finishedChan)

	ticker := time.NewTicker(l.encodingPeriod)
	defer ticker.Stop()

	for l.alive.Load() {
		select {
		case <-ticker.C:
			go l.encodeBlob()
		case <-l.ctx.Done():
			return
		}
	}
}

// encodeBlob encodes a single blob.
func (l *EncodingLoadGenerator) encodeBlob() {
	l.parallelismLimiter <- struct{}{}
	defer func() {
		<-l.parallelismLimiter
	}()

	l.metrics.startOperation()
	defer l.metrics.endOperation()

	rand := random.NewTestRandomNoPrint()

	payloadSize := int(rand.BoundedGaussian(
		l.config.AverageBlobSizeMB*units.MiB,
		l.config.BlobSizeStdDev*units.MiB,
		1.0,
		float64(l.config.MaxBlobSize+1)))

	// Print blob size and above settings
	l.logger.Info("Blob size", "size", payloadSize, "average", l.config.AverageBlobSizeMB*units.MiB, "stddev", l.config.BlobSizeStdDev*units.MiB, "max", l.config.MaxBlobSize*units.MiB)

	payloadBytes := rand.Bytes(payloadSize)
	payload := coretypes.NewPayload(payloadBytes)
	blob, err := payload.ToBlob(codecs.PolynomialFormCoeff)
	if err != nil {
		l.logger.Error("Failed to convert payload to blob", "error", err)
		return
	}

	// Use the parent context directly without timeout
	ctx := l.ctx

	// Perform the encoding operation
	// This is a simplified version - in a real implementation, you would use the actual encoding API
	err = l.performEncoding(ctx, blob)
	if err != nil {
		l.logger.Error("Failed to encode blob", "error", err)
		return
	}
}

// performEncoding performs the actual encoding operation using the encoder client v2.
func (l *EncodingLoadGenerator) performEncoding(ctx context.Context, blob *coretypes.Blob) error {
	var encoderClient disperser.EncoderClientV2
	var err error

	// Use the existing client if available, otherwise create a new one
	if l.encoderClient != nil {
		encoderClient = l.encoderClient
	} else {
		// Get the encoder client from the environment
		encoderURL := os.Getenv("ENCODER_URL")
		if encoderURL == "" {
			encoderURL = "localhost:8090" // Default encoder URL
		}

		// Create the encoder client v2
		encoderClient, err = encoder.NewEncoderClientV2(encoderURL)
		if err != nil {
			return fmt.Errorf("failed to create encoder client v2: %w", err)
		}

		// Store the client for future use
		l.encoderClient = encoderClient
	}

	// In a real scenario, we would get blob version parameters from the chain
	// For the load generator, we'll use some reasonable defaults that match
	// how the controller would calculate encoding parameters
	blobParams := &core.BlobVersionParameters{
		NumChunks:       8192,
		CodingRate:      8,
		MaxNumOperators: 3537,
	}

	// Create a mock BlobHeader with commitments
	// In a real scenario, these would be calculated from the data
	// For the load test, we'll create a simplified version with mock values
	var commitX, commitY fp.Element
	commitX = *commitX.SetBigInt(big.NewInt(1))
	commitY = *commitY.SetBigInt(big.NewInt(2))

	commitment := &encoding.G1Commitment{
		X: commitX,
		Y: commitY,
	}

	// Create mock values for G2 commitments
	var lengthXA0, lengthXA1, lengthYA0, lengthYA1 fp.Element
	lengthXA0.SetBigInt(big.NewInt(10857046999023057))
	lengthXA1.SetBigInt(big.NewInt(11559732032986387))
	lengthYA0.SetBigInt(big.NewInt(8495653923123431))
	lengthYA1.SetBigInt(big.NewInt(4082367875863433))

	var lengthG2 bn254.G2Affine
	lengthG2.X.A0 = lengthXA0
	lengthG2.X.A1 = lengthXA1
	lengthG2.Y.A0 = lengthYA0
	lengthG2.Y.A1 = lengthYA1

	// Length should be power of 2
	blobLength := uint(len(blob.Serialize()))
	symbolLength := encoding.GetBlobLengthPowerOf2(blobLength)

	blobCommitments := encoding.BlobCommitments{
		Commitment:       commitment,
		LengthCommitment: (*encoding.G2Commitment)(&lengthG2),
		LengthProof:      (*encoding.G2Commitment)(&lengthG2),
		Length:           symbolLength,
	}

	// Create a mock PaymentMetadata
	// In a real scenario, this would come from the payment system
	paymentMetadata := core.PaymentMetadata{
		AccountID:         gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"),
		Timestamp:         time.Now().UnixNano(),
		CumulativePayment: big.NewInt(1000000), // Mock payment amount
	}

	// Create a BlobHeader
	blobHeader := &corev2.BlobHeader{
		BlobVersion:     1, // Use version 1 for simplicity
		BlobCommitments: blobCommitments,
		QuorumNumbers:   []core.QuorumID{0}, // Use quorum 0 for simplicity
		PaymentMetadata: paymentMetadata,
	}

	// Generate the blob key from the header
	blobKey, err := blobHeader.BlobKey()
	if err != nil {
		return fmt.Errorf("failed to generate blob key: %w", err)
	}

	// Store the blob in the blob store
	err = l.blobStore.StoreBlob(ctx, blobKey, blob.Serialize())
	if err != nil {
		return fmt.Errorf("failed to store blob: %w", err)
	}

	// Calculate encoding parameters using the same method as the controller
	encodingParams, err := corev2.GetEncodingParams(blobCommitments.Length, blobParams)
	if err != nil {
		return fmt.Errorf("failed to get encoding params: %w", err)
	}

	// Add metadata headers similar to the controller
	md := metadata.New(map[string]string{
		"content-type": "application/grpc",
		"x-blob-size":  fmt.Sprintf("%d", blobCommitments.Length),
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Call the encoder client to encode the blob
	l.logger.Info("Encoding blob",
		"blobKey", blobKey.Hex(),
		"blobSize", blobCommitments.Length,
		"numChunks", encodingParams.NumChunks,
		"chunkLength", encodingParams.ChunkLength)

	fragmentInfo, err := encoderClient.EncodeBlob(ctx, blobKey, encodingParams, uint64(blobCommitments.Length))
	if err != nil {
		return fmt.Errorf("failed to encode blob: %w", err)
	}

	l.logger.Info("Successfully encoded blob",
		"blobKey", blobKey.Hex(),
		"fragmentSizeBytes", fragmentInfo.FragmentSizeBytes,
		"totalChunkSizeBytes", fragmentInfo.TotalChunkSizeBytes)

	return nil
}

// performValidation performs validation of the encoded data.
func (l *EncodingLoadGenerator) performValidation(ctx context.Context, data []byte) error {
	// In a real implementation, this would validate the encoded data
	// For now, we'll just simulate the validation operation with a delay
	select {
	case <-time.After(time.Duration(float64(len(data)) / (200 * units.MiB) * float64(time.Second))):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
