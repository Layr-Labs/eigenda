package common

import (
	"github.com/urfave/cli"
)

type KMSKeyConfig struct {
	// Provider specifies the KMS provider: "aws" or "oci"
	Provider string
	
	// AWS KMS fields
	KeyID   string
	Region  string
	
	// OCI KMS fields
	KeyOCID           string
	CryptoEndpoint    string
	ManagementEndpoint string
	
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
		// AWS KMS flags
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-key-id"),
			Usage:    "AWS KMS key ID that stores the private key",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_KEY_ID"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-key-region"),
			Usage:    "AWS KMS key region",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_KEY_REGION"),
		},
		// OCI KMS flags
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-key-ocid"),
			Usage:    "OCI KMS key OCID that stores the private key",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_KEY_OCID"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-crypto-endpoint"),
			Usage:    "OCI KMS crypto endpoint for signing operations",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_CRYPTO_ENDPOINT"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, "kms-management-endpoint"),
			Usage:    "OCI KMS management endpoint for key retrieval",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "KMS_MANAGEMENT_ENDPOINT"),
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
		Provider:           ctx.String(PrefixFlag(flagPrefix, "kms-provider")),
		KeyID:              ctx.String(PrefixFlag(flagPrefix, "kms-key-id")),
		Region:             ctx.String(PrefixFlag(flagPrefix, "kms-key-region")),
		KeyOCID:            ctx.String(PrefixFlag(flagPrefix, "kms-key-ocid")),
		CryptoEndpoint:     ctx.String(PrefixFlag(flagPrefix, "kms-crypto-endpoint")),
		ManagementEndpoint: ctx.String(PrefixFlag(flagPrefix, "kms-management-endpoint")),
		Disable:            ctx.Bool(PrefixFlag(flagPrefix, "kms-key-disable")),
	}
}
