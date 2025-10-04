package blobstore

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/oci"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// CreateObjectStorageClient creates an S3 client based on the backend configuration
func CreateObjectStorageClient(
	ctx context.Context,
	config Config,
	awsConfig aws.ClientConfig,
	logger logging.Logger) (s3.Client, error) {

	switch config.Backend {
	case S3Backend:
		client, err := s3.NewClient(ctx, awsConfig, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create S3 client: %w", err)
		}
		return client, nil
	case OCIBackend:
		ociConfig := oci.ObjectStorageConfig{
			BucketName:                  config.BucketName,
			FragmentParallelismConstant: awsConfig.FragmentParallelismConstant,
			FragmentParallelismFactor:   awsConfig.FragmentParallelismFactor,
		}
		client, err := oci.NewObjectStorageClient(ctx, ociConfig, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create OCI object storage client: %w", err)
		}
		return client, nil
	default:
		return nil, fmt.Errorf("unsupported object storage backend: %s", config.Backend)
	}
}
