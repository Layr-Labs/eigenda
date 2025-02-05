package s3

import (
	"github.com/urfave/cli/v2"
)

var (
	EndpointFlagName        = withFlagPrefix("endpoint")
	EnableTLSFlagName       = withFlagPrefix("enable-tls")
	CredentialTypeFlagName  = withFlagPrefix("credential-type")
	AccessKeyIDFlagName     = withFlagPrefix("access-key-id")     // #nosec G101
	AccessKeySecretFlagName = withFlagPrefix("access-key-secret") // #nosec G101
	BucketFlagName          = withFlagPrefix("bucket")
	PathFlagName            = withFlagPrefix("path")
	TimeoutFlagName         = withFlagPrefix("timeout")
)

func withFlagPrefix(s string) string {
	return "s3." + s
}

func withEnvPrefix(envPrefix, s string) []string {
	return []string{envPrefix + "_S3_" + s}
}

// CLIFlags ... used for S3 backend configuration
// category is used to group the flags in the help output (see https://cli.urfave.org/v2/examples/flags/#grouping)
func CLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     EndpointFlagName,
			Usage:    "endpoint for S3 storage",
			EnvVars:  withEnvPrefix(envPrefix, "ENDPOINT"),
			Category: category,
		},
		&cli.BoolFlag{
			Name:     EnableTLSFlagName,
			Usage:    "enable TLS connection to S3 endpoint",
			Value:    false,
			EnvVars:  withEnvPrefix(envPrefix, "ENABLE_TLS"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     CredentialTypeFlagName,
			Usage:    "the way to authenticate to S3, options are [iam, static, public]",
			EnvVars:  withEnvPrefix(envPrefix, "CREDENTIAL_TYPE"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     AccessKeyIDFlagName,
			Usage:    "access key id for S3 storage",
			EnvVars:  withEnvPrefix(envPrefix, "ACCESS_KEY_ID"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     AccessKeySecretFlagName,
			Usage:    "access key secret for S3 storage",
			EnvVars:  withEnvPrefix(envPrefix, "ACCESS_KEY_SECRET"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     BucketFlagName,
			Usage:    "bucket name for S3 storage",
			EnvVars:  withEnvPrefix(envPrefix, "BUCKET"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     PathFlagName,
			Usage:    "path for S3 storage",
			EnvVars:  withEnvPrefix(envPrefix, "PATH"),
			Category: category,
		},
		// &cli.DurationFlag{
		// 	Name:     TimeoutFlagName,
		// 	Usage:    "timeout for S3 storage operations (e.g. get, put)",
		// 	Value:    5 * time.Second,
		// 	EnvVars:  withEnvPrefix(envPrefix, "TIMEOUT"),
		// 	Category: category,
		// },
	}
}

func ReadConfig(ctx *cli.Context) Config {
	return Config{
		CredentialType:  StringToCredentialType(ctx.String(CredentialTypeFlagName)),
		Endpoint:        ctx.String(EndpointFlagName),
		EnableTLS:       ctx.Bool(EnableTLSFlagName),
		AccessKeyID:     ctx.String(AccessKeyIDFlagName),
		AccessKeySecret: ctx.String(AccessKeySecretFlagName),
		Bucket:          ctx.String(BucketFlagName),
		Path:            ctx.String(PathFlagName),
		// Timeout:         ctx.Duration(TimeoutFlagName),
	}
}
