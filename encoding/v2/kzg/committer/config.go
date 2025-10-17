package committer

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/encoding/kzgflags"
	"github.com/urfave/cli"
)

type Config struct {
	// Number of SRS points to load from SRS files. Must be a power of 2.
	// Committer will only be able to compute commitments for blobs of size up to this number of field elements.
	// e.g. if SRSNumberToLoad=2^19, then the committer can compute commitments for blobs of size up to
	// 2^19 field elements = 2^19 * 32 bytes = 16 MiB.
	SRSNumberToLoad uint64
	G1SRSPath       string
	// There are 2 ways to configure G2 points:
	// 1. Entire G2 SRS file (16GiB) is provided via G2SRSPath (G2TrailingSRSPath is not used).
	// 2. G2SRSPath and G2TrailingSRSPath both contain at least SRSNumberToLoad points,
	//    where G2SRSPath contains the first SRSNumberToLoad points of the full G2 SRS file,
	//    and G2TrailingSRSPath contains the last SRSNumberToLoad points of the G2 SRS file.
	//
	// TODO(samlaf): to prevent misconfigurations and simplify the code, we should probably
	// not multiplex G2SRSPath like this, and instead use a G2PrefixPath config.
	// Then EITHER G2SRSPath is used, OR both G2PrefixSRSPath and G2TrailingSRSPath are used.
	G2SRSPath         string
	G2TrailingSRSPath string
}

var _ config.VerifiableConfig = (*Config)(nil)

func (c *Config) Verify() error {
	if c.SRSNumberToLoad <= 0 {
		return fmt.Errorf("SRSNumberToLoad must be specified for disperser version 2")
	}
	if c.G1SRSPath == "" {
		return fmt.Errorf("G1SRSPath must be specified for disperser version 2")
	}
	if c.G2SRSPath == "" {
		return fmt.Errorf("G2SRSPath must be specified for disperser version 2")
	}
	// G2TrailingSRSPath is optional but its need depends on the content of G2SRSPath
	// so we can't check it here. It is checked inside [NewFromConfig].
	return nil
}

func ReadCLIConfig(ctx *cli.Context) Config {
	return Config{
		SRSNumberToLoad:   ctx.GlobalUint64(kzgflags.SRSLoadingNumberFlagName),
		G1SRSPath:         ctx.GlobalString(kzgflags.G1PathFlagName),
		G2SRSPath:         ctx.GlobalString(kzgflags.G2PathFlagName),
		G2TrailingSRSPath: ctx.GlobalString(kzgflags.G2TrailingPathFlagName),
	}
}
