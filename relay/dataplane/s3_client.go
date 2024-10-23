package dataplane

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
)

var _ S3Client = &s3Client{}

type s3Client struct {
	// config is the configuration for the client.
	config *S3Config
	// ctx is the context for the client.
	ctx context.Context
	// cancel is called to cancel the context.
	cancel context.CancelFunc
	// the AWS S3 session to use to talk to S3.
	session *s3.S3
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

	if config.Bucket == nil {
		return nil, errors.New("config.Bucket is required")
	}
	if config.Parallelism < 1 {
		return nil, errors.New("parameter config.Parallelism must be at least 1")
	}

	ctx, cancel := context.WithCancel(ctx)

	client := &s3Client{
		config: config,
		ctx:    ctx,
		cancel: cancel,
		tasks:  make(chan func()),
	}

	for i := 0; i < config.Parallelism; i++ {
		go client.worker()
	}

	return client, nil
}

func (s *s3Client) Upload(key string, data []byte, fragmentSize int) error {
	fragments := BreakIntoFragments(key, data, s.config.PrefixChars, fragmentSize)
	resultChannel := make(chan error, len(fragments))

	for _, fragment := range fragments {
		s.tasks <- func() {
			s.writeTask(resultChannel, fragment.FragmentKey, fragment.Data)
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

// worker is a function that processes tasks until the context is cancelled.
func (s *s3Client) worker() {
	s.wg.Add(1)
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

	ret, err := s.session.GetObject(&s3.GetObjectInput{
		Bucket: s.config.Bucket,
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
	value []byte) {

	_, err := s.session.PutObject(&s3.PutObjectInput{
		Bucket: s.config.Bucket,
		Key:    aws.String(key),
		Body:   bytes.NewReader(value),
	})

	resultChannel <- err
}
