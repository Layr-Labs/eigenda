package verify

import (
	"fmt"
	"runtime"

	"github.com/Layr-Labs/eigenda-proxy/utils"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/urfave/cli/v2"
)

var (
	BytesPerSymbol     = 31
	MaxCodingRatio     = 8
	MaxSRSPoints       = 1 << 28 // 2^28
	MaxAllowedBlobSize = uint64(MaxSRSPoints * BytesPerSymbol / MaxCodingRatio)
)

// TODO: should this live in the resources pkg?
// So that if we ever change the SRS files there we can change this value
const srsOrder = 268435456 // 2 ^ 32

var (
	// cert verification flags
	// TODO: we keep the eigenda prefix like eigenda client flags, because we
	// plan to upstream this verification logic into the eigenda client
	CertVerificationEnabledFlagName = withFlagPrefix("cert-verification-enabled")
	EthRPCFlagName                  = withFlagPrefix("eth-rpc")
	SvcManagerAddrFlagName          = withFlagPrefix("svc-manager-addr")
	EthConfirmationDepthFlagName    = withFlagPrefix("eth-confirmation-depth")

	// kzg flags
	G1PathFlagName        = withFlagPrefix("g1-path")
	G2TauFlagName         = withFlagPrefix("g2-tau-path")
	CachePathFlagName     = withFlagPrefix("cache-path")
	MaxBlobLengthFlagName = withFlagPrefix("max-blob-length")
)

func withFlagPrefix(s string) string {
	return "eigenda." + s
}

func withEnvPrefix(envPrefix, s string) []string {
	return []string{envPrefix + "_EIGENDA_" + s}
}

// CLIFlags ... used for Verifier configuration
// category is used to group the flags in the help output (see https://cli.urfave.org/v2/examples/flags/#grouping)
func CLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    CertVerificationEnabledFlagName,
			Usage:   "Whether to verify certificates received from EigenDA disperser.",
			EnvVars: withEnvPrefix(envPrefix, "CERT_VERIFICATION_ENABLED"),
			// TODO: ideally we'd want this to be turned on by default when eigenda backend is used (memstore.enabled=false)
			Value:    false,
			Category: category,
		},
		&cli.StringFlag{
			Name:     EthRPCFlagName,
			Usage:    "JSON RPC node endpoint for the Ethereum network used for finalizing DA blobs. See available list here: https://docs.eigenlayer.xyz/eigenda/networks/",
			EnvVars:  withEnvPrefix(envPrefix, "ETH_RPC"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     SvcManagerAddrFlagName,
			Usage:    "The deployed EigenDA service manager address. The list can be found here: https://github.com/Layr-Labs/eigenlayer-middleware/?tab=readme-ov-file#current-mainnet-deployment",
			EnvVars:  withEnvPrefix(envPrefix, "SERVICE_MANAGER_ADDR"),
			Category: category,
		},
		&cli.Uint64Flag{
			Name:     EthConfirmationDepthFlagName,
			Usage:    "The number of Ethereum blocks to wait before considering a submitted blob's DA batch submission confirmed. `0` means wait for inclusion only.",
			EnvVars:  withEnvPrefix(envPrefix, "ETH_CONFIRMATION_DEPTH"),
			Value:    0,
			Category: category,
		},
		// kzg flags
		&cli.StringFlag{
			Name:    G1PathFlagName,
			Usage:   "Directory path to g1.point file.",
			EnvVars: withEnvPrefix(envPrefix, "TARGET_KZG_G1_PATH"),
			// TODO: should use absolute path wrt root directory to prevent future errors
			//       in case we move this file around
			Value:    "../resources/g1.point",
			Category: category,
		},
		&cli.StringFlag{
			Name:     G2TauFlagName,
			Usage:    "Directory path to g2.point.powerOf2 file.",
			EnvVars:  withEnvPrefix(envPrefix, "TARGET_G2_TAU_PATH"),
			Value:    "../resources/g2.point.powerOf2",
			Category: category,
		},
		&cli.StringFlag{
			Name:     CachePathFlagName,
			Usage:    "Directory path to SRS tables for caching.",
			EnvVars:  withEnvPrefix(envPrefix, "TARGET_CACHE_PATH"),
			Value:    "../resources/SRSTables/",
			Category: category,
		},
		// TODO: can we use a genericFlag for this, and automatically parse the string into a uint64?
		&cli.StringFlag{
			Name:    MaxBlobLengthFlagName,
			Usage:   "Maximum blob length to be written or read from EigenDA. Determines the number of SRS points loaded into memory for KZG commitments. Example units: '30MiB', '4Kb', '30MB'. Maximum size slightly exceeds 1GB.",
			EnvVars: withEnvPrefix(envPrefix, "MAX_BLOB_LENGTH"),
			Value:   "16MiB",
			Action: func(_ *cli.Context, maxBlobLengthStr string) error {
				// parse the string to a uint64 and set the maxBlobLengthBytes var to be used by ReadConfig()
				numBytes, err := utils.ParseBytesAmount(maxBlobLengthStr)
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

// this var is set by the action in the MaxBlobLengthFlagName flag
// TODO: there's def a better way to deal with this... perhaps a generic flag that can parse the string into a uint64?
var MaxBlobLengthBytes uint64

func ReadConfig(ctx *cli.Context) Config {
	kzgCfg := &kzg.KzgConfig{
		G1Path:          ctx.String(G1PathFlagName),
		G2PowerOf2Path:  ctx.String(G2TauFlagName),
		CacheDir:        ctx.String(CachePathFlagName),
		SRSOrder:        srsOrder,
		SRSNumberToLoad: MaxBlobLengthBytes / 32,       // # of fr.Elements
		NumWorker:       uint64(runtime.GOMAXPROCS(0)), // #nosec G115
	}

	return Config{
		KzgConfig:            kzgCfg,
		VerifyCerts:          ctx.Bool(CertVerificationEnabledFlagName),
		RPCURL:               ctx.String(EthRPCFlagName),
		SvcManagerAddr:       ctx.String(SvcManagerAddrFlagName),
		EthConfirmationDepth: uint64(ctx.Int64(EthConfirmationDepthFlagName)), // #nosec G115
	}
}
