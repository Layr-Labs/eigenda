package aws

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

var (
	RegionFlagName                      = "aws.region"
	AccessKeyIdFlagName                 = "aws.access-key-id"
	SecretAccessKeyFlagName             = "aws.secret-access-key"
	EndpointURLFlagName                 = "aws.endpoint-url"
	FragmentPrefixCharsFlagName         = "aws.fragment-prefix-chars"
	FragmentParallelismFactorFlagName   = "aws.fragment-parallelism-factor"
	FragmentParallelismConstantFlagName = "aws.fragment-parallelism-constant"
	FragmentReadTimeoutFlagName         = "aws.fragment-read-timeout"
	FragmentWriteTimeoutFlagName        = "aws.fragment-write-timeout"
)

type ClientConfig struct {
	// Region is the region to use when interacting with S3. Default is "us-east-2".
	Region string
	// AccessKey to use when interacting with S3.
	AccessKey string
	// SecretAccessKey to use when interacting with S3.
	SecretAccessKey string
	// EndpointURL of the S3 endpoint to use. If this is not set then the default AWS S3 endpoint will be used.
	EndpointURL string

	// FragmentParallelismFactor helps determine the size of the pool of workers to help upload/download files.
	// A non-zero value for this parameter adds a number of workers equal to the number of cores times this value.
	// Default is 8. In general, the number of workers here can be a lot larger than the number of cores because the
	// workers will be blocked on I/O most of the time.
	FragmentParallelismFactor int
	// FragmentParallelismConstant helps determine the size of the pool of workers to help upload/download files.
	// A non-zero value for this parameter adds a constant number of workers. Default is 0.
	FragmentParallelismConstant int
	// FragmentReadTimeout is used to bound the maximum time to wait for a single fragmented read.
	// Default is 30 seconds.
	FragmentReadTimeout time.Duration
	// FragmentWriteTimeout is used to bound the maximum time to wait for a single fragmented write.
	// Default is 30 seconds.
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
		cli.IntFlag{
			Name:     common.PrefixFlag(flagPrefix, FragmentPrefixCharsFlagName),
			Usage:    "The number of characters of the key to use as the prefix for fragmented files",
			Required: false,
			Value:    3,
			EnvVar:   common.PrefixEnvVar(envPrefix, "FRAGMENT_PREFIX_CHARS"),
		},
		cli.IntFlag{
			Name:     common.PrefixFlag(flagPrefix, FragmentParallelismFactorFlagName),
			Usage:    "Add this many threads times the number of cores to the worker pool",
			Required: false,
			Value:    8,
			EnvVar:   common.PrefixEnvVar(envPrefix, "FRAGMENT_PARALLELISM_FACTOR"),
		},
		cli.IntFlag{
			Name:     common.PrefixFlag(flagPrefix, FragmentParallelismConstantFlagName),
			Usage:    "Add this many threads to the worker pool",
			Required: false,
			Value:    0,
			EnvVar:   common.PrefixEnvVar(envPrefix, "FRAGMENT_PARALLELISM_CONSTANT"),
		},
		cli.DurationFlag{
			Name:     common.PrefixFlag(flagPrefix, FragmentReadTimeoutFlagName),
			Usage:    "The maximum time to wait for a single fragmented read",
			Required: false,
			Value:    30 * time.Second,
			EnvVar:   common.PrefixEnvVar(envPrefix, "FRAGMENT_READ_TIMEOUT"),
		},
		cli.DurationFlag{
			Name:     common.PrefixFlag(flagPrefix, FragmentWriteTimeoutFlagName),
			Usage:    "The maximum time to wait for a single fragmented write",
			Required: false,
			Value:    30 * time.Second,
			EnvVar:   common.PrefixEnvVar(envPrefix, "FRAGMENT_WRITE_TIMEOUT"),
		},
	}
}

func ReadClientConfig(ctx *cli.Context, flagPrefix string) ClientConfig {
	return ClientConfig{
		Region:                      ctx.GlobalString(common.PrefixFlag(flagPrefix, RegionFlagName)),
		AccessKey:                   ctx.GlobalString(common.PrefixFlag(flagPrefix, AccessKeyIdFlagName)),
		SecretAccessKey:             ctx.GlobalString(common.PrefixFlag(flagPrefix, SecretAccessKeyFlagName)),
		EndpointURL:                 ctx.GlobalString(common.PrefixFlag(flagPrefix, EndpointURLFlagName)),
		FragmentParallelismFactor:   ctx.GlobalInt(common.PrefixFlag(flagPrefix, FragmentParallelismFactorFlagName)),
		FragmentParallelismConstant: ctx.GlobalInt(common.PrefixFlag(flagPrefix, FragmentParallelismConstantFlagName)),
		FragmentReadTimeout:         ctx.GlobalDuration(common.PrefixFlag(flagPrefix, FragmentReadTimeoutFlagName)),
		FragmentWriteTimeout:        ctx.GlobalDuration(common.PrefixFlag(flagPrefix, FragmentWriteTimeoutFlagName)),
	}
}

// DefaultClientConfig returns a new ClientConfig with default values.
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Region:                      "us-east-2",
		FragmentParallelismFactor:   8,
		FragmentParallelismConstant: 0,
		FragmentReadTimeout:         30 * time.Second,
		FragmentWriteTimeout:        30 * time.Second,
	}
}
