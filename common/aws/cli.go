package aws

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
	"time"
)

var (
	RegionFlagName          = "aws.region"
	AccessKeyIdFlagName     = "aws.access-key-id"
	SecretAccessKeyFlagName = "aws.secret-access-key"
	EndpointURLFlagName     = "aws.endpoint-url"
)

type ClientConfig struct {
	// The region to use when interacting with S3. Default is "us-east-2".
	Region string
	// The access key to use when interacting with S3.
	AccessKey string
	// The secret key to use when interacting with S3.
	SecretAccessKey string
	// The URL of the S3 endpoint to use. If this is not set then the default AWS S3 endpoint will be used.
	EndpointURL string

	// The number of characters of the key to use as the prefix for fragmented files.
	// A value of "3" for the key "ABCDEFG" will result in the prefix "ABC". Default is 3.
	FragmentPrefixChars int
	// This framework utilizes a pool of workers to help upload/download files. A non-zero value for this parameter
	// adds a number of workers equal to the number of cores times this value. Default is 8. In general, the number
	// of workers here can be a lot larger than the number of cores because the workers will be blocked on I/O most
	// of the time.
	FragmentParallelismFactor int
	// This framework utilizes a pool of workers to help upload/download files. A non-zero value for this parameter
	// adds a constant number of workers. Default is 0.
	FragmentParallelismConstant int
	// If a single fragmented read takes longer than this value then the read will be aborted. Default is 30 seconds.
	FragmentReadTimeout time.Duration
	// If a single fragmented write takes longer than this value then the write will be aborted. Default is 30 seconds.
	FragmentWriteTimeout time.Duration
}

func ClientFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, RegionFlagName),
			Usage:    "AWS Region",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "AWS_REGION"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, AccessKeyIdFlagName),
			Usage:    "AWS Access Key Id",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "AWS_ACCESS_KEY_ID"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, SecretAccessKeyFlagName),
			Usage:    "AWS Secret Access Key",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "AWS_SECRET_ACCESS_KEY"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, EndpointURLFlagName),
			Usage:    "AWS Endpoint URL",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "AWS_ENDPOINT_URL"),
		},

		// TODO add flags for new args
	}
}

func ReadClientConfig(ctx *cli.Context, flagPrefix string) ClientConfig {
	return ClientConfig{
		Region:          ctx.GlobalString(common.PrefixFlag(flagPrefix, RegionFlagName)),
		AccessKey:       ctx.GlobalString(common.PrefixFlag(flagPrefix, AccessKeyIdFlagName)),
		SecretAccessKey: ctx.GlobalString(common.PrefixFlag(flagPrefix, SecretAccessKeyFlagName)),
		EndpointURL:     ctx.GlobalString(common.PrefixFlag(flagPrefix, EndpointURLFlagName)),
	}
}

// DefaultClientConfig returns a new ClientConfig with default values.
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Region:                      "us-east-2",
		FragmentPrefixChars:         3,
		FragmentParallelismFactor:   8,
		FragmentParallelismConstant: 0,
		FragmentReadTimeout:         30 * time.Second,
		FragmentWriteTimeout:        30 * time.Second,
	}
}
