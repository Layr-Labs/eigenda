package encoder_test

import (
	"context"
	"math/big"
	"runtime"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	pb "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder/v2"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func makeTestProver(numPoint uint64) (encoding.Prover, error) {
	// We need the larger SRS for testing the encoder with 8192 chunks
	kzgConfig := &kzg.KzgConfig{
		G1Path:          "../../inabox/resources/kzg/g1.point.300000",
		G2Path:          "../../inabox/resources/kzg/g2.point.300000",
		G2PowerOf2Path:  "../../inabox/resources/kzg/g2.point.300000.powerOf2",
		CacheDir:        "../../inabox/resources/kzg/SRSTables",
		SRSOrder:        300000,
		SRSNumberToLoad: numPoint,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	p, err := prover.NewProver(kzgConfig, false)

	return p, err
}

func TestEncodeBlobToChunkStore(t *testing.T) {
	const (
		testDataSize   = 16 * 1024
		timeoutSeconds = 30
		randSeed       = uint64(42)
	)

	var (
		codingRatio = v2.ParametersMap[0].CodingRate
		numChunks   = v2.ParametersMap[0].NumChunks
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()

	createTestData := func(t *testing.T, size int) []byte {
		t.Helper()
		data := make([]byte, size)
		rand.New(rand.NewSource(randSeed)).Read(data)
		return core.PadToPowerOf2(codec.ConvertByPaddingEmptyByte(data))
	}

	// Setup test data
	data := createTestData(t, testDataSize)
	blobSize := uint(len(data))
	blobLength := encoding.GetBlobLength(blobSize)

	// Get chunk length for blob version 0
	chunkLength, err := v2.GetChunkLength(0, uint32(blobLength))
	if !assert.NoError(t, err, "Failed to get chunk length") {
		t.FailNow()
	}

	t.Logf("Test parameters: blobversion=%d, blobLength=%d, codingRatio=%d, numChunks=%d, chunkLength=%d",
		0, blobLength, codingRatio, numChunks, chunkLength)

	// Create blob header and key
	blobHeader := createTestBlobHeader(t)
	blobKey, err := blobHeader.BlobKey()
	if !assert.NoError(t, err, "Failed to create blob key") {
		t.FailNow()
	}

	// Store test data
	if err := blobStore.StoreBlob(ctx, blobKey, data); !assert.NoError(t, err, "Failed to store blob") {
		t.FailNow()
	}

	// Verify storage succeded
	t.Run("Verify Blob Storage", func(t *testing.T) {
		storedData, err := blobStore.GetBlob(ctx, blobKey)
		assert.NoError(t, err, "Failed to get stored blob")
		assert.Equal(t, data, storedData, "Stored data doesn't match original")
	})

	// Initialize encoder server
	server := initializeEncoder(t)

	// Create and execute encoding request
	req := &pb.EncodeBlobRequest{
		BlobKey: blobKey[:],
		EncodingParams: &pb.EncodingParams{
			ChunkLength: uint32(chunkLength),
			NumChunks:   uint32(numChunks),
		},
	}

	resp, err := server.EncodeBlobToChunkStore(ctx, req)
	if !assert.NoError(t, err, "EncodeBlobToChunkStore failed") {
		t.FailNow()
	}

	// Verify encoding results
	t.Run("Verify Encoding Results", func(t *testing.T) {
		assert.NotNil(t, resp, "Response should not be nil")
		assert.Equal(t, uint32(294916), resp.FragmentInfo.TotalChunkSizeBytes, "Unexpected total chunk size")
		assert.Equal(t, uint32(512*1024), resp.FragmentInfo.FragmentSizeBytes, "Unexpected fragment size")
	})

	// Verify chunk store data
	t.Run("Verify Chunk Store Data", func(t *testing.T) {
		// Check proofs
		proofs, err := chunkStoreReader.GetChunkProofs(ctx, blobKey)
		assert.NoError(t, err, "Failed to get chunk proofs")
		assert.Len(t, proofs, int(numChunks), "Unexpected number of proofs")

		// Check coefficients
		fragmentInfo := &encoding.FragmentInfo{
			TotalChunkSizeBytes: resp.FragmentInfo.TotalChunkSizeBytes,
			FragmentSizeBytes:   resp.FragmentInfo.FragmentSizeBytes,
		}
		coefficients, err := chunkStoreReader.GetChunkCoefficients(ctx, blobKey, fragmentInfo)
		assert.NoError(t, err, "Failed to get chunk coefficients")
		assert.Len(t, coefficients, int(numChunks), "Unexpected number of coefficients")
	})
}

// Helper function to create test blob header
func createTestBlobHeader(t *testing.T) *v2.BlobHeader {
	t.Helper()
	return &v2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   []core.QuorumID{0},
		BlobCommitments: mockCommitment,
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x1234",
			BinIndex:          0,
			CumulativePayment: big.NewInt(532),
		},
	}
}

// Helper function to initialize encoder
func initializeEncoder(t *testing.T) *encoder.EncoderServerV2 {
	t.Helper()
	prover, err := makeTestProver(300000)
	if !assert.NoError(t, err, "Failed to create prover") {
		t.FailNow()
	}

	metrics := encoder.NewMetrics("9000", logger)
	return encoder.NewEncoderServerV2(encoder.ServerConfig{
		GrpcPort:              "8080",
		MaxConcurrentRequests: 10,
		RequestPoolSize:       5,
	}, blobStore, chunkStoreWriter, logger, prover, metrics)
}
