package dataplane

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
	"time"
)

var _ S3Client = &s3Client{}

type s3Client struct {
	// config is the configuration for the client.
	config *S3Config
	// ctx is the context for the client.
	ctx context.Context
	// cancel is called to cancel the context.
	cancel context.CancelFunc
	// the AWS S3 handle to use to talk to S3.
	svc *s3.S3
	// tasks are placed into this channel to be processed by workers.
	tasks chan func()
	// this wait group is completed when all workers have finished.
	wg *sync.WaitGroup
}

// TODO add a timeout maybe?

// NewS3Client creates a new S3Client instance.
func NewS3Client(
	ctx context.Context,
	config *S3Config) (S3Client, error) {

	if config.Bucket == "" {
		return nil, errors.New("config.Bucket is required")
	}
	if config.Parallelism < 1 {
		return nil, errors.New("parameter config.Parallelism must be at least 1")
	}

	ctx, cancel := context.WithCancel(ctx)

	sess, err := session.NewSession(config.AWSConfig)
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)

	client := &s3Client{
		config: config,
		ctx:    ctx,
		cancel: cancel,
		tasks:  make(chan func()),
		svc:    svc,
		wg:     &sync.WaitGroup{},
	}

	err = client.createBucketIfNeeded()
	if err != nil {
		return nil, err
	}

	client.wg.Add(config.Parallelism)
	for i := 0; i < config.Parallelism; i++ {
		go client.worker()
	}

	return client, nil
}

func (s *s3Client) Upload(key string, data []byte, fragmentSize int, ttl time.Duration) error {
	fragments := BreakIntoFragments(key, data, s.config.PrefixChars, fragmentSize)
	resultChannel := make(chan error, len(fragments))

	expiryTime := time.Now().Add(ttl)

	for _, fragment := range fragments {
		s.tasks <- func() {
			s.writeTask(resultChannel, fragment.FragmentKey, fragment.Data, expiryTime)
		}
	}

	for range fragments {
		err := <-resultChannel
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *s3Client) Download(key string, fileSize int, fragmentSize int) ([]byte, error) {
	fragmentKeys := GetFragmentKeys(key, s.config.PrefixChars, GetFragmentCount(fileSize, fragmentSize))
	resultChannel := make(chan *readResult, len(fragmentKeys))

	for i, fragmentKey := range fragmentKeys {
		s.tasks <- func() {
			s.readTask(resultChannel, fragmentKey, i)
		}
	}

	fragments := make([]*Fragment, len(fragmentKeys))
	for i := 0; i < len(fragmentKeys); i++ {
		result := <-resultChannel
		if result.err != nil {
			return nil, result.err
		}
		fragments[result.fragment.Index] = result.fragment
	}

	return RecombineFragments(fragments)
}

// Close closes the S3 client.
func (s *s3Client) Close() error {
	s.ctx.Done()
	s.wg.Wait()
	return nil
}

// createBucketIfNeeded creates the bucket if it does not exist.
func (s *s3Client) createBucketIfNeeded() error {
	_, err := s.svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(s.config.Bucket),
	})
	if err == nil {
		// Bucket exists
		return nil
	} else if s.config.AutoCreateBucket == false {
		return err
	}

	_, err = s.svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(s.config.Bucket),
	})
	return err
}

// worker is a function that processes tasks until the context is cancelled.
func (s *s3Client) worker() {
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return
		case task := <-s.tasks:
			task()
		}
	}
}

// readResult is the result of a read task.
type readResult struct {
	fragment *Fragment
	err      error
}

// readTask reads a single file from S3.
func (s *s3Client) readTask(
	resultChannel chan *readResult,
	key string,
	index int) {

	result := &readResult{}
	defer func() {
		resultChannel <- result
	}()

	ret, err := s.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		result.err = err
		return
	}

	data := make([]byte, *ret.ContentLength)
	_, err = ret.Body.Read(data)

	if err != nil {
		result.err = err
		return
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
func (s *s3Client) writeTask(
	resultChannel chan error,
	key string,
	value []byte,
	expiryTime time.Time) {

	_, err := s.svc.PutObject(&s3.PutObjectInput{
		Bucket:  aws.String(s.config.Bucket),
		Key:     aws.String(key),
		Body:    bytes.NewReader(value),
		Expires: &expiryTime,
	})

	resultChannel <- err
}
