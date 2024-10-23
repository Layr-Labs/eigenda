package dataplane

import "github.com/aws/aws-sdk-go/aws"

// S3Config is the configuration for an S3Client.
type S3Config struct {
	// The AWS configuration to use when interacting with S3.
	// Default uses the aws.Config default except for region which is set to "us-east-2".
	AWSConfig *aws.Config
	// The name of the S3 bucket to use. All data written to the S3Client will be written to this bucket.
	// This is a required field.
	Bucket string
	// If true then the bucket will be created if it does not already exist. If false and the bucket does not exist
	// then the S3Client will return an error when it is created. Default is false.
	AutoCreateBucket bool
	// The number of characters of the key to use as the prefix. A value of "3" for the key "ABCDEFG" would result in
	// the prefix "ABC". Default is 3.
	PrefixChars int
	// This framework utilizes a pool of workers to help upload/download files. This value specifies the number of
	// workers to use. Default is 32.
	Parallelism int
	// The capacity of the task channel. Default is 256. It is suggested that this value exceed the number of workers.
	TaskChannelCapacity int
}

// DefaultS3Config returns a new S3Config with default values.
func DefaultS3Config() *S3Config {
	return &S3Config{
		AWSConfig: &aws.Config{
			Region: aws.String("us-east-2"),
		},
		AutoCreateBucket:    false,
		PrefixChars:         3,
		Parallelism:         32,
		TaskChannelCapacity: 256,
	}
}
