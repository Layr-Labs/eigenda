package common

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws/secretmanager"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/fireblocks"
	walletsdk "github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gcommon "github.com/ethereum/go-ethereum/common"
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

func NewFireblocksWallet(config *FireblocksConfig, ethClient EthClient, logger logging.Logger) (walletsdk.Wallet, error) {
	if config.Disable {
		logger.Info("Fireblocks wallet disabled")
		return nil, errors.New("fireblocks wallet is disabled")
	}

	validConfigflag := len(config.APIKeyName) > 0 &&
		len(config.SecretKeyName) > 0 &&
		len(config.BaseURL) > 0 &&
		len(config.VaultAccountName) > 0 &&
		len(config.WalletAddress) > 0 &&
		len(config.Region) > 0
	if !validConfigflag {
		return nil, errors.New("fireblocks config is either invalid or incomplete")
	}
	apiKey, err := secretmanager.ReadStringFromSecretManager(context.Background(), config.APIKeyName, config.Region)
	if err != nil {
		return nil, fmt.Errorf("cannot read fireblocks api key %s from secret manager: %w", config.APIKeyName, err)
	}
	secretKey, err := secretmanager.ReadStringFromSecretManager(context.Background(), config.SecretKeyName, config.Region)
	if err != nil {
		return nil, fmt.Errorf("cannot read fireblocks secret key %s from secret manager: %w", config.SecretKeyName, err)
	}
	fireblocksClient, err := fireblocks.NewClient(
		apiKey,
		[]byte(secretKey),
		config.BaseURL,
		config.APITimeout,
		logger.With("component", "FireblocksClient"),
	)
	if err != nil {
		return nil, err
	}
	wallet, err := walletsdk.NewFireblocksWallet(fireblocksClient, ethClient, config.VaultAccountName, logger.With("component", "FireblocksWallet"))
	if err != nil {
		return nil, err
	}
	sender, err := wallet.SenderAddress(context.Background())
	if err != nil {
		return nil, err
	}
	if sender.Cmp(gcommon.HexToAddress(config.WalletAddress)) != 0 {
		return nil, fmt.Errorf("configured wallet address %s does not match derived address %s", config.WalletAddress, sender.Hex())
	}
	logger.Info("Initialized Fireblocks wallet", "vaultAccountName", config.VaultAccountName, "address", sender.Hex())

	return wallet, nil
}
