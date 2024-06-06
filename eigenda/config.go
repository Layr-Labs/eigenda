package eigenda

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	opservice "github.com/ethereum-optimism/optimism/op-service"
	"github.com/urfave/cli/v2"
)

const (
	RPCFlagName                          = "eigenda-rpc"
	StatusQueryRetryIntervalFlagName     = "eigenda-status-query-retry-interval"
	StatusQueryTimeoutFlagName           = "eigenda-status-query-timeout"
	DisableTlsFlagName                   = "eigenda-disable-tls"
	ResponseTimeoutFlagName              = "eigenda-response-timeout"
	CustomQuorumIDsFlagName              = "eigenda-custom-quorum-ids"
	SignerPrivateKeyHexFlagName          = "eigenda-signer-private-key-hex"
	PutBlobEncodingVersionFlagName       = "eigenda-put-blob-encoding-version"
	DisablePointVerificationModeFlagName = "eigenda-disable-point-verification-mode"
	// Kzg flags
	G1PathFlagName        = "eigenda-g1-path"
	G2TauFlagName         = "eigenda-g2-tau-path"
	CachePathFlagName     = "eigenda-cache-path"
	MaxBlobLengthFlagName = "eigenda-max-blob-length"
)

const BytesPerSymbol = 31
const MaxCodingRatio = 8

var MaxSRSPoints = math.Pow(2, 28)

var MaxAllowedBlobSize = uint64(MaxSRSPoints * BytesPerSymbol / MaxCodingRatio)

type Config struct {
	ClientConfig clients.EigenDAClientConfig

	// The blob encoding version to use when writing blobs from the high level interface.
	PutBlobEncodingVersion codecs.BlobEncodingVersion

	// KZG vars
	CacheDir string

	G1Path string
	G2Path string

	MaxBlobLength      string
	maxBlobLengthBytes uint64

	G2PowerOfTauPath string
}

func (c *Config) GetMaxBlobLength() (uint64, error) {
	if c.maxBlobLengthBytes == 0 {
		numBytes, err := common.ParseBytesAmount(c.MaxBlobLength)
		if err != nil {
			return 0, err
		}

		if numBytes > MaxAllowedBlobSize {
			return 0, fmt.Errorf("excluding disperser constraints on max blob size, SRS points constrain the maxBlobLength configuration parameter to be less than than ~1 GB (%d bytes)", MaxAllowedBlobSize)
		}

		c.maxBlobLengthBytes = numBytes
	}

	return c.maxBlobLengthBytes, nil
}

func (c *Config) KzgConfig() *kzg.KzgConfig {
	numBytes, err := c.GetMaxBlobLength()
	if err != nil {
		panic(fmt.Errorf("Check() was not called on config object, err is not nil: %w", err))
	}

	numPointsNeeded := uint64(math.Ceil(float64(numBytes) / BytesPerSymbol))
	return &kzg.KzgConfig{
		G1Path:          c.G1Path,
		G2PowerOf2Path:  c.G2PowerOfTauPath,
		CacheDir:        c.CacheDir,
		SRSOrder:        numPointsNeeded,
		SRSNumberToLoad: numPointsNeeded,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}
}

// NewConfig parses the Config from the provided flags or environment variables.
func ReadConfig(ctx *cli.Context) Config {
	cfg := Config{
		ClientConfig: clients.EigenDAClientConfig{
			/* Required Flags */
			RPC:                          ctx.String(RPCFlagName),
			StatusQueryRetryInterval:     ctx.Duration(StatusQueryRetryIntervalFlagName),
			StatusQueryTimeout:           ctx.Duration(StatusQueryTimeoutFlagName),
			DisableTLS:                   ctx.Bool(DisableTlsFlagName),
			ResponseTimeout:              ctx.Duration(ResponseTimeoutFlagName),
			CustomQuorumIDs:              ctx.UintSlice(CustomQuorumIDsFlagName),
			SignerPrivateKeyHex:          ctx.String(SignerPrivateKeyHexFlagName),
			PutBlobEncodingVersion:       codecs.BlobEncodingVersion(ctx.Uint(PutBlobEncodingVersionFlagName)),
			DisablePointVerificationMode: ctx.Bool(DisablePointVerificationModeFlagName),
		},
		G1Path:           ctx.String(G1PathFlagName),
		G2PowerOfTauPath: ctx.String(G2TauFlagName),
		CacheDir:         ctx.String(CachePathFlagName),
		MaxBlobLength:    ctx.String(MaxBlobLengthFlagName),
	}
	return cfg
}

func (m Config) Check() error {
	_, err := m.GetMaxBlobLength()
	if err != nil {
		return err
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
			EnvVars: prefixEnvVars("RPC"),
		},
		&cli.DurationFlag{
			Name:    StatusQueryTimeoutFlagName,
			Usage:   "Timeout for aborting an EigenDA blob dispersal if the disperser does not report that the blob has been finalized dispersed.",
			Value:   30 * time.Minute,
			EnvVars: prefixEnvVars("STATUS_QUERY_TIMEOUT"),
		},
		&cli.DurationFlag{
			Name:    StatusQueryRetryIntervalFlagName,
			Usage:   "Wait time between retries of EigenDA blob status queries (made while waiting for a blob to be confirmed by).",
			Value:   5 * time.Second,
			EnvVars: prefixEnvVars("STATUS_QUERY_INTERVAL"),
		},
		&cli.BoolFlag{
			Name:    DisableTlsFlagName,
			Usage:   "Disable TLS when connecting to the EigenDA disperser.",
			Value:   false,
			EnvVars: prefixEnvVars("GRPC_DISABLE_TLS"),
		},
		&cli.DurationFlag{
			Name:    ResponseTimeoutFlagName,
			Usage:   "The total amount of time that the client will wait for a response from the EigenDA disperser.",
			Value:   10 * time.Second,
			EnvVars: prefixEnvVars("RESPONSE_TIMEOUT"),
		},
		&cli.UintSliceFlag{
			Name:    CustomQuorumIDsFlagName,
			Usage:   "The quorum IDs to write blobs to using this client. Should not include default quorums 0 or 1.",
			Value:   cli.NewUintSlice(),
			EnvVars: prefixEnvVars("CUSTOM_QUORUM_IDS"),
		},
		&cli.StringFlag{
			Name:    SignerPrivateKeyHexFlagName,
			Usage:   "Signer private key in hex encoded format. This key should not be associated with an Ethereum address holding any funds.",
			EnvVars: prefixEnvVars("SIGNER_PRIVATE_KEY_HEX"),
		},
		&cli.UintFlag{
			Name:    PutBlobEncodingVersionFlagName,
			Usage:   "The blob encoding version to use when writing blobs from the high level interface.",
			EnvVars: prefixEnvVars("PUT_BLOB_ENCODING_VERSION"),
			Value:   0,
		},
		&cli.BoolFlag{
			Name:    DisablePointVerificationModeFlagName,
			Usage:   "Point verification mode does an IFFT on data before it is written, and does an FFT on data after it is read. This makes it possible to open points on the KZG commitment to prove that the field elements correspond to the commitment. With this mode disabled, you will need to supply the entire blob to perform a verification that any part of the data matches the KZG commitment.",
			EnvVars: prefixEnvVars("DISABLE_POINT_VERIFICATION_MODE"),
			Value:   false,
		},
		&cli.StringFlag{
			Name:    MaxBlobLengthFlagName,
			Usage:   "Maximum size in string representation (e.g. \"10mb\", \"4 KiB\") of blobs to be dispersed and verified using this proxy. This impacts the number of SRS points loaded into memory.",
			EnvVars: prefixEnvVars("TARGET_KZG_G1_PATH"),
			Value:   "2MiB",
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
