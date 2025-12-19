package node

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments/reservation/reservationvalidation"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	"github.com/Layr-Labs/eigenda/node/flags"
	"github.com/docker/go-units"

	blssignerTypes "github.com/Layr-Labs/eigensdk-go/signer/bls/types"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/urfave/cli"
)

const (
	// Min number of seconds for the ExpirationPollIntervalSecFlag.
	minExpirationPollIntervalSec   = 3
	minReachabilityPollIntervalSec = 10
	AppName                        = "da-node"
)

// Config contains all of the configuration information for a DA node.
type Config struct {
	Hostname                        string
	RetrievalPort                   string
	DispersalPort                   string
	InternalRetrievalPort           string
	InternalDispersalPort           string
	V2DispersalPort                 string
	V2RetrievalPort                 string
	InternalV2DispersalPort         string
	InternalV2RetrievalPort         string
	EnableNodeApi                   bool
	NodeApiPort                     string
	EnableMetrics                   bool
	MetricsPort                     int
	OnchainMetricsInterval          int64
	Timeout                         time.Duration
	RegisterNodeAtStart             bool
	ExpirationPollIntervalSec       uint64
	LevelDBDisableSeeksCompactionV1 bool
	LevelDBSyncWritesV1             bool
	EnableTestMode                  bool
	OverrideBlockStaleMeasure       uint64
	OverrideStoreDurationBlocks     uint64
	// If set, overrides the default TTL for v2 chunks
	OverrideV2Ttl                  time.Duration
	QuorumIDList                   []core.QuorumID
	DbPath                         string
	LogPath                        string
	ID                             core.OperatorID
	EigenDADirectory               string
	PubIPProviders                 []string
	PubIPCheckInterval             time.Duration
	ChurnerUrl                     string
	DataApiUrl                     string
	NumBatchValidators             int
	NumBatchDeserializationWorkers int
	EnableGnarkBundleEncoding      bool
	ClientIPHeader                 string
	ChurnerUseSecureGrpc           bool
	RelayUseSecureGrpc             bool
	RelayMaxMessageSize            uint
	// The number of connections to establish with each relay node.
	RelayConnectionPoolSize        uint
	ReachabilityPollIntervalSec    uint64
	DisableNodeInfoResources       bool
	StoreChunksRequestMaxPastAge   time.Duration
	StoreChunksRequestMaxFutureAge time.Duration
	// Temporary blacklisting forgiveness window for dispersers that send invalid StoreChunks requests.
	// If set to 0, blacklisting is disabled.
	DisperserBlacklistForgivenessWindow time.Duration

	BlsSignerConfig blssignerTypes.SignerConfig

	EthClientConfig geth.EthClientConfig
	LoggerConfig    common.LoggerConfig
	EncoderConfig   kzg.KzgConfig

	EnableV1 bool
	EnableV2 bool

	OnchainStateRefreshInterval time.Duration
	ChunkDownloadTimeout        time.Duration
	GRPCMsgSizeLimitV2          int

	PprofHttpPort string
	EnablePprof   bool

	// the size of the cache for storing public keys of dispersers
	DispersalAuthenticationKeyCacheSize int
	// the timeout for disperser keys (after which the disperser key is reloaded from the chain)
	DisperserKeyTimeout time.Duration

	// The size of the pool where chunks are downloaded from the relay network.
	DownloadPoolSize int

	// A special test only setting. If true, then littDB will throw an error if the same data is written twice.
	LittDBDoubleWriteProtection bool

	// The percentage of the total memory to use for the write cache in littDB as a fraction of 1.0, where 1.0
	// means that all available memory will be used for the write cache (don't actually use 1.0, that leaves no buffer
	// for other stuff). Ignored if LittDBWriteCacheSizeGB is set.
	LittDBWriteCacheSizeFraction float64

	// The size of the cache for storing recently written chunks in littDB. Ignored if 0. If set,
	// this config value overrides the LittDBWriteCacheSizeFraction value.
	LittDBWriteCacheSizeBytes uint64

	// The percentage of the total memory to use for the read cache in littDB as a fraction of 1.0, where 1.0
	// means that all available memory will be used for the read cache (don't actually use 1.0, that leaves no buffer
	// for other stuff). Ignored if LittDBReadCacheSizeGB is set.
	LittDBReadCacheSizeFraction float64

	// The size of the cache for storing recently read chunks in littDB. Ignored if 0. If set,
	// this config value overrides the LittDBReadCacheSizeFraction value.
	LittDBReadCacheSizeBytes uint64

	// The list of paths to the littDB storage directories. Data is spread across these directories.
	// Directories do not need to be on the same filesystem.
	LittDBStoragePaths []string

	// If true, then LittDB will refuse to start if it can't acquire locks on the database file structure.
	//
	// Ideally, this would always be enabled. But PID reuse in common platforms such as Docker/Kubernetes can lead to
	// a breakdown in lock files being able to detect unsafe concurrent access to the database. Since many (if not most)
	// users of this software will be running in such an environment, this is disabled by default.
	LittRespectLocks bool

	// The minimum interval between littDB flushes. If zero, then there is no minimum interval.
	// Useful for "batching" flush operations when flush operations become extremely frequent.
	// Set this to zero to disable this feature.
	LittMinimumFlushInterval time.Duration

	// If set, the directory where littDB incremental snapshots are stored.
	//
	// WARNING: if snapshots are written to this directory, the responsibility of pruning those snapshots lies
	// external to the node. LittDB will write to this directory, but never delete anything from it. If data is not
	// periodically pruned, the disk will eventually fill up. It is highly suggested to use the LittDB cli
	// for managing this directory.
	LittSnapshotDirectory string

	// The rate limit for the number of bytes served by the GetChunks API if the data is in the cache.
	// Unit is in megabytes per second.
	GetChunksHotCacheReadLimitMB float64

	// The burst limit for the number of bytes served by the GetChunks API if the data is in the cache.
	// Unit is in megabytes.
	GetChunksHotBurstLimitMB float64

	// The rate limit for the number of bytes served by the GetChunks API if the data is not in the cache.
	// Unit is in megabytes per second.
	GetChunksColdCacheReadLimitMB float64

	// The burst limit for the number of bytes served by the GetChunks API if the data is not in the cache.
	// Unit is in megabytes.
	GetChunksColdBurstLimitMB float64

	// GCSafetyBufferSizeFraction is the fraction of the total memory to use as a safety buffer for the garbage
	// collector. If non-zero, the garbage collector will be instructed to aggressively garbage collect so as to
	// keep this amount of memory free. Useful for preventing kubernetes from OOM-killing the process. Ignored if
	// GCSafetyBufferSizeGB is greater than 0.
	GCSafetyBufferSizeFraction float64

	// Defines a safety buffer for the garbage collector. If non-zero, the garbage collector will be instructed
	// to aggressively garbage collect so as to keep this amount of memory free. Useful for preventing kubernetes
	// from OOM-killing the process. Overrides the GCSafetyBufferSizeFraction value if greater than 0.
	GCSafetyBufferSizeBytes uint64

	// The maximum amount of time to wait to acquire buffer capacity to store chunks in the StoreChunks() gRPC request.
	StoreChunksBufferTimeout time.Duration

	// StoreChunksBufferSizeFraction controls the maximum memory that can be used to store chunks in the
	// StoreChunks() gRPC request buffer, as a fraction of the total memory available to the process.
	// Ignored if StoreChunksBufferSizeBytes is greater than 0.
	StoreChunksBufferSizeFraction float64

	// StoreChunksBufferSizeBytes controls the maximum memory that can be used to store chunks in the
	// StoreChunks() gRPC request buffer, in bytes. If set, this config value overrides the
	// StoreChunksBufferSizeFraction value if greater than 0.
	StoreChunksBufferSizeBytes uint64

	// The size of the cache for operator states. Cache will remember operator states for this number of unique blocks.
	OperatorStateCacheSize uint64

	// Controls how often the ejection sentinel checks to see if the node is being ejected. This should be configured
	// to be smaller than the onchain ejection period.
	EjectionSentinelPeriod time.Duration

	// If true, the ejection sentinel will attempt to contest ejection by sending a transaction to cancel the ejection.
	EjectionDefenseEnabled bool

	// Under normal circumstances, honest validators should not contest an ejection if they are running software that
	// does not meet the minimum version number as defined onchain. However, if the governing body in control of
	// setting the minimum version number goes rogue, honest validators may want to contest ejection regardless of the
	// claimed minimum version number.
	IgnoreVersionForEjectionDefense bool

	// Whether the validator should perform payment validation.
	//
	// TODO(litt3): This is a temporary field, which will be removed once the new payments system is fully in place.
	// Payment validation is currently optional to make implementation and testing possible before actually shipping
	// the new payments system.
	EnablePaymentValidation        bool
	ReservationLedgerCacheConfig   reservationvalidation.ReservationLedgerCacheConfig
	EnablePerAccountPaymentMetrics bool
}

// NewConfig parses the Config from the provided flags or environment variables and
// returns a Config.
func NewConfig(ctx *cli.Context) (*Config, error) {
	timeout, err := time.ParseDuration(ctx.GlobalString(flags.TimeoutFlag.Name))
	if err != nil {
		return &Config{}, err
	}

	idsStr := strings.Split(ctx.GlobalString(flags.QuorumIDListFlag.Name), ",")
	ids := make([]core.QuorumID, 0)
	for _, id := range idsStr {
		val, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, core.QuorumID(val))
	}
	if len(ids) == 0 {
		return nil, errors.New("no quorum ids provided")
	}

	expirationPollIntervalSec := ctx.GlobalUint64(flags.ExpirationPollIntervalSecFlag.Name)
	if expirationPollIntervalSec < minExpirationPollIntervalSec {
		return nil, fmt.Errorf("the expiration-poll-interval flag must be >= %d seconds", minExpirationPollIntervalSec)
	}

	reachabilityPollIntervalSec := ctx.GlobalUint64(flags.ReachabilityPollIntervalSecFlag.Name)
	if reachabilityPollIntervalSec != 0 && reachabilityPollIntervalSec < minReachabilityPollIntervalSec {
		return nil, fmt.Errorf("the reachability-poll-interval flag must be >= %d seconds or 0 to disable", minReachabilityPollIntervalSec)
	}

	testMode := ctx.GlobalBool(flags.EnableTestModeFlag.Name)

	// Configuration options that require the Node Operator ECDSA key at runtime
	registerNodeAtStart := ctx.GlobalBool(flags.RegisterAtNodeStartFlag.Name)
	pubIPCheckInterval := ctx.GlobalDuration(flags.PubIPCheckIntervalFlag.Name)
	ejectionDefenseEnabled := ctx.GlobalBool(flags.EjectionDefenseEnabledFlag.Name)
	needECDSAKey := registerNodeAtStart || pubIPCheckInterval > 0 || ejectionDefenseEnabled
	if registerNodeAtStart && (ctx.GlobalString(flags.EcdsaKeyFileFlag.Name) == "" || ctx.GlobalString(flags.EcdsaKeyPasswordFlag.Name) == "") {
		return nil, fmt.Errorf("%s and %s are required if %s is enabled", flags.EcdsaKeyFileFlag.Name, flags.EcdsaKeyPasswordFlag.Name, flags.RegisterAtNodeStartFlag.Name)
	}
	if pubIPCheckInterval > 0 && (ctx.GlobalString(flags.EcdsaKeyFileFlag.Name) == "" || ctx.GlobalString(flags.EcdsaKeyPasswordFlag.Name) == "") {
		return nil, fmt.Errorf("%s and %s are required if %s is > 0", flags.EcdsaKeyFileFlag.Name, flags.EcdsaKeyPasswordFlag.Name, flags.PubIPCheckIntervalFlag.Name)
	}
	if ejectionDefenseEnabled && (ctx.GlobalString(flags.EcdsaKeyFileFlag.Name) == "" ||
		ctx.GlobalString(flags.EcdsaKeyPasswordFlag.Name) == "") {
		return nil, fmt.Errorf("%s and %s are required if %s is enabled",
			flags.EcdsaKeyFileFlag.Name, flags.EcdsaKeyPasswordFlag.Name,
			flags.EjectionDefenseEnabledFlag.Name)
	}

	var ethClientConfig geth.EthClientConfig
	if !testMode {
		ethClientConfig = geth.ReadEthClientConfigRPCOnly(ctx)
		if needECDSAKey {
			// Decrypt ECDSA key
			keyContents, err := os.ReadFile(ctx.GlobalString(flags.EcdsaKeyFileFlag.Name))
			if err != nil {
				return nil, fmt.Errorf("could not read ECDSA key file: %v", err)
			}
			sk, err := keystore.DecryptKey(keyContents, ctx.GlobalString(flags.EcdsaKeyPasswordFlag.Name))
			if err != nil {
				return nil, fmt.Errorf("could not decrypt the ECDSA file: %s", ctx.GlobalString(flags.EcdsaKeyFileFlag.Name))
			}
			ethClientConfig.PrivateKeyString = fmt.Sprintf("%x", crypto.FromECDSA(sk.PrivateKey))
		}
	} else {
		ethClientConfig = geth.ReadEthClientConfig(ctx)
	}

	var blsSignerConfig blssignerTypes.SignerConfig
	if testMode && ctx.GlobalString(flags.TestPrivateBlsFlag.Name) != "" {
		privateBls := ctx.GlobalString(flags.TestPrivateBlsFlag.Name)
		blsSignerConfig = blssignerTypes.SignerConfig{
			SignerType: blssignerTypes.PrivateKey,
			PrivateKey: privateBls,
		}
	} else {
		blsSignerCertFilePath := ctx.GlobalString(flags.BLSSignerCertFileFlag.Name)
		enableTLS := len(blsSignerCertFilePath) > 0
		signerType := blssignerTypes.Local

		// check if BLS remote signer configuration is provided
		blsRemoteSignerEnabled := ctx.GlobalBool(flags.BLSRemoteSignerEnabledFlag.Name)
		blsRemoteSignerUrl := ctx.GlobalString(flags.BLSRemoteSignerUrlFlag.Name)
		blsPublicKeyHex := ctx.GlobalString(flags.BLSPublicKeyHexFlag.Name)
		blsKeyFilePath := ctx.GlobalString(flags.BlsKeyFileFlag.Name)
		blsKeyPassword := ctx.GlobalString(flags.BlsKeyPasswordFlag.Name)
		blsSignerAPIKey := ctx.GlobalString(flags.BLSSignerAPIKeyFlag.Name)

		if blsRemoteSignerEnabled && (blsRemoteSignerUrl == "" || blsPublicKeyHex == "") {
			return nil, errors.New("BLS remote signer URL and Public Key Hex is required if BLS remote signer is enabled")
		}
		if !blsRemoteSignerEnabled && (blsKeyFilePath == "" || blsKeyPassword == "") {
			return nil, errors.New("BLS key file and password is required if BLS remote signer is disabled")
		}

		if blsRemoteSignerEnabled && blsSignerAPIKey == "" {
			return nil, errors.New("BLS signer API key is required if BLS remote signer is enabled")
		}

		if blsRemoteSignerEnabled {
			signerType = blssignerTypes.Cerberus
		}

		blsSignerConfig = blssignerTypes.SignerConfig{
			SignerType:       signerType,
			Path:             blsKeyFilePath,
			Password:         blsKeyPassword,
			CerberusUrl:      blsRemoteSignerUrl,
			PublicKeyHex:     blsPublicKeyHex,
			CerberusPassword: blsKeyPassword,
			EnableTLS:        enableTLS,
			TLSCertFilePath:  ctx.GlobalString(flags.BLSSignerCertFileFlag.Name),
			CerberusAPIKey:   blsSignerAPIKey,
		}
	}

	internalDispersalFlag := ctx.GlobalString(flags.InternalDispersalPortFlag.Name)
	internalRetrievalFlag := ctx.GlobalString(flags.InternalRetrievalPortFlag.Name)
	if internalDispersalFlag == "" {
		internalDispersalFlag = ctx.GlobalString(flags.DispersalPortFlag.Name)
	}
	if internalRetrievalFlag == "" {
		internalRetrievalFlag = ctx.GlobalString(flags.RetrievalPortFlag.Name)
	}

	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, err
	}

	runtimeMode := ctx.GlobalString(flags.RuntimeModeFlag.Name)
	switch runtimeMode {
	case flags.ModeV1Only, flags.ModeV2Only, flags.ModeV1AndV2:
		// Valid mode
	default:
		return nil, fmt.Errorf("invalid runtime mode %q: must be one of %s, %s, or %s", runtimeMode, flags.ModeV1Only, flags.ModeV2Only, flags.ModeV1AndV2)
	}

	// Convert mode to v1/v2 enabled flags
	v1Enabled := runtimeMode == flags.ModeV1Only || runtimeMode == flags.ModeV1AndV2
	v2Enabled := runtimeMode == flags.ModeV2Only || runtimeMode == flags.ModeV1AndV2

	// v1 ports must be defined and valid even if v1 is disabled
	dispersalPort := ctx.GlobalString(flags.DispersalPortFlag.Name)
	err = core.ValidatePort(dispersalPort)
	if err != nil {
		return nil, fmt.Errorf("invalid v1 dispersal port: %s", dispersalPort)
	}
	retrievalPort := ctx.GlobalString(flags.RetrievalPortFlag.Name)
	err = core.ValidatePort(retrievalPort)
	if err != nil {
		return nil, fmt.Errorf("invalid v1 retrieval port: %s", retrievalPort)
	}

	v2DispersalPort := ctx.GlobalString(flags.V2DispersalPortFlag.Name)
	v2RetrievalPort := ctx.GlobalString(flags.V2RetrievalPortFlag.Name)
	internalV2DispersalPort := ctx.GlobalString(flags.InternalV2DispersalPortFlag.Name)
	internalV2RetrievalPort := ctx.GlobalString(flags.InternalV2RetrievalPortFlag.Name)
	if internalV2DispersalPort == "" {
		internalV2DispersalPort = v2DispersalPort
	}
	if internalV2RetrievalPort == "" {
		internalV2RetrievalPort = v2RetrievalPort
	}

	if v2Enabled {
		if v2DispersalPort == "" {
			return nil, errors.New("v2 dispersal port (NODE_V2_DISPERSAL_PORT) must be defined when v2 is enabled")
		} else if err := core.ValidatePort(v2DispersalPort); err != nil {
			return nil, fmt.Errorf("invalid v2 dispersal port: %s", v2DispersalPort)
		}
		if v2RetrievalPort == "" {
			return nil, errors.New("v2 retrieval port (NODE_V2_RETRIEVAL_PORT) must be defined when v2 is enabled")
		} else if err := core.ValidatePort(v2RetrievalPort); err != nil {
			return nil, fmt.Errorf("invalid v2 retrieval port: %s", v2RetrievalPort)
		}
	}

	paymentValidationEnabled := ctx.GlobalBool(flags.EnablePaymentValidationFlag.Name)
	var reservationLedgerCacheConfig reservationvalidation.ReservationLedgerCacheConfig
	if paymentValidationEnabled {
		reservationLedgerCacheConfig, err = reservationvalidation.NewReservationLedgerCacheConfig(
			ctx.GlobalInt(flags.ReservationMaxLedgersFlag.Name),
			// TODO(litt3): once the checkpointed onchain config registry is ready, that should be used
			// instead of hardcoding. At that point, this field will be removed from the config struct
			// entirely, and the value will be fetched dynamically at runtime.
			120*time.Second,
			// this is hardcoded: it's a parameter just in case, but it's never expected to change
			ratelimit.OverfillOncePermitted,
			ctx.GlobalDuration(flags.PaymentVaultUpdateIntervalFlag.Name),
		)
		if err != nil {
			return nil, fmt.Errorf("new reservation ledger cache config: %w", err)
		}
	}

	return &Config{
		Hostname:                            ctx.GlobalString(flags.HostnameFlag.Name),
		DispersalPort:                       dispersalPort,
		RetrievalPort:                       retrievalPort,
		InternalDispersalPort:               internalDispersalFlag,
		InternalRetrievalPort:               internalRetrievalFlag,
		V2DispersalPort:                     v2DispersalPort,
		V2RetrievalPort:                     v2RetrievalPort,
		InternalV2DispersalPort:             internalV2DispersalPort,
		InternalV2RetrievalPort:             internalV2RetrievalPort,
		EnableNodeApi:                       ctx.GlobalBool(flags.EnableNodeApiFlag.Name),
		NodeApiPort:                         ctx.GlobalString(flags.NodeApiPortFlag.Name),
		EnableMetrics:                       ctx.GlobalBool(flags.EnableMetricsFlag.Name),
		MetricsPort:                         ctx.GlobalInt(flags.MetricsPortFlag.Name),
		OnchainMetricsInterval:              ctx.GlobalInt64(flags.OnchainMetricsIntervalFlag.Name),
		Timeout:                             timeout,
		RegisterNodeAtStart:                 registerNodeAtStart,
		ExpirationPollIntervalSec:           expirationPollIntervalSec,
		ReachabilityPollIntervalSec:         reachabilityPollIntervalSec,
		EnableTestMode:                      testMode,
		OverrideBlockStaleMeasure:           ctx.GlobalUint64(flags.OverrideBlockStaleMeasureFlag.Name),
		LevelDBDisableSeeksCompactionV1:     ctx.GlobalBool(flags.LevelDBDisableSeeksCompactionV1Flag.Name),
		LevelDBSyncWritesV1:                 ctx.GlobalBool(flags.LevelDBEnableSyncWritesV1Flag.Name),
		OverrideStoreDurationBlocks:         ctx.GlobalUint64(flags.OverrideStoreDurationBlocksFlag.Name),
		OverrideV2Ttl:                       ctx.GlobalDuration(flags.OverrideV2TtlFlag.Name),
		QuorumIDList:                        ids,
		DbPath:                              ctx.GlobalString(flags.DbPathFlag.Name),
		EthClientConfig:                     ethClientConfig,
		EncoderConfig:                       kzg.ReadCLIConfig(ctx),
		LoggerConfig:                        *loggerConfig,
		EigenDADirectory:                    ctx.GlobalString(flags.EigenDADirectoryFlag.Name),
		PubIPProviders:                      ctx.GlobalStringSlice(flags.PubIPProviderFlag.Name),
		PubIPCheckInterval:                  pubIPCheckInterval,
		ChurnerUrl:                          ctx.GlobalString(flags.ChurnerUrlFlag.Name),
		DataApiUrl:                          ctx.GlobalString(flags.DataApiUrlFlag.Name),
		NumBatchValidators:                  ctx.GlobalInt(flags.NumBatchValidatorsFlag.Name),
		NumBatchDeserializationWorkers:      ctx.GlobalInt(flags.NumBatchDeserializationWorkersFlag.Name),
		EnableGnarkBundleEncoding:           ctx.Bool(flags.EnableGnarkBundleEncodingFlag.Name),
		ClientIPHeader:                      ctx.GlobalString(flags.ClientIPHeaderFlag.Name),
		ChurnerUseSecureGrpc:                ctx.GlobalBoolT(flags.ChurnerUseSecureGRPC.Name),
		RelayUseSecureGrpc:                  ctx.GlobalBoolT(flags.RelayUseSecureGRPC.Name),
		RelayMaxMessageSize:                 uint(ctx.GlobalInt(flags.RelayMaxGRPCMessageSizeFlag.Name)),
		RelayConnectionPoolSize:             ctx.GlobalUint(flags.RelayConnectionPoolSizeFlag.Name),
		DisableNodeInfoResources:            ctx.GlobalBool(flags.DisableNodeInfoResourcesFlag.Name),
		BlsSignerConfig:                     blsSignerConfig,
		EnableV2:                            v2Enabled,
		EnableV1:                            v1Enabled,
		OnchainStateRefreshInterval:         ctx.GlobalDuration(flags.OnchainStateRefreshIntervalFlag.Name),
		ChunkDownloadTimeout:                ctx.GlobalDuration(flags.ChunkDownloadTimeoutFlag.Name),
		GRPCMsgSizeLimitV2:                  ctx.GlobalInt(flags.GRPCMsgSizeLimitV2Flag.Name),
		PprofHttpPort:                       ctx.GlobalString(flags.PprofHttpPort.Name),
		EnablePprof:                         ctx.GlobalBool(flags.EnablePprof.Name),
		DispersalAuthenticationKeyCacheSize: ctx.GlobalInt(flags.DispersalAuthenticationKeyCacheSizeFlag.Name),
		DisperserKeyTimeout:                 ctx.GlobalDuration(flags.DisperserKeyTimeoutFlag.Name),
		StoreChunksRequestMaxPastAge:        ctx.GlobalDuration(flags.StoreChunksRequestMaxPastAgeFlag.Name),
		StoreChunksRequestMaxFutureAge:      ctx.GlobalDuration(flags.StoreChunksRequestMaxFutureAgeFlag.Name),
		DisperserBlacklistForgivenessWindow: ctx.GlobalDuration(flags.DisperserBlacklistForgivenessWindowFlag.Name),
		LittDBWriteCacheSizeBytes: uint64(ctx.GlobalFloat64(
			flags.LittDBWriteCacheSizeGBFlag.Name) * units.GiB),
		LittDBWriteCacheSizeFraction:    ctx.GlobalFloat64(flags.LittDBWriteCacheSizeFractionFlag.Name),
		LittDBReadCacheSizeBytes:        uint64(ctx.GlobalFloat64(flags.LittDBReadCacheSizeGBFlag.Name) * units.GiB),
		LittDBReadCacheSizeFraction:     ctx.GlobalFloat64(flags.LittDBReadCacheSizeFractionFlag.Name),
		LittDBStoragePaths:              ctx.GlobalStringSlice(flags.LittDBStoragePathsFlag.Name),
		LittRespectLocks:                ctx.GlobalBool(flags.LittRespectLocksFlag.Name),
		LittMinimumFlushInterval:        ctx.GlobalDuration(flags.LittMinimumFlushIntervalFlag.Name),
		LittSnapshotDirectory:           ctx.GlobalString(flags.LittSnapshotDirectoryFlag.Name),
		DownloadPoolSize:                ctx.GlobalInt(flags.DownloadPoolSizeFlag.Name),
		GetChunksHotCacheReadLimitMB:    ctx.GlobalFloat64(flags.GetChunksHotCacheReadLimitMBFlag.Name),
		GetChunksHotBurstLimitMB:        ctx.GlobalFloat64(flags.GetChunksHotBurstLimitMBFlag.Name),
		GetChunksColdCacheReadLimitMB:   ctx.GlobalFloat64(flags.GetChunksColdCacheReadLimitMBFlag.Name),
		GetChunksColdBurstLimitMB:       ctx.GlobalFloat64(flags.GetChunksColdBurstLimitMBFlag.Name),
		GCSafetyBufferSizeBytes:         uint64(ctx.GlobalFloat64(flags.GCSafetyBufferSizeGBFlag.Name) * units.GiB),
		GCSafetyBufferSizeFraction:      ctx.GlobalFloat64(flags.GCSafetyBufferSizeFractionFlag.Name),
		StoreChunksBufferTimeout:        ctx.GlobalDuration(flags.StoreChunksBufferTimeoutFlag.Name),
		StoreChunksBufferSizeFraction:   ctx.GlobalFloat64(flags.StoreChunksBufferSizeFractionFlag.Name),
		StoreChunksBufferSizeBytes:      uint64(ctx.GlobalFloat64(flags.StoreChunksBufferSizeGBFlag.Name) * units.GiB),
		OperatorStateCacheSize:          ctx.GlobalUint64(flags.OperatorStateCacheSizeFlag.Name),
		EjectionSentinelPeriod:          ctx.GlobalDuration(flags.EjectionSentinelPeriodFlag.Name),
		EjectionDefenseEnabled:          ctx.GlobalBool(flags.EjectionDefenseEnabledFlag.Name),
		IgnoreVersionForEjectionDefense: ctx.GlobalBool(flags.IgnoreVersionForEjectionDefenseFlag.Name),
		EnablePaymentValidation:         paymentValidationEnabled,
		ReservationLedgerCacheConfig:    reservationLedgerCacheConfig,
		EnablePerAccountPaymentMetrics:  ctx.GlobalBool(flags.EnablePerAccountPaymentMetricsFlag.Name),
	}, nil
}
