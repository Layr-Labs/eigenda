package arbitrum_altda

import (
	"github.com/urfave/cli/v2"
)

const (
	ListenAddrFlagName = "arbitrum-da.addr"
	PortFlagName       = "arbitrum-da.port"
	Enabled            = "arbitrum-da.enabled"
)

func withEnvPrefix(prefix, s string) []string {
	return []string{prefix + "_ARB_DA_" + s}
}

func CLIFlags(envPrefix string, category string) []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     ListenAddrFlagName,
			Usage:    "Server listening address",
			Value:    "0.0.0.0",
			EnvVars:  withEnvPrefix(envPrefix, "ADDR"),
			Category: category,
		},
		&cli.IntFlag{
			Name:     PortFlagName,
			Usage:    "Server listening port",
			Value:    3101,
			EnvVars:  withEnvPrefix(envPrefix, "PORT"),
			Category: category,
		},
		&cli.BoolFlag{
			Name:     Enabled,
			Usage:    "Whether or not to enable Arbitrum Custom DA JSON RPC API",
			Value:    false,
			EnvVars:  withEnvPrefix(envPrefix, "ENABLED"),
			Category: category,
		},
	}

	return flags
}

func ReadConfig(ctx *cli.Context) Config {
	return Config{
		Host:   ctx.String(ListenAddrFlagName),
		Port:   ctx.Int(PortFlagName),
		Enable: ctx.Bool(Enabled),
	}
}
