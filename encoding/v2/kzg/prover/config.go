package prover

import (
	"github.com/Layr-Labs/eigenda/encoding/kzgflags"
	kzgv1 "github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	"github.com/urfave/cli"
)

// KzgConfig holds configuration for the V2 KZG prover.
type KzgConfig struct {
	// Number of G1 points to be loaded from the SRS file at G1Path.
	// This number times 32 bytes will be loaded.
	// Need at least as many points as the maximum blob size in field elements.
	SRSNumberToLoad uint64

	// G1Path is the path to the G1 SRS file.
	G1Path string

	// If true, SRS tables are read from CacheDir during initialization,
	// and parametrizedProvers (fka encoders) are preloaded for all supported encoding params.
	// Generating these on startup would take minutes otherwise.
	PreloadEncoder bool
	// Path to SRS Table directory. Always required even if PreloadEncoder is false,
	// because the prover will write the SRS tables to this directory if they are not already present.
	CacheDir string

	// NumWorker is used in a few places:
	// 1. Num goroutines used to parse the SRS points read from the SRS files.
	// 2. Num goroutines used by the prover for various operations.
	NumWorker uint64
}

// KzgConfigFromV1Config converts a v1 KzgConfig to a v2 prover KzgConfig.
// The V1 KzgConfig is used all over the place in multiple different structs,
// making it very hard to update, optimize, change, or remove unused fields.
// The V2 prover has its own KzgConfig, which is a subset of the V1 KzgConfig.
func KzgConfigFromV1Config(v1 *kzgv1.KzgConfig) *KzgConfig {
	return &KzgConfig{
		SRSNumberToLoad: v1.SRSNumberToLoad,
		G1Path:          v1.G1Path,
		PreloadEncoder:  v1.PreloadEncoder,
		CacheDir:        v1.CacheDir,
		NumWorker:       v1.NumWorker,
	}
}

func ReadCLIConfig(ctx *cli.Context) KzgConfig {
	cfg := KzgConfig{
		SRSNumberToLoad: ctx.GlobalUint64(kzgflags.SRSLoadingNumberFlagName),
		G1Path:          ctx.GlobalString(kzgflags.G1PathFlagName),
		CacheDir:        ctx.GlobalString(kzgflags.CachePathFlagName),
		NumWorker:       ctx.GlobalUint64(kzgflags.NumWorkerFlagName),
		PreloadEncoder:  ctx.GlobalBool(kzgflags.PreloadEncoderFlagName),
	}
	return cfg
}
