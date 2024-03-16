package verifier

import (
	"runtime"

	"github.com/urfave/cli"
)

var (
	/* Required Flags */
	G1PathFlag = cli.StringFlag{
		Name:     "g1-path",
		Usage:    "File path to SRS g1 point",
		Required: true,
		EnvVar:   "G1_PATH",
	}
	G2PathFlag = cli.StringFlag{
		Name:     "g2-path",
		Usage:    "File path to SRS g2 point",
		Required: true,
		EnvVar:   "G2_PATH",
	}

	/* Optional Flags */
	VerifierNumBatchFlag = cli.Uint64Flag{
		Name:     "verifier-num-batch",
		Usage:    "Set total number batch size for parallel parser to work on",
		Required: false,
		EnvVar:   "VERIFIER_NUM_BATCH",
		Value:    uint64(5000),
	}
	NumPointToVerifyFlag = cli.Uint64Flag{
		Name:     "verifier-num-points",
		Usage:    "Set total number of points (g1 and g2) to verify",
		Required: false,
		EnvVar:   "VERIFIER_NUM_POINT",
		Value:    uint64(268435456),
	}
	NumWorkerFlag = cli.IntFlag{
		Name:     "verifier-num-worker",
		Usage:    "Set total number of worker thread",
		Required: false,
		EnvVar:   "NUM_WORKER",
		Value:    runtime.GOMAXPROCS(0),
	}
)

var requiredFlags = []cli.Flag{
	G1PathFlag,
	G2PathFlag,
}

var optionalFlags = []cli.Flag{
	VerifierNumBatchFlag,
	NumPointToVerifyFlag,
	NumWorkerFlag,
}

func ReadCLIConfig(ctx *cli.Context) Config {
	cfg := Config{}
	cfg.G1Path = ctx.String(G1PathFlag.Name)
	cfg.G2Path = ctx.String(G2PathFlag.Name)
	cfg.NumPoint = ctx.Uint64(NumPointToVerifyFlag.Name)
	cfg.NumBatch = ctx.Uint64(VerifierNumBatchFlag.Name)
	cfg.NumWorker = ctx.Int(NumWorkerFlag.Name)

	return cfg
}

func init() {
	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag
