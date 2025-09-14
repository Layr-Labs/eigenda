package enabled_apis

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

func ReadEnabledAPIs(ctx *cli.Context) *EnabledAPIs {
	enabledAPIStrings := ctx.StringSlice(EnabledAPIsFlagName)

	apis, err := StringsToEnabledAPIs(enabledAPIStrings)
	if err != nil {
		panic(err)
	}

	return apis
}

func CLIFlags(category string, envPrefix string) []cli.Flag {
	return []cli.Flag{&cli.StringSliceFlag{
		Name: EnabledAPIsFlagName,
		Usage: fmt.Sprintf("Which proxy application APIs to enable. supported options are "+
			"%s, %s, %s, %s, %s, %s", Admin.ToString(), StandardCommitment.ToString(),
			OpGenericCommitment.ToString(), OpKeccakCommitment.ToString(),
			ArbCustomDAServer.ToString(), MetricsServer.ToString()),
		Value:    cli.NewStringSlice(),
		EnvVars:  withEnvPrefix(envPrefix, "APIS_TO_ENABLE"),
		Category: category,
	}}
}
