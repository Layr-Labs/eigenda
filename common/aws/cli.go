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
	Region string `docs:"required"`
	// AccessKey to use when interacting with S3.
	AccessKey string `docs:"required"`
	// SecretAccessKey to use when interacting with S3.
	SecretAccessKey string `docs:"required"` // TODO (cody.littley): Change to *secret.Secret
	// EndpointURL of the S3 endpoint to use. If this is not set then the default AWS S3 endpoint will be used.
	EndpointURL string

	// This is a deprecated setting and can be ignored.
	FragmentParallelismFactor int // TODO (cody.littley): Remove
	// This is a deprecated setting and can be ignored.
	FragmentParallelismConstant int // TODO (cody.littley): Remove
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
	}
}

// DefaultClientConfig returns a new ClientConfig with default values.
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Region:                      "us-east-2",
		FragmentParallelismFactor:   8,
		FragmentParallelismConstant: 0,
	}
}
