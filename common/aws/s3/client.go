package s3

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"sync"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/sync/errgroup"
)

var (
	once              sync.Once
	ref               *client
	ErrObjectNotFound = errors.New("object not found")
)

type Object struct {
	Key  string
	Size int64
}

type client struct {
	cfg      *commonaws.ClientConfig
	s3Client *s3.Client

	// concurrencyLimiter is a channel that limits the number of concurrent operations.
	concurrencyLimiter chan struct{}

	logger logging.Logger
}

var _ Client = (*client)(nil)

func NewClient(ctx context.Context, cfg commonaws.ClientConfig, logger logging.Logger) (Client, error) {
	var err error
	once.Do(func() {
		customResolver := aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if cfg.EndpointURL != "" {
					return aws.Endpoint{
						PartitionID:   "aws",
						URL:           cfg.EndpointURL,
						SigningRegion: cfg.Region,
					}, nil
				}

				// returning EndpointNotFoundError will allow the service to fallback to its default resolution
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			})

		options := [](func(*config.LoadOptions) error){
			config.WithRegion(cfg.Region),
			config.WithEndpointResolverWithOptions(customResolver),
			config.WithRetryMode(aws.RetryModeStandard),
		}
		// If access key and secret access key are not provided, use the default credential provider
		if len(cfg.AccessKey) > 0 && len(cfg.SecretAccessKey) > 0 {
			options = append(options,
				config.WithCredentialsProvider(
					credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretAccessKey, "")))
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
		if cfg.FragmentParallelismConstant > 0 {
			workers = cfg.FragmentParallelismConstant
		}
		if cfg.FragmentParallelismFactor > 0 {
			workers = cfg.FragmentParallelismFactor * runtime.NumCPU()
		}

		if workers == 0 {
			workers = 1
		}

		pool := &errgroup.Group{}
		pool.SetLimit(workers)

		ref = &client{
			cfg:                &cfg,
			s3Client:           s3Client,
			concurrencyLimiter: make(chan struct{}, workers),
			logger:             logger.With("component", "S3Client"),
		}
	})
	return ref, err
}

func (s *client) DownloadObject(ctx context.Context, bucket string, key string) ([]byte, error) {
	var partMiBs int64 = 10
	downloader := manager.NewDownloader(s.s3Client, func(d *manager.Downloader) {
		d.PartSize = partMiBs * 1024 * 1024 // 10MB per part
		d.Concurrency = 3                   //The number of goroutines to spin up in parallel per call to Upload when sending parts
	})

	buffer := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	if buffer == nil || len(buffer.Bytes()) == 0 {
		return nil, ErrObjectNotFound
	}

	return buffer.Bytes(), nil
}

func (s *client) HeadObject(ctx context.Context, bucket string, key string) (*int64, error) {
	output, err := s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var notFound *types.NotFound
		if ok := errors.As(err, &notFound); ok {
			return nil, ErrObjectNotFound
		}
		return nil, err
	}

	return output.ContentLength, nil
}

func (s *client) UploadObject(ctx context.Context, bucket string, key string, data []byte) error {
	var partMiBs int64 = 10
	uploader := manager.NewUploader(s.s3Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024 // 10MiB per part
		u.Concurrency = 3                   //The number of goroutines to spin up in parallel per call to upload when sending parts
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

func (s *client) DeleteObject(ctx context.Context, bucket string, key string) error {
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
func (s *client) ListObjects(ctx context.Context, bucket string, prefix string) ([]Object, error) {
	output, err := s.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	objects := make([]Object, 0, len(output.Contents))
	for _, object := range output.Contents {
		var size int64 = 0
		if object.Size != nil {
			size = *object.Size
		}
		objects = append(objects, Object{
			Key:  *object.Key,
			Size: size,
		})
	}
	return objects, nil
}

func (s *client) CreateBucket(ctx context.Context, bucket string) error {
	_, err := s.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *client) FragmentedUploadObject(
	ctx context.Context,
	bucket string,
	key string,
	data []byte,
	fragmentSize int) error {

	fragments, err := BreakIntoFragments(key, data, fragmentSize)
	if err != nil {
		return err
	}
	resultChannel := make(chan error, len(fragments))

	ctx, cancel := context.WithTimeout(ctx, s.cfg.FragmentWriteTimeout)
	defer cancel()

	for _, fragment := range fragments {
		fragmentCapture := fragment
		s.concurrencyLimiter <- struct{}{}
		go func() {
			defer func() {
				<-s.concurrencyLimiter
			}()
			s.fragmentedWriteTask(ctx, resultChannel, fragmentCapture, bucket)
		}()
	}

	for range fragments {
		err = <-resultChannel
		if err != nil {
			return err
		}
	}
	return ctx.Err()
}

// fragmentedWriteTask writes a single file to S3.
func (s *client) fragmentedWriteTask(
	ctx context.Context,
	resultChannel chan error,
	fragment *Fragment,
	bucket string) {

	_, err := s.s3Client.PutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fragment.FragmentKey),
			Body:   bytes.NewReader(fragment.Data),
		})

	resultChannel <- err
}

func (s *client) FragmentedDownloadObject(
	ctx context.Context,
	bucket string,
	key string,
	fileSize int,
	fragmentSize int) ([]byte, error) {
	if fileSize <= 0 {
		return nil, errors.New("fileSize must be greater than 0")
	}

	if fragmentSize <= 0 {
		return nil, errors.New("fragmentSize must be greater than 0")
	}

	fragmentKeys, err := GetFragmentKeys(key, getFragmentCount(fileSize, fragmentSize))
	if err != nil {
		return nil, err
	}
	resultChannel := make(chan *readResult, len(fragmentKeys))

	ctx, cancel := context.WithTimeout(ctx, s.cfg.FragmentWriteTimeout)
	defer cancel()

	for i, fragmentKey := range fragmentKeys {
		boundFragmentKey := fragmentKey
		boundI := i
		s.concurrencyLimiter <- struct{}{}
		go func() {
			defer func() {
				<-s.concurrencyLimiter
			}()
			s.readTask(ctx, resultChannel, bucket, boundFragmentKey, boundI)
		}()
	}

	fragments := make([]*Fragment, len(fragmentKeys))
	for i := 0; i < len(fragmentKeys); i++ {
		result := <-resultChannel
		if result.err != nil {
			return nil, result.err
		}
		fragments[result.fragment.Index] = result.fragment
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return recombineFragments(fragments)

}

// readResult is the result of a read task.
type readResult struct {
	fragment *Fragment
	err      error
}

// readTask reads a single file from S3.
func (s *client) readTask(
	ctx context.Context,
	resultChannel chan *readResult,
	bucket string,
	key string,
	index int) {

	result := &readResult{}
	defer func() {
		resultChannel <- result
	}()

	ret, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		result.err = err
		return
	}

	data := make([]byte, *ret.ContentLength)
	bytesRead := 0

	for bytesRead < len(data) && ctx.Err() == nil {
		count, err := ret.Body.Read(data[bytesRead:])
		if err != nil && err.Error() != "EOF" {
			result.err = err
			return
		}
		bytesRead += count
	}

	result.fragment = &Fragment{
		FragmentKey: key,
		Data:        data,
		Index:       index,
	}

	err = ret.Body.Close()
	if err != nil {
		result.err = err
	}
}
