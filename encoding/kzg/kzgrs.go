package kzg

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/kzgflags"
	_ "github.com/Layr-Labs/eigenda/resources/srs"
	"github.com/urfave/cli"
)

// KzgConfig holds configuration for KZG prover and verifier.
// Some of the configurations only apply to the prover or verifier.
// TODO(samlaf): split into separate Prover and Verifier configs.
type KzgConfig struct {
	// SRSOrder is the total size of SRS.
	// TODO(samlaf): this should always be 2^28. Get rid of this field and replace with hardcoded constant.
	SRSOrder uint64
	// Number of G1 (and optionally G2) points to be loaded from the SRS files:
	// G1Path, and optionally G2Path and G2TrailingPath.
	// This number times 32 bytes will be loaded from G1Path, and if LoadG2Points is true,
	// this number times 64 bytes will be loaded from G2Path and optionally G2TrailingPath.
	SRSNumberToLoad uint64

	// G1 points are needed by both the prover and verifier, so G1Path is always needed.
	G1Path string

	// G2 SRS points are only needed by the prover, since the verifier uses hardcoded G2 powers of 2.
	// See [srs.G2PowerOf2SRS] for details.
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

func ReadCLIConfig(ctx *cli.Context) KzgConfig {
	cfg := KzgConfig{}
	cfg.G1Path = ctx.GlobalString(kzgflags.G1PathFlagName)
	cfg.G2Path = ctx.GlobalString(kzgflags.G2PathFlagName)
	cfg.G2TrailingPath = ctx.GlobalString(kzgflags.G2TrailingPathFlagName)
	cfg.CacheDir = ctx.GlobalString(kzgflags.CachePathFlagName)
	cfg.SRSOrder = ctx.GlobalUint64(kzgflags.SRSOrderFlagName)
	cfg.SRSNumberToLoad = ctx.GlobalUint64(kzgflags.SRSLoadingNumberFlagName)
	cfg.NumWorker = ctx.GlobalUint64(kzgflags.NumWorkerFlagName)
	cfg.Verbose = ctx.GlobalBool(kzgflags.VerboseFlagName)
	cfg.PreloadEncoder = ctx.GlobalBool(kzgflags.PreloadEncoderFlagName)

	if ctx.GlobalString(kzgflags.DeprecatedG2PowerOf2PathFlagName) != "" {
		fmt.Printf("Warning: --%s is deprecated. "+
			"The g2.point.powerOf2 file is now embedded in the binary, so this flag is no longer needed.\n",
			kzgflags.DeprecatedG2PowerOf2PathFlagName)
	}

	return cfg
}
