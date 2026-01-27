package common

import (
	"github.com/urfave/cli"
)

type KMSKeyConfig struct {
	KeyID   string
	Region  string
	Disable bool
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
		cli.BoolFlag{
			Name:     PrefixFlag(flagPrefix, "kms-key-disable"),
			Usage:    "Disable KMS wallet",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_KEY_DISABLE"),
		},
	}
}

func ReadKMSKeyConfig(ctx *cli.Context, flagPrefix string) KMSKeyConfig {
	return KMSKeyConfig{
		KeyID:   ctx.String(PrefixFlag(flagPrefix, "kms-key-id")),
		Region:  ctx.String(PrefixFlag(flagPrefix, "kms-key-region")),
		Disable: ctx.Bool(PrefixFlag(flagPrefix, "kms-key-disable")),
	}
}
