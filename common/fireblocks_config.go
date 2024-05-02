package common

import (
	"time"

	"github.com/urfave/cli"
)

const (
	FireblocksAPIKeyNameFlagName       = "fireblocks-api-key-name"
	FireblocksAPISecretNameFlagName    = "fireblocks-api-secret-name"
	FireblocksBaseURLFlagName          = "fireblocks-api-url"
	FireblocksVaultAccountNameFlagName = "fireblocks-vault-account-name"
	FireblocksWalletAddressFlagName    = "fireblocks-wallet-address"
	FireblocksSecretManagerRegion      = "fireblocks-secret-manager-region"
	FireblocksDisable                  = "fireblocks-disable"
	FireblocksAPITimeoutFlagName       = "fireblocks-api-timeout"
)

type FireblocksConfig struct {
	APIKeyName       string
	SecretKeyName    string
	BaseURL          string
	VaultAccountName string
	WalletAddress    string
	Region           string
	Disable          bool
	APITimeout       time.Duration
}

func FireblocksCLIFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, FireblocksAPIKeyNameFlagName),
			Usage:    "Fireblocks API Key Name. To configure Fireblocks MPC wallet, this field is required. Otherwise, private key must be configured in eth client so that it can fall back to private key wallet.",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "FIREBLOCKS_API_KEY_NAME"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, FireblocksAPISecretNameFlagName),
			Usage:    "Fireblocks API Secret Name. To configure Fireblocks MPC wallet, this field is required. Otherwise, private key must be configured in eth client so that it can fall back to private key wallet.",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "FIREBLOCKS_API_SECRET_NAME"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, FireblocksBaseURLFlagName),
			Usage:    "Fireblocks API URL. To configure Fireblocks MPC wallet, this field is required. Otherwise, private key must be configured in eth client so that it can fall back to private key wallet.",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "FIREBLOCKS_API_URL"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, FireblocksVaultAccountNameFlagName),
			Usage:    "Fireblocks Vault Account Name. To configure Fireblocks MPC wallet, this field is required. Otherwise, private key must be configured in eth client so that it can fall back to private key wallet.",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "FIREBLOCKS_VAULT_ACCOUNT_NAME"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, FireblocksWalletAddressFlagName),
			Usage:    "Fireblocks Wallet Address. To configure Fireblocks MPC wallet, this field is required. Otherwise, private key must be configured in eth client so that it can fall back to private key wallet.",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "FIREBLOCKS_WALLET_ADDRESS"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, FireblocksSecretManagerRegion),
			Usage:    "Fireblocks AWS Secret Manager Region.",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "FIREBLOCKS_SECRET_MANAGER_REGION"),
		},
		cli.BoolFlag{
			Name:     PrefixFlag(flagPrefix, FireblocksDisable),
			Usage:    "Disable Fireblocks. By default, Disable is set to false.",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "FIREBLOCKS_DISABLE"),
		},
		cli.DurationFlag{
			Name:     PrefixFlag(flagPrefix, FireblocksAPITimeoutFlagName),
			Usage:    "Timeout for Fireblocks API requests",
			Required: false,
			Value:    2 * time.Minute,
			EnvVar:   PrefixEnvVar(envPrefix, "FIREBLOCKS_API_TIMEOUT"),
		},
	}
}

func ReadFireblocksCLIConfig(ctx *cli.Context, flagPrefix string) FireblocksConfig {
	return FireblocksConfig{
		APIKeyName:       ctx.GlobalString(PrefixFlag(flagPrefix, FireblocksAPIKeyNameFlagName)),
		SecretKeyName:    ctx.GlobalString(PrefixFlag(flagPrefix, FireblocksAPISecretNameFlagName)),
		BaseURL:          ctx.GlobalString(PrefixFlag(flagPrefix, FireblocksBaseURLFlagName)),
		VaultAccountName: ctx.GlobalString(PrefixFlag(flagPrefix, FireblocksVaultAccountNameFlagName)),
		WalletAddress:    ctx.GlobalString(PrefixFlag(flagPrefix, FireblocksWalletAddressFlagName)),
		Region:           ctx.GlobalString(PrefixFlag(flagPrefix, FireblocksSecretManagerRegion)),
		Disable:          ctx.GlobalBool(PrefixFlag(flagPrefix, FireblocksDisable)),
		APITimeout:       ctx.GlobalDuration(PrefixFlag(flagPrefix, FireblocksAPITimeoutFlagName)),
	}
}
