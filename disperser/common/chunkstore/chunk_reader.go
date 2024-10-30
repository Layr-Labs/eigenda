package chunkstore

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

// ChunkReader reads chunks written by ChunkWriter.
type ChunkReader interface {
	// GetChunkProofs reads a slice of proofs from the chunk store.
	GetChunkProofs(ctx context.Context, blobKey disperser.BlobKey) ([]*encoding.Proof, error)
	// GetChunkCoefficients reads a slice of frames from the chunk store.
	GetChunkCoefficients(ctx context.Context, blobKey disperser.BlobKey) ([]*rs.Frame, error)
}

var _ ChunkReader = (*chunkReader)(nil)

type chunkReader struct {
	logger             logging.Logger
	chunkMetadataStore *ChunkMetadataStore
	client             s3.Client
	bucket             string
	shards             []uint32
}

// NewChunkReader creates a new ChunkReader.
func NewChunkReader(
	logger logging.Logger,
	chunkMetadataStore *ChunkMetadataStore,
	s3Client s3.Client,
	bucketName string,
	shards []uint32) ChunkReader {

	return &chunkReader{
		logger:             logger,
		chunkMetadataStore: chunkMetadataStore,
		client:             s3Client,
		bucket:             bucketName,
		shards:             shards,
	}
}

func (r *chunkReader) GetChunkProofs(ctx context.Context, blobKey disperser.BlobKey) ([]*encoding.Proof, error) {
	s3Key := blobKey.String()

	bytes, err := r.client.DownloadObject(ctx, r.bucket, s3Key)
	if err != nil {
		r.logger.Error("Failed to download chunks from S3: %v", err)
		return nil, fmt.Errorf("failed to download chunks from S3: %w", err)
	}

	if len(bytes)%bn254.SizeOfG1AffineCompressed != 0 {
		r.logger.Error("Invalid proof size")
		return nil, fmt.Errorf("invalid proof size: %w", err)
	}

	proofCount := len(bytes) / bn254.SizeOfG1AffineCompressed
	proofs := make([]*encoding.Proof, proofCount)

	for i := 0; i < proofCount; i++ {
		proof := encoding.Proof{}
		err := proof.Unmarshal(bytes[i*bn254.SizeOfG1AffineCompressed:])
		if err != nil {
			r.logger.Error("Failed to unmarshal proof: %v", err)
			return nil, fmt.Errorf("failed to unmarshal proof: %w", err)
		}
		proofs[i] = &proof
	}

	return proofs, nil
}

func (r *chunkReader) GetChunkCoefficients(ctx context.Context, blobKey disperser.BlobKey) ([]*rs.Frame, error) {
	s3Key := blobKey.String()

	bytes, err := r.client.DownloadObject(ctx, r.bucket, s3Key)
	if err != nil {
		r.logger.Error("Failed to download chunks from S3: %v", err)
		return nil, fmt.Errorf("failed to download chunks from S3: %w", err)
	}

	frames, err := rs.DecodeFrames(bytes)
	if err != nil {
		r.logger.Error("Failed to decode frames: %v", err)
		return nil, fmt.Errorf("failed to decode frames: %w", err)
	}

	return frames, nil
}
