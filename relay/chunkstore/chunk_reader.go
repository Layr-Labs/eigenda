package chunkstore

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/tracing"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"go.opentelemetry.io/otel/attribute"
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
	ctx, span := tracing.TraceOperation(ctx, "chunkReader.GetBinaryChunkProofs")
	defer span.End()

	// Add span for S3 download operation
	downloadCtx, downloadSpan := tracing.TraceOperation(ctx, "chunkReader.GetBinaryChunkProofs.s3Download")
	downloadSpan.SetAttributes(
		attribute.String("blob_key", blobKey.Hex()),
		attribute.String("bucket", r.bucket),
	)
	bytes, err := r.client.DownloadObject(downloadCtx, r.bucket, s3.ScopedProofKey(blobKey))
	downloadSpan.End()
	if err != nil {
		r.logger.Error("failed to download proofs from S3", "blob", blobKey.Hex(), "error", err)
		return nil, fmt.Errorf("failed to download proofs from S3 for blob %s: %w", blobKey.Hex(), err)
	}

	// Add span for proof splitting operation
	_, splitSpan := tracing.TraceOperation(ctx, "chunkReader.GetBinaryChunkProofs.splitProofs")
	splitSpan.SetAttributes(attribute.Int("bytes_length", len(bytes)))
	proofs, err := rs.SplitSerializedFrameProofs(bytes)
	splitSpan.End()
	if err != nil {
		r.logger.Error("failed to split proofs", "blob", blobKey.Hex(), "error", err)
		return nil, fmt.Errorf("failed to split proofs for blob %s: %w", blobKey.Hex(), err)
	}

	span.SetAttributes(attribute.Int("proof_count", len(proofs)))
	return proofs, nil
}

func (r *chunkReader) GetBinaryChunkCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
	fragmentInfo *encoding.FragmentInfo) (uint32, [][]byte, error) {

	ctx, span := tracing.TraceOperation(ctx, "chunkReader.GetBinaryChunkCoefficients")
	defer span.End()

	// Add span for S3 download operation
	downloadCtx, downloadSpan := tracing.TraceOperation(ctx, "chunkReader.GetBinaryChunkCoefficients.s3Download")
	downloadSpan.SetAttributes(
		attribute.String("blob_key", blobKey.Hex()),
		attribute.String("bucket", r.bucket),
		attribute.Int64("total_chunk_size", int64(fragmentInfo.TotalChunkSizeBytes)),
		attribute.Int64("fragment_size", int64(fragmentInfo.FragmentSizeBytes)),
	)
	bytes, err := r.client.FragmentedDownloadObject(
		downloadCtx,
		r.bucket,
		s3.ScopedChunkKey(blobKey),
		int(fragmentInfo.TotalChunkSizeBytes),
		int(fragmentInfo.FragmentSizeBytes))
	downloadSpan.End()

	if err != nil {
		r.logger.Error("failed to download coefficients from S3",
			"blob", blobKey.Hex(),
			"totalSize", fragmentInfo.TotalChunkSizeBytes,
			"fragmentSize", fragmentInfo.FragmentSizeBytes,
			"error", err)
		return 0, nil, fmt.Errorf("failed to download coefficients from S3 for blob %s (total size: %d, fragment size: %d): %w",
			blobKey.Hex(), fragmentInfo.TotalChunkSizeBytes, fragmentInfo.FragmentSizeBytes, err)
	}

	// Add span for coefficient splitting operation
	_, splitSpan := tracing.TraceOperation(ctx, "chunkReader.GetBinaryChunkCoefficients.splitCoefficients")
	splitSpan.SetAttributes(attribute.Int("bytes_length", len(bytes)))
	elementCount, frames, err := rs.SplitSerializedFrameCoeffs(bytes)
	splitSpan.End()
	if err != nil {
		r.logger.Error("failed to split coefficient frames", "blob", blobKey.Hex(), "error", err)
		return 0, nil, fmt.Errorf("failed to split coefficient frames for blob %s: %w", blobKey.Hex(), err)
	}

	span.SetAttributes(
		attribute.Int("frame_count", len(frames)),
		attribute.Int64("element_count", int64(elementCount)),
	)
	return elementCount, frames, nil
}
