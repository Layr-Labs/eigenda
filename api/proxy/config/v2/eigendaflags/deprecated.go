package eigendaflags

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	deprecatedServiceManagerAddrFlagName        = withFlagPrefix("service-manager-addr")
	deprecatedBLSOperatorStateRetrieverFlagName = withFlagPrefix("bls-operator-state-retriever-addr")
)

func DeprecatedCLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     deprecatedServiceManagerAddrFlagName,
			Usage:    "[Deprecated: use EigenDADirectory instead] Address of the EigenDA Service Manager contract.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "SERVICE_MANAGER_ADDR")},
			Category: category,
			Required: false,
			Hidden:   true,
			Action: func(c *cli.Context, s string) error {
				return fmt.Errorf("--%s is deprecated. Contract addresses shall now be read from the "+
					"EigenDA Directory contract (via the --%s flag) instead. "+
					"See https://docs.eigencloud.xyz/products/eigenda/networks/mainnet#contract-addresses for more details",
					s, EigenDADirectoryFlagName)
			},
		},
		&cli.StringFlag{
			Name:     deprecatedBLSOperatorStateRetrieverFlagName,
			Usage:    "[Deprecated: use EigenDADirectory instead] Address of the BLS operator state retriever contract.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "BLS_OPERATOR_STATE_RETRIEVER_ADDR")},
			Category: category,
			Required: false,
			Hidden:   true,
			Action: func(c *cli.Context, s string) error {
				return fmt.Errorf("--%s is deprecated. Contract addresses shall now be read from the "+
					"EigenDA Directory contract (via the --%s flag) instead. "+
					"See https://docs.eigencloud.xyz/products/eigenda/networks/mainnet#contract-addresses for more details",
					s, EigenDADirectoryFlagName)
			},
		},
	}
}
