package kzg

import (
	"runtime"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	G1PathFlagName            = "kzg.g1-path"
	G2PathFlagName            = "kzg.g2-path"
	G2TrailingPathFlagName    = "kzg.g2-trailing-path"
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
			Usage:    "Path to G2 SRS. Either this flag or G2_POWER_OF_2_PATH needs to be specified. For operator node, if both are specified, the node uses G2_POWER_OF_2_PATH first, if failed then tries to G2_PATH",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "G2_PATH"),
		},
		cli.StringFlag{
			Name:     G2TrailingPathFlagName,
			Usage:    "Path to trailing G2 SRS file. Its intended purpose is to allow local generation the blob length proof. If you already downloaded the entire G2 SRS file which contains 268435456 G2 points with total size 16GiB, this flag is not needed. With this G2TrailingPathFlag, user can use a smaller file that contains only the trailing end of the whole G2 SRS file. Ignoring this flag, the program assumes the entire G2 SRS file is provided. With this flag, the size of the provided file must be at least SRSLoadingNumberFlagName * 64 Bytes.",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "G2_TRAILING_PATH"),
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
			Usage:    "Path to G2 SRS points that are on power of 2. Either this flag or G2_PATH needs to be specified. For operator node, if both are specified, the node uses G2_POWER_OF_2_PATH first, if failed then tries to G2_PATH",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "G2_POWER_OF_2_PATH"),
		},
	}
}

func ReadCLIConfig(ctx *cli.Context) KzgConfig {
	cfg := KzgConfig{}
	cfg.G1Path = ctx.GlobalString(G1PathFlagName)
	cfg.G2Path = ctx.GlobalString(G2PathFlagName)
	cfg.G2TrailingPath = ctx.GlobalString(G2TrailingPathFlagName)
	cfg.CacheDir = ctx.GlobalString(CachePathFlagName)
	cfg.SRSOrder = ctx.GlobalUint64(SRSOrderFlagName)
	cfg.SRSNumberToLoad = ctx.GlobalUint64(SRSLoadingNumberFlagName)
	cfg.NumWorker = ctx.GlobalUint64(NumWorkerFlagName)
	cfg.Verbose = ctx.GlobalBool(VerboseFlagName)
	cfg.PreloadEncoder = ctx.GlobalBool(PreloadEncoderFlagName)
	cfg.G2PowerOf2Path = ctx.GlobalString(G2PowerOf2PathFlagName)

	return cfg
}
