package chunkstore

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	eigens3 "github.com/Layr-Labs/eigenda/common/aws/s3"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// A wrapper around an S3 client for reading and writing chunks/proofs to/from S3.
type ChunkClient struct {
	// the S3 client to use for writing
	s3Client *s3.Client
	// the S3 bucket to write to
	bucket string
}

// NewChunkClient creates a new ChunkClient.
func NewChunkClient(
	awsUrl string,
	region string,
	accessKey string,
	secretAccessKey string,
) (*ChunkClient, error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if awsUrl != "" {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           awsUrl,
					SigningRegion: region,
				}, nil
			}

			// returning EndpointNotFoundError will allow the service to fallback to its default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

	options := [](func(*config.LoadOptions) error){
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRetryMode(aws.RetryModeStandard),
	}
	// If access key and secret access key are not provided, use the default credential provider
	if len(accessKey) > 0 && len(secretAccessKey) > 0 {
		options = append(options,
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(accessKey, secretAccessKey, "")))
	}
	awsConfig, err := config.LoadDefaultConfig(context.Background(), options...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &ChunkClient{
		s3Client: s3Client,
	}, nil
}

// Write frame proofs to S3.
func (c *ChunkClient) PutFrameProofs(
	ctx context.Context,
	blobKey corev2.BlobKey,
	proofs []*encoding.Proof,
) error {

	s3Key := eigens3.ScopedProofKey(blobKey)

	serialized, err := encoding.SerializeFrameProofs(proofs)
	if err != nil {
		return fmt.Errorf("failed to encode proofs: %v", err)
	}

	_, err = c.s3Client.PutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(s3Key),
			Body:   bytes.NewReader(serialized),
		})

	if err != nil {
		return fmt.Errorf("failed to upload chunk proofs to S3: %v", err)
	}

	return nil
}

// Write frame coefficients to S3.
func (c *ChunkClient) PutFrameCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
	frames []rs.FrameCoeffs,
) error {
	s3Key := eigens3.ScopedChunkKey(blobKey)

	serialized, err := rs.SerializeFrameCoeffsSlice(frames)
	if err != nil {
		return fmt.Errorf("failed to encode proofs: %v", err)
	}

	_, err = c.s3Client.PutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(s3Key),
			Body:   bytes.NewReader(serialized),
		})

	if err != nil {
		return fmt.Errorf("failed to upload chunks to S3: %v", err)
	}

	return nil
}

// Check to see if proofs exist in S3 for the given blob key.
func (c *ChunkClient) ProofExists(
	ctx context.Context,
	blobKey corev2.BlobKey,
) (bool, error) {

	s3Key := eigens3.ScopedProofKey(blobKey)

	_, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		var notFound *s3types.NotFound
		if ok := errors.As(err, &notFound); ok {
			return false, nil
		}
		return false, fmt.Errorf("failed to head object in S3: %v", err)
	}

	return true, nil

}

// Check to see if coefficients exist in S3 for the given blob key.
func (c *ChunkClient) CoefficientsExists(
	ctx context.Context,
	blobKey corev2.BlobKey,
) (bool, error) {

	s3Key := eigens3.ScopedChunkKey(blobKey)

	_, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		var notFound *s3types.NotFound
		if ok := errors.As(err, &notFound); ok {
			return false, nil
		}
		return false, fmt.Errorf("failed to head object in S3: %v", err)
	}

	return true, nil
}

// Read frame proofs from S3, returning them in serialized form.
func (c *ChunkClient) GetBinaryChunkProofs(
	ctx context.Context,
	blobKey corev2.BlobKey,
	firstIndex uint32,
	count uint32,
) ([][]byte, error) {

	// bytes, err := r.client.DownloadObject(ctx, r.bucket, s3.ScopedProofKey(blobKey))
	// if err != nil {
	// 	r.logger.Error("failed to download proofs from S3", "blob", blobKey.Hex(), "error", err)
	// 	return nil, fmt.Errorf("failed to download proofs from S3 for blob %s: %w", blobKey.Hex(), err)
	// }

	// proofs, err := encoding.SplitSerializedFrameProofs(bytes)
	// if err != nil {
	// 	r.logger.Error("failed to split proofs", "blob", blobKey.Hex(), "error", err)
	// 	return nil, fmt.Errorf("failed to split proofs for blob %s: %w", blobKey.Hex(), err)
	// }

	// return proofs, nil

	return nil, nil // TODO
}

// Read frame coefficients from S3, returning them in serialized form.
func (c *ChunkClient) GetBinaryChunkCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
	firstIndex uint32,
	count uint32,
) ([][]byte, error) {

	// bytes, err := r.client.FragmentedDownloadObject(
	// 	ctx,
	// 	r.bucket,
	// 	s3.ScopedChunkKey(blobKey),
	// 	int(fragmentInfo.TotalChunkSizeBytes),
	// 	int(fragmentInfo.FragmentSizeBytes))

	// if err != nil {
	// 	r.logger.Error("failed to download coefficients from S3",
	// 		"blob", blobKey.Hex(),
	// 		"totalSize", fragmentInfo.TotalChunkSizeBytes,
	// 		"fragmentSize", fragmentInfo.FragmentSizeBytes,
	// 		"error", err)
	// 	return 0, nil, fmt.Errorf("failed to download coefficients from S3 for blob %s (total size: %d, fragment size: %d): %w",
	// 		blobKey.Hex(), fragmentInfo.TotalChunkSizeBytes, fragmentInfo.FragmentSizeBytes, err)
	// }

	// elementCount, frames, err := rs.SplitSerializedFrameCoeffs(bytes)
	// if err != nil {
	// 	r.logger.Error("failed to split coefficient frames", "blob", blobKey.Hex(), "error", err)
	// 	return 0, nil, fmt.Errorf("failed to split coefficient frames for blob %s: %w", blobKey.Hex(), err)
	// }

	// return elementCount, frames, nil

	return nil, nil
}
