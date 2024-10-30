package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"golang.org/x/exp/rand"
)

func main() {
	maxSize := 16 * 1024 * 1024 // 16MB

	// Initialize AWS/S3 clients
	cfg := aws.DefaultClientConfig()
	cfg.Region = "us-east-1"

	s3Client, err := s3.NewClient(context.Background(), *cfg, logging.NewNoopLogger())
	if err != nil {
		panic("failed to create s3 client: " + err.Error())
	}

	blobStore := blobstorev2.NewBlobStore("eigenda-blobstore", s3Client, logging.NewNoopLogger())

	// Test powers of 2 starting from 1024 bytes
	for size := 1024; size <= maxSize; size *= 2 {
		fmt.Printf("\nTesting size: %d bytes (%.2f KB)\n", size, float64(size)/1024)

		// Generate data
		data := make([]byte, size)
		_, err := rand.New(rand.NewSource(42)).Read(data)
		if err != nil {
			panic(fmt.Sprintf("Failed to create test data for size %d: %v", size, err))
		}

		paddedData := core.PadToPowerOf2(codec.ConvertByPaddingEmptyByte(data))
		fmt.Printf("After padding: %d bytes (%.2f KB)\n", len(paddedData), float64(len(paddedData))/1024)

		// Create blob header
		blobHeader := createDummyBlobHeader(paddedData)

		// Validate field elements
		_, err = rs.ToFrArray(paddedData)
		if err != nil {
			fmt.Printf("❌ Size %d failed field element validation: %v\n", size, err)
			continue
		}

		blobKey, err := blobHeader.BlobKey()
		if err != nil {
			fmt.Printf("❌ Size %d failed to create blob key: %v\n", size, err)
			continue
		}

		// Store blob
		err = blobStore.StoreBlob(context.Background(), blobKey, paddedData)
		if err != nil {
			fmt.Printf("❌ Size %d failed to store: %v\n", size, err)
			continue
		}

		fmt.Printf("✅ Successfully processed and stored blob of size %d (%.2f KB)\n",
			len(paddedData), float64(len(paddedData))/1024)
		fmt.Printf("   Blob key: %s\n", blobKey.Hex())

		// Get encoding params
		encodignParams, err := blobHeader.GetEncodingParams()
		if err != nil {
			panic(fmt.Sprintf("Failed to get encoding params: %v", err))
		}
		fmt.Printf("   Encoding Params: %v %v\n", encodignParams.NumChunks, encodignParams.ChunkLength)
	}
}

func createDummyBlobHeader(data []byte) *corev2.BlobHeader {
	// Create dummy G1 commitment
	commitment := &encoding.G1Commitment{
		X: new(bn254.G1Affine).X,
		Y: new(bn254.G1Affine).Y,
	}

	// Create dummy G2 commitment for length
	lengthCommitment := &encoding.G2Commitment{
		X: new(bn254.G2Affine).X,
		Y: new(bn254.G2Affine).Y,
	}

	// Create dummy length proof
	lengthProof := &encoding.LengthProof{
		X: new(bn254.G2Affine).X,
		Y: new(bn254.G2Affine).Y,
	}

	// Create blob commitments
	blobCommitments := encoding.BlobCommitments{
		Commitment:       commitment,
		LengthCommitment: lengthCommitment,
		LengthProof:      lengthProof,
		Length:           encoding.GetBlobLength(uint(len(data))),
	}

	// Create payment metadata
	paymentMetadata := core.PaymentMetadata{
		AccountID:         "test-account",
		BinIndex:          1,
		CumulativePayment: new(big.Int).SetInt64(1000000),
	}

	return &corev2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: blobCommitments,
		QuorumNumbers:   []core.QuorumID{1, 2, 3},
		PaymentMetadata: paymentMetadata,
		Signature:       []byte("dummy-signature"),
	}
}
