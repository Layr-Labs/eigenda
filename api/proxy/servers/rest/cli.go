package rest

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	"github.com/urfave/cli/v2"
)

const (
	ListenAddrFlagName = "addr"
	PortFlagName       = "port"

	DeprecatedAPIsEnabledFlagName = "api-enabled"
	DeprecatedAdminAPIType        = "admin"
)

// We don't add any _SERVER_ middlefix to the env vars like we do for other categories
// because these flags were originally in the global namespace, and we don't want to cause
// any breaking changes to the env var names.
func withEnvPrefix(prefix, s string) []string {
	return []string{prefix + "_" + s}
}

func DeprecatedCLIFlags(envPrefix string, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:    DeprecatedAPIsEnabledFlagName,
			Usage:   "List of API types to enable (e.g. admin)",
			Value:   cli.NewStringSlice(),
			EnvVars: withEnvPrefix(envPrefix, "API_ENABLED"),
			Action: func(*cli.Context, []string) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --apis.enabled with `admin` to turn on instead",
					DeprecatedAdminAPIType, withEnvPrefix(envPrefix, "API_ENABLED"))
			},
			Category: category,
		},
	}
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
	}

	return flags
}

func ReadConfig(ctx *cli.Context, apisEnabled *enablement.RestApisEnabled) Config {
	return Config{
		Host:        ctx.String(ListenAddrFlagName),
		Port:        ctx.Int(PortFlagName),
		APIsEnabled: apisEnabled,
		// We can't set compatibility values until after configs have been read as
		// ChainID requires an ethClient connection.
		CompatibilityCfg: common.CompatibilityConfig{},
	}
}
