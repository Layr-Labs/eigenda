package server

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/utils"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/urfave/cli/v2"
)

const (
	// eigenda client flags
	EigenDADisperserRPCFlagName          = "eigenda-disperser-rpc"
	EthRPCFlagName                       = "eigenda-eth-rpc"
	SvcManagerAddrFlagName               = "eigenda-svc-manager-addr"
	EthConfirmationDepthFlagName         = "eigenda-eth-confirmation-depth"
	StatusQueryRetryIntervalFlagName     = "eigenda-status-query-retry-interval"
	StatusQueryTimeoutFlagName           = "eigenda-status-query-timeout"
	DisableTLSFlagName                   = "eigenda-disable-tls"
	ResponseTimeoutFlagName              = "eigenda-response-timeout"
	CustomQuorumIDsFlagName              = "eigenda-custom-quorum-ids"
	SignerPrivateKeyHexFlagName          = "eigenda-signer-private-key-hex"
	PutBlobEncodingVersionFlagName       = "eigenda-put-blob-encoding-version"
	DisablePointVerificationModeFlagName = "eigenda-disable-point-verification-mode"

	// kzg flags
	G1PathFlagName        = "eigenda-g1-path"
	G2TauFlagName         = "eigenda-g2-tau-path"
	CachePathFlagName     = "eigenda-cache-path"
	MaxBlobLengthFlagName = "eigenda-max-blob-length"

	// memstore flags
	MemstoreFlagName           = "memstore.enabled"
	MemstoreExpirationFlagName = "memstore.expiration"
	MemstorePutLatencyFlagName = "memstore.put-latency"
	MemstoreGetLatencyFlagName = "memstore.get-latency"

	// redis client flags
	RedisEndpointFlagName = "redis.endpoint"
	RedisPasswordFlagName = "redis.password"
	RedisDBFlagName       = "redis.db"
	RedisEvictionFlagName = "redis.eviction"

	// S3 client flags
	S3CredentialTypeFlagName  = "s3.credential-type" // #nosec G101
	S3BucketFlagName          = "s3.bucket"          // #nosec G101
	S3PathFlagName            = "s3.path"
	S3EndpointFlagName        = "s3.endpoint"
	S3AccessKeyIDFlagName     = "s3.access-key-id"     // #nosec G101
	S3AccessKeySecretFlagName = "s3.access-key-secret" // #nosec G101
	S3BackupFlagName          = "s3.backup"
	S3TimeoutFlagName         = "s3.timeout"

	// routing flags
	FallbackTargets = "routing.fallback-targets"
	CacheTargets    = "routing.cache-targets"
)

const (
	BytesPerSymbol = 31
	MaxCodingRatio = 8
)

var (
	MaxSRSPoints       = 1 << 28 // 2^28
	MaxAllowedBlobSize = uint64(MaxSRSPoints * BytesPerSymbol / MaxCodingRatio)
)

type Config struct {
	// eigenda
	ClientConfig clients.EigenDAClientConfig

	// the blob encoding version to use when writing blobs from the high level interface.
	PutBlobEncodingVersion codecs.BlobEncodingVersion

	// eth vars
	EthRPC               string
	SvcManagerAddr       string
	EthConfirmationDepth int64

	// kzg vars
	CacheDir         string
	G1Path           string
	G2Path           string
	G2PowerOfTauPath string

	// size constraints
	MaxBlobLength      string
	maxBlobLengthBytes uint64

	// memstore
	MemstoreEnabled        bool
	MemstoreBlobExpiration time.Duration
	MemstoreGetLatency     time.Duration
	MemstorePutLatency     time.Duration

	// routing
	FallbackTargets []string
	CacheTargets    []string

	// secondary storage
	RedisCfg store.RedisConfig
	S3Config store.S3Config
}

// GetMaxBlobLength ... returns the maximum blob length in bytes
func (cfg *Config) GetMaxBlobLength() (uint64, error) {
	if cfg.maxBlobLengthBytes == 0 {
		numBytes, err := utils.ParseBytesAmount(cfg.MaxBlobLength)
		if err != nil {
			return 0, err
		}

		if numBytes > MaxAllowedBlobSize {
			return 0, fmt.Errorf("excluding disperser constraints on max blob size, SRS points constrain the maxBlobLength configuration parameter to be less than than %d bytes", MaxAllowedBlobSize)
		}

		cfg.maxBlobLengthBytes = numBytes
	}

	return cfg.maxBlobLengthBytes, nil
}

// VerificationCfg ... returns certificate config used to verify blobs from eigenda
func (cfg *Config) VerificationCfg() *verify.Config {
	numBytes, err := cfg.GetMaxBlobLength()
	if err != nil {
		panic(fmt.Errorf("failed to read max blob length: %w", err))
	}

	kzgCfg := &kzg.KzgConfig{
		G1Path:          cfg.G1Path,
		G2PowerOf2Path:  cfg.G2PowerOfTauPath,
		CacheDir:        cfg.CacheDir,
		SRSOrder:        268435456,     // 2 ^ 32
		SRSNumberToLoad: numBytes / 32, // # of fp.Elements
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	if cfg.EthRPC == "" || cfg.SvcManagerAddr == "" {
		return &verify.Config{
			Verify:    false,
			KzgConfig: kzgCfg,
		}
	}

	return &verify.Config{
		Verify:               true,
		RPCURL:               cfg.EthRPC,
		SvcManagerAddr:       cfg.SvcManagerAddr,
		KzgConfig:            kzgCfg,
		EthConfirmationDepth: uint64(cfg.EthConfirmationDepth),
	}
}

// ReadConfig ... parses the Config from the provided flags or environment variables.
func ReadConfig(ctx *cli.Context) Config {
	cfg := Config{
		RedisCfg: store.RedisConfig{
			Endpoint: ctx.String(RedisEndpointFlagName),
			Password: ctx.String(RedisPasswordFlagName),
			DB:       ctx.Int(RedisDBFlagName),
			Eviction: ctx.Duration(RedisEvictionFlagName),
		},
		S3Config: store.S3Config{
			S3CredentialType: store.StringToS3CredentialType(ctx.String(S3CredentialTypeFlagName)),
			Bucket:           ctx.String(S3BucketFlagName),
			Path:             ctx.String(S3PathFlagName),
			Endpoint:         ctx.String(S3EndpointFlagName),
			AccessKeyID:      ctx.String(S3AccessKeyIDFlagName),
			AccessKeySecret:  ctx.String(S3AccessKeySecretFlagName),
			Backup:           ctx.Bool(S3BackupFlagName),
			Timeout:          ctx.Duration(S3TimeoutFlagName),
		},
		ClientConfig: clients.EigenDAClientConfig{
			RPC:                          ctx.String(EigenDADisperserRPCFlagName),
			StatusQueryRetryInterval:     ctx.Duration(StatusQueryRetryIntervalFlagName),
			StatusQueryTimeout:           ctx.Duration(StatusQueryTimeoutFlagName),
			DisableTLS:                   ctx.Bool(DisableTLSFlagName),
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
		MemstoreGetLatency:     ctx.Duration(MemstoreGetLatencyFlagName),
		MemstorePutLatency:     ctx.Duration(MemstorePutLatencyFlagName),
		FallbackTargets:        ctx.StringSlice(FallbackTargets),
		CacheTargets:           ctx.StringSlice(CacheTargets),
	}
	cfg.ClientConfig.WaitForFinalization = (cfg.EthConfirmationDepth < 0)

	return cfg
}

// checkTargets ... verifies that a backend target slice is constructed correctly
func (cfg *Config) checkTargets(targets []string) error {
	if len(targets) == 0 {
		return nil
	}

	if utils.ContainsDuplicates(targets) {
		return fmt.Errorf("duplicate targets provided: %+v", targets)
	}

	for _, t := range targets {
		if store.StringToBackendType(t) == store.Unknown {
			return fmt.Errorf("unknown fallback target provided: %s", t)
		}
	}

	return nil
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

	if cfg.S3Config.S3CredentialType == store.S3CredentialUnknown && cfg.S3Config.Endpoint != "" {
		return fmt.Errorf("s3 credential type must be set")
	}
	if cfg.S3Config.S3CredentialType == store.S3CredentialStatic {
		if cfg.S3Config.Endpoint != "" && (cfg.S3Config.AccessKeyID == "" || cfg.S3Config.AccessKeySecret == "") {
			return fmt.Errorf("s3 endpoint is set, but access key id or access key secret is not set")
		}
	}

	if cfg.RedisCfg.Endpoint == "" && cfg.RedisCfg.Password != "" {
		return fmt.Errorf("redis password is set, but endpoint is not")
	}

	if !cfg.MemstoreEnabled && cfg.ClientConfig.RPC == "" {
		return fmt.Errorf("eigenda disperser rpc url is not set")
	}

	err = cfg.checkTargets(cfg.FallbackTargets)
	if err != nil {
		return err
	}

	err = cfg.checkTargets(cfg.CacheTargets)
	if err != nil {
		return err
	}

	// verify that same target is not in both fallback and cache targets
	for _, t := range cfg.FallbackTargets {
		if utils.Contains(cfg.CacheTargets, t) {
			return fmt.Errorf("target %s is in both fallback and cache targets", t)
		}
	}

	return nil
}

// s3Flags ... used for S3 backend configuration
func s3Flags() []cli.Flag {
	return []cli.Flag{
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
		},
		&cli.StringFlag{
			Name:    S3AccessKeySecretFlagName,
			Usage:   "access key secret for S3 storage",
			Value:   "",
			EnvVars: prefixEnvVars("S3_ACCESS_KEY_SECRET"),
		},
		&cli.BoolFlag{
			Name:    S3BackupFlagName,
			Usage:   "whether to use S3 as a backup store to ensure resiliency in case of EigenDA read failure",
			Value:   false,
			EnvVars: prefixEnvVars("S3_BACKUP"),
		},
		&cli.DurationFlag{
			Name:    S3TimeoutFlagName,
			Usage:   "timeout for S3 storage operations (e.g. get, put)",
			Value:   5 * time.Second,
			EnvVars: prefixEnvVars("S3_TIMEOUT"),
		},
	}
}

// redisFlags ... used for Redis backend configuration
func redisFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    RedisEndpointFlagName,
			Usage:   "Redis endpoint",
			EnvVars: prefixEnvVars("REDIS_ENDPOINT"),
		},
		&cli.StringFlag{
			Name:    RedisPasswordFlagName,
			Usage:   "Redis password",
			EnvVars: prefixEnvVars("REDIS_PASSWORD"),
		},
		&cli.IntFlag{
			Name:    RedisDBFlagName,
			Usage:   "Redis database",
			Value:   0,
			EnvVars: prefixEnvVars("REDIS_DB"),
		},
		&cli.DurationFlag{
			Name:    RedisEvictionFlagName,
			Usage:   "Redis eviction time",
			Value:   24 * time.Hour,
			EnvVars: prefixEnvVars("REDIS_EVICTION"),
		},
	}
}

func CLIFlags() []cli.Flag {
	// TODO: Decompose all flags into constituent parts based on their respective category / usage
	flags := []cli.Flag{
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
			Name:    DisableTLSFlagName,
			Usage:   "Disable TLS for gRPC communication with the EigenDA disperser. Default is false.",
			Value:   false,
			EnvVars: prefixEnvVars("GRPC_DISABLE_TLS"),
		},
		&cli.DurationFlag{
			Name:    ResponseTimeoutFlagName,
			Usage:   "Total time to wait for a response from the EigenDA disperser. Default is 60 seconds.",
			Value:   60 * time.Second,
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
			Value:   "16MiB",
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
		&cli.DurationFlag{
			Name:    MemstorePutLatencyFlagName,
			Usage:   "Artificial latency added for memstore backend to mimic EigenDA's dispersal latency.",
			Value:   0,
			EnvVars: []string{"MEMSTORE_PUT_LATENCY"},
		},
		&cli.DurationFlag{
			Name:    MemstoreGetLatencyFlagName,
			Usage:   "Artificial latency added for memstore backend to mimic EigenDA's retrieval latency.",
			Value:   0,
			EnvVars: []string{"MEMSTORE_GET_LATENCY"},
		},
		&cli.StringSliceFlag{
			Name:    FallbackTargets,
			Usage:   "List of read fallback targets to rollover to if cert can't be read from EigenDA.",
			Value:   cli.NewStringSlice(),
			EnvVars: prefixEnvVars("FALLBACK_TARGETS"),
		},
		&cli.StringSliceFlag{
			Name:    CacheTargets,
			Usage:   "List of caching targets to use fast reads from EigenDA.",
			Value:   cli.NewStringSlice(),
			EnvVars: prefixEnvVars("CACHE_TARGETS"),
		},
	}

	flags = append(flags, s3Flags()...)
	flags = append(flags, redisFlags()...)
	return flags
}
