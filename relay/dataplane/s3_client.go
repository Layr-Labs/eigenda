package dataplane

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gammazero/workerpool"
	"runtime"
)

var _ S3Client = &client{}

type client struct {
	// config is the configuration for the client.
	config *S3Config
	// ctx is the context for the client.
	ctx context.Context
	// cancel is called to cancel the context.
	cancel context.CancelFunc
	// the S3 client to use.
	client *s3.Client
	pool   *workerpool.WorkerPool
}

// NewS3Client creates a new S3Client instance.
func NewS3Client(
	ctx context.Context,
	cfg *S3Config) (S3Client, error) {

	if cfg.Bucket == "" {
		return nil, errors.New("config.Bucket is required")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.EndpointURL != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           cfg.EndpointURL,
				SigningRegion: cfg.Region,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fall back to its default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	options := [](func(*config.LoadOptions) error){
		config.WithRegion(cfg.Region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRetryMode(aws.RetryModeStandard),
	}
	// If access key and secret access key are not provided, use the default credential provider
	if len(cfg.AccessKey) > 0 && len(cfg.SecretAccessKey) > 0 {
		options = append(options, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretAccessKey, "")))
	}
	awsConfig, err := config.LoadDefaultConfig(context.Background(), options...)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	ctx, cancel := context.WithCancel(ctx)

	workers := 0
	if cfg.ParallelismConstant > 0 {
		workers = cfg.ParallelismConstant
	}
	if cfg.ParallelismFactor > 0 {
		workers = cfg.ParallelismFactor * runtime.NumCPU()
	}

	if workers == 0 {
		workers = 1
	}
	pool := workerpool.New(workers)

	c := &client{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
		client: s3Client,
		pool:   pool,
	}

	err = c.createBucketIfNeeded()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *client) Upload(key string, data []byte, fragmentSize int) error {
	fragments, err := BreakIntoFragments(key, data, s.config.PrefixChars, fragmentSize)
	if err != nil {
		return err
	}
	resultChannel := make(chan error, len(fragments))

	ctx, cancel := context.WithTimeout(s.ctx, s.config.WriteTimeout)
	defer cancel()

	for _, fragment := range fragments {
		fragmentCapture := fragment
		s.pool.Submit(func() {
			s.writeTask(ctx, resultChannel, fragmentCapture)
		})
	}

	for range fragments {
		err := <-resultChannel
		if err != nil {
			return err
		}
	}
	return ctx.Err()
}

func (s *client) Download(key string, fileSize int, fragmentSize int) ([]byte, error) {
	if fragmentSize <= 0 {
		return nil, errors.New("fragmentSize must be greater than 0")
	}

	fragmentKeys, err := GetFragmentKeys(key, s.config.PrefixChars, GetFragmentCount(fileSize, fragmentSize))
	if err != nil {
		return nil, err
	}
	resultChannel := make(chan *readResult, len(fragmentKeys))

	ctx, cancel := context.WithTimeout(s.ctx, s.config.WriteTimeout)
	defer cancel()

	for i, fragmentKey := range fragmentKeys {
		boundFragmentKey := fragmentKey
		boundI := i
		s.pool.Submit(func() {
			s.readTask(ctx, resultChannel, boundFragmentKey, boundI)
		})
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

	return RecombineFragments(fragments)
}

// Close closes the S3 client.
func (s *client) Close() error {
	s.cancel()
	s.pool.StopWait()
	return nil
}

// createBucketIfNeeded creates the bucket if it does not exist.
func (s *client) createBucketIfNeeded() error {
	if !s.config.AutoCreateBucket {
		return nil
	}

	listBucketsOutput, err := s.client.ListBuckets(s.ctx, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	for _, bucket := range listBucketsOutput.Buckets {
		if *bucket.Name == s.config.Bucket {
			return nil
		}
	}

	_, err = s.client.CreateBucket(s.ctx,
		&s3.CreateBucketInput{
			Bucket: aws.String(s.config.Bucket),
		})

	return err
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
	key string,
	index int) {

	result := &readResult{}
	defer func() {
		resultChannel <- result
	}()

	ret, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
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

// writeTask writes a single file to S3.
func (s *client) writeTask(
	ctx context.Context,
	resultChannel chan error,
	fragment *Fragment) {

	_, err := s.client.PutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(s.config.Bucket),
			Key:    aws.String(fragment.FragmentKey),
			Body:   bytes.NewReader(fragment.Data),
		})

	resultChannel <- err
}