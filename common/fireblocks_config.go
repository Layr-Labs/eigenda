package common

import (
	"github.com/urfave/cli"
)

const (
	FireblocksAPIKeyFlagName           = "fireblocks-api-key"
	FireblocksAPISecretPathFlagName    = "fireblocks-api-secret-path"
	FireblocksBaseURLFlagName          = "fireblocks-api-url"
	FireblocksVaultAccountNameFlagName = "fireblocks-vault-account-name"
	FireblocksWalletAddressFlagName    = "fireblocks-wallet-address"
)

type FireblocksConfig struct {
	APIKey           string
	SecretKeyPath    string
	BaseURL          string
	VaultAccountName string
	WalletAddress    string
}

func FireblocksCLIFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, FireblocksAPIKeyFlagName),
			Usage:    "Fireblocks API Key. To configure Fireblocks MPC wallet, this field is required. Otherwise, private key must be configured in eth client so that it can fall back to private key wallet.",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "FIREBLOCKS_API_KEY"),
		},
		cli.StringFlag{
			Name:     PrefixFlag(flagPrefix, FireblocksAPISecretPathFlagName),
			Usage:    "Fireblocks API Secret Path. To configure Fireblocks MPC wallet, this field is required. Otherwise, private key must be configured in eth client so that it can fall back to private key wallet.",
			Required: false,
			EnvVar:   PrefixEnvVar(envPrefix, "FIREBLOCKS_API_SECRET_PATH"),
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
	}
}

func ReadFireblocksCLIConfig(ctx *cli.Context, flagPrefix string) FireblocksConfig {
	return FireblocksConfig{
		APIKey:           ctx.GlobalString(PrefixFlag(flagPrefix, FireblocksAPIKeyFlagName)),
		SecretKeyPath:    ctx.GlobalString(PrefixFlag(flagPrefix, FireblocksAPISecretPathFlagName)),
		BaseURL:          ctx.GlobalString(PrefixFlag(flagPrefix, FireblocksBaseURLFlagName)),
		VaultAccountName: ctx.GlobalString(PrefixFlag(flagPrefix, FireblocksVaultAccountNameFlagName)),
		WalletAddress:    ctx.GlobalString(PrefixFlag(flagPrefix, FireblocksWalletAddressFlagName)),
	}
}
