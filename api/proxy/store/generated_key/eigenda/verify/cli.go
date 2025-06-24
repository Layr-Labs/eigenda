package verify

import (
	"fmt"
	"runtime"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/config/eigendaflags"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/urfave/cli/v2"
)

var (
	// cert verification flags
	CertVerificationDisabledFlagName = withFlagPrefix("cert-verification-disabled")

	// kzg flags
	G1PathFlagName                   = withFlagPrefix("g1-path")
	G2PowerOf2PathFlagNameDeprecated = withFlagPrefix("g2-power-of-2-path")
	G2PathFlagName                   = withFlagPrefix("g2-path")
	G2TrailingPathFlagName           = withFlagPrefix("g2-path-trailing")
	CachePathFlagName                = withFlagPrefix("cache-path")
)

// we keep the eigenda prefix like eigenda client flags, because we
// plan to upstream this verification logic into the eigenda client
func withFlagPrefix(s string) string {
	return "eigenda." + s
}

func withEnvPrefix(envPrefix, s string) string {
	return envPrefix + "_EIGENDA_" + s
}

// CLIFlags ... used for Verifier configuration
// category is used to group the flags in the help output (see https://cli.urfave.org/v2/examples/flags/#grouping)
func VerifierCLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:     CertVerificationDisabledFlagName,
			Usage:    "Whether to verify certificates received from EigenDA disperser.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "CERT_VERIFICATION_DISABLED")},
			Value:    false,
			Category: category,
		},
	}
}

// KZGCLIFlags ... used for KZG configuration
// category is used to group the flags in the help output (see https://cli.urfave.org/v2/examples/flags/#grouping)
func KZGCLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		// we use a relative path for these so that the path works for both the binary and the docker container
		// aka we assume the binary is run from root dir, and that the resources/ dir is copied into the working dir of
		// the container
		&cli.StringFlag{
			Name:     G1PathFlagName,
			Usage:    "path to g1.point file.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "TARGET_KZG_G1_PATH")},
			Value:    "resources/g1.point",
			Category: category,
		},
		&cli.StringFlag{
			Name:    G2PowerOf2PathFlagNameDeprecated,
			Usage:   "path to g2.point.powerOf2 file. Deprecated.",
			EnvVars: []string{withEnvPrefix(envPrefix, "TARGET_KZG_G2_POWER_OF_2_PATH")},
			Action: func(_ *cli.Context, _ string) error {
				return fmt.Errorf(
					"flag --%s (env var %s) is deprecated",
					G2PowerOf2PathFlagNameDeprecated, withEnvPrefix(envPrefix, "TARGET_KZG_G2_POWER_OF_2_PATH"))
			},
			Category: category,
			Hidden:   true,
		},
		&cli.StringFlag{
			Name:     G2PathFlagName,
			Usage:    "path to g2.point file.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "TARGET_KZG_G2_PATH")},
			Value:    "resources/g2.point",
			Category: category,
		},
		&cli.StringFlag{
			Name:     G2TrailingPathFlagName,
			Usage:    "path to g2.trailing.point file.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "TARGET_KZG_G2_TRAILING_PATH")},
			Value:    "resources/g2.trailing.point",
			Category: category,
		},
		&cli.StringFlag{
			Name:     CachePathFlagName,
			Usage:    "path to SRS tables for caching. This resource is not currently used, but needed because of the shared eigenda KZG library that we use. We will eventually fix this.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "TARGET_CACHE_PATH")},
			Value:    "resources/SRSTables/",
			Category: category,
		},
	}
}

func ReadKzgConfig(ctx *cli.Context, maxBlobSizeBytes uint64) kzg.KzgConfig {
	return kzg.KzgConfig{
		G1Path:          ctx.String(G1PathFlagName),
		G2Path:          ctx.String(G2PathFlagName),
		G2TrailingPath:  ctx.String(G2TrailingPathFlagName),
		CacheDir:        ctx.String(CachePathFlagName),
		SRSOrder:        eigendaflags.SrsOrder,
		SRSNumberToLoad: maxBlobSizeBytes / 32,         // # of fr.Elements
		NumWorker:       uint64(runtime.GOMAXPROCS(0)), // #nosec G115
		// we are intentionally not setting the `LoadG2Points` field here. `LoadG2Points` has different requirements
		// for v1 vs v2. To make things foolproof, we just set this value locally prior to use, so that it can't
		// ever be set incorrectly.
	}
}

// ReadConfig takes an eigendaClientConfig as input because the verifier config reuses some configs that are already
// defined in the client config
func ReadConfig(ctx *cli.Context, clientConfigV1 common.ClientConfigV1) Config {
	return Config{
		VerifyCerts: !ctx.Bool(CertVerificationDisabledFlagName),
		// reuse some configs from the eigenda client
		RPCURL:               clientConfigV1.EdaClientCfg.EthRpcUrl,
		SvcManagerAddr:       clientConfigV1.EdaClientCfg.SvcManagerAddr,
		EthConfirmationDepth: clientConfigV1.EdaClientCfg.WaitForConfirmationDepth,
		WaitForFinalization:  clientConfigV1.EdaClientCfg.WaitForFinalization,
	}
}
