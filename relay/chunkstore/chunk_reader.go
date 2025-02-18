package chunkstore

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common/aws/s3"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// ChunkReader reads chunks written by ChunkWriter.
type ChunkReader interface {

	// GetBinaryChunkProofs reads a slice of proofs from the chunk store, similar to GetChunkProofs.
	// Unlike GetChunkProofs, this method returns the raw serialized bytes of the proofs, as opposed to
	// deserializing them into encoding.Proof structs.
	GetBinaryChunkProofs(ctx context.Context, blobKey corev2.BlobKey) ([][]byte, error)

	// GetBinaryChunkCoefficients reads a slice of frames from the chunk store, similar to GetChunkCoefficients.
	// Unlike GetChunkCoefficients, this method returns the raw serialized bytes of the frames, as opposed to
	// deserializing them into rs.FrameCoeffs structs. The returned uint32 is the number of elements in each frame.
	GetBinaryChunkCoefficients(
		ctx context.Context,
		blobKey corev2.BlobKey,
		fragmentInfo *encoding.FragmentInfo) (uint32, [][]byte, error)
}

var _ ChunkReader = (*chunkReader)(nil)

type chunkReader struct {
	logger logging.Logger
	client s3.Client
	bucket string
}

// NewChunkReader creates a new ChunkReader.
//
// This chunk reader will only return data for the shards specified in the shards parameter.
// If empty, it will return data for all shards. (Note: shard feature is not yet implemented.)
func NewChunkReader(
	logger logging.Logger,
	s3Client s3.Client,
	bucketName string) ChunkReader {

	return &chunkReader{
		logger: logger,
		client: s3Client,
		bucket: bucketName,
	}
}

func (r *chunkReader) GetBinaryChunkProofs(ctx context.Context, blobKey corev2.BlobKey) ([][]byte, error) {
	bytes, err := r.client.DownloadObject(ctx, r.bucket, s3.ScopedProofKey(blobKey))
	if err != nil {
		r.logger.Error("Failed to download chunks from S3: %v", err)
		return nil, fmt.Errorf("failed to download chunks from S3: %w", err)
	}

	proofs, err := rs.SplitSerializedFrameProofs(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to split proofs: %w", err)
	}

	return proofs, nil
}

func (r *chunkReader) GetBinaryChunkCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
	fragmentInfo *encoding.FragmentInfo) (uint32, [][]byte, error) {

	bytes, err := r.client.FragmentedDownloadObject(
		ctx,
		r.bucket,
		s3.ScopedChunkKey(blobKey),
		int(fragmentInfo.TotalChunkSizeBytes),
		int(fragmentInfo.FragmentSizeBytes))

	if err != nil {
		r.logger.Error("Failed to download chunks from S3: %v", err)
		return 0, nil, fmt.Errorf("failed to download chunks from S3: %w", err)
	}

	elementCount, frames, err := rs.SplitSerializedFrameCoeffs(bytes)
	if err != nil {
		r.logger.Error("Failed to split frames: %v", err)
		return 0, nil, fmt.Errorf("failed to split frames: %w", err)
	}

	return elementCount, frames, nil
}
