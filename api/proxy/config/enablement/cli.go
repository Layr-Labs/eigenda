package enablement

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const (
	EnabledAPIsFlagName = "apis.enabled"
)

func withEnvPrefix(envPrefix, s string) []string {
	return []string{envPrefix + "_" + s}
}

func ReadEnabledServersCfg(ctx *cli.Context) *EnabledServersConfig {
	enabledAPIStrings := ctx.StringSlice(EnabledAPIsFlagName)

	cfg, err := APIStringsToEnabledServersConfig(enabledAPIStrings)
	if err != nil {
		panic(err)
	}

	return cfg
}

func CLIFlags(category string, envPrefix string) []cli.Flag {
	return []cli.Flag{&cli.StringSliceFlag{
		Name: EnabledAPIsFlagName,
		Usage: fmt.Sprintf("Which proxy application APIs to enable. supported options are "+
			"%s", AllAPIsString()),
		Value:    cli.NewStringSlice(),
		Required: false,
		EnvVars:  withEnvPrefix(envPrefix, "APIS_TO_ENABLE"),
		Category: category,
	}}
}
