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
		return s3.NewClient(ctx, awsConfig, logger)
	case OCIBackend:
		ociConfig := oci.ObjectStorageConfig{
			Namespace:                   config.OCINamespace,
			Region:                      config.OCIRegion,
			CompartmentID:               config.OCICompartmentID,
			BucketName:                  config.BucketName,
			FragmentParallelismConstant: awsConfig.FragmentParallelismConstant,
			FragmentParallelismFactor:   awsConfig.FragmentParallelismFactor,
		}
		return oci.NewObjectStorageClient(ctx, ociConfig, logger)
	default:
		return nil, fmt.Errorf("unsupported object storage backend: %s", config.Backend)
	}
}