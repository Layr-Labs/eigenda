package eigenda

import (
	"errors"
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	opservice "github.com/ethereum-optimism/optimism/op-service"
	"github.com/urfave/cli/v2"
)

const (
	RPCFlagName                      = "eigenda-rpc"
	StatusQueryRetryIntervalFlagName = "eigenda-status-query-retry-interval"
	StatusQueryTimeoutFlagName       = "eigenda-status-query-timeout"
	UseTlsFlagName                   = "eigenda-use-tls"
	// Kzg flags
	G1PathFlagName    = "eigenda-g1-path"
	G2TauFlagName     = "eigenda-g2-tau-path"
	CachePathFlagName = "eigenda-cache-path"
)

type Config struct {
	// TODO(eigenlayer): Update quorum ID command-line parameters to support passing
	// an arbitrary number of quorum IDs.

	// RPC is the HTTP provider URL for the Data Availability node.
	RPC string

	// The total amount of time that the batcher will spend waiting for EigenDA to confirm a blob
	StatusQueryTimeout time.Duration

	// The amount of time to wait between status queries of a newly dispersed blob
	StatusQueryRetryInterval time.Duration

	// UseTLS specifies whether the client should use TLS as a transport layer when connecting to disperser.
	UseTLS bool

	// KZG vars
	CacheDir string

	G1Path string
	G2Path string

	G2PowerOfTauPath string
}

func (c *Config) KzgConfig() *kzg.KzgConfig {
	return &kzg.KzgConfig{
		G1Path:          c.G1Path,
		G2PowerOf2Path:  c.G2PowerOfTauPath,
		CacheDir:        c.CacheDir,
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}
}

// NewConfig parses the Config from the provided flags or environment variables.
func ReadConfig(ctx *cli.Context) Config {
	cfg := Config{
		/* Required Flags */
		RPC:                      ctx.String(RPCFlagName),
		StatusQueryRetryInterval: ctx.Duration(StatusQueryRetryIntervalFlagName),
		StatusQueryTimeout:       ctx.Duration(StatusQueryTimeoutFlagName),
		UseTLS:                   ctx.Bool(UseTlsFlagName),
		G1Path:                   ctx.String(G1PathFlagName),
		G2PowerOfTauPath:         ctx.String(G2TauFlagName),
		CacheDir:                 ctx.String(CachePathFlagName),
	}
	return cfg
}

func (m Config) Check() error {
	if m.StatusQueryTimeout == 0 {
		return errors.New("EigenDA status query timeout must be greater than 0")
	}
	if m.StatusQueryRetryInterval == 0 {
		return errors.New("EigenDA status query retry interval must be greater than 0")
	}
	return nil
}

func CLIFlags(envPrefix string) []cli.Flag {
	prefixEnvVars := func(name string) []string {
		return opservice.PrefixEnvVar(envPrefix, name)
	}
	return []cli.Flag{
		&cli.StringFlag{
			Name:    RPCFlagName,
			Usage:   "RPC endpoint of the EigenDA disperser.",
			EnvVars: prefixEnvVars("TARGET_RPC"),
		},
		&cli.DurationFlag{
			Name:    StatusQueryTimeoutFlagName,
			Usage:   "Timeout for aborting an EigenDA blob dispersal if the disperser does not report that the blob has been confirmed dispersed.",
			Value:   25 * time.Minute,
			EnvVars: prefixEnvVars("TARGET_STATUS_QUERY_TIMEOUT"),
		},
		&cli.DurationFlag{
			Name:    StatusQueryRetryIntervalFlagName,
			Usage:   "Wait time between retries of EigenDA blob status queries (made while waiting for a blob to be confirmed by).",
			Value:   5 * time.Second,
			EnvVars: prefixEnvVars("TARGET_STATUS_QUERY_INTERVAL"),
		},
		&cli.BoolFlag{
			Name:    UseTlsFlagName,
			Usage:   "Use TLS when connecting to the EigenDA disperser.",
			Value:   true,
			EnvVars: prefixEnvVars("TARGET_GRPC_USE_TLS"),
		},
		&cli.StringFlag{
			Name:    G1PathFlagName,
			Usage:   "Directory path to g1.point file",
			EnvVars: prefixEnvVars("TARGET_KZG_G1_PATH"),
		},
		&cli.StringFlag{
			Name:    G2TauFlagName,
			Usage:   "Directory path to g2.point.powerOf2 file",
			EnvVars: prefixEnvVars("TARGET_G2_TAU_PATH"),
		},
		&cli.StringFlag{
			Name:    CachePathFlagName,
			Usage:   "Directory path to SRS tables",
			EnvVars: prefixEnvVars("TARGET_CACHE_PATH"),
		},
	}
}
