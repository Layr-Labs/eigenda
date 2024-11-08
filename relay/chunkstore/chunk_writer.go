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

	bytes, err := rs.GnarkEncodeFrames(frames)
	if err != nil {
		c.logger.Errorf("Failed to encode frames: %v", err)
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

// func (c *chunkWriter) headObject(ctx context.Context, blobKey corev2.BlobKey) (int, error) {
// 	return c.s3Client.HeadObject(ctx, c.bucketName, s3.ScopedChunkKey(blobKey))
// }
