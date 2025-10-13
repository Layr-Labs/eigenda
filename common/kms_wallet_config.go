package common

import (
	"github.com/urfave/cli"
)

type KMSKeyConfig struct {
	// Provider specifies the KMS provider: "aws", "oci", or "local"
	Provider string

	// Shared fields
	KeyID  string // AWS KMS key ID or OCI KMS key OCID
	Region string // AWS region (not used for OCI)

	// Local private key (used when Provider is "local")
	PrivateKeyHex string // Hex-encoded private key (with or without 0x prefix)

	Disable bool
}

func KMSWalletCLIFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-provider"),
			Usage:    "KMS provider: 'aws', 'oci', or 'local' (defaults to 'aws' for backward compatibility)",
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
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "private-key"),
			Usage:    "Private key in hex format (used when kms-provider is 'local')",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "PRIVATE_KEY"),
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
		Provider:      ctx.String(PrefixFlag(flagPrefix, "kms-provider")),
		KeyID:         ctx.String(PrefixFlag(flagPrefix, "kms-key-id")),
		Region:        ctx.String(PrefixFlag(flagPrefix, "kms-key-region")),
		PrivateKeyHex: ctx.String(PrefixFlag(flagPrefix, "private-key")),
		Disable:       ctx.Bool(PrefixFlag(flagPrefix, "kms-key-disable")),
	}
}
