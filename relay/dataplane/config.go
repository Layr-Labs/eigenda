package dataplane

import (
	"time"
)

// S3Config is the configuration for an S3Client.
type S3Config struct {
	// The URL of the S3 endpoint to use. If this is not set then the default AWS S3 endpoint will be used.
	EndpointURL string
	// The region to use when interacting with S3. Default is "us-east-2".
	Region string
	// The access key to use when interacting with S3.
	AccessKey string
	// The secret key to use when interacting with S3.
	SecretAccessKey string
	// The name of the S3 bucket to use. All data written to the S3Client will be written to this bucket.
	// This is a required field.
	Bucket string
	// If true then the bucket will be created if it does not already exist. If false and the bucket does not exist
	// then the S3Client will return an error when it is created. Default is false.
	AutoCreateBucket bool
	// The number of characters of the key to use as the prefix. A value of "3" for the key "ABCDEFG" would result in
	// the prefix "ABC". Default is 3.
	PrefixChars int
	// This framework utilizes a pool of workers to help upload/download files. A non-zero value for this parameter
	// adds a number of workers equal to the number of cores times this value. Default is 8. In general, the number
	// of workers here can be a lot larger than the number of cores because the workers will be blocked on I/O most
	// of the time.
	ParallelismFactor int
	// This framework utilizes a pool of workers to help upload/download files. A non-zero value for this parameter
	// adds a constant number of workers. Default is 0.
	ParallelismConstant int
	// The capacity of the task channel. Default is 256. It is suggested that this value exceed the number of workers.
	TaskChannelCapacity int
	// If a single read takes longer than this value then the read will be aborted. Default is 30 seconds.
	ReadTimeout time.Duration
	// If a single write takes longer than this value then the write will be aborted. Default is 30 seconds.
	WriteTimeout time.Duration
}

// DefaultS3Config returns a new S3Config with default values.
func DefaultS3Config() *S3Config {
	return &S3Config{
		Region:              "us-east-2",
		AutoCreateBucket:    false,
		PrefixChars:         3,
		ParallelismFactor:   8,
		ParallelismConstant: 0,
		TaskChannelCapacity: 256,
		ReadTimeout:         30 * time.Second,
		WriteTimeout:        30 * time.Second,
	}
}
