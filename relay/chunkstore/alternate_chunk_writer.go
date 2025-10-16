package chunkstore

import (
	"context"
	"fmt"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	eigens3 "github.com/Layr-Labs/eigenda/common/aws/s3"
)

type AlternateChunkWriter struct {
	s3Client *s3.Client
}

func NewAlternateChunkWriter(
	awsUrl string,
	region string,
	accessKey string,
	secretAccessKey string,
) (*AlternateChunkWriter, error) {

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

	return &AlternateChunkWriter{
		s3Client: s3Client,
	}, nil
}

// Write frame proofs to S3.
func (a *AlternateChunkWriter) PutFrameProofs(
	ctx context.Context,
	blobKey corev2.BlobKey,
	proofs []*encoding.Proof,
) error {

	blobKey := eigens3.ScopedProofKey(blobKey)

	bytes, err := encoding.SerializeFrameProofs(proofs)
	if err != nil {
		return fmt.Errorf("failed to encode proofs: %v", err)
	}

	return nil // TODO

}

// WRite frame coefficients to S3.
func (a *AlternateChunkWriter) PutFrameCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
	frames []rs.FrameCoeffs,
) error {

	return nil // TODO

}

// Check to see if proofs exist in S3 for the given blob key.
func (a *AlternateChunkWriter) ProofExists(
	ctx context.Context,
	blobKey corev2.BlobKey,
) (bool, error) {

	return false, nil // TODO

}

// Check to see if coefficients exist in S3 for the given blob key.
func (a *AlternateChunkWriter) CoefficientsExists(
	ctx context.Context,
	blobKey corev2.BlobKey,
) (bool, error) {

	return false, nil // TODO

}
