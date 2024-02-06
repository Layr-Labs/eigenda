package encoding

import (
	"runtime"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/pkg/encoding/kzgEncoder"
	"github.com/urfave/cli"
)

const (
	G1PathFlagName            = "kzg.g1-path"
	G2PathFlagName            = "kzg.g2-path"
	CachePathFlagName         = "kzg.cache-path"
	SRSOrderFlagName          = "kzg.srs-order"
	NumWorkerFlagName         = "kzg.num-workers"
	VerboseFlagName           = "kzg.verbose"
	PreloadEncoderFlagName    = "kzg.preload-encoder"
	CacheEncodedBlobsFlagName = "cache-encoded-blobs"
	SRSLoadingNumberFlagName  = "kzg.srs-load"
	G2PowerOf2PathFlagName    = "kzg.g2-power-of-2-path"
)

func CLIFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     G1PathFlagName,
			Usage:    "Path to G1 SRS",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "G1_PATH"),
		},
		cli.StringFlag{
			Name:     G2PathFlagName,
			Usage:    "Path to G2 SRS",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "G2_PATH"),
		},
		cli.StringFlag{
			Name:     CachePathFlagName,
			Usage:    "Path to SRS Table directory",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "CACHE_PATH"),
		},
		cli.Uint64Flag{
			Name:     SRSOrderFlagName,
			Usage:    "Order of the SRS",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "SRS_ORDER"),
		},
		cli.Uint64Flag{
			Name:     SRSLoadingNumberFlagName,
			Usage:    "Number of SRS points to load into memory",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "SRS_LOAD"),
		},
		cli.Uint64Flag{
			Name:     NumWorkerFlagName,
			Usage:    "Number of workers for multithreading",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "NUM_WORKERS"),
			Value:    uint64(runtime.GOMAXPROCS(0)),
		},
		cli.BoolFlag{
			Name:     VerboseFlagName,
			Usage:    "Enable to see verbose output for encoding/decoding",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "VERBOSE"),
		},
		cli.BoolFlag{
			Name:     CacheEncodedBlobsFlagName,
			Usage:    "Enable to cache encoded results",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "CACHE_ENCODED_BLOBS"),
		},
		cli.BoolFlag{
			Name:     PreloadEncoderFlagName,
			Usage:    "Set to enable Encoder PreLoading",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "PRELOAD_ENCODER"),
		},
		cli.StringFlag{
			Name:     G2PowerOf2PathFlagName,
			Usage:    "Path to G2 SRS points that are on power of 2",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "G2_POWER_OF_2_PATH"),
		},
	}
}

func ReadCLIConfig(ctx *cli.Context) EncoderConfig {
	cfg := kzgEncoder.KzgConfig{}
	cfg.G1Path = ctx.GlobalString(G1PathFlagName)
	cfg.G2Path = ctx.GlobalString(G2PathFlagName)
	cfg.CacheDir = ctx.GlobalString(CachePathFlagName)
	cfg.SRSOrder = ctx.GlobalUint64(SRSOrderFlagName)
	cfg.SRSNumberToLoad = ctx.GlobalUint64(SRSLoadingNumberFlagName)
	cfg.NumWorker = ctx.GlobalUint64(NumWorkerFlagName)
	cfg.Verbose = ctx.GlobalBool(VerboseFlagName)
	cfg.PreloadEncoder = ctx.GlobalBool(PreloadEncoderFlagName)
	cfg.G2PowerOf2Path = ctx.GlobalString(G2PowerOf2PathFlagName)

	return EncoderConfig{
		KzgConfig:         cfg,
		CacheEncodedBlobs: ctx.GlobalBoolT(CacheEncodedBlobsFlagName),
	}
}
