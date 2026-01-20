package common

import (
	"strings"

	"github.com/urfave/cli"
)

type KMSKeyConfig struct {
	KeyID           string
	Region          string
	FallbackRegions []string
	FailFast        bool
	Disable         bool
}

func KMSWalletCLIFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-key-id"),
			Usage:    "KMS key ID that stores the private key",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_KEY_ID"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-key-region"),
			Usage:    "KMS key region",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_KEY_REGION"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-fallback-regions"),
			Usage:    "Comma-separated list of fallback KMS regions for multi-regional failover",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_FALLBACK_REGIONS"),
		},
		cli.BoolFlag{
			Name:     PrefixFlag(flagPrefix, "kms-fail-fast"),
			Usage:    "Fail immediately if primary KMS region fails, without trying fallback regions",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_FAIL_FAST"),
		},
		cli.BoolFlag{
			Name:     PrefixFlag(flagPrefix, "kms-key-disable"),
			Usage:    "Disable KMS wallet",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_KEY_DISABLE"),
		},
	}
}

func ReadKMSKeyConfig(ctx *cli.Context, flagPrefix string) KMSKeyConfig {
	fallbackRegionsStr := ctx.String(PrefixFlag(flagPrefix, "kms-fallback-regions"))
	var fallbackRegions []string
	if fallbackRegionsStr != "" {
		// Split comma-separated list and trim whitespace
		for _, region := range strings.Split(fallbackRegionsStr, ",") {
			trimmed := strings.TrimSpace(region)
			if trimmed != "" {
				fallbackRegions = append(fallbackRegions, trimmed)
			}
		}
	}

	return KMSKeyConfig{
		KeyID:           ctx.String(PrefixFlag(flagPrefix, "kms-key-id")),
		Region:          ctx.String(PrefixFlag(flagPrefix, "kms-key-region")),
		FallbackRegions: fallbackRegions,
		FailFast:        ctx.Bool(PrefixFlag(flagPrefix, "kms-fail-fast")),
		Disable:         ctx.Bool(PrefixFlag(flagPrefix, "kms-key-disable")),
	}
}
