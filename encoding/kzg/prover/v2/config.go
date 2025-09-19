package prover

import (
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/urfave/cli"
)

// KzgConfig holds configuration for the V2 KZG prover.
type KzgConfig struct {
	// Number of G1 (and optionally G2) points to be loaded from the SRS files:
	// G1Path, and optionally G2Path and G2TrailingPath.
	// This number times 32 bytes will be loaded from G1Path, and if LoadG2Points is true,
	// this number times 64 bytes will be loaded from G2Path and optionally G2TrailingPath.
	SRSNumberToLoad uint64

	// G1 points are needed by both the prover and verifier, so G1Path is always needed.
	G1Path string

	// G2 SRS points are only needed to compute length proofs.
	// Hence this should be set to true by the apiserver, and false by the encoder.
	LoadG2Points bool
	// G2Path and G2TrailingPath are only needed if LoadG2Points is true.
	// G2 points are used to generate the blob length proof.
	//
	// There are 2 ways to configure G2 points:
	// 1. Entire G2 SRS file (16GiB) is provided via G2Path
	// 2. G2Path and G2TrailingPath both contain at least SRSNumberToLoad points,
	//    where G2Path contains the first part of the G2 SRS file, and G2TrailingPath
	//    contains the trailing end of the G2 SRS file.
	// TODO(samlaf): to prevent misconfigurations and simplify the code, we should probably not multiplex G2Path like this,
	// and instead use a G2PrefixPath config. Then EITHER G2Path is used, OR both G2PrefixPath and G2TrailingPath are used.
	G2Path         string
	G2TrailingPath string

	// PreloadEncoder is only used by the prover to generate kzg multiproofs.
	// It is not needed by the clients/proxy, which only need to generate kzg commitments, not proofs.
	//
	// If true, SRS tables are read from CacheDir during initialization.
	// Generating these on startup would take hours otherwise.
	PreloadEncoder bool
	// Path to SRS Table directory. Always required even if PreloadEncoder is false,
	// because the prover will write the SRS tables to this directory if they are not already present.
	CacheDir string

	// NumWorker is used in a few places:
	// 1. Num goroutines used to parse the SRS points read from the SRS files.
	// 2. Num goroutines used by the prover and verifier.
	// TODO(samlaf): split into separate configs only specified for prover or verifier, where needed.
	NumWorker uint64
	Verbose   bool
}

// KzgConfigFromV1Config converts a v1 KzgConfig to a v2 prover KzgConfig.
// The V1 KzgConfig is used all over the place in multiple different structs,
// making it very hard to update, optimize, change, or remove unused fields.
// The V2 prover has its own KzgConfig, which is a subset of the V1 KzgConfig.
func KzgConfigFromV1Config(v1 *kzg.KzgConfig) *KzgConfig {
	return &KzgConfig{
		SRSNumberToLoad: v1.SRSNumberToLoad,
		G1Path:          v1.G1Path,
		LoadG2Points:    v1.LoadG2Points,
		G2Path:          v1.G2Path,
		G2TrailingPath:  v1.G2TrailingPath,
		PreloadEncoder:  v1.PreloadEncoder,
		CacheDir:        v1.CacheDir,
		NumWorker:       v1.NumWorker,
		Verbose:         v1.Verbose,
	}
}

func ReadCLIConfig(ctx *cli.Context) KzgConfig {
	cfg := KzgConfig{
		SRSNumberToLoad: ctx.GlobalUint64(kzg.SRSLoadingNumberFlagName),
		G1Path:          ctx.GlobalString(kzg.G1PathFlagName),
		// Set to false by default. apiserver needs to set this to true.
		LoadG2Points:   false,
		G2Path:         ctx.GlobalString(kzg.G2PathFlagName),
		G2TrailingPath: ctx.GlobalString(kzg.G2TrailingPathFlagName),
		CacheDir:       ctx.GlobalString(kzg.CachePathFlagName),
		NumWorker:      ctx.GlobalUint64(kzg.NumWorkerFlagName),
		Verbose:        ctx.GlobalBool(kzg.VerboseFlagName),
		PreloadEncoder: ctx.GlobalBool(kzg.PreloadEncoderFlagName),
	}
	return cfg
}
