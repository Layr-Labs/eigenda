package verify

import (
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/flags/eigendaflags"
	"github.com/urfave/cli/v2"
)

// TODO: we should also upstream these deprecated flags into the eigenda client
// if we upstream the changes before removing the deprecated flags
var (
	// cert verification flags
	DeprecatedCertVerificationEnabledFlagName = withDeprecatedFlagPrefix("cert-verification-enabled")
	DeprecatedEthRPCFlagName                  = withDeprecatedFlagPrefix("eth-rpc")
	DeprecatedSvcManagerAddrFlagName          = withDeprecatedFlagPrefix("svc-manager-addr")
	DeprecatedEthConfirmationDepthFlagName    = withDeprecatedFlagPrefix("eth-confirmation-depth")

	// kzg flags
	DeprecatedG1PathFlagName        = withDeprecatedFlagPrefix("g1-path")
	DeprecatedG2TauFlagName         = withDeprecatedFlagPrefix("g2-tau-path")
	DeprecatedCachePathFlagName     = withDeprecatedFlagPrefix("cache-path")
	DeprecatedMaxBlobLengthFlagName = withDeprecatedFlagPrefix("max-blob-length")
)

func withDeprecatedFlagPrefix(s string) string {
	return "eigenda-" + s
}

func withDeprecatedEnvPrefix(envPrefix, s string) string {
	return envPrefix + "_" + s
}

// CLIFlags ... used for Verifier configuration
// category is used to group the flags in the help output (see https://cli.urfave.org/v2/examples/flags/#grouping)
func DeprecatedCLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    DeprecatedEthRPCFlagName,
			Usage:   "JSON RPC node endpoint for the Ethereum network used for finalizing DA blobs. See available list here: https://docs.eigenlayer.xyz/eigenda/networks/",
			EnvVars: []string{withDeprecatedEnvPrefix(envPrefix, "ETH_RPC")},
			Action: func(_ *cli.Context, _ string) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedEthRPCFlagName, withDeprecatedEnvPrefix(envPrefix, "ETH_RPC"),
					eigendaflags.EthRPCURLFlagName, withEnvPrefix(envPrefix, "ETH_RPC"))
			},
			Category: category,
		},
		&cli.StringFlag{
			Name:    DeprecatedSvcManagerAddrFlagName,
			Usage:   "The deployed EigenDA service manager address. The list can be found here: https://github.com/Layr-Labs/eigenlayer-middleware/?tab=readme-ov-file#current-mainnet-deployment",
			EnvVars: []string{withDeprecatedEnvPrefix(envPrefix, "SERVICE_MANAGER_ADDR")},
			Action: func(_ *cli.Context, _ string) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedSvcManagerAddrFlagName, withDeprecatedEnvPrefix(envPrefix, "SERVICE_MANAGER_ADDR"),
					eigendaflags.SvcManagerAddrFlagName, withEnvPrefix(envPrefix, "SERVICE_MANAGER_ADDR"))
			},
			Category: category,
		},
		&cli.Uint64Flag{
			Name:    DeprecatedEthConfirmationDepthFlagName,
			Usage:   "The number of Ethereum blocks to wait before considering a submitted blob's DA batch submission confirmed. `0` means wait for inclusion only.",
			EnvVars: []string{withDeprecatedEnvPrefix(envPrefix, "ETH_CONFIRMATION_DEPTH")},
			Value:   0,
			Action: func(_ *cli.Context, _ uint64) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedEthConfirmationDepthFlagName, withDeprecatedEnvPrefix(envPrefix, "ETH_CONFIRMATION_DEPTH"),
					eigendaflags.ConfirmationDepthFlagName, withEnvPrefix(envPrefix, "CONFIRMATION_DEPTH"))
			},
			Category: category,
		},
		// kzg flags
		&cli.StringFlag{
			Name:    DeprecatedG1PathFlagName,
			Usage:   "Directory path to g1.point file.",
			EnvVars: []string{withDeprecatedEnvPrefix(envPrefix, "TARGET_KZG_G1_PATH")},
			// we use a relative path so that the path works for both the binary and the docker container
			// aka we assume the binary is run from root dir, and that the resources/ dir is copied into the working dir of the container
			Value: "resources/g1.point",
			Action: func(_ *cli.Context, _ string) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedG1PathFlagName, withDeprecatedEnvPrefix(envPrefix, "TARGET_KZG_G1_PATH"),
					G1PathFlagName, withEnvPrefix(envPrefix, "TARGET_KZG_G1_PATH"))
			},
			Category: category,
		},
		&cli.StringFlag{
			Name:    DeprecatedG2TauFlagName,
			Usage:   "Directory path to g2.point.powerOf2 file.",
			EnvVars: []string{withDeprecatedEnvPrefix(envPrefix, "TARGET_G2_TAU_PATH")},
			// we use a relative path so that the path works for both the binary and the docker container
			// aka we assume the binary is run from root dir, and that the resources/ dir is copied into the working dir of the container
			Value: "resources/g2.point.powerOf2",
			Action: func(_ *cli.Context, _ string) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedG2TauFlagName, withDeprecatedEnvPrefix(envPrefix, "TARGET_G2_TAU_PATH"),
					G2PowerOf2PathFlagName, withEnvPrefix(envPrefix, "TARGET_KZG_G2_POWER_OF_2_PATH"))
			},
			Category: category,
		},
		&cli.StringFlag{
			Name:    DeprecatedCachePathFlagName,
			Usage:   "Directory path to SRS tables for caching.",
			EnvVars: []string{withDeprecatedEnvPrefix(envPrefix, "TARGET_CACHE_PATH")},
			// we use a relative path so that the path works for both the binary and the docker container
			// aka we assume the binary is run from root dir, and that the resources/ dir is copied into the working dir of the container
			Value: "resources/SRSTables/",
			Action: func(_ *cli.Context, _ string) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedCachePathFlagName, withDeprecatedEnvPrefix(envPrefix, "TARGET_CACHE_PATH"),
					CachePathFlagName, withEnvPrefix(envPrefix, "TARGET_CACHE_PATH"))
			},
			Category: category,
		},
		// TODO: can we use a genericFlag for this, and automatically parse the string into a uint64?
		&cli.StringFlag{
			Name:    DeprecatedMaxBlobLengthFlagName,
			Usage:   "Maximum blob length to be written or read from EigenDA. Determines the number of SRS points loaded into memory for KZG commitments. Example units: '30MiB', '4Kb', '30MB'. Maximum size slightly exceeds 1GB.",
			EnvVars: []string{withDeprecatedEnvPrefix(envPrefix, "MAX_BLOB_LENGTH")},
			Value:   "16MiB",
			Action: func(_ *cli.Context, _ string) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedMaxBlobLengthFlagName, withDeprecatedEnvPrefix(envPrefix, "MAX_BLOB_LENGTH"),
					MaxBlobLengthFlagName, withEnvPrefix(envPrefix, "MAX_BLOB_LENGTH"))
			},
			// we also use this flag for memstore.
			// should we duplicate the flag? Or is there a better way to handle this?
			Category: category,
		},
	}
}
