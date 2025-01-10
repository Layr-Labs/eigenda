package config

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

// Config configures a traffic generator.
type Config struct {
	// Logging configuration.
	LoggingConfig common.LoggerConfig

	// Configuration for the disperser client.
	DisperserClientConfig *clients.DisperserClientConfig

	// Signer private key
	SignerPrivateKey string

	// The port at which the metrics server listens for HTTP requests.
	MetricsHTTPPort string

	// The timeout for the node client.
	NodeClientTimeout time.Duration

	// Path to the runtime configuration file that defines writer groups.
	RuntimeConfigPath string
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, FlagPrefix)
	if err != nil {
		return nil, err
	}

	config := &Config{
		DisperserClientConfig: &clients.DisperserClientConfig{
			Hostname:          ctx.GlobalString(HostnameFlag.Name),
			Port:              ctx.GlobalString(GrpcPortFlag.Name),
			UseSecureGrpcFlag: ctx.GlobalBool(UseSecureGrpcFlag.Name),
		},

		SignerPrivateKey: ctx.String(SignerPrivateKeyFlag.Name),
		LoggingConfig:    *loggerConfig,

		MetricsHTTPPort:   ctx.GlobalString(MetricsHTTPPortFlag.Name),
		NodeClientTimeout: ctx.Duration(NodeClientTimeoutFlag.Name),
		RuntimeConfigPath: ctx.GlobalString(RuntimeConfigPathFlag.Name),
	}

	return config, nil
}
