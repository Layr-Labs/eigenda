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

// ChunkWriter writes chunks that can be read by ChunkReader.
type ChunkWriter interface {
	// PutChunkProofs writes a slice of proofs to the chunk store.
	PutChunkProofs(ctx context.Context, blobKey corev2.BlobKey, proofs []*encoding.Proof) error
	// PutChunkCoefficients writes a slice of frames to the chunk store.
	PutChunkCoefficients(
		ctx context.Context,
		blobKey corev2.BlobKey,
		frames []*rs.Frame) (*encoding.FragmentInfo, error)
	// ProofExists checks if the proofs for the blob key exist in the chunk store.
	ProofExists(ctx context.Context, blobKey corev2.BlobKey) bool
	// CoefficientsExists checks if the coefficients for the blob key exist in the chunk store.
	// Returns a bool indicating if the coefficients exist and fragment info.
	CoefficientsExists(ctx context.Context, blobKey corev2.BlobKey) (bool, *encoding.FragmentInfo)
}

var _ ChunkWriter = (*chunkWriter)(nil)

type chunkWriter struct {
	logger       logging.Logger
	s3Client     s3.Client
	bucketName   string
	fragmentSize int
}

// NewChunkWriter creates a new ChunkWriter.
func NewChunkWriter(
	logger logging.Logger,
	s3Client s3.Client,
	bucketName string,
	fragmentSize int) ChunkWriter {

	return &chunkWriter{
		logger:       logger,
		s3Client:     s3Client,
		bucketName:   bucketName,
		fragmentSize: fragmentSize,
	}
}

func (c *chunkWriter) PutChunkProofs(ctx context.Context, blobKey corev2.BlobKey, proofs []*encoding.Proof) error {
	if len(proofs) == 0 {
		return fmt.Errorf("no proofs to upload")
	}

	bytes := make([]byte, 0, bn254.SizeOfG1AffineCompressed*len(proofs))
	for _, proof := range proofs {
		proofBytes := proof.Bytes()
		bytes = append(bytes, proofBytes[:]...)
	}

	err := c.s3Client.UploadObject(ctx, c.bucketName, s3.ScopedProofKey(blobKey), bytes)
	if err != nil {
		c.logger.Errorf("Failed to upload chunk proofs to S3: %v", err)
		return fmt.Errorf("failed to upload chunk proofs to S3: %v", err)
	}

	return nil
}

func (c *chunkWriter) PutChunkCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
	frames []*rs.Frame) (*encoding.FragmentInfo, error) {
	if len(frames) == 0 {
		return nil, fmt.Errorf("no frames to upload")
	}
	bytes, err := rs.GnarkEncodeFrames(frames)
	if err != nil {
		c.logger.Error("Failed to encode frames", "err", err)
		return nil, fmt.Errorf("failed to encode frames: %v", err)
	}

	err = c.s3Client.FragmentedUploadObject(ctx, c.bucketName, s3.ScopedChunkKey(blobKey), bytes, c.fragmentSize)
	if err != nil {
		c.logger.Errorf("Failed to upload chunk coefficients to S3: %v", err)
		return nil, fmt.Errorf("failed to upload chunk coefficients to S3: %v", err)
	}

	return &encoding.FragmentInfo{
		TotalChunkSizeBytes: uint32(len(bytes)),
		FragmentSizeBytes:   uint32(c.fragmentSize),
	}, nil
}

func (c *chunkWriter) ProofExists(ctx context.Context, blobKey corev2.BlobKey) bool {
	size, err := c.s3Client.HeadObject(ctx, c.bucketName, s3.ScopedProofKey(blobKey))
	if err == nil && size != nil && *size > 0 {
		return true
	}

	return false
}

func (c *chunkWriter) CoefficientsExists(ctx context.Context, blobKey corev2.BlobKey) (bool, *encoding.FragmentInfo) {
	// TODO(ian-shim): check latency
	objs, err := c.s3Client.ListObjects(ctx, c.bucketName, s3.ScopedChunkKey(blobKey))
	if err != nil {
		return false, nil
	}

	keys := make([]string, len(objs))
	totalSize := int64(0)
	for i, obj := range objs {
		keys[i] = obj.Key
		totalSize += int64(obj.Size)
	}

	return s3.SortAndCheckAllFragmentsExist(keys), &encoding.FragmentInfo{
		TotalChunkSizeBytes: uint32(totalSize),
		FragmentSizeBytes:   uint32(c.fragmentSize),
	}
}
