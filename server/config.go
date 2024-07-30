package server

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/utils"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	opservice "github.com/ethereum-optimism/optimism/op-service"
	"github.com/urfave/cli/v2"
)

const (
	EigenDADisperserRPCFlagName          = "eigenda-disperser-rpc"
	EthRPCFlagName                       = "eigenda-eth-rpc"
	SvcManagerAddrFlagName               = "eigenda-svc-manager-addr"
	EthConfirmationDepthFlagName         = "eigenda-eth-confirmation-depth"
	StatusQueryRetryIntervalFlagName     = "eigenda-status-query-retry-interval"
	StatusQueryTimeoutFlagName           = "eigenda-status-query-timeout"
	DisableTlsFlagName                   = "eigenda-disable-tls"
	ResponseTimeoutFlagName              = "eigenda-response-timeout"
	CustomQuorumIDsFlagName              = "eigenda-custom-quorum-ids"
	SignerPrivateKeyHexFlagName          = "eigenda-signer-private-key-hex"
	PutBlobEncodingVersionFlagName       = "eigenda-put-blob-encoding-version"
	DisablePointVerificationModeFlagName = "eigenda-disable-point-verification-mode"
	// Kzg flags
	G1PathFlagName             = "eigenda-g1-path"
	G2TauFlagName              = "eigenda-g2-tau-path"
	CachePathFlagName          = "eigenda-cache-path"
	MaxBlobLengthFlagName      = "eigenda-max-blob-length"
	MemstoreFlagName           = "memstore.enabled"
	MemstoreExpirationFlagName = "memstore.expiration"
	// S3 flags
	S3CredentialTypeFlagName  = "s3.credential-type"
	S3BucketFlagName          = "s3.bucket"
	S3PathFlagName            = "s3.path"
	S3EndpointFlagName        = "s3.endpoint"
	S3AccessKeyIDFlagName     = "s3.access-key-id"     // #nosec G101
	S3AccessKeySecretFlagName = "s3.access-key-secret" // #nosec G101
)

const BytesPerSymbol = 31
const MaxCodingRatio = 8

var MaxSRSPoints = math.Pow(2, 28)

var MaxAllowedBlobSize = uint64(MaxSRSPoints * BytesPerSymbol / MaxCodingRatio)

type Config struct {
	S3Config store.S3Config

	ClientConfig clients.EigenDAClientConfig

	// The blob encoding version to use when writing blobs from the high level interface.
	PutBlobEncodingVersion codecs.BlobEncodingVersion

	// ETH vars
	EthRPC               string
	SvcManagerAddr       string
	EthConfirmationDepth int64

	// KZG vars
	CacheDir string
	G1Path   string
	G2Path   string

	MaxBlobLength      string
	maxBlobLengthBytes uint64

	G2PowerOfTauPath string

	// Memstore Config params
	MemstoreEnabled        bool
	MemstoreBlobExpiration time.Duration
}

func (c *Config) GetMaxBlobLength() (uint64, error) {
	if c.maxBlobLengthBytes == 0 {
		numBytes, err := utils.ParseBytesAmount(c.MaxBlobLength)
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

func (c *Config) VerificationCfg() *verify.Config {
	numBytes, err := c.GetMaxBlobLength()
	if err != nil {
		panic(fmt.Errorf("Check() was not called on config object, err is not nil: %w", err))
	}

	kzgCfg := &kzg.KzgConfig{
		G1Path:          c.G1Path,
		G2PowerOf2Path:  c.G2PowerOfTauPath,
		CacheDir:        c.CacheDir,
		SRSOrder:        268435456,     // 2 ^ 32
		SRSNumberToLoad: numBytes / 32, // # of fp.Elements
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	if c.EthRPC == "" || c.SvcManagerAddr == "" {
		return &verify.Config{
			Verify:    false,
			KzgConfig: kzgCfg,
		}
	}

	return &verify.Config{
		Verify:               true,
		RPCURL:               c.EthRPC,
		SvcManagerAddr:       c.SvcManagerAddr,
		KzgConfig:            kzgCfg,
		EthConfirmationDepth: uint64(c.EthConfirmationDepth),
	}

}

// NewConfig parses the Config from the provided flags or environment variables.
func ReadConfig(ctx *cli.Context) Config {
	cfg := Config{
		S3Config: store.S3Config{
			S3CredentialType: toS3CredentialType(ctx.String(S3CredentialTypeFlagName)),
			Bucket:           ctx.String(S3BucketFlagName),
			Path:             ctx.String(S3PathFlagName),
			Endpoint:         ctx.String(S3EndpointFlagName),
			AccessKeyID:      ctx.String(S3AccessKeyIDFlagName),
			AccessKeySecret:  ctx.String(S3AccessKeySecretFlagName),
		},
		ClientConfig: clients.EigenDAClientConfig{
			RPC:                          ctx.String(EigenDADisperserRPCFlagName),
			StatusQueryRetryInterval:     ctx.Duration(StatusQueryRetryIntervalFlagName),
			StatusQueryTimeout:           ctx.Duration(StatusQueryTimeoutFlagName),
			DisableTLS:                   ctx.Bool(DisableTlsFlagName),
			ResponseTimeout:              ctx.Duration(ResponseTimeoutFlagName),
			CustomQuorumIDs:              ctx.UintSlice(CustomQuorumIDsFlagName),
			SignerPrivateKeyHex:          ctx.String(SignerPrivateKeyHexFlagName),
			PutBlobEncodingVersion:       codecs.BlobEncodingVersion(ctx.Uint(PutBlobEncodingVersionFlagName)),
			DisablePointVerificationMode: ctx.Bool(DisablePointVerificationModeFlagName),
		},
		G1Path:                 ctx.String(G1PathFlagName),
		G2PowerOfTauPath:       ctx.String(G2TauFlagName),
		CacheDir:               ctx.String(CachePathFlagName),
		MaxBlobLength:          ctx.String(MaxBlobLengthFlagName),
		SvcManagerAddr:         ctx.String(SvcManagerAddrFlagName),
		EthRPC:                 ctx.String(EthRPCFlagName),
		EthConfirmationDepth:   ctx.Int64(EthConfirmationDepthFlagName),
		MemstoreEnabled:        ctx.Bool(MemstoreFlagName),
		MemstoreBlobExpiration: ctx.Duration(MemstoreExpirationFlagName),
	}
	cfg.ClientConfig.WaitForFinalization = (cfg.EthConfirmationDepth < 0)

	return cfg
}

func toS3CredentialType(s string) store.S3CredentialType {
	if s == string(store.S3CredentialStatic) {
		return store.S3CredentialStatic
	} else if s == string(store.S3CredentialIAM) {
		return store.S3CredentialIAM
	}
	return store.S3CredentialUnknown
}

// Check ... verifies that configuration values are adequately set
func (cfg *Config) Check() error {
	l, err := cfg.GetMaxBlobLength()
	if err != nil {
		return err
	}

	if l == 0 {
		return fmt.Errorf("max blob length is 0")
	}

	if cfg.SvcManagerAddr != "" && cfg.EthRPC == "" {
		return fmt.Errorf("svc manager address is set, but Eth RPC is not set")
	}

	if cfg.EthRPC != "" && cfg.SvcManagerAddr == "" {
		return fmt.Errorf("eth rpc is set, but svc manager address is not set")
	}

	if cfg.EthConfirmationDepth >= 0 && (cfg.SvcManagerAddr == "" || cfg.EthRPC == "") {
		return fmt.Errorf("eth confirmation depth is set for certificate verification, but Eth RPC or SvcManagerAddr is not set")
	}

	if cfg.S3Config.S3CredentialType == store.S3CredentialUnknown {
		return fmt.Errorf("s3 credential type must be set")
	}
	if cfg.S3Config.S3CredentialType == store.S3CredentialStatic {
		if cfg.S3Config.Endpoint != "" && (cfg.S3Config.AccessKeyID == "" || cfg.S3Config.AccessKeySecret == "") {
			return fmt.Errorf("s3 endpoint is set, but access key id or access key secret is not set")
		}
	}

	if !cfg.MemstoreEnabled && cfg.ClientConfig.RPC == "" {
		return fmt.Errorf("eigenda disperser rpc url is not set")
	}

	return nil
}

func CLIFlags(envPrefix string) []cli.Flag {
	prefixEnvVars := func(name string) []string {
		return opservice.PrefixEnvVar(envPrefix, name)
	}
	return []cli.Flag{
		&cli.StringFlag{
			Name:    EigenDADisperserRPCFlagName,
			Usage:   "RPC endpoint of the EigenDA disperser.",
			EnvVars: prefixEnvVars("EIGENDA_DISPERSER_RPC"),
		},
		&cli.DurationFlag{
			Name:    StatusQueryTimeoutFlagName,
			Usage:   "Duration to wait for a blob to finalize after being sent for dispersal. Default is 30 minutes.",
			Value:   30 * time.Minute,
			EnvVars: prefixEnvVars("STATUS_QUERY_TIMEOUT"),
		},
		&cli.DurationFlag{
			Name:    StatusQueryRetryIntervalFlagName,
			Usage:   "Interval between retries when awaiting network blob finalization. Default is 5 seconds.",
			Value:   5 * time.Second,
			EnvVars: prefixEnvVars("STATUS_QUERY_INTERVAL"),
		},
		&cli.BoolFlag{
			Name:    DisableTlsFlagName,
			Usage:   "Disable TLS for gRPC communication with the EigenDA disperser. Default is false.",
			Value:   false,
			EnvVars: prefixEnvVars("GRPC_DISABLE_TLS"),
		},
		&cli.DurationFlag{
			Name:    ResponseTimeoutFlagName,
			Usage:   "Total time to wait for a response from the EigenDA disperser. Default is 10 seconds.",
			Value:   10 * time.Second,
			EnvVars: prefixEnvVars("RESPONSE_TIMEOUT"),
		},
		&cli.UintSliceFlag{
			Name:    CustomQuorumIDsFlagName,
			Usage:   "Custom quorum IDs for writing blobs. Should not include default quorums 0 or 1.",
			Value:   cli.NewUintSlice(),
			EnvVars: prefixEnvVars("CUSTOM_QUORUM_IDS"),
		},
		&cli.StringFlag{
			Name:    SignerPrivateKeyHexFlagName,
			Usage:   "Hex-encoded signer private key. This key should not be associated with an Ethereum address holding any funds.",
			EnvVars: prefixEnvVars("SIGNER_PRIVATE_KEY_HEX"),
		},
		&cli.UintFlag{
			Name:    PutBlobEncodingVersionFlagName,
			Usage:   "Blob encoding version to use when writing blobs from the high-level interface.",
			EnvVars: prefixEnvVars("PUT_BLOB_ENCODING_VERSION"),
			Value:   0,
		},
		&cli.BoolFlag{
			Name:    DisablePointVerificationModeFlagName,
			Usage:   "Disable point verification mode. This mode performs IFFT on data before writing and FFT on data after reading. Disabling requires supplying the entire blob for verification against the KZG commitment.",
			EnvVars: prefixEnvVars("DISABLE_POINT_VERIFICATION_MODE"),
			Value:   false,
		},
		&cli.StringFlag{
			Name:    MaxBlobLengthFlagName,
			Usage:   "Maximum blob length to be written or read from EigenDA. Determines the number of SRS points loaded into memory for KZG commitments. Example units: '30MiB', '4Kb', '30MB'. Maximum size slightly exceeds 1GB.",
			EnvVars: prefixEnvVars("MAX_BLOB_LENGTH"),
			Value:   "2MiB",
		},
		&cli.StringFlag{
			Name:    G1PathFlagName,
			Usage:   "Directory path to g1.point file.",
			EnvVars: prefixEnvVars("TARGET_KZG_G1_PATH"),
			Value:   "resources/g1.point",
		},
		&cli.StringFlag{
			Name:    G2TauFlagName,
			Usage:   "Directory path to g2.point.powerOf2 file.",
			EnvVars: prefixEnvVars("TARGET_G2_TAU_PATH"),
			Value:   "resources/g2.point.powerOf2",
		},
		&cli.StringFlag{
			Name:    CachePathFlagName,
			Usage:   "Directory path to SRS tables for caching.",
			EnvVars: prefixEnvVars("TARGET_CACHE_PATH"),
			Value:   "resources/SRSTables/",
		},
		&cli.StringFlag{
			Name:    EthRPCFlagName,
			Usage:   "JSON RPC node endpoint for the Ethereum network used for finalizing DA blobs. See available list here: https://docs.eigenlayer.xyz/eigenda/networks/",
			EnvVars: prefixEnvVars("ETH_RPC"),
		},
		&cli.StringFlag{
			Name:    SvcManagerAddrFlagName,
			Usage:   "The deployed EigenDA service manager address. The list can be found here: https://github.com/Layr-Labs/eigenlayer-middleware/?tab=readme-ov-file#current-mainnet-deployment",
			EnvVars: prefixEnvVars("SERVICE_MANAGER_ADDR"),
		},
		&cli.Int64Flag{
			Name:    EthConfirmationDepthFlagName,
			Usage:   "The number of Ethereum blocks of confirmation that the DA batch submission tx must have before it is assumed by the proxy to be final. The value of `0` indicates that the proxy shouldn't wait for any confirmations.",
			EnvVars: prefixEnvVars("ETH_CONFIRMATION_DEPTH"),
			Value:   -1,
		},
		&cli.BoolFlag{
			Name:    MemstoreFlagName,
			Usage:   "Whether to use mem-store for DA logic.",
			EnvVars: []string{"MEMSTORE_ENABLED"},
		},
		&cli.DurationFlag{
			Name:    MemstoreExpirationFlagName,
			Usage:   "Duration that a mem-store blob/commitment pair are allowed to live.",
			Value:   25 * time.Minute,
			EnvVars: []string{"MEMSTORE_EXPIRATION"},
		},
		&cli.StringFlag{
			Name:    S3CredentialTypeFlagName,
			Usage:   "The way to authenticate to S3, options are [iam, static]",
			EnvVars: prefixEnvVars("S3_CREDENTIAL_TYPE"),
		},
		&cli.StringFlag{
			Name:    S3BucketFlagName,
			Usage:   "bucket name for S3 storage",
			EnvVars: prefixEnvVars("S3_BUCKET"),
		},
		&cli.StringFlag{
			Name:    S3PathFlagName,
			Usage:   "path for S3 storage",
			EnvVars: prefixEnvVars("S3_PATH"),
		},
		&cli.StringFlag{
			Name:    S3EndpointFlagName,
			Usage:   "endpoint for S3 storage",
			Value:   "",
			EnvVars: prefixEnvVars("S3_ENDPOINT"),
		},
		&cli.StringFlag{
			Name:    S3AccessKeyIDFlagName,
			Usage:   "access key id for S3 storage",
			Value:   "",
			EnvVars: prefixEnvVars("S3_ACCESS_KEY_ID"),
		}, &cli.StringFlag{
			Name:    S3AccessKeySecretFlagName,
			Usage:   "access key secret for S3 storage",
			Value:   "",
			EnvVars: prefixEnvVars("S3_ACCESS_KEY_SECRET"),
		},
	}
}
