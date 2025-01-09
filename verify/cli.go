package verify

import (
	"fmt"
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
)

const (
	BytesPerSymbol     = 31
	MaxCodingRatio     = 8
	SrsOrder           = 1 << 28 // 2^28
	MaxAllowedBlobSize = uint64(SrsOrder * BytesPerSymbol / MaxCodingRatio)
)

var (
	// cert verification flags
	CertVerificationDisabledFlagName = withFlagPrefix("cert-verification-disabled")

	// kzg flags
	G1PathFlagName         = withFlagPrefix("g1-path")
	G2PowerOf2PathFlagName = withFlagPrefix("g2-power-of-2-path")
	CachePathFlagName      = withFlagPrefix("cache-path")
	MaxBlobLengthFlagName  = withFlagPrefix("max-blob-length")
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
func CLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:     CertVerificationDisabledFlagName,
			Usage:    "Whether to verify certificates received from EigenDA disperser.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "CERT_VERIFICATION_DISABLED")},
			Value:    false,
			Category: category,
		},
		// kzg flags
		&cli.StringFlag{
			Name:    G1PathFlagName,
			Usage:   "Directory path to g1.point file.",
			EnvVars: []string{withEnvPrefix(envPrefix, "TARGET_KZG_G1_PATH")},
			// we use a relative path so that the path works for both the binary and the docker container
			// aka we assume the binary is run from root dir, and that the resources/ dir is copied into the working dir of the container
			Value:    "resources/g1.point",
			Category: category,
		},
		&cli.StringFlag{
			Name:    G2PowerOf2PathFlagName,
			Usage:   "Directory path to g2.point.powerOf2 file. This resource is not currently used, but needed because of the shared eigenda KZG library that we use. We will eventually fix this.",
			EnvVars: []string{withEnvPrefix(envPrefix, "TARGET_KZG_G2_POWER_OF_2_PATH")},
			// we use a relative path so that the path works for both the binary and the docker container
			// aka we assume the binary is run from root dir, and that the resources/ dir is copied into the working dir of the container
			Value:    "resources/g2.point.powerOf2",
			Category: category,
		},
		&cli.StringFlag{
			Name:    CachePathFlagName,
			Usage:   "Directory path to SRS tables for caching. This resource is not currently used, but needed because of the shared eigenda KZG library that we use. We will eventually fix this.",
			EnvVars: []string{withEnvPrefix(envPrefix, "TARGET_CACHE_PATH")},
			// we use a relative path so that the path works for both the binary and the docker container
			// aka we assume the binary is run from root dir, and that the resources/ dir is copied into the working dir of the container
			Value:    "resources/SRSTables/",
			Category: category,
		},
		// TODO: can we use a genericFlag for this, and automatically parse the string into a uint64?
		&cli.StringFlag{
			Name:    MaxBlobLengthFlagName,
			Usage:   "Maximum blob length to be written or read from EigenDA. Determines the number of SRS points loaded into memory for KZG commitments. Example units: '30MiB', '4Kb', '30MB'. Maximum size slightly exceeds 1GB.",
			EnvVars: []string{withEnvPrefix(envPrefix, "MAX_BLOB_LENGTH")},
			Value:   "16MiB",
			// set to true to force action to run on the default Value
			// see https://github.com/urfave/cli/issues/1973
			HasBeenSet: true,
			Action: func(_ *cli.Context, maxBlobLengthStr string) error {
				// parse the string to a uint64 and set the maxBlobLengthBytes var to be used by ReadConfig()
				numBytes, err := common.ParseBytesAmount(maxBlobLengthStr)
				if err != nil {
					return fmt.Errorf("failed to parse max blob length flag: %w", err)
				}
				if numBytes == 0 {
					return fmt.Errorf("max blob length is 0")
				}
				if numBytes > MaxAllowedBlobSize {
					return fmt.Errorf("excluding disperser constraints on max blob size, SRS points constrain the maxBlobLength configuration parameter to be less than than %d bytes", MaxAllowedBlobSize)
				}
				MaxBlobLengthBytes = numBytes
				return nil
			},
			// we also use this flag for memstore.
			// should we duplicate the flag? Or is there a better way to handle this?
			Category: category,
		},
	}
}

// MaxBlobLengthBytes ... there's def a better way to deal with this... perhaps a generic flag that can parse the string into a uint64?
// this var is set by the action in the MaxBlobLengthFlagName flag
var MaxBlobLengthBytes uint64

// ReadConfig takes an eigendaClientConfig as input because the verifier config
// reuses some configs that are also used by the eigenda client.
// Not sure if urfave has a way to do flag aliases so opted for this approach.
func ReadConfig(ctx *cli.Context, edaClientConfig clients.EigenDAClientConfig) Config {
	kzgCfg := &kzg.KzgConfig{
		G1Path:          ctx.String(G1PathFlagName),
		G2PowerOf2Path:  ctx.String(G2PowerOf2PathFlagName),
		CacheDir:        ctx.String(CachePathFlagName),
		SRSOrder:        SrsOrder,
		SRSNumberToLoad: MaxBlobLengthBytes / 32,       // # of fr.Elements
		NumWorker:       uint64(runtime.GOMAXPROCS(0)), // #nosec G115
	}

	return Config{
		KzgConfig:   kzgCfg,
		VerifyCerts: !ctx.Bool(CertVerificationDisabledFlagName),
		// reuse some configs from the eigenda client
		RPCURL:               edaClientConfig.EthRpcUrl,
		SvcManagerAddr:       edaClientConfig.SvcManagerAddr,
		EthConfirmationDepth: edaClientConfig.WaitForConfirmationDepth,
		WaitForFinalization:  edaClientConfig.WaitForFinalization,
	}
}
