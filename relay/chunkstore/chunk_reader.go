package chunkstore

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common/aws/s3"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

// ChunkReader reads chunks written by ChunkWriter.
type ChunkReader interface {

	// GetChunkProofs reads a slice of proofs from the chunk store.
	GetChunkProofs(ctx context.Context, blobKey corev2.BlobKey) ([]*encoding.Proof, error)

	// GetBinaryChunkProofs reads a slice of proofs from the chunk store, similar to GetChunkProofs.
	// Unlike GetChunkProofs, this method returns the raw serialized bytes of the proofs, as opposed to
	// deserializing them into encoding.Proof structs. The serialized proofs can be deserialized individually
	// via TODO.
	GetBinaryChunkProofs(ctx context.Context, blobKey corev2.BlobKey) ([][]byte, error)

	// GetChunkCoefficients reads a slice of frames from the chunk store. The metadata parameter
	// should match the metadata returned by PutChunkCoefficients.
	GetChunkCoefficients(
		ctx context.Context,
		blobKey corev2.BlobKey,
		fragmentInfo *encoding.FragmentInfo) ([]*rs.Frame, error)

	// GetBinaryChunkCoefficients reads a slice of frames from the chunk store, similar to GetChunkCoefficients.
	// Unlike GetChunkCoefficients, this method returns the raw serialized bytes of the frames, as opposed to
	// deserializing them into rs.Frame structs. The serialized frames can be deserialized individually
	// via rs.GnarkDecodeFrame.
	GetBinaryChunkCoefficients(
		ctx context.Context,
		blobKey corev2.BlobKey,
		fragmentInfo *encoding.FragmentInfo) ([][]byte, error)
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

func (r *chunkReader) GetChunkProofs(
	ctx context.Context,
	blobKey corev2.BlobKey) ([]*encoding.Proof, error) {

	bytes, err := r.client.DownloadObject(ctx, r.bucket, s3.ScopedProofKey(blobKey))
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

// TODO unit test new methods

func (r *chunkReader) GetBinaryChunkProofs(ctx context.Context, blobKey corev2.BlobKey) ([][]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (r *chunkReader) GetChunkCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
	fragmentInfo *encoding.FragmentInfo) ([]*rs.Frame, error) {

	bytes, err := r.client.FragmentedDownloadObject(
		ctx,
		r.bucket,
		s3.ScopedChunkKey(blobKey),
		int(fragmentInfo.TotalChunkSizeBytes),
		int(fragmentInfo.FragmentSizeBytes))

	if err != nil {
		r.logger.Error("Failed to download chunks from S3: %v", err)
		return nil, fmt.Errorf("failed to download chunks from S3: %w", err)
	}

	frames, err := rs.GnarkDecodeFrames(bytes)
	if err != nil {
		r.logger.Error("Failed to decode frames: %v", err)
		return nil, fmt.Errorf("failed to decode frames: %w", err)
	}

	return frames, nil
}

func (r *chunkReader) GetBinaryChunkCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
	fragmentInfo *encoding.FragmentInfo) ([][]byte, error) {

	bytes, err := r.client.FragmentedDownloadObject(
		ctx,
		r.bucket,
		s3.ScopedChunkKey(blobKey),
		int(fragmentInfo.TotalChunkSizeBytes),
		int(fragmentInfo.FragmentSizeBytes))

	if err != nil {
		r.logger.Error("Failed to download chunks from S3: %v", err)
		return nil, fmt.Errorf("failed to download chunks from S3: %w", err)
	}

	frames, err := rs.GnarkSplitBinaryFrames(bytes)
	if err != nil {
		r.logger.Error("Failed to split frames: %v", err)
		return nil, fmt.Errorf("failed to split frames: %w", err)
	}

	return frames, nil
}
