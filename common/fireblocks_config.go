package common

import (
	"github.com/urfave/cli"
)

const (
	FireblocksAPIKeyNameFlagName       = "fireblocks-api-key-name"
	FireblocksAPISecretNameFlagName    = "fireblocks-api-secret-name"
	FireblocksBaseURLFlagName          = "fireblocks-api-url"
	FireblocksVaultAccountNameFlagName = "fireblocks-vault-account-name"
	FireblocksWalletAddressFlagName    = "fireblocks-wallet-address"
	FireblocksSecretManagerRegion      = "fireblocks-secret-manager-region"
)

type FireblocksConfig struct {
	APIKeyName       string
	SecretKeyName    string
	BaseURL          string
	VaultAccountName string
	WalletAddress    string
	Region           string
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
			EnvVar:   PrefixEnvVar(envPrefix, "FIREBLOCKS_API_SECRET_Name"),
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
	}
}
