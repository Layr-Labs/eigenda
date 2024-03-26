package parser

import (
	"runtime"

	"github.com/urfave/cli"
)

var (
	/* Required Flags */
	PtauPathFlag = cli.StringFlag{
		Name:     "ptau-path",
		Usage:    "File path to the ptau challenge file",
		Required: true,
		EnvVar:   "PTAU_PATH",
	}

	/* Optional Flags */
	ParserNumBatchFlag = cli.Uint64Flag{
		Name:     "parser-num-batch",
		Usage:    "Set total number batch size for parallel parser to work on",
		Required: false,
		EnvVar:   "PARSER_NUM_BATCH",
		Value:    uint64(50),
	}
	NumPointToParseFlag = cli.Uint64Flag{
		Name:     "parser-num-points",
		Usage:    "Set total number of points (g1 and g2) to parse",
		Required: false,
		EnvVar:   "PARSER_NUM_POINT",
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
	PtauPathFlag,
}

var optionalFlags = []cli.Flag{
	ParserNumBatchFlag,
	NumPointToParseFlag,
	NumWorkerFlag,
}

func ReadCLIConfig(ctx *cli.Context) Config {
	cfg := Config{}
	cfg.PtauPath = ctx.String(PtauPathFlag.Name)
	cfg.NumBatch = ctx.Uint64(ParserNumBatchFlag.Name)
	cfg.NumPoint = ctx.Uint64(NumPointToParseFlag.Name)
	cfg.NumWorker = ctx.Int(NumWorkerFlag.Name)

	return cfg
}

func init() {
	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag
