package common

import (
	"github.com/urfave/cli"
)

type KMSKeyConfig struct {
	// Provider specifies the KMS provider: "aws" or "oci"
	Provider string

	// Shared fields
	KeyID  string // AWS KMS key ID or OCI KMS key OCID
	Region string // AWS region (not used for OCI)

	Disable bool
}

func KMSWalletCLIFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-provider"),
			Usage:    "KMS provider: 'aws' or 'oci' (defaults to 'aws' for backward compatibility)",
			Value:    "aws",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_PROVIDER"),
		},
		// KMS flags (shared between providers)
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-key-id"),
			Usage:    "KMS key ID/OCID that stores the private key (AWS KMS key ID or OCI KMS key OCID)",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_KEY_ID"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-key-region"),
			Usage:    "AWS KMS key region (not used for OCI)",
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
		Provider: ctx.String(PrefixFlag(flagPrefix, "kms-provider")),
		KeyID:    ctx.String(PrefixFlag(flagPrefix, "kms-key-id")),
		Region:   ctx.String(PrefixFlag(flagPrefix, "kms-key-region")),
		Disable:  ctx.Bool(PrefixFlag(flagPrefix, "kms-key-disable")),
	}
}
