package aws

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"

	s3common "github.com/Layr-Labs/eigenda/common/s3"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/sync/errgroup"
)

const (
	defaultBlobBufferSizeByte = 128 * 1024
)

var (
	once sync.Once
	ref  *awsS3Client
)

// An implementation of s3common.S3Client using AWS S3.
type awsS3Client struct {
	logger logging.Logger

	// Amazon's S3 client implementation.
	s3Client *s3.Client

	// concurrencyLimiter is a channel that limits the number of concurrent operations.
	concurrencyLimiter chan struct{}
}

var _ s3common.S3Client = (*awsS3Client)(nil)

// NewAwsS3Client creates a new S3Client that talks to AWS S3.
func NewAwsS3Client(
	ctx context.Context,
	logger logging.Logger,
	endpointUrl string,
	region string,
	fragmentParallelismFactor int,
	fragmentParallelismConstant int,
	accessKey string,
	secretAccessKey string,
) (s3common.S3Client, error) {

	var err error
	once.Do(func() {
		customResolver := aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if endpointUrl != "" {
					return aws.Endpoint{
						PartitionID:   "aws",
						URL:           endpointUrl,
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
		awsConfig, errCfg := config.LoadDefaultConfig(context.Background(), options...)

		if errCfg != nil {
			err = errCfg
			return
		}

		s3Client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
			o.UsePathStyle = true
		})

		workers := 0
		if fragmentParallelismConstant > 0 {
			workers = fragmentParallelismConstant
		}
		if fragmentParallelismFactor > 0 {
			workers = fragmentParallelismFactor * runtime.NumCPU()
		}

		if workers == 0 {
			workers = 1
		}

		pool := &errgroup.Group{}
		pool.SetLimit(workers)

		ref = &awsS3Client{
			s3Client:           s3Client,
			concurrencyLimiter: make(chan struct{}, workers),
			logger:             logger.With("component", "S3Client"),
		}
	})
	return ref, err
}

func (s *awsS3Client) DownloadObject(ctx context.Context, bucket string, key string) ([]byte, bool, error) {
	objectSize := defaultBlobBufferSizeByte
	size, err := s.HeadObject(ctx, bucket, key)
	if err == nil {
		objectSize = int(*size)
	}
	buffer := manager.NewWriteAtBuffer(make([]byte, 0, objectSize))

	var partMiBs int64 = 10
	downloader := manager.NewDownloader(s.s3Client, func(d *manager.Downloader) {
		// 10MB per part
		d.PartSize = partMiBs * 1024 * 1024
		// The number of goroutines to spin up in parallel per call to Upload when sending parts
		d.Concurrency = 3
	})

	_, err = downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if errors.Is(err, &types.NoSuchKey{}) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to download object: %w", err)
	}

	if buffer == nil || len(buffer.Bytes()) == 0 {
		return nil, false, nil
	}

	if len(buffer.Bytes()) != objectSize {
		return nil, false, fmt.Errorf("downloaded object size (%d) does not match expected size (%d)",
			len(buffer.Bytes()), objectSize)
	}

	return buffer.Bytes(), true, nil
}

func (s *awsS3Client) DownloadPartialObject(
	ctx context.Context,
	bucket string,
	key string,
	startIndex int64,
	endIndex int64) ([]byte, bool, error) {

	if startIndex < 0 || endIndex <= startIndex {
		return nil, false, fmt.Errorf("invalid startIndex (%d) or endIndex (%d)", startIndex, endIndex)
	}

	rangeHeader := fmt.Sprintf("bytes=%d-%d", startIndex, endIndex-1)

	buffer := manager.NewWriteAtBuffer(make([]byte, 0, endIndex-startIndex))

	var partMiBs int64 = 10
	downloader := manager.NewDownloader(s.s3Client, func(d *manager.Downloader) {
		// 10MB per part
		d.PartSize = partMiBs * 1024 * 1024
		// The number of goroutines to spin up in parallel per call to download when sending parts
		d.Concurrency = 3
	})

	_, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Range:  aws.String(rangeHeader),
	})
	if err != nil {
		if errors.Is(err, &types.NoSuchKey{}) {
			return nil, false, s3common.ErrObjectNotFound
		}
		return nil, false, fmt.Errorf("failed to download partial object: %w", err)
	}

	if buffer == nil || len(buffer.Bytes()) == 0 {
		return nil, false, nil
	}

	return buffer.Bytes(), true, nil
}

func (s *awsS3Client) HeadObject(ctx context.Context, bucket string, key string) (*int64, error) {
	output, err := s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var notFound *types.NotFound
		if ok := errors.As(err, &notFound); ok {
			return nil, s3common.ErrObjectNotFound
		}
		return nil, err
	}

	return output.ContentLength, nil
}

func (s *awsS3Client) UploadObject(ctx context.Context, bucket string, key string, data []byte) error {
	var partMiBs int64 = 10
	uploader := manager.NewUploader(s.s3Client, func(u *manager.Uploader) {
		// 10MiB per part
		u.PartSize = partMiBs * 1024 * 1024
		// The number of goroutines to spin up in parallel per call to upload when sending parts
		u.Concurrency = 3
	})

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *awsS3Client) DeleteObject(ctx context.Context, bucket string, key string) error {
	_, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return err
}

// ListObjects lists all items metadata in a bucket with the given prefix up to 1000 items.
func (s *awsS3Client) ListObjects(ctx context.Context, bucket string, prefix string) ([]s3common.ListedObject, error) {
	output, err := s.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	objects := make([]s3common.ListedObject, 0, len(output.Contents))
	for _, object := range output.Contents {
		var size int64 = 0
		if object.Size != nil {
			size = *object.Size
		}
		objects = append(objects, s3common.ListedObject{
			Key:  *object.Key,
			Size: size,
		})
	}
	return objects, nil
}

func (s *awsS3Client) CreateBucket(ctx context.Context, bucket string) error {
	_, err := s.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	return nil
}
