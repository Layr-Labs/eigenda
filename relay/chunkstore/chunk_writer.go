package chunkstore

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common/s3"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
)

// ChunkWriter writes chunks that can be read by ChunkReader.
type ChunkWriter interface {
	// PutFrameProofs writes a slice of proofs to the chunk store.
	PutFrameProofs(ctx context.Context, blobKey corev2.BlobKey, proofs []*encoding.Proof) error
	// PutFrameCoefficients writes a slice of frames to the chunk store.
	PutFrameCoefficients(
		ctx context.Context,
		blobKey corev2.BlobKey,
		frames []rs.FrameCoeffs) (*encoding.FragmentInfo, error)
	// ProofExists checks if the proofs for the blob key exist in the chunk store.
	ProofExists(ctx context.Context, blobKey corev2.BlobKey) bool
	// CoefficientsExists checks if the coefficients for the blob key exist in the chunk store.
	// Returns a bool indicating if the coefficients exist and fragment info.
	CoefficientsExists(ctx context.Context, blobKey corev2.BlobKey) bool
}

var _ ChunkWriter = (*chunkWriter)(nil)

type chunkWriter struct {
	s3Client   s3.S3Client
	bucketName string
}

// NewChunkWriter creates a new ChunkWriter.
func NewChunkWriter(
	s3Client s3.S3Client,
	bucketName string,
) ChunkWriter {

	return &chunkWriter{
		s3Client:   s3Client,
		bucketName: bucketName,
	}
}

func (c *chunkWriter) PutFrameProofs(ctx context.Context, blobKey corev2.BlobKey, proofs []*encoding.Proof) error {
	if len(proofs) == 0 {
		return fmt.Errorf("no proofs to upload")
	}

	bytes, err := encoding.SerializeFrameProofs(proofs)
	if err != nil {
		return fmt.Errorf("failed to encode proofs: %v", err)
	}
	err = c.s3Client.UploadObject(ctx, c.bucketName, s3.ScopedProofKey(blobKey), bytes)
	if err != nil {
		return fmt.Errorf("failed to upload chunk proofs to S3: %v", err)
	}

	return nil
}

func (c *chunkWriter) PutFrameCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
	frames []rs.FrameCoeffs) (*encoding.FragmentInfo, error) {
	if len(frames) == 0 {
		return nil, fmt.Errorf("no frames to upload")
	}
	bytes, err := rs.SerializeFrameCoeffsSlice(frames)
	if err != nil {
		return nil, fmt.Errorf("failed to encode frames: %v", err)
	}

	err = c.s3Client.UploadObject(ctx, c.bucketName, s3.ScopedChunkKey(blobKey), bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to upload chunk coefficients to S3: %v", err)
	}

	return &encoding.FragmentInfo{
		SymbolsPerFrame: uint32(len(frames[0])),
	}, nil
}

func (c *chunkWriter) ProofExists(ctx context.Context, blobKey corev2.BlobKey) bool {
	size, err := c.s3Client.HeadObject(ctx, c.bucketName, s3.ScopedProofKey(blobKey))
	if err == nil && size != nil && *size > 0 {
		return true
	}

	return false
}

func (c *chunkWriter) CoefficientsExists(ctx context.Context, blobKey corev2.BlobKey) bool {
	// TODO(ian-shim): check latency
	objs, err := c.s3Client.ListObjects(ctx, c.bucketName, s3.ScopedChunkKey(blobKey))
	if err != nil {
		return false
	}

	keys := make([]string, len(objs))
	totalSize := int64(0)
	for i, obj := range objs {
		keys[i] = obj.Key
		totalSize += int64(obj.Size)
	}

	return true
}
