package chunkstore

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common/s3"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
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
	) (uint32, [][]byte, error)
}

var _ ChunkReader = (*chunkReader)(nil)

type chunkReader struct {
	client s3.S3Client
	bucket string
}

// NewChunkReader creates a new ChunkReader.
//
// This chunk reader will only return data for the shards specified in the shards parameter.
// If empty, it will return data for all shards. (Note: shard feature is not yet implemented.)
func NewChunkReader(
	s3Client s3.S3Client,
	bucketName string) ChunkReader {

	return &chunkReader{
		client: s3Client,
		bucket: bucketName,
	}
}

func (r *chunkReader) GetBinaryChunkProofs(ctx context.Context, blobKey corev2.BlobKey) ([][]byte, error) {
	bytes, err := r.client.DownloadObject(ctx, r.bucket, s3.ScopedProofKey(blobKey))
	if err != nil {
		return nil, fmt.Errorf("failed to download proofs from S3 for blob %s: %w", blobKey.Hex(), err)
	}

	proofs, err := encoding.SplitSerializedFrameProofs(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to split proofs for blob %s: %w", blobKey.Hex(), err)
	}

	return proofs, nil
}

func (r *chunkReader) GetBinaryChunkCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
) (uint32, [][]byte, error) {

	bytes, err := r.client.DownloadObject(ctx, r.bucket, s3.ScopedChunkKey(blobKey))

	if err != nil {
		return 0, nil, fmt.Errorf("failed to download coefficients from S3 for blob %s: %w", blobKey.Hex(), err)
	}

	elementCount, frames, err := rs.SplitSerializedFrameCoeffs(bytes)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to split coefficient frames for blob %s: %w", blobKey.Hex(), err)
	}

	return elementCount, frames, nil
}
