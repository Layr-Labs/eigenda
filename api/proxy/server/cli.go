package server

import (
	"github.com/urfave/cli/v2"
)

const (
	ListenAddrFlagName  = "addr"
	PortFlagName        = "port"
	APIsEnabledFlagName = "api-enabled"
	AdminAPIType        = "admin"
)

// We don't add any _SERVER_ middlefix to the env vars like we do for other categories
// because these flags were originally in the global namespace, and we don't want to cause
// any breaking changes to the env var names.
func withEnvPrefix(prefix, s string) []string {
	return []string{prefix + "_" + s}
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
			Value:    3100,
			EnvVars:  withEnvPrefix(envPrefix, "PORT"),
			Category: category,
		},
		&cli.StringSliceFlag{
			Name:     APIsEnabledFlagName,
			Usage:    "List of API types to enable (e.g. admin)",
			Value:    cli.NewStringSlice(),
			EnvVars:  withEnvPrefix(envPrefix, "API_ENABLED"),
			Category: category,
		},
	}

	return flags
}

func ReadConfig(ctx *cli.Context) Config {
	return Config{
		Host:        ctx.String(ListenAddrFlagName),
		Port:        ctx.Int(PortFlagName),
		EnabledAPIs: ctx.StringSlice(APIsEnabledFlagName),
	}
}
