package eigenda

import (
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	opservice "github.com/ethereum-optimism/optimism/op-service"
	"github.com/urfave/cli/v2"
)

const (
	RPCFlagName                      = "eigenda-rpc"
	StatusQueryRetryIntervalFlagName = "eigenda-status-query-retry-interval"
	StatusQueryTimeoutFlagName       = "eigenda-status-query-timeout"
	DisableTlsFlagName               = "eigenda-disable-tls"
	ResponseTimeoutFlagName          = "eigenda-response-timeout"
	CustomQuorumIDsFlagName          = "eigenda-custom-quorum-ids"
	SignerPrivateKeyHexFlagName      = "eigenda-signer-private-key-hex"
	PutBlobEncodingVersionFlagName   = "eigenda-put-blob-encoding-version"
	// Kzg flags
	G1PathFlagName    = "eigenda-g1-path"
	G2TauFlagName     = "eigenda-g2-tau-path"
	CachePathFlagName = "eigenda-cache-path"
)

type Config struct {
	ClientConfig clients.EigenDAClientConfig

	// The blob encoding version to use when writing blobs from the high level interface.
	PutBlobEncodingVersion codecs.BlobEncodingVersion

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
		ClientConfig: clients.EigenDAClientConfig{
			/* Required Flags */
			RPC:                      ctx.String(RPCFlagName),
			StatusQueryRetryInterval: ctx.Duration(StatusQueryRetryIntervalFlagName),
			StatusQueryTimeout:       ctx.Duration(StatusQueryTimeoutFlagName),
			DisableTLS:               ctx.Bool(DisableTlsFlagName),
			ResponseTimeout:          ctx.Duration(ResponseTimeoutFlagName),
			CustomQuorumIDs:          ctx.UintSlice(CustomQuorumIDsFlagName),
			SignerPrivateKeyHex:      ctx.String(SignerPrivateKeyHexFlagName),
			PutBlobEncodingVersion:   codecs.BlobEncodingVersion(ctx.Uint(PutBlobEncodingVersionFlagName)),
		},
		G1Path:           ctx.String(G1PathFlagName),
		G2PowerOfTauPath: ctx.String(G2TauFlagName),
		CacheDir:         ctx.String(CachePathFlagName),
	}
	return cfg
}

func (m Config) Check() error {
	return nil
}

func CLIFlags(envPrefix string) []cli.Flag {
	prefixEnvVars := func(name string) []string {
		return opservice.PrefixEnvVar(envPrefix, name)
	}
	return []cli.Flag{
		&cli.StringFlag{
			Name:     RPCFlagName,
			Usage:    "RPC endpoint of the EigenDA disperser.",
			EnvVars:  prefixEnvVars("TARGET_RPC"),
			Required: true,
		},
		&cli.DurationFlag{
			Name:     StatusQueryTimeoutFlagName,
			Usage:    "Timeout for aborting an EigenDA blob dispersal if the disperser does not report that the blob has been confirmed dispersed.",
			Value:    30 * time.Minute,
			EnvVars:  prefixEnvVars("TARGET_STATUS_QUERY_TIMEOUT"),
			Required: false,
		},
		&cli.DurationFlag{
			Name:     StatusQueryRetryIntervalFlagName,
			Usage:    "Wait time between retries of EigenDA blob status queries (made while waiting for a blob to be confirmed by).",
			Value:    5 * time.Second,
			EnvVars:  prefixEnvVars("TARGET_STATUS_QUERY_INTERVAL"),
			Required: false,
		},
		&cli.BoolFlag{
			Name:     DisableTlsFlagName,
			Usage:    "Disable TLS when connecting to the EigenDA disperser.",
			Value:    false,
			EnvVars:  prefixEnvVars("TARGET_GRPC_DISABLE_TLS"),
			Required: false,
		},
		&cli.DurationFlag{
			Name:     ResponseTimeoutFlagName,
			Usage:    "The total amount of time that the client will waiting for a response from the EigenDA disperser.",
			Value:    10 * time.Second,
			EnvVars:  prefixEnvVars("RESPONSE_TIMEOUT"),
			Required: false,
		},
		&cli.UintSliceFlag{
			Name:     CustomQuorumIDsFlagName,
			Usage:    "The quorum IDs to write blobs to using this client. Should not include default quorums 0 or 1.",
			Value:    cli.NewUintSlice(),
			EnvVars:  prefixEnvVars("CUSTOM_QUORUM_IDS"),
			Required: false,
		},
		&cli.StringFlag{
			Name:     SignerPrivateKeyHexFlagName,
			Usage:    "Signer private key in hex encoded format. This key should not be associated with an Ethereum address holding any funds.",
			EnvVars:  prefixEnvVars("SIGNER_PRIVATE_KEY_HEX"),
			Required: true,
		},
		&cli.UintFlag{
			Name:     PutBlobEncodingVersionFlagName,
			Usage:    "The blob encoding version to use when writing blobs from the high level interface.",
			EnvVars:  prefixEnvVars("PUT_BLOB_ENCODING_VERSION"),
			Value:    1,
			Required: false,
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
