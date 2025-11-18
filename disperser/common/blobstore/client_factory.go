package blobstore

import (
	"context"
	"fmt"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/s3"
	"github.com/Layr-Labs/eigenda/common/s3/aws"
	"github.com/Layr-Labs/eigenda/common/s3/oci"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// CreateObjectStorageClient creates an S3 client based on the backend configuration
func CreateObjectStorageClient(
	ctx context.Context,
	config Config,
	awsConfig commonaws.ClientConfig,
	logger logging.Logger) (s3.S3Client, error) {

	switch config.Backend {
	case S3Backend:
		client, err := aws.NewAwsS3Client(
			ctx,
			logger,
			awsConfig.EndpointURL,
			awsConfig.Region,
			awsConfig.FragmentParallelismFactor,
			awsConfig.FragmentParallelismConstant,
			awsConfig.AccessKey,
			awsConfig.SecretAccessKey,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create S3 client: %w", err)
		}
		return client, nil
	case OCIBackend:
		ociConfig := oci.ObjectStorageConfig{
			BucketName:                  config.BucketName,
			Namespace:                   config.OCINamespace,
			Region:                      config.OCIRegion,
			CompartmentID:               config.OCICompartmentID,
			FragmentParallelismConstant: awsConfig.FragmentParallelismConstant,
			FragmentParallelismFactor:   awsConfig.FragmentParallelismFactor,
		}
		client, err := oci.NewOciS3Client(ctx, ociConfig, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create OCI object storage client: %w", err)
		}
		return client, nil
	default:
		return nil, fmt.Errorf("unsupported object storage backend: %s", config.Backend)
	}
}
