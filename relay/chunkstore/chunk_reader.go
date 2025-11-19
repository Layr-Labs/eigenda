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
	// deserializing them into rs.FrameCoeffs structs. The returned uint32 is the number of symbols per frame.
	GetBinaryChunkCoefficients(
		ctx context.Context,
		blobKey corev2.BlobKey,
	) (uint32, [][]byte, error)

	// GetBinaryChunkProofsRange reads a range of proofs from the chunk store.
	GetBinaryChunkProofsRange(
		ctx context.Context,
		blobKey corev2.BlobKey,
		// The index of the first proof to fetch (inclusive).
		startIndex uint32,
		// The index of the last proof to fetch (exclusive).
		endIndex uint32,
	) ([][]byte, bool, error)

	// GetBinaryChunkCoefficientRange reads a range of chunks from the chunk store.
	GetBinaryChunkCoefficientRange(
		ctx context.Context,
		blobKey corev2.BlobKey,
		// The index of the first chunk to fetch (inclusive).
		startIndex uint32,
		// The index of the last chunk to fetch (exclusive).
		endIndex uint32,
		// The number of symbols per frame. Required to determine the exact byte range to fetch.
		symbolsPerFrame uint32,
	) ([][]byte, bool, error)
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
	bytes, found, err := r.client.DownloadObject(ctx, r.bucket, s3.ScopedProofKey(blobKey))
	if err != nil {
		return nil, fmt.Errorf("failed to download proofs from S3 for blob %s: %w", blobKey.Hex(), err)
	}

	if !found {
		return nil, fmt.Errorf("proofs not found for blob %s", blobKey.Hex())
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

	bytes, found, err := r.client.DownloadObject(ctx, r.bucket, s3.ScopedChunkKey(blobKey))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to download coefficients from S3 for blob %s: %w", blobKey.Hex(), err)
	}

	if !found {
		return 0, nil, fmt.Errorf("coefficients not found for blob %s", blobKey.Hex())
	}

	elementCount, frames, err := rs.SplitSerializedFrameCoeffs(bytes)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to split coefficient frames for blob %s: %w", blobKey.Hex(), err)
	}

	return elementCount, frames, nil
}

func (r *chunkReader) GetBinaryChunkProofsRange(
	ctx context.Context,
	blobKey corev2.BlobKey,
	firstChunkIndex uint32,
	endChunkIndex uint32,
) ([][]byte, bool, error) {

	if firstChunkIndex >= endChunkIndex {
		return nil, false, fmt.Errorf("invalid startIndex (%d) or endIndex (%d)", firstChunkIndex, endChunkIndex)
	}

	firstByteIndex := firstChunkIndex * encoding.SerializedProofLength
	count := endChunkIndex - firstChunkIndex
	size := count * encoding.SerializedProofLength

	s3Key := s3.ScopedProofKey(blobKey)

	data, found, err := r.client.DownloadPartialObject(
		ctx,
		r.bucket,
		s3Key,
		int64(firstByteIndex),
		int64(firstByteIndex+size))
	if err != nil {
		return nil, false, fmt.Errorf("failed to download proofs from S3 for blob %s: %w", blobKey.Hex(), err)
	}

	if !found {
		return nil, false, nil
	}

	proofs, err := encoding.SplitSerializedFrameProofs(data)
	if err != nil {
		return nil, false, fmt.Errorf("failed to split proofs for blob %s: %w", blobKey.Hex(), err)
	}

	return proofs, true, nil
}

func (r *chunkReader) GetBinaryChunkCoefficientRange(
	ctx context.Context,
	blobKey corev2.BlobKey,
	startIndex uint32,
	endIndex uint32,
	symbolsPerFrame uint32,
) ([][]byte, bool, error) {

	if startIndex >= endIndex {
		return nil, false, fmt.Errorf("invalid startIndex (%d) or endIndex (%d)", startIndex, endIndex)
	}

	if symbolsPerFrame == 0 {
		return nil, false, fmt.Errorf("symbolsPerFrame must be greater than 0")
	}

	bytesPerFrame := encoding.BYTES_PER_SYMBOL * symbolsPerFrame
	firstByteIndex := 4 + startIndex*bytesPerFrame
	size := (endIndex - startIndex) * bytesPerFrame

	s3Key := s3.ScopedChunkKey(blobKey)

	data, found, err := r.client.DownloadPartialObject(
		ctx,
		r.bucket,
		s3Key,
		int64(firstByteIndex),
		int64(firstByteIndex+size))
	if err != nil {
		return nil, false, fmt.Errorf("failed to download coefficients from S3 for blob %s: %w", blobKey.Hex(), err)
	}

	if !found {
		return nil, false, nil
	}

	// Deserialize the frames
	frames, err := rs.SplitSerializedFrameCoeffsWithElementCount(data, symbolsPerFrame)
	if err != nil {
		return nil, false, fmt.Errorf(
			"failed to split coefficient frames for blob %s, symbols per frame %d: %w",
			blobKey.Hex(), symbolsPerFrame, err)
	}

	return frames, true, nil
}
