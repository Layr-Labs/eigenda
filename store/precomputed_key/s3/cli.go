package s3

import (
	"time"

	"github.com/urfave/cli/v2"
)

var (
	EndpointFlagName        = withFlagPrefix("endpoint")
	CredentialTypeFlagName  = withFlagPrefix("credential-type")
	AccessKeyIDFlagName     = withFlagPrefix("access-key-id")     // #nosec G101
	AccessKeySecretFlagName = withFlagPrefix("access-key-secret") // #nosec G101
	BucketFlagName          = withFlagPrefix("bucket")
	PathFlagName            = withFlagPrefix("path")
	BackupFlagName          = withFlagPrefix("backup")
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
			EnvVars:  withEnvPrefix(envPrefix, "S3_ENDPOINT"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     CredentialTypeFlagName,
			Usage:    "The way to authenticate to S3, options are [iam, static]",
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
		&cli.BoolFlag{
			Name:     BackupFlagName,
			Usage:    "whether to use S3 as a backup store to ensure resiliency in case of EigenDA read failure",
			Value:    false,
			EnvVars:  withEnvPrefix(envPrefix, "BACKUP"),
			Category: category,
		},
		&cli.DurationFlag{
			Name:     TimeoutFlagName,
			Usage:    "timeout for S3 storage operations (e.g. get, put)",
			Value:    5 * time.Second,
			EnvVars:  withEnvPrefix(envPrefix, "TIMEOUT"),
			Category: category,
		},
	}
}

func ReadConfig(ctx *cli.Context) Config {
	return Config{
		CredentialType:  StringToCredentialType(ctx.String(CredentialTypeFlagName)),
		Endpoint:        ctx.String(EndpointFlagName),
		AccessKeyID:     ctx.String(AccessKeyIDFlagName),
		AccessKeySecret: ctx.String(AccessKeySecretFlagName),
		Bucket:          ctx.String(BucketFlagName),
		Path:            ctx.String(PathFlagName),
		Backup:          ctx.Bool(BackupFlagName),
		Timeout:         ctx.Duration(TimeoutFlagName),
	}
}
