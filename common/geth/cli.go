package geth

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

var (
	rpcUrlFlagName              = "chain.rpc"
	rpcFallbackUrlFlagName      = "chain.rpc_fallback"
	privateKeyFlagName          = "chain.private-key"
	numConfirmationsFlagName    = "chain.num-confirmations"
	numRetriesFlagName          = "chain.num-retries"
	retryDelayIncrementFlagName = "chain.retry-delay-increment"
)

// TODO(cody.littley): RPCURLs and PrivateKeyString should be converted to *secret.Secret types.

type EthClientConfig struct {
	// A list of RPC URL endpoints to connect to the Ethereum chain.
	RPCURLs []string `docs:"required"`
	// Ethereum private key in hex string format.
	PrivateKeyString string
	// Number of block confirmations to wait for.
	NumConfirmations int
	// Max number of retries for each RPC call after failure.
	NumRetries int
	// Time duration for linear retry delay increment.
	RetryDelay time.Duration
}

func EthClientFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{
			Name:     rpcUrlFlagName,
			Usage:    "Chain rpc. Disperser/Batcher can accept multiple comma separated rpc url. Node only uses the first one",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "CHAIN_RPC"),
		},
		cli.StringFlag{
			Name:     rpcFallbackUrlFlagName,
			Usage:    "Fallback chain rpc for Disperser/Batcher/Dataapi",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "CHAIN_RPC_FALLBACK"),
		},
		cli.StringFlag{
			Name:     privateKeyFlagName,
			Usage:    "Ethereum private key for disperser",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "PRIVATE_KEY"),
		},
		cli.IntFlag{
			Name:     numConfirmationsFlagName,
			Usage:    "Number of confirmations to wait for",
			Required: false,
			Value:    0,
			EnvVar:   common.PrefixEnvVar(envPrefix, "NUM_CONFIRMATIONS"),
		},
		cli.IntFlag{
			Name:     numRetriesFlagName,
			Usage:    "Number of maximal retry for each rpc call after failure",
			Required: false,
			Value:    2,
			EnvVar:   common.PrefixEnvVar(envPrefix, "NUM_RETRIES"),
		},
		cli.DurationFlag{
			Name: retryDelayIncrementFlagName,
			Usage: "Time unit for linear retry delay. For instance, if the retries count is 2 and retry delay is " +
				"1 second, then 0 second is waited for the first call; 1 seconds are waited before the next retry; " +
				"2 seconds are waited for the second retry; if the call failed, the total waited time for retry is " +
				"3 seconds. If the retry delay is 0 second, the total waited time for retry is 0 second.",
			Required: false,
			Value:    0 * time.Second,
			EnvVar:   common.PrefixEnvVar(envPrefix, "RETRY_DELAY_INCREMENT"),
		},
	}
}

func ReadEthClientConfig(ctx *cli.Context) EthClientConfig {
	cfg := EthClientConfig{}
	cfg.RPCURLs = ctx.GlobalStringSlice(rpcUrlFlagName)
	cfg.PrivateKeyString = ctx.GlobalString(privateKeyFlagName)
	cfg.NumConfirmations = ctx.GlobalInt(numConfirmationsFlagName)
	cfg.NumRetries = ctx.GlobalInt(numRetriesFlagName)

	fallbackRPCURL := ctx.GlobalString(rpcFallbackUrlFlagName)
	if len(fallbackRPCURL) > 0 {
		cfg.RPCURLs = append(cfg.RPCURLs, []string{fallbackRPCURL}...)
	}

	return cfg
}

// ReadEthClientConfigRPCOnly doesn't read private key from flag.
// The private key for Node should be read from encrypted key file.
func ReadEthClientConfigRPCOnly(ctx *cli.Context) EthClientConfig {
	cfg := EthClientConfig{}
	cfg.RPCURLs = ctx.GlobalStringSlice(rpcUrlFlagName)
	cfg.NumConfirmations = ctx.GlobalInt(numConfirmationsFlagName)
	cfg.NumRetries = ctx.GlobalInt(numRetriesFlagName)
	cfg.RetryDelay = ctx.GlobalDuration(retryDelayIncrementFlagName)

	fallbackRPCURL := ctx.GlobalString(rpcFallbackUrlFlagName)
	if len(fallbackRPCURL) > 0 {
		cfg.RPCURLs = append(cfg.RPCURLs, []string{fallbackRPCURL}...)
	}

	return cfg
}

// DefaultEthClientConfig returns the default Ethereum client configuration.
func DefaultEthClientConfig() EthClientConfig {
	return EthClientConfig{
		NumConfirmations: 0,
		NumRetries:       2,
		RetryDelay:       0 * time.Second,
	}
}

// Verify validates the Ethereum client configuration.
func (c *EthClientConfig) Verify() error {
	if len(c.RPCURLs) == 0 {
		return fmt.Errorf("at least one RPC URL must be provided")
	}
	for _, url := range c.RPCURLs {
		if url == "" {
			return fmt.Errorf("RPC URL cannot be empty")
		}
	}
	if c.NumConfirmations < 0 {
		return fmt.Errorf("number of confirmations cannot be negative")
	}
	if c.NumRetries < 0 {
		return fmt.Errorf("number of retries cannot be negative")
	}
	if c.RetryDelay < 0 {
		return fmt.Errorf("retry delay cannot be negative")
	}
	return nil
}
